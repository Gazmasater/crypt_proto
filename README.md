mx0vglmT3srN1IS19H
135bb7a7509e4421bad692415c53753b



sudo systemctl mask sleep.target suspend.target hibernate.target hybrid-sleep.target



wbs-api.mexc.com/ws 


[https://edis-global.vercel.app/ru/vps-hosting/singapore-singapore
](https://sg.edisglobal.com/)



git pull --rebase origin privat
git push origin privat


BOOK_INTERVAL=100ms
SYMBOLS_FILE=triangles_markets.csv
DEBUG=false


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







		return false
	default:
		return def
	}
}

func clamp01(v float64, def float64) float64 {
	if v <= 0 || v > 1 {
		return def
	}
	return v
}

func Load() Config {
	_ = godotenv.Load(".env")

	ex := strings.ToUpper(strings.TrimSpace(os.Getenv("EXCHANGE")))
	if ex == "" {
		ex = "MEXC"
	}

	tf := strings.TrimSpace(os.Getenv("TRIANGLES_FILE"))
	if tf == "" {
		tf = "triangles_markets.csv"
	}

	bi := strings.TrimSpace(os.Getenv("BOOK_INTERVAL"))
	if bi == "" {
		bi = "100ms"
	}

	feePct := loadEnvFloat("FEE_PCT", 0.04)
	minPct := loadEnvFloat("MIN_PROFIT_PCT", 0.5)

	minStart := loadEnvFloat("MIN_START_USDT", -1)
	if minStart < 0 {
		minStart = loadEnvFloat("MIN_START", 0)
	}

	startFraction := clamp01(loadEnvFloat("START_FRACTION", 0.5), 0.5)
	debug = loadEnvBool("DEBUG", false)

	// keys
	apiKey := strings.TrimSpace(os.Getenv("MEXC_API_KEY"))
	if apiKey == "" {
		apiKey = strings.TrimSpace(os.Getenv("API_KEY"))
	}
	apiSecret := strings.TrimSpace(os.Getenv("MEXC_API_SECRET"))
	if apiSecret == "" {
		apiSecret = strings.TrimSpace(os.Getenv("API_SECRET"))
	}

	tradeEnabled := loadEnvBool("TRADE_ENABLED", false)
	tradeAmt := loadEnvFloat("TRADE_AMOUNT_USDT", 2.0)
	if tradeAmt <= 0 {
		tradeAmt = 2.0
	}
	tradeCooldown := loadEnvInt("TRADE_COOLDOWN_MS", 800)
	if tradeCooldown < 0 {
		tradeCooldown = 0
	}

	cfg := Config{
		Exchange:        ex,
		TrianglesFile:   tf,
		BookInterval:    bi,
		FeePerLeg:       feePct / 100.0,
		MinProfit:       minPct / 100.0,
		MinStart:        minStart,
		StartFraction:   startFraction,
		Debug:           debug,
		APIKey:          apiKey,
		APISecret:       apiSecret,
		TradeEnabled:    tradeEnabled,
		TradeAmountUSDT: tradeAmt,
		TradeCooldownMs: tradeCooldown,
	}

	log.Printf("Exchange: %s", cfg.Exchange)
	log.Printf("Triangles file: %s", cfg.TrianglesFile)
	log.Printf("Book interval: %s", cfg.BookInterval)
	log.Printf("Fee per leg: %.4f %% (rate=%.6f)", feePct, cfg.FeePerLeg)
	log.Printf("Min profit per cycle: %.4f %% (rate=%.6f)", minPct, cfg.MinProfit)
	log.Printf("Min start amount (USDT): %.4f", cfg.MinStart)
	log.Printf("Start fraction: %.4f", cfg.StartFraction)
	log.Printf("Trade enabled: %v", cfg.TradeEnabled)
	log.Printf("Trade amount USDT: %.4f", cfg.TradeAmountUSDT)
	log.Printf("Trade cooldown ms: %d", cfg.TradeCooldownMs)

	if cfg.APIKey != "" && cfg.APISecret != "" {
		log.Printf("API key/secret: loaded for %s", cfg.Exchange)
	} else {
		log.Printf("API key/secret: NOT set")
	}

	return cfg
}

