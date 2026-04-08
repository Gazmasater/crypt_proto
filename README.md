scanner.go

package calculator

import (
	"math"
	"time"

	"crypt_proto/internal/queue"
)

type Scanner struct {
	mem      *queue.MemoryStore
	bySymbol map[string][]*Triangle
	cfg      Config
}

func NewScanner(mem *queue.MemoryStore, triangles []*Triangle, cfg Config) *Scanner {
	bySymbol := make(map[string][]*Triangle, 1024)
	for _, t := range triangles {
		for _, leg := range t.Legs {
			if leg.Symbol == "" {
				continue
			}
			bySymbol[leg.Symbol] = append(bySymbol[leg.Symbol], t)
		}
	}
	return &Scanner{mem: mem, bySymbol: bySymbol, cfg: cfg}
}

func (s *Scanner) CandidatesFor(mdSymbol string, triggeredAt int64) []ScanResult {
	tris := s.bySymbol[mdSymbol]
	if len(tris) == 0 {
		return nil
	}

	out := make([]ScanResult, 0, len(tris))
	for _, tri := range tris {
		out = append(out, s.scanTriangle(tri, mdSymbol, triggeredAt))
	}
	return out
}

func (s *Scanner) scanTriangle(tri *Triangle, triggeredBy string, triggeredAt int64) ScanResult {
	cand := ScanCandidate{
		Triangle:      tri,
		TriggeredBy:   triggeredBy,
		TriggeredAtMS: triggeredAt,
	}

	var q [3]queue.Quote
	var minTS int64
	var maxTS int64

	for i, leg := range tri.Legs {
		quote, ok := s.mem.Get("KuCoin", leg.Symbol)
		if !ok {
			return ScanResult{Candidate: cand, Reject: rejectNoQuote(i), OK: false}
		}

		if !quoteLooksUsable(quote) {
			return ScanResult{Candidate: cand, Reject: rejectBadQuote(i), OK: false}
		}

		q[i] = quote

		if quote.Timestamp > 0 {
			if minTS == 0 || quote.Timestamp < minTS {
				minTS = quote.Timestamp
			}
			if quote.Timestamp > maxTS {
				maxTS = quote.Timestamp
			}
		}
	}

	if minTS > 0 && maxTS > 0 {
		if s.cfg.QuoteAgeMaxMS > 0 && maxTS-minTS > s.cfg.QuoteAgeMaxMS {
			return ScanResult{Candidate: cand, Reject: "quote_skew_too_large", OK: false}
		}

		nowMS := triggeredAt
		if nowMS == 0 {
			nowMS = time.Now().UnixMilli()
		}
		if nowMS < maxTS {
			nowMS = maxTS
		}

		if s.cfg.QuoteAgeMaxMS > 0 {
			for i, quote := range q {
				if quote.Timestamp <= 0 {
					continue
				}
				if nowMS-quote.Timestamp > s.cfg.QuoteAgeMaxMS {
					return ScanResult{Candidate: cand, Reject: rejectStaleQuote(i), OK: false}
				}
			}
		}
	}

	maxStart := s.maxStartUSDT(tri, q)
	if maxStart <= 0 {
		return ScanResult{Candidate: cand, Reject: "max_start_zero", OK: false}
	}
	if maxStart+1e-12 < s.cfg.MinVolumeUSDT {
		return ScanResult{Candidate: cand, Reject: rejectMaxStart(s.cfg.MinVolumeUSDT), OK: false}
	}

	estPct, ok := estimateProfitPct(tri, q)
	if !ok {
		return ScanResult{Candidate: cand, Reject: "estimate_failed", OK: false}
	}

	cand.Quotes = q
	cand.EstimatedPct = estPct
	cand.MaxStartUSDT = maxStart
	return ScanResult{Candidate: cand, OK: true}
}

