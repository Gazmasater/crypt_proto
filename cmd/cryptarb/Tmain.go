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
	"strconv"
	"strings"
	"time"
)

const (
	baseURL   = "https://api.mexc.com"
	inCSV     = "triangles_markets.csv"
	outCSV    = "triangles_markets_ready.csv"
	httpTO    = 25 * time.Second
	safetyMul = 1.03 // запас на округления/микро-скольжение

	maxCapUSDT = 5000.0 // если требуется больше — считаем "непригодным"
)

type exchangeInfo struct {
	Symbols []symbolInfo `json:"symbols"`
}

type symbolInfo struct {
	Symbol                   string           `json:"symbol"`
	BaseAsset                string           `json:"baseAsset"`
	QuoteAsset               string           `json:"quoteAsset"`
	Status                   string           `json:"status"`
	IsSpotTradingAllowed     *bool            `json:"isSpotTradingAllowed"`
	QuoteOrderQtyMarketAllow *bool            `json:"quoteOrderQtyMarketAllowed"`
	OrderTypes               []string         `json:"orderTypes"`
	Filters                  []map[string]any `json:"filters"`
}

type bookTicker struct {
	Symbol   string `json:"symbol"`
	BidPrice string `json:"bidPrice"`
	AskPrice string `json:"askPrice"`
}

