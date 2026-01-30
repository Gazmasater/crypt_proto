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



package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"time"
)

/* ================= CONFIG ================= */

const (
	SymbolBTC = "BTC-USDT"
	SymbolETH = "ETH-USDT"

	WindowMinutes = 120

	EntryDeviationPct = 0.6
	StopLossPct       = 1.2
	MinCorrelation    = 0.85

	RestInterval = 5 * time.Second
	LogInterval  = 5 * time.Minute
)

/* ================= DATA ================= */

type Candle struct {
	Open  float64
	High  float64
	Low   float64
	Close float64
}

type MinuteStat struct {
	Mean float64
}

type RingBuffer struct {
	data []MinuteStat
	size int
}

func NewRing(size int) *RingBuffer {
	return &RingBuffer{size: size}
}

func (r *RingBuffer) Push(v MinuteStat) {
	if len(r.data) < r.size {
		r.data = append(r.data, v)
		return
	}
	copy(r.data, r.data[1:])
	r.data[len(r.data)-1] = v
}

func (r *RingBuffer) Values() []MinuteStat { return r.data }
func (r *RingBuffer) Full() bool           { return len(r.data) == r.size }

/* ================= SIGNAL ================= */

type Signal struct {
	Active     bool
	Direction  string
	EntryCoef  float64
	EntryTime  time.Time
}

/* ================= MATH ================= */

func Correlation(a, b []float64) float64 {
	n := float64(len(a))
	var sa, sb, sab, sa2, sb2 float64

	for i := range a {
		sa += a[i]
		sb += b[i]
		sab += a[i] * b[i]
		sa2 += a[i] * a[i]
		sb2 += b[i] * b[i]
	}

	num := n*sab - sa*sb
	den := math.Sqrt((n*sa2-sa*sa)*(n*sb2-sb*sb))
	if den == 0 {
		return 0
	}
	return num / den
}

/* ================= KUCOIN REST ================= */

func fetchLastPrice(symbol string) float64 {
	url := fmt.Sprintf("https://api.kucoin.com/api/v1/market/orderbook/level1?symbol=%s", symbol)
	resp, err := http.Get(url)
	if err != nil {
		return 0
	}
	defer resp.Body.Close()

	var r struct {
		Data struct {
			Price string `json:"price"`
		} `json:"data"`
	}
	json.NewDecoder(resp.Body).Decode(&r)

	var p float64
	fmt.Sscanf(r.Data.Price, "%f", &p)
	return p
}

func fetchHistory(symbol string) []MinuteStat {
	url := fmt.Sprintf(
		"https://api.kucoin.com/api/v1/market/candles?symbol=%s&type=1min",
		symbol,
	)

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var raw struct {
		Data [][]string `json:"data"`
	}
	json.NewDecoder(resp.Body).Decode(&raw)

	out := make([]MinuteStat, 0, WindowMinutes)
	for i := len(raw.Data) - WindowMinutes; i < len(raw.Data); i++ {
		c := raw.Data[i]
		var o, h, l, cl float64
		fmt.Sscanf(c[1], "%f", &o)
		fmt.Sscanf(c[2], "%f", &h)
		fmt.Sscanf(c[3], "%f", &l)
		fmt.Sscanf(c[4], "%f", &cl)

		mean := (o + h + l + cl) / 4
		out = append(out, MinuteStat{Mean: mean})
	}
	return out
}

/* ================= MAIN ================= */

