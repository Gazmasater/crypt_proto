apikey = "4333ed4b-cd83-49f5-97d1-c399e2349748"
secretkey = "E3848531135EDB4CCFDA0F1BC14CD274"
IP = ""
Название API-ключа = "Arb"
Доступы = "Чтение"



sudo systemctl mask sleep.target suspend.target hibernate.target hybrid-sleep.target



wbs-api.mexc.com/ws 


[https://edis-global.vercel.app/ru/vps-hosting/singapore-singapore
](https://sg.edisglobal.com/)



git pull --rebase origin privat
git push origin privat


BOOK_INTERVAL=100ms
SYMBOLS_FILE=triangles_markets.csv
DEBUG=false


import (
    // ...
    "net/http"
    _ "net/http/pprof"
)


   // pprof HTTP-сервер
    go func() {
        log.Println("pprof on http://localhost:6060/debug/pprof/")
        if err := http.ListenAndServe("localhost:6060", nil); err != nil {
            log.Printf("pprof server error: %v", err)
        }
    }()


	go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30


(pprof) top        # показать топ функций по CPU
(pprof) top10
(pprof) list parsePBWrapperMid   # подробный разбор одной функции
(pprof) quit


go tool pprof http://localhost:6060/debug/pprof/heap


(pprof) top
(pprof) top -cum
(pprof) list parsePBWrapperMid
(pprof) quit




package main

import (
	"encoding/csv"
	"fmt"
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
	csvPath := "mexc_triangles_usdt_routes.csv"

	symbols, err := readSymbolsFromCSV(csvPath, exchange)
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
		c = collector.NewMEXCCollector(symbols, whitelist, marketDataPool)
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
				md.Exchange, md.Symbol, md.Bid, md.Ask,
			)
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
// CSV → symbols
// ------------------------------------------------------------

func readSymbolsFromCSV(path, exchange string) ([]string, error) {
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

	required := []string{"leg1", "leg2", "leg3"}
	var idx []int
	for _, name := range required {
		i, ok := colIndex[strings.ToLower(name)]
		if !ok {
			return nil, fmt.Errorf("CSV missing column %s", name)
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
			s = extractSymbolFromLeg(s)
			if s == "" {
				continue
			}
			s = normalizeSymbolForExchange(s, exchange)
			uniq[s] = struct{}{}
		}
	}

	out := make([]string, 0, len(uniq))
	for s := range uniq {
		out = append(out, s)
	}
	return out, nil
}

func extractSymbolFromLeg(leg string) string {
	parts := strings.Fields(leg) // разделяем по пробелу
	if len(parts) == 0 {
		return ""
	}
	if len(parts) == 1 {
		return parts[0] // уже чистый символ
	}
	return parts[1] // берем второй элемент после BUY/SELL
}

func normalizeSymbolForExchange(sym, exchange string) string {
	parts := strings.Split(sym, "/")
	if len(parts) != 2 {
		return sym
	}
	base, quote := parts[0], parts[1]
	switch exchange {
	case "kucoin", "okx":
		return base + "-" + quote
	case "mexc":
		return base + "_" + quote
	default:
		return sym
	}
}


package collector

import "crypt_proto/pkg/models"

type Collector interface {
	Start(out chan<- *models.MarketData) error
	Stop() error
	Name() string
}




