cal.go

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
		"[ARB] %s→%s→%s | est=%.4f%% | ideal=%.4f%% | rounded=%.4f%% | strict=%.4f%% | tails=%.4f%% | start=%.2f USDT | minStart=%.2f USDT | strict_final=%.8f USDT | total=%.8f USDT | tail_usdt=%.8f | profit=%.8f USDT | symbol=%s | tails_map={%s}",
		opp.Triangle.A,
		opp.Triangle.B,
		opp.Triangle.C,
		opp.EstimatedPct*100,
		opp.IdealProfitPct*100,
		opp.RoundedProfitPct*100,
		opp.StrictProfitPct*100,
		opp.ProfitPct*100,
		opp.StartUSDT,
		opp.MinStartUSDT,
		opp.FinalUSDT,
		opp.TotalUSDT,
		opp.TailValueUSDT,
		opp.ProfitUSDT,
		opp.TriggeredBy,
		formatFloatMap(opp.TailUSDTBreakdown),
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
		"[STATS] ticks=%d triangles_seen=%d cand=%d exec=%d | scan_rejects={%s} | exec_rejects={%s}",
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



filter.go

package calculator

import (
	"log"
	"math"

	"crypt_proto/internal/queue"
)

type ExecutorFilter struct {
	cfg Config
}

func NewExecutorFilter(cfg Config) *ExecutorFilter {
	return &ExecutorFilter{cfg: cfg}
}

func (f *ExecutorFilter) Evaluate(cand ScanCandidate) (ExecutableOpportunity, string, bool) {
	minStart, ok := findMinStartForTriangle(cand.Triangle, cand.Quotes, f.cfg.MinVolumeUSDT, cand.MaxStartUSDT, f.cfg.SearchStepUSDT)
	if !ok {
		return ExecutableOpportunity{}, "cannot_find_valid_start", false
	}

	startUSDT := floorToStep(math.Max(f.cfg.MinVolumeUSDT, minStart), f.cfg.SearchStepUSDT)
	if startUSDT < f.cfg.MinVolumeUSDT || startUSDT > cand.MaxStartUSDT {
		return ExecutableOpportunity{}, "max_start_lt_min_volume", false
	}

	state, ok := simulateTriangle(startUSDT, cand.Triangle, cand.Quotes)
	if !ok {
		return ExecutableOpportunity{}, "simulate_failed", false
	}

	if f.cfg.LogMode == LogDebug {
		log.Printf(
			"[EXEC CMP] %s→%s→%s | est=%.4f%% | ideal=%.4f%% | rounded=%.4f%% | strict=%.4f%% | tails=%.4f%% | tail_usdt=%.8f | tails_map={%s}",
			cand.Triangle.A,
			cand.Triangle.B,
			cand.Triangle.C,
			cand.EstimatedPct*100,
			state.IdealProfitPct*100,
			state.RoundedProfitPct*100,
			state.FinalStrictPct*100,
			state.ProfitPct*100,
			state.TailValueUSDT,
			formatFloatMap(state.TailUSDTBreakdown),
		)
	}

	if state.ProfitPct < f.cfg.MinProfitPct {
		return ExecutableOpportunity{}, "profit_below_threshold", false
	}

	return ExecutableOpportunity{
		Triangle:          cand.Triangle,
		Quotes:            cand.Quotes,
		EstimatedPct:      cand.EstimatedPct,
		StartUSDT:         state.StartUSDT,
		MinStartUSDT:      minStart,
		FinalUSDT:         state.FinalUSDT,
		TailValueUSDT:     state.TailValueUSDT,
		TotalUSDT:         state.TotalUSDT,
		ProfitUSDT:        state.ProfitUSDT,
		ProfitPct:         state.ProfitPct,
		StrictProfitUSDT:  state.FinalStrictProfit,
		StrictProfitPct:   state.FinalStrictPct,
		IdealProfitPct:    state.IdealProfitPct,
		RoundedProfitPct:  state.RoundedProfitPct,
		TriggeredBy:       cand.TriggeredBy,
		TriggeredAtMS:     cand.TriggeredAtMS,
		TailAssets:        state.TailAssets,
		TailUSDTBreakdown: state.TailUSDTBreakdown,
	}, "", true
}