func Dlog(format string, args ...any) {
	if debug {
		log.Printf(format, args...)
	}
}
2) crypt_proto/arb/executor.go (НОВЫЙ ФАЙЛ)
go
Копировать код
package arb

import (
	"context"

	"crypt_proto/domain"
)

type SpotTrader interface {
	PlaceMarketOrder(ctx context.Context, symbol, side string, quantity, quoteOrderQty float64) (orderID string, err error)
	GetBalance(ctx context.Context, asset string) (free float64, err error)
}

type SymbolFilter struct {
	StepSize float64
	MinQty   float64
}

type Executor interface {
	Name() string
	Execute(ctx context.Context, t domain.Triangle, quotes map[string]domain.Quote, startUSDT float64) error
}
3) crypt_proto/arb/executor_dryrun.go (НОВЫЙ ФАЙЛ)
go
Копировать код
package arb

import (
	"context"
	"fmt"
	"io"

	"crypt_proto/domain"
)

type DryRunExecutor struct {
	out io.Writer
}

func NewDryRunExecutor(out io.Writer) *DryRunExecutor {
	return &DryRunExecutor{out: out}
}

func (e *DryRunExecutor) Name() string { return "DRY-RUN" }

func (e *DryRunExecutor) Execute(ctx context.Context, t domain.Triangle, quotes map[string]domain.Quote, startUSDT float64) error {
	fmt.Fprintf(e.out, "  [DRY RUN] start=%.6f USDT triangle=%s\n", startUSDT, t.Name)
	return nil
}
4) crypt_proto/arb/executor_real.go (ПОЛНОСТЬЮ)
go
Копировать код
package arb

import (
	"context"
	"fmt"
	"io"
	"math"
	"sync"
	"time"

	"crypt_proto/domain"
)

type RealExecutor struct {
	trader  SpotTrader
	out     io.Writer
	filters map[string]SymbolFilter

	// ban list по символам, которые не поддерживают API (10007) и т.п.
	mu       sync.Mutex
	banned   map[string]string
	lastExec map[string]time.Time // key = triangle name
}

func NewRealExecutor(tr SpotTrader, out io.Writer, filters map[string]SymbolFilter) *RealExecutor {
	return &RealExecutor{
		trader:   tr,
		out:      out,
		filters:  filters,
		banned:   make(map[string]string),
		lastExec: make(map[string]time.Time),
	}
}

func (e *RealExecutor) Name() string { return "REAL" }

func (e *RealExecutor) IsBannedSymbol(symbol string) (bool, string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	reason, ok := e.banned[symbol]
	return ok, reason
}

func (e *RealExecutor) ban(symbol, reason string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.banned[symbol] = reason
}

