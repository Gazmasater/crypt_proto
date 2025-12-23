package main

import (
	"crypt_proto/internal/collector"
	"crypt_proto/pkg/models"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	exchange := strings.ToLower(os.Getenv("EXCHANGE"))
	if exchange == "" {
		exchange = "mexc"
	}
	log.Println("EXCHANGE:", exchange)

	marketDataCh := make(chan models.MarketData, 1000)

	var c collector.Collector

	switch exchange {
	case "mexc":
		c = collector.NewMEXCCollector([]string{
			"BTCUSDT",
			"ETHUSDT",
			"ETHBTC",
		})
	case "okx":
		c = collector.NewOKXCollector([]string{
			"BTC-USDT",
			"ETH-USDT",
			"ETH-BTC",
		})
	default:
		log.Fatal("unknown exchange")
	}

	if err := c.Start(marketDataCh); err != nil {
		log.Fatal(err)
	}

	// читаем данные в фоне
	go func() {
		for data := range marketDataCh {
			log.Printf("[%s] %s bid=%.4f ask=%.4f\n",
				data.Exchange, data.Symbol, data.Bid, data.Ask)
		}
	}()

	// корректное завершение по Ctrl+C
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	log.Println("Stopping collector...")
	if err := c.Stop(); err != nil {
		log.Println("Stop error:", err)
	}
}
