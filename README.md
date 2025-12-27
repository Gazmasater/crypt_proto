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




package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const mexcExchangeInfoURL = "https://api.mexc.com/api/v3/exchangeInfo"

type ExchangeInfo struct {
	Symbols []Symbol `json:"symbols"`
}

type Symbol struct {
	Symbol               string   `json:"symbol"`
	BaseAsset            string   `json:"baseAsset"`
	QuoteAsset           string   `json:"quoteAsset"`
	IsSpotTradingAllowed bool     `json:"isSpotTradingAllowed"`
	OrderTypes           []string `json:"orderTypes"`
}

type Edge struct {
	From   string
	To     string
	Symbol string
	Side   string // BUY or SELL
}

type Graph map[string][]Edge

// -------------------- MAIN --------------------

func main() {
	info := loadExchangeInfo()

	graph := buildDirectedGraph(info)

	triangles := findTriangles(graph)

	saveCSV("triangles.csv", triangles)

	fmt.Println("✔ triangles saved:", len(triangles))
}

// -------------------- LOAD EXCHANGE INFO --------------------

func loadExchangeInfo() ExchangeInfo {
	resp, err := http.Get(mexcExchangeInfoURL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var info ExchangeInfo
	if err := json.Unmarshal(body, &info); err != nil {
		panic(err)
	}

	return info
}

// -------------------- GRAPH --------------------

func buildDirectedGraph(info ExchangeInfo) Graph {
	graph := make(Graph)

	for _, s := range info.Symbols {

		if !s.IsSpotTradingAllowed {
			continue
		}

		if !hasMarket(s.OrderTypes) {
			continue
		}

		base := s.BaseAsset
		quote := s.QuoteAsset

		// SELL: base -> quote
		graph[base] = append(graph[base], Edge{
			From:   base,
			To:     quote,
			Symbol: s.Symbol,
			Side:   "SELL",
		})

		// BUY: quote -> base
		graph[quote] = append(graph[quote], Edge{
			From:   quote,
			To:     base,
			Symbol: s.Symbol,
			Side:   "BUY",
		})
	}

	return graph
}

func hasMarket(types []string) bool {
	for _, t := range types {
		if t == "MARKET" {
			return true
		}
	}
	return false
}

// -------------------- TRIANGLES --------------------

type Triangle struct {
	A, B, C string
	Path    string
}

func findTriangles(graph Graph) []Triangle {
	var result []Triangle
	seen := map[string]bool{}

	for a, edgesAB := range graph {
		for _, ab := range edgesAB {
			b := ab.To

			for _, bc := range graph[b] {
				c := bc.To
				if c == a {
					continue
				}

				for _, ca := range graph[c] {
					if ca.To != a {
						continue
					}

					// фильтр мусора: USDT <-> USDC
					if isStableLoop(a, b, c) {
						continue
					}

					key := canonicalKey(a, b, c)
					if seen[key] {
						continue
					}
					seen[key] = true

					result = append(result, Triangle{
						A:    a,
						B:    b,
						C:    c,
						Path: fmt.Sprintf("%s -> %s -> %s -> %s", a, b, c, a),
					})
				}
			}
		}
	}

	return result
}

// -------------------- FILTERS --------------------

func isStableLoop(a, b, c string) bool {
	stable := func(x string) bool {
		return x == "USDT" || x == "USDC"
	}

	// убираем циклы типа USDT -> USDC -> X -> USDT
	return stable(a) && stable(b) ||
		stable(b) && stable(c) ||
		stable(a) && stable(c)
}

func canonicalKey(a, b, c string) string {
	if a < b && a < c {
		return a + "|" + b + "|" + c
	}
	if b < a && b < c {
		return b + "|" + c + "|" + a
	}
	return c + "|" + a + "|" + b
}

// -------------------- CSV --------------------

func saveCSV(path string, items []Triangle) {
	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	fmt.Fprintln(f, "A,B,C,PATH")

	for _, t := range items {
		fmt.Fprintf(
			f,
			"%s,%s,%s,%s\n",
			t.A, t.B, t.C, t.Path,
		)
	}
}

