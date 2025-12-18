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
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

// ===== MEXC exchangeInfo (берем только нужное) =====

type mexcExchangeInfo struct {
	Symbols []struct {
		Symbol     string `json:"symbol"`
		BaseAsset  string `json:"baseAsset"`
		QuoteAsset string `json:"quoteAsset"`
		Status     string `json:"status"`
		// На MEXC часто есть эти поля; если вдруг нет — просто будут нулевые.
		IsSpotTradingAllowed bool `json:"isSpotTradingAllowed"`
	} `json:"symbols"`
}

// ===== Модель рынка =====

type Market struct {
	Symbol string
	Base   string
	Quote  string
}

func (m Market) Key() string {
	// уникальный ключ рынка по направлению
	return m.Base + "/" + m.Quote
}

func normAsset(a string) string {
	return strings.ToUpper(strings.TrimSpace(a))
}

// ===== Треугольник в твоем CSV формате =====

type Triangle struct {
	Base1  string
	Quote1 string
	Base2  string
	Quote2 string
	Base3  string
	Quote3 string

	Symbol1 string
	Symbol2 string
	Symbol3 string
}

func (t Triangle) symbolsSortedKey() string {
	s := []string{t.Symbol1, t.Symbol2, t.Symbol3}
	sort.Strings(s)
	return strings.Join(s, "|")
}

// ===== HTTP =====

