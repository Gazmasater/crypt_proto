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
Saved profile in /home/gaz358/pprof/pprof.arb.samples.cpu.142.pb.gz
File: arb
Build ID: 4e81efd69945c21c37fff699535a53c6e0bec3d6
Type: cpu
Time: 2026-01-22 01:19:31 MSK
Duration: 30s, Total samples = 280ms ( 0.93%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 270ms, 96.43% of 280ms total
Showing top 10 nodes out of 51
      flat  flat%   sum%        cum   cum%
     150ms 53.57% 53.57%      150ms 53.57%  internal/runtime/syscall.Syscall6
      20ms  7.14% 60.71%       30ms 10.71%  runtime.pidleget
      20ms  7.14% 67.86%       30ms 10.71%  runtime.reentersyscall
      20ms  7.14% 75.00%       20ms  7.14%  runtime.write1
      10ms  3.57% 78.57%       10ms  3.57%  github.com/tidwall/gjson.parseString
      10ms  3.57% 82.14%      170ms 60.71%  net.(*conn).Read
      10ms  3.57% 85.71%       10ms  3.57%  runtime.futex
      10ms  3.57% 89.29%       10ms  3.57%  runtime.getitab
      10ms  3.57% 92.86%       10ms  3.57%  runtime.memmove
      10ms  3.57% 96.43%       10ms  3.57%  runtime.nanotime1
(pprof) 







