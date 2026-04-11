package queue

import (
	"sync"
	"time"

	"crypt_proto/pkg/models"
)

const defaultHistoryTTL = 3 * time.Second

type Quote struct {
	Bid, Ask         float64
	BidSize, AskSize float64
	Timestamp        int64
}

type symbolBuffer struct {
	items []Quote
}

type MemoryStore struct {
	mu      sync.RWMutex
	latest  map[string]Quote
	history map[string]*symbolBuffer
	ttlMS   int64
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		latest:  make(map[string]Quote),
		history: make(map[string]*symbolBuffer),
		ttlMS:   defaultHistoryTTL.Milliseconds(),
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

// Get возвращает последнюю известную котировку по символу.
func (s *MemoryStore) Get(exchange, symbol string) (Quote, bool) {
	s.mu.RLock()
	q, ok := s.latest[exchange+"|"+symbol]
	s.mu.RUnlock()
	return q, ok
}

// GetLatestBefore возвращает самую свежую котировку не позже ts и не старше maxAgeMS.
func (s *MemoryStore) GetLatestBefore(exchange, symbol string, ts int64, maxAgeMS int64) (Quote, bool) {
	if ts <= 0 {
		ts = time.Now().UnixMilli()
	}

	key := exchange + "|" + symbol
	cutoff := int64(0)
	if maxAgeMS > 0 {
		cutoff = ts - maxAgeMS
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	buf, ok := s.history[key]
	if !ok || len(buf.items) == 0 {
		return Quote{}, false
	}

	for i := len(buf.items) - 1; i >= 0; i-- {
		q := buf.items[i]
		if q.Timestamp > ts {
			continue
		}
		if maxAgeMS > 0 && q.Timestamp < cutoff {
			break
		}
		return q, true
	}

	return Quote{}, false
}

// apply обновляет тикер и добавляет его в короткую историю по символу.
func (s *MemoryStore) apply(md *models.MarketData) {
	key := md.Exchange + "|" + md.Symbol
	timestamp := md.Timestamp
	if timestamp <= 0 {
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
	defer s.mu.Unlock()

	if prev, ok := s.latest[key]; !ok || quote.Timestamp >= prev.Timestamp {
		s.latest[key] = quote
	}

	buf, ok := s.history[key]
	if !ok {
		buf = &symbolBuffer{}
		s.history[key] = buf
	}
	buf.items = append(buf.items, quote)

	cutoff := timestamp - s.ttlMS
	n := 0
	for _, item := range buf.items {
		if item.Timestamp >= cutoff {
			buf.items[n] = item
			n++
		}
	}
	buf.items = buf.items[:n]
}
