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





1.2. Новый normalizeQty

Где-нибудь рядом с RealExecutor/legExec добавь/замени функцию на такую:

// normalizeQty округляет объём вниз до допустимого количества знаков.
// Для дорогих активов, например BTC, даём больше знаков, чтобы не уйти в 0.
// Для остального оставляем более грубую точность, чтобы не ловить "quantity scale is invalid".
func normalizeQty(symbol string, q float64) float64 {
	if q <= 0 {
		return 0
	}

	// по умолчанию 3 знака (подходит для TON, GIGGLE и т.п.)
	decimals := 3

	// для BTC-пар даём 6 знаков (BTCUSDT, BTCUSDC и т.п.)
	if strings.HasPrefix(symbol, "BTC") {
		decimals = 6
	}

	factor := math.Pow10(decimals)
	return math.Floor(q*factor) / factor
}

2. Используем normalizeQty в ExecuteTriangle

В RealExecutor.ExecuteTriangle у тебя сейчас логика типа:

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


Замени вызов normalizeQty на вариант с символом:

    qty := normalizeQty(lg.Symbol, qtyRaw)


Остальное можно оставить как есть.







