main.go

package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"crypt_proto/internal/calculator"
	"crypt_proto/internal/collector"
	"crypt_proto/internal/queue"
	"crypt_proto/pkg/models"
)

func main() {
	logFile, cleanup, err := setupLogging()
	if err != nil {
		log.Fatalf("setup logging: %v", err)
	}
	defer cleanup()
	log.Printf("[Main] log file: %s", logFile)

	go func() {
		log.Println("pprof on http://localhost:6060/debug/pprof/")
		if err := http.ListenAndServe("localhost:6060", nil); err != nil {
			log.Printf("pprof server error: %v", err)
		}
	}()

	out := make(chan *models.MarketData, 100_000)
	mem := queue.NewMemoryStore()

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

	cfg := calculator.DefaultConfig()
	cfg.LogMode = calculator.LogDebug
	cfg.MinVolumeUSDT = 10
	cfg.MinProfitPct = 0
	cfg.QuoteAgeMaxMS = 400
	cfg.StatsEverySec = 5

	calc := calculator.NewCalculator(mem, triangles, cfg)
	go calc.Run(out)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("[Main] shutting down...")
	kc.Stop()
	close(out)
	log.Println("[Main] exited")
}

func setupLogging() (string, func(), error) {
	logDir := "logs"
	if err := os.MkdirAll(logDir, 0o755); err != nil {
		return "", func() {}, err
	}
	filename := filepath.Join(logDir, fmt.Sprintf("arb_%s.log", time.Now().Format("20060110_150405")))
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return "", func() {}, err
	}
	mw := io.MultiWriter(os.Stdout, f)
	log.SetOutput(mw)
	log.SetFlags(log.LstdFlags)
	cleanup := func() {
		_ = f.Close()
	}
	return filename, cleanup, nil
}


calc.go

package calculator

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"crypt_proto/internal/queue"
	"crypt_proto/pkg/models"
)

type Calculator struct {
	mem     *queue.MemoryStore
	scanner *Scanner
	filter  *ExecutorFilter
	log     *log.Logger
	cfg     Config
	stats   Stats
}

func NewCalculator(mem *queue.MemoryStore, triangles []*Triangle, cfg Config) *Calculator {
	if cfg.MinVolumeUSDT <= 0 {
		cfg.MinVolumeUSDT = defaultMinVolumeUSDT
	}
	if cfg.SearchStepUSDT <= 0 {
		cfg.SearchStepUSDT = defaultSearchStep
	}
	if cfg.StatsEverySec <= 0 {
		cfg.StatsEverySec = 5
	}

	writer := log.Writer()
	c := &Calculator{
		mem:     mem,
		scanner: NewScanner(mem, triangles, cfg),
		filter:  NewExecutorFilter(cfg),
		log:     log.New(writer, "", log.LstdFlags),
		cfg:     cfg,
		stats: Stats{
			ScanRejects:     make(map[string]int64),
			ExecRejects:     make(map[string]int64),
			TriangleMetrics: make(map[string]*TriangleMetrics),
		},
	}

	if cfg.LogMode != LogSilent {
		c.log.Printf(
			"[CALC] started | triangles indexed=%d | minVolume=%.2f USDT | minProfit=%.4f%% | quoteAgeMax=%dms | logMode=%s",
			len(triangles),
			cfg.MinVolumeUSDT,
			cfg.MinProfitPct*100,
			cfg.QuoteAgeMaxMS,
			c.logModeString(),
		)
	}

	return c
}

