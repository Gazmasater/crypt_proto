package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

// ===============================
// MEXC API structs
// ===============================

type exchangeInfo struct {
	Symbols []symbolInfo `json:"symbols"`
}

type symbolInfo struct {
	Symbol                     string         `json:"symbol"`
	BaseAsset                  string         `json:"baseAsset"`
	QuoteAsset                 string         `json:"quoteAsset"`
	Status                     string         `json:"status"`
	IsSpotTradingAllowed       bool           `json:"isSpotTradingAllowed"`
	OrderTypes                 []string       `json:"orderTypes"`
	QuoteOrderQtyMarketAllowed bool           `json:"quoteOrderQtyMarketAllowed"`
	BaseAssetPrecision         int            `json:"baseAssetPrecision"`
	QuoteAssetPrecision        int            `json:"quoteAssetPrecision"`
	Filters                    []symbolFilter `json:"filters"`
}

type symbolFilter struct {
	FilterType  string `json:"filterType"`
	StepSize    string `json:"stepSize"`
	MinQty      string `json:"minQty"`
	MinNotional string `json:"minNotional"` // иногда бывает, не всегда
}

type bookTicker struct {
	Symbol   string `json:"symbol"`
	BidPrice string `json:"bidPrice"`
	AskPrice string `json:"askPrice"`
	BidQty   string `json:"bidQty"`
	AskQty   string `json:"askQty"`
}

// ===============================
// Internal models
// ===============================

type pairMarket struct {
	Symbol string
	Base   string
	Quote  string
}

type pairKey struct {
	A, B string
}

type rules struct {
	Symbol                     string
	Base, Quote                string
	SpotAllowed                bool
	HasMarket                  bool
	QuoteOrderQtyMarketAllowed bool
	StepSize                   float64 // base step
	MinQty                     float64 // base min qty
	BasePrec                   int
	QuotePrec                  int
}

type tob struct {
	Bid, Ask       float64
	BidQty, AskQty float64
}

// ===============================
// Helpers
// ===============================

func parseFloat(s string) float64 {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	v, _ := strconv.ParseFloat(s, 64)
	return v
}

func hasMarketOrder(orderTypes []string) bool {
	for _, t := range orderTypes {
		if strings.EqualFold(strings.TrimSpace(t), "MARKET") {
			return true
		}
	}
	return false
}

func loadRulesFromSymbolInfo(si symbolInfo) rules {
	r := rules{
		Symbol: si.Symbol,
		Base:   si.BaseAsset,
		Quote:  si.QuoteAsset,

		SpotAllowed:                si.IsSpotTradingAllowed,
		HasMarket:                  hasMarketOrder(si.OrderTypes),
		QuoteOrderQtyMarketAllowed: si.QuoteOrderQtyMarketAllowed,

		BasePrec:  si.BaseAssetPrecision,
		QuotePrec: si.QuoteAssetPrecision,
	}

	// Берём step/minQty из LOT_SIZE или MARKET_LOT_SIZE
	for _, f := range si.Filters {
		ft := strings.ToUpper(strings.TrimSpace(f.FilterType))
		if ft == "LOT_SIZE" || ft == "MARKET_LOT_SIZE" {
			if r.StepSize == 0 {
				r.StepSize = parseFloat(f.StepSize)
			}
			if r.MinQty == 0 {
				r.MinQty = parseFloat(f.MinQty)
			}
		}
	}

	return r
}

// floorToStep: округление вниз к кратности шагу
func floorToStep(x, step float64) float64 {
	if step <= 0 {
		return x
	}
	return math.Floor(x/step) * step
}

func ceilToStep(x, step float64) float64 {
	if step <= 0 {
		return x
	}
	return math.Ceil(x/step) * step
}

