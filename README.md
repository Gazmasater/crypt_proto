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
Saved profile in /home/gaz358/pprof/pprof.arb.samples.cpu.026.pb.gz
File: arb
Build ID: 1f1f95b8e0ca67922f6912d7eb16f8b22a5639b1
Type: cpu
Time: 2026-01-08 00:46:11 MSK
Duration: 30.11s, Total samples = 1.90s ( 6.31%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 1330ms, 70.00% of 1900ms total
Showing top 10 nodes out of 154
      flat  flat%   sum%        cum   cum%
     950ms 50.00% 50.00%      950ms 50.00%  internal/runtime/syscall.Syscall6
     180ms  9.47% 59.47%      180ms  9.47%  runtime.futex
      30ms  1.58% 61.05%       30ms  1.58%  crypto/internal/fips140/aes/gcm.gcmAesDec
      30ms  1.58% 62.63%       70ms  3.68%  github.com/tidwall/gjson.parseObject
      30ms  1.58% 64.21%      440ms 23.16%  runtime.findRunnable
      30ms  1.58% 65.79%       70ms  3.68%  runtime.mallocgc
      20ms  1.05% 66.84%       20ms  1.05%  crypto/internal/fips140/aes.encryptBlockAsm
      20ms  1.05% 67.89%       20ms  1.05%  github.com/tidwall/gjson.parseString
      20ms  1.05% 68.95%       20ms  1.05%  internal/poll.convertErr
      20ms  1.05% 70.00%       20ms  1.05%  runtime.(*randomEnum).next (inline)
(pprof) 


gaz358@gaz358-BOD-WXX9:~/myprog/crypt_proto/cmd/arb$    go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
Fetching profile over HTTP from http://localhost:6060/debug/pprof/profile?seconds=30
Saved profile in /home/gaz358/pprof/pprof.arb.samples.cpu.028.pb.gz
File: arb
Build ID: 119e2d6e20f521b0a5339543d15161b36cc093bb
Type: cpu
Time: 2026-01-08 01:09:46 MSK
Duration: 30s, Total samples = 1.72s ( 5.73%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 1170ms, 68.02% of 1720ms total
Showing top 10 nodes out of 156
      flat  flat%   sum%        cum   cum%
     760ms 44.19% 44.19%      760ms 44.19%  internal/runtime/syscall.Syscall6
     170ms  9.88% 54.07%      170ms  9.88%  runtime.futex
      50ms  2.91% 56.98%       50ms  2.91%  runtime.nextFreeFast (inline)
      40ms  2.33% 59.30%       40ms  2.33%  runtime.nanotime
      30ms  1.74% 61.05%       80ms  4.65%  github.com/tidwall/gjson.parseObject
      30ms  1.74% 62.79%       30ms  1.74%  github.com/tidwall/gjson.parseSquash
      30ms  1.74% 64.53%       40ms  2.33%  runtime.stealWork
      20ms  1.16% 65.70%      300ms 17.44%  crypt_proto/internal/collector.(*kucoinWS).handle
      20ms  1.16% 66.86%       20ms  1.16%  crypto/internal/fips140/aes/gcm.gcmAesData
      20ms  1.16% 68.02%       20ms  1.16%  crypto/internal/fips140/aes/gcm.gcmAesDec
(pprof) 





