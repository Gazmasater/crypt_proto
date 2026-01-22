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





package calculator

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"

	"crypt_proto/internal/queue"
	"crypt_proto/pkg/models"
)

const (
	minVolumeUSDT = 20.0
	feeMul        = 0.999
)

type LegIndex struct {
	Key    string
	Symbol string
	IsBuy  bool
}

type Triangle struct {
	A, B, C  string
	Legs     [3]LegIndex
	RoughMax float64 // грубая оценка для быстрого отсева
}

type Calculator struct {
	mem      *queue.MemoryStore
	bySymbol map[string][]*Triangle
	fileLog  *log.Logger
}

// NewCalculator строит индекс symbol -> triangles
func NewCalculator(mem *queue.MemoryStore, triangles []*Triangle) *Calculator {
	f, err := os.OpenFile("arb_opportunities.log",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("failed to open log: %v", err)
	}

	bySymbol := make(map[string][]*Triangle, 1024)
	for _, t := range triangles {
		for _, leg := range t.Legs {
			bySymbol[leg.Symbol] = append(bySymbol[leg.Symbol], t)
		}
	}

	log.Printf("[Calculator] indexed %d symbols\n", len(bySymbol))

	return &Calculator{
		mem:      mem,
		bySymbol: bySymbol,
		fileLog:  log.New(f, "", log.LstdFlags),
	}
}

// Run — считаем только нужные треугольники
func (c *Calculator) Run(in <-chan *models.MarketData) {
	for md := range in {
		c.mem.Push(md)

		tris := c.bySymbol[md.Symbol]
		if len(tris) == 0 {
			continue
		}

		for _, tri := range tris {
			c.calcTriangle(tri)
		}
	}
}

func (c *Calculator) calcTriangle(tri *Triangle) {
	// достаем котировки напрямую
	q0, ok0 := c.mem.Get("KuCoin", tri.Legs[0].Symbol)
	q1, ok1 := c.mem.Get("KuCoin", tri.Legs[1].Symbol)
	q2, ok2 := c.mem.Get("KuCoin", tri.Legs[2].Symbol)
	if !ok0 || !ok1 || !ok2 {
		return
	}

	// выбираем цену для грубой проверки
	v0 := q0.Bid
	if tri.Legs[0].IsBuy {
		v0 = q0.Ask
	}
	v1 := q1.Bid
	if tri.Legs[1].IsBuy {
		v1 = q1.Ask
	}
	v2 := q2.Bid
	if tri.Legs[2].IsBuy {
		v2 = q2.Ask
	}

	// быстрый отбор треугольников
	if v0 <= v1*v2 || v0 < minVolumeUSDT || v1 < minVolumeUSDT || v2 < minVolumeUSDT {
		return
	}

	// точный расчет объёма
	var usdt0, usdt1, usdt2 float64

	if tri.Legs[0].IsBuy {
		if q0.Ask <= 0 || q0.AskSize <= 0 {
			return
		}
		usdt0 = q0.Ask * q0.AskSize
	} else {
		if q0.Bid <= 0 || q0.BidSize <= 0 {
			return
		}
		usdt0 = q0.Bid * q0.BidSize
	}

	if tri.Legs[1].IsBuy {
		if q1.Ask <= 0 || q1.AskSize <= 0 || q2.Bid <= 0 {
			return
		}
		usdt1 = q1.Ask * q1.AskSize * q2.Bid
	} else {
		if q1.Bid <= 0 || q1.BidSize <= 0 || q2.Bid <= 0 {
			return
		}
		usdt1 = q1.BidSize * q2.Bid
	}

	if q2.Bid <= 0 || q2.BidSize <= 0 {
		return
	}
	usdt2 = q2.Bid * q2.BidSize

	// проверка минимального объёма
	maxUSDT := usdt0
	if usdt1 < maxUSDT {
		maxUSDT = usdt1
	}
	if usdt2 < maxUSDT {
		maxUSDT = usdt2
	}
	if maxUSDT < minVolumeUSDT {
		return
	}

	// расчет прибыли с учетом feeMul
	amount := maxUSDT
	if tri.Legs[0].IsBuy {
		amount = amount / q0.Ask * feeMul
	} else {
		amount = amount * q0.Bid * feeMul
	}
	if tri.Legs[1].IsBuy {
		amount = amount / q1.Ask * feeMul
	} else {
		amount = amount * q1.Bid * feeMul
	}
	if tri.Legs[2].IsBuy {
		amount = amount / q2.Ask * feeMul
	} else {
		amount = amount * q2.Bid * feeMul
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

// ParseTrianglesFromCSV — загружаем треугольники из CSV
func ParseTrianglesFromCSV(path string) ([]*Triangle, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rows, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return nil, err
	}

	var res []*Triangle
	for _, row := range rows[1:] {
		if len(row) < 6 {
			continue
		}

		t := &Triangle{
			A: strings.TrimSpace(row[0]),
			B: strings.TrimSpace(row[1]),
			C: strings.TrimSpace(row[2]),
		}

		for i, leg := range []string{row[3], row[4], row[5]} {
			leg = strings.ToUpper(strings.TrimSpace(leg))
			parts := strings.Fields(leg)
			if len(parts) != 2 {
				continue
			}

			isBuy := parts[0] == "BUY"
			pair := strings.Split(parts[1], "/")
			if len(pair) != 2 {
				continue
			}

			symbol := pair[0] + "-" + pair[1]
			key := "KuCoin|" + symbol

			t.Legs[i] = LegIndex{
				Key:    key,
				Symbol: symbol,
				IsBuy:  isBuy,
			}
		}

		res = append(res, t)
	}

	return res, nil
}








