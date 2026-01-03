package okx

import (
	"crypt_proto/cmd/exchange/common"
	"encoding/json"
	"log"
	"math"
	"net/http"
	"strconv"
)

type okxSymbol struct {
	InstID    string `json:"instId"`
	BaseCcy   string `json:"baseCcy"`
	QuoteCcy  string `json:"quoteCcy"`
	InstType  string `json:"instType"`
	State     string `json:"state"`
	MinSz     string `json:"minSz"`
	TickSz    string `json:"tickSz"`
	BasePrec  int    `json:"baseCcyPrecision"`  // иногда отсутствует
	QuotePrec int    `json:"quoteCcyPrecision"` // иногда отсутствует
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
