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
	Symbol               string `json:"symbol"`
	BaseAsset            string `json:"baseAsset"`
	QuoteAsset           string `json:"quoteAsset"`
	Status               string `json:"status"`
	BaseSizePrecision    int    `json:"baseSizePrecision"`
	QuotePrecision       int    `json:"quotePrecision"`
	MinQty               string `json:"minQty"`
	IsSpotTradingAllowed bool   `json:"isSpotTradingAllowed"`
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
			Symbol:         s.Symbol,
			Base:           s.BaseAsset,
			Quote:          s.QuoteAsset,
			EnableTrading:  true,
			BaseMinSize:    minQty,
			BaseIncrement:  baseStep,
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
