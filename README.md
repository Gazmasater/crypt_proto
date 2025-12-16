mx0vglmT3srN1IS19H
135bb7a7509e4421bad692415c53753b



sudo systemctl mask sleep.target suspend.target hibernate.target hybrid-sleep.target



wbs-api.mexc.com/ws 


[https://edis-global.vercel.app/ru/vps-hosting/singapore-singapore
](https://sg.edisglobal.com/)



git pull --rebase origin privat
git push origin privat


BOOK_INTERVAL=100ms
SYMBOLS_FILE=triangles_markets.csv
DEBUG=false


import (
    // ...
    "net/http"
    _ "net/http/pprof"
)


   // pprof HTTP-сервер
    go func() {
        log.Println("pprof on http://localhost:6060/debug/pprof/")
        if err := http.ListenAndServe("localhost:6060", nil); err != nil {
            log.Printf("pprof server error: %v", err)
        }
    }()


	go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30


(pprof) top        # показать топ функций по CPU
(pprof) top10
(pprof) list parsePBWrapperMid   # подробный разбор одной функции
(pprof) quit


go tool pprof http://localhost:6060/debug/pprof/heap


(pprof) top
(pprof) top -cum
(pprof) list parsePBWrapperMid
(pprof) quit





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

	// Trading toggles (если в твоём Config этих полей нет — добавь их в config.go)
	consumer.TradeEnabled = cfg.TradeEnabled
	consumer.TradeAmountUSDT = cfg.TradeAmountUSDT
	consumer.TradeCooldown = time.Duration(cfg.TradeCooldownMs) * time.Millisecond

	// Executor
	if cfg.Exchange == "MEXC" && cfg.TradeEnabled && cfg.APIKey != "" && cfg.APISecret != "" {
		tr := mexc.NewTrader(cfg.APIKey, cfg.APISecret, cfg.Debug)

		// startUSDT — фиксированная сумма на сделку (например 2)
		startUSDT := cfg.TradeAmountUSDT
		if startUSDT <= 0 {
			startUSDT = 2.0
		}

		consumer.Executor = arb.NewRealExecutor(tr, arbOut, startUSDT)
		log.Printf("Executor: REAL (startUSDT=%.6f)", startUSDT)
	} else {
		consumer.Executor = arb.NewNoopExecutor() // безопаснее, чем несуществующий DryRun
		log.Printf("Executor: NOOP (trade disabled or missing keys)")
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



arb/executor_noop.go (новый файл)
package arb

import (
	"context"

	"crypt_proto/domain"
)

// NoopExecutor — безопасный исполнитель "ничего не делаю".
// Используется когда трейд выключен или нет ключей API.
type NoopExecutor struct{}

func NewNoopExecutor() *NoopExecutor { return &NoopExecutor{} }

func (e *NoopExecutor) Name() string { return "NOOP" }

func (e *NoopExecutor) Execute(ctx context.Context, t domain.Triangle, quotes map[string]domain.Quote, startUSDT float64) error {
	// намеренно ничего не делаем
	_ = ctx
	_ = t
	_ = quotes
	_ = startUSDT
	return nil
}




[ARB] +0.140%  USDT→DOLO→USDC→USDT  maxStart=660.1859 USDT (660.1859 USDT)  safeStart=660.1859 USDT (660.1859 USDT) (x1.00)  bottleneck=DOLOUSDC
  DOLOUSDT (DOLO/USDT): bid=0.0380900000 ask=0.0384100000  spread=0.0003200000 (0.83660%)  bidQty=20892.2400 askQty=23100.6400
  DOLOUSDC (DOLO/USDC): bid=0.0385100000 ask=0.0388000000  spread=0.0002900000 (0.75023%)  bidQty=17180.9900 askQty=18468.0300
  USDCUSDT (USDC/USDT): bid=1.0000000000 ask=1.0001000000  spread=0.0001000000 (0.01000%)  bidQty=145472.9500 askQty=27025.4700

  [REAL EXEC] start=2.000000 USDT triangle=USDT→DOLO→USDC→USDT
    [REAL EXEC] leg 1: BUY DOLOUSDT by USDT=2.000000 (ask=0.0384100000)