// Оценка минимального USDT старта для символа:
// берём minQty (base) и цену mid, переводим quote->USDT при необходимости.
func minUSDTForSymbol(sym rules, prices map[string]tob) (float64, string) {
	p, ok := prices[sym.Symbol]
	if !ok || p.Bid <= 0 || p.Ask <= 0 {
		return 0, "price not ok"
	}
	mid := (p.Bid + p.Ask) / 2
	if mid <= 0 {
		return 0, "price not ok"
	}

	// Если minQty неизвестен — это подозрительно, но не будем сразу резать.
	minQty := sym.MinQty
	if minQty <= 0 {
		// fallback: пытаемся вывести минимально “не ноль” (очень мягко)
		minQty = 0
	}

	// минимальный “размер позиции” в quote валюте:
	needQuote := 0.0
	if minQty > 0 {
		needQuote = minQty * mid
	}

	// Перевод quote -> USDT:
	quote := sym.Quote
	if quote == "USDT" {
		return needQuote, ""
	}
	if quote == "USDC" {
		// USDCUSDT или USDTUSDC
		if q, ok := prices["USDCUSDT"]; ok && q.Bid > 0 {
			return needQuote * q.Bid, ""
		}
		if q, ok := prices["USDTUSDC"]; ok && q.Ask > 0 {
			return needQuote / q.Ask, ""
		}
		return 0, "no USDC<->USDT price"
	}

	// Если quote не USDT/USDC — пробуем через quoteUSDT или USDTquote
	if q, ok := prices[quote+"USDT"]; ok && q.Bid > 0 {
		return needQuote * q.Bid, ""
	}
	if q, ok := prices["USDT"+quote]; ok && q.Ask > 0 {
		return needQuote / q.Ask, ""
	}

	// через USDC как промежуточную
	if q1, ok := prices[quote+"USDC"]; ok && q1.Bid > 0 {
		amtUSDC := needQuote * q1.Bid
		if q2, ok := prices["USDCUSDT"]; ok && q2.Bid > 0 {
			return amtUSDC * q2.Bid, ""
		}
	}
	if q1, ok := prices["USDC"+quote]; ok && q1.Ask > 0 {
		amtUSDC := needQuote / q1.Ask
		if q2, ok := prices["USDCUSDT"]; ok && q2.Bid > 0 {
			return amtUSDC * q2.Bid, ""
		}
	}

	return 0, "cannot convert quote->USDT"
}

