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
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"crypt_proto/pkg/models"

	"github.com/gorilla/websocket"
)

/* ================= STRUCT ================= */

type KuCoinCollector struct {
	ctx    context.Context
	cancel context.CancelFunc

	conn   *websocket.Conn
	wsURL  string
	symbols []string

	out chan<- *models.MarketData

	last map[string][2]float64
	mu   sync.Mutex

	ready bool
}

/* ================= CONSTRUCTOR ================= */

func NewKuCoinCollectorFromCSV(path string) (*KuCoinCollector, error) {
	symbols, err := readPairsFromCSV(path)
	if err != nil {
		return nil, err
	}
	if len(symbols) == 0 {
		return nil, fmt.Errorf("no symbols")
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &KuCoinCollector{
		ctx:     ctx,
		cancel:  cancel,
		symbols: symbols,
		last:    make(map[string][2]float64),
	}, nil
}

/* ================= INTERFACE ================= */

func (c *KuCoinCollector) Name() string { return "KuCoin" }

func (c *KuCoinCollector) Start(out chan<- *models.MarketData) error {
	c.out = out

	if err := c.initWS(); err != nil {
		return err
	}

	go c.readLoop()
	go c.subscribeBatches(15, 400*time.Millisecond)

	log.Println("[KuCoin] started")
	return nil
}

func (c *KuCoinCollector) Stop() error {
	c.cancel()
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

/* ================= WS INIT ================= */

func (c *KuCoinCollector) initWS() error {
	req, _ := http.NewRequest(
		"POST",
		"https://api.kucoin.com/api/v1/bullet-public",
		nil,
	)

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

	c.wsURL = fmt.Sprintf(
		"%s?token=%s&connectId=%d",
		r.Data.InstanceServers[0].Endpoint,
		r.Data.Token,
		time.Now().UnixNano(),
	)

	conn, _, err := websocket.DefaultDialer.Dial(c.wsURL, nil)
	if err != nil {
		return err
	}

	// deadlines
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	c.conn = conn
	log.Println("[KuCoin] WS connected")
	return nil
}

/* ================= SUBSCRIBE ================= */

func (c *KuCoinCollector) subscribeBatches(batch int, delay time.Duration) {
	// ждём welcome
	for !c.ready {
		select {
		case <-c.ctx.Done():
			return
		default:
			time.Sleep(50 * time.Millisecond)
		}
	}

	log.Println("[KuCoin] subscribing...")

	for i := 0; i < len(c.symbols); i += batch {
		end := i + batch
		if end > len(c.symbols) {
			end = len(c.symbols)
		}

		for _, s := range c.symbols[i:end] {
			_ = c.conn.WriteJSON(map[string]any{
				"id":       time.Now().UnixNano(),
				"type":     "subscribe",
				"topic":    "/market/ticker:" + s,
				"response": true,
			})
		}

		time.Sleep(delay)
	}
}

/* ================= READ LOOP ================= */

func (c *KuCoinCollector) readLoop() {
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
			c.handle(msg)
		}
	}
}

/* ================= HANDLE ================= */

func (c *KuCoinCollector) handle(msg []byte) {
	var raw map[string]any
	if err := json.Unmarshal(msg, &raw); err != nil {
		return
	}

	switch raw["type"] {

	case "welcome":
		c.ready = true
		log.Println("[KuCoin] welcome")
		return

	case "ack":
		return

	case "ping":
		_ = c.conn.WriteJSON(map[string]any{
			"id":   raw["id"],
			"type": "pong",
		})
		return

	case "message":
		// continue
	default:
		return
	}

	data, ok := raw["data"].(map[string]any)
	if !ok {
		return
	}

	bid := parseFloat(data["bestBid"])
	ask := parseFloat(data["bestAsk"])
	if bid == 0 || ask == 0 {
		return
	}

	sym, ok := data["symbol"].(string)
	if !ok {
		return
	}
	symbol := normalize(sym)

	c.mu.Lock()
	last := c.last[symbol]
	if last[0] == bid && last[1] == ask {
		c.mu.Unlock()
		return
	}
	c.last[symbol] = [2]float64{bid, ask}
	c.mu.Unlock()

	c.out <- &models.MarketData{
		Exchange: c.Name(),
		Symbol:   symbol,
		Bid:      bid,
		Ask:      ask,
	}
}

/* ================= CSV ================= */

func readPairsFromCSV(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	rows, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	set := make(map[string]struct{})
	for _, row := range rows[1:] {
		for i := 3; i <= 5 && i < len(row); i++ {
			p := parseLeg(row[i])
			if p != "" {
				set[p] = struct{}{}
			}
		}
	}

	var res []string
	for k := range set {
		res = append(res, k)
	}
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

/* ================= HELPERS ================= */

func normalize(s string) string {
	p := strings.Split(s, "-")
	return p[0] + "/" + p[1]
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



az358@gaz358-BOD-WXX9:~/myprog/crypt_proto/cmd/arb$ go run .
2026/01/05 11:11:05 [KuCoin] WS connected
2026/01/05 11:11:05 [KuCoin] started
2026/01/05 11:11:05 [Main] KuCoinCollector started. Listening for data...
2026/01/05 11:11:05 [KuCoin] welcome
2026/01/05 11:11:05 [KuCoin] subscribing...
2026/01/05 11:11:07 [KuCoin] read error: websocket: close 1000 (normal): Bye



