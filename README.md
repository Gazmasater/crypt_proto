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



# ===== INPUT / OUTPUT =====
INPUT_CSV=triangles_usdt_routes_okx.csv
OUTPUT_CSV=triangles_usdt_routes_market_okx.csv
BLACKLIST_FILE=blacklist_symbols_okx.txt

# ===== OKX API =====
OKX_INSTRUMENTS_URL=https://www.okx.com/api/v5/public/instruments?instType=SPOT
OKX_BOOKS_URL=https://www.okx.com/api/v5/market/books

# ===== RUNTIME CFG =====
HTTP_TIMEOUT=25s

START_AMT=25
START_CCY=USDT

# комиссия на 1 ногу (taker), в процентах
FEE_PCT=0.1

# минимальная прибыль по кругу (после комиссий), в процентах
MIN_PROFIT_PCT=0.3

# 1 = писать все строки + fail_reason (для дебага)
# 0 = писать только прошедшие фильтры
KEEP_FAILED=0

2025/12/20 14:00:06.159091 CFG: OKX_BOOKS_URL=https://www.okx.com/api/v5/market/books
2025/12/20 14:00:06.159092 CFG: HTTP_TIMEOUT=25s
2025/12/20 14:00:06.159102 CFG: START_AMT=25.00000000
2025/12/20 14:00:06.159110 CFG: START_CCY=USDT
2025/12/20 14:00:06.159112 CFG: KEEP_FAILED=false
2025/12/20 14:00:06.159115 OK: loaded blacklist: 0 symbols
2025/12/20 14:00:06.563296 OK: OKX instruments loaded: 693
2025/12/20 14:00:06.563394 OK: startup blacklist added=0 total=0
2025/12/20 14:00:06.563582 OK: saved blacklist -> blacklist_symbols_okx.txt
2025/12/20 14:03:37.801421 OK: read=788 written=766 skippedNoSymbol=0 skippedBlacklisted=0 skippedNotEligible=0 skippedNotFeasible(start=25.00 USDT)=22 keep_failed=false -> triangles_usdt_routes_market_okx.csv
2025/12/20 14:03:37.801514 DONE
gaz358@gaz358-BOD-WXX9:~/myprog/cr


