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




Сейчас у тебя три разные проблемы, и они реально “ломают сборку”:


SpotTrader redeclared — потому что интерфейс SpotTrader уже есть в arb/executor.go, а ты ещё раз объявил его в executor_real.go.


SmartMarketBuyUSDT undefined — потому что в arb.SpotTrader нет этих методов, там другой контракт (скорее всего PlaceMarketOrder + GetBalance).


main.go: ты передал filters туда, где по сигнатуре ожидается float64 — значит сигнатуры NewRealExecutor не совпадают с тем, как ты её вызываешь.


Ниже — нормальное решение, которое:


не трогает arb/executor.go (интерфейс SpotTrader остаётся таким, какой он у тебя уже есть),


убирает quoteOrderQty (чтобы не ловить amount scale invalid),


делает нормализацию qty по твоим arb.SymbolFilter (step/minQty + APIEnabled),


прекращает ловить 10007 заранее (через APIEnabled),


и исправляет Oversold (SELL делаем по реальному балансу после BUY).



1) Перепиши arb/executor_real.go полностью

ВАЖНО: НЕ объявляй SpotTrader здесь. Используй тот, что уже в arb/executor.go.

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

// RealExecutor использует arb.SpotTrader (из executor.go) и твои SymbolFilter для нормализации.
type RealExecutor struct {
	trader SpotTrader
	out    io.Writer

	Filters map[string]SymbolFilter

	StartUSDT float64       // фиксированный старт (например 2)
	Cooldown time.Duration  // защита от спама (например 300ms)
	safety   float64        // SELL safety чтобы не ловить Oversold (например 0.995)

	mu       sync.Mutex
	lastExec map[string]time.Time
}

func NewRealExecutor(tr SpotTrader, out io.Writer, filters map[string]SymbolFilter, startUSDT float64) *RealExecutor {
	if startUSDT <= 0 {
		startUSDT = 2
	}
	return &RealExecutor{
		trader:   tr,
		out:      out,
		Filters:  filters,
		StartUSDT: startUSDT,
		Cooldown: 300 * time.Millisecond,
		safety:   0.995,
		lastExec: make(map[string]time.Time),
	}
}

func (e *RealExecutor) Name() string { return "REAL" }

// Execute исполняет только безопасный класс треугольников вида:
// USDT -> X (BUY по quantity)
// X -> Y (SELL по quantity)
// Y -> USDT (SELL по quantity)
// и всё через реальные балансы после каждой ноги.
func (e *RealExecutor) Execute(ctx context.Context, t domain.Triangle, quotes map[string]domain.Quote, startAmount float64) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	start := e.StartUSDT
	if startAmount > 0 {
		start = startAmount
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

	// проверяем, что все 3 символа разрешены через APIEnabled
	for _, leg := range t.Legs {
		sf, ok := e.Filters[leg.Symbol]
		if ok && !sf.APIEnabled {
			return fmt.Errorf("symbol %s not api-enabled (skip)", leg.Symbol)
		}
	}

	fmt.Fprintf(e.out, "  [REAL EXEC] start=%.6f USDT triangle=%s\n", start, t.Name)

	// ---------- LEG 1: BUY (USDT -> A) ----------
	leg1 := t.Legs[0]
	q1, ok := quotes[leg1.Symbol]
	if !ok || q1.Ask <= 0 {
		return fmt.Errorf("no ask for %s", leg1.Symbol)
	}

	// хотим купить asset = leg1.To за USDT (leg1.From должен быть USDT)
	if !strings.EqualFold(leg1.From, "USDT") {
		return fmt.Errorf("leg1 must start from USDT, got %s", leg1.From)
	}

	// qty = USDT / ask, затем нормализуем по step/minQty
	rawQty1 := start / q1.Ask
	qty1 := e.normalizeQty(leg1.Symbol, rawQty1)
	if qty1 <= 0 {
		return fmt.Errorf("leg1: qty<=0 after normalize (raw=%.12f) symbol=%s", rawQty1, leg1.Symbol)
	}

	fmt.Fprintf(e.out, "    [REAL EXEC] leg 1: BUY %s qty=%.8f (raw=%.12f)\n", leg1.Symbol, qty1, rawQty1)
	if _, err := e.trader.PlaceMarketOrder(ctx, leg1.Symbol, "BUY", qty1, 0); err != nil {
		return fmt.Errorf("leg1 error: %w", err)
	}

	// после BUY берём реальный баланс купленного актива
	assetA := leg1.To
	balA, err := e.trader.GetBalance(ctx, assetA)
	if err != nil {
		return fmt.Errorf("leg1 balance error: %w", err)
	}
	if balA <= 0 {
		return fmt.Errorf("leg1 balance=0 for asset %s", assetA)
	}

	// ---------- LEG 2: SELL (A -> B) ----------
	leg2 := t.Legs[1]
	if !strings.EqualFold(leg2.From, assetA) {
		// Это важно: текущий executor рассчитан на цепочку, где leg2 продаёт A.
		// Если у тебя тут BUY-логика (A как quote), надо отдельную реализацию.
		return fmt.Errorf("leg2 expects FROM=%s, got %s (unsupported for now)", assetA, leg2.From)
	}

	sellA := balA * e.safety
	qty2 := e.normalizeQty(leg2.Symbol, sellA)
	if qty2 <= 0 {
		return fmt.Errorf("leg2: qty<=0 after normalize (raw=%.12f) symbol=%s", sellA, leg2.Symbol)
	}

	fmt.Fprintf(e.out, "    [REAL EXEC] leg 2: SELL %s qty=%.8f (raw=%.12f)\n", leg2.Symbol, qty2, sellA)
	if _, err := e.trader.PlaceMarketOrder(ctx, leg2.Symbol, "SELL", qty2, 0); err != nil {
		return fmt.Errorf("leg2 error: %w", err)
	}

	assetB := leg2.To
	balB, err := e.trader.GetBalance(ctx, assetB)
	if err != nil {
		return fmt.Errorf("leg2 balance error: %w", err)
	}
	if balB <= 0 {
		return fmt.Errorf("leg2 balance=0 for asset %s", assetB)
	}

	// ---------- LEG 3: SELL (B -> USDT) ----------
	leg3 := t.Legs[2]
	if !strings.EqualFold(leg3.From, assetB) {
		return fmt.Errorf("leg3 expects FROM=%s, got %s (unsupported for now)", assetB, leg3.From)
	}
	if !strings.EqualFold(leg3.To, "USDT") {
		return fmt.Errorf("leg3 must end in USDT, got %s", leg3.To)
	}

	sellB := balB * e.safety
	qty3 := e.normalizeQty(leg3.Symbol, sellB)
	if qty3 <= 0 {
		return fmt.Errorf("leg3: qty<=0 after normalize (raw=%.12f) symbol=%s", sellB, leg3.Symbol)
	}

	fmt.Fprintf(e.out, "    [REAL EXEC] leg 3: SELL %s qty=%.8f (raw=%.12f)\n", leg3.Symbol, qty3, sellB)
	if _, err := e.trader.PlaceMarketOrder(ctx, leg3.Symbol, "SELL", qty3, 0); err != nil {
		return fmt.Errorf("leg3 error: %w", err)
	}

	usdtBal, _ := e.trader.GetBalance(ctx, "USDT")
	fmt.Fprintf(e.out, "  [REAL EXEC] done triangle %s  USDT_balance=%.6f\n", t.Name, math.Max(usdtBal, 0))
	return nil
}

