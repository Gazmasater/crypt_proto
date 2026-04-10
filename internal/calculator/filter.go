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

	idealState, okIdeal := simulateTriangleMode(startUSDT, cand.Triangle, cand.Quotes, true, true)
	roundedState, okRounded := simulateTriangleMode(startUSDT, cand.Triangle, cand.Quotes, false, true)
	finalState, okFinal := simulateTriangleMode(startUSDT, cand.Triangle, cand.Quotes, false, false)
	if !okIdeal || !okRounded || !okFinal {
		return ExecutableOpportunity{}, "simulate_failed", false
	}

	opp := ExecutableOpportunity{
		Triangle:         cand.Triangle,
		Quotes:           cand.Quotes,
		EstimatedPct:     cand.EstimatedPct,
		StartUSDT:        finalState.StartUSDT,
		MinStartUSDT:     minStart,
		FinalUSDT:        finalState.FinalUSDT,
		ProfitUSDT:       finalState.ProfitUSDT,
		ProfitPct:        finalState.ProfitPct,
		TriggeredBy:      cand.TriggeredBy,
		TriggeredAtMS:    cand.TriggeredAtMS,
		IdealFinalUSDT:   idealState.FinalUSDT,
		IdealProfitPct:   idealState.ProfitPct,
		RoundedFinalUSDT: roundedState.FinalUSDT,
		RoundedProfitPct: roundedState.ProfitPct,
	}

	if f.cfg.LogMode == LogDebug {
		log.Printf(
			"[EXEC CMP] %s→%s→%s | est=%.4f%% | ideal=%.4f%% | rounded=%.4f%% | final=%.4f%%",
			cand.Triangle.A,
			cand.Triangle.B,
			cand.Triangle.C,
			cand.EstimatedPct*100,
			opp.IdealProfitPct*100,
			opp.RoundedProfitPct*100,
			opp.ProfitPct*100,
		)
	}

	if opp.ProfitPct < f.cfg.MinProfitPct {
		return opp, "profit_below_threshold", false
	}

	return opp, "", true
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
	return simulateTriangleMode(startUSDT, tri, q, false, false)
}

func simulateTriangleMode(startUSDT float64, tri *Triangle, q [3]queue.Quote, ignoreFees bool, ignoreRounding bool) (ExecutionResult, bool) {
	state := ExecutionResult{StartUSDT: startUSDT}
	amount := startUSDT

	for i := 0; i < 3; i++ {
		nextAmount, notional, ok := simulateLegMode(amount, tri.Legs[i], tri.Rules[i], q[i], ignoreFees, ignoreRounding)
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

func simulateLeg(inputAmount float64, leg LegIndex, rules LegRules, quote queue.Quote) (float64, float64, bool) {
	return simulateLegMode(inputAmount, leg, rules, quote, false, false)
}

func simulateLegMode(inputAmount float64, leg LegIndex, rules LegRules, quote queue.Quote, ignoreFees bool, ignoreRounding bool) (float64, float64, bool) {
	mul := feeMultiplier(rules.Fee)
	if ignoreFees {
		mul = 1
	}

	if leg.IsBuy {
		if quote.Ask <= 0 || quote.AskSize <= 0 || inputAmount <= 0 {
			return 0, 0, false
		}

		qty := inputAmount / quote.Ask
		if !ignoreRounding {
			qty = applyFloorStep(qty, rules.QtyStep)
		}
		if qty <= 0 {
			return 0, 0, false
		}
		if qty > quote.AskSize {
			qty = quote.AskSize
			if !ignoreRounding {
				qty = applyFloorStep(qty, rules.QtyStep)
			}
		}
		if qty <= 0 {
			return 0, 0, false
		}

		notional := qty * quote.Ask
		if !ignoreRounding {
			notional = applyFloorStep(notional, rules.QuoteStep)
			if !passesMinChecks(qty, notional, rules) {
				return 0, 0, false
			}
		}

		outQty := qty * mul
		if !ignoreRounding {
			outQty = applyFloorStep(outQty, rules.QtyStep)
		}
		if outQty <= 0 {
			return 0, 0, false
		}
		return outQty, notional, true
	}

	if quote.Bid <= 0 || quote.BidSize <= 0 || inputAmount <= 0 {
		return 0, 0, false
	}

	qty := inputAmount
	if !ignoreRounding {
		qty = applyFloorStep(qty, rules.QtyStep)
	}
	if qty <= 0 {
		return 0, 0, false
	}
	if qty > quote.BidSize {
		qty = quote.BidSize
		if !ignoreRounding {
			qty = applyFloorStep(qty, rules.QtyStep)
		}
	}
	if qty <= 0 {
		return 0, 0, false
	}

	notional := qty * quote.Bid
	if !ignoreRounding {
		notional = applyFloorStep(notional, rules.QuoteStep)
		if !passesMinChecks(qty, notional, rules) {
			return 0, 0, false
		}
	}

	outQuote := notional * mul
	if !ignoreRounding {
		outQuote = applyFloorStep(outQuote, rules.QuoteStep)
	}
	if outQuote <= 0 {
		return 0, 0, false
	}
	return outQuote, notional, true
}
