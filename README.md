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

// ---------------- OKX RESPONSE ----------------

type OkxInstrument struct {
	InstId   string `json:"instId"`
	BaseCcy  string `json:"baseCcy"`
	QuoteCcy string `json:"quoteCcy"`
}

type OkxResponse struct {
	Code string          `json:"code"`
	Msg  string          `json:"msg"`
	Data []OkxInstrument `json:"data"`
}

// ---------------- TRIANGLE STRUCT ----------------

type Triangle struct {
	A, B, C string
	Leg1    string
	Leg2    string
	Leg3    string
}

// ---------------- FILTERS ----------------

var forbiddenUSD = map[string]bool{
	"USDC": true,
	"USD1": true,
	// добавь сюда другие, если нужно
}

func isStable(x string) bool {
	if x == "USDT" || forbiddenUSD[x] {
		return true
	}
	return false
}

// валидность треугольника: стартует с USDT и нет других USD‑стейблов
func validTriangleUSDT(A, B, C string) bool {
	if A != "USDT" {
		return false
	}
	if forbiddenUSD[B] || forbiddenUSD[C] {
		return false
	}
	return true
}

// ---------------- FETCH OKX ----------------

func fetchOkxInstruments() ([]OkxInstrument, error) {
	url := "https://www.okx.com/api/v5/public/instruments?instType=SPOT"
	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var res OkxResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}
	return res.Data, nil
}

// ---------------- TRIANGLE LOGIC ----------------

func findTriangles(instruments []OkxInstrument) []Triangle {
	// graph[from][to] = pair
	graph := make(map[string]map[string]string)

	// строим граф
	for _, inst := range instruments {
		base := inst.BaseCcy
		quote := inst.QuoteCcy
		// прямая связь
		if _, ok := graph[base]; !ok {
			graph[base] = make(map[string]string)
		}
		graph[base][quote] = inst.InstId

		// инверсия — обозначаем явно
		if _, ok := graph[quote]; !ok {
			graph[quote] = make(map[string]string)
		}
		graph[quote][base] = inst.InstId + "_INV"
	}

	var result []Triangle

	// A→B→C→A
	for A, toB := range graph {

		// старт должен быть USDT
		if A != "USDT" {
			continue
		}

		for B, leg1 := range toB {
			// если нет переходов с B → …
			if graph[B] == nil {
				continue
			}
			for C, leg2 := range graph[B] {
				// если нет переходов с C → …
				if graph[C] == nil {
					continue
				}
				leg3, ok := graph[C][A]
				if !ok {
					continue
				}

				// фильтр: нет больших стейбл‑комбинаций
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

// ---------------- MAIN ----------------

func main() {
	instruments, err := fetchOkxInstruments()
	if err != nil {
		log.Fatalf("Ошибка получения пар OKX: %v", err)
	}

	triangles := findTriangles(instruments)
	log.Printf("Найдено треугольников на OKX (USDT, нет других USD‑стейблов): %d", len(triangles))

	// создаём CSV
	file, err := os.Create("triangles_okx.csv")
	if err != nil {
		log.Fatalf("Ошибка создания файла: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// заголовок
	writer.Write([]string{"A", "B", "C", "leg1", "leg2", "leg3"})

	for _, t := range triangles {
		writer.Write([]string{
			t.A,
			t.B,
			t.C,
			t.Leg1,
			t.Leg2,
			t.Leg3,
		})
	}

	log.Println("triangles_okx.csv успешно создан")
}


