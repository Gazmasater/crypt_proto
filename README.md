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




File: arb
Build ID: f667473fd8b6ec748bc74083d2b9bad95785f06e
Type: cpu
Time: 2026-01-09 16:33:43 MSK
Duration: 30.14s, Total samples = 3.48s (11.55%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 1960ms, 56.32% of 3480ms total
Dropped 80 nodes (cum <= 17.40ms)
Showing top 10 nodes out of 171
      flat  flat%   sum%        cum   cum%
    1010ms 29.02% 29.02%     1010ms 29.02%  internal/runtime/syscall.Syscall6
     300ms  8.62% 37.64%      300ms  8.62%  runtime.futex
     140ms  4.02% 41.67%      430ms 12.36%  runtime.scanobject
     110ms  3.16% 44.83%      110ms  3.16%  aeshashbody
      90ms  2.59% 47.41%      110ms  3.16%  runtime.typePointers.next
      80ms  2.30% 49.71%      120ms  3.45%  runtime.findObject
      70ms  2.01% 51.72%      130ms  3.74%  github.com/tidwall/gjson.parseObject
      60ms  1.72% 53.45%       60ms  1.72%  runtime.(*mspan).base (inline)
      50ms  1.44% 54.89%      510ms 14.66%  crypt_proto/internal/queue.(*MemoryStore).apply
      50ms  1.44% 56.32%      190ms  5.46%  github.com/tidwall/gjson.Get
(pprof) 



gaz358@gaz358-BOD-WXX9:~/myprog/crypt_proto$ go tool pprof http://localhost:6060/debug/pprof/heap
Fetching profile over HTTP from http://localhost:6060/debug/pprof/heap
Saved profile in /home/gaz358/pprof/pprof.arb.alloc_objects.alloc_space.inuse_objects.inuse_space.003.pb.gz
File: arb
Build ID: f667473fd8b6ec748bc74083d2b9bad95785f06e
Type: inuse_space
Time: 2026-01-09 16:35:51 MSK
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 3565.46kB, 100% of 3565.46kB total
Showing top 10 nodes out of 27
      flat  flat%   sum%        cum   cum%
    2052kB 57.55% 57.55%     2052kB 57.55%  runtime.allocm
 1000.34kB 28.06% 85.61%  1000.34kB 28.06%  main.main
  513.12kB 14.39%   100%   513.12kB 14.39%  vendor/golang.org/x/net/http2/hpack.newInternalNode (inline)
         0     0%   100%   513.12kB 14.39%  net/http.(*http2ClientConn).readLoop
         0     0%   100%   513.12kB 14.39%  net/http.(*http2Framer).ReadFrame
         0     0%   100%   513.12kB 14.39%  net/http.(*http2Framer).readMetaFrame
         0     0%   100%   513.12kB 14.39%  net/http.(*http2clientConnReadLoop).run
         0     0%   100%  1000.34kB 28.06%  runtime.main
         0     0%   100%      513kB 14.39%  runtime.mcall
         0     0%   100%     1539kB 43.16%  runtime.mstart
(pprof) 


