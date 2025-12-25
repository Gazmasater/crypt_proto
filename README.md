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

	"crypt_proto/configs"
	"crypt_proto/pkg/models"
	"crypt_proto/internal/market"

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
	// 1) Получаем токен и endpoint
	if err := c.initWS(); err != nil {
		return err
	}

	// 2) Подключаемся к WS
	conn, _, err := websocket.DefaultDialer.Dial(c.wsURL, nil)
	if err != nil {
		return err
	}
	c.conn = conn
	log.Println("[KuCoin] connected")

	// 3) Подписка на Level2/Ticker
	for _, s := range c.symbols {
		// нормализуем символ для KuCoin (например BTCUSDT → BTC-USDT)
		subSymbol := strings.ReplaceAll(market.NormalizeSymbol(s), "/", "-")
		sub := map[string]interface{}{
			"id":       time.Now().UnixNano(),
			"type":     "subscribe",
			"topic":    "level2/ticker:" + subSymbol,
			"response": true,
		}
		if err := conn.WriteJSON(sub); err != nil {
			return err
		}
		log.Println("[KuCoin] subscribed:", subSymbol)
	}

	// 4) Ping loop
	go func() {
		ticker := time.NewTicker(configs.KUCOIN_PING_INTERVAL)
		defer ticker.Stop()
		for {
			select {
			case <-c.ctx.Done():
				return
			case <-ticker.C:
				_ = conn.WriteMessage(websocket.PingMessage, nil)
			}
		}
	}()

	// 5) Read loop
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

// ================= private ==================

func (c *KuCoinCollector) initWS() error {
	resp, err := http.Post(configs.KUCOIN_REST_PUBLIC, "application/json", nil)
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
			} `json:"instanceServers"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return err
	}

	if len(r.Data.InstanceServers) == 0 {
		return fmt.Errorf("no KuCoin WS endpoints returned")
	}

	c.wsURL = fmt.Sprintf("%s?token=%s&connectId=%d",
		r.Data.InstanceServers[0].Endpoint,
		r.Data.Token,
		time.Now().UnixNano(),
	)

	return nil
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

				bid, _ := strconv.ParseFloat(body["bestBid"].(string), 64)
				ask, _ := strconv.ParseFloat(body["bestAsk"].(string), 64)

				out <- models.MarketData{
					Exchange:  "KuCoin",
					Symbol:    market.NormalizeSymbol(symbol),
					Bid:       bid,
					Ask:       ask,
					Timestamp: time.Now().UnixMilli(),
				}
			}
		}
	}
}


