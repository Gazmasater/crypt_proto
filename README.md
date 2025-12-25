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





func (c *KuCoinCollector) readLoop(out chan<- models.MarketData) {
	defer func() {
		if c.conn != nil {
			c.conn.Close()
		}
	}()

	// Храним последние данные по каждому символу
	lastData := make(map[string]struct {
		Bid, Ask, BidSize, AskSize float64
	})

	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			_, msg, err := c.conn.ReadMessage()
			if err != nil {
				log.Println("[KuCoin] read error:", err)
				return
			}

			var raw map[string]any
			if err := json.Unmarshal(msg, &raw); err != nil {
				continue
			}

			typ, _ := raw["type"].(string)
			if typ == "welcome" || typ == "ack" || typ != "message" {
				continue
			}

			topic, _ := raw["topic"].(string)
			data, ok := raw["data"].(map[string]any)
			if !ok {
				continue
			}

			rawsymbol := strings.TrimPrefix(topic, "/market/ticker:")
			symbol := market.NormalizeSymbol_Full(rawsymbol)

			bid := parseFloat(data["bestBid"])
			ask := parseFloat(data["bestAsk"])
			bidSize := parseFloat(data["sizeBid"])
			askSize := parseFloat(data["sizeAsk"])

			if bid == 0 || ask == 0 {
				continue
			}

			// фильтрация повторов
			if last, exists := lastData[symbol]; exists {
				if last.Bid == bid && last.Ask == ask && last.BidSize == bidSize && last.AskSize == askSize {
					continue // ничего не поменялось — пропускаем
				}
			}

			// обновляем последние данные
			lastData[symbol] = struct {
				Bid, Ask, BidSize, AskSize float64
			}{Bid: bid, Ask: ask, BidSize: bidSize, AskSize: askSize}

			out <- models.MarketData{
				Exchange: "KuCoin",
				Symbol:   symbol,
				Bid:      bid,
				Ask:      ask,
				// при желании можно добавить объёмы в MarketData
			}
		}
	}
}



