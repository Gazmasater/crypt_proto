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





2025-12-17 01:49:36.658
[ARB] +0.027%  USDT→USDC→VANRY→USDT  maxStart=10.4886 USDT (10.4886 USDT)  safeStart=10.4886 USDT (10.4886 USDT) (x1.00)  bottleneck=VANRYUSDC
  USDCUSDT (USDC/USDT): bid=1.0000000000 ask=1.0001000000  spread=0.0001000000 (0.01000%)  bidQty=124316.4700 askQty=213421.9500
  VANRYUSDC (VANRY/USDC): bid=0.0082710000 ask=0.0082750000  spread=0.0000040000 (0.04835%)  bidQty=6030.2500 askQty=1266.8700
  VANRYUSDT (VANRY/USDT): bid=0.0082880000 ask=0.0083000000  spread=0.0000120000 (0.14468%)  bidQty=20781.8400 askQty=1835.8500

  [REAL EXEC] start=2.000000 USDT triangle=USDT→USDC→VANRY→USDT
    [REAL EXEC] legs: sym1=USDCUSDT sym2=VANRYUSDC sym3=VANRYUSDT
    [REAL EXEC] parsed: sym1=USDCUSDT (USDC/USDT) sym2=VANRYUSDC (VANRY/USDC) sym3=VANRYUSDT (VANRY/USDT)
    [REAL EXEC] >>> GetBalance USDT (before)
2025-12-17 01:49:36.681
[ARB] -0.021%  USDT→USDC→VANRY→USDT  maxStart=10.4886 USDT (10.4886 USDT)  safeStart=10.4886 USDT (10.4886 USDT) (x1.00)  bottleneck=VANRYUSDC
  USDCUSDT (USDC/USDT): bid=1.0000000000 ask=1.0001000000  spread=0.0001000000 (0.01000%)  bidQty=124316.4700 askQty=213421.9500
  VANRYUSDC (VANRY/USDC): bid=0.0082710000 ask=0.0082750000  spread=0.0000040000 (0.04835%)  bidQty=6030.2500 askQty=1266.8700
  VANRYUSDT (VANRY/USDT): bid=0.0082840000 ask=0.0083000000  spread=0.0000160000 (0.19296%)  bidQty=43991.4400 askQty=1835.8500

2025-12-17 01:49:36.720
[ARB] +0.003%  USDT→USDC→VANRY→USDT  maxStart=10.4886 USDT (10.4886 USDT)  safeStart=10.4886 USDT (10.4886 USDT) (x1.00)  bottleneck=VANRYUSDC
  USDCUSDT (USDC/USDT): bid=1.0000000000 ask=1.0001000000  spread=0.0001000000 (0.01000%)  bidQty=124316.4700 askQty=213418.1000
  VANRYUSDC (VANRY/USDC): bid=0.0082700000 ask=0.0082750000  spread=0.0000050000 (0.06044%)  bidQty=2300.8500 askQty=1266.8700
  VANRYUSDT (VANRY/USDT): bid=0.0082860000 ask=0.0083010000  spread=0.0000150000 (0.18086%)  bidQty=20781.8400 askQty=1885.7100

    [REAL EXEC] <<< GetBalance USDT (before) (408ms)
    [REAL EXEC] BAL before: USDT=35.782830681146
    [REAL EXEC] >>> GetBalance USDC (before leg1)
    [REAL EXEC] <<< GetBalance USDC (before leg1) (241ms)
    [REAL EXEC] leg1 PRE: BUY USDCUSDT by USDT=2.000000 ask=1.0001000000 bid=1.0000000000 | USDC before=0.001070300000
    [REAL EXEC] >>> SmartMarketBuyUSDT leg1
^C2025/12/17 01:49:37.411025 shutting down...
    [REAL EXEC] <<< SmartMarketBuyUSDT leg1 (102ms)
    [REAL EXEC] leg1 PLACE ERR: Post "https://api.mexc.com/api/v3/order?quantity=1&side=BUY&signature=41c5bb5a29e7aefc0d0346b4706c4fefda8acb7a5df89ee4307f2911e9974622&symbol=USDCUSDT&timestamp=1765925377380&type=MARKET": context canceled



