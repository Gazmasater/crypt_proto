package filters

import "math"

type VolatilityStats struct {
	StdDev float64
	MaxDev float64
}

func CalcVolatility(values []float64) VolatilityStats {
	if len(values) < 3 {
		return VolatilityStats{}
	}

	var sum float64
	for _, v := range values {
		sum += v
	}
	mean := sum / float64(len(values))

	var variance float64
	var maxDev float64

	for _, v := range values {
		d := math.Abs(v - mean)
		variance += d * d
		if d > maxDev {
			maxDev = d
		}
	}

	return VolatilityStats{
		StdDev: math.Sqrt(variance / float64(len(values))),
		MaxDev: maxDev,
	}
}