func findMinStartForTriangle(tri *Triangle, q [3]queue.Quote, lowerBound, upperBound, searchStep float64) (float64, bool) {
	if upperBound <= 0 || upperBound+1e-12 < lowerBound {
		return 0, false
	}

	lo := math.Max(searchStep, floorToStep(lowerBound, searchStep))
	if lo < lowerBound {
		lo += searchStep
	}
	if lo > upperBound {
		return 0, false
	}

	if _, ok := simulateTriangle(lo, tri, q); ok {
		return lo, true
	}

	step := searchStep
	high := lo
	for high <= upperBound {
		next := floorToStep(high+step, searchStep)
		if next <= high {
			next = floorToStep(high+searchStep, searchStep)
		}
		if next > upperBound {
			break
		}
		if _, ok := simulateTriangle(next, tri, q); ok {
			left, right := high, next
			for right-left > searchStep+1e-12 {
				mid := floorToStep((left+right)/2, searchStep)
				if mid <= left {
					mid = floorToStep(left+searchStep, searchStep)
				}
				if mid >= right {
					break
				}
				if _, ok := simulateTriangle(mid, tri, q); ok {
					right = mid
				} else {
					left = mid
				}
			}
			return right, true
		}
		high = next
		step *= 2
	}

	return 0, false
}

func simulateTriangle(startUSDT float64, tri *Triangle, q [3]queue.Quote) (ExecutionResult, bool) {
	state := ExecutionResult{
		StartUSDT:         startUSDT,
		TailAssets:        make(map[string]float64),
		TailUSDTBreakdown: make(map[string]float64),
	}

	idealFinal, ok := simulateMode(startUSDT, tri, q, true, true, false)
	if !ok {
		return ExecutionResult{}, false
	}
	state.IdealFinalUSDT = idealFinal.TotalUSDT
	state.IdealProfitPct = pctChange(startUSDT, idealFinal.TotalUSDT)

	roundedFinal, ok := simulateMode(startUSDT, tri, q, false, true, false)
	if !ok {
		return ExecutionResult{}, false
	}
	state.RoundedFinalUSDT = roundedFinal.TotalUSDT
	state.RoundedProfitPct = pctChange(startUSDT, roundedFinal.TotalUSDT)

	finalStrict, ok := simulateMode(startUSDT, tri, q, false, false, false)
	if !ok {
		return ExecutionResult{}, false
	}
	state.FinalUSDT = finalStrict.FinalUSDT
	state.FinalStrictUSDT = finalStrict.FinalUSDT
	state.FinalStrictProfit = finalStrict.FinalUSDT - startUSDT
	state.FinalStrictPct = pctChange(startUSDT, finalStrict.FinalUSDT)
	state.LegNotional = finalStrict.LegNotional
	state.LegAmount = finalStrict.LegAmount

	finalWithTails, ok := simulateMode(startUSDT, tri, q, false, false, true)
	if !ok {
		return ExecutionResult{}, false
	}
	state.FinalUSDT = finalWithTails.FinalUSDT
	state.TailValueUSDT = finalWithTails.TailValueUSDT
	state.TotalUSDT = finalWithTails.TotalUSDT
	state.ProfitUSDT = finalWithTails.TotalUSDT - startUSDT
	state.ProfitPct = pctChange(startUSDT, finalWithTails.TotalUSDT)
	state.TailAssets = finalWithTails.TailAssets
	state.TailUSDTBreakdown = finalWithTails.TailUSDTBreakdown
	state.LegNotional = finalWithTails.LegNotional
	state.LegAmount = finalWithTails.LegAmount

	return state, true
}

