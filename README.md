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
	if gjson.GetBytes(msg, "type").String() != "message" {
		return
	}

	topic := gjson.GetBytes(msg, "topic").String()
	const prefix = "/market/ticker:"
	if !strings.HasPrefix(topic, prefix) {
		return
	}
	symbol := strings.TrimPrefix(topic, prefix)

	data := gjson.GetBytes(msg, "data")
	bid := data.Get("bestBid").Float()
	ask := data.Get("bestAsk").Float()
	if bid == 0 || ask == 0 {
		return
	}

	// проверяем изменения
	if last, ok := ws.last[symbol]; ok && last[0] == bid && last[1] == ask {
		return
	}

	bidSize := data.Get("bestBidSize").Float()
	askSize := data.Get("bestAskSize").Float()

	// обновляем last
	ws.last[symbol] = [2]float64{bid, ask}

	// отправляем в канал
	c.out <- &models.MarketData{
		Exchange:  "KuCoin",
		Symbol:    symbol,
		Bid:       bid,
		Ask:       ask,
		BidSize:   bidSize,
		AskSize:   askSize,
		Timestamp: time.Now().UnixMilli(),
	}

	// опционально логирование
	// log.Printf("[KuCoin WS %d] %s bid=%.6f ask=%.6f", ws.id, symbol, bid, ask)
}




gaz358@gaz358-BOD-WXX9:~/myprog/crypt_proto$    go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
Fetching profile over HTTP from http://localhost:6060/debug/pprof/profile?seconds=30
Saved profile in /home/gaz358/pprof/pprof.arb.samples.cpu.303.pb.gz
File: arb
Build ID: 82b5c9a24851fe4609ad4810d06569374d167430
Type: cpu
Time: 2026-01-28 19:53:09 MSK
Duration: 30s, Total samples = 1.10s ( 3.67%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) list crypt_proto
Total: 1.10s
ROUTINE ======================== crypt_proto/internal/collector.(*kucoinWS).handle in /home/gaz358/myprog/crypt_proto/internal/collector/kucoin_collector.go
         0      140ms (flat, cum) 12.73% of Total
         .          .    163:func (ws *kucoinWS) handle(c *KuCoinCollector, msg []byte) {
         .       50ms    164:   if gjson.GetBytes(msg, "type").String() != "message" {
         .          .    165:           return
         .          .    166:   }
         .          .    167:
         .          .    168:   topic := gjson.GetBytes(msg, "topic").String()
         .          .    169:   const prefix = "/market/ticker:"
         .          .    170:   if !strings.HasPrefix(topic, prefix) {
         .          .    171:           return
         .          .    172:   }
         .          .    173:   symbol := strings.TrimPrefix(topic, prefix)
         .          .    174:
         .       50ms    175:   data := gjson.GetBytes(msg, "data")
         .       10ms    176:   bid := data.Get("bestBid").Float()
         .          .    177:   ask := data.Get("bestAsk").Float()
         .          .    178:   if bid == 0 || ask == 0 {
         .          .    179:           return
         .          .    180:   }
         .          .    181:
         .          .    182:   // проверяем изменения
         .          .    183:   if last, ok := ws.last[symbol]; ok && last[0] == bid && last[1] == ask {
         .          .    184:           return
         .          .    185:   }
         .          .    186:
         .          .    187:   bidSize := data.Get("bestBidSize").Float()
         .          .    188:   askSize := data.Get("bestAskSize").Float()
         .          .    189:
         .          .    190:   // обновляем last
         .          .    191:   ws.last[symbol] = [2]float64{bid, ask}
         .          .    192:
         .          .    193:   // отправляем в канал
         .       20ms    194:   c.out <- &models.MarketData{
         .          .    195:           Exchange:  "KuCoin",
         .          .    196:           Symbol:    symbol,
         .          .    197:           Bid:       bid,
         .          .    198:           Ask:       ask,
         .          .    199:           BidSize:   bidSize,
         .          .    200:           AskSize:   askSize,
         .       10ms    201:           Timestamp: time.Now().UnixMilli(),
         .          .    202:   }
         .          .    203:
         .          .    204:   // опционально логирование
         .          .    205:   // log.Printf("[KuCoin WS %d] %s bid=%.6f ask=%.6f", ws.id, symbol, bid, ask)
         .          .    206:}
ROUTINE ======================== crypt_proto/internal/collector.(*kucoinWS).readLoop in /home/gaz358/myprog/crypt_proto/internal/collector/kucoin_collector.go
      10ms      740ms (flat, cum) 67.27% of Total
         .          .    152:func (ws *kucoinWS) readLoop(c *KuCoinCollector) {
         .          .    153:   for {
         .      590ms    154:           _, msg, err := ws.conn.ReadMessage()
      10ms       10ms    155:           if err != nil {
         .          .    156:                   log.Printf("[KuCoin WS %d] read error: %v\n", ws.id, err)
         .          .    157:                   return
         .          .    158:           }
         .      140ms    159:           ws.handle(c, msg)
         .          .    160:   }
         .          .    161:}
         .          .    162:
         .          .    163:func (ws *kucoinWS) handle(c *KuCoinCollector, msg []byte) {
         .          .    164:   if gjson.GetBytes(msg, "type").String() != "message" {
ROUTINE ======================== crypt_proto/internal/queue.(*MemoryStore).Run in /home/gaz358/myprog/crypt_proto/internal/queue/in_memory_queue.go
         0       30ms (flat, cum)  2.73% of Total
         .          .     30:func (s *MemoryStore) Run() {
         .          .     31:   for md := range s.batch {
         .       30ms     32:           s.apply(md)
         .          .     33:   }
         .          .     34:}
         .          .     35:
         .          .     36:func (s *MemoryStore) Push(md *models.MarketData) {
         .          .     37:   select {
ROUTINE ======================== crypt_proto/internal/queue.(*MemoryStore).apply in /home/gaz358/myprog/crypt_proto/internal/queue/in_memory_queue.go
         0       30ms (flat, cum)  2.73% of Total
         .          .     51:func (s *MemoryStore) apply(md *models.MarketData) {
         .          .     52:   key := md.Exchange + "|" + md.Symbol
         .          .     53:   quote := Quote{
         .          .     54:           Bid: md.Bid, Ask: md.Ask,
         .          .     55:           BidSize: md.BidSize, AskSize: md.AskSize,
         .          .     56:           Timestamp: time.Now().UnixMilli(),
         .          .     57:   }
         .          .     58:
         .          .     59:   oldMap := s.data.Load().(map[string]Quote)
         .       20ms     60:   newMap := make(map[string]Quote, len(oldMap)+1)
         .          .     61:   for k, v := range oldMap {
         .       10ms     62:           newMap[k] = v
         .          .     63:   }
         .          .     64:   newMap[key] = quote
         .          .     65:   s.data.Store(newMap)
         .          .     66:}
(pprof) 
