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
	A, B, C string
	Legs    [3]LegIndex
}

type Calculator struct {
	mem      *queue.MemoryStore
	bySymbol map[string][]*Triangle

	roughMax map[*Triangle]float64 // грубая оценка прибыли

	tmpQuotes [3]queue.Quote // временный массив для котировок

	fileLog *log.Logger
}

// NewCalculator — строим индекс symbol -> triangles и создаем roughMax map
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
		roughMax: make(map[*Triangle]float64, len(triangles)),
		fileLog:  log.New(f, "", log.LstdFlags),
	}
}

// Run — обрабатываем поток котировок
func (c *Calculator) Run(in <-chan *models.MarketData) {
	for md := range in {
		c.mem.Push(md)

		tris := c.bySymbol[md.Symbol]
		if len(tris) == 0 {
			continue
		}

		for _, tri := range tris {
			// достаем котировки и сохраняем в tmpQuotes
			if !c.prepareQuotes(tri) {
				continue
			}

			// обновляем грубую оценку прибыли
			c.updateRoughMaxWithQuotes(tri, &c.tmpQuotes)

			// точный расчёт прибыли
			c.calcTriangleWithQuotes(tri, &c.tmpQuotes)
		}
	}
}

// prepareQuotes — достаем котировки и сохраняем в tmpQuotes
func (c *Calculator) prepareQuotes(tri *Triangle) bool {
	for i, leg := range tri.Legs {
		q, ok := c.mem.Get("KuCoin", leg.Symbol)
		if !ok {
			return false
		}
		c.tmpQuotes[i] = q
	}
	return true
}

// updateRoughMaxWithQuotes — грубая оценка прибыли
func (c *Calculator) updateRoughMaxWithQuotes(tri *Triangle, q *[3]queue.Quote) {
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

	if v0 <= v1*v2 || v0 < minVolumeUSDT || v1 < minVolumeUSDT || v2 < minVolumeUSDT {
		c.roughMax[tri] = 0
		return
	}

	c.roughMax[tri] = v0 - v1*v2
}

// calcTriangleWithQuotes — точный расчёт объёма и прибыли
func (c *Calculator) calcTriangleWithQuotes(tri *Triangle, q *[3]queue.Quote) {
	if c.roughMax[tri] <= 0 {
		return
	}

	var usdt [3]float64

	// расчёт USDT по каждой ноге
	if tri.Legs[0].IsBuy {
		if q[0].Ask <= 0 || q[0].AskSize <= 0 {
			return
		}
		usdt[0] = q[0].Ask * q[0].AskSize
	} else {
		if q[0].Bid <= 0 || q[0].BidSize <= 0 {
			return
		}
		usdt[0] = q[0].Bid * q[0].BidSize
	}

	if tri.Legs[1].IsBuy {
		if q[1].Ask <= 0 || q[1].AskSize <= 0 || q[2].Bid <= 0 {
			return
		}
		usdt[1] = q[1].Ask * q[1].AskSize * q[2].Bid
	} else {
		if q[1].Bid <= 0 || q[1].BidSize <= 0 || q[2].Bid <= 0 {
			return
		}
		usdt[1] = q[1].BidSize * q[2].Bid
	}

	if q[2].Bid <= 0 || q[2].BidSize <= 0 {
		return
	}
	usdt[2] = q[2].Bid * q[2].BidSize

	// минимальный объём
	maxUSDT := usdt[0]
	if usdt[1] < maxUSDT {
		maxUSDT = usdt[1]
	}
	if usdt[2] < maxUSDT {
		maxUSDT = usdt[2]
	}
	if maxUSDT < minVolumeUSDT {
		return
	}

	// расчёт прибыли с учётом fee
	amount := maxUSDT
	for i := 0; i < 3; i++ {
		if tri.Legs[i].IsBuy {
			amount = amount / q[i].Ask * feeMul
		} else {
			amount = amount * q[i].Bid * feeMul
		}
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

// CSV парсинг без изменений логики, сохраняем Symbol
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






