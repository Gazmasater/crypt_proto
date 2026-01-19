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

const fee = 0.001

type LegIndex struct {
	Key   string
	IsBuy bool
}

type Triangle struct {
	A, B, C string
	Legs    [3]LegIndex
}

type Calculator struct {
	mem       *queue.MemoryStore
	triangles []*Triangle
	fileLog   *log.Logger
}

// NewCalculator полностью lock-free
func NewCalculator(mem *queue.MemoryStore, triangles []*Triangle) *Calculator {
	f, err := os.OpenFile("arb_opportunities.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("failed to open log: %v", err)
	}

	return &Calculator{
		mem:       mem,
		triangles: triangles,
		fileLog:   log.New(f, "", log.LstdFlags),
	}
}

// Run читает входящие данные и пересчитывает треугольники без блокировок
func (c *Calculator) Run(in <-chan *models.MarketData) {
	for md := range in {
		c.mem.Push(md)
		// lock-free: проверяем каждый треугольник
		for _, tri := range c.triangles {
			c.calcTriangle(tri)
		}
	}
}

func (c *Calculator) calcTriangle(tri *Triangle) {
	q := [3]queue.Quote{}
	for i, leg := range tri.Legs {
		quote, ok := c.mem.Get("KuCoin", strings.Split(leg.Key, "|")[1])
		if !ok {
			return
		}
		q[i] = quote
	}

	usdtLimits := [3]float64{}
	// LEG1
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
	// LEG2
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
	// LEG3
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
		amount = amount / q[0].Ask * (1 - fee)
	} else {
		amount = amount * q[0].Bid * (1 - fee)
	}
	if tri.Legs[1].IsBuy {
		amount = amount / q[1].Ask * (1 - fee)
	} else {
		amount = amount * q[1].Bid * (1 - fee)
	}
	if tri.Legs[2].IsBuy {
		amount = amount / q[2].Ask * (1 - fee)
	} else {
		amount = amount * q[2].Bid * (1 - fee)
	}

	profitUSDT := amount - maxUSDT
	profitPct := profitUSDT / maxUSDT

	if profitPct > 0.001 && profitUSDT > 0.02 {
		msg := fmt.Sprintf("[ARB] %s→%s→%s | %.4f%% | volume=%.2f USDT | profit=%.4f USDT",
			tri.A, tri.B, tri.C, profitPct*100, maxUSDT, profitUSDT)
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

		tri := &Triangle{
			A: strings.TrimSpace(row[0]),
			B: strings.TrimSpace(row[1]),
			C: strings.TrimSpace(row[2]),
		}

		legs := []string{row[3], row[4], row[5]}
		for i, leg := range legs {
			leg = strings.ToUpper(strings.TrimSpace(leg))
			parts := strings.Fields(leg)
			if len(parts) != 2 {
				continue
			}
			isBuy := parts[0] == "BUY"
			symbolParts := strings.Split(parts[1], "/")
			if len(symbolParts) != 2 {
				continue
			}
			key := "KuCoin|" + symbolParts[0] + "-" + symbolParts[1]
			tri.Legs[i] = LegIndex{Key: key, IsBuy: isBuy}
		}

		res = append(res, tri)
	}
	return res, nil
}
