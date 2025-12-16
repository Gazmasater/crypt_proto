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
	"strings"
	"time"

	"crypt_proto/domain"
)

// OrderInfo — то, что мы хотим видеть в логах от биржи.
type OrderInfo struct {
	OrderID             string
	Symbol              string
	Side                string // BUY/SELL
	Type                string // MARKET/LIMIT/...
	Status              string // NEW, PARTIALLY_FILLED, FILLED, CANCELED...
	ExecutedQty         float64
	CummulativeQuoteQty float64
	AvgPrice            float64
	Fee                 float64
	FeeAsset            string
	UpdateTime          time.Time
	Raw                 string // сырой ответ (json/строка) для дебага
}

// SpotTrader должен УМЕТЬ внутри себя:
// - SmartMarketBuyUSDT: корректно сформировать MARKET BUY за фиксированный USDT
// - SmartMarketSellQty: корректно SELL по quantity
// - GetBalance: читать free баланс (или доступный)
// - GetOrder: получать статус/исполнение ордера для подробных логов
type SpotTrader interface {
	SmartMarketBuyUSDT(ctx context.Context, symbol string, usdt float64, ask float64) (string, error)
	SmartMarketSellQty(ctx context.Context, symbol string, qty float64) (string, error)
	GetBalance(ctx context.Context, asset string) (float64, error)

	GetOrder(ctx context.Context, symbol, orderID string) (OrderInfo, error)
}

