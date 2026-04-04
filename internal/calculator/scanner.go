package calculator

import (
	"fmt"
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
	result := ScanResult{Candidate: ScanCandidate{Triangle: tri, TriggeredBy: triggeredBy, TriggeredAtMS: triggeredAt}}
	var q [3]queue.Quote
	nowMS := time.Now().UnixMilli()

	for i, leg := range tri.Legs {
		quote, ok := s.mem.Get("KuCoin", leg.Symbol)
		if !ok {
			result.Reject = fmt.Sprintf("no_quote_leg_%d", i+1)
			return result
		}
		if s.cfg.QuoteAgeMaxMS > 0 && quote.Timestamp > 0 && nowMS-quote.Timestamp > s.cfg.QuoteAgeMaxMS {
			result.Reject = fmt.Sprintf("stale_quote_leg_%d", i+1)
			return result
		}
		q[i] = quote
	}
	result.Candidate.Quotes = q

	maxStart := maxStartUSDT(tri, q, s.cfg.SearchStepUSDT)
	result.Candidate.MaxStartUSDT = maxStart
	if maxStart < s.cfg.MinVolumeUSDT {
		result.Reject = fmt.Sprintf("max_start_lt_%.2f", s.cfg.MinVolumeUSDT)
		return result
	}

	estPct, reason, ok := estimateProfitPct(tri, q)
	result.Candidate.EstimatedPct = estPct
	if !ok {
		result.Reject = reason
		return result
	}
	if estPct < 0 {
		result.Reject = "estimated_negative"
		return result
	}

	result.OK = true
	return result
}

func estimateProfitPct(tri *Triangle, q [3]queue.Quote) (float64, string, bool) {
	amount := 1.0
	for i := 0; i < 3; i++ {
		leg := tri.Legs[i]
		quote := q[i]
		if leg.IsBuy {
			if quote.Ask <= 0 || quote.AskSize <= 0 {
				return 0, "zero_liquidity", false
			}
			amount = (amount / quote.Ask) * feeMultiplier(tri.Rules[i].Fee)
			continue
		}
		if quote.Bid <= 0 || quote.BidSize <= 0 {
			return 0, "zero_liquidity", false
		}
		amount = (amount * quote.Bid) * feeMultiplier(tri.Rules[i].Fee)
	}
	return amount - 1.0, "", true
}

func maxStartUSDT(tri *Triangle, q [3]queue.Quote, searchStep float64) float64 {
	allowedNext := math.MaxFloat64

	for i := 2; i >= 0; i-- {
		leg := tri.Legs[i]
		quote := q[i]

		if leg.IsBuy {
			if quote.Ask <= 0 || quote.AskSize <= 0 {
				return 0
			}
			maxIn := quote.Ask * quote.AskSize
			if allowedNext != math.MaxFloat64 {
				needIn := allowedNext * quote.Ask / feeMultiplier(tri.Rules[i].Fee)
				if needIn < maxIn {
					maxIn = needIn
				}
			}
			allowedNext = maxIn
			continue
		}

		if quote.Bid <= 0 || quote.BidSize <= 0 {
			return 0
		}
		maxIn := quote.BidSize
		if allowedNext != math.MaxFloat64 {
			needIn := allowedNext / (quote.Bid * feeMultiplier(tri.Rules[i].Fee))
			if needIn < maxIn {
				maxIn = needIn
			}
		}
		allowedNext = maxIn
	}

	if allowedNext == math.MaxFloat64 || allowedNext <= 0 {
		return 0
	}
	return floorToStep(allowedNext, searchStep)
}
