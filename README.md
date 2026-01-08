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
	m     sync.Map
	batch chan *models.MarketData
}

func NewMemoryStore() *MemoryStore {
	s := &MemoryStore{
		batch: make(chan *models.MarketData, 10_000),
	}

	// Инициализация atomic.Value ОБЯЗАТЕЛЬНА
	s.data.Store(make(map[string]Quote))

	return s
}

//
// ===== Публичный API =====
//

// Run — основной цикл стора
func (s *MemoryStore) Run() {
	for md := range s.batch {
		s.apply(md)
	}
}

// Push — приём данных от коллекторов
func (s *MemoryStore) Push(md *models.MarketData) {
	select {
	case s.batch <- md:
	default:
		// защита от переполнения
	}
}

// Get — snapshot-чтение
func (s *MemoryStore) Get(exchange, symbol string) (Quote, bool) {
	m := s.data.Load().(map[string]Quote)
	q, ok := m[exchange+"|"+symbol]
	return q, ok
}

//
// ===== Внутренняя логика =====
//

func (s *MemoryStore) apply(md *models.MarketData) {
	key := md.Exchange + "|" + md.Symbol

	quote := Quote{
		Bid:       md.Bid,
		Ask:       md.Ask,
		BidSize:   md.BidSize,
		AskSize:   md.AskSize,
		Timestamp: time.Now().UnixMilli(),
	}

	// читаем старую map
	oldMap := s.data.Load().(map[string]Quote)

	// делаем copy-on-write
	newMap := make(map[string]Quote, len(oldMap)+1)
	for k, v := range oldMap {
		newMap[k] = v
	}
	newMap[key] = quote

	// атомарно подменяем snapshot
	s.data.Store(newMap)
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

	// ------------------- Канал данных от коллекторов -------------------
	out := make(chan *models.MarketData, 100_000)

	// ------------------- In-Memory Store -------------------
	mem := queue.NewMemoryStore()
	go mem.Run()

	// прокачка данных: out → mem
	go func() {
		for md := range out {
			mem.Push(md)
		}
	}()

	// ------------------- Коллектор -------------------
	kc, err := collector.NewKuCoinCollectorFromCSV(
		"../exchange/data/kucoin_triangles_usdt.csv",
	)
	if err != nil {
		log.Fatal(err)
	}

	if err := kc.Start(out); err != nil {
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
	go calc.Run()

	// ------------------- Graceful shutdown -------------------
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("[Main] shutting down...")

	kc.Stop()
	close(out)

	log.Println("[Main] exited")
}








