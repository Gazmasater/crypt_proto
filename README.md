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



[ARB] +0.149%  USDT→NAKA→USDC→USDT  maxStart=6.8235 USDT (6.8235 USDT)  safeStart=6.8235 USDT (6.8235 USDT) (x1.00)  bottleneck=NAKAUSDC
  NAKAUSDT (NAKA/USDT): bid=0.0787800000 ask=0.0788800000  spread=0.0001000000 (0.12686%)  bidQty=73.2300 askQty=822.5200
  NAKAUSDC (NAKA/USDC): bid=0.0791000000 ask=0.0792000000  spread=0.0001000000 (0.12634%)  bidQty=86.4700 askQty=36.6700
  USDCUSDT (USDC/USDT): bid=0.9999000000 ask=1.0000000000  spread=0.0001000000 (0.01000%)  bidQty=1992412.8000 askQty=69169.1500

  [REAL EXEC] start=2.000000 USDT triangle=USDT→NAKA→USDC→USDT
    [REAL EXEC] leg 1: BUY NAKAUSDT quoteOrderQty=2.000000
    [REAL EXEC] leg 1 ERROR: mexc order error: status=400 body={"msg":" amount scale is invalid","code":400}
^C2025/12/16 07:18:37.206853 shutting down...
2025/12/16 07:18:37.216017 bye







