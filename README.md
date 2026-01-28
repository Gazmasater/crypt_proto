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
	// проверка типа
	if !bytes.Contains(msg, []byte(`"type":"message"`)) {
		return
	}

	// проверка топика
	const prefix = `/market/ticker:`
	topicIdx := bytes.Index(msg, []byte(`"topic":"`))
	if topicIdx == -1 {
		return
	}
	topicStart := topicIdx + len(`"topic":"`)
	topicEnd := bytes.IndexByte(msg[topicStart:], '"')
	if topicEnd == -1 {
		return
	}
	topic := string(msg[topicStart : topicStart+topicEnd])
	if !strings.HasPrefix(topic, prefix) {
		return
	}
	symbol := strings.TrimPrefix(topic, prefix)

	// парс чисел (строковые числа KuCoin)
	bid := parseFloat(msg, `"bestBid":`)
	ask := parseFloat(msg, `"bestAsk":`)
	bidSize := parseFloat(msg, `"bestBidSize":`)
	askSize := parseFloat(msg, `"bestAskSize":`)

	if bid == 0 || ask == 0 {
		return
	}

	// Логируем после парсинга
	log.Printf("[KuCoin WS %d] parsed %s bid=%.6f ask=%.6f bidSize=%.6f askSize=%.6f",
		ws.id, symbol, bid, ask, bidSize, askSize)

	last := ws.last[symbol]
	if last[0] == bid && last[1] == ask {
		return
	}
	ws.last[symbol] = [2]float64{bid, ask}

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

// parseFloat ищет число после ключа, учитывая, что KuCoin присылает числа в кавычках
func parseFloat(msg []byte, key string) float64 {
	idx := bytes.Index(msg, []byte(key))
	if idx == -1 {
		return 0
	}
	start := idx + len(key)

	// пропускаем пробелы и кавычки
	for start < len(msg) && (msg[start] == '"' || msg[start] == ' ') {
		start++
	}
	end := start
	for end < len(msg) && ((msg[end] >= '0' && msg[end] <= '9') || msg[end] == '.' || msg[end] == 'e' || msg[end] == 'E' || msg[end] == '-') {
		end++
	}

	if end == start {
		return 0
	}

	val, err := strconv.ParseFloat(string(msg[start:end]), 64)
	if err != nil {
		return 0
	}
	return val
}




