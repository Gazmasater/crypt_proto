Если оставить только нужное:

p99 execution latency
Micro-volatility (100 мс)
Fill ratio
Capture rate
Inventory drift




Название API
9623527002

696935c42a6dcd00013273f2
b348b686-55ff-4290-897b-02d55f815f65




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





func (c *Calculator) calcTriangle(tri *Triangle) {
	var q [3]queue.Quote

	for i, leg := range tri.Legs {
		quote, ok := c.mem.Get("KuCoin", leg.Symbol)
		if !ok {
			return
		}
		q[i] = quote
	}

	// rough volume check — грубая оценка максимального объёма
	const minVolumeUSDT = 20.0

	v0 := q[0].Bid
	if tri.Legs[0].IsBuy {
		v0 = q[0].Ask
	}
	v1 := q[1].Bid
	if tri.Legs[1].IsBuy {
		v1 = q[1].Ask
	}
	v2 := q[2].Bid
	if tri.Legs[2].IsBuy {
		v2 = q[2].Ask
	}

	// грубая проверка возможности прибыли (без fee, без точного объёма)
	if v0 <= v1*v2 {
		return // явно убыточный треугольник
	}

	// точный расчет объёма
	var usdtLimits [3]float64

	// LEG 1
	if tri.Legs[0].IsBuy {
		if q[0].Ask <= 0 || q[0].AskSize <= 0 {
			return
		}
		usdtLimits[0] = q[0].Ask * q[0].AskSize
	} else {
		if q[0].Bid <= 0 || q[0].BidSize <= 0 {
			return
		}
		usdtLimits[0] = q[0].Bid * q[0].BidSize
	}

	// LEG 2
	if tri.Legs[1].IsBuy {
		if q[1].Ask <= 0 || q[1].AskSize <= 0 || q[2].Bid <= 0 {
			return
		}
		usdtLimits[1] = q[1].Ask * q[1].AskSize * q[2].Bid
	} else {
		if q[1].Bid <= 0 || q[1].BidSize <= 0 || q[2].Bid <= 0 {
			return
		}
		usdtLimits[1] = q[1].BidSize * q[2].Bid
	}

	// LEG 3
	if q[2].Bid <= 0 || q[2].BidSize <= 0 {
		return
	}
	usdtLimits[2] = q[2].Bid * q[2].BidSize

	// точная проверка минимального объёма
	maxUSDT := usdtLimits[0]
	if maxUSDT < minVolumeUSDT {
		return
	}

	if usdtLimits[1] < minVolumeUSDT || usdtLimits[1] < maxUSDT {
		maxUSDT = usdtLimits[1]
	}

	if usdtLimits[2] < minVolumeUSDT || usdtLimits[2] < maxUSDT {
		maxUSDT = usdtLimits[2]
	}

	// расчет прибыли с учетом feeMul
	amount := maxUSDT
	const feeMul = 0.999

	if tri.Legs[0].IsBuy {
		amount = amount / q[0].Ask * feeMul
	} else {
		amount = amount * q[0].Bid * feeMul
	}

	if tri.Legs[1].IsBuy {
		amount = amount / q[1].Ask * feeMul
	} else {
		amount = amount * q[1].Bid * feeMul
	}

	if tri.Legs[2].IsBuy {
		amount = amount / q[2].Ask * feeMul
	} else {
		amount = amount * q[2].Bid * feeMul
	}

	profitUSDT := amount - maxUSDT
	profitPct := profitUSDT / maxUSDT

	if profitPct > 0.001 && profitUSDT > 0.02 {
		msg := fmt.Sprintf(
			"[ARB] %s→%s→%s | %.4f%% | volume=%.2f USDT | profit=%.4f USDT",
			tri.A, tri.B, tri.C,
			profitPct*100, maxUSDT, profitUSDT,
		)
		log.Println(msg)
		c.fileLog.Println(msg)
	}
}






