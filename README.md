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

//
// ==========================
// Универсальная модель рынка
// ==========================
type Market struct {
	Symbol string
	Base   string
	Quote  string

	EnableTrading bool

	BaseMinSize    string
	QuoteMinSize   string
	BaseIncrement  string
	QuoteIncrement string
	PriceIncrement string
}

//
// ==========================
// Треугольник
// ==========================
type Triangle struct {
	A, B, C string

	Leg1 string
	Leg2 string
	Leg3 string

	BaseMin1, QuoteMin1, BaseInc1, QuoteInc1, PriceInc1 string
	BaseMin2, QuoteMin2, BaseInc2, QuoteInc2, PriceInc2 string
	BaseMin3, QuoteMin3, BaseInc3, QuoteInc3, PriceInc3 string
}

//
// ==========================
// Stable coins
// ==========================
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
	return stableCoins[strings.ToUpper(s)]
}

//
// ==========================
// KuCoin API structs
// ==========================
type KuCoinSymbol struct {
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

type KuCoinResponse struct {
	Code string         `json:"code"`
	Data []KuCoinSymbol `json:"data"`
}

//
// ==========================
// MAIN
// ==========================
func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	markets := loadKuCoinMarkets()
	log.Printf("markets loaded: %d", len(markets))

	triangles := buildTriangles(markets, "USDT")
	log.Printf("triangles found: %d", len(triangles))

	saveCSV("triangles.csv", triangles)

	log.Println("done ✅")
}

//
// ==========================
// Загрузка KuCoin
// ==========================
func loadKuCoinMarkets() map[string]Market {
	resp, err := http.Get("https://api.kucoin.com/api/v2/symbols")
	if err != nil {
		log.Fatalf("http error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		log.Fatalf("bad status %d: %s", resp.StatusCode, body)
	}

	var api KuCoinResponse
	if err := json.NewDecoder(resp.Body).Decode(&api); err != nil {
		log.Fatalf("decode error: %v", err)
	}

	markets := make(map[string]Market)

	for _, s := range api.Data {
		if !s.EnableTrading || s.BaseCurrency == "" || s.QuoteCurrency == "" {
			continue
		}

		key := s.BaseCurrency + "_" + s.QuoteCurrency

		markets[key] = Market{
			Symbol:          s.Symbol,
			Base:            s.BaseCurrency,
			Quote:           s.QuoteCurrency,
			EnableTrading:   s.EnableTrading,
			BaseMinSize:     s.BaseMinSize,
			QuoteMinSize:    s.QuoteMinSize,
			BaseIncrement:   s.BaseIncrement,
			QuoteIncrement:  s.QuoteIncrement,
			PriceIncrement:  s.PriceIncrement,
		}
	}

	return markets
}

//
// ==========================
// Генератор треугольников
// ==========================
func buildTriangles(markets map[string]Market, anchor string) []Triangle {
	var out []Triangle

	for _, m1 := range markets {
		// anchor → A
		if m1.Quote != anchor || isStable(m1.Base) {
			continue
		}
		A := m1.Base

		for _, m2 := range markets {
			// A → B
			if m2.Quote != A || isStable(m2.Base) || m2.Base == A {
				continue
			}
			B := m2.Base

			// ===== вариант: USDT → A → B → USDT
			if m3, ok := markets[B+"_"+anchor]; ok {
				out = append(out, newTriangle(anchor, A, B, m1, m2, m3))
			}
		}
	}

	return out
}

//
// ==========================
// Конструктор треугольника
// ==========================
func newTriangle(A, B, C string, l1, l2, l3 Market) Triangle {
	return Triangle{
		A: A,
		B: B,
		C: C,

		Leg1: "BUY " + l1.Base + "/" + l1.Quote,
		Leg2: "BUY " + l2.Base + "/" + l2.Quote,
		Leg3: "SELL " + l3.Base + "/" + l3.Quote,

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

//
// ==========================
// CSV
// ==========================
func saveCSV(filename string, data []Triangle) {
	f, err := os.Create(filename)
	if err != nil {
		log.Fatalf("file create error: %v", err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	w.Write([]string{
		"A", "B", "C",
		"leg1", "leg2", "leg3",
		"baseMin1", "quoteMin1", "baseInc1", "quoteInc1", "priceInc1",
		"baseMin2", "quoteMin2", "baseInc2", "quoteInc2", "priceInc2",
		"baseMin3", "quoteMin3", "baseInc3", "quoteInc3", "priceInc3",
	})

	for _, t := range data {
		w.Write([]string{
			t.A, t.B, t.C,
			t.Leg1, t.Leg2, t.Leg3,
			t.BaseMin1, t.QuoteMin1, t.BaseInc1, t.QuoteInc1, t.PriceInc1,
			t.BaseMin2, t.QuoteMin2, t.BaseInc2, t.QuoteInc2, t.PriceInc2,
			t.BaseMin3, t.QuoteMin3, t.BaseInc3, t.QuoteInc3, t.PriceInc3,
		})
	}
}
