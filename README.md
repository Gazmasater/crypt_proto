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




[ARB] +0.108%  USDT→USDC→GIGGLE→USDT  maxStart=14.4766 USDT (14.4766 USDT)  safeStart=14.4766 USDT (14.4766 USDT) (x1.00)  bottleneck=GIGGLEUSDC
  USDCUSDT (USDC/USDT): bid=0.9999000000 ask=1.0000000000  spread=0.0001000000 (0.01000%)  bidQty=108139.9300 askQty=247174.0100
  GIGGLEUSDC (GIGGLE/USDC): bid=65.5300000000 ask=65.7700000000  spread=0.2400000000 (0.36558%)  bidQty=5.4500 askQty=0.2200
  GIGGLEUSDT (GIGGLE/USDT): bid=65.9400000000 ask=66.0200000000  spread=0.0800000000 (0.12125%)  bidQty=1.5900 askQty=3.8000
  Legs execution with fees:
    leg 1: USDCUSDT  14.476638 USDT → 14.469400 USDC  fee=0.00723832 USDC
    leg 2: GIGGLEUSDC  14.469400 USDC → 0.219890 GIGGLE  fee=0.00011000 GIGGLE
    leg 3: GIGGLEUSDT  0.219890 GIGGLE → 14.492297 USDT  fee=0.00724977 USDT

  [REAL EXEC] start=2.000000 USDT triangle=USDT→USDC→GIGGLE→USDT
    [REAL EXEC] leg 1: BUY USDCUSDT qty=2.00000000
2025-12-16 02:26:07.922
[ARB] +0.108%  USDT→USDC→GIGGLE→USDT  maxStart=14.4766 USDT (14.4766 USDT)  safeStart=14.4766 USDT (14.4766 USDT) (x1.00)  bottleneck=GIGGLEUSDC
  USDCUSDT (USDC/USDT): bid=0.9999000000 ask=1.0000000000  spread=0.0001000000 (0.01000%)  bidQty=108139.9300 askQty=247174.0100
  GIGGLEUSDC (GIGGLE/USDC): bid=65.5300000000 ask=65.7700000000  spread=0.2400000000 (0.36558%)  bidQty=5.4500 askQty=0.2200
  GIGGLEUSDT (GIGGLE/USDT): bid=65.9400000000 ask=66.0300000000  spread=0.0900000000 (0.13639%)  bidQty=1.5900 askQty=4.6100
  Legs execution with fees:
    leg 1: USDCUSDT  14.476638 USDT → 14.469400 USDC  fee=0.00723832 USDC
    leg 2: GIGGLEUSDC  14.469400 USDC → 0.219890 GIGGLE  fee=0.00011000 GIGGLE
    leg 3: GIGGLEUSDT  0.219890 GIGGLE → 14.492297 USDT  fee=0.00724977 USDT

  [REAL EXEC] start=2.000000 USDT triangle=USDT→USDC→GIGGLE→USDT
    [REAL EXEC] leg 1: BUY USDCUSDT qty=2.00000000
2025-12-16 02:26:07.999
[ARB] +0.108%  USDT→USDC→GIGGLE→USDT  maxStart=14.4766 USDT (14.4766 USDT)  safeStart=14.4766 USDT (14.4766 USDT) (x1.00)  bottleneck=GIGGLEUSDC
  USDCUSDT (USDC/USDT): bid=0.9999000000 ask=1.0000000000  spread=0.0001000000 (0.01000%)  bidQty=108139.9300 askQty=247174.0100
  GIGGLEUSDC (GIGGLE/USDC): bid=65.5300000000 ask=65.7700000000  spread=0.2400000000 (0.36558%)  bidQty=5.4500 askQty=0.2200
  GIGGLEUSDT (GIGGLE/USDT): bid=65.9400000000 ask=66.0300000000  spread=0.0900000000 (0.13639%)  bidQty=1.5900 askQty=4.7600
  Legs execution with fees:
    leg 1: USDCUSDT  14.476638 USDT → 14.469400 USDC  fee=0.00723832 USDC
    leg 2: GIGGLEUSDC  14.469400 USDC → 0.219890 GIGGLE  fee=0.00011000 GIGGLE
    leg 3: GIGGLEUSDT  0.219890 GIGGLE → 14.492297 USDT  fee=0.00724977 USDT

  [REAL EXEC] start=2.000000 USDT triangle=USDT→USDC→GIGGLE→USDT
    [REAL EXEC] leg 1: BUY USDCUSDT qty=2.00000000
