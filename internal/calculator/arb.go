package calculator

import (
	"crypt_proto/internal/queue"
	"encoding/csv"
	"log"
	"os"
	"strings"
	"time"
)

// Triangle описывает один треугольный арбитраж
type Triangle struct {
	A, B, C          string // имена валют для логов
	Leg1, Leg2, Leg3 string // "BUY COTI/USDT", "SELL COTI/BTC" и т.д.
}

// Calculator считает профит по треугольникам
type Calculator struct {
	mem       *queue.MemoryStore
	triangles []Triangle
}

// NewCalculator создаёт калькулятор
func NewCalculator(mem *queue.MemoryStore, triangles []Triangle) *Calculator {
	return &Calculator{
		mem:       mem,
		triangles: triangles,
	}
}

// Run запускает цикл расчёта
func (c *Calculator) Run() {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {

		for _, tri := range c.triangles {

			// ключи MemoryStore
			k1 := "KuCoin|" + legSymbol(tri.Leg1)
			k2 := "KuCoin|" + legSymbol(tri.Leg2)
			k3 := "KuCoin|" + legSymbol(tri.Leg3)

			q1, ok1 := c.mem.Get("Kucoin", k1)
			q2, ok2 := c.mem.Get("Kucoin", k2)
			q3, ok3 := c.mem.Get("Kucoin", k3)
			if !ok1 || !ok2 || !ok3 {
				continue
			}

			amount := 1.0

			// Leg1
			if strings.HasPrefix(tri.Leg1, "BUY") {
				if q1.Ask == 0 {
					continue
				}
				amount /= q1.Ask
			} else {
				amount *= q1.Bid
			}

			// Leg2
			if strings.HasPrefix(tri.Leg2, "BUY") {
				if q2.Ask == 0 {
					continue
				}
				amount /= q2.Ask
			} else {
				amount *= q2.Bid
			}

			// Leg3
			if strings.HasPrefix(tri.Leg3, "BUY") {
				if q3.Ask == 0 {
					continue
				}
				amount /= q3.Ask
			} else {
				amount *= q3.Bid
			}

			profit := amount - 1.0

			// РЕАЛЬНЫЙ порог
			if profit > 0.001 {
				log.Printf(
					"[ARB] %s → %s → %s | profit=%.4f%%",
					tri.A, tri.B, tri.C,
					profit*100,
				)
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
