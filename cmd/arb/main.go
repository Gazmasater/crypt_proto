package main

import (
	"crypt_proto/internal/collector"
	"crypt_proto/pkg/models"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// канал для рыночных данных
	marketDataCh := make(chan models.MarketData, 1000)

	// инициализация коллектора
	okxCollector := collector.NewOKXCollector()

	// старт коллектора
	if err := okxCollector.Start(marketDataCh); err != nil {
		log.Fatal("failed to start OKX collector:", err)
	}

	log.Println("OKX collector started")

	// consumer (пока просто логируем)
	go func() {
		for md := range marketDataCh {
			log.Printf(
				"[MARKET] %s %s bid=%.4f ask=%.4f",
				md.Exchange,
				md.Symbol,
				md.Bid,
				md.Ask,
			)
		}
	}()

	// graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	<-sigCh
	log.Println("shutdown signal received")

	okxCollector.Stop()

	log.Println("collector stopped, exit")
}