func quoteLooksUsable(q queue.Quote) bool {
	if q.Bid <= 0 || q.Ask <= 0 {
		return false
	}
	if q.BidSize <= 0 || q.AskSize <= 0 {
		return false
	}
	if q.Ask < q.Bid {
		return false
	}
	if math.IsNaN(q.Bid) || math.IsNaN(q.Ask) || math.IsNaN(q.BidSize) || math.IsNaN(q.AskSize) {
		return false
	}
	if math.IsInf(q.Bid, 0) || math.IsInf(q.Ask, 0) || math.IsInf(q.BidSize, 0) || math.IsInf(q.AskSize, 0) {
		return false
	}
	return true
}

func estimateProfitPct(tri *Triangle, q [3]queue.Quote) (float64, bool) {
	amount := 1.0
	for i := 0; i < 3; i++ {
		leg := tri.Legs[i]
		quote := q[i]
		if leg.IsBuy {
			if quote.Ask <= 0 {
				return 0, false
			}
			amount = (amount / quote.Ask) * feeMultiplier(tri.Rules[i].Fee)
			continue
		}
		if quote.Bid <= 0 {
			return 0, false
		}
		amount = (amount * quote.Bid) * feeMultiplier(tri.Rules[i].Fee)
	}
	return amount - 1.0, true
}

func (s *Scanner) maxStartUSDT(tri *Triangle, q [3]queue.Quote) float64 {
	var maxInLeg3 float64
	if tri.Legs[2].IsBuy {
		if q[2].Ask <= 0 || q[2].AskSize <= 0 {
			return 0
		}
		maxInLeg3 = q[2].Ask * q[2].AskSize
	} else {
		if q[2].BidSize <= 0 {
			return 0
		}
		maxInLeg3 = q[2].BidSize
	}

	maxInLeg2 := reverseInputLimit(maxInLeg3, tri.Legs[1].IsBuy, q[1])
	if maxInLeg2 <= 0 {
		return 0
	}

	maxInLeg1 := reverseInputLimit(maxInLeg2, tri.Legs[0].IsBuy, q[0])
	if maxInLeg1 <= 0 {
		return 0
	}

	return floorToStep(maxInLeg1, s.cfg.SearchStepUSDT)
}

func reverseInputLimit(maxOutput float64, isBuy bool, q queue.Quote) float64 {
	if maxOutput <= 0 {
		return 0
	}

	if isBuy {
		if q.Ask <= 0 || q.AskSize <= 0 {
			return 0
		}
		if maxOutput > q.AskSize {
			maxOutput = q.AskSize
		}
		return q.Ask * maxOutput
	}

	if q.Bid <= 0 || q.BidSize <= 0 {
		return 0
	}
	if maxOutput/q.Bid > q.BidSize {
		return q.BidSize
	}
	return maxOutput / q.Bid
}

func rejectNoQuote(legIdx int) string {
	switch legIdx {
	case 0:
		return "no_quote_leg_1"
	case 1:
		return "no_quote_leg_2"
	default:
		return "no_quote_leg_3"
	}
}

func rejectBadQuote(legIdx int) string {
	switch legIdx {
	case 0:
		return "bad_quote_leg_1"
	case 1:
		return "bad_quote_leg_2"
	default:
		return "bad_quote_leg_3"
	}
}

func rejectStaleQuote(legIdx int) string {
	switch legIdx {
	case 0:
		return "stale_quote_leg_1"
	case 1:
		return "stale_quote_leg_2"
	default:
		return "stale_quote_leg_3"
	}
}

func rejectMaxStart(minVolume float64) string {
	return "max_start_lt_" + trimFloat(minVolume)
}



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

				c.stats.Opportunities++
				c.logOpportunity(opp)
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
		"[STATS] ticks=%d triangles_seen=%d candidates=%d opportunities=%d | scan_rejects={%s} | exec_rejects={%s}",
		c.stats.Ticks,
		c.stats.TrianglesSeen,
		c.stats.Candidates,
		c.stats.Opportunities,
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