func httpGetJSON(url string, out any) error {
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("http get: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return fmt.Errorf("http status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(out); err != nil {
		return fmt.Errorf("decode json: %w", err)
	}
	return nil
}

// ===== Построение треугольников =====

func buildMarkets(exchangeInfoURL string, onlyEnabled bool, onlySpot bool) ([]Market, error) {
	var info mexcExchangeInfo
	if err := httpGetJSON(exchangeInfoURL, &info); err != nil {
		return nil, err
	}

	markets := make([]Market, 0, len(info.Symbols))
	for _, s := range info.Symbols {
		base := normAsset(s.BaseAsset)
		quote := normAsset(s.QuoteAsset)

		if base == "" || quote == "" {
			continue
		}

		if onlyEnabled && strings.ToUpper(strings.TrimSpace(s.Status)) != "ENABLED" {
			continue
		}
		if onlySpot && !s.IsSpotTradingAllowed {
			// если у символа нет поля в ответе — оно false, тогда ты сам решаешь:
			// хочешь жестко отсеивать -> оставляем как есть.
			continue
		}

		m := Market{
			Symbol: strings.TrimSpace(s.Symbol),
			Base:   base,
			Quote: quote,
		}
		if m.Symbol == "" {
			// на всякий: если биржа вернула пустой symbol, соберем сами
			m.Symbol = base + quote
		}
		markets = append(markets, m)
	}

	return markets, nil
}

func generateTriangles(markets []Market) []Triangle {
	// 1) marketMap: base/quote -> symbol
	marketMap := make(map[string]string, len(markets))
	// 2) quotesByBase: base -> set(quotes)
	quotesByBase := make(map[string]map[string]struct{})
	// 3) exists both directions helper
	exists := func(base, quote string) (symbol string, ok bool) {
		s, ok := marketMap[base+"/"+quote]
		return s, ok
	}

	for _, m := range markets {
		marketMap[m.Base+"/"+m.Quote] = m.Symbol
		if _, ok := quotesByBase[m.Base]; !ok {
			quotesByBase[m.Base] = make(map[string]struct{})
		}
		quotesByBase[m.Base][m.Quote] = struct{}{}
	}

	out := make([]Triangle, 0, 20000)
	seen := make(map[string]struct{}, 20000)

	// Перебор базового актива A и двух котировок X,Y для рынков A/X и A/Y
	for A, quotesSet := range quotesByBase {
		quotes := make([]string, 0, len(quotesSet))
		for q := range quotesSet {
			quotes = append(quotes, q)
		}
		sort.Strings(quotes)

		for i := 0; i < len(quotes); i++ {
			for j := i + 1; j < len(quotes); j++ {
				X := quotes[i]
				Y := quotes[j]

				// должны существовать A/X и A/Y
				s1, ok1 := exists(A, X)
				s3, ok3 := exists(A, Y)
				if !ok1 || !ok3 {
					continue
				}

				// Между X и Y должен быть рынок в любом направлении
				// Вариант 1: X/Y
				if s2, ok2 := exists(X, Y); ok2 {
					t := Triangle{
						Base1:  A, Quote1: X,
						Base2:  X, Quote2: Y,
						Base3:  A, Quote3: Y,
						Symbol1: s1, Symbol2: s2, Symbol3: s3,
					}
					key := t.symbolsSortedKey()
					if _, ok := seen[key]; !ok {
						seen[key] = struct{}{}
						out = append(out, t)
					}
				}

				// Вариант 2: Y/X (если X/Y нет, либо даже если есть — оставим оба? обычно хватает одного)
				// Чтобы не плодить дубликаты, добавляем только если X/Y НЕ было.
				if _, ok2 := exists(X, Y); !ok2 {
					if s2rev, ok2rev := exists(Y, X); ok2rev {
						t := Triangle{
							Base1:  A, Quote1: X,
							Base2:  Y, Quote2: X, // обратное направление реального рынка
							Base3:  A, Quote3: Y,
							Symbol1: s1, Symbol2: s2rev, Symbol3: s3,
						}
						key := t.symbolsSortedKey()
						if _, ok := seen[key]; !ok {
							seen[key] = struct{}{}
							out = append(out, t)
						}
					}
				}
			}
		}
	}

	// Стабильная сортировка для git-диффов
	sort.Slice(out, func(i, j int) bool {
		a := out[i]
		b := out[j]
		ai := a.Base1 + a.Quote1 + "|" + a.Base2 + a.Quote2 + "|" + a.Base3 + a.Quote3
		bi := b.Base1 + b.Quote1 + "|" + b.Base2 + b.Quote2 + "|" + b.Base3 + b.Quote3
		return ai < bi
	})

	return out
}

// ===== CSV writer =====

func writeCSV(path string, triangles []Triangle, withSymbols bool) error {
	tmp := path + ".tmp"

	f, err := os.Create(tmp)
	if err != nil {
		return fmt.Errorf("create tmp: %w", err)
	}
	defer func() {
		_ = f.Close()
	}()

	w := csv.NewWriter(f)

	// header
	if withSymbols {
		_ = w.Write([]string{
			"base1", "quote1", "base2", "quote2", "base3", "quote3",
			"symbol1", "symbol2", "symbol3",
		})
	} else {
		_ = w.Write([]string{"base1", "quote1", "base2", "quote2", "base3", "quote3"})
	}

	for _, t := range triangles {
		if withSymbols {
			if err := w.Write([]string{
				t.Base1, t.Quote1, t.Base2, t.Quote2, t.Base3, t.Quote3,
				t.Symbol1, t.Symbol2, t.Symbol3,
			}); err != nil {
				return fmt.Errorf("write row: %w", err)
			}
		} else {
			if err := w.Write([]string{
				t.Base1, t.Quote1, t.Base2, t.Quote2, t.Base3, t.Quote3,
			}); err != nil {
				return fmt.Errorf("write row: %w", err)
			}
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		return fmt.Errorf("csv flush: %w", err)
	}

	if err := f.Close(); err != nil {
		return fmt.Errorf("close tmp: %w", err)
	}

	if err := os.Rename(tmp, path); err != nil {
		return fmt.Errorf("rename tmp -> out: %w", err)
	}
	return nil
}

func main() {
	var (
		outFile      = flag.String("out", "cmd/cryptarb/triangles_markets.csv", "куда писать CSV")
		exchangeInfo = flag.String("url", "https://api.mexc.com/api/v3/exchangeInfo", "URL exchangeInfo")
		onlyEnabled  = flag.Bool("enabled", true, "брать только ENABLED")
		onlySpot     = flag.Bool("spot", false, "брать только isSpotTradingAllowed=true (если поле есть)")
		withSymbols  = flag.Bool("with-symbols", false, "добавить в CSV колонки symbol1,symbol2,symbol3")
	)
	flag.Parse()

	markets, err := buildMarkets(*exchangeInfo, *onlyEnabled, *onlySpot)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERR: exchangeInfo: %v\n", err)
		os.Exit(1)
	}

	triangles := generateTriangles(markets)

	if err := writeCSV(*outFile, triangles, *withSymbols); err != nil {
		fmt.Fprintf(os.Stderr, "ERR: write csv: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("OK: markets=%d triangles=%d out=%s\n", len(markets), len(triangles), *outFile)
}


