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




gaz358@gaz358-BOD-WXX9:~/myprog/crypt_proto$    go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
Fetching profile over HTTP from http://localhost:6060/debug/pprof/profile?seconds=30
Saved profile in /home/gaz358/pprof/pprof.arb.samples.cpu.007.pb.gz
File: arb
Build ID: b7f6cbe195780e80f45cf9c0dc233b7b7862e62c
Type: cpu
Time: 2025-12-27 02:35:03 MSK
Duration: 30s, Total samples = 490ms ( 1.63%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 280ms, 57.14% of 490ms total
Showing top 10 nodes out of 117
      flat  flat%   sum%        cum   cum%
     130ms 26.53% 26.53%      130ms 26.53%  internal/runtime/syscall.Syscall6
      30ms  6.12% 32.65%       30ms  6.12%  runtime.futex
      20ms  4.08% 36.73%       60ms 12.24%  crypt_proto/internal/collector.(*MEXCCollector).handleWrapper
      20ms  4.08% 40.82%       60ms 12.24%  google.golang.org/protobuf/internal/impl.(*MessageInfo).initOneofFieldCoders.func1
      20ms  4.08% 44.90%       20ms  4.08%  runtime.nanotime
      20ms  4.08% 48.98%      150ms 30.61%  syscall.Syscall
      10ms  2.04% 51.02%       10ms  2.04%  crypto/internal/fips140/aes/gcm.gcmAesDec
      10ms  2.04% 53.06%       10ms  2.04%  gogo
      10ms  2.04% 55.10%       90ms 18.37%  google.golang.org/protobuf/internal/impl.(*MessageInfo).unmarshalPointerEager
      10ms  2.04% 57.14%       10ms  2.04%  internal/abi.(*Type).NumMethod
(pprof) 



gaz358@gaz358-BOD-WXX9:~/myprog/crypt_proto$ go tool pprof http://localhost:6060/debug/pprof/heap
Fetching profile over HTTP from http://localhost:6060/debug/pprof/heap
Saved profile in /home/gaz358/pprof/pprof.arb.alloc_objects.alloc_space.inuse_objects.inuse_space.002.pb.gz
File: arb
Build ID: b7f6cbe195780e80f45cf9c0dc233b7b7862e62c
Type: inuse_space
Time: 2025-12-27 02:36:50 MSK
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 2563.10kB, 100% of 2563.10kB total
Showing top 10 nodes out of 13
      flat  flat%   sum%        cum   cum%
    1539kB 60.04% 60.04%     1539kB 60.04%  runtime.allocm
  512.05kB 19.98% 80.02%   512.05kB 19.98%  runtime.main
  512.05kB 19.98%   100%   512.05kB 19.98%  runtime.acquireSudog
         0     0%   100%   512.05kB 19.98%  runtime.chanrecv
         0     0%   100%   512.05kB 19.98%  runtime.chanrecv1
         0     0%   100%     1539kB 60.04%  runtime.mcall
         0     0%   100%     1539kB 60.04%  runtime.newm
         0     0%   100%     1539kB 60.04%  runtime.park_m
         0     0%   100%     1539kB 60.04%  runtime.resetspinning
         0     0%   100%     1539kB 60.04%  runtime.schedule
(pprof) 





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
	ctx         context.Context
	cancel      context.CancelFunc
	conn        *websocket.Conn
	symbols     []string
	allowed     map[string]struct{} // whitelist
	lastData    sync.Map             // symbol -> *LastTick
	marketPool  *sync.Pool
	batchSize   int
	batchPeriod time.Duration
}

type LastTick struct {
	Bid, Ask, BidSize, AskSize float64
}

func NewMEXCCollector(symbols []string, whitelist []string, pool *sync.Pool) *MEXCCollector {
	ctx, cancel := context.WithCancel(context.Background())
	allowed := make(map[string]struct{}, len(whitelist))
	for _, s := range whitelist {
		allowed[market.NormalizeSymbol_Full(s)] = struct{}{}
	}
	return &MEXCCollector{
		ctx:         ctx,
		cancel:      cancel,
		symbols:     symbols,
		allowed:     allowed,
		lastData:    sync.Map{},
		marketPool:  pool,
		batchSize:   50,
		batchPeriod: 50 * time.Millisecond,
	}
}

func (c *MEXCCollector) Name() string {
	return "MEXC"
}

