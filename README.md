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
	"strings"
	"time"
)

type KucoinSymbol struct {
	Symbol       string `json:"symbol"`
	BaseCurrency string `json:"baseCurrency"`
	QuoteCurrency string `json:"quoteCurrency"`
	EnableTrading bool   `json:"enableTrading"`
}

func fetchKucoinSymbols() ([]KucoinSymbol, error) {
	url := "https://api.kucoin.com/api/v2/symbols"
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data struct {
		Code string          `json:"code"`
		Data []KucoinSymbol `json:"data"`
	}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return data.Data, nil
}

var forbiddenUSD = map[string]bool{
	"USDC": true,
	"USD1": true,
	// добавь сюда любые другие, начинающиеся на USD, которые нужно исключить
}

func isBadUsd(x string) bool {
	// USDT — это нормально
	if x == "USDT" {
		return false
	}
	// запрет для любых, начинающихся с USD
	return strings.HasPrefix(x, "USD")
}

type Triangle struct {
	A, B, C string
	Leg1    string
	Leg2    string
	Leg3    string
}

func findTriangles(symbols []KucoinSymbol) []Triangle {
	graph := make(map[string]map[string]string)

	for _, s := range symbols {
		if !s.EnableTrading {
			continue
		}
		base := s.BaseCurrency
		quote := s.QuoteCurrency

		if graph[base] == nil {
			graph[base] = make(map[string]string)
		}
		graph[base][quote] = s.Symbol

		if graph[quote] == nil {
			graph[quote] = make(map[string]string)
		}
		graph[quote][base] = s.Symbol + "_INV"
	}

	var result []Triangle
	for A, toB := range graph {
		if A != "USDT" {
			continue
		}
		for B, leg1 := range toB {
			if isBadUsd(B) {
				continue
			}
			if graph[B] == nil {
				continue
			}
			for C, leg2 := range graph[B] {
				if isBadUsd(C) {
					continue
				}
				if graph[C] == nil {
					continue
				}
				leg3, ok := graph[C][A]
				if !ok {
					continue
				}
				if isBadUsd(B) || isBadUsd(C) {
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

func main() {
	symbols, err := fetchKucoinSymbols()
	if err != nil {
		log.Fatalf("Ошибка получения пар KuCoin: %v", err)
	}

	triangles := findTriangles(symbols)
	log.Printf("Найдено треугольников на KuCoin (USDT + без USD*): %d\n", len(triangles))

	file, err := os.Create("triangles_kucoin.csv")
	if err != nil {
		log.Fatalf("Ошибка создания файла: %v", err)
	}
	defer file.Close()

	w := csv.NewWriter(file)
	defer w.Flush()

	w.Write([]string{"A", "B", "C", "leg1", "leg2", "leg3"})
	for _, t := range triangles {
		w.Write([]string{t.A, t.B, t.C, t.Leg1, t.Leg2, t.Leg3})
	}

	log.Println("triangles_kucoin.csv успешно создан")
}



