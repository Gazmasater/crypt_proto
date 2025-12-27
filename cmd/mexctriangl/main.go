package main

import (
	"encoding/csv"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
)

type Symbol struct {
	Symbol               string   `json:"symbol"`
	BaseAsset            string   `json:"baseAsset"`
	QuoteAsset           string   `json:"quoteAsset"`
	Status               string   `json:"status"`
	IsSpotTradingAllowed bool     `json:"isSpotTradingAllowed"`
	OrderTypes           []string `json:"orderTypes"`
}

type ExchangeInfo struct {
	Symbols []Symbol `json:"symbols"`
}

// canTrade[from][to] = true
type TradeGraph map[string]map[string]bool

func main() {
	canTrade := loadMEXCDirectedGraph()

	graph := loadGraphFromCSV("triangles_markets.csv")

	results := [][]string{}

	for a := range graph {
		for _, x := range graph[a] {
			for _, y := range graph[a] {

				if x == y || x == a || y == a {
					continue
				}

				// a → x → y → a
				if canTrade[a][x] && canTrade[x][y] && canTrade[y][a] {
					results = append(results, []string{a, x, y, a})
				}
			}
		}
	}

	writeCSV("real_triangles.csv", results)
	log.Printf("✅ найдено %d реальных направленных треугольников\n", len(results))
}

//////////////////////////////////////////////////////
// ЗАГРУЗКА MEXC → НАПРАВЛЕННЫЙ ГРАФ
//////////////////////////////////////////////////////

func loadMEXCDirectedGraph() TradeGraph {
	resp, err := http.Get("https://api.mexc.com/api/v3/exchangeInfo")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var info ExchangeInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		log.Fatal(err)
	}

	canTrade := make(TradeGraph)

	for _, s := range info.Symbols {

		if s.Status != "ENABLED" {
			continue
		}
		if !s.IsSpotTradingAllowed {
			continue
		}
		if !hasMarketOrder(s.OrderTypes) {
			continue
		}

		base := strings.ToUpper(s.BaseAsset)
		quote := strings.ToUpper(s.QuoteAsset)

		// sell: base → quote
		if canTrade[base] == nil {
			canTrade[base] = map[string]bool{}
		}
		canTrade[base][quote] = true

		// buy: quote → base
		if canTrade[quote] == nil {
			canTrade[quote] = map[string]bool{}
		}
		canTrade[quote][base] = true
	}

	return canTrade
}

func hasMarketOrder(types []string) bool {
	for _, t := range types {
		if t == "MARKET" {
			return true
		}
	}
	return false
}

//////////////////////////////////////////////////////
// CSV → граф валют (без направлений)
//////////////////////////////////////////////////////

func loadGraphFromCSV(filename string) map[string][]string {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	graph := map[string][]string{}

	for i, row := range rows {
		if i == 0 || len(row) < 6 {
			continue
		}

		pairs := [][2]string{
			{row[0], row[1]},
			{row[2], row[3]},
			{row[4], row[5]},
		}

		for _, p := range pairs {
			a := strings.ToUpper(p[0])
			b := strings.ToUpper(p[1])

			graph[a] = appendUnique(graph[a], b)
			graph[b] = appendUnique(graph[b], a)
		}
	}

	return graph
}

//////////////////////////////////////////////////////

func appendUnique(arr []string, v string) []string {
	for _, x := range arr {
		if x == v {
			return arr
		}
	}
	return append(arr, v)
}

//////////////////////////////////////////////////////

func writeCSV(name string, rows [][]string) {
	f, err := os.Create(name)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	for _, r := range rows {
		_ = w.Write(r)
	}
}
