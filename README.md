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






package config

import "time"

// --- MEXC ---
const (
	MEXC_WS       = "wss://wbs-api.mexc.com/ws"
	MEXC_READ_TIMEOUT  = 30 * time.Second
	MEXC_PING_INTERVAL = 10 * time.Second
	MEXC_RECONNECT_DUR = time.Second
)

// --- KuCoin ---
const (
	KUCOIN_WS        = "wss://ws.kucoin.com/endpoint"
	KUCOIN_PING_INTERVAL = 30 * time.Second
	KUCOIN_READ_TIMEOUT  = 60 * time.Second
)

// --- OKX ---
const (
	OKX_WS           = "wss://ws.okx.com:8443/ws/v5/public"
	OKX_PING_INTERVAL = 20 * time.Second
	OKX_READ_TIMEOUT  = 60 * time.Second
)




