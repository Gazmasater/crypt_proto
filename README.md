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



package arb

import (
	"context"
	"fmt"
	"io"
	"math"
	"strings"
	"sync"
	"time"

	"crypt_proto/domain"
)

type SpotTrader interface {
	SmartMarketBuyUSDT(ctx context.Context, symbol string, usdt float64, ask float64) (string, error)
	SmartMarketSellQty(ctx context.Context, symbol string, qty float64) (string, error)
	GetBalance(ctx context.Context, asset string) (float64, error)
}

type execReq struct {
	ctx       context.Context
	t         domain.Triangle
	quotes    map[string]domain.Quote // snapshot только нужных символов
	startUSDT float64
	triName   string
}

type RealExecutor struct {
	trader SpotTrader
	out    io.Writer

	StartUSDT  float64
	SellSafety float64

	// cooldown по треугольнику (по имени)
	Cooldown time.Duration

	mu       sync.Mutex
	lastExec map[string]time.Time

	// Очередь (строго последовательное исполнение)
	queue chan execReq
	wg    sync.WaitGroup

	// STOP logic
	StopAfterOne bool // остановить программу после первого треугольника (успех или ошибка)
	stopOnce     sync.Once
	onStop       func()
}

func NewRealExecutor(tr SpotTrader, out io.Writer, startUSDT float64) *RealExecutor {
	e := &RealExecutor{
		trader:     tr,
		out:        out,
		StartUSDT:  startUSDT,
		SellSafety: 0.995,
		Cooldown:   500 * time.Millisecond,
		lastExec:   make(map[string]time.Time),

		// буфер можно увеличить, но лучше небольшой, чтобы не копить “устаревшие” сделки
		queue: make(chan execReq, 16),
	}

	// worker: исполняет строго по одному
	e.wg.Add(1)
	go func() {
		defer e.wg.Done()
		for req := range e.queue {
			if err := e.executeOnce(req); err != nil {
				// ВАЖНО: логируем ошибки исполнения здесь тоже, чтобы не потерять их
				e.logf("  [REAL EXEC] EXEC ERROR: triangle=%s err=%v", req.triName, err)
			}
		}
	}()

	return e
}

func (e *RealExecutor) Name() string { return "REAL" }

type flusher interface{ Flush() error }

func (e *RealExecutor) logf(format string, args ...any) {
	fmt.Fprintf(e.out, format+"\n", args...)
	if f, ok := e.out.(flusher); ok {
		_ = f.Flush()
	}
}

func (e *RealExecutor) step(name string) func() {
	start := time.Now()
	e.logf("    [REAL EXEC] >>> %s", name)
	return func() {
		e.logf("    [REAL EXEC] <<< %s (%s)", name, time.Since(start).Truncate(time.Millisecond))
	}
}

// SetStopFunc задаёт функцию остановки (обычно cancel() из main)
func (e *RealExecutor) SetStopFunc(fn func()) {
	e.onStop = fn
}

func (e *RealExecutor) requestStop(reason string) {
	if !e.StopAfterOne || e.onStop == nil {
		return
	}
	e.stopOnce.Do(func() {
		e.logf("  [REAL EXEC] STOP_AFTER_ONE: %s", reason)
		e.onStop()
	})
}

// Execute теперь НЕ исполняет сразу.
// Он кладёт треугольник в очередь со снапшотом котировок и возвращает.
func (e *RealExecutor) Execute(ctx context.Context, t domain.Triangle, quotes map[string]domain.Quote, startUSDT float64) error {
	triName := strings.TrimSpace(t.Name)
	if triName == "" {
		triName = "triangle"
	}

	if startUSDT <= 0 {
		startUSDT = e.StartUSDT
	}
	if startUSDT <= 0 {
		return fmt.Errorf("startUSDT<=0 (startUSDT=%.6f, StartUSDT=%.6f)", startUSDT, e.StartUSDT)
	}

	// 3 символа из ног
	if len(t.Legs) < 3 {
		return fmt.Errorf("triangle %s has <3 legs", triName)
	}
	sym1 := strings.TrimSpace(t.Legs[0].Symbol)
	sym2 := strings.TrimSpace(t.Legs[1].Symbol)
	sym3 := strings.TrimSpace(t.Legs[2].Symbol)
	if sym1 == "" || sym2 == "" || sym3 == "" {
		return fmt.Errorf("triangle %s has empty leg symbols: [%q, %q, %q]", triName, sym1, sym2, sym3)
	}

	// СНАПШОТ котировок только по нужным символам
	snap := make(map[string]domain.Quote, 3)
	if q, ok := quotes[sym1]; ok {
		snap[sym1] = q
	}
	if q, ok := quotes[sym2]; ok {
		snap[sym2] = q
	}
	if q, ok := quotes[sym3]; ok {
		snap[sym3] = q
	}

	req := execReq{
		ctx:       ctx,
		t:         t,
		quotes:    snap,
		startUSDT: startUSDT,
		triName:   triName,
	}

	select {
	case e.queue <- req:
		e.logf("  [REAL EXEC] QUEUED: start=%.6f USDT triangle=%s", startUSDT, triName)
		return nil
	default:
		e.logf("  [REAL EXEC] SKIP: queue full (triangle=%s)", triName)
		return nil
	}
}

