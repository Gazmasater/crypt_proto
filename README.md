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


ms, okMS := domain.ComputeMaxStartTopOfBook(tr, quotes, c.FeePerLeg)
if okMS {
    safeStart := ms.MaxStart * sf

    // ФИЛЬТР: MIN_START_USDT сравниваем по safeStart в USDT
    if c.MinStart > 0 {
        safeUSDT, okConv := convertToUSDT(safeStart, ms.StartAsset, quotes)
        // если не смогли перевести в USDT — треугольник отбрасываем,
        // раз порог задан в USDT
        if !okConv || safeUSDT < c.MinStart {
            continue
        }
    }
}



