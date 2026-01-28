package queue

import (
	"sync"
	"time"

	"crypt_proto/pkg/models"
)

type Quote struct {
	Bid, Ask         float64
	BidSize, AskSize float64
	Timestamp        int64
}

type MemoryStore struct {
	mu    sync.RWMutex
	data  map[string]Quote
	batch chan *models.MarketData
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data:  make(map[string]Quote),
		batch: make(chan *models.MarketData, 10_000),
	}
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

// Get snapshot lock-safe
func (s *MemoryStore) Get(exchange, symbol string) (Quote, bool) {
	s.mu.RLock()
	q, ok := s.data[exchange+"|"+symbol]
	s.mu.RUnlock()
	return q, ok
}

// apply обновляет тикер на месте без копирования всей карты
func (s *MemoryStore) apply(md *models.MarketData) {
	key := md.Exchange + "|" + md.Symbol
	quote := Quote{
		Bid:       md.Bid,
		Ask:       md.Ask,
		BidSize:   md.BidSize,
		AskSize:   md.AskSize,
		Timestamp: time.Now().UnixMilli(),
	}

	s.mu.Lock()
	s.data[key] = quote
	s.mu.Unlock()
}
