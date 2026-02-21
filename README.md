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


gaz358@gaz358-BOD-WXX9:~/myprog/arb/cmd/server$ go run .
{"time":"2026-02-21T08:41:19.966673021+03:00","level":"INFO","msg":"starting scalper","version":"0.1.0","exchanges":"okx"}
{"time":"2026-02-21T08:41:19.966724375+03:00","level":"WARN","msg":"OKX API keys not configured, skipping auth check"}
{"time":"2026-02-21T08:41:19.966791017+03:00","level":"INFO","msg":"pprof server starting","addr":"localhost:6060"}
{"time":"2026-02-21T08:41:20.583535623+03:00","level":"INFO","msg":"instruments loaded","exchange":"okx","count":745}
{"time":"2026-02-21T08:41:20.583610278+03:00","level":"INFO","msg":"inventory initialized","usdt_per_exchange":10000,"max_seeds":50,"seed_per_coin":60,"stop_loss_pct":0.05,"seed_expiry":1800000000000,"cross_coins":0}
{"time":"2026-02-21T08:41:20.583767878+03:00","level":"INFO","msg":"trade log file opened","path":"trades.log"}
{"time":"2026-02-21T08:41:20.592421615+03:00","level":"INFO","msg":"triangles built","exchange":"okx","count":512}
{"time":"2026-02-21T08:41:21.383801621+03:00","level":"INFO","msg":"OKX ws connected","http_status":101}
{"time":"2026-02-21T08:41:21.385556968+03:00","level":"INFO","msg":"OKX subscribed to books5","count":745}
{"time":"2026-02-21T08:41:21.385594669+03:00","level":"INFO","msg":"ws started","exchange":"okx","instruments":745}
{"time":"2026-02-21T08:41:21.385672274+03:00","level":"INFO","msg":"startup complete","elapsed":1419012897,"exchanges":"okx","triangles":512,"cross_pairs":745}
{"time":"2026-02-21T08:41:21.385790833+03:00","level":"INFO","msg":"http server starting","addr":"0.0.0.0:8080"}
{"time":"2026-02-21T08:41:51.386201488+03:00","level":"INFO","msg":"heartbeat","uptime":"0d0h0m30s","okx_ping":"123ms(avg:157ms)"}
{"time":"2026-02-21T08:41:51.3862504+03:00","level":"INFO","msg":"triangle stats","scanned":10650,"stale":6009,"volume":1555,"opps":0,"best_pct":0,"best_tri":""}
{"time":"2026-02-21T08:41:51.386299441+03:00","level":"INFO","msg":"inventory","portfolio":10000,"unrealized":0,"free_usdt":10000,"seeds":0,"rated":0,"promoted":0,"demoted":0,"stop_loss":0,"expired":0}
{"time":"2026-02-21T08:42:21.386239862+03:00","level":"INFO","msg":"heartbeat","uptime":"0d0h1m0s","okx_ping":"121ms(avg:145ms)"}
{"time":"2026-02-21T08:42:21.386289808+03:00","level":"INFO","msg":"triangle stats","scanned":21840,"stale":11604,"volume":3145,"opps":0,"best_pct":0,"best_tri":""}
{"time":"2026-02-21T08:42:21.386328704+03:00","level":"INFO","msg":"inventory","portfolio":10000,"unrealized":0,"free_usdt":10000,"seeds":0,"rated":0,"promoted":0,"demoted":0,"stop_loss":0,"expired":0}
{"time":"2026-02-21T08:42:51.386668024+03:00","level":"INFO","msg":"heartbeat","uptime":"0d0h1m30s","okx_ping":"122ms(avg:141ms)"}
{"time":"2026-02-21T08:42:51.386709136+03:00","level":"INFO","msg":"triangle stats","scanned":31832,"stale":16175,"volume":3777,"opps":0,"best_pct":0,"best_tri":""}
{"time":"2026-02-21T08:42:51.386746052+03:00","level":"INFO","msg":"inventory","portfolio":10000,"unrealized":0,"free_usdt":10000,"seeds":0,"rated":0,"promoted":0,"demoted":0,"stop_loss":0,"expired":0}

