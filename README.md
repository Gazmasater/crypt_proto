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
	"net/http"
	"os"
	"strings"
)

const url = "https://api.mexc.com/api/v3/exchangeInfo"

type ExchangeInfo struct {
	Symbols []Symbol `json:"symbols"`
}

type Symbol struct {
	Symbol                    string   `json:"symbol"`
	BaseAsset                 string   `json:"baseAsset"`
	QuoteAsset                string   `json:"quoteAsset"`
	IsSpotTradingAllowed      bool     `json:"isSpotTradingAllowed"`
	OrderTypes                []string `json:"orderTypes"`
	BaseSizePrecision         string   `json:"baseSizePrecision"`
	QuoteAmountPrecisionMarket string  `json:"quoteAmountPrecisionMarket"`
	TakerCommission           string   `json:"takerCommission"`
}

type Pair struct {
	Symbol     string
	Base       string
	Quote      string
	BaseStep   string
	QuoteStep  string
	TakerFee   string
}

// ---------- helpers ----------

func hasMarket(orderTypes []string) bool {
	for _, t := range orderTypes {
		if t == "MARKET" {
			return true
		}
	}
	return false
}

func fetchPairs() ([]Pair, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ex ExchangeInfo
	if err := json.NewDecoder(resp.Body).Decode(&ex); err != nil {
		return nil, err
	}

	var pairs []Pair
	for _, s := range ex.Symbols {
		if !s.IsSpotTradingAllowed || !hasMarket(s.OrderTypes) {
			continue
		}

		pairs = append(pairs, Pair{
			Symbol:    s.Symbol,
			Base:      s.BaseAsset,
			Quote:     s.QuoteAsset,
			BaseStep:  s.BaseSizePrecision,
			QuoteStep: s.QuoteAmountPrecisionMarket,
			TakerFee:  s.TakerCommission,
		})
	}

	return pairs, nil
}

func main() {
	pairs, err := fetchPairs()
	if err != nil {
		panic(err)
	}

	// быстрый поиск пар
	pairMap := make(map[string]Pair)
	for _, p := range pairs {
		pairMap[p.Base+"_"+p.Quote] = p
	}

	file, err := os.Create("triangles.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	w := csv.NewWriter(file)
	defer w.Flush()

	// header
	w.Write([]string{
		"A", "B", "C",
		"pair",
		"side",
		"base",
		"quote",
		"baseStep",
		"quoteStep",
		"takerFee",
		"invert",
	})

	// A → B → C → A
	for _, p1 := range pairs {
		A := p1.Base
		B := p1.Quote

		for _, p2 := range pairs {
			if p2.Base != B {
				continue
			}
			C := p2.Quote

			if C == A {
				continue
			}

			// ищем третью ногу
			var p3 Pair
			invert := false
			found := false

			// C -> A
			if v, ok := pairMap[C+"_"+A]; ok {
				p3 = v
				invert = false
				found = true
			} else if v, ok := pairMap[A+"_"+C]; ok {
				p3 = v
				invert = true
				found = true
			}

			if !found {
				continue
			}

			// ---- LEG 1 ----
			writeLeg(w, A, B, C, p1, false)

			// ---- LEG 2 ----
			writeLeg(w, A, B, C, p2, false)

			// ---- LEG 3 ----
			writeLeg(w, A, B, C, p3, invert)
		}
	}

	fmt.Println("triangles.csv generated")
}

func writeLeg(w *csv.Writer, A, B, C string, p Pair, invert bool) {
	side := "SELL"
	base := p.Base
	quote := p.Quote

	if invert {
		side = "BUY"
		base = p.Quote
		quote = p.Base
	}

	w.Write([]string{
		A,
		B,
		C,
		p.Symbol,
		side,
		base,
		quote,
		p.BaseStep,
		p.QuoteStep,
		p.TakerFee,
		fmt.Sprintf("%v", invert),
	})
}


