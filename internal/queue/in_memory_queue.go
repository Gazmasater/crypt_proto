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

func (s *MemoryStore) Put(exchange, symbol string, q Quote) {
	old := s.data.Load().(map[string]Quote)
	newMap := make(map[string]Quote, len(old)+1)
	for k, v := range old {
		newMap[k] = v
	}
	newMap[exchange+"|"+symbol] = q
	s.data.Store(newMap)
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
