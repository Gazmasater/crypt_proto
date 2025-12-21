package main

import (
	"crypt_proto/internal/collector"
	"crypt_proto/pkg/models"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load(".env")
	exchange := strings.ToLower(os.Getenv("EXCHANGE"))
	if exchange == "" {
		exchange = "okx"
	}

	symbolEnv := os.Getenv("SYMBOLS") // допустим, SYMBOLS=BTCUSDT,ETHUSDT,XRPUSDT
	if symbolEnv == "" {
		symbolEnv = "BTCUSDT"
	}
	symbols := strings.Split(symbolEnv, ",")

	fmt.Println("EXCHANGE!!!!!!!!!", exchange)
	fmt.Println("SYMBOLS!!!!!!!!!", symbols)

	marketDataCh := make(chan models.MarketData, 1000)

	var c collector.Collector

	switch exchange {
	case "okx":
		c = collector.NewOKXCollector() // старый OKX
	case "mexc":
		c = collector.NewMEXCCollector(symbols)
	default:
		panic("unsupported exchange")
	}

	fmt.Println("Starting collector:", c.Name())
	if err := c.Start(marketDataCh); err != nil {
		panic(err)
	}

	// простой вывод данных
	for md := range marketDataCh {
		fmt.Printf("%s %s bid=%.6f ask=%.6f\n", md.Exchange, md.Symbol, md.Bid, md.Ask)
	}
}
