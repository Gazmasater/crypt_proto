package main

import (
	"crypt_proto/internal/collector"
	"crypt_proto/pkg/models"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"net/http"
	_ "net/http/pprof"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load(".env")
	go func() {
		log.Println("pprof on http://localhost:6060/debug/pprof/")
		if err := http.ListenAndServe("localhost:6060", nil); err != nil {
			log.Printf("pprof server error: %v", err)
		}
	}()
	exchange := strings.ToLower(os.Getenv("EXCHANGE"))
	if exchange == "" {
		log.Fatal("Set EXCHANGE env variable: mexc | okx | kucoin")
	}
	log.Println("EXCHANGE:", exchange)

	marketDataCh := make(chan models.MarketData, 1000)

	var c collector.Collector

	switch exchange {
	case "mexc":
		c = collector.NewMEXCCollector([]string{"BTCUSDT", "ETHUSDT", "ETHBTC"})
	case "okx":
		c = collector.NewOKXCollector([]string{"BTC-USDT", "ETH-USDT", "ETH-BTC"})
	case "kucoin":
		c = collector.NewKuCoinCollector([]string{"BTCUSDT", "ETHUSDT", "ETHBTC"})
	default:
		log.Fatal("Unknown exchange:", exchange)
	}

	// стартуем
	if err := c.Start(marketDataCh); err != nil {
		log.Fatal(err)
	}

	// выводим данные
	go func() {
		for data := range marketDataCh {
			log.Printf("[%s] %s bid=%.8f ask=%.8f\n",
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
