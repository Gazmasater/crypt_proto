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

	"crypt_proto/pkg/models"
)

type ringBuffer struct {
	buf   []*models.MarketData
	size  uint64
	write uint64
	read  uint64
}

func newRingBuffer(size int) *ringBuffer {
	return &ringBuffer{
		buf:  make([]*models.MarketData, size),
		size: uint64(size),
	}
}

func (r *ringBuffer) push(md *models.MarketData) {
	w := atomic.AddUint64(&r.write, 1) - 1
	r.buf[w%r.size] = md
}

func (r *ringBuffer) drain(fn func(*models.MarketData)) {
	for {
		rp := atomic.LoadUint64(&r.read)
		wp := atomic.LoadUint64(&r.write)
		if rp >= wp {
			return
		}
		md := r.buf[rp%r.size]
		atomic.AddUint64(&r.read, 1)
		if md != nil {
			fn(md)
		}
	}
}


package queue

import (
	"sync/atomic"
	"time"

	"crypt_proto/pkg/models"
)

type Quote struct {
	Bid, Ask         float64
	BidSize, AskSize float64
	Timestamp        int64
}

type MemoryStore struct {
	data atomic.Value // map[string]Quote
	ring *ringBuffer
}

func NewMemoryStore(bufferSize int) *MemoryStore {
	s := &MemoryStore{
		ring: newRingBuffer(bufferSize),
	}
	s.data.Store(make(map[string]Quote))
	return s
}

func (s *MemoryStore) Push(md *models.MarketData) {
	s.ring.push(md)
}

func (s *MemoryStore) Run() {
	for {
		s.ring.drain(s.apply)
		time.Sleep(1 * time.Millisecond) // минимальный yield
	}
}

func (s *MemoryStore) Get(exchange, symbol string) (Quote, bool) {
	m := s.data.Load().(map[string]Quote)
	q, ok := m[exchange+"|"+symbol]
	return q, ok
}

func (s *MemoryStore) apply(md *models.MarketData) {
	key := md.Exchange + "|" + md.Symbol

	old := s.data.Load().(map[string]Quote)
	q := Quote{
		Bid: md.Bid, Ask: md.Ask,
		BidSize: md.BidSize, AskSize: md.AskSize,
		Timestamp: md.Timestamp,
	}

	// copy-on-write snapshot
	nm := make(map[string]Quote, len(old)+1)
	for k, v := range old {
		nm[k] = v
	}
	nm[key] = q
	s.data.Store(nm)
}



type KuCoinCollector struct {
	ctx    context.Context
	cancel context.CancelFunc
	wsList []*kucoinWS
	store  *queue.MemoryStore
}


func (c *KuCoinCollector) Start(store *queue.MemoryStore) error {
	c.store = store
	for _, ws := range c.wsList {
		if err := ws.connect(); err != nil {
			return err
		}
		go ws.readLoop(c)
		go ws.subscribeLoop()
		go ws.pingLoop()
	}
	return nil
}




c.store.Push(&models.MarketData{
	Exchange:  "KuCoin",
	Symbol:    symbol,
	Bid:       bid,
	Ask:       ask,
	BidSize:   bidSize,
	AskSize:   askSize,
	Timestamp: time.Now().UnixMilli(),
})




func main() {
	mem := queue.NewMemoryStore(50_000)
	go mem.Run()

	kc, _, err := collector.NewKuCoinCollectorFromCSV(
		"../exchange/data/kucoin_triangles_usdt.csv",
	)
	if err != nil {
		log.Fatal(err)
	}
	if err := kc.Start(mem); err != nil {
		log.Fatal(err)
	}

	triangles, _ := calculator.ParseTrianglesFromCSV(
		"../exchange/data/kucoin_triangles_usdt.csv",
	)
	calc := calculator.NewCalculator(mem, triangles)
	go calc.Run()

	select {}
}



func (c *Calculator) Run() {
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		c.calculate()
	}
}



