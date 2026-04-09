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
	if triggeredAt == 0 {
		triggeredAt = time.Now().UnixMilli()
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
		quote, ok := s.mem.GetLatestBefore("KuCoin", leg.Symbol, triggeredAt, s.cfg.QuoteAgeMaxMS)
		if !ok {
			return ScanResult{Candidate: cand, Reject: rejectNoQuote(i), OK: false}
		}

		if !quoteLooksUsable(quote) {
			return ScanResult{Candidate: cand, Reject: rejectBadQuote(i), OK: false}
		}

		q[i] = quote
		if quote.Timestamp <= 0 {
			continue
		}
		if minTS == 0 || quote.Timestamp < minTS {
			minTS = quote.Timestamp
		}
		if quote.Timestamp > maxTS {
			maxTS = quote.Timestamp
		}
	}

	if minTS > 0 && maxTS > 0 && s.cfg.QuoteAgeMaxMS > 0 {
		if maxTS-minTS > s.cfg.QuoteAgeMaxMS {
			return ScanResult{Candidate: cand, Reject: "quote_skew_too_large", OK: false}
		}

		for i, quote := range q {
			if triggeredAt-quote.Timestamp > s.cfg.QuoteAgeMaxMS {
				return ScanResult{Candidate: cand, Reject: rejectStaleQuote(i), OK: false}
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
