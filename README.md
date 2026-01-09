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
	"crypt_proto/internal/queue"
	"log"
	"strings"
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
	index     map[string][]int // symbol -> индексы треугольников
}

// NewCalculator создаёт калькулятор с индексом
func NewCalculator(mem *queue.MemoryStore, triangles []Triangle) *Calculator {
	c := &Calculator{
		mem:       mem,
		triangles: triangles,
		index:     make(map[string][]int),
	}

	// формируем индекс: какая пара участвует в каких треугольниках
	for i, tri := range triangles {
		for _, leg := range []string{tri.Leg1, tri.Leg2, tri.Leg3} {
			sym := legSymbol(leg)
			if sym != "" {
				c.index[sym] = append(c.index[sym], i)
			}
		}
	}
	return c
}

// OnMarketData вызываем, когда приходит обновление котировки
func (c *Calculator) OnMarketData(symbol string) {
	triIndexes, ok := c.index[symbol]
	if !ok {
		return // эта пара не участвует в треугольниках
	}

	for _, i := range triIndexes {
		tri := c.triangles[i]

		s1, s2, s3 := legSymbol(tri.Leg1), legSymbol(tri.Leg2), legSymbol(tri.Leg3)

		q1, ok1 := c.mem.Get("KuCoin", s1)
		q2, ok2 := c.mem.Get("KuCoin", s2)
		q3, ok3 := c.mem.Get("KuCoin", s3)
		if !ok1 || !ok2 || !ok3 {
			continue
		}

		amount := 1.0 // стартуем с 1 A

		// LEG 1
		if strings.HasPrefix(tri.Leg1, "BUY") {
			if q1.Ask <= 0 || q1.AskSize <= 0 {
				continue
			}
			amount = amount / q1.Ask
			if amount > q1.AskSize {
				amount = q1.AskSize
			}
			amount *= (1 - fee)
		} else {
			if q1.Bid <= 0 || q1.BidSize <= 0 {
				continue
			}
			if amount > q1.BidSize {
				amount = q1.BidSize
			}
			amount = amount * q1.Bid * (1 - fee)
		}

		// LEG 2
		if strings.HasPrefix(tri.Leg2, "BUY") {
			if q2.Ask <= 0 || q2.AskSize <= 0 {
				continue
			}
			amount = amount / q2.Ask
			if amount > q2.AskSize {
				amount = q2.AskSize
			}
			amount *= (1 - fee)
		} else {
			if q2.Bid <= 0 || q2.BidSize <= 0 {
				continue
			}
			if amount > q2.BidSize {
				amount = q2.BidSize
			}
			amount = amount * q2.Bid * (1 - fee)
		}

		// LEG 3
		if strings.HasPrefix(tri.Leg3, "BUY") {
			if q3.Ask <= 0 || q3.AskSize <= 0 {
				continue
			}
			amount = amount / q3.Ask
			if amount > q3.AskSize {
				amount = q3.AskSize
			}
			amount *= (1 - fee)
		} else {
			if q3.Bid <= 0 || q3.BidSize <= 0 {
				continue
			}
			if amount > q3.BidSize {
				amount = q3.BidSize
			}
			amount = amount * q3.Bid * (1 - fee)
		}

		profit := amount - 1.0
		if profit > 0 {
			log.Printf("[ARB] %s → %s → %s | profit=%.4f%% | volumes: [%.2f / %.2f / %.2f]",
				tri.A, tri.B, tri.C,
				profit*100,
				q1.BidSize, q2.BidSize, q3.BidSize,
			)
		}
	}
}

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

	"crypt_proto/internal/calculator"
	"crypt_proto/internal/collector"
	"crypt_proto/internal/queue"
)

func main() {
	// --- память для котировок ---
	mem := queue.NewMemoryStore()

	// --- читаем треугольники из CSV ---
	triangles, err := calculator.ParseTrianglesFromCSV("triangles.csv")
	if err != nil {
		log.Fatal(err)
	}

	// --- создаём калькулятор с индексом по символам ---
	calc := calculator.NewCalculator(mem, triangles)

	// --- создаём коллектор ---
	kc, err := collector.NewKuCoinCollectorFromCSV("triangles.csv")
	if err != nil {
		log.Fatal(err)
	}

	// --- канал для получения MarketData ---
	out := make(chan *collector.MarketData, 1000)

	// --- старт коллектора ---
	err = kc.Start(out)
	if err != nil {
		log.Fatal(err)
	}
	defer kc.Stop()

	log.Println("Collector started, waiting for market data...")

	// --- основной цикл обработки данных ---
	for md := range out {
		// обновляем память
		mem.Set(md.Exchange, md.Symbol, md)

		// считаем только треугольники, где участвует эта пара
		calc.OnMarketData(md.Symbol)
	}
}



[{
	"resource": "/home/gaz358/myprog/crypt_proto/cmd/arb/main.go",
	"owner": "_generated_diagnostic_collection_name_#1",
	"code": {
		"value": "UndeclaredImportedName",
		"target": {
			"$mid": 1,
			"path": "/golang.org/x/tools/internal/typesinternal",
			"scheme": "https",
			"authority": "pkg.go.dev",
			"fragment": "UndeclaredImportedName"
		}
	},
	"severity": 8,
	"message": "undefined: collector.MarketData",
	"source": "compiler",
	"startLineNumber": 31,
	"startColumn": 30,
	"endLineNumber": 31,
	"endColumn": 40,
	"origin": "extHost1"
}]


[{
	"resource": "/home/gaz358/myprog/crypt_proto/cmd/arb/main.go",
	"owner": "_generated_diagnostic_collection_name_#1",
	"code": {
		"value": "MissingFieldOrMethod",
		"target": {
			"$mid": 1,
			"path": "/golang.org/x/tools/internal/typesinternal",
			"scheme": "https",
			"authority": "pkg.go.dev",
			"fragment": "MissingFieldOrMethod"
		}
	},
	"severity": 8,
	"message": "mem.Set undefined (type *queue.MemoryStore has no field or method Set)",
	"source": "compiler",
	"startLineNumber": 45,
	"startColumn": 7,
	"endLineNumber": 45,
	"endColumn": 10,
	"origin": "extHost1"
}]






