package collector

import (
	"context"
	"log"
	"strings"
	"time"

	"crypt_proto/configs"
	"crypt_proto/pkg/models"

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
	conn, _, err := websocket.DefaultDialer.Dial(configs.MEXC_WS, nil)
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
		t := time.NewTicker(configs.MEXC_PING_INTERVAL)
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
