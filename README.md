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
	const prefixLen = len("/market/ticker:")

	topic := gjson.GetBytes(msg, "topic").String()
	if len(topic) <= prefixLen {
		return
	}
	symbol := topic[prefixLen:]

	bid := gjson.GetBytes(msg, "data.bestBid").Float()
	ask := gjson.GetBytes(msg, "data.bestAsk").Float()
	if bid == 0 || ask == 0 {
		return
	}

	if last, ok := ws.last[symbol]; ok && last.Bid == bid && last.Ask == ask {
		return
	}

	// сразу извлекаем объёмы
	c.out <- &models.MarketData{
		Exchange: "KuCoin",
		Symbol:   symbol,
		Bid:      bid,
		Ask:      ask,
		BidSize:  gjson.GetBytes(msg, "data.bestBidSize").Float(),
		AskSize:  gjson.GetBytes(msg, "data.bestAskSize").Float(),
	}

	ws.last[symbol] = Last{Bid: bid, Ask: ask}
}


ROUTINE ======================== crypt_proto/internal/collector.(*kucoinWS).handle in /home/gaz358/myprog/crypt_proto/internal/collector/kucoin_collector.go
      10ms       70ms (flat, cum)  8.97% of Total
         .          .    177:func (ws *kucoinWS) handle(c *KuCoinCollector, msg []byte) {
         .          .    178:   const prefixLen = len("/market/ticker:")
         .          .    179:
         .       20ms    180:   topic := gjson.GetBytes(msg, "topic").String()
         .          .    181:   if len(topic) <= prefixLen {
         .          .    182:           return
         .          .    183:   }
         .          .    184:   symbol := topic[prefixLen:]
         .          .    185:
         .       30ms    186:   bid := gjson.GetBytes(msg, "data.bestBid").Float()
         .          .    187:   ask := gjson.GetBytes(msg, "data.bestAsk").Float()
         .          .    188:   if bid == 0 || ask == 0 {
         .          .    189:           return
         .          .    190:   }
         .          .    191:
      10ms       10ms    192:   if last, ok := ws.last[symbol]; ok && last.Bid == bid && last.Ask == ask {
         .          .    193:           return
         .          .    194:   }
         .          .    195:
         .          .    196:   // сразу извлекаем объёмы
         .       10ms    197:   c.out <- &models.MarketData{
         .          .    198:           Exchange: "KuCoin",
         .          .    199:           Symbol:   symbol,
         .          .    200:           Bid:      bid,
         .          .    201:           Ask:      ask,
         .          .    202:           BidSize:  gjson.GetBytes(msg, "data.bestBidSize").Float(),


