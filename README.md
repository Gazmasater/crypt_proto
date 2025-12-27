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
	"log"
	"net/http"
	"time"
)

// ------------------- Структуры -------------------
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

// ------------------- Фильтр -------------------
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

// ------------------- Генерация треугольников с инверсными парами -------------------
func findTriangles(symbols []Symbol) []Triangle {
	// map[валюта][валюта] = symbol
	pairMap := make(map[string]map[string]Symbol)
	for _, s := range symbols {
		if !isTradable(s) {
			continue
		}
		if _, ok := pairMap[s.BaseAsset]; !ok {
			pairMap[s.BaseAsset] = make(map[string]Symbol)
		}
		pairMap[s.BaseAsset][s.QuoteAsset] = s
		// добавляем инверсную пару для поиска
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
					// A -> B -> C -> A
					triangles = append(triangles, Triangle{
						A: A, B: B, C: C,
						Pairs: [3]string{s1.Symbol, s2.Symbol, s3.Symbol},
					})
					// A -> C -> B -> A
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

// ------------------- Получение данных с API -------------------
func fetchExchangeInfo() ([]Symbol, error) {
	url := "https://www.mexc.com/open/api/v2/market/symbols"
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

// ------------------- Основная функция -------------------
func main() {
	symbols, err := fetchExchangeInfo()
	if err != nil {
		log.Fatal("Ошибка получения данных с MEXC:", err)
	}

	triangles := findTriangles(symbols)
	fmt.Printf("Найдено %d треугольников с обеими направленностями и инверсными парами:\n", len(triangles))
	for _, t := range triangles {
		fmt.Println(t.A, "->", t.B, "->", t.C, "Pairs:", t.Pairs)
	}
}



