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

type ExchangeInfo struct {
	Symbols []struct {
		Symbol   string `json:"symbol"`
		Status   string `json:"status"`
		BaseAsset  string `json:"baseAsset"`
		QuoteAsset string `json:"quoteAsset"`
		Filters []struct {
			FilterType string `json:"filterType"`
			StepSize   string `json:"stepSize"`   // lotSize
			MinQty     string `json:"minQty"`    // минимальный объём
			MaxQty     string `json:"maxQty"`    // максимальный объём
			TickSize   string `json:"tickSize"`  // шаг цены
			MaxNotional string `json:"maxNotional"` // максимальная сумма
		} `json:"filters"`
	} `json:"symbols"`
}

type Market struct {
	Symbol     string
	Base       string
	Quote      string
	Status     string
	LotSize    string
	MinQty     string
	MaxQty     string
	TickSize   string
	MaxNotional string
}

// Определяем направление ноги
func determineDirection(from, to, base, quote string) string {
	if strings.EqualFold(from, base) && strings.EqualFold(to, quote) {
		return "SELL"
	} else if strings.EqualFold(from, quote) && strings.EqualFold(to, base) {
		return "BUY"
	}
	return "UNKNOWN"
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	inputCSV := "triangles_routes.csv"
	outputCSV := "triangles_routes_full.csv"

	// 1. Загружаем exchangeInfo с MEXC
	resp, err := http.Get("https://www.mexc.com/api/v3/exchangeInfo")
	if err != nil {
		log.Fatalf("cannot get exchangeInfo: %v", err)
	}
	defer resp.Body.Close()

	var info ExchangeInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		log.Fatalf("decode exchangeInfo: %v", err)
	}

	markets := make(map[string]Market)
	for _, s := range info.Symbols {
		if s.Status != "TRADING" && s.Status != "ENABLED" && s.Status != "1" {
			continue
		}
		m := Market{
			Symbol: s.Symbol,
			Base:   s.BaseAsset,
			Quote:  s.QuoteAsset,
			Status: s.Status,
		}
		// ищем фильтры
		for _, f := range s.Filters {
			switch strings.ToUpper(f.FilterType) {
			case "LOT_SIZE":
				m.LotSize = f.StepSize
				m.MinQty = f.MinQty
				m.MaxQty = f.MaxQty
			case "PRICE_FILTER":
				m.TickSize = f.TickSize
				if f.MaxNotional != "" {
					m.MaxNotional = f.MaxNotional
				}
			}
		}
		markets[s.Symbol] = m
	}
	log.Printf("Loaded markets: %d", len(markets))

	// 2. Читаем исходный CSV маршрутов
	inFile, err := os.Open(inputCSV)
	if err != nil {
		log.Fatalf("open input CSV: %v", err)
	}
	defer inFile.Close()

	reader := csv.NewReader(inFile)
	header, err := reader.Read()
	if err != nil {
		log.Fatalf("read header: %v", err)
	}
	if len(header) < 6 {
		log.Fatalf("CSV must have 6 columns: base1,quote1,base2,quote2,base3,quote3")
	}

	// 3. Создаём выходной CSV
	outFile, err := os.Create(outputCSV)
	if err != nil {
		log.Fatalf("create output CSV: %v", err)
	}
	defer outFile.Close()

	writer := csv.NewWriter(outFile)
	defer writer.Flush()

	// Заголовок
	writer.Write([]string{
		"leg1_symbol", "leg1_from", "leg1_to", "leg1_lotSize", "leg1_minQty", "leg1_maxQty", "leg1_tickSize", "leg1_maxNotional",
		"leg2_symbol", "leg2_from", "leg2_to", "leg2_lotSize", "leg2_minQty", "leg2_maxQty", "leg2_tickSize", "leg2_maxNotional",
		"leg3_symbol", "leg3_from", "leg3_to", "leg3_lotSize", "leg3_minQty", "leg3_maxQty", "leg3_tickSize", "leg3_maxNotional",
		"start_amt", "end_amt", "fail_reason",
	})

	// 4. Обрабатываем строки
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("read row: %v", err)
		}
		base1, quote1 := row[0], row[1]
		base2, quote2 := row[2], row[3]
		base3, quote3 := row[4], row[5]

		legs := []struct{From, To string; Base, Quote string}{
			{From: base1, To: quote1, Base: base1, Quote: quote1},
			{From: base2, To: quote2, Base: base2, Quote: quote2},
			{From: base3, To: quote3, Base: base3, Quote: quote3},
		}

		outRow := make([]string, 0, 30)
		fail := ""

		for _, leg := range legs {
			found := false
			for _, m := range markets {
				if (strings.EqualFold(leg.Base, m.Base) && strings.EqualFold(leg.Quote, m.Quote)) ||
					(strings.EqualFold(leg.Base, m.Quote) && strings.EqualFold(leg.Quote, m.Base)) {
					dir := determineDirection(leg.From, leg.To, m.Base, m.Quote)
					if dir == "UNKNOWN" {
						fail = "cannot determine direction"
					}
					outRow = append(outRow, m.Symbol, leg.From, leg.To, m.LotSize, m.MinQty, m.MaxQty, m.TickSize, m.MaxNotional)
					found = true
					break
				}
			}
			if !found {
				outRow = append(outRow, "", leg.From, leg.To, "", "", "", "", "")
				if fail == "" {
					fail = "pair not found"
				}
			}
		}

		// start_amt, end_amt, fail_reason
		outRow = append(outRow, "25.0", "", fail)

		if err := writer.Write(outRow); err != nil {
			log.Fatalf("write row: %v", err)
		}
	}

	log.Printf("Done! Output saved to %s", outputCSV)
}




