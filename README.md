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




package okx

import (
	"crypt_proto/cmd/exchange/common"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
)

type okxSymbol struct {
	InstID   string `json:"instId"`
	BaseCcy  string `json:"baseCcy"`
	QuoteCcy string `json:"quoteCcy"`
	InstType string `json:"instType"`
	State    string `json:"state"`
	MinSz    string `json:"minSz"`
	TickSz   string `json:"tickSz"`
	BasePrec int    `json:"baseCcyPrecision"`  // иногда отсутствует
	QuotePrec int   `json:"quoteCcyPrecision"` // иногда отсутствует
}

type okxResponse struct {
	Code string      `json:"code"`
	Data []okxSymbol `json:"data"`
}

// LoadMarkets загружает все активные спотовые пары OKX
func LoadMarkets() map[string]common.Market {
	resp, err := http.Get("https://www.okx.com/api/v5/public/instruments?instType=SPOT")
	if err != nil {
		log.Fatalf("HTTP error: %v", err)
	}
	defer resp.Body.Close()

	var api okxResponse
	if err := json.NewDecoder(resp.Body).Decode(&api); err != nil {
		log.Fatalf("JSON decode error: %v", err)
	}

	markets := make(map[string]common.Market)
	fmt.Println("TOTAL SYMBOLS:", len(api.Data))

	for _, s := range api.Data {
		if s.InstType != "SPOT" || s.State != "live" {
			continue
		}

		minQty := parseFloatFallback(s.MinSz, s.BasePrec)
		stepSize := parseFloatFallback(s.TickSz, s.BasePrec)

		key := s.BaseCcy + "_" + s.QuoteCcy
		markets[key] = common.Market{
			Symbol:        s.InstID,
			Base:          s.BaseCcy,
			Quote:         s.QuoteCcy,
			EnableTrading: true,
			BaseMinSize:   minQty,
			BaseIncrement: stepSize,
		}
	}

	fmt.Println("SPOT MARKETS:", len(markets))
	return markets
}

// parseFloatFallback конвертирует строку в float64 с fallback по точности
func parseFloatFallback(s string, precision int) float64 {
	if s == "" {
		return math.Pow10(-precision)
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return math.Pow10(-precision)
	}
	return f
}




package main

import (
	"log"

	"crypt_proto/cmd/exchange/builder"
	"crypt_proto/cmd/exchange/common"
	"crypt_proto/cmd/exchange/kucoin"
	"crypt_proto/cmd/exchange/mexc"
	"crypt_proto/cmd/exchange/okx"
)

func main() {

	// ---------- KUCOIN ----------
	kucoinMarkets := kucoin.LoadMarkets()
	kucoinTriangles := builder.BuildTriangles(kucoinMarkets, "USDT")
	if err := common.SaveTrianglesCSV(
		"data/kucoin_triangles_usdt.csv",
		kucoinTriangles,
	); err != nil {
		log.Fatalf("kucoin csv error: %v", err)
	}

	// ---------- MEXC ----------
	mexcMarkets := mexc.LoadMarkets()
	mexcTriangles := builder.BuildTriangles(mexcMarkets, "USDT")
	if err := common.SaveTrianglesCSV(
		"data/mexc_triangles_usdt.csv",
		mexcTriangles,
	); err != nil {
		log.Fatalf("mexc csv error: %v", err)
	}

	// ---------- OKX ----------
	okxMarkets := okx.LoadMarkets()
	okxTriangles := builder.BuildTriangles(okxMarkets, "USDT")
	if err := common.SaveTrianglesCSV(
		"data/okx_triangles_usdt.csv",
		okxTriangles,
	); err != nil {
		log.Fatalf("okx csv error: %v", err)
	}
}



