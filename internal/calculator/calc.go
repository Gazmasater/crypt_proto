package calculator

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"crypt_proto/internal/queue"
	"crypt_proto/pkg/models"
)

type Calculator struct {
	mem     *queue.MemoryStore
	scanner *Scanner
	filter  *ExecutorFilter
	log     *log.Logger
	cfg     Config
	stats   Stats
}

func NewCalculator(mem *queue.MemoryStore, triangles []*Triangle, cfg Config) *Calculator {
	if cfg.MinVolumeUSDT <= 0 {
		cfg.MinVolumeUSDT = defaultMinVolumeUSDT
	}
	if cfg.SearchStepUSDT <= 0 {
		cfg.SearchStepUSDT = defaultSearchStep
	}
	if cfg.StatsEverySec <= 0 {
		cfg.StatsEverySec = 5
	}

	writer := log.Writer()
	c := &Calculator{
		mem:     mem,
		scanner: NewScanner(mem, triangles, cfg),
		filter:  NewExecutorFilter(cfg),
		log:     log.New(writer, "", log.LstdFlags),
		cfg:     cfg,
		stats: Stats{
			ScanRejects:     make(map[string]int64),
			ExecRejects:     make(map[string]int64),
			TriangleMetrics: make(map[string]*TriangleMetrics),
		},
	}

	if cfg.LogMode != LogSilent {
		c.log.Printf(
			"[CALC] started | triangles indexed=%d | minVolume=%.2f USDT | minProfit=%.4f%% | quoteAgeMax=%dms | logMode=%s",
			len(triangles),
			cfg.MinVolumeUSDT,
			cfg.MinProfitPct*100,
			cfg.QuoteAgeMaxMS,
			c.logModeString(),
		)
	}

	return c
}

func (c *Calculator) Run(in <-chan *models.MarketData) {
	ticker := time.NewTicker(time.Duration(c.cfg.StatsEverySec) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case md, ok := <-in:
			if !ok {
				if c.cfg.LogMode != LogSilent {
					c.logStats()
				}
				return
			}
			if md == nil {
				continue
			}

			c.stats.Ticks++
			c.mem.Push(md)

			results := c.scanner.CandidatesFor(md.Symbol, md.Timestamp)
			if len(results) == 0 {
				continue
			}

			for _, res := range results {
				c.stats.TrianglesSeen++
				if !res.OK {
					c.addScanReject(res.Reject, res.Candidate.Triangle)
					continue
				}

				c.stats.Candidates++
				opp, reason, ok := c.filter.Evaluate(res.Candidate)
				c.recordTriangleMetrics(res.Candidate, opp)
				if !ok {
					c.addExecReject(reason, res.Candidate.Triangle)
					continue
				}

				c.stats.Opportunities++
				if opp.ProfitPct > 0 {
					c.stats.Positive++
					c.logOpportunity(opp)
					c.stats.Logged++
				} else {
					c.stats.Negative++
				}
			}
		case <-ticker.C:
			if c.cfg.LogMode != LogSilent {
				c.logStats()
			}
		}
	}
}

func (c *Calculator) recordTriangleMetrics(cand ScanCandidate, opp ExecutableOpportunity) {
	if cand.Triangle == nil {
		return
	}
	key := cand.Triangle.Key()
	m := c.stats.TriangleMetrics[key]
	if m == nil {
		m = &TriangleMetrics{
			Key:              key,
			BestEstimatedPct: -1e9,
			BestIdealPct:     -1e9,
			BestRoundedPct:   -1e9,
			BestFinalPct:     -1e9,
			WorstFinalPct:    1e9,
		}
		c.stats.TriangleMetrics[key] = m
	}
	m.Seen++
	m.LastEstimatedPct = cand.EstimatedPct
	m.LastIdealPct = opp.IdealProfitPct
	m.LastRoundedPct = opp.RoundedProfitPct
	m.LastFinalPct = opp.ProfitPct
	if cand.EstimatedPct > m.BestEstimatedPct {
		m.BestEstimatedPct = cand.EstimatedPct
	}
	if opp.IdealProfitPct > m.BestIdealPct {
		m.BestIdealPct = opp.IdealProfitPct
	}
	if opp.RoundedProfitPct > m.BestRoundedPct {
		m.BestRoundedPct = opp.RoundedProfitPct
	}
	if opp.ProfitPct > m.BestFinalPct {
		m.BestFinalPct = opp.ProfitPct
	}
	if opp.ProfitPct < m.WorstFinalPct {
		m.WorstFinalPct = opp.ProfitPct
	}
}

func (c *Calculator) addScanReject(reason string, tri *Triangle) {
	if reason == "" {
		reason = "unknown"
	}
	c.stats.ScanRejects[reason]++
	c.logReject("scan", reason, tri, c.stats.ScanRejects[reason])
}

