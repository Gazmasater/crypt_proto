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

	// roughMax для каждого треугольника
	roughMax map[*Triangle]float64

	// временный массив котировок для каждого треугольника
	tmpQuotes [3]queue.Quote

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
			if !c.prepareQuotes(tri) {
				continue
			}
			c.updateRoughMaxWithQuotes(tri, c.tmpQuotes[:])
			c.calcTriangleWithQuotes(tri, c.tmpQuotes[:])
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
		c.tmpQuotes[i] = *q
	}
	return true
}

// updateRoughMaxWithQuotes — грубая оценка прибыли (без объёмов и fee)
func (c *Calculator) updateRoughMaxWithQuotes(tri *Triangle, q []queue.Quote) {
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
func (c *Calculator) calcTriangleWithQuotes(tri *Triangle, q []queue.Quote) {
	if c.roughMax[tri] <= 0 {
		return
	}

	var usdt [3]float64

	// расчёт объёмов для каждой ноги
	for i, leg := range tri.Legs {
		switch i {
		case 0:
			if leg.IsBuy {
				if q[i].Ask <= 0 || q[i].AskSize <= 0 {
					return
				}
				usdt[i] = q[i].Ask * q[i].AskSize
			} else {
				if q[i].Bid <= 0 || q[i].BidSize <= 0 {
					return
				}
				usdt[i] = q[i].Bid * q[i].BidSize
			}
		case 1:
			if leg.IsBuy {
				if q[i].Ask <= 0 || q[i].AskSize <= 0 || q[2].Bid <= 0 {
					return
				}
				usdt[i] = q[i].Ask * q[i].AskSize * q[2].Bid
			} else {
				if q[i].Bid <= 0 || q[i].BidSize <= 0 || q[2].Bid <= 0 {
					return
				}
				usdt[i] = q[i].BidSize * q[2].Bid
			}
		case 2:
			if q[i].Bid <= 0 || q[i].BidSize <= 0 {
				return
			}
			usdt[i] = q[i].Bid * q[i].BidSize
		}
	}

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

	// расчет прибыли с учетом fee
	amount := maxUSDT
	for i, leg := range tri.Legs {
		if leg.IsBuy {
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
		"value": "InvalidIndirection",
		"target": {
			"$mid": 1,
			"path": "/golang.org/x/tools/internal/typesinternal",
			"scheme": "https",
			"authority": "pkg.go.dev",
			"fragment": "InvalidIndirection"
		}
	},
	"severity": 8,
	"message": "invalid operation: cannot indirect q (variable of struct type queue.Quote)",
	"source": "compiler",
	"startLineNumber": 95,
	"startColumn": 21,
	"endLineNumber": 95,
	"endColumn": 22,
	"modelVersionId": 4,
	"origin": "extHost1"
}]




