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
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

// ===================== Defaults =====================

const (
	DefaultEnvFile        = ".env"
	DefaultInputCSV       = "triangles_usdt_routes_okx.csv"
	DefaultOutputCSV      = "triangles_usdt_routes_market_okx.csv"
	DefaultBlacklistFile  = "blacklist_symbols_okx.txt"
	DefaultInstrumentsURL = "https://www.okx.com/api/v5/public/instruments?instType=SPOT"
	DefaultBooksURL       = "https://www.okx.com/api/v5/market/books"

	DefaultTimeout      = 25 * time.Second
	DefaultStartAmt     = 25.0
	DefaultStartCcy     = "USDT"
	DefaultFeePct       = 0.1
	DefaultMinProfitPct = 0.3
	DefaultKeepFailed   = false
	DefaultBookSz       = 1
)

// ===================== OKX API structs =====================

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

	LotSz  string `json:"lotSz"`  // base step
	MinSz  string `json:"minSz"`  // base min
	TickSz string `json:"tickSz"` // price tick

	MaxMktAmt string `json:"maxMktAmt"` // quote limit (mostly for BUY)
	MaxMktSz  string `json:"maxMktSz"`  // base limit (mostly for SELL)
}

type okxBooksResp struct {
	Code string        `json:"code"`
	Msg  string        `json:"msg"`
	Data []okxBookData `json:"data"`
}

type okxBookData struct {
	Asks [][]string `json:"asks"` // [price, size, ...]
	Bids [][]string `json:"bids"` // [price, size, ...]
}

type Book struct {
	Ask float64
	Bid float64
}

// ===================== .env loader =====================
// Loads ENV_FILE (default ".env") if present. Does NOT override already-set env vars.

func loadDotEnvIfExists(path string) {
	path = strings.TrimSpace(path)
	if path == "" {
		return
	}
	f, err := os.Open(path)
	if err != nil {
		// silent if not exist
		return
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		// allow "export KEY=VAL"
		if strings.HasPrefix(line, "export ") {
			line = strings.TrimSpace(strings.TrimPrefix(line, "export "))
		}
		k, v, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key := strings.TrimSpace(k)
		val := strings.TrimSpace(v)

		// strip quotes
		if len(val) >= 2 {
			if (val[0] == '"' && val[len(val)-1] == '"') || (val[0] == '\'' && val[len(val)-1] == '\'') {
				val = val[1 : len(val)-1]
			}
		}
		if key == "" {
			continue
		}
		if _, exists := os.LookupEnv(key); exists {
			continue // do not override
		}
		_ = os.Setenv(key, val)
	}
}

// ===================== ENV helpers =====================

func getenv(key, def string) string {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return def
	}
	return v
}

func getenvBool(key string, def bool) bool {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return def
	}
	switch strings.ToLower(v) {
	case "1", "true", "yes", "y", "on":
		return true
	case "0", "false", "no", "n", "off":
		return false
	default:
		return def
	}
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

func getenvFloat(key string, def float64) float64 {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return def
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil || math.IsNaN(f) || math.IsInf(f, 0) {
		log.Printf("WARN: bad %s=%q, using default %.8f", key, v, def)
		return def
	}
	return f
}

func getenvInt(key string, def int) int {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		log.Printf("WARN: bad %s=%q, using default %d", key, v, def)
		return def
	}
	return n
}

// ===================== Blacklist =====================

func loadBlacklist(path string) (map[string]string, error) {
	bl := make(map[string]string)
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return bl, nil
		}
		return nil, err
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "\t", 2)
		sym := strings.TrimSpace(parts[0])
		if sym == "" {
			continue
		}
		reason := ""
		if len(parts) == 2 {
			reason = strings.TrimSpace(parts[1])
		}
		bl[sym] = reason
	}
	return bl, sc.Err()
}

