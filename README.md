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



gaz358@gaz358-BOD-WXX9:~/myprog/crypt_proto$    go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
Fetching profile over HTTP from http://localhost:6060/debug/pprof/profile?seconds=30
Saved profile in /home/gaz358/pprof/pprof.arb.samples.cpu.299.pb.gz
File: arb
Build ID: 6ecf4917e7db0b58eb9a9f93fd4c58e027784165
Type: cpu
Time: 2026-01-28 15:37:07 MSK
Duration: 30s, Total samples = 950ms ( 3.17%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) list crypt_proto
Total: 950ms
ROUTINE ======================== crypt_proto/internal/collector.(*kucoinWS).handle in /home/gaz358/myprog/crypt_proto/internal/collector/kucoin_collector.go
         0       90ms (flat, cum)  9.47% of Total
         .          .    163:func (ws *kucoinWS) handle(c *KuCoinCollector, msg []byte) {
         .       30ms    164:   if gjson.GetBytes(msg, "type").String() != "message" {
         .          .    165:           return
         .          .    166:   }
         .          .    167:   topic := gjson.GetBytes(msg, "topic").String()
         .          .    168:   if !strings.HasPrefix(topic, "/market/ticker:") {
         .          .    169:           return
         .          .    170:   }
         .          .    171:   symbol := strings.TrimPrefix(topic, "/market/ticker:")
         .       10ms    172:   data := gjson.GetBytes(msg, "data")
         .       30ms    173:   bid, ask := data.Get("bestBid").Float(), data.Get("bestAsk").Float()
         .       10ms    174:   bidSize, askSize := data.Get("bestBidSize").Float(), data.Get("bestAskSize").Float()
         .          .    175:   if bid == 0 || ask == 0 {
         .          .    176:           return
         .          .    177:   }
         .          .    178:   last := ws.last[symbol]
         .          .    179:   if last[0] == bid && last[1] == ask {
         .          .    180:           return
         .          .    181:   }
         .          .    182:   ws.last[symbol] = [2]float64{bid, ask}
         .       10ms    183:   c.out <- &models.MarketData{
         .          .    184:           Exchange:  "KuCoin",
         .          .    185:           Symbol:    symbol,
         .          .    186:           Bid:       bid,
         .          .    187:           Ask:       ask,
         .          .    188:           BidSize:   bidSize,
ROUTINE ======================== crypt_proto/internal/collector.(*kucoinWS).readLoop in /home/gaz358/myprog/crypt_proto/internal/collector/kucoin_collector.go
         0      640ms (flat, cum) 67.37% of Total
         .          .    152:func (ws *kucoinWS) readLoop(c *KuCoinCollector) {
         .          .    153:   for {
         .      550ms    154:           _, msg, err := ws.conn.ReadMessage()
         .          .    155:           if err != nil {
         .          .    156:                   log.Printf("[KuCoin WS %d] read error: %v\n", ws.id, err)
         .          .    157:                   return
         .          .    158:           }
         .       90ms    159:           ws.handle(c, msg)
         .          .    160:   }
         .          .    161:}
         .          .    162:
         .          .    163:func (ws *kucoinWS) handle(c *KuCoinCollector, msg []byte) {
         .          .    164:   if gjson.GetBytes(msg, "type").String() != "message" {
ROUTINE ======================== crypt_proto/internal/queue.(*MemoryStore).Run in /home/gaz358/myprog/crypt_proto/internal/queue/in_memory_queue.go
         0       40ms (flat, cum)  4.21% of Total
         .          .     30:func (s *MemoryStore) Run() {
         .          .     31:   for md := range s.batch {
         .       40ms     32:           s.apply(md)
         .          .     33:   }
         .          .     34:}
         .          .     35:
         .          .     36:func (s *MemoryStore) Push(md *models.MarketData) {
         .          .     37:   select {
ROUTINE ======================== crypt_proto/internal/queue.(*MemoryStore).apply in /home/gaz358/myprog/crypt_proto/internal/queue/in_memory_queue.go
      10ms       40ms (flat, cum)  4.21% of Total
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
      10ms       10ms     61:   for k, v := range oldMap {
         .       10ms     62:           newMap[k] = v
         .          .     63:   }
         .          .     64:   newMap[key] = quote
         .          .     65:   s.data.Store(newMap)
         .          .     66:}
(pprof) 
