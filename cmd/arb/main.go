package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"crypt_proto/internal/collector"
	"crypt_proto/pkg/models"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load(".env")
	exchange := os.Getenv("EXCHANGE")
	if exchange == "" {
		log.Fatal("EXCHANGE не задан. Используй MEXC, OKX или KuCoin")
	}
	log.Println("EXCHANGE:", exchange)

	marketDataCh := make(chan models.MarketData, 1000)

	var c collector.Collector

	switch exchange {
	case "MEXC":
		c = collector.NewMEXCCollector([]string{
			"BTCUSDT",
			"ETHUSDT",
			"ETHBTC",
		})
	case "OKX":
		c = collector.NewOKXCollector([]string{
			"BTC-USDT",
			"ETH-USDT",
			"ETH-BTC",
		})
	case "KuCoin":
		c = collector.NewKuCoinCollector([]string{
			"BTC-USDT",
			"ETH-USDT",
			"ETH-BTC",
		})
	default:
		log.Fatal("Неподдерживаемый EXCHANGE:", exchange)
	}

	if err := c.Start(marketDataCh); err != nil {
		log.Fatal("Start error:", err)
	}

	// читаем данные в фоне
	go func() {
		for data := range marketDataCh {
			log.Printf("[%s] %s bid=%.6f ask=%.6f\n",
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
