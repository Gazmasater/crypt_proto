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
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"crypt_proto/configs"
	"crypt_proto/internal/market"
	pb "crypt_proto/pb"
	"crypt_proto/pkg/models"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

type MEXCCollector struct {
	ctx      context.Context
	cancel   context.CancelFunc
	conn     *websocket.Conn
	symbols  []string
	allowed  map[string]struct{} // whitelist
	lastData map[string]struct {
		Bid, Ask, BidSize, AskSize float64
	}
	mu sync.Mutex
}

func NewMEXCCollector(symbols []string, whitelist []string) *MEXCCollector {
	ctx, cancel := context.WithCancel(context.Background())
	allowed := make(map[string]struct{}, len(whitelist))
	for _, s := range whitelist {
		allowed[market.NormalizeSymbol_Full(s)] = struct{}{}
	}
	return &MEXCCollector{
		ctx:      ctx,
		cancel:   cancel,
		symbols:  symbols,
		allowed:  allowed,
		lastData: make(map[string]struct{ Bid, Ask, BidSize, AskSize float64 }),
	}
}

func (c *MEXCCollector) Name() string {
	return "MEXC"
}

func (c *MEXCCollector) Start(out chan<- models.MarketData) error {
	conn, _, err := websocket.DefaultDialer.Dial(configs.MEXC_WS, nil)
	if err != nil {
		return err
	}
	c.conn = conn
	log.Println("[MEXC] connected")

	// подписка чанками
	if err := c.subscribeChunks(25); err != nil {
		return err
	}

	go c.pingLoop()
	go c.readLoop(out)

	return nil
}

func (c *MEXCCollector) Stop() error {
	c.cancel()
	if c.conn != nil {
		_ = c.conn.Close()
	}
	return nil
}

// ----------------- приватные методы -----------------

func (c *MEXCCollector) subscribeChunks(chunkSize int) error {
	chunks := chunkSymbols(c.symbols, chunkSize)
	for _, chunk := range chunks {
		params := make([]string, 0, len(chunk))
		for _, s := range chunk {
			params = append(params, "spot@public.aggre.bookTicker.v3.api.pb@100ms@"+s)
		}
		sub := map[string]interface{}{
			"method": "SUBSCRIPTION",
			"params": params,
		}
		if err := c.conn.WriteJSON(sub); err != nil {
			return err
		}
	}
	return nil
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
			continue // игнорируем текст
		}

		if mt != websocket.BinaryMessage {
			continue
		}

		var wrap pb.PushDataV3ApiWrapper
		if err := proto.Unmarshal(raw, &wrap); err != nil {
			continue
		}

		if md := c.handleWrapper(&wrap); md != nil {
			out <- *md
		}
	}
}

func (c *MEXCCollector) handleWrapper(wrap *pb.PushDataV3ApiWrapper) *models.MarketData {
	switch body := wrap.GetBody().(type) {
	case *pb.PushDataV3ApiWrapper_PublicAggreBookTicker:
		bt := body.PublicAggreBookTicker

		bid, _ := strconv.ParseFloat(bt.GetBidPrice(), 64)
		ask, _ := strconv.ParseFloat(bt.GetAskPrice(), 64)
		bidSize, _ := strconv.ParseFloat(bt.GetBidQuantity(), 64)
		askSize, _ := strconv.ParseFloat(bt.GetAskQuantity(), 64)

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

		// фильтрация по whitelist
		if len(c.allowed) > 0 {
			if _, ok := c.allowed[symbol]; !ok {
				return nil
			}
		}

		// дедупликация
		c.mu.Lock()
		last, ok := c.lastData[symbol]
		if ok && last.Bid == bid && last.Ask == ask &&
			last.BidSize == bidSize && last.AskSize == askSize {
			c.mu.Unlock()
			return nil
		}
		c.lastData[symbol] = struct {
			Bid, Ask, BidSize, AskSize float64
		}{bid, ask, bidSize, askSize}
		c.mu.Unlock()

		ts := time.Now().UnixMilli()
		if t := wrap.GetSendTime(); t > 0 {
			ts = t
		}

		return &models.MarketData{
			Exchange:  "MEXC",
			Symbol:    symbol,
			Bid:       bid,
			Ask:       ask,
			BidSize:   bidSize,
			AskSize:   askSize,
			Timestamp: ts,
		}

	default:
		return nil
	}
}

// разбивает слайс на чанки
func chunkSymbols(src []string, size int) [][]string {
	var out [][]string
	for size < len(src) {
		out = append(out, src[:size:size])
		src = src[size:]
	}
	out = append(out, src)
	return out
}



package main

import (
	"crypt_proto/internal/collector"
	"crypt_proto/pkg/models"
	"encoding/csv"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load(".env")

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
	marketDataCh := make(chan models.MarketData, 1000)

	// === читаем whitelist из CSV ===
	csvPath := "mexc_triangles_usdt_routes.csv"

	symbols, err := readSymbolsFromCSV(csvPath)
	if err != nil {
		log.Fatalf("read CSV symbols: %v", err)
	}

	log.Printf("Loaded %d unique symbols from %s", len(symbols), csvPath)

	var c collector.Collector

	switch exchange {
	case "mexc":
		c = collector.NewMEXCCollector(symbols)

	case "okx":
		c = collector.NewOKXCollector(symbols)

	case "kucoin":
		c = collector.NewKuCoinCollector(symbols)

	default:
		log.Fatal("Unknown exchange:", exchange)
	}

	// старт
	if err := c.Start(marketDataCh); err != nil {
		log.Fatal("start collector:", err)
	}

	// consumer
	go func() {
		for md := range marketDataCh {
			log.Printf("[%s] %s bid=%.8f ask=%.8f",
				md.Exchange,
				md.Symbol,
				md.Bid,
				md.Ask,
			)
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
// CSV → symbols
// ------------------------------------------------------------

func readSymbolsFromCSV(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)

	header, err := r.Read()
	if err != nil {
		return nil, err
	}

	colIndex := make(map[string]int)
	for i, h := range header {
		colIndex[strings.ToLower(strings.TrimSpace(h))] = i
	}

	required := []string{
		"leg1_symbol",
		"leg2_symbol",
		"leg3_symbol",
	}

	var idx []int
	for _, name := range required {
		i, ok := colIndex[name]
		if !ok {
			return nil, csv.ErrFieldCount
		}
		idx = append(idx, i)
	}

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
			s := strings.TrimSpace(row[i])
			if s != "" {
				uniq[s] = struct{}{}
			}
		}
	}

	out := make([]string, 0, len(uniq))
	for s := range uniq {
		out = append(out, s)
	}

	return out, nil
}

// ------------------------------------------------------------





