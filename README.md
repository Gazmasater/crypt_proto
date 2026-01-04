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




func (c *KuCoinCollector) readLoop(out chan<- *models.MarketData) {
	defer func() {
		if c.conn != nil {
			_ = c.conn.Close()
		}
	}()

	log.Println("[KuCoin] readLoop started")

	for {
		select {
		case <-c.ctx.Done():
			log.Println("[KuCoin] context cancelled")
			return
		default:
		}

		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("[KuCoin] read error:", err)
			return
		}

		// 🔍 RAW LOG
		log.Printf("[KuCoin RAW] %s\n", msg)

		var raw map[string]any
		if err := json.Unmarshal(msg, &raw); err != nil {
			log.Println("[KuCoin] json unmarshal error:", err)
			continue
		}

		typ, _ := raw["type"].(string)
		if typ != "message" {
			log.Println("[KuCoin] skip type:", typ)
			continue
		}

		topic, _ := raw["topic"].(string)
		if !strings.HasPrefix(topic, "/market/ticker:") {
			log.Println("[KuCoin] skip topic:", topic)
			continue
		}

		data, ok := raw["data"].(map[string]any)
		if !ok {
			log.Println("[KuCoin] data is not map")
			continue
		}

		rawSymbol := strings.TrimPrefix(topic, "/market/ticker:")
		symbol := market.NormalizeSymbol_NoAlloc(rawSymbol, &c.buf)

		log.Printf("[KuCoin] symbol raw=%s normalized=%s\n", rawSymbol, symbol)

		if symbol == "" {
			log.Println("[KuCoin] empty symbol after normalize")
			continue
		}

		// whitelist
		if len(c.allowed) > 0 {
			if _, ok := c.allowed[symbol]; !ok {
				log.Println("[KuCoin] symbol not in whitelist:", symbol)
				continue
			}
		}

		bid := parseFloat(data["bestBid"])
		ask := parseFloat(data["bestAsk"])
		bidSize := parseFloat(data["sizeBid"])
		askSize := parseFloat(data["sizeAsk"])

		log.Printf(
			"[KuCoin] parsed %s bid=%f ask=%f bidSize=%f askSize=%f",
			symbol, bid, ask, bidSize, askSize,
		)

		if bid == 0 || ask == 0 {
			log.Println("[KuCoin] zero bid or ask")
			continue
		}

		// ⚠ дедупликация (ПО ЦЕНАМ)
		c.mu.Lock()
		last, exists := c.lastData[symbol]
		if exists && last.Bid == bid && last.Ask == ask {
			c.mu.Unlock()
			log.Println("[KuCoin] dedup price:", symbol)
			continue
		}

		c.lastData[symbol] = struct {
			Bid, Ask, BidSize, AskSize float64
		}{
			Bid: bid, Ask: ask,
			BidSize: bidSize, AskSize: askSize,
		}
		c.mu.Unlock()

		md := c.pool.Get().(*models.MarketData)
		md.Exchange = "KuCoin"
		md.Symbol = symbol
		md.Bid = bid
		md.Ask = ask
		md.BidSize = bidSize
		md.AskSize = askSize
		md.Timestamp = time.Now().UnixMilli()

		log.Printf("[KuCoin] PUSH %s bid=%f ask=%f\n", symbol, bid, ask)

		out <- md
	}
}

