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



package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"

	"crypt_proto/internal/calculator"
	"crypt_proto/internal/collector"
	"crypt_proto/internal/queue"
)

func main() {

	go func() {
		log.Println("pprof on http://localhost:6060/debug/pprof/")
		if err := http.ListenAndServe("localhost:6060", nil); err != nil {
			log.Printf("pprof server error: %v", err)
		}
	}()

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


package calculator

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"crypt_proto/internal/queue"
)

const feeM = 0.999

/* ===================== MODELS ===================== */

type LegIndex struct {
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
	log      *log.Logger
}

/* ===================== CONSTRUCTOR ===================== */

func NewCalculator(mem *queue.MemoryStore, triangles []*Triangle) *Calculator {
	f, err := os.OpenFile(
		"arb_opportunities.log",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	if err != nil {
		log.Fatalf("failed to open log: %v", err)
	}

	bySymbol := make(map[string][]*Triangle, 512)

	for _, t := range triangles {
		for _, leg := range t.Legs {
			bySymbol[leg.Symbol] = append(bySymbol[leg.Symbol], t)
		}
	}

	log.Printf("[Calculator] indexed %d symbols\n", len(bySymbol))

	return &Calculator{
		mem:      mem,
		bySymbol: bySymbol,
		log:      log.New(f, "", log.LstdFlags),
	}
}

/* ===================== RUN LOOP ===================== */

// Run — периодический пересчёт (pull-модель)
func (c *Calculator) Run() {
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		c.calculate()
	}
}

/* ===================== CORE ===================== */

func (c *Calculator) calculate() {
	seen := make(map[*Triangle]struct{}, 256)

	for _, tris := range c.bySymbol {
		for _, tri := range tris {
			if _, ok := seen[tri]; ok {
				continue
			}
			seen[tri] = struct{}{}
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

	maxUSDT := min(
		legLimitUSDT(tri.Legs[0], q[0], 0),
		legLimitUSDT(tri.Legs[1], q[1], q[2].Bid),
		legLimitUSDT(tri.Legs[2], q[2], 0),
	)

	if maxUSDT <= 0 {
		return
	}

	amount := maxUSDT

	amount = applyLeg(amount, tri.Legs[0], q[0])
	amount = applyLeg(amount, tri.Legs[1], q[1])
	amount = applyLeg(amount, tri.Legs[2], q[2])

	profitUSDT := amount - maxUSDT
	profitPct := profitUSDT / maxUSDT

	if profitPct > 0.001 && profitUSDT > 0.02 {
		msg := fmt.Sprintf(
			"%s [ARB] %s→%s→%s | %.4f%% | volume=%.2f USDT | profit=%.4f USDT",
			time.Now().Format("2006/01/02 15:04:05.000"),
			tri.A, tri.B, tri.C,
			profitPct*100,
			maxUSDT,
			profitUSDT,
		)

		log.Println(msg)
		c.log.Println(msg)
	}
}

/* ===================== HELPERS ===================== */

func legLimitUSDT(leg LegIndex, q queue.Quote, extra float64) float64 {
	if leg.IsBuy {
		if q.Ask <= 0 || q.AskSize <= 0 {
			return 0
		}
		limit := q.Ask * q.AskSize
		if extra > 0 {
			limit *= extra
		}
		return limit
	}

	if q.Bid <= 0 || q.BidSize <= 0 {
		return 0
	}
	limit := q.Bid * q.BidSize
	if extra > 0 {
		limit = q.BidSize * extra
	}
	return limit
}

func applyLeg(amount float64, leg LegIndex, q queue.Quote) float64 {
	if leg.IsBuy {
		if q.Ask <= 0 {
			return 0
		}
		return amount / q.Ask * feeM
	}
	if q.Bid <= 0 {
		return 0
	}
	return amount * q.Bid * feeM
}

func min(a, b, c float64) float64 {
	if a <= 0 || b <= 0 || c <= 0 {
		return 0
	}
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}

/* ===================== CSV PARSER ===================== */

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

		for i, raw := range row[3:6] {
			parts := strings.Fields(strings.ToUpper(strings.TrimSpace(raw)))
			if len(parts) != 2 {
				continue
			}

			isBuy := parts[0] == "BUY"
			pair := strings.Split(parts[1], "/")
			if len(pair) != 2 {
				continue
			}

			t.Legs[i] = LegIndex{
				Symbol: pair[0] + "-" + pair[1],
				IsBuy:  isBuy,
			}
		}

		res = append(res, t)
	}

	return res, nil
}


