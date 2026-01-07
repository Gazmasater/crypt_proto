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




func (ws *kucoinWS) handle(c *KuCoinCollector, msg []byte) {
	msgType := gjson.GetBytes(msg, "type")
	if msgType.String() != "message" {
		return
	}

	topic := gjson.GetBytes(msg, "topic").String()
	parts := strings.Split(topic, ":")
	if len(parts) != 2 {
		return
	}
	symbol := normalize(parts[1])

	data := gjson.GetBytes(msg, "data")
	bid := data.Get("bestBid").Float()
	ask := data.Get("bestAsk").Float()
	if bid == 0 || ask == 0 {
		return
	}

	// Проверяем изменения через методы lastData
	lastBid, lastAsk, ok := ws.last.Get(symbol)
	if ok && lastBid == bid && lastAsk == ask {
		return
	}
	ws.last.Set(symbol, bid, ask)

	c.out <- &models.MarketData{
		Exchange: "KuCoin",
		Symbol:   symbol,
		Bid:      bid,
		Ask:      ask,
	}
}


[{
	"resource": "/home/gaz358/myprog/crypt_proto/internal/collector/kucoin_collector.go",
	"owner": "_generated_diagnostic_collection_name_#0",
	"code": {
		"value": "MissingFieldOrMethod",
		"target": {
			"$mid": 1,
			"path": "/golang.org/x/tools/internal/typesinternal",
			"scheme": "https",
			"authority": "pkg.go.dev",
			"fragment": "MissingFieldOrMethod"
		}
	},
	"severity": 8,
	"message": "ws.last.Get undefined (type map[string]lastData has no field or method Get)",
	"source": "compiler",
	"startLineNumber": 217,
	"startColumn": 34,
	"endLineNumber": 217,
	"endColumn": 37,
	"origin": "extHost1"
}]

[{
	"resource": "/home/gaz358/myprog/crypt_proto/internal/collector/kucoin_collector.go",
	"owner": "_generated_diagnostic_collection_name_#0",
	"code": {
		"value": "MissingFieldOrMethod",
		"target": {
			"$mid": 1,
			"path": "/golang.org/x/tools/internal/typesinternal",
			"scheme": "https",
			"authority": "pkg.go.dev",
			"fragment": "MissingFieldOrMethod"
		}
	},
	"severity": 8,
	"message": "ws.last.Set undefined (type map[string]lastData has no field or method Set)",
	"source": "compiler",
	"startLineNumber": 221,
	"startColumn": 10,
	"endLineNumber": 221,
	"endColumn": 13,
	"origin": "extHost1"
}]