func (c *Calculator) Run(in <-chan *models.MarketData) {
	ticker := time.NewTicker(time.Duration(c.cfg.StatsEverySec) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case md, ok := <-in:
			if !ok {
				if c.cfg.LogMode != LogSilent {
					c.logStats()
				}
				return
			}
			if md == nil {
				continue
			}

			c.stats.Ticks++
			c.mem.Push(md)

			results := c.scanner.CandidatesFor(md.Symbol, md.Timestamp)
			if len(results) == 0 {
				continue
			}

			for _, res := range results {
				c.stats.TrianglesSeen++
				if !res.OK {
					c.addScanReject(res.Reject, res.Candidate.Triangle)
					continue
				}

				c.stats.Candidates++
				opp, reason, ok := c.filter.Evaluate(res.Candidate)
				c.recordTriangleMetrics(res.Candidate, opp)
				if !ok {
					c.addExecReject(reason, res.Candidate.Triangle)
					continue
				}

				c.stats.Opportunities++
				if opp.ProfitPct > 0 {
					c.stats.Positive++
					c.logOpportunity(opp)
					c.stats.Logged++
				} else {
					c.stats.Negative++
				}
			}
		case <-ticker.C:
			if c.cfg.LogMode != LogSilent {
				c.logStats()
			}
		}
	}
}

func (c *Calculator) recordTriangleMetrics(cand ScanCandidate, opp ExecutableOpportunity) {
	if cand.Triangle == nil {
		return
	}
	key := cand.Triangle.Key()
	m := c.stats.TriangleMetrics[key]
	if m == nil {
		m = &TriangleMetrics{
			Key:              key,
			BestEstimatedPct: -1e9,
			BestIdealPct:     -1e9,
			BestRoundedPct:   -1e9,
			BestFinalPct:     -1e9,
			WorstFinalPct:    1e9,
		}
		c.stats.TriangleMetrics[key] = m
	}
	m.Seen++
	m.LastEstimatedPct = cand.EstimatedPct
	m.LastIdealPct = opp.IdealProfitPct
	m.LastRoundedPct = opp.RoundedProfitPct
	m.LastFinalPct = opp.ProfitPct
	if cand.EstimatedPct > m.BestEstimatedPct {
		m.BestEstimatedPct = cand.EstimatedPct
	}
	if opp.IdealProfitPct > m.BestIdealPct {
		m.BestIdealPct = opp.IdealProfitPct
	}
	if opp.RoundedProfitPct > m.BestRoundedPct {
		m.BestRoundedPct = opp.RoundedProfitPct
	}
	if opp.ProfitPct > m.BestFinalPct {
		m.BestFinalPct = opp.ProfitPct
	}
	if opp.ProfitPct < m.WorstFinalPct {
		m.WorstFinalPct = opp.ProfitPct
	}
}

func (c *Calculator) addScanReject(reason string, tri *Triangle) {
	if reason == "" {
		reason = "unknown"
	}
	c.stats.ScanRejects[reason]++
	c.logReject("scan", reason, tri, c.stats.ScanRejects[reason])
}

func (c *Calculator) addExecReject(reason string, tri *Triangle) {
	if reason == "" {
		reason = "unknown"
	}
	c.stats.ExecRejects[reason]++
	c.logReject("exec", reason, tri, c.stats.ExecRejects[reason])
}

func (c *Calculator) logOpportunity(opp ExecutableOpportunity) {
	c.log.Printf(
		"[ARB] %s→%s→%s | est=%.4f%% | ideal=%.4f%% | rounded=%.4f%% | real=%.4f%% | start=%.2f USDT | minStart=%.2f USDT | final=%.4f USDT | profit=%.4f USDT | symbol=%s",
		opp.Triangle.A,
		opp.Triangle.B,
		opp.Triangle.C,
		opp.EstimatedPct*100,
		opp.IdealProfitPct*100,
		opp.RoundedProfitPct*100,
		opp.ProfitPct*100,
		opp.StartUSDT,
		opp.MinStartUSDT,
		opp.FinalUSDT,
		opp.ProfitUSDT,
		opp.TriggeredBy,
	)
}

func (c *Calculator) logReject(stage, reason string, tri *Triangle, count int64) {
	if c.cfg.LogMode != LogDebug {
		return
	}
	if count != 1 && count != 10 && count != 100 && count%1000 != 0 {
		return
	}
	if tri == nil {
		c.log.Printf("[REJECT] stage=%s reason=%s count=%d", stage, reason, count)
		return
	}
	c.log.Printf("[REJECT] stage=%s reason=%s count=%d tri=%s->%s->%s", stage, reason, count, tri.A, tri.B, tri.C)
}

