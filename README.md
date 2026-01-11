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
Saved profile in /home/gaz358/pprof/pprof.arb.samples.cpu.077.pb.gz
File: arb
Build ID: 00f359f630cea5d5eb1389920b6bee5aa91f0b5e
Type: cpu
Time: 2026-01-11 23:57:35 MSK
Duration: 30.04s, Total samples = 2.03s ( 6.76%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 1190ms, 58.62% of 2030ms total
Dropped 96 nodes (cum <= 10.15ms)
Showing top 10 nodes out of 121
      flat  flat%   sum%        cum   cum%
     690ms 33.99% 33.99%      690ms 33.99%  internal/runtime/syscall.Syscall6
     150ms  7.39% 41.38%      150ms  7.39%  runtime.futex
      80ms  3.94% 45.32%      100ms  4.93%  runtime.typePointers.next
      70ms  3.45% 48.77%       70ms  3.45%  aeshashbody
      40ms  1.97% 50.74%       40ms  1.97%  github.com/tidwall/gjson.parseSquash
      40ms  1.97% 52.71%      520ms 25.62%  runtime.findRunnable
      30ms  1.48% 54.19%      110ms  5.42%  github.com/tidwall/gjson.Get
      30ms  1.48% 55.67%       70ms  3.45%  runtime.greyobject
      30ms  1.48% 57.14%       30ms  1.48%  runtime.memmove
      30ms  1.48% 58.62%       30ms  1.48%  runtime.nanotime (inline)
(pprof) 




