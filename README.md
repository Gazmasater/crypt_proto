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





[ARB] +0.218%  USDT→KITE→USDC→USDT  maxStart=3.2600 USDT (3.2600 USDT)  safeStart=3.2600 USDT (3.2600 USDT) (x1.00)  bottleneck=KITEUSDC
  KITEUSDT (KITE/USDT): bid=0.0898700000 ask=0.0900100000  spread=0.0001400000 (0.15566%)  bidQty=42016.0800 askQty=8736.0000
  KITEUSDC (KITE/USDC): bid=0.0903600000 ask=0.0909700000  spread=0.0006100000 (0.67281%)  bidQty=36.2000 askQty=200.9000
  USDCUSDT (USDC/USDT): bid=0.9998000000 ask=0.9999000000  spread=0.0001000000 (0.01000%)  bidQty=1887866.5400 askQty=45130.3000
  Legs execution with fees:
  [REAL EXEC] start=2.000000 USDT triangle=USDT→KITE→USDC→USDT
    [REAL EXEC] leg 1: BUY KITEUSDT qty=22.21900000 (raw=22.2197533607)
    leg 1: KITEUSDT  3.259992 USDT → 36.200000 KITE  fee=0.01810905 KITE
    leg 2: KITEUSDC  36.200000 KITE → 3.269396 USDC  fee=0.00163552 USDC
    leg 3: USDCUSDT  3.269396 USDC → 3.267108 USDT  fee=0.00163437 USDT

    [REAL EXEC] leg 1 ERROR: mexc order error: status=400 body={"msg":" quantity scale is invalid","code":400}







