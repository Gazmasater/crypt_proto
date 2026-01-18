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

	var (
		dash float64
		btc  float64
		usdt float64
	)

	var (
		leg1Time time.Duration
		leg2Time time.Duration
		leg3Time time.Duration
	)

	/* ================= LEG 1 ================= */
	{
		legStart := time.Now()

		apiStart := time.Now()
		_, err := placeMarket(sym1, "buy", startUSDT)
		apiDur := time.Since(apiStart)
		if err != nil {
			log.Fatal("LEG1 BUY failed:", err)
		}

		time.Sleep(50 * time.Millisecond)

		balStart := time.Now()
		dash, err = getBalance("DASH")
		balDur := time.Since(balStart)
		if err != nil || dash <= 0 {
			log.Fatal("LEG1 balance failed:", err)
		}

		leg1Time = time.Since(legStart)

		log.Printf(
			"LEG1 USDT→DASH | DASH=%.6f | api=%s | balance=%s | total=%s",
			dash, apiDur, balDur, leg1Time,
		)
	}

	/* ================= LEG 2 ================= */
	{
		legStart := time.Now()

		dash = roundDown(dash, 0.001)

		apiStart := time.Now()
		_, err := placeMarket(sym2, "sell", dash)
		apiDur := time.Since(apiStart)
		if err != nil {
			log.Fatal("LEG2 SELL failed:", err)
		}

		time.Sleep(50 * time.Millisecond)

		balStart := time.Now()
		btc, err = getBalance("BTC")
		balDur := time.Since(balStart)
		if err != nil || btc <= 0 {
			log.Fatal("LEG2 balance failed:", err)
		}

		leg2Time = time.Since(legStart)

		log.Printf(
			"LEG2 DASH→BTC | BTC=%.8f | api=%s | balance=%s | total=%s",
			btc, apiDur, balDur, leg2Time,
		)
	}

	/* ================= LEG 3 ================= */
	{
		legStart := time.Now()

		btc = roundDown(btc, 0.0001)

		if btc >= 0.0001 {
			apiStart := time.Now()
			_, err := placeMarket(sym3, "sell", btc)
			apiDur := time.Since(apiStart)
			if err != nil {
				log.Fatal("LEG3 SELL failed:", err)
			}
			log.Printf("LEG3 API sell time=%s", apiDur)
		} else {
			log.Printf("LEG3 skipped | BTC dust=%.8f", btc)
		}

		time.Sleep(50 * time.Millisecond)

		balStart := time.Now()
		usdt, _ = getBalance("USDT")
		balDur := time.Since(balStart)

		leg3Time = time.Since(legStart)

		log.Printf(
			"LEG3 BTC→USDT | USDT=%.4f | balance=%s | total=%s",
			usdt, balDur, leg3Time,
		)
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

/* ================= API ================= */

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





