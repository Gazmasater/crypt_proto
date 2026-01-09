package calculator

import (
	"crypt_proto/internal/queue"
	"crypt_proto/pkg/models"
	"encoding/csv"
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
}

// NewCalculator создаёт калькулятор
func NewCalculator(mem *queue.MemoryStore, triangles []Triangle) *Calculator {
	return &Calculator{
		mem:       mem,
		triangles: triangles,
	}
}

// OnUpdate вызывается на каждый апдейт котировки
func (c *Calculator) OnUpdate(symbol string) {
	for _, tri := range c.triangles {
		s1 := legSymbol(tri.Leg1)
		s2 := legSymbol(tri.Leg2)
		s3 := legSymbol(tri.Leg3)

		// пересчитываем только если обновился один из символов треугольника
		if symbol != s1 && symbol != s2 && symbol != s3 {
			continue
		}

		c.calculateTriangle(tri, s1, s2, s3)
	}
}

// calculateTriangle рассчитывает прибыль одного треугольника
func (c *Calculator) calculateTriangle(tri Triangle, s1, s2, s3 string) {
	q1, ok1 := c.mem.Get("KuCoin", s1)
	q2, ok2 := c.mem.Get("KuCoin", s2)
	q3, ok3 := c.mem.Get("KuCoin", s3)

	if !ok1 || !ok2 || !ok3 {
		return
	}

	amount := 1.0 // стартуем с 1 A

	legs := []struct {
		leg string
		q   *models.MarketData
	}{
		{tri.Leg1, q1},
		{tri.Leg2, q2},
		{tri.Leg3, q3},
	}

	for _, l := range legs {
		if strings.HasPrefix(l.leg, "BUY") {
			if l.q.Ask <= 0 || l.q.AskSize <= 0 {
				return
			}
			maxBuy := l.q.AskSize
			amount = amount / l.q.Ask
			if amount > maxBuy {
				amount = maxBuy
			}
			amount *= (1 - fee)
		} else {
			if l.q.Bid <= 0 || l.q.BidSize <= 0 {
				return
			}
			maxSell := l.q.BidSize
			if amount > maxSell {
				amount = maxSell
			}
			amount = amount * l.q.Bid
			amount *= (1 - fee)
		}
	}

	profit := amount - 1.0
	if profit > 0 {
		log.Printf(
			"[ARB] %s → %s → %s | profit=%.4f%% | volumes: [%.2f / %.2f / %.2f]",
			tri.A, tri.B, tri.C,
			profit*100,
			q1.BidSize, q2.BidSize, q3.BidSize,
		)
	}
}

// ParseTrianglesFromCSV парсит треугольники из CSV
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

// legSymbol извлекает символ из Leg, например "BUY COTI/USDT" -> "COTI/USDT"
func legSymbol(leg string) string {
	parts := strings.Fields(leg)
	if len(parts) != 2 {
		return ""
	}
	return strings.ToUpper(parts[1])
}
