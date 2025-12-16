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





func (e *RealExecutor) Execute(ctx context.Context, t domain.Triangle, quotes map[string]domain.Quote, startUSDT float64) error {
	e.out.Printf("  [REAL EXEC] start=%.6f USDT triangle=%s", startUSDT, t.String())

	// --- LEG 1: BUY (quote -> base) ---
	q1 := quotes[t.Leg1.Symbol]
	e.out.Printf("    [REAL EXEC] leg1 PRE: symbol=%s side=BUY byQuote=%.8f ask=%.10f bid=%.10f",
		t.Leg1.Symbol, startUSDT, q1.Ask, q1.Bid)

	// Важно: рассчитать “сырой” baseQty и что получится после округления/фильтров
	rawBase1 := startUSDT / q1.Ask
	base1, dbg1, ok1 := e.qtyForBuyByQuote(t.Leg1.Symbol, rawBase1, q1.Ask) // см. ниже
	e.out.Printf("    [REAL EXEC] leg1 QTY: rawBase=%.12f -> adjBase=%.12f ok=%v %s",
		rawBase1, base1, ok1, dbg1)
	if !ok1 {
		return fmt.Errorf("leg1 blocked: %s", dbg1)
	}

	o1, err := e.tr.PlaceMarketOrderByQuote(ctx, t.Leg1.Symbol, "BUY", startUSDT)
	if err != nil {
		e.out.Printf("    [REAL EXEC] leg1 PLACE ERR: %v", err)
		return err
	}
	e.out.Printf("    [REAL EXEC] leg1 PLACE OK: orderId=%s status=%s", o1.OrderID, o1.Status)

	f1, err := e.tr.WaitFilled(ctx, t.Leg1.Symbol, o1.OrderID, 3*time.Second)
	if err != nil {
		e.out.Printf("    [REAL EXEC] leg1 FILL ERR: %v", err)
		return err
	}
	e.out.Printf("    [REAL EXEC] leg1 FILLED: executedBase=%.12f quoteSpent=%.12f fee=%.12f feeAsset=%s avgPrice=%.10f",
		f1.ExecutedQty, f1.CummulativeQuoteQty, f1.Fee, f1.FeeAsset, f1.AvgPrice)

	baseGot := f1.ExecutedQty
	if baseGot <= 0 {
		return fmt.Errorf("leg1: executedQty=0")
	}

	// --- LEG 2: SELL (base -> quote) ---
	q2 := quotes[t.Leg2.Symbol]
	e.out.Printf("    [REAL EXEC] leg2 PRE: symbol=%s side=SELL baseIn=%.12f bid=%.10f ask=%.10f",
		t.Leg2.Symbol, baseGot, q2.Bid, q2.Ask)

	sell2, dbg2, ok2 := e.qtyForSell(t.Leg2.Symbol, baseGot) // округление по stepSize/minQty/minNotional
	e.out.Printf("    [REAL EXEC] leg2 QTY: rawBase=%.12f -> adjBase=%.12f ok=%v %s",
		baseGot, sell2, ok2, dbg2)
	if !ok2 {
		return fmt.Errorf("leg2 blocked: %s", dbg2)
	}

	o2, err := e.tr.PlaceMarketOrder(ctx, t.Leg2.Symbol, "SELL", sell2)
	if err != nil {
		e.out.Printf("    [REAL EXEC] leg2 PLACE ERR: %v", err)
		return err
	}
	e.out.Printf("    [REAL EXEC] leg2 PLACE OK: orderId=%s status=%s", o2.OrderID, o2.Status)

	f2, err := e.tr.WaitFilled(ctx, t.Leg2.Symbol, o2.OrderID, 3*time.Second)
	if err != nil {
		e.out.Printf("    [REAL EXEC] leg2 FILL ERR: %v", err)
		return err
	}
	e.out.Printf("    [REAL EXEC] leg2 FILLED: executedBase=%.12f quoteOut=%.12f fee=%.12f feeAsset=%s avgPrice=%.10f",
		f2.ExecutedQty, f2.CummulativeQuoteQty, f2.Fee, f2.FeeAsset, f2.AvgPrice)

	quote2 := f2.CummulativeQuoteQty
	if quote2 <= 0 {
		return fmt.Errorf("leg2: quoteOut=0")
	}

	// --- LEG 3: SELL (USDC -> USDT) ---
	q3 := quotes[t.Leg3.Symbol]
	e.out.Printf("    [REAL EXEC] leg3 PRE: symbol=%s side=SELL quoteIn=%.12f bid=%.10f ask=%.10f",
		t.Leg3.Symbol, quote2, q3.Bid, q3.Ask)

	sell3, dbg3, ok3 := e.qtyForSell(t.Leg3.Symbol, quote2)
	e.out.Printf("    [REAL EXEC] leg3 QTY: raw=%.12f -> adj=%.12f ok=%v %s",
		quote2, sell3, ok3, dbg3)
	if !ok3 {
		return fmt.Errorf("leg3 blocked: %s", dbg3)
	}

	o3, err := e.tr.PlaceMarketOrder(ctx, t.Leg3.Symbol, "SELL", sell3)
	if err != nil {
		e.out.Printf("    [REAL EXEC] leg3 PLACE ERR: %v", err)
		return err
	}
	e.out.Printf("    [REAL EXEC] leg3 PLACE OK: orderId=%s status=%s", o3.OrderID, o3.Status)

	f3, err := e.tr.WaitFilled(ctx, t.Leg3.Symbol, o3.OrderID, 3*time.Second)
	if err != nil {
		e.out.Printf("    [REAL EXEC] leg3 FILL ERR: %v", err)
		return err
	}
	e.out.Printf("    [REAL EXEC] leg3 FILLED: executedBase=%.12f quoteOut=%.12f fee=%.12f feeAsset=%s avgPrice=%.10f",
		f3.ExecutedQty, f3.CummulativeQuoteQty, f3.Fee, f3.FeeAsset, f3.AvgPrice)

	finalUSDT := f3.CummulativeQuoteQty
	pnl := finalUSDT - startUSDT
	e.out.Printf("    [REAL EXEC] DONE: start=%.6f final=%.6f pnl=%.6f (%.4f%%)",
		startUSDT, finalUSDT, pnl, (pnl/startUSDT)*100)

	return nil
}




