package kucoin

import (
	"crypt_proto/cmd/exchange/common"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

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
	MinFunds       string `json:"minFunds"`
}

type kuCoinResponse struct {
	Code string         `json:"code"`
	Data []kuCoinSymbol `json:"data"`
}

func LoadMarkets() map[string]common.Market {
	client := &http.Client{Timeout: 15 * time.Second}

	req, err := http.NewRequest(http.MethodGet, "https://api.kucoin.com/api/v2/symbols", nil)
	if err != nil {
		log.Fatalf("kucoin request build error: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("kucoin http error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Fatalf("kucoin bad status %d: %s", resp.StatusCode, body)
	}

	var api kuCoinResponse
	if err := json.NewDecoder(resp.Body).Decode(&api); err != nil {
		log.Fatalf("kucoin decode error: %v", err)
	}

	if api.Code != "200000" {
		log.Fatalf("kucoin api returned code %s", api.Code)
	}

	markets := make(map[string]common.Market, len(api.Data))

	for _, s := range api.Data {
		if !s.EnableTrading || s.BaseCurrency == "" || s.QuoteCurrency == "" {
			continue
		}

		key := s.BaseCurrency + "_" + s.QuoteCurrency
		markets[key] = common.Market{
			Symbol:         s.Symbol,
			Base:           s.BaseCurrency,
			Quote:          s.QuoteCurrency,
			EnableTrading:  s.EnableTrading,
			BaseMinSize:    parseFloat(s.BaseMinSize),
			QuoteMinSize:   parseFloat(s.QuoteMinSize),
			BaseIncrement:  parseFloat(s.BaseIncrement),
			QuoteIncrement: parseFloat(s.QuoteIncrement),
			PriceIncrement: parseFloat(s.PriceIncrement),
			MinNotional:    parseFloat(s.MinFunds),
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
