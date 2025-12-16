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
	go func() {
		log.Println("pprof on http://localhost:6060/debug/pprof/")
		_ = http.ListenAndServe("localhost:6060", nil)
	}()

	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	cfg := config.Load()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

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

	var feed exchange.MarketDataFeed
	switch cfg.Exchange {
	case "MEXC":
		feed = mexc.NewFeed(cfg.Debug)
	case "KUCOIN":
		feed = kucoin.NewFeed(cfg.Debug)
	default:
		log.Fatalf("unknown EXCHANGE=%q (expected MEXC or KUCOIN)", cfg.Exchange)
	}

	logFile, logBuf, arbOut := arb.OpenLogWriter("arbitrage.log")
	defer logFile.Close()
	defer logBuf.Flush()

	events := make(chan domain.Event, 8192)
	var wg sync.WaitGroup

	consumer := arb.NewConsumer(cfg.FeePerLeg, cfg.MinProfit, cfg.MinStart, arbOut)
	consumer.StartFraction = cfg.StartFraction

	consumer.TradeEnabled = cfg.TradeEnabled
	consumer.TradeAmountUSDT = cfg.TradeAmountUSDT
	consumer.TradeCooldown = time.Duration(cfg.TradeCooldownMs) * time.Millisecond

	// ---------------------------
	// Executor selection
	// ---------------------------
	useReal := cfg.Exchange == "MEXC" && cfg.TradeEnabled && cfg.APIKey != "" && cfg.APISecret != ""
	if useReal {
		tr := mexc.NewTrader(cfg.APIKey, cfg.APISecret, cfg.Debug)

		// ВАЖНО: сигнатура NewRealExecutor(tr,out,startUSDT float64)
		exec := arb.NewRealExecutor(tr, arbOut, cfg.TradeAmountUSDT)

		// Подхватываем cooldown из конфига (если поле есть в executor'е)
		// Если у тебя RealExecutor.Cooldown already есть — ок.
		exec.Cooldown = time.Duration(cfg.TradeCooldownMs) * time.Millisecond

		// Попробуем подтянуть filters (step/minQty/apiEnabled) и закинуть в exec.Filters
		// Это нужно для нормализации qty и чтобы заранее отсекать "symbol not support api" (10007).
		caps, err := mexc.FetchSymbolCapsMEXC(ctx)
		if err != nil {
			log.Printf("WARN: cannot fetch exchangeInfo caps: %v", err)
		} else {
			filters := make(map[string]arb.SymbolFilter, len(caps))
			for sym, c := range caps {
				filters[sym] = arb.SymbolFilter{
					StepSize:   c.StepSize,
					MinQty:     c.MinQty,
					APIEnabled: c.APIEnabled,
				}
			}
			// !!! ВАЖНО !!!
			// Чтобы эта строка компилилась, в RealExecutor должно быть поле:
			// Filters map[string]arb.SymbolFilter
			exec.Filters = filters

			log.Printf("[MEXC] loaded caps: symbols=%d", len(filters))
		}

		consumer.Executor = exec
		log.Printf("Executor: REAL (amount=%.6f USDT, cooldown=%v)", cfg.TradeAmountUSDT, consumer.TradeCooldown)
	} else {
		consumer.Executor = arb.NewDryRunExecutor(arbOut)
		log.Printf("Executor: DRY-RUN (trade disabled or no keys)")
	}

	consumer.Start(ctx, events, triangles, indexBySymbol, &wg)
	feed.Start(ctx, &wg, symbols, cfg.BookInterval, events)

	<-ctx.Done()
	log.Println("shutting down...")

	// НЕ закрываем events: WS горутины могут ещё писать
	wg.Wait()
	log.Println("bye")
}



[{
	"resource": "/home/gaz358/myprog/crypt_proto/cmd/cryptarb/main.go",
	"owner": "_generated_diagnostic_collection_name_#0",
	"code": {
		"value": "MissingLitField",
		"target": {
			"$mid": 1,
			"path": "/golang.org/x/tools/internal/typesinternal",
			"scheme": "https",
			"authority": "pkg.go.dev",
			"fragment": "MissingLitField"
		}
	},
	"severity": 8,
	"message": "unknown field APIEnabled in struct literal of type arb.SymbolFilter",
	"source": "compiler",
	"startLineNumber": 96,
	"startColumn": 6,
	"endLineNumber": 96,
	"endColumn": 16,
	"origin": "extHost1"
}]


[{
	"resource": "/home/gaz358/myprog/crypt_proto/cmd/cryptarb/main.go",
	"owner": "_generated_diagnostic_collection_name_#0",
	"code": {
		"value": "MissingFieldOrMethod",
		"target": {
			"$mid": 1,
			"path": "/golang.org/x/tools/internal/typesinternal",
			"scheme": "https",
			"authority": "pkg.go.dev",
			"fragment": "MissingFieldOrMethod"
		}
	},
	"severity": 8,
	"message": "c.APIEnabled undefined (type mexc.SymbolCaps has no field or method APIEnabled)",
	"source": "compiler",
	"startLineNumber": 96,
	"startColumn": 20,
	"endLineNumber": 96,
	"endColumn": 30,
	"origin": "extHost1"
}]

[{
	"resource": "/home/gaz358/myprog/crypt_proto/cmd/cryptarb/main.go",
	"owner": "_generated_diagnostic_collection_name_#0",
	"code": {
		"value": "MissingFieldOrMethod",
		"target": {
			"$mid": 1,
			"path": "/golang.org/x/tools/internal/typesinternal",
			"scheme": "https",
			"authority": "pkg.go.dev",
			"fragment": "MissingFieldOrMethod"
		}
	},
	"severity": 8,
	"message": "exec.Filters undefined (type *arb.RealExecutor has no field or method Filters)",
	"source": "compiler",
	"startLineNumber": 102,
	"startColumn": 9,
	"endLineNumber": 102,
	"endColumn": 16,
	"origin": "extHost1"
}]

[{
	"resource": "/home/gaz358/myprog/crypt_proto/arb/executor_real.go",
	"owner": "_generated_diagnostic_collection_name_#0",
	"code": {
		"value": "default",
		"target": {
			"$mid": 1,
			"path": "/docs/checks/",
			"scheme": "https",
			"authority": "staticcheck.dev",
			"fragment": "SA6005"
		}
	},
	"severity": 4,
	"message": "should use strings.EqualFold instead",
	"source": "SA6005",
	"startLineNumber": 137,
	"startColumn": 5,
	"endLineNumber": 137,
	"endColumn": 61,
	"origin": "extHost1"
}]










