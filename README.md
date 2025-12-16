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





1. Исправить формат quantity в запросах (MEXC)
Файл: mexc/trader.go
Функция: PlaceMarket
Убедись, что вот это:
params.Set("quantity", fmt.Sprintf("%.8f", quantity))

ЗАМЕНЕНО на вот это:
qtyStr := strconv.FormatFloat(quantity, 'f', -1, 64)
params.Set("quantity", qtyStr)

И вверху файла есть импорт:
import (
    // ...
    "strconv"
    // ...
)


Важно: пересобери проект после правки (go build ./cmd/cryptarb), а не просто сохрани файл.

Это уберёт формат "2.00000000" → "2".

2. Огрубить количество до разумного числа знаков (чтобы не было 1.33511348)
Для таких монет как TON/GIGGLE биржа не любит 8 знаков после запятой в quantity (скорее всего stepSize = 0.001 или около того).
Сделаем универсальный костыль: обрезаем до 3 знаков после запятой перед отправкой ордера.
2.1. Добавь helper в arb/executor_real.go
Вверху файла добавь импорт:
import (
    // ...
    "math"
    // ...
)

Где-нибудь под структурами RealExecutor/legExec добавь функцию:
// normalizeQty округляет объём вниз до 3 знаков после запятой.
// Это грубый, но безопасный костыль против "quantity scale is invalid".
func normalizeQty(q float64) float64 {
	if q <= 0 {
		return 0
	}
	const decimals = 3
	factor := math.Pow10(decimals)
	return math.Floor(q*factor) / factor
}

2.2. Используй её в ExecuteTriangle
В func (e *RealExecutor) ExecuteTriangle(...) найди участок, где считаются side и qty:
Сейчас у тебя примерно так:
for i, lg := range execs {
    leg := t.Legs[i]

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
        ...
    }
}

Замени на:
for i, lg := range execs {
    leg := t.Legs[i]

    side := ""
    var qtyRaw float64

    if leg.Dir > 0 {
        side = "SELL"
        qtyRaw = lg.AmountIn
    } else {
        side = "BUY"
        qtyRaw = lg.AmountOut
    }

    qty := normalizeQty(qtyRaw)
    if qty <= 0 {
        e.log("    [REAL EXEC] leg %d: qty<=0 after normalize (raw=%.10f) (%s %s)", i+1, qtyRaw, side, lg.Symbol)
        return
    }

    e.log("    [REAL EXEC] leg %d: %s %s qty=%.8f (raw=%.10f)", i+1, side, lg.Symbol, qty, qtyRaw)

    if err := e.trader.PlaceMarket(ctx, lg.Symbol, side, qty); err != nil {
        e.log("    [REAL EXEC] leg %d ERROR: %v", i+1, err)
        return
    }
}

Теперь:


Для USDCUSDT: raw 2.0 → normalizeQty(2.0) = 2 → ок.


Для TONUSDC: raw 1.33511348 → 1.335 → scale=3, биржа должна принять.



3. Дальше (чуть позже)
Сейчас не трогаем, но запомни:


Реальный stepSize/minQty для каждой пары всё равно лучше брать из exchangeInfo (у тебя раньше уже была SymbolFilter и stepSize).


Многоразовое [REAL EXEC] leg 1/2 за пару миллисекунд — это отдельная тема (антиспам на исполнение, типа lastExec[triangleID] с интервалом 300–500 ms). Можно добавить после того, как ордера начнут проходить.



Итого, что сделать прямо сейчас:


В mexc/trader.go — убедиться, что quantity формируется через strconv.FormatFloat(..., -1, 64).


В arb/executor_real.go — добавить normalizeQty() и использовать его в ExecuteTriangle для qty.


После этого пересобрать и снова запустить — в логах уже должна уйти ошибка quantity scale is invalid, либо смениться на более “рыночную” (NO_SUFFICIENT_BALANCE, min notional и т.п.).






