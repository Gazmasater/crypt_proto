calc.go

package calculator

import (
	"fmt"
	"io"
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

	f, err := os.OpenFile("arb_opportunities.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		log.Fatalf("failed to open log: %v", err)
	}

	mw := io.MultiWriter(os.Stdout, f)
	c := &Calculator{
		mem:     mem,
		scanner: NewScanner(mem, triangles, cfg),
		filter:  NewExecutorFilter(cfg),
		log:     log.New(mw, "", log.LstdFlags),
		cfg:     cfg,
		stats: Stats{
			ScanRejects: make(map[string]int64),
			ExecRejects: make(map[string]int64),
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
				if !ok {
					c.addExecReject(reason, res.Candidate.Triangle)
					continue
				}

				c.stats.Executable++
				if opp.ProfitPct > 0 {
					c.stats.Positive++
				} else {
					c.stats.Negative++
				}

				if c.shouldLogOpportunity(opp) {
					c.stats.LoggedOpportunities++
					c.logOpportunity(opp)
				}
			}
		case <-ticker.C:
			if c.cfg.LogMode != LogSilent {
				c.logStats()
			}
		}
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

func (c *Calculator) shouldLogOpportunity(opp ExecutableOpportunity) bool {
	switch c.cfg.LogMode {
	case LogSilent:
		return false
	case LogDebug:
		return true
	default:
		return opp.ProfitPct > 0
	}
}

func (c *Calculator) logOpportunity(opp ExecutableOpportunity) {
	c.log.Printf(
		"[ARB] %s→%s→%s | est=%.4f%% | real=%.4f%% | start=%.2f USDT | minStart=%.2f USDT | final=%.4f USDT | profit=%.4f USDT | symbol=%s",
		opp.Triangle.A,
		opp.Triangle.B,
		opp.Triangle.C,
		opp.EstimatedPct*100,
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
		"[STATS] ticks=%d triangles=%d cand=%d exec=%d pos=%d neg=%d logged=%d | scan_rejects={%s} | exec_rejects={%s}",
		c.stats.Ticks,
		c.stats.TrianglesSeen,
		c.stats.Candidates,
		c.stats.Executable,
		c.stats.Positive,
		c.stats.Negative,
		c.stats.LoggedOpportunities,
		formatCounts(c.stats.ScanRejects),
		formatCounts(c.stats.ExecRejects),
	)
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



types.go

type Stats struct {
    Ticks           int64
    TrianglesSeen   int64
    Candidates      int64
    Opportunities   int64

    Positive        int64
    Negative        int64
    Logged          int64

    ScanRejects map[string]int64
    ExecRejects map[string]int64
}

