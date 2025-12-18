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




export TRADE_AMOUNT_USDT=100
export FEE_PCT=0.04
export SELL_SAFETY=0.995

export TRIANGLES_FILE=triangles_markets.csv
export TRIANGLES_ENRICHED_FILE=triangles_markets_enriched.csv

go run ./cmd/triangles_enrich_mexc




package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// ======================
// ENV
// ======================

func envFloat(key string, def float64) float64 {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return def
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return def
	}
	return f
}

func envString(key, def string) string {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return def
	}
	return v
}

// ======================
// MEXC API models
// ======================

type mexcExchangeInfo struct {
	Symbols []mexcSymbol `json:"symbols"`
}

type mexcSymbol struct {
	Symbol              string       `json:"symbol"`
	Status              string       `json:"status"`
	BaseAsset           string       `json:"baseAsset"`
	QuoteAsset          string       `json:"quoteAsset"`
	BaseAssetPrecision  any          `json:"baseAssetPrecision"`
	QuoteAssetPrecision any          `json:"quoteAssetPrecision"`
	OrderTypes          []string     `json:"orderTypes"`
	Permissions         []string     `json:"permissions"`
	Filters             []mexcFilter `json:"filters"`

	// важные "нестандартные" поля MEXC:
	IsSpotTradingAllowed       *bool  `json:"isSpotTradingAllowed"`       // НЕ используем как решающее!
	QuoteOrderQtyMarketAllowed *bool  `json:"quoteOrderQtyMarketAllowed"` // если есть
	BaseSizePrecision          string `json:"baseSizePrecision"`          // пример: "0.000001"
	QuoteAmountPrecisionMarket string `json:"quoteAmountPrecisionMarket"` // пример: "1"
	MaxQuoteAmountMarket       string `json:"maxQuoteAmountMarket"`       // иногда нужно
}

type mexcFilter struct {
	FilterType string `json:"filterType"`

	// LOT_SIZE / MARKET_LOT_SIZE
	StepSize string `json:"stepSize"`
	MinQty   string `json:"minQty"`

	// MIN_NOTIONAL / NOTIONAL
	MinNotional string `json:"minNotional"`
	Notional    string `json:"notional"`
}

type mexcBookTicker struct {
	Symbol string `json:"symbol"`
	Bid    string `json:"bidPrice"`
	Ask    string `json:"askPrice"`
}

// ======================
// Rules normalized
// ======================

type SymbolRules struct {
	Symbol string
	Base   string
	Quote  string

	// решающие флаги для генератора
	SpotAllowed   bool // по permissions содержит SPOT
	MarketAllowed bool // orderTypes содержит MARKET

	// "сырое" поле для диагностики
	RawIsSpotTradingAllowed bool
	RawHasIsSpotField       bool

	// quoteOrderQty market возможность (если можно купить на сумму)
	QuoteOrderQtyMarketAllowed bool

	// шаг/точности
	BaseStep float64

	// из фильтров (если есть)
	MinQty      float64
	MinNotional float64

	// точности
	BasePrecision         int
	QuotePrecision        int
	QuotePrecisionMarket  int // из quoteAmountPrecisionMarket (если есть), иначе = QuotePrecision
}

type Quote struct {
	Bid float64
	Ask float64
}

// ======================
// Utils
// ======================

