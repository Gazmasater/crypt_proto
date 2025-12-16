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

	// анти-спам, чтобы один и тот же треугольник не пытался исполняться 50 раз/сек
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

func (e *RealExecutor) Execute(ctx context.Context, t domain.Triangle, quotes map[string]domain.Quote, startUSDT float64) error {
	start := e.StartUSDT
	if startUSDT > 0 {
		start = startUSDT
	}
	if start <= 0 {
		return fmt.Errorf("start<=0")
	}

	// cooldown по имени треугольника
	if e.Cooldown > 0 {
		if last, ok := e.lastExec[t.Name]; ok && time.Since(last) < e.Cooldown {
			return nil
		}
		e.lastExec[t.Name] = time.Now()
	}

	fmt.Fprintf(e.out, "  [REAL EXEC] start=%.6f USDT triangle=%s\n", start, t.Name)

	// ожидаем, что стартовая валюта реально USDT
	curAsset := "USDT"
	curAmount := start

	// LEG 1: BUY (USDT -> A)
	leg1 := t.Legs[0]
	q1, ok := quotes[leg1.Symbol]
	if !ok || q1.Ask <= 0 {
		return fmt.Errorf("no quote/ask for %s", leg1.Symbol)
	}

	if !legMatchesFlow(leg1, curAsset) {
		return fmt.Errorf("leg1 flow mismatch: have=%s leg=%s", curAsset, leg1.Symbol)
	}

	// покупаем на 2 USDT
	fmt.Fprintf(e.out, "    [REAL EXEC] leg 1: BUY %s by USDT=%.6f\n", leg1.Symbol, curAmount)
	_, err := e.trader.SmartMarketBuyUSDT(ctx, leg1.Symbol, curAmount, q1.Ask)
	if err != nil {
		return fmt.Errorf("leg1 error: %w", err)
	}

	// после покупки: узнаём что реально купили (баланс base актива leg1)
	nextAsset1 := leg1.To // USDT->A по идее
	if nextAsset1 == "USDT" {
		nextAsset1 = leg1.From
	}
	bal1, err := e.trader.GetBalance(ctx, nextAsset1)
	if err != nil {
		return fmt.Errorf("leg1 balance error: %w", err)
	}
	curAsset = nextAsset1
	curAmount = bal1
	if curAmount <= 0 {
		return fmt.Errorf("leg1 result balance=0 asset=%s", curAsset)
	}

	// LEG 2: SELL/BUY в зависимости от Dir, но мы делаем проще:
	// Мы всегда хотим перейти curAsset -> nextAsset по leg.Dir.
	leg2 := t.Legs[1]
	if !legMatchesFlow(leg2, curAsset) {
		return fmt.Errorf("leg2 flow mismatch: have=%s leg=%s", curAsset, leg2.Symbol)
	}

	// если leg.Dir>0: From->To это SELL base->quote
	// если leg.Dir<0: From->To это BUY base<-quote (то есть curAsset=quote), мы должны BUY base за quote qty
	// У нас в SpotTrader только:
	// - SmartMarketSellQty(symbol, qty)
	// - SmartMarketBuyUSDT(symbol, usdt, ask) (это только когда quote=USDT)
	//
	// Поэтому для универсальности:
	// - Если мы на leg2 должны ПРОДАТЬ текущий актив -> используем SELL qty.
	// - Если должны КУПИТЬ base за quote и quote != USDT — пока НЕ делаем “BUY quantity” в Executor,
	//   а делаем SELL на другом направлении через правильный symbol (в домене Dir уже отражает направление).
	//
	// Практически для твоих треугольников USDT→X→USDC→USDT:
	// leg2 обычно SELL XUSDC, т.е. SELL qty — подходит.

	fmt.Fprintf(e.out, "    [REAL EXEC] leg 2: SELL %s qty=%.12f\n", leg2.Symbol, curAmount)
	sell2 := curAmount * e.SellSafety
	if sell2 <= 0 {
		return fmt.Errorf("leg2 qty<=0 after safety")
	}
	_, err = e.trader.SmartMarketSellQty(ctx, leg2.Symbol, sell2)
	if err != nil {
		return fmt.Errorf("leg2 error: %w", err)
	}

	// после продажи: баланс следующего актива
	nextAsset2 := leg2.To
	if strings.ToUpper(nextAsset2) == strings.ToUpper(curAsset) {
		nextAsset2 = leg2.From
	}
	bal2, err := e.trader.GetBalance(ctx, nextAsset2)
	if err != nil {
		return fmt.Errorf("leg2 balance error: %w", err)
	}
	curAsset = nextAsset2
	curAmount = bal2
	if curAmount <= 0 {
		return fmt.Errorf("leg2 result balance=0 asset=%s", curAsset)
	}

	// LEG 3: обычно SELL USDCUSDT
	leg3 := t.Legs[2]
	if !legMatchesFlow(leg3, curAsset) {
		return fmt.Errorf("leg3 flow mismatch: have=%s leg=%s", curAsset, leg3.Symbol)
	}

	fmt.Fprintf(e.out, "    [REAL EXEC] leg 3: SELL %s qty=%.12f\n", leg3.Symbol, curAmount)
	sell3 := curAmount * e.SellSafety
	if sell3 <= 0 {
		return fmt.Errorf("leg3 qty<=0 after safety")
	}
	_, err = e.trader.SmartMarketSellQty(ctx, leg3.Symbol, sell3)
	if err != nil {
		return fmt.Errorf("leg3 error: %w", err)
	}

	// финальный USDT баланс (опционально)
	usdtBal, _ := e.trader.GetBalance(ctx, "USDT")
	fmt.Fprintf(e.out, "  [REAL EXEC] done triangle %s  USDT_balance=%.6f\n", t.Name, math.Max(0, usdtBal))
	return nil
}

func legMatchesFlow(leg domain.Leg, have string) bool {
	// leg.From/To — это уже “логическая” цепочка в Triangle
	return strings.EqualFold(leg.From, have) || strings.EqualFold(leg.To, have)
}
