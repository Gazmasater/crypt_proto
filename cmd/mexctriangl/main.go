package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type Symbol struct {
	Symbol               string   `json:"symbol"`
	BaseAsset            string   `json:"baseAsset"`
	QuoteAsset           string   `json:"quoteAsset"`
	IsSpotTradingAllowed bool     `json:"isSpotTradingAllowed"`
	OrderTypes           []string `json:"orderTypes"`
	Permissions          []string `json:"permissions"`
}

type ExchangeInfo struct {
	Symbols []Symbol `json:"symbols"`
}

type Triangle struct {
	A, B, C string
	Pairs   [3]string
}

func isTradable(s Symbol) bool {
	if !s.IsSpotTradingAllowed {
		return false
	}
	hasMarket := false
	for _, t := range s.OrderTypes {
		if t == "MARKET" {
			hasMarket = true
			break
		}
	}
	if !hasMarket {
		return false
	}
	for _, p := range s.Permissions {
		if p == "SPOT" {
			return true
		}
	}
	return false
}

func findTriangles(symbols []Symbol) []Triangle {
	pairMap := make(map[string]map[string]Symbol)
	for _, s := range symbols {
		if !isTradable(s) {
			continue
		}
		if _, ok := pairMap[s.BaseAsset]; !ok {
			pairMap[s.BaseAsset] = make(map[string]Symbol)
		}
		pairMap[s.BaseAsset][s.QuoteAsset] = s

		// добавляем инверсную пару
		if _, ok := pairMap[s.QuoteAsset]; !ok {
			pairMap[s.QuoteAsset] = make(map[string]Symbol)
		}
		pairMap[s.QuoteAsset][s.BaseAsset] = Symbol{
			Symbol:               s.Symbol + "_INV",
			BaseAsset:            s.QuoteAsset,
			QuoteAsset:           s.BaseAsset,
			IsSpotTradingAllowed: true,
			OrderTypes:           s.OrderTypes,
			Permissions:          s.Permissions,
		}
	}

	var triangles []Triangle
	for _, s1 := range symbols {
		if !isTradable(s1) {
			continue
		}
		A := s1.BaseAsset
		B := s1.QuoteAsset

		for C, s2 := range pairMap[B] {
			if !isTradable(s2) {
				continue
			}
			if s3map, ok := pairMap[C]; ok {
				if s3, ok2 := s3map[A]; ok2 && isTradable(s3) {
					// A->B->C->A
					triangles = append(triangles, Triangle{
						A: A, B: B, C: C,
						Pairs: [3]string{s1.Symbol, s2.Symbol, s3.Symbol},
					})
					// A->C->B->A
					triangles = append(triangles, Triangle{
						A: A, B: C, C: B,
						Pairs: [3]string{s2.Symbol, s1.Symbol, s3.Symbol},
					})
				}
			}
		}
	}
	return triangles
}

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

func main() {
	symbols, err := fetchExchangeInfo()
	if err != nil {
		log.Fatal("Ошибка получения данных с MEXC:", err)
	}

	fmt.Println("!!!!!", symbols)

	triangles := findTriangles(symbols)
	log.Printf("Найдено %d треугольников", len(triangles))

	// Создаём CSV файл
	file, err := os.Create("triangles.csv")
	if err != nil {
		log.Fatal("Ошибка создания файла:", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Записываем заголовок
	writer.Write([]string{"A", "B", "C", "PAIR1", "PAIR2", "PAIR3"})

	// Записываем все треугольники
	for _, t := range triangles {
		writer.Write([]string{t.A, t.B, t.C, t.Pairs[0], t.Pairs[1], t.Pairs[2]})
	}

	log.Println("Все треугольники записаны в triangles.csv")
}
