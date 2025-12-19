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
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

const (
	InputCSV      = "triangles_usdt_routes.csv"
	OutputCSV     = "triangles_usdt_routes_market.csv"
	BlacklistFile = "blacklist_symbols.txt"
	BaseURL       = "https://api.mexc.com"
)

type exchangeInfo struct {
	Symbols []symbolInfo `json:"symbols"`
}

type symbolInfo struct {
	Symbol string `json:"symbol"`
	Status string `json:"status"`

	OrderTypes  []string `json:"orderTypes"`
	Permissions []string `json:"permissions"`
	St          bool     `json:"st"`

	BaseSizePrec string `json:"baseSizePrecision"`

	QuoteAmountPrecisionMarket string `json:"quoteAmountPrecisionMarket"`
	QuoteAmountPrecision       string `json:"quoteAmountPrecision"`
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	// 1) Load existing blacklist (or empty if missing)
	bl, err := loadBlacklist(BlacklistFile)
	if err != nil {
		log.Fatalf("ERR: load blacklist: %v", err)
	}
	log.Printf("OK: loaded blacklist: %d symbols (%s)", len(bl), BlacklistFile)

	// 2) Load exchangeInfo
	rules, err := loadRules()
	if err != nil {
		log.Fatalf("ERR: load exchangeInfo: %v", err)
	}
	log.Printf("OK: exchangeInfo symbols: %d", len(rules))

	// 3) Build "startup" blacklist from reliable fields ONLY
	added := 0
	for sym, s := range rules {
		if marketOk(s) {
			continue
		}
		reason := notOkReason(s)
		if reason == "" {
			reason = "not eligible"
		}
		if _, exists := bl[sym]; !exists {
			bl[sym] = reason
			added++
		}
	}
	log.Printf("OK: startup blacklist added=%d total=%d", added, len(bl))

	// 4) Save merged blacklist
	if err := saveBlacklist(BlacklistFile, bl); err != nil {
		log.Fatalf("ERR: save blacklist: %v", err)
	}
	log.Printf("OK: saved blacklist -> %s", BlacklistFile)

	// 5) Filter triangles and write output CSV
	if err := buildOutputCSV(rules, bl); err != nil {
		log.Fatalf("ERR: build output csv: %v", err)
	}
}

func buildOutputCSV(rules map[string]symbolInfo, bl map[string]string) error {
	in, err := os.Open(InputCSV)
	if err != nil {
		return fmt.Errorf("open %s: %w", InputCSV, err)
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
	if iL1 < 0 || iL2 < 0 || iL3 < 0 {
		return fmt.Errorf("нет колонок leg1_symbol/leg2_symbol/leg3_symbol в CSV")
	}

	// Add precision columns
	header = append(header,
		"leg1_qty_dp", "leg2_qty_dp", "leg3_qty_dp",
		"leg1_quote_dp_market", "leg2_quote_dp_market", "leg3_quote_dp_market",
	)

	out, err := os.Create(OutputCSV)
	if err != nil {
		return fmt.Errorf("create %s: %w", OutputCSV, err)
	}
	defer out.Close()

	cw := csv.NewWriter(out)
	defer cw.Flush()

	if err := cw.Write(header); err != nil {
		return fmt.Errorf("write header: %w", err)
	}

	read := 0
	written := 0
	skippedNoSymbol := 0
	skippedBlacklisted := 0
	skippedNotEligible := 0

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

		// if any leg is in blacklist -> skip
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

		// reliable eligibility (same as blacklist criteria)
		if !marketOk(r1) || !marketOk(r2) || !marketOk(r3) {
			skippedNotEligible++
			continue
		}

		// qty precision (base)
		dpQty1 := decimalsFromStep(r1.BaseSizePrec)
		dpQty2 := decimalsFromStep(r2.BaseSizePrec)
		dpQty3 := decimalsFromStep(r3.BaseSizePrec)

		// quote precision for MARKET (quote)
		dpQm1 := decimalsFromStep(quoteMarketStep(r1))
		dpQm2 := decimalsFromStep(quoteMarketStep(r2))
		dpQm3 := decimalsFromStep(quoteMarketStep(r3))

		row = append(row,
			fmt.Sprintf("%d", dpQty1), fmt.Sprintf("%d", dpQty2), fmt.Sprintf("%d", dpQty3),
			fmt.Sprintf("%d", dpQm1), fmt.Sprintf("%d", dpQm2), fmt.Sprintf("%d", dpQm3),
		)

		if err := cw.Write(row); err != nil {
			return fmt.Errorf("write row: %w", err)
		}
		written++
	}

	log.Printf(
		"OK: read=%d written=%d skippedNoSymbol=%d skippedBlacklisted=%d skippedNotEligible=%d -> %s",
		read, written, skippedNoSymbol, skippedBlacklisted, skippedNotEligible, OutputCSV,
	)

	if written == 0 {
		log.Printf("WARN: output is empty. Most likely your blacklist is too broad or input csv symbols mismatch exchangeInfo.")
		log.Printf("TIP: temporarily move %s away and rerun to see baseline written>0.", BlacklistFile)
	}

	return nil
}

