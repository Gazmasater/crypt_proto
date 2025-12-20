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




export TRADE_AMOUNT_USDT=100
export FEE_PCT=0.04
export SELL_SAFETY=0.995

export TRIANGLES_FILE=triangles_markets.csv
export TRIANGLES_ENRICHED_FILE=triangles_markets_enriched.csv

go run ./cmd/triangles_enrich_mexc



# ===== INPUT / OUTPUT =====
INPUT_CSV=triangles_usdt_routes_okx.csv
OUTPUT_CSV=triangles_usdt_routes_market_okx.csv
BLACKLIST_FILE=blacklist_symbols_okx.txt

# ===== OKX API =====
OKX_INSTRUMENTS_URL=https://www.okx.com/api/v5/public/instruments?instType=SPOT
OKX_BOOKS_URL=https://www.okx.com/api/v5/market/books

# ===== RUNTIME CFG =====
HTTP_TIMEOUT=25s

START_AMT=25
START_CCY=USDT

# комиссия на 1 ногу (taker), в процентах
FEE_PCT=0.1

# минимальная прибыль по кругу (после комиссий), в процентах
MIN_PROFIT_PCT=0.3

# 1 = писать все строки + fail_reason (для дебага)
# 0 = писать только прошедшие фильтры
KEEP_FAILED=0





package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	DefaultEnvFile        = ".env"
	DefaultOutCSV         = "triangles_routes_params_okx.csv"
	DefaultInstrumentsURL = "https://www.okx.com/api/v5/public/instruments?instType=SPOT"
	DefaultHTTPTimeout    = 25 * time.Second
	DefaultStartCcy       = "USDT"
)

// ===== OKX response =====

type okxInstrumentsResp struct {
	Code string          `json:"code"`
	Msg  string          `json:"msg"`
	Data []okxInstrument `json:"data"`
}

type okxInstrument struct {
	InstID   string `json:"instId"`
	InstType string `json:"instType"`
	BaseCcy  string `json:"baseCcy"`
	QuoteCcy string `json:"quoteCcy"`
	State    string `json:"state"`

	LotSz  string `json:"lotSz"`  // qty step (base)
	MinSz  string `json:"minSz"`  // min qty (base)
	TickSz string `json:"tickSz"` // price step

	MaxMktAmt string `json:"maxMktAmt"` // limit in quote (often for BUY)
	MaxMktSz  string `json:"maxMktSz"`  // limit in base  (often for SELL)
}

// ===== .env loader (doesn't override exported vars) =====

func loadDotEnvIfExists(path string) {
	path = strings.TrimSpace(path)
	if path == "" {
		return
	}
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "export ") {
			line = strings.TrimSpace(strings.TrimPrefix(line, "export "))
		}
		k, v, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key := strings.TrimSpace(k)
		val := strings.TrimSpace(v)
		if key == "" {
			continue
		}
		// strip quotes
		if len(val) >= 2 {
			if (val[0] == '"' && val[len(val)-1] == '"') || (val[0] == '\'' && val[len(val)-1] == '\'') {
				val = val[1 : len(val)-1]
			}
		}
		if _, exists := os.LookupEnv(key); exists {
			continue
		}
		_ = os.Setenv(key, val)
	}
}

func getenv(key, def string) string {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return def
	}
	return v
}

func getenvDuration(key string, def time.Duration) time.Duration {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return def
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		log.Printf("WARN: bad %s=%q, using default %s", key, v, def)
		return def
	}
	return d
}

// ===== numeric helpers =====

func parsePositive(s string) (float64, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, false
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, false
	}
	if f <= 0 {
		return 0, false
	}
	return f, true
}

func parseNonNegative(s string) (float64, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, false
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, false
	}
	if f < 0 {
		return 0, false
	}
	return f, true
}

func hasPositiveNumber(s string) bool {
	_, ok := parsePositive(s)
	return ok
}

func b01(b bool) string {
	if b {
		return "1"
	}
	return "0"
}

func decimalsFromStep(step string) int {
	step = strings.TrimSpace(step)
	if step == "" {
		return -1
	}
	if !strings.Contains(step, ".") {
		return 0
	}
	parts := strings.SplitN(step, ".", 2)
	frac := strings.TrimRight(parts[1], "0")
	return len(frac)
}

