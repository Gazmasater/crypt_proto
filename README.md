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





BOOK_INTERVAL=100ms
SYMBOLS_FILE=triangles_markets.csv
DEBUG=false

# Комиссия (в процентах, одна нога)
FEE_PCT=0.1          # 0.1% = 0.001

# Минимальная прибыль по кругу (в процентах)
MIN_PROFIT_PCT=0.3   # 0.3% = 0.003



// Консумер: на КАЖДОМ тике пересчитывает все треугольники
go func(tris []Triangle, ctx context.Context) {
	last := make(map[string]Quote)

	for {
		select {
		case ev, ok := <-events:
			if !ok {
				return
			}

			// обновляем последнюю котировку по символу
			last[ev.Symbol] = Quote{
				Bid:    ev.Bid,
				Ask:    ev.Ask,
				BidQty: ev.BidQty,
				AskQty: ev.AskQty,
			}

			// считаем треугольники на каждом тике
			prof := findProfitableTriangles(tris, last)
			if len(prof) == 0 {
				// прибыльных нет — вообще ничего не выводим
				continue
			}

			fmt.Printf("\nquotes known: %d symbols, profitable triangles: %d\n",
				len(last), len(prof))

			maxShow := 5
			if len(prof) < maxShow {
				maxShow = len(prof)
			}
			for i := 0; i < maxShow; i++ {
				printTriangleWithDetails(prof[i], last)
			}

		case <-ctx.Done():
			return
		}
	}
}(tris, ctx)









