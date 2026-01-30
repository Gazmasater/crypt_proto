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
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/tidwall/gjson"
)

/* =======================
   DATA STRUCTURES
======================= */

type MarketData struct {
	Symbol string
	Bid    float64
	Ask    float64
}

type Last struct {
	Bid float64
	Ask float64
}

type Signal struct {
	Dir      string
	EntryBTC float64
	EntryETH float64
	Closed   bool
}

/* =======================
   RING BUFFER PER MINUTE
======================= */

type MinuteData struct {
	Min   float64
	Max   float64
	Sum   float64
	Count int
}

func (m *MinuteData) Add(v float64) {
	if m.Count == 0 || v < m.Min {
		m.Min = v
	}
	if m.Count == 0 || v > m.Max {
		m.Max = v
	}
	m.Sum += v
	m.Count++
}

func (m *MinuteData) Avg() float64 {
	if m.Count == 0 {
		return 0
	}
	return m.Sum / float64(m.Count)
}

type RingBuffer struct {
	data []MinuteData
	size int
	idx  int
	full bool
}

func NewRingBuffer(size int) *RingBuffer {
	return &RingBuffer{
		data: make([]MinuteData, size),
		size: size,
	}
}

func (r *RingBuffer) AddMinute(min MinuteData) {
	r.data[r.idx] = min
	r.idx = (r.idx + 1) % r.size
	if r.idx == 0 {
		r.full = true
	}
}

func (r *RingBuffer) Values() []float64 {
	n := r.size
	if !r.full {
		n = r.idx
	}
	vals := make([]float64, 0, n)
	for i := 0; i < n; i++ {
		vals = append(vals, r.data[i].Avg())
	}
	return vals
}

func (r *RingBuffer) Len() int {
	if r.full {
		return r.size
	}
	return r.idx
}

/* =======================
   MATH
======================= */

func Correlation(x, y []float64) float64 {
	if len(x) != len(y) || len(x) == 0 {
		return 0
	}
	n := float64(len(x))
	var sx, sy, sxy, sx2, sy2 float64
	for i := range x {
		sx += x[i]
		sy += y[i]
		sxy += x[i] * y[i]
		sx2 += x[i] * x[i]
		sy2 += y[i] * y[i]
	}
	num := n*sxy - sx*sy
	den := math.Sqrt((n*sx2 - sx*sx) * (n*sy2 - sy*sy))
	if den == 0 {
		return 0
	}
	return num / den
}

/* =======================
   SIGNAL LOGIC
======================= */

func CheckSignal(btc, eth *RingBuffer, spreadPct, minCorr float64, f *os.File) *Signal {
	b := btc.Values()
	e := eth.Values()

	corr := Correlation(b, e)
	if corr < minCorr {
		return nil
	}

	curCoef := b[len(b)-1] / e[len(e)-1]

	minC, maxC := curCoef, curCoef
	for i := range b {
		c := b[i] / e[i]
		if c < minC {
			minC = c
		}
		if c > maxC {
			maxC = c
		}
	}

	mid := (minC + maxC) / 2
	dev := (curCoef - mid) / mid * 100

	if math.Abs(dev) < spreadPct {
		return nil
	}

	dir := "BUY BTC / SELL ETH"
	if dev > 0 {
		dir = "SELL BTC / BUY ETH"
	}

	fmt.Fprintf(f, "[SIGNAL] %s | coef=%.5f | dev=%.2f%% | corr=%.3f\n",
		dir, curCoef, dev, corr)

	return &Signal{
		Dir:      dir,
		EntryBTC: b[len(b)-1],
		EntryETH: e[len(e)-1],
	}
}

func LogPnL(sig *Signal, btc, eth *RingBuffer, stopLossPct, takeProfitPct float64, f *os.File) {
	if sig == nil || sig.Closed {
		return
	}
	curBTC := btc.Values()[len(btc.Values())-1]
	curETH := eth.Values()[len(eth.Values())-1]

	pnl := 0.0
	closeSignal := false
	closeReason := ""

	if sig.Dir == "BUY BTC / SELL ETH" {
		pnl = (curBTC - sig.EntryBTC) - (curETH - sig.EntryETH)
		if (curBTC-sig.EntryBTC)/sig.EntryBTC*100 >= takeProfitPct {
			closeSignal = true
			closeReason = "TAKE-PROFIT"
		} else if (curBTC-sig.EntryBTC)/sig.EntryBTC*100 <= -stopLossPct {
			closeSignal = true
			closeReason = "STOP-LOSS"
		}
	} else {
		pnl = (sig.EntryBTC - curBTC) - (sig.EntryETH - curETH)
		if (sig.EntryBTC-curBTC)/sig.EntryBTC*100 >= takeProfitPct {
			closeSignal = true
			closeReason = "TAKE-PROFIT"
		} else if (sig.EntryBTC-curBTC)/sig.EntryBTC*100 <= -stopLossPct {
			closeSignal = true
			closeReason = "STOP-LOSS"
		}
	}

	fmt.Fprintf(f, "[PnL] %.5f | signal=%s\n", pnl, sig.Dir)
	if closeSignal {
		fmt.Fprintf(f, "[%s] closing signal %s | pnl=%.5f\n", closeReason, sig.Dir, pnl)
		sig.Closed = true
	}
}

