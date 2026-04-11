package collector

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"crypt_proto/pkg/models"

	"github.com/gorilla/websocket"
	"github.com/tidwall/gjson"
)

const (
	maxSubsPerWS   = 126
	subRate        = 120 * time.Millisecond
	pingInterval   = 20 * time.Second
	reconnectDelay = 3 * time.Second
	readTimeout    = 45 * time.Second
)

type KuCoinCollector struct {
	ctx    context.Context
	cancel context.CancelFunc
	wsList []*kucoinWS
	out    chan<- *models.MarketData
}

type Last struct {
	Bid     float64
	Ask     float64
	BidSize float64
	AskSize float64
}

type kucoinWS struct {
	id      int
	conn    *websocket.Conn
	symbols []string
	last    map[string]Last
	writeMu sync.Mutex
}

func NewKuCoinCollectorFromCSV(path string) (*KuCoinCollector, []string, error) {
	symbols, err := readPairsFromCSV(path)
	if err != nil {
		return nil, nil, err
	}
	if len(symbols) == 0 {
		return nil, nil, fmt.Errorf("no symbols")
	}

	ctx, cancel := context.WithCancel(context.Background())

	var wsList []*kucoinWS
	for i := 0; i < len(symbols); i += maxSubsPerWS {
		end := i + maxSubsPerWS
		if end > len(symbols) {
			end = len(symbols)
		}

		wsList = append(wsList, &kucoinWS{
			id:      len(wsList),
			symbols: symbols[i:end],
			last:    make(map[string]Last),
		})
	}

	c := &KuCoinCollector{ctx: ctx, cancel: cancel, wsList: wsList}
	return c, symbols, nil
}

func (c *KuCoinCollector) Name() string { return "KuCoin" }

func (c *KuCoinCollector) Start(out chan<- *models.MarketData) error {
	c.out = out
	for _, ws := range c.wsList {
		go ws.run(c)
	}
	log.Printf("[KuCoin] started with %d WS\n", len(c.wsList))
	return nil
}

func (c *KuCoinCollector) Stop() error {
	c.cancel()
	for _, ws := range c.wsList {
		ws.closeConn()
	}
	return nil
}

func (ws *kucoinWS) run(c *KuCoinCollector) {
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
		}

		if err := ws.connect(); err != nil {
			log.Printf("[KuCoin WS %d] connect error: %v, retry in %v\n", ws.id, err, reconnectDelay)
			time.Sleep(reconnectDelay)
			continue
		}

		connDone := make(chan struct{})
		go ws.subscribeLoop(c.ctx, connDone)
		go ws.pingLoop(c.ctx, connDone)

		if err := ws.readLoop(c); err != nil {
			log.Printf("[KuCoin WS %d] readLoop ended: %v\n", ws.id, err)
		}

		close(connDone)
		ws.closeConn()
		log.Printf("[KuCoin WS %d] reconnecting in %v...\n", ws.id, reconnectDelay)
		time.Sleep(reconnectDelay)
	}
}

func (ws *kucoinWS) connect() error {
	req, _ := http.NewRequest("POST", "https://api.kucoin.com/api/v1/bullet-public", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

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
		return fmt.Errorf("no instance servers")
	}

	url := fmt.Sprintf("%s?token=%s&connectId=%d", r.Data.InstanceServers[0].Endpoint, r.Data.Token, time.Now().UnixNano())
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return err
	}

	_ = conn.SetReadDeadline(time.Now().Add(readTimeout))
	conn.SetPongHandler(func(string) error {
		return conn.SetReadDeadline(time.Now().Add(readTimeout))
	})

	ws.conn = conn
	log.Printf("[KuCoin WS %d] connected\n", ws.id)
	return nil
}

func (ws *kucoinWS) subscribeLoop(ctx context.Context, connDone <-chan struct{}) {
	t := time.NewTicker(subRate)
	defer t.Stop()
	for _, s := range ws.symbols {
		select {
		case <-ctx.Done():
			return
		case <-connDone:
			return
		case <-t.C:
			if err := ws.writeJSON(map[string]any{
				"id":       time.Now().UnixNano(),
				"type":     "subscribe",
				"topic":    "/market/ticker:" + s,
				"response": true,
			}); err != nil {
				return
			}
		}
	}
}

