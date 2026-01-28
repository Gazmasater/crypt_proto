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

	// достаём топик сразу
	topic := gjson.GetBytes(msg, "topic").String()
	if len(topic) <= prefixLen {
		return
	}
	symbol := topic[prefixLen:]

	// извлекаем bid, ask и размеры одним вызовом
	values := gjson.GetManyBytes(msg,
		"data.bestBid",
		"data.bestAsk",
		"data.bestBidSize",
		"data.bestAskSize",
	)
	bid := values[0].Float()
	ask := values[1].Float()
	bidSize := values[2].Float()
	askSize := values[3].Float()

	if bid == 0 || ask == 0 {
		return
	}

	// проверяем, изменились ли цены
	if last, ok := ws.last[symbol]; ok && last.Bid == bid && last.Ask == ask {
		return
	}

	// обновляем last
	ws.last[symbol] = Last{Bid: bid, Ask: ask}

	// отправляем в канал
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
         0       60ms (flat, cum)  7.14% of Total
         .          .    177:func (ws *kucoinWS) handle(c *KuCoinCollector, msg []byte) {
         .          .    178:   const prefixLen = len("/market/ticker:")
         .          .    179:
         .          .    180:   // достаём топик сразу
         .       30ms    181:   topic := gjson.GetBytes(msg, "topic").String()
         .          .    182:   if len(topic) <= prefixLen {
         .          .    183:           return
         .          .    184:   }
         .          .    185:   symbol := topic[prefixLen:]
         .          .    186:
         .          .    187:   // извлекаем bid, ask и размеры одним вызовом
         .       20ms    188:   values := gjson.GetManyBytes(msg,
         .          .    189:           "data.bestBid",
         .          .    190:           "data.bestAsk",
         .          .    191:           "data.bestBidSize",
         .          .    192:           "data.bestAskSize",
         .          .    193:   )
         .          .    194:   bid := values[0].Float()
         .          .    195:   ask := values[1].Float()
         .          .    196:   bidSize := values[2].Float()
         .          .    197:   askSize := values[3].Float()
         .          .    198:
         .          .    199:   if bid == 0 || ask == 0 {
         .          .    200:           return
         .          .    201:   }
         .          .    202:
         .          .    203:   // проверяем, изменились ли цены
         .       10ms    204:   if last, ok := ws.last[symbol]; ok && last.Bid == bid && last.Ask == ask {
         .          .    205:           return
         .          .    206:   }
         .          .    207:
         .          .    208:   // обновляем last
         .          .    209:   ws.last[symbol] = Last{Bid: bid, Ask: ask}


