package collector

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"time"

	"crypt_proto/configs"
	"crypt_proto/pkg/models"

	"github.com/gorilla/websocket"
)

type KuCoinCollector struct {
	ctx     context.Context
	cancel  context.CancelFunc
	symbols []string
	conn    *websocket.Conn
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
	c.conn = conn
	log.Println("[KuCoin] connected")

	// subscribe
	for _, s := range c.symbols {
		sub := map[string]interface{}{
			"id":       time.Now().Unix(),
			"type":     "subscribe",
			"topic":    "level2/ticker:" + s,
			"response": true,
		}
		if err := conn.WriteJSON(sub); err != nil {
			return err
		}
		log.Println("[KuCoin] subscribed:", s)
	}

	// ping loop
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

				var data map[string]interface{}
				if err := json.Unmarshal(msg, &data); err != nil {
					continue
				}

				// проверяем, что это update
				if data["type"] == "message" {
					topic := data["topic"].(string)
					symbol := strings.Split(topic, ":")[1]

					body := data["data"].(map[string]interface{})
					bid, _ := parseStringToFloat(body["bestBid"].(string))
					ask, _ := parseStringToFloat(body["bestAsk"].(string))

					out <- models.MarketData{
						Exchange: "KuCoin",
						Symbol:   symbol,
						Bid:      bid,
						Ask:      ask,
					}
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

// вспомогательная функция
func parseStringToFloat(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}
