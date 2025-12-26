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
	"strconv"
	"strings"
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
		lastData: make(map[string]struct{ Bid, Ask, BidSize, AskSize float64 }),
	}
}

func (c *MEXCCollector) Name() string {
	return "MEXC"
}

func (c *MEXCCollector) Start(out chan<- models.MarketData) error {
	conn, _, err := websocket.DefaultDialer.Dial(configs.MEXC_WS, nil)
	if err != nil {
		return err
	}
	c.conn = conn
	log.Println("[MEXC] connected")

	if err := c.subscribe(); err != nil {
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

// --- приватные методы ---

func (c *MEXCCollector) subscribe() error {
	params := make([]string, 0, len(c.symbols))
	for _, s := range c.symbols {
		params = append(params, "spot@public.aggre.bookTicker.v3.api.pb@100ms@"+s)
	}

	sub := map[string]interface{}{
		"method": "SUBSCRIPTION",
		"params": params,
	}

	return c.conn.WriteJSON(sub)
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

func (c *MEXCCollector) readLoop(out chan<- models.MarketData) {
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
			log.Printf("[MEXC] text: %s\n", raw)
			continue
		}
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

func (c *MEXCCollector) handleWrapper(wrap *pb.PushDataV3ApiWrapper) *models.MarketData {
	switch body := wrap.GetBody().(type) {
	case *pb.PushDataV3ApiWrapper_PublicAggreBookTicker:
		bt := body.PublicAggreBookTicker

		bid, _ := strconv.ParseFloat(bt.GetBidPrice(), 64)
		ask, _ := strconv.ParseFloat(bt.GetAskPrice(), 64)
		bidSize, _ := strconv.ParseFloat(bt.GetBidQty(), 64)
		askSize, _ := strconv.ParseFloat(bt.GetAskQty(), 64)

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

		// фильтрация повторов
		if last, exists := c.lastData[symbol]; exists {
			if last.Bid == bid && last.Ask == ask && last.BidSize == bidSize && last.AskSize == askSize {
				return nil
			}
		}

		// обновляем последние данные
		c.lastData[symbol] = struct {
			Bid, Ask, BidSize, AskSize float64
		}{Bid: bid, Ask: ask, BidSize: bidSize, AskSize: askSize}

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
	default:
		return nil
	}
}


[{
	"resource": "/home/gaz358/myprog/crypt_proto/internal/collector/mexc_collector.go",
	"owner": "_generated_diagnostic_collection_name_#0",
	"code": {
		"value": "MissingLitField",
		"target": {
			"$mid": 1,
			"path": "/golang.org/x/tools/internal/typesinternal",
			"scheme": "https",
			"authority": "pkg.go.dev",
			"fragment": "MissingLitField"
		}
	},
	"severity": 8,
	"message": "unknown field BidSize in struct literal of type models.MarketData",
	"source": "compiler",
	"startLineNumber": 188,
	"startColumn": 4,
	"endLineNumber": 188,
	"endColumn": 11,
	"origin": "extHost1"
}]


[{
	"resource": "/home/gaz358/myprog/crypt_proto/internal/collector/mexc_collector.go",
	"owner": "_generated_diagnostic_collection_name_#0",
	"code": {
		"value": "MissingLitField",
		"target": {
			"$mid": 1,
			"path": "/golang.org/x/tools/internal/typesinternal",
			"scheme": "https",
			"authority": "pkg.go.dev",
			"fragment": "MissingLitField"
		}
	},
	"severity": 8,
	"message": "unknown field AskSize in struct literal of type models.MarketData",
	"source": "compiler",
	"startLineNumber": 189,
	"startColumn": 4,
	"endLineNumber": 189,
	"endColumn": 11,
	"origin": "extHost1"
}]


