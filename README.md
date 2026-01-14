Название API
9623527002

6966b78122ca320001d2acae
fa1e37ae-21ff-4257-844d-3dcd21d26ccd





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



package queue

import (
	"sync/atomic"
	"time"

	"crypt_proto/pkg/models"
)

type Quote struct {
	Bid, Ask       float64
	BidSize, AskSize float64
	Timestamp      int64
}

type MemoryStore struct {
	data  atomic.Value // map[string]Quote
	batch chan *models.MarketData
}

func NewMemoryStore() *MemoryStore {
	s := &MemoryStore{
		batch: make(chan *models.MarketData, 10_000),
	}
	s.data.Store(make(map[string]Quote))
	return s
}

// Run обрабатывает входящие данные и обновляет snapshot
func (s *MemoryStore) Run() {
	for md := range s.batch {
		s.apply(md)
	}
}

func (s *MemoryStore) Push(md *models.MarketData) {
	select {
	case s.batch <- md:
	default:
		// drop if full
	}
}

// Get snapshot lock-free
func (s *MemoryStore) Get(exchange, symbol string) (Quote, bool) {
	m := s.data.Load().(map[string]Quote)
	q, ok := m[exchange+"|"+symbol]
	return q, ok
}

func (s *MemoryStore) apply(md *models.MarketData) {
	key := md.Exchange + "|" + md.Symbol
	quote := Quote{
		Bid: md.Bid, Ask: md.Ask,
		BidSize: md.BidSize, AskSize: md.AskSize,
		Timestamp: time.Now().UnixMilli(),
	}

	oldMap := s.data.Load().(map[string]Quote)
	newMap := make(map[string]Quote, len(oldMap)+1)
	for k, v := range oldMap {
		newMap[k] = v
	}
	newMap[key] = quote
	s.data.Store(newMap)
}



package collector

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"crypt_proto/pkg/models"

	"github.com/gorilla/websocket"
	"github.com/tidwall/gjson"
)

const (
	maxSubsPerWS = 110
	subRate      = 120 * time.Millisecond
	pingInterval = 20 * time.Second
)

type KuCoinCollector struct {
	ctx    context.Context
	cancel context.CancelFunc
	wsList []*kucoinWS
	out    chan<- *models.MarketData
}

type kucoinWS struct {
	id      int
	conn    *websocket.Conn
	symbols []string
	last    map[string][2]float64
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
			last:    make(map[string][2]float64),
		})
	}

	return &KuCoinCollector{ctx: ctx, cancel: cancel, wsList: wsList}, symbols, nil
}

func (c *KuCoinCollector) Name() string { return "KuCoin" }

func (c *KuCoinCollector) Start(out chan<- *models.MarketData) error {
	c.out = out
	for _, ws := range c.wsList {
		if err := ws.connect(); err != nil {
			return err
		}
		go ws.readLoop(c)
		go ws.subscribeLoop()
		go ws.pingLoop()
	}
	log.Printf("[KuCoin] started with %d WS\n", len(c.wsList))
	return nil
}

func (c *KuCoinCollector) Stop() error {
	c.cancel()
	for _, ws := range c.wsList {
		if ws.conn != nil {
			_ = ws.conn.Close()
		}
	}
	return nil
}

func (ws *kucoinWS) connect() error {
	// минимальный connect, реальная логика API
	return nil
}

func (ws *kucoinWS) subscribeLoop() {
	t := time.NewTicker(subRate)
	defer t.Stop()
	for _, s := range ws.symbols {
		<-t.C
		_ = ws.conn.WriteJSON(map[string]any{
			"id":       time.Now().UnixNano(),
			"type":     "subscribe",
			"topic":    "/market/ticker:" + s,
			"response": true,
		})
	}
}

func (ws *kucoinWS) pingLoop() {
	t := time.NewTicker(pingInterval)
	defer t.Stop()
	for range t.C {
		_ = ws.conn.WriteJSON(map[string]any{"id": time.Now().UnixNano(), "type": "ping"})
	}
}

func (ws *kucoinWS) readLoop(c *KuCoinCollector) {
	for {
		_, msg, err := ws.conn.ReadMessage()
		if err != nil {
			log.Printf("[KuCoin WS %d] read error: %v\n", ws.id, err)
			return
		}
		ws.handle(c, msg)
	}
}

