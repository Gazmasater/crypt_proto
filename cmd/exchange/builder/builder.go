package builder

import "crypt_proto/cmd/exchange/common"

// BuildTriangles строит треугольники из всех доступных рынков с учётом anchor.
// Пропускаются стейблкоины, кроме anchor.
// Возвращает все варианты: anchor → B → C → anchor и anchor → C → B → anchor
func BuildTriangles(
	markets map[string]common.Market,
	anchor string,
) []common.Triangle {

	var result []common.Triangle

	for _, m1 := range markets {
		if !m1.EnableTrading {
			continue
		}

		var B string
		if m1.Base == anchor {
			B = m1.Quote
		} else if m1.Quote == anchor {
			B = m1.Base
		} else {
			continue
		}

		if common.IsStable(B) && B != anchor {
			continue
		}

		for _, m2 := range markets {
			if !m2.EnableTrading {
				continue
			}

			var C string
			if m2.Base == B {
				C = m2.Quote
			} else if m2.Quote == B {
				C = m2.Base
			} else {
				continue
			}

			if C == anchor || C == B || (common.IsStable(C) && C != anchor) {
				continue
			}

			// Первая нога: anchor → B
			l1, ok1 := common.FindLeg(anchor, B, markets)
			// Вторая нога: B → C
			l2, ok2 := common.FindLeg(B, C, markets)
			// Третья нога: C → anchor
			l3, ok3 := common.FindLeg(C, anchor, markets)

			if ok1 && ok2 && ok3 {
				t := common.NewTriangle(anchor, B, C, l1, l2, l3)
				result = append(result, t)
			}

			// Вариант в обратном порядке: anchor → C → B → anchor
			l1r, ok1r := common.FindLeg(anchor, C, markets)
			l2r, ok2r := common.FindLeg(C, B, markets)
			l3r, ok3r := common.FindLeg(B, anchor, markets)

			if ok1r && ok2r && ok3r {
				t := common.NewTriangle(anchor, C, B, l1r, l2r, l3r)
				result = append(result, t)
			}
		}
	}

	return result
}
