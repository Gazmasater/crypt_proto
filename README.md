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

type Step int

const (
	StepIdle Step = iota
	StepLeg1
	StepLeg2
	StepLeg3
)

var step = StepIdle

type WSMsg struct {
	Type string                 `json:"type"`
	Topic string                `json:"topic"`
	Data map[string]interface{} `json:"data"`
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

type WSToken struct {
	Token    string
	Endpoint string
}

func getWSToken() WSToken {
	req, _ := http.NewRequest("POST", baseURL+"/api/v1/bullet-private", nil)
	ts := strconv.FormatInt(time.Now().UnixMilli(), 10)
	sig := sign(ts, "POST", "/api/v1/bullet-private", "")
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
	defer ws.Close()

	// запускаем горутину WS listener
	done := make(chan struct{})
	var leg1Size, leg2Size, leg3Size float64
	var leg1Funds, leg2Funds, leg3Funds float64

	go func() {
		for {
			_, msg, err := ws.ReadMessage()
			if err != nil {
				log.Fatal(err)
			}

			var m WSMsg
			_ = json.Unmarshal(msg, &m)
			if m.Type != "message" || m.Topic != "/spotMarket/tradeOrders" {
				continue
			}

			status := m.Data["status"].(string)
			if status != "done" {
				continue
			}

			symbol := m.Data["symbol"].(string)
			side := m.Data["side"].(string)
			size, _ := strconv.ParseFloat(m.Data["filledSize"].(string), 64)
			funds, _ := strconv.ParseFloat(m.Data["filledFunds"].(string), 64)

			switch step {
			case StepLeg1:
				if symbol == sym1 && side == "buy" {
					leg1Size = size
					leg1Funds = funds
					log.Printf("[FILL] LEG1 USDT→X size=%.8f funds=%.8f", size, funds)
					step = StepLeg2
				}
			case StepLeg2:
				if symbol == sym2 && side == "sell" {
					leg2Size = size
					leg2Funds = funds
					log.Printf("[FILL] LEG2 X→BTC size=%.8f funds=%.8f", size, funds)
					step = StepLeg3
				}
			case StepLeg3:
				if symbol == sym3 && side == "sell" {
					leg3Size = size
					leg3Funds = funds
					log.Printf("[FILL] LEG3 BTC→USDT size=%.8f funds=%.8f", size, funds)
					close(done)
				}
			}
		}
	}()

	// STEP 1
	step = StepLeg1
	o1, err := placeMarket(sym1, "buy", startUSDT)
	if err != nil {
		log.Println("[FAIL] LEG1 USDT→X", err)
		return
	}
	log.Printf("[OK] LEG1 orderId=%s", o1)

	// ждем WS событие о заполнении
	<-done

	// STEP 2
	step = StepLeg2
	o2, err := placeMarket(sym2, "sell", leg1Size)
	if err != nil {
		log.Println("[FAIL] LEG2 X→BTC", err)
		return
	}
	log.Printf("[OK] LEG2 orderId=%s", o2)

	<-done

	// STEP 3
	step = StepLeg3
	o3, err := placeMarket(sym3, "sell", leg2Funds)
	if err != nil {
		log.Println("[FAIL] LEG3 BTC→USDT", err)
		return
	}
	log.Printf("[OK] LEG3 orderId=%s", o3)

	<-done

	profit := leg3Funds - startUSDT
	pct := profit / startUSDT * 100
	log.Println("====== RESULT ======")
	log.Printf("START: %.4f USDT", startUSDT)
	log.Printf("END:   %.4f USDT", leg3Funds)
	log.Printf("PNL:   %.6f USDT (%.4f%%)", profit, pct)
} 


[{
	"resource": "/home/gaz358/myprog/crypt_proto/test/main.go",
	"owner": "_generated_diagnostic_collection_name_#0",
	"code": {
		"value": "UnusedVar",
		"target": {
			"$mid": 1,
			"path": "/golang.org/x/tools/internal/typesinternal",
			"scheme": "https",
			"authority": "pkg.go.dev",
			"fragment": "UnusedVar"
		}
	},
	"severity": 8,
	"message": "declared and not used: sig",
	"source": "compiler",
	"startLineNumber": 132,
	"startColumn": 2,
	"endLineNumber": 132,
	"endColumn": 5,
	"modelVersionId": 4,
	"tags": [
		1
	],
	"origin": "extHost1"
}]

[{
	"resource": "/home/gaz358/myprog/crypt_proto/test/main.go",
	"owner": "_generated_diagnostic_collection_name_#0",
	"code": {
		"value": "UnusedVar",
		"target": {
			"$mid": 1,
			"path": "/golang.org/x/tools/internal/typesinternal",
			"scheme": "https",
			"authority": "pkg.go.dev",
			"fragment": "UnusedVar"
		}
	},
	"severity": 8,
	"message": "declared and not used: leg2Size",
	"source": "compiler",
	"startLineNumber": 192,
	"startColumn": 16,
	"endLineNumber": 192,
	"endColumn": 24,
	"modelVersionId": 4,
	"origin": "extHost1"
}]

[{
	"resource": "/home/gaz358/myprog/crypt_proto/test/main.go",
	"owner": "_generated_diagnostic_collection_name_#0",
	"code": {
		"value": "UnusedVar",
		"target": {
			"$mid": 1,
			"path": "/golang.org/x/tools/internal/typesinternal",
			"scheme": "https",
			"authority": "pkg.go.dev",
			"fragment": "UnusedVar"
		}
	},
	"severity": 8,
	"message": "declared and not used: leg3Size",
	"source": "compiler",
	"startLineNumber": 192,
	"startColumn": 26,
	"endLineNumber": 192,
	"endColumn": 34,
	"modelVersionId": 4,
	"origin": "extHost1"
}]

[{
	"resource": "/home/gaz358/myprog/crypt_proto/test/main.go",
	"owner": "_generated_diagnostic_collection_name_#0",
	"code": {
		"value": "UnusedVar",
		"target": {
			"$mid": 1,
			"path": "/golang.org/x/tools/internal/typesinternal",
			"scheme": "https",
			"authority": "pkg.go.dev",
			"fragment": "UnusedVar"
		}
	},
	"severity": 8,
	"message": "declared and not used: leg1Funds",
	"source": "compiler",
	"startLineNumber": 193,
	"startColumn": 6,
	"endLineNumber": 193,
	"endColumn": 15,
	"modelVersionId": 4,
	"origin": "extHost1"
}]