type RealExecutor struct {
	trader SpotTrader
	out    io.Writer

	// фиксированный старт в USDT (например 2 или 10)
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

func (e *RealExecutor) logf(format string, args ...any) {
	fmt.Fprintf(e.out, format+"\n", args...)
}

func (e *RealExecutor) logOrder(prefix string, o OrderInfo) {
	e.logf(
		"    %s order: id=%s sym=%s side=%s type=%s status=%s execQty=%.12f quoteQty=%.12f avg=%.10f fee=%.12f feeAsset=%s raw=%s",
		prefix,
		o.OrderID, o.Symbol, o.Side, o.Type, o.Status,
		o.ExecutedQty, o.CummulativeQuoteQty, o.AvgPrice,
		o.Fee, o.FeeAsset, o.Raw,
	)
}

// Execute исполняет ТОЛЬКО безопасный класс треугольников:
// USDT -> A -> B -> USDT
func (e *RealExecutor) Execute(ctx context.Context, t domain.Triangle, quotes map[string]domain.Quote, startUSDT float64) error {
	key := t.String()
	now := time.Now()

	// cooldown
	if last, ok := e.lastExec[key]; ok && e.Cooldown > 0 {
		if now.Sub(last) < e.Cooldown {
			e.logf("  [REAL EXEC] SKIP cooldown triangle=%s left=%s", key, (e.Cooldown-now.Sub(last)).Truncate(time.Millisecond))
			return nil
		}
	}

	// старт: аргумент приоритетнее, иначе StartUSDT
	if startUSDT <= 0 {
		startUSDT = e.StartUSDT
	}
	if startUSDT <= 0 {
		return fmt.Errorf("startUSDT<=0 (startUSDT=%.6f, StartUSDT=%.6f)", startUSDT, e.StartUSDT)
	}

	e.logf("  [REAL EXEC] start=%.6f USDT triangle=%s", startUSDT, key)

	// Разбор активов (в логах полезно)
	a1Base, a1Quote := parseBaseQuote(t.Leg1.Symbol)
	a2Base, a2Quote := parseBaseQuote(t.Leg2.Symbol)
	a3Base, a3Quote := parseBaseQuote(t.Leg3.Symbol)

	e.logf("    [REAL EXEC] symbols: leg1=%s (%s/%s) leg2=%s (%s/%s) leg3=%s (%s/%s)",
		t.Leg1.Symbol, a1Base, a1Quote,
		t.Leg2.Symbol, a2Base, a2Quote,
		t.Leg3.Symbol, a3Base, a3Quote,
	)

	// Баланс USDT до
	usdt0, err := e.trader.GetBalance(ctx, "USDT")
	if err != nil {
		e.logf("    [REAL EXEC] BAL ERR: get USDT before: %v", err)
		return err
	}
	e.logf("    [REAL EXEC] BAL before: USDT=%.12f", usdt0)

	if usdt0+1e-9 < startUSDT {
		return fmt.Errorf("insufficient USDT balance: have=%.12f need=%.12f", usdt0, startUSDT)
	}

	// ===== LEG 1: BUY A за USDT =====
	q1, ok := quotes[t.Leg1.Symbol]
	if !ok {
		return fmt.Errorf("no quote for leg1 symbol=%s", t.Leg1.Symbol)
	}
	e.logf("    [REAL EXEC] leg1 PRE: BUY %s by USDT=%.6f ask=%.10f bid=%.10f", t.Leg1.Symbol, startUSDT, q1.Ask, q1.Bid)

	order1, err := e.trader.SmartMarketBuyUSDT(ctx, t.Leg1.Symbol, startUSDT, q1.Ask)
	if err != nil {
		e.logf("    [REAL EXEC] leg1 PLACE ERR: %v", err)
		return err
	}
	e.logf("    [REAL EXEC] leg1 PLACE OK: orderId=%s", order1)

	if oi, err := e.trader.GetOrder(ctx, t.Leg1.Symbol, order1); err == nil {
		e.logOrder("[REAL EXEC] leg1 GET", oi)
	} else {
		e.logf("    [REAL EXEC] leg1 GET ERR: %v", err)
	}

	o1, err := e.waitFilledByPolling(ctx, t.Leg1.Symbol, order1, 3*time.Second, 200*time.Millisecond)
	if err != nil {
		e.logf("    [REAL EXEC] leg1 FILL ERR: %v", err)
		e.logOrder("[REAL EXEC] leg1 LAST", o1)
		return err
	}
	e.logOrder("[REAL EXEC] leg1 FILLED", o1)

	// Ждём баланс A (на случай задержек учёта)
	aBal1, err := e.waitBalanceIncrease(ctx, a1Base, 0, 3*time.Second, 150*time.Millisecond)
	if err != nil {
		e.logf("    [REAL EXEC] leg1 WAIT BAL ERR (%s): %v", a1Base, err)
		return err
	}
	e.logf("    [REAL EXEC] leg1 BAL: %s=%.12f", a1Base, aBal1)

	// ===== LEG 2: SELL A -> B =====
	q2, ok := quotes[t.Leg2.Symbol]
	if !ok {
		return fmt.Errorf("no quote for leg2 symbol=%s", t.Leg2.Symbol)
	}

	sellA := aBal1 * e.SellSafety
	if sellA <= 0 {
		return fmt.Errorf("leg2: computed sell qty <=0: %s=%.12f safety=%.6f", a1Base, aBal1, e.SellSafety)
	}

	e.logf("    [REAL EXEC] leg2 PRE: SELL %s qty=%s=%.12f (safety x%.6f) bid=%.10f ask=%.10f",
		t.Leg2.Symbol, a2Base, sellA, e.SellSafety, q2.Bid, q2.Ask)

	b0, err := e.trader.GetBalance(ctx, a2Quote) // quote 2-й пары = B
	if err != nil {
		e.logf("    [REAL EXEC] BAL ERR: get %s before leg2: %v", a2Quote, err)
		return err
	}
	e.logf("    [REAL EXEC] BAL before leg2: %s=%.12f", a2Quote, b0)

	order2, err := e.trader.SmartMarketSellQty(ctx, t.Leg2.Symbol, sellA)
	if err != nil {
		e.logf("    [REAL EXEC] leg2 PLACE ERR: %v", err)
		return err
	}
	e.logf("    [REAL EXEC] leg2 PLACE OK: orderId=%s", order2)

	if oi, err := e.trader.GetOrder(ctx, t.Leg2.Symbol, order2); err == nil {
		e.logOrder("[REAL EXEC] leg2 GET", oi)
	} else {
		e.logf("    [REAL EXEC] leg2 GET ERR: %v", err)
	}

	o2, err := e.waitFilledByPolling(ctx, t.Leg2.Symbol, order2, 3*time.Second, 200*time.Millisecond)
	if err != nil {
		e.logf("    [REAL EXEC] leg2 FILL ERR: %v", err)
		e.logOrder("[REAL EXEC] leg2 LAST", o2)
		return err
	}
	e.logOrder("[REAL EXEC] leg2 FILLED", o2)

	b1, err := e.waitBalanceIncrease(ctx, a2Quote, b0, 3*time.Second, 150*time.Millisecond)
	if err != nil {
		e.logf("    [REAL EXEC] leg2 WAIT BAL ERR (%s): %v", a2Quote, err)
		return err
	}
	e.logf("    [REAL EXEC] leg2 BAL: %s=%.12f (delta=%.12f)", a2Quote, b1, b1-b0)

	// ===== LEG 3: SELL B -> USDT =====
	q3, ok := quotes[t.Leg3.Symbol]
	if !ok {
		return fmt.Errorf("no quote for leg3 symbol=%s", t.Leg3.Symbol)
	}

	bToSell := b1 * e.SellSafety
	if bToSell <= 0 {
		return fmt.Errorf("leg3: computed sell qty <=0: %s=%.12f safety=%.6f", a3Base, b1, e.SellSafety)
	}

	e.logf("    [REAL EXEC] leg3 PRE: SELL %s qty=%s=%.12f (safety x%.6f) bid=%.10f ask=%.10f",
		t.Leg3.Symbol, a3Base, bToSell, e.SellSafety, q3.Bid, q3.Ask)

	usdtBefore3, err := e.trader.GetBalance(ctx, "USDT")
	if err != nil {
		e.logf("    [REAL EXEC] BAL ERR: get USDT before leg3: %v", err)
		return err
	}
	e.logf("    [REAL EXEC] BAL before leg3: USDT=%.12f", usdtBefore3)

	order3, err := e.trader.SmartMarketSellQty(ctx, t.Leg3.Symbol, bToSell)
	if err != nil {
		e.logf("    [REAL EXEC] leg3 PLACE ERR: %v", err)
		return err
	}
	e.logf("    [REAL EXEC] leg3 PLACE OK: orderId=%s", order3)

	if oi, err := e.trader.GetOrder(ctx, t.Leg3.Symbol, order3); err == nil {
		e.logOrder("[REAL EXEC] leg3 GET", oi)
	} else {
		e.logf("    [REAL EXEC] leg3 GET ERR: %v", err)
	}

	o3, err := e.waitFilledByPolling(ctx, t.Leg3.Symbol, order3, 3*time.Second, 200*time.Millisecond)
	if err != nil {
		e.logf("    [REAL EXEC] leg3 FILL ERR: %v", err)
		e.logOrder("[REAL EXEC] leg3 LAST", o3)
		return err
	}
	e.logOrder("[REAL EXEC] leg3 FILLED", o3)

	usdtAfter, err := e.waitBalanceIncrease(ctx, "USDT", usdtBefore3, 3*time.Second, 150*time.Millisecond)
	if err != nil {
		e.logf("    [REAL EXEC] leg3 WAIT BAL ERR (USDT): %v", err)
		return err
	}

	pnlLeg3 := usdtAfter - usdtBefore3
	pnlTotal := usdtAfter - usdt0

	e.logf("    [REAL EXEC] DONE: USDT before3=%.12f after=%.12f pnl(leg3)=%.12f | USDT start=%.12f end=%.12f pnl(total)=%.12f",
		usdtBefore3, usdtAfter, pnlLeg3,
		usdt0, usdtAfter, pnlTotal,
	)

	e.lastExec[key] = now
	return nil
}

// waitFilledByPolling опрашивает GetOrder до FILLED/ошибки/таймаута.
func (e *RealExecutor) waitFilledByPolling(ctx context.Context, symbol, orderID string, timeout, interval time.Duration) (OrderInfo, error) {
	deadline := time.NewTimer(timeout)
	tick := time.NewTicker(interval)
	defer deadline.Stop()
	defer tick.Stop()

	for {
		o, err := e.trader.GetOrder(ctx, symbol, orderID)
		if err == nil {
			switch o.Status {
			case "FILLED":
				return o, nil
			case "CANCELED", "REJECTED", "EXPIRED":
				return o, fmt.Errorf("order finished not filled: id=%s sym=%s status=%s execQty=%.12f quoteQty=%.12f",
					orderID, symbol, o.Status, o.ExecutedQty, o.CummulativeQuoteQty)
			}
		}

		select {
		case <-ctx.Done():
			return OrderInfo{}, ctx.Err()
		case <-deadline.C:
			last, err := e.trader.GetOrder(ctx, symbol, orderID)
			if err != nil {
				return OrderInfo{}, fmt.Errorf("timeout waiting fill; last GetOrder err: %v", err)
			}
			return last, fmt.Errorf("timeout waiting fill: id=%s sym=%s status=%s execQty=%.12f quoteQty=%.12f",
				orderID, symbol, last.Status, last.ExecutedQty, last.CummulativeQuoteQty)
		case <-tick.C:
			// next poll
		}
	}
}

func (e *RealExecutor) waitBalanceIncrease(ctx context.Context, asset string, baseline float64, timeout, interval time.Duration) (float64, error) {
	deadline := time.NewTimer(timeout)
	tick := time.NewTicker(interval)
	defer deadline.Stop()
	defer tick.Stop()

	cur, err := e.trader.GetBalance(ctx, asset)
	if err == nil {
		if baseline == 0 {
			if cur > 0 {
				return cur, nil
			}
		} else if cur > baseline+1e-12 {
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
			return 0, fmt.Errorf("timeout waiting %s balance increase: baseline=%.12f last=%.12f", asset, baseline, last)
		case <-tick.C:
			cur, err := e.trader.GetBalance(ctx, asset)
			if err != nil {
				continue
			}
			if baseline == 0 {
				if cur > 0 {
					return cur, nil
				}
			} else if cur > baseline+1e-12 {
				return cur, nil
			}
		}
	}
}

// parseBaseQuote — простой парсер BASE/QUOTE по суффиксу (USDT/USDC достаточно).
func parseBaseQuote(symbol string) (base, quote string) {
	quotes := []string{"USDT", "USDC", "BTC", "ETH", "EUR", "TRY", "BRL", "RUB"}
	for _, q := range quotes {
		if strings.HasSuffix(symbol, q) && len(symbol) > len(q) {
			return symbol[:len(symbol)-len(q)], q
		}
	}
	return symbol, ""
}




