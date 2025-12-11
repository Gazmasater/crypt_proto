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









cmd/
  cryptarb/
    main.go            // запуск: конфиг, DI

internal/
  config/
    config.go          // чтение .env, флагов

  domain/
    triangle.go        // Triangle, Leg, Pair, buildTriangleFromPairs
    quote.go           // Quote, Event
    eval.go            // evalTriangle

  usecase/
    arbitrage/
      service.go       // ядро мониторинга треугольников

  ports/
    ticker_source.go   // интерфейс источника котировок (биржа)
    arb_sink.go        // интерфейс для вывода арбитражей (лог/файл/бот)

  adapters/
    mexc/
      ws_source.go     // реализация TickerSource для MEXC
    kucoin/
      ws_source.go     // реализация TickerSource для KuCoin (позже)
    sink/
      stdout_file.go   // реализация ArbSink (stdout + файл)

  triangles/
    loader.go          // загрузка triangles_markets.csv