// канал принимает указатели
func (c *MEXCCollector) Start(out chan<- *models.MarketData) error {
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

func (c *MEXCCollector) readLoop(out chan<- *models.MarketData) {
	batch := make([]*models.MarketData, 0, c.batchSize)
	timer := time.NewTimer(c.batchPeriod)
	defer timer.Stop()

	flush := func() {
		for _, md := range batch {
			out <- md
		}
		batch = batch[:0]
		timer.Reset(c.batchPeriod)
	}

	for {
		select {
		case <-c.ctx.Done():
			return
		default:
		}

		_, raw, err := c.conn.ReadMessage()
		if err != nil {
			log.Printf("[MEXC] read error: %v\n", err)
			return
		}

		var wrap pb.PushDataV3ApiWrapper
		if err := proto.Unmarshal(raw, &wrap); err != nil {
			continue
		}

		md := c.handleWrapper(&wrap)
		if md != nil {
			batch = append(batch, md)
		}

		if len(batch) >= c.batchSize {
			flush()
		} else if !timer.Stop() {
			<-timer.C
			flush()
		}
	}
}

func (c *MEXCCollector) handleWrapper(wrap *pb.PushDataV3ApiWrapper) *models.MarketData {
	body, ok := wrap.GetBody().(*pb.PushDataV3ApiWrapper_PublicAggreBookTicker)
	if !ok {
		return nil
	}

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

	// дедупликация через sync.Map
	lastRaw, _ := c.lastData.Load(symbol)
	if lastRaw != nil {
		last := lastRaw.(*LastTick)
		if last.Bid == bid && last.Ask == ask && last.BidSize == bidSize && last.AskSize == askSize {
			return nil
		}
	}

	c.lastData.Store(symbol, &LastTick{Bid: bid, Ask: ask, BidSize: bidSize, AskSize: askSize})

	// берём MarketData из пула
	md := c.marketPool.Get().(*models.MarketData)
	md.Exchange = "MEXC"
	md.Symbol = symbol
	md.Bid = bid
	md.Ask = ask
	md.BidSize = bidSize
	md.AskSize = askSize
	md.Timestamp = wrap.GetSendTime()
	if md.Timestamp == 0 {
		md.Timestamp = time.Now().UnixMilli()
	}

	return md
}

func chunkSymbols(src []string, size int) [][]string {
	var out [][]string
	for size < len(src) {
		out = append(out, src[:size:size])
		src = src[size:]
	}
	if len(src) > 0 {
		out = append(out, src)
	}
	return out
}




package main

import (
	"encoding/csv"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"crypt_proto/internal/collector"
	"crypt_proto/pkg/models"

	_ "net/http/pprof"
	"github.com/joho/godotenv"
)

func main() {
	// загружаем .env
	_ = godotenv.Load(".env")

	// запускаем pprof
	go func() {
		log.Println("pprof running on http://localhost:6060/debug/pprof/")
		if err := http.ListenAndServe("localhost:6060", nil); err != nil {
			log.Printf("pprof server error: %v", err)
		}
	}()

	exchange := strings.ToLower(os.Getenv("EXCHANGE"))
	if exchange == "" {
		log.Fatal("Set EXCHANGE env variable: mexc | okx | kucoin")
	}
	log.Println("EXCHANGE:", exchange)

	// создаём пул MarketData
	marketPool := &sync.Pool{
		New: func() any {
			return &models.MarketData{}
		},
	}

	// канал маркет-данных (указатели)
	marketDataCh := make(chan *models.MarketData, 1000)

	// читаем whitelist из CSV
	csvPath := "mexc_triangles_usdt_routes.csv"
	symbols, err := readSymbolsFromCSV(csvPath)
	if err != nil {
		log.Fatalf("read CSV symbols: %v", err)
	}
	log.Printf("Loaded %d unique symbols from %s", len(symbols), csvPath)

	// создаём whitelist
	whitelist := make([]string, len(symbols))
	copy(whitelist, symbols)

	var c collector.Collector

	switch exchange {
	case "mexc":
		c = collector.NewMEXCCollector(symbols, whitelist, marketPool)
	case "okx":
		log.Fatal("OKX collector not implemented yet")
	case "kucoin":
		log.Fatal("KuCoin collector not implemented yet")
	default:
		log.Fatal("Unknown exchange:", exchange)
	}

	// старт Collector
	if err := c.Start(marketDataCh); err != nil {
		log.Fatal("start collector:", err)
	}

	// consumer
	go func() {
		for md := range marketDataCh {
			log.Printf("[%s] %s bid=%.8f ask=%.8f bidSize=%.8f askSize=%.8f",
				md.Exchange,
				md.Symbol,
				md.Bid,
				md.Ask,
				md.BidSize,
				md.AskSize,
			)
			// возвращаем объект в пул
			marketPool.Put(md)
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


