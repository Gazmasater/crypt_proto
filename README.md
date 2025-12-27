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
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Symbol struct {
	BaseAsset  string `json:"baseAsset"`
	QuoteAsset string `json:"quoteAsset"`
	Status     string `json:"status"`
}

type ExchangeInfo struct {
	Symbols []Symbol `json:"symbols"`
}

type Graph map[string][]string

func main() {
	// --- Шаг 1: Получаем реальные пары с MEXC ---
	pairs := fetchMEXCPairs()

	// --- Шаг 2: Загружаем CSV и строим граф ---
	graph := loadGraph("triangles.csv")

	// --- Шаг 3: Находим все треугольники ---
	triangles := findTriangles(graph)

	// --- Шаг 4: Фильтруем реальные ---
	realTriangles := [][]string{}
	for _, t := range triangles {
		if isRealTriangle(t, pairs) {
			realTriangles = append(realTriangles, []string{t[0], t[1], t[2], t[0]})
		}
	}

	// --- Шаг 5: Записываем в CSV ---
	writeTriangles("real_triangles.csv", realTriangles)
	fmt.Printf("Найдено %d реальных треугольников\n", len(realTriangles))
}

// --- Функции ---

func fetchMEXCPairs() map[string]map[string]bool {
	resp, err := http.Get("https://api.mexc.com/api/v3/exchangeInfo")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var info ExchangeInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		log.Fatal(err)
	}

	pairs := make(map[string]map[string]bool)
	for _, s := range info.Symbols {
		if s.Status != "ENABLED" {
			continue
		}
		if pairs[s.BaseAsset] == nil {
			pairs[s.BaseAsset] = make(map[string]bool)
		}
		pairs[s.BaseAsset][s.QuoteAsset] = true
	}
	return pairs
}

func loadGraph(filename string) Graph {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	graph := make(Graph)
	for i, row := range records {
		if i == 0 || len(row) < 6 {
			continue
		}
		base1, quote1 := row[0], row[1]
		base2, quote2 := row[2], row[3]
		base3, quote3 := row[4], row[5]

		graph[base1] = appendUnique(graph[base1], quote1)
		graph[quote1] = appendUnique(graph[quote1], base1)
		graph[base2] = appendUnique(graph[base2], quote2)
		graph[quote2] = appendUnique(graph[quote2], base2)
		graph[base3] = appendUnique(graph[base3], quote3)
		graph[quote3] = appendUnique(graph[quote3], base3)
	}
	return graph
}

func appendUnique(slice []string, val string) []string {
	for _, s := range slice {
		if s == val {
			return slice
		}
	}
	return append(slice, val)
}

func findTriangles(graph Graph) [][3]string {
	var triangles [][3]string
	for a := range graph {
		for _, b := range graph[a] {
			if b == a {
				continue
			}
			for _, c := range graph[b] {
				if c == a || c == b {
					continue
				}
				if contains(graph[c], a) {
					triangles = append(triangles, [3]string{a, b, c})
					triangles = append(triangles, [3]string{a, c, b})
				}
			}
		}
	}
	return triangles
}

func contains(slice []string, val string) bool {
	for _, s := range slice {
		if s == val {
			return true
		}
	}
	return false
}

func isRealTriangle(tri [3]string, pairs map[string]map[string]bool) bool {
	a, b, c := tri[0], tri[1], tri[2]
	return (pairs[a][b] || pairs[b][a]) &&
		(pairs[b][c] || pairs[c][b]) &&
		(pairs[c][a] || pairs[a][c])
}

func writeTriangles(filename string, triangles [][]string) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, t := range triangles {
		if err := writer.Write(t); err != nil {
			log.Fatal(err)
		}
	}
}



