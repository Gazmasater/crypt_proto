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



type mexcResponse struct {
	Symbols []mexcSymbol `json:"symbols"`
}


func LoadMarkets() map[string]common.Market {
	resp, err := http.Get("https://api.mexc.com/api/v3/exchangeInfo")
	if err != nil {
		log.Fatalf("http error: %v", err)
	}
	defer resp.Body.Close()

	var api mexcResponse
	if err := json.NewDecoder(resp.Body).Decode(&api); err != nil {
		log.Fatalf("decode error: %v", err)
	}

	markets := make(map[string]common.Market)

	fmt.Println("TOTAL SYMBOLS:", len(api.Symbols))

	for _, s := range api.Symbols {
		if (s.Status != "ENABLED" && s.Status != "TRADING") || !s.IsSpotTradingAllowed {
			continue
		}

		var minQty, stepSize float64

		for _, f := range s.Filters {
			if f.FilterType == "LOT_SIZE" {
				minQty, _ = strconv.ParseFloat(f.MinQty, 64)
				stepSize, _ = strconv.ParseFloat(f.StepSize, 64)
			}
		}

		key := s.BaseAsset + "_" + s.QuoteAsset

		markets[key] = common.Market{
			Symbol:        s.Symbol, // BTCUSDT
			Base:          s.BaseAsset,
			Quote:         s.QuoteAsset,
			EnableTrading: true,

			BaseMinSize:   minQty,
			BaseIncrement: stepSize,
		}
	}

	fmt.Println("SPOT MARKETS:", len(markets))
	return markets
}


az358@gaz358-BOD-WXX9:~/myprog/crypt_proto/cmd/exchange$ go run .
TOTAL SYMBOLS: 2502
SPOT MARKETS: 0



