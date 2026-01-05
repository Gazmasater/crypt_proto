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
	go func() {
		for data := range out {
			log.Printf("[MarketData] %s %s bid=%.6f ask=%.6f",
				data.Exchange, data.Symbol, data.Bid, data.Ask)
		}
	}()

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
