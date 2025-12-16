package mexc

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type SymbolRules struct {
	Symbol string

	IsSpotTradingAllowed       bool
	QuoteOrderQtyMarketAllowed bool

	// baseSizePrecision приходит строкой "0.0001" — это step для quantity
	BaseStepStr string
	BaseStep    float64
	QtyDecimals int

	// точность для quoteOrderQty (amount)
	QuoteAssetPrecision int
	QuotePrecision      int

	// минималки (если есть)
	MinQty         float64
	MinNotional    float64
	MinOrderAmount float64 // quoteAmountPrecision (по сути min amount), если пригодится
}

type exchangeInfoResp struct {
	Symbols []struct {
		Symbol string `json:"symbol"`

		Status string `json:"status"`

		IsSpotTradingAllowed       bool `json:"isSpotTradingAllowed"`
		QuoteOrderQtyMarketAllowed bool `json:"quoteOrderQtyMarketAllowed"`

		BaseSizePrecision string `json:"baseSizePrecision"`

		BaseAssetPrecision  int `json:"baseAssetPrecision"`
		QuoteAssetPrecision int `json:"quoteAssetPrecision"`
		QuotePrecision      int `json:"quotePrecision"`

		// Иногда присутствуют (зависит от версии ответа)
		MinQty               string `json:"minQty"`
		MinNotional          string `json:"minNotional"`
		QuoteAmountPrecision string `json:"quoteAmountPrecision"`
	} `json:"symbols"`
}

func LoadSymbolRules(ctx context.Context, baseURL string, client *http.Client) (map[string]SymbolRules, error) {
	if baseURL == "" {
		baseURL = "https://api.mexc.com"
	}
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+"/api/v3/exchangeInfo", nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, _ := io.ReadAll(resp.Body)
	if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("exchangeInfo error: status=%d body=%s", resp.StatusCode, string(b))
	}

	var root exchangeInfoResp
	if err := json.Unmarshal(b, &root); err != nil {
		return nil, err
	}

	out := make(map[string]SymbolRules, len(root.Symbols))
	for _, s := range root.Symbols {
		sym := strings.TrimSpace(s.Symbol)
		if sym == "" {
			continue
		}

		step := parseStep(s.BaseSizePrecision)
		dec := decimalsFromStepStr(s.BaseSizePrecision)

		r := SymbolRules{
			Symbol:                     sym,
			IsSpotTradingAllowed:       s.IsSpotTradingAllowed,
			QuoteOrderQtyMarketAllowed: s.QuoteOrderQtyMarketAllowed,
			BaseStepStr:                s.BaseSizePrecision,
			BaseStep:                   step,
			QtyDecimals:                dec,
			QuoteAssetPrecision:        s.QuoteAssetPrecision,
			QuotePrecision:             s.QuotePrecision,
			MinQty:                     parseFloatSafe(s.MinQty),
			MinNotional:                parseFloatSafe(s.MinNotional),
			MinOrderAmount:             parseFloatSafe(s.QuoteAmountPrecision),
		}

		out[sym] = r
	}

	return out, nil
}

func parseFloatSafe(s string) float64 {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return v
}

func parseStep(stepStr string) float64 {
	stepStr = strings.TrimSpace(stepStr)
	if stepStr == "" {
		return 0
	}
	v, err := strconv.ParseFloat(stepStr, 64)
	if err != nil {
		return 0
	}
	return v
}

func decimalsFromStepStr(step string) int {
	step = strings.TrimSpace(step)
	if step == "" || step == "1" {
		return 0
	}
	if i := strings.IndexByte(step, '.'); i >= 0 {
		frac := step[i+1:]
		frac = strings.TrimRight(frac, "0")
		return len(frac)
	}
	return 0
}

func floorToStep(x, step float64) float64 {
	if step <= 0 {
		return x
	}
	return math.Floor(x/step) * step
}

func truncToDecimals(x float64, decimals int) float64 {
	if decimals <= 0 {
		return math.Floor(x)
	}
	p := math.Pow10(decimals)
	return math.Floor(x*p) / p
}
