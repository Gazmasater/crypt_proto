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

	"net/http"
	_ "net/http/pprof"

	"crypt_proto/internal/collector"
	"crypt_proto/pkg/models"
)

func main() {

	go func() {
		log.Println("pprof on http://localhost:6060/debug/pprof/")
		if err := http.ListenAndServe("localhost:6060", nil); err != nil {
			log.Printf("pprof server error: %v", err)
		}
	}()

	// ------------------- Канал для данных -------------------
	out := make(chan *models.MarketData, 100)

	// ------------------- Создание коллектора -------------------
	kc, err := collector.NewKuCoinCollectorFromCSV("../exchange/data/kucoin_triangles_usdt.csv")
	if err != nil {
		log.Fatal("Failed to create KuCoinCollector:", err)
	}

	// ------------------- Запуск коллектора -------------------
	if err := kc.Start(out); err != nil {
		log.Fatal("Failed to start KuCoinCollector:", err)
	}

	log.Println("[Main] KuCoinCollector started. Listening for data...")

	// ------------------- Обработка данных -------------------
	//go func() {
	//	for data := range out {
	//		log.Printf("[MarketData] %s %s bid=%.6f bidsize=%.6f ask=%.6f  asksize=%.6f",
	//			data.Exchange, data.Symbol, data.Bid, data.BidSize, data.Ask, data.AskSize)
	//	}
	//}()

	// ------------------- Завершение при SIGINT / SIGTERM -------------------
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("[Main] Stopping KuCoinCollector...")
	if err := kc.Stop(); err != nil {
		log.Println("Error stopping collector:", err)
	}
	close(out)
	log.Println("[Main] Exited.")
}







