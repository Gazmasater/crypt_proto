package common

import "sort"

type Triangle struct {
	A string
	B string
	C string

	Leg1 string
	Leg2 string
	Leg3 string
}

func TriangleKey(a, b, c string) string {
	x := []string{b, c}
	sort.Strings(x)
	return a + "|" + x[0] + "|" + x[1]
}
