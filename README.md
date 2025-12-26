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




gaz358@gaz358-BOD-WXX9:~/myprog/crypt_proto$ go tool pprof http://localhost:6060/debug/pprof/heap
Fetching profile over HTTP from http://localhost:6060/debug/pprof/heap
Saved profile in /home/gaz358/pprof/pprof.arb.alloc_objects.alloc_space.inuse_objects.inuse_space.001.pb.gz
File: arb
Build ID: f408ee82a0b1fc807761d706ef2d6f9d43e572ce
Type: inuse_space
Time: 2025-12-27 01:55:29 MSK
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 2600.84kB, 100% of 2600.84kB total
Showing top 10 nodes out of 13
      flat  flat%   sum%        cum   cum%
    2052kB 78.90% 78.90%     2052kB 78.90%  runtime.allocm
  548.84kB 21.10%   100%   548.84kB 21.10%  main.main
         0     0%   100%   548.84kB 21.10%  runtime.main
         0     0%   100%     1539kB 59.17%  runtime.mcall
         0     0%   100%      513kB 19.72%  runtime.mstart
         0     0%   100%      513kB 19.72%  runtime.mstart0
         0     0%   100%      513kB 19.72%  runtime.mstart1
         0     0%   100%     2052kB 78.90%  runtime.newm
         0     0%   100%     1539kB 59.17%  runtime.park_m
         0     0%   100%     2052kB 78.90%  runtime.resetspinning
(pprof) 



gaz358@gaz358-BOD-WXX9:~/myprog/crypt_proto$    go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
Fetching profile over HTTP from http://localhost:6060/debug/pprof/profile?seconds=30
Saved profile in /home/gaz358/pprof/pprof.arb.samples.cpu.006.pb.gz
File: arb
Build ID: f408ee82a0b1fc807761d706ef2d6f9d43e572ce
Type: cpu
Time: 2025-12-27 01:54:17 MSK
Duration: 30s, Total samples = 500ms ( 1.67%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 320ms, 64.00% of 500ms total
Showing top 10 nodes out of 144
      flat  flat%   sum%        cum   cum%
     160ms 32.00% 32.00%      160ms 32.00%  internal/runtime/syscall.Syscall6
      50ms 10.00% 42.00%       50ms 10.00%  runtime.futex
      20ms  4.00% 46.00%      340ms 68.00%  crypt_proto/internal/collector.(*MEXCCollector).readLoop
      20ms  4.00% 50.00%       70ms 14.00%  runtime.netpoll
      20ms  4.00% 54.00%       20ms  4.00%  runtime.scanblock
      10ms  2.00% 56.00%       20ms  4.00%  crypto/internal/fips140/aes/gcm.(*GCM).Open
      10ms  2.00% 58.00%       10ms  2.00%  crypto/internal/fips140/alias.InexactOverlap (inline)
      10ms  2.00% 60.00%       40ms  8.00%  crypto/tls.(*halfConn).decrypt
      10ms  2.00% 62.00%       10ms  2.00%  crypto/tls.(*halfConn).explicitNonceLen
      10ms  2.00% 64.00%       10ms  2.00%  gogo
(pprof) 




