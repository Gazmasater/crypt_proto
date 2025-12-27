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
	"log"
	"os"
)

type Graph map[string][]string

func main() {
	// Читаем CSV с парами
	file, err := os.Open("triangles.csv")
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
	availablePairs := make(map[string]map[string]bool)

	// Пропускаем заголовок
	for i, row := range records {
		if i == 0 {
			continue
		}
		if len(row) < 6 {
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

		// Создаем доступные пары для проверки реальности треугольника
		addPair(availablePairs, base1, quote1)
		addPair(availablePairs, base2, quote2)
		addPair(availablePairs, base3, quote3)
	}

	// Ищем все возможные треугольники
	triangles := findTriangles(graph)

	// Фильтруем только реальные треугольники
	realTriangles := [][]string{}
	for _, tri := range triangles {
		if isRealTriangle(tri, availablePairs) {
			realTriangles = append(realTriangles, []string{tri[0], tri[1], tri[2]})
		}
	}

	// Сохраняем в CSV
	outFile, err := os.Create("real_triangles.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

	writer := csv.NewWriter(outFile)
	defer writer.Flush()

	for _, tri := range realTriangles {
		if err := writer.Write([]string{tri[0], tri[1], tri[2], tri[0]}); err != nil {
			log.Fatal(err)
		}
	}
}

// appendUnique добавляет значение в слайс, если его там ещё нет
func appendUnique(slice []string, val string) []string {
	for _, s := range slice {
		if s == val {
			return slice
		}
	}
	return append(slice, val)
}

// addPair добавляет пару в карту доступных пар
func addPair(m map[string]map[string]bool, a, b string) {
	if m[a] == nil {
		m[a] = make(map[string]bool)
	}
	m[a][b] = true
}

// findTriangles ищет все треугольники с обеими последовательностями обхода
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

// isRealTriangle проверяет, что все три пары существуют
func isRealTriangle(tri [3]string, pairs map[string]map[string]bool) bool {
	a, b, c := tri[0], tri[1], tri[2]
	return (pairs[a][b] || pairs[b][a]) &&
		(pairs[b][c] || pairs[c][b]) &&
		(pairs[c][a] || pairs[a][c])
}

func contains(slice []string, val string) bool {
	for _, s := range slice {
		if s == val {
			return true
		}
	}
	return false
}


