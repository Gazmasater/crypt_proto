package main

import (
	"encoding/csv"
	"log"
	"os"
)

type Graph map[string][]string

func main() {
	file, err := os.Open("triangles_markets.csv")
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

	// Создаём файл для записи результата
	outFile, err := os.Create("triangles_output.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

	writer := csv.NewWriter(outFile)
	defer writer.Flush()

	// Заголовок
	writer.Write([]string{"start", "mid1", "mid2", "end"})

	for _, t := range triangles {
		writer.Write([]string{t[0], t[1], t[2], t[0]})
	}
}

// appendUnique добавляет элемент в срез, если его там нет
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

// contains проверяет, есть ли элемент в срезе
func contains(slice []string, val string) bool {
	for _, s := range slice {
		if s == val {
			return true
		}
	}
	return false
}