/* =======================
   KUCOIN WS
======================= */

type KuCoin struct {
	ctx    context.Context
	cancel context.CancelFunc
	conn   *websocket.Conn
	last   map[string]Last
	out    chan<- MarketData
}

func NewKuCoin(out chan<- MarketData) *KuCoin {
	ctx, cancel := context.WithCancel(context.Background())
	return &KuCoin{
		ctx:    ctx,
		cancel: cancel,
		last:   make(map[string]Last),
		out:    out,
	}
}

func (k *KuCoin) Start() error {
	req, _ := http.NewRequest("POST", "https://api.kucoin.com/api/v1/bullet-public", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var r struct {
		Data struct {
			Token string `json:"token"`
			IS    []struct {
				Endpoint string `json:"endpoint"`
			} `json:"instanceServers"`
		} `json:"data"`
	}
	json.NewDecoder(resp.Body).Decode(&r)

	url := fmt.Sprintf("%s?token=%s&connectId=%d", r.Data.IS[0].Endpoint, r.Data.Token, time.Now().UnixNano())
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return err
	}
	k.conn = conn

	sub := func(sym string) {
		conn.WriteJSON(map[string]any{
			"id":       time.Now().UnixNano(),
			"type":     "subscribe",
			"topic":    "/market/ticker:" + sym,
			"response": true,
		})
	}
	sub("BTC-USDT")
	sub("ETH-USDT")

	go k.read()
	go k.ping()
	return nil
}

func (k *KuCoin) ping() {
	t := time.NewTicker(20 * time.Second)
	for range t.C {
		k.conn.WriteJSON(map[string]any{"type": "ping"})
	}
}

func (k *KuCoin) read() {
	for {
		_, msg, err := k.conn.ReadMessage()
		if err != nil {
			log.Println("ws error:", err)
			return
		}

		topic := gjson.GetBytes(msg, "topic").String()
		if topic == "" {
			continue
		}

		sym := topic[len("/market/ticker:"):]
		d := gjson.GetBytes(msg, "data")
		bid := d.Get("bestBid").Float()
		ask := d.Get("bestAsk").Float()

		if bid == 0 || ask == 0 {
			continue
		}

		last := k.last[sym]
		if last.Bid == bid && last.Ask == ask {
			continue
		}
		k.last[sym] = Last{bid, ask}

		k.out <- MarketData{Symbol: sym, Bid: bid, Ask: ask}
	}
}

/* =======================
   MAIN
======================= */

func main() {
	btcBuf := NewRingBuffer(120) // 120 минут
	ethBuf := NewRingBuffer(120)

	f, _ := os.OpenFile("signals.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	defer f.Close()

	out := make(chan MarketData, 1000)
	kc := NewKuCoin(out)
	if err := kc.Start(); err != nil {
		panic(err)
	}

	var activeSignal *Signal
	var curMinuteBTC, curMinuteETH MinuteData
	minTicker := time.NewTicker(5 * time.Second) // обновляем текущую минуту каждые 5 сек
	logTicker := time.NewTicker(5 * time.Minute) // вывод метрик

	const (
		spreadPct     = 0.6
		minCorr       = 0.85
		stopLossPct   = 0.5
		takeProfitPct = 0.5
	)

	for {
		select {
		case md := <-out:
			if md.Symbol == "BTC-USDT" {
				curMinuteBTC.Add(md.Bid)
			} else if md.Symbol == "ETH-USDT" {
				curMinuteETH.Add(md.Bid)
			}

		case <-minTicker.C:
			// добавляем минуту в ринг-буфер
			if curMinuteBTC.Count > 0 && curMinuteETH.Count > 0 {
				btcBuf.AddMinute(curMinuteBTC)
				ethBuf.AddMinute(curMinuteETH)

				if activeSignal == nil && btcBuf.Len() == 120 && ethBuf.Len() == 120 {
					activeSignal = CheckSignal(btcBuf, ethBuf, spreadPct, minCorr, f)
				}

				if activeSignal != nil && !activeSignal.Closed {
					LogPnL(activeSignal, btcBuf, ethBuf, stopLossPct, takeProfitPct, f)
				}

				curMinuteBTC = MinuteData{}
				curMinuteETH = MinuteData{}
			}

		case <-logTicker.C:
			if btcBuf.Len() > 0 && ethBuf.Len() > 0 {
				bVals := btcBuf.Values()
				eVals := ethBuf.Values()
				corr := Correlation(bVals, eVals)
				coef := bVals[len(bVals)-1] / eVals[len(eVals)-1]
				fmt.Printf("[METRICS] corr=%.3f | coef=%.5f\n", corr, coef)
			}
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
			"path": "/docs/checks/",
			"scheme": "https",
			"authority": "staticcheck.dev",
			"fragment": "QF1003"
		}
	},
	"severity": 2,
	"message": "could use tagged switch on md.Symbol",
	"source": "QF1003",
	"startLineNumber": 354,
	"startColumn": 4,
	"endLineNumber": 354,
	"endColumn": 30,
	"modelVersionId": 3,
	"origin": "extHost1"
}]