func (e *RealExecutor) executeOnce(req execReq) (retErr error) {
	now := time.Now()

	// cooldown по имени треугольника
	e.mu.Lock()
	if last, ok := e.lastExec[req.triName]; ok && e.Cooldown > 0 && now.Sub(last) < e.Cooldown {
		left := (e.Cooldown - now.Sub(last)).Truncate(time.Millisecond)
		e.mu.Unlock()
		e.logf("  [REAL EXEC] SKIP cooldown triangle=%s left=%s", req.triName, left)
		return nil
	}
	e.mu.Unlock()

	triName := req.triName
	t := req.t
	quotes := req.quotes
	startUSDT := req.startUSDT

	// Обязательная остановка после первого треугольника:
	// - при успехе: остановить
	// - при любой ошибке: тоже остановить
	defer func() {
		if retErr != nil {
			e.logf("  [REAL EXEC] TRIANGLE FAILED: %s err=%v", triName, retErr)
			e.requestStop(fmt.Sprintf("triangle=%s failed: %v", triName, retErr))
			return
		}
		e.logf("  [REAL EXEC] TRIANGLE SUCCESS: %s", triName)
		e.requestStop(fmt.Sprintf("triangle=%s done (stop after one)", triName))
	}()

	// 3 символа из ног
	sym1 := strings.TrimSpace(t.Legs[0].Symbol)
	sym2 := strings.TrimSpace(t.Legs[1].Symbol)
	sym3 := strings.TrimSpace(t.Legs[2].Symbol)

	e.logf("  [REAL EXEC] start=%.6f USDT triangle=%s", startUSDT, triName)
	e.logf("    [REAL EXEC] legs: sym1=%s sym2=%s sym3=%s", sym1, sym2, sym3)

	// base/quote для логов
	base1, quote1 := parseBaseQuote(sym1)
	base2, quote2 := parseBaseQuote(sym2)
	base3, quote3 := parseBaseQuote(sym3)
	e.logf("    [REAL EXEC] parsed: sym1=%s (%s/%s) sym2=%s (%s/%s) sym3=%s (%s/%s)",
		sym1, base1, quote1,
		sym2, base2, quote2,
		sym3, base3, quote3,
	)

	// ===== balances before =====
	var usdt0 float64
	{
		done := e.step("GetBalance USDT (before)")
		v, err := e.trader.GetBalance(req.ctx, "USDT")
		done()
		if err != nil {
			e.logf("    [REAL EXEC] BAL ERR: get USDT before: %v", err)
			return err
		}
		usdt0 = v
		e.logf("    [REAL EXEC] BAL before: USDT=%.12f", usdt0)
		if usdt0+1e-9 < startUSDT {
			return fmt.Errorf("insufficient USDT: have=%.12f need=%.12f", usdt0, startUSDT)
		}
	}

	// ===== LEG 1: BUY sym1 by USDT =====
	q1, ok := quotes[sym1]
	if !ok {
		return fmt.Errorf("no quote snapshot for sym1=%s", sym1)
	}

	var aBefore1 float64
	{
		done := e.step(fmt.Sprintf("GetBalance %s (before leg1)", base1))
		v, err := e.trader.GetBalance(req.ctx, base1)
		done()
		if err != nil {
			e.logf("    [REAL EXEC] BAL ERR: get %s before leg1: %v", base1, err)
			return err
		}
		aBefore1 = v
	}

	e.logf("    [REAL EXEC] leg1 PRE: BUY %s by %s=%.6f ask=%.10f bid=%.10f | %s before=%.12f",
		sym1, quote1, startUSDT, q1.Ask, q1.Bid, base1, aBefore1)

	var ord1 string
	{
		orderCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		done := e.step("SmartMarketBuyUSDT leg1")
		id, err := e.trader.SmartMarketBuyUSDT(orderCtx, sym1, startUSDT, q1.Ask)
		done()
		if err != nil {
			e.logf("    [REAL EXEC] leg1 PLACE ERR: %v", err)
			return err
		}
		ord1 = id
	}
	e.logf("    [REAL EXEC] leg1 PLACE OK: orderId=%s", ord1)

	var aAfter1 float64
	{
		done := e.step(fmt.Sprintf("waitBalanceChange %s (after leg1)", base1))
		v, err := e.waitBalanceChange(req.ctx, base1, aBefore1, 3*time.Second, 150*time.Millisecond)
		done()
		if err != nil {
			e.logf("    [REAL EXEC] leg1 WAIT BAL ERR (%s): %v", base1, err)
			return err
		}
		aAfter1 = v
	}
	dA := aAfter1 - aBefore1
	e.logf("    [REAL EXEC] leg1 BAL after: %s=%.12f delta=%.12f", base1, aAfter1, dA)
	if dA <= 0 {
		return fmt.Errorf("leg1: %s did not increase (before=%.12f after=%.12f)", base1, aBefore1, aAfter1)
	}

	// ===== LEG 2: SELL sym2 =====
	q2, ok := quotes[sym2]
	if !ok {
		return fmt.Errorf("no quote snapshot for sym2=%s", sym2)
	}

	var base2Bal float64
	{
		done := e.step(fmt.Sprintf("GetBalance %s (before leg2)", base2))
		v, err := e.trader.GetBalance(req.ctx, base2)
		done()
		if err != nil {
			e.logf("    [REAL EXEC] BAL ERR: get %s before leg2: %v", base2, err)
			return err
		}
		base2Bal = v
	}

	sellA := base2Bal * e.SellSafety
	if sellA <= 0 {
		return fmt.Errorf("leg2: sell qty <=0 (%s=%.12f safety=%.6f)", base2, base2Bal, e.SellSafety)
	}

	var bBefore2 float64
	{
		done := e.step(fmt.Sprintf("GetBalance %s (before leg2)", quote2))
		v, err := e.trader.GetBalance(req.ctx, quote2)
		done()
		if err != nil {
			e.logf("    [REAL EXEC] BAL ERR: get %s before leg2: %v", quote2, err)
			return err
		}
		bBefore2 = v
	}

	e.logf("    [REAL EXEC] leg2 PRE: SELL %s qty=%s=%.12f (safety x%.6f) bid=%.10f ask=%.10f | %s before=%.12f %s before=%.12f",
		sym2, base2, sellA, e.SellSafety, q2.Bid, q2.Ask,
		base2, base2Bal,
		quote2, bBefore2,
	)

	var ord2 string
	{
		orderCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		done := e.step("SmartMarketSellQty leg2")
		id, err := e.trader.SmartMarketSellQty(orderCtx, sym2, sellA)
		done()
		if err != nil {
			e.logf("    [REAL EXEC] leg2 PLACE ERR: %v", err)
			return err
		}
		ord2 = id
	}
	e.logf("    [REAL EXEC] leg2 PLACE OK: orderId=%s", ord2)

	var bAfter2 float64
	{
		done := e.step(fmt.Sprintf("waitBalanceChange %s (after leg2)", quote2))
		v, err := e.waitBalanceChange(req.ctx, quote2, bBefore2, 3*time.Second, 150*time.Millisecond)
		done()
		if err != nil {
			e.logf("    [REAL EXEC] leg2 WAIT BAL ERR (%s): %v", quote2, err)
			return err
		}
		bAfter2 = v
	}
	dB := bAfter2 - bBefore2
	e.logf("    [REAL EXEC] leg2 BAL after: %s=%.12f delta=%.12f", quote2, bAfter2, dB)
	if dB <= 0 {
		return fmt.Errorf("leg2: %s did not increase (before=%.12f after=%.12f)", quote2, bBefore2, bAfter2)
	}

	// ===== LEG 3: SELL sym3 (base3 -> USDT) =====
	q3, ok := quotes[sym3]
	if !ok {
		return fmt.Errorf("no quote snapshot for sym3=%s", sym3)
	}

	var base3Bal float64
	{
		done := e.step(fmt.Sprintf("GetBalance %s (before leg3)", base3))
		v, err := e.trader.GetBalance(req.ctx, base3)
		done()
		if err != nil {
			e.logf("    [REAL EXEC] BAL ERR: get %s before leg3: %v", base3, err)
			return err
		}
		base3Bal = v
	}

	sellB := base3Bal * e.SellSafety
	if sellB <= 0 {
		return fmt.Errorf("leg3: sell qty <=0 (%s=%.12f safety=%.6f)", base3, base3Bal, e.SellSafety)
	}

	var usdtBefore3 float64
	{
		done := e.step("GetBalance USDT (before leg3)")
		v, err := e.trader.GetBalance(req.ctx, "USDT")
		done()
		if err != nil {
			e.logf("    [REAL EXEC] BAL ERR: get USDT before leg3: %v", err)
			return err
		}
		usdtBefore3 = v
	}

	e.logf("    [REAL EXEC] leg3 PRE: SELL %s qty=%s=%.12f (safety x%.6f) bid=%.10f ask=%.10f | %s before=%.12f USDT before=%.12f",
		sym3, base3, sellB, e.SellSafety, q3.Bid, q3.Ask,
		base3, base3Bal, usdtBefore3)

	var ord3 string
	{
		orderCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		done := e.step("SmartMarketSellQty leg3")
		id, err := e.trader.SmartMarketSellQty(orderCtx, sym3, sellB)
		done()
		if err != nil {
			e.logf("    [REAL EXEC] leg3 PLACE ERR: %v", err)
			return err
		}
		ord3 = id
	}
	e.logf("    [REAL EXEC] leg3 PLACE OK: orderId=%s", ord3)

	var usdtAfter float64
	{
		done := e.step("waitBalanceChange USDT (after leg3)")
		v, err := e.waitBalanceChange(req.ctx, "USDT", usdtBefore3, 3*time.Second, 150*time.Millisecond)
		done()
		if err != nil {
			e.logf("    [REAL EXEC] leg3 WAIT BAL ERR (USDT): %v", err)
			return err
		}
		usdtAfter = v
	}

	dUSDT3 := usdtAfter - usdtBefore3
	dUSDTTotal := usdtAfter - usdt0

	e.logf("    [REAL EXEC] leg3 BAL after: USDT=%.12f delta=%.12f", usdtAfter, dUSDT3)
	e.logf("    [REAL EXEC] DONE: USDT start=%.12f end=%.12f pnl(total)=%.12f (%.4f%%)",
		usdt0, usdtAfter, dUSDTTotal, pct(dUSDTTotal, startUSDT))

	e.mu.Lock()
	e.lastExec[triName] = time.Now()
	e.mu.Unlock()

	return nil
}