func (e *RealExecutor) Execute(ctx context.Context, t domain.Triangle, quotes map[string]domain.Quote, startUSDT float64) error {
	// 1) Быстрая проверка на бан ног
	for _, leg := range t.Legs {
		if ok, reason := e.IsBannedSymbol(leg.Symbol); ok {
			fmt.Fprintf(e.out, "  [REAL EXEC] skip banned symbol=%s reason=%s\n", leg.Symbol, reason)
			return nil
		}
	}

	// 2) Симуляция (чтобы понять порядок from/to и расчётные qty)
	execs, ok := simulateTriangleExecution(t, quotes, "USDT", startUSDT, 0) // fee=0 тут не важно, мы будем по балансу
	if !ok || len(execs) != 3 {
		fmt.Fprintf(e.out, "  [REAL EXEC] simulate failed for %s\n", t.Name)
		return nil
	}

	fmt.Fprintf(e.out, "  [REAL EXEC] start=%.6f USDT triangle=%s\n", startUSDT, t.Name)

	// Нога 1: BUY (тратим ровно startUSDT) через quoteOrderQty
	leg1 := execs[0]
	if leg1.From != "USDT" {
		fmt.Fprintf(e.out, "  [REAL EXEC] leg1 expected FROM=USDT, got %s\n", leg1.From)
		return nil
	}

	if err := e.placeBuyByQuote(ctx, leg1.Symbol, startUSDT); err != nil {
		if isAPINotSupported(err) {
			e.ban(leg1.Symbol, "symbol not support api")
		}
		fmt.Fprintf(e.out, "    [REAL EXEC] leg 1 ERROR: %v\n", err)
		return nil
	}

	// Подтянуть баланс полученной монеты (base)
	assetAfter1 := leg1.To
	amtAfter1, err := e.waitBalance(ctx, assetAfter1, 6, 60*time.Millisecond)
	if err != nil {
		fmt.Fprintf(e.out, "    [REAL EXEC] leg 1 balance ERROR: %v\n", err)
		return nil
	}
	if amtAfter1 <= 0 {
		fmt.Fprintf(e.out, "    [REAL EXEC] leg 1 balance=%s is zero\n", assetAfter1)
		return nil
	}

	// Нога 2: обычно SELL base->quote (продаём то, что реально есть)
	leg2 := execs[1]
	if leg2.From != assetAfter1 {
		fmt.Fprintf(e.out, "    [REAL EXEC] leg2 expected FROM=%s, got %s\n", assetAfter1, leg2.From)
		return nil
	}

	if err := e.placeSellByBalance(ctx, leg2.Symbol, leg2.From); err != nil {
		if isAPINotSupported(err) {
			e.ban(leg2.Symbol, "symbol not support api")
		}
		fmt.Fprintf(e.out, "    [REAL EXEC] leg 2 ERROR: %v\n", err)
		return nil
	}

	assetAfter2 := leg2.To
	amtAfter2, err := e.waitBalance(ctx, assetAfter2, 6, 60*time.Millisecond)
	if err != nil {
		fmt.Fprintf(e.out, "    [REAL EXEC] leg 2 balance ERROR: %v\n", err)
		return nil
	}
	if amtAfter2 <= 0 {
		fmt.Fprintf(e.out, "    [REAL EXEC] leg 2 balance=%s is zero\n", assetAfter2)
		return nil
	}

	// Нога 3: SELL полученную валюту в USDT по фактическому балансу (важно! иначе Oversold)
	leg3 := execs[2]
	if leg3.To != "USDT" {
		fmt.Fprintf(e.out, "    [REAL EXEC] leg3 expected TO=USDT, got %s\n", leg3.To)
		return nil
	}
	if leg3.From != assetAfter2 {
		fmt.Fprintf(e.out, "    [REAL EXEC] leg3 expected FROM=%s, got %s\n", assetAfter2, leg3.From)
		return nil
	}

	if err := e.placeSellByBalance(ctx, leg3.Symbol, leg3.From); err != nil {
		if isAPINotSupported(err) {
			e.ban(leg3.Symbol, "symbol not support api")
		}
		fmt.Fprintf(e.out, "    [REAL EXEC] leg 3 ERROR: %v\n", err)
		return nil
	}

	fmt.Fprintf(e.out, "  [REAL EXEC] done triangle %s\n", t.Name)
	return nil
}

func (e *RealExecutor) placeBuyByQuote(ctx context.Context, symbol string, quoteUSDT float64) error {
	// quoteOrderQty: тратим ровно quoteUSDT USDT
	_, err := e.trader.PlaceMarketOrder(ctx, symbol, "BUY", 0, quoteUSDT)
	if err != nil {
		return err
	}
	fmt.Fprintf(e.out, "    [REAL EXEC] leg 1: BUY %s quoteOrderQty=%.6f\n", symbol, quoteUSDT)
	return nil
}

func (e *RealExecutor) placeSellByBalance(ctx context.Context, symbol string, asset string) error {
	free, err := e.trader.GetBalance(ctx, asset)
	if err != nil {
		return fmt.Errorf("balance %s: %w", asset, err)
	}

	// зазор от микродрифтов/комиссий, чтобы не ловить Oversold
	free = free * 0.995

	qty := e.normalizeQtyDown(symbol, free)
	if qty <= 0 {
		return fmt.Errorf("qty<=0 after normalize (symbol=%s asset=%s free=%.10f)", symbol, asset, free)
	}

	_, err = e.trader.PlaceMarketOrder(ctx, symbol, "SELL", qty, 0)
	if err != nil {
		return err
	}

	fmt.Fprintf(e.out, "    [REAL EXEC] SELL %s qty=%.10f (by balance)\n", symbol, qty)
	return nil
}

