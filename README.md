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




export TRADE_AMOUNT_USDT=10
export FEE_PCT=0.04
export SELL_SAFETY=0.995

export TRIANGLES_FILE=triangles_markets.csv
export TRIANGLES_ENRICHED_FILE=triangles_markets_enriched.csv

go run ./cmd/triangles_enrich_mexc




// outOK — чистый файл ТОЛЬКО проходящих треугольников (без мусора)
outOK := "triangles_markets_enriched.csv"
fOK, err := os.Create(outOK)
if err != nil { log.Fatalf("create %s: %v", outOK, err) }
defer fOK.Close()
wOK := csv.NewWriter(fOK)
defer wOK.Flush()

// outBad — отдельный файл с отбраковкой (почему не прошёл)
outBad := "triangles_markets_excluded.csv"
fBad, err := os.Create(outBad)
if err != nil { log.Fatalf("create %s: %v", outBad, err) }
defer fBad.Close()
wBad := csv.NewWriter(fBad)
defer wBad.Flush()

// Заголовок чистого — ровно как входной triangles_markets.csv
if err := wOK.Write([]string{"base1","quote1","base2","quote2","base3","quote3"}); err != nil {
    log.Fatalf("write ok header: %v", err)
}

// Заголовок excluded — можно шире (минимум причина)
if err := wBad.Write([]string{
    "base1","quote1","base2","quote2","base3","quote3",
    "reason",
}); err != nil {
    log.Fatalf("write excluded header: %v", err)
}

okCount := 0
badCount := 0

for _, row := range rows { // rows — твои прочитанные из triangles_markets.csv строки (6 полей)
    // ======= ТВОЯ ПРОВЕРКА =======
    // ok := true/false
    // reason := "ok" / "minQty" / "minNotional" / "spot not allowed" / "no book" / ...
    ok, reason := checkTriangle(row, tradeAmountUSDT, feePct, sellSafety, rules, books)
    // ============================

    if ok {
        // пишем ТОЛЬКО 6 полей — без мусора
        if err := wOK.Write(row[:6]); err != nil {
            log.Fatalf("write ok row: %v", err)
        }
        okCount++
    } else {
        rec := append(row[:6], reason)
        if err := wBad.Write(rec); err != nil {
            log.Fatalf("write excluded row: %v", err)
        }
        badCount++
    }
}

log.Printf("DONE: ok=%d excluded=%d -> %s (+%s)", okCount, badCount, outOK, outBad)
fmt.Println("Готово, файл:", outOK)





