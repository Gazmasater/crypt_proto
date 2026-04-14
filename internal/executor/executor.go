package executor

import (
	"fmt"
	"log"
	"math"
	"os"
	"strings"
	"sync"
	"time"

	"crypt_proto/internal/collector"
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
	Books    [3]collector.BookSnapshot
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
	RejectBadInputProfit     RejectReason = "bad_input_profit"
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

	// Что вообще пускать в executor из calculator.
	// Всё что ниже этого значения - сразу reject без лишней работы.
	MinInputProfitPct float64

	DedupWindow  time.Duration
	SummaryEvery time.Duration
}

func DefaultConfig() Config {
	return Config{
		SafeStartFactor:   0.70,
		MinSafeStart:      20.0,
		MaxQuoteAgeMS:     300,
		MaxAgeSpreadMS:    200,
		MinProfitUSDT:     0.001,
		MinInputProfitPct: 0.0, // executor по умолчанию ждёт только неотрицательные окна
		DedupWindow:       250 * time.Millisecond,
		SummaryEvery:      10 * time.Second,
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

	// Главная правка: executor не должен тратить время на явный минус.
	if op.ProfitPct < e.cfg.MinInputProfitPct {
		e.reject(RejectBadInputProfit, op, fmt.Sprintf("input_profit_pct=%.8f", op.ProfitPct))
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
	if ok && hasAllBooks(op.Books) {
		if depthFinal, depthDiag, depthOK := simulateTriangleDepth(safeStart, op.Triangle, op.Books); depthOK && depthFinal > 0 {
			finalAmount = depthFinal
			diag = depthDiag
		} else {
			ok = false
		}
	}
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

func hasAllBooks(books [3]collector.BookSnapshot) bool {
	for _, b := range books {
		if len(b.Bids) == 0 || len(b.Asks) == 0 {
			return false
		}
	}
	return true
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
		op.AnchorTS/100, // корзина 100мс
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

	// Не шумим в stdout по каждому reject, только в файл.
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
		"[Executor] summary accepted=%d reject_duplicate=%d reject_bad_input_profit=%d reject_stale=%d reject_desync=%d reject_small_safe=%d reject_non_positive=%d reject_too_small_profit=%d reject_sim_failed=%d",
		e.accepted,
		e.rejects[RejectDuplicate],
		e.rejects[RejectBadInputProfit],
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

func simulateTriangleDepth(startUSDT float64, tri *Triangle, books [3]collector.BookSnapshot) (float64, [3]legExecution, bool) {
	var diag [3]legExecution
	amount := startUSDT
	for i := 0; i < 3; i++ {
		out, d, ok := executeLegByDepth(amount, tri.Legs[i], books[i])
		if !ok {
			return 0, diag, false
		}
		diag[i] = d
		amount = out
	}
	return amount, diag, true
}

func executeLegByDepth(in float64, leg LegRule, book collector.BookSnapshot) (float64, legExecution, bool) {
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
		if len(book.Asks) == 0 {
			return 0, legExecution{}, false
		}
		remainingQuote := in
		totalQty := 0.0
		totalNotional := 0.0
		bookLimitIn := 0.0
		for _, lvl := range book.Asks {
			if lvl.Price <= 0 || lvl.Size <= 0 {
				continue
			}
			bookLimitIn += lvl.Price * lvl.Size
			if remainingQuote <= eps() {
				break
			}
			maxSpend := lvl.Price * lvl.Size
			spend := math.Min(remainingQuote, maxSpend)
			qty := spend / lvl.Price
			totalQty += qty
			totalNotional += qty * lvl.Price
			remainingQuote -= qty * lvl.Price
		}
		tradeQty := floorToStep(totalQty, qtyStep)
		if tradeQty <= 0 {
			return 0, legExecution{}, false
		}
		avgPrice := totalNotional / totalQty
		tradeNotional := tradeQty * avgPrice
		if tradeNotional <= 0 || tradeNotional > in+1e-9 {
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
		return outBase, legExecution{In: in, Out: outBase, Price: avgPrice, BookLimitIn: bookLimitIn, TradeQty: tradeQty, TradeNotional: tradeNotional}, true
	case "SELL":
		if len(book.Bids) == 0 {
			return 0, legExecution{}, false
		}
		remainingBase := in
		totalQty := 0.0
		totalNotional := 0.0
		bookLimitIn := 0.0
		for _, lvl := range book.Bids {
			if lvl.Price <= 0 || lvl.Size <= 0 {
				continue
			}
			bookLimitIn += lvl.Size
			if remainingBase <= eps() {
				break
			}
			qty := math.Min(remainingBase, lvl.Size)
			totalQty += qty
			totalNotional += qty * lvl.Price
			remainingBase -= qty
		}
		tradeQty := floorToStep(totalQty, qtyStep)
		if tradeQty <= 0 {
			return 0, legExecution{}, false
		}
		avgPrice := totalNotional / totalQty
		tradeNotional := tradeQty * avgPrice
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
		return outQuote, legExecution{In: in, Out: outQuote, Price: avgPrice, BookLimitIn: bookLimitIn, TradeQty: tradeQty, TradeNotional: tradeNotional}, true
	}
	return 0, legExecution{}, false
}
