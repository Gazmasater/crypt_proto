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




package calculator

import (
	"crypt_proto/internal/queue"
	"encoding/csv"
	"log"
	"os"
	"strings"
)

const fee = 0.001 // 0.1%

// Triangle описывает один треугольный арбитраж
type Triangle struct {
	A, B, C          string // имена валют для логов
	Leg1, Leg2, Leg3 string // "BUY COTI/USDT", "SELL COTI/BTC" и т.д.
}

// Calculator считает профит по треугольникам
type Calculator struct {
	mem       *queue.MemoryStore
	triangles []Triangle
}

// NewCalculator создаёт калькулятор
func NewCalculator(mem *queue.MemoryStore, triangles []Triangle) *Calculator {
	return &Calculator{
		mem:       mem,
		triangles: triangles,
	}
}

// OnUpdate вызывается на каждый апдейт котировки
func (c *Calculator) OnUpdate(symbol string) {
	for _, tri := range c.triangles {
		s1 := legSymbol(tri.Leg1)
		s2 := legSymbol(tri.Leg2)
		s3 := legSymbol(tri.Leg3)

		// пересчитываем только если обновился один из символов треугольника
		if symbol != s1 && symbol != s2 && symbol != s3 {
			continue
		}

		c.calculateTriangle(tri, s1, s2, s3)
	}
}

// calculateTriangle рассчитывает прибыль одного треугольника
func (c *Calculator) calculateTriangle(tri Triangle, s1, s2, s3 string) {
	q1, ok1 := c.mem.Get("KuCoin", s1)
	q2, ok2 := c.mem.Get("KuCoin", s2)
	q3, ok3 := c.mem.Get("KuCoin", s3)

	if !ok1 || !ok2 || !ok3 {
		return
	}

	amount := 1.0 // стартуем с 1 A

	legs := []struct {
		leg string
		q   *queue.MarketData
	}{
		{tri.Leg1, q1},
		{tri.Leg2, q2},
		{tri.Leg3, q3},
	}

	for _, l := range legs {
		if strings.HasPrefix(l.leg, "BUY") {
			if l.q.Ask <= 0 || l.q.AskSize <= 0 {
				return
			}
			maxBuy := l.q.AskSize
			amount = amount / l.q.Ask
			if amount > maxBuy {
				amount = maxBuy
			}
			amount *= (1 - fee)
		} else {
			if l.q.Bid <= 0 || l.q.BidSize <= 0 {
				return
			}
			maxSell := l.q.BidSize
			if amount > maxSell {
				amount = maxSell
			}
			amount = amount * l.q.Bid
			amount *= (1 - fee)
		}
	}

	profit := amount - 1.0
	if profit > 0 {
		log.Printf(
			"[ARB] %s → %s → %s | profit=%.4f%% | volumes: [%.2f / %.2f / %.2f]",
			tri.A, tri.B, tri.C,
			profit*100,
			q1.BidSize, q2.BidSize, q3.BidSize,
		)
	}
}

// ParseTrianglesFromCSV парсит треугольники из CSV
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

// legSymbol извлекает символ из Leg, например "BUY COTI/USDT" -> "COTI/USDT"
func legSymbol(leg string) string {
	parts := strings.Fields(leg)
	if len(parts) != 2 {
		return ""
	}
	return strings.ToUpper(parts[1])
}





package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"crypt_proto/internal/calculator"
	"crypt_proto/internal/collector"
	"crypt_proto/internal/queue"
)

func main() {
	// -------------------- Создаем память --------------------
	mem := queue.NewMemoryStore()

	// -------------------- Загружаем треугольники --------------------
	triangles, err := calculator.ParseTrianglesFromCSV("triangles.csv")
	if err != nil {
		log.Fatalf("failed to load triangles: %v", err)
	}

	// -------------------- Создаем калькулятор --------------------
	calc := calculator.NewCalculator(mem, triangles)

	// -------------------- Создаем коллектор KuCoin --------------------
	ku := collector.NewKuCoinCollector(mem)

	// -------------------- Запуск коллектора --------------------
	// Передаем функцию обратного вызова на апдейт котировки
	ku.OnUpdate = func(symbol string) {
		calc.OnUpdate(symbol)
	}

	if err := ku.Start(); err != nil {
		log.Fatalf("failed to start KuCoin collector: %v", err)
	}

	log.Println("Calculator and KuCoin collector started")

	// -------------------- Ждем завершения --------------------
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down...")

	if err := ku.Stop(); err != nil {
		log.Printf("Error stopping collector: %v", err)
	}
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
	"sync"
	"time"

	"github.com/tidwall/gjson"
	"github.com/gorilla/websocket"

	"crypt_proto/internal/calculator"
	"crypt_proto/internal/queue"
	"crypt_proto/pkg/models"
)