func (e *RealExecutor) normalizeQtyDown(symbol string, qty float64) float64 {
	if qty <= 0 {
		return 0
	}
	f, ok := e.filters[symbol]
	if !ok {
		// если фильтров нет — оставим 8 знаков вниз
		return floorToDecimals(qty, 8)
	}

	if f.StepSize > 0 {
		qty = math.Floor(qty/f.StepSize) * f.StepSize
	}

	if f.MinQty > 0 && qty < f.MinQty {
		return 0
	}

	// подчистим "хвосты"
	return trimFloat(qty, 12)
}

func floorToDecimals(v float64, dec int) float64 {
	if dec <= 0 {
		return math.Floor(v)
	}
	p := math.Pow(10, float64(dec))
	return math.Floor(v*p) / p
}

func trimFloat(v float64, maxDec int) float64 {
	return floorToDecimals(v, maxDec)
}

func (e *RealExecutor) waitBalance(ctx context.Context, asset string, tries int, sleep time.Duration) (float64, error) {
	var last float64
	var err error
	for i := 0; i < tries; i++ {
		last, err = e.trader.GetBalance(ctx, asset)
		if err != nil {
			time.Sleep(sleep)
			continue
		}
		if last > 0 {
			return last, nil
		}
		time.Sleep(sleep)
	}
	return last, err
}

// ---- ошибки ----

func isAPINotSupported(err error) bool {
	if err == nil {
		return false
	}
	s := err.Error()
	return stringsContainsAny(s,
		`"code":10007`,
		"symbol not support api",
	)
}

func stringsContainsAny(s string, subs ...string) bool {
	for _, sub := range subs {
		if sub != "" && contains(s, sub) {
			return true
		}
	}
	return false
}

func contains(s, sub string) bool {
	return len(sub) > 0 && (len(s) >= len(sub)) && (indexOf(s, sub) >= 0)
}

// маленькие inline-реализации, чтобы не тащить strings/index
func indexOf(s, sub string) int {
	// простая реализация
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
Тут есть маленькая “самодельная” строковая логика, чтобы файл был самодостаточным. Если хочешь — заменим на strings.Contains.

5) crypt_proto/arb/arb.go (ПОЛНОСТЬЮ)
go
Копировать код
package arb

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"crypt_proto/domain"
)

const BaseAsset = "USDT"

type Consumer struct {
	FeePerLeg     float64
	MinProfit     float64
	MinStart      float64
	StartFraction float64

	TradeEnabled    bool
	TradeAmountUSDT float64
	TradeCooldown   time.Duration

	Executor Executor

	writer io.Writer

	// анти-спам/анти-ддос по трейду
	mu         sync.Mutex
	inFlight   bool
	lastTrade  map[int]time.Time
	lastPrint  map[int]time.Time
	minPrint   time.Duration
}

func NewConsumer(feePerLeg, minProfit, minStart float64, out io.Writer) *Consumer {
	return &Consumer{
		FeePerLeg:       feePerLeg,
		MinProfit:       minProfit,
		MinStart:        minStart,
		StartFraction:   0.5,
		TradeEnabled:    false,
		TradeAmountUSDT: 2.0,
		TradeCooldown:   800 * time.Millisecond,
		writer:          out,
		lastTrade:       make(map[int]time.Time),
		lastPrint:       make(map[int]time.Time),
		minPrint:        5 * time.Millisecond,
	}
}

func (c *Consumer) Start(
	ctx context.Context,
	events <-chan domain.Event,
	triangles []domain.Triangle,
	indexBySymbol map[string][]int,
	wg *sync.WaitGroup,
) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		c.run(ctx, events, triangles, indexBySymbol)
	}()
}

