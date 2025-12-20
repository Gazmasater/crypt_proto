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



package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
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

// ============ OKX API structs ============

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

	LotSz  string `json:"lotSz"`
	MinSz  string `json:"minSz"`
	TickSz string `json:"tickSz"`

	MaxMktAmt string `json:"maxMktAmt"` // quote
	MaxMktSz  string `json:"maxMktSz"`  // base
}

type okxBooksResp struct {
	Code string        `json:"code"`
	Msg  string        `json:"msg"`
	Data []okxBookData `json:"data"`
}

type okxBookData struct {
	Asks [][]string `json:"asks"`
	Bids [][]string `json:"bids"`
}

type Book struct {
	Ask float64
	Bid float64
}

// ============ Routes CSV leg ============

type LegRow struct {
	Symbol    string
	Action    string // BUY/SELL
	PriceSide string // ASK/BID (not used, just carried)
	From      string
	To        string
}

// ============ ENV helpers ============

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

// ============ Parsing helpers ============

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

// "0.001" -> 3, "1" -> 0, "" -> -1
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

func floorByDP(x float64, dp int) float64 {
	if dp <= 0 {
		return math.Floor(x)
	}
	pow := math.Pow10(dp)
	return math.Floor(x*pow) / pow
}

// ============ Blacklist I/O ============

// File format: SYMBOL<TAB>reason
// Lines starting with # are comments.
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

// ============ FS helpers ============

func ensureDirForFile(path string) error {
	dir := filepath.Dir(path)
	if dir == "." || dir == "/" || dir == "" {
		return nil
	}
	return os.MkdirAll(dir, 0o755)
}

// ============ OKX instruments/rules ============

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

// Eligibility:
// instType==SPOT
// state==live
// lotSz/minSz/tickSz not empty
// AND (maxMktAmt OR maxMktSz) not empty
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

// ============ OKX top-of-book ============

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

// ============ Route simulation for START_AMT (e.g. 25 USDT) ============

// BUY: spend quote at ASK -> receive base (floor to lotSz), check minSz, check maxMktAmt/maxMktSz
// SELL: sell base at BID (floor to lotSz), check minSz, check maxMktSz
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

		bk, ok := bookCache[leg.Symbol]
		if !ok {
			b, err := fetchTopBook(client, booksURL, leg.Symbol)
			if err != nil {
				return 0, "book err: " + err.Error()
			}
			bookCache[leg.Symbol] = b
			bk = b
		}

		lotDP := decimalsFromStep(inst.LotSz)
		minSz, okMin := parseF(inst.MinSz)
		maxMktAmt, okMaxAmt := parseF(inst.MaxMktAmt)
		maxMktSz, okMaxSz := parseF(inst.MaxMktSz)

		switch strings.ToUpper(strings.TrimSpace(leg.Action)) {
		case "BUY":
			// Expect current currency to be quote
			if ccy != inst.QuoteCcy {
				return 0, "buy wrong ccy (need quote)"
			}
			// If OKX restricts market amount in quote
			if okMaxAmt && amt > maxMktAmt {
				return 0, "exceeds maxMktAmt"
			}

			base := amt / bk.Ask
			base = floorByDP(base, lotDP)

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
			// Expect current currency to be base
			if ccy != inst.BaseCcy {
				return 0, "sell wrong ccy (need base)"
			}
			if !okMin {
				return 0, "minSz parse"
			}

			base := floorByDP(amt, lotDP)
			if base < minSz {
				return 0, "below minSz"
			}
			if okMaxSz && base > maxMktSz {
				return 0, "exceeds maxMktSz"
			}

			quote := base * bk.Bid
			amt = quote
			ccy = inst.QuoteCcy

		default:
			return 0, "bad action"
		}
	}

	if ccy != startCcy {
		return 0, "end ccy mismatch"
	}
	return amt, ""
}

// ============ CSV processing ============

func colIndex(header []string, name string) int {
	name = strings.ToLower(strings.TrimSpace(name))
	for i, h := range header {
		if strings.ToLower(strings.TrimSpace(h)) == name {
			return i
		}
	}
	return -1
}

