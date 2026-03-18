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

			t := common.Triangle{
				A: anchor,
				B: b,
				C: c,

				Leg1: common.ResolveSide(anchor, b, l1) + " " + l1.Base + "/" + l1.Quote,
				Leg2: common.ResolveSide(b, c, l2) + " " + l2.Base + "/" + l2.Quote,
				Leg3: common.ResolveSide(c, anchor, l3) + " " + l3.Base + "/" + l3.Quote,

				Step1:        l1.BaseIncrement,
				MinQty1:      l1.BaseMinSize,
				MinNotional1: l1.MinNotional,

				Step2:        l2.BaseIncrement,
				MinQty2:      l2.BaseMinSize,
				MinNotional2: l2.MinNotional,

				Step3:        l3.BaseIncrement,
				MinQty3:      l3.BaseMinSize,
				MinNotional3: l3.MinNotional,
			}

			result = append(result, t)
		}
	}

	return result
}