Fetching profile over HTTP from http://localhost:6060/debug/pprof/profile?seconds=30
Saved profile in /home/gaz358/pprof/pprof.arb.samples.cpu.301.pb.gz
File: arb
Build ID: 161efe5a1fc7f6dcfc1b03c833c1c1cb43b6632b
Type: cpu
Time: 2026-01-28 19:30:49 MSK
Duration: 30s, Total samples = 1.17s ( 3.90%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) list crypt_proto
Total: 1.17s
ROUTINE ======================== crypt_proto/internal/calculator.(*Calculator).Run in /home/gaz358/myprog/crypt_proto/internal/calculator/arb.go
         0       30ms (flat, cum)  2.56% of Total
         .          .     59:func (c *Calculator) Run(in <-chan *models.MarketData) {
         .       10ms     60:   for md := range in {
         .       10ms     61:           c.mem.Push(md)
         .          .     62:
         .       10ms     63:           tris := c.bySymbol[md.Symbol]
         .          .     64:           if len(tris) == 0 {
         .          .     65:                   continue
         .          .     66:           }
         .          .     67:
         .          .     68:           for _, tri := range tris {
ROUTINE ======================== crypt_proto/internal/collector.(*kucoinWS).handle in /home/gaz358/myprog/crypt_proto/internal/collector/kucoin_collector.go
         0       70ms (flat, cum)  5.98% of Total
         .          .    164:func (ws *kucoinWS) handle(c *KuCoinCollector, msg []byte) {
         .          .    165:   // проверка типа
         .          .    166:   if !bytes.Contains(msg, []byte(`"type":"message"`)) {
         .          .    167:           return
         .          .    168:   }
         .          .    169:
         .          .    170:   // проверка топика
         .          .    171:   const prefix = `/market/ticker:`
         .          .    172:   topicIdx := bytes.Index(msg, []byte(`"topic":"`))
         .          .    173:   if topicIdx == -1 {
         .          .    174:           return
         .          .    175:   }
         .          .    176:   topicStart := topicIdx + len(`"topic":"`)
         .          .    177:   topicEnd := bytes.IndexByte(msg[topicStart:], '"')
         .          .    178:   if topicEnd == -1 {
         .          .    179:           return
         .          .    180:   }
         .          .    181:   topic := string(msg[topicStart : topicStart+topicEnd])
         .          .    182:   if !strings.HasPrefix(topic, prefix) {
         .          .    183:           return
         .          .    184:   }
         .          .    185:   symbol := strings.TrimPrefix(topic, prefix)
         .          .    186:
         .          .    187:   // парс чисел (строковые числа KuCoin)
         .       10ms    188:   bid := parseFloat(msg, `"bestBid":`)
         .          .    189:   ask := parseFloat(msg, `"bestAsk":`)
         .          .    190:   bidSize := parseFloat(msg, `"bestBidSize":`)
         .       20ms    191:   askSize := parseFloat(msg, `"bestAskSize":`)
         .          .    192:
         .          .    193:   if bid == 0 || ask == 0 {
         .          .    194:           return
         .          .    195:   }
         .          .    196:
         .          .    197:   // Логируем после парсинга
         .          .    198:   //log.Printf("[KuCoin WS %d] parsed %s bid=%.6f ask=%.6f bidSize=%.6f askSize=%.6f",
         .          .    199:   //      ws.id, symbol, bid, ask, bidSize, askSize)
         .          .    200:
         .          .    201:   last := ws.last[symbol]
         .          .    202:   if last[0] == bid && last[1] == ask {
         .          .    203:           return
         .          .    204:   }
         .          .    205:   ws.last[symbol] = [2]float64{bid, ask}
         .          .    206:
         .       20ms    207:   c.out <- &models.MarketData{
         .          .    208:           Exchange:  "KuCoin",
         .          .    209:           Symbol:    symbol,
         .          .    210:           Bid:       bid,
         .          .    211:           Ask:       ask,
         .          .    212:           BidSize:   bidSize,
         .          .    213:           AskSize:   askSize,
         .       20ms    214:           Timestamp: time.Now().UnixMilli(),
         .          .    215:   }
         .          .    216:}
         .          .    217:
         .          .    218:// parseFloat ищет число после ключа, учитывая, что KuCoin присылает числа в кавычках
         .          .    219:func parseFloat(msg []byte, key string) float64 {
ROUTINE ======================== crypt_proto/internal/collector.(*kucoinWS).readLoop in /home/gaz358/myprog/crypt_proto/internal/collector/kucoin_collector.go
         0      680ms (flat, cum) 58.12% of Total
         .          .    153:func (ws *kucoinWS) readLoop(c *KuCoinCollector) {
         .          .    154:   for {
         .      610ms    155:           _, msg, err := ws.conn.ReadMessage()
         .          .    156:           if err != nil {
         .          .    157:                   log.Printf("[KuCoin WS %d] read error: %v\n", ws.id, err)
         .          .    158:                   return
         .          .    159:           }
         .       70ms    160:           ws.handle(c, msg)
         .          .    161:   }
         .          .    162:}
         .          .    163:
         .          .    164:func (ws *kucoinWS) handle(c *KuCoinCollector, msg []byte) {
         .          .    165:   // проверка типа
ROUTINE ======================== crypt_proto/internal/collector.parseFloat in /home/gaz358/myprog/crypt_proto/internal/collector/kucoin_collector.go
      10ms       30ms (flat, cum)  2.56% of Total
         .          .    219:func parseFloat(msg []byte, key string) float64 {
         .       10ms    220:   idx := bytes.Index(msg, []byte(key))
         .          .    221:   if idx == -1 {
         .          .    222:           return 0
         .          .    223:   }
         .          .    224:   start := idx + len(key)
         .          .    225:
         .          .    226:   // пропускаем пробелы и кавычки
         .          .    227:   for start < len(msg) && (msg[start] == '"' || msg[start] == ' ') {
         .          .    228:           start++
         .          .    229:   }
         .          .    230:   end := start
      10ms       10ms    231:   for end < len(msg) && ((msg[end] >= '0' && msg[end] <= '9') || msg[end] == '.' || msg[end] == 'e' || msg[end] == 'E' || msg[end] == '-') {
         .          .    232:           end++
         .          .    233:   }
         .          .    234:
         .          .    235:   if end == start {
         .          .    236:           return 0
         .          .    237:   }
         .          .    238:
         .       10ms    239:   val, err := strconv.ParseFloat(string(msg[start:end]), 64)
         .          .    240:   if err != nil {
         .          .    241:           return 0
         .          .    242:   }
         .          .    243:   return val
         .          .    244:}
ROUTINE ======================== crypt_proto/internal/queue.(*MemoryStore).Push in /home/gaz358/myprog/crypt_proto/internal/queue/in_memory_queue.go
      10ms       10ms (flat, cum)  0.85% of Total
         .          .     36:func (s *MemoryStore) Push(md *models.MarketData) {
         .          .     37:   select {
      10ms       10ms     38:   case s.batch <- md:
         .          .     39:   default:
         .          .     40:           // drop if full
         .          .     41:   }
         .          .     42:}
         .          .     43:
ROUTINE ======================== crypt_proto/internal/queue.(*MemoryStore).Run in /home/gaz358/myprog/crypt_proto/internal/queue/in_memory_queue.go
         0      130ms (flat, cum) 11.11% of Total
         .          .     30:func (s *MemoryStore) Run() {
         .          .     31:   for md := range s.batch {
         .      130ms     32:           s.apply(md)
         .          .     33:   }
         .          .     34:}
         .          .     35:
         .          .     36:func (s *MemoryStore) Push(md *models.MarketData) {
         .          .     37:   select {
ROUTINE ======================== crypt_proto/internal/queue.(*MemoryStore).apply in /home/gaz358/myprog/crypt_proto/internal/queue/in_memory_queue.go
         0      130ms (flat, cum) 11.11% of Total
         .          .     51:func (s *MemoryStore) apply(md *models.MarketData) {
         .          .     52:   key := md.Exchange + "|" + md.Symbol
         .          .     53:   quote := Quote{
         .          .     54:           Bid: md.Bid, Ask: md.Ask,
         .          .     55:           BidSize: md.BidSize, AskSize: md.AskSize,
         .          .     56:           Timestamp: time.Now().UnixMilli(),
         .          .     57:   }
         .          .     58:
         .          .     59:   oldMap := s.data.Load().(map[string]Quote)
         .       40ms     60:   newMap := make(map[string]Quote, len(oldMap)+1)
         .       20ms     61:   for k, v := range oldMap {
         .       60ms     62:           newMap[k] = v
         .          .     63:   }
         .          .     64:   newMap[key] = quote
         .       10ms     65:   s.data.Store(newMap)
         .          .     66:}
(pprof) 

