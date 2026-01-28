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



package queue

import (
	"sync/atomic"
	"time"

	"crypt_proto/pkg/models"
)

/* =========================
   Quote
========================= */

type Quote struct {
	Bid, Ask         float64
	BidSize, AskSize float64
	Timestamp        int64
}

/* =========================
   Lock-free RingBuffer
   1 writer / many readers
========================= */

type RingBuffer struct {
	data []Quote
	size uint64
	pos  uint64 // atomic
}

func NewRingBuffer(size int) *RingBuffer {
	return &RingBuffer{
		data: make([]Quote, size),
		size: uint64(size),
	}
}

// Push — lock-free, no heap, single writer
func (r *RingBuffer) Push(q Quote) {
	i := atomic.AddUint64(&r.pos, 1) - 1
	r.data[i%r.size] = q
}

// GetLast — snapshot read
func (r *RingBuffer) GetLast() (Quote, bool) {
	p := atomic.LoadUint64(&r.pos)
	if p == 0 {
		return Quote{}, false
	}
	return r.data[(p-1)%r.size], true
}

/* =========================
   MemoryStore
========================= */

type MemoryStore struct {
	buffers map[string]*RingBuffer
	batch   chan *models.MarketData
	bufSize int
}

/* =========================
   Constructor
========================= */

func NewMemoryStore(bufSize int) *MemoryStore {
	return &MemoryStore{
		buffers: make(map[string]*RingBuffer),
		batch:   make(chan *models.MarketData, 10_000),
		bufSize: bufSize,
	}
}

/* =========================
   Writer loop (single goroutine)
========================= */

func (s *MemoryStore) Run() {
	for md := range s.batch {
		s.apply(md)
	}
}

/* =========================
   Non-blocking push
========================= */

func (s *MemoryStore) Push(md *models.MarketData) {
	select {
	case s.batch <- md:
	default:
		// drop if overloaded
	}
}

/* =========================
   Read API
========================= */

func (s *MemoryStore) Get(exchange, symbol string) (Quote, bool) {
	key := exchange + "|" + symbol
	buf, ok := s.buffers[key]
	if !ok {
		return Quote{}, false
	}
	return buf.GetLast()
}

/* =========================
   Internal apply
========================= */

func (s *MemoryStore) apply(md *models.MarketData) {
	key := md.Exchange + "|" + md.Symbol

	buf, ok := s.buffers[key]
	if !ok {
		buf = NewRingBuffer(s.bufSize)
		s.buffers[key] = buf
	}

	buf.Push(Quote{
		Bid:       md.Bid,
		Ask:       md.Ask,
		BidSize:   md.BidSize,
		AskSize:   md.AskSize,
		Timestamp: time.Now().UnixMilli(),
	})
}



