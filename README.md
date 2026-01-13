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
Saved profile in /home/gaz358/pprof/pprof.arb.samples.cpu.102.pb.gz
File: arb
Build ID: 88f36eeb9337c7213d316c4fb407e31ceb0f85a7
Type: cpu
Time: 2026-01-12 23:15:21 MSK
Duration: 30s, Total samples = 1.38s ( 4.60%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 820ms, 59.42% of 1380ms total
Showing top 10 nodes out of 167
      flat  flat%   sum%        cum   cum%
     550ms 39.86% 39.86%      550ms 39.86%  internal/runtime/syscall.Syscall6
      70ms  5.07% 44.93%       70ms  5.07%  runtime.futex
      40ms  2.90% 47.83%       40ms  2.90%  aeshashbody
      30ms  2.17% 50.00%       30ms  2.17%  bytes.(*Buffer).Len
      30ms  2.17% 52.17%       30ms  2.17%  runtime.nextFreeFast
      20ms  1.45% 53.62%       50ms  3.62%  crypto/tls.(*xorNonceAEAD).Open
      20ms  1.45% 55.07%       20ms  1.45%  internal/runtime/maps.(*ctrlGroup).setEmpty
      20ms  1.45% 56.52%      320ms 23.19%  runtime.mcall
      20ms  1.45% 57.97%      170ms 12.32%  runtime.netpoll
      20ms  1.45% 59.42%       20ms  1.45%  runtime.save
(pprof) gaz358@gaz358-BOD-WXX9:~/myprog/crypt_proto$ ^C
gaz358@gaz358-BOD-WXX9:~/myprog/crypt_proto$ ^C
gaz358@gaz358-BOD-WXX9:~/myprog/crypt_proto$    go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
Fetching profile over HTTP from http://localhost:6060/debug/pprof/profile?seconds=30
Saved profile in /home/gaz358/pprof/pprof.arb.samples.cpu.103.pb.gz
File: arb
Build ID: 88f36eeb9337c7213d316c4fb407e31ceb0f85a7
Type: cpu
Time: 2026-01-13 12:04:38 MSK
Duration: 30s, Total samples = 1.13s ( 3.77%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 740ms, 65.49% of 1130ms total
Showing top 10 nodes out of 157
      flat  flat%   sum%        cum   cum%
     420ms 37.17% 37.17%      420ms 37.17%  internal/runtime/syscall.Syscall6
     130ms 11.50% 48.67%      130ms 11.50%  runtime.futex
      30ms  2.65% 51.33%       30ms  2.65%  runtime.memmove
      30ms  2.65% 53.98%       40ms  3.54%  runtime.scanobject
      30ms  2.65% 56.64%       40ms  3.54%  strconv.atof64
      20ms  1.77% 58.41%       20ms  1.77%  aeshashbody
      20ms  1.77% 60.18%       20ms  1.77%  internal/runtime/maps.ctrlGroup.matchH2
      20ms  1.77% 61.95%       20ms  1.77%  memeqbody
      20ms  1.77% 63.72%       20ms  1.77%  runtime.(*mspan).refillAllocCache
      20ms  1.77% 65.49%       20ms  1.77%  runtime.ifaceeq
(pprof) 







