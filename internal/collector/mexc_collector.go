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

const mexcWS = "wss://wbs.mexc.com/ws"

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

	subscribe := map[string]interface{}{
		"method": "SUBSCRIPTION",
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
