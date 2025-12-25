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
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

// MarketData хранит данные стакана
type MarketData struct {
	Exchange string
	Symbol   string
	Bid      float64
	Ask      float64
}

// KuCoinCollector управляет подключением к WS
type KuCoinCollector struct {
	ctx     context.Context
	cancel  context.CancelFunc
	symbols []string
	conn    *websocket.Conn
	wsURL   string
}

// NewKuCoinCollector создаёт новый объект
func NewKuCoinCollector(symbols []string) *KuCoinCollector {
	ctx, cancel := context.WithCancel(context.Background())
	return &KuCoinCollector{
		ctx:     ctx,
		cancel:  cancel,
		symbols: symbols,
	}
}

func (c *KuCoinCollector) Name() string { return "KuCoin" }

// Start подключается к KuCoin WS
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

	// подписка на level2/ticker
	for _, s := range c.symbols {
		topicSymbol := strings.ReplaceAll(s, "/", "-") // BTC/USDT -> BTC-USDT
		sub := map[string]interface{}{
			"id":       time.Now().UnixNano(),
			"type":     "subscribe",
			"topic":    "level2/ticker:" + topicSymbol,
			"response": true,
		}
		if err := conn.WriteJSON(sub); err != nil {
			return err
		}
		log.Println("[KuCoin] Subscribed:", topicSymbol)
	}

	// ping loop
	go func() {
		t := time.NewTicker(10 * time.Second)
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

	// read loop
	go c.readLoop(out)
	return nil
}

// Stop закрывает WS
func (c *KuCoinCollector) Stop() error {
	c.cancel()
	if c.conn != nil {
		c.conn.Close()
	}
	return nil
}

// initWS получает bullet-public токен и endpoint
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

// readLoop обрабатывает входящие сообщения
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

				bid, _ := parseStringToFloat(body["bestBid"].(string))
				ask, _ := parseStringToFloat(body["bestAsk"].(string))

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

func parseStringToFloat(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

// ======= main =======
func main() {
	out := make(chan MarketData, 10)
	collector := NewKuCoinCollector([]string{"BTC/USDT", "ETH/USDT", "XRP/USDT", "DOGE/USDT"})

	err := collector.Start(out)
	if err != nil {
		log.Fatal(err)
	}
	defer collector.Stop()

	timeout := time.After(10 * time.Second)
	for {
		select {
		case data := <-out:
			log.Printf("Ticker: %s Bid: %.4f Ask: %.4f", data.Symbol, data.Bid, data.Ask)
		case <-timeout:
			log.Println("No ticks received in 10 seconds")
			return
		}
	}
}



