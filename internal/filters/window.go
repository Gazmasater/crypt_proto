package filters

import "time"

type MidPoint struct {
	Price     float64
	Timestamp int64
}

type MidWindow struct {
	data     []MidPoint
	duration int64 // ms
}

func NewMidWindow(duration time.Duration) *MidWindow {
	return &MidWindow{
		duration: duration.Milliseconds(),
	}
}

func (w *MidWindow) Push(price float64, ts int64) {
	w.data = append(w.data, MidPoint{price, ts})

	cutoff := ts - w.duration
	i := 0
	for ; i < len(w.data); i++ {
		if w.data[i].Timestamp >= cutoff {
			break
		}
	}
	w.data = w.data[i:]
}

func (w *MidWindow) Values() []float64 {
	out := make([]float64, 0, len(w.data))
	for _, p := range w.data {
		out = append(out, p.Price)
	}
	return out
}