func saveBlacklist(path string, bl map[string]string) error {
	if err := ensureDirForFile(path); err != nil {
		return err
	}
	tmp := path + ".tmp"
	f, err := os.Create(tmp)
	if err != nil {
		return err
	}
	defer f.Close()

	bw := bufio.NewWriter(f)
	fmt.Fprintf(bw, "# generated: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Fprintf(bw, "# format: SYMBOL<TAB>reason\n")

	keys := make([]string, 0, len(bl))
	for k := range bl {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Fprintf(bw, "%s\t%s\n", k, bl[k])
	}
	if err := bw.Flush(); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

func isBlacklisted(bl map[string]string, sym string) bool {
	_, ok := bl[sym]
	return ok
}

// ===================== CSV helpers =====================

func ensureDirForFile(path string) error {
	dir := filepath.Dir(path)
	if dir == "." || dir == "/" || dir == "" {
		return nil
	}
	return os.MkdirAll(dir, 0o755)
}

func colIndex(header []string, name string) int {
	want := strings.ToLower(strings.TrimSpace(name))
	for i, h := range header {
		if strings.ToLower(strings.TrimSpace(h)) == want {
			return i
		}
	}
	return -1
}

// ===================== parse helpers =====================

func parseF(s string) (float64, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, false
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil || math.IsNaN(v) || math.IsInf(v, 0) {
		return 0, false
	}
	return v, true
}

// ===================== quantize by exact step (lotSz) =====================

func ratFromFloat(x float64) (*big.Rat, bool) {
	if math.IsNaN(x) || math.IsInf(x, 0) {
		return nil, false
	}
	// 18 decimals is enough for OKX steps; we remove trailing zeros.
	s := strconv.FormatFloat(x, 'f', 18, 64)
	s = strings.TrimRight(s, "0")
	s = strings.TrimRight(s, ".")
	if s == "" {
		s = "0"
	}
	r := new(big.Rat)
	if _, ok := r.SetString(s); !ok {
		return nil, false
	}
	return r, true
}

// floor(value/step) * step, for value>=0
func quantizeDownStep(value float64, stepStr string) (float64, bool) {
	stepStr = strings.TrimSpace(stepStr)
	if stepStr == "" {
		return value, true
	}

	step := new(big.Rat)
	if _, ok := step.SetString(stepStr); !ok {
		return 0, false
	}
	if step.Sign() <= 0 {
		return 0, false
	}

	v, ok := ratFromFloat(value)
	if !ok || v.Sign() < 0 {
		return 0, false
	}

	qRat := new(big.Rat).Quo(v, step)
	q := new(big.Int).Quo(qRat.Num(), qRat.Denom()) // floor for positive

	res := new(big.Rat).Mul(new(big.Rat).SetInt(q), step)
	f, _ := res.Float64()
	return f, true
}

// ===================== OKX rules =====================

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

	m := make(map[string]okxInstrument, len(r.Data))
	for _, it := range r.Data {
		id := strings.TrimSpace(it.InstID)
		if id == "" {
			continue
		}
		m[id] = it
	}
	return m, nil
}

// your criteria:
// instType==SPOT, state==live
// lotSz/minSz/tickSz not empty
// and (maxMktAmt or maxMktSz) exists
func marketOkOKX(s okxInstrument) bool {
	if !strings.EqualFold(strings.TrimSpace(s.InstType), "SPOT") && strings.TrimSpace(s.InstType) != "" {
		return false
	}
	if !strings.EqualFold(strings.TrimSpace(s.State), "live") {
		return false
	}
	if strings.TrimSpace(s.LotSz) == "" || strings.TrimSpace(s.MinSz) == "" || strings.TrimSpace(s.TickSz) == "" {
		return false
	}
	if strings.TrimSpace(s.MaxMktAmt) == "" && strings.TrimSpace(s.MaxMktSz) == "" {
		return false
	}
	return true
}

func notOkReasonOKX(s okxInstrument) string {
	var reasons []string
	if s.InstType != "" && !strings.EqualFold(strings.TrimSpace(s.InstType), "SPOT") {
		reasons = append(reasons, "instType!=SPOT")
	}
	if !strings.EqualFold(strings.TrimSpace(s.State), "live") {
		reasons = append(reasons, "state!=live")
	}
	if strings.TrimSpace(s.LotSz) == "" {
		reasons = append(reasons, "lotSz empty")
	}
	if strings.TrimSpace(s.MinSz) == "" {
		reasons = append(reasons, "minSz empty")
	}
	if strings.TrimSpace(s.TickSz) == "" {
		reasons = append(reasons, "tickSz empty")
	}
	if strings.TrimSpace(s.MaxMktAmt) == "" && strings.TrimSpace(s.MaxMktSz) == "" {
		reasons = append(reasons, "maxMktAmt/maxMktSz empty")
	}
	return strings.Join(reasons, ", ")
}

// ===================== OKX books =====================

func fetchTopBook(client *http.Client, booksURL, instId string, sz int) (Book, error) {
	url := fmt.Sprintf("%s?instId=%s&sz=%d", strings.TrimRight(booksURL, "/"), instId, sz)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return Book{}, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "cryptarb/1.0 (+okx books)")

	resp, err := client.Do(req)
	if err != nil {
		return Book{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		return Book{}, fmt.Errorf("books http %d: %s", resp.StatusCode, strings.TrimSpace(string(b)))
	}

	var r okxBooksResp
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return Book{}, err
	}
	if strings.TrimSpace(r.Code) != "" && strings.TrimSpace(r.Code) != "0" {
		return Book{}, fmt.Errorf("books code=%s msg=%s", strings.TrimSpace(r.Code), strings.TrimSpace(r.Msg))
	}
	if len(r.Data) == 0 {
		return Book{}, fmt.Errorf("empty data")
	}
	d := r.Data[0]
	if len(d.Asks) == 0 || len(d.Bids) == 0 {
		return Book{}, fmt.Errorf("empty asks/bids")
	}

	ask, okA := parseF(d.Asks[0][0])
	bid, okB := parseF(d.Bids[0][0])
	if !okA || !okB || ask <= 0 || bid <= 0 {
		return Book{}, fmt.Errorf("bad ask/bid")
	}
	return Book{Ask: ask, Bid: bid}, nil
}

