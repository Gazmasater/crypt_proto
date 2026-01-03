package collector

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"sync"
	"time"

	"crypt_proto/configs"
	"crypt_proto/internal/market"
	"crypt_proto/pkg/models"

	"github.com/gorilla/websocket"
)

type OKXCollector struct {
	ctx      context.Context
	cancel   context.CancelFunc
	conn     *websocket.Conn
	symbols  []string
	allowed  map[string]struct{} // whitelist
	lastData map[string]struct {
		Bid, Ask, BidSize, AskSize float64
	}
	mu   sync.Mutex
	pool *sync.Pool
	buf  []byte // для NormalizeSymbol_NoAlloc
}

// Конструктор с whitelist
func NewOKXCollector(symbols []string, whitelist []string, pool *sync.Pool) *OKXCollector {
	ctx, cancel := context.WithCancel(context.Background())
	allowed := make(map[string]struct{}, len(whitelist))
	for _, s := range whitelist {
		allowed[market.NormalizeSymbol_Full(s)] = struct{}{}
	}
	return &OKXCollector{
		ctx:      ctx,
		cancel:   cancel,
		symbols:  symbols,
		allowed:  allowed,
		lastData: make(map[string]struct{ Bid, Ask, BidSize, AskSize float64 }),
		pool:     pool,
		buf:      make([]byte, 0, 32),
	}
}

func (c *OKXCollector) Name() string { return "OKX" }

func (c *OKXCollector) Start(out chan<- *models.MarketData) error {
	conn, _, err := websocket.DefaultDialer.Dial(configs.OKX_WS, nil)
	if err != nil {
		return err
	}
	c.conn = conn
	log.Println("[OKX] connected")

	if err := c.subscribe(); err != nil {
		return err
	}

	go c.pingLoop()
	go c.readLoop(out)

	return nil
}

func (c *OKXCollector) Stop() error {
	c.cancel()
	if c.conn != nil {
		_ = c.conn.Close()
	}
	return nil
}

// ----------------- приватные методы -----------------

func (c *OKXCollector) subscribe() error {
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
	return c.conn.WriteJSON(sub)
}

func (c *OKXCollector) pingLoop() {
	t := time.NewTicker(configs.OKX_PING_INTERVAL)
	defer t.Stop()

	for {
		select {
		case <-c.ctx.Done():
			return
		case <-t.C:
			_ = c.conn.WriteMessage(websocket.PingMessage, nil)
		}
	}
}

func (c *OKXCollector) readLoop(out chan<- *models.MarketData) {
	_ = c.conn.SetReadDeadline(time.Now().Add(configs.OKX_READ_TIMEOUT))

	for {
		select {
		case <-c.ctx.Done():
			return
		default:
		}

		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("[OKX] read error:", err)
			return
		}

		_ = c.conn.SetReadDeadline(time.Now().Add(configs.OKX_READ_TIMEOUT))

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

		md := c.handleData(resp.Arg.InstID, resp.Data[0])
		if md != nil {
			out <- md
		}
	}
}

func (c *OKXCollector) handleData(instID string, data struct {
	Asks [][]string `json:"asks"`
	Bids [][]string `json:"bids"`
}) *models.MarketData {

	bid, ask, bidSize, askSize := 0.0, 0.0, 0.0, 0.0

	if len(data.Bids) > 0 {
		bid, _ = strconv.ParseFloat(data.Bids[0][0], 64)
		bidSize, _ = strconv.ParseFloat(data.Bids[0][1], 64)
	}
	if len(data.Asks) > 0 {
		ask, _ = strconv.ParseFloat(data.Asks[0][0], 64)
		askSize, _ = strconv.ParseFloat(data.Asks[0][1], 64)
	}

	symbol := market.NormalizeSymbol_NoAlloc(instID, &c.buf)
	if symbol == "" || bid == 0 || ask == 0 {
		return nil
	}

	// фильтрация по whitelist
	if len(c.allowed) > 0 {
		if _, ok := c.allowed[symbol]; !ok {
			return nil
		}
	}

	// дедупликация
	c.mu.Lock()
	last, ok := c.lastData[symbol]
	if ok && last.Bid == bid && last.Ask == ask && last.BidSize == bidSize && last.AskSize == askSize {
		c.mu.Unlock()
		return nil
	}
	c.lastData[symbol] = struct{ Bid, Ask, BidSize, AskSize float64 }{
		Bid: bid, Ask: ask, BidSize: bidSize, AskSize: askSize,
	}
	c.mu.Unlock()

	md := c.pool.Get().(*models.MarketData)
	md.Exchange = "OKX"
	md.Symbol = symbol
	md.Bid = bid
	md.Ask = ask
	md.BidSize = bidSize
	md.AskSize = askSize
	md.Timestamp = time.Now().UnixMilli()

	return md
}
