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




1. Обновлённый simulateTriangleExecution в arb.go
В файле arb/arb.go замени определение legExec и функции simulateTriangleExecution на этот вариант:
type legExec struct {
	Symbol    string
	From      string
	To        string
	AmountIn  float64
	AmountOut float64
	FeeAmount float64
	FeeAsset  string
}

// simulateTriangleExecution:
// - идём строго по направлению треугольника (leg.From -> leg.To),
// - Dir > 0: base -> quote, сделка по bid (sell base),
// - Dir < 0: quote -> base, сделка по ask (buy base),
// - считаем реальные объёмы на каждой ноге.
func simulateTriangleExecution(
	t domain.Triangle,
	quotes map[string]domain.Quote,
	startAsset string,
	startAmount float64,
	feePerLeg float64,
) ([]legExec, bool) {
	if startAmount <= 0 {
		return nil, false
	}

	curAsset := startAsset
	curAmount := startAmount
	res := make([]legExec, 0, len(t.Legs))

	for i := 0; i < len(t.Legs); i++ {
		leg := t.Legs[i]
		q, ok := quotes[leg.Symbol]
		if !ok || q.Bid <= 0 || q.Ask <= 0 {
			return nil, false
		}

		// Должны двигаться по направлению треугольника.
		if curAsset != leg.From {
			return nil, false
		}

		amountIn := curAmount
		var gross, feeAmount, amountOut float64
		var feeAsset string

		if leg.Dir > 0 {
			// Dir > 0: leg.From = base, leg.To = quote
			// Продаём base -> получаем quote по bid
			gross = amountIn * q.Bid
			feeAmount = gross * feePerLeg
			amountOut = gross - feeAmount
			feeAsset = leg.To
		} else {
			// Dir < 0: leg.From = quote, leg.To = base
			// Покупаем base за quote по ask
			gross = amountIn / q.Ask
			feeAmount = gross * feePerLeg
			amountOut = gross - feeAmount
			feeAsset = leg.To
		}

		if amountOut <= 0 {
			return nil, false
		}

		res = append(res, legExec{
			Symbol:    leg.Symbol,
			From:      leg.From,
			To:        leg.To,
			AmountIn:  amountIn,
			AmountOut: amountOut,
			FeeAmount: feeAmount,
			FeeAsset:  feeAsset,
		})

		curAsset = leg.To
		curAmount = amountOut
	}

	// В идеале curAsset должен совпасть со стартовым активом (замкнутый треугольник),
	// но жёстко не проверяем, чтобы не ломать треугольники с "грязными" данными.
	return res, true
}


Важно: больше не используем detectBaseQuote внутри симуляции — работаем строго по leg.From/leg.To и Dir.

DryRunExecutor можно оставить как есть — он просто печатает AmountIn/AmountOut и не влияет на торговлю.

2. Обновлённый RealExecutor (если хочешь без сюрпризов)
Если ты вставлял executor_real.go по моему прошлому варианту, там логика завязана на detectBaseQuote и старый смысл legExec. Лучше её тоже синхронизировать с новой симуляцией.
Файл arb/executor_real.go перепиши целиком вот так:
package arb

import (
	"context"
	"fmt"
	"io"
	"time"

	"crypt_proto/domain"
)

// SpotTrader — минимальный интерфейс для трейдера биржи.
type SpotTrader interface {
	PlaceMarket(ctx context.Context, symbol, side string, quantity float64) error
}

// RealExecutor — реальный исполнитель треугольника через SpotTrader.
type RealExecutor struct {
	trader         SpotTrader
	out            io.Writer
	FixedStartUSDT float64 // если >0 и StartAsset=USDT — ограничиваем старт этим значением
}

func NewRealExecutor(trader SpotTrader, out io.Writer, fixedStartUSDT float64) *RealExecutor {
	return &RealExecutor{
		trader:         trader,
		out:            out,
		FixedStartUSDT: fixedStartUSDT,
	}
}

func (e *RealExecutor) log(format string, args ...any) {
	if e.out == nil {
		return
	}
	fmt.Fprintf(e.out, format+"\n", args...)
}

