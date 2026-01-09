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
	_ = godotenv.Load(".env")

	// --- память для котировок ---
	mem := queue.NewMemoryStore()

	// --- читаем треугольники из CSV ---
	triangles, err := calculator.ParseTrianglesFromCSV("triangles.csv")
	if err != nil {
		log.Fatal(err)
	}

	// --- калькулятор ---
	calc := calculator.NewCalculator(mem, triangles)
	go calc.RunAsync() // запускаем асинхронный расчёт

	// --- коллектор ---
	kc, err := collector.NewKuCoinCollectorFromCSV("triangles.csv")
	if err != nil {
		log.Fatal(err)
	}

	out := make(chan *models.MarketData, 1000)
	go func() {
		for md := range out {
			// --- конвертация в Quote и запись в память ---
			q := queue.Quote{
				Bid:     md.Bid,
				Ask:     md.Ask,
				BidSize: md.BidSize,
				AskSize: md.AskSize,
			}
			mem.Put(md.Exchange, md.Symbol, q)

			// --- считаем только треугольники с этой парой ---
			calc.OnMarketData(md.Symbol)
		}
	}()

	if err := kc.Start(out); err != nil {
		log.Fatal(err)
	}

	// --- graceful shutdown ---
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("stopping...")
	_ = kc.Stop()
	close(out)
}



[{
	"resource": "/home/gaz358/myprog/crypt_proto/cmd/arb/main.go",
	"owner": "_generated_diagnostic_collection_name_#0",
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
	"message": "calc.RunAsync undefined (type *calculator.Calculator has no field or method RunAsync)",
	"source": "compiler",
	"startLineNumber": 31,
	"startColumn": 10,
	"endLineNumber": 31,
	"endColumn": 18,
	"origin": "extHost1"
}]





