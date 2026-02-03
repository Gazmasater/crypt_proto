package filters

func SpreadPct(bid, ask float64) float64 {
	if bid <= 0 || ask <= 0 {
		return 1
	}
	mid := (bid + ask) / 2
	return (ask - bid) / mid
}
