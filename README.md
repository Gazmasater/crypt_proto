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
	"sync"
	"sync/atomic"
	"time"

	"crypt_proto/pkg/models"
)

type Quote struct {
	Bid       float64
	Ask       float64
	BidSize   float64
	AskSize   float64
	Timestamp int64
}

type MemoryStore struct {
	data  atomic.Value // map[string]Quote
	batch chan *models.MarketData
	// Индексация для быстрого поиска Leg → symbol
	index map[string]string
	mu    sync.Mutex
}

func NewMemoryStore() *MemoryStore {
	s := &MemoryStore{
		batch: make(chan *models.MarketData, 100_000),
		index: make(map[string]string),
	}
	s.data.Store(make(map[string]Quote))
	return s
}

func (s *MemoryStore) Run() {
	for md := range s.batch {
		s.apply(md)
	}
}

func (s *MemoryStore) Push(md *models.MarketData) {
	select {
	case s.batch <- md:
	default:
		// защита от переполнения
	}
}

func (s *MemoryStore) Get(exchange, symbol string) (Quote, bool) {
	m := s.data.Load().(map[string]Quote)
	q, ok := m[exchange+"|"+symbol]
	return q, ok
}

// Индекс символа → Leg
func (s *MemoryStore) Index(leg string) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	if v, ok := s.index[leg]; ok {
		return v
	}
	s.index[leg] = legSymbol(leg)
	return s.index[leg]
}

func (s *MemoryStore) apply(md *models.MarketData) {
	key := md.Exchange + "|" + md.Symbol
	quote := Quote{
		Bid:       md.Bid,
		Ask:       md.Ask,
		BidSize:   md.BidSize,
		AskSize:   md.AskSize,
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

// Вспомогательная функция
func legSymbol(leg string) string {
	parts := strings.Fields(leg)
	if len(parts) < 2 {
		return ""
	}
	return parts[1]
}




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
	"time"

	"github.com/gorilla/websocket"
	"github.com/tidwall/gjson"

	"crypt_proto/pkg/models"
	"crypt_proto/internal/queue"
	"crypt_proto/internal/calculator"
)

const (
	maxSubsPerWS = 90
	subRate      = 120 * time.Millisecond
	pingInterval = 20 * time.Second
)

type KuCoinCollector struct {
	ctx       context.Context
	cancel    context.CancelFunc
	mem       *queue.MemoryStore
	wsList    []*kucoinWS
	out       chan<- *models.MarketData
	triangles []calculator.TriangleFast
}

type kucoinWS struct {
	id      int
	conn    *websocket.Conn
	symbols []string
	last    map[string][2]float64
}

func NewKuCoinCollectorFromCSV(path string, mem *queue.MemoryStore) (*KuCoinCollector, error) {
	symbols, err := readPairsFromCSV(path)
	if err != nil {
		return nil, err
	}
	if len(symbols) == 0 {
		return nil, fmt.Errorf("no symbols")
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

	return &KuCoinCollector{
		ctx:    ctx,
		cancel: cancel,
		mem:    mem,
		wsList: wsList,
	}, nil
}

func (kc *KuCoinCollector) Name() string {
	return "KuCoin"
}

func (kc *KuCoinCollector) Start(out chan<- *models.MarketData) error {
	kc.out = out
	for _, ws := range kc.wsList {
		if err := ws.connect(); err != nil {
			return err
		}
		go ws.readLoop(kc)
		go ws.subscribeLoop()
		go ws.pingLoop()
	}
	log.Printf("[KuCoin] started with %d WS\n", len(kc.wsList))
	return nil
}

func (kc *KuCoinCollector) Stop() error {
	kc.cancel()
	for _, ws := range kc.wsList {
		if ws.conn != nil {
			_ = ws.conn.Close()
		}
	}
	return nil
}

func (ws *kucoinWS) connect() error {
	req, _ := http.NewRequest("POST", "https://api.kucoin.com/api/v1/bullet-public", nil)
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

	url := fmt.Sprintf("%s?token=%s&connectId=%d",
		r.Data.InstanceServers[0].Endpoint,
		r.Data.Token,
		time.Now().UnixNano(),
	)

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return err
	}

	ws.conn = conn
	log.Printf("[KuCoin WS %d] connected\n", ws.id)
	return nil
}

func (ws *kucoinWS) subscribeLoop() {
	time.Sleep(time.Second)
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
		_ = ws.conn.WriteJSON(map[string]any{
			"id":   time.Now().UnixNano(),
			"type": "ping",
		})
	}
}

func (ws *kucoinWS) readLoop(kc *KuCoinCollector) {
	for {
		_, msg, err := ws.conn.ReadMessage()
		if err != nil {
			log.Printf("[KuCoin WS %d] read error: %v\n", ws.id, err)
			return
		}
		ws.handle(kc, msg)
	}
}

