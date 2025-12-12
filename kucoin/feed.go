package kucoin

import (
	"context"
	"log"
	"sync"

	"crypt_proto/domain"
)

type Feed struct {
	debug bool
}

func NewFeed(debug bool) *Feed { return &Feed{debug: debug} }

func (f *Feed) Name() string { return "KuCoin" }

func (f *Feed) Start(
	ctx context.Context,
	wg *sync.WaitGroup,
	symbols []string,
	interval string,
	out chan<- domain.Event,
) {
	// TODO: здесь нужно реализовать KuCoin WebSocket подписку
	// Сейчас просто логируем, чтобы не падало.
	log.Printf("[KuCoin] Start called, but not implemented yet. symbols=%d interval=%s", len(symbols), interval)
}