func (c *Consumer) run(
	ctx context.Context,
	events <-chan domain.Event,
	triangles []domain.Triangle,
	indexBySymbol map[string][]int,
) {
	quotes := make(map[string]domain.Quote)

	sf := c.StartFraction
	if sf <= 0 || sf > 1 {
		sf = 0.5
	}

	for {
		select {
		case <-ctx.Done():
			return

		case ev, ok := <-events:
			if !ok {
				return
			}

			prev, okPrev := quotes[ev.Symbol]
			if okPrev &&
				prev.Bid == ev.Bid && prev.Ask == ev.Ask &&
				prev.BidQty == ev.BidQty && prev.AskQty == ev.AskQty {
				continue
			}

			quotes[ev.Symbol] = domain.Quote{Bid: ev.Bid, Ask: ev.Ask, BidQty: ev.BidQty, AskQty: ev.AskQty}

			trIDs := indexBySymbol[ev.Symbol]
			if len(trIDs) == 0 {
				continue
			}

			now := time.Now()

			for _, id := range trIDs {
				tr := triangles[id]

				prof, ok := domain.EvalTriangle(tr, quotes, c.FeePerLeg)
				if !ok || prof < c.MinProfit {
					continue
				}

				ms, okMS := domain.ComputeMaxStartTopOfBook(tr, quotes, c.FeePerLeg)
				if !okMS {
					continue
				}

				safeStart := ms.MaxStart * sf

				// фильтр по MIN_START_USDT (по safeStart)
				if c.MinStart > 0 {
					safeUSDT, okConv := convertToUSDT(safeStart, ms.StartAsset, quotes)
					if !okConv || safeUSDT < c.MinStart {
						continue
					}
				}

				// анти-спам печати
				if last, okLast := c.lastPrint[id]; okLast && now.Sub(last) < c.minPrint {
					// но торговлю можем всё равно запускать (ниже), печать — отдельно
				} else {
					c.lastPrint[id] = now
					msCopy := ms
					c.printTriangle(now, tr, prof, quotes, &msCopy, sf)
				}

				// торговля
				c.tryTrade(ctx, now, id, tr, quotes)
			}
		}
	}
}

func (c *Consumer) tryTrade(ctx context.Context, now time.Time, id int, tr domain.Triangle, quotes map[string]domain.Quote) {
	if !c.TradeEnabled || c.Executor == nil {
		return
	}

	// cooldown per triangle
	c.mu.Lock()
	if last, ok := c.lastTrade[id]; ok && now.Sub(last) < c.TradeCooldown {
		c.mu.Unlock()
		return
	}
	if c.inFlight {
		c.mu.Unlock()
		return
	}
	c.inFlight = true
	c.lastTrade[id] = now
	c.mu.Unlock()

	go func() {
		defer func() {
			c.mu.Lock()
			c.inFlight = false
			c.mu.Unlock()
		}()

		_ = c.Executor.Execute(ctx, tr, quotes, c.TradeAmountUSDT)
	}()
}

func (c *Consumer) printTriangle(
	ts time.Time,
	t domain.Triangle,
	profit float64,
	quotes map[string]domain.Quote,
	ms *domain.MaxStartInfo,
	startFraction float64,
) {
	w := c.writer
	fmt.Fprintf(w, "%s\n", ts.Format("2006-01-02 15:04:05.000"))

	if ms == nil {
		fmt.Fprintf(w, "[ARB] %+0.3f%%  %s\n", profit*100, t.Name)
		fmt.Fprintln(w)
		return
	}

	bneckSym := ""
	if ms.BottleneckLeg >= 0 && ms.BottleneckLeg < len(t.Legs) {
		bneckSym = t.Legs[ms.BottleneckLeg].Symbol
	}

	safeStart := ms.MaxStart * startFraction
	maxUSDT, okMax := convertToUSDT(ms.MaxStart, ms.StartAsset, quotes)
	safeUSDT, okSafe := convertToUSDT(safeStart, ms.StartAsset, quotes)

	maxUSDTStr, safeUSDTStr := "?", "?"
	if okMax {
		maxUSDTStr = fmt.Sprintf("%.4f", maxUSDT)
	}
	if okSafe {
		safeUSDTStr = fmt.Sprintf("%.4f", safeUSDT)
	}

	fmt.Fprintf(w,
		"[ARB] %+0.3f%%  %s  maxStart=%.4f %s (%s USDT)  safeStart=%.4f %s (%s USDT) (x%.2f)  bottleneck=%s\n",
		profit*100, t.Name,
		ms.MaxStart, ms.StartAsset, maxUSDTStr,
		safeStart, ms.StartAsset, safeUSDTStr,
		startFraction,
		bneckSym,
	)

	for _, leg := range t.Legs {
		q := quotes[leg.Symbol]
		mid := (q.Bid + q.Ask) / 2
		spreadAbs := q.Ask - q.Bid
		spreadPct := 0.0
		if mid > 0 {
			spreadPct = spreadAbs / mid * 100
		}
		side := ""
		if leg.Dir > 0 {
			side = fmt.Sprintf("%s/%s", leg.From, leg.To)
		} else {
			side = fmt.Sprintf("%s/%s", leg.To, leg.From)
		}
		fmt.Fprintf(w, "  %s (%s): bid=%.10f ask=%.10f  spread=%.10f (%.5f%%)  bidQty=%.4f askQty=%.4f\n",
			leg.Symbol, side,
			q.Bid, q.Ask,
			spreadAbs, spreadPct,
			q.BidQty, q.AskQty,
		)
	}
	fmt.Fprintln(w)
}

