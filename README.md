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
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
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

	// Шаги ордеров (получены через API/Postman)
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

/* ================= WS TOKEN ================= */

type wsTokenResp struct {
	Code string `json:"code"`
	Data struct {
		Token           string `json:"token"`
		InstanceServers []struct {
			Endpoint string `json:"endpoint"`
		} `json:"instanceServers"`
	} `json:"data"`
}

func getWsToken() (string, string, error) {
	req, _ := http.NewRequest("POST", baseURL+"/api/v1/bullet-private", nil)
	req.Header = headers("POST", "/api/v1/bullet-private", "")
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	var r wsTokenResp
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return "", "", err
	}
	if r.Code != "200000" || len(r.Data.InstanceServers) == 0 {
		return "", "", fmt.Errorf("failed to get WS token")
	}
	return r.Data.Token, r.Data.InstanceServers[0].Endpoint, nil
}

/* ================= ORDERS ================= */

type wsOrder struct {
	ID    string      `json:"id"`
	Type  string      `json:"type"`
	Topic string      `json:"topic"`
	Data  interface{} `json:"data"`
}

type orderData struct {
	Symbol   string `json:"symbol"`
	Side     string `json:"side"`
	Type     string `json:"type"`
	Size     string `json:"size,omitempty"`
	Funds    string `json:"funds,omitempty"`
	ClientID string `json:"clientOid"`
}

/* ================= MAIN ================= */

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.Printf("START TRIANGLE %.2f USDT", startUSDT)

	// Получаем WS токен
	token, endpoint, err := getWsToken()
	if err != nil {
		log.Fatal("WS token failed:", err)
	}
	wsURL := fmt.Sprintf("%s?token=%s&connectId=%s", endpoint, token, uuid.NewString())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Подключаемся к WS
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		log.Fatal("WS connect failed:", err)
	}
	defer conn.Close()

	// Канал событий для отслеживания заполнения ордеров
	events := make(chan string, 3)

	// Запускаем слушатель WS
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				_, msg, err := conn.ReadMessage()
				if err != nil {
					log.Println("WS read error:", err)
					continue
				}
				var ev map[string]interface{}
				if err := json.Unmarshal(msg, &ev); err == nil {
					if ev["type"] == "message" {
						data := ev["data"].(map[string]interface{})
						if filled, ok := data["filledSize"]; ok && filled.(string) != "0" {
							events <- data["symbol"].(string)
						}
					}
				}
			}
		}
	}()

	// LEG1: USDT → DASH
	ord1 := wsOrder{
		ID:    uuid.NewString(),
		Type:  "order",
		Topic: "/spotMarket/tradeOrders",
		Data: orderData{
			Symbol:   sym1,
			Side:     "buy",
			Type:     "market",
			Funds:    fmt.Sprintf("%.8f", startUSDT),
			ClientID: uuid.NewString(),
		},
	}
	conn.WriteJSON(ord1)
	log.Println("LEG1 order sent, waiting for fill...")
	<-events
	log.Println("LEG1 filled ✅")

	// LEG2: DASH → BTC
	ord2 := wsOrder{
		ID:    uuid.NewString(),
		Type:  "order",
		Topic: "/spotMarket/tradeOrders",
		Data: orderData{
			Symbol:   sym2,
			Side:     "sell",
			Type:     "market",
			Size:     fmt.Sprintf("%.8f", step1), // минимальный размер
			ClientID: uuid.NewString(),
		},
	}
	conn.WriteJSON(ord2)
	log.Println("LEG2 order sent, waiting for fill...")
	<-events
	log.Println("LEG2 filled ✅")

	// LEG3: BTC → USDT
	ord3 := wsOrder{
		ID:    uuid.NewString(),
		Type:  "order",
		Topic: "/spotMarket/tradeOrders",
		Data: orderData{
			Symbol:   sym3,
			Side:     "sell",
			Type:     "market",
			Size:     fmt.Sprintf("%.8f", step3),
			ClientID: uuid.NewString(),
		},
	}
	conn.WriteJSON(ord3)
	log.Println("LEG3 order sent, waiting for fill...")
	<-events
	log.Println("LEG3 filled ✅")

	log.Println("====== TRIANGLE COMPLETE ======")
}



gaz358@gaz358-BOD-WXX9:~/myprog/crypt_proto/test$ go run .
2026/01/18 20:44:51.097912 START TRIANGLE 12.00 USDT
2026/01/18 20:44:52.992756 LEG1 order sent, waiting for fill...




