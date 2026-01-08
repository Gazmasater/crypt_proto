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

	"crypt_proto/pkg/models"
	"crypt_proto/pkg/store"
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
	mem       *store.MemoryStore
}

// NewCalculator создаёт калькулятор
func NewCalculator(triangles []Triangle, mem *store.MemoryStore) *Calculator {
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





package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"net/http"
	_ "net/http/pprof"

	"crypt_proto/internal/collector"
	"crypt_proto/pkg/calculator"
	"crypt_proto/pkg/models"
	"crypt_proto/pkg/store"
)

func main() {
	go func() {
		log.Println("pprof on http://localhost:6060/debug/pprof/")
		_ = http.ListenAndServe("localhost:6060", nil)
	}()

	// ------------------- Канал для данных -------------------
	out := make(chan *models.MarketData, 100_000)

	// ------------------- In-Memory Store -------------------
	mem := store.NewMemoryStore()
	go mem.Run(out)

	// ------------------- Коллектор -------------------
	kc, err := collector.NewKuCoinCollectorFromCSV("../exchange/data/kucoin_triangles_usdt.csv")
	if err != nil {
		log.Fatal(err)
	}

	if err := kc.Start(out); err != nil {
		log.Fatal(err)
	}
	log.Println("[Main] KuCoinCollector started")

	// ------------------- Загружаем треугольники из CSV -------------------
	triangles, err := collector.LoadTrianglesFromCSV("../exchange/data/kucoin_triangles_usdt.csv")
	if err != nil {
		log.Fatal(err)
	}

	// ------------------- Запуск калькулятора -------------------
	calc := calculator.NewCalculator(triangles, mem)
	go calc.Run()

	// ------------------- Вывод snapshot для проверки -------------------
	go func() {
		for {
			snap := mem.Snapshot()
			log.Printf("[Store] quotes=%d", len(snap))
			time.Sleep(5 * time.Second)
		}
	}()

	// ------------------- Завершение при SIGINT / SIGTERM -------------------
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	kc.Stop()
	close(out)
	log.Println("[Main] Exited.")
}


