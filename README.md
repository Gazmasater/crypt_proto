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




BOOK_INTERVAL=100ms
SYMBOLS_FILE=triangles_markets.csv
DEBUG=false

# Комиссия (в процентах, одна нога)
FEE_PCT=0.1

# Минимальная прибыль по кругу (в процентах)
MIN_PROFIT_PCT=0.5






package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

// =================== CONFIG ===================

type Config struct {
	BookInterval string  // "100ms" / "10ms" — строка для MEXC топика
	SymbolsFile  string  // CSV с рынками из треугольников
	Debug        bool    // включать ли подробный лог
	FeeRate      float64 // комиссия за одну ногу, ДОЛЯ (0.001 = 0.1%)
	MinProfit    float64 // минимальная прибыль за цикл, ДОЛЯ (0.005 = 0.5%)
}

func loadConfig() (Config, error) {
	_ = godotenv.Load(".env")

	var cfg Config

	// BOOK_INTERVAL
	cfg.BookInterval = os.Getenv("BOOK_INTERVAL")
	if cfg.BookInterval == "" {
		cfg.BookInterval = "100ms"
	}

	// SYMBOLS_FILE
	cfg.SymbolsFile = os.Getenv("SYMBOLS_FILE")
	if cfg.SymbolsFile == "" {
		cfg.SymbolsFile = "triangles_markets.csv"
	}

	// DEBUG
	if strings.ToLower(os.Getenv("DEBUG")) == "true" {
		cfg.Debug = true
	}

	// FEE_PCT — комиссия в процентах, одна нога (например "0.1" = 0.1%)
	feeStr := os.Getenv("FEE_PCT")
	if feeStr == "" {
		feeStr = "0.1" // дефолт: 0.1%
	}
	feePct, err := strconv.ParseFloat(strings.TrimSpace(feeStr), 64)
	if err != nil {
		return cfg, fmt.Errorf("bad FEE_PCT=%q: %w", feeStr, err)
	}
	cfg.FeeRate = feePct / 100.0

	// MIN_PROFIT_PCT — минимальная прибыль по кругу в процентах ("0.5" = 0.5%)
	minStr := os.Getenv("MIN_PROFIT_PCT")
	if minStr == "" {
		minStr = "0.5" // дефолт: 0.5%
	}
	minPct, err := strconv.ParseFloat(strings.TrimSpace(minStr), 64)
	if err != nil {
		return cfg, fmt.Errorf("bad MIN_PROFIT_PCT=%q: %w", minStr, err)
	}
	cfg.MinProfit = minPct / 100.0

	return cfg, nil
}

// глобальный флаг для dlog
var debug bool

func dlog(format string, args ...any) {
	if debug {
		log.Printf(format, args...)
	}
}

// =============== MAIN ===============