func (c *Calculator) logStats() {
	c.log.Printf(
		"[STATS] ticks=%d triangles_seen=%d cand=%d exec=%d pos=%d neg=%d logged=%d | scan_rejects={%s} | exec_rejects={%s}",
		c.stats.Ticks,
		c.stats.TrianglesSeen,
		c.stats.Candidates,
		c.stats.Opportunities,
		c.stats.Positive,
		c.stats.Negative,
		c.stats.Logged,
		formatCounts(c.stats.ScanRejects),
		formatCounts(c.stats.ExecRejects),
	)
	c.logTopTriangleMetrics(5)
}

func (c *Calculator) logTopTriangleMetrics(n int) {
	if len(c.stats.TriangleMetrics) == 0 {
		return
	}
	items := make([]*TriangleMetrics, 0, len(c.stats.TriangleMetrics))
	for _, m := range c.stats.TriangleMetrics {
		items = append(items, m)
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].BestIdealPct == items[j].BestIdealPct {
			return items[i].Seen > items[j].Seen
		}
		return items[i].BestIdealPct > items[j].BestIdealPct
	})
	if n > len(items) {
		n = len(items)
	}
	parts := make([]string, 0, n)
	for _, m := range items[:n] {
		parts = append(parts, fmt.Sprintf(
			"%s seen=%d best(est=%.4f%% ideal=%.4f%% rounded=%.4f%% final=%.4f%%) last(ideal=%.4f%% rounded=%.4f%% final=%.4f%%)",
			m.Key,
			m.Seen,
			m.BestEstimatedPct*100,
			m.BestIdealPct*100,
			m.BestRoundedPct*100,
			m.BestFinalPct*100,
			m.LastIdealPct*100,
			m.LastRoundedPct*100,
			m.LastFinalPct*100,
		))
	}
	c.log.Printf("[TOP] %s", strings.Join(parts, " | "))
}

func formatCounts(m map[string]int64) string {
	if len(m) == 0 {
		return "none"
	}
	type kv struct {
		k string
		v int64
	}
	items := make([]kv, 0, len(m))
	for k, v := range m {
		items = append(items, kv{k: k, v: v})
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].v == items[j].v {
			return items[i].k < items[j].k
		}
		return items[i].v < items[j].v
	})
	parts := make([]string, 0, len(items))
	for _, item := range items {
		parts = append(parts, fmt.Sprintf("%s=%d", item.k, item.v))
	}
	return strings.Join(parts, ", ")
}

func (c *Calculator) logModeString() string {
	switch c.cfg.LogMode {
	case LogSilent:
		return "silent"
	case LogDebug:
		return "debug"
	default:
		return "normal"
	}
}

func trimFloat(v float64) string {
	s := fmt.Sprintf("%.2f", v)
	s = strings.TrimRight(s, "0")
	s = strings.TrimRight(s, ".")
	if s == "" {
		return "0"
	}
	return s
}

func ensureDir(path string) {
	if path == "" {
		return
	}
	_ = os.MkdirAll(path, 0o755)
}


filter.go

package calculator

import (
	"log"
	"math"

	"crypt_proto/internal/queue"
)

type ExecutorFilter struct {
	cfg Config
}

func NewExecutorFilter(cfg Config) *ExecutorFilter {
	return &ExecutorFilter{cfg: cfg}
}

