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

// SpotTrader должен уметь:
// - SmartMarketBuyUSDT: MARKET BUY на фиксированную сумму USDT
// - SmartMarketSellQty: MARKET SELL по quantity (qty)
// - GetBalance: получить free (или доступный) баланс
type SpotTrader interface {
	SmartMarketBuyUSDT(ctx context.Context, symbol string, usdt float64, ask float64) (string, error)
	SmartMarketSellQty(ctx context.Context, symbol string, qty float64) (string, error)
	GetBalance(ctx context.Context, asset string) (float64, error)
}

type RealExecutor struct {
	trader SpotTrader
	out    io.Writer

	// фиксированный старт (например 2 или 10)
	StartUSDT float64

	// safety чтобы не словить Oversold на SELL
	SellSafety float64

	// анти-спам
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

func (e *RealExecutor) logf(format string, args ...any) {
	// io.Writer не имеет Printf — только fmt.Fprintf
	fmt.Fprintf(e.out, format+"\n", args...)
}

func (e *RealExecutor) Execute(ctx context.Context, t domain.Triangle, quotes map[string]domain.Quote, startUSDT float64) error {
	key := t.String()
	now := time.Now()

	// cooldown
	if last, ok := e.lastExec[key]; ok && e.Cooldown > 0 {
		if now.Sub(last) < e.Cooldown {
			e.logf("  [REAL EXEC] SKIP cooldown triangle=%s left=%s",
				key, (e.Cooldown-now.Sub(last)).Truncate(time.Millisecond))
			return nil
		}
	}

	if startUSDT <= 0 {
		startUSDT = e.StartUSDT
	}
	if startUSDT <= 0 {
		return fmt.Errorf("startUSDT<=0 (startUSDT=%.6f, StartUSDT=%.6f)", startUSDT, e.StartUSDT)
	}

	// разберём активы из символов (для понятных логов)
	a1Base, a1Quote := parseBaseQuote(t.Leg1.Symbol)
	a2Base, a2Quote := parseBaseQuote(t.Leg2.Symbol)
	a3Base, a3Quote := parseBaseQuote(t.Leg3.Symbol)

	e.logf("  [REAL EXEC] start=%.6f USDT triangle=%s", startUSDT, key)
	e.logf("    [REAL EXEC] symbols: leg1=%s (%s/%s) leg2=%s (%s/%s) leg3=%s (%s/%s)",
		t.Leg1.Symbol, a1Base, a1Quote,
		t.Leg2.Symbol, a2Base, a2Quote,
		t.Leg3.Symbol, a3Base, a3Quote,
	)

	// ===== балансы до =====
	usdt0, err := e.trader.GetBalance(ctx, "USDT")
	if err != nil {
		e.logf("    [REAL EXEC] BAL ERR: get USDT before: %v", err)
		return err
	}
	e.logf("    [REAL EXEC] BAL before: USDT=%.12f", usdt0)

	if usdt0+1e-9 < startUSDT {
		return fmt.Errorf("insufficient USDT: have=%.12f need=%.12f", usdt0, startUSDT)
	}

	// ===== LEG 1: BUY A за USDT =====
	q1, ok := quotes[t.Leg1.Symbol]
	if !ok {
		return fmt.Errorf("no quote for leg1 symbol=%s", t.Leg1.Symbol)
	}
	e.logf("    [REAL EXEC] leg1 PRE: BUY %s by USDT=%.6f ask=%.10f bid=%.10f",
		t.Leg1.Symbol, startUSDT, q1.Ask, q1.Bid)

	// До leg1: баланс base A
	aBefore1, _ := e.trader.GetBalance(ctx, a1Base)
	e.logf("    [REAL EXEC] leg1 BAL before: %s=%.12f", a1Base, aBefore1)

	ord1, err := e.trader.SmartMarketBuyUSDT(ctx, t.Leg1.Symbol, startUSDT, q1.Ask)
	if err != nil {
		e.logf("    [REAL EXEC] leg1 PLACE ERR: %v", err)
		return err
	}
	e.logf("    [REAL EXEC] leg1 PLACE OK: orderId=%s", ord1)

	// Ждём изменения баланса A (главный индикатор что биржа реально исполнила)
	aAfter1, err := e.waitBalanceChange(ctx, a1Base, aBefore1, 3*time.Second, 150*time.Millisecond)
	if err != nil {
		e.logf("    [REAL EXEC] leg1 WAIT BAL ERR (%s): %v", a1Base, err)
		return err
	}
	dA := aAfter1 - aBefore1
	e.logf("    [REAL EXEC] leg1 BAL after: %s=%.12f delta=%.12f", a1Base, aAfter1, dA)

	if dA <= 0 {
		return fmt.Errorf("leg1: %s did not increase (before=%.12f after=%.12f)", a1Base, aBefore1, aAfter1)
	}

	// ===== LEG 2: SELL A -> B (quote 2-й пары) =====
	q2, ok := quotes[t.Leg2.Symbol]
	if !ok {
		return fmt.Errorf("no quote for leg2 symbol=%s", t.Leg2.Symbol)
	}

	// Safety на oversold: продаём чуть меньше фактического
	sellA := aAfter1 * e.SellSafety
	sellA = clampNonNegative(sellA)

	e.logf("    [REAL EXEC] leg2 PRE: SELL %s qty=%s=%.12f (safety x%.6f) bid=%.10f ask=%.10f",
		t.Leg2.Symbol, a2Base, sellA, e.SellSafety, q2.Bid, q2.Ask)

	// Баланс B до leg2 (это quote второй пары)
	bBefore2, err := e.trader.GetBalance(ctx, a2Quote)
	if err != nil {
		e.logf("    [REAL EXEC] BAL ERR: get %s before leg2: %v", a2Quote, err)
		return err
	}
	e.logf("    [REAL EXEC] leg2 BAL before: %s=%.12f", a2Quote, bBefore2)

	if sellA <= 0 {
		return fmt.Errorf("leg2: sell qty <=0 after safety (%s=%.12f)", a1Base, sellA)
	}

	ord2, err := e.trader.SmartMarketSellQty(ctx, t.Leg2.Symbol, sellA)
	if err != nil {
		e.logf("    [REAL EXEC] leg2 PLACE ERR: %v", err)
		return err
	}
	e.logf("    [REAL EXEC] leg2 PLACE OK: orderId=%s", ord2)

	bAfter2, err := e.waitBalanceChange(ctx, a2Quote, bBefore2, 3*time.Second, 150*time.Millisecond)
	if err != nil {
		e.logf("    [REAL EXEC] leg2 WAIT BAL ERR (%s): %v", a2Quote, err)
		return err
	}
	dB := bAfter2 - bBefore2
	e.logf("    [REAL EXEC] leg2 BAL after: %s=%.12f delta=%.12f", a2Quote, bAfter2, dB)

	if dB <= 0 {
		return fmt.Errorf("leg2: %s did not increase (before=%.12f after=%.12f)", a2Quote, bBefore2, bAfter2)
	}

	// ===== LEG 3: SELL B -> USDT (третья пара) =====
	q3, ok := quotes[t.Leg3.Symbol]
	if !ok {
		return fmt.Errorf("no quote for leg3 symbol=%s", t.Leg3.Symbol)
	}

	// Что продавать на leg3: base третьей пары = a3Base (обычно это и есть B)
	// Баланс base третьей пары до leg3:
	bBaseBefore3, err := e.trader.GetBalance(ctx, a3Base)
	if err != nil {
		e.logf("    [REAL EXEC] BAL ERR: get %s before leg3: %v", a3Base, err)
		return err
	}
	e.logf("    [REAL EXEC] leg3 BAL before: %s=%.12f", a3Base, bBaseBefore3)

	sellB := bBaseBefore3 * e.SellSafety
	sellB = clampNonNegative(sellB)

	e.logf("    [REAL EXEC] leg3 PRE: SELL %s qty=%s=%.12f (safety x%.6f) bid=%.10f ask=%.10f",
		t.Leg3.Symbol, a3Base, sellB, e.SellSafety, q3.Bid, q3.Ask)

	usdtBefore3, err := e.trader.GetBalance(ctx, "USDT")
	if err != nil {
		e.logf("    [REAL EXEC] BAL ERR: get USDT before leg3: %v", err)
		return err
	}
	e.logf("    [REAL EXEC] leg3 BAL before: USDT=%.12f", usdtBefore3)

	if sellB <= 0 {
		return fmt.Errorf("leg3: sell qty <=0 after safety (%s=%.12f)", a3Base, sellB)
	}

	ord3, err := e.trader.SmartMarketSellQty(ctx, t.Leg3.Symbol, sellB)
	if err != nil {
		e.logf("    [REAL EXEC] leg3 PLACE ERR: %v", err)
		return err
	}
	e.logf("    [REAL EXEC] leg3 PLACE OK: orderId=%s", ord3)

	usdtAfter, err := e.waitBalanceChange(ctx, "USDT", usdtBefore3, 3*time.Second, 150*time.Millisecond)
	if err != nil {
		e.logf("    [REAL EXEC] leg3 WAIT BAL ERR (USDT): %v", err)
		return err
	}

	dUSDT3 := usdtAfter - usdtBefore3
	dUSDTTotal := usdtAfter - usdt0
	e.logf("    [REAL EXEC] leg3 BAL after: USDT=%.12f delta=%.12f", usdtAfter, dUSDT3)
	e.logf("    [REAL EXEC] DONE: USDT start=%.12f end=%.12f pnl(total)=%.12f (%.4f%%)",
		usdt0, usdtAfter, dUSDTTotal, pct(dUSDTTotal, startUSDT))

	e.lastExec[key] = now
	return nil
}

// waitBalanceChange ждёт, пока баланс станет != baseline (с небольшим порогом).
// Это наш “аналог waitFilled”, раз нет GetOrder().
func (e *RealExecutor) waitBalanceChange(ctx context.Context, asset string, baseline float64, timeout, interval time.Duration) (float64, error) {
	const tol = 1e-12

	deadline := time.NewTimer(timeout)
	tick := time.NewTicker(interval)
	defer deadline.Stop()
	defer tick.Stop()

	// быстрый первый замер
	cur, err := e.trader.GetBalance(ctx, asset)
	if err == nil {
		if math.Abs(cur-baseline) > tol {
			return cur, nil
		}
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

func parseBaseQuote(symbol string) (base, quote string) {
	quotes := []string{"USDT", "USDC", "BTC", "ETH", "EUR", "TRY", "BRL", "RUB"}
	for _, q := range quotes {
		if strings.HasSuffix(symbol, q) && len(symbol) > len(q) {
			return symbol[:len(symbol)-len(q)], q
		}
	}
	return symbol, ""
}

func clampNonNegative(x float64) float64 {
	if x < 0 {
		return 0
	}
	return x
}

func pct(delta, denom float64) float64 {
	if denom == 0 {
		return 0
	}
	return (delta / denom) * 100
}


[{
	"resource": "/home/gaz358/myprog/crypt_proto/arb/executor_real.go",
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
	"message": "t.Leg1 undefined (type domain.Triangle has no field or method Leg1)",
	"source": "compiler",
	"startLineNumber": 78,
	"startColumn": 38,
	"endLineNumber": 78,
	"endColumn": 42,
	"origin": "extHost1"
}]


