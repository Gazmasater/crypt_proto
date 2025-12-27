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
	"fmt"
	"log"
	"os"
)

type Graph map[string][]string

func main() {
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
	}

	triangles := findTriangles(graph)

	for _, t := range triangles {
		fmt.Printf("%s -> %s -> %s -> %s\n", t[0], t[1], t[2], t[0])
	}
}

func appendUnique(slice []string, val string) []string {
	for _, s := range slice {
		if s == val {
			return slice
		}
	}
	return append(slice, val)
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
					// Добавляем оба обхода: a-b-c-a и a-c-b-a
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



