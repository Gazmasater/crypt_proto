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
	}
}

func (c *KuCoinCollector) Start(out chan<- MarketData) error {
	// 1) Получаем WS endpoint и token
	if err := c.initWS(); err != nil {
		return err
	}

	// 2) Делаем REST snapshot для первой котировки
	for _, s := range c.symbols {
		data, err := c.snapshotREST(s)
		if err != nil {
			log.Printf("[KuCoin] snapshot error: %v", err)
			continue
		}
		out <- data
		log.Printf("[KuCoin] Snapshot sent: %s", s)
	}

	// 3) Подключаемся к WS
	conn, _, err := websocket.DefaultDialer.Dial(c.wsURL, nil)
	if err != nil {
		return err
	}
	c.conn = conn
	log.Println("[KuCoin] Connected to WS")

	// 4) Подписка на Level2 ticker
	for _, s := range c.symbols {
		topic := "level2/ticker:" + strings.ReplaceAll(s, "/", "-")
		sub := map[string]interface{}{
			"id":       time.Now().UnixNano(),
			"type":     "subscribe",
			"topic":    topic,
			"response": true,
		}
		if err := conn.WriteJSON(sub); err != nil {
			return err
		}
		log.Printf("[KuCoin] Subscribed: %s", s)
	}

	// 5) Ping loop
	go func() {
		t := time.NewTicker(18 * time.Second) // KuCoin рекомендует ~18s
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

	// 6) Read loop
	go c.readLoop(out)

	return nil
}

func (c *KuCoinCollector) Stop() error {
	c.cancel()
	if c.conn != nil {
		c.conn.Close()
	}
	return nil
}

// ===================== private =====================

func (c *KuCoinCollector) initWS() error {
	resp, err := http.Post("https://api.kucoin.com/api/v1/bullet-public", "application/json", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var r struct {
		Code string `json:"code"`
		Data struct {
			Token           string `json:"token"`
			InstanceServers []struct {
				Endpoint string `json:"endpoint"`
				Encrypt  bool   `json:"encrypt"`
			} `json:"instanceServers"`
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

func (c *KuCoinCollector) snapshotREST(symbol string) (MarketData, error) {
	url := fmt.Sprintf("https://api.kucoin.com/api/v1/market/orderbook/level2_1?symbol=%s", strings.ReplaceAll(symbol, "/", "-"))
	resp, err := http.Get(url)
	if err != nil {
		return MarketData{}, err
	}
	defer resp.Body.Close()

	var r struct {
		Code string `json:"code"`
		Data struct {
			Bids [][]string `json:"bids"`
			Asks [][]string `json:"asks"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return MarketData{}, err
	}

	if len(r.Data.Bids) == 0 || len(r.Data.Asks) == 0 {
		return MarketData{}, fmt.Errorf("empty snapshot for %s", symbol)
	}

	bid, _ := strconv.ParseFloat(r.Data.Bids[0][0], 64)
	ask, _ := strconv.ParseFloat(r.Data.Asks[0][0], 64)

	return MarketData{
		Exchange: "KuCoin",
		Symbol:   symbol,
		Bid:      bid,
		Ask:      ask,
	}, nil
}

func (c *KuCoinCollector) readLoop(out chan<- MarketData) {
	defer func() {
		if c.conn != nil {
			c.conn.Close()
		}
		log.Println("[KuCoin] readLoop stopped")
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

				bidStr, _ := body["bestBid"].(string)
				askStr, _ := body["bestAsk"].(string)

				bid, _ := strconv.ParseFloat(bidStr, 64)
				ask, _ := strconv.ParseFloat(askStr, 64)

				out <- MarketData{
					Exchange: "KuCoin",
					Symbol:   symbol,
					Bid:      bid,
					Ask:      ask,
				}
			}
		}
	}
}

// ===================== main =====================

func main() {
	tickers := []string{"BTC/USDT", "ETH/USDT", "XRP/USDT", "DOGE/USDT"}
	out := make(chan MarketData)

	collector := NewKuCoinCollector(tickers)
	if err := collector.Start(out); err != nil {
		log.Fatal(err)
	}

	// Читаем данные 30 секунд
	timeout := time.After(30 * time.Second)
	for {
		select {
		case md := <-out:
			log.Printf("[Tick] %s %s bid=%.4f ask=%.4f", md.Exchange, md.Symbol, md.Bid, md.Ask)
		case <-timeout:
			log.Println("No more ticks received in 30s, exiting")
			collector.Stop()
			return
		}
	}
}




az358@gaz358-BOD-WXX9:~/myprog/crypt_proto/cmd/arb/arb_test$ go run .
2025/12/25 18:00:15 [KuCoin] snapshot error: empty snapshot for BTC/USDT
2025/12/25 18:00:16 [KuCoin] snapshot error: empty snapshot for ETH/USDT
2025/12/25 18:00:16 [KuCoin] snapshot error: empty snapshot for XRP/USDT
2025/12/25 18:00:16 [KuCoin] snapshot error: empty snapshot for DOGE/USDT
2025/12/25 18:00:18 [KuCoin] Connected to WS
2025/12/25 18:00:18 [KuCoin] Subscribed: BTC/USDT
2025/12/25 18:00:18 [KuCoin] Subscribed: ETH/USDT
2025/12/25 18:00:18 [KuCoin] Subscribed: XRP/USDT
2025/12/25 18:00:18 [KuCoin] Subscribed: DOGE/USDT



