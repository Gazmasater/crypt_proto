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
	apiKey    = "KUCOIN_API_KEY"
	apiSecret = "KUCOIN_API_SECRET"
	apiPass   = "KUCOIN_API_PASSPHRASE"

	startUSDT = 20.0
	symbol1   = "X-USDT"
	symbol2   = "X-BTC"
	symbol3   = "BTC-USDT"
	baseURL   = "https://api.kucoin.com"
	fee       = 0.001
)

/* ================= STATE ================= */

type Step uint8

const (
	Idle Step = iota
	Buy1
	Sell2
	Sell3
)

var step = Idle
var lastAmount float64

/* ================= MAIN ================= */

func main() {
	log.Println("START")

	// Получаем Private WS токен
	wsToken := getPrivateWSToken()

	// Подключаем WS
	conn := connectPrivateWS(wsToken)

	// Аутентификация через WS
	loginPrivateWS(conn)
	log.Println("Private WS connected & authenticated ✅")

	// Читаем события WS
	go readWS(conn)

	// Стартуем треугольник
	step = Buy1
	lastAmount = startUSDT
	log.Printf("[STEP 1] Buying %s for %.2f USDT\n", symbol1, startUSDT)
	placeMarketWS(conn, symbol1, "buy", startUSDT)

	select {} // блок main
}

/* ================= PRIVATE WS ================= */

func connectPrivateWS(token WSToken) *websocket.Conn {
	url := token.Endpoint + "?token=" + token.Token
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Подписка на ордера
	sub := map[string]any{
		"id":    strconv.FormatInt(time.Now().UnixNano(), 10),
		"type":  "subscribe",
		"topic": "/spotMarket/tradeOrders",
	}
	if err := conn.WriteJSON(sub); err != nil {
		log.Fatal(err)
	}
	return conn
}

func readWS(conn *websocket.Conn) {
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Fatal(err)
		}
		handleOrderEvent(msg, conn)
	}
}

func loginPrivateWS(conn *websocket.Conn) {
	ts := strconv.FormatInt(time.Now().UnixMilli(), 10)
	sign := hmacSHA256Base64(ts+"POST"+"/users/self/verify", apiSecret)
	pass := hmacSHA256Base64(apiPass, apiSecret)

	loginMsg := map[string]any{
		"id":         strconv.FormatInt(time.Now().UnixNano(), 10),
		"type":       "login",
		"clientOid":  strconv.FormatInt(time.Now().UnixNano(), 10),
		"apiKey":     apiKey,
		"passphrase": pass,
		"timestamp":  ts,
		"signature":  sign,
	}
	if err := conn.WriteJSON(loginMsg); err != nil {
		log.Fatal(err)
	}
}

/* ================= WS HANDLER ================= */

func handleOrderEvent(msg []byte, conn *websocket.Conn) {
	var m map[string]any
	if err := json.Unmarshal(msg, &m); err != nil {
		return
	}

	if m["type"] != "message" {
		return
	}

	data, ok := m["data"].(map[string]any)
	if !ok || data["status"].(string) != "done" {
		return
	}

	filledSize, _ := strconv.ParseFloat(data["filledSize"].(string), 64)
	lastPrice, _ := strconv.ParseFloat(data["price"].(string), 64)

	switch step {
	case Buy1:
		log.Printf("[STEP 1 DONE] Bought %.8f %s at price %.8f\n", filledSize, symbol1, lastPrice)
		step = Sell2
		log.Printf("[STEP 2] Selling %s for %s\n", symbol2, symbol3)
		placeMarketWS(conn, symbol2, "sell", filledSize)

	case Sell2:
		log.Printf("[STEP 2 DONE] Sold %.8f %s at price %.8f\n", filledSize, symbol2, lastPrice)
		step = Sell3
		placeMarketWS(conn, symbol3, "sell", filledSize)

	case Sell3:
		log.Printf("[STEP 3 DONE] Sold %.8f %s at price %.8f\n", filledSize, symbol3, lastPrice)
		profit := filledSize - startUSDT
		profitPct := (profit / startUSDT) * 100
		log.Printf("[ARB DONE ✅] Profit: %.6f USDT (%.4f%%)\n", profit, profitPct)
		step = Idle
	}
}

/* ================= MARKET ORDER WS ================= */

func placeMarketWS(conn *websocket.Conn, symbol, side string, size float64) {
	msg := map[string]any{
		"id":        strconv.FormatInt(time.Now().UnixNano(), 10),
		"type":      "order",
		"clientOid": strconv.FormatInt(time.Now().UnixNano(), 10),
		"side":      side,
		"symbol":    symbol,
		"size":      fmt.Sprintf("%.8f", size),
		"orderType": "market",
	}
	if err := conn.WriteJSON(msg); err != nil {
		log.Println("WS order error:", err)
	}
}

/* ================= AUTH HELPERS ================= */

func hmacSHA256Base64(msg, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(msg))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

/* ================= WS TOKEN ================= */

type WSToken struct {
	Token    string
	Endpoint string
}

func getPrivateWSToken() WSToken {
	req, _ := http.NewRequest("POST", baseURL+"/api/v1/bullet-private", nil)

	ts := strconv.FormatInt(time.Now().UnixMilli(), 10)
	sign := hmacSHA256Base64(ts+"POST"+"/api/v1/bullet-private", apiSecret)

	req.Header.Set("KC-API-KEY", apiKey)
	req.Header.Set("KC-API-SIGN", sign)
	req.Header.Set("KC-API-TIMESTAMP", ts)
	req.Header.Set("KC-API-PASSPHRASE", hmacSHA256Base64(apiPass, apiSecret))
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
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		log.Fatal(err)
	}

	return WSToken{
		Token:    r.Data.Token,
		Endpoint: r.Data.InstanceServers[0].Endpoint,
	}
}

