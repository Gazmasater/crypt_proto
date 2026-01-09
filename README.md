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

type Triangle struct {
	A, B, C          string
	Leg1, Leg2, Leg3 string
}

type Calculator struct {
	mem       *queue.MemoryStore
	triangles []Triangle
	index     map[string][]int // символ → индексы треугольников
}

// NewCalculator создаёт калькулятор и индекс треугольников
func NewCalculator(mem *queue.MemoryStore, triangles []Triangle) *Calculator {
	c := &Calculator{
		mem:       mem,
		triangles: triangles,
		index:     make(map[string][]int),
	}

	// строим индекс: какая пара участвует в каких треугольниках
	for i, tri := range triangles {
		for _, leg := range []string{tri.Leg1, tri.Leg2, tri.Leg3} {
			s := legSymbol(leg)
			c.index[s] = append(c.index[s], i)
		}
	}

	return c
}

// OnMarketData пересчитывает только треугольники с этой парой
func (c *Calculator) OnMarketData(symbol string) {
	triIndexes, ok := c.index[symbol]
	if !ok {
		return
	}

	for _, i := range triIndexes {
		tri := c.triangles[i]

		s1 := legSymbol(tri.Leg1)
		s2 := legSymbol(tri.Leg2)
		s3 := legSymbol(tri.Leg3)

		q1, ok1 := c.mem.Get("KuCoin", s1)
		q2, ok2 := c.mem.Get("KuCoin", s2)
		q3, ok3 := c.mem.Get("KuCoin", s3)

		if !ok1 || !ok2 || !ok3 {
			continue
		}

		amount := 1.0

		legs := []struct {
			leg  string
			quote queue.Quote
		}{
			{tri.Leg1, q1},
			{tri.Leg2, q2},
			{tri.Leg3, q3},
		}

		for _, l := range legs {
			if strings.HasPrefix(l.leg, "BUY") {
				if l.quote.Ask <= 0 || l.quote.AskSize <= 0 {
					amount = 0
					break
				}
				maxBuy := l.quote.AskSize
				amount = amount / l.quote.Ask
				if amount > maxBuy {
					amount = maxBuy
				}
				amount *= (1 - fee)
			} else {
				if l.quote.Bid <= 0 || l.quote.BidSize <= 0 {
					amount = 0
					break
				}
				maxSell := l.quote.BidSize
				if amount > maxSell {
					amount = maxSell
				}
				amount = amount * l.quote.Bid
				amount *= (1 - fee)
			}
		}

		profit := amount - 1.0
		if profit > 0 {
			log.Printf("[ARB] %s → %s → %s | profit=%.4f%% | volumes: [%.2f / %.2f / %.2f]",
				tri.A, tri.B, tri.C, profit*100, q1.BidSize, q2.BidSize, q3.BidSize)
		}
	}
}

// legSymbol: "BUY COTI/USDT" -> "COTI/USDT"
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
	"crypt_proto/pkg/models"

	"github.com/joho/godotenv"
)

func main() {
	// --- env ---
	_ = godotenv.Load(".env")

	// --- память ---
	mem := queue.NewMemoryStore()

	// --- создаём коллектор из CSV ---
	kc, err := collector.NewKuCoinCollectorFromCSV("triangles.csv")
	if err != nil {
		log.Fatal(err)
	}

	// --- парсим треугольники ---
	triangles := kc.Triangles()

	// --- калькулятор ---
	calc := calculator.NewCalculator(mem, triangles)

	// --- канал данных с коллектора ---
	out := make(chan *models.MarketData, 1000)

	// --- запуск коллектора ---
	if err := kc.Start(out); err != nil {
		log.Fatal(err)
	}
	defer kc.Stop()

	// --- ловим Ctrl+C ---
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// --- основной цикл ---
	go func() {
		for md := range out {
			// конвертируем MarketData в queue.Quote
			q := queue.Quote{
				Bid:     md.Bid,
				Ask:     md.Ask,
				BidSize: md.BidSize,
				AskSize: md.AskSize,
			}

			// обновляем память
			mem.Put(md.Exchange, md.Symbol, q)

			// пересчёт только треугольников с этой парой
			calc.OnMarketData(md.Symbol)
		}
	}()

	// --- ждём сигнал на выход ---
	<-sigs
	log.Println("Exiting...")
}








