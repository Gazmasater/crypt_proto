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



func placeMarket(symbol, side string, value float64) (filled float64, err error) {
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
		return 0, err
	}
	defer resp.Body.Close()

	// читаем тело полностью для логирования
	respBytes, _ := io.ReadAll(resp.Body)
	
	var r struct {
		Code string `json:"code"`
		Msg  string `json:"msg"` // <- здесь сообщение об ошибке
		Data struct {
			DealFunds string `json:"dealFunds"`
			DealSize  string `json:"dealSize"`
		} `json:"data"`
	}

	if err := json.Unmarshal(respBytes, &r); err != nil {
		return 0, fmt.Errorf("decode error: %v; raw: %s", err, string(respBytes))
	}

	if r.Code != "200000" {
		return 0, fmt.Errorf("order rejected: %s; raw: %s", r.Msg, string(respBytes))
	}

	if side == "buy" {
		filled, _ = strconv.ParseFloat(r.Data.DealSize, 64)
	} else {
		filled, _ = strconv.ParseFloat(r.Data.DealFunds, 64)
	}
	return filled, nil
}

