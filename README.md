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
	"io"
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

type MinuteData struct {
	Min, Max, Sum float64
	Count         int
}

type RingBuffer struct {
	data []MinuteData
	size int
}

func NewRingBuffer(size int) *RingBuffer {
	return &RingBuffer{
		data: make([]MinuteData, 0, size),
		size: size,
	}
}

// Добавляем новую минуту
func (r *RingBuffer) AddMinute(m MinuteData) {
	if len(r.data) < r.size {
		r.data = append(r.data, m)
		return
	}
	copy(r.data, r.data[1:])
	r.data[len(r.data)-1] = m
}

// Последние N минут, усреднённые
func (r *RingBuffer) AvgValues() []float64 {
	res := make([]float64, len(r.data))
	for i, m := range r.data {
		if m.Count > 0 {
			res[i] = m.Sum / float64(m.Count)
		}
	}
	return res
}

func (r *RingBuffer) Len() int { return len(r.data) }

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
	den := math.Sqrt((n*sx2 - sx*sx) * (n*sy2 - sy*sy))
	if den == 0 {
		return 0
	}
	return num / den
}

/* =======================
   SIGNAL LOGIC
======================= */

type Signal struct {
	Direction   string
	EntryCoef   float64
	EntryTime   time.Time
	StopLossPct float64
	TakeProfit  float64
	Active      bool
}

func CheckSignal(btc, eth *RingBuffer, spreadPct, minCorr float64, f *os.File, sig *Signal) {
	b := btc.AvgValues()
	e := eth.AvgValues()
	if len(b) == 0 || len(e) == 0 {
		return
	}

	corr := Correlation(b, e)
	if corr < minCorr || sig.Active {
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

	sig.Direction = dir
	sig.EntryCoef = curCoef
	sig.EntryTime = time.Now()
	sig.StopLossPct = spreadPct // фиксируем stop-loss = spread
	sig.TakeProfit = 0
	sig.Active = true

	fmt.Fprintf(f, "[SIGNAL] %s | coef=%.5f | dev=%.2f%% | corr=%.3f\n",
		dir, curCoef, dev, corr)
	f.Sync()
}

/* =======================
   PAPER PnL
======================= */

func UpdatePnL(btc, eth *RingBuffer, sig *Signal, f *os.File) {
	if !sig.Active || sig.EntryCoef == 0 {
		return
	}

	curCoef := btc.AvgValues()[len(btc.data)-1] / eth.AvgValues()[len(eth.data)-1]
	var pnl float64
	if sig.Direction == "BUY BTC / SELL ETH" {
		pnl = (curCoef - sig.EntryCoef) / sig.EntryCoef * 100
	} else {
		pnl = (sig.EntryCoef - curCoef) / sig.EntryCoef * 100
	}

	fmt.Fprintf(f, "[PnL] %.3f%% | curCoef=%.5f | entry=%.5f\n", pnl, curCoef, sig.EntryCoef)
	f.Sync()

	// Stop-loss или take-profit
	if math.Abs(pnl) >= sig.StopLossPct {
		fmt.Fprintf(f, "[CLOSE] Signal closed due to stop-loss | pnl=%.3f%%\n", pnl)
		sig.Active = false
	}
}

/* =======================
   KUCOIN WS
======================= */

type KuCoin struct {
	ctx    context.Context
	cancel context.CancelFunc
	conn   *websocket.Conn
	last   map[string]float64
	out    chan<- MinuteData
}

func NewKuCoin(out chan<- MinuteData) *KuCoin {
	ctx, cancel := context.WithCancel(context.Background())
	return &KuCoin{
		ctx:    ctx,
		cancel: cancel,
		last:   make(map[string]float64),
		out:    out,
	}
}

func (k *KuCoin) Start(symbols []string) error {
	req, _ := http.NewRequest("POST", "https://api.kucoin.com/api/v1/bullet-public", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var r struct {
		Data struct {
			Token           string `json:"token"`
			InstanceServers []struct{ Endpoint string } `json:"instanceServers"`
		} `json:"data"`
	}
	json.NewDecoder(resp.Body).Decode(&r)

	url := fmt.Sprintf("%s?token=%s&connectId=%d", r.Data.InstanceServers[0].Endpoint, r.Data.Token, time.Now().UnixNano())
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return err
	}
	k.conn = conn

	for _, s := range symbols {
		conn.WriteJSON(map[string]any{
			"id":       time.Now().UnixNano(),
			"type":     "subscribe",
			"topic":    "/market/candles:1min:" + s,
			"response": true,
		})
	}

	go k.read(symbols)
	go k.ping()
	return nil
}

func (k *KuCoin) ping() {
	t := time.NewTicker(20 * time.Second)
	for range t.C {
		k.conn.WriteJSON(map[string]any{"type": "ping"})
	}
}

func (k *KuCoin) read(symbols []string) {
	for {
		_, msg, err := k.conn.ReadMessage()
		if err != nil {
			log.Println("ws error:", err)
			return
		}

		for _, sym := range symbols {
			data := gjson.GetBytes(msg, "data")
			if !data.Exists() {
				continue
			}

			// Берём close каждой свечи
			closePrice := data.Get("close").Float()
			if closePrice == 0 {
				continue
			}
			k.out <- MinuteData{Sum: closePrice, Count: 1}
		}
	}
}

/* =======================
   LOAD HISTORY
======================= */

func loadHistory(symbol string, buf *RingBuffer) error {
	end := time.Now()
	start := end.Add(-120 * time.Minute)

	url := fmt.Sprintf("https://api.kucoin.com/api/v1/market/candles?type=1min&symbol=%s&startAt=%d&endAt=%d",
		symbol, start.Unix(), end.Unix())

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	arr := gjson.ParseBytes(body).Array()
	for i := len(arr) - 1; i >= 0; i-- {
		c := arr[i].Array()
		closePrice := c[2].Float()
		high := c[3].Float()
		low := c[4].Float()
		buf.AddMinute(MinuteData{Min: low, Max: high, Sum: closePrice, Count: 1})
	}
	return nil
}

/* =======================
   MAIN
======================= */

func main() {
	btcBuf := NewRingBuffer(120)
	ethBuf := NewRingBuffer(120)

	fmt.Println("loading history...")
	if err := loadHistory("BTC-USDT", btcBuf); err != nil {
		log.Println("BTC history error:", err)
	}
	if err := loadHistory("ETH-USDT", ethBuf); err != nil {
		log.Println("ETH history error:", err)
	}
	fmt.Println("history loaded")

	f, _ := os.OpenFile("signals.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	defer f.Close()

	out := make(chan MinuteData, 1000)
	sig := &Signal{}

	kc := NewKuCoin(out)
	if err := kc.Start([]string{"BTC-USDT", "ETH-USDT"}); err != nil {
		panic(err)
	}

	logTicker := time.NewTicker(5 * time.Minute)
	spreadPct := 0.6
	minCorr := 0.85

	for {
		select {
		case md := <-out:
			btcBuf.AddMinute(md)
			ethBuf.AddMinute(md)

			if btcBuf.Len() == 120 && ethBuf.Len() == 120 {
				CheckSignal(btcBuf, ethBuf, spreadPct, minCorr, f, sig)
			}
			UpdatePnL(btcBuf, ethBuf, sig, f)

		case <-logTicker.C:
			fmt.Printf("[METRICS] BTC avg=%.5f ETH avg=%.5f\n",
				btcBuf.AvgValues()[len(btcBuf.data)-1],
				ethBuf.AvgValues()[len(ethBuf.data)-1])
		}
	}
}



