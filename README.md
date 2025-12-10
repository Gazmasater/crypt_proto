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
Saved profile in /home/gaz358/pprof/pprof.crypt_proto.samples.cpu.006.pb.gz
File: crypt_proto
Build ID: e917f6089adf5df784007f2200e913a21bb75414
Type: cpu
Time: 2025-12-10 15:20:13 MSK
Duration: 30.16s, Total samples = 7.12s (23.61%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 3460ms, 48.60% of 7120ms total
Dropped 151 nodes (cum <= 35.60ms)
Showing top 10 nodes out of 158
      flat  flat%   sum%        cum   cum%
    1650ms 23.17% 23.17%     1650ms 23.17%  internal/runtime/syscall.Syscall6
     760ms 10.67% 33.85%      760ms 10.67%  runtime.futex
     180ms  2.53% 36.38%      180ms  2.53%  runtime.procyield
     160ms  2.25% 38.62%      160ms  2.25%  runtime.nextFreeFast (inline)
     150ms  2.11% 40.73%      190ms  2.67%  runtime.unlock2
     130ms  1.83% 42.56%      570ms  8.01%  runtime.selectgo
     120ms  1.69% 44.24%      330ms  4.63%  runtime.lock2
     110ms  1.54% 45.79%      110ms  1.54%  strconv.readFloat
     100ms  1.40% 47.19%     1000ms 14.04%  google.golang.org/protobuf/internal/impl.(*MessageInfo).unmarshalPointerEager
     100ms  1.40% 48.60%      140ms  1.97%  runtime.findObject
(pprof) 


