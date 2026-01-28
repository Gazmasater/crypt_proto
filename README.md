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
	// проверка типа сообщения
	if !bytes.Contains(msg, []byte(`"type":"message"`)) {
		return
	}

	// ищем topic
	const prefix = `/market/ticker:`
	topicKey := []byte(`"topic":"`)
	topicIdx := bytes.Index(msg, topicKey)
	if topicIdx == -1 {
		return
	}
	topicStart := topicIdx + len(topicKey)
	topicEnd := bytes.IndexByte(msg[topicStart:], '"')
	if topicEnd == -1 {
		return
	}
	topic := msg[topicStart : topicStart+topicEnd]
	if !bytes.HasPrefix(topic, []byte(prefix)) {
		return
	}
	symbol := string(topic[len(prefix):])

	// ищем поле "data": {...}
	dataKey := []byte(`"data":`)
	dataIdx := bytes.Index(msg, dataKey)
	if dataIdx == -1 {
		return
	}
	dataStart := dataIdx + len(dataKey)

	// парсим числа из блока data
	bid := parseFloat(msg[dataStart:], []byte(`"bestBid":`))
	ask := parseFloat(msg[dataStart:], []byte(`"bestAsk":`))
	bidSize := parseFloat(msg[dataStart:], []byte(`"bestBidSize":`))
	askSize := parseFloat(msg[dataStart:], []byte(`"bestAskSize":`))

	if bid == 0 || ask == 0 {
		return
	}

	// проверка на изменения
	last := ws.last[symbol]
	if last[0] == bid && last[1] == ask {
		return
	}
	ws.last[symbol] = [2]float64{bid, ask}

	// отправляем в канал
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

// parseFloat ищет число после ключа, например `"bestBid":123.45`
// работает с []byte, без аллокаций кроме strconv.ParseFloat
func parseFloat(msg []byte, key []byte) float64 {
	idx := bytes.Index(msg, key)
	if idx == -1 {
		return 0
	}
	start := idx + len(key)

	// ищем конец числа — любой символ, который не цифра и не точка
	end := start
	for end < len(msg) && ((msg[end] >= '0' && msg[end] <= '9') || msg[end] == '.') {
		end++
	}
	if end == start {
		return 0
	}

	val, err := strconv.ParseFloat(string(msg[start:end]), 64)
	if err != nil {
		return 0
	}
	return val
}



