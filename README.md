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
}

func NewTrader(apiKey, apiSecret string, debug bool) *Trader {
	return &Trader{
		apiKey:    strings.TrimSpace(apiKey),
		apiSecret: strings.TrimSpace(apiSecret),
		debug:     debug,
		baseURL:   "https://api.mexc.com",
		client:    &http.Client{Timeout: 10 * time.Second},
	}
}

// PlaceMarketOrder:
// ВАЖНО для MEXC: параметры отправляем В QUERY STRING, тело запроса пустое.
// Пример из доков/практики: POST /api/v3/order?symbol=...&...&signature=...
//
// - BUY: передавай quoteOrderQty (сколько QUOTE потратить, например USDT)
// - SELL: передавай quantity (сколько BASE продаём)
func (t *Trader) PlaceMarketOrder(ctx context.Context, symbol, side string, quantity, quoteOrderQty float64) (string, error) {
	symbol = strings.TrimSpace(symbol)
	side = strings.ToUpper(strings.TrimSpace(side))
	if symbol == "" {
		return "", fmt.Errorf("empty symbol")
	}
	if side != "BUY" && side != "SELL" {
		return "", fmt.Errorf("bad side=%s", side)
	}

	// собираем params
	params := url.Values{}
	params.Set("symbol", symbol)
	params.Set("side", side)
	params.Set("type", "MARKET")
	params.Set("timestamp", strconv.FormatInt(time.Now().UnixMilli(), 10))

	if side == "BUY" {
		if quoteOrderQty <= 0 {
			return "", fmt.Errorf("BUY requires quoteOrderQty>0")
		}
		// Оставь 6-8 знаков — чаще проходит.
		params.Set("quoteOrderQty", formatFloat(quoteOrderQty, 8))
	} else {
		if quantity <= 0 {
			return "", fmt.Errorf("SELL requires quantity>0")
		}
		// Для quantity обычно нужно меньше знаков, иначе "quantity scale invalid".
		// Точное число знаков зависит от символа (нужно exchangeInfo), но начнём с 8.
		params.Set("quantity", formatFloat(quantity, 8))
	}

	// signature считается по query string (без signature), затем signature добавляется
	queryToSign := params.Encode()
	params.Set("signature", t.sign(queryToSign))

	reqURL := t.baseURL + "/api/v3/order" + "?" + params.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, nil) // <— BODY NIL!
	if err != nil {
		return "", err
	}
	req.Header.Set("X-MEXC-APIKEY", t.apiKey)
	// Content-Type не обязателен, потому что тела нет.
	// Но можно явно:
	// req.Header.Set("Content-Type", "application/json")

	resp, err := t.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	b, _ := io.ReadAll(resp.Body)
	if resp.StatusCode/100 != 2 {
		return "", fmt.Errorf("mexc order error: status=%d body=%s", resp.StatusCode, string(b))
	}

	// попробуем достать orderId
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

		// free может быть строкой
		if s, ok := m["free"].(string); ok {
			v, _ := strconv.ParseFloat(s, 64)
			return v, nil
		}
		// или числом
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

func formatFloat(v float64, decimals int) string {
	if decimals < 0 {
		decimals = 0
	}
	// FormatFloat с 'f' не даёт экспоненты.
	return strconv.FormatFloat(v, 'f', decimals, 64)
}
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
}

func NewTrader(apiKey, apiSecret string, debug bool) *Trader {
	return &Trader{
		apiKey:    strings.TrimSpace(apiKey),
		apiSecret: strings.TrimSpace(apiSecret),
		debug:     debug,
		baseURL:   "https://api.mexc.com",
		client:    &http.Client{Timeout: 10 * time.Second},
	}
}

// PlaceMarketOrder:
// ВАЖНО для MEXC: параметры отправляем В QUERY STRING, тело запроса пустое.
// Пример из доков/практики: POST /api/v3/order?symbol=...&...&signature=...
//
// - BUY: передавай quoteOrderQty (сколько QUOTE потратить, например USDT)
// - SELL: передавай quantity (сколько BASE продаём)
func (t *Trader) PlaceMarketOrder(ctx context.Context, symbol, side string, quantity, quoteOrderQty float64) (string, error) {
	symbol = strings.TrimSpace(symbol)
	side = strings.ToUpper(strings.TrimSpace(side))
	if symbol == "" {
		return "", fmt.Errorf("empty symbol")
	}
	if side != "BUY" && side != "SELL" {
		return "", fmt.Errorf("bad side=%s", side)
	}

	// собираем params
	params := url.Values{}
	params.Set("symbol", symbol)
	params.Set("side", side)
	params.Set("type", "MARKET")
	params.Set("timestamp", strconv.FormatInt(time.Now().UnixMilli(), 10))

	if side == "BUY" {
		if quoteOrderQty <= 0 {
			return "", fmt.Errorf("BUY requires quoteOrderQty>0")
		}
		// Оставь 6-8 знаков — чаще проходит.
		params.Set("quoteOrderQty", formatFloat(quoteOrderQty, 8))
	} else {
		if quantity <= 0 {
			return "", fmt.Errorf("SELL requires quantity>0")
		}
		// Для quantity обычно нужно меньше знаков, иначе "quantity scale invalid".
		// Точное число знаков зависит от символа (нужно exchangeInfo), но начнём с 8.
		params.Set("quantity", formatFloat(quantity, 8))
	}

	// signature считается по query string (без signature), затем signature добавляется
	queryToSign := params.Encode()
	params.Set("signature", t.sign(queryToSign))

	reqURL := t.baseURL + "/api/v3/order" + "?" + params.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, nil) // <— BODY NIL!
	if err != nil {
		return "", err
	}
	req.Header.Set("X-MEXC-APIKEY", t.apiKey)
	// Content-Type не обязателен, потому что тела нет.
	// Но можно явно:
	// req.Header.Set("Content-Type", "application/json")

	resp, err := t.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	b, _ := io.ReadAll(resp.Body)
	if resp.StatusCode/100 != 2 {
		return "", fmt.Errorf("mexc order error: status=%d body=%s", resp.StatusCode, string(b))
	}

	// попробуем достать orderId
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

		// free может быть строкой
		if s, ok := m["free"].(string); ok {
			v, _ := strconv.ParseFloat(s, 64)
			return v, nil
		}
		// или числом
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

func formatFloat(v float64, decimals int) string {
	if decimals < 0 {
		decimals = 0
	}
	// FormatFloat с 'f' не даёт экспоненты.
	return strconv.FormatFloat(v, 'f', decimals, 64)
}





