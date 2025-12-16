package mexc

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Trader struct {
	apiKey    string
	apiSecret string
	debug     bool
	baseURL   string
	client    *http.Client

	rules map[string]SymbolRules
}

func NewTrader(apiKey, apiSecret string, debug bool) *Trader {
	return &Trader{
		apiKey:    strings.TrimSpace(apiKey),
		apiSecret: strings.TrimSpace(apiSecret),
		debug:     debug,
		baseURL:   "https://api.mexc.com",
		client:    &http.Client{Timeout: 10 * time.Second},
		rules:     map[string]SymbolRules{},
	}
}

func (t *Trader) SetRules(r map[string]SymbolRules) {
	if r == nil {
		t.rules = map[string]SymbolRules{}
		return
	}
	t.rules = r
}

func (t *Trader) Rules(symbol string) (SymbolRules, bool) {
	r, ok := t.rules[strings.TrimSpace(symbol)]
	return r, ok
}

// SmartMarketBuyUSDT:
// Покупка по USDT (quote), нормализует amount/qty по правилам symbol.
func (t *Trader) SmartMarketBuyUSDT(ctx context.Context, symbol string, usdt float64, ask float64) (string, error) {
	symbol = strings.TrimSpace(symbol)
	if usdt <= 0 {
		return "", fmt.Errorf("usdt<=0")
	}
	r, ok := t.Rules(symbol)
	if !ok {
		return "", fmt.Errorf("no rules for symbol=%s", symbol)
	}
	if !r.IsSpotTradingAllowed {
		return "", fmt.Errorf("symbol not allowed for spot/api: %s", symbol)
	}

	// если биржа разрешает quoteOrderQty — используем amount (и режем precision)
	if r.QuoteOrderQtyMarketAllowed {
		dec := r.QuoteAssetPrecision
		if dec <= 0 {
			// fallback
			dec = 2
		}
		amount := truncToDecimals(usdt, dec)
		if amount <= 0 {
			return "", fmt.Errorf("amount<=0 after trunc (dec=%d)", dec)
		}
		return t.placeMarket(ctx, symbol, "BUY", 0, amount)
	}

	// иначе — покупаем quantity по ask и режем step
	if ask <= 0 {
		return "", fmt.Errorf("ask<=0 for %s", symbol)
	}

	qtyRaw := usdt / ask
	qty := qtyRaw
	if r.BaseStep > 0 {
		qty = floorToStep(qtyRaw, r.BaseStep)
	} else if r.QtyDecimals >= 0 {
		qty = truncToDecimals(qtyRaw, r.QtyDecimals)
	}

	if r.MinQty > 0 && qty < r.MinQty {
		return "", fmt.Errorf("qty<minQty (qty=%.12f minQty=%.12f)", qty, r.MinQty)
	}
	if qty <= 0 {
		return "", fmt.Errorf("qty<=0 after normalize (raw=%.12f)", qtyRaw)
	}

	return t.placeMarket(ctx, symbol, "BUY", qty, 0)
}

// SmartMarketSellQty:
// Продажа base quantity, нормализует по step/precision.
func (t *Trader) SmartMarketSellQty(ctx context.Context, symbol string, qtyRaw float64) (string, error) {
	symbol = strings.TrimSpace(symbol)
	if qtyRaw <= 0 {
		return "", fmt.Errorf("qty<=0")
	}
	r, ok := t.Rules(symbol)
	if !ok {
		return "", fmt.Errorf("no rules for symbol=%s", symbol)
	}
	if !r.IsSpotTradingAllowed {
		return "", fmt.Errorf("symbol not allowed for spot/api: %s", symbol)
	}

	qty := qtyRaw
	if r.BaseStep > 0 {
		qty = floorToStep(qtyRaw, r.BaseStep)
	} else if r.QtyDecimals >= 0 {
		qty = truncToDecimals(qtyRaw, r.QtyDecimals)
	}

	if r.MinQty > 0 && qty < r.MinQty {
		return "", fmt.Errorf("qty<minQty (qty=%.12f minQty=%.12f)", qty, r.MinQty)
	}
	if qty <= 0 {
		return "", fmt.Errorf("qty<=0 after normalize (raw=%.12f)", qtyRaw)
	}

	return t.placeMarket(ctx, symbol, "SELL", qty, 0)
}

