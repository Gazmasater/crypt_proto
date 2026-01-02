package common

func FindLeg(a, b string, markets map[string]Market) (Market, bool) {
	if m, ok := markets[a+"_"+b]; ok {
		return m, true
	}
	if m, ok := markets[b+"_"+a]; ok {
		return m, true
	}
	return Market{}, false
}

func ResolveSide(from, to string, m Market) string {
	if m.Base == to && m.Quote == from {
		return "BUY"
	}
	if m.Base == from && m.Quote == to {
		return "SELL"
	}
	return ""
}
