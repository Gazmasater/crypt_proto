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




2025/12/10 14:24:21.785430 [WS #10] Pong через 274.903338ms
2025/12/10 14:24:21.785556 [WS #3] Pong через 275.037248ms
2025/12/10 14:24:21.785987 [WS #7] Pong через 275.5602ms
2025/12/10 14:24:21.799730 [WS #1] Pong через 289.169191ms
2025/12/10 14:24:21.803034 [WS #6] Pong через 292.581285ms
2025/12/10 14:24:21.804146 [WS #4] Pong через 293.675262ms
2025/12/10 14:24:21.804488 [WS #8] Pong через 294.073703ms
2025/12/10 14:24:21.805072 [WS #11] Pong через 294.60639ms\



gaz358@gaz358-BOD-WXX9:~/myprog/crypt_proto$    go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
Fetching profile over HTTP from http://localhost:6060/debug/pprof/profile?seconds=30
Saved profile in /home/gaz358/pprof/pprof.crypt_proto.samples.cpu.008.pb.gz
File: crypt_proto
Build ID: e917f6089adf5df784007f2200e913a21bb75414
Type: cpu
Time: 2025-12-10 16:21:00 MSK
Duration: 30.18s, Total samples = 18.12s (60.04%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 9000ms, 49.67% of 18120ms total
Dropped 288 nodes (cum <= 90.60ms)
Showing top 10 nodes out of 148
      flat  flat%   sum%        cum   cum%
    4290ms 23.68% 23.68%     4290ms 23.68%  internal/runtime/syscall.Syscall6
    1590ms  8.77% 32.45%     1590ms  8.77%  runtime.futex
     640ms  3.53% 35.98%      640ms  3.53%  strconv.readFloat
     460ms  2.54% 38.52%      460ms  2.54%  runtime.procyield
     410ms  2.26% 40.78%     2820ms 15.56%  google.golang.org/protobuf/internal/impl.(*MessageInfo).unmarshalPointerEager
     370ms  2.04% 42.83%      370ms  2.04%  runtime.usleep
     360ms  1.99% 44.81%     1660ms  9.16%  runtime.selectgo
     330ms  1.82% 46.63%      330ms  1.82%  runtime.nextFreeFast (inline)
     310ms  1.71% 48.34%      860ms  4.75%  runtime.lock2
     240ms  1.32% 49.67%      660ms  3.64%  runtime.scanobject


