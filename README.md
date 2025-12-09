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




gaz358@gaz358-BOD-WXX9:~/myprog/crypt_proto$    go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
Fetching profile over HTTP from http://localhost:6060/debug/pprof/profile?seconds=30
Saved profile in /home/gaz358/pprof/pprof.crypt_proto.samples.cpu.003.pb.gz
File: crypt_proto
Build ID: 2fb330639396ed0792a55855b01f6138ff5f7e9d
Type: cpu
Time: 2025-12-09 17:48:05 MSK
Duration: 30.14s, Total samples = 35.09s (116.42%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top 
Showing nodes accounting for 21.31s, 60.73% of 35.09s total
Dropped 230 nodes (cum <= 0.18s)
Showing top 10 nodes out of 81
      flat  flat%   sum%        cum   cum%
     5.31s 15.13% 15.13%      5.31s 15.13%  aeshashbody
     2.94s  8.38% 23.51%     28.59s 81.48%  main.evalTriangle
     2.20s  6.27% 29.78%      2.82s  8.04%  internal/runtime/maps.(*Iter).Next
     1.78s  5.07% 34.85%      1.78s  5.07%  internal/runtime/maps.ctrlGroup.matchH2 (inline)
     1.69s  4.82% 39.67%      1.69s  4.82%  runtime.futex
     1.66s  4.73% 44.40%      4.65s 13.25%  internal/runtime/maps.(*Map).putSlotSmallFastStr
     1.60s  4.56% 48.96%      8.37s 23.85%  runtime.mapassign_faststr
     1.45s  4.13% 53.09%      4.78s 13.62%  runtime.mapaccess2_faststr
     1.41s  4.02% 57.11%      1.41s  4.02%  memeqbody
     1.27s  3.62% 60.73%      1.27s  3.62%  internal/chacha8rand.block
(pprof) 