func main() {
	// pprof HTTP-сервер
	go func() {
		log.Println("pprof on http://localhost:6060/debug/pprof/")
		if err := http.ListenAndServe("localhost:6060", nil); err != nil {
			log.Printf("pprof server error: %v", err)
		}
	}()

	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	cfg, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}
	debug = cfg.Debug

	log.Printf("Triangles file: %s", cfg.SymbolsFile)
	log.Printf("Book interval: %s", cfg.BookInterval)
	log.Printf("Fee per leg: %.4f %% (rate=%.6f)", cfg.FeeRate*100, cfg.FeeRate)
	log.Printf("Min profit per cycle: %.4f %% (rate=%.6f)", cfg.MinProfit*100, cfg.MinProfit)

	// ---------- загрузка треугольников / рынков ----------

	// тут подставь свои функции
	triangles, err := loadTriangles(cfg.SymbolsFile)
	if err != nil {
		log.Fatalf("loadTriangles: %v", err)
	}
	log.Printf("треугольников всего: %d", len(triangles))

	trisBySymbol := buildIndex(triangles) // map[string][]Triangle
	log.Printf("символов в индексе треугольников: %d", len(trisBySymbol))

	symbols := collectSymbols(trisBySymbol) // []string всех market-символов
	log.Printf("символов для подписки всего: %d", len(symbols))

	// ---------- контекст + каналы ----------

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	events := make(chan Event, 8192) // Event{Symbol, Bid, Ask, BidQty, AskQty}

	var wg sync.WaitGroup

	// ---- Чанкуем символы по 50 на одно WS-подключение ----
	const maxPerConn = 50
	chunks := make([][]string, 0)
	for i := 0; i < len(symbols); i += maxPerConn {
		j := i + maxPerConn
		if j > len(symbols) {
			j = len(symbols)
		}
		chunks = append(chunks, symbols[i:j])
	}
	log.Printf("будем использовать %d WS-подключений", len(chunks))

	for idx, chunk := range chunks {
		wg.Add(1)
		go func(i int, syms []string) {
			log.Printf("[WS #%d] symbols in this conn: %d", i, len(syms))
			runPublicBookTicker(ctx, &wg, syms, cfg.BookInterval, events)
		}(idx, chunk)
	}

	// ---------- консумер: котировки + треугольники ----------

	lastQuotes := make(map[string]Quote)   // последняя котировка по символу
	live       := make(map[string]*LiveInfo) // живые прибыльные треугольники

	go func() {
		for {
			select {
			case ev, ok := <-events:
				if !ok {
					return
				}

				// Обновили котировку
				lastQuotes[ev.Symbol] = Quote{
					Bid:    ev.Bid,
					Ask:    ev.Ask,
					BidQty: ev.BidQty,
					AskQty: ev.AskQty,
				}

				now := time.Now()

				// Все треугольники, где участвует этот символ
				trisForSym := trisBySymbol[ev.Symbol]
				if len(trisForSym) == 0 {
					continue
				}

				// Считаем прибыльные ТОЛЬКО по этим треугольникам
				prof := findProfitableTrianglesForSymbol(ev.Symbol, trisForSym, lastQuotes, cfg.FeeRate, cfg.MinProfit)

				// какие ключи сейчас прибыльны
				currentProfitable := make(map[string]ProfitResult, len(prof))
				for _, pr := range prof {
					k := triKey(pr.Tri)
					currentProfitable[k] = pr

					if lt, ok := live[k]; ok {
						// уже живой — обновляем
						lt.LastSeen = now
						if pr.ProfitNet > lt.MaxProfit {
							lt.MaxProfit = pr.ProfitNet
						}
					} else {
						// новый прибыльный треугольник
						live[k] = &LiveInfo{
							FirstSeen: now,
							LastSeen:  now,
							MaxProfit: pr.ProfitNet,
						}
						fmt.Printf("\n[NEW] triangle %s стал прибыльным (%.3f%%)\n",
							k, pr.ProfitNet*100)
					}
				}

				// Найдём треугольники с этим символом,
				// которые БЫЛИ прибыльными, но на этом тике уже нет в prof
				for _, tri := range trisForSym {
					k := triKey(tri)
					lt, wasLive := live[k]
					if !wasLive {
						continue
					}
					if _, stillProfitable := currentProfitable[k]; stillProfitable {
						continue
					}

					// Перестал быть прибыльным
					lifetime := lt.LastSeen.Sub(lt.FirstSeen)
					fmt.Printf(
						"\n[END] triangle %s перестал быть прибыльным\n"+
							"     Жизнь: %v   maxProfit: %.3f%%\n",
						k, lifetime, lt.MaxProfit*100,
					)
					delete(live, k)
				}

				// Опционально: выводим актуальные прибыльные по этому символу
				if len(prof) > 0 {
					fmt.Printf(
						"\nquotes known: %d symbols, profitable triangles (on %s update): %d\n",
						len(lastQuotes), ev.Symbol, len(prof),
					)
					maxShow := 5
					if len(prof) < maxShow {
						maxShow = len(prof)
					}
					for i := 0; i < maxShow; i++ {
						printTriangleWithDetails(prof[i], lastQuotes)
					}
				}

			case <-ctx.Done():
				return
			}
		}
	}()

	<-ctx.Done()
	log.Println("shutting down...")

	time.Sleep(300 * time.Millisecond)
	close(events)
	wg.Wait()
	log.Println("bye")
}



