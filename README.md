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



package mexc

import (
	"crypt_proto/cmd/exchange/common"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
)

type mexcSymbol struct {
	Symbol              string `json:"symbol"`
	BaseAsset           string `json:"baseAsset"`
	QuoteAsset          string `json:"quoteAsset"`
	Status              string `json:"status"`
	BaseSizePrecision   int    `json:"baseSizePrecision"`
	QuotePrecision      int    `json:"quotePrecision"`
	MinQty              string `json:"minQty"`
	IsSpotTradingAllowed bool  `json:"isSpotTradingAllowed"`
}

type mexcResponse struct {
	Data []mexcSymbol `json:"data"`
}

func LoadMarkets() map[string]common.Market {
	resp, err := http.Get("https://api.mexc.com/api/v3/exchangeInfo")
	if err != nil {
		log.Fatalf("http error: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var api mexcResponse
	if err := json.Unmarshal(body, &api); err != nil {
		log.Fatalf("decode error: %v", err)
	}

	markets := make(map[string]common.Market)

	for _, s := range api.Data {
		if s.Status != "ENABLED" || !s.IsSpotTradingAllowed {
			continue
		}

		minQty, _ := strconv.ParseFloat(s.MinQty, 64)

		baseStep := precisionToStep(s.BaseSizePrecision)
		quoteStep := precisionToStep(s.QuotePrecision)

		key := s.BaseAsset + "_" + s.QuoteAsset

		markets[key] = common.Market{
			Symbol:        s.Symbol,
			Base:          s.BaseAsset,
			Quote:         s.QuoteAsset,
			EnableTrading: true,
			BaseMinSize:   minQty,
			BaseIncrement: baseStep,
			QuoteIncrement: quoteStep,
		}
	}

	return markets
}

func precisionToStep(p int) float64 {
	if p <= 0 {
		return 1
	}
	step := 1.0
	for i := 0; i < p; i++ {
		step /= 10
	}
	return step
}



