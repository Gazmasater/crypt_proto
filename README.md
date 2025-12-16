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





package arb

import (
	"context"
	"fmt"
	"io"
	"math"
	"strings"
	"time"

	"crypt_proto/domain"
)

// SpotTrader должен УМЕТЬ внутри себя:
// - SmartMarketBuyUSDT: корректно сформировать MARKET BUY за фиксированный USDT
//   (на MEXC это часто либо quoteOrderQty, либо quantity — зависит от пары/правил;
//    нормализацию/точность делай внутри mexc.Trader)
// - SmartMarketSellQty: корректно SELL по quantity (нормализация/точность внутри)
// - GetBalance: читать free баланс (или доступный)
type SpotTrader interface {
	SmartMarketBuyUSDT(ctx context.Context, symbol string, usdt float64, ask float64) (string, error)
	SmartMarketSellQty(ctx context.Context, symbol string, qty float64) (string, error)
	GetBalance(ctx context.Context, asset string) (float64, error)
}

type RealExecutor struct {
	trader SpotTrader
	out    io.Writer

	// фиксированный старт в USDT (например 2)
	StartUSDT float64

	// safety чтобы не словить Oversold на SELL
	SellSafety float64

	// анти-спам: один и тот же треугольник не исполняем чаще Cooldown
	Cooldown time.Duration
	lastExec map[string]time.Time
}

func NewRealExecutor(tr SpotTrader, out io.Writer, startUSDT float64) *RealExecutor {
	return &RealExecutor{
		trader:     tr,
		out:        out,
		StartUSDT:  startUSDT,
		SellSafety: 0.995,
		Cooldown:   500 * time.Millisecond,
		lastExec:   make(map[string]time.Time),
	}
}

func (e *RealExecutor) Name() string { return "REAL" }

