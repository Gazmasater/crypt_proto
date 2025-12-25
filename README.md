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
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

// ---------------------- types ----------------------

type MarketData struct {
	Exchange  string
	Symbol    string
	Bid       float64
	Ask       float64
	Timestamp int64
}

// ---------------------- KuCoinCollector ----------------------

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
	}
}

func (c *KuCoinCollector) Start(out chan<- MarketData) error {
	if err := c.initWS(); err != nil {
		return err
	}

	conn, _, err := websocket.DefaultDialer.Dial(c.wsURL, nil)
	if err != nil {
		return err
	}
	c.conn = conn
	log.Println("[KuCoin] Connected to WS")

	// 1) Snapshot
	for _, s := range c.symbols {
		snapSub := map[string]interface{}{
			"id":       time.Now().UnixNano(),
			"type":     "subscribe",
			"topic":    "level2/snapshot:" + s,
			"response": true,
		}
		if err := conn.WriteJSON(snapSub); err != nil {
			return err
		}
		log.Println("[KuCoin] Snapshot subscribed:", s)
	}

	// 2) Level2/Ticker
	for _, s := range c.symbols {
		tickerSub := map[string]interface{}{
			"id":       time.Now().UnixNano(),
			"type":     "subscribe",
			"topic":    "level2/ticker:" + s,
			"response": true,
		}
		if err := conn.WriteJSON(tickerSub); err != nil {
			return err
		}
		log.Println("[KuCoin] Subscribed:", s)
	}

	// 3) Ping loop
	go func() {
		t := time.NewTicker(18 * time.Second)
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

	// 4) Read loop
	go c.readLoop(out)

	return nil
}

func (c *KuCoinCollector) Stop() {
	c.cancel()
	if c.conn != nil {
		c.conn.Close()
	}
}

// ---------------------- private ----------------------

func (c *KuCoinCollector) initWS() error {
	const bulletPublic = "https://api.kucoin.com/api/v1/bullet-public"
	resp, err := http.Post(bulletPublic, "application/json", strings.NewReader("{}"))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var r struct {
		Code string `json:"code"`
		Data struct {
			InstanceServers []struct {
				Endpoint string `json:"endpoint"`
				Encrypt  bool   `json:"encrypt"`
			} `json:"instanceServers"`
			Token string `json:"token"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return err
	}

	if len(r.Data.InstanceServers) == 0 {
		return fmt.Errorf("no KuCoin WS endpoints returned")
	}

	endpoint := r.Data.InstanceServers[0].Endpoint
	c.wsURL = fmt.Sprintf("%s?token=%s&connectId=%d", endpoint, r.Data.Token, time.Now().UnixNano())
	return nil
}

func (c *KuCoinCollector) readLoop(out chan<- MarketData) {
	defer func() {
		if c.conn != nil {
			c.conn.Close()
		}
	}()

	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			_, msg, err := c.conn.ReadMessage()
			if err != nil {
				log.Println("[KuCoin] read error:", err)
				return
			}

			var data map[string]interface{}
			if err := json.Unmarshal(msg, &data); err != nil {
				continue
			}

			if data["type"] != "message" {
				continue
			}

			topic, ok := data["topic"].(string)
			if !ok {
				continue
			}
			symbol := strings.Split(topic, ":")[1]

			body, ok := data["data"].(map[string]interface{})
			if !ok {
				continue
			}

			var bid, ask float64

			if b, ok := body["bestBid"].(string); ok {
				bid, _ = strconv.ParseFloat(b, 64)
			} else if b, ok := body["bids"].([]interface{}); ok && len(b) > 0 {
				bid, _ = strconv.ParseFloat(b[0].([]interface{})[0].(string), 64)
			}

			if a, ok := body["bestAsk"].(string); ok {
				ask, _ = strconv.ParseFloat(a, 64)
			} else if a, ok := body["asks"].([]interface{}); ok && len(a) > 0 {
				ask, _ = strconv.ParseFloat(a[0].([]interface{})[0].(string), 64)
			}

			if bid > 0 && ask > 0 {
				out <- MarketData{
					Exchange:  "KuCoin",
					Symbol:    symbol,
					Bid:       bid,
					Ask:       ask,
					Timestamp: time.Now().UnixMilli(),
				}
			}
		}
	}
}

// ---------------------- main ----------------------

func main() {
	out := make(chan MarketData, 10)
	collector := NewKuCoinCollector([]string{"BTC-USDT", "ETH-USDT", "XRP-USDT", "DOGE-USDT"})

	if err := collector.Start(out); err != nil {
		log.Fatal(err)
	}
	defer collector.Stop()

	timeout := time.After(30 * time.Second)
	for {
		select {
		case tick := <-out:
			fmt.Printf("[%s] %s: Bid=%.4f Ask=%.4f\n", tick.Exchange, tick.Symbol, tick.Bid, tick.Ask)
		case <-timeout:
			log.Println("No more ticks received in 30s, exiting")
			return
		}
	}
}


2025/12/25 17:52:46 [KuCoin] Connected to WS
2025/12/25 17:52:46 [KuCoin] Snapshot subscribed: BTC-USDT
2025/12/25 17:52:46 [KuCoin] Snapshot subscribed: ETH-USDT
2025/12/25 17:52:46 [KuCoin] Snapshot subscribed: XRP-USDT
2025/12/25 17:52:46 [KuCoin] Snapshot subscribed: DOGE-USDT
2025/12/25 17:52:46 [KuCoin] Subscribed: BTC-USDT
2025/12/25 17:52:46 [KuCoin] Subscribed: ETH-USDT
2025/12/25 17:52:46 [KuCoin] Subscribed: XRP-USDT
2025/12/25 17:52:46 [KuCoin] Subscribed: DOGE-USDT
2025/12/25 17:53:16 No more ticks received in 30s, exiting
2025/12/25 17:53:16 [KuCoin] read error: read tcp 192.168.1.71:42788->108.157.229.12:443: use of closed network connection
gaz358@gaz358-BOD-WXX9:~/myprog/crypt_proto/cmd/arb/arb_test$ 



