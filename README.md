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
	apiKey        = "696935c42a6dcd00013273f2"
	apiSecret     = "b348b686-55ff-4290-897b-02d55f815f65"
	apiPassphrase = "Gazmaster_358"

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
func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	startUSDT := 12.0
	start := time.Now()

	log.Printf("START TRIANGLE %.2f USDT", startUSDT)

	// ================= LEG 1 =================
	t1 := time.Now()

	_, err := placeMarket("DASH-USDT", "buy", startUSDT)
	if err != nil {
		log.Fatal("LEG1 BUY failed:", err)
	}

	time.Sleep(50 * time.Millisecond)

	dash, err := getBalance("DASH")
	if err != nil || dash <= 0 {
		log.Fatal("LEG1 balance failed:", err)
	}

	dur1 := time.Since(t1)
	log.Printf("LEG1 USDT→DASH done | DASH=%.6f | time=%s", dash, dur1)

	// ================= LEG 2 =================
	t2 := time.Now()

	dash = roundDown(dash, 0.001)

	_, err = placeMarket("DASH-BTC", "sell", dash)
	if err != nil {
		log.Fatal("LEG2 SELL failed:", err)
	}

	time.Sleep(50 * time.Millisecond)

	btc, err := getBalance("BTC")
	if err != nil || btc <= 0 {
		log.Fatal("LEG2 balance failed:", err)
	}

	dur2 := time.Since(t2)
	log.Printf("LEG2 DASH→BTC done | BTC=%.8f | time=%s", btc, dur2)

	// ================= LEG 3 =================
	t3 := time.Now()

	btc = roundDown(btc, 0.0001)

	if btc >= 0.0001 {
		_, err = placeMarket("BTC-USDT", "sell", btc)
		if err != nil {
			log.Fatal("LEG3 SELL failed:", err)
		}
	} else {
		log.Printf("LEG3 skipped | BTC dust=%.8f", btc)
	}

	time.Sleep(50 * time.Millisecond)

	usdt, err := getBalance("USDT")
	if err != nil {
		log.Fatal("LEG3 balance failed:", err)
	}

	dur3 := time.Since(t3)
	total := time.Since(start)

	log.Println("====== RESULT ======")
	log.Printf("LEG3 BTC→USDT done | USDT=%.4f | time=%s", usdt, dur3)
	log.Printf("TOTAL TRIANGLE TIME: %s", total)
	log.Printf("PNL: %.6f USDT", usdt-startUSDT)
}

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


2026/01/18 15:22:53.465428 START TRIANGLE 12.00 USDT
2026/01/18 15:22:54.699725 LEG1 USDT→DASH done | DASH=0.150000 | time=1.234217344s
2026/01/18 15:22:55.621540 LEG2 DASH→BTC done | BTC=0.00017771 | time=921.748853ms
2026/01/18 15:22:56.445419 ====== RESULT ======
2026/01/18 15:22:56.445448 LEG3 BTC→USDT done | USDT=25.9971 | time=823.845029ms
2026/01/18 15:22:56.445479 TOTAL TRIANGLE TIME: 2.97999153s
2026/01/18 15:22:56.445488 PNL: 13.997108 USDT
gaz358@gaz358-BOD-WXX9:~/myprog/crypt_proto/test$




