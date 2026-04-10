package builder

import (
	"sort"

	"crypt_proto/cmd/exchange/common"
)

// BuildTriangles строит anchored-треугольники от anchor по расширенной логике:
// anchor -> B -> C -> anchor.
//
// В отличие от предыдущей версии, здесь сохраняются ОБЕ валидные ориентации
// для одной и той же пары активов B/C, если они реально дают разные маршруты:
//
//	USDT -> BTC -> ALT -> USDT
//	USDT -> ALT -> BTC -> USDT
//
// Это делает новый генератор ближе к фактическому покрытию старого генератора,
// который из-за недетерминированного обхода map мог сохранять то одну, то другую
// ориентацию. Теперь мы не теряем маршруты вроде USDT->ALGO->BTC->USDT.
func BuildTriangles(markets map[string]common.Market, anchor string) []common.Triangle {
	neighbors := make(map[string]map[string]struct{}, 512)
	for _, m := range markets {
		if !m.EnableTrading || m.Base == "" || m.Quote == "" {
			continue
		}
		if neighbors[m.Base] == nil {
			neighbors[m.Base] = make(map[string]struct{}, 8)
		}
		if neighbors[m.Quote] == nil {
			neighbors[m.Quote] = make(map[string]struct{}, 8)
		}
		neighbors[m.Base][m.Quote] = struct{}{}
		neighbors[m.Quote][m.Base] = struct{}{}
	}

	firstHop := sortedKeys(neighbors[anchor])
	result := make([]common.Triangle, 0, 2048)
	seen := make(map[string]bool, 2048)

	for _, b := range firstHop {
		if b == anchor || common.StableCoins[b] {
			continue
		}

		secondHop := sortedKeys(neighbors[b])
		for _, c := range secondHop {
			if c == "" || c == anchor || c == b || common.StableCoins[c] {
				continue
			}

			// Должно быть реальное замыкание обратно в anchor.
			if _, ok := common.FindLeg(c, anchor, markets); !ok {
				continue
			}

			if t, ok := buildTriangle(anchor, b, c, markets); ok {
				key := triangleRouteKey(t)
				if !seen[key] {
					seen[key] = true
					result = append(result, t)
				}
			}
		}
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].A != result[j].A {
			return result[i].A < result[j].A
		}
		if result[i].B != result[j].B {
			return result[i].B < result[j].B
		}
		return result[i].C < result[j].C
	})

	return result
}

func buildTriangle(anchor, b, c string, markets map[string]common.Market) (common.Triangle, bool) {
	l1, ok1 := common.FindLeg(anchor, b, markets)
	l2, ok2 := common.FindLeg(b, c, markets)
	l3, ok3 := common.FindLeg(c, anchor, markets)
	if !ok1 || !ok2 || !ok3 {
		return common.Triangle{}, false
	}

	side1 := common.ResolveSide(anchor, b, l1)
	side2 := common.ResolveSide(b, c, l2)
	side3 := common.ResolveSide(c, anchor, l3)
	if side1 == "" || side2 == "" || side3 == "" {
		return common.Triangle{}, false
	}

	return common.Triangle{
		A: anchor,
		B: b,
		C: c,

		Leg1: side1 + " " + l1.Base + "/" + l1.Quote,
		Leg2: side2 + " " + l2.Base + "/" + l2.Quote,
		Leg3: side3 + " " + l3.Base + "/" + l3.Quote,

		Step1:        l1.BaseIncrement,
		MinQty1:      l1.BaseMinSize,
		MinNotional1: l1.MinNotional,
		Step2:        l2.BaseIncrement,
		MinQty2:      l2.BaseMinSize,
		MinNotional2: l2.MinNotional,
		Step3:        l3.BaseIncrement,
		MinQty3:      l3.BaseMinSize,
		MinNotional3: l3.MinNotional,

		Leg1Symbol:      l1.Symbol,
		Leg1Side:        side1,
		Leg1Base:        l1.Base,
		Leg1Quote:       l1.Quote,
		Leg1QtyStep:     l1.BaseIncrement,
		Leg1QuoteStep:   l1.QuoteIncrement,
		Leg1PriceStep:   l1.PriceIncrement,
		Leg1MinQty:      l1.BaseMinSize,
		Leg1MinQuote:    l1.QuoteMinSize,
		Leg1MinNotional: l1.MinNotional,

		Leg2Symbol:      l2.Symbol,
		Leg2Side:        side2,
		Leg2Base:        l2.Base,
		Leg2Quote:       l2.Quote,
		Leg2QtyStep:     l2.BaseIncrement,
		Leg2QuoteStep:   l2.QuoteIncrement,
		Leg2PriceStep:   l2.PriceIncrement,
		Leg2MinQty:      l2.BaseMinSize,
		Leg2MinQuote:    l2.QuoteMinSize,
		Leg2MinNotional: l2.MinNotional,

		Leg3Symbol:      l3.Symbol,
		Leg3Side:        side3,
		Leg3Base:        l3.Base,
		Leg3Quote:       l3.Quote,
		Leg3QtyStep:     l3.BaseIncrement,
		Leg3QuoteStep:   l3.QuoteIncrement,
		Leg3PriceStep:   l3.PriceIncrement,
		Leg3MinQty:      l3.BaseMinSize,
		Leg3MinQuote:    l3.QuoteMinSize,
		Leg3MinNotional: l3.MinNotional,
	}, true
}

func sortedKeys(m map[string]struct{}) []string {
	if len(m) == 0 {
		return nil
	}
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func triangleRouteKey(t common.Triangle) string {
	return t.A + "|" + t.B + "|" + t.C + "|" + t.Leg1Symbol + "|" + t.Leg2Symbol + "|" + t.Leg3Symbol + "|" + t.Leg1Side + "|" + t.Leg2Side + "|" + t.Leg3Side
}
