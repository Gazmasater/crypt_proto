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
	"strings"
	"sync"
	"time"

	"crypt_proto/pkg/models"

	"github.com/gorilla/websocket"
)

// ------------------- Структуры -------------------

type KuCoinCollector struct {
	ctx      context.Context
	cancel   context.CancelFunc
	symbols  []string
	conn     *websocket.Conn
	wsURL    string
	lastData map[string]struct {
		Bid, Ask, BidSize, AskSize float64
	}
	mu sync.Mutex
}

// ------------------- Чтение CSV -------------------

func readPairsFromCSV(csvFile string) ([][2]string, error) {
	file, err := os.Open(csvFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	r := csv.NewReader(file)
	rows, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(rows) < 2 {
		return nil, nil
	}

	pairsMap := make(map[string][2]string)
	for _, row := range rows[1:] { // пропускаем заголовок
		if len(row) < 6 {
			continue
		}
		for i := 3; i <= 5; i++ { // Leg1..Leg3
			pair := splitLegPair(row[i])
			if pair[0] != "" && pair[1] != "" {
				key := pair[0] + "-" + pair[1] // формат KuCoin BASE-QUOTE
				pairsMap[key] = pair
			}
		}
	}

	pairs := make([][2]string, 0, len(pairsMap))
	for _, p := range pairsMap {
		pairs = append(pairs, p)
	}
	return pairs, nil
}

func splitLegPair(s string) [2]string {
	s = strings.ToUpper(strings.TrimSpace(s))
	fields := strings.Fields(s)
	if len(fields) < 2 {
		return [2]string{}
	}
	parts := strings.Split(fields[1], "/")
	if len(parts) != 2 {
		return [2]string{}
	}
	return [2]string{parts[0], parts[1]}
}

// ------------------- API KuCoin -------------------

type KuCoinSymbol struct {
	Symbol        string `json:"symbol"`
	EnableTrading bool   `json:"enableTrading"`
	Market        string `json:"market"`
}

func FetchKuCoinSymbols() ([]string, error) {
	resp, err := http.Get("https://api.kucoin.com/api/v2/symbols")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Code string         `json:"code"`
		Data []KuCoinSymbol `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if result.Code != "200000" {
		return nil, fmt.Errorf("kucoin API error: %s", result.Code)
	}

	symbols := []string{}
	for _, s := range result.Data {
		if s.EnableTrading && s.Market == "USDS" { // фильтр спота
			symbols = append(symbols, s.Symbol)
		}
	}

	return symbols, nil
}

// ------------------- Фильтрация по треугольникам -------------------

func FilterPairsForTriangles(pairs [][2]string, kucoinSymbols []string) []string {
	existing := make(map[string]struct{}, len(kucoinSymbols))
	for _, s := range kucoinSymbols {
		existing[s] = struct{}{}
	}

	symbols := []string{}
	for _, p := range pairs {
		key := p[0] + "-" + p[1]
		if _, ok := existing[key]; ok {
			symbols = append(symbols, key)
		}
	}
	return symbols
}

// ------------------- Создание коллектора -------------------

func NewKuCoinCollectorFromCSV(csvFile string) (*KuCoinCollector, error) {
	pairs, err := readPairsFromCSV(csvFile)
	if err != nil {
		return nil, err
	}

	allSymbols, err := FetchKuCoinSymbols()
	if err != nil {
		return nil, err
	}

	symbols := FilterPairsForTriangles(pairs, allSymbols)
	if len(symbols) == 0 {
		return nil, fmt.Errorf("no valid symbols to subscribe")
	}

	ctx, cancel := context.WithCancel(context.Background())
	return &KuCoinCollector{
		ctx:     ctx,
		cancel:  cancel,
		symbols: symbols,
		lastData: make(map[string]struct {
			Bid, Ask, BidSize, AskSize float64
		}),
	}, nil
}

func (c *KuCoinCollector) Name() string { return "KuCoin" }

// ------------------- WS -------------------

func (c *KuCoinCollector) Start(out chan<- *models.MarketData) error {
	if err := c.initWS(); err != nil {
		return err
	}

	conn, _, err := websocket.DefaultDialer.Dial(c.wsURL, nil)
	if err != nil {
		return err
	}
	c.conn = conn
	log.Println("[KuCoin] Connected to WS")

	for _, s := range c.symbols {
		sub := map[string]any{
			"id":       time.Now().UnixNano(),
			"type":     "subscribe",
			"topic":    "/market/ticker:" + s,
			"response": true,
		}
		if err := conn.WriteJSON(sub); err != nil {
			return err
		}
		log.Println("[KuCoin] Subscribed:", s)
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

func (c *KuCoinCollector) pingLoop() {
	t := time.NewTicker(15 * time.Second)
	defer t.Stop()
	for {
		select {
		case <-c.ctx.Done():
			return
		case <-t.C:
			_ = c.conn.WriteJSON(map[string]string{"type": "ping"})
		}
	}
}

func (c *KuCoinCollector) readLoop(out chan<- *models.MarketData) {
	defer func() {
		if c.conn != nil {
			_ = c.conn.Close()
		}
	}()
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
			if typ != "message" {
				continue
			}

			topic, _ := raw["topic"].(string)
			data, ok := raw["data"].(map[string]any)
			if !ok {
				continue
			}

			rawSymbol := strings.TrimPrefix(topic, "/market/ticker:")
			symbol := normalizeSymbol(rawSymbol)
			if symbol == "" {
				continue
			}

			bid := parseFloat(data["bestBid"])
			ask := parseFloat(data["bestAsk"])

			if bid == 0 || ask == 0 {
				continue
			}

			c.mu.Lock()
			if last, ok := c.lastData[symbol]; ok {
				if last.Bid == bid && last.Ask == ask {
					c.mu.Unlock()
					continue
				}
			}
			c.lastData[symbol] = struct{ Bid, Ask, BidSize, AskSize float64 }{Bid: bid, Ask: ask}
			c.mu.Unlock()

			out <- &models.MarketData{
				Exchange: "KuCoin",
				Symbol:   symbol,
				Bid:      bid,
				Ask:      ask,
			}
		}
	}
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

func normalizeSymbol(s string) string {
	parts := strings.Split(s, "-")
	if len(parts) != 2 {
		return ""
	}
	return parts[0] + "/" + parts[1]
}

func (c *KuCoinCollector) initWS() error {
	req, err := http.NewRequest("POST", "https://api.kucoin.com/api/v1/bullet-public", nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
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
		return fmt.Errorf("no kucoin ws endpoints")
	}

	c.wsURL = fmt.Sprintf("%s?token=%s&connectId=%d",
		r.Data.InstanceServers[0].Endpoint,
		r.Data.Token,
		time.Now().UnixNano(),
	)
	return nil
}



package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"crypt_proto/collector"
	"crypt_proto/pkg/models"
)

func main() {
	// ------------------- Канал для данных -------------------
	out := make(chan *models.MarketData, 100)

	// ------------------- Создание коллектора -------------------
	kc, err := collector.NewKuCoinCollectorFromCSV("triangles.csv")
	if err != nil {
		log.Fatal("Failed to create KuCoinCollector:", err)
	}

	// ------------------- Запуск коллектора -------------------
	if err := kc.Start(out); err != nil {
		log.Fatal("Failed to start KuCoinCollector:", err)
	}

	log.Println("[Main] KuCoinCollector started. Listening for data...")

	// ------------------- Обработка данных -------------------
	go func() {
		for data := range out {
			log.Printf("[MarketData] %s %s bid=%.6f ask=%.6f",
				data.Exchange, data.Symbol, data.Bid, data.Ask)
		}
	}()

	// ------------------- Завершение при SIGINT / SIGTERM -------------------
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("[Main] Stopping KuCoinCollector...")
	if err := kc.Stop(); err != nil {
		log.Println("Error stopping collector:", err)
	}
	close(out)
	log.Println("[Main] Exited.")
}




