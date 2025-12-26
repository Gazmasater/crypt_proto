package collector

import (
	"context"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"crypt_proto/configs"
	"crypt_proto/internal/market"
	pb "crypt_proto/pb"
	"crypt_proto/pkg/models"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

const mexcChunkSize = 25

type MEXCCollector struct {
	ctx    context.Context
	cancel context.CancelFunc

	symbols []string

	conns []*websocket.Conn
	mu    sync.Mutex

	lastData map[string]struct {
		Bid, Ask, BidSize, AskSize float64
	}
}

func NewMEXCCollector(symbols []string) *MEXCCollector {
	ctx, cancel := context.WithCancel(context.Background())

	return &MEXCCollector{
		ctx:      ctx,
		cancel:   cancel,
		symbols:  symbols,
		conns:    make([]*websocket.Conn, 0),
		lastData: make(map[string]struct{ Bid, Ask, BidSize, AskSize float64 }),
	}
}

func (c *MEXCCollector) Name() string {
	return "MEXC"
}

// ------------------------------------------------------------
// START
// ------------------------------------------------------------

func (c *MEXCCollector) Start(out chan<- models.MarketData) error {
	chunks := chunkSymbols(c.symbols, mexcChunkSize)

	for i, chunk := range chunks {
		conn, _, err := websocket.DefaultDialer.Dial(configs.MEXC_WS, nil)
		if err != nil {
			return err
		}

		log.Printf("[MEXC] connected chunk %d/%d (%d symbols)", i+1, len(chunks), len(chunk))

		c.mu.Lock()
		c.conns = append(c.conns, conn)
		c.mu.Unlock()

		if err := c.subscribe(conn, chunk); err != nil {
			return err
		}

		go c.pingLoop(conn)
		go c.readLoop(conn, out)
	}

	return nil
}

// ------------------------------------------------------------
// STOP
// ------------------------------------------------------------

func (c *MEXCCollector) Stop() error {
	c.cancel()

	c.mu.Lock()
	defer c.mu.Unlock()

	for _, conn := range c.conns {
		_ = conn.Close()
	}

	return nil
}

// ------------------------------------------------------------
// SUBSCRIBE
// ------------------------------------------------------------

func (c *MEXCCollector) subscribe(conn *websocket.Conn, symbols []string) error {
	params := make([]string, 0, len(symbols))

	for _, s := range symbols {
		params = append(params, "spot@public.aggre.bookTicker.v3.api.pb@100ms@"+s)
	}

	req := map[string]any{
		"method": "SUBSCRIPTION",
		"params": params,
	}

	return conn.WriteJSON(req)
}

// ------------------------------------------------------------
// PING
// ------------------------------------------------------------

func (c *MEXCCollector) pingLoop(conn *websocket.Conn) {
	t := time.NewTicker(configs.MEXC_PING_INTERVAL)
	defer t.Stop()

	for {
		select {
		case <-c.ctx.Done():
			return
		case <-t.C:
			_ = conn.WriteMessage(websocket.PingMessage, []byte("hb"))
		}
	}
}

// ------------------------------------------------------------
// READ LOOP
// ------------------------------------------------------------

func (c *MEXCCollector) readLoop(conn *websocket.Conn, out chan<- models.MarketData) {
	_ = conn.SetReadDeadline(time.Now().Add(configs.MEXC_READ_TIMEOUT))

	for {
		select {
		case <-c.ctx.Done():
			return
		default:
		}

		mt, raw, err := conn.ReadMessage()
		if err != nil {
			log.Printf("[MEXC] read error: %v", err)
			return
		}

		_ = conn.SetReadDeadline(time.Now().Add(configs.MEXC_READ_TIMEOUT))

		if mt != websocket.BinaryMessage {
			continue
		}

		var wrap pb.PushDataV3ApiWrapper
		if err := proto.Unmarshal(raw, &wrap); err != nil {
			continue
		}

		md := c.handleWrapper(&wrap)
		if md != nil {
			out <- *md
		}
	}
}

// ------------------------------------------------------------
// MESSAGE HANDLER
// ------------------------------------------------------------

func (c *MEXCCollector) handleWrapper(wrap *pb.PushDataV3ApiWrapper) *models.MarketData {
	switch body := wrap.GetBody().(type) {

	case *pb.PushDataV3ApiWrapper_PublicAggreBookTicker:
		bt := body.PublicAggreBookTicker

		bid, _ := strconv.ParseFloat(bt.GetBidPrice(), 64)
		ask, _ := strconv.ParseFloat(bt.GetAskPrice(), 64)
		bidSize, _ := strconv.ParseFloat(bt.GetBidQuantity(), 64)
		askSize, _ := strconv.ParseFloat(bt.GetAskQuantity(), 64)

		symbol := wrap.GetSymbol()
		if symbol == "" {
			ch := wrap.GetChannel()
			if ch != "" {
				parts := strings.Split(ch, "@")
				symbol = parts[len(parts)-1]
			}
		}

		symbol = market.NormalizeSymbol_Full(symbol)
		if symbol == "" {
			return nil
		}

		// дедупликация
		c.mu.Lock()
		last, ok := c.lastData[symbol]
		if ok &&
			last.Bid == bid &&
			last.Ask == ask &&
			last.BidSize == bidSize &&
			last.AskSize == askSize {
			c.mu.Unlock()
			return nil
		}

		c.lastData[symbol] = struct {
			Bid, Ask, BidSize, AskSize float64
		}{bid, ask, bidSize, askSize}
		c.mu.Unlock()

		ts := time.Now().UnixMilli()
		if t := wrap.GetSendTime(); t > 0 {
			ts = t
		}

		return &models.MarketData{
			Exchange:  "MEXC",
			Symbol:    symbol,
			Bid:       bid,
			Ask:       ask,
			BidSize:   bidSize,
			AskSize:   askSize,
			Timestamp: ts,
		}
	}

	return nil
}

// ------------------------------------------------------------
// HELPERS
// ------------------------------------------------------------

func chunkSymbols(src []string, size int) [][]string {
	if size <= 0 {
		return nil
	}

	var out [][]string
	for len(src) > size {
		out = append(out, src[:size])
		src = src[size:]
	}

	if len(src) > 0 {
		out = append(out, src)
	}

	return out
}