func parseFloat(s string) float64 {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func asInt(v any) (int, bool) {
	switch x := v.(type) {
	case float64:
		return int(x), true
	case int:
		return x, true
	case int64:
		return int(x), true
	case string:
		i, err := strconv.Atoi(x)
		if err != nil {
			return 0, false
		}
		return i, true
	default:
		return 0, false
	}
}

func hasMarket(orderTypes []string) bool {
	for _, s := range orderTypes {
		if strings.EqualFold(strings.TrimSpace(s), "MARKET") {
			return true
		}
	}
	return false
}

func hasPermission(perms []string, want string) bool {
	for _, p := range perms {
		if strings.EqualFold(strings.TrimSpace(p), want) {
			return true
		}
	}
	return false
}

func decimalsFromStep(step float64) int {
	if step <= 0 {
		return -1
	}
	s := strconv.FormatFloat(step, 'f', 18, 64)
	s = strings.TrimRight(s, "0")
	if strings.Contains(s, ".") {
		parts := strings.SplitN(s, ".", 2)
		return len(parts[1])
	}
	return 0
}

func floorToStep(x, step float64) float64 {
	if step <= 0 {
		return x
	}
	return math.Floor(x/step) * step
}

func truncToDecimals(x float64, dec int) float64 {
	if dec < 0 {
		return x
	}
	p := math.Pow10(dec)
	return math.Trunc(x*p) / p
}

// ======================
// Build rules from exchangeInfo
// ======================

func buildRules(info mexcExchangeInfo) map[string]SymbolRules {
	out := make(map[string]SymbolRules, len(info.Symbols))

	for _, s := range info.Symbols {
		sym := strings.TrimSpace(s.Symbol)
		if sym == "" {
			continue
		}

		r := SymbolRules{
			Symbol: sym,
			Base:   strings.TrimSpace(s.BaseAsset),
			Quote:  strings.TrimSpace(s.QuoteAsset),
		}

		// status: на MEXC чаще "1"
		statusOK := (strings.TrimSpace(s.Status) == "" || strings.TrimSpace(s.Status) == "1" ||
			strings.EqualFold(strings.TrimSpace(s.Status), "ENABLED") ||
			strings.EqualFold(strings.TrimSpace(s.Status), "TRADING"))

		// РЕШАЮЩЕЕ: SpotAllowed — только через permissions содержит "SPOT"
		// (isSpotTradingAllowed у BTCUSDT = false, но permissions=["SPOT"] и MARKET есть)
		r.SpotAllowed = statusOK && hasPermission(s.Permissions, "SPOT")

		// MarketAllowed
		r.MarketAllowed = hasMarket(s.OrderTypes)

		// raw isSpotTradingAllowed (только для диагностики)
		if s.IsSpotTradingAllowed != nil {
			r.RawHasIsSpotField = true
			r.RawIsSpotTradingAllowed = *s.IsSpotTradingAllowed
		}

		// precision
		if v, ok := asInt(s.BaseAssetPrecision); ok {
			r.BasePrecision = v
		}
		if v, ok := asInt(s.QuoteAssetPrecision); ok {
			r.QuotePrecision = v
		} else {
			r.QuotePrecision = 2
		}

		// quote precision for MARKET quote amount
		r.QuotePrecisionMarket = r.QuotePrecision
		if qpm := strings.TrimSpace(s.QuoteAmountPrecisionMarket); qpm != "" {
			if iv, err := strconv.Atoi(qpm); err == nil && iv >= 0 {
				r.QuotePrecisionMarket = iv
			}
		}

		// quoteOrderQtyMarketAllowed
		// 1) если биржа прислала булево — используем
		// 2) иначе "эвристика": если есть quoteAmountPrecisionMarket и maxQuoteAmountMarket — считаем что можно
		if s.QuoteOrderQtyMarketAllowed != nil {
			r.QuoteOrderQtyMarketAllowed = *s.QuoteOrderQtyMarketAllowed
		} else {
			if strings.TrimSpace(s.QuoteAmountPrecisionMarket) != "" && strings.TrimSpace(s.MaxQuoteAmountMarket) != "" {
				r.QuoteOrderQtyMarketAllowed = true
			}
		}

		// filters (если есть) — minQty/step/minNotional
		var lotStep, lotMin float64
		var mktStep, mktMin float64
		var minNotional float64

		for _, f := range s.Filters {
			ft := strings.ToUpper(strings.TrimSpace(f.FilterType))
			switch ft {
			case "LOT_SIZE":
				lotStep = parseFloat(f.StepSize)
				lotMin = parseFloat(f.MinQty)
			case "MARKET_LOT_SIZE":
				mktStep = parseFloat(f.StepSize)
				mktMin = parseFloat(f.MinQty)
			case "MIN_NOTIONAL":
				if mn := parseFloat(f.MinNotional); mn > 0 {
					minNotional = mn
				}
			case "NOTIONAL":
				if mn := parseFloat(f.Notional); mn > 0 {
					minNotional = mn
				}
				if mn := parseFloat(f.MinNotional); mn > 0 {
					minNotional = mn
				}
			}
		}

		// ВАЖНО для MEXC: часто нет LOT_SIZE в filters, но есть baseSizePrecision
		// Пример BTCUSDT: baseSizePrecision="0.000001"
		baseStepFromPrecision := parseFloat(s.BaseSizePrecision)

		// выбираем step/minQty: MARKET_LOT_SIZE > LOT_SIZE > baseSizePrecision
		if mktStep > 0 {
			r.BaseStep = mktStep
		} else if lotStep > 0 {
			r.BaseStep = lotStep
		} else if baseStepFromPrecision > 0 {
			r.BaseStep = baseStepFromPrecision
		}

		if mktMin > 0 {
			r.MinQty = mktMin
		} else if lotMin > 0 {
			r.MinQty = lotMin
		} else {
			r.MinQty = 0
		}

		r.MinNotional = minNotional

		out[sym] = r
	}

	return out
}

// ======================
// API calls
// ======================

func httpGetJSON(url string, dst any) error {
	c := &http.Client{Timeout: 20 * time.Second}
	resp, err := c.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("http %d: %s", resp.StatusCode, string(b))
	}
	return json.NewDecoder(resp.Body).Decode(dst)
}

