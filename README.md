Если оставить только нужное:

p99 execution latency
Micro-volatility (100 мс)
Fill ratio
Capture rate
Inventory drift




Название API
9623527002

696935c42a6dcd00013273f2
b348b686-55ff-4290-897b-02d55f815f65




apikey = "4333ed4b-cd83-49f5-97d1-c399e2349748"
secretkey = "E3848531135EDB4CCFDA0F1BC14CD274"
IP = ""
Название API-ключа = "Arb"
Доступы = "Чтение"



sudo systemctl mask sleep.target suspend.target hibernate.target hybrid-sleep.target



wbs-api.mexc.com/ws 


[https://edis-global.vercel.app/ru/vps-hosting/singapore-singapore
](https://sg.edisglobal.com/)



git pull --rebase origin privat
git push origin privat


import (
    // ...
    "net/http"
    _ "net/http/pprof"
)


   // pprof HTTP-сервер
    go func() {
        log.Println("pprof on http://localhost:6060/debug/pprof/")
        if err := http.ListenAndServe("localhost:6060", nil); err != nil {
            log.Printf("pprof server error: %v", err)
        }
    }()


	go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30


(pprof) top        # показать топ функций по CPU
(pprof) top10
(pprof) list parsePBWrapperMid   # подробный разбор одной функции
(pprof) quit


go tool pprof http://localhost:6060/debug/pprof/heap


(pprof) top
(pprof) top -cum
(pprof) list parsePBWrapperMid
(pprof) quit




go run -race main.go


GOMAXPROCS=8 go run -race main.go


window.go

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


volatility.go

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


spread.go

package filters

func SpreadPct(bid, ask float64) float64 {
	if bid <= 0 || ask <= 0 {
		return 1
	}
	mid := (bid + ask) / 2
	return (ask - bid) / mid
}


sync.go

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


result.go

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





