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

type LegIndex struct {
	Key    string
	Symbol string
	IsBuy  bool
}

type Triangle struct {
	A, B, C string
	Legs    [3]LegIndex
}

type Calculator struct {
	mem      *queue.MemoryStore
	bySymbol map[string][]*Triangle
	fileLog  *log.Logger
	quotes   [3]queue.Quote // единый массив для переиспользования
}

const (
	minVolumeUSDT = 20.0
	feeMul        = 0.999
	profitPctMin  = 0.001
	profitUSDTMin = 0.02
)

// NewCalculator — строим индекс symbol -> triangles
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

// Run — считаем ТОЛЬКО нужные треугольники
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
	// достаем котировки напрямую в единый массив
	for i, leg := range tri.Legs {
		quote, ok := c.mem.Get("KuCoin", leg.Symbol)
		if !ok {
			return
		}
		c.quotes[i] = quote
	}

	// Rough check — выбрать цену для предварительной проверки
	v := [3]float64{c.quotes[0].Bid, c.quotes[1].Bid, c.quotes[2].Bid}
	for i, leg := range tri.Legs {
		if leg.IsBuy {
			v[i] = c.quotes[i].Ask
		}
	}

	// грубая проверка на потенциальную прибыль и минимальный объём
	if v[0] <= v[1]*v[2] || v[0] < minVolumeUSDT || v[1] < minVolumeUSDT || v[2] < minVolumeUSDT {
		return // треугольник явно невыгодный или слишком мал
	}

	// точный расчёт объёма
	usdt := [3]float64{}
	for i, leg := range tri.Legs {
		q := c.quotes[i]
		switch i {
		case 0:
			if leg.IsBuy {
				if q.Ask <= 0 || q.AskSize <= 0 {
					return
				}
				usdt[0] = q.Ask * q.AskSize
			} else {
				if q.Bid <= 0 || q.BidSize <= 0 {
					return
				}
				usdt[0] = q.Bid * q.BidSize
			}
		case 1:
			if leg.IsBuy {
				if q.Ask <= 0 || q.AskSize <= 0 || c.quotes[2].Bid <= 0 {
					return
				}
				usdt[1] = q.Ask * q.AskSize * c.quotes[2].Bid
			} else {
				if q.Bid <= 0 || q.BidSize <= 0 || c.quotes[2].Bid <= 0 {
					return
				}
				usdt[1] = q.BidSize * c.quotes[2].Bid
			}
		case 2:
			if q.Bid <= 0 || q.BidSize <= 0 {
				return
			}
			usdt[2] = q.Bid * q.BidSize
		}
	}

	// проверка минимального объёма
	maxUSDT := min3(usdt[0], usdt[1], usdt[2])
	if maxUSDT < minVolumeUSDT {
		return
	}

	// расчёт прибыли с учётом fee
	amount := maxUSDT
	for i, leg := range tri.Legs {
		q := c.quotes[i]
		if leg.IsBuy {
			amount = amount / q.Ask * feeMul
		} else {
			amount = amount * q.Bid * feeMul
		}
	}

	profitUSDT := amount - maxUSDT
	profitPct := profitUSDT / maxUSDT

	if profitPct > profitPctMin && profitUSDT > profitUSDTMin {
		msg := fmt.Sprintf(
			"[ARB] %s→%s→%s | %.4f%% | volume=%.2f USDT | profit=%.4f USDT",
			tri.A, tri.B, tri.C,
			profitPct*100, maxUSDT, profitUSDT,
		)
		log.Println(msg)
		c.fileLog.Println(msg)
	}
}

func min3(a, b, c float64) float64 {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}

// CSV без изменений логики, но сразу сохраняем Symbol
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








