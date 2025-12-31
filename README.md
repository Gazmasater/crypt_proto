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
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type kucoinExchangeInfo struct {
	Code string `json:"code"`
	Data []struct {
		Symbol        string `json:"symbol"`
		BaseCurrency  string `json:"baseCurrency"`
		QuoteCurrency string `json:"quoteCurrency"`
		EnableTrading bool   `json:"enableTrading"`
	} `json:"data"`
}

type pairMarket struct {
	Symbol string
	Base   string
	Quote  string
}

type Triangle struct {
	A    string
	B    string
	C    string
	Leg1 string
	Leg2 string
	Leg3 string
}

// список стейблкоинов
var stableCoins = map[string]bool{
	"USDT":  true,
	"USDC":  true,
	"BUSD":  true,
	"DAI":   true,
	"TUSD":  true,
	"FDUSD": true,
	"USDP":  true,
}

func isStable(s string) bool {
	_, ok := stableCoins[strings.ToUpper(s)]
	return ok
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	resp, err := http.Get("https://api.kucoin.com/api/v2/symbols")
	if err != nil {
		log.Fatalf("get symbols: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		log.Fatalf("status %d: %s", resp.StatusCode, string(b))
	}

	var info kucoinExchangeInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		log.Fatalf("decode: %v", err)
	}

	log.Printf("symbols from API: %d", len(info.Data))

	// фильтруем активные пары
	pairMap := make(map[string]pairMarket)
	for _, s := range info.Data {
		if !s.EnableTrading || s.BaseCurrency == "" || s.QuoteCurrency == "" {
			continue
		}
		m := pairMarket{
			Symbol: s.Symbol,
			Base:   s.BaseCurrency,
			Quote:  s.QuoteCurrency,
		}
		key := m.Base + "_" + m.Quote
		pairMap[key] = m
	}

	log.Printf("filtered markets: %d", len(pairMap))

	// формируем треугольники
	triangles := buildTriangles(pairMap)

	// сохраняем CSV
	saveCSV("triangles_kucoin.csv", triangles)
	log.Println("Готово: triangles_kucoin.csv")
}

// =======================================================
// Формирование всех треугольников с якорем USDT
// =======================================================
func buildTriangles(pairMap map[string]pairMarket) []Triangle {
	var result []Triangle
	anchor := "USDT"

	for _, m1 := range pairMap {
		// leg1: anchor -> A
		if m1.Quote != anchor {
			continue
		}
		A := m1.Base
		if isStable(A) {
			continue
		}

		for _, m2 := range pairMap {
			// leg2: A -> B
			if m2.Quote != A {
				continue
			}
			B := m2.Base
			if isStable(B) || B == A {
				continue
			}

			// 🔹 Вариант 1: USDT -> A -> B -> USDT
			if m3, ok := pairMap[B+"_"+anchor]; ok {
				result = append(result, Triangle{
					A:    anchor,
					B:    A,
					C:    B,
					Leg1: "BUY " + A + "/" + anchor,
					Leg2: "BUY " + B + "/" + A,
					Leg3: "SELL " + m3.Base + "/" + m3.Quote,
				})
			}

			// 🔹 Вариант 2: USDT -> B -> A -> USDT
			if mBA, okBA := pairMap[B+"_"+A]; okBA {
				if mAU, okAU := pairMap[A+"_"+anchor]; okAU {
					result = append(result, Triangle{
						A:    anchor,
						B:    B,
						C:    A,
						Leg1: "BUY " + B + "/" + anchor,
						Leg2: "BUY " + A + "/" + B,
						Leg3: "SELL " + mAU.Base + "/" + mAU.Quote,
					})
				}
			}
		}
	}

	log.Printf("found triangles: %d", len(result))
	return result
}

// =======================================================
// Запись CSV
// =======================================================
func saveCSV(filename string, data []Triangle) {
	f, err := os.Create(filename)
	if err != nil {
		log.Fatalf("create file: %v", err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	// заголовок
	w.Write([]string{"A", "B", "C", "leg1", "leg2", "leg3"})

	for _, t := range data {
		w.Write([]string{t.A, t.B, t.C, t.Leg1, t.Leg2, t.Leg3})
	}
}



