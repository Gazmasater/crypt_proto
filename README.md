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
	"sync"
	"time"

	"crypt_proto/pkg/models"
)

type Quote struct {
	Bid, Ask         float64
	BidSize, AskSize float64
	Timestamp        int64
}

type RingBuffer struct {
	data []Quote
	size int
	pos  int
	full bool
	mu   sync.RWMutex
}

func NewRingBuffer(size int) *RingBuffer {
	return &RingBuffer{
		data: make([]Quote, size),
		size: size,
	}
}

func (r *RingBuffer) Push(q Quote) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[r.pos] = q
	r.pos++
	if r.pos >= r.size {
		r.pos = 0
		r.full = true
	}
}

func (r *RingBuffer) GetLast() (Quote, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if !r.full && r.pos == 0 {
		return Quote{}, false
	}
	idx := r.pos - 1
	if idx < 0 {
		idx = r.size - 1
	}
	return r.data[idx], true
}

// MemoryStore хранит рингбуфер для каждой пары
type MemoryStore struct {
	buffers map[string]*RingBuffer
	batch   chan *models.MarketData
	mu      sync.RWMutex
	BufSize int
}

func NewMemoryStore(bufSize int) *MemoryStore {
	return &MemoryStore{
		buffers: make(map[string]*RingBuffer),
		batch:   make(chan *models.MarketData, 100_000),
		BufSize: bufSize,
	}
}

func (s *MemoryStore) Run() {
	for md := range s.batch {
		s.apply(md)
	}
}

func (s *MemoryStore) Push(md *models.MarketData) {
	select {
	case s.batch <- md:
	default:
		// drop if full
	}
}

func (s *MemoryStore) Get(exchange, symbol string) (Quote, bool) {
	key := exchange + "|" + symbol
	s.mu.RLock()
	buf, ok := s.buffers[key]
	s.mu.RUnlock()
	if !ok {
		return Quote{}, false
	}
	return buf.GetLast()
}

func (s *MemoryStore) apply(md *models.MarketData) {
	key := md.Exchange + "|" + md.Symbol

	s.mu.RLock()
	buf, ok := s.buffers[key]
	s.mu.RUnlock()

	if !ok {
		buf = NewRingBuffer(s.BufSize)
		s.mu.Lock()
		s.buffers[key] = buf
		s.mu.Unlock()
	}

	buf.Push(Quote{
		Bid:       md.Bid,
		Ask:       md.Ask,
		BidSize:   md.BidSize,
		AskSize:   md.AskSize,
		Timestamp: time.Now().UnixMilli(),
	})
}




package calculator

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"

	"crypt_proto/internal/queue"
	"crypt_proto/pkg/models"
)

const fee = 0.001

type LegIndex struct {
	Key    string
	Symbol string
	IsBuy  bool
}

type Triangle struct {
	A, B, C string
	Legs    [3]LegIndex
}

type Calculator struct {
	mem      *queue.MemoryStore
	bySymbol map[string][]*Triangle
	fileLog  *log.Logger
}

func NewCalculator(mem *queue.MemoryStore, triangles []*Triangle) *Calculator {
	f, err := os.OpenFile("arb_opportunities.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("failed to open log: %v", err)
	}

	bySymbol := make(map[string][]*Triangle, 1024)
	for _, t := range triangles {
		for _, leg := range t.Legs {
			bySymbol[leg.Symbol] = append(bySymbol[leg.Symbol], t)
		}
	}

	log.Printf("[Calculator] indexed %d symbols\n", len(bySymbol))

	return &Calculator{
		mem:      mem,
		bySymbol: bySymbol,
		fileLog:  log.New(f, "", log.LstdFlags),
	}
}

func (c *Calculator) Run(in <-chan *models.MarketData) {
	for md := range in {
		c.mem.Push(md)

		tris := c.bySymbol[md.Symbol]
		if len(tris) == 0 {
			continue
		}

		for _, tri := range tris {
			c.calcTriangle(tri)
		}
	}
}

