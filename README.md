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


go run .
2025/12/15 04:27:14.150917 pprof on http://localhost:6060/debug/pprof/
2025/12/15 04:27:14.150920 Exchange: MEXC
2025/12/15 04:27:14.151038 Triangles file: triangles_markets.csv
2025/12/15 04:27:14.151044 Book interval: 10ms
2025/12/15 04:27:14.151047 Fee per leg: 0.0400 % (rate=0.000400)
2025/12/15 04:27:14.151051 Min profit per cycle: 0.1000 % (rate=0.001000)
2025/12/15 04:27:14.151054 Min start amount: 10.0000
2025/12/15 04:27:14.151056 Start fraction: 0.5000
2025/12/15 04:27:14.151900 треугольников всего: 415
2025/12/15 04:27:14.151904 символов в индексе треугольников: 594
2025/12/15 04:27:14.151911 символов для подписки всего: 594
2025/12/15 04:27:14.151914 Using exchange: MEXC
2025/12/15 04:27:14.152151 [MEXC] будем использовать 24 WS-подключений
2025/12/15 04:27:20.147405 [MEXC WS #20] connected to wss://wbs-api.mexc.com/ws (symbols: 25)
2025/12/15 04:27:20.147917 [MEXC WS #1] connected to wss://wbs-api.mexc.com/ws (symbols: 25)
2025/12/15 04:27:20.148095 [MEXC WS #19] connected to wss://wbs-api.mexc.com/ws (symbols: 25)
2025/12/15 04:27:20.148201 [MEXC WS #15] connected to wss://wbs-api.mexc.com/ws (symbols: 25)
2025/12/15 04:27:20.148221 [MEXC WS #1] SUB -> 25 topics
2025/12/15 04:27:20.148290 [MEXC WS #19] SUB -> 25 topics
2025/12/15 04:27:20.148098 [MEXC WS #17] connected to wss://wbs-api.mexc.com/ws (symbols: 25)
2025/12/15 04:27:20.148340 [MEXC WS #15] SUB -> 25 topics
2025/12/15 04:27:20.148360 [MEXC WS #21] connected to wss://wbs-api.mexc.com/ws (symbols: 25)
2025/12/15 04:27:20.148407 [MEXC WS #17] SUB -> 25 topics
2025/12/15 04:27:20.148425 [MEXC WS #23] connected to wss://wbs-api.mexc.com/ws (symbols: 19)
2025/12/15 04:27:20.148445 [MEXC WS #18] connected to wss://wbs-api.mexc.com/ws (symbols: 25)
2025/12/15 04:27:20.148457 [MEXC WS #6] connected to wss://wbs-api.mexc.com/ws (symbols: 25)
2025/12/15 04:27:20.147387 [MEXC WS #22] connected to wss://wbs-api.mexc.com/ws (symbols: 25)
2025/12/15 04:27:20.148479 [MEXC WS #5] connected to wss://wbs-api.mexc.com/ws (symbols: 25)
2025/12/15 04:27:20.148453 [MEXC WS #21] SUB -> 25 topics
2025/12/15 04:27:20.148541 [MEXC WS #18] SUB -> 25 topics
2025/12/15 04:27:20.148560 [MEXC WS #20] SUB -> 25 topics
2025/12/15 04:27:20.148604 [MEXC WS #6] SUB -> 25 topics
2025/12/15 04:27:20.148607 [MEXC WS #23] SUB -> 19 topics
2025/12/15 04:27:20.148605 [MEXC WS #16] connected to wss://wbs-api.mexc.com/ws (symbols: 25)
2025/12/15 04:27:20.148723 [MEXC WS #22] SUB -> 25 topics
2025/12/15 04:27:20.148772 [MEXC WS #5] SUB -> 25 topics
2025/12/15 04:27:20.148772 [MEXC WS #12] connected to wss://wbs-api.mexc.com/ws (symbols: 25)
2025/12/15 04:27:20.148797 [MEXC WS #16] SUB -> 25 topics
2025/12/15 04:27:20.148923 [MEXC WS #12] SUB -> 25 topics
2025/12/15 04:27:20.149205 [MEXC WS #3] connected to wss://wbs-api.mexc.com/ws (symbols: 25)
2025/12/15 04:27:20.149269 [MEXC WS #8] connected to wss://wbs-api.mexc.com/ws (symbols: 25)
2025/12/15 04:27:20.149310 [MEXC WS #3] SUB -> 25 topics
2025/12/15 04:27:20.149322 [MEXC WS #0] connected to wss://wbs-api.mexc.com/ws (symbols: 25)
2025/12/15 04:27:20.149352 [MEXC WS #8] SUB -> 25 topics
2025/12/15 04:27:20.149439 [MEXC WS #0] SUB -> 25 topics
2025/12/15 04:27:20.352332 [MEXC WS #4] connected to wss://wbs-api.mexc.com/ws (symbols: 25)
2025/12/15 04:27:20.352332 [MEXC WS #9] connected to wss://wbs-api.mexc.com/ws (symbols: 25)
2025/12/15 04:27:20.352332 [MEXC WS #13] connected to wss://wbs-api.mexc.com/ws (symbols: 25)
2025/12/15 04:27:20.352332 [MEXC WS #10] connected to wss://wbs-api.mexc.com/ws (symbols: 25)
2025/12/15 04:27:20.352508 [MEXC WS #4] SUB -> 25 topics
2025/12/15 04:27:20.352332 [MEXC WS #14] connected to wss://wbs-api.mexc.com/ws (symbols: 25)
2025/12/15 04:27:20.352548 [MEXC WS #9] SUB -> 25 topics
2025/12/15 04:27:20.352332 [MEXC WS #11] connected to wss://wbs-api.mexc.com/ws (symbols: 25)
2025/12/15 04:27:20.352578 [MEXC WS #10] SUB -> 25 topics
2025/12/15 04:27:20.352620 [MEXC WS #13] SUB -> 25 topics
2025/12/15 04:27:20.352635 [MEXC WS #14] SUB -> 25 topics
2025/12/15 04:27:20.352672 [MEXC WS #11] SUB -> 25 topics
2025/12/15 04:27:20.355191 [MEXC WS #7] connected to wss://wbs-api.mexc.com/ws (symbols: 25)
2025/12/15 04:27:20.355304 [MEXC WS #7] SUB -> 25 topics
2025/12/15 04:27:20.369004 [MEXC WS #2] connected to wss://wbs-api.mexc.com/ws (symbols: 25)
2025/12/15 04:27:20.369215 [MEXC WS #2] SUB -> 25 topics



[ARB] +0.106%  FHE→USDC→USDT→FHE  maxStart=27.1500 FHE (2.1408 USDT)  safeStart=13.5750 FHE (1.0704 USDT) (x0.50)  bottleneck=FHEUSDC
  FHEUSDC (FHE/USDC): bid=0.0791500000 ask=0.0791800000  spread=0.0000300000 (0.03790%)  bidQty=27.1500 askQty=17.2400
  USDCUSDT (USDC/USDT): bid=0.9996000000 ask=0.9997000000  spread=0.0001000000 (0.01000%)  bidQty=45731.2000 askQty=68806.2500
  FHEUSDT (FHE/USDT): bid=0.0788500000 ask=0.0789400000  spread=0.0000900000 (0.11408%)  bidQty=16.1800 askQty=2292.1200


