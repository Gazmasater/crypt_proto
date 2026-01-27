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
	"os"
	"strings"
	"time"

	"crypt_proto/internal/queue"
)

const feeM = 0.999

/* ===================== MODELS ===================== */

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
	log      *log.Logger
}

/* ===================== CONSTRUCTOR ===================== */

func NewCalculator(mem *queue.MemoryStore, triangles []*Triangle) *Calculator {
	f, err := os.OpenFile(
		"arb_opportunities.log",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	if err != nil {
		log.Fatalf("failed to open log: %v", err)
	}

	bySymbol := make(map[string][]*Triangle, 512)

	for _, t := range triangles {
		for _, leg := range t.Legs {
			bySymbol[leg.Symbol] = append(bySymbol[leg.Symbol], t)
		}
	}

	log.Printf("[Calculator] indexed %d symbols\n", len(bySymbol))

	return &Calculator{
		mem:      mem,
		bySymbol: bySymbol,
		log:      log.New(f, "", log.LstdFlags),
	}
}

/* ===================== RUN LOOP ===================== */

// Run — периодический пересчёт (pull-модель)
func (c *Calculator) Run() {
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		c.calculate()
	}
}

/* ===================== CORE ===================== */

func (c *Calculator) calculate() {
	seen := make(map[*Triangle]struct{}, 256)

	for _, tris := range c.bySymbol {
		for _, tri := range tris {
			if _, ok := seen[tri]; ok {
				continue
			}
			seen[tri] = struct{}{}
			c.calcTriangle(tri)
		}
	}
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

	maxUSDT := min(
		legLimitUSDT(tri.Legs[0], q[0], 0),
		legLimitUSDT(tri.Legs[1], q[1], q[2].Bid),
		legLimitUSDT(tri.Legs[2], q[2], 0),
	)

	if maxUSDT <= 0 {
		return
	}

	amount := maxUSDT

	amount = applyLeg(amount, tri.Legs[0], q[0])
	amount = applyLeg(amount, tri.Legs[1], q[1])
	amount = applyLeg(amount, tri.Legs[2], q[2])

	profitUSDT := amount - maxUSDT
	profitPct := profitUSDT / maxUSDT

	if profitPct > 0.001 && profitUSDT > 0.02 {
		msg := fmt.Sprintf(
			"%s [ARB] %s→%s→%s | %.4f%% | volume=%.2f USDT | profit=%.4f USDT",
			time.Now().Format("2006/01/02 15:04:05.000"),
			tri.A, tri.B, tri.C,
			profitPct*100,
			maxUSDT,
			profitUSDT,
		)

		log.Println(msg)
		c.log.Println(msg)
	}
}

/* ===================== HELPERS ===================== */

func legLimitUSDT(leg LegIndex, q queue.Quote, extra float64) float64 {
	if leg.IsBuy {
		if q.Ask <= 0 || q.AskSize <= 0 {
			return 0
		}
		limit := q.Ask * q.AskSize
		if extra > 0 {
			limit *= extra
		}
		return limit
	}

	if q.Bid <= 0 || q.BidSize <= 0 {
		return 0
	}
	limit := q.Bid * q.BidSize
	if extra > 0 {
		limit = q.BidSize * extra
	}
	return limit
}

func applyLeg(amount float64, leg LegIndex, q queue.Quote) float64 {
	if leg.IsBuy {
		if q.Ask <= 0 {
			return 0
		}
		return amount / q.Ask * feeM
	}
	if q.Bid <= 0 {
		return 0
	}
	return amount * q.Bid * feeM
}

func min(a, b, c float64) float64 {
	if a <= 0 || b <= 0 || c <= 0 {
		return 0
	}
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

/* ===================== CSV PARSER ===================== */

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

		for i, raw := range row[3:6] {
			parts := strings.Fields(strings.ToUpper(strings.TrimSpace(raw)))
			if len(parts) != 2 {
				continue
			}

			isBuy := parts[0] == "BUY"
			pair := strings.Split(parts[1], "/")
			if len(pair) != 2 {
				continue
			}

			t.Legs[i] = LegIndex{
				Symbol: pair[0] + "-" + pair[1],
				IsBuy:  isBuy,
			}
		}

		res = append(res, t)
	}

	return res, nil
}




