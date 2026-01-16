Название API
9623527002

696935c42a6dcd00013273f2
b348b686-55ff-4290-897b-02d55f815f65




apikey = "4333ed4b-cd83-49f5-97d1-c399e2349748"
secretkey = "E3848531135EDB4CCFDA0F1BC14CD274"
IP = ""
Название API-ключа = "Arb"
Доступы = "Чтение"



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



package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
)

/* ================= CONFIG ================= */
const (
	apiKey        = "YOUR_API_KEY"
	apiSecret     = "YOUR_API_SECRET"
	apiPassphrase = "YOUR_API_PASSPHRASE"

	baseURL   = "https://api.kucoin.com"
	startUSDT = 12.0

	sym1 = "DASH-USDT"
	sym2 = "DASH-BTC"
	sym3 = "BTC-USDT"
)

/* ================= AUTH ================= */
func sign(ts, method, path, body string) string {
	mac := hmac.New(sha256.New, []byte(apiSecret))
	mac.Write([]byte(ts + method + path + body))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func passphrase() string {
	mac := hmac.New(sha256.New, []byte(apiSecret))
	mac.Write([]byte(apiPassphrase))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func headers(method, path, body string) http.Header {
	ts := strconv.FormatInt(time.Now().UnixMilli(), 10)
	sig := sign(ts, method, path, body)
	h := http.Header{}
	h.Set("KC-API-KEY", apiKey)
	h.Set("KC-API-SIGN", sig)
	h.Set("KC-API-TIMESTAMP", ts)
	h.Set("KC-API-PASSPHRASE", passphrase())
	h.Set("KC-API-KEY-VERSION", "2")
	h.Set("Content-Type", "application/json")
	return h
}

/* ================= UTILS ================= */
// округление вниз до ближайшего шага
func roundDown(value, step float64) float64 {
	return math.Floor(value/step) * step
}

/* ================= PLACE MARKET ================= */
func placeMarket(symbol, side string, value float64) (filled float64, err error) {
	clientOid := uuid.NewString()
	body := map[string]string{
		"symbol":    symbol,
		"type":      "market",
		"side":      side,
		"clientOid": clientOid,
	}

	if side == "buy" {
		body["funds"] = fmt.Sprintf("%.8f", value)
	} else {
		body["size"] = fmt.Sprintf("%.8f", value)
	}

	rawBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", baseURL+"/api/v1/orders", bytes.NewReader(rawBody))
	req.Header = headers("POST", "/api/v1/orders", string(rawBody))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	respBytes, _ := io.ReadAll(resp.Body)

	var r struct {
		Code string `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			DealFunds string `json:"dealFunds"`
			DealSize  string `json:"dealSize"`
			OrderId   string `json:"orderId"`
		} `json:"data"`
	}

	if err := json.Unmarshal(respBytes, &r); err != nil {
		return 0, fmt.Errorf("decode error: %v; raw: %s", err, string(respBytes))
	}

	if r.Code != "200000" {
		return 0, fmt.Errorf("order rejected: %s; raw: %s", r.Msg, string(respBytes))
	}

	if side == "buy" {
		filled, _ = strconv.ParseFloat(r.Data.DealSize, 64)
	} else {
		filled, _ = strconv.ParseFloat(r.Data.DealFunds, 64)
	}

	log.Printf("[OK] %s %s orderId=%s filled=%.8f", side, symbol, r.Data.OrderId, filled)
	return filled, nil
}

/* ================= MAIN ================= */
func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.Printf("START TRIANGLE %.2f USDT", startUSDT)

	// ------------------- STEP 1 -------------------
	leg1Size, err := placeMarket(sym1, "buy", startUSDT)
	if err != nil {
		log.Println("[FAIL] LEG1 USDT→DASH", err)
		return
	}

	// получаем шаг для DASH-BTC (например 0.001)
	stepLeg2 := 0.001
	leg2Size := roundDown(leg1Size, stepLeg2)

	// ------------------- STEP 2 -------------------
	leg2Funds, err := placeMarket(sym2, "sell", leg2Size)
	if err != nil {
		log.Println("[FAIL] LEG2 DASH→BTC", err)
		return
	}

	// получаем шаг для BTC-USDT (например 0.0001)
	stepLeg3 := 0.0001
	leg3Size := roundDown(leg2Funds, stepLeg3)

	// ------------------- STEP 3 -------------------
	leg3Funds, err := placeMarket(sym3, "sell", leg3Size)
	if err != nil {
		log.Println("[FAIL] LEG3 BTC→USDT", err)
		return
	}

	// ------------------- RESULT -------------------
	profit := leg3Funds - startUSDT
	pct := profit / startUSDT * 100
	log.Println("====== RESULT ======")
	log.Printf("START: %.4f USDT", startUSDT)
	log.Printf("END:   %.4f USDT", leg3Funds)
	log.Printf("PNL:   %.6f USDT (%.4f%%)", profit, pct)
}


gaz358@gaz358-BOD-WXX9:~/myprog/crypt_proto/test$ go run .
2026/01/16 04:37:55.162871 START TRIANGLE 12.00 USDT
2026/01/16 04:37:55.871028 [OK] buy DASH-USDT orderId=696996738333910007e03fec filled=0.00000000
2026/01/16 04:37:56.177800 [FAIL] LEG2 DASH→BTC order rejected: The quantity is invalid.; raw: {"msg":"The quantity is invalid.","code":"300000"}


