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
	"crypt_proto/pkg/models"
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

const (
	mexcWS        = "wss://wbs.mexc.com/ws"
	readTimeout  = 30 * time.Second
	reconnectDur = time.Second
)

type MEXCCollector struct {
	ctx    context.Context
	cancel context.CancelFunc
	symbol string
}

func NewMEXCCollector(symbol string) *MEXCCollector {
	ctx, cancel := context.WithCancel(context.Background())
	return &MEXCCollector{
		ctx:    ctx,
		cancel: cancel,
		symbol: strings.ToUpper(symbol),
	}
}

func (c *MEXCCollector) Name() string {
	return "MEXC"
}

func (c *MEXCCollector) Start(out chan<- models.MarketData) error {
	go c.run(out)
	return nil
}

func (c *MEXCCollector) Stop() error {
	c.cancel()
	return nil
}

func (c *MEXCCollector) run(out chan<- models.MarketData) {
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			log.Println("[MEXC] connecting...")
			c.connectAndRead(out)
			log.Println("[MEXC] reconnect in 1s...")
			time.Sleep(reconnectDur)
		}
	}
}

func (c *MEXCCollector) connectAndRead(out chan<- models.MarketData) {
	conn, _, err := websocket.DefaultDialer.Dial(mexcWS, nil)
	if err != nil {
		log.Println("[MEXC] dial error:", err)
		return
	}
	defer conn.Close()

	// heartbeat
	conn.SetReadDeadline(time.Now().Add(readTimeout))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(readTimeout))
		return nil
	})

	// subscribe
	subscribe := map[string]interface{}{
		"method": "SUBSCRIBE",
		"params": []string{
			"spot@public.bookTicker." + c.symbol,
		},
	}

	if err := conn.WriteJSON(subscribe); err != nil {
		log.Println("[MEXC] subscribe error:", err)
		return
	}

	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println("[MEXC] read error:", err)
				return
			}
			c.handleMessage(msg, out)
		}
	}
}

func (c *MEXCCollector) handleMessage(msg []byte, out chan<- models.MarketData) {
	var raw struct {
		Data struct {
			Symbol string `json:"s"`
			Bid    string `json:"b"`
			Ask    string `json:"a"`
		} `json:"d"`
	}

	if err := json.Unmarshal(msg, &raw); err != nil {
		return
	}

	if raw.Data.Symbol == "" {
		return
	}

	bid, err1 := strconv.ParseFloat(raw.Data.Bid, 64)
	ask, err2 := strconv.ParseFloat(raw.Data.Ask, 64)
	if err1 != nil || err2 != nil {
		return
	}

	out <- models.MarketData{
		Exchange:  "MEXC",
		Symbol:    raw.Data.Symbol,
		Bid:       bid,
		Ask:       ask,
		Timestamp: time.Now().UnixMilli(),
	}
}



