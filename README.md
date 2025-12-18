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
	"context"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type exchangeInfo struct {
	Symbols []struct {
		Symbol     string `json:"symbol"`
		BaseAsset  string `json:"baseAsset"`
		QuoteAsset string `json:"quoteAsset"`
		Status     string `json:"status"`
	} `json:"symbols"`
}

type Market struct {
	Symbol string
	Base   string
	Quote  string
}

type TriangleRow struct {
	Base1, Quote1 string
	Base2, Quote2 string
	Base3, Quote3 string
	Symbol1       string
	Symbol2       string
	Symbol3       string
}

func main() {
	var (
		apiURL   = flag.String("api", "https://api.mexc.com/api/v3/exchangeInfo", "exchangeInfo endpoint")
		outPath  = flag.String("out", "triangles_markets.csv", "output csv path (default: ./triangles_markets.csv)")
		timeout  = flag.Duration("timeout", 20*time.Second, "http timeout")
		limitA   = flag.String("base", "", "optional: generate only for this base asset (e.g. BTC)")
		statusOK = flag.String("status", "", "optional: allow only this status (e.g. 1 or TRADING). empty = allow common statuses")
	)
	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	markets, err := fetchMarkets(ctx, *apiURL, *statusOK)
	if err != nil {
		fatalf("fetch exchangeInfo: %v", err)
	}

	rows := buildTriangles(markets, strings.TrimSpace(*limitA))

	if err := writeCSVAtomic(*outPath, rows); err != nil {
		fatalf("write csv: %v", err)
	}

	fmt.Printf("OK: triangles=%d -> %s\n", len(rows), *outPath)
}

func fetchMarkets(ctx context.Context, url string, onlyStatus string) ([]Market, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	cl := &http.Client{Timeout: 20 * time.Second}
	resp, err := cl.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		return nil, fmt.Errorf("http %d: %s", resp.StatusCode, strings.TrimSpace(string(b)))
	}

	var info exchangeInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, err
	}

	allowStatus := func(s string) bool {
		s = strings.TrimSpace(s)
		if onlyStatus != "" {
			return s == onlyStatus
		}
		// MEXC часто отдаёт "1", у некоторых бирж "TRADING"
		return s == "" || s == "1" || strings.EqualFold(s, "TRADING")
	}

	out := make([]Market, 0, len(info.Symbols))
	for _, it := range info.Symbols {
		if it.Symbol == "" || it.BaseAsset == "" || it.QuoteAsset == "" {
			continue
		}
		if !allowStatus(it.Status) {
			continue
		}
		out = append(out, Market{
			Symbol: it.Symbol,
			Base:   it.BaseAsset,
			Quote: it.QuoteAsset,
		})
	}

	return out, nil
}

