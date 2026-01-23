Если оставить только нужное:

p99 execution latency
Micro-volatility (100 мс)
Fill ratio
Capture rate
Inventory drift




Название API
9623527002

696935c42a6dcd00013273f2
b348b686-55ff-4290-897b-02d55f815f65




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
Saved profile in /home/gaz358/pprof/pprof.arb.samples.cpu.223.pb.gz
File: arb
Build ID: e1f97f19b4005e7c00459e4ff590268a861df570
Type: cpu
Time: 2026-01-23 16:09:12 MSK
Duration: 30.09s, Total samples = 1.70s ( 5.65%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top 
Showing nodes accounting for 950ms, 55.88% of 1700ms total
Showing top 10 nodes out of 179
      flat  flat%   sum%        cum   cum%
     540ms 31.76% 31.76%      540ms 31.76%  internal/runtime/syscall.Syscall6
      90ms  5.29% 37.06%       90ms  5.29%  runtime.futex
      80ms  4.71% 41.76%      130ms  7.65%  github.com/tidwall/gjson.parseObject
      50ms  2.94% 44.71%       50ms  2.94%  aeshashbody
      40ms  2.35% 47.06%       40ms  2.35%  runtime.casgstatus
      30ms  1.76% 48.82%      150ms  8.82%  github.com/tidwall/gjson.getBytes
      30ms  1.76% 50.59%       40ms  2.35%  internal/runtime/maps.(*Iter).Next
      30ms  1.76% 52.35%       30ms  1.76%  runtime.pMask.read (inline)
      30ms  1.76% 54.12%      120ms  7.06%  runtime.scanobject
      30ms  1.76% 55.88%       50ms  2.94%  runtime.stealWork
(pprof) 