// waitBalanceChange ждёт, пока баланс станет отличаться от baseline.
func (e *RealExecutor) waitBalanceChange(ctx context.Context, asset string, baseline float64, timeout, interval time.Duration) (float64, error) {
	const tol = 1e-12

	deadline := time.NewTimer(timeout)
	tick := time.NewTicker(interval)
	defer deadline.Stop()
	defer tick.Stop()

	cur, err := e.trader.GetBalance(ctx, asset)
	if err == nil && math.Abs(cur-baseline) > tol {
		return cur, nil
	}

	for {
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		case <-deadline.C:
			last, err := e.trader.GetBalance(ctx, asset)
			if err != nil {
				return 0, fmt.Errorf("timeout, last balance read error for %s: %v", asset, err)
			}
			return 0, fmt.Errorf("timeout waiting %s balance change: baseline=%.12f last=%.12f", asset, baseline, last)
		case <-tick.C:
			cur, err := e.trader.GetBalance(ctx, asset)
			if err != nil {
				continue
			}
			if math.Abs(cur-baseline) > tol {
				return cur, nil
			}
		}
	}
}

// parseBaseQuote — простой парсер BASE/QUOTE по суффиксу.
func parseBaseQuote(symbol string) (base, quote string) {
	quotes := []string{"USDT", "USDC", "BTC", "ETH", "EUR", "TRY", "BRL", "RUB"}
	for _, q := range quotes {
		if strings.HasSuffix(symbol, q) && len(symbol) > len(q) {
			return symbol[:len(symbol)-len(q)], q
		}
	}
	return symbol, ""
}

func pct(delta, denom float64) float64 {
	if denom == 0 {
		return 0
	}
	return (delta / denom) * 100
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