func (c *Calculator) calcTriangle(tri *Triangle) {
	var q [3]queue.Quote

	for i, leg := range tri.Legs {
		quote, ok := c.mem.Get("KuCoin", leg.Symbol)
		if !ok {
			return
		}
		q[i] = quote
	}

	var usdtLimits [3]float64

	// LEG 1
	if tri.Legs[0].IsBuy {
		if q[0].Ask <= 0 || q[0].AskSize <= 0 {
			return
		}
		usdtLimits[0] = q[0].Ask * q[0].AskSize
	} else {
		if q[0].Bid <= 0 || q[0].BidSize <= 0 {
			return
		}
		usdtLimits[0] = q[0].Bid * q[0].BidSize
	}

	// LEG 2
	if tri.Legs[1].IsBuy {
		if q[1].Ask <= 0 || q[1].AskSize <= 0 || q[2].Bid <= 0 {
			return
		}
		usdtLimits[1] = q[1].Ask * q[1].AskSize * q[2].Bid
	} else {
		if q[1].Bid <= 0 || q[1].BidSize <= 0 || q[2].Bid <= 0 {
			return
		}
		usdtLimits[1] = q[1].BidSize * q[2].Bid
	}

	// LEG 3
	if q[2].Bid <= 0 || q[2].BidSize <= 0 {
		return
	}
	usdtLimits[2] = q[2].Bid * q[2].BidSize

	maxUSDT := usdtLimits[0]
	if usdtLimits[1] < maxUSDT {
		maxUSDT = usdtLimits[1]
	}
	if usdtLimits[2] < maxUSDT {
		maxUSDT = usdtLimits[2]
	}
	if maxUSDT <= 0 {
		return
	}

	amount := maxUSDT

	if tri.Legs[0].IsBuy {
		amount = amount / q[0].Ask * (1 - fee)
	} else {
		amount = amount * q[0].Bid * (1 - fee)
	}

	if tri.Legs[1].IsBuy {
		amount = amount / q[1].Ask * (1 - fee)
	} else {
		amount = amount * q[1].Bid * (1 - fee)
	}

	if tri.Legs[2].IsBuy {
		amount = amount / q[2].Ask * (1 - fee)
	} else {
		amount = amount * q[2].Bid * (1 - fee)
	}

	profitUSDT := amount - maxUSDT
	profitPct := profitUSDT / maxUSDT

	if profitPct > 0.001 && profitUSDT > 0.02 {
		msg := fmt.Sprintf(
			"[ARB] %s→%s→%s | %.4f%% | volume=%.2f USDT | profit=%.4f USDT",
			tri.A, tri.B, tri.C,
			profitPct*100, maxUSDT, profitUSDT,
		)
		log.Println(msg)
		c.fileLog.Println(msg)
	}
}

func ParseTrianglesFromCSV(path string) ([]*Triangle, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rows, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return nil, err
	}

	var res []*Triangle
	for _, row := range rows[1:] {
		if len(row) < 6 {
			continue
		}

		t := &Triangle{
			A: strings.TrimSpace(row[0]),
			B: strings.TrimSpace(row[1]),
			C: strings.TrimSpace(row[2]),
		}

		for i, leg := range []string{row[3], row[4], row[5]} {
			parts := strings.Fields(strings.ToUpper(strings.TrimSpace(leg)))
			if len(parts) != 2 {
				continue
			}
			isBuy := parts[0] == "BUY"
			pair := strings.Split(parts[1], "/")
			if len(pair) != 2 {
				continue
			}
			symbol := pair[0] + "-" + pair[1]
			key := "KuCoin|" + symbol

			t.Legs[i] = LegIndex{
				Key:    key,
				Symbol: symbol,
				IsBuy:  isBuy,
			}
		}

		res = append(res, t)
	}

	return res, nil
}




package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"crypt_proto/internal/calculator"
	"crypt_proto/internal/collector"
	"crypt_proto/internal/queue"
	"crypt_proto/pkg/models"
)

func main() {
	go func() {
		log.Println("pprof on http://localhost:6060/debug/pprof/")
		_ = http.ListenAndServe("localhost:6060", nil)
	}()

	out := make(chan *models.MarketData, 100_000)
	mem := queue.NewMemoryStore(1000) // размер рингбуфера 1000 котировок
	go mem.Run()

	kc, _, err := collector.NewKuCoinCollectorFromCSV("../exchange/data/kucoin_triangles_usdt.csv")
	if err != nil {
		log.Fatal(err)
	}
	if err := kc.Start(out); err != nil {
		log.Fatal(err)
	}
	log.Println("[Main] KuCoinCollector started")

	triangles, _ := calculator.ParseTrianglesFromCSV("../exchange/data/kucoin_triangles_usdt.csv")
	calc := calculator.NewCalculator(mem, triangles)
	go calc.Run(out)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("[Main] shutting down...")
	kc.Stop()
	close(out)
	log.Println("[Main] exited")
}



gaz358@gaz358-BOD-WXX9:~/myprog/crypt_proto$    go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
Fetching profile over HTTP from http://localhost:6060/debug/pprof/profile?seconds=30
Saved profile in /home/gaz358/pprof/pprof.arb.samples.cpu.285.pb.gz
File: arb
Build ID: b6b2e5615b7516cbc5189db116aa14d0d394ed12
Type: cpu
Time: 2026-01-28 02:33:34 MSK
Duration: 30s, Total samples = 790ms ( 2.63%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 590ms, 74.68% of 790ms total
Showing top 10 nodes out of 96
      flat  flat%   sum%        cum   cum%
     380ms 48.10% 48.10%      380ms 48.10%  internal/runtime/syscall.Syscall6
      30ms  3.80% 51.90%      350ms 44.30%  bufio.(*Reader).fill
      30ms  3.80% 55.70%       70ms  8.86%  github.com/tidwall/gjson.getBytes
      30ms  3.80% 59.49%       30ms  3.80%  runtime.futex
      20ms  2.53% 62.03%       40ms  5.06%  github.com/tidwall/gjson.parseObject
      20ms  2.53% 64.56%       20ms  2.53%  github.com/tidwall/gjson.parseSquash
      20ms  2.53% 67.09%       20ms  2.53%  runtime.gopark
      20ms  2.53% 69.62%       20ms  2.53%  runtime.memclrNoHeapPointers
      20ms  2.53% 72.15%       20ms  2.53%  runtime.nextFreeFast
      20ms  2.53% 74.68%      210ms 26.58%  runtime.schedule
(pprof) 
