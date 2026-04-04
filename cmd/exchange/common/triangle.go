package common

import "sort"

type Triangle struct {
	A string
	B string
	C string

	Leg1 string
	Leg2 string
	Leg3 string

	Step1        float64
	MinQty1      float64
	MinNotional1 float64

	Step2        float64
	MinQty2      float64
	MinNotional2 float64

	Step3        float64
	MinQty3      float64
	MinNotional3 float64

	Leg1Symbol      string
	Leg1Side        string
	Leg1Base        string
	Leg1Quote       string
	Leg1QtyStep     float64
	Leg1QuoteStep   float64
	Leg1PriceStep   float64
	Leg1MinQty      float64
	Leg1MinQuote    float64
	Leg1MinNotional float64

	Leg2Symbol      string
	Leg2Side        string
	Leg2Base        string
	Leg2Quote       string
	Leg2QtyStep     float64
	Leg2QuoteStep   float64
	Leg2PriceStep   float64
	Leg2MinQty      float64
	Leg2MinQuote    float64
	Leg2MinNotional float64

	Leg3Symbol      string
	Leg3Side        string
	Leg3Base        string
	Leg3Quote       string
	Leg3QtyStep     float64
	Leg3QuoteStep   float64
	Leg3PriceStep   float64
	Leg3MinQty      float64
	Leg3MinQuote    float64
	Leg3MinNotional float64
}

func TriangleKey(a, b, c string) string {
	x := []string{b, c}
	sort.Strings(x)
	return a + "|" + x[0] + "|" + x[1]
}