2025-12-16 02:26:08.253
[ARB] +0.108%  USDT→USDC→GIGGLE→USDT  maxStart=14.4766 USDT (14.4766 USDT)  safeStart=14.4766 USDT (14.4766 USDT) (x1.00)  bottleneck=GIGGLEUSDC
  USDCUSDT (USDC/USDT): bid=0.9999000000 ask=1.0000000000  spread=0.0001000000 (0.01000%)  bidQty=108139.9300 askQty=247174.0100
  GIGGLEUSDC (GIGGLE/USDC): bid=65.5300000000 ask=65.7700000000  spread=0.2400000000 (0.36558%)  bidQty=5.4500 askQty=0.2200
  GIGGLEUSDT (GIGGLE/USDT): bid=65.9400000000 ask=66.0200000000  spread=0.0800000000 (0.12125%)  bidQty=1.5900 askQty=0.7600
  Legs execution with fees:
    leg 1: USDCUSDT  14.476638 USDT → 14.469400 USDC  fee=0.00723832 USDC
    leg 2: GIGGLEUSDC  14.469400 USDC → 0.219890 GIGGLE  fee=0.00011000 GIGGLE
    leg 3: GIGGLEUSDT  0.219890 GIGGLE → 14.492297 USDT  fee=0.00724977 USDT

  [REAL EXEC] start=2.000000 USDT triangle=USDT→USDC→GIGGLE→USDT
    [REAL EXEC] leg 1: BUY USDCUSDT qty=2.00000000
    [REAL EXEC] leg 1 ERROR: mexc order error: status=400 body={"msg":" quantity scale is invalid","code":400}
    [REAL EXEC] leg 1 ERROR: mexc order error: status=400 body={"msg":" quantity scale is invalid","code":400}
    [REAL EXEC] leg 1 ERROR: mexc order error: status=400 body={"msg":" quantity scale is invalid","code":400}
2025-12-16 02:26:08.350
[ARB] +0.123%  USDT→USDC→GIGGLE→USDT  maxStart=14.4766 USDT (14.4766 USDT)  safeStart=14.4766 USDT (14.4766 USDT) (x1.00)  bottleneck=GIGGLEUSDC
  USDCUSDT (USDC/USDT): bid=0.9999000000 ask=1.0000000000  spread=0.0001000000 (0.01000%)  bidQty=108139.9300 askQty=247174.0100
  GIGGLEUSDC (GIGGLE/USDC): bid=65.5300000000 ask=65.7700000000  spread=0.2400000000 (0.36558%)  bidQty=5.4500 askQty=0.2200
  GIGGLEUSDT (GIGGLE/USDT): bid=65.9500000000 ask=66.0200000000  spread=0.0700000000 (0.10608%)  bidQty=3.0100 askQty=0.7600
  Legs execution with fees:
    leg 1: USDCUSDT  14.476638 USDT → 14.469400 USDC  fee=0.00723832 USDC
    leg 2: GIGGLEUSDC  14.469400 USDC → 0.219890 GIGGLE  fee=0.00011000 GIGGLE
    leg 3: GIGGLEUSDT  0.219890 GIGGLE → 14.494495 USDT  fee=0.00725087 USDT

  [REAL EXEC] start=2.000000 USDT triangle=USDT→USDC→GIGGLE→USDT
    [REAL EXEC] leg 1: BUY USDCUSDT qty=2.00000000
    [REAL EXEC] leg 1 ERROR: mexc order error: status=400 body={"msg":" quantity scale is invalid","code":400}
    [REAL EXEC] leg 1 ERROR: mexc order error: status=400 body={"msg":" quantity scale is invalid","code":400}
^C2025/12/16 02:26:18.515957 shutting down...
panic: send on closed channel

goroutine 34 [running]:







