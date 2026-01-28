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
	// проверка типа
	if !bytes.Contains(msg, []byte(`"type":"message"`)) {
		return
	}

	// проверка топика
	const prefix = `/market/ticker:`
	topicIdx := bytes.Index(msg, []byte(`"topic":"`))
	if topicIdx == -1 {
		return
	}
	topicStart := topicIdx + len(`"topic":"`)
	topicEnd := bytes.IndexByte(msg[topicStart:], '"')
	if topicEnd == -1 {
		return
	}
	topic := string(msg[topicStart : topicStart+topicEnd])
	if !strings.HasPrefix(topic, prefix) {
		return
	}
	symbol := strings.TrimPrefix(topic, prefix)

	// парс чисел
	bid := parseFloat(msg, `"bestBid":`)
	ask := parseFloat(msg, `"bestAsk":`)
	bidSize := parseFloat(msg, `"bestBidSize":`)
	askSize := parseFloat(msg, `"bestAskSize":`)

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

// parseFloat ищет число после ключа, например `"bestBid":123.45`
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


[{
	"resource": "/home/gaz358/myprog/crypt_proto/internal/collector/kucoin_collector.go",
	"owner": "_generated_diagnostic_collection_name_#1",
	"code": {
		"value": "IncompatibleAssign",
		"target": {
			"$mid": 1,
			"path": "/golang.org/x/tools/internal/typesinternal",
			"scheme": "https",
			"authority": "pkg.go.dev",
			"fragment": "IncompatibleAssign"
		}
	},
	"severity": 8,
	"message": "cannot use `\"bestBid\":` (untyped string constant \"\\\"bestBid\\\":\") as []byte value in argument to parseFloat",
	"source": "compiler",
	"startLineNumber": 188,
	"startColumn": 25,
	"endLineNumber": 188,
	"endColumn": 37,
	"modelVersionId": 3,
	"origin": "extHost1"
}]


