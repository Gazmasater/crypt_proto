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
	"time"
)

type Symbol struct {
	Symbol               string   `json:"symbol"`
	BaseAsset            string   `json:"baseAsset"`
	QuoteAsset           string   `json:"quoteAsset"`
	IsSpotTradingAllowed bool     `json:"isSpotTradingAllowed"`
	OrderTypes           []string `json:"orderTypes"`
	Permissions          []string `json:"permissions"`
	BaseSizePrecision    string   `json:"baseSizePrecision"`
	QuoteAmountPrecision string   `json:"quoteAmountPrecisionMarket"`
	TakerCommission      string   `json:"takerCommission"`
}

type ExchangeInfo struct {
	Symbols []Symbol `json:"symbols"`
}

type Triangle struct {
	A, B, C string
	Legs    [3]string
}

type Leg struct {
	Pair      string
	Side      string
	BaseStep  string
	QuoteStep string
	TakerFee  string
	Invert    bool
}

// Проверка торгуемости пары
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

// Получение exchangeInfo
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

// Генерация всех треугольников
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

		// Добавляем инверсную пару
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
			BaseSizePrecision:    s.BaseSizePrecision,
			QuoteAmountPrecision: s.QuoteAmountPrecision,
			TakerCommission:      s.TakerCommission,
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
					triangles = append(triangles, Triangle{
						A: A, B: B, C: C,
						Legs: [3]string{s1.Symbol, s2.Symbol, s3.Symbol},
					})
				}
			}
		}
	}
	return triangles
}

// Создаём leg с BUY/SELL и invert
func makeLeg(s Symbol, name string) Leg {
	invert := false
	side := "SELL"
	if len(name) >= 4 && name[len(name)-4:] == "_INV" {
		invert = true
		side = "BUY"
	}
	return Leg{
		Pair:      name,
		Side:      side,
		BaseStep:  s.BaseSizePrecision,
		QuoteStep: s.QuoteAmountPrecision,
		TakerFee:  s.TakerCommission,
		Invert:    invert,
	}
}

// Строим legs треугольника с возвратом в стартовую валюту
func buildTriangleExec(tr Triangle, start string, pairMap map[string]map[string]Symbol) []Leg {
	var legs []Leg

	// если стартовая валюта есть в треугольнике
	if tr.A == start || tr.B == start || tr.C == start {
		// стандартный цикл 3 leg
		order := []string{tr.Legs[0], tr.Legs[1], tr.Legs[2]}
		for _, n := range order {
			for _, m := range pairMap {
				if sym, ok := m[n]; ok {
					legs = append(legs, makeLeg(sym, n))
					break
				}
			}
		}
		return legs
	}

	// иначе добавляем вход start->трёхугольник
	found := false
	var entry Symbol
	var entryName string
	for _, vertex := range []string{tr.A, tr.B, tr.C} {
		if sym, ok := pairMap[start][vertex]; ok {
			entry = sym
			entryName = sym.Symbol
			found = true
			break
		}
		if sym, ok := pairMap[vertex][start]; ok {
			entry = sym
			entryName = sym.Symbol + "_INV"
			found = true
			break
		}
	}
	if !found {
		return nil
	}
	legs = append(legs, makeLeg(entry, entryName))

	// добавляем стандартные 3 leg треугольника
	for _, n := range tr.Legs {
		for _, m := range pairMap {
			if sym, ok := m[n]; ok {
				legs = append(legs, makeLeg(sym, n))
				break
			}
		}
	}

	// добавляем выход в стартовую валюту
	last := legs[len(legs)-1]
	// если последняя валюта != стартовая, ищем пару для выхода
	lastBase := last.Pair
	for _, vertex := range []string{tr.A, tr.B, tr.C} {
		if sym, ok := pairMap[vertex][start]; ok {
			legs = append(legs, makeLeg(sym, sym.Symbol))
			break
		}
		if sym, ok := pairMap[start][vertex]; ok {
			legs = append(legs, makeLeg(sym, sym.Symbol))
			break
		}
	}

	return legs
}

func main() {
	symbols, err := fetchExchangeInfo()
	if err != nil {
		log.Fatal(err)
	}

	triangles := findTriangles(symbols)
	log.Printf("Найдено %d треугольников", len(triangles))

	// создаём map для быстрого поиска
	pairMap := make(map[string]map[string]Symbol)
	for _, s := range symbols {
		if _, ok := pairMap[s.BaseAsset]; !ok {
			pairMap[s.BaseAsset] = make(map[string]Symbol)
		}
		pairMap[s.BaseAsset][s.QuoteAsset] = s
	}

	startAsset := "USDT"

	// CSV
	file, err := os.Create("triangles_exec.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"StartAsset", "LegPair", "Side", "BaseStep", "QuoteStep", "TakerFee", "Invert"})

	for _, tr := range triangles {
		legs := buildTriangleExec(tr, startAsset, pairMap)
		if len(legs) == 0 {
			continue
		}
		for _, l := range legs {
			writer.Write([]string{
				startAsset,
				l.Pair,
				l.Side,
				l.BaseStep,
				l.QuoteStep,
				l.TakerFee,
				fmt.Sprintf("%v", l.Invert),
			})
		}
	}

	log.Println("Треугольники с возвратом в стартовую валюту записаны в triangles_exec.csv")
}

[{
	"resource": "/home/gaz358/myprog/crypt_proto/cmd/mexctriangl/main.go",
	"owner": "_generated_diagnostic_collection_name_#1",
	"code": {
		"value": "UnusedVar",
		"target": {
			"$mid": 1,
			"path": "/golang.org/x/tools/internal/typesinternal",
			"scheme": "https",
			"authority": "pkg.go.dev",
			"fragment": "UnusedVar"
		}
	},
	"severity": 8,
	"message": "declared and not used: lastBase",
	"source": "compiler",
	"startLineNumber": 210,
	"startColumn": 2,
	"endLineNumber": 210,
	"endColumn": 10,
	"tags": [
		1
	],
	"origin": "extHost1"
}]

