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

	"crypt_proto/pkg/models"

	"github.com/gorilla/websocket"
)

const mexcWS = "wss://wbs-api.mexc.com/ws"

type MEXCCollector struct {
	ctx     context.Context
	cancel  context.CancelFunc
	symbols []string
}

func NewMEXCCollector(symbols []string) *MEXCCollector {
	ctx, cancel := context.WithCancel(context.Background())

	up := make([]string, 0, len(symbols))
	for _, s := range symbols {
		up = append(up, strings.ToUpper(s))
	}

	return &MEXCCollector{
		ctx:     ctx,
		cancel:  cancel,
		symbols: up,
	}
}

func (c *MEXCCollector) Name() string {
	return "mexc"
}

func (c *MEXCCollector) Start(out chan<- models.MarketData) error {
	conn, _, err := websocket.DefaultDialer.Dial(mexcWS, nil)
	if err != nil {
		return err
	}

	log.Println("[MEXC] connected")

	// --- subscribe ---
	params := make([]string, 0, len(c.symbols))
	for _, s := range c.symbols {
		params = append(params, "spot@public.bookTicker.batch@"+s)
	}

	sub := map[string]interface{}{
		"method": "SUBSCRIPTION",
		"params": params,
	}

	if err := conn.WriteJSON(sub); err != nil {
		return err
	}

	log.Println("[MEXC] subscribed:", params)

	// --- ping loop ---
	go func() {
		t := time.NewTicker(20 * time.Second)
		defer t.Stop()
		for {
			select {
			case <-c.ctx.Done():
				return
			case <-t.C:
				_ = conn.WriteJSON(map[string]string{
					"method": "PING",
				})
			}
		}
	}()

	// --- read loop ---
	go func() {
		defer conn.Close()
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
	}()

	return nil
}

func (c *MEXCCollector) Stop() error {
	c.cancel()
	return nil
}

/*
MEXC bookTicker.batch пример:
{
  "d": {
    "s": "BTCUSDT",
    "b": "43500.1",
    "a": "43500.2"
  }
}
*/
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

	out <- models.MarketData{
		Exchange:  "MEXC",
		Symbol:    raw.Data.Symbol,
		Bid:       parseFloat(raw.Data.Bid),
		Ask:       parseFloat(raw.Data.Ask),
		Timestamp: time.Now().UnixMilli(),
	}
}

func parseFloat(s string) float64 {
	v, _ := strconv.ParseFloat(s, 64)
	return v
}


az358@gaz358-BOD-WXX9:~/myprog/crypt_proto/cmd/arb$ go run .
EXCHANGE: mexc
Starting collector: mexc
2025/12/23 02:39:10 [MEXC] connected
2025/12/23 02:39:10 [MEXC] subscribed: [spot@public.bookTicker.batch@BTCUSDT spot@public.bookTicker.batch@ETHUSDT]
2025/12/23 02:39:41 [MEXC] read error: websocket: close 1005 (no status)
^Csignal: interrupt