type simulationModeResult struct {
	FinalUSDT         float64
	TailValueUSDT     float64
	TotalUSDT         float64
	LegNotional       [3]float64
	LegAmount         [3]float64
	TailAssets        map[string]float64
	TailUSDTBreakdown map[string]float64
}

func simulateMode(startUSDT float64, tri *Triangle, q [3]queue.Quote, ignoreRounding, ignoreFees, includeTails bool) (simulationModeResult, bool) {
	amount := startUSDT
	tails := make(map[string]float64)
	result := simulationModeResult{TailAssets: tails, TailUSDTBreakdown: make(map[string]float64)}

	for i := 0; i < 3; i++ {
		nextAmount, notional, legTails, ok := simulateLegDetailed(amount, tri.Legs[i], tri.Rules[i], q[i], ignoreRounding, ignoreFees)
		if !ok {
			return simulationModeResult{}, false
		}
		result.LegNotional[i] = notional
		result.LegAmount[i] = nextAmount
		for asset, tailAmt := range legTails {
			if tailAmt > 0 {
				tails[asset] += tailAmt
			}
		}
		amount = nextAmount
	}

	result.FinalUSDT = amount
	result.TotalUSDT = amount
	if includeTails {
		tailValue := 0.0
		for asset, tailAmt := range tails {
			value, ok := valueAssetInUSDT(asset, tailAmt, tri, q)
			if !ok {
				continue
			}
			if value > 0 {
				result.TailUSDTBreakdown[asset] = value
				tailValue += value
			}
		}
		result.TailValueUSDT = tailValue
		result.TotalUSDT += tailValue
	}
	return result, true
}

func simulateLegDetailed(inputAmount float64, leg LegIndex, rules LegRules, quote queue.Quote, ignoreRounding, ignoreFees bool) (float64, float64, map[string]float64, bool) {
	tails := make(map[string]float64)
	mul := 1.0
	if !ignoreFees {
		mul = feeMultiplier(rules.Fee)
	}

	if leg.IsBuy {
		if quote.Ask <= 0 || quote.AskSize <= 0 || inputAmount <= 0 {
			return 0, 0, nil, false
		}

		rawQty := inputAmount / quote.Ask
		qty := rawQty
		if !ignoreRounding {
			qty = applyFloorStep(rawQty, rules.QtyStep)
		}
		if qty <= 0 {
			return 0, 0, nil, false
		}
		if qty > quote.AskSize {
			qty = quote.AskSize
			if !ignoreRounding {
				qty = applyFloorStep(qty, rules.QtyStep)
			}
		}
		if qty <= 0 {
			return 0, 0, nil, false
		}

		notional := qty * quote.Ask
		if !ignoreRounding {
			notional = applyFloorStep(notional, rules.QuoteStep)
		}
		if notional <= 0 {
			return 0, 0, nil, false
		}
		if !ignoreRounding && !passesMinChecks(qty, notional, rules) {
			return 0, 0, nil, false
		}

		leftoverQuote := inputAmount - notional
		if leftoverQuote > 1e-12 {
			tails[rules.Quote] += leftoverQuote
		}

		grossOut := qty * mul
		tradableOut := grossOut
		if !ignoreRounding {
			tradableOut = applyFloorStep(grossOut, rules.QtyStep)
		}
		if tradableOut <= 0 {
			return 0, 0, nil, false
		}
		if !ignoreRounding && grossOut-tradableOut > 1e-12 {
			tails[rules.Base] += grossOut - tradableOut
		}
		return tradableOut, notional, tails, true
	}

	if quote.Bid <= 0 || quote.BidSize <= 0 || inputAmount <= 0 {
		return 0, 0, nil, false
	}

	rawQty := inputAmount
	qty := rawQty
	if !ignoreRounding {
		qty = applyFloorStep(rawQty, rules.QtyStep)
	}
	if qty <= 0 {
		return 0, 0, nil, false
	}
	if qty > quote.BidSize {
		qty = quote.BidSize
		if !ignoreRounding {
			qty = applyFloorStep(qty, rules.QtyStep)
		}
	}
	if qty <= 0 {
		return 0, 0, nil, false
	}

	if !ignoreRounding && rawQty-qty > 1e-12 {
		tails[rules.Base] += rawQty - qty
	}

	notional := qty * quote.Bid
	if !ignoreRounding {
		notional = applyFloorStep(notional, rules.QuoteStep)
	}
	if notional <= 0 {
		return 0, 0, nil, false
	}
	if !ignoreRounding && !passesMinChecks(qty, notional, rules) {
		return 0, 0, nil, false
	}

	grossOut := notional * mul
	tradableOut := grossOut
	if !ignoreRounding {
		tradableOut = applyFloorStep(grossOut, rules.QuoteStep)
	}
	if tradableOut <= 0 {
		return 0, 0, nil, false
	}
	if !ignoreRounding && grossOut-tradableOut > 1e-12 {
		tails[rules.Quote] += grossOut - tradableOut
	}

	return tradableOut, notional, tails, true
}

func valueAssetInUSDT(asset string, amount float64, tri *Triangle, q [3]queue.Quote) (float64, bool) {
	if amount <= 0 {
		return 0, false
	}
	if asset == tri.A || asset == "USDT" {
		return amount, true
	}

	for i, rules := range tri.Rules {
		quote := q[i]
		if rules.Base == asset && rules.Quote == tri.A {
			if quote.Bid <= 0 {
				return 0, false
			}
			qty := applyFloorStep(amount, rules.QtyStep)
			if qty <= 0 {
				return 0, false
			}
			notional := applyFloorStep(qty*quote.Bid, rules.QuoteStep)
			if notional <= 0 || !passesMinChecks(qty, notional, rules) {
				return 0, false
			}
			return applyFloorStep(notional*feeMultiplier(rules.Fee), rules.QuoteStep), true
		}
		if rules.Base == tri.A && rules.Quote == asset {
			if quote.Ask <= 0 {
				return 0, false
			}
			qtyUSDT := applyFloorStep(amount/quote.Ask, rules.QtyStep)
			if qtyUSDT <= 0 {
				return 0, false
			}
			notional := applyFloorStep(qtyUSDT*quote.Ask, rules.QuoteStep)
			if notional <= 0 || notional > amount+1e-12 || !passesMinChecks(qtyUSDT, notional, rules) {
				return 0, false
			}
			return applyFloorStep(qtyUSDT*feeMultiplier(rules.Fee), rules.QtyStep), true
		}
	}
	return 0, false
}

func pctChange(start, end float64) float64 {
	if start <= 0 {
		return 0
	}
	return (end / start) - 1.0
}



types.go

package calculator

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"crypt_proto/internal/queue"
)

const (
	defaultTakerFee      = 0.001
	defaultMinVolumeUSDT = 50.0
	defaultMinProfitPct  = 0.001
	defaultSearchStep    = 0.01
)

var triangleLegColumns = [3]int{3, 7, 11}

type LogMode int

const (
	LogSilent LogMode = iota
	LogNormal
	LogDebug
)

type Config struct {
	MinVolumeUSDT  float64
	MinProfitPct   float64
	SearchStepUSDT float64
	QuoteAgeMaxMS  int64
	StatsEverySec  int
	LogMode        LogMode
}

func DefaultConfig() Config {
	return Config{
		MinVolumeUSDT:  defaultMinVolumeUSDT,
		MinProfitPct:   defaultMinProfitPct,
		SearchStepUSDT: defaultSearchStep,
		QuoteAgeMaxMS:  2500,
		StatsEverySec:  5,
		LogMode:        LogNormal,
	}
}

type LegIndex struct {
	Key    string
	Symbol string
	IsBuy  bool
}

type LegRules struct {
	Symbol      string
	Side        string
	Base        string
	Quote       string
	QtyStep     float64
	QuoteStep   float64
	PriceStep   float64
	MinQty      float64
	MinQuote    float64
	MinNotional float64
	Fee         float64
}

