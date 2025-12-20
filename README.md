apikey = "4333ed4b-cd83-49f5-97d1-c399e2349748"
secretkey = "E3848531135EDB4CCFDA0F1BC14CD274"
IP = ""
Название API-ключа = "Arb"
Доступы = "Чтение"



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



025/12/20 13:12:54.522816 CFG: OKX_BOOKS_URL=https://www.okx.com/api/v5/market/books
2025/12/20 13:12:54.522818 CFG: HTTP_TIMEOUT=25s
2025/12/20 13:12:54.522828 CFG: START_AMT=25.00000000
2025/12/20 13:12:54.522835 CFG: START_CCY=USDT
2025/12/20 13:12:54.522836 CFG: KEEP_FAILED=false
2025/12/20 13:12:54.522867 OK: loaded blacklist: 0 symbols
2025/12/20 13:12:54.837208 OK: OKX instruments loaded: 693
2025/12/20 13:12:54.837300 OK: startup blacklist added=0 total=0
2025/12/20 13:12:54.837600 OK: saved blacklist -> blacklist_symbols_okx.txt
2025/12/20 13:16:13.171776 OK: read=788 written=766 skippedNoSymbol=0 skippedBlacklisted=0 skippedNotEligible=0 keep_failed=false -> triangles_usdt_routes_market_okx.csv
2025/12/20 13:16:13.171829 DONE