func (ws *kucoinWS) pingLoop(ctx context.Context, connDone <-chan struct{}) {
	t := time.NewTicker(pingInterval)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-connDone:
			return
		case <-t.C:
			if err := ws.writeJSON(map[string]any{"id": time.Now().UnixNano(), "type": "ping"}); err != nil {
				return
			}
		}
	}
}

func (ws *kucoinWS) readLoop(c *KuCoinCollector) error {
	for {
		if ws.conn == nil {
			return fmt.Errorf("connection closed")
		}
		_, msg, err := ws.conn.ReadMessage()
		if err != nil {
			return err
		}
		_ = ws.conn.SetReadDeadline(time.Now().Add(readTimeout))
		ws.handle(c, msg)
	}
}

func (ws *kucoinWS) handle(c *KuCoinCollector, msg []byte) {
	const prefix = "/market/ticker:"
	const prefixLen = len(prefix)

	topicRes := gjson.GetBytes(msg, "topic")
	if !topicRes.Exists() {
		return
	}

	raw := topicRes.Raw
	if len(raw) <= prefixLen+2 {
		return
	}
	if raw[1:1+prefixLen] != prefix {
		return
	}

	symbol := raw[1+prefixLen : len(raw)-1]
	data := gjson.GetBytes(msg, "data")
	bid := data.Get("bestBid").Float()
	ask := data.Get("bestAsk").Float()
	if bid == 0 || ask == 0 {
		return
	}
	bidSize := data.Get("bestBidSize").Float()
	askSize := data.Get("bestAskSize").Float()
	if bidSize == 0 || askSize == 0 {
		return
	}

	if last, ok := ws.last[symbol]; ok && last.Bid == bid && last.Ask == ask && last.BidSize == bidSize && last.AskSize == askSize {
		return
	}
	ws.last[symbol] = Last{Bid: bid, Ask: ask, BidSize: bidSize, AskSize: askSize}

	md := &models.MarketData{
		Exchange:  "KuCoin",
		Symbol:    symbol,
		Bid:       bid,
		Ask:       ask,
		BidSize:   bidSize,
		AskSize:   askSize,
		Timestamp: time.Now().UnixMilli(),
	}

	select {
	case c.out <- md:
	case <-c.ctx.Done():
	}
}

func (ws *kucoinWS) writeJSON(v any) error {
	ws.writeMu.Lock()
	defer ws.writeMu.Unlock()
	if ws.conn == nil {
		return fmt.Errorf("connection closed")
	}
	return ws.conn.WriteJSON(v)
}

func (ws *kucoinWS) closeConn() {
	ws.writeMu.Lock()
	defer ws.writeMu.Unlock()
	if ws.conn != nil {
		_ = ws.conn.Close()
		ws.conn = nil
	}
}

func readPairsFromCSV(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rows, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return nil, err
	}
	if len(rows) < 2 {
		return nil, nil
	}

	header := make(map[string]int, len(rows[0]))
	for i, col := range rows[0] {
		header[strings.TrimSpace(col)] = i
	}

	set := make(map[string]struct{})
	for _, row := range rows[1:] {
		if len(strings.TrimSpace(strings.Join(row, ""))) == 0 {
			continue
		}

		foundInRow := false
		for _, key := range []string{"Leg1Symbol", "Leg2Symbol", "Leg3Symbol"} {
			if idx, ok := header[key]; ok && idx < len(row) {
				symbol := strings.ToUpper(strings.TrimSpace(row[idx]))
				if symbol != "" {
					set[symbol] = struct{}{}
					foundInRow = true
				}
			}
		}

		if !foundInRow {
			for _, key := range []string{"Leg1", "Leg2", "Leg3"} {
				if idx, ok := header[key]; ok && idx < len(row) {
					if p := parseLeg(row[idx]); p != "" {
						set[p] = struct{}{}
					}
				}
			}
		}
	}

	res := make([]string, 0, len(set))
	for k := range set {
		res = append(res, k)
	}
	sort.Strings(res)
	return res, nil
}

func parseLeg(s string) string {
	parts := strings.Fields(strings.ToUpper(strings.TrimSpace(s)))
	if len(parts) < 2 {
		return ""
	}
	p := strings.Split(parts[1], "/")
	if len(p) != 2 {
		return ""
	}
	return p[0] + "-" + p[1]
}