package collector

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"crypt_proto/internal/queue"
	"crypt_proto/pkg/models"

	"github.com/gorilla/websocket"
	"github.com/tidwall/gjson"
)

const (
	maxSubsPerWS = 126
	subRate      = 120 * time.Millisecond
	pingInterval = 18 * time.Second
)

type KuCoinCollector struct {
	ctx    context.Context
	cancel context.CancelFunc
	wsList []*kucoinWS
	store  *queue.MemoryStore
}

type kucoinWS struct {
	id      int
	conn    *websocket.Conn
	symbols []string
	last    map[string][2]float64
}

func NewKuCoinCollectorFromCSV(path string) (*KuCoinCollector, []string, error) {
	symbols, err := readPairsFromCSV(path)
	if err != nil {
		return nil, nil, err
	}
	if len(symbols) == 0 {
		return nil, nil, fmt.Errorf("no symbols")
	}

	ctx, cancel := context.WithCancel(context.Background())
	var wsList []*kucoinWS
	for i := 0; i < len(symbols); i += maxSubsPerWS {
		end := i + maxSubsPerWS
		if end > len(symbols) {
			end = len(symbols)
		}
		wsList = append(wsList, &kucoinWS{
			id:      len(wsList),
			symbols: symbols[i:end],
			last:    make(map[string][2]float64),
		})
	}

	return &KuCoinCollector{ctx: ctx, cancel: cancel, wsList: wsList}, symbols, nil
}

func (c *KuCoinCollector) Name() string { return "KuCoin" }

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

func (c *KuCoinCollector) Stop() error {
	c.cancel()
	for _, ws := range c.wsList {
		if ws.conn != nil {
			_ = ws.conn.Close()
		}
	}
	return nil
}

