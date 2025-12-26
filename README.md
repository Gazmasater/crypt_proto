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
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type MEXCExchangeInfo struct {
	Data []struct {
		Symbol            string `json:"symbol"`
		BaseAsset         string `json:"baseAsset"`
		QuoteAsset        string `json:"quoteAsset"`
		State             string `json:"state"`
		PricePrecision    int    `json:"pricePrecision"`
		QuantityPrecision int    `json:"quantityPrecision"`
		MinOrderQty       string `json:"minOrderQty"`
		MaxOrderQty       string `json:"maxOrderQty"`
	} `json:"data"`
}

type Market struct {
	Symbol   string
	Base     string
	Quote    string
	LotSize  string
	MinQty   string
	MaxQty   string
	TickSize string
}

func pow10neg(n int) string {
	if n <= 0 {
		return "1"
	}
	return "0." + strings.Repeat("0", n-1) + "1"
}

// BUY если from=quote → to=base
// SELL если from=base → to=quote
func direction(from, to, base, quote string) string {
	switch {
	case strings.EqualFold(from, base) && strings.EqualFold(to, quote):
		return "SELL"
	case strings.EqualFold(from, quote) && strings.EqualFold(to, base):
		return "BUY"
	default:
		return "UNKNOWN"
	}
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	inputCSV := "triangles_markets.csv"
	outputCSV := "triangles_routes_full.csv"

	// === HTTP REQUEST WITH HEADERS (IMPORTANT) ===
	req, err := http.NewRequest("GET", "https://www.mexc.com/api/v2/market/symbols", nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		log.Fatalf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	var info MEXCExchangeInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		log.Fatalf("decode exchangeInfo: %v", err)
	}

	// === BUILD MARKET MAP ===
	markets := make(map[string]Market)

	for _, s := range info.Data {
		if s.State != "ENABLED" {
			continue
		}

		markets[s.Symbol] = Market{
			Symbol:   s.Symbol,
			Base:     s.BaseAsset,
			Quote:    s.QuoteAsset,
			LotSize:  pow10neg(s.QuantityPrecision),
			TickSize: pow10neg(s.PricePrecision),
			MinQty:   s.MinOrderQty,
			MaxQty:   s.MaxOrderQty,
		}
	}

	log.Printf("Loaded active markets: %d", len(markets))

	// === READ TRIANGLE ROUTES ===
	inFile, err := os.Open(inputCSV)
	if err != nil {
		log.Fatalf("open input csv: %v", err)
	}
	defer inFile.Close()

	reader := csv.NewReader(inFile)
	header, err := reader.Read()
	if err != nil {
		log.Fatalf("read header: %v", err)
	}

	if len(header) < 6 {
		log.Fatalf("CSV must contain 6 columns")
	}

	// === OUTPUT FILE ===
	outFile, err := os.Create(outputCSV)
	if err != nil {
		log.Fatalf("create output csv: %v", err)
	}
	defer outFile.Close()

	writer := csv.NewWriter(outFile)
	defer writer.Flush()

	writer.Write([]string{
		"leg1_symbol", "leg1_from", "leg1_to", "leg1_lot", "leg1_min", "leg1_max", "leg1_tick",
		"leg2_symbol", "leg2_from", "leg2_to", "leg2_lot", "leg2_min", "leg2_max", "leg2_tick",
		"leg3_symbol", "leg3_from", "leg3_to", "leg3_lot", "leg3_min", "leg3_max", "leg3_tick",
		"start_amount",
		"fail_reason",
	})

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		base1, quote1 := row[0], row[1]
		base2, quote2 := row[2], row[3]
		base3, quote3 := row[4], row[5]

		legs := []struct {
			From  string
			To    string
			Base  string
			Quote string
		}{
			{base1, quote1, base1, quote1},
			{base2, quote2, base2, quote2},
			{base3, quote3, base3, quote3},
		}

		out := make([]string, 0, 30)
		fail := ""

		for _, leg := range legs {
			found := false

			for _, m := range markets {
				if (m.Base == leg.Base && m.Quote == leg.Quote) ||
					(m.Base == leg.Quote && m.Quote == leg.Base) {

					dir := direction(leg.From, leg.To, m.Base, m.Quote)
					if dir == "UNKNOWN" {
						fail = "direction_error"
					}

					out = append(out,
						m.Symbol,
						leg.From,
						leg.To,
						m.LotSize,
						m.MinQty,
						m.MaxQty,
						m.TickSize,
					)
					found = true
					break
				}
			}

			if !found {
				out = append(out, "", leg.From, leg.To, "", "", "", "")
				if fail == "" {
					fail = "pair_not_found"
				}
			}
		}

		out = append(out, "25.0", fail)

		if err := writer.Write(out); err != nil {
			log.Fatal(err)
		}
	}

	log.Printf("DONE → %s", outputCSV)
}




