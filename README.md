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



package collector

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"crypt_proto/pkg/models"

	"github.com/gorilla/websocket"
	"github.com/tidwall/gjson"
)

type Last struct {
	Bid float64
	Ask float64
}

type KuCoinCollector struct {
	ctx    context.Context
	cancel context.CancelFunc
	ws     *kucoinWS
	out    chan<- *models.MarketData
}

type kucoinWS struct {
	conn    *websocket.Conn
	symbols []string
	last    map[string]Last
}

func NewKuCoinCollector(symbols []string) (*KuCoinCollector, error) {
	ctx, cancel := context.WithCancel(context.Background())
	ws := &kucoinWS{
		symbols: symbols,
		last:    make(map[string]Last),
	}
	return &KuCoinCollector{
		ctx:    ctx,
		cancel: cancel,
		ws:     ws,
	}, nil
}

func (c *KuCoinCollector) Start(out chan<- *models.MarketData) error {
	c.out = out
	if err := c.ws.connect(); err != nil {
		return err
	}
	go c.ws.readLoop(c)
	go c.ws.subscribeLoop()
	go c.ws.pingLoop()
	log.Println("[KuCoin] WS started")
	return nil
}

func (c *KuCoinCollector) Stop() error {
	c.cancel()
	if c.ws.conn != nil {
		return c.ws.conn.Close()
	}
	return nil
}

// ======== WS методы ========
func (ws *kucoinWS) connect() error {
	req, _ := http.NewRequest("POST", "https://api.kucoin.com/api/v1/bullet-public", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var r struct {
		Data struct {
			Token           string `json:"token"`
			InstanceServers []struct {
				Endpoint string `json:"endpoint"`
			} `json:"instanceServers"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return err
	}

	url := fmt.Sprintf(
		"%s?token=%s&connectId=%d",
		r.Data.InstanceServers[0].Endpoint,
		r.Data.Token,
		time.Now().UnixNano(),
	)

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return err
	}
	ws.conn = conn
	log.Println("[KuCoin WS] connected")
	return nil
}

func (ws *kucoinWS) subscribeLoop() {
	for _, s := range ws.symbols {
		_ = ws.conn.WriteJSON(map[string]any{
			"id":       time.Now().UnixNano(),
			"type":     "subscribe",
			"topic":    "/market/ticker:" + s,
			"response": true,
		})
	}
}

func (ws *kucoinWS) pingLoop() {
	t := time.NewTicker(20 * time.Second)
	defer t.Stop()
	for range t.C {
		_ = ws.conn.WriteJSON(map[string]any{"id": time.Now().UnixNano(), "type": "ping"})
	}
}

func (ws *kucoinWS) readLoop(c *KuCoinCollector) {
	for {
		_, msg, err := ws.conn.ReadMessage()
		if err != nil {
			log.Println("[KuCoin WS] read error:", err)
			return
		}
		ws.handle(c, msg)
	}
}

func (ws *kucoinWS) handle(c *KuCoinCollector, msg []byte) {
	topic := gjson.GetBytes(msg, "topic").String()
	if topic == "" {
		return
	}

	symbol := topic[len("/market/ticker:"):]
	data := gjson.GetBytes(msg, "data")
	bid := data.Get("bestBid").Float()
	ask := data.Get("bestAsk").Float()
	if bid == 0 || ask == 0 {
		return
	}

	last, ok := ws.last[symbol]
	if ok && last.Bid == bid && last.Ask == ask {
		return
	}

	ws.last[symbol] = Last{Bid: bid, Ask: ask}

	c.out <- &models.MarketData{
		Exchange: "KuCoin",
		Symbol:   symbol,
		Bid:      bid,
		Ask:      ask,
		BidSize:  data.Get("bestBidSize").Float(),
		AskSize:  data.Get("bestAskSize").Float(),
	}
}



package main

import (
	"log"
	"os"
	"time"

	"crypt_proto/pkg/collector"
	"crypt_proto/pkg/models"
)

func main() {
	btcBuf := NewRingBuffer(120)
	ethBuf := NewRingBuffer(120)

	f, err := os.OpenFile("signals.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	out := make(chan *models.MarketData, 1000)

	coll, err := collector.NewKuCoinCollector([]string{"BTC-USDT", "ETH-USDT"})
	if err != nil {
		panic(err)
	}

	if err := coll.Start(out); err != nil {
		panic(err)
	}
	defer coll.Stop()

	spreadThresholdPct := 0.6
	minCorr := 0.85
	logTicker := time.NewTicker(5 * time.Minute)

	for {
		select {
		case md := <-out:
			switch md.Symbol {
			case "BTC-USDT":
				btcBuf.Add(md.Bid)
			case "ETH-USDT":
				ethBuf.Add(md.Bid)
			}
			if btcBuf.Len() >= 120 && ethBuf.Len() >= 120 {
				CheckSignal(btcBuf, ethBuf, spreadThresholdPct, minCorr, f)
			}
		case <-logTicker.C:
			if btcBuf.Len() >= 120 && ethBuf.Len() >= 120 {
				LogMetrics(btcBuf, ethBuf, f)
			}
		}
	}
}

