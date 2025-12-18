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

	// Мы используем это поле как "допущен к торговле спот MARKET через API".
	// Для MEXC оно НЕ равно json-полю isSpotTradingAllowed (которое часто false даже на BTCUSDT).
	IsSpotTradingAllowed bool

	// Историческое поле. На MEXC оно часто отсутствует/ложное,
	// поэтому в торговой логике на него НЕ опираемся (только TRY->fallback).
	QuoteOrderQtyMarketAllowed bool

	// baseSizePrecision приходит строкой "0.0001" — это step для quantity
	BaseStepStr string
	BaseStep    float64
	QtyDecimals int

	// точность для quoteOrderQty (amount)
	QuoteAssetPrecision int
	QuotePrecision      int

	// Точность для MARKET BUY через quoteOrderQty.
	// У MEXC приходит как строка (например "1" или "0.01").
	QuoteMarketStepStr  string
	QuoteMarketStep     float64
	QuoteMarketDecimals int

	// минималки (если есть)
	MinQty         float64
	MinNotional    float64
	MinOrderAmount float64 // quoteAmountPrecision (по сути min amount), если пригодится
}

type exchangeInfoResp struct {
	Symbols []struct {
		Symbol string `json:"symbol"`

		Status string `json:"status"`

		IsSpotTradingAllowed       bool     `json:"isSpotTradingAllowed"`
		QuoteOrderQtyMarketAllowed bool     `json:"quoteOrderQtyMarketAllowed"`
		OrderTypes                 []string `json:"orderTypes"`
		Permissions                []string `json:"permissions"`
		St                         bool     `json:"st"`

		BaseSizePrecision string `json:"baseSizePrecision"`

		BaseAssetPrecision  int `json:"baseAssetPrecision"`
		QuoteAssetPrecision int `json:"quoteAssetPrecision"`
		QuotePrecision      int `json:"quotePrecision"`

		QuoteAmountPrecisionMarket string `json:"quoteAmountPrecisionMarket"`

		// Иногда присутствуют (зависит от версии ответа)
		MinQty               string `json:"minQty"`
		MinNotional          string `json:"minNotional"`
		QuoteAmountPrecision string `json:"quoteAmountPrecision"`
	} `json:"symbols"`
}

func hasStr(a []string, want string) bool {
	for _, v := range a {
		if strings.EqualFold(strings.TrimSpace(v), want) {
			return true
		}
	}
	return false
}

// marketEligibleMEXC: критерий "пара подходит для торговли" по твоему требованию:
// status=="1", st==false, permissions содержит "SPOT", orderTypes содержит "MARKET".
func marketEligibleMEXC(status string, st bool, permissions, orderTypes []string) bool {
	if strings.TrimSpace(status) != "1" {
		return false
	}
	if st {
		return false
	}
	if !hasStr(permissions, "SPOT") {
		return false
	}
	if !hasStr(orderTypes, "MARKET") {
		return false
	}
	return true
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

		baseStep := parseStep(s.BaseSizePrecision)
		baseDec := decimalsFromStepStr(s.BaseSizePrecision)

		qmStepStr := strings.TrimSpace(s.QuoteAmountPrecisionMarket)
		qmStep := parseStep(qmStepStr)
		qmDec := decimalsFromStepStr(qmStepStr)
		// В ответе MEXC quoteAmountPrecisionMarket часто бывает "1" (т.е. decimals=0)
		// Если поля нет — оставим -1 и будем фолбэчить на QuoteAssetPrecision.
		if qmStepStr == "" {
			qmDec = -1
		}

		r := SymbolRules{
			Symbol: sym,

			IsSpotTradingAllowed: marketEligibleMEXC(s.Status, s.St, s.Permissions, s.OrderTypes),
			// На это поле не опираемся, но сохраним если пришло.
			QuoteOrderQtyMarketAllowed: s.QuoteOrderQtyMarketAllowed,

			BaseStepStr: s.BaseSizePrecision,
			BaseStep:    baseStep,
			QtyDecimals: baseDec,

			QuoteAssetPrecision: s.QuoteAssetPrecision,
			QuotePrecision:      s.QuotePrecision,

			QuoteMarketStepStr:  qmStepStr,
			QuoteMarketStep:     qmStep,
			QuoteMarketDecimals: qmDec,

			MinQty:         parseFloatSafe(s.MinQty),
			MinNotional:    parseFloatSafe(s.MinNotional),
			MinOrderAmount: parseFloatSafe(s.QuoteAmountPrecision),
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
