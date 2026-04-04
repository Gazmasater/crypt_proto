package calculator

import (
	"math"

	"crypt_proto/internal/queue"
)

type ExecutorFilter struct{}

func NewExecutorFilter() *ExecutorFilter {
	return &ExecutorFilter{}
}

func (f *ExecutorFilter) Evaluate(cand ScanCandidate) (ExecutableOpportunity, bool) {
	minStart, ok := findMinStartForTriangle(cand.Triangle, cand.Quotes, minVolumeUSDT, cand.MaxStartUSDT)
	if !ok {
		return ExecutableOpportunity{}, false
	}

	startUSDT := floorToStep(math.Max(minVolumeUSDT, minStart), searchStepUSDT)
	if startUSDT < minVolumeUSDT || startUSDT > cand.MaxStartUSDT {
		return ExecutableOpportunity{}, false
	}

	state, ok := simulateTriangle(startUSDT, cand.Triangle, cand.Quotes)
	if !ok || state.ProfitPct < minProfitPct {
		return ExecutableOpportunity{}, false
	}

	return ExecutableOpportunity{
		Triangle:      cand.Triangle,
		Quotes:        cand.Quotes,
		EstimatedPct:  cand.EstimatedPct,
		StartUSDT:     state.StartUSDT,
		MinStartUSDT:  minStart,
		FinalUSDT:     state.FinalUSDT,
		ProfitUSDT:    state.ProfitUSDT,
		ProfitPct:     state.ProfitPct,
		TriggeredBy:   cand.TriggeredBy,
		TriggeredAtMS: cand.TriggeredAtMS,
	}, true
}

func findMinStartForTriangle(tri *Triangle, q [3]queue.Quote, lowerBound, upperBound float64) (float64, bool) {
	if upperBound <= 0 || upperBound+1e-12 < lowerBound {
		return 0, false
	}

	lo := math.Max(searchStepUSDT, floorToStep(lowerBound, searchStepUSDT))
	if lo < lowerBound {
		lo += searchStepUSDT
	}
	if lo > upperBound {
		return 0, false
	}

	if _, ok := simulateTriangle(lo, tri, q); ok {
		return lo, true
	}

	high := lo
	for high <= upperBound {
		high = floorToStep(high+searchStepUSDT, searchStepUSDT)
		if high > upperBound {
			break
		}
		if _, ok := simulateTriangle(high, tri, q); ok {
			left, right := lo, high
			for right-left > searchStepUSDT+1e-12 {
				mid := floorToStep((left+right)/2, searchStepUSDT)
				if mid <= left {
					mid = floorToStep(left+searchStepUSDT, searchStepUSDT)
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
		lo = high
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

func simulateLeg(inputAmount float64, leg LegIndex, rules LegRules, quote queue.Quote) (float64, float64, bool) {
	mul := feeMultiplier(rules.Fee)

	if leg.IsBuy {
		if quote.Ask <= 0 || quote.AskSize <= 0 || inputAmount <= 0 {
			return 0, 0, false
		}

		qty := applyFloorStep(inputAmount/quote.Ask, rules.QtyStep)
		if qty <= 0 {
			return 0, 0, false
		}
		if qty > quote.AskSize {
			qty = applyFloorStep(quote.AskSize, rules.QtyStep)
		}
		if qty <= 0 {
			return 0, 0, false
		}

		notional := applyFloorStep(qty*quote.Ask, rules.QuoteStep)
		if !passesMinChecks(qty, notional, rules) {
			return 0, 0, false
		}

		outQty := applyFloorStep(qty*mul, rules.QtyStep)
		if outQty <= 0 {
			return 0, 0, false
		}
		return outQty, notional, true
	}

	if quote.Bid <= 0 || quote.BidSize <= 0 || inputAmount <= 0 {
		return 0, 0, false
	}

	qty := applyFloorStep(inputAmount, rules.QtyStep)
	if qty <= 0 {
		return 0, 0, false
	}
	if qty > quote.BidSize {
		qty = applyFloorStep(quote.BidSize, rules.QtyStep)
	}
	if qty <= 0 {
		return 0, 0, false
	}

	notional := applyFloorStep(qty*quote.Bid, rules.QuoteStep)
	if !passesMinChecks(qty, notional, rules) {
		return 0, 0, false
	}

	outQuote := applyFloorStep(notional*mul, rules.QuoteStep)
	if outQuote <= 0 {
		return 0, 0, false
	}
	return outQuote, notional, true
}
