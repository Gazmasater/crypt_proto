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
	ctx    context.Context
	cancel context.CancelFunc

	conn         *websocket.Conn
	wsURL        string
	pingInterval time.Duration

	symbols []string
	allowed map[string]struct{}

	last map[string]lastTick
	mu   sync.Mutex

	pool *sync.Pool
	buf  []byte
}

type lastTick struct {
	Bid, Ask, BidSize, AskSize float64
}

// ------------------------------------------------------------
// Constructor
// ------------------------------------------------------------
func NewKuCoinCollector(
	symbols []string,
	whitelist []string,
	pool *sync.Pool,
) *KuCoinCollector {

	ctx, cancel := context.WithCancel(context.Background())

	allowed := make(map[string]struct{}, len(whitelist))
	for _, s := range whitelist {
		allowed[market.NormalizeSymbol_Full(s)] = struct{}{}
	}

	return &KuCoinCollector{
		ctx:     ctx,
		cancel:  cancel,
		symbols: symbols,
		allowed: allowed,
		last:    make(map[string]lastTick),
		pool:    pool,
		buf:     make([]byte, 0, 32),
	}
}

func (c *KuCoinCollector) Name() string { return "KuCoin" }

// ------------------------------------------------------------
// Start / Stop
// ------------------------------------------------------------
func (c *KuCoinCollector) Start(out chan<- *models.MarketData) error {
	log.Println("[KuCoin] init WS")

	if err := c.initWS(); err != nil {
		return err
	}

	log.Println("[KuCoin] connect:", c.wsURL)
	conn, _, err := websocket.DefaultDialer.Dial(c.wsURL, nil)
	if err != nil {
		return err
	}
	c.conn = conn
	log.Println("[KuCoin] connected")

	// 🔴 ОБЯЗАТЕЛЬНО ждём welcome
	if err := c.waitWelcome(); err != nil {
		return err
	}
	log.Println("[KuCoin] welcome received")

	if err := c.subscribeBatch(20); err != nil {
		return err
	}

	if c.pingInterval > 0 {
		go c.pingLoop()
		log.Printf("[KuCoin] pingLoop started (%s)\n", c.pingInterval)
	}

	go c.readLoop(out)
	return nil
}

func (c *KuCoinCollector) Stop() error {
	log.Println("[KuCoin] stopping")
	c.cancel()
	if c.conn != nil {
		_ = c.conn.Close()
	}
	return nil
}

