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

		// A -> B
		var B string
		if m1.Base == anchor {
			B = m1.Quote
		} else if m1.Quote == anchor {
			B = m1.Base
		} else {
			continue
		}

		if common.IsStable(B) {
			continue
		}

		for _, m2 := range markets {
			if !m2.EnableTrading {
				continue
			}

			// B -> C
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

			if common.IsStable(C) {
				continue
			}

			// проверяем замыкание C -> A
			l3, ok := common.FindLeg(C, anchor, markets)
			if !ok {
				continue
			}

			l1, ok1 := common.FindLeg(anchor, B, markets)
			l2, ok2 := common.FindLeg(B, C, markets)
			if !ok1 || !ok2 {
				continue
			}

			// ===== дедуп =====
			key := common.CanonicalKey(anchor, B, C)
			if seen[key] {
				continue
			}
			seen[key] = true

			result = append(result, common.NewTriangle(
				anchor,
				B,
				C,
				l1,
				l2,
				l3,
			))
		}
	}

	return result
}