func fetchExchangeInfo(baseURL string) (mexcExchangeInfo, error) {
	var info mexcExchangeInfo
	err := httpGetJSON(baseURL+"/api/v3/exchangeInfo", &info)
	return info, err
}

func fetchBookTickerAll(baseURL string) (map[string]Quote, error) {
	var raw any
	if err := httpGetJSON(baseURL+"/api/v3/ticker/bookTicker", &raw); err != nil {
		return nil, err
	}

	out := map[string]Quote{}

	switch v := raw.(type) {
	case []any:
		for _, it := range v {
			m, ok := it.(map[string]any)
			if !ok {
				continue
			}
			sym, _ := m["symbol"].(string)
			if strings.TrimSpace(sym) == "" {
				continue
			}
			bs, _ := m["bidPrice"].(string)
			as, _ := m["askPrice"].(string)
			if bs == "" {
				if bf, ok := m["bidPrice"].(float64); ok {
					bs = strconv.FormatFloat(bf, 'f', 18, 64)
				}
			}
			if as == "" {
				if af, ok := m["askPrice"].(float64); ok {
					as = strconv.FormatFloat(af, 'f', 18, 64)
				}
			}
			bid := parseFloat(bs)
			ask := parseFloat(as)
			if bid <= 0 || ask <= 0 {
				continue
			}
			out[sym] = Quote{Bid: bid, Ask: ask}
		}
	default:
		b, _ := json.Marshal(raw)
		var arr []mexcBookTicker
		if err := json.Unmarshal(b, &arr); err != nil {
			return nil, fmt.Errorf("unexpected bookTicker format")
		}
		for _, t := range arr {
			bid := parseFloat(t.Bid)
			ask := parseFloat(t.Ask)
			if bid <= 0 || ask <= 0 {
				continue
			}
			out[t.Symbol] = Quote{Bid: bid, Ask: ask}
		}
	}

	if len(out) == 0 {
		return nil, fmt.Errorf("bookTicker: empty result")
	}
	return out, nil
}

// ======================
// Triangle simulation
// ======================

type Edge struct {
	Symbol string
	From   string
	To     string
	Side   string // BUY or SELL (market)
}

type Cycle struct {
	E1, E2, E3 Edge
}

type SimCfg struct {
	FeePct     float64 // percent, e.g. 0.04
	SellSafety float64 // 0.995
}