func (f *ExecutorFilter) Evaluate(cand ScanCandidate) (ExecutableOpportunity, string, bool) {
	minStart, ok := findMinStartForTriangle(cand.Triangle, cand.Quotes, f.cfg.MinVolumeUSDT, cand.MaxStartUSDT, f.cfg.SearchStepUSDT)
	if !ok {
		return ExecutableOpportunity{}, "cannot_find_valid_start", false
	}

	startUSDT := floorToStep(math.Max(f.cfg.MinVolumeUSDT, minStart), f.cfg.SearchStepUSDT)
	if startUSDT < f.cfg.MinVolumeUSDT || startUSDT > cand.MaxStartUSDT {
		return ExecutableOpportunity{}, "max_start_lt_min_volume", false
	}

	idealState, okIdeal := simulateTriangleMode(startUSDT, cand.Triangle, cand.Quotes, true, true)
	roundedState, okRounded := simulateTriangleMode(startUSDT, cand.Triangle, cand.Quotes, false, true)
	finalState, okFinal := simulateTriangleMode(startUSDT, cand.Triangle, cand.Quotes, false, false)
	if !okIdeal || !okRounded || !okFinal {
		return ExecutableOpportunity{}, "simulate_failed", false
	}

	opp := ExecutableOpportunity{
		Triangle:         cand.Triangle,
		Quotes:           cand.Quotes,
		EstimatedPct:     cand.EstimatedPct,
		StartUSDT:        finalState.StartUSDT,
		MinStartUSDT:     minStart,
		FinalUSDT:        finalState.FinalUSDT,
		ProfitUSDT:       finalState.ProfitUSDT,
		ProfitPct:        finalState.ProfitPct,
		TriggeredBy:      cand.TriggeredBy,
		TriggeredAtMS:    cand.TriggeredAtMS,
		IdealFinalUSDT:   idealState.FinalUSDT,
		IdealProfitPct:   idealState.ProfitPct,
		RoundedFinalUSDT: roundedState.FinalUSDT,
		RoundedProfitPct: roundedState.ProfitPct,
	}

	if f.cfg.LogMode == LogDebug {
		log.Printf(
			"[EXEC CMP] %s→%s→%s | est=%.4f%% | ideal=%.4f%% | rounded=%.4f%% | final=%.4f%%",
			cand.Triangle.A,
			cand.Triangle.B,
			cand.Triangle.C,
			cand.EstimatedPct*100,
			opp.IdealProfitPct*100,
			opp.RoundedProfitPct*100,
			opp.ProfitPct*100,
		)
	}

	if opp.ProfitPct < f.cfg.MinProfitPct {
		return opp, "profit_below_threshold", false
	}

	return opp, "", true
}

func findMinStartForTriangle(tri *Triangle, q [3]queue.Quote, lowerBound, upperBound, searchStep float64) (float64, bool) {
	if upperBound <= 0 || upperBound+1e-12 < lowerBound {
		return 0, false
	}

	lo := math.Max(searchStep, floorToStep(lowerBound, searchStep))
	if lo < lowerBound {
		lo += searchStep
	}
	if lo > upperBound {
		return 0, false
	}

	if _, ok := simulateTriangle(lo, tri, q); ok {
		return lo, true
	}

	step := searchStep
	high := lo
	for high <= upperBound {
		next := floorToStep(high+step, searchStep)
		if next <= high {
			next = floorToStep(high+searchStep, searchStep)
		}
		if next > upperBound {
			break
		}
		if _, ok := simulateTriangle(next, tri, q); ok {
			left, right := high, next
			for right-left > searchStep+1e-12 {
				mid := floorToStep((left+right)/2, searchStep)
				if mid <= left {
					mid = floorToStep(left+searchStep, searchStep)
				}
				if mid >= right {
					break
				}
				if _, ok := simulateTriangle(mid, tri, q); ok {
					right = mid
				} else {
					left = mid
				}
			}
			return right, true
		}
		high = next
		step *= 2
	}

	return 0, false
}

func simulateTriangle(startUSDT float64, tri *Triangle, q [3]queue.Quote) (ExecutionResult, bool) {
	return simulateTriangleMode(startUSDT, tri, q, false, false)
}

