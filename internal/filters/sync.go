package filters

import "math"

func LegsLag(timestamps ...int64) int64 {
	if len(timestamps) == 0 {
		return 0
	}

	minTs := timestamps[0]
	maxTs := timestamps[0]

	for _, ts := range timestamps {
		minTs = int64(math.Min(float64(minTs), float64(ts)))
		maxTs = int64(math.Max(float64(maxTs), float64(ts)))
	}
	return maxTs - minTs
}
