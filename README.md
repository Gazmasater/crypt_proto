Если оставить только нужное:

p99 execution latency
Micro-volatility (100 мс)
Fill ratio
Capture rate
Inventory drift




Название API
9623527002

696935c42a6dcd00013273f2
b348b686-55ff-4290-897b-02d55f815f65




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




go run -race main.go


GOMAXPROCS=8 go run -race main.go



package queue

import (
	"sync/atomic"
	"time"

	"crypt_proto/pkg/models"
)

type Quote struct {
	Bid, Ask         float64
	BidSize, AskSize float64
	Timestamp        int64
}

type MemoryStore struct {
	data  atomic.Value           // map[string]Quote
	batch chan *models.MarketData
}

// NewMemoryStore создаёт MemoryStore с буфером batch
func NewMemoryStore() *MemoryStore {
	s := &MemoryStore{
		batch: make(chan *models.MarketData, 10_000),
	}
	s.data.Store(make(map[string]Quote))
	return s
}

// Run обрабатывает входящие данные с минимальной задержкой
func (s *MemoryStore) Run() {
	const microBatch = 5 // применяем батч каждые 5 котировок

	batch := make([]*models.MarketData, 0, microBatch)

	for md := range s.batch {
		// фильтруем котировки с нулевым объемом
		if md.BidSize == 0 && md.AskSize == 0 {
			continue
		}

		batch = append(batch, md)

		// применяем батч
		if len(batch) >= microBatch {
			s.applyBatch(batch)
			batch = batch[:0]
		}
	}

	// оставшиеся котировки при завершении
	if len(batch) > 0 {
		s.applyBatch(batch)
	}
}

// Push добавляет MarketData в очередь
func (s *MemoryStore) Push(md *models.MarketData) {
	select {
	case s.batch <- md:
	default:
		// drop if full
	}
}

// Get возвращает snapshot lock-free
func (s *MemoryStore) Get(exchange, symbol string) (Quote, bool) {
	m := s.data.Load().(map[string]Quote)
	q, ok := m[exchange+"|"+symbol]
	return q, ok
}

// applyBatch применяет несколько котировок сразу
func (s *MemoryStore) applyBatch(batch []*models.MarketData) {
	oldMap := s.data.Load().(map[string]Quote)
	newMap := make(map[string]Quote, len(oldMap)+len(batch))

	// копируем старые данные
	for k, v := range oldMap {
		newMap[k] = v
	}

	// применяем новые котировки
	now := time.Now().UnixMilli()
	for _, md := range batch {
		key := md.Exchange + "|" + md.Symbol
		newMap[key] = Quote{
			Bid: md.Bid,
			Ask: md.Ask,
			BidSize: md.BidSize,
			AskSize: md.AskSize,
			Timestamp: now,
		}
	}

	s.data.Store(newMap)
}




