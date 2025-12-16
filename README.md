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





  [REAL EXEC] start=2.000000 USDT triangle=USDT→FHE→USDC→USDT
[ARB] +0.135%  USDT→FHE→USDC→USDT  maxStart=37.6690 USDT (37.6690 USDT)  safeStart=37.6690 USDT (37.6690 USDT) (x1.00)  bottleneck=FHEUSDC
  FHEUSDT (FHE/USDT): bid=0.1250100000 ask=0.1252500000  spread=0.0002400000 (0.19180%)  bidQty=385.5800 askQty=2572.3300
    [REAL EXEC] leg 1: BUY FHEUSDT qty=15.96000000 (raw=15.9680638723)
  FHEUSDC (FHE/USDC): bid=0.1256200000 ask=0.1266700000  spread=0.0010500000 (0.83238%)  bidQty=300.6000 askQty=549.1400
  USDCUSDT (USDC/USDT): bid=0.9999000000 ask=1.0000000000  spread=0.0001000000 (0.01000%)  bidQty=107899.4100 askQty=273768.8700
  Legs execution with fees:
    leg 1: FHEUSDT  37.668984 USDT → 300.600000 FHE  fee=0.15037519 FHE
    leg 2: FHEUSDC  300.600000 FHE → 37.742491 USDC  fee=0.01888069 USDC
    leg 3: USDCUSDT  37.742491 USDC → 37.719848 USDT  fee=0.01886936 USDT

2025-12-16 04:08:07.525
[ARB] +0.127%  USDT→FHE→USDC→USDT  maxStart=37.6690 USDT (37.6690 USDT)  safeStart=37.6690 USDT (37.6690 USDT) (x1.00)  bottleneck=FHEUSDC
  FHEUSDT (FHE/USDT): bid=0.1250100000 ask=0.1252500000  spread=0.0002400000 (0.19180%)  bidQty=385.5800 askQty=4386.8500
  FHEUSDC (FHE/USDC): bid=0.1256100000 ask=0.1263900000  spread=0.0007800000 (0.61905%)  bidQty=300.6000 askQty=90.3500
  USDCUSDT (USDC/USDT): bid=0.9999000000 ask=1.0000000000  spread=0.0001000000 (0.01000%)  bidQty=107899.4100 askQty=273768.8700
  Legs execution with fees:
    leg 1: FHEUSDT  37.668984 USDT → 300.600000 FHE  fee=0.15037519 FHE
    leg 2: FHEUSDC  300.600000 FHE → 37.739487 USDC  fee=0.01887918 USDC
    leg 3: USDCUSDT  37.739487 USDC → 37.716845 USDT  fee=0.01886786 USDT

  [REAL EXEC] start=2.000000 USDT triangle=USDT→FHE→USDC→USDT
    [REAL EXEC] leg 1: BUY FHEUSDT qty=15.96000000 (raw=15.9680638723)
    [REAL EXEC] leg 2: SELL FHEUSDC qty=15.96000000 (raw=15.9680638723)
    [REAL EXEC] leg 2: SELL FHEUSDC qty=15.96000000 (raw=15.9680638723)
    [REAL EXEC] leg 2 ERROR: mexc order error: status=400 body={"code":10007,"msg":"symbol not support api"}
    [REAL EXEC] leg 2 ERROR: mexc order error: status=400 body={"code":10007,"msg":"symbol not support api"}
^C2025/12/16 04:08:17.516971 shutting down...







