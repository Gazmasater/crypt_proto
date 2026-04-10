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

	idealFinalUSDT, ok := simulateTriangleMode(startUSDT, cand.Triangle, cand.Quotes, true, true, false)
	if !ok {
		return ExecutableOpportunity{}, "simulate_ideal_failed", false
	}

	roundedFinalUSDT, ok := simulateTriangleMode(startUSDT, cand.Triangle, cand.Quotes, true, false, false)
	if !ok {
		return ExecutableOpportunity{}, "simulate_rounded_failed", false
	}

	state, ok := simulateTriangle(startUSDT, cand.Triangle, cand.Quotes)
	if !ok {
		return ExecutableOpportunity{}, "simulate_failed", false
	}

	idealProfitPct := (idealFinalUSDT / startUSDT) - 1.0
	roundedProfitPct := (roundedFinalUSDT / startUSDT) - 1.0

	if f.cfg.LogMode == LogDebug {
		log.Printf(
			"[EXEC CMP] %s→%s→%s | est=%.4f%% | ideal=%.4f%% | rounded=%.4f%% | final=%.4f%%",
			cand.Triangle.A,
			cand.Triangle.B,
			cand.Triangle.C,
			cand.EstimatedPct*100,
			idealProfitPct*100,
			roundedProfitPct*100,
			state.ProfitPct*100,
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
		ProfitUSDT:        state.ProfitUSDT,
		ProfitPct:         state.ProfitPct,
		IdealFinalUSDT:    idealFinalUSDT,
		IdealProfitPct:    idealProfitPct,
		RoundedFinalUSDT:  roundedFinalUSDT,
		RoundedProfitPct:  roundedProfitPct,
		TriggeredBy:       cand.TriggeredBy,
		TriggeredAtMS:     cand.TriggeredAtMS,
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
	state := ExecutionResult{StartUSDT: startUSDT}
	amount := startUSDT

	for i := 0; i < 3; i++ {
		nextAmount, notional, ok := simulateLeg(amount, tri.Legs[i], tri.Rules[i], q[i])
		if !ok {
			return ExecutionResult{}, false
		}
		state.LegNotional[i] = notional
		state.LegAmount[i] = nextAmount
		amount = nextAmount
	}

	state.FinalUSDT = amount
	state.ProfitUSDT = state.FinalUSDT - state.StartUSDT
	if state.StartUSDT <= 0 {
		return ExecutionResult{}, false
	}
	state.ProfitPct = state.ProfitUSDT / state.StartUSDT
	return state, true
}

func simulateTriangleMode(startUSDT float64, tri *Triangle, q [3]queue.Quote, ignoreFees, ignoreRounding, ignoreMinChecks bool) (float64, bool) {
	amount := startUSDT
	for i := 0; i < 3; i++ {
		nextAmount, _, ok := simulateLegMode(amount, tri.Legs[i], tri.Rules[i], q[i], ignoreFees, ignoreRounding, ignoreMinChecks)
		if !ok {
			return 0, false
		}
		amount = nextAmount
	}
	return amount, true
}

func simulateLeg(inputAmount float64, leg LegIndex, rules LegRules, quote queue.Quote) (float64, float64, bool) {
	return simulateLegMode(inputAmount, leg, rules, quote, false, false, false)
}

func simulateLegMode(inputAmount float64, leg LegIndex, rules LegRules, quote queue.Quote, ignoreFees, ignoreRounding, ignoreMinChecks bool) (float64, float64, bool) {
	mul := feeMultiplier(rules.Fee)
	if ignoreFees {
		mul = 1
	}

	qtyStep := rules.QtyStep
	quoteStep := rules.QuoteStep
	if ignoreRounding {
		qtyStep = 0
		quoteStep = 0
	}

	if leg.IsBuy {
		if quote.Ask <= 0 || quote.AskSize <= 0 || inputAmount <= 0 {
			return 0, 0, false
		}

		qty := applyFloorStep(inputAmount/quote.Ask, qtyStep)
		if qty <= 0 {
			return 0, 0, false
		}
		if qty > quote.AskSize {
			qty = applyFloorStep(quote.AskSize, qtyStep)
		}
		if qty <= 0 {
			return 0, 0, false
		}

		notional := applyFloorStep(qty*quote.Ask, quoteStep)
		if !ignoreMinChecks && !passesMinChecks(qty, notional, rules) {
			return 0, 0, false
		}

		outQty := applyFloorStep(qty*mul, qtyStep)
		if outQty <= 0 {
			return 0, 0, false
		}
		return outQty, notional, true
	}

	if quote.Bid <= 0 || quote.BidSize <= 0 || inputAmount <= 0 {
		return 0, 0, false
	}

	qty := applyFloorStep(inputAmount, qtyStep)
	if qty <= 0 {
		return 0, 0, false
	}
	if qty > quote.BidSize {
		qty = applyFloorStep(quote.BidSize, qtyStep)
	}
	if qty <= 0 {
		return 0, 0, false
	}

	notional := applyFloorStep(qty*quote.Bid, quoteStep)
	if !ignoreMinChecks && !passesMinChecks(qty, notional, rules) {
		return 0, 0, false
	}

	outQuote := applyFloorStep(notional*mul, quoteStep)
	if outQuote <= 0 {
		return 0, 0, false
	}
	return outQuote, notional, true
}





package calculator

import (
	"fmt"
	"math"
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
	StartUSDT   float64
	MinStart    float64
	FinalUSDT   float64
	ProfitUSDT  float64
	ProfitPct   float64
	LegNotional [3]float64
	LegAmount   [3]float64
}

type ExecutableOpportunity struct {
	Triangle         *Triangle
	Quotes           [3]queue.Quote
	EstimatedPct     float64
	StartUSDT        float64
	MinStartUSDT     float64
	FinalUSDT        float64
	ProfitUSDT       float64
	ProfitPct        float64
	IdealFinalUSDT   float64
	IdealProfitPct   float64
	RoundedFinalUSDT float64
	RoundedProfitPct float64
	TriggeredBy      string
	TriggeredAtMS    int64
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

	Positive int64
	Negative int64
	Logged   int64

	ScanRejects map[string]int64
	ExecRejects map[string]int64
}

func feeMultiplier(fee float64) float64 {
	if fee > 0 && fee < 1 {
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
