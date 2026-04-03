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

const feeM = 0.999

var triangleLegColumns = [3]int{3, 7, 11}

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
			if leg.Symbol == "" {
				continue
			}
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
	var q [3]queue.Quote

	for i, leg := range tri.Legs {
		if leg.Symbol == "" {
			log.Printf("skip triangle %s->%s->%s: empty symbol for leg=%d", tri.A, tri.B, tri.C, i)
			return
		}

		quote, ok := c.mem.Get("KuCoin", leg.Symbol)
		if !ok {
			log.Printf("skip triangle %s->%s->%s: no quote for leg=%d symbol=%s", tri.A, tri.B, tri.C, i, leg.Symbol)
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

	msg := fmt.Sprintf(
		"[ARB] %s→%s→%s | %.4f%% | volume=%.2f USDT | profit=%.4f USDT",
		tri.A, tri.B, tri.C,
		profitPct*100, maxUSDT, profitUSDT,
	)
	log.Println(msg)
	c.fileLog.Println(msg)
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

	for rowIdx, row := range rows[1:] {
		if len(row) < 15 {
			continue
		}

		t := &Triangle{
			A: strings.TrimSpace(row[0]),
			B: strings.TrimSpace(row[1]),
			C: strings.TrimSpace(row[2]),
		}

		for i, col := range triangleLegColumns {
			leg, err := parseTriangleLeg(row[col])
			if err != nil {
				return nil, fmt.Errorf("row %d leg %d: %w", rowIdx+2, i+1, err)
			}
			t.Legs[i] = leg
		}

		res = append(res, t)
	}

	return res, nil
}

func parseTriangleLeg(raw string) (LegIndex, error) {
	leg := strings.ToUpper(strings.TrimSpace(raw))
	parts := strings.Fields(leg)
	if len(parts) != 2 {
		return LegIndex{}, fmt.Errorf("bad leg format: %q", raw)
	}

	isBuy := parts[0] == "BUY"
	if parts[0] != "BUY" && parts[0] != "SELL" {
		return LegIndex{}, fmt.Errorf("bad leg side: %q", raw)
	}

	pair := strings.Split(parts[1], "/")
	if len(pair) != 2 {
		return LegIndex{}, fmt.Errorf("bad pair format: %q", raw)
	}

	symbol := pair[0] + "-" + pair[1]
	return LegIndex{
		Key:    "KuCoin|" + symbol,
		Symbol: symbol,
		IsBuy:  isBuy,
	}, nil
}
