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
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Trader struct {
	apiKey    string
	apiSecret string
	debug     bool
	baseURL   string
	client    *http.Client

	mu    sync.RWMutex
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

func (t *Trader) logf(format string, args ...any) {
	if !t.debug {
		return
	}
	fmt.Printf(time.Now().Format("2006-01-02 15:04:05.000 ")+"[MEXC TRADER] "+format+"\n", args...)
}

func (t *Trader) SetRules(r map[string]SymbolRules) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if r == nil {
		t.rules = map[string]SymbolRules{}
		return
	}
	t.rules = r
}

func (t *Trader) Rules(symbol string) (SymbolRules, bool) {
	symbol = strings.TrimSpace(symbol)
	t.mu.RLock()
	defer t.mu.RUnlock()
	r, ok := t.rules[symbol]
	return r, ok
}

func (t *Trader) ensureRules(ctx context.Context, symbol string) (SymbolRules, bool, error) {
	symbol = strings.TrimSpace(symbol)

	if r, ok := t.Rules(symbol); ok {
		return r, true, nil
	}

	t.mu.RLock()
	cnt := len(t.rules)
	sample := sampleKeys(t.rules, 10)
	t.mu.RUnlock()
	t.logf("rules miss: symbol=%s rules_count=%d sample_keys=%v", symbol, cnt, sample)

	r, err := t.fetchRulesForSymbol(ctx, symbol)
	if err != nil {
		t.logf("fetchRulesForSymbol err: symbol=%s err=%v", symbol, err)
		return SymbolRules{}, false, err
	}

	t.mu.Lock()
	t.rules[symbol] = r
	t.mu.Unlock()

	return r, true, nil
}

func sampleKeys(m map[string]SymbolRules, n int) []string {
	if n <= 0 {
		return nil
	}
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	if len(keys) > n {
		keys = keys[:n]
	}
	return keys
}

// SmartMarketBuyUSDT:
// Backward-compatible wrapper around SmartMarketBuyQuote.
// Покупка по USDT (quote), предпочитает quoteOrderQty.
func (t *Trader) SmartMarketBuyUSDT(ctx context.Context, symbol string, usdt float64, ask float64) (string, error) {
	return t.SmartMarketBuyQuote(ctx, symbol, usdt, ask)
}

