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
	"log"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

const (
	kucoinWS       = "wss://ws.kucoin.com/endpoint"
	pingInterval   = 30 * time.Second
	readTimeout    = 60 * time.Second
)

type KuCoinCollector struct {
	ctx     context.Context
	cancel  context.CancelFunc
	symbols []string
	conn    *websocket.Conn
}

func NewKuCoinCollector(symbols []string) *KuCoinCollector {
	ctx, cancel := context.WithCancel(context.Background())
	return &KuCoinCollector{
		ctx:     ctx,
		cancel:  cancel,
		symbols: symbols,
	}
}

func (c *KuCoinCollector) Start() error {
	conn, _, err := websocket.DefaultDialer.Dial(kucoinWS, nil)
	if err != nil {
		return err
	}
	c.conn = conn
	log.Println("[KuCoin] connected")

	if err := c.subscribe(); err != nil {
		return err
	}

	go c.pingLoop()
	go c.readLoop()

	return nil
}

func (c *KuCoinCollector) Stop() {
	c.cancel()
	if c.conn != nil {
		_ = c.conn.Close()
	}
}

func (c *KuCoinCollector) subscribe() error {
	// KuCoin требует init message для публичного ws
	initMsg := map[string]interface{}{
		"id":     time.Now().UnixMilli(),
		"type":   "subscribe",
		"topic":  "/market/ticker:all",
		"response": true,
	}

	if err := c.conn.WriteJSON(initMsg); err != nil {
		return err
	}

	log.Println("[KuCoin] subscribed to ticker:all")
	return nil
}

func (c *KuCoinCollector) pingLoop() {
	t := time.NewTicker(pingInterval)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			_ = c.conn.WriteMessage(websocket.PingMessage, []byte("hb"))
		case <-c.ctx.Done():
			return
		}
	}
}

func (c *KuCoinCollector) readLoop() {
	_ = c.conn.SetReadDeadline(time.Now().Add(readTimeout))

	for {
		select {
		case <-c.ctx.Done():
			return
		default:
		}

		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			log.Printf("[KuCoin] read error: %v", err)
			return
		}

		_ = c.conn.SetReadDeadline(time.Now().Add(readTimeout))

		var data map[string]interface{}
		if err := json.Unmarshal(msg, &data); err != nil {
			continue
		}

		if d, ok := data["data"].(map[string]interface{}); ok {
			symbol := strings.ToUpper(d["s"].(string))
			bid := d["b"].(string)
			ask := d["a"].(string)
			log.Printf("[KuCoin] %s bid=%s ask=%s", symbol, bid, ask)
		}
	}
}





package main

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"crypt_proto/pkg/collector"
)

func main() {
	exchange := strings.ToLower(os.Getenv("EXCHANGE"))
	if exchange == "" {
		exchange = "mexc"
	}
	log.Println("EXCHANGE:", exchange)

	var c collector.Collector

	switch exchange {
	case "mexc":
		c = collector.NewMEXCCollector([]string{
			"BTCUSDT",
			"ETHUSDT",
			"ETHBTC",
		})
	case "kucoin":
		c = collector.NewKuCoinCollector([]string{
			"BTC-USDT",
			"ETH-USDT",
			"ETH-BTC",
		})
	case "okx":
		c = collector.NewOKXCollector([]string{
			"BTC-USDT",
			"ETH-USDT",
			"ETH-BTC",
		})
	default:
		log.Fatalf("Unknown exchange: %s", exchange)
	}

	if err := c.Start(); err != nil {
		log.Fatal(err)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	log.Println("Stopping collector...")
	c.Stop()
}
