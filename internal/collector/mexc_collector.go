package collector

import (
	"context"
	"log"
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

type MEXCCollector struct {
	ctx      context.Context
	cancel   context.CancelFunc
	conn     *websocket.Conn
	symbols  []string
	allowed  map[string]struct{} // whitelist
	lastData map[string]struct {
		Bid, Ask, BidSize, AskSize float64
	}
	mu  sync.Mutex
	buf []byte // отдельный буфер для NormalizeSymbol

	pool *sync.Pool
}

// Конструктор с пулом
func NewMEXCCollector(symbols []string, whitelist []string, pool *sync.Pool) *MEXCCollector {
	ctx, cancel := context.WithCancel(context.Background())
	allowed := make(map[string]struct{}, len(whitelist))
	for _, s := range whitelist {
		allowed[market.NormalizeSymbol_Full(s)] = struct{}{}
	}
	return &MEXCCollector{
		ctx:      ctx,
		cancel:   cancel,
		symbols:  symbols,
		allowed:  allowed,
		lastData: make(map[string]struct{ Bid, Ask, BidSize, AskSize float64 }),
		pool:     pool,
		buf:      make([]byte, 0, 32),
	}
}

func (c *MEXCCollector) Name() string {
	return "MEXC"
}

func (c *MEXCCollector) Start(out chan<- *models.MarketData) error {
	conn, _, err := websocket.DefaultDialer.Dial(configs.MEXC_WS, nil)
	if err != nil {
		return err
	}
	c.conn = conn
	log.Println("[MEXC] connected")

	if err := c.subscribeChunks(25); err != nil {
		return err
	}

	go c.pingLoop()
	go c.readLoop(out)

	return nil
}

func (c *MEXCCollector) Stop() error {
	c.cancel()
	if c.conn != nil {
		_ = c.conn.Close()
	}
	return nil
}

// ----------------- приватные методы -----------------

func (c *MEXCCollector) subscribeChunks(chunkSize int) error {
	chunks := chunkSymbols(c.symbols, chunkSize)
	for _, chunk := range chunks {
		params := make([]string, 0, len(chunk))
		for _, s := range chunk {
			params = append(params, "spot@public.aggre.bookTicker.v3.api.pb@100ms@"+s)
		}
		sub := map[string]interface{}{
			"method": "SUBSCRIPTION",
			"params": params,
		}
		if err := c.conn.WriteJSON(sub); err != nil {
			return err
		}
	}
	return nil
}

func (c *MEXCCollector) pingLoop() {
	t := time.NewTicker(configs.MEXC_PING_INTERVAL)
	defer t.Stop()

	for {
		select {
		case <-c.ctx.Done():
			return
		case <-t.C:
			_ = c.conn.WriteMessage(websocket.PingMessage, []byte("hb"))
		}
	}
}

func (c *MEXCCollector) readLoop(out chan<- *models.MarketData) {
	_ = c.conn.SetReadDeadline(time.Now().Add(configs.MEXC_READ_TIMEOUT))

	for {
		select {
		case <-c.ctx.Done():
			return
		default:
		}

		mt, raw, err := c.conn.ReadMessage()
		if err != nil {
			log.Printf("[MEXC] read error: %v\n", err)
			return
		}

		_ = c.conn.SetReadDeadline(time.Now().Add(configs.MEXC_READ_TIMEOUT))

		if mt == websocket.TextMessage {
			continue // игнорируем текст
		}
		if mt != websocket.BinaryMessage {
			continue
		}

		var wrap pb.PushDataV3ApiWrapper
		if err := proto.Unmarshal(raw, &wrap); err != nil {
			continue
		}

		if md := c.handleWrapper(&wrap); md != nil {
			out <- md
		}
	}
}

func (c *MEXCCollector) handleWrapper(wrap *pb.PushDataV3ApiWrapper) *models.MarketData {
	body, ok := wrap.GetBody().(*pb.PushDataV3ApiWrapper_PublicAggreBookTicker)
	if !ok {
		return nil
	}
	bt := body.PublicAggreBookTicker

	bid := fastParseFloat(bt.GetBidPrice())
	ask := fastParseFloat(bt.GetAskPrice())
	bidSize := fastParseFloat(bt.GetBidQuantity())
	askSize := fastParseFloat(bt.GetAskQuantity())

	symbol := wrap.GetSymbol()
	if symbol == "" {
		ch := wrap.GetChannel()
		if ch != "" {
			parts := strings.Split(ch, "@")
			symbol = parts[len(parts)-1]
		}
	}

	symbol = market.NormalizeSymbol_NoAlloc(symbol, &c.buf)
	if symbol == "" {
		return nil
	}

	// фильтрация по whitelist
	if len(c.allowed) > 0 {
		if _, ok := c.allowed[symbol]; !ok {
			return nil
		}
	}

	// дедупликация struct
	c.mu.Lock()
	last, ok := c.lastData[symbol]
	if ok && last.Bid == bid && last.Ask == ask && last.BidSize == bidSize && last.AskSize == askSize {
		c.mu.Unlock()
		return nil
	}
	c.lastData[symbol] = struct {
		Bid, Ask, BidSize, AskSize float64
	}{Bid: bid, Ask: ask, BidSize: bidSize, AskSize: askSize}
	c.mu.Unlock()

	ts := time.Now().UnixMilli()
	if t := wrap.GetSendTime(); t > 0 {
		ts = t
	}

	// объект из пула
	md := c.pool.Get().(*models.MarketData)
	md.Exchange = "MEXC"
	md.Symbol = symbol
	md.Bid = bid
	md.Ask = ask
	md.BidSize = bidSize
	md.AskSize = askSize
	md.Timestamp = ts

	return md
}

func fastParseFloat(s string) float64 {
	// очень простой парсер для формата "123.456"
	i := 0
	var intPart int64
	for ; i < len(s) && s[i] != '.'; i++ {
		intPart = intPart*10 + int64(s[i]-'0')
	}

	var fracPart int64
	fracDiv := 1.0
	if i < len(s) && s[i] == '.' {
		i++
		for ; i < len(s); i++ {
			fracPart = fracPart*10 + int64(s[i]-'0')
			fracDiv *= 10
		}
	}

	return float64(intPart) + float64(fracPart)/fracDiv
}

// разбивает слайс на чанки
func chunkSymbols(src []string, size int) [][]string {
	var out [][]string
	for size < len(src) {
		out = append(out, src[:size:size])
		src = src[size:]
	}
	out = append(out, src)
	return out
}
