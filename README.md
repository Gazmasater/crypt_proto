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
	root := gjson.ParseBytes(msg)

	topic := root.Get("topic").String()
	const prefix = "/market/ticker:"
	if len(topic) <= len(prefix) || topic[:len(prefix)] != prefix {
		return
	}
	symbol := topic[len(prefix):]

	data := root.Get("data")
	bid := data.Get("bestBid").Float()
	ask := data.Get("bestAsk").Float()
	if bid == 0 || ask == 0 {
		return
	}

	last := ws.last[symbol]
	if last.Bid == bid && last.Ask == ask {
		return
	}

	// если реально нужны
	bidSize := data.Get("bestBidSize").Float()
	askSize := data.Get("bestAskSize").Float()

	ws.last[symbol] = Last{Bid: bid, Ask: ask}

	c.out <- &models.MarketData{
		Exchange: "KuCoin",
		Symbol:   symbol,
		Bid:      bid,
		Ask:      ask,
		BidSize:  bidSize,
		AskSize:  askSize,
	}
}



ROUTINE ======================== crypt_proto/internal/collector.(*kucoinWS).handle in /home/gaz358/myprog/crypt_proto/internal/collector/kucoin_collector.go
         0      150ms (flat, cum) 15.00% of Total
         .          .    177:func (ws *kucoinWS) handle(c *KuCoinCollector, msg []byte) {
         .          .    178:   root := gjson.ParseBytes(msg)
         .          .    179:
         .       50ms    180:   topic := root.Get("topic").String()
         .          .    181:   const prefix = "/market/ticker:"
         .          .    182:   if len(topic) <= len(prefix) || topic[:len(prefix)] != prefix {
         .          .    183:           return
         .          .    184:   }
         .          .    185:   symbol := topic[len(prefix):]
         .          .    186:
         .       10ms    187:   data := root.Get("data")
         .       30ms    188:   bid := data.Get("bestBid").Float()
         .       20ms    189:   ask := data.Get("bestAsk").Float()
         .          .    190:   if bid == 0 || ask == 0 {
         .          .    191:           return
         .          .    192:   }
         .          .    193:
         .       10ms    194:   last := ws.last[symbol]
         .          .    195:   if last.Bid == bid && last.Ask == ask {
         .          .    196:           return
         .          .    197:   }
         .          .    198:
         .          .    199:   bidSize := data.Get("bestBidSize").Float()
         .       10ms    200:   askSize := data.Get("bestAskSize").Float()
         .          .    201:
         .          .    202:   ws.last[symbol] = Last{Bid: bid, Ask: ask}
         .          .    203:
         .       20ms    204:   c.out <- &models.MarketData{
         .          .    205:           Exchange: "KuCoin",
         .          .    206:           Symbol:   symbol,
         .          .    207:           Bid:      bid,
         .          .    208:           Ask:      ask,
         .          .    209:           BidSize:  bidSize,


