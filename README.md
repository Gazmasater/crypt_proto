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




type LiveTriangle struct {
	Tri    Triangle
	Start  int64   // unix ms, когда треугольник стал прибыльным
	Stop   int64   // unix ms, когда треугольник перестал быть прибыльным
	MaxProfit float64
}

type Calculator struct {
	mem       *queue.MemoryStore
	triangles []Triangle
	bySymbol  map[string][]Triangle
	fileLog   *log.Logger

	live map[string]*LiveTriangle // ключ = Leg1|Leg2|Leg3
}

func NewCalculator(mem *queue.MemoryStore, triangles []Triangle) *Calculator {
	f, err := os.OpenFile(
		"arb_opportunities.log",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	if err != nil {
		log.Fatalf("failed to open arb log file: %v", err)
	}

	c := &Calculator{
		mem:       mem,
		triangles: triangles,
		bySymbol:  make(map[string][]Triangle),
		fileLog:   log.New(f, "", log.LstdFlags),
		live:      make(map[string]*LiveTriangle),
	}

	for _, tri := range triangles {
		s1 := legSymbol(tri.Leg1)
		s2 := legSymbol(tri.Leg2)
		s3 := legSymbol(tri.Leg3)

		if s1 != "" {
			c.bySymbol[s1] = append(c.bySymbol[s1], tri)
		}
		if s2 != "" {
			c.bySymbol[s2] = append(c.bySymbol[s2], tri)
		}
		if s3 != "" {
			c.bySymbol[s3] = append(c.bySymbol[s3], tri)
		}
	}

	return c
}

func (c *Calculator) calcTriangle(tri *Triangle) {
	s1 := legSymbol(tri.Leg1)
	s2 := legSymbol(tri.Leg2)
	s3 := legSymbol(tri.Leg3)

	q1, ok1 := c.mem.Get("KuCoin", s1)
	q2, ok2 := c.mem.Get("KuCoin", s2)
	q3, ok3 := c.mem.Get("KuCoin", s3)
	if !ok1 || !ok2 || !ok3 {
		return
	}

	// ===== USDT LIMITS =====
	var usdtLimits [3]float64
	i := 0

	if strings.HasPrefix(tri.Leg1, "BUY") {
		if q1.Ask <= 0 || q1.AskSize <= 0 {
			return
		}
		usdtLimits[i] = q1.Ask * q1.AskSize
	} else {
		if q1.Bid <= 0 || q1.BidSize <= 0 {
			return
		}
		usdtLimits[i] = q1.Bid * q1.BidSize
	}
	i++

	if strings.HasPrefix(tri.Leg2, "BUY") {
		if q2.Ask <= 0 || q2.AskSize <= 0 || q3.Bid <= 0 {
			return
		}
		usdtLimits[i] = q2.Ask * q2.AskSize * q3.Bid
	} else {
		if q2.Bid <= 0 || q2.BidSize <= 0 || q3.Bid <= 0 {
			return
		}
		usdtLimits[i] = q2.BidSize * q3.Bid
	}
	i++

	if q3.Bid <= 0 || q3.BidSize <= 0 {
		return
	}
	usdtLimits[i] = q3.Bid * q3.BidSize

	// ===== MIN LIMIT =====
	maxUSDT := usdtLimits[0]
	if usdtLimits[1] < maxUSDT {
		maxUSDT = usdtLimits[1]
	}
	if usdtLimits[2] < maxUSDT {
		maxUSDT = usdtLimits[2]
	}
	if maxUSDT <= 0 {
		return
	}

	// ===== ПРОГОН =====
	amount := maxUSDT

	if strings.HasPrefix(tri.Leg1, "BUY") {
		amount = amount / q1.Ask * (1 - fee)
	} else {
		amount = amount * q1.Bid * (1 - fee)
	}

	if strings.HasPrefix(tri.Leg2, "BUY") {
		amount = amount / q2.Ask * (1 - fee)
	} else {
		amount = amount * q2.Bid * (1 - fee)
	}

	if strings.HasPrefix(tri.Leg3, "BUY") {
		amount = amount / q3.Ask * (1 - fee)
	} else {
		amount = amount * q3.Bid * (1 - fee)
	}

	profitUSDT := amount - maxUSDT
	profitPct := profitUSDT / maxUSDT

	key := tri.Leg1 + "|" + tri.Leg2 + "|" + tri.Leg3

	// ===== Пороговая прибыль =====
	if profitPct > 0.001 && profitUSDT > 0.02 {
		lt, ok := c.live[key]
		if !ok {
			// создаём новый живой треугольник
			c.live[key] = &LiveTriangle{
				Tri:       *tri,
				Start:     time.Now().UnixMilli(),
				MaxProfit: profitUSDT,
			}
		} else {
			// обновляем макс. прибыль
			if profitUSDT > lt.MaxProfit {
				lt.MaxProfit = profitUSDT
			}
		}

		msg := fmt.Sprintf(
			"[ARB] %s → %s → %s | %.4f%% | volume=%.2f USDT | profit=%.4f USDT",
			tri.A, tri.B, tri.C, profitPct*100, maxUSDT, profitUSDT,
		)
		log.Println(msg)
		c.fileLog.Println(msg)

	} else {
		// треугольник перестал быть прибыльным
		if lt, ok := c.live[key]; ok {
			lt.Stop = time.Now().UnixMilli()
			liveTime := lt.Stop - lt.Start
			log.Printf("[LIVE] %s → %s → %s прожил %d ms | max profit %.4f USDT",
				lt.Tri.A, lt.Tri.B, lt.Tri.C, liveTime, lt.MaxProfit)
			delete(c.live, key)
		}
	}
}