func (ws *kucoinWS) connect() error {
	req, _ := http.NewRequest("POST", "https://api.kucoin.com/api/v1/bullet-public", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var r struct {
		Data struct {
			Token           string `json:"token"`
			InstanceServers []struct {
				Endpoint string `json:"endpoint"`
			} `json:"instanceServers"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return err
	}

	url := fmt.Sprintf(
		"%s?token=%s&connectId=%d",
		r.Data.InstanceServers[0].Endpoint,
		r.Data.Token,
		time.Now().UnixNano(),
	)

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return err
	}

	ws.conn = conn
	log.Printf("[KuCoin WS %d] connected\n", ws.id)
	return nil
}

func (ws *kucoinWS) subscribeLoop() {
	t := time.NewTicker(subRate)
	defer t.Stop()
	for _, s := range ws.symbols {
		<-t.C
		_ = ws.conn.WriteJSON(map[string]any{
			"id":       time.Now().UnixNano(),
			"type":     "subscribe",
			"topic":    "/market/ticker:" + s,
			"response": true,
		})
	}
}

func (ws *kucoinWS) pingLoop() {
	t := time.NewTicker(pingInterval)
	defer t.Stop()
	for range t.C {
		_ = ws.conn.WriteJSON(map[string]any{"id": time.Now().UnixNano(), "type": "ping"})
	}
}

func (ws *kucoinWS) readLoop(c *KuCoinCollector) {
	for {
		_, msg, err := ws.conn.ReadMessage()
		if err != nil {
			log.Printf("[KuCoin WS %d] read error: %v\n", ws.id, err)
			return
		}
		ws.handle(c, msg)
	}
}

func (ws *kucoinWS) handle(c *KuCoinCollector, msg []byte) {
	if gjson.GetBytes(msg, "type").String() != "message" {
		return
	}
	topic := gjson.GetBytes(msg, "topic").String()
	if !strings.HasPrefix(topic, "/market/ticker:") {
		return
	}
	symbol := strings.TrimPrefix(topic, "/market/ticker:")
	data := gjson.GetBytes(msg, "data")
	bid, ask := data.Get("bestBid").Float(), data.Get("bestAsk").Float()
	bidSize, askSize := data.Get("bestBidSize").Float(), data.Get("bestAskSize").Float()
	if bid == 0 || ask == 0 {
		return
	}
	last := ws.last[symbol]
	if last[0] == bid && last[1] == ask {
		return
	}
	ws.last[symbol] = [2]float64{bid, ask}
	c.store.Push(&models.MarketData{
		Exchange:  "KuCoin",
		Symbol:    symbol,
		Bid:       bid,
		Ask:       ask,
		BidSize:   bidSize,
		AskSize:   askSize,
		Timestamp: time.Now().UnixMilli(),
	})

}

func readPairsFromCSV(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	rows, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return nil, err
	}
	set := make(map[string]struct{})
	for _, row := range rows[1:] {
		for i := 3; i <= 5 && i < len(row); i++ {
			if p := parseLeg(row[i]); p != "" {
				set[p] = struct{}{}
			}
		}
	}
	res := make([]string, 0, len(set))
	for k := range set {
		res = append(res, k)
	}
	return res, nil
}

func parseLeg(s string) string {
	parts := strings.Fields(strings.ToUpper(strings.TrimSpace(s)))
	if len(parts) < 2 {
		return ""
	}
	p := strings.Split(parts[1], "/")
	if len(p) != 2 {
		return ""
	}
	return p[0] + "-" + p[1]
}



package queue

import (
	"sync/atomic"
	"time"

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


gaz358@gaz358-BOD-WXX9:~/myprog/crypt_proto$    go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
Fetching profile over HTTP from http://localhost:6060/debug/pprof/profile?seconds=30
Saved profile in /home/gaz358/pprof/pprof.arb.samples.cpu.283.pb.gz
File: arb
Build ID: 62de704e11fbabe21b59d9dcb99f95c19fb9dac8
Type: cpu
Time: 2026-01-28 00:24:55 MSK
Duration: 30s, Total samples = 1.74s ( 5.80%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 1250ms, 71.84% of 1740ms total
Showing top 10 nodes out of 158
      flat  flat%   sum%        cum   cum%
     630ms 36.21% 36.21%      630ms 36.21%  internal/runtime/syscall.Syscall6
     380ms 21.84% 58.05%      380ms 21.84%  runtime.futex
      60ms  3.45% 61.49%      930ms 53.45%  runtime.findRunnable
      40ms  2.30% 63.79%       70ms  4.02%  runtime.(*timers).check
      30ms  1.72% 65.52%       90ms  5.17%  runtime.scanobject
      30ms  1.72% 67.24%       80ms  4.60%  runtime.stealWork
      20ms  1.15% 68.39%       20ms  1.15%  crypt_proto/internal/queue.(*ringBuffer).drain
      20ms  1.15% 69.54%      320ms 18.39%  crypto/tls.(*Conn).readRecordOrCCS
      20ms  1.15% 70.69%       40ms  2.30%  github.com/tidwall/gjson.parseObject
      20ms  1.15% 71.84%       20ms  1.15%  internal/runtime/maps.ctrlGroup.matchH2
(pprof) 




gaz358@gaz358-BOD-WXX9:~/myprog/crypt_proto$    go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
Fetching profile over HTTP from http://localhost:6060/debug/pprof/profile?seconds=30
Saved profile in /home/gaz358/pprof/pprof.arb.samples.cpu.284.pb.gz
File: arb
Build ID: e1f97f19b4005e7c00459e4ff590268a861df570
Type: cpu
Time: 2026-01-28 00:33:34 MSK
Duration: 30s, Total samples = 970ms ( 3.23%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top 
Showing nodes accounting for 660ms, 68.04% of 970ms total
Showing top 10 nodes out of 102
      flat  flat%   sum%        cum   cum%
     380ms 39.18% 39.18%      380ms 39.18%  internal/runtime/syscall.Syscall6
      90ms  9.28% 48.45%       90ms  9.28%  runtime.futex
      60ms  6.19% 54.64%       90ms  9.28%  github.com/tidwall/gjson.parseObject
      40ms  4.12% 58.76%       40ms  4.12%  runtime.nextFreeFast
      20ms  2.06% 60.82%      120ms 12.37%  github.com/tidwall/gjson.Get
      20ms  2.06% 62.89%       20ms  2.06%  runtime.nanotime
      20ms  2.06% 64.95%      150ms 15.46%  runtime.netpoll
      10ms  1.03% 65.98%       10ms  1.03%  bytes.(*Buffer).Len
      10ms  1.03% 67.01%      370ms 38.14%  bytes.(*Buffer).ReadFrom
      10ms  1.03% 68.04%       10ms  1.03%  bytes.(*Reader).Read
(pprof) 