// ===================== simulation =====================

type LegRow struct {
	Symbol string
	From   string
	To     string
}

func legKind(from, to string, inst okxInstrument) (string, bool) {
	from = strings.TrimSpace(from)
	to = strings.TrimSpace(to)

	// BUY: quote -> base
	if strings.EqualFold(from, inst.QuoteCcy) && strings.EqualFold(to, inst.BaseCcy) {
		return "BUY", true
	}
	// SELL: base -> quote
	if strings.EqualFold(from, inst.BaseCcy) && strings.EqualFold(to, inst.QuoteCcy) {
		return "SELL", true
	}
	return "", false
}

type SimResult struct {
	EndAmt    float64
	ProfitPct float64
	Fail      string
}

func simulateRoute(
	startAmt float64,
	startCcy string,
	feePct float64,
	minProfitPct float64,
	legs [3]LegRow,
	rules map[string]okxInstrument,
	bookCache map[string]Book,
	client *http.Client,
	booksURL string,
	bookSz int,
) SimResult {

	if startAmt <= 0 {
		return SimResult{Fail: "bad startAmt"}
	}
	feeMul := 1.0 - feePct/100.0
	if feeMul <= 0 || feeMul > 1 {
		return SimResult{Fail: "bad feePct"}
	}

	amt := startAmt
	ccy := startCcy

	for i := 0; i < 3; i++ {
		leg := legs[i]

		inst, ok := rules[leg.Symbol]
		if !ok {
			return SimResult{Fail: "no instrument"}
		}
		if !marketOkOKX(inst) {
			return SimResult{Fail: "not eligible"}
		}

		kind, ok := legKind(leg.From, leg.To, inst)
		if !ok {
			return SimResult{Fail: "from/to mismatch"}
		}

		// get book
		bk, ok := bookCache[leg.Symbol]
		if !ok {
			b, err := fetchTopBook(client, booksURL, leg.Symbol, bookSz)
			if err != nil {
				return SimResult{Fail: "book err: " + err.Error()}
			}
			bookCache[leg.Symbol] = b
			bk = b
		}

		minSz, okMin := parseF(inst.MinSz)
		if !okMin {
			return SimResult{Fail: "minSz parse"}
		}

		maxMktAmt, hasMaxAmt := parseF(inst.MaxMktAmt)
		maxMktSz, hasMaxSz := parseF(inst.MaxMktSz)

		switch kind {
		case "BUY":
			// spend quote, receive base at ASK
			if !strings.EqualFold(ccy, inst.QuoteCcy) {
				return SimResult{Fail: "buy wrong ccy"}
			}
			// OKX: market BUY limited by maxMktAmt (quote)
			if hasMaxAmt && amt > maxMktAmt {
				return SimResult{Fail: "exceeds maxMktAmt"}
			}

			baseRaw := amt / bk.Ask
			base, okQ := quantizeDownStep(baseRaw, inst.LotSz)
			if !okQ {
				return SimResult{Fail: "quantize lotSz"}
			}
			if base < minSz {
				return SimResult{Fail: "below minSz"}
			}
			// Some markets still have maxMktSz. Apply if present.
			if hasMaxSz && base > maxMktSz {
				return SimResult{Fail: "exceeds maxMktSz"}
			}

			// fee after trade
			base *= feeMul

			amt = base
			ccy = inst.BaseCcy

		case "SELL":
			// sell base, receive quote at BID
			if !strings.EqualFold(ccy, inst.BaseCcy) {
				return SimResult{Fail: "sell wrong ccy"}
			}

			base, okQ := quantizeDownStep(amt, inst.LotSz)
			if !okQ {
				return SimResult{Fail: "quantize lotSz"}
			}
			if base < minSz {
				return SimResult{Fail: "below minSz"}
			}

			// OKX: market SELL limited by maxMktSz (base)
			if hasMaxSz && base > maxMktSz {
				return SimResult{Fail: "exceeds maxMktSz"}
			}

			quote := base * bk.Bid

			// Sometimes maxMktAmt is also meaningful; don't overblock, but if present we can check output quote too.
			if hasMaxAmt && quote > maxMktAmt {
				return SimResult{Fail: "exceeds maxMktAmt"}
			}

			// fee after trade
			quote *= feeMul

			amt = quote
			ccy = inst.QuoteCcy
		}
	}

	if !strings.EqualFold(ccy, startCcy) {
		return SimResult{Fail: "end ccy mismatch"}
	}

	profitPct := (amt/startAmt - 1.0) * 100.0
	if profitPct < minProfitPct {
		return SimResult{EndAmt: amt, ProfitPct: profitPct, Fail: "profit<min"}
	}

	return SimResult{EndAmt: amt, ProfitPct: profitPct, Fail: ""}
}

