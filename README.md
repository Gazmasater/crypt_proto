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

	baseURL   = "https://api.kucoin.com"
	startUSDT = 12.0

	// Пары треугольника
	sym1 = "DASH-USDT"
	sym2 = "DASH-BTC"
	sym3 = "BTC-USDT"

	// Шаги ордеров (можно получить через REST /symbols)
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

/* ================= PRIVATE WS TOKEN ================= */

type BulletResp struct {
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

	var r BulletResp
	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &r); err != nil {
		return "", "", err
	}

	if r.Code != "200000" || len(r.Data.InstanceServers) == 0 {
		return "", "", fmt.Errorf("failed to get WS token")
	}

	return r.Data.Token, r.Data.InstanceServers[0].Endpoint, nil
}

/* ================= WS HANDLER ================= */

func connectWS(token, endpoint string) (*websocket.Conn, error) {
	wsURL := endpoint + "?token=" + token
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func subscribeOrders(ws *websocket.Conn) error {
	sub := map[string]interface{}{
		"type":           "subscribe",
		"topic":          "/spotMarket/tradeOrders",
		"privateChannel": true,
		"response":       true,
	}
	return ws.WriteJSON(sub)
}

func heartbeat(ws *websocket.Conn) {
	t := time.NewTicker(30 * time.Second)
	go func() {
		for range t.C {
			ws.WriteJSON(map[string]string{"id": uuid.NewString(), "type": "ping"})
		}
	}()
}

func waitForFill(ws *websocket.Conn, clientOid string) {
	for {
		var msg map[string]interface{}
		if err := ws.ReadJSON(&msg); err != nil {
			log.Println("WS read error:", err)
			continue
		}
		if data, ok := msg["data"].(map[string]interface{}); ok {
			if data["clientOid"] == clientOid && data["status"] == "done" {
				log.Printf("Order %s filled", clientOid)
				return
			}
		}
	}
}

/* ================= API PLACE ORDER ================= */

func placeMarket(symbol, side string, value float64) (string, error) {
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
	return clientOid, nil
}

/* ================= GET BALANCE ================= */

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

/* ================= UTILS ================= */

func roundDown(v, step float64) float64 {
	if step == 0 {
		return v
	}
	return float64(int(v/step)) * step
}

/* ================= MAIN ================= */

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.Printf("START TRIANGLE %.2f USDT", startUSDT)

	// Получаем private WS токен
	token, endpoint, err := getPrivateWSToken()
	if err != nil {
		log.Fatal("Failed to get WS token:", err)
	}

	ws, err := connectWS(token, endpoint)
	if err != nil {
		log.Fatal("WS connect error:", err)
	}
	defer ws.Close()

	heartbeat(ws)
	if err := subscribeOrders(ws); err != nil {
		log.Fatal("WS subscribe error:", err)
	}

	/* ================= LEG 1: USDT → DASH ================= */
	leg1Oid, err := placeMarket(sym1, "buy", startUSDT)
	if err != nil {
		log.Fatal("LEG1 BUY failed:", err)
	}
	log.Println("LEG1 order sent, waiting for fill...")
	waitForFill(ws, leg1Oid)
	dash, _ := getBalance("DASH")
	dash = roundDown(dash, step1)
	log.Printf("LEG1 OK | DASH=%.6f", dash)

	/* ================= LEG 2: DASH → BTC ================= */
	leg2Oid, err := placeMarket(sym2, "sell", dash)
	if err != nil {
		log.Fatal("LEG2 SELL failed:", err)
	}
	log.Println("LEG2 order sent, waiting for fill...")
	waitForFill(ws, leg2Oid)
	btc, _ := getBalance("BTC")
	btc = roundDown(btc, step2)
	log.Printf("LEG2 OK | BTC=%.8f", btc)

	/* ================= LEG 3: BTC → USDT ================= */
	leg3Oid, err := placeMarket(sym3, "sell", btc)
	if err != nil {
		log.Fatal("LEG3 SELL failed:", err)
	}
	log.Println("LEG3 order sent, waiting for fill...")
	waitForFill(ws, leg3Oid)
	usdt, _ := getBalance("USDT")
	usdt = roundDown(usdt, step3)
	log.Printf("LEG3 OK | USDT=%.6f", usdt)

	/* ================= SUMMARY ================= */
	log.Println("====== SUMMARY ======")
	log.Printf("PNL: %.6f USDT", usdt-startUSDT)
}



026/01/18 23:55:07.197640 START TRIANGLE 12.00 USDT
2026/01/18 23:55:09.780668 LEG1 order sent, waiting for fill...
2026/01/18 23:55:09.780836 Order 6ce46516-2f44-4de3-bd39-3e21b44803ea filled
2026/01/18 23:55:10.236562 LEG1 OK | DASH=0.286900
2026/01/18 23:55:10.649091 LEG2 order sent, waiting for fill...
2026/01/18 23:55:10.649282 Order 4e49ab8d-a3a1-498d-b4c5-c87ad379e6db filled
2026/01/18 23:55:11.055396 LEG2 OK | BTC=0.00020000
2026/01/18 23:55:11.465079 LEG3 order sent, waiting for fill...
2026/01/18 23:55:11.465235 Order d744c14e-ed99-4b93-998a-4e535924148d filled
2026/01/18 23:55:11.874653 LEG3 OK | USDT=29.723960
2026/01/18 23:55:11.874685 ====== SUMMARY ======
2026/01/18 23:55:11.874693 PNL: 17.723960 USDT






