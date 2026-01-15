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
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

const (
	apiKey        = "KUCOIN_API_KEY"
	apiSecret     = "KUCOIN_API_SECRET"
	apiPassphrase = "KUCOIN_API_PASSPHRASE"

	baseURL = "https://api.kucoin.com"
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
			Funds string `json:"funds"`
			Fee   string `json:"fee"`
		} `json:"items"`
	} `json:"data"`
}

func sign(ts, method, path, body string) string {
	msg := ts + method + path + body
	mac := hmac.New(sha256.New, []byte(apiSecret))
	mac.Write([]byte(msg))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func headers(method, path, body string) http.Header {
	ts := strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
	sig := sign(ts, method, path, body)

	h := http.Header{}
	h.Set("KC-API-KEY", apiKey)
	h.Set("KC-API-SIGN", sig)
	h.Set("KC-API-TIMESTAMP", ts)
	h.Set("KC-API-PASSPHRASE", apiPassphrase)
	h.Set("KC-API-KEY-VERSION", "2")
	h.Set("Content-Type", "application/json")
	return h
}

func placeMarket(symbol, side string, funds float64) string {
	bodyMap := map[string]any{
		"symbol": symbol,
		"side":   side,
		"type":   "market",
	}

	if side == "buy" {
		bodyMap["funds"] = fmt.Sprintf("%.8f", funds)
	} else {
		bodyMap["size"] = fmt.Sprintf("%.8f", funds)
	}

	body, _ := json.Marshal(bodyMap)

	req, _ := http.NewRequest(
		"POST",
		baseURL+"/api/v1/orders",
		bytes.NewReader(body),
	)
	req.Header = headers("POST", "/api/v1/orders", string(body))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("ORDER ERROR: %v", err)
	}
	defer resp.Body.Close()

	raw, _ := io.ReadAll(resp.Body)
	var r OrderResp
	json.Unmarshal(raw, &r)

	if r.Code != "200000" {
		log.Fatalf("ORDER FAIL %s: %s", symbol, string(raw))
	}

	log.Printf("ORDER OK %s %s id=%s", side, symbol, r.Data.OrderId)
	return r.Data.OrderId
}

func waitFill(orderID string) (float64, float64) {
	time.Sleep(200 * time.Millisecond) // KuCoin latency

	path := "/api/v1/fills?orderId=" + orderID
	req, _ := http.NewRequest("GET", baseURL+path, nil)
	req.Header = headers("GET", path, "")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("FILL ERROR: %v", err)
	}
	defer resp.Body.Close()

	raw, _ := io.ReadAll(resp.Body)
	var r FillResp
	json.Unmarshal(raw, &r)

	if len(r.Data.Items) == 0 {
		log.Fatalf("NO FILLS %s", orderID)
	}

	var size, funds float64
	for _, f := range r.Data.Items {
		s, _ := strconv.ParseFloat(f.Size, 64)
		fd, _ := strconv.ParseFloat(f.Funds, 64)
		size += s
		funds += fd
	}

	return size, funds
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	startUSDT := 20.0
	usdt := startUSDT

	log.Println("START TRIANGLE WITH", usdt, "USDT")

	// 1) USDT → X
	o1 := placeMarket("X-USDT", "buy", usdt)
	xQty, usdtSpent := waitFill(o1)
	log.Printf("LEG1 OK: bought X=%.8f spent USDT=%.8f", xQty, usdtSpent)

	// 2) X → BTC
	o2 := placeMarket("X-BTC", "sell", xQty)
	xSold, btcGot := waitFill(o2)
	log.Printf("LEG2 OK: sold X=%.8f got BTC=%.8f", xSold, btcGot)

	// 3) BTC → USDT
	o3 := placeMarket("BTC-USDT", "sell", btcGot)
	btcSold, usdtFinal := waitFill(o3)
	log.Printf("LEG3 OK: sold BTC=%.8f got USDT=%.8f", btcSold, usdtFinal)

	profit := usdtFinal - startUSDT
	pct := profit / startUSDT * 100

	log.Println("------ RESULT ------")
	log.Printf("START: %.4f USDT", startUSDT)
	log.Printf("END:   %.4f USDT", usdtFinal)
	log.Printf("PNL:   %.6f USDT (%.4f%%)", profit, pct)
}


