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
	InputCSV  = "triangles_markets_usdt_routes.csv"
	OutputCSV = "triangles_markets_usdt_routes_market.csv"
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
		if strings.EqualFold(t, "MARKET") {
			return true
		}
	}
	return false
}

func loadRules() (map[string]symbolInfo, error) {
	client := &http.Client{Timeout: 15 * time.Second}
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

func getSymbolsFromRow(header, row []string) (string, string, string, error) {
	// вариант 1: symbol1,symbol2,symbol3
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

	// вариант 2: base/quote пары
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

// восстановим маршрут по базам/квотам: стартуем с USDT, идём 3 ноги, должны вернуться в USDT
// возвращаем X, Y и признак ориентации (XY или YX относительно USDT)
func routeXY(r1, r2, r3 symbolInfo) (x string, y string, ok bool) {
	curr := "USDT"
	next := func(curr string, base string, quote string) (string, bool) {
		if curr == quote {
			return base, true
		}
		if curr == base {
			return quote, true
		}
		return "", false
	}

	a1, ok1 := next(curr, r1.BaseAsset, r1.QuoteAsset)
	if !ok1 {
		return "", "", false
	}
	a2, ok2 := next(a1, r2.BaseAsset, r2.QuoteAsset)
	if !ok2 {
		return "", "", false
	}
	a3, ok3 := next(a2, r3.BaseAsset, r3.QuoteAsset)
	if !ok3 || a3 != "USDT" {
		return "", "", false
	}

	// curr=USDT -> a1=X -> a2=Y -> USDT
	if a1 == "USDT" || a2 == "USDT" || a1 == "" || a2 == "" {
		return "", "", false
	}
	return a1, a2, true
}

type triRow struct {
	row     []string
	s1, s2, s3 string
	x, y    string // ориентация USDT->x->y->USDT
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
	// orientation: "X->Y" or "Y->X" (строго по тому, что восстановили)
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

		x, y, ok := routeXY(r1, r2, r3)
		if !ok {
			continue
		}
		parsed++

		// фильтр MARKET по 3 символам (в обе стороны торговля маркетом подразумевается самим наличием MARKET)
		if !marketOk(r1) || !marketOk(r2) || !marketOk(r3) {
			continue
		}
		eligible++

		a := []string{x, y}
		sort.Strings(a)
		key := a[0] + "|" + a[1]
		orient := x + "->" + y

		if groups[key] == nil {
			groups[key] = orientMap{}
		}
		// сохраняем первую встреченную строку на эту ориентацию
		if _, exists := groups[key][orient]; !exists {
			groups[key][orient] = triRow{row: row, s1: s1, s2: s2, s3: s3, x: x, y: y}
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
		// хотим именно обе ориентации
		// (внутри om ключи типа "X->Y" и "Y->X")
		if len(om) < 2 {
			continue
		}

		// найдём две разные ориентации
		var rows []triRow
		for _, tr := range om {
			rows = append(rows, tr)
		}
		if len(rows) < 2 {
			continue
		}

		// на всякий: убедимся что это ровно X<->Y разворот
		// берём первые две уникальные
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

		// пишем обе строки
		if err := cw.Write(tr1.row); err != nil {
			log.Fatalf("ERR: write csv: %v", err)
		}
		if err := cw.Write(tr2.row); err != nil {
			log.Fatalf("ERR: write csv: %v", err)
		}
		written += 2
		pairs++
	}

	log.Printf("OK: read=%d parsedUSDT=%d marketEligible=%d paired=%d writtenRows=%d -> %s",
		total, parsed, eligible, pairs, written, OutputCSV)
}

gaz358@gaz358-BOD-WXX9:~/myprog/crypt_proto/cmd/cryptarb$ go run .
2025/12/18 17:33:34.899074 OK: read=764 parsedUSDT=0 marketEligible=0 paired=0 writtenRows=0 -> triangles_usdt_routes_market.csv
