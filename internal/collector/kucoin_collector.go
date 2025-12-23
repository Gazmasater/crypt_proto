package collector

import (
	"context"
	"log"
	"time"

	"crypt_proto/configs"
	"crypt_proto/pkg/models"

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
	conn, _, err := websocket.DefaultDialer.Dial(configs.KUCOIN_WS, nil)
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
		"id":             1,
		"type":           "subscribe",
		"topic":          params,
		"privateChannel": false,
		"response":       true,
	}
	if err := conn.WriteJSON(sub); err != nil {
		return err
	}

	// ping
	go func() {
		t := time.NewTicker(configs.KUCOIN_PING_INTERVAL)
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