// normalize numeric string (strip spaces, parse->format without trailing zeros where possible)
func normNumString(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return ""
	}
	// Use 'g' to avoid scientific for typical OKX steps, but keep precision.
	// 16 is safe for float64 round-tripping for these sizes.
	out := strconv.FormatFloat(f, 'g', 16, 64)
	return out
}

// ===== filters =====
// criteria from you:
// instType==SPOT, state==live, lotSz/minSz/tickSz not empty, and (maxMktAmt or maxMktSz) present (positive)
func marketOkOKX(s okxInstrument) bool {
	if s.InstType != "" && !strings.EqualFold(strings.TrimSpace(s.InstType), "SPOT") {
		return false
	}
	if !strings.EqualFold(strings.TrimSpace(s.State), "live") {
		return false
	}
	if strings.TrimSpace(s.InstID) == "" || strings.TrimSpace(s.BaseCcy) == "" || strings.TrimSpace(s.QuoteCcy) == "" {
		return false
	}
	if strings.TrimSpace(s.LotSz) == "" || strings.TrimSpace(s.MinSz) == "" || strings.TrimSpace(s.TickSz) == "" {
		return false
	}
	// IMPORTANT: require at least one positive limit field
	if !hasPositiveNumber(s.MaxMktAmt) && !hasPositiveNumber(s.MaxMktSz) {
		return false
	}
	return true
}

// ===== network =====

func loadRulesOKX(url string, timeout time.Duration) (map[string]okxInstrument, error) {
	client := &http.Client{Timeout: timeout}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "cryptarb/1.0 (+okx instruments)")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return nil, fmt.Errorf("OKX instruments http %d: %s", resp.StatusCode, strings.TrimSpace(string(b)))
	}

	var r okxInstrumentsResp
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}
	if strings.TrimSpace(r.Code) != "" && strings.TrimSpace(r.Code) != "0" {
		return nil, fmt.Errorf("OKX instruments code=%s msg=%s", strings.TrimSpace(r.Code), strings.TrimSpace(r.Msg))
	}

	out := make(map[string]okxInstrument, len(r.Data))
	for _, it := range r.Data {
		id := strings.TrimSpace(it.InstID)
		if id == "" {
			continue
		}
		out[id] = it
	}
	return out, nil
}

// ===== triangle build using X<->Y explicitly =====

type Leg struct {
	From, To  string
	InstID    string
	Action    string // BUY / SELL
	PriceSide string // ASK / BID

	// execution helpers
	MktLimitKind  string // QUOTE for BUY, BASE for SELL
	HasMaxMktAmt  string // "1" or "0"
	HasMaxMktSz   string // "1" or "0"
	MktLimitField string // "maxMktAmt" or "maxMktSz" or ""
	MktLimitValue string // numeric string (the chosen limit) or ""

	// normalized numeric fields (so runtime doesn't parse strings)
	QtyStep   string // lotSz as normalized float string
	MinQty    string // minSz as normalized float string
	PriceStep string // tickSz as normalized float string

	BaseCcy   string
	QuoteCcy  string
	LotSz     string
	MinSzRaw  string
	TickSz    string
	MaxMktAmt string
	MaxMktSz  string
	QtyDP     int
	PriceDP   int
}

// choose the limit to apply for this leg, returning:
// kind: QUOTE/BASE, field: maxMktAmt/maxMktSz, value string (or "")
func chooseMktLimit(action string, inst okxInstrument) (kind, field, value string) {
	amt, okAmt := parsePositive(inst.MaxMktAmt)
	sz, okSz := parsePositive(inst.MaxMktSz)

	amtStr := ""
	szStr := ""
	if okAmt {
		_ = amt
		amtStr = strings.TrimSpace(inst.MaxMktAmt)
	}
	if okSz {
		_ = sz
		szStr = strings.TrimSpace(inst.MaxMktSz)
	}

	switch strings.ToUpper(strings.TrimSpace(action)) {
	case "BUY":
		if okAmt {
			return "QUOTE", "maxMktAmt", amtStr
		}
		if okSz {
			return "QUOTE", "maxMktSz", szStr // fallback
		}
		return "QUOTE", "", ""
	case "SELL":
		if okSz {
			return "BASE", "maxMktSz", szStr
		}
		if okAmt {
			return "BASE", "maxMktAmt", amtStr // fallback
		}
		return "BASE", "", ""
	default:
		return "", "", ""
	}
}

