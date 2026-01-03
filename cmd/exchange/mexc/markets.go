package mexc

import (
	"crypt_proto/cmd/exchange/common"
	"encoding/json"
	"log"
	"math"
	"net/http"
)

// Структуры для ответа MEXC
type mexcSymbol struct {
	Symbol      string   `json:"symbol"`
	BaseAsset   string   `json:"baseAsset"`
	QuoteAsset  string   `json:"quoteAsset"`
	Status      string   `json:"status"`
	Permissions []string `json:"permissions"`

	BaseAssetPrecision  int `json:"baseAssetPrecision"`
	QuoteAssetPrecision int `json:"quoteAssetPrecision"`
}

type mexcResponse struct {
	Symbols []mexcSymbol `json:"symbols"`
}

// LoadMarkets загружает все спотовые рынки MEXC
func LoadMarkets() map[string]common.Market {
	resp, err := http.Get("https://api.mexc.com/api/v3/exchangeInfo")
	if err != nil {
		log.Fatalf("HTTP error: %v", err)
	}
	defer resp.Body.Close()

	var api mexcResponse
	if err := json.NewDecoder(resp.Body).Decode(&api); err != nil {
		log.Fatalf("JSON decode error: %v", err)
	}

	markets := make(map[string]common.Market)

	for _, s := range api.Symbols {
		// Фильтруем только активные спотовые пары
		if s.Status != "1" || !hasSpotPermission(s.Permissions) {
			continue
		}

		// Используем fallback через BaseAssetPrecision
		minQty := math.Pow10(-s.BaseAssetPrecision)
		stepSize := minQty

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

// проверка, есть ли permission "SPOT"
func hasSpotPermission(perms []string) bool {
	for _, p := range perms {
		if p == "SPOT" {
			return true
		}
	}
	return false
}
