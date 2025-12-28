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

type Leg struct {
	Symbol string
}

type Triangle struct {
	A, B, C string
	Leg1    string // A -> B
	Leg2    string // B -> C
	Leg3    string // C -> A
}

// ---------- helpers ----------

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

// ---------- main logic ----------

func findTriangles(symbols []Symbol) []Triangle {

	// map[from][to] = symbolName (or *_INV)
	graph := make(map[string]map[string]string)

	for _, s := range symbols {
		if !isTradable(s) {
			continue
		}

		// прямое направление
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

	// A → B → C → A
	for A, toB := range graph {
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

				result = append(result, Triangle{
					A:    A,
					B:    B,
					C:    C,
					Leg1: leg1, // A → B
					Leg2: leg2, // B → C
					Leg3: leg3, // C → A
				})
			}
		}
	}

	return result
}

// ---------- HTTP ----------

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

// ---------- main ----------

func main() {
	symbols, err := fetchExchangeInfo()
	if err != nil {
		log.Fatal(err)
	}

	triangles := findTriangles(symbols)
	log.Printf("Найдено треугольников: %d\n", len(triangles))

	file, err := os.Create("triangles.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	w := csv.NewWriter(file)
	defer w.Flush()

	// строго фиксированный формат
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

	log.Println("triangles.csv успешно создан")
}

