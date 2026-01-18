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





const (
	apiKey        = "XXX"
	apiSecret     = "YYY"
	apiPassphrase = "ZZZ"

	baseURL   = "https://api.kucoin.com"
	startUSDT = 12.0

	sym1 = "DASH-USDT"
	sym2 = "DASH-BTC"
	sym3 = "BTC-USDT"

	step1 = 0.0001
	step2 = 0.0001
	step3 = 0.00001

	orderTimeout = 3 * time.Second
)

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

	json.NewDecoder(resp.Body).Decode(&r)
	if r.Code != "200000" {
		return "", fmt.Errorf("order rejected: %s", r.Msg)
	}
	return r.Data.OrderId, nil
}



func getPrivateWS() string {
	req, _ := http.NewRequest("POST", baseURL+"/api/v1/bullet-private", nil)
	req.Header = headers("POST", "/api/v1/bullet-private", "")

	resp, _ := http.DefaultClient.Do(req)
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
	json.NewDecoder(resp.Body).Decode(&r)

	return r.Data.InstanceServers[0].Endpoint + "?token=" + r.Data.Token
}

func wsLogin(conn *websocket.Conn) {
	ts := strconv.FormatInt(time.Now().UnixMilli(), 10)
	signStr := ts + "GET" + "/users/self/verify"

	mac := hmac.New(sha256.New, []byte(apiSecret))
	mac.Write([]byte(signStr))
	sign := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	login := map[string]string{
		"id":         "login",
		"type":       "login",
		"apiKey":     apiKey,
		"passphrase": passphrase(),
		"timestamp":  ts,
		"sign":       sign,
	}

	raw, _ := json.Marshal(login)
	conn.WriteMessage(websocket.TextMessage, raw)

	_, msg, _ := conn.ReadMessage()
	if !bytes.Contains(msg, []byte(`"success":true`)) {
		log.Fatal("WS login failed:", string(msg))
	}
}



type OrderEvent struct {
	OrderID string
	Filled  float64
}

func wsReader(conn *websocket.Conn, out chan<- OrderEvent) {
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Fatal("WS read error:", err)
		}

		var evt struct {
			Type string `json:"type"`
			Data struct {
				OrderId    string `json:"orderId"`
				FilledSize string `json:"filledSize"`
			} `json:"data"`
		}

		if json.Unmarshal(msg, &evt) != nil {
			continue
		}

		if evt.Type == "order.done" {
			val, _ := strconv.ParseFloat(evt.Data.FilledSize, 64)
			out <- OrderEvent{
				OrderID: evt.Data.OrderId,
				Filled:  val,
			}
		}
	}
}



func waitOrder(ctx context.Context, ch <-chan OrderEvent, orderID string, step float64) (float64, error) {
	for {
		select {
		case evt := <-ch:
			if evt.OrderID == orderID {
				return math.Floor(evt.Filled/step) * step, nil
			}
		case <-ctx.Done():
			return 0, fmt.Errorf("order %s timeout", orderID)
		}
	}
}


func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.Println("START TRIANGLE", startUSDT)

	wsURL := getPrivateWS()
	conn, _, _ := websocket.DefaultDialer.Dial(wsURL, nil)
	defer conn.Close()

	wsLogin(conn)

	sub := `{"id":"1","type":"subscribe","topic":"/spot/tradeOrders","privateChannel":true}`
	conn.WriteMessage(websocket.TextMessage, []byte(sub))

	events := make(chan OrderEvent, 16)
	go wsReader(conn, events)

	// LEG1
	o1, _ := placeMarket(sym1, "buy", startUSDT)
	ctx1, _ := context.WithTimeout(context.Background(), orderTimeout)
	dash, err := waitOrder(ctx1, events, o1, step1)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("LEG1 OK", dash)

	// LEG2
	o2, _ := placeMarket(sym2, "sell", dash)
	ctx2, _ := context.WithTimeout(context.Background(), orderTimeout)
	btc, err := waitOrder(ctx2, events, o2, step2)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("LEG2 OK", btc)

	// LEG3
	o3, _ := placeMarket(sym3, "sell", btc)
	ctx3, _ := context.WithTimeout(context.Background(), orderTimeout)
	usdt, err := waitOrder(ctx3, events, o3, step3)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("PNL:", usdt-startUSDT)
}








