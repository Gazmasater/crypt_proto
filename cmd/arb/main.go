package main

import (
	"crypt_proto/internal/collector"
	"crypt_proto/pkg/models"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	exchange := strings.ToLower(os.Getenv("EXCHANGE"))
	if exchange == "" {
		exchange = "okx"
	}

	marketDataCh := make(chan models.MarketData, 1000)

	var c collector.Collector

	fmt.Println("EXCHANGE!!!!!!!!!", exchange)

	switch exchange {
	case "okx":
		c = collector.NewOKXCollector()

	case "mexc":
		c = collector.NewMEXCCollector("BTCUSDT")
	default:
		log.Fatalf("unknown exchange: %s", exchange)
	}

	log.Printf("Starting collector: %s\n", c.Name())

	if err := c.Start(marketDataCh); err != nil {
		log.Fatal(err)
	}

	go func() {
		for md := range marketDataCh {
			log.Printf(
				"[MARKET] %s %s bid=%.6f ask=%.6f",
				md.Exchange,
				md.Symbol,
				md.Bid,
				md.Ask,
			)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	<-sigCh
	log.Println("shutdown signal")

	c.Stop()
}
