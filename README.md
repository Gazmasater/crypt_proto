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



func roundDown(v, step float64) float64 {
	return math.Floor(v/step) * step
}



func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	startUSDT := 12.0
	log.Printf("START %.2f USDT", startUSDT)

	// -------- LEG 1: USDT → DASH --------
	_, err := placeMarket("DASH-USDT", "buy", startUSDT)
	if err != nil {
		log.Fatal("LEG1 failed:", err)
	}

	time.Sleep(50 * time.Millisecond)

	dash, _ := getBalance("DASH")
	log.Printf("DASH balance: %.8f", dash)

	if dash <= 0 {
		log.Fatal("DASH not received")
	}

	// -------- LEG 2: DASH → BTC --------
	dash = roundDown(dash, 0.001)

	_, err = placeMarket("DASH-BTC", "sell", dash)
	if err != nil {
		log.Fatal("LEG2 failed:", err)
	}

	time.Sleep(50 * time.Millisecond)

	btc, _ := getBalance("BTC")
	log.Printf("BTC balance: %.8f", btc)

	// -------- LEG 3: BTC → USDT --------
	btc = roundDown(btc, 0.0001)

	_, err = placeMarket("BTC-USDT", "sell", btc)
	if err != nil {
		log.Fatal("LEG3 failed:", err)
	}

	time.Sleep(50 * time.Millisecond)

	usdt, _ := getBalance("USDT")
	log.Println("====== RESULT ======")
	log.Printf("END USDT: %.4f", usdt)
	log.Printf("PNL: %.4f", usdt-startUSDT)
}

