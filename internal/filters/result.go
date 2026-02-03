package filters

type FilterResult struct {
	OK     bool
	Reason string
	Data   map[string]float64
}

func Pass(data map[string]float64) FilterResult {
	return FilterResult{
		OK:   true,
		Data: data,
	}
}

func Reject(reason string, data map[string]float64) FilterResult {
	return FilterResult{
		OK:     false,
		Reason: reason,
		Data:   data,
	}
}
