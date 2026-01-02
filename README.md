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



internal/
├── exchange/
│   ├── exchange.go          ← интерфейсы + общие типы
│   ├── registry.go          ← фабрика (выбор биржи)
│   │
│   ├── kucoin/
│   │   ├── client.go        ← HTTP / WS
│   │   ├── markets.go       ← загрузка raw рынков
│   │   ├── normalizer.go    ← KuCoin → Market
│   │   └── adapter.go       ← implements Exchange
│   │
│   ├── mexc/
│   │   ├── client.go
│   │   ├── markets.go
│   │   ├── normalizer.go
│   │   └── adapter.go
│   │
│   ├── okx/
│   │   ├── client.go
│   │   ├── markets.go
│   │   ├── normalizer.go
│   │   └── adapter.go
│
├── arbitrage/
│   ├── triangle.go          ← Triangle, Leg
│   ├── generator.go         ← построение A-B-C
│   ├── filters.go           ← стейблы, направления
│   └── csv.go               ← экспорт
│
├── model/
│   ├── market.go            ← универсальный Market
│   ├── symbol_filter.go
│   └── enums.go
│
└── main.go



internal/
├── exchange/
├── arbitrage/
│   ├── triangle.go
│   ├── generator.go
│   ├── csv.go        ← логика сохранения
├── output/
│   └── triangles/
│       ├── kucoin.csv
│       ├── okx.csv
│       └── mexc.csv