// normalizeQty — главный фикс “quantity/amount scale is invalid”.
// Использует твой SymbolFilter (StepSize/MinQty).
func (e *RealExecutor) normalizeQty(symbol string, raw float64) float64 {
	if raw <= 0 {
		return 0
	}

	sf, ok := e.Filters[symbol]
	if !ok {
		// если нет фильтра — режем до 8 знаков, чтобы не швырять мусор в API
		return trunc(raw, 8)
	}

	step := sf.StepSize
	minq := sf.MinQty

	q := raw

	if step > 0 {
		q = math.Floor(q/step) * step
		// ещё раз подрежем хвост, чтобы избежать 1.0000000000002
		dec := decimalsFromStep(step)
		q = trunc(q, dec)
	} else {
		q = trunc(q, 8)
	}

	if minq > 0 && q < minq {
		return 0
	}
	return q
}

func decimalsFromStep(step float64) int {
	// грубо: для step=0.001 -> 3, для 0.0001 -> 4 и т.д.
	if step <= 0 {
		return 8
	}
	dec := 0
	for dec < 12 {
		step *= 10
		if math.Abs(step-math.Round(step)) < 1e-9 {
			return dec + 1
		}
		dec++
	}
	return 8
}

func trunc(x float64, dec int) float64 {
	if dec <= 0 {
		return math.Floor(x)
	}
	p := math.Pow10(dec)
	return math.Floor(x*p) / p
}

✅ Это исправит:


SpotTrader redeclared (потому что мы его не объявляем),


SmartMarketBuyUSDT undefined (мы используем PlaceMarketOrder),


Oversold (SELL идёт по балансу),


amount scale invalid (мы вообще не используем quoteOrderQty),


quantity scale invalid (qty нормализуется step’ом).



2) Исправь вызов в cmd/cryptarb/main.go
У тебя сейчас ошибка:

cannot use filters as float64 in argument to arb.NewRealExecutor

Это значит: твоя функция NewRealExecutor ожидает float64 третьим аргументом.
После переписывания выше — сигнатура такая:
arb.NewRealExecutor(tr, arbOut, filters, cfg.TradeAmountUSDT)

То есть в main должно быть так:
consumer.Executor = arb.NewRealExecutor(tr, arbOut, filters, cfg.TradeAmountUSDT)

А не наоборот.

3) Почему компилятор пишет “*mexc.Trader does not implement arb.SpotTrader (missing PlaceMarketOrder)”
Это почти всегда значит, что подпись метода PlaceMarketOrder в интерфейсе arb.SpotTrader отличается от подписи в mexc.Trader.
Твой mexc.Trader.PlaceMarketOrder сейчас:
PlaceMarketOrder(ctx, symbol, side string, quantity, quoteOrderQty float64) (string, error)

Сделай так:


Открой arb/executor.go, посмотри точную сигнатуру интерфейса SpotTrader.


Приведи mexc.Trader к той же сигнатуре.


Чтобы быстро понять что именно не совпадает — добавь в main.go временно:
var _ arb.SpotTrader = (*mexc.Trader)(nil)

Компилятор сам скажет, какой именно метод/сигнатура не совпадает.

Если покажешь содержимое arb/executor.go (кусок с type SpotTrader interface { ... }), я под твою реальную сигнатуру точно подгоню mexc/trader.go (без гаданий).










