package calculator

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"crypt_proto/internal/executor"
	"crypt_proto/internal/queue"
	"crypt_proto/pkg/models"
)

const (
	feeM                 = 0.9992
	defaultMaxQuoteAgeMS = int64(300)
)

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

type metricsConfig struct {
	MinProfitPct     float64
	MinVolumeUSDT    float64
	MaxAgeMS         int64
	NearProfitPct    float64
	SummaryEvery     time.Duration
	DedupWindow      time.Duration
	FlushEveryWrites int
}

type metricsWriter struct {
	mu          sync.Mutex
	file        *os.File
	csv         *csv.Writer
	enabled     bool
	writeCount  int
	lastKey     string
	lastWriteAt time.Time
	cfg         metricsConfig
}

type metricsSummary struct {
	mu         sync.Mutex
	startedAt  time.Time
	lastLogAt  time.Time
	checked    uint64
	written    uint64
	profitable uint64
	bestPct    float64
	bestUSDT   float64
	bestTri    string
}

func newMetricsSummary() *metricsSummary {
	now := time.Now()
	return &metricsSummary{
		startedAt: now,
		lastLogAt: now,
		bestPct:   math.Inf(-1),
		bestUSDT:  math.Inf(-1),
	}
}

func (ms *metricsSummary) Observe(tri string, profitPct, profitUSDT float64, written bool) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	ms.checked++
	if written {
		ms.written++
	}
	if profitPct > 0 {
		ms.profitable++
	}
	if profitPct > ms.bestPct || (profitPct == ms.bestPct && profitUSDT > ms.bestUSDT) {
		ms.bestPct = profitPct
		ms.bestUSDT = profitUSDT
		ms.bestTri = tri
	}
}

func (ms *metricsSummary) MaybeLog(now time.Time, every time.Duration) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	if every <= 0 || now.Sub(ms.lastLogAt) < every {
		return
	}
	ms.lastLogAt = now

	bestPct := 0.0
	bestUSDT := 0.0
	bestTri := ""
	if ms.checked > 0 && !math.IsInf(ms.bestPct, -1) {
		bestPct = ms.bestPct * 100
		bestUSDT = ms.bestUSDT
		bestTri = ms.bestTri
	}

	log.Printf("[Calculator] summary checked=%d written=%d profitable=%d best_pct=%.4f%% best_usdt=%.6f best_tri=%s\n",
		ms.checked, ms.written, ms.profitable, bestPct, bestUSDT, bestTri)
}

func defaultMetricsConfig() metricsConfig {
	return metricsConfig{
		MinProfitPct:     0.0,
		MinVolumeUSDT:    50.0,
		MaxAgeMS:         defaultMaxQuoteAgeMS,
		NearProfitPct:    -0.0002, // -0.02%
		SummaryEvery:     10 * time.Second,
		DedupWindow:      250 * time.Millisecond,
		FlushEveryWrites: 1,
	}
}

func newMetricsWriter(path string, cfg metricsConfig) (*metricsWriter, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return nil, err
	}

	mw := &metricsWriter{
		file:    f,
		csv:     csv.NewWriter(f),
		enabled: true,
		cfg:     cfg,
	}

	stat, err := f.Stat()
	if err != nil {
		_ = f.Close()
		return nil, err
	}
	if stat.Size() == 0 {
		header := []string{
			"ts_unix_ms", "tri", "A", "B", "C",
			"profit_pct", "profit_usdt", "volume_usdt", "final_usdt",
			"opportunity_strength", "age_min_ms", "age_max_ms", "age_spread_ms",
			"leg1_symbol", "leg1_side", "leg1_bid", "leg1_ask", "leg1_bid_size", "leg1_ask_size", "leg1_age_ms", "leg1_in", "leg1_out", "leg1_trade_qty", "leg1_trade_notional", "leg1_book_limit_in",
			"leg2_symbol", "leg2_side", "leg2_bid", "leg2_ask", "leg2_bid_size", "leg2_ask_size", "leg2_age_ms", "leg2_in", "leg2_out", "leg2_trade_qty", "leg2_trade_notional", "leg2_book_limit_in",
			"leg3_symbol", "leg3_side", "leg3_bid", "leg3_ask", "leg3_bid_size", "leg3_ask_size", "leg3_age_ms", "leg3_in", "leg3_out", "leg3_trade_qty", "leg3_trade_notional", "leg3_book_limit_in",
		}
		if err := mw.csv.Write(header); err != nil {
			_ = f.Close()
			return nil, err
		}
		mw.csv.Flush()
		if err := mw.csv.Error(); err != nil {
			_ = f.Close()
			return nil, err
		}
	}
	return mw, nil
}

