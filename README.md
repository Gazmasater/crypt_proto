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
	Symbol                 string       `json:"symbol"`
	BaseAsset              string       `json:"baseAsset"`
	QuoteAsset             string       `json:"quoteAsset"`
	Status                 string       `json:"status"`
	IsSpotTradingAllowed   *bool        `json:"isSpotTradingAllowed"`
	QuoteOrderQtyMarket    *bool        `json:"quoteOrderQtyMarketAllowed"`
	OrderTypes             []string     `json:"orderTypes"`
	BaseAssetPrecision     any          `json:"baseAssetPrecision"`
	QuoteAssetPrecision    any          `json:"quoteAssetPrecision"`
	Filters                []mexcFilter `json:"filters"`
	Permissions            []string     `json:"permissions"`
	IsMarginTradingAllowed *bool        `json:"isMarginTradingAllowed"`
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

	SpotAllowed   bool
	MarketAllowed bool

	QuoteOrderQtyMarketAllowed bool

	BaseStep       float64
	MinQty         float64
	MinNotional    float64 // assumed in QUOTE units (USDT/USDC/etc)
	QuotePrecision int
}

type Quote struct {
	Bid float64
	Ask float64
}

// ======================
// Utils: parsing
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

func decimalsFromStep(step float64) int {
	if step <= 0 {
		return -1
	}
	// step like 0.001 => 3 decimals
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

		// spot allowed
		if s.IsSpotTradingAllowed != nil {
			r.SpotAllowed = *s.IsSpotTradingAllowed
		} else {
			// fallback: best-effort
			r.SpotAllowed = true
		}

		// market allowed (by orderTypes)
		r.MarketAllowed = hasMarket(s.OrderTypes)

		// quoteOrderQtyMarketAllowed
		if s.QuoteOrderQtyMarket != nil {
			r.QuoteOrderQtyMarketAllowed = *s.QuoteOrderQtyMarket
		}

		// quote precision
		if v, ok := asInt(s.QuoteAssetPrecision); ok {
			r.QuotePrecision = v
		} else {
			r.QuotePrecision = 2
		}

		// filters: prefer MARKET_LOT_SIZE for market orders; else LOT_SIZE
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

		// choose market lot size if present
		if mktStep > 0 || mktMin > 0 {
			r.BaseStep = mktStep
			r.MinQty = mktMin
		} else {
			r.BaseStep = lotStep
			r.MinQty = lotMin
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
	c := &http.Client{Timeout: 15 * time.Second}
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
		if rules.QuoteOrderQtyMarketAllowed {
			amtQuote := truncToDecimals(amount, rules.QuotePrecision)
			if amtQuote <= 0 {
				return 0, fmt.Errorf("buy amount<=0 after quote precision")
			}
			if minNotional > 0 && amtQuote+1e-12 < minNotional {
				return 0, fmt.Errorf("buy minNotional: need>=%.12f have=%.12f", minNotional, amtQuote)
			}
			gotBase := amtQuote / q.Ask
			gotBase *= (1.0 - fee)
			if minQty > 0 && gotBase+1e-12 < minQty {
				return 0, fmt.Errorf("buy minQty: need>=%.12f got=%.12f", minQty, gotBase)
			}
			return gotBase, nil
		}

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
	var err error
	amt, err = applyEdge(amt, c.E1, r1, q1, cfg)
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

	best := Cycle{}
	found := false

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
				best = Cycle{E1: n.path[0], E2: n.path[1], E3: n.path[2]}
				found = true
				break
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

	if !found {
		return Cycle{}, false, "no USDT 3-leg cycle"
	}
	return best, true, "ok"
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

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	baseURL := envString("MEXC_BASE_URL", "https://api.mexc.com")
	inFile := envString("TRIANGLES_FILE", "triangles_markets.csv")

	// оставляем название как ты сказал
	outFile := envString("TRIANGLES_ENRICHED_FILE", "triangles_markets_enriched.csv")

	// доп. файл с мусором/причинами (чтобы не терять диагностику)
	excludedFile := envString("TRIANGLES_EXCLUDED_FILE", "triangles_markets_excluded.csv")

	tradeAmountUSDT := envFloat("TRADE_AMOUNT_USDT", 10)
	feePct := envFloat("FEE_PCT", 0.04)
	sellSafety := envFloat("SELL_SAFETY", 0.995)

	log.Printf("input=%s output=%s excluded=%s TRADE_AMOUNT_USDT=%.6f FEE_PCT=%.6f SELL_SAFETY=%.6f",
		inFile, outFile, excludedFile, tradeAmountUSDT, feePct, sellSafety)

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

	// ===== writers: OK + EXCLUDED =====

	fOK, err := os.Create(outFile)
	if err != nil {
		log.Fatalf("create out: %v", err)
	}
	defer fOK.Close()

	wOK := csv.NewWriter(fOK)
	defer wOK.Flush()

	fBad, err := os.Create(excludedFile)
	if err != nil {
		log.Fatalf("create excluded out: %v", err)
	}
	defer fBad.Close()

	wBad := csv.NewWriter(fBad)
	defer wBad.Flush()

	header := []string{
		"base1", "quote1", "symbol1", "step1", "minQty1", "minNotional1", "spot1", "market1", "qoq1", "quotePrec1",
		"base2", "quote2", "symbol2", "step2", "minQty2", "minNotional2", "spot2", "market2", "qoq2", "quotePrec2",
		"base3", "quote3", "symbol3", "step3", "minQty3", "minNotional3", "spot3", "market3", "qoq3", "quotePrec3",
		"cycle_leg1", "cycle_leg2", "cycle_leg3",
		"min_start_usdt", "trade_amount_ok", "reason",
	}
	if err := wOK.Write(header); err != nil {
		log.Fatalf("write header ok: %v", err)
	}
	if err := wBad.Write(header); err != nil {
		log.Fatalf("write header bad: %v", err)
	}

	cfg := SimCfg{FeePct: feePct, SellSafety: sellSafety}

	okCount := 0
	badCount := 0

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
	f64 := func(x float64) string {
		return strconv.FormatFloat(x, 'f', 12, 64)
	}

	for _, r := range rows {
		s1 := r.Base1 + r.Quote1
		s2 := r.Base2 + r.Quote2
		s3 := r.Base3 + r.Quote3

		symbols := [3]string{s1, s2, s3}

		rr1, ok1 := rulesMap[s1]
		rr2, ok2 := rulesMap[s2]
		rr3, ok3 := rulesMap[s3]

		reason := "ok"
		if !ok1 || !ok2 || !ok3 {
			reason = "missing rules"
		} else if !rr1.SpotAllowed || !rr2.SpotAllowed || !rr3.SpotAllowed {
			reason = "spot not allowed"
		} else if !rr1.MarketAllowed || !rr2.MarketAllowed || !rr3.MarketAllowed {
			reason = "market not allowed"
		} else if _, ok := quotes[s1]; !ok {
			reason = "price not ok"
		} else if _, ok := quotes[s2]; !ok {
			reason = "price not ok"
		} else if _, ok := quotes[s3]; !ok {
			reason = "price not ok"
		}

		cycleStr1, cycleStr2, cycleStr3 := "", "", ""
		minStart := 0.0
		tradeOK := "0"

		if reason == "ok" {
			cycle, ok, why := findBestUSDT3Cycle(symbols, rulesMap)
			if !ok {
				reason = why
			} else {
				cycleStr1 = fmt.Sprintf("%s:%s %s→%s", cycle.E1.Symbol, cycle.E1.Side, cycle.E1.From, cycle.E1.To)
				cycleStr2 = fmt.Sprintf("%s:%s %s→%s", cycle.E2.Symbol, cycle.E2.Side, cycle.E2.From, cycle.E2.To)
				cycleStr3 = fmt.Sprintf("%s:%s %s→%s", cycle.E3.Symbol, cycle.E3.Side, cycle.E3.From, cycle.E3.To)

				ms, why2 := minStartByBinarySearch(cycle, rulesMap, quotes, cfg)
				minStart = ms
				if why2 != "ok" {
					reason = why2
				} else {
					if minStart > 0 && tradeAmountUSDT+1e-9 >= minStart {
						tradeOK = "1"
					} else {
						reason = fmt.Sprintf("needUSDT>=%.6f", minStart)
					}
				}
			}
		}

		spot1 := "0"
		market1 := "0"
		qoq1 := "0"
		qp1 := "0"
		if ok1 {
			if rr1.SpotAllowed {
				spot1 = "1"
			}
			if rr1.MarketAllowed {
				market1 = "1"
			}
			if rr1.QuoteOrderQtyMarketAllowed {
				qoq1 = "1"
			}
			qp1 = strconv.Itoa(rr1.QuotePrecision)
		}

		spot2 := "0"
		market2 := "0"
		qoq2 := "0"
		qp2 := "0"
		if ok2 {
			if rr2.SpotAllowed {
				spot2 = "1"
			}
			if rr2.MarketAllowed {
				market2 = "1"
			}
			if rr2.QuoteOrderQtyMarketAllowed {
				qoq2 = "1"
			}
			qp2 = strconv.Itoa(rr2.QuotePrecision)
		}

		spot3 := "0"
		market3 := "0"
		qoq3 := "0"
		qp3 := "0"
		if ok3 {
			if rr3.SpotAllowed {
				spot3 = "1"
			}
			if rr3.MarketAllowed {
				market3 = "1"
			}
			if rr3.QuoteOrderQtyMarketAllowed {
				qoq3 = "1"
			}
			qp3 = strconv.Itoa(rr3.QuotePrecision)
		}

		out := []string{
			r.Base1, r.Quote1, s1, stepFmt(rr1.BaseStep), f64(rr1.MinQty), f64(rr1.MinNotional), spot1, market1, qoq1, qp1,
			r.Base2, r.Quote2, s2, stepFmt(rr2.BaseStep), f64(rr2.MinQty), f64(rr2.MinNotional), spot2, market2, qoq2, qp2,
			r.Base3, r.Quote3, s3, stepFmt(rr3.BaseStep), f64(rr3.MinQty), f64(rr3.MinNotional), spot3, market3, qoq3, qp3,
			cycleStr1, cycleStr2, cycleStr3,
			strconv.FormatFloat(minStart, 'f', 6, 64),
			tradeOK,
			reason,
		}

		// ===== ВАЖНО: мусор не пишем в enriched =====
		if tradeOK == "1" {
			if err := wOK.Write(out); err != nil {
				log.Fatalf("write ok out: %v", err)
			}
			okCount++
		} else {
			if err := wBad.Write(out); err != nil {
				log.Fatalf("write bad out: %v", err)
			}
			badCount++
		}
	}

	log.Printf("DONE: ok triangles written: %d, excluded written: %d", okCount, badCount)
	fmt.Println("Готово, OK файл:", outFile)
	fmt.Println("Готово, EXCLUDED файл:", excludedFile)
}




az358@gaz358-BOD-WXX9:~/myprog/crypt_proto/cmd/cryptarb$ go run .
2025/12/17 18:28:24.375529 input=triangles_markets.csv output=triangles_markets_enriched.csv excluded=triangles_markets_excluded.csv TRADE_AMOUNT_USDT=100.000000 FEE_PCT=0.040000 SELL_SAFETY=0.995000
2025/12/17 18:28:24.375755 triangles rows: 415
2025/12/17 18:28:25.007060 rules loaded: 2538 symbols
2025/12/17 18:28:25.124686 bookTicker loaded: 2533 symbols
2025/12/17 18:28:25.128289 DONE: ok triangles written: 44, excluded written: 371
Готово, OK файл: triangles_markets_enriched.csv
Готово, EXCLUDED файл: triangles_markets_excluded.csv






