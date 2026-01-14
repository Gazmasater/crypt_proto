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
	data  atomic.Value       // snapshot: []Quote
	index map[string]int     // symbol -> индекс
	batch chan *models.MarketData
}

// NewMemoryStore создаёт store и индекс для всех символов
func NewMemoryStore(symbols []string) *MemoryStore {
	s := &MemoryStore{
		batch: make(chan *models.MarketData, 100_000),
		index: make(map[string]int),
	}

	snapshot := make([]Quote, len(symbols))
	s.data.Store(snapshot)

	for i, sym := range symbols {
		s.index[sym] = i
	}

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
	}
}

// быстрый доступ по индексу
func (s *MemoryStore) GetByIndex(idx int) (Quote, bool) {
	snap := s.data.Load().([]Quote)
	if idx < 0 || idx >= len(snap) {
		return Quote{}, false
	}
	return snap[idx], true
}

// быстрый доступ по символу
func (s *MemoryStore) Get(symbol string) (Quote, bool) {
	idx, ok := s.index[symbol]
	if !ok {
		return Quote{}, false
	}
	return s.GetByIndex(idx)
}

// получение индекса символа
func (s *MemoryStore) Index(symbol string) int {
	if idx, ok := s.index[symbol]; ok {
		return idx
	}
	return -1
}

// обновление snapshot
func (s *MemoryStore) apply(md *models.MarketData) {
	idx, ok := s.index[md.Symbol]
	if !ok {
		return
	}

	old := s.data.Load().([]Quote)
	newSnap := make([]Quote, len(old))
	copy(newSnap, old)

	newSnap[idx] = Quote{
		Bid:       md.Bid,
		Ask:       md.Ask,
		BidSize:   md.BidSize,
		AskSize:   md.AskSize,
		Timestamp: time.Now().UnixMilli(),
	}

	s.data.Store(newSnap)
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
	A, B, C         string
	Leg1Idx, Leg2Idx, Leg3Idx int
	Buy1, Buy2, Buy3 bool
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

func (c *CalculatorFast) RunFast(in <-chan *queue.Quote) {
	for _, tri := range c.triangles {
		_ = tri // заглушка для первой версии
	}
}

func (c *CalculatorFast) calcTriangleFast(tri TriangleFast) {
	q1, ok1 := c.mem.GetByIndex(tri.Leg1Idx)
	q2, ok2 := c.mem.GetByIndex(tri.Leg2Idx)
	q3, ok3 := c.mem.GetByIndex(tri.Leg3Idx)
	if !ok1 || !ok2 || !ok3 {
		return
	}

	// LEG1
	amount := 1.0
	if tri.Buy1 {
		if q1.Ask <= 0 || q1.AskSize <= 0 {
			return
		}
		amount = amount / q1.Ask * (1 - fee)
	} else {
		if q1.Bid <= 0 || q1.BidSize <= 0 {
			return
		}
		amount = amount * q1.Bid * (1 - fee)
	}

	// LEG2
	if tri.Buy2 {
		if q2.Ask <= 0 || q2.AskSize <= 0 {
			return
		}
		amount = amount / q2.Ask * (1 - fee)
	} else {
		amount = amount * q2.Bid * (1 - fee)
	}

	// LEG3
	if tri.Buy3 {
		if q3.Ask <= 0 || q3.AskSize <= 0 {
			return
		}
		amount = amount / q3.Ask * (1 - fee)
	} else {
		amount = amount * q3.Bid * (1 - fee)
	}

	profit := amount - 1.0
	if profit > 0.001 {
		msg := fmt.Sprintf("[ARB] %s→%s→%s profit=%.4f", tri.A, tri.B, tri.C, profit)
		fmt.Println(msg)
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
	"crypt_proto/pkg/models"
)

func main() {
	// ------------------- pprof -------------------
	go func() {
		log.Println("pprof on http://localhost:6060/debug/pprof/")
		_ = http.ListenAndServe("localhost:6060", nil)
	}()

	// ------------------- Коллектор -------------------
	kc, symbols, err := collector.NewKuCoinCollectorFromCSV("../exchange/data/kucoin_triangles_usdt.csv")
	if err != nil {
		log.Fatal(err)
	}

	mem := queue.NewMemoryStore(symbols)
	go mem.Run()

	out := make(chan *models.MarketData, 100_000)
	go func() {
		for md := range out {
			mem.Push(md)
		}
	}()

	if err := kc.Start(out); err != nil {
		log.Fatal(err)
	}

	trianglesCSV, _ := calculator.ParseTrianglesFromCSV("../exchange/data/kucoin_triangles_usdt.csv")
	var triangles []calculator.TriangleFast
	for _, t := range trianglesCSV {
		leg1, leg2, leg3 := t.Leg1, t.Leg2, t.Leg3
		triangles = append(triangles, calculator.TriangleFast{
			A: t.A, B: t.B, C: t.C,
			Leg1Idx: mem.Index(leg1),
			Leg2Idx: mem.Index(leg2),
			Leg3Idx: mem.Index(leg3),
			Buy1: leg1[:3] == "BUY",
			Buy2: leg2[:3] == "BUY",
			Buy3: leg3[:3] == "BUY",
		})
	}

	calc := calculator.NewCalculatorFast(mem, triangles)
	go func() {
		for md := range out {
			_ = md
			for _, tri := range triangles {
				calc.calcTriangleFast(tri)
			}
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	log.Println("[Main] shutting down...")
	kc.Stop()
	close(out)
	log.Println("[Main] exited")
}




