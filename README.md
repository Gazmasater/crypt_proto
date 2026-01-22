Если оставить только нужное:

p99 execution latency
Micro-volatility (100 мс)
Fill ratio
Capture rate
Inventory drift




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





func (ws *kucoinWS) handle(c *KuCoinCollector, msg []byte) {
	if gjson.GetBytes(msg, "type").String() != "message" {
		return
	}

	topic := gjson.GetBytes(msg, "topic").String()
	if !strings.HasPrefix(topic, "/market/ticker:") {
		return
	}
	symbol := strings.TrimPrefix(topic, "/market/ticker:")

	data := gjson.GetBytes(msg, "data")
	if !data.Exists() {
		return
	}

	bidStr := data.Get("bestBid").String()
	askStr := data.Get("bestAsk").String()
	if bidStr == "" || askStr == "" {
		return
	}

	bid, err1 := strconv.ParseFloat(bidStr, 64)
	ask, err2 := strconv.ParseFloat(askStr, 64)
	if err1 != nil || err2 != nil || bid <= 0 || ask <= 0 {
		return
	}

	last := ws.last[symbol]
	if last[0] == bid && last[1] == ask {
		return
	}
	ws.last[symbol] = [2]float64{bid, ask}

	bidSize, _ := strconv.ParseFloat(data.Get("bestBidSize").String(), 64)
	askSize, _ := strconv.ParseFloat(data.Get("bestAskSize").String(), 64)

	select {
	case c.out <- &models.MarketData{
		Exchange:  "KuCoin",
		Symbol:    symbol,
		Bid:       bid,
		Ask:       ask,
		BidSize:   bidSize,
		AskSize:   askSize,
		Timestamp: time.Now().UnixMilli(),
	}:
	default:
	}
}