func buildOutputCSVOKX(
	inputCSV, outputCSV string,
	startAmt float64,
	startCcy string,
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

	iA1 := colIndex(header, "leg1_action")
	iA2 := colIndex(header, "leg2_action")
	iA3 := colIndex(header, "leg3_action")

	iP1 := colIndex(header, "leg1_price")
	iP2 := colIndex(header, "leg2_price")
	iP3 := colIndex(header, "leg3_price")

	iF1 := colIndex(header, "leg1_from")
	iF2 := colIndex(header, "leg2_from")
	iF3 := colIndex(header, "leg3_from")

	iT1 := colIndex(header, "leg1_to")
	iT2 := colIndex(header, "leg2_to")
	iT3 := colIndex(header, "leg3_to")

	if iL1 < 0 || iL2 < 0 || iL3 < 0 ||
		iA1 < 0 || iA2 < 0 || iA3 < 0 ||
		iF1 < 0 || iF2 < 0 || iF3 < 0 ||
		iT1 < 0 || iT2 < 0 || iT3 < 0 {
		return fmt.Errorf("CSV must contain leg*_symbol, leg*_action, leg*_from, leg*_to")
	}
	_ = iP1
	_ = iP2
	_ = iP3

	header = append(header,
		"start_amt", "end_amt",
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
			continue
		}

		r1, ok1 := rules[s1]
		r2, ok2 := rules[s2]
		r3, ok3 := rules[s3]
		if !ok1 || !ok2 || !ok3 {
			skippedNoSymbol++
			continue
		}

		if !marketOkOKX(r1) || !marketOkOKX(r2) || !marketOkOKX(r3) {
			skippedNotEligible++
			continue
		}

		legs := [3]LegRow{
			{
				Symbol:    s1,
				Action:    strings.TrimSpace(row[iA1]),
				PriceSide: strings.TrimSpace(row[iP1]),
				From:      strings.TrimSpace(row[iF1]),
				To:        strings.TrimSpace(row[iT1]),
			},
			{
				Symbol:    s2,
				Action:    strings.TrimSpace(row[iA2]),
				PriceSide: strings.TrimSpace(row[iP2]),
				From:      strings.TrimSpace(row[iF2]),
				To:        strings.TrimSpace(row[iT2]),
			},
			{
				Symbol:    s3,
				Action:    strings.TrimSpace(row[iA3]),
				PriceSide: strings.TrimSpace(row[iP3]),
				From:      strings.TrimSpace(row[iF3]),
				To:        strings.TrimSpace(row[iT3]),
			},
		}

		endAmt, fail := simulateRoute(startAmt, startCcy, legs, rules, bookCache, client, booksURL)
		if fail != "" {
			skippedNotFeasible++
			continue
		}

		row = append(row,
			fmt.Sprintf("%.8f", startAmt),
			fmt.Sprintf("%.8f", endAmt),

			strings.TrimSpace(r1.LotSz), strings.TrimSpace(r2.LotSz), strings.TrimSpace(r3.LotSz),
			strings.TrimSpace(r1.MinSz), strings.TrimSpace(r2.MinSz), strings.TrimSpace(r3.MinSz),
			strings.TrimSpace(r1.TickSz), strings.TrimSpace(r2.TickSz), strings.TrimSpace(r3.TickSz),
			strings.TrimSpace(r1.MaxMktAmt), strings.TrimSpace(r2.MaxMktAmt), strings.TrimSpace(r3.MaxMktAmt),
			strings.TrimSpace(r1.MaxMktSz), strings.TrimSpace(r2.MaxMktSz), strings.TrimSpace(r3.MaxMktSz),
		)

		if err := cw.Write(row); err != nil {
			return fmt.Errorf("write row: %w", err)
		}
		written++
	}

	log.Printf(
		"OK: read=%d written=%d skippedNoSymbol=%d skippedBlacklisted=%d skippedNotEligible=%d skippedNotFeasible(%.2f %s)=%d -> %s",
		read, written, skippedNoSymbol, skippedBlacklisted, skippedNotEligible, startAmt, startCcy, skippedNotFeasible, outputCSV,
	)

	return nil
}

// ============ main ============

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	// Read from ENV
	inputCSV := getenv("INPUT_CSV", DefaultInputCSV)
	outputCSV := getenv("OUTPUT_CSV", DefaultOutputCSV)
	blacklistFile := getenv("BLACKLIST_FILE", DefaultBlacklistFile)

	instrumentsURL := getenv("OKX_INSTRUMENTS_URL", DefaultInstrumentsURL)
	booksURL := getenv("OKX_BOOKS_URL", DefaultBooksURL)

	timeout := getenvDuration("HTTP_TIMEOUT", DefaultTimeout)
	startAmt := getenvFloat("START_AMT", DefaultStartAmt)
	startCcy := getenv("START_CCY", DefaultStartCcy)

	log.Printf("CFG: INPUT_CSV=%s", inputCSV)
	log.Printf("CFG: OUTPUT_CSV=%s", outputCSV)
	log.Printf("CFG: BLACKLIST_FILE=%s", blacklistFile)
	log.Printf("CFG: OKX_INSTRUMENTS_URL=%s", instrumentsURL)
	log.Printf("CFG: OKX_BOOKS_URL=%s", booksURL)
	log.Printf("CFG: HTTP_TIMEOUT=%s", timeout)
	log.Printf("CFG: START_AMT=%.8f", startAmt)
	log.Printf("CFG: START_CCY=%s", startCcy)

	// Load blacklist
	bl, err := loadBlacklist(blacklistFile)
	if err != nil {
		log.Fatalf("ERR: load blacklist: %v", err)
	}
	log.Printf("OK: loaded blacklist: %d symbols", len(bl))

	// Load OKX rules
	rules, err := loadRulesOKX(instrumentsURL, timeout)
	if err != nil {
		log.Fatalf("ERR: load OKX instruments: %v", err)
	}
	log.Printf("OK: OKX instruments loaded: %d", len(rules))

	// Build startup blacklist from reliable fields
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

	// Filter routes for START_AMT (e.g. 25 USDT)
	if err := buildOutputCSVOKX(inputCSV, outputCSV, startAmt, strings.TrimSpace(startCcy), rules, bl, booksURL, timeout); err != nil {
		log.Fatalf("ERR: build output csv: %v", err)
	}

	log.Printf("DONE")
}





INPUT_CSV=triangles_usdt_routes_okx.csv
OUTPUT_CSV=triangles_usdt_routes_market_okx.csv
BLACKLIST_FILE=blacklist_symbols_okx.txt

OKX_INSTRUMENTS_URL=https://www.okx.com/api/v5/public/instruments?instType=SPOT
OKX_BOOKS_URL=https://www.okx.com/api/v5/market/books

HTTP_TIMEOUT=25s
START_AMT=25
START_CCY=USDT




