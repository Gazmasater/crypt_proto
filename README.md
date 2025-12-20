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




