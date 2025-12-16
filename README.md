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




exchangeinfo.go


[ARB] +0.264%  USDT→KEKIUS→USDC→USDT  maxStart=23.9905 USDT (23.9905 USDT)  safeStart=23.9905 USDT (23.9905 USDT) (x1.00)  bottleneck=KEKIUSUSDC
  KEKIUSUSDT (KEKIUS/USDT): bid=0.0087080000 ask=0.0087090000  spread=0.0000010000 (0.01148%)  bidQty=3484.8900 askQty=1045811.6600
  [REAL EXEC] start=2.000000 USDT triangle=USDT→KEKIUS→USDC→USDT
    [REAL EXEC] leg 1: BUY KEKIUSUSDT qty=229.64000000 (raw=229.6474911012)
  KEKIUSUSDC (KEKIUS/USDC): bid=0.0087460000 ask=0.0087930000  spread=0.0000470000 (0.53595%)  bidQty=2753.3000 askQty=21943.0300
  USDCUSDT (USDC/USDT): bid=0.9999000000 ask=1.0000000000  spread=0.0001000000 (0.01000%)  bidQty=173992.3400 askQty=147018.2800
  Legs execution with fees:
    leg 1: KEKIUSUSDT  23.990485 USDT → 2753.300000 KEKIUS  fee=1.37733867 KEKIUS
    leg 2: KEKIUSUSDC  2753.300000 KEKIUS → 24.068322 USDC  fee=0.01204018 USDC
    leg 3: USDCUSDT  24.068322 USDC → 24.053882 USDT  fee=0.01203296 USDT

    [REAL EXEC] leg 2: SELL KEKIUSUSDC qty=229.64000000 (raw=229.6474911012)
    [REAL EXEC] leg 3: SELL USDCUSDT qty=2.00000000 (raw=2.0084969572)
    [REAL EXEC] leg 3 ERROR: mexc order error: status=400 body={"msg":"Oversold","code":30005}
2025-12-16 05:14:40.706
[ARB] +0.244%  USDT→KEKIUS→USDC→USDT  maxStart=40.9358 USDT (40.9358 USDT)  safeStart=40.9358 USDT (40.9358 USDT) (x1.00)  bottleneck=KEKIUSUSDC
  KEKIUSUSDT (KEKIUS/USDT): bid=0.0086450000 ask=0.0086510000  spread=0.0000060000 (0.06938%)  bidQty=2866.6900 askQty=1002468.9500
  KEKIUSUSDC (KEKIUS/USDC): bid=0.0086860000 ask=0.0087230000  spread=0.0000370000 (0.42507%)  bidQty=4729.5500 askQty=195.9700
  USDCUSDT (USDC/USDT): bid=0.9999000000 ask=1.0000000000  spread=0.0001000000 (0.01000%)  bidQty=173968.3300 askQty=148539.2900
  Legs execution with fees:
    leg 1: KEKIUSUSDT  40.935805 USDT → 4729.550000 KEKIUS  fee=2.36595798 KEKIUS
    leg 2: KEKIUSUSDC  4729.550000 KEKIUS → 41.060331 USDC  fee=0.02054044 USDC
    leg 3: USDCUSDT  41.060331 USDC → 41.035697 USDT  fee=0.02052811 USDT

  [REAL EXEC] start=2.000000 USDT triangle=USDT→KEKIUS→USDC→USDT
    [REAL EXEC] leg 1: BUY KEKIUSUSDT qty=231.18000000 (raw=231.1871459947)
    [REAL EXEC] leg 2: SELL KEKIUSUSDC qty=231.18000000 (raw=231.1871459947)
    [REAL EXEC] leg 3: SELL USDCUSDT qty=2.00000000 (raw=2.0080915501)
  [REAL EXEC] done triangle USDT→KEKIUS→USDC→USDT
^C2025/12/16 05:15:05.056785 shutting down...