Fetching profile over HTTP from http://localhost:6060/debug/pprof/heap
Saved profile in /home/gaz358/pprof/pprof.arb.alloc_objects.alloc_space.inuse_objects.inuse_space.013.pb.gz
File: arb
Build ID: 2d9348127d2f3102690db75c2474cd12f2fed507
Type: inuse_space
Time: 2026-01-28 11:24:23 MSK
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) list queue
Total: 6.89MB
ROUTINE ======================== crypt_proto/internal/queue.(*MemoryStore).Run in /home/gaz358/myprog/crypt_proto/internal/queue/in_memory_queue.go
  512.05kB     3.41MB (flat, cum) 49.50% of Total
         .          .     79:func (s *MemoryStore) Run() {
  512.05kB   512.05kB     80:   for md := range s.batch {
         .     2.91MB     81:           s.apply(md)
         .          .     82:   }
         .          .     83:}
         .          .     84:
         .          .     85:/* =========================
         .          .     86:   Non-blocking push
ROUTINE ======================== crypt_proto/internal/queue.(*MemoryStore).apply in /home/gaz358/myprog/crypt_proto/internal/queue/in_memory_queue.go
         0     2.91MB (flat, cum) 42.24% of Total
         .          .    114:func (s *MemoryStore) apply(md *models.MarketData) {
         .          .    115:   key := md.Exchange + "|" + md.Symbol
         .          .    116:
         .          .    117:   buf, ok := s.buffers[key]
         .          .    118:   if !ok {
         .     2.91MB    119:           buf = NewRingBuffer(s.bufSize)
         .          .    120:           s.buffers[key] = buf
         .          .    121:   }
         .          .    122:
         .          .    123:   buf.Push(Quote{
         .          .    124:           Bid:       md.Bid,
ROUTINE ======================== crypt_proto/internal/queue.NewRingBuffer in /home/gaz358/myprog/crypt_proto/internal/queue/in_memory_queue.go
    2.91MB     2.91MB (flat, cum) 42.24% of Total
         .          .     31:func NewRingBuffer(size int) *RingBuffer {
         .          .     32:   return &RingBuffer{
    2.91MB     2.91MB     33:           data: make([]Quote, size),
         .          .     34:           size: uint64(size),
         .          .     35:   }
         .          .     36:}
         .          .     37:
         .          .     38:// Push — lock-free, no heap, single writer
(pprof) 




gaz358@gaz358-BOD-WXX9:~/myprog/crypt_proto$    go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
Fetching profile over HTTP from http://localhost:6060/debug/pprof/profile?seconds=30
Saved profile in /home/gaz358/pprof/pprof.arb.samples.cpu.291.pb.gz
File: arb
Build ID: 2d9348127d2f3102690db75c2474cd12f2fed507
Type: cpu
Time: 2026-01-28 11:25:50 MSK
Duration: 30s, Total samples = 770ms ( 2.57%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) list crypt_proto
Total: 770ms
ROUTINE ======================== crypt_proto/internal/calculator.(*Calculator).Run in /home/gaz358/myprog/crypt_proto/internal/calculator/arb.go
         0       10ms (flat, cum)  1.30% of Total
         .          .     56:func (c *Calculator) Run(in <-chan *models.MarketData) {
         .          .     57:   for md := range in {
         .          .     58:           c.mem.Push(md)
         .          .     59:
         .          .     60:           tris := c.bySymbol[md.Symbol]
         .          .     61:           if len(tris) == 0 {
         .          .     62:                   continue
         .          .     63:           }
         .          .     64:
         .          .     65:           for _, tri := range tris {
         .       10ms     66:                   c.calcTriangle(tri)
         .          .     67:           }
         .          .     68:   }
         .          .     69:}
         .          .     70:
         .          .     71:func (c *Calculator) calcTriangle(tri *Triangle) {
ROUTINE ======================== crypt_proto/internal/calculator.(*Calculator).calcTriangle in /home/gaz358/myprog/crypt_proto/internal/calculator/arb.go
         0       10ms (flat, cum)  1.30% of Total
         .          .     71:func (c *Calculator) calcTriangle(tri *Triangle) {
         .          .     72:   var q [3]queue.Quote
         .          .     73:
         .          .     74:   for i, leg := range tri.Legs {
         .       10ms     75:           quote, ok := c.mem.Get("KuCoin", leg.Symbol)
         .          .     76:           if !ok {
         .          .     77:                   return
         .          .     78:           }
         .          .     79:           q[i] = quote
         .          .     80:   }
ROUTINE ======================== crypt_proto/internal/collector.(*kucoinWS).handle in /home/gaz358/myprog/crypt_proto/internal/collector/kucoin_collector.go
         0      130ms (flat, cum) 16.88% of Total
         .          .    163:func (ws *kucoinWS) handle(c *KuCoinCollector, msg []byte) {
         .       10ms    164:   if gjson.GetBytes(msg, "type").String() != "message" {
         .          .    165:           return
         .          .    166:   }
         .       10ms    167:   topic := gjson.GetBytes(msg, "topic").String()
         .          .    168:   if !strings.HasPrefix(topic, "/market/ticker:") {
         .          .    169:           return
         .          .    170:   }
         .          .    171:   symbol := strings.TrimPrefix(topic, "/market/ticker:")
         .       30ms    172:   data := gjson.GetBytes(msg, "data")
         .       10ms    173:   bid, ask := data.Get("bestBid").Float(), data.Get("bestAsk").Float()
         .       30ms    174:   bidSize, askSize := data.Get("bestBidSize").Float(), data.Get("bestAskSize").Float()
         .          .    175:   if bid == 0 || ask == 0 {
         .          .    176:           return
         .          .    177:   }
         .          .    178:   last := ws.last[symbol]
         .          .    179:   if last[0] == bid && last[1] == ask {
         .          .    180:           return
         .          .    181:   }
         .          .    182:   ws.last[symbol] = [2]float64{bid, ask}
         .       40ms    183:   c.out <- &models.MarketData{
         .          .    184:           Exchange:  "KuCoin",
         .          .    185:           Symbol:    symbol,
         .          .    186:           Bid:       bid,
         .          .    187:           Ask:       ask,
         .          .    188:           BidSize:   bidSize,
ROUTINE ======================== crypt_proto/internal/collector.(*kucoinWS).readLoop in /home/gaz358/myprog/crypt_proto/internal/collector/kucoin_collector.go
         0      600ms (flat, cum) 77.92% of Total
         .          .    152:func (ws *kucoinWS) readLoop(c *KuCoinCollector) {
         .          .    153:   for {
         .      470ms    154:           _, msg, err := ws.conn.ReadMessage()
         .          .    155:           if err != nil {
         .          .    156:                   log.Printf("[KuCoin WS %d] read error: %v\n", ws.id, err)
         .          .    157:                   return
         .          .    158:           }
         .      130ms    159:           ws.handle(c, msg)
         .          .    160:   }
         .          .    161:}
         .          .    162:
         .          .    163:func (ws *kucoinWS) handle(c *KuCoinCollector, msg []byte) {
         .          .    164:   if gjson.GetBytes(msg, "type").String() != "message" {
ROUTINE ======================== crypt_proto/internal/queue.(*MemoryStore).Get in /home/gaz358/myprog/crypt_proto/internal/queue/in_memory_queue.go
         0       10ms (flat, cum)  1.30% of Total
         .          .    101:func (s *MemoryStore) Get(exchange, symbol string) (Quote, bool) {
         .          .    102:   key := exchange + "|" + symbol
         .       10ms    103:   buf, ok := s.buffers[key]
         .          .    104:   if !ok {
         .          .    105:           return Quote{}, false
         .          .    106:   }
         .          .    107:   return buf.GetLast()
         .          .    108:}
(pprof) 

