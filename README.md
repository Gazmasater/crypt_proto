mx0vglmT3srN1IS19H
135bb7a7509e4421bad692415c53753b



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




exchangeinfo.go
package mexc

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type SymbolCaps struct {
	Symbol      string
	Status      string
	APIEnabled  bool
	HasMarket   bool
	StepSize    float64
	MinQty      float64
	MinNotional float64
}

// FetchSymbolCapsMEXC читает https://api.mexc.com/api/v3/exchangeInfo и строит карту возможностей по символам.
// Мы используем это, чтобы заранее выкинуть пары, которые НЕ поддерживают API или MARKET (чтобы не ловить 10007).
func FetchSymbolCapsMEXC(ctx context.Context) (map[string]SymbolCaps, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.mexc.com/api/v3/exchangeInfo", nil)
	if err != nil {
		return nil, err
	}

	cl := &http.Client{Timeout: 12 * time.Second}
	resp, err := cl.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("exchangeInfo status=%d", resp.StatusCode)
	}

	var root map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&root); err != nil {
		return nil, err
	}

	rawSyms, _ := root["symbols"].([]any)
	out := make(map[string]SymbolCaps, len(rawSyms))

	for _, item := range rawSyms {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}

		symbol, _ := m["symbol"].(string)
		if symbol == "" {
			continue
		}

		status, _ := m["status"].(string)

		// orderTypes: ["LIMIT","MARKET",...]
		hasMarket := false
		if ots, ok := m["orderTypes"].([]any); ok {
			for _, v := range ots {
				if s, ok := v.(string); ok && strings.EqualFold(s, "MARKET") {
					hasMarket = true
					break
				}
			}
		}

		// apiEnabled / isSpotTradingAllowed / apiAllowed — у MEXC встречается по-разному
		apiEnabled := false
		for _, key := range []string{"apiEnabled", "isSpotTradingAllowed", "apiAllowed"} {
			if b, ok := m[key].(bool); ok && b {
				apiEnabled = true
				break
			}
		}

		// filters: LOT_SIZE(stepSize/minQty), MIN_NOTIONAL
		stepSize, minQty, minNotional := 0.0, 0.0, 0.0
		if flt, ok := m["filters"].([]any); ok {
			for _, f := range flt {
				fm, ok := f.(map[string]any)
				if !ok {
					continue
				}
				ft, _ := fm["filterType"].(string)
				switch ft {
				case "LOT_SIZE":
					stepSize = readFloatAny(fm["stepSize"])
					minQty = readFloatAny(fm["minQty"])
				case "MIN_NOTIONAL":
					minNotional = readFloatAny(fm["minNotional"])
				}
			}
		}

		out[symbol] = SymbolCaps{
			Symbol:      symbol,
			Status:      status,
			APIEnabled:  apiEnabled,
			HasMarket:   hasMarket,
			StepSize:    stepSize,
			MinQty:      minQty,
			MinNotional: minNotional,
		}
	}

	return out, nil
}

func readFloatAny(v any) float64 {
	switch t := v.(type) {
	case string:
		f, _ := strconv.ParseFloat(t, 64)
		return f
	case float64:
		return t
	default:
		return 0
	}
}




package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"crypt_proto/arb"
	"crypt_proto/config"
	"crypt_proto/domain"
	"crypt_proto/exchange"
	"crypt_proto/kucoin"
	"crypt_proto/mexc"

	_ "net/http/pprof"
)

