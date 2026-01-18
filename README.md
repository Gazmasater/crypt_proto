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
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

/* ================= CONFIG ================= */

const (
	apiKey        = "YOUR_API_KEY"
	apiSecret     = "YOUR_API_SECRET"
	apiPassphrase = "YOUR_API_PASSPHRASE"

	baseURL     = "https://api.kucoin.com"
	startUSDT   = 12.0
	sym1        = "DASH-USDT"
	sym2        = "DASH-BTC"
	sym3        = "BTC-USDT"

	// Шаги ордеров (можно взять через getSymbolStep или константы)
	step1 = 0.0001
	step2 = 0.0001
	step3 = 0.00001
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

/* ================= GET WS TOKEN ================= */

func getWSToken() (string, string, error) {
	req, _ := http.NewRequest("POST", baseURL+"/api/v1/bullet-public", nil)
	req.Header = headers("POST", "/api/v1/bullet-public", "")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	var r struct {
		Code string `json:"code"`
		Data struct {
			Token string `json:"token"`
			InstanceServers []struct {
				Endpoint string `json:"endpoint"`
			} `json:"instanceServers"`
		} `json:"data"`
	}
	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &r); err != nil {
		return "", "", err
	}
	if r.Code != "200000" || len(r.Data.InstanceServers) == 0 {
		return "", "", fmt.Errorf("failed to get WS token")
	}
	return r.Data.InstanceServers[0].Endpoint, r.Data.Token, nil
}

/* ================= LEG STRUCT ================= */

type Leg struct {
	Name      string
	Symbol    string
	Side      string
	Amount    float64
	ClientOid string
	Done      bool
	Filled    float64
}

/* ================= UTILS ================= */

func roundDown(v, step float64) float64 {
	return float64(int(v/step)) * step
}

/* ================= EXECUTION ================= */

func sendOrder(conn *websocket.Conn, leg *Leg) error {
	data := map[string]string{
		"symbol":    leg.Symbol,
		"side":      leg.Side,
		"type":      "market",
		"clientOid": leg.ClientOid,
	}
	if leg.Side == "buy" {
		data["funds"] = fmt.Sprintf("%.8f", leg.Amount)
	} else {
		data["size"] = fmt.Sprintf("%.8f", leg.Amount)
	}
	msg := map[string]interface{}{
		"id":    uuid.NewString(),
		"type":  "order",
		"topic": "/spotMarket/tradeOrders",
		"data":  data,
	}
	return conn.WriteJSON(msg)
}

func waitForFill(conn *websocket.Conn, leg *Leg) {
	for {
		var ev map[string]interface{}
		if err := conn.ReadJSON(&ev); err != nil {
			log.Println("WS read error:", err)
			continue
		}
		if ev["type"] == "message" {
			data := ev["data"].(map[string]interface{})
			if data["clientOid"] == leg.ClientOid && data["status"] == "done" {
				leg.Filled, _ = strconv.ParseFloat(data["filledSize"].(string), 64)
				log.Printf("%s filled: %.8f", leg.Name, leg.Filled)
				return
			}
		}
	}
}

/* ================= MAIN ================= */

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.Printf("START TRIANGLE %.2f USDT", startUSDT)

	endpoint, token, err := getWSToken()
	if err != nil {
		log.Fatal("WS token error:", err)
	}
	wsURL := fmt.Sprintf("%s?token=%s", endpoint, token)

	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		log.Fatal("WS dial error:", err)
	}
	defer conn.Close()

	legs := []Leg{
		{Name: "LEG1", Symbol: sym1, Side: "buy", Amount: startUSDT, ClientOid: uuid.NewString()},
		{Name: "LEG2", Symbol: sym2, Side: "sell", ClientOid: uuid.NewString()},
		{Name: "LEG3", Symbol: sym3, Side: "sell", ClientOid: uuid.NewString()},
	}

	// LEG1
	if err := sendOrder(conn, &legs[0]); err != nil {
		log.Fatal("LEG1 send error:", err)
	}
	waitForFill(conn, &legs[0])

	// LEG2
	legs[1].Amount = legs[0].Filled
	if err := sendOrder(conn, &legs[1]); err != nil {
		log.Fatal("LEG2 send error:", err)
	}
	waitForFill(conn, &legs[1])

	// LEG3
	legs[2].Amount = legs[1].Filled
	if err := sendOrder(conn, &legs[2]); err != nil {
		log.Fatal("LEG3 send error:", err)
	}
	waitForFill(conn, &legs[2])

	log.Println("====== TRIANGLE COMPLETE ======")
	log.Printf("USDT final: %.6f | PNL: %.6f", legs[2].Filled, legs[2].Filled-startUSDT)
}