func simulateTriangleMode(startUSDT float64, tri *Triangle, q [3]queue.Quote, ignoreFees bool, ignoreRounding bool) (ExecutionResult, bool) {
	state := ExecutionResult{StartUSDT: startUSDT}
	amount := startUSDT

	for i := 0; i < 3; i++ {
		nextAmount, notional, ok := simulateLegMode(amount, tri.Legs[i], tri.Rules[i], q[i], ignoreFees, ignoreRounding)
		if !ok {
			return ExecutionResult{}, false
		}
		state.LegNotional[i] = notional
		state.LegAmount[i] = nextAmount
		amount = nextAmount
	}

	state.FinalUSDT = amount
	state.ProfitUSDT = state.FinalUSDT - state.StartUSDT
	if state.StartUSDT <= 0 {
		return ExecutionResult{}, false
	}
	state.ProfitPct = state.ProfitUSDT / state.StartUSDT
	return state, true
}

func simulateLeg(inputAmount float64, leg LegIndex, rules LegRules, quote queue.Quote) (float64, float64, bool) {
	return simulateLegMode(inputAmount, leg, rules, quote, false, false)
}

func simulateLegMode(inputAmount float64, leg LegIndex, rules LegRules, quote queue.Quote, ignoreFees bool, ignoreRounding bool) (float64, float64, bool) {
	mul := feeMultiplier(rules.Fee)
	if ignoreFees {
		mul = 1
	}

	if leg.IsBuy {
		if quote.Ask <= 0 || quote.AskSize <= 0 || inputAmount <= 0 {
			return 0, 0, false
		}

		qty := inputAmount / quote.Ask
		if !ignoreRounding {
			qty = applyFloorStep(qty, rules.QtyStep)
		}
		if qty <= 0 {
			return 0, 0, false
		}
		if qty > quote.AskSize {
			qty = quote.AskSize
			if !ignoreRounding {
				qty = applyFloorStep(qty, rules.QtyStep)
			}
		}
		if qty <= 0 {
			return 0, 0, false
		}

		notional := qty * quote.Ask
		if !ignoreRounding {
			notional = applyFloorStep(notional, rules.QuoteStep)
			if !passesMinChecks(qty, notional, rules) {
				return 0, 0, false
			}
		}

		outQty := qty * mul
		if !ignoreRounding {
			outQty = applyFloorStep(outQty, rules.QtyStep)
		}
		if outQty <= 0 {
			return 0, 0, false
		}
		return outQty, notional, true
	}

	if quote.Bid <= 0 || quote.BidSize <= 0 || inputAmount <= 0 {
		return 0, 0, false
	}

	qty := inputAmount
	if !ignoreRounding {
		qty = applyFloorStep(qty, rules.QtyStep)
	}
	if qty <= 0 {
		return 0, 0, false
	}
	if qty > quote.BidSize {
		qty = quote.BidSize
		if !ignoreRounding {
			qty = applyFloorStep(qty, rules.QtyStep)
		}
	}
	if qty <= 0 {
		return 0, 0, false
	}

	notional := qty * quote.Bid
	if !ignoreRounding {
		notional = applyFloorStep(notional, rules.QuoteStep)
		if !passesMinChecks(qty, notional, rules) {
			return 0, 0, false
		}
	}

	outQuote := notional * mul
	if !ignoreRounding {
		outQuote = applyFloorStep(outQuote, rules.QuoteStep)
	}
	if outQuote <= 0 {
		return 0, 0, false
	}
	return outQuote, notional, true
}


types.go

package calculator

import (
	"fmt"
	"math"
	"strings"

	"crypt_proto/internal/queue"
)

const (
	defaultTakerFee      = 0.001
	defaultMinVolumeUSDT = 50.0
	defaultMinProfitPct  = 0.001
	defaultSearchStep    = 0.01
)

var triangleLegColumns = [3]int{3, 7, 11}

type LogMode int

const (
	LogSilent LogMode = iota
	LogNormal
	LogDebug
)

type Config struct {
	MinVolumeUSDT  float64
	MinProfitPct   float64
	SearchStepUSDT float64
	QuoteAgeMaxMS  int64
	StatsEverySec  int
	LogMode        LogMode
}