func (c *Calculator) addExecReject(reason string, tri *Triangle) {
	if reason == "" {
		reason = "unknown"
	}
	c.stats.ExecRejects[reason]++
	c.logReject("exec", reason, tri, c.stats.ExecRejects[reason])
}

func (c *Calculator) logOpportunity(opp ExecutableOpportunity) {
	c.log.Printf(
		"[ARB] %s→%s→%s | est=%.4f%% | ideal=%.4f%% | rounded=%.4f%% | real=%.4f%% | start=%.2f USDT | minStart=%.2f USDT | final=%.4f USDT | profit=%.4f USDT | symbol=%s",
		opp.Triangle.A,
		opp.Triangle.B,
		opp.Triangle.C,
		opp.EstimatedPct*100,
		opp.IdealProfitPct*100,
		opp.RoundedProfitPct*100,
		opp.ProfitPct*100,
		opp.StartUSDT,
		opp.MinStartUSDT,
		opp.FinalUSDT,
		opp.ProfitUSDT,
		opp.TriggeredBy,
	)
}

func (c *Calculator) logReject(stage, reason string, tri *Triangle, count int64) {
	if c.cfg.LogMode != LogDebug {
		return
	}
	if count != 1 && count != 10 && count != 100 && count%1000 != 0 {
		return
	}
	if tri == nil {
		c.log.Printf("[REJECT] stage=%s reason=%s count=%d", stage, reason, count)
		return
	}
	c.log.Printf("[REJECT] stage=%s reason=%s count=%d tri=%s->%s->%s", stage, reason, count, tri.A, tri.B, tri.C)
}

func (c *Calculator) logStats() {
	c.log.Printf(
		"[STATS] ticks=%d triangles_seen=%d cand=%d exec=%d pos=%d neg=%d logged=%d | scan_rejects={%s} | exec_rejects={%s}",
		c.stats.Ticks,
		c.stats.TrianglesSeen,
		c.stats.Candidates,
		c.stats.Opportunities,
		c.stats.Positive,
		c.stats.Negative,
		c.stats.Logged,
		formatCounts(c.stats.ScanRejects),
		formatCounts(c.stats.ExecRejects),
	)
	c.logTopTriangleMetrics(5)
}

func (c *Calculator) logTopTriangleMetrics(n int) {
	if len(c.stats.TriangleMetrics) == 0 {
		return
	}
	items := make([]*TriangleMetrics, 0, len(c.stats.TriangleMetrics))
	for _, m := range c.stats.TriangleMetrics {
		items = append(items, m)
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].BestIdealPct == items[j].BestIdealPct {
			return items[i].Seen > items[j].Seen
		}
		return items[i].BestIdealPct > items[j].BestIdealPct
	})
	if n > len(items) {
		n = len(items)
	}
	parts := make([]string, 0, n)
	for _, m := range items[:n] {
		parts = append(parts, fmt.Sprintf(
			"%s seen=%d best(est=%.4f%% ideal=%.4f%% rounded=%.4f%% final=%.4f%%) last(ideal=%.4f%% rounded=%.4f%% final=%.4f%%)",
			m.Key,
			m.Seen,
			m.BestEstimatedPct*100,
			m.BestIdealPct*100,
			m.BestRoundedPct*100,
			m.BestFinalPct*100,
			m.LastIdealPct*100,
			m.LastRoundedPct*100,
			m.LastFinalPct*100,
		))
	}
	c.log.Printf("[TOP] %s", strings.Join(parts, " | "))
}

func formatCounts(m map[string]int64) string {
	if len(m) == 0 {
		return "none"
	}
	type kv struct {
		k string
		v int64
	}
	items := make([]kv, 0, len(m))
	for k, v := range m {
		items = append(items, kv{k: k, v: v})
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].v == items[j].v {
			return items[i].k < items[j].k
		}
		return items[i].v < items[j].v
	})
	parts := make([]string, 0, len(items))
	for _, item := range items {
		parts = append(parts, fmt.Sprintf("%s=%d", item.k, item.v))
	}
	return strings.Join(parts, ", ")
}

func (c *Calculator) logModeString() string {
	switch c.cfg.LogMode {
	case LogSilent:
		return "silent"
	case LogDebug:
		return "debug"
	default:
		return "normal"
	}
}

func trimFloat(v float64) string {
	s := fmt.Sprintf("%.2f", v)
	s = strings.TrimRight(s, "0")
	s = strings.TrimRight(s, ".")
	if s == "" {
		return "0"
	}
	return s
}

func ensureDir(path string) {
	if path == "" {
		return
	}
	_ = os.MkdirAll(path, 0o755)
}
