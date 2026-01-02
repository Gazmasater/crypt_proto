package common

func FindLeg(from, to string, markets map[string]Market) (Market, bool) {
	if m, ok := markets[from+"_"+to]; ok {
		return m, true
	}
	if m, ok := markets[to+"_"+from]; ok {
		return m, true
	}
	return Market{}, false
}

func ResolveSide(from, to string, m Market) string {
	if m.Quote == from && m.Base == to {
		return "BUY"
	}
	if m.Base == from && m.Quote == to {
		return "SELL"
	}
	return ""
}