func main() {
	// pprof
	go func() {
		log.Println("pprof on http://localhost:6060/debug/pprof/")
		if err := http.ListenAndServe("localhost:6060", nil); err != nil {
			log.Printf("pprof server error: %v", err)
		}
	}()

	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	cfg := config.Load()

	// ВАЖНО: ctx создаём сразу, чтобы можно было дергать exchangeInfo ДО запуска фида/консьюмера
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	triangles, _, _, err := domain.LoadTriangles(cfg.TrianglesFile)
	if err != nil {
		log.Fatalf("load triangles: %v", err)
	}
	if len(triangles) == 0 {
		log.Fatal("нет треугольников, нечего мониторить")
	}

	// ===== Вариант A: заранее фильтруем пары по exchangeInfo (чтобы не ловить 10007) =====
	// Работает сейчас для MEXC.
	if cfg.Exchange == "MEXC" {
		caps, err := mexc.FetchSymbolCapsMEXC(ctx)
		if err != nil {
			log.Printf("WARN: cannot fetch MEXC exchangeInfo: %v (будет риск 10007, оставь runtime-ban как страховку)", err)
		} else {
			before := len(triangles)
			triangles = filterTrianglesByCaps(triangles, caps)
			after := len(triangles)
			log.Printf("[MEXC] triangles filtered by exchangeInfo: before=%d after=%d", before, after)
		}
	}

	// пересобираем symbols + indexBySymbol ПОСЛЕ фильтрации
	symbols, indexBySymbol := rebuildSymbolsAndIndex(triangles)
	if len(symbols) == 0 {
		log.Fatal("нет символов для подписки после фильтрации")
	}
	log.Printf("символов для подписки всего: %d", len(symbols))

	var feed exchange.MarketDataFeed
	switch cfg.Exchange {
	case "MEXC":
		feed = mexc.NewFeed(cfg.Debug)
	case "KUCOIN":
		feed = kucoin.NewFeed(cfg.Debug)
	default:
		log.Fatalf("unknown EXCHANGE=%q (expected MEXC or KUCOIN)", cfg.Exchange)
	}
	log.Printf("Using exchange: %s", feed.Name())

	// лог-файл для арбитража
	logFile, logBuf, arbOut := arb.OpenLogWriter("arbitrage.log")
	defer logFile.Close()
	defer logBuf.Flush()

	events := make(chan domain.Event, 8192)
	var wg sync.WaitGroup

	// consumer
	consumer := arb.NewConsumer(cfg.FeePerLeg, cfg.MinProfit, cfg.MinStart, arbOut)
	consumer.StartFraction = cfg.StartFraction

	// тут у тебя уже есть executor_real.go
	// фиксированный старт = 2 USDT (или из env, если ты добавил TRADE_AMOUNT_USDT в конфиг)
	consumer.Executor = arb.NewRealExecutor(
		mexc.NewTrader(cfg.APIKey, cfg.APISecret, cfg.Debug),
		arbOut,
		2.0,
	)

	consumer.Start(ctx, events, triangles, indexBySymbol, &wg)

	// feed
	feed.Start(ctx, &wg, symbols, cfg.BookInterval, events)

	<-ctx.Done()
	log.Println("shutting down...")

	wg.Wait()
	log.Println("bye")
}

// filterTrianglesByCaps оставляет только те треугольники, у которых все 3 ноги:
// - есть в caps
// - status ENABLED
// - APIEnabled=true
// - MARKET доступен
func filterTrianglesByCaps(tris []domain.Triangle, caps map[string]mexc.SymbolCaps) []domain.Triangle {
	out := make([]domain.Triangle, 0, len(tris))

	for _, t := range tris {
		ok := true
		for _, leg := range t.Legs {
			c, okCap := caps[leg.Symbol]
			if !okCap {
				ok = false
				break
			}
			if !strings.EqualFold(c.Status, "ENABLED") {
				ok = false
				break
			}
			if !c.APIEnabled {
				ok = false
				break
			}
			if !c.HasMarket {
				ok = false
				break
			}
		}
		if ok {
			out = append(out, t)
		}
	}

	return out
}

func rebuildSymbolsAndIndex(tris []domain.Triangle) ([]string, map[string][]int) {
	symbolSet := make(map[string]struct{})
	for _, t := range tris {
		for _, leg := range t.Legs {
			symbolSet[leg.Symbol] = struct{}{}
		}
	}

	symbols := make([]string, 0, len(symbolSet))
	for s := range symbolSet {
		symbols = append(symbols, s)
	}

	index := make(map[string][]int)
	for i, t := range tris {
		for _, leg := range t.Legs {
			index[leg.Symbol] = append(index[leg.Symbol], i)
		}
	}
	return symbols, index
}







