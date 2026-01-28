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
	// быстро проверяем тип сообщения
	if !bytes.Contains(msg, []byte(`"type":"message"`)) {
		return
	}

	// ищем topic
	const prefix = `/market/ticker:`
	topic := parseString(msg, `"topic":"`)
	if topic == "" || !strings.HasPrefix(topic, prefix) {
		return
	}
	symbol := strings.TrimSpace(topic[len(prefix):])

	// парсим числа
	bid := parseNumber(msg, `"bestBid":`)
	ask := parseNumber(msg, `"bestAsk":`)
	bidSize := parseNumber(msg, `"bestBidSize":`)
	askSize := parseNumber(msg, `"bestAskSize":`)

	if bid == 0 || ask == 0 {
		return
	}

	// проверка изменений
	last := ws.last[symbol]
	if last[0] == bid && last[1] == ask {
		return
	}
	ws.last[symbol] = [2]float64{bid, ask}

	// формируем структуру
	md := &models.MarketData{
		Exchange:  "KuCoin",
		Symbol:    symbol,
		Bid:       bid,
		Ask:       ask,
		BidSize:   bidSize,
		AskSize:   askSize,
		Timestamp: time.Now().UnixMilli(),
	}

	// отправка в канал без блокировки
	select {
	case c.out <- md:
	default:
		log.Printf("[KuCoin WS %d] drop MarketData %s\n", ws.id, symbol)
	}
}

// parseString возвращает строку после ключа, до следующей кавычки
func parseString(msg []byte, key string) string {
	idx := bytes.Index(msg, []byte(key))
	if idx == -1 {
		return ""
	}
	start := idx + len(key)
	end := bytes.IndexByte(msg[start:], '"')
	if end == -1 {
		return ""
	}
	return string(msg[start : start+end])
}

// parseNumber возвращает число float64 после ключа
func parseNumber(msg []byte, key string) float64 {
	idx := bytes.Index(msg, []byte(key))
	if idx == -1 {
		return 0
	}
	start := idx + len(key)
	end := start
	for end < len(msg) && ((msg[end] >= '0' && msg[end] <= '9') || msg[end] == '.') {
		end++
	}
	if end == start {
		return 0
	}
	val, _ := strconv.ParseFloat(string(msg[start:end]), 64)
	return val
}



