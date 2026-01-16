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
	apiKey        = "696935c42a6dcd00013273f2"
	apiSecret     = "b348b686-55ff-4290-897b-02d55f815f65"
	apiPassphrase = "Gazmaster_358"

	baseURL = "https://api.kucoin.com"

	startUSDT = 12.0

	sym1 = "DASH-USDT"
	sym2 = "DASH-BTC"
	sym3 = "BTC-USDT"
)

/* ================= TYPES ================= */

type Step int

const (
	StepIdle Step = iota
	StepLeg1
	StepLeg2
	StepLeg3
)

type WSMsg struct {
	Type  string                 `json:"type"`
	Topic string                 `json:"topic"`
	Data  map[string]interface{} `json:"data"`
}

type WSToken struct {
	Token    string
	Endpoint string
}

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

/* ================= REST ================= */

func placeMarket(symbol, side string, value float64) (string, error) {
	body := map[string]string{
		"symbol": symbol,
		"type":   "market",
		"side":   side,
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
		return "", err
	}
	defer resp.Body.Close()

	var r struct {
		Code string `json:"code"`
		Data struct {
			OrderId string `json:"orderId"`
		} `json:"data"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&r)

	if r.Code != "200000" {
		return "", fmt.Errorf("order rejected")
	}

	return r.Data.OrderId, nil
}

/* ================= PRIVATE WS ================= */

func getWSToken() WSToken {
	req, _ := http.NewRequest("POST", baseURL+"/api/v1/bullet-private", nil)
	req.Header = headers("POST", "/api/v1/bullet-private", "")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var r struct {
		Code string `json:"code"`
		Data struct {
			Token           string `json:"token"`
			InstanceServers []struct {
				Endpoint string `json:"endpoint"`
			} `json:"instanceServers"`
		} `json:"data"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&r)

	if len(r.Data.InstanceServers) == 0 {
		log.Fatal("No WS endpoint")
	}

	return WSToken{
		Token:    r.Data.Token,
		Endpoint: r.Data.InstanceServers[0].Endpoint,
	}
}

func connectPrivateWS() *websocket.Conn {
	token := getWSToken()
	url := token.Endpoint + "?token=" + token.Token

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal(err)
	}

	sub := map[string]interface{}{
		"id":    strconv.FormatInt(time.Now().UnixNano(), 10),
		"type":  "subscribe",
		"topic": "/spotMarket/tradeOrders",
	}
	_ = conn.WriteJSON(sub)
	log.Println("Private WS connected")
	return conn
}

/* ================= MAIN ================= */

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.Printf("START TRIANGLE %.2f USDT", startUSDT)

	ws := connectPrivateWS()

	leg1Done := make(chan struct{})
	leg2Done := make(chan struct{})
	leg3Done := make(chan struct{})
	errorChan := make(chan string, 1)

	var leg1Size, leg2Funds, leg3Funds float64
	step := StepIdle

	// WS listener
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			_, msg, err := ws.ReadMessage()
			if err != nil {
				log.Println("[WS ERROR]", err)
				return
			}

			var m WSMsg
			_ = json.Unmarshal(msg, &m)
			if m.Type != "message" || m.Topic != "/spotMarket/tradeOrders" {
				continue
			}

			status, _ := m.Data["status"].(string)
			symbol, _ := m.Data["symbol"].(string)
			side, _ := m.Data["side"].(string)
			size, _ := strconv.ParseFloat(m.Data["filledSize"].(string), 64)
			funds, _ := strconv.ParseFloat(m.Data["filledFunds"].(string), 64)

			if status == "done" {
				switch step {
				case StepLeg1:
					if symbol == sym1 && side == "buy" {
						leg1Size = size
						log.Printf("[FILL] LEG1 USDT→DASH size=%.8f funds=%.8f", size, funds)
						step = StepLeg2
						close(leg1Done)
					}
				case StepLeg2:
					if symbol == sym2 && side == "sell" {
						leg2Funds = funds
						log.Printf("[FILL] LEG2 DASH→BTC size=%.8f funds=%.8f", size, funds)
						step = StepLeg3
						close(leg2Done)
					}
				case StepLeg3:
					if symbol == sym3 && side == "sell" {
						leg3Funds = funds
						log.Printf("[FILL] LEG3 BTC→USDT size=%.8f funds=%.8f", size, funds)
						close(leg3Done)
					}
				}
			} else if status == "rejected" {
				log.Printf("[REJECTED] %s %s", symbol, side)
				select {
				case errorChan <- fmt.Sprintf("%s %s rejected", symbol, side):
				default:
				}
			}
		}
	}()

	// STEP 1
	step = StepLeg1
	o1, err := placeMarket(sym1, "buy", startUSDT)
	if err != nil {
		log.Println("[FAIL] LEG1 USDT→DASH", err)
		ws.Close()
		return
	}
	log.Printf("[OK] LEG1 orderId=%s", o1)
	select {
	case <-leg1Done:
	case errMsg := <-errorChan:
		log.Println("[ORDER ERROR]", errMsg)
		ws.Close()
		return
	}

	// STEP 2
	step = StepLeg2
	o2, err := placeMarket(sym2, "sell", leg1Size)
	if err != nil {
		log.Println("[FAIL] LEG2 DASH→BTC", err)
		ws.Close()
		return
	}
	log.Printf("[OK] LEG2 orderId=%s", o2)
	select {
	case <-leg2Done:
	case errMsg := <-errorChan:
		log.Println("[ORDER ERROR]", errMsg)
		ws.Close()
		return
	}

	// STEP 3
	step = StepLeg3
	o3, err := placeMarket(sym3, "sell", leg2Funds)
	if err != nil {
		log.Println("[FAIL] LEG3 BTC→USDT", err)
		ws.Close()
		return
	}
	log.Printf("[OK] LEG3 orderId=%s", o3)
	select {
	case <-leg3Done:
	case errMsg := <-errorChan:
		log.Println("[ORDER ERROR]", errMsg)
		ws.Close()
		return
	}

	// Итог
	profit := leg3Funds - startUSDT
	pct := profit / startUSDT * 100
	log.Println("====== RESULT ======")
	log.Printf("START: %.4f USDT", startUSDT)
	log.Printf("END:   %.4f USDT", leg3Funds)
	log.Printf("PNL:   %.6f USDT (%.4f%%)", profit, pct)

	// Закрываем WS после всех шагов
	ws.Close()
	<-done
}
