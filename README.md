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
	parsed := gjson.ParseBytes(msg)
	if parsed.Get("type").String() != "message" {
		return
	}

	topic := parsed.Get("topic").String()
	const prefix = "/market/ticker:"
	if !strings.HasPrefix(topic, prefix) {
		return
	}
	symbol := strings.TrimPrefix(topic, prefix)

	bid := parsed.Get("data.bestBid").Float()
	ask := parsed.Get("data.bestAsk").Float()
	if bid == 0 || ask == 0 {
		return
	}

	// проверяем изменения
	if last, ok := ws.last[symbol]; ok && last[0] == bid && last[1] == ask {
		return
	}

	// вычисляем размеры только если цены изменились
	bidSize := parsed.Get("data.bestBidSize").Float()
	askSize := parsed.Get("data.bestAskSize").Float()

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
}



gaz358@gaz358-BOD-WXX9:~/myprog/crypt_proto$    go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
Fetching profile over HTTP from http://localhost:6060/debug/pprof/profile?seconds=30
Saved profile in /home/gaz358/pprof/pprof.arb.samples.cpu.304.pb.gz
File: arb
Build ID: 33c38191a9482af8df18ca4c46eaf82c98d8b6e0
Type: cpu
Time: 2026-01-28 20:04:58 MSK
Duration: 30s, Total samples = 930ms ( 3.10%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) list crypt_proto
Total: 930ms
ROUTINE ======================== crypt_proto/internal/calculator.(*Calculator).Run in /home/gaz358/myprog/crypt_proto/internal/calculator/arb.go
         0       10ms (flat, cum)  1.08% of Total
         .          .     59:func (c *Calculator) Run(in <-chan *models.MarketData) {
         .          .     60:   for md := range in {
         .          .     61:           c.mem.Push(md)
         .          .     62:
         .       10ms     63:           tris := c.bySymbol[md.Symbol]
         .          .     64:           if len(tris) == 0 {
         .          .     65:                   continue
         .          .     66:           }
         .          .     67:
         .          .     68:           for _, tri := range tris {
ROUTINE ======================== crypt_proto/internal/collector.(*kucoinWS).handle in /home/gaz358/myprog/crypt_proto/internal/collector/kucoin_collector.go
         0       70ms (flat, cum)  7.53% of Total
         .          .    163:func (ws *kucoinWS) handle(c *KuCoinCollector, msg []byte) {
         .          .    164:   parsed := gjson.ParseBytes(msg)
         .          .    165:   if parsed.Get("type").String() != "message" {
         .          .    166:           return
         .          .    167:   }
         .          .    168:
         .       10ms    169:   topic := parsed.Get("topic").String()
         .          .    170:   const prefix = "/market/ticker:"
         .          .    171:   if !strings.HasPrefix(topic, prefix) {
         .          .    172:           return
         .          .    173:   }
         .          .    174:   symbol := strings.TrimPrefix(topic, prefix)
         .          .    175:
         .       40ms    176:   bid := parsed.Get("data.bestBid").Float()
         .          .    177:   ask := parsed.Get("data.bestAsk").Float()
         .          .    178:   if bid == 0 || ask == 0 {
         .          .    179:           return
         .          .    180:   }
         .          .    181:
         .          .    182:   // проверяем изменения
         .       10ms    183:   if last, ok := ws.last[symbol]; ok && last[0] == bid && last[1] == ask {
         .          .    184:           return
         .          .    185:   }
         .          .    186:
         .          .    187:   // вычисляем размеры только если цены изменились
         .          .    188:   bidSize := parsed.Get("data.bestBidSize").Float()
         .          .    189:   askSize := parsed.Get("data.bestAskSize").Float()
         .          .    190:
         .          .    191:   // обновляем last
         .          .    192:   ws.last[symbol] = [2]float64{bid, ask}
         .          .    193:
         .          .    194:   // отправляем в канал
         .       10ms    195:   c.out <- &models.MarketData{
         .          .    196:           Exchange:  "KuCoin",
         .          .    197:           Symbol:    symbol,
         .          .    198:           Bid:       bid,
         .          .    199:           Ask:       ask,
         .          .    200:           BidSize:   bidSize,
ROUTINE ======================== crypt_proto/internal/collector.(*kucoinWS).readLoop in /home/gaz358/myprog/crypt_proto/internal/collector/kucoin_collector.go
         0      520ms (flat, cum) 55.91% of Total
         .          .    152:func (ws *kucoinWS) readLoop(c *KuCoinCollector) {
         .          .    153:   for {
         .      450ms    154:           _, msg, err := ws.conn.ReadMessage()
         .          .    155:           if err != nil {
         .          .    156:                   log.Printf("[KuCoin WS %d] read error: %v\n", ws.id, err)
         .          .    157:                   return
         .          .    158:           }
         .       70ms    159:           ws.handle(c, msg)
         .          .    160:   }
         .          .    161:}
         .          .    162:
         .          .    163:func (ws *kucoinWS) handle(c *KuCoinCollector, msg []byte) {
         .          .    164:   parsed := gjson.ParseBytes(msg)
ROUTINE ======================== crypt_proto/internal/queue.(*MemoryStore).Run in /home/gaz358/myprog/crypt_proto/internal/queue/in_memory_queue.go
         0       60ms (flat, cum)  6.45% of Total
         .          .     30:func (s *MemoryStore) Run() {
         .          .     31:   for md := range s.batch {
         .       60ms     32:           s.apply(md)
         .          .     33:   }
         .          .     34:}
         .          .     35:
         .          .     36:func (s *MemoryStore) Push(md *models.MarketData) {
         .          .     37:   select {
ROUTINE ======================== crypt_proto/internal/queue.(*MemoryStore).apply in /home/gaz358/myprog/crypt_proto/internal/queue/in_memory_queue.go
         0       60ms (flat, cum)  6.45% of Total
         .          .     51:func (s *MemoryStore) apply(md *models.MarketData) {
         .          .     52:   key := md.Exchange + "|" + md.Symbol
         .          .     53:   quote := Quote{
         .          .     54:           Bid: md.Bid, Ask: md.Ask,
         .          .     55:           BidSize: md.BidSize, AskSize: md.AskSize,
         .          .     56:           Timestamp: time.Now().UnixMilli(),
         .          .     57:   }
         .          .     58:
         .          .     59:   oldMap := s.data.Load().(map[string]Quote)
         .       10ms     60:   newMap := make(map[string]Quote, len(oldMap)+1)
         .       30ms     61:   for k, v := range oldMap {
         .       10ms     62:           newMap[k] = v
         .          .     63:   }
         .          .     64:   newMap[key] = quote
         .       10ms     65:   s.data.Store(newMap)
         .          .     66:}
(pprof) 