func (ws *kucoinWS) handle(c *KuCoinCollector, msg []byte) {
	if gjson.GetBytes(msg, "type").String() != "message" {
		return
	}
	topic := gjson.GetBytes(msg, "topic").String()
	if !strings.HasPrefix(topic, "/market/ticker:") {
		return
	}
	symbol := strings.TrimPrefix(topic, "/market/ticker:")
	data := gjson.GetBytes(msg, "data")
	bid, ask := data.Get("bestBid").Float(), data.Get("bestAsk").Float()
	bidSize, askSize := data.Get("bestBidSize").Float(), data.Get("bestAskSize").Float()
	if bid == 0 || ask == 0 {
		return
	}
	last := ws.last[symbol]
	if last[0] == bid && last[1] == ask {
		return
	}
	ws.last[symbol] = [2]float64{bid, ask}
	c.out <- &models.MarketData{
		Exchange:  "KuCoin",
		Symbol:    symbol,
		Bid:       bid,
		Ask:       ask,
		BidSize:   bidSize,
		AskSize:   askSize,
		Timestamp: time.Now().UnixMilli(),
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
	set := make(map[string]struct{})
	for _, row := range rows[1:] {
		for i := 3; i <= 5 && i < len(row); i++ {
			if p := parseLeg(row[i]); p != "" {
				set[p] = struct{}{}
			}
		}
	}
	res := make([]string, 0, len(set))
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




package calculator

import (
	"fmt"
	"log"
	"os"
	"strings"

	"crypt_proto/internal/queue"
	"crypt_proto/pkg/models"
)

const fee = 0.001

type LegIndex struct {
	Key   string
	IsBuy bool
}

type Triangle struct {
	A, B, C string
	Legs    [3]LegIndex
}

type Calculator struct {
	mem       *queue.MemoryStore
	triangles []*Triangle
	fileLog   *log.Logger
}

// NewCalculator полностью lock-free
func NewCalculator(mem *queue.MemoryStore, triangles []*Triangle) *Calculator {
	f, err := os.OpenFile("arb_opportunities.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("failed to open log: %v", err)
	}

	return &Calculator{
		mem:       mem,
		triangles: triangles,
		fileLog:   log.New(f, "", log.LstdFlags),
	}
}

// Run читает входящие данные и пересчитывает треугольники без блокировок
func (c *Calculator) Run(in <-chan *models.MarketData) {
	for md := range in {
		c.mem.Push(md)
		// lock-free: проверяем каждый треугольник
		for _, tri := range c.triangles {
			c.calcTriangle(tri)
		}
	}
}

func (c *Calculator) calcTriangle(tri *Triangle) {
	q := [3]queue.Quote{}
	for i, leg := range tri.Legs {
		quote, ok := c.mem.Get("KuCoin", strings.Split(leg.Key, "|")[1])
		if !ok {
			return
		}
		q[i] = quote
	}

	usdtLimits := [3]float64{}
	// LEG1
	if tri.Legs[0].IsBuy {
		if q[0].Ask <= 0 || q[0].AskSize <= 0 {
			return
		}
		usdtLimits[0] = q[0].Ask * q[0].AskSize
	} else {
		if q[0].Bid <= 0 || q[0].BidSize <= 0 {
			return
		}
		usdtLimits[0] = q[0].Bid * q[0].BidSize
	}
	// LEG2
	if tri.Legs[1].IsBuy {
		if q[1].Ask <= 0 || q[1].AskSize <= 0 || q[2].Bid <= 0 {
			return
		}
		usdtLimits[1] = q[1].Ask * q[1].AskSize * q[2].Bid
	} else {
		if q[1].Bid <= 0 || q[1].BidSize <= 0 || q[2].Bid <= 0 {
			return
		}
		usdtLimits[1] = q[1].BidSize * q[2].Bid
	}
	// LEG3
	if q[2].Bid <= 0 || q[2].BidSize <= 0 {
		return
	}
	usdtLimits[2] = q[2].Bid * q[2].BidSize

	maxUSDT := usdtLimits[0]
	if usdtLimits[1] < maxUSDT {
		maxUSDT = usdtLimits[1]
	}
	if usdtLimits[2] < maxUSDT {
		maxUSDT = usdtLimits[2]
	}
	if maxUSDT <= 0 {
		return
	}

	amount := maxUSDT
	if tri.Legs[0].IsBuy {
		amount = amount / q[0].Ask * (1 - fee)
	} else {
		amount = amount * q[0].Bid * (1 - fee)
	}
	if tri.Legs[1].IsBuy {
		amount = amount / q[1].Ask * (1 - fee)
	} else {
		amount = amount * q[1].Bid * (1 - fee)
	}
	if tri.Legs[2].IsBuy {
		amount = amount / q[2].Ask * (1 - fee)
	} else {
		amount = amount * q[2].Bid * (1 - fee)
	}

	profitUSDT := amount - maxUSDT
	profitPct := profitUSDT / maxUSDT

	if profitPct > 0.001 && profitUSDT > 0.02 {
		msg := fmt.Sprintf("[ARB] %s→%s→%s | %.4f%% | volume=%.2f USDT | profit=%.4f USDT",
			tri.A, tri.B, tri.C, profitPct*100, maxUSDT, profitUSDT)
		log.Println(msg)
		c.fileLog.Println(msg)
	}
}





package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"crypt_proto/internal/calculator"
	"crypt_proto/internal/collector"
	"crypt_proto/internal/queue"
)

func main() {
	go func() {
		log.Println("pprof on http://localhost:6060/debug/pprof/")
		_ = http.ListenAndServe("localhost:6060", nil)
	}()

	out := make(chan *queue.Quote, 100_000)
	mem := queue.NewMemoryStore()
	go mem.Run()

	kc, _, err := collector.NewKuCoinCollectorFromCSV("../exchange/data/kucoin_triangles_usdt.csv")
	if err != nil {
		log.Fatal(err)
	}
	if err := kc.Start(out); err != nil {
		log.Fatal(err)
	}
	log.Println("[Main] KuCoinCollector started")

	triangles, _ := calculator.ParseTrianglesFromCSV("../exchange/data/kucoin_triangles_usdt.csv")
	calc := calculator.NewCalculator(mem, triangles)
	go calc.Run(out)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("[Main] shutting down...")
	kc.Stop()
	close(out)
	log.Println("[Main] exited")
}

