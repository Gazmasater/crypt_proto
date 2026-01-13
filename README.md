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
Saved profile in /home/gaz358/pprof/pprof.arb.samples.cpu.110.pb.gz
File: arb
Build ID: 4200a135dcaa2d81002edd0dcb42db1ef4801138
Type: cpu
Time: 2026-01-13 23:36:51 MSK
Duration: 30.08s, Total samples = 1.49s ( 4.95%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 940ms, 63.09% of 1490ms total
Showing top 10 nodes out of 158
      flat  flat%   sum%        cum   cum%
     590ms 39.60% 39.60%      590ms 39.60%  internal/runtime/syscall.Syscall6
     100ms  6.71% 46.31%      100ms  6.71%  runtime.futex
      60ms  4.03% 50.34%       60ms  4.03%  aeshashbody
      30ms  2.01% 52.35%       30ms  2.01%  github.com/tidwall/gjson.parseString
      30ms  2.01% 54.36%      350ms 23.49%  runtime.findRunnable
      30ms  2.01% 56.38%       30ms  2.01%  runtime.memmove
      30ms  2.01% 58.39%      200ms 13.42%  runtime.netpoll
      30ms  2.01% 60.40%       30ms  2.01%  runtime.typePointers.next
      20ms  1.34% 61.74%      910ms 61.07%  crypt_proto/internal/collector.(*kucoinWS).readLoop
      20ms  1.34% 63.09%       20ms  1.34%  crypto/tls.(*Conn).handshakeContext
(pprof) 



