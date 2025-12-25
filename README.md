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


BOOK_INTERVAL=100ms
SYMBOLS_FILE=triangles_markets.csv
DEBUG=false


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





package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type MarketData struct {
	Exchange string
	Symbol   string
	Bid      float64
	Ask      float64
}

type KuCoinCollector struct {
	ctx     context.Context
	cancel  context.CancelFunc
	symbols []string
	conn    *websocket.Conn
	wsURL   string
}

func NewKuCoinCollector(symbols []string) *KuCoinCollector {
	ctx, cancel := context.WithCancel(context.Background())
	return &KuCoinCollector{
		ctx:     ctx,
		cancel:  cancel,
		symbols: symbols,
		wsURL:   "wss://ws-api-spot.kucoin.com/endpoint", // песочница без REST token
	}
}

func (c *KuCoinCollector) Start(out chan<- MarketData) error {
	// Подключаемся
	conn, _, err := websocket.DefaultDialer.Dial(c.wsURL, nil)
	if err != nil {
		return err
	}
	c.conn = conn
	log.Println("Connected to KuCoin WS")

	// Подписка на пары
	for _, s := range c.symbols {
		topic := "level2/ticker:" + strings.ToUpper(s)
		sub := map[string]interface{}{
			"id":       time.Now().UnixNano(),
			"type":     "subscribe",
			"topic":    topic,
			"response": true,
		}
		if err := conn.WriteJSON(sub); err != nil {
			return err
		}
		log.Println("subscribed:", topic)
	}

	// Ping
	go func() {
		t := time.NewTicker(20 * time.Second)
		defer t.Stop()
		for {
			select {
			case <-c.ctx.Done():
				return
			case <-t.C:
				_ = conn.WriteMessage(websocket.PingMessage, nil)
			}
		}
	}()

	// Read loop
	go func() {
		defer log.Println("readLoop stopped")
		for {
			select {
			case <-c.ctx.Done():
				return
			default:
				_, msg, err := conn.ReadMessage()
				if err != nil {
					if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
						log.Println("WS closed normally")
						return
					}
					log.Println("read error:", err)
					return
				}

				var data map[string]interface{}
				if err := json.Unmarshal(msg, &data); err != nil {
					continue
				}

				if data["type"] == "message" {
					topic, ok := data["topic"].(string)
					if !ok {
						continue
					}
					symbol := strings.Split(topic, ":")[1]

					body, ok := data["data"].(map[string]interface{})
					if !ok {
						continue
					}

					bid := parseFloat(body["bestBid"])
					ask := parseFloat(body["bestAsk"])

					out <- MarketData{
						Exchange: "KuCoin",
						Symbol:   symbol,
						Bid:      bid,
						Ask:      ask,
					}
				}
			}
		}
	}()

	return nil
}

func (c *KuCoinCollector) Stop() {
	c.cancel()
	if c.conn != nil {
		c.conn.Close()
	}
}

func parseFloat(v interface{}) float64 {
	switch t := v.(type) {
	case string:
		f, _ := strconv.ParseFloat(t, 64)
		return f
	case float64:
		return t
	}
	return 0
}

// ====== пример запуска для песочницы ======
func main() {
	out := make(chan MarketData)
	collector := NewKuCoinCollector([]string{"BTC-USDT", "ETH-USDT"})

	if err := collector.Start(out); err != nil {
		log.Fatal(err)
	}

	go func() {
		for md := range out {
			fmt.Printf("MarketData: %+v\n", md)
		}
	}()

	// Run 30 секунд
	time.Sleep(30 * time.Second)
	collector.Stop()
	close(out)
}