// ===================== main filter =====================

func buildOutputCSV(
	inputCSV, outputCSV string,
	startAmt float64,
	startCcy string,
	feePct float64,
	minProfitPct float64,
	keepFailed bool,
	rules map[string]okxInstrument,
	bl map[string]string,
	booksURL string,
	timeout time.Duration,
	bookSz int,
) error {
	in, err := os.Open(inputCSV)
	if err != nil {
		return fmt.Errorf("open %s: %w", inputCSV, err)
	}
	defer in.Close()

	cr := csv.NewReader(in)
	header, err := cr.Read()
	if err != nil {
		return fmt.Errorf("read header: %w", err)
	}

	iL1 := colIndex(header, "leg1_symbol")
	iL2 := colIndex(header, "leg2_symbol")
	iL3 := colIndex(header, "leg3_symbol")
	iF1 := colIndex(header, "leg1_from")
	iF2 := colIndex(header, "leg2_from")
	iF3 := colIndex(header, "leg3_from")
	iT1 := colIndex(header, "leg1_to")
	iT2 := colIndex(header, "leg2_to")
	iT3 := colIndex(header, "leg3_to")

	if iL1 < 0 || iL2 < 0 || iL3 < 0 || iF1 < 0 || iF2 < 0 || iF3 < 0 || iT1 < 0 || iT2 < 0 || iT3 < 0 {
		return errors.New("CSV must contain: leg1_symbol leg2_symbol leg3_symbol leg1_from leg2_from leg3_from leg1_to leg2_to leg3_to")
	}

	// add output columns
	header = append(header,
		"start_amt", "end_amt", "profit_pct", "fee_pct", "min_profit_pct", "fail_reason",
		"leg1_lotSz", "leg2_lotSz", "leg3_lotSz",
		"leg1_minSz", "leg2_minSz", "leg3_minSz",
		"leg1_maxMktAmt", "leg2_maxMktAmt", "leg3_maxMktAmt",
		"leg1_maxMktSz", "leg2_maxMktSz", "leg3_maxMktSz",
	)

	if err := ensureDirForFile(outputCSV); err != nil {
		return err
	}
	out, err := os.Create(outputCSV)
	if err != nil {
		return fmt.Errorf("create %s: %w", outputCSV, err)
	}
	defer out.Close()

	cw := csv.NewWriter(out)
	defer cw.Flush()

	if err := cw.Write(header); err != nil {
		return fmt.Errorf("write header: %w", err)
	}

	client := &http.Client{Timeout: timeout}
	bookCache := make(map[string]Book, 4096)

	read := 0
	written := 0

	skippedNoSymbol := 0
	skippedBlacklisted := 0
	skippedNotEligible := 0
	skippedNotFeasible := 0
	skippedBelowProfit := 0

	for {
		row, err := cr.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return fmt.Errorf("read csv: %w", err)
		}
		read++

		s1 := strings.TrimSpace(row[iL1])
		s2 := strings.TrimSpace(row[iL2])
		s3 := strings.TrimSpace(row[iL3])
		if s1 == "" || s2 == "" || s3 == "" {
			skippedNoSymbol++
			continue
		}

		fail := ""
		endAmt := 0.0
		profitPct := 0.0

		if isBlacklisted(bl, s1) || isBlacklisted(bl, s2) || isBlacklisted(bl, s3) {
			skippedBlacklisted++
			fail = "blacklisted"
		}

		r1, ok1 := rules[s1]
		r2, ok2 := rules[s2]
		r3, ok3 := rules[s3]
		if fail == "" && (!ok1 || !ok2 || !ok3) {
			skippedNoSymbol++
			fail = "no instrument"
		}

		if fail == "" && (!marketOkOKX(r1) || !marketOkOKX(r2) || !marketOkOKX(r3)) {
			skippedNotEligible++
			fail = "not eligible"
		}

		if fail == "" {
			legs := [3]LegRow{
				{Symbol: s1, From: strings.TrimSpace(row[iF1]), To: strings.TrimSpace(row[iT1])},
				{Symbol: s2, From: strings.TrimSpace(row[iF2]), To: strings.TrimSpace(row[iT2])},
				{Symbol: s3, From: strings.TrimSpace(row[iF3]), To: strings.TrimSpace(row[iT3])},
			}

			res := simulateRoute(
				startAmt, startCcy,
				feePct, minProfitPct,
				legs, rules, bookCache,
				client, booksURL, bookSz,
			)
			fail = res.Fail
			endAmt = res.EndAmt
			profitPct = res.ProfitPct

			if fail == "profit<min" {
				skippedBelowProfit++
			} else if fail != "" {
				skippedNotFeasible++
			}
		} else {
			skippedNotFeasible++
		}

		if fail != "" && !keepFailed {
			continue
		}

		// append data columns
		row = append(row,
			fmt.Sprintf("%.8f", startAmt),
			fmt.Sprintf("%.8f", endAmt),
			fmt.Sprintf("%.6f", profitPct),
			fmt.Sprintf("%.6f", feePct),
			fmt.Sprintf("%.6f", minProfitPct),
			fail,

			strings.TrimSpace(r1.LotSz), strings.TrimSpace(r2.LotSz), strings.TrimSpace(r3.LotSz),
			strings.TrimSpace(r1.MinSz), strings.TrimSpace(r2.MinSz), strings.TrimSpace(r3.MinSz),
			strings.TrimSpace(r1.MaxMktAmt), strings.TrimSpace(r2.MaxMktAmt), strings.TrimSpace(r3.MaxMktAmt),
			strings.TrimSpace(r1.MaxMktSz), strings.TrimSpace(r2.MaxMktSz), strings.TrimSpace(r3.MaxMktSz),
		)

		if err := cw.Write(row); err != nil {
			return fmt.Errorf("write row: %w", err)
		}
		written++
	}

	log.Printf("OK: read=%d written=%d skippedNoSymbol=%d skippedBlacklisted=%d skippedNotEligible=%d skippedNotFeasible=%d skippedBelowProfit=%d -> %s",
		read, written, skippedNoSymbol, skippedBlacklisted, skippedNotEligible, skippedNotFeasible, skippedBelowProfit, outputCSV)

	return nil
}

