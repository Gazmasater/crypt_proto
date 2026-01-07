package kucoin

import (
	"crypt_proto/cmd/exchange/common"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// Структуры API KuCoin
type kuCoinSymbol struct {
	Symbol         string `json:"symbol"`
	BaseCurrency   string `json:"baseCurrency"`
	QuoteCurrency  string `json:"quoteCurrency"`
	EnableTrading  bool   `json:"enableTrading"`
	BaseMinSize    string `json:"baseMinSize"`
	QuoteMinSize   string `json:"quoteMinSize"`
	BaseIncrement  string `json:"baseIncrement"`
	QuoteIncrement string `json:"quoteIncrement"`
	PriceIncrement string `json:"priceIncrement"`
}

type kuCoinResponse struct {
	Code string         `json:"code"`
	Data []kuCoinSymbol `json:"data"`
}

// LoadMarkets загружает все рынки KuCoin и возвращает map[string]common.Market
func LoadMarkets() map[string]common.Market {
	resp, err := http.Get("https://api.kucoin.com/api/v2/symbols")
	if err != nil {
		log.Fatalf("http error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		log.Fatalf("bad status %d: %s", resp.StatusCode, body)
	}

	var api kuCoinResponse
	if err := json.NewDecoder(resp.Body).Decode(&api); err != nil {
		log.Fatalf("decode error: %v", err)
	}

	markets := make(map[string]common.Market)
	for _, s := range api.Data {
		if !s.EnableTrading || s.BaseCurrency == "" || s.QuoteCurrency == "" {
			continue
		}

		// Преобразуем строки в float64
		bMin := parseFloat(s.BaseMinSize)
		qMin := parseFloat(s.QuoteMinSize)
		bInc := parseFloat(s.BaseIncrement)
		qInc := parseFloat(s.QuoteIncrement)
		pInc := parseFloat(s.PriceIncrement)

		key := s.BaseCurrency + "_" + s.QuoteCurrency
		markets[key] = common.Market{
			Symbol:         s.Symbol,
			Base:           s.BaseCurrency,
			Quote:          s.QuoteCurrency,
			EnableTrading:  s.EnableTrading,
			BaseMinSize:    bMin,
			QuoteMinSize:   qMin,
			BaseIncrement:  bInc,
			QuoteIncrement: qInc,
			PriceIncrement: pInc,
		}
	}

	return markets
}

func parseFloat(s string) float64 {
	var f float64
	if s == "" {
		return 0
	}
	_, err := fmt.Sscan(s, &f)
	if err != nil {
		return 0
	}
	return f
}