// Формат как у тебя на скрине:
// base1/quote1 , base2/quote2 , base3/quote3
// где base1==base3 (A), а quote1 и quote3 — две “котировки” (Q1,Q2),
// и есть рынок между Q1 и Q2.
func buildTriangles(markets []Market, onlyBase string) []TriangleRow {
	// marketByPair: unordered pair key -> market (единственный спот-символ на пару)
	type pairKey struct{ a, b string }
	pk := func(x, y string) pairKey {
		if x < y {
			return pairKey{x, y}
		}
		return pairKey{y, x}
	}

	marketByPair := make(map[pairKey]Market, len(markets))
	quotesByBase := make(map[string][]string) // A -> list of quotes Q where A/Q exists

	for _, m := range markets {
		marketByPair[pk(m.Base, m.Quote)] = m
		quotesByBase[m.Base] = append(quotesByBase[m.Base], m.Quote)
	}

	// нормализуем quotesByBase (unique + sort)
	for base, qs := range quotesByBase {
		sort.Strings(qs)
		qs = uniqueStrings(qs)
		quotesByBase[base] = qs
	}

	var bases []string
	for b := range quotesByBase {
		if onlyBase != "" && b != onlyBase {
			continue
		}
		bases = append(bases, b)
	}
	sort.Strings(bases)

	seen := make(map[string]struct{}, 1024)
	rows := make([]TriangleRow, 0, 4096)

	for _, A := range bases {
		qs := quotesByBase[A]
		if len(qs) < 2 {
			continue
		}

		// перебор пар котировок
		for i := 0; i < len(qs); i++ {
			for j := i + 1; j < len(qs); j++ {
				qx, qy := qs[i], qs[j]

				mMid, ok := marketByPair[pk(qx, qy)]
				if !ok {
					continue
				}

				// делаем quote1=base middle рынка, quote3=quote middle рынка
				quote1 := mMid.Base
				quote3 := mMid.Quote

				// проверяем что A/quote1 и A/quote3 существуют (они есть из qs, но на всякий)
				m1, ok1 := marketByPair[pk(A, quote1)]
				m3, ok3 := marketByPair[pk(A, quote3)]
				if !(ok1 && ok3) {
					continue
				}

				// уникальный ключ (чтобы не плодить дубли)
				key := A + "|" + quote1 + "|" + quote3 + "|" + mMid.Symbol
				if _, exists := seen[key]; exists {
					continue
				}
				seen[key] = struct{}{}

				rows = append(rows, TriangleRow{
					Base1:  m1.Base,
					Quote1: m1.Quote,
					Base2:  mMid.Base,
					Quote2: mMid.Quote,
					Base3:  m3.Base,
					Quote3: m3.Quote,

					// на всякий случай делаем symbol = base+quote
					Symbol1: m1.SymbolIfEmptyConcat(),
					Symbol2: mMid.SymbolIfEmptyConcat(),
					Symbol3: m3.SymbolIfEmptyConcat(),
				})
			}
		}
	}

	// детерминированная сортировка
	sort.Slice(rows, func(i, j int) bool {
		a := rows[i]
		b := rows[j]
		if a.Base1 != b.Base1 {
			return a.Base1 < b.Base1
		}
		if a.Quote1 != b.Quote1 {
			return a.Quote1 < b.Quote1
		}
		if a.Quote3 != b.Quote3 {
			return a.Quote3 < b.Quote3
		}
		return a.Symbol2 < b.Symbol2
	})

	return rows
}

func (m Market) SymbolIfEmptyConcat() string {
	if m.Symbol != "" {
		return m.Symbol
	}
	return m.Base + m.Quote
}

func writeCSVAtomic(path string, rows []TriangleRow) error {
	dir := filepath.Dir(path)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("mkdir: %w", err)
		}
	}

	// temp-файл рядом с целевым (чтобы rename был атомарный)
	tmp, err := os.CreateTemp(dir, filepath.Base(path)+".tmp.*")
	if err != nil {
		return fmt.Errorf("create tmp: %w", err)
	}
	tmpName := tmp.Name()

	cw := csv.NewWriter(tmp)
	// Header
	if err := cw.Write([]string{
		"base1", "quote1",
		"base2", "quote2",
		"base3", "quote3",
		"symbol1", "symbol2", "symbol3",
	}); err != nil {
		tmp.Close()
		_ = os.Remove(tmpName)
		return err
	}

	for _, r := range rows {
		if err := cw.Write([]string{
			r.Base1, r.Quote1,
			r.Base2, r.Quote2,
			r.Base3, r.Quote3,
			r.Symbol1, r.Symbol2, r.Symbol3,
		}); err != nil {
			tmp.Close()
			_ = os.Remove(tmpName)
			return err
		}
	}

	cw.Flush()
	if err := cw.Error(); err != nil {
		tmp.Close()
		_ = os.Remove(tmpName)
		return err
	}

	if err := tmp.Sync(); err != nil {
		tmp.Close()
		_ = os.Remove(tmpName)
		return err
	}
	if err := tmp.Close(); err != nil {
		_ = os.Remove(tmpName)
		return err
	}

	// атомарная подмена
	if err := os.Rename(tmpName, path); err != nil {
		_ = os.Remove(tmpName)
		return err
	}
	return nil
}

func uniqueStrings(in []string) []string {
	if len(in) == 0 {
		return in
	}
	out := make([]string, 0, len(in))
	prev := ""
	for i, s := range in {
		if i == 0 || s != prev {
			out = append(out, s)
		}
		prev = s
	}
	return out
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "ERR: "+format+"\n", args...)
	os.Exit(1)
}


gaz358@gaz358-BOD-WXX9:~/myprog/crypt_proto/cmd/cryptarb$ go run .
OK: triangles=467 -> triangles_markets.csv

