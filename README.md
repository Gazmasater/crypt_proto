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
	HasMarket   bool
	StepSize    float64
	MinQty      float64
	MinNotional float64
}

// FetchSymbolCapsMEXC читает https://api.mexc.com/api/v3/exchangeInfo.
// ВАЖНО: MEXC может не всегда отдавать orderTypes/статусы в ожидаемом формате,
// поэтому HasMarket делаем "безопасно":
// - если orderTypes отсутствует или неожиданного типа => считаем HasMarket=true, чтобы не убить всё фильтром
// - если orderTypes есть как массив строк => HasMarket=true только если найден "MARKET"
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

		// safe default: если нет orderTypes — не режем
		hasMarket := true

		if otsAny, ok := m["orderTypes"]; ok {
			// если orderTypes есть — пытаемся понять, есть ли MARKET
			hasMarket = false

			if ots, ok := otsAny.([]any); ok {
				for _, v := range ots {
					if s, ok := v.(string); ok && strings.EqualFold(s, "MARKET") {
						hasMarket = true
						break
					}
				}
			} else {
				// неожиданный формат — не ломаем запуск
				hasMarket = true
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
	"sort"
	"strings"
	"sync"
	"syscall"

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

	// ctx нужен до запуска фида: здесь же дергаем exchangeInfo для диагностики/фильтра
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	triangles, _, _, err := domain.LoadTriangles(cfg.TrianglesFile)
	if err != nil {
		log.Fatalf("load triangles: %v", err)
	}
	if len(triangles) == 0 {
		log.Fatal("нет треугольников, нечего мониторить")
	}

	// ===== Prefilter + диагностика (MEXC) =====
	if cfg.Exchange == "MEXC" {
		orig := triangles

		caps, err := mexc.FetchSymbolCapsMEXC(ctx)
		if err != nil {
			log.Printf("WARN: cannot fetch MEXC exchangeInfo: %v (продолжаю без фильтрации)", err)
		} else {
			// диагностика: что реально пришло
			dumpMEXCExchangeInfoStats(caps, triangles)

			before := len(triangles)
			triangles = filterTrianglesByCaps(triangles, caps)
			after := len(triangles)
			log.Printf("[MEXC] triangles filtered by exchangeInfo: before=%d after=%d", before, after)

			// если фильтр убил всё — не ломаем запуск, просто отключаем prefilter
			if after == 0 {
				log.Printf("WARN: exchangeInfo filtering removed everything. Disabling prefilter and using original triangles.")
				triangles = orig
			}
		}
	}

	// пересобираем symbols + indexBySymbol после фильтрации/фоллбека
	symbols, indexBySymbol := rebuildSymbolsAndIndex(triangles)
	if len(symbols) == 0 {
		log.Fatal("нет символов для подписки")
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

	// лог-файл
	logFile, logBuf, arbOut := arb.OpenLogWriter("arbitrage.log")
	defer logFile.Close()
	defer logBuf.Flush()

	events := make(chan domain.Event, 8192)
	var wg sync.WaitGroup

	consumer := arb.NewConsumer(cfg.FeePerLeg, cfg.MinProfit, cfg.MinStart, arbOut)
	consumer.StartFraction = cfg.StartFraction

	// Реальный исполнитель (фикс-старт 2 USDT)
	consumer.Executor = arb.NewRealExecutor(
		mexc.NewTrader(cfg.APIKey, cfg.APISecret, cfg.Debug),
		arbOut,
		2.0,
	)

	consumer.Start(ctx, events, triangles, indexBySymbol, &wg)

	feed.Start(ctx, &wg, symbols, cfg.BookInterval, events)

	<-ctx.Done()
	log.Println("shutting down...")

	wg.Wait()
	log.Println("bye")
}

// Мягкий фильтр:
// - символ должен быть в exchangeInfo (иначе выкидываем)
// - HasMarket должен быть true
// - status проверяем мягко: если пустой — не режем; если известный “торговый” — оставляем; иначе режем
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

			if !c.HasMarket {
				ok = false
				break
			}

			st := strings.ToUpper(strings.TrimSpace(c.Status))
			if st != "" && st != "ENABLED" && st != "TRADING" && st != "1" {
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

func dumpMEXCExchangeInfoStats(caps map[string]mexc.SymbolCaps, tris []domain.Triangle) {
	total := len(caps)

	statusCnt := map[string]int{}
	marketYes := 0
	marketNo := 0

	for _, c := range caps {
		st := strings.ToUpper(strings.TrimSpace(c.Status))
		if st == "" {
			st = "<EMPTY>"
		}
		statusCnt[st]++

		if c.HasMarket {
			marketYes++
		} else {
			marketNo++
		}
	}

	type kv struct {
		K string
		V int
	}
	var ss []kv
	for k, v := range statusCnt {
		ss = append(ss, kv{k, v})
	}
	sort.Slice(ss, func(i, j int) bool { return ss[i].V > ss[j].V })
	if len(ss) > 10 {
		ss = ss[:10]
	}

	log.Printf("[MEXC DIAG] exchangeInfo symbols=%d marketYes=%d marketNo=%d", total, marketYes, marketNo)
	log.Printf("[MEXC DIAG] top statuses:")
	for _, it := range ss {
		log.Printf("  - %s: %d", it.K, it.V)
	}

	need := map[string]struct{}{}
	for _, t := range tris {
		for _, leg := range t.Legs {
			need[leg.Symbol] = struct{}{}
		}
	}

	missing := 0
	var missingList []string
	for s := range need {
		if _, ok := caps[s]; !ok {
			missing++
			if len(missingList) < 50 {
				missingList = append(missingList, s)
			}
		}
	}
	sort.Strings(missingList)

	log.Printf("[MEXC DIAG] symbols in triangles=%d missing in exchangeInfo=%d", len(need), missing)
	if missing > 0 {
		log.Printf("[MEXC DIAG] first missing symbols (up to 50): %s", strings.Join(missingList, ","))
	}

	noMarketInNeed := 0
	var noMarketList []string
	for s := range need {
		if c, ok := caps[s]; ok && !c.HasMarket {
			noMarketInNeed++
			if len(noMarketList) < 50 {
				noMarketList = append(noMarketList, s)
			}
		}
	}
	sort.Strings(noMarketList)

	log.Printf("[MEXC DIAG] symbols in triangles with HasMarket=false: %d", noMarketInNeed)
	if noMarketInNeed > 0 {
		log.Printf("[MEXC DIAG] first no-market symbols (up to 50): %s", strings.Join(noMarketList, ","))
	}

	needStatusCnt := map[string]int{}
	for s := range need {
		if c, ok := caps[s]; ok {
			st := strings.ToUpper(strings.TrimSpace(c.Status))
			if st == "" {
				st = "<EMPTY>"
			}
			needStatusCnt[st]++
		}
	}
	var ns []kv
	for k, v := range needStatusCnt {
		ns = append(ns, kv{k, v})
	}
	sort.Slice(ns, func(i, j int) bool { return ns[i].V > ns[j].V })
	if len(ns) > 10 {
		ns = ns[:10]
	}

	log.Printf("[MEXC DIAG] top statuses among triangle symbols:")
	for _, it := range ns {
		log.Printf("  - %s: %d", it.K, it.V)
	}

	passing := 0
	for s := range need {
		c, ok := caps[s]
		if !ok {
			continue
		}
		st := strings.ToUpper(strings.TrimSpace(c.Status))
		if st != "" && st != "ENABLED" && st != "TRADING" && st != "1" {
			continue
		}
		if !c.HasMarket {
			continue
		}
		passing++
	}
	log.Printf("[MEXC DIAG] single-symbol pass (status+market) among triangle symbols: %d/%d", passing, len(need))
}