func (ws *kucoinWS) handle(kc *KuCoinCollector, msg []byte) {
	if gjson.GetBytes(msg, "type").String() != "message" {
		return
	}
	topic := gjson.GetBytes(msg, "topic").String()
	if !strings.HasPrefix(topic, "/market/ticker:") {
		return
	}
	symbol := strings.TrimPrefix(topic, "/market/ticker:")
	data := gjson.GetBytes(msg, "data")
	bid := data.Get("bestBid").Float()
	ask := data.Get("bestAsk").Float()
	bidSize := data.Get("bestBidSize").Float()
	askSize := data.Get("bestAskSize").Float()
	if bid == 0 || ask == 0 {
		return
	}

	last := ws.last[symbol]
	if last[0] == bid && last[1] == ask {
		return
	}
	ws.last[symbol] = [2]float64{bid, ask}

	kc.out <- &models.MarketData{
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

	r := csv.NewReader(f)
	rows, err := r.ReadAll()
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





package calculator

import (
	"fmt"
	"log"
	"os"

	"crypt_proto/internal/queue"
)

const fee = 0.001

type TriangleFast struct {
	A, B, C       string
	Leg1Idx, Leg2Idx, Leg3Idx string
	Buy1, Buy2, Buy3           bool
}

type CalculatorFast struct {
	mem       *queue.MemoryStore
	triangles []TriangleFast
	fileLog   *log.Logger
}

func NewCalculatorFast(mem *queue.MemoryStore, triangles []TriangleFast) *CalculatorFast {
	f, err := os.OpenFile("arb_opportunities.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("failed to open arb log file: %v", err)
	}
	return &CalculatorFast{
		mem:       mem,
		triangles: triangles,
		fileLog:   log.New(f, "", log.LstdFlags),
	}
}

func (c *CalculatorFast) CalcTriangleFast(tri TriangleFast) {
	q1, ok1 := c.mem.Get("KuCoin", tri.Leg1Idx)
	q2, ok2 := c.mem.Get("KuCoin", tri.Leg2Idx)
	q3, ok3 := c.mem.Get("KuCoin", tri.Leg3Idx)
	if !ok1 || !ok2 || !ok3 {
		return
	}

	// Рассчёт лимитов
	usdt1 := q1.Ask*q1.AskSize
	if !tri.Buy1 {
		usdt1 = q1.Bid * q1.BidSize
	}
	usdt2 := q2.Ask*q2.AskSize
	if !tri.Buy2 {
		usdt2 = q2.Bid * q2.BidSize
	}
	usdt3 := q3.Ask * q3.AskSize
	if !tri.Buy3 {
		usdt3 = q3.Bid * q3.BidSize
	}

	maxUSDT := usdt1
	if usdt2 < maxUSDT {
		maxUSDT = usdt2
	}
	if usdt3 < maxUSDT {
		maxUSDT = usdt3
	}
	if maxUSDT <= 0 {
		return
	}

	amount := maxUSDT
	if tri.Buy1 {
		amount = amount / q1.Ask * (1 - fee)
	} else {
		amount = amount * q1.Bid * (1 - fee)
	}
	if tri.Buy2 {
		amount = amount / q2.Ask * (1 - fee)
	} else {
		amount = amount * q2.Bid * (1 - fee)
	}
	if tri.Buy3 {
		amount = amount / q3.Ask * (1 - fee)
	} else {
		amount = amount * q3.Bid * (1 - fee)
	}

	profitUSDT := amount - maxUSDT
	profitPct := profitUSDT / maxUSDT

	if profitPct > 0.001 && profitUSDT > 0.02 {
		msg := fmt.Sprintf("[ARB] %s → %s → %s | %.4f%% | volume=%.2f USDT | profit=%.4f USDT",
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
	"time"

	"crypt_proto/internal/collector"
	"crypt_proto/internal/calculator"
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

	kc, err := collector.NewKuCoinCollectorFromCSV("../exchange/data/kucoin_triangles_usdt.csv", mem)
	if err != nil {
		log.Fatal(err)
	}
	if err := kc.Start(out); err != nil {
		log.Fatal(err)
	}

	trianglesCSV, err := collector.ReadTrianglesCSV("../exchange/data/kucoin_triangles_usdt.csv")
	if err != nil {
		log.Fatal(err)
	}

	var triangles []calculator.TriangleFast
	for _, t := range trianglesCSV {
		triangles = append(triangles, calculator.TriangleFast{
			A: t.A, B: t.B, C: t.C,
			Leg1Idx: mem.Index(t.Leg1),
			Leg2Idx: mem.Index(t.Leg2),
			Leg3Idx: mem.Index(t.Leg3),
			Buy1: t.Leg1[:3] == "BUY",
			Buy2: t.Leg2[:3] == "BUY",
			Buy3: t.Leg3[:3] == "BUY",
		})
	}

	calc := calculator.NewCalculatorFast(mem, triangles)
	go func() {
		ticker := time.NewTicker(10 * time.Millisecond)
		defer ticker.Stop()
		for range ticker.C {
			for _, tri := range triangles {
				calc.CalcTriangleFast(tri)
			}
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("[Main] shutting down...")
	_ = kc.Stop()
	close(out)
	log.Println("[Main] exited")
}