// SmartMarketBuyQuote:
// MARKET BUY с использованием quoteOrderQty (сумма в quote) — это то, что нужно для маленьких депозитов.
// Логика:
//  1. TRY: BUY через quoteOrderQty (округляем по quoteAmountPrecisionMarket если есть)
//  2. FALLBACK: если биржа ругается на quoteOrderQty — BUY через quantity (нужен ask)
func (t *Trader) SmartMarketBuyQuote(ctx context.Context, symbol string, quoteAmount float64, ask float64) (string, error) {
	symbol = strings.TrimSpace(symbol)
	if quoteAmount <= 0 {
		return "", fmt.Errorf("quoteAmount<=0")
	}

	r, ok, err := t.ensureRules(ctx, symbol)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", fmt.Errorf("no rules for symbol=%s", symbol)
	}
	if !r.IsSpotTradingAllowed {
		return "", fmt.Errorf("symbol not allowed for spot MARKET/api: %s", symbol)
	}

	// 1) TRY quoteOrderQty
	dec := r.QuoteMarketDecimals
	if dec < 0 {
		dec = r.QuoteAssetPrecision
		if dec <= 0 {
			dec = 2
		}
	}
	amount := truncToDecimals(quoteAmount, dec)
	if amount <= 0 {
		return "", fmt.Errorf("amount<=0 after trunc (dec=%d)", dec)
	}

	t.logf("BUY TRY by QUOTE: symbol=%s quoteAmount=%.8f amount=%.8f dec=%d (qmStep=%q)",
		symbol, quoteAmount, amount, dec, r.QuoteMarketStepStr)

	id, err := t.placeMarket(ctx, symbol, "BUY", 0, amount)
	if err == nil {
		return id, nil
	}

	// 2) FALLBACK to qty if quoteOrderQty not supported / rejected as param
	if ask <= 0 {
		return "", err
	}
	if !isQuoteOrderQtyParamError(err) {
		// если ошибка не похожа на "quoteOrderQty не поддерживается" — не делаем второй ордер автоматически
		return "", err
	}

	qtyRaw := quoteAmount / ask
	qty := qtyRaw
	if r.BaseStep > 0 {
		qty = floorToStep(qtyRaw, r.BaseStep)
	} else if r.QtyDecimals >= 0 {
		qty = truncToDecimals(qtyRaw, r.QtyDecimals)
	}
	if qty <= 0 {
		needQuote := 0.0
		if r.MinQty > 0 {
			needQuote = r.MinQty * ask
		}
		return "", fmt.Errorf(
			"quoteOrderQty rejected, and qty<=0 after normalize (raw=%.12f norm=%.12f step=%.12f minQty=%.12f ask=%.10f tradeQuote=%.6f needQuote>=%.6f): firstErr=%v",
			qtyRaw, qty, r.BaseStep, r.MinQty, ask, quoteAmount, needQuote, err,
		)
	}
	if r.MinQty > 0 && qty < r.MinQty {
		needQuote := r.MinQty * ask
		return "", fmt.Errorf(
			"quoteOrderQty rejected, and qty<minQty (raw=%.12f norm=%.12f minQty=%.12f step=%.12f ask=%.10f tradeQuote=%.6f needQuote>=%.6f): firstErr=%v",
			qtyRaw, qty, r.MinQty, r.BaseStep, ask, quoteAmount, needQuote, err,
		)
	}

	t.logf("BUY FALLBACK by QTY: symbol=%s quote=%.8f ask=%.10f qtyRaw=%.12f qty=%.12f step=%.12f minQty=%.12f firstErr=%v",
		symbol, quoteAmount, ask, qtyRaw, qty, r.BaseStep, r.MinQty, err)

	return t.placeMarket(ctx, symbol, "BUY", qty, 0)
}

func isQuoteOrderQtyParamError(err error) bool {
	if err == nil {
		return false
	}
	s := strings.ToLower(err.Error())
	// Варианты сообщений от MEXC бывают разные, поэтому матчим грубо.
	if !strings.Contains(s, "quoteorderqty") {
		return false
	}
	// типичные формулировки
	badHints := []string{
		"not support",
		"not supported",
		"illegal",
		"invalid",
		"parameter",
		"mandatory",
		"required",
	}
	for _, h := range badHints {
		if strings.Contains(s, h) {
			return true
		}
	}
	// если уже содержит quoteOrderQty — почти наверняка параметрная проблема
	return true
}

// SmartMarketSellQty:
// Продажа base quantity, нормализует по step/precision.
func (t *Trader) SmartMarketSellQty(ctx context.Context, symbol string, qtyRaw float64) (string, error) {
	symbol = strings.TrimSpace(symbol)
	if qtyRaw <= 0 {
		return "", fmt.Errorf("qty<=0")
	}

	r, ok, err := t.ensureRules(ctx, symbol)
	if err != nil {
		return "", err
	}
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

	if qty <= 0 {
		return "", fmt.Errorf("qty<=0 after normalize (raw=%.12f norm=%.12f step=%.12f minQty=%.12f)", qtyRaw, qty, r.BaseStep, r.MinQty)
	}
	if r.MinQty > 0 && qty < r.MinQty {
		return "", fmt.Errorf("qty<minQty (raw=%.12f norm=%.12f minQty=%.12f step=%.12f)", qtyRaw, qty, r.MinQty, r.BaseStep)
	}

	t.logf("SELL: symbol=%s qtyRaw=%.12f qty=%.12f step=%.12f minQty=%.12f", symbol, qtyRaw, qty, r.BaseStep, r.MinQty)
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
	t.logf("placeMarket ok but no orderId in body=%s", string(b))
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
	s := strconv.FormatFloat(v, 'f', 12, 64)
	s = strings.TrimRight(s, "0")
	s = strings.TrimRight(s, ".")
	if s == "" {
		return "0"
	}
	return s
}

