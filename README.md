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




go run -race main.go


GOMAXPROCS=8 go run -race main.go




package calculator

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
	"strings"

	"crypt_proto/internal/queue"
	"crypt_proto/pkg/models"
)

const feeM = 0.9992

type LegIndex struct {
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
}

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

// maxLegInput возвращает максимальный входной объем, который может принять нога.
// BUY A/B  -> вход в quote, лимит = Ask * AskSize
// SELL A/B -> вход в base,  лимит = BidSize
func maxLegInput(leg LegIndex, q queue.Quote) (float64, bool) {
	if leg.IsBuy {
		if q.Ask <= 0 || q.AskSize <= 0 {
			return 0, false
		}
		return q.Ask * q.AskSize, true
	}

	if q.Bid <= 0 || q.BidSize <= 0 {
		return 0, false
	}
	return q.BidSize, true
}

// forwardAmount считает, сколько выйдет после исполнения ноги.
func forwardAmount(in float64, leg LegIndex, q queue.Quote) (float64, bool) {
	if in <= 0 {
		return 0, false
	}

	if leg.IsBuy {
		if q.Ask <= 0 {
			return 0, false
		}
		return (in / q.Ask) * feeM, true
	}

	if q.Bid <= 0 {
		return 0, false
	}
	return (in * q.Bid) * feeM, true
}

// backwardAmount переводит допустимый выход ноги назад во вход ноги.
// То есть отвечает на вопрос:
// "какой вход нужен, чтобы получить out после этой ноги?"
func backwardAmount(out float64, leg LegIndex, q queue.Quote) (float64, bool) {
	if out <= 0 {
		return 0, false
	}

	if leg.IsBuy {
		if q.Ask <= 0 || feeM <= 0 {
			return 0, false
		}
		return (out * q.Ask) / feeM, true
	}

	if q.Bid <= 0 || feeM <= 0 {
		return 0, false
	}
	return out / (q.Bid * feeM), true
}

func (c *Calculator) calcTriangle(tri *Triangle) {
	var q [3]queue.Quote

	for i, leg := range tri.Legs {
		quote, ok := c.mem.Get("KuCoin", leg.Symbol)
		if !ok {
			return
		}
		q[i] = quote
	}

	// 1) Лимит по 1-й ноге уже в стартовом активе
	limit1, ok := maxLegInput(tri.Legs[0], q[0])
	if !ok {
		return
	}

	// 2) Лимит по 2-й ноге -> назад через 1-ю ногу
	leg2InputCap, ok := maxLegInput(tri.Legs[1], q[1])
	if !ok {
		return
	}

	limit2, ok := backwardAmount(leg2InputCap, tri.Legs[0], q[0])
	if !ok {
		return
	}

	// 3) Лимит по 3-й ноге -> назад через 2-ю и 1-ю ноги
	leg3InputCap, ok := maxLegInput(tri.Legs[2], q[2])
	if !ok {
		return
	}

	beforeLeg2, ok := backwardAmount(leg3InputCap, tri.Legs[1], q[1])
	if !ok {
		return
	}

	limit3, ok := backwardAmount(beforeLeg2, tri.Legs[0], q[0])
	if !ok {
		return
	}

	maxStart := math.Min(limit1, math.Min(limit2, limit3))
	if maxStart <= 0 || math.IsNaN(maxStart) || math.IsInf(maxStart, 0) {
		return
	}

	amount := maxStart

	amount, ok = forwardAmount(amount, tri.Legs[0], q[0])
	if !ok {
		return
	}

	amount, ok = forwardAmount(amount, tri.Legs[1], q[1])
	if !ok {
		return
	}

	amount, ok = forwardAmount(amount, tri.Legs[2], q[2])
	if !ok {
		return
	}

	profit := amount - maxStart
	profitPct := profit / maxStart

	if profitPct > 0.001 && maxStart > 50 {
		msg := fmt.Sprintf(
			"[ARB] %s→%s→%s | %.4f%% | volume=%.2f | profit=%.4f",
			tri.A, tri.B, tri.C,
			profitPct*100, maxStart, profit,
		)
		log.Println(msg)
		c.fileLog.Println(msg)
	}
}

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

			t.Legs[i] = LegIndex{
				Symbol: symbol,
				IsBuy:  isBuy,
			}
		}

		res = append(res, t)
	}

	return res, nil
}




if profitPct > 0.0 && maxStart > 0 {
		msg := fmt.Sprintf(
			"[ARB] %s→%s→%s | %.4f%% | volume=%.2f | profit=%.4f",
			tri.A, tri.B, tri.C,
			profitPct*100, maxStart, profit,
		)
		log.Println(msg)
		c.fileLog.Println(msg)
	}
