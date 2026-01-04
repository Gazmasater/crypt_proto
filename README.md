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
	"strconv"
	"strings"
	"time"

	"crypt_proto/pkg/models"

	"github.com/gorilla/websocket"
)

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

func (c *KuCoinCollector) Name() string { return "KuCoin" }

func (c *KuCoinCollector) Start(out chan<- models.MarketData) error {
	if err := c.initWS(); err != nil {
		return err
	}

	conn, _, err := websocket.DefaultDialer.Dial(c.wsURL, nil)
	if err != nil {
		return err
	}
	c.conn = conn

	log.Println("[KuCoin] Connected to WS")

	// subscribe
	for _, s := range c.symbols {
		sym := normalizeKucoinSymbol(s)

		sub := map[string]any{
			"id":       time.Now().UnixNano(),
			"type":     "subscribe",
			"topic":    "/market/ticker:" + sym,
			"response": true,
		}

		if err := conn.WriteJSON(sub); err != nil {
			return err
		}

		log.Println("[KuCoin] Subscribed:", sym)
	}

	go c.pingLoop()
	go c.readLoop(out)

	return nil
}

func (c *KuCoinCollector) Stop() error {
	c.cancel()
	if c.conn != nil {
		_ = c.conn.Close()
	}
	return nil
}

func (c *KuCoinCollector) initWS() error {
	req, err := http.NewRequest(
		"POST",
		"https://api.kucoin.com/api/v1/bullet-public",
		nil,
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("kucoin bullet status: %s", resp.Status)
	}

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

	if len(r.Data.InstanceServers) == 0 {
		return fmt.Errorf("no kucoin ws endpoints")
	}

	c.wsURL = fmt.Sprintf(
		"%s?token=%s&connectId=%d",
		r.Data.InstanceServers[0].Endpoint,
		r.Data.Token,
		time.Now().UnixNano(),
	)

	return nil
}

func (c *KuCoinCollector) pingLoop() {
	t := time.NewTicker(15 * time.Second)
	defer t.Stop()

	for {
		select {
		case <-c.ctx.Done():
			return
		case <-t.C:
			_ = c.conn.WriteJSON(map[string]string{
				"type": "ping",
			})
		}
	}
}

func (c *KuCoinCollector) readLoop(out chan<- models.MarketData) {
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

			var raw map[string]any
			if err := json.Unmarshal(msg, &raw); err != nil {
				continue
			}

			typ, _ := raw["type"].(string)

			// служебные
			if typ == "welcome" || typ == "ack" {
				continue
			}

			if typ != "message" {
				continue
			}

			topic, _ := raw["topic"].(string)
			data, ok := raw["data"].(map[string]any)
			if !ok {
				continue
			}

			symbol := strings.TrimPrefix(topic, "/market/ticker:")

			bid := parseFloat(data["bestBid"])
			ask := parseFloat(data["bestAsk"])

			if bid == 0 || ask == 0 {
				continue
			}

			out <- models.MarketData{
				Exchange: "KuCoin",
				Symbol:   symbol,
				Bid:      bid,
				Ask:      ask,
			}
		}
	}
}

func normalizeKucoinSymbol(s string) string {
	if strings.Contains(s, "-") {
		return s
	}
	if strings.HasSuffix(s, "USDT") {
		return strings.Replace(s, "USDT", "-USDT", 1)
	}
	return s
}

func parseFloat(v any) float64 {
	switch t := v.(type) {
	case string:
		f, _ := strconv.ParseFloat(t, 64)
		return f
	case float64:
		return t
	default:
		return 0
	}
}



		c = collector.NewKuCoinCollector([]string{"BTCUSDT", "ETHUSDT", "ETHBTC"})



