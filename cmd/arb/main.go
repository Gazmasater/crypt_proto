package main

import (
	"crypt_proto/internal/collector"
	"crypt_proto/pkg/models"
	"encoding/csv"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"strings"
	"syscall"

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

	// канал маркет-данных
	marketDataCh := make(chan models.MarketData, 1000)

	// === читаем whitelist из CSV ===
	csvPath := "mexc_triangles_usdt_routes.csv"

	// symbols для подписки
	symbols, err := readSymbolsFromCSV(csvPath)
	if err != nil {
		log.Fatalf("read CSV symbols: %v", err)
	}
	log.Printf("Loaded %d unique symbols from %s", len(symbols), csvPath)

	// создаём whitelist
	whitelist := make([]string, len(symbols))
	copy(whitelist, symbols)

	var c collector.Collector

	switch exchange {
	case "mexc":
		// передаём whitelist
		c = collector.NewMEXCCollector(symbols, whitelist)

	case "okx":
		c = collector.NewOKXCollector(symbols)

	case "kucoin":
		c = collector.NewKuCoinCollector(symbols)

	default:
		log.Fatal("Unknown exchange:", exchange)
	}

	// старт
	if err := c.Start(marketDataCh); err != nil {
		log.Fatal("start collector:", err)
	}

	// consumer
	go func() {
		for md := range marketDataCh {
			log.Printf("[%s] %s bid=%.8f ask=%.8f",
				md.Exchange,
				md.Symbol,
				md.Bid,
				md.Ask,
			)
		}
	}()

	// graceful shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	log.Println("Stopping collector...")
	if err := c.Stop(); err != nil {
		log.Println("Stop error:", err)
	}
}

// ------------------------------------------------------------
// CSV → symbols
// ------------------------------------------------------------

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

	colIndex := make(map[string]int)
	for i, h := range header {
		colIndex[strings.ToLower(strings.TrimSpace(h))] = i
	}

	required := []string{
		"leg1_symbol",
		"leg2_symbol",
		"leg3_symbol",
	}

	var idx []int
	for _, name := range required {
		i, ok := colIndex[name]
		if !ok {
			return nil, csv.ErrFieldCount
		}
		idx = append(idx, i)
	}

	uniq := make(map[string]struct{})

	for {
		row, err := r.Read()
		if err != nil {
			break
		}

		for _, i := range idx {
			if i >= len(row) {
				continue
			}
			s := strings.TrimSpace(row[i])
			if s != "" {
				uniq[s] = struct{}{}
			}
		}
	}

	out := make([]string, 0, len(uniq))
	for s := range uniq {
		out = append(out, s)
	}

	return out, nil
}
