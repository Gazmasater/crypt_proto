package main

import (
	"crypt_proto/internal/collector"
	"crypt_proto/pkg/models"
	"encoding/csv"
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

	csvPath := "mexc_triangles_usdt_routes.csv" // твой CSV с белым списком
	symbols, err := readSymbolsFromCSV(csvPath)
	if err != nil {
		log.Fatal("read CSV symbols:", err)
	}
	log.Printf("Subscribing to %d symbols", len(symbols))

	csvPath = "mexc_triangles_usdt_routes.csv" // твой CSV с белым списком
	symbols, err = readSymbolsFromCSV(csvPath)
	if err != nil {
		log.Fatal("read CSV symbols:", err)
	}
	log.Printf("Subscribing to %d symbols", len(symbols))

	switch exchange {
	case "mexc":
		c = collector.NewMEXCCollector(symbols)
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

func readSymbolsFromCSV(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	header, err := r.Read()
	if err != nil {
		return nil, err
	}

	// находим индексы колонок leg1/2/3
	idx := map[string]int{}
	for i, h := range header {
		idx[strings.ToLower(strings.TrimSpace(h))] = i
	}

	legs := []string{"leg1_symbol", "leg2_symbol", "leg3_symbol"}
	var indices []int
	for _, l := range legs {
		if i, ok := idx[l]; ok {
			indices = append(indices, i)
		} else {
			return nil, &csv.ParseError{StartLine: 0, Err: csv.ErrFieldCount}
		}
	}

	set := map[string]struct{}{}
	for {
		row, err := r.Read()
		if err != nil {
			break
		}
		for _, i := range indices {
			s := strings.TrimSpace(row[i])
			if s != "" {
				set[s] = struct{}{}
			}
		}
	}

	var symbols []string
	for s := range set {
		symbols = append(symbols, s)
	}

	return symbols, nil
}
