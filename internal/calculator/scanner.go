package calculator

import "crypt_proto/internal/queue"

type Scanner struct {
	mem      *queue.MemoryStore
	bySymbol map[string][]*Triangle
}

func NewScanner(mem *queue.MemoryStore, triangles []*Triangle) *Scanner {
	bySymbol := make(map[string][]*Triangle, 1024)
	for _, t := range triangles {
		for _, leg := range t.Legs {
			if leg.Symbol == "" {
				continue
			}
			bySymbol[leg.Symbol] = append(bySymbol[leg.Symbol], t)
		}
	}
	return &Scanner{mem: mem, bySymbol: bySymbol}
}

func (s *Scanner) CandidatesFor(mdSymbol string, triggeredAt int64) []ScanCandidate {
	tris := s.bySymbol[mdSymbol]
	if len(tris) == 0 {
		return nil
	}

	out := make([]ScanCandidate, 0, len(tris))
	for _, tri := range tris {
		cand, ok := s.scanTriangle(tri, mdSymbol, triggeredAt)
		if !ok {
			continue
		}
		out = append(out, cand)
	}
	return out
}

func (s *Scanner) scanTriangle(tri *Triangle, triggeredBy string, triggeredAt int64) (ScanCandidate, bool) {
	var q [3]queue.Quote
	for i, leg := range tri.Legs {
		quote, ok := s.mem.Get("KuCoin", leg.Symbol)
		if !ok {
			return ScanCandidate{}, false
		}
		q[i] = quote
	}

	maxStart := maxStartUSDT(tri, q)
	if maxStart < minVolumeUSDT {
		return ScanCandidate{}, false
	}

	estPct, ok := estimateProfitPct(tri, q)
	if !ok || estPct < minProfitPct {
		return ScanCandidate{}, false
	}

	return ScanCandidate{
		Triangle:      tri,
		Quotes:        q,
		EstimatedPct:  estPct,
		MaxStartUSDT:  maxStart,
		TriggeredBy:   triggeredBy,
		TriggeredAtMS: triggeredAt,
	}, true
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

func maxStartUSDT(tri *Triangle, q [3]queue.Quote) float64 {
	var usdtLimits [3]float64

	if tri.Legs[0].IsBuy {
		if q[0].Ask <= 0 || q[0].AskSize <= 0 {
			return 0
		}
		usdtLimits[0] = q[0].Ask * q[0].AskSize
	} else {
		if q[0].Bid <= 0 || q[0].BidSize <= 0 {
			return 0
		}
		usdtLimits[0] = q[0].Bid * q[0].BidSize
	}

	if tri.Legs[1].IsBuy {
		if q[1].Ask <= 0 || q[1].AskSize <= 0 || q[2].Bid <= 0 {
			return 0
		}
		usdtLimits[1] = q[1].Ask * q[1].AskSize * q[2].Bid
	} else {
		if q[1].Bid <= 0 || q[1].BidSize <= 0 || q[2].Bid <= 0 {
			return 0
		}
		usdtLimits[1] = q[1].BidSize * q[2].Bid
	}

	if q[2].Bid <= 0 || q[2].BidSize <= 0 {
		return 0
	}
	usdtLimits[2] = q[2].Bid * q[2].BidSize

	maxUSDT := usdtLimits[0]
	if usdtLimits[1] < maxUSDT {
		maxUSDT = usdtLimits[1]
	}
	if usdtLimits[2] < maxUSDT {
		maxUSDT = usdtLimits[2]
	}
	if maxUSDT <= 0 {
		return 0
	}
	return floorToStep(maxUSDT, searchStepUSDT)
}
