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
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"sync"
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

	step1 = 0.0001
	step2 = 0.00001
	step3 = 0.01

	orderTimeout = 6 * time.Second
	safetyMargin = 0.999 // учитываем комиссию 0.1%
)

/* ================= AUTH ================= */

func sign(ts, method, path, body string) string {
	mac := hmac.New(sha256.New, []byte(apiSecret))
	mac.Write([]byte(ts + method + path + body))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func signPassphrase() string {
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
	h.Set("KC-API-PASSPHRASE", signPassphrase())
	h.Set("KC-API-KEY-VERSION", "2")
	h.Set("Content-Type", "application/json")
	return h
}

/* ================= REST ================= */

func placeMarketREST(symbol, side string, value float64) (float64, error) {
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
		return 0, err
	}
	defer resp.Body.Close()

	var r struct {
		Code string `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			FilledSize string `json:"filledSize"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return 0, err
	}
	if r.Code != "200000" {
		return 0, fmt.Errorf("order rejected: %s", r.Msg)
	}
	filled, _ := strconv.ParseFloat(r.Data.FilledSize, 64)
	return filled, nil
}

/* ================= WEBSOCKET ================= */

func getPrivateWSURL() (string, error) {
	req, _ := http.NewRequest("POST", baseURL+"/api/v1/bullet-private", nil)
	req.Header = headers("POST", "/api/v1/bullet-private", "")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
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
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return "", err
	}
	return r.Data.InstanceServers[0].Endpoint + "?token=" + r.Data.Token, nil
}

/* ================= ORDER ROUTER ================= */

type OrderRouter struct {
	mu   sync.Mutex
	wait map[string]chan float64
}

func NewOrderRouter() *OrderRouter {
	return &OrderRouter{wait: make(map[string]chan float64)}
}

func (r *OrderRouter) Register(oid string) chan float64 {
	r.mu.Lock()
	defer r.mu.Unlock()
	ch := make(chan float64, 1)
	r.wait[oid] = ch
	return ch
}

func (r *OrderRouter) Resolve(oid string, filled float64) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if ch, ok := r.wait[oid]; ok {
		ch <- filled
		delete(r.wait, oid)
	}
}

/* ================= UTILS ================= */

func roundDown(v, step float64) float64 {
	return math.Floor(v/step) * step
}

/* ================= MAIN ================= */

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.Println("START TRIANGLE", startUSDT)

	// ===== LEG 1 через REST (быстро) =====
	dash, err := placeMarketREST(sym1, "buy", startUSDT)
	if err != nil {
		log.Fatal("LEG1 error:", err)
	}
	dash = roundDown(dash, step1)
	log.Println("LEG1 OK", dash)

	// ===== LEG2 и LEG3 через WS =====
	wsURL, err := getPrivateWSURL()
	if err != nil {
		log.Fatal("WS token error:", err)
	}
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		log.Fatal("WS connect error:", err)
	}
	defer conn.Close()

	// heartbeat
	go func() {
		t := time.NewTicker(25 * time.Second)
		defer t.Stop()
		for range t.C {
			conn.WriteMessage(websocket.TextMessage, []byte(`{"type":"ping"}`))
		}
	}()

	// subscribe
	sub := map[string]interface{}{
		"id":             uuid.NewString(),
		"type":           "subscribe",
		"topic":          "/spotMarket/tradeOrders",
		"privateChannel": true,
	}
	raw, _ := json.Marshal(sub)
	conn.WriteMessage(websocket.TextMessage, raw)

	router := NewOrderRouter()

	// WS reader
	go func() {
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Fatal("WS read:", err)
			}
			var evt struct {
				Topic string `json:"topic"`
				Data  struct {
					Type       string `json:"type"`
					ClientOid  string `json:"clientOid"`
					FilledSize string `json:"filledSize"`
				} `json:"data"`
			}
			if json.Unmarshal(msg, &evt) != nil {
				continue
			}
			if evt.Topic == "/spotMarket/tradeOrders" && evt.Data.Type == "done" {
				v, _ := strconv.ParseFloat(evt.Data.FilledSize, 64)
				router.Resolve(evt.Data.ClientOid, v)
			}
		}
	}()

	// ===== LEG 2 =====
	oid2 := uuid.NewString()
	ch2 := router.Register(oid2)
	if err := placeMarketWithOid(sym2, "sell", dash, oid2); err != nil {
		log.Fatal("LEG2 error:", err)
	}

	ctx2, cancel2 := context.WithTimeout(context.Background(), orderTimeout)
	defer cancel2()
	var btc float64
	select {
	case v := <-ch2:
		btc = roundDown(v*safetyMargin, step2)
	case <-ctx2.Done():
		log.Fatal("LEG2 timeout")
	}
	log.Println("LEG2 OK", btc)

	// ===== LEG 3 =====
	if btc < step2 {
		log.Fatal("LEG3 skipped, not enough BTC")
	}
	oid3 := uuid.NewString()
	ch3 := router.Register(oid3)
	if err := placeMarketWithOid(sym3, "sell", btc, oid3); err != nil {
		log.Fatal("LEG3 error:", err)
	}
	ctx3, cancel3 := context.WithTimeout(context.Background(), orderTimeout)
	defer cancel3()
	var usdt float64
	select {
	case v := <-ch3:
		usdt = roundDown(v*safetyMargin, step3)
	case <-ctx3.Done():
		log.Fatal("LEG3 timeout")
	}
	log.Println("LEG3 OK", usdt)
	log.Println("PNL:", usdt-startUSDT)
}


[{
	"resource": "/home/gaz358/myprog/crypt_proto/test/main.go",
	"owner": "_generated_diagnostic_collection_name_#0",
	"code": {
		"value": "UndeclaredName",
		"target": {
			"$mid": 1,
			"path": "/golang.org/x/tools/internal/typesinternal",
			"scheme": "https",
			"authority": "pkg.go.dev",
			"fragment": "UndeclaredName"
		}
	},
	"severity": 8,
	"message": "undefined: placeMarketWithOid",
	"source": "compiler",
	"startLineNumber": 247,
	"startColumn": 12,
	"endLineNumber": 247,
	"endColumn": 30,
	"modelVersionId": 4,
	"origin": "extHost1"
}]
