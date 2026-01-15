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
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

/* ================= CONFIG ================= */

const (
	apiKey        = "KUCOIN_API_KEY"
	apiSecret     = "KUCOIN_API_SECRET"
	apiPassphrase = "KUCOIN_API_PASSPHRASE"

	baseURL = "https://api.kucoin.com"

	startUSDT = 20.0

	sym1 = "X-USDT"
	sym2 = "X-BTC"
	sym3 = "BTC-USDT"
)

/* ================= TYPES ================= */

type Leg string

const (
	Leg1 Leg = "LEG1 USDT→X"
	Leg2 Leg = "LEG2 X→BTC"
	Leg3 Leg = "LEG3 BTC→USDT"
)

type OrderResp struct {
	Code string `json:"code"`
	Data struct {
		OrderId string `json:"orderId"`
	} `json:"data"`
}

type FillResp struct {
	Code string `json:"code"`
	Data struct {
		Items []struct {
			Size   string `json:"size"`
			Funds  string `json:"funds"`
			Fee    string `json:"fee"`
		} `json:"items"`
	} `json:"data"`
}

type WSToken struct {
	Token    string
	Endpoint string
}

/* ================= AUTH ================= */

func sign(ts, method, path, body string) string {
	msg := ts + method + path + body
	mac := hmac.New(sha256.New, []byte(apiSecret))
	mac.Write([]byte(msg))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func signPassphrase() string {
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
	h.Set("KC-API-PASSPHRASE", signPassphrase())
	h.Set("KC-API-KEY-VERSION", "2")
	h.Set("Content-Type", "application/json")
	return h
}

/* ================= PRIVATE WS ================= */

func getPrivateToken() WSToken {
	req, _ := http.NewRequest("POST", baseURL+"/api/v1/bullet-private", nil)
	ts := strconv.FormatInt(time.Now().UnixMilli(), 10)
	signature := sign(ts, "POST", "/api/v1/bullet-private", "")
	req.Header.Set("KC-API-KEY", apiKey)
	req.Header.Set("KC-API-SIGN", signature)
	req.Header.Set("KC-API-TIMESTAMP", ts)
	req.Header.Set("KC-API-PASSPHRASE", signPassphrase())
	req.Header.Set("KC-API-KEY-VERSION", "2")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var r struct {
		Data struct {
			Token           string `json:"token"`
			InstanceServers []struct {
				Endpoint string `json:"endpoint"`
			} `json:"instanceServers"`
		} `json:"data"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&r)

	return WSToken{
		Token:    r.Data.Token,
		Endpoint: r.Data.InstanceServers[0].Endpoint,
	}
}

func connectPrivateWS() *websocket.Conn {
	token := getPrivateToken()
	url := token.Endpoint + "?token=" + token.Token
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal(err)
	}

	sub := map[string]any{
		"id":    time.Now().UnixNano(),
		"type":  "subscribe",
		"topic": "/spotMarket/tradeOrders",
	}
	_ = conn.WriteJSON(sub)

	log.Println("Private WS connected")
	return conn
}

/* ================= MAIN ================= */

func placeMarket(symbol, side string, funds float64) string {
	body := map[string]any{
		"symbol": symbol,
		"type":   "market",
		"side":   side,
	}

	if side == "buy" {
		body["funds"] = fmt.Sprintf("%.8f", funds)
	} else {
		body["size"] = fmt.Sprintf("%.8f", funds)
	}

	rawBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", baseURL+"/api/v1/orders", bytes.NewReader(rawBody))
	req.Header = headers("POST", "/api/v1/orders", string(rawBody))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("ORDER ERROR: %v", err)
	}
	defer resp.Body.Close()

	raw, _ := io.ReadAll(resp.Body)
	var r OrderResp
	_ = json.Unmarshal(raw, &r)

	if r.Code != "200000" {
		log.Fatalf("ORDER FAIL %s %s resp=%s", side, symbol, raw)
	}

	log.Printf("ORDER OK %s %s orderId=%s", side, symbol, r.Data.OrderId)
	return r.Data.OrderId
}

func waitFillWS(conn *websocket.Conn, orderID string) (float64, float64) {
	for {
		_, msg, _ := conn.ReadMessage()
		var m map[string]any
		_ = json.Unmarshal(msg, &m)
		if m["type"] != "message" {
			continue
		}
		data := m["data"].(map[string]any)
		if data["orderId"] == orderID && data["status"] == "done" {
			size, _ := strconv.ParseFloat(data["filledSize"].(string), 64)
			funds, _ := strconv.ParseFloat(data["filledFunds"].(string), 64)
			log.Printf("[FILL WS] orderId=%s size=%.8f funds=%.8f", orderID, size, funds)
			return size, funds
		}
	}
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.Printf("START TRIANGLE %.2f USDT", startUSDT)

	ws := connectPrivateWS()

	// LEG1
	o1 := placeMarket(sym1, "buy", startUSDT)
	xQty, _ := waitFillWS(ws, o1)

	// LEG2
	o2 := placeMarket(sym2, "sell", xQty)
	_, btcGot := waitFillWS(ws, o2)

	// LEG3
	o3 := placeMarket(sym3, "sell", btcGot)
	_, usdtFinal := waitFillWS(ws, o3)

	// RESULT
	profit := usdtFinal - startUSDT
	pct := profit / startUSDT * 100

	log.Println("====== RESULT ======")
	log.Printf("START: %.4f USDT", startUSDT)
	log.Printf("END:   %.4f USDT", usdtFinal)
	log.Printf("PNL:   %.6f USDT (%.4f%%)", profit, pct)
}


az358@gaz358-BOD-WXX9:~/myprog/crypt_proto/test$ go run .
2026/01/16 01:48:31.442785 START TRIANGLE %!f(int=11) USDT
2026/01/16 01:48:33.049883 Private WS connected
2026/01/16 01:48:33.357082 ORDER FAIL buy DASH-USDT resp={"msg":"validation.createOrder.clientOidIsRequired","code":"400100"}
exit status 1

