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




func (c *KuCoinCollector) Connect() error {
	conn, _, err := websocket.DefaultDialer.Dial(c.wsURL, nil)
	if err != nil {
		return err
	}

	c.conn = conn
	log.Println("[KuCoin] Connected to WS")

	// ВАЖНО: подписка
	c.subscribe()

	go c.readLoop()
	return nil
}


func (c *KuCoinCollector) subscribe() {
	for _, symbol := range c.symbols {
		topic := "/market/ticker:" + symbol

		msg := map[string]interface{}{
			"id":              fmt.Sprintf("sub-%s", symbol),
			"type":            "subscribe",
			"topic":           topic,
			"privateChannel":  false,
			"response":        true,
		}

		b, _ := json.Marshal(msg)

		log.Printf("[KuCoin SUBSCRIBE] %s", topic)

		if err := c.conn.WriteMessage(websocket.TextMessage, b); err != nil {
			log.Printf("[KuCoin] subscribe error %s: %v", symbol, err)
		}
	}
}




