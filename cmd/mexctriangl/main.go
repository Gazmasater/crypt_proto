package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
)

// ================= MEXC STRUCTS =================

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

// ================= GRAPH =================

type TradeGraph map[string]map[string]bool

func addEdge(g TradeGraph, from, to string) {
	if g[from] == nil {
		g[from] = map[string]bool{}
	}
	g[from][to] = true
}

// ================= UTILS =================

func hasMarket(orderTypes []string) bool {
	for _, t := range orderTypes {
		if t == "MARKET" {
			return true
		}
	}
	return false
}

func triangleKey(a, b, c string) string {
	x := []string{a, b, c}
	sort.Strings(x)
	return x[0] + "|" + x[1] + "|" + x[2]
}

// ================= MAIN =================

func main() {
	fmt.Println("Loading MEXC exchangeInfo...")

	resp, err := http.Get("https://api.mexc.com/api/v3/exchangeInfo")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var info ExchangeInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		log.Fatal(err)
	}

	graph := TradeGraph{}

	// ---------- BUILD DIRECTED GRAPH ----------
	for _, s := range info.Symbols {
		if !s.IsSpotTradingAllowed {
			continue
		}
		if !hasMarket(s.OrderTypes) {
			continue
		}

		// SELL
		addEdge(graph, s.BaseAsset, s.QuoteAsset)
		// BUY
		addEdge(graph, s.QuoteAsset, s.BaseAsset)
	}

	fmt.Printf("Assets in graph: %d\n", len(graph))

	// ---------- FIND TRIANGLES ----------
	seen := map[string]bool{}
	var triangles [][3]string

	for a := range graph {
		for b := range graph[a] {
			for c := range graph[b] {

				if a == b || b == c || a == c {
					continue
				}

				if graph[c][a] {
					key := triangleKey(a, b, c)
					if !seen[key] {
						seen[key] = true
						triangles = append(triangles, [3]string{a, b, c})
					}
				}
			}
		}
	}

	fmt.Printf("Found %d real arbitrage triangles\n", len(triangles))

	// ---------- WRITE CSV ----------
	file, err := os.Create("triangles_real.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"A", "B", "C", "PATH"})

	for _, t := range triangles {
		writer.Write([]string{
			t[0],
			t[1],
			t[2],
			fmt.Sprintf("%s -> %s -> %s -> %s", t[0], t[1], t[2], t[0]),
		})
	}

	fmt.Println("Saved to triangles_real.csv")
}
