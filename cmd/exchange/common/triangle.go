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
}

func TriangleKey(a, b, c string) string {
	x := []string{b, c}
	sort.Strings(x)
	return a + "|" + x[0] + "|" + x[1]
}
