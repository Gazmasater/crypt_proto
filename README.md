apikey = "4333ed4b-cd83-49f5-97d1-c399e2349748"
secretkey = "E3848531135EDB4CCFDA0F1BC14CD274"
IP = ""
Название API-ключа = "Arb"
Доступы = "Чтение"



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




gaz358@gaz358-BOD-WXX9:~/myprog/crypt_proto$    go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
Fetching profile over HTTP from http://localhost:6060/debug/pprof/profile?seconds=30
Saved profile in /home/gaz358/pprof/pprof.arb.samples.cpu.004.pb.gz
File: arb
Build ID: 8d556bf00aa7f5ef0291dc0f823e64c039de5d63
Type: cpu
Time: 2025-12-27 01:06:44 MSK
Duration: 30.06s, Total samples = 3.34s (11.11%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 1670ms, 50.00% of 3340ms total
Dropped 120 nodes (cum <= 16.70ms)
Showing top 10 nodes out of 190
      flat  flat%   sum%        cum   cum%
    1060ms 31.74% 31.74%     1060ms 31.74%  internal/runtime/syscall.Syscall6
     200ms  5.99% 37.72%      200ms  5.99%  runtime.futex
      90ms  2.69% 40.42%       90ms  2.69%  strconv.readFloat
      50ms  1.50% 41.92%      380ms 11.38%  google.golang.org/protobuf/internal/impl.(*MessageInfo).unmarshalPointerEager
      50ms  1.50% 43.41%       50ms  1.50%  runtime.memclrNoHeapPointers
      50ms  1.50% 44.91%       50ms  1.50%  runtime.scanblock
      50ms  1.50% 46.41%       80ms  2.40%  runtime.stealWork
      40ms  1.20% 47.60%      400ms 11.98%  crypt_proto/internal/collector.(*MEXCCollector).handleWrapper
      40ms  1.20% 48.80%       40ms  1.20%  runtime.findObject
      40ms  1.20% 50.00%       40ms  1.20%  runtime.getMCache (inline)
(pprof) 





