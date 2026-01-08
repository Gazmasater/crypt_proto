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




func (c *Calculator) Run() {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {

		for _, tri := range c.triangles {

			// ключи MemoryStore
			k1 := "KuCoin|" + legSymbol(tri.Leg1)
			k2 := "KuCoin|" + legSymbol(tri.Leg2)
			k3 := "KuCoin|" + legSymbol(tri.Leg3)

			q1, ok1 := c.mem.Get("Kucoin", k1)
			q2, ok2 := c.mem.Get("Kucoin", k2)
			q3, ok3 := c.mem.Get("Kucoin", k3)
			if !ok1 || !ok2 || !ok3 {
				continue
			}

			amount := 1.0

			// Leg1
			if strings.HasPrefix(tri.Leg1, "BUY") {
				if q1.Ask == 0 {
					continue
				}
				amount /= q1.Ask
			} else {
				amount *= q1.Bid
			}

			// Leg2
			if strings.HasPrefix(tri.Leg2, "BUY") {
				if q2.Ask == 0 {
					continue
				}
				amount /= q2.Ask
			} else {
				amount *= q2.Bid
			}

			// Leg3
			if strings.HasPrefix(tri.Leg3, "BUY") {
				if q3.Ask == 0 {
					continue
				}
				amount /= q3.Ask
			} else {
				amount *= q3.Bid
			}

			profit := amount - 1.0

			// РЕАЛЬНЫЙ порог
			//	if profit > 0.001 {
			log.Printf(
				"[ARB] %s → %s → %s | profit=%.4f%%",
				tri.A, tri.B, tri.C,
				profit*100,
			)
		}
		//}
	}
}













