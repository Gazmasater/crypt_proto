apikey = "4333ed4b-cd83-49f5-97d1-c399e2349748"
secretkey = "E3848531135EDB4CCFDA0F1BC14CD274"
IP = ""
Название API-ключа = "Arb"
Доступы = "Чтение"



sudo systemctl mask sleep.target suspend.target hibernate.target hybrid-sleep.target



wbs-api.mexc.com/ws 


[https://edis-global.vercel.app/ru/vps-hosting/singapore-singapore
](https://sg.edisglobal.com/)



git pull --rebase origin privat
git push origin privat


BOOK_INTERVAL=100ms
SYMBOLS_FILE=triangles_markets.csv
DEBUG=false


import (
    // ...
    "net/http"
    _ "net/http/pprof"
)


   // pprof HTTP-сервер
    go func() {
        log.Println("pprof on http://localhost:6060/debug/pprof/")
        if err := http.ListenAndServe("localhost:6060", nil); err != nil {
            log.Printf("pprof server error: %v", err)
        }
    }()


	go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30


(pprof) top        # показать топ функций по CPU
(pprof) top10
(pprof) list parsePBWrapperMid   # подробный разбор одной функции
(pprof) quit


go tool pprof http://localhost:6060/debug/pprof/heap


(pprof) top
(pprof) top -cum
(pprof) list parsePBWrapperMid
(pprof) quit




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
	ctx    context.Context
	cancel context.CancelFunc

	conn *websocket.Conn
	mu   sync.Mutex

	symbols []string
	allowed map[string]struct{}

	lastData map[string]struct {
		Bid, Ask, BidSize, AskSize float64
	}

	out  chan<- *models.MarketData
	pool *sync.Pool
	buf  []byte
}

// ------------------------------------------------------------

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
		lastData: make(map[string]struct{}),
		pool:     pool,
		buf:      make([]byte, 0, 32),
	}
}

func (c *MEXCCollector) Name() string { return "MEXC" }

// ------------------------------------------------------------

func (c *MEXCCollector) Start(out chan<- *models.MarketData) error {
	c.out = out
	return c.connect()
}

func (c *MEXCCollector) Stop() error {
	c.cancel()
	c.mu.Lock()
	if c.conn != nil {
		_ = c.conn.Close()
	}
	c.mu.Unlock()
	return nil
}

// ------------------------------------------------------------
// connection lifecycle
// ------------------------------------------------------------

func (c *MEXCCollector) connect() error {
	conn, _, err := websocket.DefaultDialer.Dial(configs.MEXC_WS, nil)
	if err != nil {
		return err
	}

	c.mu.Lock()
	c.conn = conn
	c.mu.Unlock()

	log.Println("[MEXC] connected")

	conn.SetReadDeadline(time.Now().Add(configs.MEXC_READ_TIMEOUT))
	conn.SetPongHandler(func(string) error {
		return conn.SetReadDeadline(time.Now().Add(configs.MEXC_READ_TIMEOUT))
	})

	if err := c.subscribeChunks(25); err != nil {
		return err
	}

	go c.pingLoop()
	go c.readLoop()

	return nil
}

func (c *MEXCCollector) reconnect() {
	c.mu.Lock()
	if c.conn != nil {
		_ = c.conn.Close()
		c.conn = nil
	}
	c.mu.Unlock()

	log.Println("[MEXC] reconnecting...")

	for {
		select {
		case <-c.ctx.Done():
			return
		default:
		}

		if err := c.connect(); err != nil {
			time.Sleep(time.Second)
			continue
		}

		log.Println("[MEXC] reconnected")
		return
	}
}

// ------------------------------------------------------------
// loops
// ------------------------------------------------------------

func (c *MEXCCollector) pingLoop() {
	t := time.NewTicker(configs.MEXC_PING_INTERVAL)
	defer t.Stop()

	for {
		select {
		case <-c.ctx.Done():
			return
		case <-t.C:
			c.mu.Lock()
			err := c.conn.WriteMessage(websocket.PingMessage, nil)
			c.mu.Unlock()

			if err != nil {
				log.Println("[MEXC] ping error:", err)
				c.reconnect()
				return
			}
		}
	}
}

func (c *MEXCCollector) readLoop() {
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
		}

		_, raw, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("[MEXC] read error:", err)
			c.reconnect()
			return
		}

		var wrap pb.PushDataV3ApiWrapper
		if err := proto.Unmarshal(raw, &wrap); err != nil {
			continue
		}

		if md := c.handleWrapper(&wrap); md != nil {
			c.out <- md
		}
	}
}

// ------------------------------------------------------------
// parsing
// ------------------------------------------------------------

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

	if _, ok := c.allowed[symbol]; !ok {
		return nil
	}

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

	md := c.pool.Get().(*models.MarketData)
	md.Exchange = "MEXC"
	md.Symbol = symbol
	md.Bid = bid
	md.Ask = ask
	md.BidSize = bidSize
	md.AskSize = askSize
	md.Timestamp = time.Now().UnixMilli()

	return md
}

// ------------------------------------------------------------
// helpers
// ------------------------------------------------------------

func (c *MEXCCollector) subscribeChunks(size int) error {
	for size < len(c.symbols) {
		if err := c.subscribe(c.symbols[:size]); err != nil {
			return err
		}
		c.symbols = c.symbols[size:]
	}
	return c.subscribe(c.symbols)
}

func (c *MEXCCollector) subscribe(chunk []string) error {
	params := make([]string, 0, len(chunk))
	for _, s := range chunk {
		params = append(params, "spot@public.aggre.bookTicker.v3.api.pb@100ms@"+s)
	}

	return c.conn.WriteJSON(map[string]any{
		"method": "SUBSCRIPTION",
		"params": params,
	})
}

func fastParseFloat(s string) float64 {
	var intPart, fracPart int64
	div := 1.0
	i := 0

	for i < len(s) && s[i] != '.' {
		intPart = intPart*10 + int64(s[i]-'0')
		i++
	}

	if i < len(s) && s[i] == '.' {
		i++
		for i < len(s) {
			fracPart = fracPart*10 + int64(s[i]-'0')
			div *= 10
			i++
		}
	}

	return float64(intPart) + float64(fracPart)/div
}