// ------------------------------------------------------------
// WS Init
// ------------------------------------------------------------
func (c *KuCoinCollector) initWS() error {
	req, err := http.NewRequest(
		"POST",
		"https://api.kucoin.com/api/v1/bullet-public",
		nil,
	)
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bullet status: %s", resp.Status)
	}

	var r struct {
		Data struct {
			Token           string `json:"token"`
			InstanceServers []struct {
				Endpoint     string `json:"endpoint"`
				PingInterval int    `json:"pingInterval"`
			} `json:"instanceServers"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return err
	}

	if len(r.Data.InstanceServers) == 0 {
		return fmt.Errorf("no ws endpoints")
	}

	s := r.Data.InstanceServers[0]
	c.wsURL = fmt.Sprintf(
		"%s?token=%s&connectId=%d",
		s.Endpoint,
		r.Data.Token,
		time.Now().UnixNano(),
	)

	c.pingInterval = time.Duration(s.PingInterval) * time.Millisecond
	log.Printf("[KuCoin] wsURL ready, pingInterval=%s\n", c.pingInterval)
	return nil
}

// ------------------------------------------------------------
// Welcome
// ------------------------------------------------------------
func (c *KuCoinCollector) waitWelcome() error {
	_ = c.conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	defer c.conn.SetReadDeadline(time.Time{})

	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			return err
		}

		var raw map[string]any
		if err := json.Unmarshal(msg, &raw); err != nil {
			continue
		}

		if raw["type"] == "welcome" {
			return nil
		}
	}
}

// ------------------------------------------------------------
// Subscribe
// ------------------------------------------------------------
func (c *KuCoinCollector) subscribeBatch(batchSize int) error {
	total := len(c.symbols)
	sent := 0

	for i := 0; i < total; i += batchSize {
		end := i + batchSize
		if end > total {
			end = total
		}

		for _, s := range c.symbols[i:end] {
			sym := normalizeKucoinSymbol(s)
			norm := market.NormalizeSymbol_Full(sym)

			if len(c.allowed) > 0 {
				if _, ok := c.allowed[norm]; !ok {
					continue
				}
			}

			msg := map[string]any{
				"id":       fmt.Sprintf("sub-%s", sym),
				"type":     "subscribe",
				"topic":    "/market/ticker:" + sym,
				"response": true,
			}

			if err := c.conn.WriteJSON(msg); err != nil {
				return err
			}

			sent++
		}

		log.Printf("[KuCoin] subscribed %d / %d\n", sent, total)
		time.Sleep(300 * time.Millisecond) // 🔴 важно
	}

	log.Printf("[KuCoin] subscribed TOTAL: %d symbols\n", sent)
	return nil
}

// ------------------------------------------------------------
// Ping
// ------------------------------------------------------------
func (c *KuCoinCollector) pingLoop() {
	ticker := time.NewTicker(c.pingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-c.ctx.Done():
			return
		case <-ticker.C:
			_ = c.conn.WriteJSON(map[string]string{"type": "ping"})
		}
	}
}

// ------------------------------------------------------------
// Read
// ------------------------------------------------------------
func (c *KuCoinCollector) readLoop(out chan<- *models.MarketData) {
	defer func() {
		log.Println("[KuCoin] readLoop stopped")
		_ = c.conn.Close()
	}()

	log.Println("[KuCoin] readLoop started")

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

			if raw["type"] != "message" {
				continue
			}

			data, ok := raw["data"].(map[string]any)
			if !ok {
				continue
			}

			topic, _ := raw["topic"].(string)
			rawSym := strings.TrimPrefix(topic, "/market/ticker:")
			symbol := market.NormalizeSymbol_NoAlloc(rawSym, &c.buf)
			if symbol == "" {
				continue
			}

			bid := parseFloat(data["bestBid"])
			ask := parseFloat(data["bestAsk"])
			if bid == 0 || ask == 0 {
				continue
			}

			bidSize := parseFloat(data["sizeBid"])
			askSize := parseFloat(data["sizeAsk"])

			c.mu.Lock()
			prev, ok := c.last[symbol]
			if ok &&
				prev.Bid == bid &&
				prev.Ask == ask &&
				prev.BidSize == bidSize &&
				prev.AskSize == askSize {
				c.mu.Unlock()
				continue
			}
			c.last[symbol] = lastTick{bid, ask, bidSize, askSize}
			c.mu.Unlock()

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

// ------------------------------------------------------------
// Helpers
// ------------------------------------------------------------
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



2026/01/04 16:06:36 EXCHANGE: kucoin
2026/01/04 16:06:36 pprof on http://localhost:6060/debug/pprof/
2026/01/04 16:06:36 Loaded 246 unique symbols from ../exchange/data/kucoin_triangles_usdt.csv
2026/01/04 16:06:36 [KuCoin] init WS
2026/01/04 16:06:37 [KuCoin] wsURL ready, pingInterval=18s
2026/01/04 16:06:37 [KuCoin] connect: wss://ws-api-spot.kucoin.com/?token=2neAiuYvAU61ZDXANAGAsiL4-iAExhsBXZxftpOeh_55i3Ysy2q2LEsEWU64mdzUOPusi34M_wGoSf7iNyEWJ3qC8zu5Yz5AZ3_oRjxFoa6_k1pAuYXVudiYB9J6i9GjsxUuhPw3Blq6rhZlGykT3Vp1phUafnulOOpts-MEmEHtGZR-Jl-TQoCqCnec6Y_sJBvJHl5Vs9Y=.2SsaV2nQPU-90Qa8eu9N_g==&connectId=1767531997763402774
2026/01/04 16:06:38 [KuCoin] connected
2026/01/04 16:06:38 [KuCoin] welcome received
2026/01/04 16:06:38 [KuCoin] subscribed 20 / 246
2026/01/04 16:06:39 [KuCoin] subscribed 40 / 246
2026/01/04 16:06:39 [KuCoin] subscribed 60 / 246
2026/01/04 16:06:39 [KuCoin] subscribed 80 / 246
2026/01/04 16:06:40 [KuCoin] subscribed 100 / 246
2026/01/04 16:06:40 [KuCoin] subscribed 120 / 246
2026/01/04 16:06:40 [KuCoin] subscribed 140 / 246
2026/01/04 16:06:40 start collector:write tcp 192.168.1.71:44250->108.157.229.57:443: write: connection reset by peer
exit status 1







