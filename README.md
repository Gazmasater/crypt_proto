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



gaz358@gaz358-BOD-WXX9:~/myprog/crypt_proto$ go tool pprof http://localhost:6060/debug/pprof/heap
Fetching profile over HTTP from http://localhost:6060/debug/pprof/heap
Saved profile in /home/gaz358/pprof/pprof.arb.alloc_objects.alloc_space.inuse_objects.inuse_space.012.pb.gz
File: arb
Build ID: 2d9348127d2f3102690db75c2474cd12f2fed507
Type: inuse_space
Time: 2026-01-28 11:14:37 MSK
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 9692.30kB, 100% of 9692.30kB total
Showing top 10 nodes out of 18
      flat  flat%   sum%        cum   cum%
 7153.92kB 73.81% 73.81%  7153.92kB 73.81%  crypt_proto/internal/queue.NewRingBuffer (inline)
    1026kB 10.59% 84.40%     1026kB 10.59%  runtime.allocm
 1000.34kB 10.32% 94.72%  1000.34kB 10.32%  main.main
  512.05kB  5.28%   100%   512.05kB  5.28%  runtime.acquireSudog
         0     0%   100%  7153.92kB 73.81%  crypt_proto/internal/queue.(*MemoryStore).Run
         0     0%   100%  7153.92kB 73.81%  crypt_proto/internal/queue.(*MemoryStore).apply
         0     0%   100%   512.05kB  5.28%  runtime.chanrecv
         0     0%   100%   512.05kB  5.28%  runtime.chanrecv1
         0     0%   100%  1000.34kB 10.32%  runtime.main
         0     0%   100%     1026kB 10.59%  runtime.mstart
(pprof) 


Fetching profile over HTTP from http://localhost:6060/debug/pprof/profile?seconds=30
Saved profile in /home/gaz358/pprof/pprof.arb.samples.cpu.290.pb.gz
File: arb
Build ID: 2d9348127d2f3102690db75c2474cd12f2fed507
Type: cpu
Time: 2026-01-28 11:14:24 MSK
Duration: 30s, Total samples = 720ms ( 2.40%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 560ms, 77.78% of 720ms total
Showing top 10 nodes out of 110
      flat  flat%   sum%        cum   cum%
     360ms 50.00% 50.00%      360ms 50.00%  internal/runtime/syscall.Syscall6
      50ms  6.94% 56.94%       50ms  6.94%  runtime.futex
      30ms  4.17% 61.11%      280ms 38.89%  runtime.findRunnable
      20ms  2.78% 63.89%       50ms  6.94%  github.com/tidwall/gjson.parseObject
      20ms  2.78% 66.67%       20ms  2.78%  github.com/tidwall/gjson.parseSquash
      20ms  2.78% 69.44%       20ms  2.78%  runtime.nextFreeFast
      20ms  2.78% 72.22%       20ms  2.78%  runtime.pMask.read (inline)
      20ms  2.78% 75.00%       20ms  2.78%  runtime.stealWork
      10ms  1.39% 76.39%       10ms  1.39%  crypt_proto/internal/queue.(*RingBuffer).Push
      10ms  1.39% 77.78%       10ms  1.39%  crypto/tls.(*Conn).handshakeContext
(pprof) 

