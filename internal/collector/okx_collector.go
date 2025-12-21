package collector

import (
	"context"
	"crypt_proto/pkg/models"
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

func (c *OKXCollector) Start(out chan<- models.MarketData) error {
	go c.run(out)
	return nil
}

func (c *OKXCollector) Stop() error {
	c.cancel()
	return nil
}

func (c *OKXCollector) run(out chan<- models.MarketData) {
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			log.Println("[OKX] connecting...")
			c.connectAndRead(out)
			log.Println("[OKX] reconnect in 1s...")
			time.Sleep(time.Second)
		}
	}
}

func (c *OKXCollector) connectAndRead(out chan<- models.MarketData) {
	conn, _, err := websocket.DefaultDialer.Dial(okxWS, nil)
	if err != nil {
		log.Println("[OKX] dial error:", err)
		return
	}
	defer conn.Close()

	subscribe := map[string]interface{}{
		"op": "subscribe",
		"args": []map[string]string{
			{
				"channel": "tickers",
				"instId":  "BTC-USDT",
			},
		},
	}

	if err := conn.WriteJSON(subscribe); err != nil {
		log.Println("[OKX] subscribe error:", err)
		return
	}

	// keepalive ping
	go func() {
		ticker := time.NewTicker(20 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-c.ctx.Done():
				return
			case <-ticker.C:
				_ = conn.WriteMessage(websocket.PingMessage, nil)
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
				log.Println("[OKX] read error:", err)
				return
			}
			c.handleMessage(msg, out)
		}
	}
}

func (c *OKXCollector) handleMessage(msg []byte, out chan<- models.MarketData) {
	var raw struct {
		Data []struct {
			InstId string `json:"instId"`
			BidPx  string `json:"bidPx"`
			AskPx  string `json:"askPx"`
		} `json:"data"`
	}

	if err := json.Unmarshal(msg, &raw); err != nil {
		return
	}

	for _, d := range raw.Data {
		bid, err1 := strconv.ParseFloat(d.BidPx, 64)
		ask, err2 := strconv.ParseFloat(d.AskPx, 64)
		if err1 != nil || err2 != nil {
			continue
		}

		out <- models.MarketData{
			Exchange:  "OKX",
			Symbol:    d.InstId,
			Bid:       bid,
			Ask:       ask,
			Timestamp: time.Now().UnixMilli(),
		}
	}
}
