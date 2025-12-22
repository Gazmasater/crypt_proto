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
	_ = godotenv.Load(".env")
	exchange := strings.ToLower(os.Getenv("EXCHANGE"))
	if exchange == "" {
		exchange = "mexc"
	}

	fmt.Println("EXCHANGE:", exchange)

	marketDataCh := make(chan models.MarketData, 1000)

	if exchange != "mexc" && exchange != "okx" {
		panic("unknown exchange")
	}

	// два символа
	symbols := []string{"BTCUSDT", "ETHUSDT"}

	var c collector.Collector
	if exchange == "mexc" {
		c = collector.NewMEXCCollector(symbols)
	} else {
		c = collector.NewOKXCollector()
	}

	fmt.Println("Starting collector:", c.Name())
	if err := c.Start(marketDataCh); err != nil {
		panic(err)
	}

	// consumer
	go func() {
		for data := range marketDataCh {
			fmt.Printf("[%s] %s bid=%.4f ask=%.4f\n",
				data.Exchange, data.Symbol, data.Bid, data.Ask)
		}
	}()

	select {}
}
