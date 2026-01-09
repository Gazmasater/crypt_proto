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




Saved profile in /home/gaz358/pprof/pprof.arb.samples.cpu.063.pb.gz
File: arb
Build ID: 661cf4e16f8eb544a152de46fc45312b5200894b
Type: cpu
Time: 2026-01-09 18:54:16 MSK
Duration: 30.03s, Total samples = 3.78s (12.59%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 2060ms, 54.50% of 3780ms total
Dropped 91 nodes (cum <= 18.90ms)
Showing top 10 nodes out of 173
      flat  flat%   sum%        cum   cum%
    1160ms 30.69% 30.69%     1160ms 30.69%  internal/runtime/syscall.Syscall6
     250ms  6.61% 37.30%      250ms  6.61%  runtime.futex
     130ms  3.44% 40.74%      130ms  3.44%  aeshashbody
     100ms  2.65% 43.39%      190ms  5.03%  github.com/tidwall/gjson.parseObject
     100ms  2.65% 46.03%      430ms 11.38%  runtime.scanobject
      90ms  2.38% 48.41%      130ms  3.44%  runtime.typePointers.next
      80ms  2.12% 50.53%      110ms  2.91%  runtime.findObject
      50ms  1.32% 51.85%       50ms  1.32%  github.com/tidwall/gjson.parseString
      50ms  1.32% 53.17%       50ms  1.32%  runtime.(*mspan).base (inline)
      50ms  1.32% 54.50%       50ms  1.32%  runtime.nextFreeFast
(pprof) 





