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
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Tri struct {
	Base1, Quote1, Base2, Quote2, Base3, Quote3 string
	Sym1, Sym2, Sym3                            string
}

type Market struct {
	Symbol     string
	Base, Quote string
}

type Leg struct {
	From, To   string
	Symbol     string
	Action     string // BUY or SELL
	PriceSide  string // ASK for BUY, BID for SELL (под твой стакан)
}

// BUY: quote -> base по ASK
// SELL: base -> quote по BID
func makeLeg(from, to string, m Market) (Leg, bool) {
	if from == m.Quote && to == m.Base {
		return Leg{From: from, To: to, Symbol: m.Symbol, Action: "BUY", PriceSide: "ASK"}, true
	}
	if from == m.Base && to == m.Quote {
		return Leg{From: from, To: to, Symbol: m.Symbol, Action: "SELL", PriceSide: "BID"}, true
	}
	return Leg{}, false
}

// DFS по 3 ребрам (рынкам) без повторов ребер: USDT -> ? -> ? -> USDT
func findUSDTRoutes(markets []Market) [][]Leg {
	const start = "USDT"

	var routes [][]Leg

	used := make([]bool, len(markets))
	var path []Leg

	var dfs func(curr string, depth int)
	dfs = func(curr string, depth int) {
		if depth == 3 {
			if curr == start {
				cp := make([]Leg, len(path))
				copy(cp, path)
				routes = append(routes, cp)
			}
			return
		}

		for i, mk := range markets {
			if used[i] {
				continue
			}
			// этот рынок должен соединять curr с каким-то next
			// next может быть base или quote
			nexts := [2]string{mk.Base, mk.Quote}
			for _, next := range nexts {
				if next == curr {
					continue
				}
				// рынок должен содержать curr
				if mk.Base != curr && mk.Quote != curr {
					continue
				}
				leg, ok := makeLeg(curr, next, mk)
				if !ok {
					continue
				}
				used[i] = true
				path = append(path, leg)
				dfs(next, depth+1)
				path = path[:len(path)-1]
				used[i] = false
			}
		}
	}

	dfs(start, 0)
	return routes
}

func headerIndex(h []string) map[string]int {
	m := map[string]int{}
	for i, s := range h {
		m[strings.ToLower(strings.TrimSpace(s))] = i
	}
	return m
}

func needCol(idx map[string]int, name string) (int, error) {
	i, ok := idx[name]
	if !ok {
		return 0, fmt.Errorf("нет колонки %q в CSV", name)
	}
	return i, nil
}

func readTriangles(path string) ([]Tri, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.ReuseRecord = true

	h, err := r.Read()
	if err != nil {
		return nil, err
	}
	idx := headerIndex(h)

	get := func(row []string, col string) (string, error) {
		i, err := needCol(idx, col)
		if err != nil {
			return "", err
		}
		if i >= len(row) {
			return "", fmt.Errorf("строка короче заголовка: col=%s", col)
		}
		return strings.TrimSpace(row[i]), nil
	}

	var out []Tri
	for {
		row, err := r.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}

		b1, err := get(row, "base1"); if err != nil { return nil, err }
		q1, err := get(row, "quote1"); if err != nil { return nil, err }
		b2, err := get(row, "base2"); if err != nil { return nil, err }
		q2, err := get(row, "quote2"); if err != nil { return nil, err }
		b3, err := get(row, "base3"); if err != nil { return nil, err }
		q3, err := get(row, "quote3"); if err != nil { return nil, err }

		s1, err := get(row, "symbol1"); if err != nil { return nil, err }
		s2, err := get(row, "symbol2"); if err != nil { return nil, err }
		s3, err := get(row, "symbol3"); if err != nil { return nil, err }

		out = append(out, Tri{
			Base1: b1, Quote1: q1, Base2: b2, Quote2: q2, Base3: b3, Quote3: q3,
			Sym1: s1, Sym2: s2, Sym3: s3,
		})
	}
	return out, nil
}

func ensureDirForFile(path string) error {
	dir := filepath.Dir(path)
	if dir == "." || dir == "/" || dir == "" {
		return nil
	}
	return os.MkdirAll(dir, 0o755)
}

func writeRoutes(outPath string, rows [][]string) error {
	if err := ensureDirForFile(outPath); err != nil {
		return err
	}
	f, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	for _, row := range rows {
		if err := w.Write(row); err != nil {
			return err
		}
	}
	w.Flush()
	return w.Error()
}

func main() {
	in := "triangles_markets.csv"
	out := "triangles_usdt_routes.csv"

	tris, err := readTriangles(in)
	if err != nil {
		fmt.Println("ERR:", err)
		os.Exit(1)
	}

	// Заголовок итогового CSV
	rows := [][]string{{
		"start", "mid1", "mid2", "end",
		"leg1_symbol", "leg1_action", "leg1_price", "leg1_from", "leg1_to",
		"leg2_symbol", "leg2_action", "leg2_price", "leg2_from", "leg2_to",
		"leg3_symbol", "leg3_action", "leg3_price", "leg3_from", "leg3_to",
	}}

	countRoutes := 0

	for _, t := range tris {
		markets := []Market{
			{Symbol: t.Sym1, Base: t.Base1, Quote: t.Quote1},
			{Symbol: t.Sym2, Base: t.Base2, Quote: t.Quote2},
			{Symbol: t.Sym3, Base: t.Base3, Quote: t.Quote3},
		}

		routes := findUSDTRoutes(markets)
		if len(routes) == 0 {
			continue // треугольник не даёт цикл USDT->...->USDT (или странные активы)
		}

		for _, rt := range routes {
			// rt длина 3, начинается USDT и заканчивает USDT
			if len(rt) != 3 {
				continue
			}
			row := []string{
				"USDT", rt[0].To, rt[1].To, "USDT",

				rt[0].Symbol, rt[0].Action, rt[0].PriceSide, rt[0].From, rt[0].To,
				rt[1].Symbol, rt[1].Action, rt[1].PriceSide, rt[1].From, rt[1].To,
				rt[2].Symbol, rt[2].Action, rt[2].PriceSide, rt[2].From, rt[2].To,
			}
			rows = append(rows, row)
			countRoutes++
		}
	}

	if err := writeRoutes(out, rows); err != nil {
		fmt.Println("ERR: write csv:", err)
		os.Exit(1)
	}

	fmt.Printf("OK: routes=%d -> %s\n", countRoutes, out)
}

