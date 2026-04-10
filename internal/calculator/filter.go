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
		Triangle:         cand.Triangle,
		Quotes:           cand.Quotes,
		EstimatedPct:     cand.EstimatedPct,
		StartUSDT:        state.StartUSDT,
		MinStartUSDT:     minStart,
		FinalUSDT:        state.FinalUSDT,
		ProfitUSDT:       state.ProfitUSDT,
		ProfitPct:        state.ProfitPct,
		IdealFinalUSDT:   idealFinalUSDT,
		IdealProfitPct:   idealProfitPct,
		RoundedFinalUSDT: roundedFinalUSDT,
		RoundedProfitPct: roundedProfitPct,
		TriggeredBy:      cand.TriggeredBy,
		TriggeredAtMS:    cand.TriggeredAtMS,
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
