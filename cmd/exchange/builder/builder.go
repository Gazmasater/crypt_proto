package builder

import "crypt_proto/cmd/exchange/common"

func BuildTriangles(
	markets map[string]common.Market,
	anchor string,
) []common.Triangle {
	result := make([]common.Triangle, 0, 1024)
	seen := make(map[string]bool, 1024)

	for _, m1 := range markets {
		if !m1.EnableTrading {
			continue
		}

		var b string
		switch {
		case m1.Base == anchor:
			b = m1.Quote
		case m1.Quote == anchor:
			b = m1.Base
		default:
			continue
		}

		if common.StableCoins[b] {
			continue
		}

		for _, m2 := range markets {
			if !m2.EnableTrading {
				continue
			}

			var c string
			switch {
			case m2.Base == b:
				c = m2.Quote
			case m2.Quote == b:
				c = m2.Base
			default:
				continue
			}

			if c == anchor || c == b {
				continue
			}

			if common.StableCoins[c] {
				continue
			}

			l1, ok1 := common.FindLeg(anchor, b, markets)
			l2, ok2 := common.FindLeg(b, c, markets)
			l3, ok3 := common.FindLeg(c, anchor, markets)
			if !ok1 || !ok2 || !ok3 {
				continue
			}

			key := common.TriangleKey(anchor, b, c)
			if seen[key] {
				continue
			}
			seen[key] = true

			side1 := common.ResolveSide(anchor, b, l1)
			side2 := common.ResolveSide(b, c, l2)
			side3 := common.ResolveSide(c, anchor, l3)

			t := common.Triangle{
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
			}

			result = append(result, t)
		}
	}

	return result
}
