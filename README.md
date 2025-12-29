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

// ---------- OKX response structs ----------

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

// ---------- Triangle struct ----------

type Triangle struct {
	A, B, C string
	Leg1    string
	Leg2    string
	Leg3    string
}

// ---------- Helpers ----------

// Нужно, чтобы B и C не начинались на "USD", кроме точного "USDT"
func isBadUsd(x string) bool {
	// если это USDT — нормально
	if x == "USDT" {
		return false
	}
	// если начинается на "USD" — плохой
	return strings.HasPrefix(x, "USD")
}

// Проверка на треугольник, где A == "USDT" и B,C не начинаются с "USD"
func validTriangleUSDT(A, B, C string) bool {
	if A != "USDT" {
		return false
	}
	if isBadUsd(B) || isBadUsd(C) {
		return false
	}
	return true
}

// ---------- Fetch OKX ----------

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

// ---------- Triangle finder ----------

func findTriangles(instruments []OkxInstrument) []Triangle {
	graph := make(map[string]map[string]string)

	// строим граф (и инверсии)
	for _, inst := range instruments {
		base := inst.BaseCcy
		quote := inst.QuoteCcy

		// прямая
		if graph[base] == nil {
			graph[base] = make(map[string]string)
		}
		graph[base][quote] = inst.InstId

		// инверт (обратное направление)
		if graph[quote] == nil {
			graph[quote] = make(map[string]string)
		}
		graph[quote][base] = inst.InstId + "_INV"
	}

	var result []Triangle

	for A, toB := range graph {
		// сразу фильтр по стартовой валюте
		if A != "USDT" {
			continue
		}

		for B, leg1 := range toB {
			// B не должен быть плохим USD‑токеном
			if isBadUsd(B) {
				continue
			}
			if graph[B] == nil {
				continue
			}

			for C, leg2 := range graph[B] {
				// тоже запрещаем
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

				// финальный фильтр
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

// ---------- Main ----------

func main() {
	instruments, err := fetchOkxInstruments()
	if err != nil {
		log.Fatalf("ошибка получения OKX instruments: %v", err)
	}

	triangles := findTriangles(instruments)
	log.Printf("Найдено треугольников (USDT start, без USD*): %d\n", len(triangles))

	file, err := os.Create("triangles_okx_usdt_only.csv")
	if err != nil {
		log.Fatalf("ошибка создания csv: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

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

	log.Println("triangles_okx_usdt_only.csv успешно создан.")
}



