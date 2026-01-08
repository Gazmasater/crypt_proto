package calculator

import (
	"crypt_proto/internal/queue"
	"log"
)

// Triangle описывает один треугольник для арбитража
type Triangle struct {
	A, B, C string
	Leg1    string
	Leg2    string
	Leg3    string
}

// Calculator считает профит по треугольникам
type Calculator struct {
	triangles []Triangle
	mem       *queue.MemoryStore
}

// NewCalculator создаёт калькулятор
func NewCalculator(triangles []Triangle, mem *queue.MemoryStore) *Calculator {
	return &Calculator{
		triangles: triangles,
		mem:       mem,
	}
}

// Run запускает постоянный цикл вычислений
func (c *Calculator) Run() {
	for {
		snapshot := c.mem.Snapshot()

		for _, tri := range c.triangles {
			// Берём цены с MemoryStore
			leg1, ok1 := snapshot["KuCoin|"+tri.Leg1]
			leg2, ok2 := snapshot["KuCoin|"+tri.Leg2]
			leg3, ok3 := snapshot["KuCoin|"+tri.Leg3]

			if !ok1 || !ok2 || !ok3 {
				continue // ждём, пока будут все цены
			}

			// Простейший расчет "профита" (пример)
			// Начинаем с 1 единицы A
			amount := 1.0

			// Leg1: BUY
			amount /= leg1.Ask // тратим Ask для покупки B
			// Leg2: SELL
			amount *= leg2.Bid // продаём B за C
			// Leg3: SELL
			amount *= leg3.Bid // продаём C за A

			profit := amount - 1.0
			if profit > 0 {
				log.Printf("[Arb] Triangle %s-%s-%s Profit=%.6f", tri.A, tri.B, tri.C, profit)
			}
		}
		// Частота вычислений
		// Можно добавить time.Sleep(100 * time.Millisecond) если много треугольников
	}
}
