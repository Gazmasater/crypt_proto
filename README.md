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




[ARB] +0.189%  USDT→ULTIMA→EUR→USDT  maxStart=2.2651 USDT (2.2651 USDT)  safeStart=2.2651 USDT (2.2651 USDT) (x1.00)  bottleneck=ULTIMAEUR
  ULTIMAUSDT (ULTIMA/USDT): bid=5645.1400000000 ask=5659.9100000000  spread=14.7700000000 (0.26130%)  bidQty=0.0055 askQty=0.0005
  ULTIMAEUR (ULTIMA/EUR): bid=4834.1300000000 ask=4858.7700000000  spread=24.6400000000 (0.50841%)  bidQty=0.0004 askQty=0.0048
  EURUSDT (EUR/USDT): bid=1.1748000000 ask=1.1749000000  spread=0.0001000000 (0.00851%)  bidQty=48228.6200 askQty=51617.4300
  Legs execution with fees:
    leg 1: ULTIMAUSDT  2.265097 USDT → 0.000400 ULTIMA  fee=0.00000020 ULTIMA
    leg 2: ULTIMAEUR  0.000400 ULTIMA → 1.932685 EUR  fee=0.00096683 EUR
    leg 3: EURUSDT  1.932685 EUR → 2.269383 USDT  fee=0.00113526 USDT

  [REAL EXEC] start=2.000000 USDT triangle=USDT→ULTIMA→EUR→USDT
    [REAL EXEC] leg 1: BUY ULTIMAUSDT qty=0.00035336
2025-12-16 01:44:07.389
[ARB] +0.189%  USDT→ULTIMA→EUR→USDT  maxStart=2.2651 USDT (2.2651 USDT)  safeStart=2.2651 USDT (2.2651 USDT) (x1.00)  bottleneck=ULTIMAEUR
  ULTIMAUSDT (ULTIMA/USDT): bid=5652.1400000000 ask=5659.9100000000  spread=7.7700000000 (0.13738%)  bidQty=0.0052 askQty=0.0005
  ULTIMAEUR (ULTIMA/EUR): bid=4834.1300000000 ask=4858.7700000000  spread=24.6400000000 (0.50841%)  bidQty=0.0004 askQty=0.0048
  EURUSDT (EUR/USDT): bid=1.1748000000 ask=1.1749000000  spread=0.0001000000 (0.00851%)  bidQty=48228.6200 askQty=51617.4300
  Legs execution with fees:
    leg 1: ULTIMAUSDT  2.265097 USDT → 0.000400 ULTIMA  fee=0.00000020 ULTIMA
    leg 2: ULTIMAEUR  0.000400 ULTIMA → 1.932685 EUR  fee=0.00096683 EUR
    leg 3: EURUSDT  1.932685 EUR → 2.269383 USDT  fee=0.00113526 USDT

  [REAL EXEC] start=2.000000 USDT triangle=USDT→ULTIMA→EUR→USDT
    [REAL EXEC] leg 1: BUY ULTIMAUSDT qty=0.00035336
    [REAL EXEC] leg 1 ERROR: mexc order error: status=400 body={"code":700004,"msg":"Mandatory parameter 'signature' was not sent, was empty/null, or malformed."}
    [REAL EXEC] leg 1 ERROR: mexc order error: status=400 body={"code":700004,"msg":"Mandatory parameter 'signature' was not sent, was empty/null, or malformed."}
^C2025/12/16 01:44:13.321650 shutting down...





