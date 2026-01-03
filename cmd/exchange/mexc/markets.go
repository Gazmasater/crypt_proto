package mexc

import (
	"crypt_proto/cmd/exchange/common"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type mexcSymbol struct {
	Symbol     string `json:"symbol"`
	BaseAsset  string `json:"baseAsset"`
	QuoteAsset string `json:"quoteAsset"`
	Status     string `json:"status"`

	IsSpotTradingAllowed bool `json:"isSpotTradingAllowed"`

	BaseAssetPrecision  int `json:"baseAssetPrecision"`
	QuoteAssetPrecision int `json:"quoteAssetPrecision"`

	Filters []struct {
		FilterType string `json:"filterType"`
		MinQty     string `json:"minQty"`
		StepSize   string `json:"stepSize"`
	} `json:"filters"`
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

	var api mexcResponse
	if err := json.NewDecoder(resp.Body).Decode(&api); err != nil {
		log.Fatalf("decode error: %v", err)
	}

	markets := make(map[string]common.Market)
	fmt.Println("STATUS", api.Data)

	for _, s := range api.Data {
		fmt.Println("STATUS", s.Symbol)
		if s.Status != "ENABLED" || !s.IsSpotTradingAllowed {
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
			Symbol:        s.Symbol,
			Base:          s.BaseAsset,
			Quote:         s.QuoteAsset,
			EnableTrading: true,

			BaseMinSize:   minQty,
			BaseIncrement: stepSize,
		}
	}

	return markets
}
