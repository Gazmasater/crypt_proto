package calculator

import (
	"crypt_proto/internal/queue"
	"encoding/csv"
	"fmt"
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

			amount := 1.0 // стартуем с 1 A

			// ---------- LEG 1 ----------
			if strings.HasPrefix(tri.Leg1, "BUY") {
				if q1.Ask <= 0 || q1.AskSize <= 0 {
					continue
				}

				maxBuy := q1.AskSize     // сколько B можно купить
				amount = amount / q1.Ask // A -> B
				if amount > maxBuy {
					amount = maxBuy
				}
				amount *= (1 - fee)

			} else {
				if q1.Bid <= 0 || q1.BidSize <= 0 {
					continue
				}

				maxSell := q1.BidSize // сколько A можно продать
				if amount > maxSell {
					amount = maxSell
				}
				amount = amount * q1.Bid // A -> B
				amount *= (1 - fee)
			}

			// ---------- LEG 2 ----------
			if strings.HasPrefix(tri.Leg2, "BUY") {
				if q2.Ask <= 0 || q2.AskSize <= 0 {
					continue
				}

				maxBuy := q2.AskSize
				amount = amount / q2.Ask
				if amount > maxBuy {
					amount = maxBuy
				}
				amount *= (1 - fee)

			} else {
				if q2.Bid <= 0 || q2.BidSize <= 0 {
					continue
				}

				maxSell := q2.BidSize
				if amount > maxSell {
					amount = maxSell
				}
				amount = amount * q2.Bid
				amount *= (1 - fee)
			}

			// ---------- LEG 3 ----------
			if strings.HasPrefix(tri.Leg3, "BUY") {
				if q3.Ask <= 0 || q3.AskSize <= 0 {
					continue
				}

				maxBuy := q3.AskSize
				amount = amount / q3.Ask
				if amount > maxBuy {
					amount = maxBuy
				}
				amount *= (1 - fee)

			} else {
				if q3.Bid <= 0 || q3.BidSize <= 0 {
					continue
				}

				maxSell := q3.BidSize
				if amount > maxSell {
					amount = maxSell
				}
				amount = amount * q3.Bid
				amount *= (1 - fee)
			}

			profit := amount - 1.0

			if profit > 0.001 {
				msg := fmt.Sprintf(
					"[ARB] %s → %s → %s | profit=%.4f%% | volumes: [%.4f / %.4f / %.4f]",
					tri.A, tri.B, tri.C,
					profit*100,
					q1.BidSize, q2.BidSize, q3.BidSize,
				)

				// консоль
				log.Println(msg)

				// файл
				c.fileLog.Println(msg)
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

func legSymbol(leg string) string {
	// "BUY COTI/USDT" -> "COTI/USDT"
	parts := strings.Fields(leg)
	if len(parts) != 2 {
		return ""
	}
	return strings.ToUpper(parts[1])
}
