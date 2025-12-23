package collector

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"

	pb "crypt_proto/pb"
)

const (
	mexcWS       = "wss://wbs-api.mexc.com/ws"
	pingInterval = 30 * time.Second
	readTimeout  = 60 * time.Second
)

type MEXCCollector struct {
	ctx     context.Context
	cancel  context.CancelFunc
	symbols []string
	conn    *websocket.Conn
}

func NewMEXCCollector(symbols []string) *MEXCCollector {
	ctx, cancel := context.WithCancel(context.Background())
	return &MEXCCollector{
		ctx:     ctx,
		cancel:  cancel,
		symbols: symbols,
	}
}

func (c *MEXCCollector) Start() error {
	conn, _, err := websocket.DefaultDialer.Dial(mexcWS, nil)
	if err != nil {
		return err
	}
	c.conn = conn
	log.Println("[MEXC] connected")

	if err := c.subscribe(); err != nil {
		return err
	}

	go c.pingLoop()
	go c.readLoop()

	return nil
}

func (c *MEXCCollector) Stop() {
	c.cancel()
	if c.conn != nil {
		_ = c.conn.Close()
	}
}

func (c *MEXCCollector) subscribe() error {
	params := make([]string, 0, len(c.symbols))
	for _, s := range c.symbols {
		s = strings.ToUpper(s)
		params = append(params,
			"spot@public.aggre.bookTicker.v3.api.pb@100ms@"+s,
		)
	}

	sub := map[string]any{
		"method": "SUBSCRIPTION",
		"params": params,
	}

	if err := c.conn.WriteJSON(sub); err != nil {
		return err
	}

	b, _ := json.Marshal(sub)
	log.Printf("[MEXC] subscribed: %s\n", b)
	return nil
}

func (c *MEXCCollector) pingLoop() {
	t := time.NewTicker(pingInterval)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			_ = c.conn.WriteMessage(websocket.PingMessage, []byte("hb"))
		case <-c.ctx.Done():
			return
		}
	}
}

func (c *MEXCCollector) readLoop() {
	_ = c.conn.SetReadDeadline(time.Now().Add(readTimeout))

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

		_ = c.conn.SetReadDeadline(time.Now().Add(readTimeout))

		// ACK / ошибки
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

		c.handleWrapper(&wrap)
	}
}

func (c *MEXCCollector) handleWrapper(w *pb.PushDataV3ApiWrapper) {
	symbol := w.GetSymbol()
	if symbol == "" {
		ch := w.GetChannel()
		if ch != "" {
			parts := strings.Split(ch, "@")
			symbol = parts[len(parts)-1]
		}
	}

	ts := time.Now()
	if t := w.GetSendTime(); t > 0 {
		ts = time.UnixMilli(t)
	}

	switch body := w.GetBody().(type) {

	case *pb.PushDataV3ApiWrapper_PublicAggreBookTicker:
		bt := body.PublicAggreBookTicker

		bid, _ := strconv.ParseFloat(bt.GetBidPrice(), 64)
		ask, _ := strconv.ParseFloat(bt.GetAskPrice(), 64)
		bq, _ := strconv.ParseFloat(bt.GetBidQuantity(), 64)
		aq, _ := strconv.ParseFloat(bt.GetAskQuantity(), 64)

		log.Printf(
			"[MEXC] %s bid=%.8f(%.6f) ask=%.8f(%.6f) ts=%s",
			symbol, bid, bq, ask, aq, ts.Format(time.RFC3339Nano),
		)
	}
}
