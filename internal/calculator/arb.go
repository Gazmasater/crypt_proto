package calculator

import (
	"crypt_proto/internal/queue"
	"encoding/csv"
	"log"
	"os"
	"strings"
	"time"
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

	return &Calculator{
		mem:       mem,
		triangles: triangles,
		fileLog:   log.New(f, "", log.LstdFlags),
	}
}

// Run запускает цикл расчёта
func (c *Calculator) Run() {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		for _, tri := range c.triangles {

			s1 := legSymbol(tri.Leg1)
			s2 := legSymbol(tri.Leg2)
			s3 := legSymbol(tri.Leg3)

			q1, ok1 := c.mem.Get("KuCoin", s1)
			q2, ok2 := c.mem.Get("KuCoin", s2)
			q3, ok3 := c.mem.Get("KuCoin", s3)

			if !ok1 || !ok2 || !ok3 {
				continue
			}

			// ===============================
			// 1. СЧИТАЕМ ЛИМИТ КАЖДОЙ НОГИ В USDT
			// ===============================

			usdtLimits := make([]float64, 0, 3)

			// ---------- LEG 1 ----------
			if strings.HasPrefix(tri.Leg1, "BUY") {
				// BUY XXX/USDT
				if q1.Ask <= 0 || q1.AskSize <= 0 {
					continue
				}
				usdtLimits = append(usdtLimits, q1.Ask*q1.AskSize)
			} else {
				// SELL XXX/USDT
				if q1.Bid <= 0 || q1.BidSize <= 0 {
					continue
				}
				usdtLimits = append(usdtLimits, q1.Bid*q1.BidSize)
			}

			// ---------- LEG 2 ----------
			if strings.HasPrefix(tri.Leg2, "BUY") {
				// BUY A/B → лимит в B → переводим в USDT через LEG3
				if q2.Ask <= 0 || q2.AskSize <= 0 || q3.Bid <= 0 {
					continue
				}
				usdtLimits = append(usdtLimits, q2.Ask*q2.AskSize*q3.Bid)
			} else {
				// SELL A/B → A * (B/USDT)
				if q2.Bid <= 0 || q2.BidSize <= 0 || q3.Bid <= 0 {
					continue
				}
				usdtLimits = append(usdtLimits, q2.BidSize*q3.Bid)
			}

			// ---------- LEG 3 ----------
			// всегда SELL ???/USDT
			if q3.Bid <= 0 || q3.BidSize <= 0 {
				continue
			}
			usdtLimits = append(usdtLimits, q3.Bid*q3.BidSize)

			// ===============================
			// 2. МИНИМАЛЬНЫЙ ИСПОЛНИМЫЙ ОБЪЁМ
			// ===============================

			maxUSDT := usdtLimits[0]
			for _, v := range usdtLimits {
				if v < maxUSDT {
					maxUSDT = v
				}
			}

			if maxUSDT <= 0 {
				continue
			}

			// ===============================
			// 3. ПРОГОН АРБИТРАЖА С maxUSDT
			// ===============================

			amount := maxUSDT // стартуем в USDT

			// ---------- LEG 1 ----------
			if strings.HasPrefix(tri.Leg1, "BUY") {
				amount = amount / q1.Ask
				amount *= (1 - fee)
			} else {
				amount = amount * q1.Bid
				amount *= (1 - fee)
			}

			// ---------- LEG 2 ----------
			if strings.HasPrefix(tri.Leg2, "BUY") {
				amount = amount / q2.Ask
				amount *= (1 - fee)
			} else {
				amount = amount * q2.Bid
				amount *= (1 - fee)
			}

			// ---------- LEG 3 ----------
			if strings.HasPrefix(tri.Leg3, "BUY") {
				amount = amount / q3.Ask
				amount *= (1 - fee)
			} else {
				amount = amount * q3.Bid
				amount *= (1 - fee)
			}

			profitUSDT := amount - maxUSDT
			profitPct := profitUSDT / maxUSDT

			if (profitPct > 0.001) && (profitUSDT > 0.02) {
				//msg := fmt.Sprintf(
				//	"[ARB] %s → %s → %s | profit=%.4f%% | volume=%.2f USDT | profit=%.4f USDT",
				//	tri.A, tri.B, tri.C,
				//	profitPct*100,
				//	maxUSDT,
				//	profitUSDT,
				//)

				//log.Println(msg)
				//c.fileLog.Println(msg)
			}
		}
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
