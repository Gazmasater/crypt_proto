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
Saved profile in /home/gaz358/pprof/pprof.arb.samples.cpu.044.pb.gz
File: arb
Build ID: 661cf4e16f8eb544a152de46fc45312b5200894b
Type: cpu
Time: 2026-01-09 13:28:49 MSK
Duration: 30s, Total samples = 110ms ( 0.37%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 110ms, 100% of 110ms total
Showing top 10 nodes out of 52
      flat  flat%   sum%        cum   cum%
      20ms 18.18% 18.18%       20ms 18.18%  internal/runtime/syscall.Syscall6
      10ms  9.09% 27.27%       20ms 18.18%  github.com/tidwall/gjson.parseObject
      10ms  9.09% 36.36%       10ms  9.09%  github.com/tidwall/gjson.parseSquash
      10ms  9.09% 45.45%       10ms  9.09%  runtime.acquirem (inline)
      10ms  9.09% 54.55%       10ms  9.09%  runtime.duffcopy
      10ms  9.09% 63.64%       10ms  9.09%  runtime.findRunnable
      10ms  9.09% 72.73%       10ms  9.09%  runtime.getitab
      10ms  9.09% 81.82%       10ms  9.09%  runtime.mapaccess2_faststr
      10ms  9.09% 90.91%       10ms  9.09%  runtime.nanotime1
      10ms  9.09%   100%       20ms 18.18%  strings.Fields
(pprof) 







