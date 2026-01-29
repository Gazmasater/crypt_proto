package collector_test

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
