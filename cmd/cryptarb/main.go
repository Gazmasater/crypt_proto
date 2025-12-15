package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"crypt_proto/arb"
	"crypt_proto/config"
	"crypt_proto/domain"
	"crypt_proto/exchange"
	"crypt_proto/kucoin"
	"crypt_proto/mexc"

	_ "net/http/pprof"
)

func main() {
	// pprof сервер
	go func() {
		log.Println("pprof on http://localhost:6060/debug/pprof/")
		if err := http.ListenAndServe("localhost:6060", nil); err != nil {
			log.Printf("pprof server error: %v", err)
		}
	}()

	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	// --------- конфиг ---------
	cfg := config.Load()

	// --------- треугольники ---------
	triangles, symbols, indexBySymbol, err := domain.LoadTriangles(cfg.TrianglesFile)
	if err != nil {
		log.Fatalf("load triangles: %v", err)
	}
	if len(triangles) == 0 {
		log.Fatal("нет треугольников, нечего мониторить")
	}
	if len(symbols) == 0 {
		log.Fatal("нет символов для подписки")
	}
	log.Printf("символов для подписки всего: %d", len(symbols))

	// --------- выбор биржи (market data feed) ---------
	var feed exchange.MarketDataFeed

	switch cfg.Exchange {
	case "MEXC":
		feed = mexc.NewFeed(cfg.Debug)
	case "KUCOIN":
		feed = kucoin.NewFeed(cfg.Debug)
	default:
		log.Fatalf("unknown EXCHANGE=%q (expected MEXC or KUCOIN)", cfg.Exchange)
	}

	log.Printf("Using exchange: %s", feed.Name())

	// --------- лог-файл для арбитража ---------
	logFile, logBuf, arbOut := arb.OpenLogWriter("arbitrage.log")
	defer logFile.Close()
	defer logBuf.Flush()

	// --------- контекст с остановкой по сигналу ---------
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	events := make(chan domain.Event, 8192)
	var wg sync.WaitGroup

	// --------- потребитель арбитража ---------
	consumer := arb.NewConsumer(cfg.FeePerLeg, cfg.MinProfit, cfg.MinStart, arbOut)
	consumer.StartFraction = cfg.StartFraction

	// --------- выбор исполнителя (DRY-RUN или REAL) ---------
	hasKeys := cfg.APIKey != "" && cfg.APISecret != ""

	if cfg.TradeEnabled && hasKeys {
		log.Printf("[EXEC] REAL TRADING ENABLED on %s", cfg.Exchange)

		switch cfg.Exchange {
		case "MEXC":
			trader := mexc.NewTrader(cfg.APIKey, cfg.APISecret, cfg.Debug)
			consumer.Executor = arb.NewRealExecutor(trader, arbOut)
		default:
			log.Printf("[EXEC] Real trading not implemented for %s, fallback to DRY-RUN", cfg.Exchange)
			consumer.Executor = arb.NewDryRunExecutor(arbOut)
		}
	} else {
		log.Printf(
			"[EXEC] DRY-RUN MODE (TRADE_ENABLED=%v, hasKeys=%v) — реальные ордера отправляться не будут",
			cfg.TradeEnabled, hasKeys,
		)
		consumer.Executor = arb.NewDryRunExecutor(arbOut)
	}

	// Стартуем потребителя
	consumer.Start(ctx, events, triangles, indexBySymbol, &wg)

	// --------- фид биржи (котировки) ---------
	feed.Start(ctx, &wg, symbols, cfg.BookInterval, events)

	// --------- ждём сигнал остановки ---------
	<-ctx.Done()
	log.Println("shutting down...")

	// немного времени на дообработку
	time.Sleep(200 * time.Millisecond)
	close(events)
	wg.Wait()
	log.Println("bye")
}
