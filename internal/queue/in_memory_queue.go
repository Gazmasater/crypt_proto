package queue

import (
	"sync"
	"sync/atomic"

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
	return &MemoryStore{
		batch: make(chan *models.MarketData, 10_000),
	}
}

// Snapshot возвращает актуальные котировки

func (s *MemoryStore) Get(exchange, symbol string) (Quote, bool) {
	m := s.data.Load().(map[string]Quote)
	q, ok := m[exchange+"|"+symbol]
	return q, ok
}
