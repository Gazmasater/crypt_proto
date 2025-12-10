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




2025/12/10 14:24:21.785430 [WS #10] Pong через 274.903338ms
2025/12/10 14:24:21.785556 [WS #3] Pong через 275.037248ms
2025/12/10 14:24:21.785987 [WS #7] Pong через 275.5602ms
2025/12/10 14:24:21.799730 [WS #1] Pong через 289.169191ms
2025/12/10 14:24:21.803034 [WS #6] Pong через 292.581285ms
2025/12/10 14:24:21.804146 [WS #4] Pong через 293.675262ms
2025/12/10 14:24:21.804488 [WS #8] Pong через 294.073703ms
2025/12/10 14:24:21.805072 [WS #11] Pong через 294.60639ms


