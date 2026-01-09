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


BOOK_INTERVAL=100ms
SYMBOLS_FILE=triangles_markets.csv
DEBUG=false


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
	"log"
	"strings"

	"crypt_proto/internal/queue"
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

	amount := 1.0 // стартуем с 1 единицы валюты A

	legs := []struct {
		leg string
		q   queue.Quote
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
		} else { // SELL
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

// legSymbol извлекает символ из Leg, например "BUY COTI/USDT" -> "COTI/USDT"
func legSymbol(leg string) string {
	parts := strings.Fields(leg)
	if len(parts) != 2 {
		return ""
	}
	return strings.ToUpper(parts[1])
}



package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"crypt_proto/internal/calculator"
	"crypt_proto/internal/collector"
	"crypt_proto/internal/queue"
)

func main() {
	// -------------------- Создаем память --------------------
	mem := queue.NewMemoryStore()
	go mem.Run() // запускаем основной цикл MemoryStore

	// -------------------- Загружаем треугольники --------------------
	triangles, err := calculator.ParseTrianglesFromCSV("triangles.csv")
	if err != nil {
		log.Fatalf("failed to load triangles: %v", err)
	}

	// -------------------- Создаем калькулятор --------------------
	calc := calculator.NewCalculator(mem, triangles)

	// -------------------- Создаем коллектор KuCoin --------------------
	ku, err := collector.NewKuCoinCollectorFromCSV("triangles.csv")
	if err != nil {
		log.Fatalf("failed to create KuCoin collector: %v", err)
	}

	// -------------------- Передаем функцию обратного вызова на апдейт котировки --------------------
	ku.OnUpdate = func(symbol string) {
		calc.OnUpdate(symbol)
	}

	// -------------------- Запуск коллектора --------------------
	if err := ku.Start(mem); err != nil {
		log.Fatalf("failed to start KuCoin collector: %v", err)
	}

	log.Println("Calculator and KuCoin collector started")

	// -------------------- Ждем завершения --------------------
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down...")

	if err := ku.Stop(); err != nil {
		log.Printf("Error stopping collector: %v", err)
	}
}





