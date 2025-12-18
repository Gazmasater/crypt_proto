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
	Symbol      string   `json:"symbol"`
	Status      string   `json:"status"`
	OrderTypes  []string `json:"orderTypes"`
	Permissions []string `json:"permissions"`
	St          bool     `json:"st"`
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

func hasSPOT(perms []string) bool {
	for _, p := range perms {
		if strings.EqualFold(strings.TrimSpace(p), "SPOT") {
			return true
		}
	}
	return false
}

func marketOk(s symbolInfo) bool {
	if strings.TrimSpace(s.Status) != "1" {
		return false
	}
	if !hasSPOT(s.Permissions) {
		return false
	}
	if !hasMarket(s.OrderTypes) {
		return false
	}
	// опционально: если st=true означает special treatment — можно отрезать
	if s.St {
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

		if err := cw.Write(row); err != nil {
			log.Fatalf("ERR: write row: %v", err)
		}
		written++
	}

	log.Printf("OK: read=%d written=%d skippedNoSymbol=%d skippedNotEligible=%d -> %s",
		read, written, skippedNoSymbol, skippedNotEligible, OutputCSV)
}