// Execute исполняет ТОЛЬКО безопасный класс треугольников:
// USDT -> A -> B -> USDT,
// где:
// leg1: BUY A за USDT
// leg2: SELL A в B
// leg3: SELL B в USDT
func (e *RealExecutor) Execute(ctx context.Context, t domain.Triangle, quotes map[string]domain.Quote, startUSDT float64) error {
	start := e.StartUSDT
	if startUSDT > 0 {
		start = startUSDT
	}
	if start <= 0 {
		return fmt.Errorf("start<=0")
	}

	// cooldown
	if e.Cooldown > 0 {
		if last, ok := e.lastExec[t.Name]; ok && time.Since(last) < e.Cooldown {
			return nil
		}
		e.lastExec[t.Name] = time.Now()
	}

	fmt.Fprintf(e.out, "  [REAL EXEC] start=%.6f USDT triangle=%s\n", start, t.Name)

	// Требуем строго USDT старт
	if !strings.EqualFold(t.Legs[0].From, "USDT") {
		return fmt.Errorf("unsupported triangle start asset: leg1.From=%s (expected USDT)", t.Legs[0].From)
	}

	// -----------------------------------
	// LEG 1: BUY A за USDT
	// -----------------------------------
	leg1 := t.Legs[0]

	// Должно быть направление USDT -> A
	if !strings.EqualFold(leg1.From, "USDT") {
		return fmt.Errorf("leg1 must start from USDT, got %s", leg1.From)
	}
	assetA := leg1.To
	if strings.EqualFold(assetA, "USDT") {
		return fmt.Errorf("leg1 invalid: To=USDT")
	}

	q1, ok := quotes[leg1.Symbol]
	if !ok || q1.Ask <= 0 {
		return fmt.Errorf("no quote/ask for %s", leg1.Symbol)
	}

	// Баланс A до покупки
	aBefore, err := e.trader.GetBalance(ctx, assetA)
	if err != nil {
		return fmt.Errorf("leg1 balance(before) error: %w", err)
	}

	fmt.Fprintf(e.out, "    [REAL EXEC] leg 1: BUY %s by USDT=%.6f (ask=%.10f)\n", leg1.Symbol, start, q1.Ask)
	if _, err := e.trader.SmartMarketBuyUSDT(ctx, leg1.Symbol, start, q1.Ask); err != nil {
		return fmt.Errorf("leg1 error: %w", err)
	}

	// Баланс A после покупки, берём дельту
	aAfter, err := e.trader.GetBalance(ctx, assetA)
	if err != nil {
		return fmt.Errorf("leg1 balance(after) error: %w", err)
	}
	aDelta := aAfter - aBefore
	if aDelta <= 0 {
		return fmt.Errorf("leg1 bought delta<=0: before=%.12f after=%.12f delta=%.12f %s", aBefore, aAfter, aDelta, assetA)
	}

	// -----------------------------------
	// LEG 2: SELL A -> B
	// -----------------------------------
	leg2 := t.Legs[1]

	// Требуем: leg2.From == A (мы должны иметь A)
	if !strings.EqualFold(leg2.From, assetA) {
		return fmt.Errorf("unsupported leg2 flow: have=%s need leg2.From=%s (triangle=%s)", assetA, leg2.From, t.Name)
	}
	assetB := leg2.To
	if strings.EqualFold(assetB, assetA) {
		return fmt.Errorf("leg2 invalid: To == From (%s)", assetB)
	}

	// Для текущего executor — leg2 должен быть SELL (Dir>0), чтобы можно было продать quantity
	// Если Dir<0 — это BUY base за quote (а у нас нет универсального buy для B!=USDT)
	if leg2.Dir <= 0 {
		return fmt.Errorf("leg2 requires BUY (Dir<0) which is not supported safely yet: %s", leg2.Symbol)
	}

	// Баланс B до
	bBefore, err := e.trader.GetBalance(ctx, assetB)
	if err != nil {
		return fmt.Errorf("leg2 balance(before) error: %w", err)
	}

	sellA := aDelta * e.SellSafety
	if sellA <= 0 {
		return fmt.Errorf("leg2 qty<=0 after safety (aDelta=%.12f)", aDelta)
	}

	fmt.Fprintf(e.out, "    [REAL EXEC] leg 2: SELL %s qty=%.12f (raw=%.12f)\n", leg2.Symbol, sellA, aDelta)
	if _, err := e.trader.SmartMarketSellQty(ctx, leg2.Symbol, sellA); err != nil {
		return fmt.Errorf("leg2 error: %w", err)
	}

	// Дельта B
	bAfter, err := e.trader.GetBalance(ctx, assetB)
	if err != nil {
		return fmt.Errorf("leg2 balance(after) error: %w", err)
	}
	bDelta := bAfter - bBefore
	if bDelta <= 0 {
		return fmt.Errorf("leg2 got delta<=0: before=%.12f after=%.12f delta=%.12f %s", bBefore, bAfter, bDelta, assetB)
	}

	// -----------------------------------
	// LEG 3: SELL B -> USDT
	// -----------------------------------
	leg3 := t.Legs[2]

	// Требуем: leg3.From == B и leg3.To == USDT
	if !strings.EqualFold(leg3.From, assetB) {
		return fmt.Errorf("unsupported leg3 flow: have=%s need leg3.From=%s (triangle=%s)", assetB, leg3.From, t.Name)
	}
	if !strings.EqualFold(leg3.To, "USDT") {
		return fmt.Errorf("leg3 must end in USDT, got %s", leg3.To)
	}
	if leg3.Dir <= 0 {
		return fmt.Errorf("leg3 requires BUY (Dir<0) which is not supported safely yet: %s", leg3.Symbol)
	}

	usdtBefore, _ := e.trader.GetBalance(ctx, "USDT")

	sellB := bDelta * e.SellSafety
	if sellB <= 0 {
		return fmt.Errorf("leg3 qty<=0 after safety (bDelta=%.12f)", bDelta)
	}

	fmt.Fprintf(e.out, "    [REAL EXEC] leg 3: SELL %s qty=%.12f (raw=%.12f)\n", leg3.Symbol, sellB, bDelta)
	if _, err := e.trader.SmartMarketSellQty(ctx, leg3.Symbol, sellB); err != nil {
		return fmt.Errorf("leg3 error: %w", err)
	}

	usdtAfter, _ := e.trader.GetBalance(ctx, "USDT")
	usdtDelta := usdtAfter - usdtBefore

	fmt.Fprintf(
		e.out,
		"  [REAL EXEC] done triangle %s  USDT_before=%.6f USDT_after=%.6f delta=%.6f\n",
		t.Name,
		math.Max(0, usdtBefore),
		math.Max(0, usdtAfter),
		usdtDelta,
	)

	return nil
}



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










