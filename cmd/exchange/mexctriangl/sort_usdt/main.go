package main

import (
	"encoding/csv"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

type Symbol struct {
	Symbol      string   `json:"symbol"`
	BaseAsset   string   `json:"baseAsset"`
	QuoteAsset  string   `json:"quoteAsset"`
	OrderTypes  []string `json:"orderTypes"`
	Permissions []string `json:"permissions"`
}

type ExchangeInfo struct {
	Symbols []Symbol `json:"symbols"`
}

type Triangle struct {
	A, B, C string
	Leg1    string
	Leg2    string
	Leg3    string
}

// ----------------- helpers -----------------

var forbiddenStable = map[string]bool{
	"USDC": true,
	"USD1": true,
}

func has(list []string, v string) bool {
	for _, x := range list {
		if x == v {
			return true
		}
	}
	return false
}

func isTradable(s Symbol) bool {
	return has(s.Permissions, "SPOT") && has(s.OrderTypes, "MARKET")
}

// проверка треугольника: старт USDT, запрещённые стейблы
func validTriangleUSDT(A, B, C string) bool {
	if A != "USDT" {
		return false
	}
	if forbiddenStable[B] || forbiddenStable[C] {
		return false
	}
	return true
}

// ----------------- main logic -----------------

func findTriangles(symbols []Symbol) []Triangle {
	graph := make(map[string]map[string]string)

	for _, s := range symbols {
		if !isTradable(s) {
			continue
		}
		if _, ok := graph[s.BaseAsset]; !ok {
			graph[s.BaseAsset] = make(map[string]string)
		}
		graph[s.BaseAsset][s.QuoteAsset] = s.Symbol

		// инвертированное направление
		if _, ok := graph[s.QuoteAsset]; !ok {
			graph[s.QuoteAsset] = make(map[string]string)
		}
		graph[s.QuoteAsset][s.BaseAsset] = s.Symbol + "_INV"
	}

	var result []Triangle

	for A, toB := range graph {
		if A != "USDT" {
			continue
		}

		for B, leg1 := range toB {
			if graph[B] == nil {
				continue
			}
			for C, leg2 := range graph[B] {
				if graph[C] == nil {
					continue
				}
				leg3, ok := graph[C][A]
				if !ok {
					continue
				}
				if !validTriangleUSDT(A, B, C) {
					continue
				}
				result = append(result, Triangle{
					A:    A,
					B:    B,
					C:    C,
					Leg1: leg1,
					Leg2: leg2,
					Leg3: leg3,
				})
			}
		}
	}
	return result
}

// ----------------- HTTP -----------------

func fetchExchangeInfo() ([]Symbol, error) {
	url := "https://api.mexc.com/api/v3/exchangeInfo"
	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var info ExchangeInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, err
	}

	return info.Symbols, nil
}

// ----------------- main -----------------

func main() {
	symbols, err := fetchExchangeInfo()
	if err != nil {
		log.Fatal(err)
	}

	triangles := findTriangles(symbols)
	log.Printf("Найдено треугольников: %d\n", len(triangles))

	file, err := os.Create("triangles_usdt.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	w := csv.NewWriter(file)
	defer w.Flush()

	// CSV заголовок
	w.Write([]string{"A", "B", "C", "leg1", "leg2", "leg3"})

	for _, t := range triangles {
		w.Write([]string{
			t.A,
			t.B,
			t.C,
			t.Leg1,
			t.Leg2,
			t.Leg3,
		})
	}

	log.Println("triangles_usdt.csv успешно создан")
}