type rules struct {
	Symbol string
	Base   string
	Quote  string

	SpotAllowed bool
	HasMarket   bool

	StepSize    float64
	MinQty      float64
	MinNotional float64

	Bid float64
	Ask float64
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	c := &http.Client{Timeout: httpTO}

	info, err := fetchExchangeInfo(c)
	if err != nil {
		log.Fatalf("exchangeInfo: %v", err)
	}
	ticks, err := fetchBookTickers(c)
	if err != nil {
		log.Fatalf("bookTicker: %v", err)
	}

	price := map[string][2]float64{}
	for _, t := range ticks {
		bid, _ := strconv.ParseFloat(t.BidPrice, 64)
		ask, _ := strconv.ParseFloat(t.AskPrice, 64)
		price[strings.TrimSpace(t.Symbol)] = [2]float64{bid, ask}
	}

	bySymbol := map[string]rules{}
	for _, s := range info.Symbols {
		sym := strings.TrimSpace(s.Symbol)
		base := strings.TrimSpace(s.BaseAsset)
		quote := strings.TrimSpace(s.QuoteAsset)
		if sym == "" || base == "" || quote == "" {
			continue
		}
		if !tradableStatus(s.Status) {
			continue
		}

		r := rules{Symbol: sym, Base: base, Quote: quote}

		// spot allow
		if s.IsSpotTradingAllowed != nil {
			r.SpotAllowed = *s.IsSpotTradingAllowed
		} else {
			r.SpotAllowed = true
		}
		r.HasMarket = hasOrderType(s.OrderTypes, "MARKET")

		// filters
		for _, f := range s.Filters {
			ft := strings.ToUpper(strings.TrimSpace(asString(f["filterType"])))
			switch ft {
			case "LOT_SIZE", "MARKET_LOT_SIZE":
				r.StepSize = parseFloatAny(f["stepSize"])
				r.MinQty = parseFloatAny(f["minQty"])
			case "MIN_NOTIONAL", "NOTIONAL":
				if v := parseFloatAny(f["minNotional"]); v > 0 {
					r.MinNotional = v
				} else if v := parseFloatAny(f["notional"]); v > 0 {
					r.MinNotional = v
				} else if v := parseFloatAny(f["minQuote"]); v > 0 {
					r.MinNotional = v
				}
			}
		}

		if p, ok := price[sym]; ok {
			r.Bid, r.Ask = p[0], p[1]
		}

		bySymbol[sym] = r
	}

	in, err := os.Open(inCSV)
	if err != nil {
		log.Fatalf("open %s: %v", inCSV, err)
	}
	defer in.Close()

	out, err := os.Create(outCSV)
	if err != nil {
		log.Fatalf("create %s: %v", outCSV, err)
	}
	defer out.Close()

	cr := csv.NewReader(in)
	cw := csv.NewWriter(out)
	defer cw.Flush()

	header, err := cr.Read()
	if err != nil {
		log.Fatalf("read header: %v", err)
	}

	// ждём минимум 6 колонок как у тебя
	if len(header) < 6 {
		log.Fatalf("bad header: need >=6 cols, got %d", len(header))
	}

	// пишем новый header: добавляем min_start_usdt
	newHeader := append([]string{}, header...)
	newHeader = append(newHeader, "min_start_usdt", "reason")
	if err := cw.Write(newHeader); err != nil {
		log.Fatalf("write header: %v", err)
	}

	written := 0
	for {
		rec, err := cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("read csv: %v", err)
		}
		if len(rec) < 6 {
			continue
		}

		base1, quote1 := strings.TrimSpace(rec[0]), strings.TrimSpace(rec[1])
		base2, quote2 := strings.TrimSpace(rec[2]), strings.TrimSpace(rec[3])
		base3, quote3 := strings.TrimSpace(rec[4]), strings.TrimSpace(rec[5])

		// из твоего CSV восстанавливаем 3 символа
		sym1 := base1 + quote1
		sym2 := base2 + quote2
		sym3 := base3 + quote3

		r1, ok1 := bySymbol[sym1]
		r2, ok2 := bySymbol[sym2]
		r3, ok3 := bySymbol[sym3]

		minStart := 0.0
		reason := ""

		if !ok1 || !ok2 || !ok3 {
			reason = "missing rules for one of symbols"
		} else if !good(r1) || !good(r2) || !good(r3) {
			reason = "spot/market/price not ok"
		} else {
			// Пытаемся посчитать minStart именно для цикла USDT→X→Y→USDT.
			// Берём множество ассетов треугольника и пробуем подобрать X,Y.
			assets := uniq3(base1, quote1, base2, quote2, base3, quote3)
			if !contains(assets, "USDT") {
				reason = "triangle has no USDT"
			} else {
				// выбираем X и Y среди остальных 2 ассетов, если их ровно 3 уникальных
				if len(assets) != 3 {
					reason = "not 3 unique assets"
				} else {
					var x, y string
					for _, a := range assets {
						if a != "USDT" {
							if x == "" {
								x = a
							} else {
								y = a
							}
						}
					}

					// строим 3 ноги в нужном направлении (USDT->x, x->y, y->USDT),
					// если не получается — пробуем x<->y местами.
					min1, ok, why := minStartForCycle("USDT", x, y, r1, r2, r3)
					if ok {
						minStart = min1
					} else {
						min2, ok2, why2 := minStartForCycle("USDT", y, x, r1, r2, r3)
						if ok2 {
							minStart = min2
						} else {
							reason = "cannot build USDT cycle: " + why + " | " + why2
						}
					}
				}
			}
		}

		newRec := append([]string{}, rec...)
		newRec = append(newRec, fmt.Sprintf("%.6f", minStart), reason)

		if err := cw.Write(newRec); err != nil {
			log.Fatalf("write rec: %v", err)
		}
		written++
	}

	log.Printf("done: written=%d -> %s", written, outCSV)
	fmt.Println("Готово, файл:", outCSV)
}

func good(r rules) bool {
	if !r.SpotAllowed || !r.HasMarket {
		return false
	}
	if r.Bid <= 0 || r.Ask <= 0 {
		return false
	}
	return true
}

// minStartForCycle считает минимальный старт USDT, чтобы пройти 3 ноги MARKET.
// Мы знаем только 3 символа (r1,r2,r3), поэтому подбираем какой символ соответствует каждой ноге по (from,to).
func minStartForCycle(usdt, x, y string, r1, r2, r3 rules) (float64, bool, string) {
	legs := []struct {
		from string
		to   string
	}{
		{usdt, x},
		{x, y},
		{y, usdt},
	}

	syms := []rules{r1, r2, r3}
	used := make([]bool, 3)

	// подбираем соответствие "нога -> символ" без повторов
	var picked [3]rules
	for i, leg := range legs {
		found := false
		for j := 0; j < 3; j++ {
			if used[j] {
				continue
			}
			if canConvert(syms[j], leg.from, leg.to) {
				picked[i] = syms[j]
				used[j] = true
				found = true
				break
			}
		}
		if !found {
			return 0, false, fmt.Sprintf("no market for leg %s->%s", leg.from, leg.to)
		}
	}

	// Теперь считаем минимальный старт так, чтобы каждая нога прошла с округлениями.
	// Делаем итеративно: подбираем старт, прогоняем, если где-то fail — увеличиваем требование.
	start := 0.01
	for start <= maxCapUSDT {
		amt := start

		ok := true
		for i, leg := range legs {
			amt, ok = apply(picked[i], leg.from, leg.to, amt)
			if !ok {
				break
			}
		}
		if ok {
			return round2(start * safetyMul), true, ""
		}
		start += 0.5
	}
	return 0, false, "need start > cap"
}