// ===== utils =====

func OpenLogWriter(path string) (io.WriteCloser, *bufio.Writer, io.Writer) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		log.Fatalf("open %s: %v", path, err)
	}
	buf := bufio.NewWriter(f)
	out := io.MultiWriter(os.Stdout, buf)
	return f, buf, out
}

func convertToUSDT(amount float64, asset string, quotes map[string]domain.Quote) (float64, bool) {
	if amount <= 0 {
		return 0, false
	}
	if asset == BaseAsset {
		return amount, true
	}
	if q, ok := quotes[asset+"USDT"]; ok && q.Bid > 0 {
		return amount * q.Bid, true
	}
	if q, ok := quotes["USDT"+asset]; ok && q.Ask > 0 {
		return amount / q.Ask, true
	}
	if amtUSDC, ok1 := convertViaQuote(amount, asset, "USDC", quotes); ok1 {
		if amtUSDT, ok2 := convertViaQuote(amtUSDC, "USDC", "USDT", quotes); ok2 {
			return amtUSDT, true
		}
	}
	return 0, false
}

func convertViaQuote(amount float64, from, to string, quotes map[string]domain.Quote) (float64, bool) {
	if amount <= 0 {
		return 0, false
	}
	if from == to {
		return amount, true
	}
	if q, ok := quotes[from+to]; ok && q.Bid > 0 {
		return amount * q.Bid, true
	}
	if q, ok := quotes[to+from]; ok && q.Ask > 0 {
		return amount / q.Ask, true
	}
	return 0, false
}
6) crypt_proto/mexc/exchangeinfo.go (ПОЛНОСТЬЮ) — фильтры stepSize/minQty
go
Копировать код
package mexc

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type SymbolCaps struct {
	Symbol      string
	Status      string
	HasMarket   bool
	StepSize    float64
	MinQty      float64
	MinNotional float64
}

func FetchSymbolCapsMEXC(ctx context.Context) (map[string]SymbolCaps, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.mexc.com/api/v3/exchangeInfo", nil)
	if err != nil {
		return nil, err
	}

	cl := &http.Client{Timeout: 12 * time.Second}
	resp, err := cl.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("exchangeInfo status=%d", resp.StatusCode)
	}

	var root map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&root); err != nil {
		return nil, err
	}

	rawSyms, _ := root["symbols"].([]any)
	out := make(map[string]SymbolCaps, len(rawSyms))

	for _, item := range rawSyms {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}

		symbol, _ := m["symbol"].(string)
		if symbol == "" {
			continue
		}

		status, _ := m["status"].(string)

		// safe default: если orderTypes нет — не режем
		hasMarket := true
		if otsAny, ok := m["orderTypes"]; ok {
			hasMarket = false
			if ots, ok := otsAny.([]any); ok {
				for _, v := range ots {
					if s, ok := v.(string); ok && strings.EqualFold(s, "MARKET") {
						hasMarket = true
						break
					}
				}
			} else {
				hasMarket = true
			}
		}

		stepSize, minQty, minNotional := 0.0, 0.0, 0.0
		if flt, ok := m["filters"].([]any); ok {
			for _, f := range flt {
				fm, ok := f.(map[string]any)
				if !ok {
					continue
				}
				ft, _ := fm["filterType"].(string)
				switch ft {
				case "LOT_SIZE":
					stepSize = readFloatAny(fm["stepSize"])
					minQty = readFloatAny(fm["minQty"])
				case "MIN_NOTIONAL":
					minNotional = readFloatAny(fm["minNotional"])
				}
			}
		}

		out[symbol] = SymbolCaps{
			Symbol:      symbol,
			Status:      status,
			HasMarket:   hasMarket,
			StepSize:    stepSize,
			MinQty:      minQty,
			MinNotional: minNotional,
		}
	}

	return out, nil
}

