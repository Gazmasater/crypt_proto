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




Fetching profile over HTTP from http://localhost:6060/debug/pprof/heap
Saved profile in /home/gaz358/pprof/pprof.arb.alloc_objects.alloc_space.inuse_objects.inuse_space.005.pb.gz
File: arb
Build ID: 4200a135dcaa2d81002edd0dcb42db1ef4801138
Type: inuse_space
Time: 2026-01-13 17:00:45 MSK
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 4629.87kB, 100% of 4629.87kB total
Showing top 10 nodes out of 39
      flat  flat%   sum%        cum   cum%
    1539kB 33.24% 33.24%     1539kB 33.24%  runtime.allocm
 1000.34kB 21.61% 54.85%  1553.38kB 33.55%  main.main
  553.04kB 11.95% 66.79%   553.04kB 11.95%  crypt_proto/internal/queue.NewMemoryStore (inline)
  513.12kB 11.08% 77.87%   513.12kB 11.08%  vendor/golang.org/x/net/http2/hpack.newInternalNode
  512.31kB 11.07% 88.94%   512.31kB 11.07%  encoding/pem.Decode
  512.05kB 11.06%   100%  2065.43kB 44.61%  runtime.main
         0     0%   100%   512.31kB 11.07%  crypto/tls.(*Conn).HandshakeContext
         0     0%   100%   512.31kB 11.07%  crypto/tls.(*Conn).clientHandshake
         0     0%   100%   512.31kB 11.07%  crypto/tls.(*Conn).handshakeContext
         0     0%   100%   512.31kB 11.07%  crypto/tls.(*Conn).verifyServerCertificate
(pprof) 


gaz358@gaz358-BOD-WXX9:~/myprog/crypt_proto$    go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
Fetching profile over HTTP from http://localhost:6060/debug/pprof/profile?seconds=30
Saved profile in /home/gaz358/pprof/pprof.arb.samples.cpu.109.pb.gz
File: arb
Build ID: 4200a135dcaa2d81002edd0dcb42db1ef4801138
Type: cpu
Time: 2026-01-13 16:59:39 MSK
Duration: 30s, Total samples = 1.04s ( 3.47%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 640ms, 61.54% of 1040ms total
Showing top 10 nodes out of 150
      flat  flat%   sum%        cum   cum%
     390ms 37.50% 37.50%      390ms 37.50%  internal/runtime/syscall.Syscall6
      70ms  6.73% 44.23%       70ms  6.73%  runtime.futex
      30ms  2.88% 47.12%       30ms  2.88%  aeshashbody
      30ms  2.88% 50.00%       40ms  3.85%  github.com/tidwall/gjson.parseObject
      20ms  1.92% 51.92%       20ms  1.92%  internal/runtime/atomic.(*Uintptr).CompareAndSwap (inline)
      20ms  1.92% 53.85%       20ms  1.92%  internal/runtime/maps.ctrlGroup.matchH2
      20ms  1.92% 55.77%       20ms  1.92%  runtime.(*mspan).base
      20ms  1.92% 57.69%      230ms 22.12%  runtime.findRunnable
      20ms  1.92% 59.62%       50ms  4.81%  runtime.mapassign_faststr
      20ms  1.92% 61.54%       20ms  1.92%  runtime.memmove
(pprof) 