// BUY: quote -> base (ASK) => kind QUOTE
// SELL: base  -> quote (BID) => kind BASE
func makeLeg(from, to string, inst okxInstrument) (Leg, bool) {
	from = strings.ToUpper(strings.TrimSpace(from))
	to = strings.ToUpper(strings.TrimSpace(to))

	base := strings.ToUpper(strings.TrimSpace(inst.BaseCcy))
	quote := strings.ToUpper(strings.TrimSpace(inst.QuoteCcy))

	hasAmt := hasPositiveNumber(inst.MaxMktAmt)
	hasSz := hasPositiveNumber(inst.MaxMktSz)

	// normalize steps for runtime
	qtyStep := normNumString(inst.LotSz)
	minQty := normNumString(inst.MinSz)
	priceStep := normNumString(inst.TickSz)

	// if steps can't be parsed (shouldn't happen because filters require not empty),
	// keep raw string for debugging, but set normalized empty
	if qtyStep == "" {
		qtyStep = strings.TrimSpace(inst.LotSz)
	}
	if minQty == "" {
		minQty = strings.TrimSpace(inst.MinSz)
	}
	if priceStep == "" {
		priceStep = strings.TrimSpace(inst.TickSz)
	}

	if from == quote && to == base {
		kind, field, val := chooseMktLimit("BUY", inst)

		return Leg{
			From: from, To: to, InstID: inst.InstID,
			Action: "BUY", PriceSide: "ASK",

			MktLimitKind:  kind,
			HasMaxMktAmt:  b01(hasAmt),
			HasMaxMktSz:   b01(hasSz),
			MktLimitField: field,
			MktLimitValue: val,

			QtyStep: qtyStep,
			MinQty:  minQty,
			PriceStep: priceStep,

			BaseCcy: base, QuoteCcy: quote,
			LotSz: inst.LotSz, MinSzRaw: inst.MinSz, TickSz: inst.TickSz,
			MaxMktAmt: inst.MaxMktAmt, MaxMktSz: inst.MaxMktSz,
			QtyDP: decimalsFromStep(inst.LotSz),
			PriceDP: decimalsFromStep(inst.TickSz),
		}, true
	}

	if from == base && to == quote {
		kind, field, val := chooseMktLimit("SELL", inst)

		return Leg{
			From: from, To: to, InstID: inst.InstID,
			Action: "SELL", PriceSide: "BID",

			MktLimitKind:  kind,
			HasMaxMktAmt:  b01(hasAmt),
			HasMaxMktSz:   b01(hasSz),
			MktLimitField: field,
			MktLimitValue: val,

			QtyStep: qtyStep,
			MinQty:  minQty,
			PriceStep: priceStep,

			BaseCcy: base, QuoteCcy: quote,
			LotSz: inst.LotSz, MinSzRaw: inst.MinSz, TickSz: inst.TickSz,
			MaxMktAmt: inst.MaxMktAmt, MaxMktSz: inst.MaxMktSz,
			QtyDP: decimalsFromStep(inst.LotSz),
			PriceDP: decimalsFromStep(inst.TickSz),
		}, true
	}

	return Leg{}, false
}

func pairKey(a, b string) string {
	a = strings.ToUpper(strings.TrimSpace(a))
	b = strings.ToUpper(strings.TrimSpace(b))
	if a < b {
		return a + "|" + b
	}
	return b + "|" + a
}

func ensureDirForFile(path string) error {
	dir := filepath.Dir(path)
	if dir == "." || dir == "/" || dir == "" {
		return nil
	}
	return os.MkdirAll(dir, 0o755)
}

