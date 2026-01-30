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

type MinuteData struct {
	Min   float64
	Max   float64
	Sum   float64
	Count int
}

func (m *MinuteData) Add(v float64) {
	if m.Count == 0 {
		m.Min = v
		m.Max = v
		m.Sum = v
		m.Count = 1
		return
	}
	if v < m.Min {
		m.Min = v
	}
	if v > m.Max {
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

/* =======================
   RING BUFFER
======================= */

type RingBuffer struct {
	data []*MinuteData
	size int
	pos  int
	full bool
}

func NewRingBuffer(size int) *RingBuffer {
	r := &RingBuffer{
		data: make([]*MinuteData, size),
		size: size,
	}
	for i := range r.data {
		r.data[i] = &MinuteData{}
	}
	return r
}

func (r *RingBuffer) Add(v float64) {
	m := r.data[r.pos]
	m.Add(v)
}

func (r *RingBuffer) NextMinute() {
	r.pos = (r.pos + 1) % r.size
	r.data[r.pos] = &MinuteData{}
	if r.pos == 0 {
		r.full = true
	}
}

func (r *RingBuffer) Values() []float64 {
	var vals []float64
	for _, m := range r.data {
		if m.Count > 0 {
			vals = append(vals, m.Avg())
		}
	}
	return vals
}

func (r *RingBuffer) IsFull() bool {
	return r.full || r.pos == r.size-1
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
   SIGNAL & PnL
======================= */

type Signal struct {
	Direction string
	Coef      float64
	Dev       float64
	Corr      float64
	EntryBTC  float64
	EntryETH  float64
	EntryTime time.Time
	Closed    bool
}

func CheckSignal(btc, eth *RingBuffer, spreadPct, minCorr float64, f *os.File, active **Signal, md MarketData) {
	if *active != nil && !(*active).Closed {
		// уже есть активный сигнал, ничего не делаем
		return
	}

	b := btc.Values()
	e := eth.Values()

	corr := Correlation(b, e)
	if corr < minCorr {
		return
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
		return
	}

	dir := "BUY BTC / SELL ETH"
	if dev > 0 {
		dir = "SELL BTC / BUY ETH"
	}

	sig := &Signal{
		Direction: dir,
		Coef:      curCoef,
		Dev:       dev,
		Corr:      corr,
		EntryBTC:  md.Bid,
		EntryETH:  md.Bid, // для упрощения, берем Bid для обоих
		EntryTime: time.Now(),
	}

	*active = sig

	fmt.Fprintf(f, "[SIGNAL] %s | coef=%.5f | dev=%.2f%% | corr=%.3f | entryBTC=%.2f | entryETH=%.2f\n",
		dir, curCoef, dev, corr, sig.EntryBTC, sig.EntryETH)
	f.Sync()
}

/* =======================
   KuCoin WS
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

	url := fmt.Sprintf("%s?token=%s&connectId=%d",
		r.Data.IS[0].Endpoint, r.Data.Token, time.Now().UnixNano())

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
	btcBuf := NewRingBuffer(120)
	ethBuf := NewRingBuffer(120)

	f, _ := os.OpenFile("signals.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	defer f.Close()

	out := make(chan MarketData, 1000)

	kc := NewKuCoin(out)
	if err := kc.Start(); err != nil {
		panic(err)
	}

	logTicker := time.NewTicker(5 * time.Minute)
	minTicker := time.NewTicker(1 * time.Minute)

	const (
		spreadPct   = 0.6
		minCorr     = 0.85
		stopLossPct = 0.3
	)

	var activeSignal *Signal
	var curMinuteBTC, curMinuteETH = &MinuteData{}, &MinuteData{}

	for {
		select {
		case md := <-out:
			switch md.Symbol {
			case "BTC-USDT":
				curMinuteBTC.Add(md.Bid)
			case "ETH-USDT":
				curMinuteETH.Add(md.Bid)
			}

			if btcBuf.IsFull() && ethBuf.IsFull() {
				CheckSignal(btcBuf, ethBuf, spreadPct, minCorr, f, &activeSignal, md)
			}

			// после сигнала проверяем стоп-лосс
			if activeSignal != nil && !activeSignal.Closed {
				curCoef := md.Bid / md.Bid // упрощенно для примера
				dev := (curCoef - activeSignal.Coef) / activeSignal.Coef * 100
				if math.Abs(dev) >= stopLossPct {
					activeSignal.Closed = true
					fmt.Fprintf(f, "[STOPLOSS] Closed | profit/loss: %.2f%%\n", dev)
					f.Sync()
				}
			}

		case <-minTicker.C:
			// сохраняем минуту в ринг
			btcBuf.NextMinute()
			btcBuf.data[btcBuf.pos] = curMinuteBTC
			ethBuf.NextMinute()
			ethBuf.data[ethBuf.pos] = curMinuteETH
			curMinuteBTC, curMinuteETH = &MinuteData{}, &MinuteData{}

		case <-logTicker.C:
			if btcBuf.IsFull() && ethBuf.IsFull() {
				b := btcBuf.Values()
				e := ethBuf.Values()
				corr := Correlation(b, e)
				coef := b[len(b)-1] / e[len(e)-1]
				fmt.Printf("[METRICS] corr=%.3f | coef=%.5f\n", corr, coef)
			}
		}
	}
}