func (mw *metricsWriter) Close() error {
	if mw == nil || !mw.enabled {
		return nil
	}
	mw.mu.Lock()
	defer mw.mu.Unlock()

	mw.csv.Flush()
	if err := mw.csv.Error(); err != nil {
		_ = mw.file.Close()
		return err
	}
	return mw.file.Close()
}

// Прибыльные сигналы пишем всегда.
// Для near-profit сигналов применяем фильтры и dedup.
func (mw *metricsWriter) ShouldWrite(anchorTS int64, tri string, profitPct, volumeUSDT float64, maxAgeMS int64) bool {
	if mw == nil || !mw.enabled {
		return false
	}

	// Всегда пишем прибыльные окна.
	if profitPct >= mw.cfg.MinProfitPct {
		return true
	}

	// Ниже — только почти-прибыльные кандидаты.
	if volumeUSDT < mw.cfg.MinVolumeUSDT {
		return false
	}
	if maxAgeMS < 0 || maxAgeMS > mw.cfg.MaxAgeMS {
		return false
	}
	if profitPct < mw.cfg.NearProfitPct {
		return false
	}

	now := time.UnixMilli(anchorTS)
	key := fmt.Sprintf("%s|%.4f|%.2f", tri, profitPct*100, volumeUSDT)

	mw.mu.Lock()
	defer mw.mu.Unlock()

	if mw.lastKey == key && now.Sub(mw.lastWriteAt) < mw.cfg.DedupWindow {
		return false
	}
	mw.lastKey = key
	mw.lastWriteAt = now
	return true
}

func (mw *metricsWriter) Write(record []string) error {
	if mw == nil || !mw.enabled {
		return nil
	}

	mw.mu.Lock()
	defer mw.mu.Unlock()

	if err := mw.csv.Write(record); err != nil {
		return err
	}
	mw.writeCount++

	if mw.cfg.FlushEveryWrites <= 1 || mw.writeCount%mw.cfg.FlushEveryWrites == 0 {
		mw.csv.Flush()
		return mw.csv.Error()
	}
	return nil
}

type Calculator struct {
	mem           *queue.MemoryStore
	bySymbol      map[string][]*Triangle
	fileLog       *log.Logger
	metrics       *metricsWriter
	summary       *metricsSummary
	maxQuoteAgeMS int64
	metricsCfg    metricsConfig
	oppOut        chan<- *executor.Opportunity
}

func NewCalculator(mem *queue.MemoryStore, triangles []*Triangle, oppOut chan<- *executor.Opportunity) *Calculator {
	f, err := os.OpenFile("arb_opportunities.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		log.Fatalf("failed to open log: %v", err)
	}

	cfg := defaultMetricsConfig()
	metrics, err := newMetricsWriter("arb_metrics.csv", cfg)
	if err != nil {
		log.Fatalf("failed to open metrics file: %v", err)
	}

	bySymbol := make(map[string][]*Triangle, 1024)
	for _, t := range triangles {
		for _, leg := range t.Legs {
			if leg.Symbol == "" {
				continue
			}
			bySymbol[leg.Symbol] = append(bySymbol[leg.Symbol], t)
		}
	}

	log.Printf("[Calculator] indexed %d symbols\n", len(bySymbol))

	return &Calculator{
		mem:           mem,
		bySymbol:      bySymbol,
		fileLog:       log.New(f, "", log.LstdFlags),
		metrics:       metrics,
		summary:       newMetricsSummary(),
		maxQuoteAgeMS: defaultMaxQuoteAgeMS,
		metricsCfg:    cfg,
		oppOut:        oppOut,
	}
}

