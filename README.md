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




type MarketData struct {
	Exchange string
	Symbol   string

	Bid float64
	Ask float64

	BidSize float64
	AskSize float64
}



func (ws *kucoinWS) handle(c *KuCoinCollector, msg []byte) {
	// быстрый фильтр
	if gjson.GetBytes(msg, "type").String() != "message" {
		return
	}

	topic := gjson.GetBytes(msg, "topic").String()
	if !strings.HasPrefix(topic, "/market/ticker:") {
		return
	}

	symbol := normalize(strings.TrimPrefix(topic, "/market/ticker:"))

	data := gjson.GetBytes(msg, "data")

	bid := data.Get("bestBid").Float()
	ask := data.Get("bestAsk").Float()
	if bid == 0 || ask == 0 {
		return
	}

	// 🔥 ГЛАВНЫЙ CPU-ФИЛЬТР
	ws.mu.Lock()
	last := ws.last[symbol]
	if last[0] == bid && last[1] == ask {
		ws.mu.Unlock()
		return
	}
	ws.last[symbol] = [2]float64{bid, ask}
	ws.mu.Unlock()

	// объёмы парсим ТОЛЬКО если цена изменилась
	c.out <- &models.MarketData{
		Exchange: "KuCoin",
		Symbol:   symbol,
		Bid:      bid,
		Ask:      ask,
		BidSize:  data.Get("bestBidSize").Float(),
		AskSize:  data.Get("bestAskSize").Float(),
	}
}



package models

// MarketData хранит данные с биржи для одного инструмента
type MarketData struct {
	Exchange  string  `json:"exchange"`  // название биржи
	Symbol    string  `json:"symbol"`    // торговая пара, например BTC-USDT
	Bid       float64 `json:"bid"`       // лучшая цена покупки
	Ask       float64 `json:"ask"`       // лучшая цена продажи
	BidSize   float64 `json:"bid_size"`  // объём на лучшей цене покупки
	AskSize   float64 `json:"ask_size"`  // объём на лучшей цене продажи
	Timestamp int64   `json:"timestamp"` // метка времени в миллисекундах
}




