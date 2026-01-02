package builder

import "crypt_proto/cmd/exchange/common"

func BuildTriangles(
	markets map[string]common.Market,
	anchor string,
) []common.Triangle {

	result := []common.Triangle{}
	seen := map[string]bool{}

	for _, m1 := range markets {
		if !m1.EnableTrading {
			continue
		}

		// шаг 1: A -> B
		var B string
		if m1.Base == anchor {
			B = m1.Quote
		} else if m1.Quote == anchor {
			B = m1.Base
		} else {
			continue
		}

		// ❌ нельзя стейб в середине
		if common.StableCoins[B] {
			continue
		}

		for _, m2 := range markets {
			if !m2.EnableTrading {
				continue
			}

			// шаг 2: B -> C
			var C string
			if m2.Base == B {
				C = m2.Quote
			} else if m2.Quote == B {
				C = m2.Base
			} else {
				continue
			}

			if C == anchor || C == B {
				continue
			}

			if common.StableCoins[C] {
				continue
			}

			// шаг 3: C -> A должен существовать
			l3, ok := common.FindLeg(C, anchor, markets)
			if !ok {
				continue
			}

			l1, ok1 := common.FindLeg(anchor, B, markets)
			l2, ok2 := common.FindLeg(B, C, markets)
			if !ok1 || !ok2 {
				continue
			}

			// 🔒 дедупликация A-X-Y-A / A-Y-X-A
			key := common.TriangleKey(anchor, B, C)
			if seen[key] {
				continue
			}
			seen[key] = true

			t := common.Triangle{
				A: anchor,
				B: B,
				C: C,

				Leg1: common.ResolveSide(anchor, B, l1) + " " + l1.Base + "/" + l1.Quote,
				Leg2: common.ResolveSide(B, C, l2) + " " + l2.Base + "/" + l2.Quote,
				Leg3: common.ResolveSide(C, anchor, l3) + " " + l3.Base + "/" + l3.Quote,
			}

			result = append(result, t)
		}
	}

	return result
}