func applyEdge(amount float64, e Edge, rules SymbolRules, q Quote, cfg SimCfg) (float64, error) {
	if amount <= 0 {
		return 0, errors.New("amount<=0")
	}
	fee := cfg.FeePct / 100.0
	if fee < 0 {
		fee = 0
	}
	sellSafety := cfg.SellSafety
	if sellSafety <= 0 || sellSafety > 1 {
		sellSafety = 0.995
	}

	step := rules.BaseStep
	minQty := rules.MinQty
	minNotional := rules.MinNotional

	switch e.Side {
	case "BUY":
		// from QUOTE -> to BASE, pay quote amount = amount
		// если можно BUY по quoteOrderQty — используем amount с ограничением по quoteAmountPrecisionMarket
		if rules.QuoteOrderQtyMarketAllowed {
			dec := rules.QuotePrecisionMarket
			amtQuote := truncToDecimals(amount, dec)
			if amtQuote <= 0 {
				return 0, fmt.Errorf("buy amount<=0 after quote precisionMarket dec=%d", dec)
			}
			if minNotional > 0 && amtQuote+1e-12 < minNotional {
				return 0, fmt.Errorf("buy minNotional: need>=%.12f have=%.12f", minNotional, amtQuote)
			}
			gotBase := (amtQuote / q.Ask) * (1.0 - fee)

			// проверка minQty (если она известна)
			if minQty > 0 && gotBase+1e-12 < minQty {
				return 0, fmt.Errorf("buy minQty: need>=%.12f got=%.12f", minQty, gotBase)
			}
			return gotBase, nil
		}

		// иначе BUY quantity по ask и режем step
		qtyRaw := amount / q.Ask
		qty := qtyRaw
		if step > 0 {
			qty = floorToStep(qtyRaw, step)
		}
		if qty <= 0 {
			return 0, fmt.Errorf("buy qty<=0 after step (raw=%.12f step=%.12f)", qtyRaw, step)
		}
		if minQty > 0 && qty+1e-12 < minQty {
			return 0, fmt.Errorf("buy minQty: need>=%.12f got=%.12f", minQty, qty)
		}
		notional := qty * q.Ask
		if minNotional > 0 && notional+1e-12 < minNotional {
			return 0, fmt.Errorf("buy minNotional: need>=%.12f got=%.12f", minNotional, notional)
		}
		gotBase := qty * (1.0 - fee)
		return gotBase, nil

	case "SELL":
		// from BASE -> to QUOTE, sell base amount with safety
		qtyRaw := amount * sellSafety
		qty := qtyRaw
		if step > 0 {
			qty = floorToStep(qtyRaw, step)
		}
		if qty <= 0 {
			return 0, fmt.Errorf("sell qty<=0 after step (raw=%.12f step=%.12f)", qtyRaw, step)
		}
		if minQty > 0 && qty+1e-12 < minQty {
			return 0, fmt.Errorf("sell minQty: need>=%.12f got=%.12f", minQty, qty)
		}
		notional := qty * q.Bid
		if minNotional > 0 && notional+1e-12 < minNotional {
			return 0, fmt.Errorf("sell minNotional: need>=%.12f got=%.12f", minNotional, notional)
		}
		gotQuote := qty * q.Bid * (1.0 - fee)
		return gotQuote, nil

	default:
		return 0, fmt.Errorf("unknown side=%s", e.Side)
	}
}

