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

const (
	kucoinBatchSize = 50
)

type KuCoinCollector struct {
	ctx    context.Context
	cancel context.CancelFunc

	conn  *websocket.Conn
	wsURL string

	pingInterval time.Duration

	symbols []string

	last map[string]lastTick
	mu   sync.Mutex

	pool *sync.Pool
	buf  []byte
}

type lastTick struct {
	Bid, Ask, BidSize, AskSize float64
}

func NewKuCoinCollector(symbols []string, pool *sync.Pool) *KuCoinCollector {
	ctx, cancel := context.WithCancel(context.Background())

	return &KuCoinCollector{
		ctx:     ctx,
		cancel:  cancel,
		symbols: symbols,
		last:    make(map[string]lastTick),
		pool:    pool,
		buf:     make([]byte, 0, 32),
	}
}

func (c *KuCoinCollector) Name() string { return "KuCoin" }

// -------------------- START / STOP --------------------

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

	if err := c.subscribe(); err != nil {
		return err
	}

	go c.pingLoop()
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

// -------------------- SUBSCRIBE --------------------

func (c *KuCoinCollector) subscribe() error {
	var batch []string

	for _, s := range c.symbols {
		batch = append(batch, normalizeKucoinSymbol(s))

		if len(batch) == kucoinBatchSize {
			if err := c.subscribeBatch(batch); err != nil {
				return err
			}
			batch = batch[:0]
			time.Sleep(300 * time.Millisecond)
		}
	}

	if len(batch) > 0 {
		return c.subscribeBatch(batch)
	}

	return nil
}

func (c *KuCoinCollector) subscribeBatch(symbols []string) error {
	for _, sym := range symbols {
		msg := map[string]any{
			"id":       fmt.Sprintf("sub-%s", sym),
			"type":     "subscribe",
			"topic":    "/market/ticker:" + sym,
			"response": true,
		}

		if err := c.conn.WriteJSON(msg); err != nil {
			return err
		}
	}

	log.Printf("[KuCoin] subscribed batch: %d symbols\n", len(symbols))
	return nil
}

// -------------------- PING --------------------

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

// -------------------- READ LOOP --------------------

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

			typ, _ := raw["type"].(string)

			switch typ {
			case "welcome", "ack", "pong":
				continue
			case "message":
			default:
				continue
			}

			topic, _ := raw["topic"].(string)
			data, ok := raw["data"].(map[string]any)
			if !ok {
				continue
			}

			rawSym := strings.TrimPrefix(topic, "/market/ticker:")
			symbol := market.NormalizeSymbol_NoAlloc(rawSym, &c.buf)
			if symbol == "" {
				continue
			}

			bid := parseFloat(data["bestBid"])
			ask := parseFloat(data["bestAsk"])
			bidSize := parseFloat(data["sizeBid"])
			askSize := parseFloat(data["sizeAsk"])

			if bid == 0 || ask == 0 {
				continue
			}

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

// -------------------- INIT WS --------------------

func (c *KuCoinCollector) initWS() error {
	log.Println("[KuCoin] request bullet-public")

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
			Token        string `json:"token"`
			PingInterval int64  `json:"pingInterval"`
			PingTimeout  int64  `json:"pingTimeout"`
			Instance     []struct {
				Endpoint string `json:"endpoint"`
			} `json:"instanceServers"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return err
	}

	if len(r.Data.Instance) == 0 {
		return fmt.Errorf("no ws endpoints")
	}

	c.wsURL = fmt.Sprintf(
		"%s?token=%s&connectId=%d",
		r.Data.Instance[0].Endpoint,
		r.Data.Token,
		time.Now().UnixNano(),
	)

	c.pingInterval = time.Duration(r.Data.PingInterval) * time.Millisecond

	log.Printf("[KuCoin] wsURL ready, pingInterval=%v\n", c.pingInterval)
	return nil
}

// -------------------- HELPERS --------------------

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


2026/01/04 11:25:34 [KuCoin] connect: wss://ws-api-spot.kucoin.com/?token=2neAiuYvAU61ZDXANAGAsiL4-iAExhsBXZxftpOeh_55i3Ysy2q2LEsEWU64mdzUOPusi34M_wGoSf7iNyEWJzP9evvOrxtv30MhfOvDMxlwLvwvfH4J2diYB9J6i9GjsxUuhPw3Blq6rhZlGykT3Vp1phUafnulOOpts-MEmEHtGZR-Jl-TQtzCQxbLF0kVJBvJHl5Vs9Y=.nmAnl7bXIz7cIP4_r_Dbew==&connectId=1767515134453591035
2026/01/04 11:25:35 [KuCoin] connected
2026/01/04 11:25:35 [KuCoin] subscribed batch: 50 symbols
2026/01/04 11:25:35 [KuCoin] subscribed batch: 50 symbols
2026/01/04 11:25:36 [KuCoin] subscribed batch: 50 symbols
2026/01/04 11:25:36 [KuCoin] subscribed batch: 50 symbols
2026/01/04 11:25:36 start collector:write tcp 192.168.1.71:52536->108.157.229.57:443: write: broken pipe
exit status 1
gaz358@gaz358-BOD-WXX9:~/myprog/crypt_proto/cmd/arb$ 




