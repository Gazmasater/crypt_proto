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

// OKX: GET /api/v5/public/instruments?instType=SPOT
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
}

func main() {
	var (
		apiURL  = flag.String("api", "https://www.okx.com/api/v5/public/instruments?instType=SPOT", "OKX instruments endpoint (public)")
		outPath = flag.String("out", "triangles_markets_okx.csv", "output csv path")
		timeout = flag.Duration("timeout", 20*time.Second, "http timeout")
		limitA  = flag.String("base", "", "optional: generate only for this base asset A (e.g. BTC)")
		state   = flag.String("state", "", "optional: allow only this state (e.g. live). empty = allow common states")
	)
	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	markets, err := fetchMarketsOKX(ctx, *apiURL, *state, *timeout)
	if err != nil {
		fatalf("fetch OKX instruments: %v", err)
	}

	rows := buildTrianglesOKX(markets, strings.TrimSpace(*limitA))

	if err := writeCSVAtomic(*outPath, rows); err != nil {
		fatalf("write csv: %v", err)
	}

	fmt.Printf("OK: markets=%d triangles=%d -> %s\n", len(markets), len(rows), *outPath)
}

func fetchMarketsOKX(ctx context.Context, url string, onlyState string, timeout time.Duration) ([]Market, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	cl := &http.Client{Timeout: timeout}
	resp, err := cl.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return nil, fmt.Errorf("http %d: %s", resp.StatusCode, strings.TrimSpace(string(b)))
	}

	var r okxInstrumentsResp
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}

	// OKX обычно: code="0" = success
	if strings.TrimSpace(r.Code) != "" && strings.TrimSpace(r.Code) != "0" {
		return nil, fmt.Errorf("okx code=%s msg=%s", strings.TrimSpace(r.Code), strings.TrimSpace(r.Msg))
	}

	allowState := func(s string) bool {
		s = strings.TrimSpace(s)
		if onlyState != "" {
			return strings.EqualFold(s, onlyState)
		}
		// чаще всего нужны только "live"
		return s == "" || strings.EqualFold(s, "live")
	}

	out := make([]Market, 0, len(r.Data))
	seen := make(map[string]struct{}, len(r.Data))

	for _, it := range r.Data {
		base := strings.TrimSpace(it.BaseCcy)
		quote := strings.TrimSpace(it.QuoteCcy)
		sym := strings.TrimSpace(it.InstID)

		if sym == "" || base == "" || quote == "" {
			continue
		}
		if !allowState(it.State) {
			continue
		}
		// instType=SPOT уже в query, но на всякий случай:
		if it.InstType != "" && !strings.EqualFold(it.InstType, "SPOT") {
			continue
		}

		// instId уникален
		if _, ok := seen[sym]; ok {
			continue
		}
		seen[sym] = struct{}{}

		out = append(out, Market{
			Symbol: sym,   // например: BTC-USDT
			Base:   base,  // BTC
			Quote:  quote, // USDT
		})
	}

	return out, nil
}

// Как у тебя: base1==base3==A, а quote1/quote3 — две котировки,
// между которыми должен существовать рынок (quote1/quote3) ИЛИ (quote3/quote1).
func buildTrianglesOKX(markets []Market, onlyBaseA string) []TriangleRow {
	// dir[base][quote] = Market  (ровно как на OKX: base-quote)
	dir := make(map[string]map[string]Market, 4096)
	quotesByBase := make(map[string][]string, 4096) // A -> list of quotes (A/quote)

	for _, m := range markets {
		if _, ok := dir[m.Base]; !ok {
			dir[m.Base] = make(map[string]Market, 64)
		}
		// если внезапно дубликаты base/quote, оставим первый (детерминируем ниже сортировкой входа не будем)
		if _, exists := dir[m.Base][m.Quote]; !exists {
			dir[m.Base][m.Quote] = m
			quotesByBase[m.Base] = append(quotesByBase[m.Base], m.Quote)
		}
	}

	// unique + sort quotes
	for base, qs := range quotesByBase {
		sort.Strings(qs)
		quotesByBase[base] = uniqueStrings(qs)
	}

	var bases []string
	for b := range quotesByBase {
		if onlyBaseA != "" && b != onlyBaseA {
			continue
		}
		bases = append(bases, b)
	}
	sort.Strings(bases)

	seen := make(map[string]struct{}, 4096)
	rows := make([]TriangleRow, 0, 16384)

	for _, A := range bases {
		qs := quotesByBase[A]
		if len(qs) < 2 {
			continue
		}

		for i := 0; i < len(qs); i++ {
			for j := i + 1; j < len(qs); j++ {
				qx, qy := qs[i], qs[j]

				// middle market: либо qx->qy, либо qy->qx
				var (
					mMid          Market
					quote1, quote3 string
					okMid         bool
				)

				if mm, ok := dir[qx][qy]; ok {
					mMid = mm
					quote1, quote3 = qx, qy
					okMid = true
				} else if mm, ok := dir[qy][qx]; ok {
					mMid = mm
					quote1, quote3 = qy, qx
					okMid = true
				}

				if !okMid {
					continue
				}

				// outer legs гарантированно есть, т.к. quote1/quote3 взяты из quotesByBase[A]
				m1, ok1 := dir[A][quote1]
				m3, ok3 := dir[A][quote3]
				if !(ok1 && ok3) {
					continue
				}

				// дедуп (на всякий) — фиксируем порядок quote1/quote3 ровно как выбран mid
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
					Symbol1: m1.Symbol,
					Symbol2: mMid.Symbol,
					Symbol3: m3.Symbol,
				})
			}
		}
	}

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

func writeCSVAtomic(path string, rows []TriangleRow) error {
	dir := filepath.Dir(path)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("mkdir: %w", err)
		}
	}

	tmp, err := os.CreateTemp(dir, filepath.Base(path)+".tmp.*")
	if err != nil {
		return fmt.Errorf("create tmp: %w", err)
	}
	tmpName := tmp.Name()

	cw := csv.NewWriter(tmp)
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


