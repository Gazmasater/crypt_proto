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




2026/01/04 11:33:05 EXCHANGE: kucoin
2026/01/04 11:33:05 pprof on http://localhost:6060/debug/pprof/
2026/01/04 11:33:05 Loaded 246 unique symbols from ../exchange/data/kucoin_triangles_usdt.csv
2026/01/04 11:33:05 [KuCoin] init WS
2026/01/04 11:33:05 [KuCoin] request bullet-public
2026/01/04 11:33:06 [KuCoin] wsURL ready, pingInterval=0s
2026/01/04 11:33:06 [KuCoin] connect: wss://ws-api-spot.kucoin.com/?token=2neAiuYvAU61ZDXANAGAsiL4-iAExhsBXZxftpOeh_55i3Ysy2q2LEsEWU64mdzUOPusi34M_wGoSf7iNyEWJyRQQkns52iq8lDQcWAtencR3R2PfzVokdiYB9J6i9GjsxUuhPw3Blq6rhZlGykT3Vp1phUafnulOOpts-MEmEHtGZR-Jl-TQpU-gjzmv_-6JBvJHl5Vs9Y=.56dJLGSKb9deoowHoR7FKQ==&connectId=1767515586042933615
2026/01/04 11:33:07 [KuCoin] connected
2026/01/04 11:33:07 [KuCoin] subscribed 1 / 246
2026/01/04 11:33:11 [KuCoin] subscribed 11 / 246
2026/01/04 11:33:15 [KuCoin] subscribed 21 / 246
2026/01/04 11:33:19 [KuCoin] subscribed 31 / 246
2026/01/04 11:33:23 [KuCoin] subscribed 41 / 246
2026/01/04 11:33:27 [KuCoin] subscribed 51 / 246
2026/01/04 11:33:31 [KuCoin] subscribed 61 / 246
2026/01/04 11:33:35 [KuCoin] subscribed 71 / 246
2026/01/04 11:33:39 [KuCoin] subscribed 81 / 246
2026/01/04 11:33:43 [KuCoin] subscribed 91 / 246
2026/01/04 11:33:47 [KuCoin] subscribed 101 / 246
2026/01/04 11:33:51 [KuCoin] subscribed 111 / 246
2026/01/04 11:33:55 [KuCoin] subscribed 121 / 246
2026/01/04 11:33:59 [KuCoin] subscribed 131 / 246
2026/01/04 11:34:03 [KuCoin] subscribed 141 / 246
2026/01/04 11:34:07 [KuCoin] subscribed 151 / 246
2026/01/04 11:34:11 [KuCoin] subscribed 161 / 246
2026/01/04 11:34:15 [KuCoin] subscribed 171 / 246
2026/01/04 11:34:19 [KuCoin] subscribed 181 / 246
2026/01/04 11:34:23 [KuCoin] subscribed 191 / 246
2026/01/04 11:34:27 [KuCoin] subscribed 201 / 246
2026/01/04 11:34:31 [KuCoin] subscribed 211 / 246
2026/01/04 11:34:35 [KuCoin] subscribed 221 / 246
2026/01/04 11:34:39 [KuCoin] subscribed 231 / 246
2026/01/04 11:34:43 [KuCoin] subscribed 241 / 246
2026/01/04 11:34:45 [KuCoin] subscribed TOTAL: 246 symbols
2026/01/04 11:34:45 [KuCoin] readLoop started
panic: non-positive interval for NewTicker

goroutine 52 [running]:
time.NewTicker(0xc000305720?)
        /usr/local/go/src/time/tick.go:38 +0xbb
crypt_proto/internal/collector.(*KuCoinCollector).pingLoop(0xc0000ea200)
        /home/gaz358/myprog/crypt_proto/internal/collector/kucoin_collector.go:149 +0x3a
created by crypt_proto/internal/collector.(*KuCoinCollector).Start in goroutine 1
        /home/gaz358/myprog/crypt_proto/internal/collector/kucoin_collector.go:84 +0x27e
exit status 2
gaz358@gaz358-BOD-WXX9:~/myprog/crypt_proto/cmd/arb$ 








