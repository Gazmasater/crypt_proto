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
// - BUY: отправляем quoteOrderQty (потратить ровно USDT)
// - SELL: отправляем quantity (сколько base продаём)
//
// ВАЖНО: MEXC private POST лучше делать с form-urlencoded BODY,
// иначе можно ловить 700013 Invalid content Type.
func (t *Trader) PlaceMarketOrder(ctx context.Context, symbol, side string, quantity, quoteOrderQty float64) (string, error) {
	side = strings.ToUpper(strings.TrimSpace(side))
	symbol = strings.TrimSpace(symbol)

	if side != "BUY" && side != "SELL" {
		return "", fmt.Errorf("bad side=%s", side)
	}
	if symbol == "" {
		return "", fmt.Errorf("empty symbol")
	}

	params := url.Values{}
	params.Set("symbol", symbol)
	params.Set("side", side)
	params.Set("type", "MARKET")

	// только одно из двух
	if side == "BUY" {
		if quoteOrderQty <= 0 {
			return "", fmt.Errorf("BUY requires quoteOrderQty>0")
		}
		// quoteOrderQty = сколько потратить quote-валюты (обычно USDT)
		params.Set("quoteOrderQty", fmt.Sprintf("%.8f", quoteOrderQty))
	} else {
		if quantity <= 0 {
			return "", fmt.Errorf("SELL requires quantity>0")
		}
		params.Set("quantity", fmt.Sprintf("%.12f", quantity))
	}

	// подписываем и отправляем
	bodyBytes, code, err := t.doPrivatePOST(ctx, "/api/v3/order", params)
	if err != nil {
		return "", err
	}
	if code/100 != 2 {
		return "", fmt.Errorf("mexc order error: status=%d body=%s", code, string(bodyBytes))
	}

	// orderId достаем максимально мягко
	var m map[string]any
	_ = json.Unmarshal(bodyBytes, &m)
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
	bodyBytes, code, err := t.doPrivateGET(ctx, "/api/v3/account", params)
	if err != nil {
		return 0, err
	}
	if code/100 != 2 {
		return 0, fmt.Errorf("mexc account error: status=%d body=%s", code, string(bodyBytes))
	}

	var root map[string]any
	if err := json.Unmarshal(bodyBytes, &root); err != nil {
		return 0, err
	}

	bals, _ := root["balances"].([]any)
	for _, it := range bals {
		m, ok := it.(map[string]any)
		if !ok {
			continue
		}
		a, _ := m["asset"].(string)
		if strings.ToUpper(strings.TrimSpace(a)) != asset {
			continue
		}

		switch v := m["free"].(type) {
		case string:
			if v == "" {
				return 0, nil
			}
			f, _ := strconv.ParseFloat(v, 64)
			return f, nil
		case float64:
			return v, nil
		default:
			return 0, nil
		}
	}
	return 0, nil
}

// -------- low-level signed requests --------

func (t *Trader) doPrivateGET(ctx context.Context, path string, params url.Values) ([]byte, int, error) {
	ts := strconv.FormatInt(time.Now().UnixMilli(), 10)
	params.Set("timestamp", ts)

	qs := params.Encode()
	params.Set("signature", t.sign(qs))

	u := t.baseURL + path + "?" + params.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("X-MEXC-APIKEY", t.apiKey)

	resp, err := t.client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	b, _ := io.ReadAll(resp.Body)
	return b, resp.StatusCode, nil
}

func (t *Trader) doPrivatePOST(ctx context.Context, path string, params url.Values) ([]byte, int, error) {
	ts := strconv.FormatInt(time.Now().UnixMilli(), 10)
	params.Set("timestamp", ts)

	qs := params.Encode()
	params.Set("signature", t.sign(qs))

	// POST: параметры в BODY (form-urlencoded)
	body := strings.NewReader(params.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, t.baseURL+path, body)
	if err != nil {
		return nil, 0, err
	}

	req.Header.Set("X-MEXC-APIKEY", t.apiKey)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := t.client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	b, _ := io.ReadAll(resp.Body)
	return b, resp.StatusCode, nil
}

func (t *Trader) sign(query string) string {
	mac := hmac.New(sha256.New, []byte(t.apiSecret))
	_, _ = mac.Write([]byte(query))
	return hex.EncodeToString(mac.Sum(nil))
}



  NAKAUSDT (NAKA/USDT): bid=0.0787400000 ask=0.0789200000  spread=0.0001800000 (0.22834%)  bidQty=64.9100 askQty=473.8900
  NAKAUSDC (NAKA/USDC): bid=0.0791000000 ask=0.0792000000  spread=0.0001000000 (0.12634%)  bidQty=64.6800 askQty=36.6700
  USDCUSDT (USDC/USDT): bid=1.0000000000 ask=1.0001000000  spread=0.0001000000 (0.01000%)  bidQty=54974.8400 askQty=310476.4100

  [REAL EXEC] start=2.000000 USDT triangle=USDT→NAKA→USDC→USDT
    [REAL EXEC] leg 1: BUY NAKAUSDT quoteOrderQty=2.000000
    [REAL EXEC] leg 1 ERROR: mexc order error: status=400 body={"code":700013,"msg":"Invalid content Type."}
^C2025/12/16 07:03:20.213809 shutting down...