// ===============================

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	// 1) exchangeInfo
	resp, err := http.Get("https://api.mexc.com/api/v3/exchangeInfo")
	if err != nil {
		log.Fatalf("get exchangeInfo: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		log.Fatalf("exchangeInfo status %d: %s", resp.StatusCode, string(b))
	}

	var info exchangeInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		log.Fatalf("decode exchangeInfo: %v", err)
	}
	log.Printf("total symbols from API: %d", len(info.Symbols))

	// 2) bookTicker (цены)
	prices := make(map[string]tob, 100000)
	{
		r2, err := http.Get("https://api.mexc.com/api/v3/ticker/bookTicker")
		if err != nil {
			log.Fatalf("get bookTicker: %v", err)
		}
		defer r2.Body.Close()
		if r2.StatusCode != http.StatusOK {
			b, _ := io.ReadAll(r2.Body)
			log.Fatalf("bookTicker status %d: %s", r2.StatusCode, string(b))
		}

		var arr []bookTicker
		if err := json.NewDecoder(r2.Body).Decode(&arr); err != nil {
			log.Fatalf("decode bookTicker: %v", err)
		}
		for _, it := range arr {
			s := strings.TrimSpace(it.Symbol)
			if s == "" {
				continue
			}
			prices[s] = tob{
				Bid:    parseFloat(it.BidPrice),
				Ask:    parseFloat(it.AskPrice),
				BidQty: parseFloat(it.BidQty),
				AskQty: parseFloat(it.AskQty),
			}
		}
		log.Printf("bookTicker loaded: %d symbols", len(prices))
	}

	// 3) фильтрация “торгуемых” маркетов + сбор правил
	markets := make([]pairMarket, 0, len(info.Symbols))
	rulesBySymbol := make(map[string]rules, len(info.Symbols))

	for _, s := range info.Symbols {
		base := strings.TrimSpace(s.BaseAsset)
		quote := strings.TrimSpace(s.QuoteAsset)
		if base == "" || quote == "" {
			continue
		}

		// Мягкий фильтр по статусу (как у тебя было)
		if s.Status != "" && s.Status != "ENABLED" && s.Status != "TRADING" && s.Status != "1" {
			continue
		}

		pm := pairMarket{Symbol: strings.TrimSpace(s.Symbol), Base: base, Quote: quote}
		if pm.Symbol == "" {
			continue
		}

		markets = append(markets, pm)
		rulesBySymbol[pm.Symbol] = loadRulesFromSymbolInfo(s)
	}

	log.Printf("filtered markets: %d", len(markets))

	// 4) pairmap + граф валют
	pairmap := make(map[pairKey][]pairMarket)
	adj := make(map[string]map[string]struct{})

	addEdge := func(a, b string) {
		if adj[a] == nil {
			adj[a] = make(map[string]struct{})
		}
		adj[a][b] = struct{}{}
	}

	for _, m := range markets {
		a, b := m.Base, m.Quote
		key := pairKey{A: a, B: b}
		if a > b {
			key = pairKey{A: b, B: a}
		}
		pairmap[key] = append(pairmap[key], m)

		addEdge(a, b)
		addEdge(b, a)
	}

	log.Printf("currencies (vertices): %d, pair keys: %d", len(adj), len(pairmap))

	// 5) индекс валют
	coins := make([]string, 0, len(adj))
	for c := range adj {
		coins = append(coins, c)
	}
	sort.Strings(coins)

	idx := make(map[string]int, len(coins))
	for i, c := range coins {
		idx[c] = i
	}

	neighbors := make([]map[int]struct{}, len(coins))
	for i := range neighbors {
		neighbors[i] = make(map[int]struct{})
	}
	for c, neighs := range adj {
		i := idx[c]
		for nb := range neighs {
			j := idx[nb]
			neighbors[i][j] = struct{}{}
		}
	}

	// 6) поиск валютных треугольников
	type triangle struct{ A, B, C string }
	triangles := make([]triangle, 0)

	for i := 0; i < len(coins); i++ {
		ni := neighbors[i]
		for j := range ni {
			if j <= i {
				continue
			}
			nj := neighbors[j]
			for k := range ni {
				if k <= j {
					continue
				}
				if _, ok := nj[k]; ok {
					triangles = append(triangles, triangle{A: coins[i], B: coins[j], C: coins[k]})
				}
			}
		}
	}
	log.Printf("found currency triangles: %d", len(triangles))

	// 7) output CSV (расширенный)
	outFile := "triangles_markets_enriched.csv"
	f, err := os.Create(outFile)
	if err != nil {
		log.Fatalf("create %s: %v", outFile, err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	header := []string{
		"base1", "quote1", "symbol1", "step1", "minQty1", "spot1", "market1",
		"base2", "quote2", "symbol2", "step2", "minQty2", "spot2", "market2",
		"base3", "quote3", "symbol3", "step3", "minQty3", "spot3", "market3",
		"min_start_usdt", "reason",
	}
	if err := w.Write(header); err != nil {
		log.Fatalf("write header: %v", err)
	}

	// берём первый маркет для пары (как у тебя), но теперь ещё тянем rules
	pick := func(x, y string) (pairMarket, bool) {
		key := pairKey{A: x, B: y}
		if x > y {
			key = pairKey{A: y, B: x}
		}
		list := pairmap[key]
		if len(list) == 0 {
			return pairMarket{}, false
		}
		return list[0], true
	}

	// форматирование
	ff := func(v float64) string { return fmt.Sprintf("%.12f", v) }
	bf := func(b bool) string {
		if b {
			return "1"
		}
		return "0"
	}

	okCount := 0
	total := 0

	for _, t := range triangles {
		m1, ok1 := pick(t.A, t.B)
		m2, ok2 := pick(t.B, t.C)
		m3, ok3 := pick(t.A, t.C)
		if !ok1 || !ok2 || !ok3 {
			continue
		}
		total++

		r1, okr1 := rulesBySymbol[m1.Symbol]
		r2, okr2 := rulesBySymbol[m2.Symbol]
		r3, okr3 := rulesBySymbol[m3.Symbol]

		reason := "ok"
		minStartUSDT := 0.0

		// базовая валидация rules
		if !okr1 || !okr2 || !okr3 {
			reason = "no rules"
		} else {
			// spot/market checks
			if !r1.SpotAllowed || !r2.SpotAllowed || !r3.SpotAllowed {
				reason = "spot not allowed"
			} else if !r1.HasMarket || !r2.HasMarket || !r3.HasMarket {
				reason = "market not allowed"
			} else {
				// price checks + min start
				a1, e1 := minUSDTForSymbol(r1, prices)
				a2, e2 := minUSDTForSymbol(r2, prices)
				a3, e3 := minUSDTForSymbol(r3, prices)

				if e1 != "" || e2 != "" || e3 != "" {
					// соберём конкретику
					reason = "price/convert not ok"
					if e1 != "" {
						reason += " sym1=" + e1
					}
					if e2 != "" {
						reason += " sym2=" + e2
					}
					if e3 != "" {
						reason += " sym3=" + e3
					}
				} else {
					// берём максимум как “минимально безопасный старт” (чтобы любая нога не упала в minQty)
					minStartUSDT = math.Max(a1, math.Max(a2, a3))
					okCount++
				}
			}
		}

		rec := []string{
			m1.Base, m1.Quote, m1.Symbol, ff(r1.StepSize), ff(r1.MinQty), bf(r1.SpotAllowed), bf(r1.HasMarket),
			m2.Base, m2.Quote, m2.Symbol, ff(r2.StepSize), ff(r2.MinQty), bf(r2.SpotAllowed), bf(r2.HasMarket),
			m3.Base, m3.Quote, m3.Symbol, ff(r3.StepSize), ff(r3.MinQty), bf(r3.SpotAllowed), bf(r3.HasMarket),
			fmt.Sprintf("%.6f", minStartUSDT), reason,
		}

		if err := w.Write(rec); err != nil {
			log.Fatalf("write record: %v", err)
		}
	}

	log.Printf("triangles total written: %d (ok=%d)", total, okCount)
	log.Printf("file: %s", outFile)
	fmt.Println("Готово, файл:", outFile)

	// маленькая пауза, чтобы flush успел на некоторых FS
	time.Sleep(50 * time.Millisecond)
}
