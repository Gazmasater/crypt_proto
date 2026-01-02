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



Я покажу:

📁 структуру папок

📦 общий пакет exchange

🔌 интерфейс биржи

🧱 модель Market

🔍 универсальный поиск leg

🔄 определение BUY / SELL

🔺 генератор треугольников

💾 где и как хранить CSV

✅ итог: что у тебя получается

✅ 1. Структура проекта (рекомендую)
exchange/
├── common/
│   ├── market.go
│   ├── triangle.go
│   ├── resolver.go
│   └── csv.go
│
├── kucoin/
│   ├── client.go
│   └── markets.go
│
├── okx/
│   ├── client.go
│   └── markets.go
│
├── mexc/
│   ├── client.go
│   └── markets.go
│
├── builder/
│   └── triangles.go
│
└── main.go

✅ 2. exchange/common/market.go
package common

type Market struct {
	Symbol string

	Base  string
	Quote string

	EnableTrading bool

	BaseMinSize  float64
	QuoteMinSize float64

	BaseIncrement  float64
	QuoteIncrement float64
	PriceIncrement float64
}

✅ 3. Универсальный поиск пары (в любую сторону)

📄 exchange/common/resolver.go

package common

func FindLeg(from, to string, markets map[string]Market) (Market, bool) {
	if m, ok := markets[from+"_"+to]; ok {
		return m, true
	}
	if m, ok := markets[to+"_"+from]; ok {
		return m, true
	}
	return Market{}, false
}

✅ 4. Определение BUY / SELL
package common

func ResolveSide(from, to string, m Market) string {
	if m.Quote == from && m.Base == to {
		return "BUY"
	}
	if m.Base == from && m.Quote == to {
		return "SELL"
	}
	return ""
}

✅ 5. Структура треугольника

📄 exchange/common/triangle.go

package common

type Triangle struct {
	A string
	B string
	C string

	Leg1 string
	Leg2 string
	Leg3 string

	BaseMin1  float64
	QuoteMin1 float64
	BaseInc1  float64
	QuoteInc1 float64
	PriceInc1 float64

	BaseMin2  float64
	QuoteMin2 float64
	BaseInc2  float64
	QuoteInc2 float64
	PriceInc2 float64

	BaseMin3  float64
	QuoteMin3 float64
	BaseInc3  float64
	QuoteInc3 float64
	PriceInc3 float64
}

✅ 6. Универсальный конструктор треугольника
package common

func NewTriangle(A, B, C string, l1, l2, l3 Market) Triangle {
	return Triangle{
		A: A,
		B: B,
		C: C,

		Leg1: ResolveSide(A, B, l1) + " " + l1.Base + "/" + l1.Quote,
		Leg2: ResolveSide(B, C, l2) + " " + l2.Base + "/" + l2.Quote,
		Leg3: ResolveSide(C, A, l3) + " " + l3.Base + "/" + l3.Quote,

		BaseMin1:  l1.BaseMinSize,
		QuoteMin1: l1.QuoteMinSize,
		BaseInc1:  l1.BaseIncrement,
		QuoteInc1: l1.QuoteIncrement,
		PriceInc1: l1.PriceIncrement,

		BaseMin2:  l2.BaseMinSize,
		QuoteMin2: l2.QuoteMinSize,
		BaseInc2:  l2.BaseIncrement,
		QuoteInc2: l2.QuoteIncrement,
		PriceInc2: l2.PriceIncrement,

		BaseMin3:  l3.BaseMinSize,
		QuoteMin3: l3.QuoteMinSize,
		BaseInc3:  l3.BaseIncrement,
		QuoteInc3: l3.QuoteIncrement,
		PriceInc3: l3.PriceIncrement,
	}
}

✅ 7. Генератор треугольников

📄 exchange/builder/triangles.go

package builder

import "exchange/common"

