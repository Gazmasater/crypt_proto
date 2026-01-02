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

// ----------------------------
// Универсальные структуры
// ----------------------------
type Market struct {
	Symbol         string
	Base           string
	Quote          string
	EnableTrading  bool
	BaseMinSize    string
	QuoteMinSize   string
	BaseIncrement  string
	QuoteIncrement string
	PriceIncrement string
}

type Triangle struct {
	A, B, C                                             string
	Leg1, Leg2, Leg3                                    string
	BaseMin1, QuoteMin1, BaseInc1, QuoteInc1, PriceInc1 string
	BaseMin2, QuoteMin2, BaseInc2, QuoteInc2, PriceInc2 string
	BaseMin3, QuoteMin3, BaseInc3, QuoteInc3, PriceInc3 string
}

// ----------------------------
// Стейблкоины
// ----------------------------
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

// ----------------------------
// KuCoin API структуры
// ----------------------------
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

type KuCoinExchangeInfo struct {
	Code string         `json:"code"`
	Data []KuCoinSymbol `json:"data"`
}

// ----------------------------
// Main
// ----------------------------
func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	// 1) Загружаем пары KuCoin
	kucoinMarkets := fetchKuCoin()
	log.Printf("KuCoin markets loaded: %d", len(kucoinMarkets))

	// 2) Построение треугольников
	triangles := BuildTriangles(kucoinMarkets, "USDT")
	log.Printf("Triangles found: %d", len(triangles))

	// 3) Сохраняем CSV
	SaveCSV("triangles.csv", triangles)
	log.Println("Готово: triangles.csv")
}

// ----------------------------
// Функции для загрузки и нормализации KuCoin
// ----------------------------
func fetchKuCoin() map[string]Market {
	resp, err := http.Get("https://api.kucoin.com/api/v2/symbols")
	if err != nil {
		log.Fatalf("get symbols: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		b, _ := io.ReadAll(resp.Body)
		log.Fatalf("status %d: %s", resp.StatusCode, string(b))
	}

	var info KuCoinExchangeInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		log.Fatalf("decode: %v", err)
	}

	return KuCoinToMarketMap(info.Data)
}

func KuCoinToMarketMap(data []KuCoinSymbol) map[string]Market {
	m := make(map[string]Market)
	for _, s := range data {
		if !s.EnableTrading || s.BaseCurrency == "" || s.QuoteCurrency == "" {
			continue
		}
		key := s.BaseCurrency + "_" + s.QuoteCurrency
		m[key] = Market{
			Symbol:         s.Symbol,
			Base:           s.BaseCurrency,
			Quote:          s.QuoteCurrency,
			EnableTrading:  s.EnableTrading,
			BaseMinSize:    s.BaseMinSize,
			QuoteMinSize:   s.QuoteMinSize,
			BaseIncrement:  s.BaseIncrement,
			QuoteIncrement: s.QuoteIncrement,
			PriceIncrement: s.PriceIncrement,
		}
	}
	return m
}

// ----------------------------
// Заглушки для OKX и MEXC (будет подключение аналогично)
// ----------------------------
/*
func OKXToMarketMap(data []OKXSymbol) map[string]Market { ... }
func MEXCToMarketMap(data []MEXCSymbol) map[string]Market { ... }
*/

// ----------------------------
// Генератор треугольников
// ----------------------------
func BuildTriangles(pairMap map[string]Market, anchor string) []Triangle {
	var result []Triangle

	for _, m1 := range pairMap {
		if m1.Quote != anchor || isStable(m1.Base) {
			continue
		}
		A := m1.Base

		for _, m2 := range pairMap {
			if m2.Quote != A || isStable(m2.Base) || m2.Base == A {
				continue
			}
			B := m2.Base

			// Вариант 1: anchor -> A -> B -> anchor
			if m3, ok := pairMap[B+"_"+anchor]; ok {
				result = append(result, NewTriangle(anchor, A, B, m1, m2, m3))
			}

			// Вариант 2: anchor -> B -> A -> anchor
			if _, okBA := pairMap[B+"_"+A]; okBA {
				if mAU, okAU := pairMap[A+"_"+anchor]; okAU {
					result = append(result, NewTriangle(anchor, B, A, pairMap[B+"_"+anchor], pairMap[A+"_"+B], mAU))
				}
			}
		}
	}

	return result
}

// ----------------------------
// Создание треугольника
// ----------------------------
func NewTriangle(A, B, C string, leg1, leg2, leg3 Market) Triangle {
	return Triangle{
		A: A, B: B, C: C,
		Leg1: "BUY " + leg1.Base + "/" + leg1.Quote,
		Leg2: "BUY " + leg2.Base + "/" + leg2.Quote,
		Leg3: "SELL " + leg3.Base + "/" + leg3.Quote,

		BaseMin1:  leg1.BaseMinSize,
		QuoteMin1: leg1.QuoteMinSize,
		BaseInc1:  leg1.BaseIncrement,
		QuoteInc1: leg1.QuoteIncrement,
		PriceInc1: leg1.PriceIncrement,

		BaseMin2:  leg2.BaseMinSize,
		QuoteMin2: leg2.QuoteMinSize,
		BaseInc2:  leg2.BaseIncrement,
		QuoteInc2: leg2.QuoteIncrement,
		PriceInc2: leg2.PriceIncrement,

		BaseMin3:  leg3.BaseMinSize,
		QuoteMin3: leg3.QuoteMinSize,
		BaseInc3:  leg3.BaseIncrement,
		QuoteInc3: leg3.QuoteIncrement,
		PriceInc3: leg3.PriceIncrement,
	}
}

// ----------------------------
// Сохранение CSV
// ----------------------------
func SaveCSV(filename string, data []Triangle) {
	f, err := os.Create(filename)
	if err != nil {
		log.Fatalf("create file: %v", err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

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