func main() {
	log.SetFlags(log.Ltime)

	log.Println("loading history...")

	btcRing := NewRing(WindowMinutes)
	ethRing := NewRing(WindowMinutes)

	for _, v := range fetchHistory(SymbolBTC) {
		btcRing.Push(v)
	}
	for _, v := range fetchHistory(SymbolETH) {
		ethRing.Push(v)
	}

	log.Println("history loaded")

	file, _ := os.OpenFile("signals.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	defer file.Close()

	var (
		curSumBTC float64
		curSumETH float64
		curCount  int

		activeSignal Signal
	)

	restTicker := time.NewTicker(RestInterval)
	logTicker := time.NewTicker(LogInterval)
	minuteTicker := time.NewTicker(time.Minute)

	for {
		select {

		case <-restTicker.C:
			curSumBTC += fetchLastPrice(SymbolBTC)
			curSumETH += fetchLastPrice(SymbolETH)
			curCount++

		case <-minuteTicker.C:
			if curCount == 0 {
				continue
			}

			btcRing.Push(MinuteStat{Mean: curSumBTC / float64(curCount)})
			ethRing.Push(MinuteStat{Mean: curSumETH / float64(curCount)})

			curSumBTC, curSumETH, curCount = 0, 0, 0

			if !btcRing.Full() {
				continue
			}

			btcVals := btcRing.Values()
			ethVals := ethRing.Values()

			btcArr := make([]float64, WindowMinutes)
			ethArr := make([]float64, WindowMinutes)
			coefArr := make([]float64, WindowMinutes)

			for i := 0; i < WindowMinutes; i++ {
				btcArr[i] = btcVals[i].Mean
				ethArr[i] = ethVals[i].Mean
				coefArr[i] = btcArr[i] / ethArr[i]
			}

			corr := Correlation(btcArr, ethArr)
			if corr < MinCorrelation {
				continue
			}

			curCoef := coefArr[len(coefArr)-1]
			minC, maxC := coefArr[0], coefArr[0]
			for _, c := range coefArr {
				if c < minC {
					minC = c
				}
				if c > maxC {
					maxC = c
				}
			}

			mid := (minC + maxC) / 2
			dev := (curCoef - mid) / mid * 100

			if !activeSignal.Active && math.Abs(dev) >= EntryDeviationPct {
				activeSignal = Signal{
					Active:    true,
					EntryCoef: curCoef,
					EntryTime: time.Now(),
					Direction: map[bool]string{
						true:  "SELL BTC / BUY ETH",
						false: "BUY BTC / SELL ETH",
					}[dev > 0],
				}

				fmt.Fprintf(file,
					"[OPEN] %s coef=%.5f dev=%.2f%% corr=%.2f\n",
					activeSignal.Direction, curCoef, dev, corr,
				)
			}

			if activeSignal.Active {
				pnl := (activeSignal.EntryCoef - curCoef) / activeSignal.EntryCoef * 100
				fmt.Fprintf(file, "[PNL] %.3f%%\n", pnl)

				if math.Abs(dev) <= 0.1 {
					fmt.Fprintln(file, "[CLOSE] TAKE PROFIT")
					activeSignal.Active = false
				}

				if math.Abs(dev) >= StopLossPct {
					fmt.Fprintln(file, "[CLOSE] STOP LOSS")
					activeSignal.Active = false
				}
			}

		case <-logTicker.C:
			if !btcRing.Full() {
				continue
			}
			log.Printf("corr=%.2f coef=%.4f dev=%.2f%%",
				Correlation(
					[]float64{btcRing.data[len(btcRing.data)-1].Mean},
					[]float64{ethRing.data[len(ethRing.data)-1].Mean},
				),
				btcRing.data[len(btcRing.data)-1].Mean/
					ethRing.data[len(ethRing.data)-1].Mean,
				0.0,
			)
		}
	}
}


[{
	"resource": "/home/gaz358/myprog/crypt_proto/cmd/stat_arb/stat_arb.go",
	"owner": "_generated_diagnostic_collection_name_#0",
	"code": {
		"value": "default",
		"target": {
			"$mid": 1,
			"path": "/golang.org/x/tools/go/analysis/passes/unusedwrite",
			"scheme": "https",
			"authority": "pkg.go.dev"
		}
	},
	"severity": 2,
	"message": "unused write to field EntryTime",
	"source": "unusedwrite",
	"startLineNumber": 240,
	"startColumn": 15,
	"endLineNumber": 240,
	"endColumn": 15,
	"modelVersionId": 4,
	"origin": "extHost1"
}]


