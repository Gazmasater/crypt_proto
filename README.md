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



gjson.GetBytes(msg, "data.bestBid")
gjson.GetBytes(msg, "data.bestAsk")


data := gjson.GetBytes(msg, "data")
bid := data.Get("bestBid").Float()
ask := data.Get("bestAsk").Float()

ROUTINE ======================== crypt_proto/internal/collector.(*kucoinWS).handle in /home/gaz358/myprog/crypt_proto/internal/collector/kucoin_collector.go
         0      100ms (flat, cum) 12.35% of Total
         .          .    177:func (ws *kucoinWS) handle(c *KuCoinCollector, msg []byte) {
         .          .    178:   //if gjson.GetBytes(msg, "type").String() != "message" {
         .          .    179:   //      return
         .          .    180:   //}
         .          .    181:
         .       20ms    182:   topic := gjson.GetBytes(msg, "topic").String()
         .          .    183:   const prefix = "/market/ticker:"
         .          .    184:   if len(topic) <= len(prefix) || topic[:len(prefix)] != prefix {
         .          .    185:           return
         .          .    186:   }
         .          .    187:   symbol := topic[len(prefix):]
         .          .    188:
         .       30ms    189:   data := gjson.GetBytes(msg, "data")
         .       30ms    190:   bid := data.Get("bestBid").Float()
         .          .    191:   ask := data.Get("bestAsk").Float()
         .          .    192:   if bid == 0 || ask == 0 {
         .          .    193:           return
         .          .    194:   }
         .          .    195:
         .          .    196:   last := ws.last[symbol]
         .          .    197:   if last.Bid == bid && last.Ask == ask {
         .          .    198:           return
         .          .    199:   }
         .          .    200:
         .          .    201:   // если реально нужны
         .          .    202:   bidSize := gjson.GetBytes(msg, "data.bestBidSize").Float()
         .          .    203:   askSize := gjson.GetBytes(msg, "data.bestAskSize").Float()
         .          .    204:
         .       10ms    205:   ws.last[symbol] = Last{Bid: bid, Ask: ask}
         .          .    206:
         .       10ms    207:   c.out <- &models.MarketData{
         .          .    208:           Exchange: "KuCoin",
         .          .    209:           Symbol:   symbol,
         .          .    210:           Bid:      bid,
         .          .    211:           Ask:      ask,
         .          .    212:           BidSize:  bidSize,


