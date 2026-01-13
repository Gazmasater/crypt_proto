package calculator

import (
	"crypt_proto/internal/queue"
	"crypt_proto/pkg/models"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
)

const fee = 0.001 // 0.1%

// Triangle описывает один треугольный арбитраж
type Triangle struct {
	A, B, C          string // имена валют для логов
	Leg1, Leg2, Leg3 string // "BUY COTI/USDT", "SELL COTI/BTC" и т.д.
}

// Calculator считает профит по треугольникам
type Calculator struct {
	mem       *queue.MemoryStore
	triangles []Triangle
	bySymbol  map[string][]Triangle
	fileLog   *log.Logger
}

// NewCalculator создаёт калькулятор
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

// Run запускает цикл расчёта
func (c *Calculator) Run(in <-chan *models.MarketData) {
	for md := range in {

		// сохраняем маркет
		c.mem.Push(md)

		tris := c.bySymbol[md.Symbol]
		if len(tris) == 0 {
			continue
		}

		for _, tri := range tris {
			c.calcTriangle(&tri)
		}
	}
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

	// ===== 1. USDT LIMITS =====

	var usdtLimits [3]float64
	i := 0

	// LEG 1
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

	// LEG 2
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

	// LEG 3
	if q3.Bid <= 0 || q3.BidSize <= 0 {
		return
	}
	usdtLimits[i] = q3.Bid * q3.BidSize

	// ===== 2. MIN LIMIT =====

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

	// ===== 3. PROGON =====

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

	// ===== 4. LOG =====

	if profitPct > 0.001 && profitUSDT > 0.02 {

		msg := fmt.Sprintf(
			"[ARB] %s → %s → %s | %.4f%% | volume=%.2f USDT | profit=%.4f USDT",
			tri.A,
			tri.B,
			tri.C,
			profitPct*100,
			maxUSDT,
			profitUSDT,
		)

		log.Println(msg)
		c.fileLog.Println(msg)
	}
}

func ParseTrianglesFromCSV(path string) ([]Triangle, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	rows, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	var res []Triangle
	for _, row := range rows[1:] {
		if len(row) < 6 {
			continue
		}

		res = append(res, Triangle{
			A:    strings.TrimSpace(row[0]),
			B:    strings.TrimSpace(row[1]),
			C:    strings.TrimSpace(row[2]),
			Leg1: strings.TrimSpace(row[3]),
			Leg2: strings.TrimSpace(row[4]),
			Leg3: strings.TrimSpace(row[5]),
		})
	}

	return res, nil
}

//func legSymbol(leg string) string {
//	// "BUY COTI/USDT" -> "COTI/USDT"
//	parts := strings.Fields(leg)
//	if len(parts) != 2 {
//		return ""
//	}
//	return strings.ToUpper(parts[1])
//}

func legSymbol(leg string) string {
	parts := strings.Fields(strings.ToUpper(strings.TrimSpace(leg)))
	if len(parts) < 2 {
		return ""
	}
	p := strings.Split(parts[1], "/")
	if len(p) != 2 {
		return ""
	}
	return p[0] + "-" + p[1] // формат BASE-QUOTE
}
