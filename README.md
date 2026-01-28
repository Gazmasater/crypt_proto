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




go run -race main.go


GOMAXPROCS=8 go run -race main.go



func (ws *kucoinWS) handle(c *KuCoinCollector, msg []byte) {
	// проверяем, что тип "message"
	if !bytes.Contains(msg, []byte(`"type":"message"`)) {
		return
	}

	// получаем symbol после "/market/ticker:"
	topicKey := []byte(`"topic":"/market/ticker:`)
	topicIdx := bytes.Index(msg, topicKey)
	if topicIdx < 0 {
		return
	}
	symbolStart := topicIdx + len(topicKey)
	symbolEnd := bytes.IndexByte(msg[symbolStart:], '"')
	if symbolEnd < 0 {
		return
	}
	symbol := string(msg[symbolStart : symbolStart+symbolEnd])

	// получаем данные "data"
	dataKey := []byte(`"data":{`)
	dataIdx := bytes.Index(msg, dataKey)
	if dataIdx < 0 {
		return
	}
	dataStart := dataIdx + len(dataKey)
	dataEnd := bytes.IndexByte(msg[dataStart:], '}')
	if dataEnd < 0 {
		return
	}
	data := msg[dataStart : dataStart+dataEnd]

	// парсим числа
	bid := parseFloat(data, "bestBid")
	ask := parseFloat(data, "bestAsk")
	bidSize := parseFloat(data, "bestBidSize")
	askSize := parseFloat(data, "bestAskSize")

	if bid == 0 || ask == 0 {
		return
	}

	last := ws.last[symbol]
	if last[0] == bid && last[1] == ask {
		return
	}

	ws.last[symbol] = [2]float64{bid, ask}

	c.out <- &models.MarketData{
		Exchange:  "KuCoin",
		Symbol:    symbol,
		Bid:       bid,
		Ask:       ask,
		BidSize:   bidSize,
		AskSize:   askSize,
		Timestamp: time.Now().UnixMilli(),
	}
}

// минимальный быстрый парсер float из []byte
func parseFloat(data []byte, key string) float64 {
	keyBytes := []byte(`"` + key + `":`)
	idx := bytes.Index(data, keyBytes)
	if idx < 0 {
		return 0
	}
	start := idx + len(keyBytes)
	// ищем конец числа (',' или конец данных)
	end := start
	for end < len(data) && (data[end] == '.' || data[end] == '-' || (data[end] >= '0' && data[end] <= '9')) {
		end++
	}
	f, _ := strconv.ParseFloat(string(data[start:end]), 64)
	return f
}

