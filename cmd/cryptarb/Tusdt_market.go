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
	InputCSV    = "triangles_usdt_routes.csv"
	OutputCSV   = "triangles_usdt_routes_market.csv"
	Blacklist   = "blacklist_symbols.txt"
	BaseURL     = "https://api.mexc.com"
	HTTPTimeout = 25 * time.Second
)

type exchangeInfo struct {
	Symbols []symbolInfo `json:"symbols"`
}

type symbolInfo struct {
	Symbol string `json:"symbol"`
	Status string `json:"status"`

	// фильтр "разрешено ли спот-торговать через API"
	IsSpotTradingAllowed bool `json:"isSpotTradingAllowed"`

	OrderTypes  []string `json:"orderTypes"`
	Permissions []string `json:"permissions"`
	St          bool     `json:"st"`

	// точности
	BaseSizePrec string `json:"baseSizePrecision"`

	QuoteAmountPrecisionMarket string `json:"quoteAmountPrecisionMarket"`
	QuoteAmountPrecision       string `json:"quoteAmountPrecision"`

	// Если ты BUY делаешь через quoteQty/quoteOrderQty на MARKET — это КРИТИЧНО:
	// если false, то биржа может не принять quoteOrderQty.
	QuoteOrderQtyMarketAllowed bool `json:"quoteOrderQtyMarketAllowed"`
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
	frac := parts[1]
	frac = strings.TrimRight(frac, "0")
	return len(frac)
}

// Выбор точности quote для MARKET: сначала quoteAmountPrecisionMarket, иначе fallback quoteAmountPrecision
func quoteMarketStep(s symbolInfo) string {
	if strings.TrimSpace(s.QuoteAmountPrecisionMarket) != "" {
		return s.QuoteAmountPrecisionMarket
	}
	return s.QuoteAmountPrecision
}

// ТВОИ условия + фильтр API.
// И ещё важный флаг под BUY через quoteQty: QuoteOrderQtyMarketAllowed.
// Если ты не используешь quoteQty для BUY — можешь убрать этот пункт.
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
	if !s.IsSpotTradingAllowed {
		return false
	}
	// BUY через quoteQty: чтобы реально работало без сюрпризов
	if !s.QuoteOrderQtyMarketAllowed {
		return false
	}
	return true
}

// Сформировать reason для блэклиста (почему НЕ подходит)
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
	if !s.IsSpotTradingAllowed {
		reasons = append(reasons, "isSpotTradingAllowed=false")
	}
	if !s.QuoteOrderQtyMarketAllowed {
		reasons = append(reasons, "quoteOrderQtyMarketAllowed=false")
	}

	if len(reasons) == 0 {
		return "unknown"
	}
	return strings.Join(reasons, ", ")
}

func loadRules() (map[string]symbolInfo, error) {
	client := &http.Client{Timeout: HTTPTimeout}
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
		if s.Symbol == "" {
			continue
		}
		m[strings.TrimSpace(s.Symbol)] = s
	}
	return m, nil
}

// Шаг 1: создать "стартовый" блэклист на основе exchangeInfo.
// Возвращает map[symbol]reason и счётчики.
func buildInitialBlacklist(rules map[string]symbolInfo, path string) (map[string]string, int, int, error) {
	bl := make(map[string]string, 4096)

	total := 0
	bad := 0
	for sym, s := range rules {
		total++
		if marketOk(s) {
			continue
		}
		bad++
		bl[sym] = notOkReason(s)
	}

	// Записываем файл с нуля (truncate), чтобы это был именно "стартовый" блэклист
	f, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return nil, total, bad, err
	}
	defer f.Close()

	bw := bufio.NewWriter(f)
	defer bw.Flush()

	fmt.Fprintf(bw, "# generated: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Fprintf(bw, "# format: SYMBOL<TAB>reason\n")

	// чтобы было стабильно — отсортируем
	keys := make([]string, 0, len(bl))
	for sym := range bl {
		keys = append(keys, sym)
	}
	sort.Strings(keys)

	for _, sym := range keys {
		fmt.Fprintf(bw, "%s\t%s\n", sym, bl[sym])
	}

	return bl, total, bad, nil
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	// 1) rules
	rules, err := loadRules()
	if err != nil {
		log.Fatalf("ERR: load exchangeInfo: %v", err)
	}

	// 2) СНАЧАЛА: стартовый блэклист
	bl, total, bad, err := buildInitialBlacklist(rules, Blacklist)
	if err != nil {
		log.Fatalf("ERR: build blacklist: %v", err)
	}
	log.Printf("OK: exchangeInfo symbols=%d blacklisted=%d -> %s", total, bad, Blacklist)

	// 3) читаем входной CSV
	in, err := os.Open(InputCSV)
	if err != nil {
		log.Fatalf("ERR: open %s: %v", InputCSV, err)
	}
	defer in.Close()

	cr := csv.NewReader(in)
	header, err := cr.Read()
	if err != nil {
		log.Fatalf("ERR: read header: %v", err)
	}

	iL1 := colIndex(header, "leg1_symbol")
	iL2 := colIndex(header, "leg2_symbol")
	iL3 := colIndex(header, "leg3_symbol")
	if iL1 < 0 || iL2 < 0 || iL3 < 0 {
		log.Fatalf("ERR: нет колонок leg1_symbol/leg2_symbol/leg3_symbol в CSV")
	}

	// 4) готовим выход
	header = append(header,
		"leg1_qty_dp", "leg2_qty_dp", "leg3_qty_dp",
		"leg1_quote_dp_market", "leg2_quote_dp_market", "leg3_quote_dp_market",
	)

	out, err := os.Create(OutputCSV)
	if err != nil {
		log.Fatalf("ERR: create %s: %v", OutputCSV, err)
	}
	defer out.Close()

	cw := csv.NewWriter(out)
	defer cw.Flush()

	if err := cw.Write(header); err != nil {
		log.Fatalf("ERR: write header: %v", err)
	}

	read := 0
	written := 0
	skippedNoSymbol := 0
	skippedNotEligible := 0
	skippedBlacklisted := 0

	for {
		row, err := cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("ERR: read csv: %v", err)
		}
		read++

		s1 := strings.TrimSpace(row[iL1])
		s2 := strings.TrimSpace(row[iL2])
		s3 := strings.TrimSpace(row[iL3])
		if s1 == "" || s2 == "" || s3 == "" {
			skippedNoSymbol++
			continue
		}

		// если символа нет в exchangeInfo — пропускаем
		r1, ok1 := rules[s1]
		r2, ok2 := rules[s2]
		r3, ok3 := rules[s3]
		if !ok1 || !ok2 || !ok3 {
			skippedNoSymbol++
			continue
		}

		// если попал в стартовый блэклист хотя бы один символ — пропускаем треугольник
		if _, ok := bl[s1]; ok {
			skippedBlacklisted++
			continue
		}
		if _, ok := bl[s2]; ok {
			skippedBlacklisted++
			continue
		}
		if _, ok := bl[s3]; ok {
			skippedBlacklisted++
			continue
		}

		// на всякий случай — повторно проверяем marketOk на 3 символа
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
			log.Fatalf("ERR: write row: %v", err)
		}
		written++
	}

	log.Printf(
		"OK: read=%d written=%d skippedNoSymbol=%d skippedBlacklisted=%d skippedNotEligible=%d -> %s",
		read, written, skippedNoSymbol, skippedBlacklisted, skippedNotEligible, OutputCSV,
	)
}
