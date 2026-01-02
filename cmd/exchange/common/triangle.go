package common

type Triangle struct {
	A string
	B string
	C string

	Leg1 string
	Leg2 string
	Leg3 string

	BaseMin1  float64
	QuoteMin1 float64
	BaseInc1  float64
	QuoteInc1 float64
	PriceInc1 float64

	BaseMin2  float64
	QuoteMin2 float64
	BaseInc2  float64
	QuoteInc2 float64
	PriceInc2 float64

	BaseMin3  float64
	QuoteMin3 float64
	BaseInc3  float64
	QuoteInc3 float64
	PriceInc3 float64
}

func NewTriangle(A, B, C string, l1, l2, l3 Market) Triangle {
	return Triangle{
		A: A,
		B: B,
		C: C,

		Leg1: ResolveSide(A, B, l1) + " " + l1.Base + "/" + l1.Quote,
		Leg2: ResolveSide(B, C, l2) + " " + l2.Base + "/" + l2.Quote,
		Leg3: ResolveSide(C, A, l3) + " " + l3.Base + "/" + l3.Quote,

		BaseMin1:  l1.BaseMinSize,
		QuoteMin1: l1.QuoteMinSize,
		BaseInc1:  l1.BaseIncrement,
		QuoteInc1: l1.QuoteIncrement,
		PriceInc1: l1.PriceIncrement,

		BaseMin2:  l2.BaseMinSize,
		QuoteMin2: l2.QuoteMinSize,
		BaseInc2:  l2.BaseIncrement,
		QuoteInc2: l2.QuoteIncrement,
		PriceInc2: l2.PriceIncrement,

		BaseMin3:  l3.BaseMinSize,
		QuoteMin3: l3.QuoteMinSize,
		BaseInc3:  l3.BaseIncrement,
		QuoteInc3: l3.QuoteIncrement,
		PriceInc3: l3.PriceIncrement,
	}
}
