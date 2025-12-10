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





BOOK_INTERVAL=100ms
SYMBOLS_FILE=triangles_markets.csv
DEBUG=false

# Комиссия (в процентах, одна нога)
FEE_PCT=0.1          # 0.1% = 0.001

# Минимальная прибыль по кругу (в процентах)
MIN_PROFIT_PCT=0.3   # 0.3% = 0.003



[ARB] -0.345%  USDT→ETH→BTC→USDT
  ETHUSDT (ETH/USDT): bid=3315.4400000000 ask=3315.5400000000  spread=0.1000000000 (0.00302%)  bidQty=0.4540 askQty=19.6715
  ETHBTC (ETH/BTC): bid=0.0358660000 ask=0.0358910000  spread=0.0000250000 (0.06968%)  bidQty=0.3320 askQty=0.1290
  BTCUSDT (BTC/USDT): bid=92400.5300000000 ask=92400.6300000000  spread=0.1000000000 (0.00011%)  bidQty=1.8139 askQty=0.0272


