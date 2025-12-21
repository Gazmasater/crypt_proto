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

const mexcWS = "wss://wbs-api.mexc.com/ws"

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
		symbol: strings.ToUpper(symbol), // BTCUSDT
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
			time.Sleep(time.Second)
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

	// Подписка на тикер
	subscribe := map[string]interface{}{
		"method": "sub.ticker",
		"params": []string{"spot@public.ticker.v3.api@" + c.symbol},
		"id":     1,
	}

	if err := conn.WriteJSON(subscribe); err != nil {
		log.Println("[MEXC] subscribe error:", err)
		return
	}

	// Ping каждые 20 секунд
	go func() {
		ticker := time.NewTicker(20 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-c.ctx.Done():
				return
			case <-ticker.C:
				_ = conn.WriteJSON(map[string]interface{}{"method": "ping"})
			}
		}
	}()

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
	// Ответ от сервера может быть разных типов, фильтруем тикер
	var raw struct {
		Method string `json:"method"`
		Params []struct {
			Symbol string `json:"s"`
			Bid    string `json:"b"`
			Ask    string `json:"a"`
		} `json:"params"`
	}

	if err := json.Unmarshal(msg, &raw); err != nil {
		return
	}

	if raw.Method != "ticker.update" || len(raw.Params) == 0 {
		return
	}

	for _, d := range raw.Params {
		bid, err1 := strconv.ParseFloat(d.Bid, 64)
		ask, err2 := strconv.ParseFloat(d.Ask, 64)
		if err1 != nil || err2 != nil {
			continue
		}

		out <- models.MarketData{
			Exchange:  "MEXC",
			Symbol:    d.Symbol,
			Bid:       bid,
			Ask:       ask,
			Timestamp: time.Now().UnixMilli(),
		}
	}
}


az358@gaz358-BOD-WXX9:~/myprog/crypt_proto/cmd/arb$ go run .
EXCHANGE!!!!!!!!! mexc
2025/12/22 00:32:46 Starting collector: MEXC
2025/12/22 00:32:46 [MEXC] connecting...
2025/12/22 00:33:19 [MEXC] read error: websocket: close 1005 (no status)
2025/12/22 00:33:19 [MEXC] reconnect in 1s...
2025/12/22 00:33:20 [MEXC] connecting...
^C2025/12/22 00:33:23 shutdown signal

