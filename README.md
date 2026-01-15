Название API
9623527002

6966b78122ca320001d2acae
fa1e37ae-21ff-4257-844d-3dcd21d26ccd





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
	apiKey    = "KUCOIN_API_KEY"
	apiSecret = "KUCOIN_API_SECRET"
	apiPass   = "KUCOIN_API_PASSPHRASE"

	baseURL = "https://api.kucoin.com"

	startUSDT = 20.0
	symbol1   = "X-USDT"
	symbol2   = "X-BTC"
	symbol3   = "BTC-USDT"
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

/* ================= MAIN ================= */

func main() {
	log.Println("START")

	ws := connectPrivateWS()
	go readWS(ws)

	// старт треугольника
	step = Buy1
	placeMarketFunds(symbol1, "buy", startUSDT)

	select {} // блокируем main, приложение живёт
}

/* ================= PRIVATE WS ================= */

func connectPrivateWS() *websocket.Conn {
	token := getWSToken()

	url := token.Endpoint + "?token=" + token.Token
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal(err)
	}

	sub := map[string]any{
		"id":   time.Now().UnixNano(),
		"type": "subscribe",
		"topic": "/spotMarket/tradeOrders",
	}
	_ = conn.WriteJSON(sub)

	log.Println("Private WS connected")
	return conn
}

func readWS(conn *websocket.Conn) {
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Fatal(err)
		}
		handleOrderEvent(msg)
	}
}

/* ================= WS HANDLER ================= */

func handleOrderEvent(msg []byte) {
	var m map[string]any
	_ = json.Unmarshal(msg, &m)

	if m["type"] != "message" {
		return
	}

	data := m["data"].(map[string]any)
	if data["status"].(string) != "done" {
		return
	}

	filledSize, _ := strconv.ParseFloat(data["filledSize"].(string), 64)

	switch step {
	case Buy1:
		step = Sell2
		placeMarketSize(symbol2, "sell", filledSize)

	case Sell2:
		step = Sell3
		placeMarketSize(symbol3, "sell", filledSize)

	case Sell3:
		log.Println("ARB DONE")
		step = Idle
	}
}

/* ================= REST ================= */

func placeMarketFunds(symbol, side string, funds float64) {
	body := map[string]any{
		"symbol": symbol,
		"type":   "market",
		"side":   side,
		"funds":  fmt.Sprintf("%.2f", funds),
	}
	fireREST("/api/v1/orders", body)
}

func placeMarketSize(symbol, side string, size float64) {
	body := map[string]any{
		"symbol": symbol,
		"type":   "market",
		"side":   side,
		"size":   fmt.Sprintf("%.8f", size),
	}
	fireREST("/api/v1/orders", body)
}

func fireREST(path string, body any) {
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", baseURL+path, bytes.NewBuffer(b))

	ts := strconv.FormatInt(time.Now().UnixMilli(), 10)
	signature := sign(ts + "POST" + path + string(b))

	req.Header.Set("KC-API-KEY", apiKey)
	req.Header.Set("KC-API-SIGN", signature)
	req.Header.Set("KC-API-TIMESTAMP", ts)
	req.Header.Set("KC-API-PASSPHRASE", passphrase())
	req.Header.Set("KC-API-KEY-VERSION", "2")
	req.Header.Set("Content-Type", "application/json")

	go http.DefaultClient.Do(req) // fire & forget
}

/* ================= AUTH ================= */

func sign(msg string) string {
	mac := hmac.New(sha256.New, []byte(apiSecret))
	mac.Write([]byte(msg))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func passphrase() string {
	mac := hmac.New(sha256.New, []byte(apiSecret))
	mac.Write([]byte(apiPass))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

/* ================= WS TOKEN ================= */

type WSToken struct {
	Token    string
	Endpoint string
}

func getWSToken() WSToken {
	req, _ := http.NewRequest("POST", baseURL+"/api/v1/bullet-private", nil)

	ts := strconv.FormatInt(time.Now().UnixMilli(), 10)
	signature := sign(ts + "POST" + "/api/v1/bullet-private")

	req.Header.Set("KC-API-KEY", apiKey)
	req.Header.Set("KC-API-SIGN", signature)
	req.Header.Set("KC-API-TIMESTAMP", ts)
	req.Header.Set("KC-API-PASSPHRASE", passphrase())
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