func simulateCycle(startUSDT float64, c Cycle, rulesMap map[string]SymbolRules, quotes map[string]Quote, cfg SimCfg) (float64, error) {
	amt := startUSDT

	r1, ok := rulesMap[c.E1.Symbol]
	if !ok {
		return 0, fmt.Errorf("no rules %s", c.E1.Symbol)
	}
	q1, ok := quotes[c.E1.Symbol]
	if !ok {
		return 0, fmt.Errorf("no price %s", c.E1.Symbol)
	}
	amt, err := applyEdge(amt, c.E1, r1, q1, cfg)
	if err != nil {
		return 0, fmt.Errorf("leg1 %s: %v", c.E1.Symbol, err)
	}

	r2, ok := rulesMap[c.E2.Symbol]
	if !ok {
		return 0, fmt.Errorf("no rules %s", c.E2.Symbol)
	}
	q2, ok := quotes[c.E2.Symbol]
	if !ok {
		return 0, fmt.Errorf("no price %s", c.E2.Symbol)
	}
	amt, err = applyEdge(amt, c.E2, r2, q2, cfg)
	if err != nil {
		return 0, fmt.Errorf("leg2 %s: %v", c.E2.Symbol, err)
	}

	r3, ok := rulesMap[c.E3.Symbol]
	if !ok {
		return 0, fmt.Errorf("no rules %s", c.E3.Symbol)
	}
	q3, ok := quotes[c.E3.Symbol]
	if !ok {
		return 0, fmt.Errorf("no price %s", c.E3.Symbol)
	}
	amt, err = applyEdge(amt, c.E3, r3, q3, cfg)
	if err != nil {
		return 0, fmt.Errorf("leg3 %s: %v", c.E3.Symbol, err)
	}

	return amt, nil
}

func findBestUSDT3Cycle(symbols [3]string, rulesMap map[string]SymbolRules) (Cycle, bool, string) {
	edges := []Edge{}
	for _, sym := range symbols {
		r, ok := rulesMap[sym]
		if !ok {
			return Cycle{}, false, "missing rules"
		}
		if !r.SpotAllowed {
			return Cycle{}, false, "spot not allowed"
		}
		if !r.MarketAllowed {
			return Cycle{}, false, "market not allowed"
		}
		if r.Base == "" || r.Quote == "" {
			return Cycle{}, false, "base/quote empty"
		}
		edges = append(edges,
			Edge{Symbol: sym, From: r.Quote, To: r.Base, Side: "BUY"},
			Edge{Symbol: sym, From: r.Base, To: r.Quote, Side: "SELL"},
		)
	}

	type node struct {
		asset   string
		usedSym map[string]bool
		path    []Edge
	}

	stack := []node{{
		asset:   "USDT",
		usedSym: map[string]bool{},
		path:    []Edge{},
	}}

	for len(stack) > 0 {
		n := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if len(n.path) == 3 {
			if n.asset == "USDT" {
				return Cycle{E1: n.path[0], E2: n.path[1], E3: n.path[2]}, true, "ok"
			}
			continue
		}

		for _, e := range edges {
			if n.usedSym[e.Symbol] {
				continue
			}
			if e.From != n.asset {
				continue
			}
			used := make(map[string]bool, len(n.usedSym)+1)
			for k, v := range n.usedSym {
				used[k] = v
			}
			used[e.Symbol] = true

			path := append(append([]Edge{}, n.path...), e)
			stack = append(stack, node{asset: e.To, usedSym: used, path: path})
		}
	}

	return Cycle{}, false, "no USDT 3-leg cycle"
}

func minStartByBinarySearch(c Cycle, rulesMap map[string]SymbolRules, quotes map[string]Quote, cfg SimCfg) (float64, string) {
	low := 0.0
	high := 1.0
	ok := false
	var lastErr error

	for i := 0; i < 40; i++ {
		_, err := simulateCycle(high, c, rulesMap, quotes, cfg)
		if err == nil {
			ok = true
			break
		}
		lastErr = err
		high *= 2
		if high > 1_000_000 {
			break
		}
	}
	if !ok {
		if lastErr != nil {
			return 0, lastErr.Error()
		}
		return 0, "not feasible"
	}

	for i := 0; i < 60; i++ {
		mid := (low + high) / 2
		if mid <= 0 {
			low = mid
			continue
		}
		_, err := simulateCycle(mid, c, rulesMap, quotes, cfg)
		if err == nil {
			high = mid
		} else {
			low = mid
		}
	}

	return high, "ok"
}

// ======================
// CSV I/O
// ======================

type InRow struct {
	Base1, Quote1 string
	Base2, Quote2 string
	Base3, Quote3 string
}

