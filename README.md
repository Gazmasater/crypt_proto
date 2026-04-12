git rm --cached cmd/arb/metrics/arb_metrics.csv

echo "cmd/arb/metrics/*.csv" >> .gitignore

git add .gitignore
git commit --amend --no-edit


git push origin new_arh --force



git filter-branch --force --index-filter \
'git rm --cached --ignore-unmatch cmd/arb/metrics/arb_metrics.csv' \
--prune-empty --tag-name-filter cat -- new_arh


rm -rf .git/refs/original/
git reflog expire --expire=now --all
git gc --prune=now --aggressive


git push origin new_arh --force


Fetching profile over HTTP from http://localhost:6060/debug/pprof/profile?seconds=30
Saved profile in /home/gaz358/pprof/pprof.arb.samples.cpu.358.pb.gz
File: arb
Build ID: 0ba99bd9ae2047c4bad9c3d8309b6d2d1b541df0
Type: cpu
Time: 2026-04-12 16:56:50 MSK
Duration: 30s, Total samples = 3.32s (11.07%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 1900ms, 57.23% of 3320ms total
Dropped 76 nodes (cum <= 16.60ms)
Showing top 10 nodes out of 149
      flat  flat%   sum%        cum   cum%
     820ms 24.70% 24.70%      820ms 24.70%  internal/runtime/syscall.Syscall6
     310ms  9.34% 34.04%      430ms 12.95%  internal/runtime/maps.(*Iter).Next
     170ms  5.12% 39.16%      450ms 13.55%  github.com/tidwall/gjson.parseObject
     170ms  5.12% 44.28%      170ms  5.12%  runtime.futex
     110ms  3.31% 47.59%      110ms  3.31%  github.com/tidwall/gjson.parseSquash
      90ms  2.71% 50.30%      570ms 17.17%  crypt_proto/internal/collector.computeTop
      80ms  2.41% 52.71%       80ms  2.41%  runtime.duffcopy
      50ms  1.51% 54.22%       50ms  1.51%  github.com/tidwall/gjson.parseObjectPath
      50ms  1.51% 55.72%       50ms  1.51%  github.com/tidwall/gjson.parseString
      50ms  1.51% 57.23%       50ms  1.51%  internal/runtime/maps.ctrlGroup.matchFull (inline)
(pprof) 





Saved profile in /home/gaz358/pprof/pprof.arb.samples.cpu.359.pb.gz
File: arb
Build ID: ec679a626462caba32cf0b7059a06bc6a77ccf33
Type: cpu
Time: 2026-04-12 17:15:14 MSK
Duration: 30s, Total samples = 3.13s (10.43%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 1.76s, 56.23% of 3.13s total
Dropped 93 nodes (cum <= 0.02s)
Showing top 10 nodes out of 145
      flat  flat%   sum%        cum   cum%
     1.04s 33.23% 33.23%      1.04s 33.23%  internal/runtime/syscall.Syscall6
     0.22s  7.03% 40.26%      0.22s  7.03%  runtime.futex
     0.12s  3.83% 44.09%      0.31s  9.90%  github.com/tidwall/gjson.parseObject
     0.07s  2.24% 46.33%      0.07s  2.24%  github.com/tidwall/gjson.parseSquash
     0.07s  2.24% 48.56%      0.07s  2.24%  runtime.duffcopy
     0.07s  2.24% 50.80%      0.07s  2.24%  runtime.nextFreeFast
     0.06s  1.92% 52.72%      0.06s  1.92%  strconv.readFloat
     0.05s  1.60% 54.31%      0.05s  1.60%  runtime.casgstatus
     0.03s  0.96% 55.27%      0.03s  0.96%  crypto/tls.(*halfConn).explicitNonceLen
     0.03s  0.96% 56.23%      0.03s  0.96%  github.com/tidwall/gjson.Result.String
(pprof) 






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
	"crypt_proto/internal/executor"
	"crypt_proto/internal/queue"
	"crypt_proto/pkg/models"
)

func main() {
	go func() {
		log.Println("pprof on http://localhost:6060/debug/pprof/")
		if err := http.ListenAndServe("localhost:6060", nil); err != nil {
			log.Printf("pprof server stopped: %v", err)
		}
	}()

	out := make(chan *models.MarketData, 100_000)
	oppCh := make(chan *executor.Opportunity, 4096)

	mem := queue.NewMemoryStore()
	go mem.Run()

	kc, _, err := collector.NewKuCoinCollectorFromCSV("../exchange/data/kucoin_triangles_usdt.csv")
	if err != nil {
		log.Fatal(err)
	}
	if err := kc.Start(out); err != nil {
		log.Fatal(err)
	}
	log.Println("[Main] KuCoinCollector started")

	triangles, err := calculator.ParseTrianglesFromCSV("../exchange/data/kucoin_triangles_usdt.csv")
	if err != nil {
		log.Fatal(err)
	}

	calc := calculator.NewCalculator(mem, triangles, oppCh)
	go calc.Run(out)

	exec := executor.NewExecutor(executor.DefaultConfig())
	go func() {
		for op := range oppCh {
			exec.Handle(op)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("[Main] shutting down...")

	if err := kc.Stop(); err != nil {
		log.Printf("[Main] collector stop error: %v", err)
	}

	close(out)
	close(oppCh)

	log.Println("[Main] exited")
}





package executor

import (
	"fmt"
	"log"
	"math"
	"os"
	"strings"
	"sync"
	"time"

	"crypt_proto/internal/queue"
)

const feeM = 0.9992

type LegRule struct {
	Index int

	RawLeg      string
	Step        float64
	MinQty      float64
	MinNotional float64

	Symbol      string
	Side        string
	Base        string
	Quote       string
	QtyStep     float64
	QuoteStep   float64
	PriceStep   float64
	LegMinQty   float64
	LegMinQuote float64
	LegMinNotnl float64

	Key string
}

type Triangle struct {
	A, B, C string
	Legs    [3]LegRule
}

type Opportunity struct {
	Exchange string
	Triangle *Triangle
	AnchorTS int64

	Quotes   [3]queue.Quote
	AgesMS   [3]int64
	MaxStart float64

	ProfitPct  float64
	ProfitUSDT float64
	FinalUSDT  float64
}

type RejectReason string

const (
	RejectDuplicate          RejectReason = "duplicate"
	RejectNilTriangle        RejectReason = "nil_triangle"
	RejectBadAnchor          RejectReason = "bad_anchor"
	RejectNoQuotes           RejectReason = "no_quotes"
	RejectStale              RejectReason = "stale"
	RejectDesync             RejectReason = "desync"
	RejectSmallMaxStart      RejectReason = "small_max_start"
	RejectSmallSafeStart     RejectReason = "small_safe_start"
	RejectSimulationFailed   RejectReason = "simulation_failed"
	RejectNonPositiveProfit  RejectReason = "non_positive_profit"
	RejectTooSmallProfitUSDT RejectReason = "too_small_profit_usdt"
)

type Config struct {
	SafeStartFactor float64
	MinSafeStart    float64
	MaxQuoteAgeMS   int64
	MaxAgeSpreadMS  int64
	MinProfitUSDT   float64
	DedupWindow     time.Duration
	SummaryEvery    time.Duration
}

func DefaultConfig() Config {
	return Config{
		SafeStartFactor: 0.70,
		MinSafeStart:    20.0,
		MaxQuoteAgeMS:   300,
		MaxAgeSpreadMS:  120,
		MinProfitUSDT:   0.001,
		DedupWindow:     250 * time.Millisecond,
		SummaryEvery:    10 * time.Second,
	}
}

type Executor struct {
	cfg Config

	mu          sync.Mutex
	lastKey     string
	lastSeenAt  time.Time
	fileLog     *log.Logger
	rejects     map[RejectReason]uint64
	accepted    uint64
	lastSummary time.Time
}

func NewExecutor(cfg Config) *Executor {
	f, err := os.OpenFile("executor_paper.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		log.Fatalf("failed to open executor log: %v", err)
	}

	now := time.Now()

	return &Executor{
		cfg:         cfg,
		fileLog:     log.New(f, "", log.LstdFlags),
		rejects:     make(map[RejectReason]uint64),
		lastSummary: now,
	}
}

func (e *Executor) Handle(op *Opportunity) {
	if op == nil || op.Triangle == nil {
		e.reject(RejectNilTriangle, nil, "")
		return
	}
	if op.AnchorTS <= 0 {
		e.reject(RejectBadAnchor, op, "")
		return
	}
	if !e.hasValidQuotes(op) {
		e.reject(RejectNoQuotes, op, "")
		return
	}

	minAge, maxAge, spreadAge := minMaxSpread(op.AgesMS)
	if maxAge < 0 || maxAge > e.cfg.MaxQuoteAgeMS {
		e.reject(RejectStale, op, fmt.Sprintf("max_age=%d", maxAge))
		return
	}
	if spreadAge > e.cfg.MaxAgeSpreadMS {
		e.reject(RejectDesync, op, fmt.Sprintf("spread_age=%d", spreadAge))
		return
	}

	if op.MaxStart <= 0 {
		e.reject(RejectSmallMaxStart, op, "")
		return
	}

	safeStart := floorToReasonable(op.MaxStart*e.cfg.SafeStartFactor, 1e-8)
	if safeStart < e.cfg.MinSafeStart {
		e.reject(RejectSmallSafeStart, op, fmt.Sprintf("safe_start=%.8f", safeStart))
		return
	}

	if e.isDuplicate(op, safeStart) {
		e.reject(RejectDuplicate, op, "")
		return
	}

	finalAmount, diag, ok := simulateTriangle(safeStart, op.Triangle, op.Quotes)
	if !ok || finalAmount <= 0 {
		e.reject(RejectSimulationFailed, op, fmt.Sprintf("safe_start=%.8f", safeStart))
		return
	}

	profitUSDT := finalAmount - safeStart
	profitPct := profitUSDT / safeStart

	if profitPct <= 0 {
		e.reject(RejectNonPositiveProfit, op, fmt.Sprintf("safe_start=%.8f profit_pct=%.8f", safeStart, profitPct))
		return
	}
	if profitUSDT < e.cfg.MinProfitUSDT {
		e.reject(RejectTooSmallProfitUSDT, op, fmt.Sprintf("profit_usdt=%.8f", profitUSDT))
		return
	}

	e.accept(op, safeStart, finalAmount, profitPct, profitUSDT, diag, minAge, maxAge, spreadAge)
}

func (e *Executor) hasValidQuotes(op *Opportunity) bool {
	for _, q := range op.Quotes {
		if q.Bid <= 0 || q.Ask <= 0 || q.BidSize <= 0 || q.AskSize <= 0 {
			return false
		}
	}
	return true
}

func (e *Executor) isDuplicate(op *Opportunity, safeStart float64) bool {
	key := fmt.Sprintf("%s->%s->%s|%d|%.4f",
		op.Triangle.A, op.Triangle.B, op.Triangle.C,
		op.AnchorTS/100,
		round4(safeStart),
	)

	now := time.UnixMilli(op.AnchorTS)

	e.mu.Lock()
	defer e.mu.Unlock()

	if e.lastKey == key && now.Sub(e.lastSeenAt) < e.cfg.DedupWindow {
		return true
	}
	e.lastKey = key
	e.lastSeenAt = now
	return false
}

func (e *Executor) accept(
	op *Opportunity,
	safeStart float64,
	finalAmount float64,
	profitPct float64,
	profitUSDT float64,
	diag [3]legExecution,
	minAge, maxAge, spreadAge int64,
) {
	msg := fmt.Sprintf(
		"[EXEC PAPER OK] %s→%s→%s | safe=%.2f USDT | final=%.6f | profit=%.6f USDT | %.4f%% | age[min=%d max=%d spread=%d] | l1=%s %s out=%.8f | l2=%s %s out=%.8f | l3=%s %s out=%.8f",
		op.Triangle.A, op.Triangle.B, op.Triangle.C,
		safeStart, finalAmount, profitUSDT, profitPct*100,
		minAge, maxAge, spreadAge,
		op.Triangle.Legs[0].Symbol, op.Triangle.Legs[0].Side, diag[0].Out,
		op.Triangle.Legs[1].Symbol, op.Triangle.Legs[1].Side, diag[1].Out,
		op.Triangle.Legs[2].Symbol, op.Triangle.Legs[2].Side, diag[2].Out,
	)

	log.Println(msg)
	e.fileLog.Println(msg)

	e.mu.Lock()
	e.accepted++
	e.maybeSummaryLocked()
	e.mu.Unlock()
}

func (e *Executor) reject(reason RejectReason, op *Opportunity, extra string) {
	e.mu.Lock()
	e.rejects[reason]++
	e.maybeSummaryLocked()
	e.mu.Unlock()

	if op == nil || op.Triangle == nil {
		return
	}

	msg := fmt.Sprintf("[EXEC PAPER REJECT] reason=%s tri=%s->%s->%s %s",
		reason, op.Triangle.A, op.Triangle.B, op.Triangle.C, extra)
	e.fileLog.Println(msg)
}

func (e *Executor) maybeSummaryLocked() {
	now := time.Now()
	if now.Sub(e.lastSummary) < e.cfg.SummaryEvery {
		return
	}
	e.lastSummary = now

	msg := fmt.Sprintf(
		"[Executor] summary accepted=%d reject_duplicate=%d reject_stale=%d reject_desync=%d reject_small_safe=%d reject_non_positive=%d reject_too_small_profit=%d reject_sim_failed=%d",
		e.accepted,
		e.rejects[RejectDuplicate],
		e.rejects[RejectStale],
		e.rejects[RejectDesync],
		e.rejects[RejectSmallSafeStart],
		e.rejects[RejectNonPositiveProfit],
		e.rejects[RejectTooSmallProfitUSDT],
		e.rejects[RejectSimulationFailed],
	)

	log.Println(msg)
	e.fileLog.Println(msg)
}

type legExecution struct {
	In            float64
	Out           float64
	Price         float64
	BookLimitIn   float64
	TradeQty      float64
	TradeNotional float64
}

func simulateTriangle(startUSDT float64, tri *Triangle, q [3]queue.Quote) (float64, [3]legExecution, bool) {
	var diag [3]legExecution
	amount := startUSDT

	for i := 0; i < 3; i++ {
		out, d, ok := executeLeg(amount, tri.Legs[i], q[i])
		if !ok {
			return 0, diag, false
		}
		diag[i] = d
		amount = out
	}

	return amount, diag, true
}

func executeLeg(in float64, leg LegRule, q queue.Quote) (float64, legExecution, bool) {
	side := strings.ToUpper(strings.TrimSpace(leg.Side))
	if side != "BUY" && side != "SELL" {
		return 0, legExecution{}, false
	}
	if in <= 0 {
		return 0, legExecution{}, false
	}

	qtyStep := firstPositive(leg.QtyStep, leg.Step)
	minQty := firstPositive(leg.LegMinQty, leg.MinQty)
	minQuote := leg.LegMinQuote
	minNotional := firstPositive(leg.LegMinNotnl, leg.MinNotional)

	switch side {
	case "BUY":
		if q.Ask <= 0 || q.AskSize <= 0 {
			return 0, legExecution{}, false
		}

		bookLimitIn := q.Ask * q.AskSize
		spendQuote := math.Min(in, bookLimitIn)
		if spendQuote <= 0 {
			return 0, legExecution{}, false
		}

		rawQty := spendQuote / q.Ask
		tradeQty := floorToStep(rawQty, qtyStep)
		if tradeQty <= 0 {
			return 0, legExecution{}, false
		}

		tradeNotional := tradeQty * q.Ask
		if tradeNotional <= 0 || tradeNotional > spendQuote+eps() {
			return 0, legExecution{}, false
		}
		if minQty > 0 && tradeQty+eps() < minQty {
			return 0, legExecution{}, false
		}
		if minQuote > 0 && tradeNotional+eps() < minQuote {
			return 0, legExecution{}, false
		}
		if minNotional > 0 && tradeNotional+eps() < minNotional {
			return 0, legExecution{}, false
		}

		outBase := tradeQty * feeM
		if outBase <= 0 {
			return 0, legExecution{}, false
		}

		return outBase, legExecution{
			In:            in,
			Out:           outBase,
			Price:         q.Ask,
			BookLimitIn:   bookLimitIn,
			TradeQty:      tradeQty,
			TradeNotional: tradeNotional,
		}, true

	case "SELL":
		if q.Bid <= 0 || q.BidSize <= 0 {
			return 0, legExecution{}, false
		}

		bookLimitIn := q.BidSize
		sellBase := math.Min(in, bookLimitIn)
		tradeQty := floorToStep(sellBase, qtyStep)
		if tradeQty <= 0 {
			return 0, legExecution{}, false
		}

		tradeNotional := tradeQty * q.Bid
		if tradeNotional <= 0 {
			return 0, legExecution{}, false
		}
		if minQty > 0 && tradeQty+eps() < minQty {
			return 0, legExecution{}, false
		}
		if minQuote > 0 && tradeNotional+eps() < minQuote {
			return 0, legExecution{}, false
		}
		if minNotional > 0 && tradeNotional+eps() < minNotional {
			return 0, legExecution{}, false
		}

		outQuote := tradeNotional * feeM
		if outQuote <= 0 {
			return 0, legExecution{}, false
		}

		return outQuote, legExecution{
			In:            in,
			Out:           outQuote,
			Price:         q.Bid,
			BookLimitIn:   bookLimitIn,
			TradeQty:      tradeQty,
			TradeNotional: tradeNotional,
		}, true
	}

	return 0, legExecution{}, false
}

func firstPositive(vals ...float64) float64 {
	for _, v := range vals {
		if v > 0 {
			return v
		}
	}
	return 0
}

func floorToStep(v, step float64) float64 {
	if v <= 0 {
		return 0
	}
	if step <= 0 {
		return v
	}
	n := math.Floor((v + eps()) / step)
	if n <= 0 {
		return 0
	}
	return n * step
}

func floorToReasonable(v, step float64) float64 {
	if v <= 0 {
		return 0
	}
	if step <= 0 {
		return v
	}
	n := math.Floor(v / step)
	if n <= 0 {
		return 0
	}
	return n * step
}

func minMaxSpread(ages [3]int64) (int64, int64, int64) {
	minAge := ages[0]
	maxAge := ages[0]
	for _, age := range ages[1:] {
		if age < minAge {
			minAge = age
		}
		if age > maxAge {
			maxAge = age
		}
	}
	return minAge, maxAge, maxAge - minAge
}

func round4(v float64) float64 {
	return math.Round(v*1e4) / 1e4
}

func eps() float64 { return 1e-12 }






Что нужно переписать в internal/calculator/arb.go

Ниже только те места, которые обязательно заменить.

Импорт

Добавь:

"crypt_proto/internal/executor"
Поле в Calculator
oppOut chan<- *executor.Opportunity
Сигнатура конструктора

Замени:

func NewCalculator(mem *queue.MemoryStore, triangles []*Triangle) *Calculator {

на:

func NewCalculator(mem *queue.MemoryStore, triangles []*Triangle, oppOut chan<- *executor.Opportunity) *Calculator {
Возврат из NewCalculator

Добавь:

oppOut: oppOut,
Конвертер

Добавь в конец файла:

func toExecutorTriangle(t *Triangle) *executor.Triangle {
	if t == nil {
		return nil
	}

	out := &executor.Triangle{
		A: t.A,
		B: t.B,
		C: t.C,
	}

	for i, leg := range t.Legs {
		out.Legs[i] = executor.LegRule{
			Index:       leg.Index,
			RawLeg:      leg.RawLeg,
			Step:        leg.Step,
			MinQty:      leg.MinQty,
			MinNotional: leg.MinNotional,

			Symbol:      leg.Symbol,
			Side:        leg.Side,
			Base:        leg.Base,
			Quote:       leg.Quote,
			QtyStep:     leg.QtyStep,
			QuoteStep:   leg.QuoteStep,
			PriceStep:   leg.PriceStep,
			LegMinQty:   leg.LegMinQty,
			LegMinQuote: leg.LegMinQuote,
			LegMinNotnl: leg.LegMinNotnl,

			Key: leg.Key,
		}
	}

	return out
}
Отправка opportunity в calcTriangle()

Сразу после:

profitUSDT := finalAmount - maxStart
profitPct := profitUSDT / maxStart
strength := computeOpportunityStrength(profitPct, maxStart, spreadAge, maxAge)
triName := fmt.Sprintf("%s->%s->%s", tri.A, tri.B, tri.C)

вставь:

	if c.oppOut != nil {
		op := &executor.Opportunity{
			Exchange: md.Exchange,
			Triangle: toExecutorTriangle(tri),
			AnchorTS: anchorTS,
			Quotes: [3]queue.Quote{
				q[0],
				q[1],
				q[2],
			},
			AgesMS: [3]int64{
				ages[0],
				ages[1],
				ages[2],
			},
			MaxStart:   maxStart,
			ProfitPct:  profitPct,
			ProfitUSDT: profitUSDT,
			FinalUSDT:  finalAmount,
		}

		select {
		case c.oppOut <- op:
		default:
		}
	}
Что получишь после запуска

В логах появятся строки вида:

[Executor] summary accepted=...

и при хороших сигналах:

[EXEC PAPER OK] ...
