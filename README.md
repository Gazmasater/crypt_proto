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



package collector

import (
	"testing"

	"github.com/tidwall/gjson"
)

// структуры, как в collector.go
type Last struct {
	Bid float64
	Ask float64
}

type kucoinWS struct {
	last map[string]Last
}

// пример сообщения от KuCoin
var sampleMsg = []byte(`{
	"topic": "/market/ticker:BTC-USDT",
	"data": {
		"bestBid": 30000.5,
		"bestAsk": 30001.5,
		"bestBidSize": 0.123,
		"bestAskSize": 0.456
	}
}`)

func BenchmarkHandleOld(b *testing.B) {
	ws := &kucoinWS{last: make(map[string]Last)}
	for i := 0; i < b.N; i++ {
		symbol := "BTC-USDT"
		bid := gjson.GetBytes(sampleMsg, "data.bestBid").Float()
		ask := gjson.GetBytes(sampleMsg, "data.bestAsk").Float()
		bidSize := gjson.GetBytes(sampleMsg, "data.bestBidSize").Float()
		askSize := gjson.GetBytes(sampleMsg, "data.bestAskSize").Float()

		last := ws.last[symbol]
		if last.Bid == bid && last.Ask == ask {
			continue
		}
		ws.last[symbol] = Last{Bid: bid, Ask: ask}

		_ = bidSize
		_ = askSize
	}
}

func BenchmarkHandleWithData(b *testing.B) {
	ws := &kucoinWS{last: make(map[string]Last)}
	for i := 0; i < b.N; i++ {
		symbol := "BTC-USDT"
		root := gjson.ParseBytes(sampleMsg)
		data := root.Get("data")
		bid := data.Get("bestBid").Float()
		ask := data.Get("bestAsk").Float()
		bidSize := data.Get("bestBidSize").Float()
		askSize := data.Get("bestAskSize").Float()

		last := ws.last[symbol]
		if last.Bid == bid && last.Ask == ask {
			continue
		}
		ws.last[symbol] = Last{Bid: bid, Ask: ask}

		_ = bidSize
		_ = askSize
	}
}

func BenchmarkHandleGetMany(b *testing.B) {
	ws := &kucoinWS{last: make(map[string]Last)}
	for i := 0; i < b.N; i++ {
		symbol := "BTC-USDT"
		values := gjson.GetManyBytes(sampleMsg,
			"data.bestBid",
			"data.bestAsk",
			"data.bestBidSize",
			"data.bestAskSize",
		)
		bid := values[0].Float()
		ask := values[1].Float()
		bidSize := values[2].Float()
		askSize := values[3].Float()

		last := ws.last[symbol]
		if last.Bid == bid && last.Ask == ask {
			continue
		}
		ws.last[symbol] = Last{Bid: bid, Ask: ask}

		_ = bidSize
		_ = askSize
	}
}




go test -bench=. ./internal/collector/collector_test



