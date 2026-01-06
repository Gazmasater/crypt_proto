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

	"crypt_proto/internal/collector"
	"crypt_proto/pkg/models"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	// ===== OUTPUT CHANNEL =====
	out := make(chan *models.MarketData, 10_000)

	// ===== CSV WITH TRIANGLES =====
	csvPath := "./triangles.csv" // путь к твоему CSV
	symbols, err := collector.LoadSymbolsFromCSV(csvPath)
	if err != nil {
		log.Fatal("csv error:", err)
	}

	log.Printf("[Main] loaded %d unique symbols", len(symbols))

	// ===== KUCOIN COLLECTOR (POOL) =====
	kucoin := collector.NewKuCoinCollector(symbols)

	if err := kucoin.Start(out); err != nil {
		log.Fatal("kucoin start error:", err)
	}

	log.Println("[Main] KuCoinCollector started")

	// ===== DATA CONSUMER =====
	go func() {
		for md := range out {
			// временно просто логируем
			log.Printf(
				"[MD] %s %s bid=%.8f ask=%.8f",
				md.Exchange,
				md.Symbol,
				md.Bid,
				md.Ask,
			)
		}
	}()

	// ===== GRACEFUL SHUTDOWN =====
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	<-sig
	log.Println("[Main] shutdown")

	_ = kucoin.Stop()
}



