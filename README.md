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





[ARB] +0.716%  BOMB→USDC→USDT→BOMB  maxStart=6896.0000 BOMB (1.1861 USDT)  safeStart=3448.0000 BOMB (0.5931 USDT) (x0.50)  bottleneck=BOMBUSDC
  BOMBUSDC (BOMB/USDC): bid=0.0001744000 ask=0.0001776000  spread=0.0000032000 (1.81818%)
    bidQty=6896.0000 BOMB (≈1.1861 USDT, notional 1.2027 USDC ≈1.2023 USDT)
    askQty=7568.6500 BOMB (≈1.3018 USDT, notional 1.3442 USDC ≈1.3438 USDT)
  USDCUSDT (USDC/USDT): bid=0.9997000000 ask=0.9998000000  spread=0.0001000000 (0.01000%)
    bidQty=74113.8500 USDC (≈74091.6158 USDT, notional 74091.6158 USDT ≈74091.6158 USDT)
    askQty=69572.1400 USDC (≈69551.2684 USDT, notional 69558.2256 USDT ≈69558.2256 USDT)
  BOMBUSDT (BOMB/USDT): bid=0.0001720000 ask=0.0001729000  spread=0.0000009000 (0.52189%)
    bidQty=162787.8700 BOMB (≈27.9995 USDT, notional 27.9995 USDT ≈27.9995 USDT)
    askQty=98595.8600 BOMB (≈16.9585 USDT, notional 17.0472 USDT ≈17.0472 USDT)