func loadRules() (map[string]symbolInfo, error) {
	client := &http.Client{Timeout: 25 * time.Second}
	resp, err := client.Get(BaseURL + "/api/v3/exchangeInfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("exchangeInfo %d: %s", resp.StatusCode, string(b))
	}

	var info exchangeInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, err
	}

	m := make(map[string]symbolInfo, len(info.Symbols))
	for _, s := range info.Symbols {
		sym := strings.TrimSpace(s.Symbol)
		if sym == "" {
			continue
		}
		m[sym] = s
	}
	return m, nil
}

// ======= Eligibility / reasons =======

// Reliable filter you already validated with Postman:
// status=="1", st==false, permissions contains "SPOT", orderTypes contains "MARKET"
func marketOk(s symbolInfo) bool {
	if strings.TrimSpace(s.Status) != "1" {
		return false
	}
	if s.St {
		return false
	}
	if !hasPerm(s.Permissions, "SPOT") {
		return false
	}
	if !hasMarket(s.OrderTypes) {
		return false
	}
	return true
}

func notOkReason(s symbolInfo) string {
	var reasons []string
	if strings.TrimSpace(s.Status) != "1" {
		reasons = append(reasons, "status!=1")
	}
	if s.St {
		reasons = append(reasons, "st=true")
	}
	if !hasPerm(s.Permissions, "SPOT") {
		reasons = append(reasons, "no SPOT perm")
	}
	if !hasMarket(s.OrderTypes) {
		reasons = append(reasons, "no MARKET type")
	}
	return strings.Join(reasons, ", ")
}

func hasMarket(orderTypes []string) bool {
	for _, t := range orderTypes {
		if strings.EqualFold(strings.TrimSpace(t), "MARKET") {
			return true
		}
	}
	return false
}

func hasPerm(perms []string, want string) bool {
	for _, p := range perms {
		if strings.EqualFold(strings.TrimSpace(p), want) {
			return true
		}
	}
	return false
}

// ======= Precision helpers =======

func quoteMarketStep(s symbolInfo) string {
	if strings.TrimSpace(s.QuoteAmountPrecisionMarket) != "" {
		return s.QuoteAmountPrecisionMarket
	}
	return s.QuoteAmountPrecision
}

// "0.000001" -> 6; "1" -> 0; "" -> -1
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

// ======= CSV helpers =======

func colIndex(header []string, name string) int {
	name = strings.ToLower(strings.TrimSpace(name))
	for i, h := range header {
		if strings.ToLower(strings.TrimSpace(h)) == name {
			return i
		}
	}
	return -1
}

// ======= Blacklist I/O =======

// File format: SYMBOL<TAB>reason
// Lines starting with # are comments.
func loadBlacklist(path string) (map[string]string, error) {
	bl := make(map[string]string)

	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			// first run: empty blacklist
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
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return bl, nil
}

func saveBlacklist(path string, bl map[string]string) error {
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


