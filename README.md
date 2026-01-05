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





package collector

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"crypt_proto/pkg/models"

	"github.com/gorilla/websocket"
)

type KuCoinCollector struct {
	ctx     context.Context
	cancel  context.CancelFunc
	conn    *websocket.Conn
	wsURL   string
	symbols []string

	out chan<- *models.MarketData
	mu  sync.Mutex
}

func NewKuCoinCollector(symbols []string) *KuCoinCollector {
	ctx, cancel := context.WithCancel(context.Background())

	return &KuCoinCollector{
		ctx:     ctx,
		cancel:  cancel,
		wsURL:   "wss://ws-api-spot.kucoin.com/?token=public",
		symbols: symbols,
	}
}

func (c *KuCoinCollector) Name() string {
	return "kucoin"
}

func (c *KuCoinCollector) Start(out chan<- *models.MarketData) error {
	c.out = out

	if err := c.initWS(); err != nil {
		return err
	}

	if err := c.subscribe(); err != nil {
		return err
	}

	go c.readLoop()
	go c.pingLoop()

	return nil
}

func (c *KuCoinCollector) Stop() error {
	c.cancel()
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

//
// ------------------- WS INIT -------------------
//

func (c *KuCoinCollector) initWS() error {
	dialer := websocket.Dialer{
		HandshakeTimeout: 10 * time.Second,
	}

	headers := http.Header{}
	headers.Set("Origin", "https://www.kucoin.com")
	headers.Set("User-Agent", "Mozilla/5.0")

	conn, _, err := dialer.Dial(c.wsURL, headers)
	if err != nil {
		return fmt.Errorf("kucoin ws dial error: %w", err)
	}

	c.conn = conn
	log.Println("[KUCOIN] WS connected")

	return nil
}

//
// ------------------- SUBSCRIBE -------------------
//

func (c *KuCoinCollector) subscribe() error {
	msg := map[string]interface{}{
		"id":              "sub-ticker",
		"type":            "subscribe",
		"topic":           "/market/ticker:" + joinSymbols(c.symbols),
		"privateChannel":  false,
		"response":        true,
	}

	return c.conn.WriteJSON(msg)
}

//
// ------------------- PING -------------------
//

func (c *KuCoinCollector) pingLoop() {
	ticker := time.NewTicker(25 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-c.ctx.Done():
			return
		case <-ticker.C:
			_ = c.conn.WriteJSON(map[string]string{
				"type": "ping",
			})
		}
	}
}

//
// ------------------- READ -------------------
//

func (c *KuCoinCollector) readLoop() {
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			_, msg, err := c.conn.ReadMessage()
			if err != nil {
				log.Printf("[KUCOIN] read error: %v", err)
				return
			}
			c.handleMessage(msg)
		}
	}
}

//
// ------------------- HANDLE -------------------
//

func (c *KuCoinCollector) handleMessage(msg []byte) {
	var raw map[string]interface{}
	if err := json.Unmarshal(msg, &raw); err != nil {
		return
	}

	switch raw["type"] {
	case "welcome", "pong", "ack":
		return
	case "message":
		data, ok := raw["data"].(map[string]interface{})
		if !ok {
			return
		}

		symbol, _ := data["symbol"].(string)
		priceStr, _ := data["price"].(string)
		price, err := parseFloat(priceStr)
		if err != nil {
			return
		}

		c.out <- &models.MarketData{
			Exchange: c.Name(),
			Symbol:   symbol,
			Price:    price,
			Time:     time.Now(),
		}
	}
}

//
// ------------------- HELPERS -------------------
//

func joinSymbols(symbols []string) string {
	res := ""
	for i, s := range symbols {
		if i > 0 {
			res += ","
		}
		res += s
	}
	return res
}

func parseFloat(v string) (float64, error) {
	var f float64
	_, err := fmt.Sscan(v, &f)
	return f, err
}




