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




(pprof) gaz358@gaz358-BOD-WXX9:~/myprog/crypt_proto$    go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
Fetching profile over HTTP from http://localhost:6060/debug/pprof/profile?seconds=30
Saved profile in /home/gaz358/pprof/pprof.arb.samples.cpu.081.pb.gz
File: arb
Build ID: 00f359f630cea5d5eb1389920b6bee5aa91f0b5e
Type: cpu
Time: 2026-01-12 10:59:06 MSK
Duration: 30.04s, Total samples = 1.98s ( 6.59%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 1130ms, 57.07% of 1980ms total
Showing top 10 nodes out of 209
      flat  flat%   sum%        cum   cum%
     670ms 33.84% 33.84%      670ms 33.84%  internal/runtime/syscall.Syscall6
     100ms  5.05% 38.89%      100ms  5.05%  runtime.futex
      70ms  3.54% 42.42%       70ms  3.54%  aeshashbody
      60ms  3.03% 45.45%      130ms  6.57%  runtime.scanobject
      50ms  2.53% 47.98%       90ms  4.55%  github.com/tidwall/gjson.parseObject
      50ms  2.53% 50.51%      130ms  6.57%  runtime.mapassign_faststr
      40ms  2.02% 52.53%       50ms  2.53%  runtime.typePointers.next
      30ms  1.52% 54.04%      820ms 41.41%  bufio.(*Reader).fill
      30ms  1.52% 55.56%       30ms  1.52%  memeqbody
      30ms  1.52% 57.07%       60ms  3.03%  runtime.mapaccess1_faststr
(pprof) 



