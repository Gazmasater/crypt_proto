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
Saved profile in /home/gaz358/pprof/pprof.arb.samples.cpu.336.pb.gz
File: arb
Build ID: 7248eb28f7d3b6f697e0b3af6076a1f584b8acc7
Type: cpu
Time: 2026-01-29 14:20:31 MSK
Duration: 30s, Total samples = 120ms (  0.4%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) gaz358@gaz358-BOD-WXX9:~/myprog/crypt_proto$    go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
Fetching profile over HTTP from http://localhost:6060/debug/pprof/profile?seconds=30
Saved profile in /home/gaz358/pprof/pprof.arb.samples.cpu.337.pb.gz
File: arb
Build ID: 7248eb28f7d3b6f697e0b3af6076a1f584b8acc7
Type: cpu
Time: 2026-01-29 14:55:42 MSK
Duration: 30s, Total samples = 160ms ( 0.53%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) gaz358@gaz358-BOD-WXX9:~/myprog/crypt_proto$ ^C
gaz358@gaz358-BOD-WXX9:~/myprog/crypt_proto$ ^C
gaz358@gaz358-BOD-WXX9:~/myprog/crypt_proto$ ^C
gaz358@gaz358-BOD-WXX9:~/myprog/crypt_proto$ ^C
gaz358@gaz358-BOD-WXX9:~/myprog/crypt_proto$    go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
Fetching profile over HTTP from http://localhost:6060/debug/pprof/profile?seconds=30
Saved profile in /home/gaz358/pprof/pprof.arb.samples.cpu.338.pb.gz
File: arb
Build ID: 7248eb28f7d3b6f697e0b3af6076a1f584b8acc7
Type: cpu
Time: 2026-01-29 14:58:16 MSK
Duration: 30s, Total samples = 260ms ( 0.87%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) 



