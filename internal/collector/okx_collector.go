package collector

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"time"

	"crypt_proto/configs"
	"crypt_proto/internal/market"
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

				// === НОРМАЛИЗАЦИЯ ===
				symbol := market.NormalizeSymbol_Full(resp.Arg.InstID)
				if symbol == "" {
					continue
				}

				out <- models.MarketData{
					Exchange: "OKX",
					Symbol:   symbol,
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
