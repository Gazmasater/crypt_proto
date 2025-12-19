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




export TRADE_AMOUNT_USDT=100
export FEE_PCT=0.04
export SELL_SAFETY=0.995

export TRIANGLES_FILE=triangles_markets.csv
export TRIANGLES_ENRICHED_FILE=triangles_markets_enriched.csv

go run ./cmd/triangles_enrich_mexc



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

		consumer.Executor = arb.NewRealExecutor(tr, arbOut, startUSDT)
		log.Printf("Executor: REAL (startUSDT=%.6f) NOTE: BUY uses quoteQty; if triangle has BUY with quote!=USDT, extend SpotTrader with SmartMarketBuyQuote",
			startUSDT,
		)
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





# ======================
# EXCHANGE
# ======================
EXCHANGE=MEXC
TRIANGLES_FILE=triangles_usdt_routes_market.csv
BOOK_INTERVAL=10ms

# ======================
# ARBITRAGE LOGIC
# ======================
# комиссия одной ноги (%)
# 0.04 = 0.04%
FEE_PCT=0.04

# минимальная прибыль треугольника (%)
# Для 35 USDT лучше 1.5%, иначе округление/проскальзывание часто съедает профит
MIN_PROFIT_PCT=1.5

# минимальный старт (в USDT)
MIN_START_USDT=35

# доля от maxStart (1 = брать максимум)
# На 35 USDT безопаснее меньше, чтобы не ловить проскальзывание
START_FRACTION=0.3

# ======================
# TRADING
# ======================
TRADE_ENABLED=true
TRADE_AMOUNT_USDT=35

# анти-флуд между сделками (мс)
TRADE_COOLDOWN_MS=2000

# ======================
# MEXC API
# ======================
MEXC_API_KEY=ВАШ_РЕАЛЬНЫЙ_KEY
MEXC_API_SECRET=ВАШ_РЕАЛЬНЫЙ_SECRET

# ======================
# DEBUG
# ======================
DEBUG=false


