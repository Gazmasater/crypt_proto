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





func (c *Consumer) printTriangle(
	ts time.Time,
	t domain.Triangle,
	profit float64,
	quotes map[string]domain.Quote,
	ms *domain.MaxStartInfo,
	startFraction float64,
) {
	w := c.writer
	fmt.Fprintf(w, "%s\n", ts.Format("2006-01-02 15:04:05.000"))

	if ms != nil {
		bneckSym := ""
		if ms.BottleneckLeg >= 0 && ms.BottleneckLeg < len(t.Legs) {
			bneckSym = t.Legs[ms.BottleneckLeg].Symbol
		}

		// safeStart = доля от maxStart
		safeStart := ms.MaxStart * startFraction

		// конвертация maxStart/safeStart в USDT для вывода
		maxUSDT, okMax := convertToUSDT(ms.MaxStart, ms.StartAsset, quotes)
		safeUSDT, okSafe := convertToUSDT(safeStart, ms.StartAsset, quotes)

		maxUSDTStr := "?"
		safeUSDTStr := "?"
		if okMax {
			maxUSDTStr = fmt.Sprintf("%.4f", maxUSDT)
		}
		if okSafe {
			safeUSDTStr = fmt.Sprintf("%.4f", safeUSDT)
		}

		fmt.Fprintf(w,
			"[ARB] %+0.3f%%  %s  maxStart=%.4f %s (%s USDT)  safeStart=%.4f %s (%s USDT) (x%.2f)  bottleneck=%s\n",
			profit*100, t.Name,
			ms.MaxStart, ms.StartAsset, maxUSDTStr,
			safeStart, ms.StartAsset, safeUSDTStr,
			startFraction,
			bneckSym,
		)
	} else {
		fmt.Fprintf(w, "[ARB] %+0.3f%%  %s\n", profit*100, t.Name)
	}

	for _, leg := range t.Legs {
		q := quotes[leg.Symbol]
		mid := (q.Bid + q.Ask) / 2
		spreadAbs := q.Ask - q.Bid
		spreadPct := 0.0
		if mid > 0 {
			spreadPct = spreadAbs / mid * 100
		}

		// side для красоты
		side := ""
		if leg.Dir > 0 {
			side = fmt.Sprintf("%s/%s", leg.From, leg.To)
		} else {
			side = fmt.Sprintf("%s/%s", leg.To, leg.From)
		}

		// Определяем base/quote для стакана
		var baseAsset, quoteAsset string
		if leg.Dir > 0 {
			// From = base, To = quote
			baseAsset = leg.From
			quoteAsset = leg.To
		} else {
			// To = base, From = quote
			baseAsset = leg.To
			quoteAsset = leg.From
		}

		// ===== bidQty / askQty в USDT (по base) =====
		bidQtyUSDT, okBidBase := convertToUSDT(q.BidQty, baseAsset, quotes)
		askQtyUSDT, okAskBase := convertToUSDT(q.AskQty, baseAsset, quotes)

		bidQtyUSDTStr := "?"
		askQtyUSDTStr := "?"
		if okBidBase {
			bidQtyUSDTStr = fmt.Sprintf("%.4f", bidQtyUSDT)
		}
		if okAskBase {
			askQtyUSDTStr = fmt.Sprintf("%.4f", askQtyUSDT)
		}

		// ===== Денежный объём (notional) в котируемой валюте и в USDT =====
		// bid: сколько котируемой валюты можно получить за весь bidQty по bid
		bidNotionalQuote := q.BidQty * q.Bid
		askNotionalQuote := q.AskQty * q.Ask

		bidNotionalUSDT, okBidNot := convertToUSDT(bidNotionalQuote, quoteAsset, quotes)
		askNotionalUSDT, okAskNot := convertToUSDT(askNotionalQuote, quoteAsset, quotes)

		bidNotionalQuoteStr := fmt.Sprintf("%.4f", bidNotionalQuote)
		askNotionalQuoteStr := fmt.Sprintf("%.4f", askNotionalQuote)

		bidNotionalUSDTStr := "?"
		askNotionalUSDTStr := "?"
		if okBidNot {
			bidNotionalUSDTStr = fmt.Sprintf("%.4f", bidNotionalUSDT)
		}
		if okAskNot {
			askNotionalUSDTStr = fmt.Sprintf("%.4f", askNotionalUSDT)
		}

		fmt.Fprintf(
			w,
			"  %s (%s): bid=%.10f ask=%.10f  spread=%.10f (%.5f%%)\n"+
				"    bidQty=%.4f %s (≈%s USDT, notional %s %s ≈%s USDT)\n"+
				"    askQty=%.4f %s (≈%s USDT, notional %s %s ≈%s USDT)\n",
			leg.Symbol, side,
			q.Bid, q.Ask,
			spreadAbs, spreadPct,
			q.BidQty, baseAsset, bidQtyUSDTStr,
			bidNotionalQuoteStr, quoteAsset, bidNotionalUSDTStr,
			q.AskQty, baseAsset, askQtyUSDTStr,
			askNotionalQuoteStr, quoteAsset, askNotionalUSDTStr,
		)
	}

	fmt.Fprintln(w)
}