func BuildTriangles(
	markets map[string]common.Market,
	anchor string,
) []common.Triangle {

	var result []common.Triangle

	for _, m1 := range markets {
		if !m1.EnableTrading {
			continue
		}

		var B string
		if m1.Base == anchor {
			B = m1.Quote
		} else if m1.Quote == anchor {
			B = m1.Base
		} else {
			continue
		}

		for _, m2 := range markets {
			if !m2.EnableTrading {
				continue
			}

			var C string
			if m2.Base == B {
				C = m2.Quote
			} else if m2.Quote == B {
				C = m2.Base
			} else {
				continue
			}

			if C == anchor || C == B {
				continue
			}

			l3, ok := common.FindLeg(C, anchor, markets)
			if !ok {
				continue
			}

			l1, ok1 := common.FindLeg(anchor, B, markets)
			l2, ok2 := common.FindLeg(B, C, markets)

			if !ok1 || !ok2 {
				continue
			}

			t := common.NewTriangle(anchor, B, C, l1, l2, l3)
			result = append(result, t)
		}
	}

	return result
}

✅ 8. CSV — где хранить и как

📁 рекомендую:

data/
├── kucoin_triangles.csv
├── okx_triangles.csv
└── mexc_triangles.csv


📄 exchange/common/csv.go

package common

import (
	"encoding/csv"
	"os"
)

func SaveTrianglesCSV(path string, list []Triangle) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	w.Write([]string{
		"A", "B", "C",
		"Leg1", "Leg2", "Leg3",
	})

	for _, t := range list {
		w.Write([]string{
			t.A, t.B, t.C,
			t.Leg1, t.Leg2, t.Leg3,
		})
	}

	return nil
}

✅ 9. Пример использования (main.go)
package main

import (
	"exchange/builder"
	"exchange/common"
	"exchange/kucoin"
)

func main() {
	markets := kucoin.LoadMarkets()

	triangles := builder.BuildTriangles(markets, "USDT")

	common.SaveTrianglesCSV("data/kucoin_triangles.csv", triangles)
}



package kucoin

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"exchange/common"
)

// Структуры API KuCoin
type kuCoinSymbol struct {
	Symbol         string `json:"symbol"`
	BaseCurrency   string `json:"baseCurrency"`
	QuoteCurrency  string `json:"quoteCurrency"`
	EnableTrading  bool   `json:"enableTrading"`
	BaseMinSize    string `json:"baseMinSize"`
	QuoteMinSize   string `json:"quoteMinSize"`
	BaseIncrement  string `json:"baseIncrement"`
	QuoteIncrement string `json:"quoteIncrement"`
	PriceIncrement string `json:"priceIncrement"`
}

type kuCoinResponse struct {
	Code string         `json:"code"`
	Data []kuCoinSymbol `json:"data"`
}

// LoadMarkets загружает все рынки KuCoin и возвращает map[string]common.Market
func LoadMarkets() map[string]common.Market {
	resp, err := http.Get("https://api.kucoin.com/api/v2/symbols")
	if err != nil {
		log.Fatalf("http error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		log.Fatalf("bad status %d: %s", resp.StatusCode, body)
	}

	var api kuCoinResponse
	if err := json.NewDecoder(resp.Body).Decode(&api); err != nil {
		log.Fatalf("decode error: %v", err)
	}

	markets := make(map[string]common.Market)
	for _, s := range api.Data {
		if !s.EnableTrading || s.BaseCurrency == "" || s.QuoteCurrency == "" {
			continue
		}

		// Преобразуем строки в float64
		bMin := parseFloat(s.BaseMinSize)
		qMin := parseFloat(s.QuoteMinSize)
		bInc := parseFloat(s.BaseIncrement)
		qInc := parseFloat(s.QuoteIncrement)
		pInc := parseFloat(s.PriceIncrement)

		key := s.BaseCurrency + "_" + s.QuoteCurrency
		markets[key] = common.Market{
			Symbol:        s.Symbol,
			Base:          s.BaseCurrency,
			Quote:         s.QuoteCurrency,
			EnableTrading: s.EnableTrading,
			BaseMinSize:   bMin,
			QuoteMinSize:  qMin,
			BaseIncrement: bInc,
			QuoteIncrement: qInc,
			PriceIncrement: pInc,
		}
	}

	return markets
}

func parseFloat(s string) float64 {
	var f float64
	if s == "" {
		return 0
	}
	_, err := fmt.Sscan(s, &f)
	if err != nil {
		return 0
	}
	return f
}