func apply(m rules, from, to string, amount float64) (float64, bool) {
	// если from==quote & to==base => BUY: qty=amount/ask
	if from == m.Quote && to == m.Base {
		qtyRaw := amount / m.Ask
		qty := floorToStep(qtyRaw, m.StepSize)
		if qty <= 0 {
			return 0, false
		}
		if m.MinQty > 0 && qty+1e-18 < m.MinQty {
			return 0, false
		}
		notional := qty * m.Ask
		if m.MinNotional > 0 && notional+1e-12 < m.MinNotional {
			return 0, false
		}
		return qty, true
	}

	// если from==base & to==quote => SELL: quote=qty*bid
	if from == m.Base && to == m.Quote {
		qty := floorToStep(amount, m.StepSize)
		if qty <= 0 {
			return 0, false
		}
		if m.MinQty > 0 && qty+1e-18 < m.MinQty {
			return 0, false
		}
		notional := qty * m.Bid
		if m.MinNotional > 0 && notional+1e-12 < m.MinNotional {
			return 0, false
		}
		return notional, true
	}

	return 0, false
}

func canConvert(m rules, from, to string) bool {
	return (from == m.Quote && to == m.Base) || (from == m.Base && to == m.Quote)
}

func floorToStep(qty, step float64) float64 {
	if qty <= 0 {
		return 0
	}
	if step <= 0 {
		return qty
	}
	n := math.Floor(qty/step) * step
	if n < step {
		return 0
	}
	return n
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}

func fetchExchangeInfo(c *http.Client) (*exchangeInfo, error) {
	resp, err := c.Get(baseURL + "/api/v3/exchangeInfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("exchangeInfo status=%d body=%s", resp.StatusCode, string(b))
	}
	var info exchangeInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, err
	}
	return &info, nil
}

func fetchBookTickers(c *http.Client) ([]bookTicker, error) {
	resp, err := c.Get(baseURL + "/api/v3/ticker/bookTicker")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("bookTicker status=%d body=%s", resp.StatusCode, string(b))
	}
	var arr []bookTicker
	if err := json.NewDecoder(resp.Body).Decode(&arr); err != nil {
		return nil, err
	}
	return arr, nil
}

func tradableStatus(s string) bool {
	s = strings.ToUpper(strings.TrimSpace(s))
	if s == "" {
		return true
	}
	return s == "ENABLED" || s == "TRADING" || s == "1"
}

func hasOrderType(list []string, want string) bool {
	want = strings.ToUpper(strings.TrimSpace(want))
	if len(list) == 0 {
		return true
	}
	for _, it := range list {
		if strings.ToUpper(strings.TrimSpace(it)) == want {
			return true
		}
	}
	return false
}

func asString(v any) string {
	switch x := v.(type) {
	case string:
		return x
	case float64:
		return strconv.FormatFloat(x, 'f', -1, 64)
	default:
		return ""
	}
}

func parseFloatAny(v any) float64 {
	switch x := v.(type) {
	case string:
		f, _ := strconv.ParseFloat(strings.TrimSpace(x), 64)
		return f
	case float64:
		return x
	default:
		return 0
	}
}

func uniq3(vals ...string) []string {
	m := map[string]struct{}{}
	for _, v := range vals {
		v = strings.TrimSpace(v)
		if v != "" {
			m[v] = struct{}{}
		}
	}
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}

func contains(list []string, x string) bool {
	for _, v := range list {
		if v == x {
			return true
		}
	}
	return false
}