func DefaultConfig() Config {
	return Config{
		MinVolumeUSDT:  defaultMinVolumeUSDT,
		MinProfitPct:   defaultMinProfitPct,
		SearchStepUSDT: defaultSearchStep,
		QuoteAgeMaxMS:  2500,
		StatsEverySec:  5,
		LogMode:        LogNormal,
	}
}

type LegIndex struct {
	Key    string
	Symbol string
	IsBuy  bool
}

type LegRules struct {
	Symbol      string
	Side        string
	Base        string
	Quote       string
	QtyStep     float64
	QuoteStep   float64
	PriceStep   float64
	MinQty      float64
	MinQuote    float64
	MinNotional float64
	Fee         float64
}

type Triangle struct {
	A, B, C string
	Legs    [3]LegIndex
	Rules   [3]LegRules
}

func (t *Triangle) Key() string {
	if t == nil {
		return "<nil>"
	}
	return t.A + "->" + t.B + "->" + t.C
}

type ScanCandidate struct {
	Triangle      *Triangle
	Quotes        [3]queue.Quote
	EstimatedPct  float64
	MaxStartUSDT  float64
	TriggeredBy   string
	TriggeredAtMS int64
}

type ExecutionResult struct {
	StartUSDT   float64
	MinStart    float64
	FinalUSDT   float64
	ProfitUSDT  float64
	ProfitPct   float64
	LegNotional [3]float64
	LegAmount   [3]float64
}

type ExecutableOpportunity struct {
	Triangle      *Triangle
	Quotes        [3]queue.Quote
	EstimatedPct  float64
	StartUSDT     float64
	MinStartUSDT  float64
	FinalUSDT     float64
	ProfitUSDT    float64
	ProfitPct     float64
	TriggeredBy   string
	TriggeredAtMS int64

	IdealFinalUSDT   float64
	IdealProfitPct   float64
	RoundedFinalUSDT float64
	RoundedProfitPct float64
}

type TriangleMetrics struct {
	Key              string
	Seen             int64
	BestEstimatedPct float64
	BestIdealPct     float64
	BestRoundedPct   float64
	BestFinalPct     float64
	WorstFinalPct    float64
	LastEstimatedPct float64
	LastIdealPct     float64
	LastRoundedPct   float64
	LastFinalPct     float64
}

type ScanResult struct {
	Candidate ScanCandidate
	Reject    string
	OK        bool
}

type Stats struct {
	Ticks           int64
	TrianglesSeen   int64
	Candidates      int64
	Opportunities   int64
	Positive        int64
	Negative        int64
	Logged          int64
	ScanRejects     map[string]int64
	ExecRejects     map[string]int64
	TriangleMetrics map[string]*TriangleMetrics
}

func feeMultiplier(fee float64) float64 {
	if fee >= 0 && fee < 1 {
		return 1 - fee
	}
	return 1 - defaultTakerFee
}

func applyFloorStep(value, step float64) float64 {
	if value <= 0 {
		return 0
	}
	if step <= 0 {
		return value
	}
	return floorToStep(value, step)
}

func floorToStep(value, step float64) float64 {
	if step <= 0 {
		return value
	}
	units := math.Floor((value + 1e-12) / step)
	if units <= 0 {
		return 0
	}
	result := units * step
	precision := decimalsFromStep(step)
	pow := math.Pow10(precision)
	return math.Floor(result*pow+1e-9) / pow
}

func decimalsFromStep(step float64) int {
	if step <= 0 {
		return 8
	}
	s := strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.12f", step), "0"), ".")
	idx := strings.IndexByte(s, '.')
	if idx == -1 {
		return 0
	}
	return len(s) - idx - 1
}

func passesMinChecks(qty, notional float64, rules LegRules) bool {
	if rules.MinQty > 0 && qty+1e-12 < rules.MinQty {
		return false
	}
	if rules.MinQuote > 0 && notional+1e-12 < rules.MinQuote {
		return false
	}
	if rules.MinNotional > 0 && notional+1e-12 < rules.MinNotional {
		return false
	}
	return true
}


