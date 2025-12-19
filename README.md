mx0vglmT3srN1IS19H
135bb7a7509e4421bad692415c53753b



sudo systemctl mask sleep.target suspend.target hibernate.target hybrid-sleep.target



wbs-api.mexc.com/ws 


[https://edis-global.vercel.app/ru/vps-hosting/singapore-singapore
](https://sg.edisglobal.com/)



git pull --rebase origin privat
git push origin privat


BOOK_INTERVAL=100ms
SYMBOLS_FILE=triangles_markets.csv
DEBUG=false


import (
    // ...
    "net/http"
    _ "net/http/pprof"
)


   // pprof HTTP-сервер
    go func() {
        log.Println("pprof on http://localhost:6060/debug/pprof/")
        if err := http.ListenAndServe("localhost:6060", nil); err != nil {
            log.Printf("pprof server error: %v", err)
        }
    }()


	go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30


(pprof) top        # показать топ функций по CPU
(pprof) top10
(pprof) list parsePBWrapperMid   # подробный разбор одной функции
(pprof) quit


go tool pprof http://localhost:6060/debug/pprof/heap


(pprof) top
(pprof) top -cum
(pprof) list parsePBWrapperMid
(pprof) quit




export TRADE_AMOUNT_USDT=100
export FEE_PCT=0.04
export SELL_SAFETY=0.995

export TRIANGLES_FILE=triangles_markets.csv
export TRIANGLES_ENRICHED_FILE=triangles_markets_enriched.csv

go run ./cmd/triangles_enrich_mexc



rules.go
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



trader.go

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
//   1) TRY: BUY через quoteOrderQty (округляем по quoteAmountPrecisionMarket если есть)
//   2) FALLBACK: если биржа ругается на quoteOrderQty — BUY через quantity (нужен ask)
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



executor_real.go

package arb

import (
	"context"
	"fmt"
	"io"
	"math"
	"strings"
	"sync"
	"time"

	"crypt_proto/domain"
)

type SpotTrader interface {
	// ВНИМАНИЕ: в текущем интерфейсе BUY поддержан только "на сумму USDT".
	// Executor ниже использует BUY через quoteQty, но если leg.From != USDT — он вернёт ошибку.
	SmartMarketBuyUSDT(ctx context.Context, symbol string, usdt float64, ask float64) (string, error)
	SmartMarketSellQty(ctx context.Context, symbol string, qty float64) (string, error)
	GetBalance(ctx context.Context, asset string) (float64, error)
}

type execReq struct {
	ctx       context.Context
	t         domain.Triangle
	quotes    map[string]domain.Quote // snapshot только нужных символов
	startUSDT float64
	triName   string
}

type RealExecutor struct {
	trader SpotTrader
	out    io.Writer

	StartUSDT  float64
	SellSafety float64

	// cooldown по треугольнику (по имени)
	Cooldown time.Duration

	mu       sync.Mutex
	lastExec map[string]time.Time

	// Очередь (строго последовательное исполнение)
	queue chan execReq
	wg    sync.WaitGroup
}

func NewRealExecutor(tr SpotTrader, out io.Writer, startUSDT float64) *RealExecutor {
	e := &RealExecutor{
		trader:     tr,
		out:        out,
		StartUSDT:  startUSDT,
		SellSafety: 0.995,
		Cooldown:   500 * time.Millisecond,
		lastExec:   make(map[string]time.Time),

		// буфер можно увеличить, но лучше небольшой, чтобы не копить “устаревшие” сделки
		queue: make(chan execReq, 16),
	}

	// worker: исполняет строго по одному
	e.wg.Add(1)
	go func() {
		defer e.wg.Done()
		for req := range e.queue {
			_ = e.executeOnce(req)
		}
	}()

	return e
}

func (e *RealExecutor) Name() string { return "REAL" }

type flusher interface{ Flush() error }

func (e *RealExecutor) logf(format string, args ...any) {
	fmt.Fprintf(e.out, format+"\n", args...)
	if f, ok := e.out.(flusher); ok {
		_ = f.Flush()
	}
}

func (e *RealExecutor) step(name string) func() {
	start := time.Now()
	e.logf("    [REAL EXEC] >>> %s", name)
	return func() {
		e.logf("    [REAL EXEC] <<< %s (%s)", name, time.Since(start).Truncate(time.Millisecond))
	}
}

// Execute теперь НЕ исполняет сразу.
// Он кладёт треугольник в очередь со снапшотом котировок и возвращает.
func (e *RealExecutor) Execute(ctx context.Context, t domain.Triangle, quotes map[string]domain.Quote, startUSDT float64) error {
	triName := strings.TrimSpace(t.Name)
	if triName == "" {
		triName = "triangle"
	}

	if startUSDT <= 0 {
		startUSDT = e.StartUSDT
	}
	if startUSDT <= 0 {
		return fmt.Errorf("startUSDT<=0 (startUSDT=%.6f, StartUSDT=%.6f)", startUSDT, e.StartUSDT)
	}

	if len(t.Legs) < 3 {
		return fmt.Errorf("triangle %s has <3 legs", triName)
	}
	sym1 := strings.TrimSpace(t.Legs[0].Symbol)
	sym2 := strings.TrimSpace(t.Legs[1].Symbol)
	sym3 := strings.TrimSpace(t.Legs[2].Symbol)
	if sym1 == "" || sym2 == "" || sym3 == "" {
		return fmt.Errorf("triangle %s has empty leg symbols: [%q, %q, %q]", triName, sym1, sym2, sym3)
	}

	// СНАПШОТ котировок только по нужным символам
	snap := make(map[string]domain.Quote, 3)
	if q, ok := quotes[sym1]; ok {
		snap[sym1] = q
	}
	if q, ok := quotes[sym2]; ok {
		snap[sym2] = q
	}
	if q, ok := quotes[sym3]; ok {
		snap[sym3] = q
	}

	req := execReq{
		ctx:       ctx,
		t:         t,
		quotes:    snap,
		startUSDT: startUSDT,
		triName:   triName,
	}

	select {
	case e.queue <- req:
		e.logf("  [REAL EXEC] QUEUED: start=%.6f USDT triangle=%s", startUSDT, triName)
		return nil
	default:
		e.logf("  [REAL EXEC] SKIP: queue full (triangle=%s)", triName)
		return nil
	}
}

func (e *RealExecutor) executeOnce(req execReq) error {
	now := time.Now()

	// cooldown по имени треугольника
	e.mu.Lock()
	if last, ok := e.lastExec[req.triName]; ok && e.Cooldown > 0 && now.Sub(last) < e.Cooldown {
		left := (e.Cooldown - now.Sub(last)).Truncate(time.Millisecond)
		e.mu.Unlock()
		e.logf("  [REAL EXEC] SKIP cooldown triangle=%s left=%s", req.triName, left)
		return nil
	}
	e.mu.Unlock()

	t := req.t
	quotes := req.quotes
	startUSDT := req.startUSDT
	triName := req.triName

	e.logf("  [REAL EXEC] start=%.6f USDT triangle=%s", startUSDT, triName)

	// Покажем ноги
	for i, leg := range t.Legs {
		e.logf("    [REAL EXEC] leg%d: sym=%s dir=%d from=%s to=%s", i+1, leg.Symbol, leg.Dir, leg.From, leg.To)
	}

	// ===== balances before =====
	var usdt0 float64
	{
		done := e.step("GetBalance USDT (before)")
		v, err := e.trader.GetBalance(req.ctx, "USDT")
		done()
		if err != nil {
			e.logf("    [REAL EXEC] BAL ERR: get USDT before: %v", err)
			return err
		}
		usdt0 = v
		e.logf("    [REAL EXEC] BAL before: USDT=%.12f", usdt0)
		if usdt0+1e-9 < startUSDT {
			return fmt.Errorf("insufficient USDT: have=%.12f need=%.12f", usdt0, startUSDT)
		}
	}

	// Текущий “поток”: с какой валютой и суммой идём по ногам.
	curAsset := "USDT"
	curAmount := startUSDT

	// Исполняем 3 ноги по dir/from/to
	for i := 0; i < 3; i++ {
		leg := t.Legs[i]
		sym := strings.TrimSpace(leg.Symbol)
		if sym == "" {
			return fmt.Errorf("leg%d: empty symbol", i+1)
		}

		q, ok := quotes[sym]
		if !ok {
			return fmt.Errorf("leg%d: no quote snapshot for %s", i+1, sym)
		}
		if q.Ask <= 0 || q.Bid <= 0 {
			return fmt.Errorf("leg%d: bad quote for %s (ask=%.10f bid=%.10f)", i+1, sym, q.Ask, q.Bid)
		}

		from := strings.ToUpper(strings.TrimSpace(leg.From))
		to := strings.ToUpper(strings.TrimSpace(leg.To))
		if from == "" || to == "" {
			return fmt.Errorf("leg%d: empty from/to (from=%q to=%q)", i+1, leg.From, leg.To)
		}

		if curAsset != from {
			// не фейлим сразу — но это почти всегда признак рассинхрона описания треугольника/исполнения
			e.logf("    [REAL EXEC] WARN leg%d: curAsset=%s curAmount=%.12f but leg.From=%s",
				i+1, curAsset, curAmount, from)
		}

		// Балансы до
		var fromBefore, toBefore float64
		{
			done := e.step(fmt.Sprintf("GetBalance %s (before leg%d)", from, i+1))
			v, err := e.trader.GetBalance(req.ctx, from)
			done()
			if err != nil {
				e.logf("    [REAL EXEC] BAL ERR: get %s before leg%d: %v", from, i+1, err)
				return err
			}
			fromBefore = v
		}
		{
			done := e.step(fmt.Sprintf("GetBalance %s (before leg%d)", to, i+1))
			v, err := e.trader.GetBalance(req.ctx, to)
			done()
			if err != nil {
				e.logf("    [REAL EXEC] BAL ERR: get %s before leg%d: %v", to, i+1, err)
				return err
			}
			toBefore = v
		}

		if leg.Dir < 0 {
			// BUY: тратим quote (from), получаем base (to)
			// В текущем интерфейсе SpotTrader BUY поддержан только в USDT.
			spend := curAmount
			if spend <= 0 {
				return fmt.Errorf("leg%d BUY: spend<=0 (%s)", i+1, from)
			}
			if fromBefore+1e-9 < spend {
				return fmt.Errorf("leg%d BUY: insufficient %s: have=%.12f need=%.12f", i+1, from, fromBefore, spend)
			}
			if from != "USDT" {
				// ВАЖНО: пока твой трейдер умеет BUY только за USDT.
				// Чтобы торговать BUY за USDC/прочее — надо расширить интерфейс на SmartMarketBuyQuote.
				return fmt.Errorf("leg%d BUY: quote asset is %s, but SpotTrader supports BUY only by USDT (need SmartMarketBuyQuote)", i+1, from)
			}

			e.logf("    [REAL EXEC] leg%d PRE: BUY %s spend=%s=%.6f ask=%.10f bid=%.10f | %s before=%.12f %s before=%.12f",
				i+1, sym, from, spend, q.Ask, q.Bid, from, fromBefore, to, toBefore)

			var ord string
			{
				orderCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

				done := e.step(fmt.Sprintf("SmartMarketBuyUSDT leg%d", i+1))
				id, err := e.trader.SmartMarketBuyUSDT(orderCtx, sym, spend, q.Ask)
				done()
				if err != nil {
					e.logf("    [REAL EXEC] leg%d PLACE ERR (BUY): %v", i+1, err)
					return err
				}
				ord = id
			}
			e.logf("    [REAL EXEC] leg%d PLACE OK: orderId=%s", i+1, ord)

			var toAfter float64
			{
				done := e.step(fmt.Sprintf("waitBalanceChange %s (after leg%d)", to, i+1))
				v, err := e.waitBalanceChange(req.ctx, to, toBefore, 3*time.Second, 150*time.Millisecond)
				done()
				if err != nil {
					e.logf("    [REAL EXEC] leg%d WAIT BAL ERR (%s): %v", i+1, to, err)
					return err
				}
				toAfter = v
			}

			delta := toAfter - toBefore
			e.logf("    [REAL EXEC] leg%d BAL after: %s=%.12f delta=%.12f", i+1, to, toAfter, delta)
			if delta <= 0 {
				return fmt.Errorf("leg%d BUY: %s did not increase (before=%.12f after=%.12f)", i+1, to, toBefore, toAfter)
			}

			curAsset = to
			curAmount = delta
			continue
		}

		// SELL: продаём base qty (from), получаем quote (to)
		qtyRaw := curAmount
		// На SELL берём баланс from, а не curAmount, потому что фактический qty после BUY может отличаться,
		// и самый надёжный путь — продать то, что есть на балансе (с safety).
		qty := fromBefore * e.SellSafety
		if qty <= 0 {
			return fmt.Errorf("leg%d SELL: qty<=0 (%s=%.12f safety=%.6f)", i+1, from, fromBefore, e.SellSafety)
		}

		e.logf("    [REAL EXEC] leg%d PRE: SELL %s qty=%s=%.12f (curAmount=%.12f raw=%.12f) bid=%.10f ask=%.10f | %s before=%.12f %s before=%.12f",
			i+1, sym, from, qty, curAmount, qtyRaw, q.Bid, q.Ask, from, fromBefore, to, toBefore)

		var ord string
		{
			orderCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			done := e.step(fmt.Sprintf("SmartMarketSellQty leg%d", i+1))
			id, err := e.trader.SmartMarketSellQty(orderCtx, sym, qty)
			done()
			if err != nil {
				e.logf("    [REAL EXEC] leg%d PLACE ERR (SELL): %v", i+1, err)
				return err
			}
			ord = id
		}
		e.logf("    [REAL EXEC] leg%d PLACE OK: orderId=%s", i+1, ord)

		var toAfter float64
		{
			done := e.step(fmt.Sprintf("waitBalanceChange %s (after leg%d)", to, i+1))
			v, err := e.waitBalanceChange(req.ctx, to, toBefore, 3*time.Second, 150*time.Millisecond)
			done()
			if err != nil {
				e.logf("    [REAL EXEC] leg%d WAIT BAL ERR (%s): %v", i+1, to, err)
				return err
			}
			toAfter = v
		}

		delta := toAfter - toBefore
		e.logf("    [REAL EXEC] leg%d BAL after: %s=%.12f delta=%.12f", i+1, to, toAfter, delta)
		if delta <= 0 {
			return fmt.Errorf("leg%d SELL: %s did not increase (before=%.12f after=%.12f)", i+1, to, toBefore, toAfter)
		}

		curAsset = to
		curAmount = delta
	}

	// Финальный баланс USDT
	var usdtAfter float64
	{
		done := e.step("GetBalance USDT (after)")
		v, err := e.trader.GetBalance(req.ctx, "USDT")
		done()
		if err != nil {
			e.logf("    [REAL EXEC] BAL ERR: get USDT after: %v", err)
			return err
		}
		usdtAfter = v
	}

	dUSDTTotal := usdtAfter - usdt0
	e.logf("    [REAL EXEC] DONE: curAsset=%s curAmount=%.12f", curAsset, curAmount)
	e.logf("    [REAL EXEC] DONE: USDT start=%.12f end=%.12f pnl(total)=%.12f (%.4f%%)",
		usdt0, usdtAfter, dUSDTTotal, pct(dUSDTTotal, startUSDT))

	e.mu.Lock()
	e.lastExec[triName] = time.Now()
	e.mu.Unlock()

	return nil
}

// waitBalanceChange ждёт, пока баланс станет отличаться от baseline.
func (e *RealExecutor) waitBalanceChange(ctx context.Context, asset string, baseline float64, timeout, interval time.Duration) (float64, error) {
	const tol = 1e-12

	deadline := time.NewTimer(timeout)
	tick := time.NewTicker(interval)
	defer deadline.Stop()
	defer tick.Stop()

	cur, err := e.trader.GetBalance(ctx, asset)
	if err == nil && math.Abs(cur-baseline) > tol {
		return cur, nil
	}

	for {
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		case <-deadline.C:
			last, err := e.trader.GetBalance(ctx, asset)
			if err != nil {
				return 0, fmt.Errorf("timeout, last balance read error for %s: %v", asset, err)
			}
			return 0, fmt.Errorf("timeout waiting %s balance change: baseline=%.12f last=%.12f", asset, baseline, last)
		case <-tick.C:
			cur, err := e.trader.GetBalance(ctx, asset)
			if err != nil {
				continue
			}
			if math.Abs(cur-baseline) > tol {
				return cur, nil
			}
		}
	}
}

// parseBaseQuote — простой парсер BASE/QUOTE по суффиксу.
func parseBaseQuote(symbol string) (base, quote string) {
	quotes := []string{"USDT", "USDC", "BTC", "ETH", "EUR", "TRY", "BRL", "RUB"}
	for _, q := range quotes {
		if strings.HasSuffix(symbol, q) && len(symbol) > len(q) {
			return symbol[:len(symbol)-len(q)], q
		}
	}
	return symbol, ""
}

func pct(delta, denom float64) float64 {
	if denom == 0 {
		return 0
	}
	return (delta / denom) * 100
}

