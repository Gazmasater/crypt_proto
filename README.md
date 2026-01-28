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



package collector_test

import (
	"testing"

	"crypt_proto/internal/collector"
	"github.com/tidwall/gjson"
)

// имитируем структуру collector
type Last struct {
	Bid float64
	Ask float64
}

type kucoinWS struct {
	last map[string]Last
}

type KuCoinCollector struct {
	out chan<- interface{}
}

var testMsg = []byte(`{
	"topic": "/market/ticker:BTC-USDT",
	"data": {
		"bestBid": 50000.12,
		"bestAsk": 50001.34,
		"bestBidSize": 0.5,
		"bestAskSize": 0.6
	}
}`)

func BenchmarkHandleOld(b *testing.B) {
	ws := &kucoinWS{last: make(map[string]Last)}
	c := &KuCoinCollector{}

	for i := 0; i < b.N; i++ {
		// старый способ — отдельные вызовы GetBytes
		data := gjson.GetBytes(testMsg, "data")
		bid := data.Get("bestBid").Float()
		ask := data.Get("bestAsk").Float()
		bidSize := data.Get("bestBidSize").Float()
		askSize := data.Get("bestAskSize").Float()

		// обновление last
		ws.last["BTC-USDT"] = Last{Bid: bid, Ask: ask}

		// имитация вывода
		_ = bidSize
		_ = askSize
		_ = c
	}
}

func BenchmarkHandleMany(b *testing.B) {
	ws := &kucoinWS{last: make(map[string]Last)}
	c := &KuCoinCollector{}

	for i := 0; i < b.N; i++ {
		// новый способ — GetManyBytes
		values := gjson.GetManyBytes(testMsg,
			"data.bestBid",
			"data.bestAsk",
			"data.bestBidSize",
			"data.bestAskSize",
		)
		bid := values[0].Float()
		ask := values[1].Float()
		bidSize := values[2].Float()
		askSize := values[3].Float()

		ws.last["BTC-USDT"] = Last{Bid: bid, Ask: ask}

		// имитация вывода
		_ = bidSize
		_ = askSize
		_ = c
	}
}



go test -bench=. ./internal/collector


gaz358@gaz358-BOD-WXX9:~/myprog/crypt_proto$ go test -bench=. ./internal/collector/collector_test
goos: linux
goarch: amd64
pkg: crypt_proto/internal/collector/collector_test
cpu: 11th Gen Intel(R) Core(TM) i5-1135G7 @ 2.40GHz
BenchmarkHandleOld-8             1902133               647.3 ns/op
BenchmarkHandleMany-8            1283359               913.7 ns/op
PASS
ok      crypt_proto/internal/collector/collector_test   3.994s


