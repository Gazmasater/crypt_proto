package collector

import "sort"

type BookLevel struct {
	Price float64
	Size  float64
}

type BookSnapshot struct {
	Symbol      string
	Bids        []BookLevel
	Asks        []BookLevel
	Sequence    int64
	TimestampMS int64
}

type DepthBookSource interface {
	GetBookSnapshot(symbol string, depth int) (BookSnapshot, bool)
}

func buildSnapshot(symbol string, book *bookState, depth int) (BookSnapshot, bool) {
	if book == nil {
		return BookSnapshot{}, false
	}

	book.mu.Lock()
	defer book.mu.Unlock()

	if !book.ready || book.bestBid <= 0 || book.bestAsk <= 0 {
		return BookSnapshot{}, false
	}

	bids := make([]BookLevel, 0, len(book.bids))
	for p, s := range book.bids {
		if p > 0 && s > 0 {
			bids = append(bids, BookLevel{Price: p, Size: s})
		}
	}

	asks := make([]BookLevel, 0, len(book.asks))
	for p, s := range book.asks {
		if p > 0 && s > 0 {
			asks = append(asks, BookLevel{Price: p, Size: s})
		}
	}

	if len(bids) == 0 || len(asks) == 0 {
		return BookSnapshot{}, false
	}

	sort.Slice(bids, func(i, j int) bool { return bids[i].Price > bids[j].Price })
	sort.Slice(asks, func(i, j int) bool { return asks[i].Price < asks[j].Price })

	if depth > 0 {
		if len(bids) > depth {
			bids = bids[:depth]
		}
		if len(asks) > depth {
			asks = asks[:depth]
		}
	}

	return BookSnapshot{
		Symbol:      symbol,
		Bids:        bids,
		Asks:        asks,
		Sequence:    book.sequence,
		TimestampMS: book.lastEventTS,
	}, true
}
