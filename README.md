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
)

// Quote описывает один тик
type Quote struct {
	Bid       float64
	Ask       float64
	BidSize   float64
	AskSize   float64
	Timestamp int64
}

// AtomicQuote — lock-free snapshot для одного символа
type AtomicQuote struct {
	v atomic.Pointer[Quote]
}

func (q *AtomicQuote) Store(n Quote) {
	q.v.Store(&n)
}

func (q *AtomicQuote) Load() *Quote {
	return q.v.Load()
}

// MemoryStore хранит все символы
type MemoryStore struct {
	data sync.Map // key "Exchange|Symbol" -> *AtomicQuote
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{}
}

// Update обновляет один тик
func (s *MemoryStore) Update(exchange, symbol string, q Quote) {
	key := exchange + "|" + symbol
	v, _ := s.data.LoadOrStore(key, &AtomicQuote{})
	v.(*AtomicQuote).Store(q)
}

// Get возвращает snapshot
func (s *MemoryStore) Get(exchange, symbol string) (*Quote, bool) {
	key := exchange + "|" + symbol
	v, ok := s.data.Load(key)
	if !ok {
		return nil, false
	}
	q := v.(*AtomicQuote).Load()
	return q, q != nil
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

	"github.com/tidwall/gjson"

	"crypt_proto/internal/calculator"
	"crypt_proto/internal/queue"
)

const (
	maxSubsPerWS = 90
	subRate      = 120 * time.Millisecond
	pingInterval = 20 * time.Second
)

/* ================= POOL ================= */

type KuCoinCollector struct {
	ctx       context.Context
	cancel    context.CancelFunc
	wsList    []*kucoinWS
	mem       *queue.MemoryStore
	triangles []calculator.Triangle
}

/* ================= WS ================= */

type kucoinWS struct {
	id      int
	conn    *WebsocketConn
	symbols []string
	last    map[string][2]float64
}

/* ================= CONSTRUCTOR ================= */

func NewKuCoinCollectorFromCSV(path string, mem *queue.MemoryStore) (*KuCoinCollector, error) {
	symbols, err := readPairsFromCSV(path)
	if err != nil {
		return nil, err
	}
	if len(symbols) == 0 {
		return nil, fmt.Errorf("no symbols")
	}

	triangles, _ := calculator.ParseTrianglesFromCSV(path)
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
		ctx:       ctx,
		cancel:    cancel,
		wsList:    wsList,
		mem:       mem,
		triangles: triangles,
	}, nil
}

