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
	// pprof
	go func() {
		log.Println("pprof on http://localhost:6060/debug/pprof/")
		_ = http.ListenAndServe("localhost:6060", nil)
	}()

	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	cfg := config.Load()

	// Context / signals
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// Triangles
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
	log.Printf("треугольников: %d", len(triangles))
	log.Printf("символов для подписки: %d", len(symbols))

	// Exchange feed
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

	// Log output
	logFile, logBuf, arbOut := arb.OpenLogWriter("arbitrage.log")
	defer logFile.Close()
	defer logBuf.Flush()

	// Events channel
	events := make(chan domain.Event, 8192)

	var wg sync.WaitGroup

	// Consumer
	consumer := arb.NewConsumer(cfg.FeePerLeg, cfg.MinProfit, cfg.MinStart, arbOut)
	consumer.StartFraction = cfg.StartFraction

	// Trading toggles
	consumer.TradeEnabled = cfg.TradeEnabled
	consumer.TradeAmountUSDT = cfg.TradeAmountUSDT
	consumer.TradeCooldown = time.Duration(cfg.TradeCooldownMs) * time.Millisecond

	log.Printf(
		"TRADE: enabled=%v amountUSDT=%.6f cooldown=%s feePerLeg=%.6f minProfit=%.6f minStart=%.6f startFraction=%.4f exchange=%s debug=%v",
		consumer.TradeEnabled,
		consumer.TradeAmountUSDT,
		consumer.TradeCooldown,
		cfg.FeePerLeg,
		cfg.MinProfit,
		cfg.MinStart,
		consumer.StartFraction,
		cfg.Exchange,
		cfg.Debug,
	)

	// Executor
	if cfg.Exchange == "MEXC" && cfg.TradeEnabled && cfg.APIKey != "" && cfg.APISecret != "" {
		tr := mexc.NewTrader(cfg.APIKey, cfg.APISecret, cfg.Debug)

		// startUSDT — фиксированная сумма на сделку
		startUSDT := cfg.TradeAmountUSDT
		if startUSDT <= 0 {
			startUSDT = 35.0
		}

		re := arb.NewRealExecutor(tr, arbOut, startUSDT)
		re.StopAfterOne = true
		re.SetStopFunc(cancel)

		consumer.Executor = re
		log.Printf("Executor: REAL (startUSDT=%.6f) STOP_AFTER_ONE=true", startUSDT)
	} else {
		consumer.Executor = arb.NewNoopExecutor()
		log.Printf("Executor: NOOP (trade disabled, non-MEXC exchange, or missing keys)")
	}

	// Start consumer
	consumer.Start(ctx, events, triangles, indexBySymbol, &wg)

	// Start feed
	feed.Start(ctx, &wg, symbols, cfg.BookInterval, events)

	// Wait stop
	<-ctx.Done()
	log.Println("shutting down...")

	// ВАЖНО: events не закрываем — WS-горутин(ы) могут ещё писать и словить panic
	wg.Wait()
	log.Println("bye")
}