func (c *Calculator) Run(in <-chan *models.MarketData) {
	defer func() {
		if c.metrics != nil {
			_ = c.metrics.Close()
		}
	}()

	for md := range in {
		if md == nil {
			continue
		}
		c.mem.Push(md)

		if c.summary != nil {
			c.summary.MaybeLog(time.Now(), c.metricsCfg.SummaryEvery)
		}

		tris := c.bySymbol[md.Symbol]
		if len(tris) == 0 {
			continue
		}

		for _, tri := range tris {
			c.calcTriangle(md, tri)
		}
	}
}

func (c *Calculator) calcTriangle(md *models.MarketData, tri *Triangle) {
	anchorTS := md.Timestamp
	if anchorTS <= 0 {
		return
	}

	var q [3]queue.Quote
	for i, leg := range tri.Legs {
		quote, ok := c.mem.GetLatestBefore(md.Exchange, leg.Symbol, anchorTS, c.maxQuoteAgeMS)
		if !ok {
			return
		}
		if quote.Bid <= 0 || quote.Ask <= 0 || quote.BidSize <= 0 || quote.AskSize <= 0 {
			return
		}
		q[i] = quote
	}

	ages := [3]int64{
		quoteLagMS(anchorTS, q[0]),
		quoteLagMS(anchorTS, q[1]),
		quoteLagMS(anchorTS, q[2]),
	}
	minAge, maxAge, spreadAge := minMaxSpread(ages)

	maxStart, ok := computeMaxStartTopOfBook(tri, q)
	if !ok || maxStart <= 0 {
		return
	}

	if maxStart < 50 {
		return
	}

	finalAmount, diag, ok := simulateTriangle(maxStart, tri, q)
	if !ok || finalAmount <= 0 {
		return
	}

	profitUSDT := finalAmount - maxStart
	profitPct := profitUSDT / maxStart

	if profitPct < 0 {
		return
	}

	strength := computeOpportunityStrength(profitPct, maxStart, spreadAge, maxAge)
	triName := fmt.Sprintf("%s->%s->%s", tri.A, tri.B, tri.C)

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

	written := false

	shouldWrite := c.metrics != nil && c.metrics.ShouldWrite(anchorTS, triName, profitPct, maxStart, maxAge)
	if shouldWrite {
		record := []string{
			strconv.FormatInt(anchorTS, 10),
			triName,
			tri.A, tri.B, tri.C,
			fmtFloat(profitPct), fmtFloat(profitUSDT), fmtFloat(maxStart), fmtFloat(finalAmount),
			fmtFloat(strength), strconv.FormatInt(minAge, 10), strconv.FormatInt(maxAge, 10), strconv.FormatInt(spreadAge, 10),

			tri.Legs[0].Symbol, tri.Legs[0].Side,
			fmtFloat(q[0].Bid), fmtFloat(q[0].Ask), fmtFloat(q[0].BidSize), fmtFloat(q[0].AskSize),
			strconv.FormatInt(ages[0], 10),
			fmtFloat(diag[0].In), fmtFloat(diag[0].Out), fmtFloat(diag[0].TradeQty), fmtFloat(diag[0].TradeNotional), fmtFloat(diag[0].BookLimitIn),

			tri.Legs[1].Symbol, tri.Legs[1].Side,
			fmtFloat(q[1].Bid), fmtFloat(q[1].Ask), fmtFloat(q[1].BidSize), fmtFloat(q[1].AskSize),
			strconv.FormatInt(ages[1], 10),
			fmtFloat(diag[1].In), fmtFloat(diag[1].Out), fmtFloat(diag[1].TradeQty), fmtFloat(diag[1].TradeNotional), fmtFloat(diag[1].BookLimitIn),

			tri.Legs[2].Symbol, tri.Legs[2].Side,
			fmtFloat(q[2].Bid), fmtFloat(q[2].Ask), fmtFloat(q[2].BidSize), fmtFloat(q[2].AskSize),
			strconv.FormatInt(ages[2], 10),
			fmtFloat(diag[2].In), fmtFloat(diag[2].Out), fmtFloat(diag[2].TradeQty), fmtFloat(diag[2].TradeNotional), fmtFloat(diag[2].BookLimitIn),
		}

		if err := c.metrics.Write(record); err != nil {
			log.Printf("[Calculator] metrics write error: %v", err)
		} else {
			written = true
		}
	}

	if c.summary != nil {
		c.summary.Observe(triName, profitPct, profitUSDT, written)
	}

	if profitPct > 0.0 && maxStart > 50 {
		msg := fmt.Sprintf(
			"[ARB] %s→%s→%s | %.4f%% | volume=%.2f USDT | profit=%.6f USDT | "+
				"anchor=%d | l1=%s %s out=%.8f age=%dms | "+
				"l2=%s %s out=%.8f age=%dms | "+
				"l3=%s %s out=%.8f age=%dms",
			tri.A, tri.B, tri.C,
			profitPct*100, maxStart, profitUSDT,
			anchorTS,
			tri.Legs[0].Symbol, tri.Legs[0].Side, diag[0].Out, ages[0],
			tri.Legs[1].Symbol, tri.Legs[1].Side, diag[1].Out, ages[1],
			tri.Legs[2].Symbol, tri.Legs[2].Side, diag[2].Out, ages[2],
		)
		log.Println(msg)
		c.fileLog.Println(msg)
	}
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
	if side == "" {
		side = detectSideFromRawLeg(leg.RawLeg)
	}
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

func computeMaxStartTopOfBook(tri *Triangle, q [3]queue.Quote) (float64, bool) {
	kIn := 1.0
	maxStart := math.MaxFloat64

	for i := 0; i < 3; i++ {
		leg := tri.Legs[i]
		side := strings.ToUpper(strings.TrimSpace(leg.Side))
		if side == "" {
			side = detectSideFromRawLeg(leg.RawLeg)
		}
		if side != "BUY" && side != "SELL" {
			return 0, false
		}

		var limitIn float64
		switch side {
		case "BUY":
			if q[i].Ask <= 0 || q[i].AskSize <= 0 {
				return 0, false
			}
			limitIn = q[i].Ask * q[i].AskSize
			maxByThis := limitIn / kIn
			if maxByThis < maxStart {
				maxStart = maxByThis
			}
			kIn *= (1.0 / q[i].Ask) * feeM

		case "SELL":
			if q[i].Bid <= 0 || q[i].BidSize <= 0 {
				return 0, false
			}
			limitIn = q[i].BidSize
			maxByThis := limitIn / kIn
			if maxByThis < maxStart {
				maxStart = maxByThis
			}
			kIn *= q[i].Bid * feeM
		}

		if kIn <= 0 {
			return 0, false
		}
	}

	if maxStart <= 0 || !isFinite(maxStart) {
		return 0, false
	}
	return maxStart, true
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
	if len(rows) < 2 {
		return nil, nil
	}

	header := make(map[string]int, len(rows[0]))
	for i, col := range rows[0] {
		header[strings.TrimSpace(col)] = i
	}

	var res []*Triangle
	for _, row := range rows[1:] {
		if len(strings.TrimSpace(strings.Join(row, ""))) == 0 {
			continue
		}

		t := &Triangle{
			A: getString(row, header, "A"),
			B: getString(row, header, "B"),
			C: getString(row, header, "C"),
		}

		for i := 1; i <= 3; i++ {
			idx := i - 1
			leg := LegRule{
				Index:       idx,
				RawLeg:      getString(row, header, fmt.Sprintf("Leg%d", i)),
				Step:        getFloat(row, header, fmt.Sprintf("Step%d", i)),
				MinQty:      getFloat(row, header, fmt.Sprintf("MinQty%d", i)),
				MinNotional: getFloat(row, header, fmt.Sprintf("MinNotional%d", i)),
				Symbol:      getString(row, header, fmt.Sprintf("Leg%dSymbol", i)),
				Side:        strings.ToUpper(getString(row, header, fmt.Sprintf("Leg%dSide", i))),
				Base:        getString(row, header, fmt.Sprintf("Leg%dBase", i)),
				Quote:       getString(row, header, fmt.Sprintf("Leg%dQuote", i)),
				QtyStep:     getFloat(row, header, fmt.Sprintf("Leg%dQtyStep", i)),
				QuoteStep:   getFloat(row, header, fmt.Sprintf("Leg%dQuoteStep", i)),
				PriceStep:   getFloat(row, header, fmt.Sprintf("Leg%dPriceStep", i)),
				LegMinQty:   getFloat(row, header, fmt.Sprintf("Leg%dMinQty", i)),
				LegMinQuote: getFloat(row, header, fmt.Sprintf("Leg%dMinQuote", i)),
				LegMinNotnl: getFloat(row, header, fmt.Sprintf("Leg%dMinNotional", i)),
			}
			if leg.Symbol == "" {
				leg.Symbol = symbolFromRawLeg(leg.RawLeg)
			}
			if leg.Side == "" {
				leg.Side = detectSideFromRawLeg(leg.RawLeg)
			}
			if leg.Symbol == "" || leg.Side == "" {
				continue
			}
			leg.Key = "KuCoin|" + leg.Symbol
			t.Legs[idx] = leg
		}

		if t.Legs[0].Symbol == "" || t.Legs[1].Symbol == "" || t.Legs[2].Symbol == "" {
			continue
		}
		res = append(res, t)
	}
	return res, nil
}

func symbolFromRawLeg(raw string) string {
	raw = strings.ToUpper(strings.TrimSpace(raw))
	if raw == "" {
		return ""
	}
	parts := strings.Fields(raw)
	if len(parts) != 2 {
		return ""
	}
	pair := strings.Split(parts[1], "/")
	if len(pair) != 2 {
		return ""
	}
	return strings.TrimSpace(pair[0]) + "-" + strings.TrimSpace(pair[1])
}

func detectSideFromRawLeg(raw string) string {
	raw = strings.ToUpper(strings.TrimSpace(raw))
	if strings.HasPrefix(raw, "BUY ") {
		return "BUY"
	}
	if strings.HasPrefix(raw, "SELL ") {
		return "SELL"
	}
	return ""
}

func getString(row []string, header map[string]int, key string) string {
	idx, ok := header[key]
	if !ok || idx < 0 || idx >= len(row) {
		return ""
	}
	return strings.TrimSpace(row[idx])
}

func getFloat(row []string, header map[string]int, key string) float64 {
	s := getString(row, header, key)
	if s == "" {
		return 0
	}
	v, err := strconv.ParseFloat(strings.ReplaceAll(s, ",", "."), 64)
	if err != nil {
		return 0
	}
	return v
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

func firstPositive(vals ...float64) float64 {
	for _, v := range vals {
		if v > 0 {
			return v
		}
	}
	return 0
}

func quoteLagMS(anchorTS int64, q queue.Quote) int64 {
	if q.Timestamp <= 0 || anchorTS <= 0 {
		return -1
	}
	if q.Timestamp > anchorTS {
		return 0
	}
	return anchorTS - q.Timestamp
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

func computeOpportunityStrength(profitPct, volumeUSDT float64, ageSpreadMS, maxAgeMS int64) float64 {
	base := profitPct * math.Log1p(math.Max(volumeUSDT, 0))
	freshPenalty := 1.0 / (1.0 + math.Max(float64(maxAgeMS), 0)/180.0)
	spreadPenalty := 1.0 / (1.0 + math.Max(float64(ageSpreadMS), 0)/180.0)
	return base * freshPenalty * spreadPenalty
}

func fmtFloat(v float64) string {
	return strconv.FormatFloat(v, 'f', 12, 64)
}

func isFinite(v float64) bool { return !math.IsNaN(v) && !math.IsInf(v, 0) }
func eps() float64            { return 1e-12 }

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
