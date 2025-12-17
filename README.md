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




export TRADE_AMOUNT_USDT=100
export FEE_PCT=0.04
export SELL_SAFETY=0.995

export TRIANGLES_FILE=triangles_markets.csv
export TRIANGLES_ENRICHED_FILE=triangles_markets_enriched.csv

go run ./cmd/triangles_enrich_mexc




BTC,BRL,BTCBRL,0,0.000000000000,0.000000000000,0,1,0,1,BTC,USDC,BTCUSDC,0,0.000000000000,0.000000000000,1,1,0,2,USDC,BRL,USDCBRL,0,0.000000000000,0.000000000000,0,1,0,3,,,,0.000000,0,spot not allowed
BTC,BRL,BTCBRL,0,0.000000000000,0.000000000000,0,1,0,1,XRP,BTC,XRPBTC,0,0.000000000000,0.000000000000,1,1,0,9,XRP,BRL,XRPBRL,0,0.000000000000,0.000000000000,0,1,0,3,,,,0.000000,0,spot not allowed
BTC,BRL,BTCBRL,0,0.000000000000,0.000000000000,0,1,0,1,SOL,BTC,SOLBTC,0,0.000000000000,0.000000000000,1,1,0,8,SOL,BRL,SOLBRL,0,0.000000000000,0.000000000000,0,1,0,2,,,,0.000000,0,spot not allowed
BTC,BRL,BTCBRL,0,0.000000000000,0.000000000000,0,1,0,1,ETH,BTC,ETHBTC,0,0.000000000000,0.000000000000,1,1,0,6,ETH,BRL,ETHBRL,0,0.000000000000,0.000000000000,0,1,0,2,,,,0.000000,0,spot not allowed
BTC,BRL,BTCBRL,0,0.000000000000,0.000000000000,0,1,0,1,BTC,USDT,BTCUSDT,0,0.000000000000,0.000000000000,0,1,0,2,BRL,USDT,BRLUSDT,0,0.000000000000,0.000000000000,0,1,0,4,,,,0.000000,0,spot not allowed
BTC,BRL,BTCBRL,0,0.000000000000,0.000000000000,0,1,0,1,MX,BTC,MXBTC,0,0.000000000000,0.000000000000,0,1,0,8,MX,BRL,MXBRL,0,0.000000000000,0.000000000000,0,1,0,2,,,,0.000000,0,spot not allowed
USDC,BRL,USDCBRL,0,0.000000000000,0.000000000000,0,1,0,3,XRP,USDC,XRPUSDC,0,0.000000000000,0.000000000000,1,0,0,4,XRP,BRL,XRPBRL,0,0.000000000000,0.000000000000,0,1,0,3,,,,0.000000,0,spot not allowed







