package calculator

import (
	"math"
	"time"

	"crypt_proto/internal/queue"
)

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

	var minTS int64
	var maxTS int64

	for i, leg := range tri.Legs {
		quote, ok := s.mem.Get("KuCoin", leg.Symbol)
		if !ok {
			return ScanCandidate{}, false
		}

		if !quoteLooksUsable(quote) {
			return ScanCandidate{}, false
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

	// Важный фикс:
	// раньше отбрасывание часто делалось по схеме now - quote.Timestamp для каждой ноги.
	// Для треугольника это слишком жёстко: одна нога может не обновляться несколько секунд,
	// но три котировки всё ещё согласованы между собой. Поэтому проверяем:
	// 1) разброс времени между ногами,
	// 2) очень мягкий абсолютный idle-фильтр на случай совсем мёртвого рынка.
	if minTS > 0 && maxTS > 0 {
		if maxTS-minTS > 2500 {
			return ScanCandidate{}, false
		}

		nowMS := time.Now().UnixMilli()
		if nowMS-maxTS > 30000 {
			return ScanCandidate{}, false
		}
	}

	maxStart := maxStartUSDT(tri, q)
	if maxStart < minVolumeUSDT {
		return ScanCandidate{}, false
	}

	// На scan-этапе НЕ режем кандидата по отрицательной оценке.
	// Здесь нужен только быстрый coarse filter: котировки валидны и хватает top-of-book ликвидности.
	// Реальные minQty/minNotional/fee/profit потом проверит ExecutorFilter.
	estPct, ok := estimateProfitPct(tri, q)
	if !ok {
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

func maxStartUSDT(tri *Triangle, q [3]queue.Quote) float64 {
	// Обратное распространение ограничений ликвидности от 3-й ноги к старту в USDT.
	// Это корректнее старой приближённой оценки, где единицы измерения местами смешивались.

	// 3-я нога всегда должна вернуть нас в USDT.
	// Максимум входа в 3-ю ногу ограничен доступным объёмом на best bid/ask.
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

	return floorToStep(maxInLeg1, searchStepUSDT)
}

func reverseInputLimit(maxOutput float64, isBuy bool, q queue.Quote) float64 {
	if maxOutput <= 0 {
		return 0
	}

	if isBuy {
		// input = quote, output = base
		if q.Ask <= 0 || q.AskSize <= 0 {
			return 0
		}
		if maxOutput > q.AskSize {
			maxOutput = q.AskSize
		}
		return q.Ask * maxOutput
	}

	// input = base, output = quote
	if q.BidSize <= 0 {
		return 0
	}
	if maxOutput/q.Bid > q.BidSize && q.Bid > 0 {
		return q.BidSize
	}
	if q.Bid <= 0 {
		return 0
	}
	return maxOutput / q.Bid
}
