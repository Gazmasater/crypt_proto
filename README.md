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




[ARB] +0.117%  USDT→USDC→TON→USDT  maxStart=69.5719 USDT (69.5719 USDT)  safeStart=69.5719 USDT (69.5719 USDT) (x1.00)  bottleneck=TONUSDC
  USDCUSDT (USDC/USDT): bid=0.9999000000 ask=1.0000000000  spread=0.0001000000 (0.01000%)  bidQty=99760.5300 askQty=106715.9200
  TONUSDC (TON/USDC): bid=1.4960000000 ask=1.4980000000  spread=0.0020000000 (0.13360%)  bidQty=0.9300 askQty=46.4200
  TONUSDT (TON/USDT): bid=1.5020000000 ask=1.5040000000  spread=0.0020000000 (0.13307%)  bidQty=68.6100 askQty=6209.0800
  Legs execution with fees:
    leg 1: USDCUSDT  69.571946 USDT → 69.537160 USDC  fee=0.03478597 USDC
    leg 2: TONUSDC  69.537160 USDC → 46.396790 TON  fee=0.02321000 TON
    leg 3: TONUSDT  46.396790 TON → 69.653135 USDT  fee=0.03484399 USDT

  [REAL EXEC] start=2.000000 USDT triangle=USDT→USDC→TON→USDT
    [REAL EXEC] leg 1: BUY USDCUSDT qty=2.00000000
2025-12-16 03:09:47.189
[ARB] +0.117%  USDT→USDC→TON→USDT  maxStart=66.9791 USDT (66.9791 USDT)  safeStart=66.9791 USDT (66.9791 USDT) (x1.00)  bottleneck=TONUSDC
  USDCUSDT (USDC/USDT): bid=0.9999000000 ask=1.0000000000  spread=0.0001000000 (0.01000%)  bidQty=99760.5300 askQty=106715.9200
  TONUSDC (TON/USDC): bid=1.4970000000 ask=1.4980000000  spread=0.0010000000 (0.06678%)  bidQty=54333.6700 askQty=44.6900
  TONUSDT (TON/USDT): bid=1.5020000000 ask=1.5040000000  spread=0.0020000000 (0.13307%)  bidQty=68.6100 askQty=6209.0800
  Legs execution with fees:
    leg 1: USDCUSDT  66.979110 USDT → 66.945620 USDC  fee=0.03348955 USDC
    leg 2: TONUSDC  66.945620 USDC → 44.667655 TON  fee=0.02234500 TON
    leg 3: TONUSDT  44.667655 TON → 67.057272 USDT  fee=0.03354541 USDT

  [REAL EXEC] start=2.000000 USDT triangle=USDT→USDC→TON→USDT
    [REAL EXEC] leg 1: BUY USDCUSDT qty=2.00000000
2025-12-16 03:09:47.203
  [REAL EXEC] start=2.000000 USDT triangle=USDT→USDC→TON→USDT
    [REAL EXEC] leg 1: BUY USDCUSDT qty=2.00000000
[ARB] +0.117%  USDT→USDC→TON→USDT  maxStart=65.7651 USDT (65.7651 USDT)  safeStart=65.7651 USDT (65.7651 USDT) (x1.00)  bottleneck=TONUSDC
  USDCUSDT (USDC/USDT): bid=0.9999000000 ask=1.0000000000  spread=0.0001000000 (0.01000%)  bidQty=99760.5300 askQty=106715.9200
  TONUSDC (TON/USDC): bid=1.4970000000 ask=1.4980000000  spread=0.0010000000 (0.06678%)  bidQty=54333.6700 askQty=43.8800
  TONUSDT (TON/USDT): bid=1.5020000000 ask=1.5040000000  spread=0.0020000000 (0.13307%)  bidQty=68.6100 askQty=6209.0800
  Legs execution with fees:
    leg 1: USDCUSDT  65.765123 USDT → 65.732240 USDC  fee=0.03288256 USDC
    leg 2: TONUSDC  65.732240 USDC → 43.858060 TON  fee=0.02194000 TON
    leg 3: TONUSDT  43.858060 TON → 65.841869 USDT  fee=0.03293740 USDT

2025-12-16 03:09:47.209
[ARB] +0.117%  USDT→USDC→TON→USDT  maxStart=33.3189 USDT (33.3189 USDT)  safeStart=33.3189 USDT (33.3189 USDT) (x1.00)  bottleneck=TONUSDT
  USDCUSDT (USDC/USDT): bid=0.9999000000 ask=1.0000000000  spread=0.0001000000 (0.01000%)  bidQty=99760.5300 askQty=106715.9200
  TONUSDC (TON/USDC): bid=1.4970000000 ask=1.4980000000  spread=0.0010000000 (0.06678%)  bidQty=54333.6700 askQty=43.8800
  TONUSDT (TON/USDT): bid=1.5020000000 ask=1.5040000000  spread=0.0020000000 (0.13307%)  bidQty=22.2200 askQty=6209.0800
  Legs execution with fees:
    leg 1: USDCUSDT  33.318871 USDT → 33.302211 USDC  fee=0.01665944 USDC
    leg 2: TONUSDC  33.302211 USDC → 22.220000 TON  fee=0.01111556 TON
    leg 3: TONUSDT  22.220000 TON → 33.357753 USDT  fee=0.01668722 USDT

  [REAL EXEC] start=2.000000 USDT triangle=USDT→USDC→TON→USDT
    [REAL EXEC] leg 1: BUY USDCUSDT qty=2.00000000
    [REAL EXEC] leg 2: BUY TONUSDC qty=1.33511348
    [REAL EXEC] leg 2: BUY TONUSDC qty=1.33511348
    [REAL EXEC] leg 2: BUY TONUSDC qty=1.33511348
    [REAL EXEC] leg 2: BUY TONUSDC qty=1.33511348
    [REAL EXEC] leg 2 ERROR: mexc order error: status=400 body={"msg":" quantity scale is invalid","code":400}
    [REAL EXEC] leg 2 ERROR: mexc order error: status=400 body={"msg":" quantity scale is invalid","code":400}
    [REAL EXEC] leg 2 ERROR: mexc order error: status=400 body={"msg":" quantity scale is invalid","code":400}
    [REAL EXEC] leg 2 ERROR: mexc order error: status=400 body={"msg":" quantity scale is invalid","code":400}
^C2025/12/16 03:09:50.683814 shutting down...





