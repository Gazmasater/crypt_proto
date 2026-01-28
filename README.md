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
Saved profile in /home/gaz358/pprof/pprof.arb.samples.cpu.310.pb.gz
File: arb
Build ID: 5f8362b383d14632fc6dc54b7372b7368fcfe544
Type: cpu
Time: 2026-01-28 21:30:50 MSK
Duration: 30s, Total samples = 1.04s ( 3.47%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) list crypt_proto
Total: 1.04s
ROUTINE ======================== crypt_proto/internal/calculator.(*Calculator).Run in /home/gaz358/myprog/crypt_proto/internal/calculator/arb.go
         0       10ms (flat, cum)  0.96% of Total
         .          .     59:func (c *Calculator) Run(in <-chan *models.MarketData) {
         .          .     60:   for md := range in {
         .          .     61:           c.mem.Push(md)
         .          .     62:
         .          .     63:           tris := c.bySymbol[md.Symbol]
         .          .     64:           if len(tris) == 0 {
         .          .     65:                   continue
         .          .     66:           }
         .          .     67:
         .          .     68:           for _, tri := range tris {
         .       10ms     69:                   c.calcTriangle(tri)
         .          .     70:           }
         .          .     71:   }
         .          .     72:}
         .          .     73:
         .          .     74:func (c *Calculator) calcTriangle(tri *Triangle) {
ROUTINE ======================== crypt_proto/internal/calculator.(*Calculator).calcTriangle in /home/gaz358/myprog/crypt_proto/internal/calculator/arb.go
         0       10ms (flat, cum)  0.96% of Total
         .          .     74:func (c *Calculator) calcTriangle(tri *Triangle) {
         .          .     75:   var q [3]queue.Quote
         .          .     76:
         .          .     77:   for i, leg := range tri.Legs {
         .       10ms     78:           quote, ok := c.mem.Get("KuCoin", leg.Symbol)
         .          .     79:           if !ok {
         .          .     80:                   return
         .          .     81:           }
         .          .     82:           q[i] = quote
         .          .     83:   }
ROUTINE ======================== crypt_proto/internal/collector.(*kucoinWS).handle in /home/gaz358/myprog/crypt_proto/internal/collector/kucoin_collector.go
         0      150ms (flat, cum) 14.42% of Total
         .          .    163:func (ws *kucoinWS) handle(c *KuCoinCollector, msg []byte) {
         .       20ms    164:   parsed := gjson.ParseBytes(msg)
         .       30ms    165:   if parsed.Get("type").String() != "message" {
         .          .    166:           return
         .          .    167:   }
         .          .    168:
         .          .    169:   topic := parsed.Get("topic").String()
         .          .    170:   const prefix = "/market/ticker:"
         .          .    171:   if !strings.HasPrefix(topic, prefix) {
         .          .    172:           return
         .          .    173:   }
         .          .    174:   symbol := strings.TrimPrefix(topic, prefix)
         .          .    175:
         .       20ms    176:   bid := parsed.Get("data.bestBid").Float()
         .       10ms    177:   ask := parsed.Get("data.bestAsk").Float()
         .          .    178:   if bid == 0 || ask == 0 {
         .          .    179:           return
         .          .    180:   }
         .          .    181:
         .          .    182:   // проверяем изменения
         .       30ms    183:   if last, ok := ws.last[symbol]; ok && last[0] == bid && last[1] == ask {
         .          .    184:           return
         .          .    185:   }
         .          .    186:
         .          .    187:   // вычисляем размеры только если цены изменились
         .       10ms    188:   bidSize := parsed.Get("data.bestBidSize").Float()
         .       20ms    189:   askSize := parsed.Get("data.bestAskSize").Float()
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
         0      740ms (flat, cum) 71.15% of Total
         .          .    152:func (ws *kucoinWS) readLoop(c *KuCoinCollector) {
         .          .    153:   for {
         .      590ms    154:           _, msg, err := ws.conn.ReadMessage()
         .          .    155:           if err != nil {
         .          .    156:                   log.Printf("[KuCoin WS %d] read error: %v\n", ws.id, err)
         .          .    157:                   return
         .          .    158:           }
         .      150ms    159:           ws.handle(c, msg)
         .          .    160:   }
         .          .    161:}
         .          .    162:
         .          .    163:func (ws *kucoinWS) handle(c *KuCoinCollector, msg []byte) {
         .          .    164:   parsed := gjson.ParseBytes(msg)
ROUTINE ======================== crypt_proto/internal/queue.(*MemoryStore).Get in /home/gaz358/myprog/crypt_proto/internal/queue/in_memory_queue.go
         0       10ms (flat, cum)  0.96% of Total
         .          .     45:func (s *MemoryStore) Get(exchange, symbol string) (Quote, bool) {
         .       10ms     46:   s.mu.RLock()
         .          .     47:   q, ok := s.data[exchange+"|"+symbol]
         .          .     48:   s.mu.RUnlock()
         .          .     49:   return q, ok
         .          .     50:}
         .          .     51:
(pprof) 
