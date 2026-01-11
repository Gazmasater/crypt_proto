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

			amount := 1.0 // стартуем с 1 USDT

			// ---------- LEG 1 ----------
			if strings.HasPrefix(tri.Leg1, "BUY") {
				if q1.Ask <= 0 || q1.AskSize <= 0 {
					continue
				}

				amount = amount / q1.Ask
				if amount > q1.AskSize {
					amount = q1.AskSize
				}
				amount *= (1 - fee)

			} else {
				if q1.Bid <= 0 || q1.BidSize <= 0 {
					continue
				}

				if amount > q1.BidSize {
					amount = q1.BidSize
				}
				amount = amount * q1.Bid
				amount *= (1 - fee)
			}

			// ---------- LEG 2 ----------
			if strings.HasPrefix(tri.Leg2, "BUY") {
				if q2.Ask <= 0 || q2.AskSize <= 0 {
					continue
				}

				amount = amount / q2.Ask
				if amount > q2.AskSize {
					amount = q2.AskSize
				}
				amount *= (1 - fee)

			} else {
				if q2.Bid <= 0 || q2.BidSize <= 0 {
					continue
				}

				if amount > q2.BidSize {
					amount = q2.BidSize
				}
				amount = amount * q2.Bid
				amount *= (1 - fee)
			}

			// ---------- LEG 3 ----------
			if strings.HasPrefix(tri.Leg3, "BUY") {
				if q3.Ask <= 0 || q3.AskSize <= 0 {
					continue
				}

				amount = amount / q3.Ask
				if amount > q3.AskSize {
					amount = q3.AskSize
				}
				amount *= (1 - fee)

			} else {
				if q3.Bid <= 0 || q3.BidSize <= 0 {
					continue
				}

				if amount > q3.BidSize {
					amount = q3.BidSize
				}
				amount = amount * q3.Bid
				amount *= (1 - fee)
			}

			profit := amount - 1.0

			if profit <= 0.001 {
				continue
			}

			// ---------- volumes in USDT ----------
			v1 := c.volumeToUSDT("KuCoin", s1, q1.BidSize, q1.Bid)
			v2 := c.volumeToUSDT("KuCoin", s2, q2.BidSize, q2.Bid)
			v3 := c.volumeToUSDT("KuCoin", s3, q3.BidSize, q3.Bid)

			// если хотя бы один leg не приводится к USDT — пропускаем
			if v1 == 0 || v2 == 0 || v3 == 0 {
				continue
			}

			msg := fmt.Sprintf(
				"[ARB] %s → %s → %s | profit=%.3f%% | volumes USDT: [%.0f / %.0f / %.0f]",
				tri.A, tri.B, tri.C,
				profit*100,
				v1, v2, v3,
			)

			log.Println(msg)
			c.fileLog.Println(msg)
		}
	}
}



func (c *Calculator) volumeToUSDT(
	exchange string,
	symbol string,
	size float64,
	price float64,
) float64 {

	parts := strings.Split(symbol, "-")
	if len(parts) != 2 {
		return 0
	}

	base := parts[0]
	quote := parts[1]

	// BASE-QUOTE → объём в QUOTE
	quoteAmount := size * price

	// 1) QUOTE == USDT
	if quote == "USDT" {
		return quoteAmount
	}

	// 2) есть QUOTE-USDT
	if q, ok := c.mem.Get(exchange, quote+"-USDT"); ok && q.Bid > 0 {
		return quoteAmount * q.Bid
	}

	// 3) есть USDT-QUOTE
	if q, ok := c.mem.Get(exchange, "USDT-"+quote); ok && q.Ask > 0 {
		return quoteAmount / q.Ask
	}

	// нет пути в USDT
	return 0
}




2026/01/11 11:05:32 [ARB] USDT → ONT → BTC | profit=0.1461% | volumes: [1197.2179 / 224.5700 / 0.2203]
2026/01/11 13:51:38 [ARB] USDT → TRVL → BTC | profit=1.2632% | volumes: [3561.0000 / 43864.0000 / 1.0466]
2026/01/11 14:28:36 [ARB] USDT → EWT → BTC | profit=0.1783% | volumes: [56.3700 / 34.2900 / 0.8615]
2026/01/11 14:54:19 [ARB] USDT → VRA → BTC | profit=0.3496% | volumes: [479156.0000 / 249166.0000 / 0.8113]
2026/01/11 15:04:34 [ARB] USDT → EWT → BTC | profit=0.3458% | volumes: [256.2600 / 13.3300 / 0.5853]
2026/01/11 15:37:56 [ARB] USDT → EWT → BTC | profit=1.2844% | volumes: [320.1000 / 34.1000 / 0.5285]






