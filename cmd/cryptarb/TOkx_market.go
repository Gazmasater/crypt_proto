package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
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

const (
	DefaultInputCSV      = "triangles_usdt_routes_okx.csv"
	DefaultOutputCSV     = "triangles_usdt_routes_market_okx.csv"
	DefaultBlacklistFile = "blacklist_symbols_okx.txt"

	DefaultInstrumentsURL = "https://www.okx.com/api/v5/public/instruments?instType=SPOT"
	DefaultBooksURL       = "https://www.okx.com/api/v5/market/books"

	DefaultTimeout  = 25 * time.Second
	DefaultStartAmt = 25.0
	DefaultStartCcy = "USDT"
)

// ===================== OKX API =====================

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
	TickSz string `json:"tickSz"` // price step

	MaxMktAmt string `json:"maxMktAmt"` // quote
	MaxMktSz  string `json:"maxMktSz"`  // base
}

type okxBooksResp struct {
	Code string        `json:"code"`
	Msg  string        `json:"msg"`
	Data []okxBookData `json:"data"`
}

type okxBookData struct {
	Asks [][]string `json:"asks"` // [price, size, ...]
	Bids [][]string `json:"bids"`
}

type Book struct {
	Ask float64
	Bid float64
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
	tmp := path + ".tmp"
	if err := ensureDirForFile(path); err != nil {
		return err
	}

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

// ===================== Utils =====================

func ensureDirForFile(path string) error {
	dir := filepath.Dir(path)
	if dir == "." || dir == "/" || dir == "" {
		return nil
	}
	return os.MkdirAll(dir, 0o755)
}

func colIndex(header []string, name string) int {
	name = strings.ToLower(strings.TrimSpace(name))
	for i, h := range header {
		if strings.ToLower(strings.TrimSpace(h)) == name {
			return i
		}
	}
	return -1
}

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

// ======= IMPORTANT FIX: quantize by exact step (lotSz), not by decimals =======

// ratFromFloat makes decimal string without exponent for big.Rat.
func ratFromFloat(x float64) (*big.Rat, bool) {
	if math.IsNaN(x) || math.IsInf(x, 0) {
		return nil, false
	}
	// fixed precision string; OKX steps usually fit well in 18 decimals
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

// quantizeDownStep returns floor(value/step)*step for positive numbers.
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
	if !ok {
		return 0, false
	}
	if v.Sign() < 0 {
		return 0, false
	}

	// q = floor(v/step)
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
	req.Header.Set("User-Agent", "cryptarb/1.0 (+okx instruments loader)")

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
	for _, s := range r.Data {
		id := strings.TrimSpace(s.InstID)
		if id == "" {
			continue
		}
		m[id] = s
	}
	return m, nil
}

func marketOkOKX(s okxInstrument) bool {
	if s.InstType != "" && !strings.EqualFold(strings.TrimSpace(s.InstType), "SPOT") {
		return false
	}
	if !strings.EqualFold(strings.TrimSpace(s.State), "live") {
		return false
	}
	if strings.TrimSpace(s.LotSz) == "" || strings.TrimSpace(s.MinSz) == "" || strings.TrimSpace(s.TickSz) == "" {
		return false
	}
	// либо maxMktAmt, либо maxMktSz
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

func fetchTopBook(client *http.Client, booksBaseURL, instId string) (Book, error) {
	url := fmt.Sprintf("%s?instId=%s&sz=1", strings.TrimRight(booksBaseURL, "/"), instId)
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
		return Book{}, fmt.Errorf("empty books data")
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

// ===================== Simulation =====================

func legKind(from, to string, inst okxInstrument) (kind string, ok bool) {
	from = strings.TrimSpace(from)
	to = strings.TrimSpace(to)

	if strings.EqualFold(from, inst.QuoteCcy) && strings.EqualFold(to, inst.BaseCcy) {
		return "BUY", true
	}
	if strings.EqualFold(from, inst.BaseCcy) && strings.EqualFold(to, inst.QuoteCcy) {
		return "SELL", true
	}
	return "", false
}

type LegRow struct {
	Symbol string
	From   string
	To     string
}

func simulateRoute(
	startAmt float64,
	startCcy string,
	legs [3]LegRow,
	rules map[string]okxInstrument,
	bookCache map[string]Book,
	client *http.Client,
	booksURL string,
) (endAmt float64, fail string) {

	amt := startAmt
	ccy := startCcy

	for i := 0; i < 3; i++ {
		leg := legs[i]

		inst, ok := rules[leg.Symbol]
		if !ok {
			return 0, "no instrument"
		}
		if !marketOkOKX(inst) {
			return 0, "not eligible"
		}

		kind, ok := legKind(leg.From, leg.To, inst)
		if !ok {
			return 0, "from/to mismatch instrument base/quote"
		}

		bk, ok := bookCache[leg.Symbol]
		if !ok {
			b, err := fetchTopBook(client, booksURL, leg.Symbol)
			if err != nil {
				return 0, "book err: " + err.Error()
			}
			bookCache[leg.Symbol] = b
			bk = b
		}

		minSz, okMin := parseF(inst.MinSz)
		maxMktAmt, okMaxAmt := parseF(inst.MaxMktAmt)
		maxMktSz, okMaxSz := parseF(inst.MaxMktSz)

		switch kind {
		case "BUY":
			if !strings.EqualFold(ccy, inst.QuoteCcy) {
				return 0, "buy wrong ccy (need quote)"
			}
			if okMaxAmt && amt > maxMktAmt {
				return 0, "exceeds maxMktAmt"
			}

			baseRaw := amt / bk.Ask
			base, okQ := quantizeDownStep(baseRaw, inst.LotSz)
			if !okQ {
				return 0, "quantize lotSz"
			}

			if !okMin {
				return 0, "minSz parse"
			}
			if base < minSz {
				return 0, "below minSz"
			}
			if okMaxSz && base > maxMktSz {
				return 0, "exceeds maxMktSz"
			}

			amt = base
			ccy = inst.BaseCcy

		case "SELL":
			if !strings.EqualFold(ccy, inst.BaseCcy) {
				return 0, "sell wrong ccy (need base)"
			}
			if !okMin {
				return 0, "minSz parse"
			}

			base, okQ := quantizeDownStep(amt, inst.LotSz)
			if !okQ {
				return 0, "quantize lotSz"
			}
			if base < minSz {
				return 0, "below minSz"
			}
			if okMaxSz && base > maxMktSz {
				return 0, "exceeds maxMktSz"
			}

			quote := base * bk.Bid
			amt = quote
			ccy = inst.QuoteCcy
		}
	}

	if !strings.EqualFold(ccy, startCcy) {
		return 0, "end ccy mismatch"
	}
	return amt, ""
}

// ===================== CSV filter =====================

func buildOutputCSVOKX(
	inputCSV, outputCSV string,
	startAmt float64,
	startCcy string,
	keepFailed bool,
	rules map[string]okxInstrument,
	bl map[string]string,
	booksURL string,
	timeout time.Duration,
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
		return fmt.Errorf("CSV must contain leg*_symbol, leg*_from, leg*_to")
	}

	header = append(header,
		"start_amt", "end_amt", "fail_reason",
		"leg1_lot_sz", "leg2_lot_sz", "leg3_lot_sz",
		"leg1_min_sz", "leg2_min_sz", "leg3_min_sz",
		"leg1_tick_sz", "leg2_tick_sz", "leg3_tick_sz",
		"leg1_max_mkt_amt", "leg2_max_mkt_amt", "leg3_max_mkt_amt",
		"leg1_max_mkt_sz", "leg2_max_mkt_sz", "leg3_max_mkt_sz",
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

	for {
		row, err := cr.Read()
		if err == io.EOF {
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

		if isBlacklisted(bl, s1) || isBlacklisted(bl, s2) || isBlacklisted(bl, s3) {
			skippedBlacklisted++
			if !keepFailed {
				continue
			}
		}

		r1, ok1 := rules[s1]
		r2, ok2 := rules[s2]
		r3, ok3 := rules[s3]
		if !ok1 || !ok2 || !ok3 {
			skippedNoSymbol++
			if !keepFailed {
				continue
			}
		}

		eligible := ok1 && ok2 && ok3 && marketOkOKX(r1) && marketOkOKX(r2) && marketOkOKX(r3)
		if !eligible {
			skippedNotEligible++
			if !keepFailed {
				continue
			}
		}

		legs := [3]LegRow{
			{Symbol: s1, From: strings.TrimSpace(row[iF1]), To: strings.TrimSpace(row[iT1])},
			{Symbol: s2, From: strings.TrimSpace(row[iF2]), To: strings.TrimSpace(row[iT2])},
			{Symbol: s3, From: strings.TrimSpace(row[iF3]), To: strings.TrimSpace(row[iT3])},
		}

		endAmt, fail := 0.0, ""
		if eligible && !isBlacklisted(bl, s1) && !isBlacklisted(bl, s2) && !isBlacklisted(bl, s3) {
			endAmt, fail = simulateRoute(startAmt, startCcy, legs, rules, bookCache, client, booksURL)
		} else {
			if isBlacklisted(bl, s1) || isBlacklisted(bl, s2) || isBlacklisted(bl, s3) {
				fail = "blacklisted"
			} else if !eligible {
				fail = "not eligible"
			} else {
				fail = "no instrument"
			}
		}

		if fail != "" && !keepFailed {
			skippedNotFeasible++
			continue
		}

		getOrEmpty := func(sym string) okxInstrument {
			if v, ok := rules[sym]; ok {
				return v
			}
			return okxInstrument{}
		}
		rr1 := getOrEmpty(s1)
		rr2 := getOrEmpty(s2)
		rr3 := getOrEmpty(s3)

		row = append(row,
			fmt.Sprintf("%.8f", startAmt),
			fmt.Sprintf("%.8f", endAmt),
			fail,

			strings.TrimSpace(rr1.LotSz), strings.TrimSpace(rr2.LotSz), strings.TrimSpace(rr3.LotSz),
			strings.TrimSpace(rr1.MinSz), strings.TrimSpace(rr2.MinSz), strings.TrimSpace(rr3.MinSz),
			strings.TrimSpace(rr1.TickSz), strings.TrimSpace(rr2.TickSz), strings.TrimSpace(rr3.TickSz),
			strings.TrimSpace(rr1.MaxMktAmt), strings.TrimSpace(rr2.MaxMktAmt), strings.TrimSpace(rr3.MaxMktAmt),
			strings.TrimSpace(rr1.MaxMktSz), strings.TrimSpace(rr2.MaxMktSz), strings.TrimSpace(rr3.MaxMktSz),
		)

		if err := cw.Write(row); err != nil {
			return fmt.Errorf("write row: %w", err)
		}
		written++
	}

	log.Printf(
		"OK: read=%d written=%d skippedNoSymbol=%d skippedBlacklisted=%d skippedNotEligible=%d skippedNotFeasible(start=%.2f %s)=%d keep_failed=%v -> %s",
		read, written, skippedNoSymbol, skippedBlacklisted, skippedNotEligible, startAmt, startCcy, skippedNotFeasible, keepFailed, outputCSV,
	)

	return nil
}

// ===================== main =====================

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	inputCSV := getenv("INPUT_CSV", DefaultInputCSV)
	outputCSV := getenv("OUTPUT_CSV", DefaultOutputCSV)
	blacklistFile := getenv("BLACKLIST_FILE", DefaultBlacklistFile)

	instrumentsURL := getenv("OKX_INSTRUMENTS_URL", DefaultInstrumentsURL)
	booksURL := getenv("OKX_BOOKS_URL", DefaultBooksURL)

	timeout := getenvDuration("HTTP_TIMEOUT", DefaultTimeout)
	startAmt := getenvFloat("START_AMT", DefaultStartAmt)
	startCcy := getenv("START_CCY", DefaultStartCcy)

	keepFailed := getenvBool("KEEP_FAILED", false)

	log.Printf("CFG: INPUT_CSV=%s", inputCSV)
	log.Printf("CFG: OUTPUT_CSV=%s", outputCSV)
	log.Printf("CFG: BLACKLIST_FILE=%s", blacklistFile)
	log.Printf("CFG: OKX_INSTRUMENTS_URL=%s", instrumentsURL)
	log.Printf("CFG: OKX_BOOKS_URL=%s", booksURL)
	log.Printf("CFG: HTTP_TIMEOUT=%s", timeout)
	log.Printf("CFG: START_AMT=%.8f", startAmt)
	log.Printf("CFG: START_CCY=%s", startCcy)
	log.Printf("CFG: KEEP_FAILED=%v", keepFailed)

	bl, err := loadBlacklist(blacklistFile)
	if err != nil {
		log.Fatalf("ERR: load blacklist: %v", err)
	}
	log.Printf("OK: loaded blacklist: %d symbols", len(bl))

	rules, err := loadRulesOKX(instrumentsURL, timeout)
	if err != nil {
		log.Fatalf("ERR: load OKX instruments: %v", err)
	}
	log.Printf("OK: OKX instruments loaded: %d", len(rules))

	// Startup blacklist
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

	if err := buildOutputCSVOKX(inputCSV, outputCSV, startAmt, strings.TrimSpace(startCcy), keepFailed, rules, bl, booksURL, timeout); err != nil {
		log.Fatalf("ERR: build output csv: %v", err)
	}

	log.Printf("DONE")
}