func readFloatAny(v any) float64 {
	switch t := v.(type) {
	case string:
		f, _ := strconv.ParseFloat(t, 64)
		return f
	case float64:
		return t
	default:
		return 0
	}
}
7) crypt_proto/mexc/trader.go (ПОЛНОСТЬЮ) — подпись, ордера, баланс
go
Копировать код
package mexc

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Trader struct {
	apiKey    string
	apiSecret string
	debug     bool
	baseURL   string
	client    *http.Client
}

func NewTrader(apiKey, apiSecret string, debug bool) *Trader {
	return &Trader{
		apiKey:    apiKey,
		apiSecret: apiSecret,
		debug:     debug,
		baseURL:   "https://api.mexc.com",
		client:    &http.Client{Timeout: 10 * time.Second},
	}
}

// PlaceMarketOrder:
// - BUY: лучше передавать quoteOrderQty (потратить ровно USDT)
// - SELL: передавать quantity (сколько base продаем)
func (t *Trader) PlaceMarketOrder(ctx context.Context, symbol, side string, quantity, quoteOrderQty float64) (string, error) {
	side = strings.ToUpper(strings.TrimSpace(side))
	if side != "BUY" && side != "SELL" {
		return "", fmt.Errorf("bad side=%s", side)
	}

	params := url.Values{}
	params.Set("symbol", symbol)
	params.Set("side", side)
	params.Set("type", "MARKET")
	params.Set("timestamp", strconv.FormatInt(time.Now().UnixMilli(), 10))

	if side == "BUY" {
		if quoteOrderQty <= 0 {
			return "", fmt.Errorf("BUY requires quoteOrderQty>0")
		}
		// MEXC (часто) принимает quoteOrderQty как Binance-style
		params.Set("quoteOrderQty", fmt.Sprintf("%.8f", quoteOrderQty))
	} else {
		if quantity <= 0 {
			return "", fmt.Errorf("SELL requires quantity>0")
		}
		params.Set("quantity", fmt.Sprintf("%.12f", quantity))
	}

	sig := t.sign(params.Encode())
	params.Set("signature", sig)

	endpoint := t.baseURL + "/api/v3/order"
	reqURL := endpoint + "?" + params.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("X-MEXC-APIKEY", t.apiKey)
	// чтобы не ловить "Invalid content Type." (700013)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := t.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	b, _ := io.ReadAll(resp.Body)

	if resp.StatusCode/100 != 2 {
		return "", fmt.Errorf("mexc order error: status=%d body=%s", resp.StatusCode, string(b))
	}

	// ответ может быть разный — попробуем вытащить orderId
	var m map[string]any
	_ = json.Unmarshal(b, &m)
	if v, ok := m["orderId"]; ok {
		return fmt.Sprintf("%v", v), nil
	}
	if v, ok := m["orderIdStr"]; ok {
		return fmt.Sprintf("%v", v), nil
	}
	return "", nil
}

func (t *Trader) GetBalance(ctx context.Context, asset string) (float64, error) {
	params := url.Values{}
	params.Set("timestamp", strconv.FormatInt(time.Now().UnixMilli(), 10))
	sig := t.sign(params.Encode())
	params.Set("signature", sig)

	endpoint := t.baseURL + "/api/v3/account"
	reqURL := endpoint + "?" + params.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return 0, err
	}

	req.Header.Set("X-MEXC-APIKEY", t.apiKey)

	resp, err := t.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	b, _ := io.ReadAll(resp.Body)
	if resp.StatusCode/100 != 2 {
		return 0, fmt.Errorf("mexc account error: status=%d body=%s", resp.StatusCode, string(b))
	}

	// Ответ похож на Binance: balances: [{asset:"USDT", free:"...", locked:"..."}]
	var root map[string]any
	if err := json.Unmarshal(b, &root); err != nil {
		return 0, err
	}

	balAny, _ := root["balances"].([]any)
	asset = strings.ToUpper(strings.TrimSpace(asset))

	for _, it := range balAny {
		m, ok := it.(map[string]any)
		if !ok {
			continue
		}
		a, _ := m["asset"].(string)
		if strings.ToUpper(strings.TrimSpace(a)) != asset {
			continue
		}

		freeStr := ""
		if s, ok := m["free"].(string); ok {
			freeStr = s
		} else if f, ok := m["free"].(float64); ok {
			return f, nil
		}
		if freeStr == "" {
			return 0, nil
		}
		v, _ := strconv.ParseFloat(freeStr, 64)
		return v, nil
	}

	return 0, nil
}

