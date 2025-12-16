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





[ARB] +0.141%  USDT→BTC→BDX→USDT  maxStart=64.3597 USDT (64.3597 USDT)  safeStart=64.3597 USDT (64.3597 USDT) (x1.00)  bottleneck=BDXUSDT
  BTCUSDT (BTC/USDT): bid=86213.0000000000 ask=86217.2200000000  spread=4.2200000000 (0.00489%)  bidQty=0.0023 askQty=5.2859
  BDXBTC (BDX/BTC): bid=0.0000010394 ask=0.0000010413  spread=0.0000000019 (0.18263%)  bidQty=586.0000 askQty=2890.0000
  BDXUSDT (BDX/USDT): bid=0.0900400000 ask=0.0902600000  spread=0.0002200000 (0.24404%)  bidQty=716.1600 askQty=2629.8400
  Legs execution with fees:
    leg 1: BTCUSDT  64.359750 USDT → 0.000746 BTC  fee=0.00000037 BTC
    leg 2: BDXBTC  0.000746 BTC → 716.160000 BDX  fee=0.35825913 BDX
    leg 3: BDXUSDT  716.160000 BDX → 64.450805 USDT  fee=0.03224152 USDT

  [REAL EXEC] start=2.000000 USDT triangle=USDT→BTC→BDX→USDT
    [REAL EXEC] leg 1: qty<=0 after normalize (raw=0.0000231972) (BUY BTCUSDT)







