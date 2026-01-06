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




csv.go
package collector

import (
	"encoding/csv"
	"os"
	"strings"
)

// LoadSymbolsFromCSV читает треугольники и возвращает уникальные пары
func LoadSymbolsFromCSV(path string) ([]string, error) {
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

	// пропускаем header
	for _, row := range rows[1:] {
		// Leg1, Leg2, Leg3 = 3,4,5
		for i := 3; i <= 5 && i < len(row); i++ {
			symbol := parseLeg(row[i])
			if symbol != "" {
				set[symbol] = struct{}{}
			}
		}
	}

	var res []string
	for s := range set {
		res = append(res, s)
	}
	return res, nil
}

func parseLeg(s string) string {
	parts := strings.Fields(strings.TrimSpace(s))
	if len(parts) < 2 {
		return ""
	}

	pair := strings.Split(parts[1], "/")
	if len(pair) != 2 {
		return ""
	}

	return strings.ToUpper(pair[0] + "-" + pair[1])
}



kucoin_pool.go

package collector

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"crypt_proto/pkg/models"

	"github.com/gorilla/websocket"
)

const (
	maxSubsPerWS   = 50
	subscribeDelay = 100 * time.Millisecond
	pingInterval   = 20 * time.Second
)

/* ================= POOL ================= */

type KuCoinCollector struct {
	ctx    context.Context
	cancel context.CancelFunc

	workers []*kucoinWS
	out     chan<- *models.MarketData
}

func NewKuCoinCollector(symbols []string) *KuCoinCollector {
	ctx, cancel := context.WithCancel(context.Background())

	var workers []*kucoinWS
	chunks := chunk(symbols, maxSubsPerWS)

	for i, c := range chunks {
		workers = append(workers, newKucoinWS(ctx, i, c))
	}

	return &KuCoinCollector{
		ctx:     ctx,
		cancel:  cancel,
		workers: workers,
	}
}

func (c *KuCoinCollector) Name() string { return "KuCoin" }

func (c *KuCoinCollector) Start(out chan<- *models.MarketData) error {
	c.out = out

	for _, w := range c.workers {
		w.out = out
		go w.run()
	}

	log.Printf("[KuCoin] started (%d ws connections)", len(c.workers))
	return nil
}

func (c *KuCoinCollector) Stop() error {
	c.cancel()
	return nil
}






kucoin_ws.go

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

	"crypt_proto/pkg/models"

	"github.com/gorilla/websocket"
)

type kucoinWS struct {
	ctx     context.Context
	id      int
	symbols []string
	out     chan<- *models.MarketData

	conn  *websocket.Conn
	last  map[string][2]float64
	mutex sync.Mutex
}

func newKucoinWS(ctx context.Context, id int, symbols []string) *kucoinWS {
	return &kucoinWS{
		ctx:     ctx,
		id:      id,
		symbols: symbols,
		last:    make(map[string][2]float64),
	}
}

func (w *kucoinWS) run() {
	if err := w.connect(); err != nil {
		log.Println("[KuCoin WS]", w.id, "connect error:", err)
		return
	}

	go w.pingLoop()
	go w.readLoop()

	w.subscribe()
}

func (w *kucoinWS) connect() error {
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
			Token string `json:"token"`
			InstanceServers []struct {
				Endpoint string `json:"endpoint"`
			} `json:"instanceServers"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return err
	}

	wsURL := fmt.Sprintf(
		"%s?token=%s&connectId=%d",
		r.Data.InstanceServers[0].Endpoint,
		r.Data.Token,
		time.Now().UnixNano(),
	)

	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		return err
	}

	w.conn = conn
	log.Printf("[KuCoin WS %d] connected (%d symbols)", w.id, len(w.symbols))
	return nil
}

func (w *kucoinWS) subscribe() {
	for _, s := range w.symbols {
		select {
		case <-w.ctx.Done():
			return
		default:
		}

		topic := "/market/level2Depth5:" + s

		_ = w.conn.WriteJSON(map[string]any{
			"id":       time.Now().UnixNano(),
			"type":     "subscribe",
			"topic":    topic,
			"response": true,
		})

		log.Printf("[KuCoin WS %d] subscribed %s", w.id, s)
		time.Sleep(subscribeDelay)
	}
}

func (w *kucoinWS) readLoop() {
	for {
		select {
		case <-w.ctx.Done():
			return
		default:
		}

		_, msg, err := w.conn.ReadMessage()
		if err != nil {
			log.Printf("[KuCoin WS %d] read error: %v", w.id, err)
			return
		}

		w.handle(msg)
	}
}

func (w *kucoinWS) pingLoop() {
	t := time.NewTicker(pingInterval)
	defer t.Stop()

	for {
		select {
		case <-w.ctx.Done():
			return
		case <-t.C:
			_ = w.conn.WriteJSON(map[string]any{
				"id":   time.Now().UnixNano(),
				"type": "ping",
			})
		}
	}
}

func (w *kucoinWS) handle(msg []byte) {
	var raw map[string]any
	if err := json.Unmarshal(msg, &raw); err != nil {
		return
	}

	if raw["type"] != "message" {
		return
	}

	data, ok := raw["data"].(map[string]any)
	if !ok {
		return
	}

	bids := data["bids"].([]any)
	asks := data["asks"].([]any)
	if len(bids) == 0 || len(asks) == 0 {
		return
	}

	bid := parseFloat(bids[0].([]any)[0])
	ask := parseFloat(asks[0].([]any)[0])
	symbol := normalize(data["symbol"].(string))

	w.mutex.Lock()
	last := w.last[symbol]
	if last[0] == bid && last[1] == ask {
		w.mutex.Unlock()
		return
	}
	w.last[symbol] = [2]float64{bid, ask}
	w.mutex.Unlock()

	w.out <- &models.MarketData{
		Exchange: "KuCoin",
		Symbol:   symbol,
		Bid:      bid,
		Ask:      ask,
	}
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




chunk.go

package collector

func chunk[T any](arr []T, size int) [][]T {
	var res [][]T
	for size < len(arr) {
		arr, res = arr[size:], append(res, arr[0:size:size])
	}
	res = append(res, arr)
	return res
}






package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"crypt_proto/internal/collector"
	"crypt_proto/pkg/models"
)

func main() {

	// 1️⃣ Загружаем символы из CSV
	symbols, err := collector.LoadSymbolsFromCSV(
		"../exchange/data/kucoin_triangles_usdt.csv",
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("loaded %d symbols from csv", len(symbols))

	// 2️⃣ Создаём коллектор (ОДНО значение!)
	kc := collector.NewKuCoinCollector(symbols)

	out := make(chan *models.MarketData, 1000)

	// 3️⃣ Запуск
	go func() {
		if err := kc.Start(out); err != nil {
			log.Fatal(err)
		}
	}()

	// 4️⃣ Чтение котировок
	go func() {
		for md := range out {
			log.Printf(
				"[KUCOIN] %s bid=%.8f ask=%.8f",
				md.Symbol, md.Bid, md.Ask,
			)
		}
	}()

	// 5️⃣ Graceful shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	log.Println("shutdown...")
	_ = kc.Stop()
}









