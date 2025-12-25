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
	ctx     context.Context
	cancel  context.CancelFunc
	conn    *websocket.Conn
	symbols []string
}

func NewMEXCCollector(symbols []string) *MEXCCollector {
	ctx, cancel := context.WithCancel(context.Background())
	return &MEXCCollector{
		ctx:     ctx,
		cancel:  cancel,
		symbols: symbols,
	}
}

// Имя биржи
func (c *MEXCCollector) Name() string {
	return "MEXC"
}

// Запуск коллектора
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

// Остановка
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

// Преобразуем protobuf в MarketData
func (c *MEXCCollector) handleWrapper(wrap *pb.PushDataV3ApiWrapper) *models.MarketData {
	switch body := wrap.GetBody().(type) {
	case *pb.PushDataV3ApiWrapper_PublicAggreBookTicker:
		bt := body.PublicAggreBookTicker
		bid, _ := strconv.ParseFloat(bt.GetBidPrice(), 64)
		ask, _ := strconv.ParseFloat(bt.GetAskPrice(), 64)

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

		ts := time.Now().UnixMilli()
		if t := wrap.GetSendTime(); t > 0 {
			ts = t
		}

		return &models.MarketData{
			Exchange:  "MEXC",
			Symbol:    symbol,
			Bid:       bid,
			Ask:       ask,
			Timestamp: ts,
		}
	default:
		return nil
	}
}

// простой split (чтобы не импортировать strings, если уже не использовано)
//func split(s, sep string) []string {
//	return []string{}
//}
