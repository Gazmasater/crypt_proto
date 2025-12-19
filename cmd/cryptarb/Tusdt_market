package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	InputCSV  = "triangles_usdt_routes.csv"
	OutputCSV = "triangles_usdt_routes_market.csv"
	BaseURL   = "https://api.mexc.com"
)

type exchangeInfo struct {
	Symbols []symbolInfo `json:"symbols"`
}

type symbolInfo struct {
	Symbol string `json:"symbol"`
	Status string `json:"status"`

	OrderTypes   []string `json:"orderTypes"`
	Permissions  []string `json:"permissions"`
	St           bool     `json:"st"`
	BaseSizePrec string   `json:"baseSizePrecision"`

	QuoteAmountPrecisionMarket string `json:"quoteAmountPrecisionMarket"`
	// на всякий случай: иногда может быть полезно как fallback
	QuoteAmountPrecision string `json:"quoteAmountPrecision"`
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

// Фильтр "пара подходит для торговли MARKET" по твоему требованию:
// status=="1", st==false, permissions содержит "SPOT", orderTypes содержит "MARKET"
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
		m[s.Symbol] = s
	}
	return m, nil
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	rules, err := loadRules()
	if err != nil {
		log.Fatalf("ERR: load exchangeInfo: %v", err)
	}

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

	// Добавляем два вида точности:
	// qty_dp (по baseSizePrecision) и quote_dp_market (по quoteAmountPrecisionMarket / fallback quoteAmountPrecision)
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

		r1, ok1 := rules[s1]
		r2, ok2 := rules[s2]
		r3, ok3 := rules[s3]
		if !ok1 || !ok2 || !ok3 {
			skippedNoSymbol++
			continue
		}

		if !marketOk(r1) || !marketOk(r2) || !marketOk(r3) {
			skippedNotEligible++
			continue
		}

		// qty precision (base)
		dpQty1 := decimalsFromStep(r1.BaseSizePrec)
		dpQty2 := decimalsFromStep(r2.BaseSizePrec)
		dpQty3 := decimalsFromStep(r3.BaseSizePrec)

		// quoteAmount precision for MARKET (quote)
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

	log.Printf("OK: read=%d written=%d skippedNoSymbol=%d skippedNotEligible=%d -> %s",
		read, written, skippedNoSymbol, skippedNotEligible, OutputCSV)
}
