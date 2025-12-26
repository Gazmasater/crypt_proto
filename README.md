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
		Status            string `json:"state"` // "ENABLED" / "DISABLED"
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
	Status   string
	LotSize  string
	MinQty   string
	MaxQty   string
	TickSize string
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

	inputCSV := "triangles_markets.csv"
	outputCSV := "triangles_routes_full.csv"

	// 1. Получаем данные с MEXC
	resp, err := http.Get("https://www.mexc.com/api/v2/market/symbols")
	if err != nil {
		log.Fatalf("cannot get exchangeInfo: %v", err)
	}
	defer resp.Body.Close()

	var info MEXCExchangeInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		log.Fatalf("decode exchangeInfo: %v", err)
	}

	markets := make(map[string]Market)
	for _, s := range info.Data {
		if s.Status != "ENABLED" {
			continue
		}
		lotStep := fmt.Sprintf("%.*f", s.QuantityPrecision, 1.0/float64Pow(10, s.QuantityPrecision))
		priceStep := fmt.Sprintf("%.*f", s.PricePrecision, 1.0/float64Pow(10, s.PricePrecision))
		m := Market{
			Symbol:   s.Symbol,
			Base:     s.BaseAsset,
			Quote:    s.QuoteAsset,
			Status:   s.Status,
			LotSize:  lotStep,
			MinQty:   s.MinOrderQty,
			MaxQty:   s.MaxOrderQty,
			TickSize: priceStep,
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
		"leg1_symbol", "leg1_from", "leg1_to", "leg1_lotSize", "leg1_minQty", "leg1_maxQty", "leg1_tickSize",
		"leg2_symbol", "leg2_from", "leg2_to", "leg2_lotSize", "leg2_minQty", "leg2_maxQty", "leg2_tickSize",
		"leg3_symbol", "leg3_from", "leg3_to", "leg3_lotSize", "leg3_minQty", "leg3_maxQty", "leg3_tickSize",
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

		legs := []struct{ From, To, Base, Quote string }{
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
					outRow = append(outRow, m.Symbol, leg.From, leg.To, m.LotSize, m.MinQty, m.MaxQty, m.TickSize)
					found = true
					break
				}
			}
			if !found {
				outRow = append(outRow, "", leg.From, leg.To, "", "", "", "")
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

// Вспомогательная функция степени для float64
func float64Pow(x, y int) float64 {
	res := 1.0
	for i := 0; i < y; i++ {
		res /= 10
	}
	return res
}



gaz358@gaz358-BOD-WXX9:~/myprog/crypt_proto/cmd/mexctriangl$ go run .
2025/12/26 19:25:07.198175 decode exchangeInfo: invalid character '<' looking for beginning of value
exit status 1



