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
}

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
func (c *Calculator) Run() {
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		c.calculate()
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

	amount := maxUSDT

	if tri.Legs[0].IsBuy {
		amount = amount / q[0].Ask * feeM
	} else {
		amount = amount * q[0].Bid * feeM
	}

	if tri.Legs[1].IsBuy {
		amount = amount / q[1].Ask * feeM
	} else {
		amount = amount * q[1].Bid * feeM
	}

	if tri.Legs[2].IsBuy {
		amount = amount / q[2].Ask * feeM
	} else {
		amount = amount * q[2].Bid * feeM
	}

	profitUSDT := amount - maxUSDT
	profitPct := profitUSDT / maxUSDT

	if profitPct > 0.001 && profitUSDT > 0.02 {
		ts := time.Now().Format("2006/01/02 15:04:05.000")
		msg := fmt.Sprintf(
			"%s [ARB] %s→%s→%s | %.4f%% | volume=%.2f USDT | profit=%.4f USDT",
			ts,
			tri.A, tri.B, tri.C,
			profitPct*100, maxUSDT, profitUSDT,
		)

		log.Println(msg)
		c.fileLog.Println(msg)

	}
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


[{
	"resource": "/home/gaz358/myprog/crypt_proto/internal/calculator/arb.go",
	"owner": "_generated_diagnostic_collection_name_#0",
	"code": {
		"value": "MissingFieldOrMethod",
		"target": {
			"$mid": 1,
			"path": "/golang.org/x/tools/internal/typesinternal",
			"scheme": "https",
			"authority": "pkg.go.dev",
			"fragment": "MissingFieldOrMethod"
		}
	},
	"severity": 8,
	"message": "c.calculate undefined (type *Calculator has no field or method calculate)",
	"source": "compiler",
	"startLineNumber": 64,
	"startColumn": 5,
	"endLineNumber": 64,
	"endColumn": 14,
	"modelVersionId": 3,
	"origin": "extHost1"
}]