func (t *Trader) placeMarket(ctx context.Context, symbol, side string, quantity, quoteOrderQty float64) (string, error) {
	side = strings.ToUpper(strings.TrimSpace(side))

	params := url.Values{}
	params.Set("symbol", symbol)
	params.Set("side", side)
	params.Set("type", "MARKET")
	params.Set("timestamp", strconv.FormatInt(time.Now().UnixMilli(), 10))

	if side == "BUY" {
		if quoteOrderQty > 0 {
			params.Set("quoteOrderQty", stripZeros(quoteOrderQty))
		} else {
			params.Set("quantity", stripZeros(quantity))
		}
	} else {
		params.Set("quantity", stripZeros(quantity))
	}

	queryToSign := params.Encode()
	params.Set("signature", t.sign(queryToSign))

	reqURL := t.baseURL + "/api/v3/order" + "?" + params.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("X-MEXC-APIKEY", t.apiKey)

	resp, err := t.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	b, _ := io.ReadAll(resp.Body)
	if resp.StatusCode/100 != 2 {
		return "", fmt.Errorf("mexc order error: status=%d body=%s", resp.StatusCode, string(b))
	}

	var m map[string]any
	_ = json.Unmarshal(b, &m)
	if v, ok := m["orderId"]; ok {
		return fmt.Sprintf("%v", v), nil
	}
	if v, ok := m["orderIdStr"]; ok {
		return fmt.Sprintf("%v", v), nil
	}
	return "", nil
}

func (t *Trader) GetBalance(ctx context.Context, asset string) (float64, error) {
	asset = strings.ToUpper(strings.TrimSpace(asset))
	if asset == "" {
		return 0, fmt.Errorf("empty asset")
	}

	params := url.Values{}
	params.Set("timestamp", strconv.FormatInt(time.Now().UnixMilli(), 10))

	queryToSign := params.Encode()
	params.Set("signature", t.sign(queryToSign))

	reqURL := t.baseURL + "/api/v3/account" + "?" + params.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("X-MEXC-APIKEY", t.apiKey)

	resp, err := t.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	b, _ := io.ReadAll(resp.Body)
	if resp.StatusCode/100 != 2 {
		return 0, fmt.Errorf("mexc account error: status=%d body=%s", resp.StatusCode, string(b))
	}

	var root map[string]any
	if err := json.Unmarshal(b, &root); err != nil {
		return 0, err
	}

	balAny, _ := root["balances"].([]any)
	for _, it := range balAny {
		m, ok := it.(map[string]any)
		if !ok {
			continue
		}
		a, _ := m["asset"].(string)
		if strings.ToUpper(strings.TrimSpace(a)) != asset {
			continue
		}

		if s, ok := m["free"].(string); ok {
			v, _ := strconv.ParseFloat(s, 64)
			return v, nil
		}
		if f, ok := m["free"].(float64); ok {
			return f, nil
		}
		return 0, nil
	}

	return 0, nil
}

func (t *Trader) sign(query string) string {
	mac := hmac.New(sha256.New, []byte(t.apiSecret))
	_, _ = mac.Write([]byte(query))
	return hex.EncodeToString(mac.Sum(nil))
}

func stripZeros(v float64) string {
	// много знаков не надо, всё равно выше мы нормализуем
	s := strconv.FormatFloat(v, 'f', 12, 64)
	s = strings.TrimRight(s, "0")
	s = strings.TrimRight(s, ".")
	if s == "" {
		return "0"
	}
	return s
}
