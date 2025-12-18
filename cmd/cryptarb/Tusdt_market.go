package main

import (
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
	InputCSV  = "triangles_usdt_routes.csv"
	OutputCSV = "triangles_usdt_routes_market.csv"
	BaseURL   = "https://api.mexc.com"
)

type exchangeInfo struct {
	Symbols []symbolInfo `json:"symbols"`
}

type symbolInfo struct {
	Symbol               string   `json:"symbol"`
	Status               string   `json:"status"`
	BaseAsset            string   `json:"baseAsset"`
	QuoteAsset           string   `json:"quoteAsset"`
	IsSpotTradingAllowed bool     `json:"isSpotTradingAllowed"`
	OrderTypes           []string `json:"orderTypes"`
}

func hasMarket(orderTypes []string) bool {
	for _, t := range orderTypes {
		if strings.EqualFold(strings.TrimSpace(t), "MARKET") {
			return true
		}
	}
	return false
}

func marketOk(r symbolInfo) bool {
	if strings.ToUpper(r.Status) != "TRADING" {
		return false
	}
	if !r.IsSpotTradingAllowed {
		return false
	}
	if !hasMarket(r.OrderTypes) {
		return false
	}
	return true
}

func loadRules() (map[string]symbolInfo, error) {
	client := &http.Client{Timeout: 20 * time.Second}
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

func colIndex(header []string, name string) int {
	name = strings.ToLower(strings.TrimSpace(name))
	for i, h := range header {
		if strings.ToLower(strings.TrimSpace(h)) == name {
			return i
		}
	}
	return -1
}

// Достаём 3 символа из строки: либо symbol1..3, либо base/quote -> symbol
func getSymbolsFromRow(header, row []string) (string, string, string, error) {
	iS1 := colIndex(header, "symbol1")
	iS2 := colIndex(header, "symbol2")
	iS3 := colIndex(header, "symbol3")
	if iS1 >= 0 && iS2 >= 0 && iS3 >= 0 {
		s1 := strings.TrimSpace(row[iS1])
		s2 := strings.TrimSpace(row[iS2])
		s3 := strings.TrimSpace(row[iS3])
		if s1 == "" || s2 == "" || s3 == "" {
			return "", "", "", fmt.Errorf("empty symbol in row")
		}
		return s1, s2, s3, nil
	}

	iB1 := colIndex(header, "base1")
	iQ1 := colIndex(header, "quote1")
	iB2 := colIndex(header, "base2")
	iQ2 := colIndex(header, "quote2")
	iB3 := colIndex(header, "base3")
	iQ3 := colIndex(header, "quote3")
	if iB1 >= 0 && iQ1 >= 0 && iB2 >= 0 && iQ2 >= 0 && iB3 >= 0 && iQ3 >= 0 {
		s1 := strings.TrimSpace(row[iB1]) + strings.TrimSpace(row[iQ1])
		s2 := strings.TrimSpace(row[iB2]) + strings.TrimSpace(row[iQ2])
		s3 := strings.TrimSpace(row[iB3]) + strings.TrimSpace(row[iQ3])
		if s1 == "" || s2 == "" || s3 == "" {
			return "", "", "", fmt.Errorf("empty base/quote -> symbol in row")
		}
		return s1, s2, s3, nil
	}

	return "", "", "", fmt.Errorf("CSV must have either symbol1..3 or base1..quote3")
}

// Переход по маркету: если curr совпал с base или quote, вернём вторую валюту
func nextAsset(curr, base, quote string) (string, bool) {
	if curr == quote {
		return base, true
	}
	if curr == base {
		return quote, true
	}
	return "", false
}

// Ищем любую перестановку 3 ног, которая даёт цикл USDT->X->Y->USDT.
// Возвращаем X,Y и найденный порядок индексов.
func findUSDTcycle(legs [3]symbolInfo) (x, y string, order [3]int, ok bool) {
	perms := [][3]int{
		{0, 1, 2},
		{0, 2, 1},
		{1, 0, 2},
		{1, 2, 0},
		{2, 0, 1},
		{2, 1, 0},
	}

	for _, p := range perms {
		curr := "USDT"

		a1, ok1 := nextAsset(curr, legs[p[0]].BaseAsset, legs[p[0]].QuoteAsset)
		if !ok1 {
			continue
		}
		a2, ok2 := nextAsset(a1, legs[p[1]].BaseAsset, legs[p[1]].QuoteAsset)
		if !ok2 {
			continue
		}
		a3, ok3 := nextAsset(a2, legs[p[2]].BaseAsset, legs[p[2]].QuoteAsset)
		if !ok3 || a3 != "USDT" {
			continue
		}

		if a1 == "" || a2 == "" || a1 == "USDT" || a2 == "USDT" {
			continue
		}

		return a1, a2, p, true
	}
	return "", "", [3]int{}, false
}

type triRow struct {
	row  []string
	x, y string // ориентация USDT->x->y->USDT
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

	// key = sorted(X,Y) -> map[orientation]triRow
	type orientMap map[string]triRow
	groups := map[string]orientMap{}

	total := 0
	parsed := 0
	eligible := 0

	for {
		row, err := cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("ERR: read csv: %v", err)
		}
		total++

		s1, s2, s3, err := getSymbolsFromRow(header, row)
		if err != nil {
			continue
		}

		r1, ok1 := rules[s1]
		r2, ok2 := rules[s2]
		r3, ok3 := rules[s3]
		if !ok1 || !ok2 || !ok3 {
			continue
		}

		// MARKET/spot/trading фильтр на все 3 ноги
		if !marketOk(r1) || !marketOk(r2) || !marketOk(r3) {
			continue
		}
		eligible++

		legs := [3]symbolInfo{r1, r2, r3}
		x, y, _, ok := findUSDTcycle(legs)
		if !ok {
			continue
		}
		parsed++

		// группировка по паре монет X/Y
		a := []string{x, y}
		sort.Strings(a)
		key := a[0] + "|" + a[1]
		orient := x + "->" + y

		if groups[key] == nil {
			groups[key] = orientMap{}
		}
		if _, exists := groups[key][orient]; !exists {
			groups[key][orient] = triRow{row: row, x: x, y: y}
		}
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

	written := 0
	pairs := 0

	for _, om := range groups {
		if len(om) < 2 {
			continue
		}

		// найдём разворот X<->Y
		var rows []triRow
		for _, tr := range om {
			rows = append(rows, tr)
		}

		// выбираем пару tr1/tr2 где tr2 = (y->x)
		tr1 := rows[0]
		var tr2 *triRow
		for i := 1; i < len(rows); i++ {
			if rows[i].x == tr1.y && rows[i].y == tr1.x {
				tr2 = &rows[i]
				break
			}
		}
		if tr2 == nil {
			continue
		}

		if err := cw.Write(tr1.row); err != nil {
			log.Fatalf("ERR: write csv: %v", err)
		}
		if err := cw.Write(tr2.row); err != nil {
			log.Fatalf("ERR: write csv: %v", err)
		}
		written += 2
		pairs++
	}

	log.Printf("OK: read=%d marketEligible=%d parsedUSDT=%d paired=%d writtenRows=%d -> %s",
		total, eligible, parsed, pairs, written, OutputCSV)
}
