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




func buildTriangles(markets map[string]Market, anchor string) []Triangle {
	var out []Triangle

	// все монеты, которые торгуются с anchor
	var coins []string
	for _, m := range markets {
		if m.Quote == anchor && !isStable(m.Base) {
			coins = append(coins, m.Base)
		}
	}

	// перебор всех пар (X, Y)
	for i := 0; i < len(coins); i++ {
		for j := 0; j < len(coins); j++ {
			if i == j {
				continue
			}

			X := coins[i]
			Y := coins[j]

			// =============================
			// ВАРИАНТ 1: USDT → X → Y → USDT
			// =============================
			m1, ok1 := markets[X+"_"+anchor]
			m2, ok2 := markets[Y+"_"+X]
			m3, ok3 := markets[Y+"_"+anchor]

			if ok1 && ok2 && ok3 {
				out = append(out, newTriangle(anchor, X, Y, m1, m2, m3))
			}

			// =============================
			// ВАРИАНТ 2: USDT → Y → X → USDT
			// =============================
			m1b, ok1b := markets[Y+"_"+anchor]
			m2b, ok2b := markets[X+"_"+Y]
			m3b, ok3b := markets[X+"_"+anchor]

			if ok1b && ok2b && ok3b {
				out = append(out, newTriangle(anchor, Y, X, m1b, m2b, m3b))
			}
		}
	}

	return out
}

