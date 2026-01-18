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





gaz358@gaz358-BOD-WXX9:~/myprog/crypt_proto/test$ go run .
2026/01/18 19:11:12.540545 START TRIANGLE 12.00 USDT
2026/01/18 19:11:15.204222 LEG1 USDT→DASH | DASH=0.286400 | total=818.177249ms
2026/01/18 19:11:16.023274 LEG2 DASH→BTC | BTC=0.00026000 | total=819.003358ms
2026/01/18 19:11:16.434442 LEG3 BTC→USDT | BTC sold=0.00026000
2026/01/18 19:11:16.843108 LEG3 BTC→USDT | USDT=34.1307 | total=819.808459ms
2026/01/18 19:11:16.843143 ====== TRIANGLE SUMMARY ======
2026/01/18 19:11:16.843148 LEG1 time: 818.177249ms
2026/01/18 19:11:16.843152 LEG2 time: 819.003358ms
2026/01/18 19:11:16.843156 LEG3 time: 819.808459ms
2026/01/18 19:11:16.843160 TOTAL time: 4.302597573s
2026/01/18 19:11:16.843164 PNL: 22.130750 USDT




gaz358@gaz358-BOD-WXX9:~/myprog/crypt_proto/step$ go run .
2026/01/18 19:28:47 Getting steps for triangle: USDT → DASH → BTC → USDT
====== TRADING STEPS ======
DASH-USDT step: 0.00010000
DASH-BTC step: 0.00010000
BTC-USDT step: 0.00001000


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
	apiKey        = "696935c42a6dcd00013273f2"
	apiSecret     = "b348b686-55ff-4290-897b-02d55f815f65"
	apiPassphrase = "Gazmaster_358"

	baseURL   = "https://api.kucoin.com"
	startUSDT = 12.0

	sym1 = "DASH-USDT"
	sym2 = "DASH-BTC"
	sym3 = "BTC-USDT"

	// Шаги ордеров для треугольника (получены через Postman или API)
	step1 = 0.0001  // DASH-USDT
	step2 = 0.0001  // DASH-BTC
	step3 = 0.00001 // BTC-USDT
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
	h := http.Header{}
	h.Set("KC-API-KEY", apiKey)
	h.Set("KC-API-SIGN", sign(ts, method, path, body))
	h.Set("KC-API-TIMESTAMP", ts)
	h.Set("KC-API-PASSPHRASE", passphrase())
	h.Set("KC-API-KEY-VERSION", "2")
	h.Set("Content-Type", "application/json")
	return h
}

/* ================= UTILS ================= */

func roundDown(v, step float64) float64 {
	return math.Floor(v/step) * step
}

/* ================= MAIN ================= */

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	totalStart := time.Now()
	log.Printf("START TRIANGLE %.2f USDT", startUSDT)

	var dash, btc, usdt float64
	var leg1Time, leg2Time, leg3Time time.Duration

	/* ================= LEG 1: USDT → DASH ================= */
	{
		legStart := time.Now()
		_, err := placeMarket(sym1, "buy", startUSDT)
		if err != nil {
			log.Fatal("LEG1 BUY failed:", err)
		}

		dash, err := getBalance("DASH")
		if err != nil || dash < step1 {
			log.Fatal("LEG1 balance failed or below step")
		}
		dash = roundDown(dash, step1)

		leg1Time = time.Since(legStart)
		log.Printf("LEG1 USDT→DASH | DASH=%.6f | total=%s", dash, leg1Time)
	}

	/* ================= LEG 2: DASH → BTC ================= */
	{
		legStart := time.Now()
		if dash < step2 {
			log.Fatalf("LEG2 skipped | DASH below min step: %.6f", dash)
		}
		dash = roundDown(dash, step2)

		_, err := placeMarket(sym2, "sell", dash)
		if err != nil {
			log.Fatal("LEG2 SELL failed:", err)
		}

		btc, err := getBalance("BTC")
		if err != nil || btc < step3 {
			log.Fatal("LEG2 balance failed or BTC below min step")
		}
		btc = roundDown(btc, step3)

		leg2Time = time.Since(legStart)
		log.Printf("LEG2 DASH→BTC | BTC=%.8f | total=%s", btc, leg2Time)
	}

	/* ================= LEG 3: BTC → USDT ================= */
	{
		legStart := time.Now()
		if btc < step3 {
			log.Printf("LEG3 skipped | BTC below min step: %.8f", btc)
		} else {
			btcToSell := roundDown(btc, step3)
			_, err := placeMarket(sym3, "sell", btcToSell)
			if err != nil {
				log.Fatal("LEG3 SELL failed:", err)
			}
			log.Printf("LEG3 BTC→USDT | BTC sold=%.8f", btcToSell)
		}

		usdt, _ = getBalance("USDT")
		leg3Time = time.Since(legStart)
		log.Printf("LEG3 BTC→USDT | USDT=%.4f | total=%s", usdt, leg3Time)
	}

	/* ================= SUMMARY ================= */
	totalTime := time.Since(totalStart)
	log.Println("====== TRIANGLE SUMMARY ======")
	log.Printf("LEG1 time: %s", leg1Time)
	log.Printf("LEG2 time: %s", leg2Time)
	log.Printf("LEG3 time: %s", leg3Time)
	log.Printf("TOTAL time: %s", totalTime)
	log.Printf("PNL: %.6f USDT", usdt-startUSDT)
}

/* ================= API FUNCTIONS ================= */

func placeMarket(symbol, side string, value float64) (string, error) {
	body := map[string]string{
		"symbol":    symbol,
		"type":      "market",
		"side":      side,
		"clientOid": uuid.NewString(),
	}
	if side == "buy" {
		body["funds"] = fmt.Sprintf("%.8f", value)
	} else {
		body["size"] = fmt.Sprintf("%.8f", value)
	}
	raw, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", baseURL+"/api/v1/orders", bytes.NewReader(raw))
	req.Header = headers("POST", "/api/v1/orders", string(raw))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var r struct {
		Code string `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			OrderId string `json:"orderId"`
		} `json:"data"`
	}
	b, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(b, &r); err != nil {
		return "", err
	}
	if r.Code != "200000" {
		return "", fmt.Errorf("order rejected: %s", r.Msg)
	}
	return r.Data.OrderId, nil
}

func getBalance(currency string) (float64, error) {
	req, _ := http.NewRequest("GET", baseURL+"/api/v1/accounts", nil)
	req.Header = headers("GET", "/api/v1/accounts", "")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var r struct {
		Code string `json:"code"`
		Data []struct {
			Currency  string `json:"currency"`
			Type      string `json:"type"`
			Available string `json:"available"`
		} `json:"data"`
	}
	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &r); err != nil {
		return 0, err
	}

	for _, acc := range r.Data {
		if acc.Currency == currency && acc.Type == "trade" {
			return strconv.ParseFloat(acc.Available, 64)
		}
	}
	return 0, fmt.Errorf("balance %s not found", currency)
}


gaz358@gaz358-BOD-WXX9:~/myprog/crypt_proto/test$ go run .
2026/01/18 19:28:08.045527 START TRIANGLE 12.00 USDT
2026/01/18 19:28:09.147994 LEG1 USDT→DASH | DASH=0.140400 | total=1.102340872s
2026/01/18 19:28:09.148061 LEG2 skipped | DASH below min step: 0.000000
exit status 1





