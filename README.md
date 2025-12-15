mx0vglmT3srN1IS19H
135bb7a7509e4421bad692415c53753b



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




func (t *Trader) PlaceMarket(
	ctx context.Context,
	symbol string,
	side string,
	quantity float64,
) error {
	if quantity <= 0 {
		return fmt.Errorf("quantity must be > 0, got %f", quantity)
	}
	side = strings.ToUpper(strings.TrimSpace(side))
	if side != "BUY" && side != "SELL" {
		return fmt.Errorf("invalid side %q", side)
	}

	endpoint := "/api/v3/order"

	// 1) Собираем параметры
	params := url.Values{}
	params.Set("symbol", symbol)
	params.Set("side", side)
	params.Set("type", "MARKET")
	params.Set("quantity", fmt.Sprintf("%.8f", quantity))
	params.Set("recvWindow", "5000")
	params.Set("timestamp", fmt.Sprintf("%d", time.Now().UnixMilli()))

	// 2) Подпись по queryString БЕЗ signature
	queryString := params.Encode()
	signature := t.sign(queryString)

	// 3) Добавляем signature как параметр
	params.Set("signature", signature)

	// 4) Формируем полный URL с query
	fullURL := t.baseURL + endpoint + "?" + params.Encode()

	// 5) POST без тела
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fullURL,
		nil,
	)
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}

	req.Header.Set("X-MEXC-APIKEY", t.apiKey)
	// тело пустое, Content-Type можно не ставить вообще
	// req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	t.dlog("PlaceMarket %s %s qty=%.8f url=%s", side, symbol, quantity, fullURL)

	resp, err := t.client.Do(req)
	if err != nil {
		return fmt.Errorf("http do: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("mexc order error: status=%d body=%s", resp.StatusCode, string(respBody))
	}

	t.dlog("PlaceMarket OK: %s", string(respBody))
	return nil
}





