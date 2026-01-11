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

			s1 := legSymbol(tri.Leg1)
			s2 := legSymbol(tri.Leg2)
			s3 := legSymbol(tri.Leg3)

			q1, ok1 := c.mem.Get("KuCoin", s1)
			q2, ok2 := c.mem.Get("KuCoin", s2)
			q3, ok3 := c.mem.Get("KuCoin", s3)

			if !ok1 || !ok2 || !ok3 {
				continue
			}

			// ===============================
			// 1. СЧИТАЕМ ЛИМИТ КАЖДОЙ НОГИ В USDT
			// ===============================

			usdtLimits := make([]float64, 0, 3)

			// ---------- LEG 1 ----------
			if strings.HasPrefix(tri.Leg1, "BUY") {
				// BUY XXX/USDT
				if q1.Ask <= 0 || q1.AskSize <= 0 {
					continue
				}
				usdtLimits = append(usdtLimits, q1.Ask*q1.AskSize)
			} else {
				// SELL XXX/USDT
				if q1.Bid <= 0 || q1.BidSize <= 0 {
					continue
				}
				usdtLimits = append(usdtLimits, q1.Bid*q1.BidSize)
			}

			// ---------- LEG 2 ----------
			if strings.HasPrefix(tri.Leg2, "BUY") {
				// BUY A/B → лимит в B → переводим в USDT через LEG3
				if q2.Ask <= 0 || q2.AskSize <= 0 || q3.Bid <= 0 {
					continue
				}
				usdtLimits = append(usdtLimits, q2.Ask*q2.AskSize*q3.Bid)
			} else {
				// SELL A/B → A * (B/USDT)
				if q2.Bid <= 0 || q2.BidSize <= 0 || q3.Bid <= 0 {
					continue
				}
				usdtLimits = append(usdtLimits, q2.BidSize*q3.Bid)
			}

			// ---------- LEG 3 ----------
			// всегда SELL ???/USDT
			if q3.Bid <= 0 || q3.BidSize <= 0 {
				continue
			}
			usdtLimits = append(usdtLimits, q3.Bid*q3.BidSize)

			// ===============================
			// 2. МИНИМАЛЬНЫЙ ИСПОЛНИМЫЙ ОБЪЁМ
			// ===============================

			maxUSDT := usdtLimits[0]
			for _, v := range usdtLimits {
				if v < maxUSDT {
					maxUSDT = v
				}
			}

			if maxUSDT <= 0 {
				continue
			}

			// ===============================
			// 3. ПРОГОН АРБИТРАЖА С maxUSDT
			// ===============================

			amount := maxUSDT // стартуем в USDT

			// ---------- LEG 1 ----------
			if strings.HasPrefix(tri.Leg1, "BUY") {
				amount = amount / q1.Ask
				amount *= (1 - fee)
			} else {
				amount = amount * q1.Bid
				amount *= (1 - fee)
			}

			// ---------- LEG 2 ----------
			if strings.HasPrefix(tri.Leg2, "BUY") {
				amount = amount / q2.Ask
				amount *= (1 - fee)
			} else {
				amount = amount * q2.Bid
				amount *= (1 - fee)
			}

			// ---------- LEG 3 ----------
			if strings.HasPrefix(tri.Leg3, "BUY") {
				amount = amount / q3.Ask
				amount *= (1 - fee)
			} else {
				amount = amount * q3.Bid
				amount *= (1 - fee)
			}

			profitUSDT := amount - maxUSDT
			profitPct := profitUSDT / maxUSDT

			if profitPct > 0.001 {
				msg := fmt.Sprintf(
					"[ARB] %s → %s → %s | profit=%.4f%% | volume=%.2f USDT | profit=%.4f USDT",
					tri.A, tri.B, tri.C,
					profitPct*100,
					maxUSDT,
					profitUSDT,
				)

				log.Println(msg)
				c.fileLog.Println(msg)
			}
		}
	}
}





