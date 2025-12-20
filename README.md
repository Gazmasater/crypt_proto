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
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const (
	DefaultInputCSV      = "triangles_usdt_routes_okx.csv"
	DefaultOutputCSV     = "triangles_usdt_routes_market_okx.csv"
	DefaultBlacklistFile = "blacklist_symbols_okx.txt"

	OKXInstrumentsURL = "https://www.okx.com/api/v5/public/instruments?instType=SPOT"
)

type okxInstrumentsResp struct {
	Code string          `json:"code"`
	Msg  string          `json:"msg"`
	Data []okxInstrument `json:"data"`
}

type okxInstrument struct {
	InstID   string `json:"instId"`   // e.g. BTC-USDT
	InstType string `json:"instType"` // SPOT
	BaseCcy  string `json:"baseCcy"`
	QuoteCcy string `json:"quoteCcy"`
	State    string `json:"state"` // live / suspend / ...

	// Trading rules (strings):
	LotSz  string `json:"lotSz"`  // size increment (base)
	MinSz  string `json:"minSz"`  // min size (base)
	TickSz string `json:"tickSz"` // price increment
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	in := flag.String("in", DefaultInputCSV, "input routes CSV (triangles_usdt_routes*.csv)")
	out := flag.String("out", DefaultOutputCSV, "output filtered CSV")
	blFile := flag.String("blacklist", DefaultBlacklistFile, "blacklist file path")
	api := flag.String("api", OKXInstrumentsURL, "OKX instruments endpoint")
	timeout := flag.Duration("timeout", 25*time.Second, "HTTP timeout")
	flag.Parse()

	// 1) Load existing blacklist (or empty if missing)
	bl, err := loadBlacklist(*blFile)
	if err != nil {
		log.Fatalf("ERR: load blacklist: %v", err)
	}
	log.Printf("OK: loaded blacklist: %d symbols (%s)", len(bl), *blFile)

	// 2) Load OKX instruments (rules)
	rules, err := loadRulesOKX(*api, *timeout)
	if err != nil {
		log.Fatalf("ERR: load OKX instruments: %v", err)
	}
	log.Printf("OK: OKX instruments loaded: %d", len(rules))

	// 3) Build "startup" blacklist from reliable fields ONLY
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

	// 4) Save merged blacklist
	if err := saveBlacklist(*blFile, bl); err != nil {
		log.Fatalf("ERR: save blacklist: %v", err)
	}
	log.Printf("OK: saved blacklist -> %s", *blFile)

	// 5) Filter routes and write output CSV
	if err := buildOutputCSVOKX(*in, *out, rules, bl); err != nil {
		log.Fatalf("ERR: build output csv: %v", err)
	}
}

func loadRulesOKX(url string, timeout time.Duration) (map[string]okxInstrument, error) {
	client := &http.Client{Timeout: timeout}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	// Иногда помогает от странных блоков/проксей:
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
		return nil, fmt.Errorf("OKX code=%s msg=%s", strings.TrimSpace(r.Code), strings.TrimSpace(r.Msg))
	}

	m := make(map[string]okxInstrument, len(r.Data))
	for _, s := range r.Data {
		instID := strings.TrimSpace(s.InstID)
		if instID == "" {
			continue
		}
		m[instID] = s
	}
	return m, nil
}

func buildOutputCSVOKX(inputCSV, outputCSV string, rules map[string]okxInstrument, bl map[string]string) error {
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
	if iL1 < 0 || iL2 < 0 || iL3 < 0 {
		return fmt.Errorf("нет колонок leg1_symbol/leg2_symbol/leg3_symbol в CSV")
	}

	// Add OKX rule/precision columns
	header = append(header,
		"leg1_lot_sz", "leg2_lot_sz", "leg3_lot_sz",
		"leg1_min_sz", "leg2_min_sz", "leg3_min_sz",
		"leg1_tick_sz", "leg2_tick_sz", "leg3_tick_sz",
		"leg1_qty_dp", "leg2_qty_dp", "leg3_qty_dp",
		"leg1_price_dp", "leg2_price_dp", "leg3_price_dp",
	)

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

		// eligibility
		if !marketOkOKX(r1) || !marketOkOKX(r2) || !marketOkOKX(r3) {
			skippedNotEligible++
			continue
		}

		// qty precision from lotSz
		dpQty1 := decimalsFromStep(r1.LotSz)
		dpQty2 := decimalsFromStep(r2.LotSz)
		dpQty3 := decimalsFromStep(r3.LotSz)

		// price precision from tickSz
		dpPx1 := decimalsFromStep(r1.TickSz)
		dpPx2 := decimalsFromStep(r2.TickSz)
		dpPx3 := decimalsFromStep(r3.TickSz)

		row = append(row,
			strings.TrimSpace(r1.LotSz), strings.TrimSpace(r2.LotSz), strings.TrimSpace(r3.LotSz),
			strings.TrimSpace(r1.MinSz), strings.TrimSpace(r2.MinSz), strings.TrimSpace(r3.MinSz),
			strings.TrimSpace(r1.TickSz), strings.TrimSpace(r2.TickSz), strings.TrimSpace(r3.TickSz),
			fmt.Sprintf("%d", dpQty1), fmt.Sprintf("%d", dpQty2), fmt.Sprintf("%d", dpQty3),
			fmt.Sprintf("%d", dpPx1), fmt.Sprintf("%d", dpPx2), fmt.Sprintf("%d", dpPx3),
		)

		if err := cw.Write(row); err != nil {
			return fmt.Errorf("write row: %w", err)
		}
		written++
	}

	log.Printf(
		"OK: read=%d written=%d skippedNoSymbol=%d skippedBlacklisted=%d skippedNotEligible=%d -> %s",
		read, written, skippedNoSymbol, skippedBlacklisted, skippedNotEligible, outputCSV,
	)

	if written == 0 {
		log.Printf("WARN: output is empty. Possible reasons:")
		log.Printf(" - input CSV symbols are not OKX instId (must be like BTC-USDT)")
		log.Printf(" - too broad blacklist (%s)", filepath.Base(outputCSV))
		log.Printf(" - OKX instruments fetch blocked (then rules map is small/empty)")
	}

	return nil
}

// ======= Eligibility / reasons (OKX) =======

func marketOkOKX(s okxInstrument) bool {
	if s.InstType != "" && !strings.EqualFold(strings.TrimSpace(s.InstType), "SPOT") {
		return false
	}
	if !strings.EqualFold(strings.TrimSpace(s.State), "live") {
		return false
	}
	// Чтобы потом не ломаться на округлениях
	if strings.TrimSpace(s.LotSz) == "" {
		return false
	}
	if strings.TrimSpace(s.MinSz) == "" {
		return false
	}
	if strings.TrimSpace(s.TickSz) == "" {
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
	return strings.Join(reasons, ", ")
}

// ======= Precision helpers =======

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
			return bl, nil // first run: empty
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

func ensureDirForFile(path string) error {
	dir := filepath.Dir(path)
	if dir == "." || dir == "/" || dir == "" {
		return nil
	}
	return os.MkdirAll(dir, 0o755)
}



