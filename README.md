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

/* =======================
   RING BUFFER
======================= */

type RingBuffer struct {
	data []float64
	size int
}

func NewRingBuffer(size int) *RingBuffer {
	return &RingBuffer{
		data: make([]float64, 0, size),
		size: size,
	}
}

func (r *RingBuffer) Add(v float64) {
	if len(r.data) < r.size {
		r.data = append(r.data, v)
		return
	}
	copy(r.data, r.data[1:])
	r.data[len(r.data)-1] = v
}

func (r *RingBuffer) Values() []float64 { return r.data }
func (r *RingBuffer) Len() int          { return len(r.data) }

/* =======================
   MATH
======================= */

func Correlation(x, y []float64) float64 {
	if len(x) != len(y) {
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
	den := math.Sqrt((n*sx2-sx*sx)*(n*sy2-sy*sy))
	if den == 0 {
		return 0
	}
	return num / den
}

/* =======================
   SIGNAL LOGIC
======================= */

func CheckSignal(btc, eth *RingBuffer, spreadPct, minCorr float64, f *os.File) {
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

	fmt.Fprintf(
		f,
		"[SIGNAL] %s | coef=%.5f | dev=%.2f%% | corr=%.3f\n",
		dir, curCoef, dev, corr,
	)
}

func LogMetrics(btc, eth *RingBuffer, f *os.File) {
	b := btc.Values()
	e := eth.Values()

	corr := Correlation(b, e)
	coef := b[len(b)-1] / e[len(e)-1]

	fmt.Fprintf(
		f,
		"[METRICS] corr=%.3f | coef=%.5f\n",
		corr, coef,
	)
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
	btc := NewRingBuffer(120)
	eth := NewRingBuffer(120)

	f, _ := os.OpenFile("signals.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	defer f.Close()

	out := make(chan MarketData, 1000)

	kc := NewKuCoin(out)
	if err := kc.Start(); err != nil {
		panic(err)
	}

	logTicker := time.NewTicker(5 * time.Minute)

	const (
		spreadPct = 0.6
		minCorr   = 0.85
	)

	for {
		select {
		case md := <-out:
			if md.Symbol == "BTC-USDT" {
				btc.Add(md.Bid)
			}
			if md.Symbol == "ETH-USDT" {
				eth.Add(md.Bid)
			}

			if btc.Len() == 120 && eth.Len() == 120 {
				CheckSignal(btc, eth, spreadPct, minCorr, f)
			}

		case <-logTicker.C:
			if btc.Len() == 120 && eth.Len() == 120 {
				LogMetrics(btc, eth, f)
			}
		}
	}
}

