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



/arb_project
 ├─ main.go                  # Точка входа, запуск Collector, Calculator, Queue, Executor
 ├─ collector/               # Пакет для сбора данных с бирж
 │   ├─ collector.go         # Интерфейс Collector
 │   ├─ okx_collector.go     # Реализация Collector для OKX
 │   ├─ mexc_collector.go    # Реализация Collector для MEXC
 │   └─ kucoin_collector.go  # Реализация Collector для KuCoin
 ├─ calculator/              # Пакет для расчёта арбитража
 │   ├─ calculator.go        # Интерфейс Calculator
 │   └─ arb.go               # Реализация arb
 ├─ queue/                   # Пакет для очереди сигналов
 │   ├─ in_memory_queue.go   # InMemoryQueue для теста
 │   └─ redis_queue.go       # RedisQueue для продакшн
 ├─ executor/                # Пакет для исполнения сигналов
 │   ├─ executor.go          # Интерфейс Executor
 │   └─ executor_impl.go     # Реализация Executor (реальные ордера)
 ├─ models/                  # Общие структуры данных
 │   ├─ market_data.go       # Структура MarketData
 │   └─ signal.go            # Структура Signal
 └─ utils/                   # Вспомогательные функции
     └─ helpers.go           # parseFloat, логирование, конвертация



