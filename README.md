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
	"crypt_proto/models"
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

const okxWS = "wss://ws.okx.com:8443/ws/v5/public"

type OKXCollector struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func NewOKXCollector() *OKXCollector {
	ctx, cancel := context.WithCancel(context.Background())
	return &OKXCollector{
		ctx:    ctx,
		cancel: cancel,
	}
}

func (c *OKXCollector) Name() string {
	return "OKX"
}

func (c *OKXCollector) Start(out chan<- models.MarketTick) error {
	go c.run(out)
	return nil
}

func (c *OKXCollector) run(out chan<- models.MarketTick) {
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			log.Println("OKX connecting...")
			c.connectAndRead(out)
			log.Println("OKX reconnect in 1s...")
			time.Sleep(time.Second)
		}
	}
}

func (c *OKXCollector) connectAndRead(out chan<- models.MarketTick) {
	conn, _, err := websocket.DefaultDialer.Dial(okxWS, nil)
	if err != nil {
		log.Println("OKX dial error:", err)
		return
	}
	defer conn.Close()

	// subscribe
	sub := map[string]interface{}{
		"op": "subscribe",
		"args": []map[string]string{
			{
				"channel": "tickers",
				"instId":  "BTC-USDT",
			},
		},
	}

	if err := conn.WriteJSON(sub); err != nil {
		log.Println("OKX subscribe error:", err)
		return
	}

	// ping loop
	go func() {
		ticker := time.NewTicker(20 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-c.ctx.Done():
				return
			case <-ticker.C:
				_ = conn.WriteMessage(websocket.PingMessage, []byte("ping"))
			}
		}
	}()

	// read loop
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println("OKX read error:", err)
				return
			}
			c.handleMessage(msg, out)
		}
	}
}

func (c *OKXCollector) handleMessage(msg []byte, out chan<- models.MarketTick) {
	var raw struct {
		Data []struct {
			InstId string `json:"instId"`
			BidPx  string `json:"bidPx"`
			BidSz  string `json:"bidSz"`
			AskPx  string `json:"askPx"`
			AskSz  string `json:"askSz"`
		} `json:"data"`
	}

	if err := json.Unmarshal(msg, &raw); err != nil {
		return
	}

	for _, d := range raw.Data {
		bidPx, _ := strconv.ParseFloat(d.BidPx, 64)
		bidSz, _ := strconv.ParseFloat(d.BidSz, 64)
		askPx, _ := strconv.ParseFloat(d.AskPx, 64)
		askSz, _ := strconv.ParseFloat(d.AskSz, 64)

		out <- models.MarketTick{
			Exchange:  "OKX",
			Symbol:    d.InstId,
			BidPrice: bidPx,
			BidQty:   bidSz,
			AskPrice: askPx,
			AskQty:   askSz,
			Timestamp: time.Now(),
		}
	}
}

func (c *OKXCollector) Stop() error {
	c.cancel()
	return nil
}





package main

import (
	"crypt_proto/internal/collector"
	"crypt_proto/models"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ticks := make(chan models.MarketTick, 1000)

	okx := collector.NewOKXCollector()
	okx.Start(ticks)

	go func() {
		for t := range ticks {
			log.Printf("%s %s bid=%.2f ask=%.2f",
				t.Exchange, t.Symbol, t.BidPrice, t.AskPrice)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	log.Println("shutdown")
	okx.Stop()
}




