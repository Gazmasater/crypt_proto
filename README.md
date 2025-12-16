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
	// без log-пакета, чтобы не менять твой стиль проекта
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

// ensureRules: если rules для символа нет — пробуем подгрузить с биржи и закешировать.
func (t *Trader) ensureRules(ctx context.Context, symbol string) (SymbolRules, bool, error) {
	symbol = strings.TrimSpace(symbol)

	// быстрый путь: уже в кеше
	if r, ok := t.Rules(symbol); ok {
		return r, true, nil
	}

	// логируем состояние кеша
	t.mu.RLock()
	cnt := len(t.rules)
	sample := sampleKeys(t.rules, 10)
	t.mu.RUnlock()
	t.logf("rules miss: symbol=%s rules_count=%d sample_keys=%v", symbol, cnt, sample)

	// тянем exchangeInfo только по одному символу (быстро и безопасно)
	r, err := t.fetchRulesForSymbol(ctx, symbol)
	if err != nil {
		t.logf("fetchRulesForSymbol err: symbol=%s err=%v", symbol, err)
		return SymbolRules{}, false, err
	}

	// кладём в кеш
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

// ====== PUBLIC METHODS ======

// SmartMarketBuyUSDT:
// Покупка по USDT (quote), нормализует amount/qty по правилам symbol.
func (t *Trader) SmartMarketBuyUSDT(ctx context.Context, symbol string, usdt float64, ask float64) (string, error) {
	symbol = strings.TrimSpace(symbol)
	if usdt <= 0 {
		return "", fmt.Errorf("usdt<=0")
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

	// Если биржа разрешает quoteOrderQty — используем amount (и режем precision)
	if r.QuoteOrderQtyMarketAllowed {
		dec := r.QuoteAssetPrecision
		if dec <= 0 {
			dec = 2
		}
		amount := truncToDecimals(usdt, dec)
		if amount <= 0 {
			return "", fmt.Errorf("amount<=0 after trunc (dec=%d)", dec)
		}
		t.logf("BUY by QUOTE: symbol=%s usdt=%.8f amount=%.8f dec=%d", symbol, usdt, amount, dec)
		return t.placeMarket(ctx, symbol, "BUY", 0, amount)
	}

	// Иначе — покупаем quantity по ask и режем step
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

	t.logf("BUY by QTY: symbol=%s usdt=%.8f ask=%.10f qtyRaw=%.12f qty=%.12f step=%.12f minQty=%.12f",
		symbol, usdt, ask, qtyRaw, qty, r.BaseStep, r.MinQty)

	return t.placeMarket(ctx, symbol, "BUY", qty, 0)
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

	if r.MinQty > 0 && qty < r.MinQty {
		return "", fmt.Errorf("qty<minQty (qty=%.12f minQty=%.12f)", qty, r.MinQty)
	}
	if qty <= 0 {
		return "", fmt.Errorf("qty<=0 after normalize (raw=%.12f)", qtyRaw)
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
	// На всякий случай лог raw
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

//
// ===== RULES LOADER (exchangeInfo) =====
//

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

	// Парсим максимально гибко (через map), чтобы не зависеть от точной схемы.
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

	r := SymbolRules{}

	// spot/api разрешение: стараемся найти isSpotTradingAllowed и/или orderTypes MARKET
	if v, ok := m["isSpotTradingAllowed"].(bool); ok {
		r.IsSpotTradingAllowed = v
	}
	// если поля нет — пробуем по orderTypes
	if r.IsSpotTradingAllowed == false {
		if ots, ok := m["orderTypes"].([]any); ok {
			hasMarket := false
			for _, it := range ots {
				if s, ok := it.(string); ok && strings.ToUpper(s) == "MARKET" {
					hasMarket = true
					break
				}
			}
			// это не идеально, но лучше чем падать
			if hasMarket {
				r.IsSpotTradingAllowed = true
			}
		}
	}

	// precision
	if v, ok := asInt(m["baseAssetPrecision"]); ok {
		// у тебя это может быть не нужно, но оставим
		_ = v
	}
	if v, ok := asInt(m["quoteAssetPrecision"]); ok {
		r.QuoteAssetPrecision = v
	}

	// filters
	if flt, ok := m["filters"].([]any); ok {
		for _, it := range flt {
			fm, ok := it.(map[string]any)
			if !ok {
				continue
			}
			ft, _ := fm["filterType"].(string)
			ft = strings.ToUpper(strings.TrimSpace(ft))

			switch ft {
			case "LOT_SIZE", "MARKET_LOT_SIZE":
				// stepSize / minQty
				if s, ok := fm["stepSize"].(string); ok {
					if v, err := strconv.ParseFloat(s, 64); err == nil {
						r.BaseStep = v
					}
				}
				if s, ok := fm["minQty"].(string); ok {
					if v, err := strconv.ParseFloat(s, 64); err == nil {
						r.MinQty = v
					}
				}
			}
		}
	}

	// quoteOrderQtyMarketAllowed (если есть в ответе — отлично)
	if v, ok := m["quoteOrderQtyMarketAllowed"].(bool); ok {
		r.QuoteOrderQtyMarketAllowed = v
	}

	// Доп. лог: что загрузили
	t.logf("rules loaded: symbol=%s spotAllowed=%v quoteOrderQtyMarketAllowed=%v quotePrec=%d step=%.12f minQty=%.12f",
		symbol, r.IsSpotTradingAllowed, r.QuoteOrderQtyMarketAllowed, r.QuoteAssetPrecision, r.BaseStep, r.MinQty)

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


[ARB] -0.072%  USDT→MX→USDC→USDT  maxStart=6.3055 USDT (6.3055 USDT)  safeStart=6.3055 USDT (6.3055 USDT) (x1.00)  bottleneck=MXUSDC
  MXUSDT (MX/USDT): bid=2.1007000000 ask=2.1010000000  spread=0.0003000000 (0.01428%)  bidQty=1.2900 askQty=6.1800
  MXUSDC (MX/USDC): bid=2.1020000000 ask=2.1040000000  spread=0.0020000000 (0.09510%)  bidQty=3.0000 askQty=0.6600
  USDCUSDT (USDC/USDT): bid=1.0000000000 ask=1.0001000000  spread=0.0001000000 (0.01000%)  bidQty=140527.7200 askQty=224530.4700

  [REAL EXEC] start=2.000000 USDT triangle=USDT→MX→USDC→USDT
    [REAL EXEC] legs: sym1=MXUSDT sym2=MXUSDC sym3=USDCUSDT
    [REAL EXEC] parsed: sym1=MXUSDT (MX/USDT) sym2=MXUSDC (MX/USDC) sym3=USDCUSDT (USDC/USDT)
2025-12-17 01:40:02.315


