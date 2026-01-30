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
	Min float64
	Max float64
	Avg float64
}

/* =======================
   RING BUFFER
======================= */

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

func (r *RingBuffer) Add(minute MinuteData) {
	if len(r.data) < r.size {
		r.data = append(r.data, minute)
		return
	}
	copy(r.data, r.data[1:])
	r.data[len(r.data)-1] = minute
}

func (r *RingBuffer) Values() []MinuteData { return r.data }
func (r *RingBuffer) Len() int             { return len(r.data) }

/* =======================
   SIGNAL
======================= */

type Signal struct {
	Active     bool
	Direction  string
	EntryCoef  float64
	EntryTime  time.Time
	EntryAvg   float64
}

func Correlation(x, y []MinuteData) float64 {
	if len(x) != len(y) {
		return 0
	}
	n := float64(len(x))
	var sx, sy, sxy, sx2, sy2 float64
	for i := range x {
		sx += x[i].Avg
		sy += y[i].Avg
		sxy += x[i].Avg * y[i].Avg
		sx2 += x[i].Avg * x[i].Avg
		sy2 += y[i].Avg * y[i].Avg
	}
	num := n*sxy - sx*sy
	den := math.Sqrt((n*sx2 - sx*sx) * (n*sy2 - sy*sy))
	if den == 0 {
		return 0
	}
	return num / den
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
	btcBuf := NewRingBuffer(120)
	ethBuf := NewRingBuffer(120)

	out := make(chan MarketData, 1000)

	kc := NewKuCoin(out)
	if err := kc.Start(); err != nil {
		panic(err)
	}

	f, _ := os.OpenFile("signals.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	defer f.Close()

	logTicker := time.NewTicker(5 * time.Minute)

	var btcTicks, ethTicks []float64
	signal := &Signal{}

	const (
		spreadPct   = 0.6
		stopLossPct = 0.5
		minCorr     = 0.85
	)

	for {
		select {
		case md := <-out:
			// собираем тики для текущей минуты
			if md.Symbol == "BTC-USDT" {
				btcTicks = append(btcTicks, md.Bid)
			} else if md.Symbol == "ETH-USDT" {
				ethTicks = append(ethTicks, md.Bid)
			}

			// каждая минута формируем MinuteData
			now := time.Now()
			if now.Second() == 0 && len(btcTicks) > 0 && len(ethTicks) > 0 {
				minData := func(ticks []float64) MinuteData {
					min, max, sum := ticks[0], ticks[0], 0.0
					for _, v := range ticks {
						if v < min {
							min = v
						}
						if v > max {
							max = v
						}
						sum += v
					}
					return MinuteData{
						Min: min,
						Max: max,
						Avg: sum / float64(len(ticks)),
					}
				}
				btcBuf.Add(minData(btcTicks))
				ethBuf.Add(minData(ethTicks))
				btcTicks = nil
				ethTicks = nil
			}

			// Проверка сигнала
			if btcBuf.Len() == 120 && ethBuf.Len() == 120 && !signal.Active {
				bVals := btcBuf.Values()
				eVals := ethBuf.Values()
				curCoef := bVals[len(bVals)-1].Avg / eVals[len(eVals)-1].Avg
				minC, maxC := curCoef, curCoef
				for i := range bVals {
					c := bVals[i].Avg / eVals[i].Avg
					if c < minC {
						minC = c
					}
					if c > maxC {
						maxC = c
					}
				}
				mid := (minC + maxC) / 2
				dev := (curCoef - mid) / mid * 100
				corr := Correlation(bVals, eVals)

				if corr >= minCorr && math.Abs(dev) >= spreadPct {
					// вход
					dir := "BUY BTC / SELL ETH"
					if dev > 0 {
						dir = "SELL BTC / BUY ETH"
					}
					signal.Active = true
					signal.Direction = dir
					signal.EntryCoef = curCoef
					signal.EntryTime = time.Now()
					signal.EntryAvg = mid
					fmt.Fprintf(f, "[SIGNAL] %s | coef=%.5f | dev=%.2f%% | corr=%.3f\n",
						dir, curCoef, dev, corr)
				}
			}

			// после сигнала PnL и стоп-лосс
			if signal.Active && btcBuf.Len() == 120 && ethBuf.Len() == 120 {
				bVals := btcBuf.Values()
				eVals := ethBuf.Values()
				curCoef := bVals[len(bVals)-1].Avg / eVals[len(eVals)-1].Avg
				pnl := 0.0
				if signal.Direction == "BUY BTC / SELL ETH" {
					pnl = (curCoef - signal.EntryCoef) / signal.EntryCoef * 100
				} else {
					pnl = (signal.EntryCoef - curCoef) / signal.EntryCoef * 100
				}

				// закрытие по возврату к среднему или стоп-лосс
				mid := signal.EntryAvg
				dev := (curCoef - mid) / mid * 100
				if (signal.Direction == "BUY BTC / SELL ETH" && dev <= 0) ||
					(signal.Direction == "SELL BTC / BUY ETH" && dev >= 0) ||
					math.Abs(pnl) >= stopLossPct {
					fmt.Fprintf(f, "[CLOSE] %s | coef=%.5f | PnL=%.3f%% | time=%s\n",
						signal.Direction, curCoef, pnl, time.Now().Format("15:04"))
					signal.Active = false
				}
			}

		case <-logTicker.C:
			// Логируем PnL в файл
			if signal.Active && btcBuf.Len() == 120 && ethBuf.Len() == 120 {
				bVals := btcBuf.Values()
				eVals := ethBuf.Values()
				curCoef := bVals[len(bVals)-1].Avg / eVals[len(eVals)-1].Avg
				pnl := 0.0
				if signal.Direction == "BUY BTC / SELL ETH" {
					pnl = (curCoef - signal.EntryCoef) / signal.EntryCoef * 100
				} else {
					pnl = (signal.EntryCoef - curCoef) / signal.EntryCoef * 100
				}
				fmt.Fprintf(f, "[PnL] %s | coef=%.5f | PnL=%.3f%% | time=%s\n",
					signal.Direction, curCoef, pnl, time.Now().Format("15:04"))
			}

			// Выводим в консоль каждую 5 минут живую статистику
			if btcBuf.Len() == 120 && ethBuf.Len() == 120 {
				bVals := btcBuf.Values()
				eVals := ethBuf.Values()
				curCoef := bVals[len(bVals)-1].Avg / eVals[len(eVals)-1].Avg
				corr := Correlation(bVals, eVals)
				fmt.Printf("[LIVE] coef=%.5f | corr=%.3f | time=%s\n",
					curCoef, corr, time.Now().Format("15:04"))
			}
		}
	}
}


