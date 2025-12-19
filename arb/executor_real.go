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
	// ВНИМАНИЕ: в текущем интерфейсе BUY поддержан только "на сумму USDT".
	// Executor ниже использует BUY через quoteQty, но если leg.From != USDT — он вернёт ошибку.
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
			_ = e.executeOnce(req)
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

func (e *RealExecutor) executeOnce(req execReq) error {
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

	t := req.t
	quotes := req.quotes
	startUSDT := req.startUSDT
	triName := req.triName

	e.logf("  [REAL EXEC] start=%.6f USDT triangle=%s", startUSDT, triName)

	// Покажем ноги
	for i, leg := range t.Legs {
		e.logf("    [REAL EXEC] leg%d: sym=%s dir=%d from=%s to=%s", i+1, leg.Symbol, leg.Dir, leg.From, leg.To)
	}

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

	// Текущий “поток”: с какой валютой и суммой идём по ногам.
	curAsset := "USDT"
	curAmount := startUSDT

	// Исполняем 3 ноги по dir/from/to
	for i := 0; i < 3; i++ {
		leg := t.Legs[i]
		sym := strings.TrimSpace(leg.Symbol)
		if sym == "" {
			return fmt.Errorf("leg%d: empty symbol", i+1)
		}

		q, ok := quotes[sym]
		if !ok {
			return fmt.Errorf("leg%d: no quote snapshot for %s", i+1, sym)
		}
		if q.Ask <= 0 || q.Bid <= 0 {
			return fmt.Errorf("leg%d: bad quote for %s (ask=%.10f bid=%.10f)", i+1, sym, q.Ask, q.Bid)
		}

		from := strings.ToUpper(strings.TrimSpace(leg.From))
		to := strings.ToUpper(strings.TrimSpace(leg.To))
		if from == "" || to == "" {
			return fmt.Errorf("leg%d: empty from/to (from=%q to=%q)", i+1, leg.From, leg.To)
		}

		if curAsset != from {
			// не фейлим сразу — но это почти всегда признак рассинхрона описания треугольника/исполнения
			e.logf("    [REAL EXEC] WARN leg%d: curAsset=%s curAmount=%.12f but leg.From=%s",
				i+1, curAsset, curAmount, from)
		}

		// Балансы до
		var fromBefore, toBefore float64
		{
			done := e.step(fmt.Sprintf("GetBalance %s (before leg%d)", from, i+1))
			v, err := e.trader.GetBalance(req.ctx, from)
			done()
			if err != nil {
				e.logf("    [REAL EXEC] BAL ERR: get %s before leg%d: %v", from, i+1, err)
				return err
			}
			fromBefore = v
		}
		{
			done := e.step(fmt.Sprintf("GetBalance %s (before leg%d)", to, i+1))
			v, err := e.trader.GetBalance(req.ctx, to)
			done()
			if err != nil {
				e.logf("    [REAL EXEC] BAL ERR: get %s before leg%d: %v", to, i+1, err)
				return err
			}
			toBefore = v
		}

		if leg.Dir < 0 {
			// BUY: тратим quote (from), получаем base (to)
			// В текущем интерфейсе SpotTrader BUY поддержан только в USDT.
			spend := curAmount
			if spend <= 0 {
				return fmt.Errorf("leg%d BUY: spend<=0 (%s)", i+1, from)
			}
			if fromBefore+1e-9 < spend {
				return fmt.Errorf("leg%d BUY: insufficient %s: have=%.12f need=%.12f", i+1, from, fromBefore, spend)
			}
			if from != "USDT" {
				// ВАЖНО: пока твой трейдер умеет BUY только за USDT.
				// Чтобы торговать BUY за USDC/прочее — надо расширить интерфейс на SmartMarketBuyQuote.
				return fmt.Errorf("leg%d BUY: quote asset is %s, but SpotTrader supports BUY only by USDT (need SmartMarketBuyQuote)", i+1, from)
			}

			e.logf("    [REAL EXEC] leg%d PRE: BUY %s spend=%s=%.6f ask=%.10f bid=%.10f | %s before=%.12f %s before=%.12f",
				i+1, sym, from, spend, q.Ask, q.Bid, from, fromBefore, to, toBefore)

			var ord string
			{
				orderCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

				done := e.step(fmt.Sprintf("SmartMarketBuyUSDT leg%d", i+1))
				id, err := e.trader.SmartMarketBuyUSDT(orderCtx, sym, spend, q.Ask)
				done()
				if err != nil {
					e.logf("    [REAL EXEC] leg%d PLACE ERR (BUY): %v", i+1, err)
					return err
				}
				ord = id
			}
			e.logf("    [REAL EXEC] leg%d PLACE OK: orderId=%s", i+1, ord)

			var toAfter float64
			{
				done := e.step(fmt.Sprintf("waitBalanceChange %s (after leg%d)", to, i+1))
				v, err := e.waitBalanceChange(req.ctx, to, toBefore, 3*time.Second, 150*time.Millisecond)
				done()
				if err != nil {
					e.logf("    [REAL EXEC] leg%d WAIT BAL ERR (%s): %v", i+1, to, err)
					return err
				}
				toAfter = v
			}

			delta := toAfter - toBefore
			e.logf("    [REAL EXEC] leg%d BAL after: %s=%.12f delta=%.12f", i+1, to, toAfter, delta)
			if delta <= 0 {
				return fmt.Errorf("leg%d BUY: %s did not increase (before=%.12f after=%.12f)", i+1, to, toBefore, toAfter)
			}

			curAsset = to
			curAmount = delta
			continue
		}

		// SELL: продаём base qty (from), получаем quote (to)
		qtyRaw := curAmount
		// На SELL берём баланс from, а не curAmount, потому что фактический qty после BUY может отличаться,
		// и самый надёжный путь — продать то, что есть на балансе (с safety).
		qty := fromBefore * e.SellSafety
		if qty <= 0 {
			return fmt.Errorf("leg%d SELL: qty<=0 (%s=%.12f safety=%.6f)", i+1, from, fromBefore, e.SellSafety)
		}

		e.logf("    [REAL EXEC] leg%d PRE: SELL %s qty=%s=%.12f (curAmount=%.12f raw=%.12f) bid=%.10f ask=%.10f | %s before=%.12f %s before=%.12f",
			i+1, sym, from, qty, curAmount, qtyRaw, q.Bid, q.Ask, from, fromBefore, to, toBefore)

		var ord string
		{
			orderCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			done := e.step(fmt.Sprintf("SmartMarketSellQty leg%d", i+1))
			id, err := e.trader.SmartMarketSellQty(orderCtx, sym, qty)
			done()
			if err != nil {
				e.logf("    [REAL EXEC] leg%d PLACE ERR (SELL): %v", i+1, err)
				return err
			}
			ord = id
		}
		e.logf("    [REAL EXEC] leg%d PLACE OK: orderId=%s", i+1, ord)

		var toAfter float64
		{
			done := e.step(fmt.Sprintf("waitBalanceChange %s (after leg%d)", to, i+1))
			v, err := e.waitBalanceChange(req.ctx, to, toBefore, 3*time.Second, 150*time.Millisecond)
			done()
			if err != nil {
				e.logf("    [REAL EXEC] leg%d WAIT BAL ERR (%s): %v", i+1, to, err)
				return err
			}
			toAfter = v
		}

		delta := toAfter - toBefore
		e.logf("    [REAL EXEC] leg%d BAL after: %s=%.12f delta=%.12f", i+1, to, toAfter, delta)
		if delta <= 0 {
			return fmt.Errorf("leg%d SELL: %s did not increase (before=%.12f after=%.12f)", i+1, to, toBefore, toAfter)
		}

		curAsset = to
		curAmount = delta
	}

	// Финальный баланс USDT
	var usdtAfter float64
	{
		done := e.step("GetBalance USDT (after)")
		v, err := e.trader.GetBalance(req.ctx, "USDT")
		done()
		if err != nil {
			e.logf("    [REAL EXEC] BAL ERR: get USDT after: %v", err)
			return err
		}
		usdtAfter = v
	}

	dUSDTTotal := usdtAfter - usdt0
	e.logf("    [REAL EXEC] DONE: curAsset=%s curAmount=%.12f", curAsset, curAmount)
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
