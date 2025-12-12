package exchange

import (
	"context"
	"sync"

	"crypt_proto/domain"
)

// MarketDataFeed — общий интерфейс для биржи.
type MarketDataFeed interface {
	Name() string
	Start(ctx context.Context, wg *sync.WaitGroup, symbols []string, interval string, out chan<- domain.Event)
}
