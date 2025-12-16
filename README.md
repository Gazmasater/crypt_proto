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
	"sync"
	"time"

	"crypt_proto/domain"
)

type SpotTrader interface {
	SmartMarketBuyUSDT(ctx context.Context, symbol string, usdt float64, ask float64) (string, error)
	SmartMarketSellQty(ctx context.Context, symbol string, qty float64) (string, error)
	GetBalance(ctx context.Context, asset string) (float64, error)
}

type RealExecutor struct {
	trader SpotTrader
	out    io.Writer

	StartUSDT   float64
	SellSafety float64

	Cooldown time.Duration
	lastExec map[string]time.Time

	// чтобы не исполнять несколько треугольников одновременно
	mu       sync.Mutex
	inFlight bool
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

func (e *RealExecutor) Execute(ctx context.Context, t domain.Triangle, quotes map[string]domain.Quote, startUSDT float64) error {
	// не даём запускать real-exec параллельно
	e.mu.Lock()
	if e.inFlight {
		e.mu.Unlock()
		e.logf("  [REAL EXEC] SKIP: inFlight")
		return nil
	}
	e.inFlight = true
	e.mu.Unlock()
	defer func() {
		e.mu.Lock()
		e.inFlight = false
		e.mu.Unlock()
	}()

	now := time.Now()

	triName := strings.TrimSpace(t.Name)
	if triName == "" {
		triName = "triangle"
	}

	// cooldown по имени треугольника
	if last, ok := e.lastExec[triName]; ok && e.Cooldown > 0 {
		if now.Sub(last) < e.Cooldown {
			e.logf("  [REAL EXEC] SKIP cooldown triangle=%s left=%s",
				triName, (e.Cooldown-now.Sub(last)).Truncate(time.Millisecond))
			return nil
		}
	}

	if startUSDT <= 0 {
		startUSDT = e.StartUSDT
	}
	if startUSDT <= 0 {
		return fmt.Errorf("startUSDT<=0 (startUSDT=%.6f, StartUSDT=%.6f)", startUSDT, e.StartUSDT)
	}

	// 3 символа из ног
	sym1 := strings.TrimSpace(t.Legs[0].Symbol)
	sym2 := strings.TrimSpace(t.Legs[1].Symbol)
	sym3 := strings.TrimSpace(t.Legs[2].Symbol)
	if sym1 == "" || sym2 == "" || sym3 == "" {
		return fmt.Errorf("triangle %s has empty leg symbols: [%q, %q, %q]", triName, sym1, sym2, sym3)
	}

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
		v, err := e.trader.GetBalance(ctx, "USDT")
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
		return fmt.Errorf("no quote for sym1=%s", sym1)
	}

	var aBefore1 float64
	{
		done := e.step(fmt.Sprintf("GetBalance %s (before leg1)", base1))
		v, err := e.trader.GetBalance(ctx, base1)
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
		// отдельный контекст на ордер — не падаем от Ctrl+C сразу
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
		v, err := e.waitBalanceChange(ctx, base1, aBefore1, 3*time.Second, 150*time.Millisecond)
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

	// ===== LEG 2: SELL sym2 (SELL base2 -> quote2) =====
	q2, ok := quotes[sym2]
	if !ok {
		return fmt.Errorf("no quote for sym2=%s", sym2)
	}

	var base2Bal float64
	{
		done := e.step(fmt.Sprintf("GetBalance %s (before leg2)", base2))
		v, err := e.trader.GetBalance(ctx, base2)
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
		v, err := e.trader.GetBalance(ctx, quote2)
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
		v, err := e.waitBalanceChange(ctx, quote2, bBefore2, 3*time.Second, 150*time.Millisecond)
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

	// ===== LEG 3: SELL sym3 (SELL base3 -> USDT) =====
	q3, ok := quotes[sym3]
	if !ok {
		return fmt.Errorf("no quote for sym3=%s", sym3)
	}

	var base3Bal float64
	{
		done := e.step(fmt.Sprintf("GetBalance %s (before leg3)", base3))
		v, err := e.trader.GetBalance(ctx, base3)
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
		v, err := e.trader.GetBalance(ctx, "USDT")
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
		v, err := e.waitBalanceChange(ctx, "USDT", usdtBefore3, 3*time.Second, 150*time.Millisecond)
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

	e.lastExec[triName] = now
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


[ARB] -0.098%  USDT→ATOM→BTC→USDT  maxStart=16.8000 USDT (16.8000 USDT)  safeStart=16.8000 USDT (16.8000 USDT) (x1.00)  bottleneck=ATOMBTC
  ATOMUSDT (ATOM/USDT): bid=2.0150000000 ask=2.0160000000  spread=0.0010000000 (0.04962%)  bidQty=191.1000 askQty=23.1700
  ATOMBTC (ATOM/BTC): bid=0.0000229700 ask=0.0000230100  spread=0.0000000400 (0.17399%)  bidQty=8.3300 askQty=11.4300
  BTCUSDT (BTC/USDT): bid=87786.2900000000 ask=87786.8600000000  spread=0.5700000000 (0.00065%)  bidQty=0.0026 askQty=0.0024

  [REAL EXEC] start=2.000000 USDT triangle=USDT→ATOM→BTC→USDT
    [REAL EXEC] legs: sym1=ATOMUSDT sym2=ATOMBTC sym3=BTCUSDT
    [REAL EXEC] parsed: sym1=ATOMUSDT (ATOM/USDT) sym2=ATOMBTC (ATOM/BTC) sym3=BTCUSDT (BTC/USDT)
    [REAL EXEC] >>> GetBalance USDT (before)
2025-12-17 01:58:41.801
[ARB] -0.096%  USDT→ATOM→BTC→USDT  maxStart=16.8000 USDT (16.8000 USDT)  safeStart=16.8000 USDT (16.8000 USDT) (x1.00)  bottleneck=ATOMBTC
  ATOMUSDT (ATOM/USDT): bid=2.0150000000 ask=2.0160000000  spread=0.0010000000 (0.04962%)  bidQty=191.1000 askQty=23.1700
  ATOMBTC (ATOM/BTC): bid=0.0000229700 ask=0.0000230100  spread=0.0000000400 (0.17399%)  bidQty=8.3300 askQty=11.4300
  BTCUSDT (BTC/USDT): bid=87788.0900000000 ask=87788.1500000000  spread=0.0600000000 (0.00007%)  bidQty=0.2090 askQty=0.0026

2025-12-17 01:58:41.811
[ARB] -0.096%  USDT→ATOM→BTC→USDT  maxStart=16.8000 USDT (16.8000 USDT)  safeStart=16.8000 USDT (16.8000 USDT) (x1.00)  bottleneck=ATOMBTC
  ATOMUSDT (ATOM/USDT): bid=2.0150000000 ask=2.0160000000  spread=0.0010000000 (0.04962%)  bidQty=191.1000 askQty=23.1700
  ATOMBTC (ATOM/BTC): bid=0.0000229700 ask=0.0000230100  spread=0.0000000400 (0.17399%)  bidQty=8.3300 askQty=11.4300
  BTCUSDT (BTC/USDT): bid=87788.0900000000 ask=87788.1600000000  spread=0.0700000000 (0.00008%)  bidQty=0.2090 askQty=0.0023

2025-12-17 01:58:41.822
[ARB] -0.096%  USDT→ATOM→BTC→USDT  maxStart=16.8000 USDT (16.8000 USDT)  safeStart=16.8000 USDT (16.8000 USDT) (x1.00)  bottleneck=ATOMBTC
  ATOMUSDT (ATOM/USDT): bid=2.0150000000 ask=2.0160000000  spread=0.0010000000 (0.04962%)  bidQty=191.1000 askQty=23.1700
  ATOMBTC (ATOM/BTC): bid=0.0000229700 ask=0.0000230100  spread=0.0000000400 (0.17399%)  bidQty=8.3300 askQty=11.4300
  BTCUSDT (BTC/USDT): bid=87788.0900000000 ask=87793.0900000000  spread=5.0000000000 (0.00570%)  bidQty=6.0328 askQty=0.0324

2025-12-17 01:58:41.830
[ARB] -0.096%  USDT→ATOM→BTC→USDT  maxStart=16.8000 USDT (16.8000 USDT)  safeStart=16.8000 USDT (16.8000 USDT) (x1.00)  bottleneck=ATOMBTC
  ATOMUSDT (ATOM/USDT): bid=2.0150000000 ask=2.0160000000  spread=0.0010000000 (0.04962%)  bidQty=191.1000 askQty=23.1700
  ATOMBTC (ATOM/BTC): bid=0.0000229700 ask=0.0000230100  spread=0.0000000400 (0.17399%)  bidQty=8.3300 askQty=11.4300
  BTCUSDT (BTC/USDT): bid=87788.0900000000 ask=87792.8900000000  spread=4.8000000000 (0.00547%)  bidQty=6.0328 askQty=0.0324

2025-12-17 01:58:41.840
[ARB] -0.090%  USDT→ATOM→BTC→USDT  maxStart=16.8000 USDT (16.8000 USDT)  safeStart=16.8000 USDT (16.8000 USDT) (x1.00)  bottleneck=ATOMBTC
  ATOMUSDT (ATOM/USDT): bid=2.0150000000 ask=2.0160000000  spread=0.0010000000 (0.04962%)  bidQty=191.1000 askQty=23.1700
  ATOMBTC (ATOM/BTC): bid=0.0000229700 ask=0.0000230100  spread=0.0000000400 (0.17399%)  bidQty=8.3300 askQty=11.4300
  BTCUSDT (BTC/USDT): bid=87792.6400000000 ask=87792.9100000000  spread=0.2700000000 (0.00031%)  bidQty=0.0025 askQty=0.0023

2025-12-17 01:58:41.850
[ARB] -0.090%  USDT→ATOM→BTC→USDT  maxStart=16.8000 USDT (16.8000 USDT)  safeStart=16.8000 USDT (16.8000 USDT) (x1.00)  bottleneck=ATOMBTC
  ATOMUSDT (ATOM/USDT): bid=2.0150000000 ask=2.0160000000  spread=0.0010000000 (0.04962%)  bidQty=191.1000 askQty=23.1700
  ATOMBTC (ATOM/BTC): bid=0.0000229700 ask=0.0000230100  spread=0.0000000400 (0.17399%)  bidQty=8.3300 askQty=11.4300
  BTCUSDT (BTC/USDT): bid=87792.8800000000 ask=87799.1000000000  spread=6.2200000000 (0.00708%)  bidQty=0.0044 askQty=2.1562

2025-12-17 01:58:41.860
[ARB] -0.088%  USDT→ATOM→BTC→USDT  maxStart=16.8000 USDT (16.8000 USDT)  safeStart=16.8000 USDT (16.8000 USDT) (x1.00)  bottleneck=ATOMBTC
  ATOMUSDT (ATOM/USDT): bid=2.0150000000 ask=2.0160000000  spread=0.0010000000 (0.04962%)  bidQty=191.1000 askQty=23.1700
  ATOMBTC (ATOM/BTC): bid=0.0000229700 ask=0.0000230100  spread=0.0000000400 (0.17399%)  bidQty=8.3300 askQty=11.4300
  BTCUSDT (BTC/USDT): bid=87794.3000000000 ask=87801.2300000000  spread=6.9300000000 (0.00789%)  bidQty=6.0328 askQty=0.0025

2025-12-17 01:58:41.870
[ARB] -0.088%  USDT→ATOM→BTC→USDT  maxStart=16.8000 USDT (16.8000 USDT)  safeStart=16.8000 USDT (16.8000 USDT) (x1.00)  bottleneck=ATOMBTC
  ATOMUSDT (ATOM/USDT): bid=2.0150000000 ask=2.0160000000  spread=0.0010000000 (0.04962%)  bidQty=191.1000 askQty=23.1700
  ATOMBTC (ATOM/BTC): bid=0.0000229700 ask=0.0000230100  spread=0.0000000400 (0.17399%)  bidQty=8.3300 askQty=11.4300
  BTCUSDT (BTC/USDT): bid=87794.4400000000 ask=87801.8200000000  spread=7.3800000000 (0.00841%)  bidQty=6.0328 askQty=0.0055

2025-12-17 01:58:41.879
[ARB] -0.088%  USDT→ATOM→BTC→USDT  maxStart=16.8000 USDT (16.8000 USDT)  safeStart=16.8000 USDT (16.8000 USDT) (x1.00)  bottleneck=ATOMBTC
  ATOMUSDT (ATOM/USDT): bid=2.0150000000 ask=2.0160000000  spread=0.0010000000 (0.04962%)  bidQty=191.1000 askQty=23.1700
  ATOMBTC (ATOM/BTC): bid=0.0000229700 ask=0.0000230100  spread=0.0000000400 (0.17399%)  bidQty=8.3300 askQty=11.4300
  BTCUSDT (BTC/USDT): bid=87794.4400000000 ask=87797.6700000000  spread=3.2300000000 (0.00368%)  bidQty=6.0328 askQty=0.0529

2025-12-17 01:58:41.890
[ARB] -0.084%  USDT→ATOM→BTC→USDT  maxStart=16.8000 USDT (16.8000 USDT)  safeStart=16.8000 USDT (16.8000 USDT) (x1.00)  bottleneck=ATOMBTC
  ATOMUSDT (ATOM/USDT): bid=2.0150000000 ask=2.0160000000  spread=0.0010000000 (0.04962%)  bidQty=191.1000 askQty=23.1700
  ATOMBTC (ATOM/BTC): bid=0.0000229700 ask=0.0000230100  spread=0.0000000400 (0.17399%)  bidQty=8.3300 askQty=11.4300
  BTCUSDT (BTC/USDT): bid=87797.9600000000 ask=87803.7700000000  spread=5.8100000000 (0.00662%)  bidQty=3.2624 askQty=0.0023

2025-12-17 01:58:41.900
[ARB] -0.084%  USDT→ATOM→BTC→USDT  maxStart=16.8000 USDT (16.8000 USDT)  safeStart=16.8000 USDT (16.8000 USDT) (x1.00)  bottleneck=ATOMBTC
  ATOMUSDT (ATOM/USDT): bid=2.0150000000 ask=2.0160000000  spread=0.0010000000 (0.04962%)  bidQty=191.1000 askQty=23.1700
  ATOMBTC (ATOM/BTC): bid=0.0000229700 ask=0.0000230100  spread=0.0000000400 (0.17399%)  bidQty=8.3300 askQty=11.4300
  BTCUSDT (BTC/USDT): bid=87797.9600000000 ask=87809.8200000000  spread=11.8600000000 (0.01351%)  bidQty=3.2624 askQty=0.0348

2025-12-17 01:58:41.909
[ARB] -0.086%  USDT→ATOM→BTC→USDT  maxStart=16.8000 USDT (16.8000 USDT)  safeStart=16.8000 USDT (16.8000 USDT) (x1.00)  bottleneck=ATOMBTC
  ATOMUSDT (ATOM/USDT): bid=2.0150000000 ask=2.0160000000  spread=0.0010000000 (0.04962%)  bidQty=191.1000 askQty=23.1700
  ATOMBTC (ATOM/BTC): bid=0.0000229700 ask=0.0000230100  spread=0.0000000400 (0.17399%)  bidQty=8.3300 askQty=11.4300
  BTCUSDT (BTC/USDT): bid=87796.2700000000 ask=87801.1900000000  spread=4.9200000000 (0.00560%)  bidQty=0.0439 askQty=0.0529

2025-12-17 01:58:41.920
[ARB] -0.088%  USDT→ATOM→BTC→USDT  maxStart=16.8000 USDT (16.8000 USDT)  safeStart=16.8000 USDT (16.8000 USDT) (x1.00)  bottleneck=ATOMBTC
  ATOMUSDT (ATOM/USDT): bid=2.0150000000 ask=2.0160000000  spread=0.0010000000 (0.04962%)  bidQty=191.1000 askQty=23.1700
  ATOMBTC (ATOM/BTC): bid=0.0000229700 ask=0.0000230100  spread=0.0000000400 (0.17399%)  bidQty=8.3300 askQty=11.4300
  BTCUSDT (BTC/USDT): bid=87794.3100000000 ask=87801.1900000000  spread=6.8800000000 (0.00784%)  bidQty=0.0050 askQty=0.0529

2025-12-17 01:58:41.930
[ARB] -0.084%  USDT→ATOM→BTC→USDT  maxStart=16.8000 USDT (16.8000 USDT)  safeStart=16.8000 USDT (16.8000 USDT) (x1.00)  bottleneck=ATOMBTC
  ATOMUSDT (ATOM/USDT): bid=2.0150000000 ask=2.0160000000  spread=0.0010000000 (0.04962%)  bidQty=191.1000 askQty=23.1700
  ATOMBTC (ATOM/BTC): bid=0.0000229700 ask=0.0000230100  spread=0.0000000400 (0.17399%)  bidQty=8.3300 askQty=11.4300
  BTCUSDT (BTC/USDT): bid=87797.9600000000 ask=87801.1900000000  spread=3.2300000000 (0.00368%)  bidQty=3.2624 askQty=0.0529

2025-12-17 01:58:41.978
[ARB] -0.084%  USDT→ATOM→BTC→USDT  maxStart=16.8000 USDT (16.8000 USDT)  safeStart=16.8000 USDT (16.8000 USDT) (x1.00)  bottleneck=ATOMBTC
  ATOMUSDT (ATOM/USDT): bid=2.0150000000 ask=2.0160000000  spread=0.0010000000 (0.04962%)  bidQty=191.1000 askQty=23.1700
  ATOMBTC (ATOM/BTC): bid=0.0000229700 ask=0.0000230100  spread=0.0000000400 (0.17399%)  bidQty=8.3300 askQty=11.4300
  BTCUSDT (BTC/USDT): bid=87797.9600000000 ask=87802.3900000000  spread=4.4300000000 (0.00505%)  bidQty=3.2624 askQty=0.0024

2025-12-17 01:58:42.070
[ARB] -0.084%  USDT→ATOM→BTC→USDT  maxStart=16.8000 USDT (16.8000 USDT)  safeStart=16.8000 USDT (16.8000 USDT) (x1.00)  bottleneck=ATOMBTC
  ATOMUSDT (ATOM/USDT): bid=2.0150000000 ask=2.0160000000  spread=0.0010000000 (0.04962%)  bidQty=191.1000 askQty=23.1700
  ATOMBTC (ATOM/BTC): bid=0.0000229700 ask=0.0000230100  spread=0.0000000400 (0.17399%)  bidQty=8.3300 askQty=11.4300
  BTCUSDT (BTC/USDT): bid=87797.9600000000 ask=87802.6900000000  spread=4.7300000000 (0.00539%)  bidQty=3.2624 askQty=0.0022

    [REAL EXEC] <<< GetBalance USDT (before) (302ms)
    [REAL EXEC] BAL before: USDT=34.782730681146
    [REAL EXEC] >>> GetBalance ATOM (before leg1)
2025-12-17 01:58:42.171
[ARB] -0.084%  USDT→ATOM→BTC→USDT  maxStart=16.8000 USDT (16.8000 USDT)  safeStart=16.8000 USDT (16.8000 USDT) (x1.00)  bottleneck=ATOMBTC
  ATOMUSDT (ATOM/USDT): bid=2.0150000000 ask=2.0160000000  spread=0.0010000000 (0.04962%)  bidQty=191.1000 askQty=23.1700
  ATOMBTC (ATOM/BTC): bid=0.0000229700 ask=0.0000230100  spread=0.0000000400 (0.17399%)  bidQty=8.3300 askQty=11.4300
  BTCUSDT (BTC/USDT): bid=87797.9600000000 ask=87802.5600000000  spread=4.6000000000 (0.00524%)  bidQty=3.2624 askQty=0.0024

    [REAL EXEC] <<< GetBalance ATOM (before leg1) (237ms)
    [REAL EXEC] leg1 PRE: BUY ATOMUSDT by USDT=2.000000 ask=2.0160000000 bid=2.0150000000 | ATOM before=0.000000000000
    [REAL EXEC] >>> SmartMarketBuyUSDT leg1
2025-12-17 01:58:42.331
[ARB] -0.084%  USDT→ATOM→BTC→USDT  maxStart=16.8000 USDT (16.8000 USDT)  safeStart=16.8000 USDT (16.8000 USDT) (x1.00)  bottleneck=ATOMBTC
  ATOMUSDT (ATOM/USDT): bid=2.0150000000 ask=2.0160000000  spread=0.0010000000 (0.04962%)  bidQty=191.1000 askQty=23.1700
  ATOMBTC (ATOM/BTC): bid=0.0000229700 ask=0.0000230100  spread=0.0000000400 (0.17399%)  bidQty=8.3300 askQty=11.4300
  BTCUSDT (BTC/USDT): bid=87797.9600000000 ask=87802.4100000000  spread=4.4500000000 (0.00507%)  bidQty=3.2624 askQty=0.4550

2025-12-17 01:58:42.360
[ARB] -0.084%  USDT→ATOM→BTC→USDT  maxStart=16.8000 USDT (16.8000 USDT)  safeStart=16.8000 USDT (16.8000 USDT) (x1.00)  bottleneck=ATOMBTC
  ATOMUSDT (ATOM/USDT): bid=2.0150000000 ask=2.0160000000  spread=0.0010000000 (0.04962%)  bidQty=191.1000 askQty=23.1700
  ATOMBTC (ATOM/BTC): bid=0.0000229700 ask=0.0000230100  spread=0.0000000400 (0.17399%)  bidQty=8.3300 askQty=11.4300
  BTCUSDT (BTC/USDT): bid=87797.9600000000 ask=87802.1700000000  spread=4.2100000000 (0.00479%)  bidQty=3.2624 askQty=0.4550

2025-12-17 01:58:42.370
[ARB] -0.084%  USDT→ATOM→BTC→USDT  maxStart=16.8000 USDT (16.8000 USDT)  safeStart=16.8000 USDT (16.8000 USDT) (x1.00)  bottleneck=ATOMBTC
  ATOMUSDT (ATOM/USDT): bid=2.0150000000 ask=2.0160000000  spread=0.0010000000 (0.04962%)  bidQty=191.1000 askQty=23.1700
  ATOMBTC (ATOM/BTC): bid=0.0000229700 ask=0.0000230100  spread=0.0000000400 (0.17399%)  bidQty=8.3300 askQty=11.4300
  BTCUSDT (BTC/USDT): bid=87797.9600000000 ask=87802.0200000000  spread=4.0600000000 (0.00462%)  bidQty=3.2624 askQty=0.4555

2025-12-17 01:58:42.442
[ARB] -0.084%  USDT→ATOM→BTC→USDT  maxStart=16.8000 USDT (16.8000 USDT)  safeStart=16.8000 USDT (16.8000 USDT) (x1.00)  bottleneck=ATOMBTC
  ATOMUSDT (ATOM/USDT): bid=2.0150000000 ask=2.0160000000  spread=0.0010000000 (0.04962%)  bidQty=191.1000 askQty=23.1700
  ATOMBTC (ATOM/BTC): bid=0.0000229700 ask=0.0000230000  spread=0.0000000300 (0.13052%)  bidQty=8.3300 askQty=11.5900
  BTCUSDT (BTC/USDT): bid=87797.9600000000 ask=87802.0200000000  spread=4.0600000000 (0.00462%)  bidQty=3.2624 askQty=0.4555

2025-12-17 01:58:42.472
[ARB] -0.084%  USDT→ATOM→BTC→USDT  maxStart=16.8000 USDT (16.8000 USDT)  safeStart=16.8000 USDT (16.8000 USDT) (x1.00)  bottleneck=ATOMBTC
  ATOMUSDT (ATOM/USDT): bid=2.0150000000 ask=2.0160000000  spread=0.0010000000 (0.04962%)  bidQty=191.1000 askQty=23.1700
  ATOMBTC (ATOM/BTC): bid=0.0000229700 ask=0.0000230000  spread=0.0000000300 (0.13052%)  bidQty=8.3300 askQty=11.5900
  BTCUSDT (BTC/USDT): bid=87797.9600000000 ask=87801.1900000000  spread=3.2300000000 (0.00368%)  bidQty=3.2624 askQty=0.4556

2025-12-17 01:58:42.520
[ARB] -0.084%  USDT→ATOM→BTC→USDT  maxStart=16.8000 USDT (16.8000 USDT)  safeStart=16.8000 USDT (16.8000 USDT) (x1.00)  bottleneck=ATOMBTC
  ATOMUSDT (ATOM/USDT): bid=2.0150000000 ask=2.0160000000  spread=0.0010000000 (0.04962%)  bidQty=191.1000 askQty=23.1700
  ATOMBTC (ATOM/BTC): bid=0.0000229700 ask=0.0000230000  spread=0.0000000300 (0.13052%)  bidQty=8.3300 askQty=11.5900
  BTCUSDT (BTC/USDT): bid=87797.9600000000 ask=87800.8200000000  spread=2.8600000000 (0.00326%)  bidQty=3.2624 askQty=0.4556

2025-12-17 01:58:42.541
[ARB] -0.084%  USDT→ATOM→BTC→USDT  maxStart=16.8000 USDT (16.8000 USDT)  safeStart=16.8000 USDT (16.8000 USDT) (x1.00)  bottleneck=ATOMBTC
  ATOMUSDT (ATOM/USDT): bid=2.0150000000 ask=2.0160000000  spread=0.0010000000 (0.04962%)  bidQty=191.1000 askQty=23.1700
  ATOMBTC (ATOM/BTC): bid=0.0000229700 ask=0.0000230000  spread=0.0000000300 (0.13052%)  bidQty=8.3300 askQty=11.5900
  BTCUSDT (BTC/USDT): bid=87797.9600000000 ask=87800.7100000000  spread=2.7500000000 (0.00313%)  bidQty=3.2624 askQty=0.4556

2025-12-17 01:58:42.569
[ARB] -0.084%  USDT→ATOM→BTC→USDT  maxStart=16.8000 USDT (16.8000 USDT)  safeStart=16.8000 USDT (16.8000 USDT) (x1.00)  bottleneck=ATOMBTC
  ATOMUSDT (ATOM/USDT): bid=2.0150000000 ask=2.0160000000  spread=0.0010000000 (0.04962%)  bidQty=191.1000 askQty=23.1700
  ATOMBTC (ATOM/BTC): bid=0.0000229700 ask=0.0000230000  spread=0.0000000300 (0.13052%)  bidQty=8.3300 askQty=11.5900
  BTCUSDT (BTC/USDT): bid=87797.9600000000 ask=87800.5400000000  spread=2.5800000000 (0.00294%)  bidQty=3.2624 askQty=0.4556

    [REAL EXEC] <<< SmartMarketBuyUSDT leg1 (295ms)
    [REAL EXEC] leg1 PLACE ERR: qty<=0 after normalize (raw=0.992063492063)
2025-12-17 01:58:42.650
[ARB] -0.084%  USDT→ATOM→BTC→USDT  maxStart=16.8000 USDT (16.8000 USDT)  safeStart=16.8000 USDT (16.8000 USDT) (x1.00)  bottleneck=ATOMBTC
  ATOMUSDT (ATOM/USDT): bid=2.0150000000 ask=2.0160000000  spread=0.0010000000 (0.04962%)  bidQty=191.1000 askQty=23.1700
  ATOMBTC (ATOM/BTC): bid=0.0000229700 ask=0.0000230000  spread=0.0000000300 (0.13052%)  bidQty=8.3300 askQty=11.5900
  BTCUSDT (BTC/USDT): bid=87797.9600000000 ask=87800.5500000000  spread=2.5900000000 (0.00295%)  bidQty=3.2624 askQty=0.4580

  [REAL EXEC] start=2.000000 USDT triangle=USDT→ATOM→BTC→USDT
    [REAL EXEC] legs: sym1=ATOMUSDT sym2=ATOMBTC sym3=BTCUSDT
    [REAL EXEC] parsed: sym1=ATOMUSDT (ATOM/USDT) sym2=ATOMBTC (ATOM/BTC) sym3=BTCUSDT (BTC/USDT)
    [REAL EXEC] >>> GetBalance USDT (before)
2025-12-17 01:58:42.659
[ARB] -0.084%  USDT→ATOM→BTC→USDT  maxStart=16.8000 USDT (16.8000 USDT)  safeStart=16.8000 USDT (16.8000 USDT) (x1.00)  bottleneck=ATOMBTC
  ATOMUSDT (ATOM/USDT): bid=2.0150000000 ask=2.0160000000  spread=0.0010000000 (0.04962%)  bidQty=191.1000 askQty=23.1700
  ATOMBTC (ATOM/BTC): bid=0.0000229700 ask=0.0000230000  spread=0.0000000300 (0.13052%)  bidQty=8.3300 askQty=11.5900
  BTCUSDT (BTC/USDT): bid=87797.9600000000 ask=87800.0900000000  spread=2.1300000000 (0.00243%)  bidQty=3.2624 askQty=0.4556

2025-12-17 01:58:42.670
[ARB] -0.084%  USDT→ATOM→BTC→USDT  maxStart=16.8000 USDT (16.8000 USDT)  safeStart=16.8000 USDT (16.8000 USDT) (x1.00)  bottleneck=ATOMBTC
  ATOMUSDT (ATOM/USDT): bid=2.0150000000 ask=2.0160000000  spread=0.0010000000 (0.04962%)  bidQty=191.1000 askQty=23.1700
  ATOMBTC (ATOM/BTC): bid=0.0000229700 ask=0.0000230000  spread=0.0000000300 (0.13052%)  bidQty=8.3300 askQty=11.5900
  BTCUSDT (BTC/USDT): bid=87797.9600000000 ask=87800.0400000000  spread=2.0800000000 (0.00237%)  bidQty=3.2624 askQty=0.4556

2025-12-17 01:58:42.770
[ARB] -0.084%  USDT→ATOM→BTC→USDT  maxStart=16.8000 USDT (16.8000 USDT)  safeStart=16.8000 USDT (16.8000 USDT) (x1.00)  bottleneck=ATOMBTC
  ATOMUSDT (ATOM/USDT): bid=2.0150000000 ask=2.0160000000  spread=0.0010000000 (0.04962%)  bidQty=191.1000 askQty=23.1700
  ATOMBTC (ATOM/BTC): bid=0.0000229700 ask=0.0000230000  spread=0.0000000300 (0.13052%)  bidQty=8.3300 askQty=11.5900
  BTCUSDT (BTC/USDT): bid=87797.9600000000 ask=87799.6300000000  spread=1.6700000000 (0.00190%)  bidQty=3.2624 askQty=0.4556

2025-12-17 01:58:42.821
[ARB] -0.084%  USDT→ATOM→BTC→USDT  maxStart=16.8000 USDT (16.8000 USDT)  safeStart=16.8000 USDT (16.8000 USDT) (x1.00)  bottleneck=ATOMBTC
  ATOMUSDT (ATOM/USDT): bid=2.0150000000 ask=2.0160000000  spread=0.0010000000 (0.04962%)  bidQty=191.1000 askQty=23.1700
  ATOMBTC (ATOM/BTC): bid=0.0000229700 ask=0.0000230000  spread=0.0000000300 (0.13052%)  bidQty=8.3300 askQty=11.5900
  BTCUSDT (BTC/USDT): bid=87797.9600000000 ask=87799.4400000000  spread=1.4800000000 (0.00169%)  bidQty=3.2624 askQty=0.4556

2025-12-17 01:58:42.841
[ARB] -0.084%  USDT→ATOM→BTC→USDT  maxStart=16.8000 USDT (16.8000 USDT)  safeStart=16.8000 USDT (16.8000 USDT) (x1.00)  bottleneck=ATOMBTC
  ATOMUSDT (ATOM/USDT): bid=2.0150000000 ask=2.0160000000  spread=0.0010000000 (0.04962%)  bidQty=191.1000 askQty=23.1700
  ATOMBTC (ATOM/BTC): bid=0.0000229700 ask=0.0000230000  spread=0.0000000300 (0.13052%)  bidQty=8.3300 askQty=11.5900
  BTCUSDT (BTC/USDT): bid=87797.9600000000 ask=87799.3700000000  spread=1.4100000000 (0.00161%)  bidQty=3.2624 askQty=0.4556

2025-12-17 01:58:42.851
[ARB] -0.084%  USDT→ATOM→BTC→USDT  maxStart=16.8000 USDT (16.8000 USDT)  safeStart=16.8000 USDT (16.8000 USDT) (x1.00)  bottleneck=ATOMBTC
  ATOMUSDT (ATOM/USDT): bid=2.0150000000 ask=2.0160000000  spread=0.0010000000 (0.04962%)  bidQty=191.1000 askQty=23.1700
  ATOMBTC (ATOM/BTC): bid=0.0000229700 ask=0.0000230000  spread=0.0000000300 (0.13052%)  bidQty=8.3300 askQty=11.5900
  BTCUSDT (BTC/USDT): bid=87797.9600000000 ask=87799.3400000000  spread=1.3800000000 (0.00157%)  bidQty=3.2624 askQty=0.4556

2025-12-17 01:58:42.861
[ARB] -0.084%  USDT→ATOM→BTC→USDT  maxStart=16.8000 USDT (16.8000 USDT)  safeStart=16.8000 USDT (16.8000 USDT) (x1.00)  bottleneck=ATOMBTC
  ATOMUSDT (ATOM/USDT): bid=2.0150000000 ask=2.0160000000  spread=0.0010000000 (0.04962%)  bidQty=191.1000 askQty=23.1700
  ATOMBTC (ATOM/BTC): bid=0.0000229700 ask=0.0000230000  spread=0.0000000300 (0.13052%)  bidQty=8.3300 askQty=11.5900
  BTCUSDT (BTC/USDT): bid=87797.9600000000 ask=87799.3300000000  spread=1.3700000000 (0.00156%)  bidQty=3.2624 askQty=0.4556

2025-12-17 01:58:42.880
[ARB] -0.084%  USDT→ATOM→BTC→USDT  maxStart=16.8000 USDT (16.8000 USDT)  safeStart=16.8000 USDT (16.8000 USDT) (x1.00)  bottleneck=ATOMBTC
  ATOMUSDT (ATOM/USDT): bid=2.0150000000 ask=2.0160000000  spread=0.0010000000 (0.04962%)  bidQty=191.1000 askQty=23.1700
  ATOMBTC (ATOM/BTC): bid=0.0000229700 ask=0.0000230000  spread=0.0000000300 (0.13052%)  bidQty=8.3300 askQty=11.5900
  BTCUSDT (BTC/USDT): bid=87797.9600000000 ask=87799.2900000000  spread=1.3300000000 (0.00151%)  bidQty=3.2624 askQty=0.0317

2025-12-17 01:58:42.911
[ARB] -0.084%  USDT→ATOM→BTC→USDT  maxStart=16.8000 USDT (16.8000 USDT)  safeStart=16.8000 USDT (16.8000 USDT) (x1.00)  bottleneck=ATOMBTC
  ATOMUSDT (ATOM/USDT): bid=2.0150000000 ask=2.0160000000  spread=0.0010000000 (0.04962%)  bidQty=191.1000 askQty=23.1700
  ATOMBTC (ATOM/BTC): bid=0.0000229700 ask=0.0000230000  spread=0.0000000300 (0.13052%)  bidQty=8.3300 askQty=11.5900
  BTCUSDT (BTC/USDT): bid=87797.9600000000 ask=87799.1800000000  spread=1.2200000000 (0.00139%)  bidQty=3.2624 askQty=0.0317

    [REAL EXEC] <<< GetBalance USDT (before) (265ms)
    [REAL EXEC] BAL before: USDT=34.782730681146
    [REAL EXEC] >>> GetBalance ATOM (before leg1)
2025-12-17 01:58:42.925
[ARB] -0.084%  USDT→ATOM→BTC→USDT  maxStart=16.8000 USDT (16.8000 USDT)  safeStart=16.8000 USDT (16.8000 USDT) (x1.00)  bottleneck=ATOMBTC
  ATOMUSDT (ATOM/USDT): bid=2.0150000000 ask=2.0160000000  spread=0.0010000000 (0.04962%)  bidQty=191.1000 askQty=23.1700
  ATOMBTC (ATOM/BTC): bid=0.0000229700 ask=0.0000230000  spread=0.0000000300 (0.13052%)  bidQty=8.3300 askQty=11.5900
  BTCUSDT (BTC/USDT): bid=87797.9600000000 ask=87799.1600000000  spread=1.2000000000 (0.00137%)  bidQty=3.2624 askQty=0.0317