func (t *Trader) sign(query string) string {
	mac := hmac.New(sha256.New, []byte(t.apiSecret))
	_, _ = mac.Write([]byte(query))
	return hex.EncodeToString(mac.Sum(nil))
}
8) crypt_proto/cmd/cryptarb/main.go (ПОЛНОСТЬЮ)
go
Копировать код
package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"sync"
	"syscall"

	"crypt_proto/arb"
	"crypt_proto/config"
	"crypt_proto/domain"
	"crypt_proto/exchange"
	"crypt_proto/kucoin"
	"crypt_proto/mexc"

	_ "net/http/pprof"
)

func main() {
	go func() {
		log.Println("pprof on http://localhost:6060/debug/pprof/")
		_ = http.ListenAndServe("localhost:6060", nil)
	}()

	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	cfg := config.Load()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	triangles, symbols, indexBySymbol, err := domain.LoadTriangles(cfg.TrianglesFile)
	if err != nil {
		log.Fatalf("load triangles: %v", err)
	}
	if len(triangles) == 0 {
		log.Fatal("нет треугольников, нечего мониторить")
	}
	if len(symbols) == 0 {
		log.Fatal("нет символов для подписки")
	}
	log.Printf("треугольников: %d", len(triangles))
	log.Printf("символов для подписки: %d", len(symbols))

	var feed exchange.MarketDataFeed
	switch cfg.Exchange {
	case "MEXC":
		feed = mexc.NewFeed(cfg.Debug)
	case "KUCOIN":
		feed = kucoin.NewFeed(cfg.Debug)
	default:
		log.Fatalf("unknown EXCHANGE=%q (expected MEXC or KUCOIN)", cfg.Exchange)
	}

	log.Printf("Using exchange: %s", feed.Name())

	logFile, logBuf, arbOut := arb.OpenLogWriter("arbitrage.log")
	defer logFile.Close()
	defer logBuf.Flush()

	events := make(chan domain.Event, 8192)
	var wg sync.WaitGroup

	consumer := arb.NewConsumer(cfg.FeePerLeg, cfg.MinProfit, cfg.MinStart, arbOut)
	consumer.StartFraction = cfg.StartFraction
	consumer.TradeEnabled = cfg.TradeEnabled
	consumer.TradeAmountUSDT = cfg.TradeAmountUSDT
	consumer.TradeCooldown = time.Duration(cfg.TradeCooldownMs) * time.Millisecond

	// Executor
	if cfg.Exchange == "MEXC" && cfg.TradeEnabled && cfg.APIKey != "" && cfg.APISecret != "" {
		// грузим фильтры step/minQty
		caps, err := mexc.FetchSymbolCapsMEXC(ctx)
		if err != nil {
			log.Printf("WARN: cannot fetch exchangeInfo for filters: %v (будем округлять грубо)", err)
		}
		filters := make(map[string]arb.SymbolFilter)
		for sym, c := range caps {
			filters[sym] = arb.SymbolFilter{StepSize: c.StepSize, MinQty: c.MinQty}
		}

		tr := mexc.NewTrader(cfg.APIKey, cfg.APISecret, cfg.Debug)
		consumer.Executor = arb.NewRealExecutor(tr, arbOut, filters)
		log.Printf("Executor: REAL (trade enabled)")
	} else {
		consumer.Executor = arb.NewDryRunExecutor(arbOut)
		log.Printf("Executor: DRY-RUN (trade disabled or no keys)")
	}

	consumer.Start(ctx, events, triangles, indexBySymbol, &wg)
	feed.Start(ctx, &wg, symbols, cfg.BookInterval, events)

	<-ctx.Done()
	log.Println("shutting down...")

	// ВАЖНО: НЕ закрываем events, иначе можно словить send-on-closed-channel от фида.
	wg.Wait()
	log.Println("bye")
}
Обрати внимание: тут нужен time импорт. Если IDE ругнётся — добавь import "time" сверху. (Я