type Triangle struct {
	A, B, C string
	Legs    [3]LegIndex
	Rules   [3]LegRules
}

type ScanCandidate struct {
	Triangle      *Triangle
	Quotes        [3]queue.Quote
	EstimatedPct  float64
	MaxStartUSDT  float64
	TriggeredBy   string
	TriggeredAtMS int64
}

type ExecutionResult struct {
	StartUSDT          float64
	MinStart           float64
	FinalUSDT          float64
	TailValueUSDT      float64
	TotalUSDT          float64
	ProfitUSDT         float64
	ProfitPct          float64
	StrictProfitUSDT   float64
	StrictProfitPct    float64
	IdealFinalUSDT     float64
	IdealProfitPct     float64
	RoundedFinalUSDT   float64
	RoundedProfitPct   float64
	FinalStrictUSDT    float64
	FinalStrictProfit  float64
	FinalStrictPct     float64
	LegNotional        [3]float64
	LegAmount          [3]float64
	TailAssets         map[string]float64
	TailUSDTBreakdown  map[string]float64
}

type ExecutableOpportunity struct {
	Triangle          *Triangle
	Quotes            [3]queue.Quote
	EstimatedPct      float64
	StartUSDT         float64
	MinStartUSDT      float64
	FinalUSDT         float64
	TailValueUSDT     float64
	TotalUSDT         float64
	ProfitUSDT        float64
	ProfitPct         float64
	StrictProfitUSDT  float64
	StrictProfitPct   float64
	IdealProfitPct    float64
	RoundedProfitPct  float64
	TriggeredBy       string
	TriggeredAtMS     int64
	TailAssets        map[string]float64
	TailUSDTBreakdown map[string]float64
}

type ScanResult struct {
	Candidate ScanCandidate
	Reject    string
	OK        bool
}

type Stats struct {
	Ticks         int64
	TrianglesSeen int64
	Candidates    int64
	Opportunities int64
	ScanRejects   map[string]int64
	ExecRejects   map[string]int64
}

func feeMultiplier(fee float64) float64 {
	if fee >= 0 && fee < 1 {
		return 1 - fee
	}
	return 1 - defaultTakerFee
}

func applyFloorStep(value, step float64) float64 {
	if value <= 0 {
		return 0
	}
	if step <= 0 {
		return value
	}
	return floorToStep(value, step)
}

func floorToStep(value, step float64) float64 {
	if step <= 0 {
		return value
	}
	units := math.Floor((value + 1e-12) / step)
	if units <= 0 {
		return 0
	}
	result := units * step
	precision := decimalsFromStep(step)
	pow := math.Pow10(precision)
	return math.Floor(result*pow+1e-9) / pow
}

func decimalsFromStep(step float64) int {
	if step <= 0 {
		return 8
	}
	s := strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.12f", step), "0"), ".")
	idx := strings.IndexByte(s, '.')
	if idx == -1 {
		return 0
	}
	return len(s) - idx - 1
}

func passesMinChecks(qty, notional float64, rules LegRules) bool {
	if rules.MinQty > 0 && qty+1e-12 < rules.MinQty {
		return false
	}
	if rules.MinQuote > 0 && notional+1e-12 < rules.MinQuote {
		return false
	}
	if rules.MinNotional > 0 && notional+1e-12 < rules.MinNotional {
		return false
	}
	return true
}

func formatFloatMap(m map[string]float64) string {
	if len(m) == 0 {
		return "none"
	}
	type kv struct {
		k string
		v float64
	}
	items := make([]kv, 0, len(m))
	for k, v := range m {
		items = append(items, kv{k: k, v: v})
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].v == items[j].v {
			return items[i].k < items[j].k
		}
		return items[i].v > items[j].v
	})
	parts := make([]string, 0, len(items))
	for _, item := range items {
		parts = append(parts, fmt.Sprintf("%s=%.8f", item.k, item.v))
	}
	return strings.Join(parts, ", ")
}

