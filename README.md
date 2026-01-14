Название API
9623527002

6966b78122ca320001d2acae
fa1e37ae-21ff-4257-844d-3dcd21d26ccd





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



Fetching profile over HTTP from http://localhost:6060/debug/pprof/profile?seconds=30
Saved profile in /home/gaz358/pprof/pprof.arb.samples.cpu.128.pb.gz
File: arb
Build ID: d8095edd0e7b84fc1bd6776bcad2be691a5b7dcc
Type: cpu
Time: 2026-01-15 00:53:04 MSK
Duration: 30s, Total samples = 820ms ( 2.73%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 540ms, 65.85% of 820ms total
Showing top 10 nodes out of 120
      flat  flat%   sum%        cum   cum%
     370ms 45.12% 45.12%      370ms 45.12%  internal/runtime/syscall.Syscall6
      30ms  3.66% 48.78%       30ms  3.66%  runtime.mapaccess2_faststr
      20ms  2.44% 51.22%       20ms  2.44%  aeshashbody
      20ms  2.44% 53.66%       20ms  2.44%  github.com/tidwall/gjson.parseSquash
      20ms  2.44% 56.10%       20ms  2.44%  github.com/tidwall/gjson.parseString
      20ms  2.44% 58.54%       20ms  2.44%  runtime.futex
      20ms  2.44% 60.98%       20ms  2.44%  runtime.memmove
      20ms  2.44% 63.41%       20ms  2.44%  strconv.atof64exact
      10ms  1.22% 64.63%      340ms 41.46%  bytes.(*Buffer).ReadFrom
      10ms  1.22% 65.85%       80ms  9.76%  crypt_proto/internal/calculator.(*Calculator).Run
(pprof) 