// fetchRulesForSymbol вытягивает exchangeInfo по одному символу и строит SymbolRules.
func (t *Trader) fetchRulesForSymbol(ctx context.Context, symbol string) (SymbolRules, error) {
	u := t.baseURL + "/api/v3/exchangeInfo?symbol=" + url.QueryEscape(symbol)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return SymbolRules{}, err
	}

	resp, err := t.client.Do(req)
	if err != nil {
		return SymbolRules{}, err
	}
	defer resp.Body.Close()

	b, _ := io.ReadAll(resp.Body)
	if resp.StatusCode/100 != 2 {
		return SymbolRules{}, fmt.Errorf("exchangeInfo error: status=%d body=%s", resp.StatusCode, string(b))
	}

	var root map[string]any
	if err := json.Unmarshal(b, &root); err != nil {
		return SymbolRules{}, fmt.Errorf("exchangeInfo unmarshal: %w body=%s", err, string(b))
	}

	syms, _ := root["symbols"].([]any)
	if len(syms) == 0 {
		return SymbolRules{}, fmt.Errorf("exchangeInfo: no symbols in response for %s body=%s", symbol, string(b))
	}

	m, ok := syms[0].(map[string]any)
	if !ok {
		return SymbolRules{}, fmt.Errorf("exchangeInfo: invalid symbol format for %s", symbol)
	}

	r := SymbolRules{Symbol: symbol}

	// ---- Eligibility (как в rules.go) ----
	status, _ := m["status"].(string)
	st, _ := m["st"].(bool)
	var perms []string
	if arr, ok := m["permissions"].([]any); ok {
		for _, it := range arr {
			if s, ok := it.(string); ok {
				perms = append(perms, s)
			}
		}
	}
	var orderTypes []string
	if arr, ok := m["orderTypes"].([]any); ok {
		for _, it := range arr {
			if s, ok := it.(string); ok {
				orderTypes = append(orderTypes, s)
			}
		}
	}
	r.IsSpotTradingAllowed = (strings.TrimSpace(status) == "1") && !st && hasStr(perms, "SPOT") && hasStr(orderTypes, "MARKET")

	// ---- Steps / decimals ----
	if s, ok := m["baseSizePrecision"].(string); ok {
		r.BaseStepStr = s
		r.BaseStep = parseStep(s)
		r.QtyDecimals = decimalsFromStepStr(s)
	}

	// quote precisions
	if v, ok := asInt(m["quoteAssetPrecision"]); ok {
		r.QuoteAssetPrecision = v
	}
	if v, ok := asInt(m["quotePrecision"]); ok {
		r.QuotePrecision = v
	}
	if s, ok := m["quoteAmountPrecisionMarket"].(string); ok {
		r.QuoteMarketStepStr = strings.TrimSpace(s)
		r.QuoteMarketStep = parseStep(r.QuoteMarketStepStr)
		r.QuoteMarketDecimals = decimalsFromStepStr(r.QuoteMarketStepStr)
	}
	if r.QuoteMarketStepStr == "" {
		r.QuoteMarketDecimals = -1
	}

	// минималки (если есть в корне)
	if s, ok := m["minQty"].(string); ok {
		r.MinQty = parseFloatSafe(s)
	}
	if s, ok := m["minNotional"].(string); ok {
		r.MinNotional = parseFloatSafe(s)
	}
	if s, ok := m["quoteAmountPrecision"].(string); ok {
		r.MinOrderAmount = parseFloatSafe(s)
	}
	if v, ok := m["quoteOrderQtyMarketAllowed"].(bool); ok {
		r.QuoteOrderQtyMarketAllowed = v
	}

	t.logf("rules loaded: symbol=%s eligible=%v status=%q st=%v perms=%v orderTypes=%v baseStep=%q qmStep=%q",
		symbol, r.IsSpotTradingAllowed, status, st, perms, orderTypes, r.BaseStepStr, r.QuoteMarketStepStr)

	return r, nil
}

func asInt(v any) (int, bool) {
	switch x := v.(type) {
	case float64:
		return int(x), true
	case int:
		return x, true
	case int64:
		return int(x), true
	case string:
		i, err := strconv.Atoi(x)
		if err != nil {
			return 0, false
		}
		return i, true
	default:
		return 0, false
	}
}
