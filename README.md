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

	"github.com/gorilla/websocket"
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

	// Шаги ордеров (можно получать через API)
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

/* ================= GET PRIVATE WS TOKEN ================= */
type BulletPrivateResp struct {
	Code string `json:"code"`
	Data struct {
		Token           string `json:"token"`
		InstanceServers []struct {
			Endpoint string `json:"endpoint"`
		} `json:"instanceServers"`
	} `json:"data"`
}

func getPrivateWSToken() (string, string, error) {
	req, _ := http.NewRequest("POST", baseURL+"/api/v1/bullet-private", nil)
	req.Header = headers("POST", "/api/v1/bullet-private", "")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	var r BulletPrivateResp
	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &r); err != nil {
		return "", "", err
	}
	if r.Code != "200000" || len(r.Data.InstanceServers) == 0 {
		return "", "", fmt.Errorf("failed to get private ws token")
	}

	return r.Data.InstanceServers[0].Endpoint, r.Data.Token, nil
}

/* ================= WEBSOCKET ================= */
type WSMessage struct {
	Type  string `json:"type"`
	Topic string `json:"topic"`
	Data  struct {
		OrderId    string `json:"orderId"`
		Status     string `json:"status"`
		FilledSize string `json:"filledSize"`
	} `json:"data"`
}

func connectWS(endpoint, token string) (*websocket.Conn, error) {
	wsURL := fmt.Sprintf("%s?token=%s", endpoint, token)
	c, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func subscribeOrders(c *websocket.Conn) error {
	msg := map[string]interface{}{
		"id":       uuid.NewString(),
		"type":     "subscribe",
		"topic":    "/spotMarket/tradeOrders",
		"response": true,
	}
	raw, _ := json.Marshal(msg)
	return c.WriteMessage(websocket.TextMessage, raw)
}

/* ================= PLACE ORDER ================= */
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

/* ================= WAIT FOR FILL ================= */
func waitForFill(c *websocket.Conn, orderId string) error {
	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			return err
		}
		var m WSMessage
		if err := json.Unmarshal(msg, &m); err != nil {
			continue
		}
		if m.Type == "message" && m.Topic == "/spotMarket/tradeOrders" && m.Data.OrderId == orderId && m.Data.Status == "done" {
			return nil
		}
	}
}

/* ================= MAIN ================= */
func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.Printf("START TRIANGLE %.2f USDT", startUSDT)

	endpoint, token, err := getPrivateWSToken()
	if err != nil {
		log.Fatal("Get WS token failed:", err)
	}

	wsConn, err := connectWS(endpoint, token)
	if err != nil {
		log.Fatal("WS connect failed:", err)
	}
	defer wsConn.Close()

	if err := subscribeOrders(wsConn); err != nil {
		log.Fatal("Subscribe failed:", err)
	}

	// LEG1: USDT → DASH
	order1, err := placeMarket(sym1, "buy", startUSDT)
	if err != nil {
		log.Fatal("LEG1 BUY failed:", err)
	}
	log.Println("LEG1 order sent, waiting for fill...")
	if err := waitForFill(wsConn, order1); err != nil {
		log.Fatal(err)
	}
	log.Println("LEG1 filled ✅")

	// LEG2: DASH → BTC
	order2, err := placeMarket(sym2, "sell", 0.1365) // пример значения
	if err != nil {
		log.Fatal("LEG2 SELL failed:", err)
	}
	log.Println("LEG2 order sent, waiting for fill...")
	if err := waitForFill(wsConn, order2); err != nil {
		log.Fatal(err)
	}
	log.Println("LEG2 filled ✅")

	// LEG3: BTC → USDT
	order3, err := placeMarket(sym3, "sell", 0.00012) // пример значения
	if err != nil {
		log.Fatal("LEG3 SELL failed:", err)
	}
	log.Println("LEG3 order sent, waiting for fill...")
	if err := waitForFill(wsConn, order3); err != nil {
		log.Fatal(err)
	}
	log.Println("LEG3 filled ✅")

	log.Println("TRIANGLE COMPLETE ✅")
}






