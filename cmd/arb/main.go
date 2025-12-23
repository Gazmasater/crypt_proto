package main

import (
	"fmt"
	"os"
	"strings"

	"crypt_proto/internal/collector"
	"crypt_proto/pkg/models"

	"github.com/joho/godotenv"
)

func main() {
	// Загружаем .env
	_ = godotenv.Load(".env")

	exchange := strings.ToLower(os.Getenv("EXCHANGE"))
	if exchange == "" {
		exchange = "mexc"
	}
	fmt.Println("EXCHANGE:", exchange)

	// Канал для получения рыночных данных
	marketDataCh := make(chan models.MarketData, 1000)

	var c collector.Collector
	switch exchange {
	case "mexc":
		c = collector.NewMEXCCollector([]string{"BTCUSDT", "ETHUSDT"})
	case "kucoin":
		c = collector.NewKuCoinCollector([]string{"BTC-USDT", "ETH-USDT"})
	case "okx":
		c = collector.NewOKXCollector([]string{"BTC-USDT", "ETH-USDT"})
	default:
		panic("unknown exchange: " + exchange)
	}

	fmt.Println("Starting collector:", c.Name())
	if err := c.Start(marketDataCh); err != nil {
		panic(err)
	}
	defer c.Stop()

	// Consumer
	go func() {
		for data := range marketDataCh {
			fmt.Printf("[%s] %s bid=%.8f ask=%.8f\n",
				data.Exchange, data.Symbol, data.Bid, data.Ask)
		}
	}()

	select {}
}
