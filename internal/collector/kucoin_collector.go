package collector

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"crypt_proto/internal/market"
	"crypt_proto/pkg/models"

	"github.com/gorilla/websocket"
)

type KuCoinCollector struct {
	ctx      context.Context
	cancel   context.CancelFunc
	conn     *websocket.Conn
	symbols  []string
	allowed  map[string]struct{} // whitelist
	lastData map[string]struct {
		Bid, Ask, BidSize, AskSize float64
	}
	mu    sync.Mutex
	pool  *sync.Pool
	buf   []byte
	wsURL string
}

func NewKuCoinCollector(symbols []string, whitelist []string, pool *sync.Pool) *KuCoinCollector {
	ctx, cancel := context.WithCancel(context.Background())
	allowed := make(map[string]struct{}, len(whitelist))
	for _, s := range whitelist {
		allowed[market.NormalizeSymbol_Full(s)] = struct{}{}
	}
	return &KuCoinCollector{
		ctx:      ctx,
		cancel:   cancel,
		symbols:  symbols,
		allowed:  allowed,
		lastData: make(map[string]struct{ Bid, Ask, BidSize, AskSize float64 }),
		pool:     pool,
		buf:      make([]byte, 0, 32),
	}
}

func (c *KuCoinCollector) Name() string { return "KuCoin" }

func (c *KuCoinCollector) Start(out chan<- *models.MarketData) error {
	if err := c.initWS(); err != nil {
		return err
	}

	conn, _, err := websocket.DefaultDialer.Dial(c.wsURL, nil)
	if err != nil {
		return err
	}
	c.conn = conn
	log.Println("[KuCoin] Connected to WS")

	// subscribe
	for _, s := range c.symbols {
		sym := normalizeKucoinSymbol(s)

		// фильтрация по whitelist
		if len(c.allowed) > 0 {
			if _, ok := c.allowed[sym]; !ok {
				continue
			}
		}

		sub := map[string]any{
			"id":       time.Now().UnixNano(),
			"type":     "subscribe",
			"topic":    "/market/ticker:" + sym,
			"response": true,
		}

		if err := conn.WriteJSON(sub); err != nil {
			return err
		}

		log.Println("[KuCoin] Subscribed:", sym)
	}

	go c.pingLoop()
	go c.readLoop(out)

	return nil
}

func (c *KuCoinCollector) Stop() error {
	c.cancel()
	if c.conn != nil {
		_ = c.conn.Close()
	}
	return nil
}

func (c *KuCoinCollector) pingLoop() {
	t := time.NewTicker(15 * time.Second)
	defer t.Stop()

	for {
		select {
		case <-c.ctx.Done():
			return
		case <-t.C:
			_ = c.conn.WriteJSON(map[string]string{"type": "ping"})
		}
	}
}

func (c *KuCoinCollector) readLoop(out chan<- *models.MarketData) {
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

			var raw map[string]any
			if err := json.Unmarshal(msg, &raw); err != nil {
				continue
			}

			typ, _ := raw["type"].(string)
			if typ == "welcome" || typ == "ack" || typ != "message" {
				continue
			}

			topic, _ := raw["topic"].(string)
			data, ok := raw["data"].(map[string]any)
			if !ok {
				continue
			}

			rawsymbol := strings.TrimPrefix(topic, "/market/ticker:")
			symbol := market.NormalizeSymbol_NoAlloc(rawsymbol, &c.buf)
			if symbol == "" {
				continue
			}

			// фильтрация по whitelist
			if len(c.allowed) > 0 {
				if _, ok := c.allowed[symbol]; !ok {
					continue
				}
			}

			bid := parseFloat(data["bestBid"])
			ask := parseFloat(data["bestAsk"])
			bidSize := parseFloat(data["sizeBid"])
			askSize := parseFloat(data["sizeAsk"])

			if bid == 0 || ask == 0 {
				continue
			}

			// дедупликация
			c.mu.Lock()
			last, exists := c.lastData[symbol]
			if exists && last.Bid == bid && last.Ask == ask && last.BidSize == bidSize && last.AskSize == askSize {
				c.mu.Unlock()
				continue
			}
			c.lastData[symbol] = struct{ Bid, Ask, BidSize, AskSize float64 }{Bid: bid, Ask: ask, BidSize: bidSize, AskSize: askSize}
			c.mu.Unlock()

			// объект из пула
			md := c.pool.Get().(*models.MarketData)
			md.Exchange = "KuCoin"
			md.Symbol = symbol
			md.Bid = bid
			md.Ask = ask
			md.BidSize = bidSize
			md.AskSize = askSize
			md.Timestamp = time.Now().UnixMilli()

			out <- md
		}
	}
}

func normalizeKucoinSymbol(s string) string {
	if strings.Contains(s, "-") {
		return s
	}
	if strings.HasSuffix(s, "USDT") {
		return strings.Replace(s, "USDT", "-USDT", 1)
	}
	return s
}

func parseFloat(v any) float64 {
	switch t := v.(type) {
	case string:
		f, _ := strconv.ParseFloat(t, 64)
		return f
	case float64:
		return t
	default:
		return 0
	}
}
func (c *KuCoinCollector) initWS() error {
	req, err := http.NewRequest(
		"POST",
		"https://api.kucoin.com/api/v1/bullet-public",
		nil,
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("kucoin bullet status: %s", resp.Status)
	}

	var r struct {
		Data struct {
			Token           string `json:"token"`
			InstanceServers []struct {
				Endpoint string `json:"endpoint"`
			} `json:"instanceServers"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return err
	}

	if len(r.Data.InstanceServers) == 0 {
		return fmt.Errorf("no kucoin ws endpoints")
	}

	c.wsURL = fmt.Sprintf(
		"%s?token=%s&connectId=%d",
		r.Data.InstanceServers[0].Endpoint,
		r.Data.Token,
		time.Now().UnixNano(),
	)

	return nil
}
