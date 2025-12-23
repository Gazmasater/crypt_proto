package collector

import (
	"context"
	"log"
	"time"

	"crypt_proto/configs"
	"crypt_proto/pkg/models"

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
	conn, _, err := websocket.DefaultDialer.Dial(configs.OKX_WS, nil)
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
