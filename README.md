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



az358@gaz358-BOD-WXX9:~/myprog/crypt_proto/test$ go run .
2026/01/16 02:02:04.251466 START TRIANGLE 11.00 USDT
panic: runtime error: index out of range [0] with length 0

goroutine 1 [running]:
main.connectPrivateWS()
        /home/gaz358/myprog/crypt_proto/test/main.go:185 +0x3d2
main.main()
        /home/gaz358/myprog/crypt_proto/test/main.go:201 +0xea
exit status 2


type Leg string

const (
	Leg1 Leg = "LEG1 USDT→DASH"
	Leg2 Leg = "LEG2 DASH→BTC"
	Leg3 Leg = "LEG3 BTC→USDT"
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
			Size  string `json:"size"`
			Funds string `json:"funds"`
			Fee   string `json:"fee"`
		} `json:"items"`
	} `json:"data"`
}

/* ================= AUTH ================= */

func sign(ts, method, path, body string) string {
	msg := ts + method + path + body
	mac := hmac.New(sha256.New, []byte(apiSecret))
	mac.Write([]byte(msg))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func headers(method, path, body string) http.Header {
	ts := strconv.FormatInt(time.Now().UnixMilli(), 10)
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

/* ================= REST ================= */

func placeMarket(leg Leg, symbol, side string, value float64) (string, error) {
	body := map[string]any{
		"symbol":    symbol,
		"type":      "market",
		"side":      side,
		"clientOid": fmt.Sprintf("%d", time.Now().UnixNano()), // уникальный идентификатор
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
		log.Printf("[FAIL] %s HTTP error: %v", leg, err)
		return "", err
	}
	defer resp.Body.Close()

	raw, _ := io.ReadAll(resp.Body)
	var r OrderResp
	_ = json.Unmarshal(raw, &r)

	if r.Code != "200000" {
		log.Printf("[FAIL] %s %s %s value=%.8f resp=%s", leg, side, symbol, value, raw)
		return "", fmt.Errorf("order rejected")
	}

	log.Printf("[OK] %s %s %s orderId=%s", leg, side, symbol, r.Data.OrderId)
	return r.Data.OrderId, nil
}

func waitFill(leg Leg, orderID string) (float64, float64, error) {
	time.Sleep(200 * time.Millisecond)

	path := "/api/v1/fills?orderId=" + orderID
	req, _ := http.NewRequest("GET", baseURL+path, nil)
	req.Header = headers("GET", path, "")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("[FAIL] %s fills HTTP error: %v", leg, err)
		return 0, 0, err
	}
	defer resp.Body.Close()

	raw, _ := io.ReadAll(resp.Body)
	var r FillResp
	_ = json.Unmarshal(raw, &r)

	if len(r.Data.Items) == 0 {
		log.Printf("[FAIL] %s no fills orderId=%s resp=%s", leg, orderID, raw)
		return 0, 0, fmt.Errorf("no fills")
	}

	var size, funds float64
	for _, f := range r.Data.Items {
		s, _ := strconv.ParseFloat(f.Size, 64)
		fd, _ := strconv.ParseFloat(f.Funds, 64)
		size += s
		funds += fd
	}

	log.Printf("[FILL] %s size=%.8f funds=%.8f", leg, size, funds)
	return size, funds, nil
}

/* ================= PRIVATE WS ================= */

func connectPrivateWS() *websocket.Conn {
	// Получение токена KuCoin для приватного WS
	req, _ := http.NewRequest("POST", baseURL+"/api/v1/bullet-private", nil)
	_ = strconv.FormatInt(time.Now().UnixMilli(), 10)
	req.Header = headers("POST", "/api/v1/bullet-private", "")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal("Private WS token error:", err)
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

	wsURL := r.Data.InstanceServers[0].Endpoint + "?token=" + r.Data.Token
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		log.Fatal("Private WS connect error:", err)
	}

	log.Println("Private WS connected")
	return conn
}

/* ================= MAIN ================= */

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.Printf("START TRIANGLE %.2f USDT", startUSDT)

	ws := connectPrivateWS()
	defer ws.Close()

	// ---------- LEG 1 ----------
	o1, err := placeMarket(Leg1, sym1, "buy", startUSDT)
	if err != nil {
		log.Println("ABORT AFTER LEG1")
		return
	}

	xQty, _, err := waitFill(Leg1, o1)
	if err != nil {
		log.Println("ABORT AFTER LEG1 FILL")
		return
	}

	// ---------- LEG 2 ----------
	o2, err := placeMarket(Leg2, sym2, "sell", xQty)
	if err != nil {
		log.Println("ABORT AFTER LEG2")
		return
	}

	_, btcGot, err := waitFill(Leg2, o2)
	if err != nil {
		log.Println("ABORT AFTER LEG2 FILL")
		return
	}

	// ---------- LEG 3 ----------
	o3, err := placeMarket(Leg3, sym3, "sell", btcGot)
	if err != nil {
		log.Println("ABORT AFTER LEG3")
		return
	}

	_, usdtFinal, err := waitFill(Leg3, o3)
	if err != nil {
		log.Println("ABORT AFTER LEG3 FILL")
		return
	}

	// ---------- RESULT ----------
	profit := usdtFinal - startUSDT
	pct := profit / startUSDT * 100

	log.Println("====== RESULT ======")
	log.Printf("START: %.4f USDT", startUSDT)
	log.Printf("END:   %.4f USDT", usdtFinal)
	log.Printf("PNL:   %.6f USDT (%.4f%%)", profit, pct)
}

