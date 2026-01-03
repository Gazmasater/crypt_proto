package main

import (
	"encoding/csv"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"crypt_proto/internal/collector"
	"crypt_proto/pkg/models"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load(".env")

	// pprof
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
	marketDataCh := make(chan *models.MarketData, 1000)

	// пул MarketData
	marketDataPool := &sync.Pool{
		New: func() interface{} {
			return new(models.MarketData)
		},
	}

	// === читаем whitelist из CSV ===
	csvPath := "../exchange/data/okx_triangles_usdt.csv"

	symbols, err := readSymbolsFromCSV(csvPath, exchange)
	if err != nil {
		log.Fatalf("read CSV symbols: %v", err)
	}
	log.Printf("Loaded %d unique symbols from %s", len(symbols), csvPath)

	// создаём whitelist
	whitelist := make([]string, len(symbols))
	copy(whitelist, symbols)

	var c collector.Collector

	// создаём collector в зависимости от биржи
	switch exchange {
	case "mexc":
		c = collector.NewMEXCCollector(symbols, whitelist, marketDataPool)
	case "okx":
		//
		c = collector.NewOKXCollector(symbols, marketDataPool)
	//case "kucoin":
	//	c = collector.NewKuCoinCollector(symbols, whitelist, marketDataPool)
	default:
		log.Fatal("Unknown exchange:", exchange)
	}

	// старт collector
	if err := c.Start(marketDataCh); err != nil {
		log.Fatal("start collector:", err)
	}

	// consumer маркет-данных
	go func() {
		for md := range marketDataCh {
			log.Printf("[%s] %s bid=%.8f ask=%.8f",
				md.Exchange, md.Symbol, md.Bid, md.Ask,
			)
			// возвращаем объект обратно в пул
			marketDataPool.Put(md)
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
// CSV → symbols, нормализуем под биржу
// ------------------------------------------------------------
func readSymbolsFromCSV(path, exchange string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)

	// читаем заголовок
	header, err := r.Read()
	if err != nil {
		return nil, err
	}

	// ищем колонки Leg1, Leg2, Leg3
	colIndex := make(map[string]int)
	for i, h := range header {
		colIndex[strings.ToLower(strings.TrimSpace(h))] = i
	}

	required := []string{"leg1", "leg2", "leg3"}
	var idx []int
	for _, name := range required {
		i, ok := colIndex[strings.ToLower(name)]
		if !ok {
			return nil, csv.ErrFieldCount
		}
		idx = append(idx, i)
	}

	// множество уникальных символов
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

			raw := strings.TrimSpace(row[i])
			if raw == "" {
				continue
			}

			// raw = "BUY PEPE/USDT" → вытаскиваем символ
			parts := strings.Fields(raw)
			if len(parts) < 2 {
				continue
			}
			symbol := parts[1] // "PEPE/USDT"

			// нормализуем под биржу
			switch exchange {
			case "mexc":
				symbol = strings.ReplaceAll(symbol, "/", "") // PEPEUSDT
			case "okx", "kucoin":
				symbol = strings.ReplaceAll(symbol, "/", "-") // PEPE-USDT
			default:
				// оставляем как есть
			}

			uniq[symbol] = struct{}{}
		}
	}

	// формируем срез
	out := make([]string, 0, len(uniq))
	for s := range uniq {
		out = append(out, s)
	}

	return out, nil
}