type Row struct {
	Start, Mid1, Mid2, End string
	L1, L2, L3             Leg
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	// Load .env if exists
	envFile := getenv("ENV_FILE", DefaultEnvFile)
	loadDotEnvIfExists(envFile)

	instrumentsURL := getenv("OKX_INSTRUMENTS_URL", DefaultInstrumentsURL)
	outCSV := getenv("OUTPUT_CSV", DefaultOutCSV)
	startCcy := strings.ToUpper(strings.TrimSpace(getenv("START_CCY", DefaultStartCcy)))
	timeout := getenvDuration("HTTP_TIMEOUT", DefaultHTTPTimeout)

	log.Printf("CFG: ENV_FILE=%s", envFile)
	log.Printf("CFG: OKX_INSTRUMENTS_URL=%s", instrumentsURL)
	log.Printf("CFG: OUTPUT_CSV=%s", outCSV)
	log.Printf("CFG: START_CCY=%s", startCcy)
	log.Printf("CFG: HTTP_TIMEOUT=%s", timeout)

	// 1) Load instruments
	all, err := loadRulesOKX(instrumentsURL, timeout)
	if err != nil {
		log.Fatalf("ERR: load instruments: %v", err)
	}
	log.Printf("OK: instruments loaded: %d", len(all))

	// 2) Filter eligible + build maps:
	marketByPair := make(map[string]okxInstrument, 8192)
	neighbors := make([]string, 0, 2048)
	neighborSet := make(map[string]struct{}, 2048)

	eligibleCount := 0
	for _, inst := range all {
		if !marketOkOKX(inst) {
			continue
		}
		eligibleCount++

		base := strings.ToUpper(strings.TrimSpace(inst.BaseCcy))
		quote := strings.ToUpper(strings.TrimSpace(inst.QuoteCcy))
		if base == "" || quote == "" || base == quote {
			continue
		}

		marketByPair[pairKey(base, quote)] = inst

		if base == startCcy && quote != startCcy {
			if _, ok := neighborSet[quote]; !ok {
				neighborSet[quote] = struct{}{}
				neighbors = append(neighbors, quote)
			}
		} else if quote == startCcy && base != startCcy {
			if _, ok := neighborSet[base]; !ok {
				neighborSet[base] = struct{}{}
				neighbors = append(neighbors, base)
			}
		}
	}
	sort.Strings(neighbors)
	log.Printf("OK: eligible instruments: %d", eligibleCount)
	log.Printf("OK: neighbors(%s): %d", startCcy, len(neighbors))

	if len(neighbors) < 2 {
		log.Fatalf("ERR: not enough neighbors for %s", startCcy)
	}

	// 3) Generate triangles: require X<->Y market, emit both directions
	rows := make([]Row, 0, 200000)
	seen := make(map[string]struct{}, 1<<20)

	for i := 0; i < len(neighbors); i++ {
		for j := i + 1; j < len(neighbors); j++ {
			X := neighbors[i]
			Y := neighbors[j]

			instSX, ok1 := marketByPair[pairKey(startCcy, X)]
			instSY, ok2 := marketByPair[pairKey(startCcy, Y)]
			instXY, ok3 := marketByPair[pairKey(X, Y)]
			if !(ok1 && ok2 && ok3) {
				continue
			}

			emit := func(mid1, mid2 string, a, b, c okxInstrument) {
				l1, ok := makeLeg(startCcy, mid1, a)
				if !ok {
					return
				}
				l2, ok := makeLeg(mid1, mid2, b)
				if !ok {
					return
				}
				l3, ok := makeLeg(mid2, startCcy, c)
				if !ok {
					return
				}

				if l1.InstID == l2.InstID || l1.InstID == l3.InstID || l2.InstID == l3.InstID {
					return
				}

				key := startCcy + ">" + mid1 + ">" + mid2 + ">" + startCcy + "|" + l1.InstID + "|" + l2.InstID + "|" + l3.InstID
				if _, exists := seen[key]; exists {
					return
				}
				seen[key] = struct{}{}

				rows = append(rows, Row{
					Start: startCcy, Mid1: mid1, Mid2: mid2, End: startCcy,
					L1: l1, L2: l2, L3: l3,
				})
			}

			emit(X, Y, instSX, instXY, instSY)
			emit(Y, X, instSY, instXY, instSX)
		}
	}

	sort.Slice(rows, func(i, j int) bool {
		a := rows[i]
		b := rows[j]
		if a.Mid1 != b.Mid1 {
			return a.Mid1 < b.Mid1
		}
		if a.Mid2 != b.Mid2 {
			return a.Mid2 < b.Mid2
		}
		if a.L1.InstID != b.L1.InstID {
			return a.L1.InstID < b.L1.InstID
		}
		if a.L2.InstID != b.L2.InstID {
			return a.L2.InstID < b.L2.InstID
		}
		return a.L3.InstID < b.L3.InstID
	})

	// 4) Write CSV
	if err := ensureDirForFile(outCSV); err != nil {
		log.Fatalf("ERR: ensure dir: %v", err)
	}

	f, err := os.Create(outCSV)
	if err != nil {
		log.Fatalf("ERR: create csv: %v", err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	// added: qty_step/min_qty/price_step
	header := []string{
		"start", "mid1", "mid2", "end",

		"leg1_instId", "leg1_action", "leg1_price_side", "leg1_from", "leg1_to",
		"leg1_mkt_limit_kind", "leg1_has_max_mkt_amt", "leg1_has_max_mkt_sz", "leg1_mkt_limit_field", "leg1_mkt_limit_value",
		"leg1_qty_step", "leg1_min_qty", "leg1_price_step",
		"leg1_base", "leg1_quote", "leg1_lotSz", "leg1_minSz", "leg1_tickSz", "leg1_maxMktAmt", "leg1_maxMktSz", "leg1_qty_dp", "leg1_price_dp",

		"leg2_instId", "leg2_action", "leg2_price_side", "leg2_from", "leg2_to",
		"leg2_mkt_limit_kind", "leg2_has_max_mkt_amt", "leg2_has_max_mkt_sz", "leg2_mkt_limit_field", "leg2_mkt_limit_value",
		"leg2_qty_step", "leg2_min_qty", "leg2_price_step",
		"leg2_base", "leg2_quote", "leg2_lotSz", "leg2_minSz", "leg2_tickSz", "leg2_maxMktAmt", "leg2_maxMktSz", "leg2_qty_dp", "leg2_price_dp",

		"leg3_instId", "leg3_action", "leg3_price_side", "leg3_from", "leg3_to",
		"leg3_mkt_limit_kind", "leg3_has_max_mkt_amt", "leg3_has_max_mkt_sz", "leg3_mkt_limit_field", "leg3_mkt_limit_value",
		"leg3_qty_step", "leg3_min_qty", "leg3_price_step",
		"leg3_base", "leg3_quote", "leg3_lotSz", "leg3_minSz", "leg3_tickSz", "leg3_maxMktAmt", "leg3_maxMktSz", "leg3_qty_dp", "leg3_price_dp",
	}
	if err := w.Write(header); err != nil {
		log.Fatalf("ERR: write header: %v", err)
	}

	itoa := func(x int) string { return strconv.Itoa(x) }

	for _, r := range rows {
		row := []string{
			r.Start, r.Mid1, r.Mid2, r.End,

			r.L1.InstID, r.L1.Action, r.L1.PriceSide, r.L1.From, r.L1.To,
			r.L1.MktLimitKind, r.L1.HasMaxMktAmt, r.L1.HasMaxMktSz, r.L1.MktLimitField, r.L1.MktLimitValue,
			r.L1.QtyStep, r.L1.MinQty, r.L1.PriceStep,
			r.L1.BaseCcy, r.L1.QuoteCcy, r.L1.LotSz, r.L1.MinSzRaw, r.L1.TickSz, r.L1.MaxMktAmt, r.L1.MaxMktSz, itoa(r.L1.QtyDP), itoa(r.L1.PriceDP),

			r.L2.InstID, r.L2.Action, r.L2.PriceSide, r.L2.From, r.L2.To,
			r.L2.MktLimitKind, r.L2.HasMaxMktAmt, r.L2.HasMaxMktSz, r.L2.MktLimitField, r.L2.MktLimitValue,
			r.L2.QtyStep, r.L2.MinQty, r.L2.PriceStep,
			r.L2.BaseCcy, r.L2.QuoteCcy, r.L2.LotSz, r.L2.MinSzRaw, r.L2.TickSz, r.L2.MaxMktAmt, r.L2.MaxMktSz, itoa(r.L2.QtyDP), itoa(r.L2.PriceDP),

			r.L3.InstID, r.L3.Action, r.L3.PriceSide, r.L3.From, r.L3.To,
			r.L3.MktLimitKind, r.L3.HasMaxMktAmt, r.L3.HasMaxMktSz, r.L3.MktLimitField, r.L3.MktLimitValue,
			r.L3.QtyStep, r.L3.MinQty, r.L3.PriceStep,
			r.L3.BaseCcy, r.L3.QuoteCcy, r.L3.LotSz, r.L3.MinSzRaw, r.L3.TickSz, r.L3.MaxMktAmt, r.L3.MaxMktSz, itoa(r.L3.QtyDP), itoa(r.L3.PriceDP),
		}
		if err := w.Write(row); err != nil {
			log.Fatalf("ERR: write row: %v", err)
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		log.Fatalf("ERR: csv flush: %v", err)
	}

	log.Printf("OK: triangles=%d (both directions, X<->Y required) -> %s", len(rows), outCSV)
}