func readTrianglesCSV(path string) ([]InRow, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.FieldsPerRecord = -1

	head, err := r.Read()
	if err != nil {
		return nil, err
	}
	if len(head) < 6 {
		return nil, fmt.Errorf("bad header: %v", head)
	}

	var out []InRow
	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if len(rec) < 6 {
			continue
		}
		out = append(out, InRow{
			Base1:  strings.TrimSpace(rec[0]),
			Quote1: strings.TrimSpace(rec[1]),
			Base2:  strings.TrimSpace(rec[2]),
			Quote2: strings.TrimSpace(rec[3]),
			Base3:  strings.TrimSpace(rec[4]),
			Quote3: strings.TrimSpace(rec[5]),
		})
	}
	return out, nil
}

// ======================
// Main
// ======================

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	baseURL := envString("MEXC_BASE_URL", "https://api.mexc.com")
	inFile := envString("TRIANGLES_FILE", "triangles_markets.csv")
	outFile := envString("TRIANGLES_ENRICHED_FILE", "triangles_markets_enriched.csv")

	tradeAmountUSDT := envFloat("TRADE_AMOUNT_USDT", 10)
	feePct := envFloat("FEE_PCT", 0.04)
	sellSafety := envFloat("SELL_SAFETY", 0.995)

	log.Printf("input=%s output=%s TRADE_AMOUNT_USDT=%.6f FEE_PCT=%.6f SELL_SAFETY=%.6f",
		inFile, outFile, tradeAmountUSDT, feePct, sellSafety)

	rows, err := readTrianglesCSV(inFile)
	if err != nil {
		log.Fatalf("read triangles csv: %v", err)
	}
	log.Printf("triangles rows: %d", len(rows))

	info, err := fetchExchangeInfo(baseURL)
	if err != nil {
		log.Fatalf("fetch exchangeInfo: %v", err)
	}
	rulesMap := buildRules(info)
	log.Printf("rules loaded: %d symbols", len(rulesMap))

	quotes, err := fetchBookTickerAll(baseURL)
	if err != nil {
		log.Fatalf("fetch bookTicker: %v", err)
	}
	log.Printf("bookTicker loaded: %d symbols", len(quotes))

	f, err := os.Create(outFile)
	if err != nil {
		log.Fatalf("create out: %v", err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	// Пишем только проходящие треугольники, но хедер — полный
	header := []string{
		"base1", "quote1", "symbol1", "step1", "minQty1", "minNotional1",
		"spot1", "market1", "qoq1",
		"basePrec1", "quotePrec1", "quotePrecMkt1",
		"rawIsSpot1",

		"base2", "quote2", "symbol2", "step2", "minQty2", "minNotional2",
		"spot2", "market2", "qoq2",
		"basePrec2", "quotePrec2", "quotePrecMkt2",
		"rawIsSpot2",

		"base3", "quote3", "symbol3", "step3", "minQty3", "minNotional3",
		"spot3", "market3", "qoq3",
		"basePrec3", "quotePrec3", "quotePrecMkt3",
		"rawIsSpot3",

		"cycle_leg1", "cycle_leg2", "cycle_leg3",
		"min_start_usdt",
	}
	if err := w.Write(header); err != nil {
		log.Fatalf("write header: %v", err)
	}

	cfg := SimCfg{FeePct: feePct, SellSafety: sellSafety}

	okCount := 0
	skipped := 0

	stepFmt := func(x float64) string {
		if x == 0 {
			return "0"
		}
		dec := decimalsFromStep(x)
		if dec < 0 {
			return strconv.FormatFloat(x, 'f', 12, 64)
		}
		return strconv.FormatFloat(x, 'f', dec, 64)
	}
	f64 := func(x float64) string { return strconv.FormatFloat(x, 'f', 12, 64) }
	b2s := func(b bool) string {
		if b {
			return "1"
		}
		return "0"
	}

	for _, r := range rows {
		s1 := r.Base1 + r.Quote1
		s2 := r.Base2 + r.Quote2
		s3 := r.Base3 + r.Quote3

		rr1, ok1 := rulesMap[s1]
		rr2, ok2 := rulesMap[s2]
		rr3, ok3 := rulesMap[s3]

		if !ok1 || !ok2 || !ok3 {
			skipped++
			continue
		}
		// нужны цены
		if _, ok := quotes[s1]; !ok {
			skipped++
			continue
		}
		if _, ok := quotes[s2]; !ok {
			skipped++
			continue
		}
		if _, ok := quotes[s3]; !ok {
			skipped++
			continue
		}

		// структурный цикл USDT→...→USDT
		symbols := [3]string{s1, s2, s3}
		cycle, ok, _ := findBestUSDT3Cycle(symbols, rulesMap)
		if !ok {
			skipped++
			continue
		}

		ms, why := minStartByBinarySearch(cycle, rulesMap, quotes, cfg)
		if why != "ok" || ms <= 0 {
			skipped++
			continue
		}

		// ФИЛЬТР: пишем только проходящие под вход
		if tradeAmountUSDT+1e-9 < ms {
			skipped++
			continue
		}

		cycleStr1 := fmt.Sprintf("%s:%s %s→%s", cycle.E1.Symbol, cycle.E1.Side, cycle.E1.From, cycle.E1.To)
		cycleStr2 := fmt.Sprintf("%s:%s %s→%s", cycle.E2.Symbol, cycle.E2.Side, cycle.E2.From, cycle.E2.To)
		cycleStr3 := fmt.Sprintf("%s:%s %s→%s", cycle.E3.Symbol, cycle.E3.Side, cycle.E3.From, cycle.E3.To)

		raw1 := ""
		raw2 := ""
		raw3 := ""
		if rr1.RawHasIsSpotField {
			raw1 = b2s(rr1.RawIsSpotTradingAllowed)
		}
		if rr2.RawHasIsSpotField {
			raw2 = b2s(rr2.RawIsSpotTradingAllowed)
		}
		if rr3.RawHasIsSpotField {
			raw3 = b2s(rr3.RawIsSpotTradingAllowed)
		}

		out := []string{
			r.Base1, r.Quote1, s1, stepFmt(rr1.BaseStep), f64(rr1.MinQty), f64(rr1.MinNotional),
			b2s(rr1.SpotAllowed), b2s(rr1.MarketAllowed), b2s(rr1.QuoteOrderQtyMarketAllowed),
			strconv.Itoa(rr1.BasePrecision), strconv.Itoa(rr1.QuotePrecision), strconv.Itoa(rr1.QuotePrecisionMarket),
			raw1,

			r.Base2, r.Quote2, s2, stepFmt(rr2.BaseStep), f64(rr2.MinQty), f64(rr2.MinNotional),
			b2s(rr2.SpotAllowed), b2s(rr2.MarketAllowed), b2s(rr2.QuoteOrderQtyMarketAllowed),
			strconv.Itoa(rr2.BasePrecision), strconv.Itoa(rr2.QuotePrecision), strconv.Itoa(rr2.QuotePrecisionMarket),
			raw2,

			r.Base3, r.Quote3, s3, stepFmt(rr3.BaseStep), f64(rr3.MinQty), f64(rr3.MinNotional),
			b2s(rr3.SpotAllowed), b2s(rr3.MarketAllowed), b2s(rr3.QuoteOrderQtyMarketAllowed),
			strconv.Itoa(rr3.BasePrecision), strconv.Itoa(rr3.QuotePrecision), strconv.Itoa(rr3.QuotePrecisionMarket),
			raw3,

			cycleStr1, cycleStr2, cycleStr3,
			strconv.FormatFloat(ms, 'f', 6, 64),
		}

		if err := w.Write(out); err != nil {
			log.Fatalf("write out: %v", err)
		}
		okCount++
	}

	log.Printf("DONE: wrote ok triangles: %d (skipped=%d)", okCount, skipped)
	fmt.Println("Готово, файл:", outFile)
}


