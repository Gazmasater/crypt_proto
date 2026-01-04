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

import "crypt_proto/pkg/models"

type Collector interface {
	Start(out chan<- *models.MarketData) error
	Stop() error
	Name() string
}



package collector

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"crypt_proto/internal/market"
	"crypt_proto/pkg/models"

	"github.com/gorilla/websocket"
)

type KuCoinCollector struct {
	ctx      context.Context
	cancel   context.CancelFunc
	symbols  []string
	conn     *websocket.Conn
	wsURL    string
	lastData map[string]struct {
		Bid, Ask, BidSize, AskSize float64
	}
}

func NewKuCoinCollector(symbols []string) *KuCoinCollector {
	ctx, cancel := context.WithCancel(context.Background())
	return &KuCoinCollector{
		ctx:      ctx,
		cancel:   cancel,
		symbols:  symbols,
		lastData: make(map[string]struct{ Bid, Ask, BidSize, AskSize float64 }),
	}
}

func (c *KuCoinCollector) Name() string { return "KuCoin" }

func (c *KuCoinCollector) Start(out chan<- models.MarketData) error {
	if err := c.initWS(); err != nil {
		return err
	}

	conn, _, err := websocket.DefaultDialer.Dial(c.wsURL, nil)
	if err != nil {
		return err
	}
	c.conn = conn
	log.Println("[KuCoin] Connected to WS")

	// subscribe
	for _, s := range c.symbols {
		sym := normalizeKucoinSymbol(s)

		sub := map[string]any{
			"id":       time.Now().UnixNano(),
			"type":     "subscribe",
			"topic":    "/market/ticker:" + sym,
			"response": true,
		}

		if err := conn.WriteJSON(sub); err != nil {
			return err
		}

		log.Println("[KuCoin] Subscribed:", sym)
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

func (c *KuCoinCollector) readLoop(out chan<- models.MarketData) {
	defer func() {
		if c.conn != nil {
			c.conn.Close()
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
			if typ == "welcome" || typ == "ack" || typ != "message" {
				continue
			}

			topic, _ := raw["topic"].(string)
			data, ok := raw["data"].(map[string]any)
			if !ok {
				continue
			}

			rawsymbol := strings.TrimPrefix(topic, "/market/ticker:")
			symbol := market.NormalizeSymbol_Full(rawsymbol)
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

			// фильтрация повторов
			if last, exists := c.lastData[symbol]; exists {
				if last.Bid == bid && last.Ask == ask && last.BidSize == bidSize && last.AskSize == askSize {
					continue // ничего не поменялось — пропускаем
				}
			}

			// обновляем последние данные
			c.lastData[symbol] = struct {
				Bid, Ask, BidSize, AskSize float64
			}{Bid: bid, Ask: ask, BidSize: bidSize, AskSize: askSize}

			out <- models.MarketData{
				Exchange: "KuCoin",
				Symbol:   symbol,
				Bid:      bid,
				Ask:      ask,
				// при желании можно добавить объёмы
			}
		}
	}
}

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

func (c *KuCoinCollector) initWS() error {
	req, err := http.NewRequest(
		"POST",
		"https://api.kucoin.com/api/v1/bullet-public",
		nil,
	)
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

	if resp.StatusCode != 200 {
		return fmt.Errorf("kucoin bullet status: %s", resp.Status)
	}

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

	c.wsURL = fmt.Sprintf(
		"%s?token=%s&connectId=%d",
		r.Data.InstanceServers[0].Endpoint,
		r.Data.Token,
		time.Now().UnixNano(),
	)

	return nil
}



package main

import (
	"encoding/csv"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"crypt_proto/internal/collector"
	"crypt_proto/pkg/models"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load(".env")

	// pprof
	go func() {
		log.Println("pprof on http://localhost:6060/debug/pprof/")
		if err := http.ListenAndServe("localhost:6060", nil); err != nil {
			log.Printf("pprof server error: %v", err)
		}
	}()

	exchange := strings.ToLower(os.Getenv("EXCHANGE"))
	if exchange == "" {
		log.Fatal("Set EXCHANGE env variable: mexc | okx | kucoin")
	}
	log.Println("EXCHANGE:", exchange)

	// канал маркет-данных
	marketDataCh := make(chan *models.MarketData, 1000)

	// пул MarketData
	marketDataPool := &sync.Pool{
		New: func() interface{} {
			return new(models.MarketData)
		},
	}

	// === читаем whitelist из CSV ===
	csvPath := "../exchange/data/kucoin_triangles_usdt.csv"

	symbols, err := readSymbolsFromCSV(csvPath, exchange)
	if err != nil {
		log.Fatalf("read CSV symbols: %v", err)
	}
	log.Printf("Loaded %d unique symbols from %s", len(symbols), csvPath)

	// создаём whitelist
	whitelist := make([]string, len(symbols))
	copy(whitelist, symbols)

	var c collector.Collector

	// создаём collector в зависимости от биржи
	switch exchange {
	case "mexc":
		c = collector.NewMEXCCollector(symbols, whitelist, marketDataPool)
	case "okx":
		//
		c = collector.NewOKXCollector(symbols, whitelist, marketDataPool)
	case "kucoin":
		c = collector.NewKuCoinCollector([]string{"BTCUSDT", "ETHUSDT", "ETHBTC", "ETHMANA", " MANAUSDT", " KLVUSDT"})
	default:
		log.Fatal("Unknown exchange:", exchange)
	}

	// старт collector
	if err := c.Start(marketDataCh); err != nil {
		log.Fatal("start collector:", err)
	}

	// consumer маркет-данных
	go func() {
		for md := range marketDataCh {
			log.Printf("[%s] %s bid=%.8f ask=%.8f",
				md.Exchange, md.Symbol, md.Bid, md.Ask,
			)
			// возвращаем объект обратно в пул
			marketDataPool.Put(md)
		}
	}()

	// graceful shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	log.Println("Stopping collector...")
	if err := c.Stop(); err != nil {
		log.Println("Stop error:", err)
	}
}

// ------------------------------------------------------------
// CSV → symbols, нормализуем под биржу
// ------------------------------------------------------------
func readSymbolsFromCSV(path, exchange string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)

	// читаем заголовок
	header, err := r.Read()
	if err != nil {
		return nil, err
	}

	// ищем колонки Leg1, Leg2, Leg3
	colIndex := make(map[string]int)
	for i, h := range header {
		colIndex[strings.ToLower(strings.TrimSpace(h))] = i
	}

	required := []string{"leg1", "leg2", "leg3"}
	var idx []int
	for _, name := range required {
		i, ok := colIndex[strings.ToLower(name)]
		if !ok {
			return nil, csv.ErrFieldCount
		}
		idx = append(idx, i)
	}

	// множество уникальных символов
	uniq := make(map[string]struct{})

	for {
		row, err := r.Read()
		if err != nil {
			break
		}

		for _, i := range idx {
			if i >= len(row) {
				continue
			}

			raw := strings.TrimSpace(row[i])
			if raw == "" {
				continue
			}

			// raw = "BUY PEPE/USDT" → вытаскиваем символ
			parts := strings.Fields(raw)
			if len(parts) < 2 {
				continue
			}
			symbol := parts[1] // "PEPE/USDT"

			// нормализуем под биржу
			switch exchange {
			case "mexc":
				symbol = strings.ReplaceAll(symbol, "/", "") // PEPEUSDT
			case "okx", "kucoin":
				symbol = strings.ReplaceAll(symbol, "/", "-") // PEPE-USDT
			default:
				// оставляем как есть
			}

			uniq[symbol] = struct{}{}
		}
	}

	// формируем срез
	out := make([]string, 0, len(uniq))
	for s := range uniq {
		out = append(out, s)
	}

	return out, nil
}



[{
	"resource": "/home/gaz358/myprog/crypt_proto/cmd/arb/main.go",
	"owner": "_generated_diagnostic_collection_name_#0",
	"code": {
		"value": "InvalidIfaceAssign",
		"target": {
			"$mid": 1,
			"path": "/golang.org/x/tools/internal/typesinternal",
			"scheme": "https",
			"authority": "pkg.go.dev",
			"fragment": "InvalidIfaceAssign"
		}
	},
	"severity": 8,
	"message": "cannot use collector.NewKuCoinCollector([]string{…}) (value of type *collector.KuCoinCollector) as collector.Collector value in assignment: *collector.KuCoinCollector does not implement collector.Collector (wrong type for method Start)\n\t\thave Start(chan<- models.MarketData) error\n\t\twant Start(chan<- *models.MarketData) error",
	"source": "compiler",
	"startLineNumber": 70,
	"startColumn": 7,
	"endLineNumber": 70,
	"endColumn": 113,
	"origin": "extHost1"
}]



