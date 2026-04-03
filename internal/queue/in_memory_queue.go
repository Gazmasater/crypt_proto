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
	mu   sync.RWMutex
	data map[string]Quote
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data: make(map[string]Quote),
	}
}

// Run оставлен для совместимости со старым кодом.
func (s *MemoryStore) Run() {}

func (s *MemoryStore) Push(md *models.MarketData) {
	if md == nil {
		return
	}
	s.apply(md)
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
	timestamp := md.Timestamp
	if timestamp == 0 {
		timestamp = time.Now().UnixMilli()
	}
	quote := Quote{
		Bid:       md.Bid,
		Ask:       md.Ask,
		BidSize:   md.BidSize,
		AskSize:   md.AskSize,
		Timestamp: timestamp,
	}

	s.mu.Lock()
	s.data[key] = quote
	s.mu.Unlock()
}
