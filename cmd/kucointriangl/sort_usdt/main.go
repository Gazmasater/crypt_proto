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

		BaseMinSize    string `json:"baseMinSize"`
		QuoteMinSize   string `json:"quoteMinSize"`
		BaseIncrement  string `json:"baseIncrement"`
		QuoteIncrement string `json:"quoteIncrement"`
		PriceIncrement string `json:"priceIncrement"`
	} `json:"data"`
}

type pairMarket struct {
	Symbol         string
	Base           string
	Quote          string
	BaseMinSize    string
	QuoteMinSize   string
	BaseIncrement  string
	QuoteIncrement string
	PriceIncrement string
}

type Triangle struct {
	A    string
	B    string
	C    string
	Leg1 string
	Leg2 string
	Leg3 string
	// Добавим параметры каждой ноги
	BaseMin1, QuoteMin1, BaseInc1, QuoteInc1, PriceInc1 string
	BaseMin2, QuoteMin2, BaseInc2, QuoteInc2, PriceInc2 string
	BaseMin3, QuoteMin3, BaseInc3, QuoteInc3, PriceInc3 string
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
	"USD1":  true,
	"USDG":  true,
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
			Symbol:         s.Symbol,
			Base:           s.BaseCurrency,
			Quote:          s.QuoteCurrency,
			BaseMinSize:    s.BaseMinSize,
			QuoteMinSize:   s.QuoteMinSize,
			BaseIncrement:  s.BaseIncrement,
			QuoteIncrement: s.QuoteIncrement,
			PriceIncrement: s.PriceIncrement,
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

					BaseMin1:  m1.BaseMinSize,
					QuoteMin1: m1.QuoteMinSize,
					BaseInc1:  m1.BaseIncrement,
					QuoteInc1: m1.QuoteIncrement,
					PriceInc1: m1.PriceIncrement,
					BaseMin2:  m2.BaseMinSize,
					QuoteMin2: m2.QuoteMinSize,
					BaseInc2:  m2.BaseIncrement,
					QuoteInc2: m2.QuoteIncrement,
					PriceInc2: m2.PriceIncrement,
					BaseMin3:  m3.BaseMinSize,
					QuoteMin3: m3.QuoteMinSize,
					BaseInc3:  m3.BaseIncrement,
					QuoteInc3: m3.QuoteIncrement,
					PriceInc3: m3.PriceIncrement,
				})
			}

			// 🔹 Вариант 2: USDT -> B -> A -> USDT
			if _, okBA := pairMap[B+"_"+A]; okBA {
				if mAU, okAU := pairMap[A+"_"+anchor]; okAU {
					result = append(result, Triangle{
						A:    anchor,
						B:    B,
						C:    A,
						Leg1: "BUY " + B + "/" + anchor,
						Leg2: "BUY " + A + "/" + B,
						Leg3: "SELL " + mAU.Base + "/" + mAU.Quote,

						BaseMin1:  m2.BaseMinSize,
						QuoteMin1: m2.QuoteMinSize,
						BaseInc1:  m2.BaseIncrement,
						QuoteInc1: m2.QuoteIncrement,
						PriceInc1: m2.PriceIncrement,
						BaseMin2:  m1.BaseMinSize,
						QuoteMin2: m1.QuoteMinSize,
						BaseInc2:  m1.BaseIncrement,
						QuoteInc2: m1.QuoteIncrement,
						PriceInc2: m1.PriceIncrement,
						BaseMin3:  mAU.BaseMinSize,
						QuoteMin3: mAU.QuoteMinSize,
						BaseInc3:  mAU.BaseIncrement,
						QuoteInc3: mAU.QuoteIncrement,
						PriceInc3: mAU.PriceIncrement,
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
	w.Write([]string{
		"A", "B", "C", "leg1", "leg2", "leg3",
		"baseMin1", "quoteMin1", "baseInc1", "quoteInc1", "priceInc1",
		"baseMin2", "quoteMin2", "baseInc2", "quoteInc2", "priceInc2",
		"baseMin3", "quoteMin3", "baseInc3", "quoteInc3", "priceInc3",
	})

	for _, t := range data {
		w.Write([]string{
			t.A, t.B, t.C, t.Leg1, t.Leg2, t.Leg3,
			t.BaseMin1, t.QuoteMin1, t.BaseInc1, t.QuoteInc1, t.PriceInc1,
			t.BaseMin2, t.QuoteMin2, t.BaseInc2, t.QuoteInc2, t.PriceInc2,
			t.BaseMin3, t.QuoteMin3, t.BaseInc3, t.QuoteInc3, t.PriceInc3,
		})
	}
}
