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



parts := strings.Split(ch, "@")
symbol = parts[len(parts)-1]





package collector

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"time"

	"crypt_proto/configs"
	"crypt_proto/pkg/models"

	"github.com/gorilla/websocket"
)

type OKXCollector struct {
	ctx     context.Context
	cancel  context.CancelFunc
	symbols []string
	conn    *websocket.Conn
}

func NewOKXCollector(symbols []string) *OKXCollector {
	ctx, cancel := context.WithCancel(context.Background())
	return &OKXCollector{ctx: ctx, cancel: cancel, symbols: symbols}
}

func (c *OKXCollector) Name() string { return "OKX" }

func (c *OKXCollector) Start(out chan<- models.MarketData) error {
	conn, _, err := websocket.DefaultDialer.Dial(configs.OKX_WS, nil)
	if err != nil {
		return err
	}
	c.conn = conn
	log.Println("[OKX] connected")

	// Подписка на книги заявок
	args := make([]map[string]string, 0, len(c.symbols))
	for _, s := range c.symbols {
		args = append(args, map[string]string{
			"channel": "books5",
			"instId":  s,
		})
	}
	sub := map[string]interface{}{
		"op":   "subscribe",
		"args": args,
	}
	if err := conn.WriteJSON(sub); err != nil {
		return err
	}

	// Ping
	go func() {
		t := time.NewTicker(configs.OKX_PING_INTERVAL)
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

	// Read loop
	go func() {
		defer conn.Close()
		for {
			select {
			case <-c.ctx.Done():
				return
			default:
				_, msg, err := conn.ReadMessage()
				if err != nil {
					log.Println("[OKX] read error:", err)
					return
				}

				var resp struct {
					Arg struct {
						InstID string `json:"instId"`
					} `json:"arg"`
					Data []struct {
						Asks [][]string `json:"asks"`
						Bids [][]string `json:"bids"`
					} `json:"data"`
				}

				if err := json.Unmarshal(msg, &resp); err != nil {
					continue
				}

				if len(resp.Data) == 0 {
					continue
				}

				var bid, ask float64
				if len(resp.Data[0].Bids) > 0 {
					bid, _ = strconv.ParseFloat(resp.Data[0].Bids[0][0], 64)
				}
				if len(resp.Data[0].Asks) > 0 {
					ask, _ = strconv.ParseFloat(resp.Data[0].Asks[0][0], 64)
				}

				out <- models.MarketData{
					Exchange: "OKX",
					Symbol:   resp.Arg.InstID,
					Bid:      bid,
					Ask:      ask,
				}
			}
		}
	}()

	return nil
}

func (c *OKXCollector) Stop() error {
	c.cancel()
	return nil
}