func (kc *KuCoinCollector) Start() error {
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

/* ================= CONNECT ================= */

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

/* ================= SUBSCRIBE ================= */

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

/* ================= PING ================= */

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

/* ================= READ ================= */

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

/* ================= HANDLE ================= */

func (ws *kucoinWS) handle(kc *KuCoinCollector, msg []byte) {
	if gjson.GetBytes(msg, "type").Str != "message" {
		return
	}

	topic := gjson.GetBytes(msg, "topic").Str
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

	kc.mem.Update("KuCoin", symbol, queue.Quote{
		Bid:       bid,
		Ask:       ask,
		BidSize:   bidSize,
		AskSize:   askSize,
		Timestamp: time.Now().UnixMilli(),
	})
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

/* ================= HELPERS ================= */

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
	"crypt_proto/internal/queue"
	"fmt"
	"log"
	"os"
	"strings"
)

const fee = 0.001 // 0.1%

type Triangle struct {
	A, B, C          string
	Leg1, Leg2, Leg3 string
}

type Calculator struct {
	mem       *queue.MemoryStore
	triangles []Triangle
	fileLog   *log.Logger
}

func NewCalculator(mem *queue.MemoryStore, triangles []Triangle) *Calculator {
	f, err := os.OpenFile("arb_opportunities.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("failed to open arb log file: %v", err)
	}
	return &Calculator{
		mem:       mem,
		triangles: triangles,
		fileLog:   log.New(f, "", log.LstdFlags),
	}
}

// Run — polling loop
func (c *Calculator) Run() {
	for {
		for _, tri := range c.triangles {
			c.calcTriangle(&tri)
		}
	}
}

func (c *Calculator) calcTriangle(tri *Triangle) {
	s1, s2, s3 := legSymbol(tri.Leg1), legSymbol(tri.Leg2), legSymbol(tri.Leg3)
	q1, ok1 := c.mem.Get("KuCoin", s1)
	q2, ok2 := c.mem.Get("KuCoin", s2)
	q3, ok3 := c.mem.Get("KuCoin", s3)
	if !ok1 || !ok2 || !ok3 {
		return
	}

	var usdtLimits [3]float64
	i := 0
	if strings.HasPrefix(tri.Leg1, "BUY") {
		usdtLimits[i] = q1.Ask * q1.AskSize
	} else {
		usdtLimits[i] = q1.Bid * q1.BidSize
	}
	i++
	if strings.HasPrefix(tri.Leg2, "BUY") {
		usdtLimits[i] = q2.Ask * q2.AskSize * q3.Bid
	} else {
		usdtLimits[i] = q2.BidSize * q3.Bid
	}
	i++
	usdtLimits[i] = q3.Bid * q3.BidSize

	maxUSDT := usdtLimits[0]
	for _, v := range usdtLimits[1:] {
		if v < maxUSDT {
			maxUSDT = v
		}
	}
	if maxUSDT <= 0 {
		return
	}

	amount := maxUSDT
	if strings.HasPrefix(tri.Leg1, "BUY") {
		amount = amount / q1.Ask * (1 - fee)
	} else {
		amount = amount * q1.Bid * (1 - fee)
	}
	if strings.HasPrefix(tri.Leg2, "BUY") {
		amount = amount / q2.Ask * (1 - fee)
	} else {
		amount = amount * q2.Bid * (1 - fee)
	}
	if strings.HasPrefix(tri.Leg3, "BUY") {
		amount = amount / q3.Ask * (1 - fee)
	} else {
		amount = amount * q3.Bid * (1 - fee)
	}

	profitUSDT := amount - maxUSDT
	profitPct := profitUSDT / maxUSDT

	if profitPct > 0.001 && profitUSDT > 0.02 {
		msg := fmt.Sprintf(
			"[ARB] %s → %s → %s | %.4f%% | volume=%.2f USDT | profit=%.4f USDT",
			tri.A, tri.B, tri.C, profitPct*100, maxUSDT, profitUSDT,
		)
		fmt.Println(msg)
		c.fileLog.Println(msg)
	}
}

func ParseTrianglesFromCSV(path string) ([]Triangle, error) {
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

	var res []Triangle
	for _, row := range rows[1:] {
		if len(row) < 6 {
			continue
		}
		res = append(res, Triangle{
			A:    strings.TrimSpace(row[0]),
			B:    strings.TrimSpace(row[1]),
			C:    strings.TrimSpace(row[2]),
			Leg1: strings.TrimSpace(row[3]),
			Leg2: strings.TrimSpace(row[4]),
			Leg3: strings.TrimSpace(row[5]),
		})
	}

	return res, nil
}

func legSymbol(leg string) string {
	parts := strings.Fields(strings.ToUpper(strings.TrimSpace(leg)))
	if len(parts) < 2 {
		return ""
	}
	p := strings.Split(parts[1], "/")
	if len(p) != 2 {
		return ""
	}
	return p[0] + "-" + p[1]
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

	"crypt_proto/internal/calculator"
	"crypt_proto/internal/collector"
	"crypt_proto/internal/queue"
)

func main() {
	// ------------------- pprof -------------------
	go func() {
		log.Println("pprof on http://localhost:6060/debug/pprof/")
		if err := http.ListenAndServe("localhost:6060", nil); err != nil {
			log.Fatalf("pprof listen error: %v", err)
		}
	}()

	// ------------------- In-Memory Store -------------------
	mem := queue.NewMemoryStore()

	// ------------------- Коллектор -------------------
	kc, err := collector.NewKuCoinCollectorFromCSV(
		"../exchange/data/kucoin_triangles_usdt.csv",
		mem, // передаем память напрямую
	)
	if err != nil {
		log.Fatal(err)
	}

	if err := kc.Start(); err != nil {
		log.Fatal(err)
	}
	log.Println("[Main] KuCoinCollector started")

	// ------------------- Треугольники -------------------
	triangles, err := calculator.ParseTrianglesFromCSV(
		"../exchange/data/kucoin_triangles_usdt.csv",
	)
	if err != nil {
		log.Fatal(err)
	}

	// ------------------- Калькулятор -------------------
	calc := calculator.NewCalculator(mem, triangles)
	go func() {
		for {
			calc.Run() // polling loop
			time.Sleep(1 * time.Millisecond) // минимальная пауза, чтобы CPU не 100%
		}
	}()

	// ------------------- Graceful shutdown -------------------
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("[Main] shutting down...")

	// останавливаем коллектора
	if err := kc.Stop(); err != nil {
		log.Printf("error stopping collector: %v", err)
	}

	log.Println("[Main] exited")
}




import "github.com/gorilla/websocket"

type kucoinWS struct {
	id      int
	conn    *websocket.Conn  // <- исправлено
	symbols []string
	last    map[string][2]float64
}





