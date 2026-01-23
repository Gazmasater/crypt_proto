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
Saved profile in /home/gaz358/pprof/pprof.arb.samples.cpu.249.pb.gz
File: arb
Build ID: 7540ef780f428f30ea5490419d6c0666ea347d2a
Type: cpu
Time: 2026-01-24 02:15:33 MSK
Duration: 30s, Total samples = 630ms ( 2.10%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 420ms, 66.67% of 630ms total
Showing top 10 nodes out of 118
      flat  flat%   sum%        cum   cum%
     270ms 42.86% 42.86%      270ms 42.86%  internal/runtime/syscall.Syscall6
      30ms  4.76% 47.62%       30ms  4.76%  runtime.futex
      30ms  4.76% 52.38%       30ms  4.76%  runtime.typePointers.next
      20ms  3.17% 55.56%       30ms  4.76%  github.com/tidwall/gjson.parseObject
      20ms  3.17% 58.73%       20ms  3.17%  runtime.memclrNoHeapPointers
      10ms  1.59% 60.32%       10ms  1.59%  bufio.(*Reader).Discard
      10ms  1.59% 61.90%       10ms  1.59%  bytes.(*Buffer).Bytes (inline)
      10ms  1.59% 63.49%       10ms  1.59%  crypto/internal/fips140/aes/gcm.gcmAesData
      10ms  1.59% 65.08%       20ms  3.17%  crypto/tls.(*halfConn).decrypt
      10ms  1.59% 66.67%      250ms 39.68%  github.com/gorilla/websocket.(*Conn).read
(pprof) 



