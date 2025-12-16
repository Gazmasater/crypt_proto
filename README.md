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





[ARB] +0.025%  USDT→USDC→VANRY→USDT  maxStart=11.3402 USDT (11.3402 USDT)  safeStart=11.3402 USDT (11.3402 USDT) (x1.00)  bottleneck=VANRYUSDC
  USDCUSDT (USDC/USDT): bid=0.9999000000 ask=1.0000000000  spread=0.0001000000 (0.01000%)  bidQty=2044982.7500 askQty=114582.8800
  VANRYUSDC (VANRY/USDC): bid=0.0082680000 ask=0.0082690000  spread=0.0000010000 (0.01209%)  bidQty=1938.6200 askQty=1370.8600
  VANRYUSDT (VANRY/USDT): bid=0.0082810000 ask=0.0082980000  spread=0.0000170000 (0.20508%)  bidQty=3949.6400 askQty=2421.8300

  [REAL EXEC] start=3.000000 USDT triangle=USDT→USDC→VANRY→USDT
    [REAL EXEC] legs: sym1=USDCUSDT sym2=VANRYUSDC sym3=VANRYUSDT
    [REAL EXEC] parsed: sym1=USDCUSDT (USDC/USDT) sym2=VANRYUSDC (VANRY/USDC) sym3=VANRYUSDT (VANRY/USDT)
    [REAL EXEC] >>> GetBalance USDT (before)
    [REAL EXEC] <<< GetBalance USDT (before) (392ms)
    [REAL EXEC] BAL before: USDT=34.782730681146
    [REAL EXEC] >>> GetBalance USDC (before leg1)
    [REAL EXEC] <<< GetBalance USDC (before leg1) (243ms)
    [REAL EXEC] leg1 PRE: BUY USDCUSDT by USDT=3.000000 ask=1.0000000000 bid=0.9999000000 | USDC before=1.001070300000
    [REAL EXEC] >>> SmartMarketBuyUSDT leg1
    [REAL EXEC] <<< SmartMarketBuyUSDT leg1 (363ms)
    [REAL EXEC] leg1 PLACE OK: orderId=C02__629750727731863552022
    [REAL EXEC] >>> waitBalanceChange USDC (after leg1)
    [REAL EXEC] <<< waitBalanceChange USDC (after leg1) (247ms)
    [REAL EXEC] leg1 BAL after: USDC=4.001070300000 delta=3.000000000000
    [REAL EXEC] >>> GetBalance VANRY (before leg2)
    [REAL EXEC] <<< GetBalance VANRY (before leg2) (257ms)
^C2025/12/17 02:08:30.560216 shutting down...