const (
	maxSubsPerWS = 90
	subRate      = 120 * time.Millisecond // ~8/sec
	pingInterval = 20 * time.Second
)

/* ================= POOL ================= */
type KuCoinCollector struct {
	ctx       context.Context
	cancel    context.CancelFunc
	wsList    []*kucoinWS
	triangles []calculator.Triangle
	mem       *queue.MemoryStore // <- теперь работаем напрямую с MemoryStore
}

/* ================= WS ================= */
type kucoinWS struct {
	id      int
	conn    *websocket.Conn
	symbols []string
	last    map[string][2]float64
	mu      sync.Mutex
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

	// формируем треугольники
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
		triangles: triangles,
	}, nil
}

func (kc *KuCoinCollector) Triangles() []calculator.Triangle {
	return kc.triangles
}

/* ================= INTERFACE ================= */
func (c *KuCoinCollector) Name() string {
	return "KuCoin"
}

// Start принимает MemoryStore вместо канала
func (c *KuCoinCollector) Start(mem *queue.MemoryStore) error {
	c.mem = mem
	for _, ws := range c.wsList {
		if err := ws.connect(); err != nil {
			return err
		}
		go ws.readLoop(c)
		go ws.pingLoop()
		go ws.subscribeLoop()
	}
	log.Printf("[KuCoin] started with %d WS\n", len(c.wsList))
	return nil
}

func (c *KuCoinCollector) Stop() error {
	c.cancel()
	for _, ws := range c.wsList {
		if ws.conn != nil {
			ws.conn.Close()
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
	time.Sleep(1 * time.Second) // ждём welcome
	t := time.NewTicker(subRate)
	defer t.Stop()

	for _, s := range ws.symbols {
		<-t.C
		topic := "/market/ticker:" + s
		err := ws.conn.WriteJSON(map[string]any{
			"id":       time.Now().UnixNano(),
			"type":     "subscribe",
			"topic":    topic,
			"response": true,
		})
		if err != nil {
			log.Printf("[KuCoin WS %d] subscribe error %s\n", ws.id, s)
		} else {
			log.Printf("[KuCoin WS %d] subscribed %s\n", ws.id, s)
		}
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

/* ================= HANDLE ================= */
func (ws *kucoinWS) handle(c *KuCoinCollector, msg []byte) {
	if gjson.GetBytes(msg, "type").String() != "message" {
		return
	}

	topic := gjson.GetBytes(msg, "topic").String()
	if !strings.HasPrefix(topic, "/market/ticker:") {
		return
	}

	symbol := normalize(strings.TrimPrefix(topic, "/market/ticker:"))

	data := gjson.GetBytes(msg, "data")

	bid := data.Get("bestBid").Float()
	ask := data.Get("bestAsk").Float()
	if bid == 0 || ask == 0 {
		return
	}

	ws.mu.Lock()
	last := ws.last[symbol]
	if last[0] == bid && last[1] == ask {
		ws.mu.Unlock()
		return
	}
	ws.last[symbol] = [2]float64{bid, ask}
	ws.mu.Unlock()

	// ---- Push в MemoryStore ----
	c.mem.Push(&models.MarketData{
		Exchange:  "KuCoin",
		Symbol:    symbol,
		Bid:       bid,
		Ask:       ask,
		BidSize:   data.Get("bestBidSize").Float(),
		AskSize:   data.Get("bestAskSize").Float(),
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

func normalize(s string) string {
	parts := strings.Split(s, "-")
	return parts[0] + "/" + parts[1]
}





