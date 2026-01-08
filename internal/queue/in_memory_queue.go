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
	return &MemoryStore{
		batch: make(chan *models.MarketData, 10_000),
	}
}

func (s *MemoryStore) Run(in <-chan *models.MarketData) {
	const batchSize = 100
	buffer := make([]*models.MarketData, 0, batchSize)

	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	flush := func() {
		for _, md := range buffer {
			key := md.Exchange + "|" + md.Symbol
			s.m.Store(key, Quote{
				Bid: md.Bid, Ask: md.Ask,
				BidSize: md.BidSize, AskSize: md.AskSize,
				Timestamp: md.Timestamp,
			})
		}
		buffer = buffer[:0]
	}

	for {
		select {
		case md, ok := <-in:
			if !ok {
				flush()
				return
			}
			buffer = append(buffer, md)
			if len(buffer) >= batchSize {
				flush()
			}
		case <-ticker.C:
			if len(buffer) > 0 {
				flush()
			}
		}
	}
}

// Snapshot возвращает актуальные котировки
func (s *MemoryStore) Snapshot() map[string]Quote {
	snapshot := make(map[string]Quote)
	s.m.Range(func(k, v any) bool {
		snapshot[k.(string)] = v.(Quote)
		return true
	})
	return snapshot
}
func (s *MemoryStore) Get(exchange, symbol string) (Quote, bool) {
	m := s.data.Load().(map[string]Quote)
	q, ok := m[exchange+"|"+symbol]
	return q, ok
}
