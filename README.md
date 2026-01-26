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




go run -race main.go


GOMAXPROCS=8 go run -race main.go



gaz358@gaz358-BOD-WXX9:~/myprog/crypt_proto$    go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
Fetching profile over HTTP from http://localhost:6060/debug/pprof/profile?seconds=30
Saved profile in /home/gaz358/pprof/pprof.arb.samples.cpu.253.pb.gz
File: arb
Build ID: 92d6ab1c7fab5e78cad537019b70c12c249448f3
Type: cpu
Time: 2026-01-26 16:31:09 MSK
Duration: 30s, Total samples = 690ms ( 2.30%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 520ms, 75.36% of 690ms total
Showing top 10 nodes out of 86
      flat  flat%   sum%        cum   cum%
     350ms 50.72% 50.72%      350ms 50.72%  internal/runtime/syscall.Syscall6
      60ms  8.70% 59.42%       60ms  8.70%  runtime.futex
      20ms  2.90% 62.32%      410ms 59.42%  bufio.(*Reader).fill
      20ms  2.90% 65.22%      150ms 21.74%  runtime.findRunnable
      20ms  2.90% 68.12%       70ms 10.14%  runtime.reentersyscall
      10ms  1.45% 69.57%      350ms 50.72%  bytes.(*Buffer).ReadFrom
      10ms  1.45% 71.01%       10ms  1.45%  crypt_proto/internal/queue.(*MemoryStore).apply
      10ms  1.45% 72.46%       10ms  1.45%  crypto/internal/fips140/aes.encryptBlock
      10ms  1.45% 73.91%       20ms  2.90%  crypto/internal/fips140/aes/gcm.open
      10ms  1.45% 75.36%       50ms  7.25%  github.com/tidwall/gjson.Get
(pprof) 




