package collector

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
	wsURL   string
}

func NewKuCoinCollector(symbols []string) *KuCoinCollector {
	ctx, cancel := context.WithCancel(context.Background())
	return &KuCoinCollector{
		ctx:     ctx,
		cancel:  cancel,
		symbols: symbols,
	}
}

func (c *KuCoinCollector) Name() string { return "KuCoin" }

func (c *KuCoinCollector) Start(out chan<- models.MarketData) error {
	// 1) Получаем токен для WS
	if err := c.initWS(); err != nil {
		return err
	}

	// 2) Подключаемся
	conn, _, err := websocket.DefaultDialer.Dial(c.wsURL, nil)
	if err != nil {
		return err
	}
	c.conn = conn
	log.Println("[KuCoin] connected")

	// 3) ping loop
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

	// 4) сначала REST-снапшот для каждой пары
	for _, s := range c.symbols {
		if err := c.fetchSnapshot(s, out); err != nil {
			log.Println("[KuCoin] snapshot error:", err)
		}
	}

	// 5) подписка на инкременты Level2
	for _, s := range c.symbols {
		sub := map[string]interface{}{
			"id":       time.Now().UnixNano(),
			"type":     "subscribe",
			"topic":    "level2:" + strings.ToUpper(s),
			"response": true,
		}
		if err := conn.WriteJSON(sub); err != nil {
			return err
		}
		log.Println("[KuCoin] subscribed:", s)
	}

	// 6) read loop
	go c.readLoop(out)

	return nil
}

func (c *KuCoinCollector) Stop() error {
	c.cancel()
	if c.conn != nil {
		c.conn.Close()
	}
	return nil
}

// ================= private ==================

func (c *KuCoinCollector) initWS() error {
	resp, err := http.Post(configs.KUCOIN_REST_PUBLIC, "application/json", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var r struct {
		Code string `json:"code"`
		Data struct {
			InstanceServers []struct {
				Endpoint string `json:"endpoint"`
			} `json:"instanceServers"`
			Token string `json:"token"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return err
	}

	if len(r.Data.InstanceServers) == 0 {
		return fmt.Errorf("no KuCoin WS endpoints returned")
	}

	endpoint := r.Data.InstanceServers[0].Endpoint
	c.wsURL = fmt.Sprintf("%s?token=%s&connectId=%d", endpoint, r.Data.Token, time.Now().UnixNano())
	return nil
}

// REST snapshot для Level2 top100
func (c *KuCoinCollector) fetchSnapshot(symbol string, out chan<- models.MarketData) error {
	url := fmt.Sprintf("https://api.kucoin.com/api/v1/market/orderbook/level2_100?symbol=%s", strings.ToUpper(symbol))
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var snap struct {
		Code string `json:"code"`
		Data struct {
			Time int64      `json:"time"`
			Bids [][]string `json:"bids"`
			Asks [][]string `json:"asks"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&snap); err != nil {
		return err
	}

	if len(snap.Data.Bids) == 0 || len(snap.Data.Asks) == 0 {
		return fmt.Errorf("empty snapshot for %s", symbol)
	}

	// Берем лучший bid/ask
	bid, _ := strconv.ParseFloat(snap.Data.Bids[0][0], 64)
	ask, _ := strconv.ParseFloat(snap.Data.Asks[0][0], 64)

	out <- models.MarketData{
		Exchange:  "KuCoin",
		Symbol:    symbol,
		Bid:       bid,
		Ask:       ask,
		Timestamp: snap.Data.Time,
	}

	return nil
}

// readLoop обрабатывает инкрементальные обновления
func (c *KuCoinCollector) readLoop(out chan<- models.MarketData) {
	defer func() {
		if c.conn != nil {
			c.conn.Close()
		}
	}()

	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			_, msg, err := c.conn.ReadMessage()
			if err != nil {
				log.Println("[KuCoin] read error:", err)
				return
			}

			var data map[string]interface{}
			if err := json.Unmarshal(msg, &data); err != nil {
				continue
			}

			if data["type"] == "message" {
				topic, ok := data["topic"].(string)
				if !ok {
					continue
				}
				symbol := strings.Split(topic, ":")[1]

				body, ok := data["data"].(map[string]interface{})
				if !ok {
					continue
				}

				bids, okB := body["bids"].([]interface{})
				asks, okA := body["asks"].([]interface{})
				if !okB || !okA || len(bids) == 0 || len(asks) == 0 {
					continue
				}

				bestBid := bids[0].([]interface{})[0].(string)
				bestAsk := asks[0].([]interface{})[0].(string)
				bid, _ := strconv.ParseFloat(bestBid, 64)
				ask, _ := strconv.ParseFloat(bestAsk, 64)

				out <- models.MarketData{
					Exchange:  "KuCoin",
					Symbol:    symbol,
					Bid:       bid,
					Ask:       ask,
					Timestamp: time.Now().UnixMilli(),
				}
			}
		}
	}
}