// ExecuteTriangle запускает реальную торговлю по треугольнику.
// Предполагаем, что StartAsset = USDT (BaseAsset) и фильтры уже пройдены.
func (e *RealExecutor) ExecuteTriangle(
	ctx context.Context,
	t domain.Triangle,
	quotes map[string]domain.Quote,
	ms *domain.MaxStartInfo,
	startFraction float64,
) {
	if e.trader == nil || ms == nil {
		return
	}

	if ms.StartAsset != BaseAsset {
		e.log("  [REAL EXEC] skip triangle %s: StartAsset=%s != %s",
			t.Name, ms.StartAsset, BaseAsset)
		return
	}

	// Базовый старт из ликвидности
	startAmount := ms.MaxStart * startFraction
	if startAmount <= 0 {
		return
	}

	// Если задан фиксированный старт — ограничиваем сверху
	if e.FixedStartUSDT > 0 && ms.StartAsset == BaseAsset && startAmount > e.FixedStartUSDT {
		startAmount = e.FixedStartUSDT
	}

	// Считаем путь исполнения для выбранного стартового объёма
	execs, ok := simulateTriangleExecution(t, quotes, ms.StartAsset, startAmount, 0)
	if !ok || len(execs) == 0 {
		e.log("  [REAL EXEC] simulate failed for %s", t.Name)
		return
	}

	e.log("  [REAL EXEC] start=%.6f %s triangle=%s", startAmount, ms.StartAsset, t.Name)

	// Локальный таймаут на всю цепочку
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	for i, lg := range execs {
		leg := t.Legs[i]

		// По новой семантике:
		//  Dir > 0: From=base, To=quote, мы продаём base (SELL), qty = AmountIn
		//  Dir < 0: From=quote, To=base, мы покупаем base (BUY), qty = AmountOut
		side := ""
		var qty float64

		if leg.Dir > 0 {
			side = "SELL"
			qty = lg.AmountIn
		} else {
			side = "BUY"
			qty = lg.AmountOut
		}

		if qty <= 0 {
			e.log("    [REAL EXEC] leg %d: qty<=0, skip (%s %s)", i+1, side, lg.Symbol)
			return
		}

		e.log("    [REAL EXEC] leg %d: %s %s qty=%.8f", i+1, side, lg.Symbol, qty)

		if err := e.trader.PlaceMarket(ctx, lg.Symbol, side, qty); err != nil {
			e.log("    [REAL EXEC] leg %d ERROR: %v", i+1, err)
			// На первой же ошибке останавливаем весь треугольник.
			return
		}
	}

	e.log("  [REAL EXEC] done triangle %s", t.Name)
}

Здесь:


RealExecutor полностью опирается на новую семантику simulateTriangleExecution;


side и qty считаются по Dir и по AmountIn/AmountOut;


мы больше не лезем в detectBaseQuote / base / quote — всё берём из треугольника.



3. Почему теперь SUPRA не должен падать
Для треугольника:


USDT→SUPRA (SUPRAUSDT, Dir = -1),


SUPRA→USDC (SUPRAUSDC, Dir = +1),


USDC→USDT (USDCUSDT, Dir = +1),


новая логика:


Нога 1:


curAsset = USDT, leg.From = USDT → ок


Dir < 0 → покупаем SUPRA за USDT по ask SUPRAUSDT


AmountIn = USDT, AmountOut = SUPRA




Нога 2:


curAsset = SUPRA, leg.From = SUPRA → ок


Dir > 0 → продаём SUPRA за USDC по bid SUPRAUSDC




Нога 3:


curAsset = USDC, leg.From = USDC → ок


Dir > 0 → продаём USDC за USDT по bid USDCUSDT




Ошибки simulate failed уже не будет — если все котировки есть и положительные.

Если хочешь, следующий шаг можем сделать:


маленький лог в simulateTriangleExecution, чтобы при false было видно leg, на которой всё сломалось — но пока достаточно этой правки, чтобы SUPRA и похожие треугольники перестали отваливаться.



