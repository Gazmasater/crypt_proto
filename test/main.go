package main

import (
	"bytes"
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

	baseURL = "https://api.kucoin.com"

	startUSDT = 12.0

	sym1 = "DASH-USDT"
	sym2 = "DASH-BTC"
	sym3 = "BTC-USDT"

	stepDash = 0.0001
	stepBTC  = 0.00001
	stepUSDT = 0.01

	fee = 0.999 // 0.1% комиссия
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

/* ================= UTILS ================= */

func roundDown(v, step float64) float64 {
	return math.Floor(v/step) * step
}

/* ================= FAST REST ================= */

var fastHTTP = &http.Client{
	Transport: &http.Transport{
		MaxIdleConns:        200,
		MaxIdleConnsPerHost: 200,
		IdleConnTimeout:     90 * time.Second,
		DisableCompression:  true,
	},
}

func sendMarket(symbol, side string, value float64, oid string) {
	body := map[string]string{
		"symbol":    symbol,
		"type":      "market",
		"side":      side,
		"clientOid": oid,
	}

	if side == "buy" {
		body["funds"] = fmt.Sprintf("%.8f", value)
	} else {
		body["size"] = fmt.Sprintf("%.8f", value)
	}

	raw, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", baseURL+"/api/v1/orders", bytes.NewReader(raw))
	req.Header = headers("POST", "/api/v1/orders", string(raw))

	go fastHTTP.Do(req) // 🔥 FIRE & FORGET
}

/* ================= ORDER ROUTER ================= */

type OrderRouter struct {
	mu   sync.Mutex
	wait map[string]chan float64
}

func NewOrderRouter() *OrderRouter {
	return &OrderRouter{wait: map[string]chan float64{}}
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

/* ================= EXECUTOR ================= */

type Executor struct {
	router *OrderRouter
}

func NewExecutor(r *OrderRouter) *Executor {
	return &Executor{router: r}
}

func (e *Executor) Start(usdt float64) {
	go e.execute(usdt)
}

func (e *Executor) execute(usdt float64) {
	// ===== LEG1 =====
	oid1 := uuid.NewString()
	ch1 := e.router.Register(oid1)
	log.Println("SEND LEG1")
	sendMarket(sym1, "buy", usdt, oid1)

	filledDash := <-ch1
	dash := roundDown(filledDash, stepDash)
	log.Println("LEG1 FILLED:", dash)

	// ===== LEG2 =====
	oid2 := uuid.NewString()
	ch2 := e.router.Register(oid2)
	log.Println("SEND LEG2")
	sendMarket(sym2, "sell", dash, oid2)

	filledBTC := <-ch2
	// учитываем комиссию takerFee 0.1% при продаже BTC
	btc := roundDown(filledBTC*(1-0.001), stepBTC)
	log.Println("LEG2 FILLED (after fee):", btc)

	// ===== LEG3 =====
	if btc < stepBTC {
		log.Println("LEG3 skipped, BTC меньше минимального шага после комиссии")
		return
	}

	oid3 := uuid.NewString()
	ch3 := e.router.Register(oid3)
	log.Println("SEND LEG3")
	sendMarket(sym3, "sell", btc, oid3)

	filledUSDT := <-ch3
	usdtFinal := roundDown(filledUSDT, stepUSDT)
	log.Println("LEG3 FILLED:", usdtFinal)
	log.Println("PNL:", usdtFinal-startUSDT)
}

/* ================= WEBSOCKET ================= */

func getPrivateWS() (string, error) {
	req, _ := http.NewRequest("POST", baseURL+"/api/v1/bullet-private", nil)
	req.Header = headers("POST", "/api/v1/bullet-private", "")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
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
	json.NewDecoder(resp.Body).Decode(&r)

	return r.Data.InstanceServers[0].Endpoint + "?token=" + r.Data.Token, nil
}

/* ================= MAIN ================= */

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.Println("START TRIANGLE", startUSDT)

	wsURL, err := getPrivateWS()
	if err != nil {
		log.Fatal(err)
	}

	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// ping каждые 25 секунд
	go func() {
		t := time.NewTicker(25 * time.Second)
		for range t.C {
			conn.WriteMessage(websocket.TextMessage, []byte(`{"type":"ping"}`))
		}
	}()

	// подписка на private tradeOrders
	sub := fmt.Sprintf(`{
		"id":"%s",
		"type":"subscribe",
		"topic":"/spotMarket/tradeOrders",
		"privateChannel":true
	}`, uuid.NewString())
	conn.WriteMessage(websocket.TextMessage, []byte(sub))

	router := NewOrderRouter()
	exec := NewExecutor(router)

	// WS reader (критично для fill)
	go func() {
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Fatal(err)
			}

			// лог всех сообщений для дебага
			log.Println("WS MSG:", string(msg))

			if !bytes.Contains(msg, []byte("tradeOrders")) {
				continue
			}

			var evt struct {
				Topic string `json:"topic"`
				Data  struct {
					ClientOid   string `json:"clientOid"`
					FilledSize  string `json:"filledSize"`
					FilledFunds string `json:"filledFunds"`
				} `json:"data"`
			}

			if err := json.Unmarshal(msg, &evt); err != nil {
				log.Println("WS unmarshal error:", err)
				continue
			}

			if evt.Data.ClientOid == "" {
				continue
			}

			var filled float64
			if evt.Data.FilledFunds != "" {
				filled, _ = strconv.ParseFloat(evt.Data.FilledFunds, 64)
			} else if evt.Data.FilledSize != "" {
				filled, _ = strconv.ParseFloat(evt.Data.FilledSize, 64)
			}

			if filled > 0 {
				log.Printf("RESOLVE %s filled=%f\n", evt.Data.ClientOid, filled)
				router.Resolve(evt.Data.ClientOid, filled)
			}
		}
	}()
	exec.Start(startUSDT)

	select {} // держим main
}