// ===================== main =====================

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	// 0) Load .env (optional)
	envFile := getenv("ENV_FILE", DefaultEnvFile)
	loadDotEnvIfExists(envFile)

	// 1) Read cfg
	inputCSV := getenv("INPUT_CSV", DefaultInputCSV)
	outputCSV := getenv("OUTPUT_CSV", DefaultOutputCSV)
	blacklistFile := getenv("BLACKLIST_FILE", DefaultBlacklistFile)

	instrumentsURL := getenv("OKX_INSTRUMENTS_URL", DefaultInstrumentsURL)
	booksURL := getenv("OKX_BOOKS_URL", DefaultBooksURL)

	timeout := getenvDuration("HTTP_TIMEOUT", DefaultTimeout)
	startAmt := getenvFloat("START_AMT", DefaultStartAmt)
	startCcy := getenv("START_CCY", DefaultStartCcy)
	feePct := getenvFloat("FEE_PCT", DefaultFeePct)
	minProfitPct := getenvFloat("MIN_PROFIT_PCT", DefaultMinProfitPct)
	keepFailed := getenvBool("KEEP_FAILED", DefaultKeepFailed)
	bookSz := getenvInt("BOOK_SZ", DefaultBookSz)

	// 2) Print cfg (so you ALWAYS see which version runs)
	log.Printf("CFG: ENV_FILE=%s", envFile)
	log.Printf("CFG: INPUT_CSV=%s", inputCSV)
	log.Printf("CFG: OUTPUT_CSV=%s", outputCSV)
	log.Printf("CFG: BLACKLIST_FILE=%s", blacklistFile)
	log.Printf("CFG: OKX_INSTRUMENTS_URL=%s", instrumentsURL)
	log.Printf("CFG: OKX_BOOKS_URL=%s", booksURL)
	log.Printf("CFG: HTTP_TIMEOUT=%s", timeout)
	log.Printf("CFG: START_AMT=%.8f", startAmt)
	log.Printf("CFG: START_CCY=%s", startCcy)
	log.Printf("CFG: FEE_PCT=%.6f", feePct)
	log.Printf("CFG: MIN_PROFIT_PCT=%.6f", minProfitPct)
	log.Printf("CFG: KEEP_FAILED=%v", keepFailed)
	log.Printf("CFG: BOOK_SZ=%d", bookSz)

	// 3) Load blacklist
	bl, err := loadBlacklist(blacklistFile)
	if err != nil {
		log.Fatalf("ERR: load blacklist: %v", err)
	}
	log.Printf("OK: loaded blacklist: %d symbols", len(bl))

	// 4) Load OKX instruments
	rules, err := loadRulesOKX(instrumentsURL, timeout)
	if err != nil {
		log.Fatalf("ERR: load OKX instruments: %v", err)
	}
	log.Printf("OK: OKX instruments loaded: %d", len(rules))

	// 5) Startup blacklist based on eligibility
	added := 0
	for instID, s := range rules {
		if marketOkOKX(s) {
			continue
		}
		reason := notOkReasonOKX(s)
		if reason == "" {
			reason = "not eligible"
		}
		if _, exists := bl[instID]; !exists {
			bl[instID] = reason
			added++
		}
	}
	log.Printf("OK: startup blacklist added=%d total=%d", added, len(bl))

	if err := saveBlacklist(blacklistFile, bl); err != nil {
		log.Fatalf("ERR: save blacklist: %v", err)
	}
	log.Printf("OK: saved blacklist -> %s", blacklistFile)

	// 6) Filter CSV
	if err := buildOutputCSV(
		inputCSV, outputCSV,
		startAmt, strings.TrimSpace(startCcy),
		feePct, minProfitPct,
		keepFailed,
		rules, bl,
		booksURL, timeout,
		bookSz,
	); err != nil {
		log.Fatalf("ERR: build output csv: %v", err)
	}

	log.Printf("DONE")
}
