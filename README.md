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




[ARB] +0.130%  USDT→RIO→EUR→USDT  maxStart=57.0318 USDT (57.0318 USDT)  safeStart=17.1095 USDT (17.1095 USDT) (x0.30)  bottleneck=RIOUSDT
  RIOUSDT (RIO/USDT): bid=0.1474000000 ask=0.1480000000  spread=0.0006000000 (0.40623%)  bidQty=6769.6700 askQty=385.3500
  RIOEUR (RIO/EUR): bid=0.1263000000 ask=0.1264000000  spread=0.0001000000 (0.07915%)  bidQty=803.0500 askQty=101.2700
  EURUSDT (EUR/USDT): bid=1.1751000000 ask=1.1752000000  spread=0.0001000000 (0.00851%)  bidQty=38364.0300 askQty=44448.0300
  Legs execution with fees:
    leg 1: RIOUSDT  17.109540 USDT → 115.547198 RIO  fee=0.05780250 RIO
    leg 2: RIOEUR  115.547198 RIO → 14.586314 EUR  fee=0.00729681 EUR
  [REAL EXEC] start=2.000000 USDT triangle=USDT→RIO→EUR→USDT
    [REAL EXEC] leg 1: BUY RIOUSDT qty=13.51351351
    leg 3: EURUSDT  14.586314 EUR → 17.131808 USDT  fee=0.00857019 USDT

2025-12-16 00:59:16.215
[ARB] +0.130%  USDT→RIO→EUR→USDT  maxStart=57.0318 USDT (57.0318 USDT)  safeStart=17.1095 USDT (17.1095 USDT) (x0.30)  bottleneck=RIOUSDT
  RIOUSDT (RIO/USDT): bid=0.1474000000 ask=0.1480000000  spread=0.0006000000 (0.40623%)  bidQty=6769.6700 askQty=385.3500
  RIOEUR (RIO/EUR): bid=0.1263000000 ask=0.1264000000  spread=0.0001000000 (0.07915%)  bidQty=517.4400 askQty=101.2700
  EURUSDT (EUR/USDT): bid=1.1751000000 ask=1.1752000000  spread=0.0001000000 (0.00851%)  bidQty=38364.0300 askQty=44448.0300
  Legs execution with fees:
    leg 1: RIOUSDT  17.109540 USDT → 115.547198 RIO  fee=0.05780250 RIO
    leg 2: RIOEUR  115.547198 RIO → 14.586314 EUR  fee=0.00729681 EUR
    leg 3: EURUSDT  14.586314 EUR → 17.131808 USDT  fee=0.00857019 USDT

  [REAL EXEC] start=2.000000 USDT triangle=USDT→RIO→EUR→USDT
    [REAL EXEC] leg 1: BUY RIOUSDT qty=13.51351351
    [REAL EXEC] leg 1 ERROR: mexc order error: status=400 body={"code":700013,"msg":"Invalid content Type."}
    [REAL EXEC] leg 1 ERROR: mexc order error: status=400 body={"code":700013,"msg":"Invalid content Type."}
^C2025/12/16 00:59:25.286374 shutting down...
2025/12/16 00:59:25.486689 bye






