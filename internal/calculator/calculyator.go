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
	mem      *queue.MemoryStore
	scanner  *Scanner
	filter   *ExecutorFilter
	oppLog   *log.Logger
	debugLog *log.Logger
	stats    *pipelineStats
	lastStat time.Time
}

func NewCalculator(mem *queue.MemoryStore, triangles []*Triangle) *Calculator {
	oppFile, err := os.OpenFile("arb_opportunities.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		log.Fatalf("failed to open opportunities log: %v", err)
	}
	debugFile, err := os.OpenFile("arb_debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		log.Fatalf("failed to open debug log: %v", err)
	}

	oppMW := io.MultiWriter(os.Stdout, oppFile)
	debugMW := io.MultiWriter(os.Stdout, debugFile)

	return &Calculator{
		mem:      mem,
		scanner:  NewScanner(mem, triangles),
		filter:   NewExecutorFilter(),
		oppLog:   log.New(oppMW, "", log.LstdFlags),
		debugLog: log.New(debugMW, "", log.LstdFlags),
		stats:    newPipelineStats(),
		lastStat: time.Now(),
	}
}

func (c *Calculator) Run(in <-chan *models.MarketData) {
	c.debugLog.Printf("[CALC] started | triangles indexed=%d | minVolume=%.2f USDT | minProfit=%.4f%% | quoteAgeMax=%dms",
		countTriangles(c.scanner.bySymbol), minVolumeUSDT, minProfitPct*100, maxQuoteAgeMS)

	for md := range in {
		c.mem.Push(md)
		c.stats.Ticks++

		cands, scanRejects, triSeen := c.scanner.CandidatesFor(md.Symbol, md.Timestamp)
		c.stats.TrianglesSeen += uint64(triSeen)
		for reason, n := range scanRejects {
			c.stats.ScanRejects[reason] += uint64(n)
		}

		if len(cands) == 0 {
			c.flushStatsIfNeeded(md)
			continue
		}

		c.stats.Candidates += uint64(len(cands))
		for _, cand := range cands {
			opp, reason, ok := c.filter.Evaluate(cand)
			if !ok {
				c.stats.ExecRejects[reason]++
				continue
			}
			c.stats.Opportunities++
			c.logOpportunity(opp)
		}

		c.flushStatsIfNeeded(md)
	}

	c.flushStats(true)
}

func (c *Calculator) logOpportunity(opp ExecutableOpportunity) {
	c.oppLog.Printf(
		"[ARB] %s→%s→%s | est=%.4f%% | real=%.4f%% | start=%.2f USDT | minStart=%.2f USDT | final=%.4f USDT | profit=%.4f USDT | trigger=%s",
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

func (c *Calculator) flushStatsIfNeeded(md *models.MarketData) {
	if time.Since(c.lastStat) < statsFlushEvery*time.Second {
		return
	}
	c.flushStats(false)
	c.lastStat = time.Now()
	if md != nil {
		c.debugLog.Printf("[TICK] exchange=%s symbol=%s bid=%.8f ask=%.8f bidSize=%.8f askSize=%.8f ts=%d",
			md.Exchange, md.Symbol, md.Bid, md.Ask, md.BidSize, md.AskSize, md.Timestamp)
	}
}

func (c *Calculator) flushStats(final bool) {
	prefix := "[STATS]"
	if final {
		prefix = "[STATS_FINAL]"
	}
	c.debugLog.Printf("%s ticks=%d triangles_seen=%d candidates=%d opportunities=%d | scan_rejects={%s} | exec_rejects={%s}",
		prefix,
		c.stats.Ticks,
		c.stats.TrianglesSeen,
		c.stats.Candidates,
		c.stats.Opportunities,
		formatReasonMap(c.stats.ScanRejects),
		formatReasonMap(c.stats.ExecRejects),
	)
}

func formatReasonMap(m map[string]uint64) string {
	if len(m) == 0 {
		return "none"
	}
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%d", k, m[k]))
	}
	return strings.Join(parts, ", ")
}

func countTriangles(bySymbol map[string][]*Triangle) int {
	seen := make(map[*Triangle]struct{})
	for _, tris := range bySymbol {
		for _, tri := range tris {
			seen[tri] = struct{}{}
		}
	}
	return len(seen)
}
