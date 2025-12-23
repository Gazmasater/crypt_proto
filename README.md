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






1️⃣ Интерфейс Collector

pkg/collector/collector.go

package collector

import "crypt_proto/pkg/models"

type Collector interface {
	Name() string
	Start(out chan<- models.MarketData) error
	Stop() error
}

2️⃣ MEXCCollector

pkg/collector/mexc_collector.go

package collector

import (
	"context"
	"encoding/json"
	"log"
	"strings"
	"time"

	"crypt_proto/pkg/models"
	"crypt_proto/pkg/config"

	"github.com/gorilla/websocket"
)

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

func (c *MEXCCollector) Name() string { return "MEXC" }

func (c *MEXCCollector) Start(out chan<- models.MarketData) error {
	conn, _, err := websocket.DefaultDialer.Dial(config.MEXC_WS, nil)
	if err != nil {
		return err
	}

	// --- subscribe ---
	params := make([]string, 0, len(c.symbols))
	for _, s := range c.symbols {
		params = append(params, "spot@public.aggre.bookTicker.v3.api.pb@100ms@"+s)
	}
	sub := map[string]interface{}{
		"method": "SUBSCRIPTION",
		"params": params,
	}
	if err := conn.WriteJSON(sub); err != nil {
		return err
	}

	// ping
	go func() {
		t := time.NewTicker(config.MEXC_PING_INTERVAL)
		defer t.Stop()
		for {
			select {
			case <-c.ctx.Done():
				return
			case <-t.C:
				_ = conn.WriteJSON(map[string]string{"method": "PING"})
			}
		}
	}()

	// read loop
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
				// просто отправляем RAW message в канал
				out <- models.MarketData{
					Exchange: "MEXC",
					Symbol:   "", // для упрощения можно распарсить msg при желании
					Bid:      0,
					Ask:      0,
				}
				_ = msg
			}
		}
	}()

	return nil
}

func (c *MEXCCollector) Stop() error {
	c.cancel()
	return nil
}

3️⃣ KuCoinCollector

pkg/collector/kucoin_collector.go

package collector

import (
	"context"
	"log"
	"time"

	"crypt_proto/pkg/models"
	"crypt_proto/pkg/config"

	"github.com/gorilla/websocket"
)

type KuCoinCollector struct {
	ctx     context.Context
	cancel  context.CancelFunc
	symbols []string
}

func NewKuCoinCollector(symbols []string) *KuCoinCollector {
	ctx, cancel := context.WithCancel(context.Background())
	return &KuCoinCollector{ctx: ctx, cancel: cancel, symbols: symbols}
}

func (c *KuCoinCollector) Name() string { return "KuCoin" }

func (c *KuCoinCollector) Start(out chan<- models.MarketData) error {
	conn, _, err := websocket.DefaultDialer.Dial(config.KUCOIN_WS, nil)
	if err != nil {
		return err
	}

	// subscribe
	// KuCoin требует "subscribe": "level2/ticker:BTC-USDT" и т.д
	params := make([]string, 0, len(c.symbols))
	for _, s := range c.symbols {
		params = append(params, "level2/ticker:"+s)
	}
	sub := map[string]interface{}{
		"id":      1,
		"type":    "subscribe",
		"topic":   params,
		"privateChannel": false,
		"response": true,
	}
	if err := conn.WriteJSON(sub); err != nil {
		return err
	}

	// ping
	go func() {
		t := time.NewTicker(config.KUCOIN_PING_INTERVAL)
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
	go func() {
		defer conn.Close()
		for {
			select {
			case <-c.ctx.Done():
				return
			default:
				_, msg, err := conn.ReadMessage()
				if err != nil {
					log.Println("[KuCoin] read error:", err)
					return
				}
				_ = msg
				// parse if needed
				out <- models.MarketData{
					Exchange: "KuCoin",
					Symbol:   "",
					Bid:      0,
					Ask:      0,
				}
			}
		}
	}()

	return nil
}

func (c *KuCoinCollector) Stop() error {
	c.cancel()
	return nil
}

4️⃣ OKXCollector

pkg/collector/okx_collector.go

package collector

import (
	"context"
	"log"
	"time"

	"crypt_proto/pkg/models"
	"crypt_proto/pkg/config"

	"github.com/gorilla/websocket"
)

type OKXCollector struct {
	ctx     context.Context
	cancel  context.CancelFunc
	symbols []string
}

func NewOKXCollector(symbols []string) *OKXCollector {
	ctx, cancel := context.WithCancel(context.Background())
	return &OKXCollector{ctx: ctx, cancel: cancel, symbols: symbols}
}

func (c *OKXCollector) Name() string { return "OKX" }

func (c *OKXCollector) Start(out chan<- models.MarketData) error {
	conn, _, err := websocket.DefaultDialer.Dial(config.OKX_WS, nil)
	if err != nil {
		return err
	}

	// subscribe
	params := make([]map[string]string, 0, len(c.symbols))
	for _, s := range c.symbols {
		params = append(params, map[string]string{
			"channel": "books5",
			"instId":  s,
		})
	}
	sub := map[string]interface{}{
		"op":   "subscribe",
		"args": params,
	}
	if err := conn.WriteJSON(sub); err != nil {
		return err
	}

	// ping
	go func() {
		t := time.NewTicker(config.OKX_PING_INTERVAL)
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
				_ = msg
				out <- models.MarketData{
					Exchange: "OKX",
					Symbol:   "",
					Bid:      0,
					Ask:      0,
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



