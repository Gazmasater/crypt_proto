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

type KuCoinCollector struct {
	ctx    context.Context
	cancel context.CancelFunc

	conn   *websocket.Conn
	wsURL  string
	symbols []string

	last map[string][2]float64
	mu   sync.Mutex
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
		for i := 3; i <= 5; i++ {
			if i >= len(row) {
				continue
			}
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

/* ================= INIT ================= */

func NewKuCoinCollectorFromCSV(csv string) (*KuCoinCollector, error) {
	symbols, err := readPairsFromCSV(csv)
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

func (c *KuCoinCollector) Name() string { return "KuCoin" }

/* ================= START ================= */

func (c *KuCoinCollector) Start(out chan<- *models.MarketData) error {
	if err := c.initWS(); err != nil {
		return err
	}

	conn, _, err := websocket.DefaultDialer.Dial(c.wsURL, nil)
	if err != nil {
		return err
	}

	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	c.conn = conn
	log.Println("[KuCoin] WS connected")

	if err := c.waitWelcome(); err != nil {
		return err
	}

	for _, s := range c.symbols {
		c.subscribe(s)
	}

	go c.pingLoop()
	go c.readLoop(out)

	return nil
}

func (c *KuCoinCollector) Stop() error {
	c.cancel()
	if c.conn != nil {
		_ = c.conn.Close()
	}
	return nil
}

/* ================= WS ================= */

func (c *KuCoinCollector) initWS() error {
	req, _ := http.NewRequest("POST", "https://api.kucoin.com/api/v1/bullet-public", nil)
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var r struct {
		Data struct {
			Token string `json:"token"`
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

	return nil
}

func (c *KuCoinCollector) waitWelcome() error {
	_, msg, err := c.conn.ReadMessage()
	if err != nil {
		return err
	}

	var m map[string]any
	if err := json.Unmarshal(msg, &m); err != nil {
		return err
	}

	if m["type"] != "welcome" {
		return fmt.Errorf("no welcome from kucoin")
	}

	return nil
}

func (c *KuCoinCollector) subscribe(symbol string) {
	_ = c.conn.WriteJSON(map[string]any{
		"id":       time.Now().UnixNano(),
		"type":     "subscribe",
		"topic":    "/market/ticker:" + symbol,
		"response": true,
	})
	log.Println("[KuCoin] Subscribed:", symbol)
}

/* ================= LOOPS ================= */

func (c *KuCoinCollector) pingLoop() {
	t := time.NewTicker(10 * time.Second)
	defer t.Stop()

	for {
		select {
		case <-c.ctx.Done():
			return
		case <-t.C:
			_ = c.conn.WriteJSON(map[string]any{
				"id":   time.Now().UnixNano(),
				"type": "ping",
			})
		}
	}
}

func (c *KuCoinCollector) readLoop(out chan<- *models.MarketData) {
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

			switch raw["type"] {
			case "pong", "ack":
				continue
			case "message":
				c.handleTicker(raw, out)
			}
		}
	}
}

/* ================= DATA ================= */

func (c *KuCoinCollector) handleTicker(raw map[string]any, out chan<- *models.MarketData) {
	data, ok := raw["data"].(map[string]any)
	if !ok {
		return
	}

	bid := parseFloat(data["bestBid"])
	ask := parseFloat(data["bestAsk"])
	if bid == 0 || ask == 0 {
		return
	}

	symbol := normalize(strings.TrimPrefix(raw["topic"].(string), "/market/ticker:"))

	c.mu.Lock()
	last := c.last[symbol]
	if last[0] == bid && last[1] == ask {
		c.mu.Unlock()
		return
	}
	c.last[symbol] = [2]float64{bid, ask}
	c.mu.Unlock()

	out <- &models.MarketData{
		Exchange: "KuCoin",
		Symbol:   symbol,
		Bid:      bid,
		Ask:      ask,
	}
}

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


026/01/05 03:35:06 [KuCoin] Subscribed: RSR-BTC
2026/01/05 03:35:06 [KuCoin] Subscribed: KNC-USDT
2026/01/05 03:35:06 [KuCoin] Subscribed: AAVE-USDT
2026/01/05 03:35:06 [KuCoin] Subscribed: USDT-BRL
2026/01/05 03:35:06 [KuCoin] Subscribed: XDC-ETH
2026/01/05 03:35:06 [KuCoin] Subscribed: MOVR-ETH
2026/01/05 03:35:06 [Main] KuCoinCollector started. Listening for data...
2026/01/05 03:35:07 [KuCoin] read error: websocket: close 1006 (abnormal closure): unexpected EOF




