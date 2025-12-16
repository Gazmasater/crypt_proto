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



package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Exchange      string
	TrianglesFile string
	BookInterval  string

	FeePerLeg float64 // доля, 0.0004 = 0.04%
	MinProfit float64 // доля

	MinStart      float64 // USDT, 0 = выключено
	StartFraction float64 // 0..1

	// Trading
	TradeEnabled    bool
	TradeAmountUSDT float64
	TradeCooldownMs int

	// API keys (for selected exchange)
	APIKey    string
	APISecret string

	Debug bool
}

var debug bool

func SetDebug(v bool) { debug = v }

func Dlog(format string, args ...any) {
	if debug {
		log.Printf(format, args...)
	}
}

func loadEnvFloat(name string, def float64) float64 {
	raw := strings.TrimSpace(os.Getenv(name))
	if raw == "" {
		return def
	}
	v, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		log.Printf("bad %s=%q: %v, using default %f", name, raw, err, def)
		return def
	}
	return v
}

func loadEnvInt(name string, def int) int {
	raw := strings.TrimSpace(os.Getenv(name))
	if raw == "" {
		return def
	}
	v, err := strconv.Atoi(raw)
	if err != nil {
		log.Printf("bad %s=%q: %v, using default %d", name, raw, err, def)
		return def
	}
	return v
}

func loadEnvBool(name string, def bool) bool {
	raw := strings.TrimSpace(os.Getenv(name))
	if raw == "" {
		return def
	}
	switch strings.ToLower(raw) {
	case "1", "true", "yes", "y", "on":
		return true
	case "0", "false", "no", "n", "off":
		return false
	default:
		log.Printf("bad %s=%q: using default %v", name, raw, def)
		return def
	}
}

func clamp01(v, def float64) float64 {
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

	// проценты -> доли
	feePct := loadEnvFloat("FEE_PCT", 0.04)
	minPct := loadEnvFloat("MIN_PROFIT_PCT", 0.1)

	// MIN_START_USDT (предпочтительно) или MIN_START
	minStart := loadEnvFloat("MIN_START_USDT", -1)
	if minStart < 0 {
		minStart = loadEnvFloat("MIN_START", 0)
	}

	startFraction := clamp01(loadEnvFloat("START_FRACTION", 0.5), 0.5)

	// trading flags
	tradeEnabled := loadEnvBool("TRADE_ENABLED", false)
	tradeAmount := loadEnvFloat("TRADE_AMOUNT_USDT", 2.0)
	tradeCooldown := loadEnvInt("TRADE_COOLDOWN_MS", 800)

	// debug
	dbg := loadEnvBool("DEBUG", false)
	SetDebug(dbg)

	// keys: exchange-specific first, then fallback
	apiKey, apiSecret := "", ""
	switch ex {
	case "MEXC":
		apiKey = strings.TrimSpace(os.Getenv("MEXC_API_KEY"))
		apiSecret = strings.TrimSpace(os.Getenv("MEXC_API_SECRET"))
	case "KUCOIN":
		// если позже добавишь KuCoin trader — будет удобно
		apiKey = strings.TrimSpace(os.Getenv("KUCOIN_API_KEY"))
		apiSecret = strings.TrimSpace(os.Getenv("KUCOIN_API_SECRET"))
	case "OKX":
		apiKey = strings.TrimSpace(os.Getenv("OKX_API_KEY"))
		apiSecret = strings.TrimSpace(os.Getenv("OKX_API_SECRET"))
	}

	// fallback для совместимости
	if apiKey == "" {
		apiKey = strings.TrimSpace(os.Getenv("API_KEY"))
	}
	if apiSecret == "" {
		apiSecret = strings.TrimSpace(os.Getenv("API_SECRET"))
	}

	cfg := Config{
		Exchange:        ex,
		TrianglesFile:   tf,
		BookInterval:    bi,
		FeePerLeg:       feePct / 100.0,
		MinProfit:       minPct / 100.0,
		MinStart:        minStart,
		StartFraction:   startFraction,
		TradeEnabled:    tradeEnabled,
		TradeAmountUSDT: tradeAmount,
		TradeCooldownMs: tradeCooldown,
		APIKey:          apiKey,
		APISecret:       apiSecret,
		Debug:           dbg,
	}

	log.Printf("Exchange: %s", cfg.Exchange)
	log.Printf("Triangles file: %s", cfg.TrianglesFile)
	log.Printf("Book interval: %s", cfg.BookInterval)
	log.Printf("Fee per leg: %.4f %% (rate=%.6f)", feePct, cfg.FeePerLeg)
	log.Printf("Min profit per cycle: %.4f %% (rate=%.6f)", minPct, cfg.MinProfit)
	log.Printf("Min start amount (USDT): %.4f", cfg.MinStart)
	log.Printf("Start fraction: %.4f", cfg.StartFraction)
	log.Printf("Trade enabled: %v", cfg.TradeEnabled)
	log.Printf("Trade amount (USDT): %.4f", cfg.TradeAmountUSDT)
	log.Printf("Trade cooldown (ms): %d", cfg.TradeCooldownMs)

	if cfg.TradeEnabled {
		if cfg.APIKey != "" && cfg.APISecret != "" {
			log.Printf("API key/secret: loaded for %s", cfg.Exchange)
		} else {
			log.Printf("API key/secret: MISSING for %s (will fall back to DRY-RUN if main.go checks keys)", cfg.Exchange)
		}
	}

	return cfg
}



executor_dryrun.go

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


executor_real.go

package arb

import (
	"context"
	"fmt"
	"io"
	"math"
	"strings"
	"sync"
	"time"

	"crypt_proto/domain"
)

type RealExecutor struct {
	trader  SpotTrader
	out     io.Writer
	filters map[string]SymbolFilter

	mu     sync.Mutex
	banned map[string]string
}

func NewRealExecutor(tr SpotTrader, out io.Writer, filters map[string]SymbolFilter) *RealExecutor {
	return &RealExecutor{
		trader:  tr,
		out:     out,
		filters: filters,
		banned:  make(map[string]string),
	}
}

func (e *RealExecutor) Name() string { return "REAL" }

func (e *RealExecutor) isBanned(symbol string) (bool, string) {
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

// Execute: старт всегда USDT; BUY по quoteOrderQty; SELL по балансу (анти-Oversold)
func (e *RealExecutor) Execute(ctx context.Context, t domain.Triangle, quotes map[string]domain.Quote, startUSDT float64) error {
	// бан ног
	for _, leg := range t.Legs {
		if ok, reason := e.isBanned(leg.Symbol); ok {
			fmt.Fprintf(e.out, "  [REAL EXEC] skip banned symbol=%s reason=%s\n", leg.Symbol, reason)
			return nil
		}
	}

	curAsset := "USDT"

	fmt.Fprintf(e.out, "  [REAL EXEC] start=%.6f USDT triangle=%s\n", startUSDT, t.Name)

	// --- leg1 ---
	if err := e.execLeg(ctx, 1, t.Legs[0], curAsset, startUSDT); err != nil {
		e.handleExecErr(t.Legs[0].Symbol, err, 1)
		return nil
	}
	curAsset = t.Legs[0].To

	// дождаться баланса после 1 ноги
	if _, err := e.waitBalance(ctx, curAsset, 6, 60*time.Millisecond); err != nil {
		fmt.Fprintf(e.out, "    [REAL EXEC] leg 1 balance ERROR: %v\n", err)
		return nil
	}

	// --- leg2 ---
	if err := e.execLeg(ctx, 2, t.Legs[1], curAsset, 0); err != nil {
		e.handleExecErr(t.Legs[1].Symbol, err, 2)
		return nil
	}
	curAsset = t.Legs[1].To

	if _, err := e.waitBalance(ctx, curAsset, 6, 60*time.Millisecond); err != nil {
		fmt.Fprintf(e.out, "    [REAL EXEC] leg 2 balance ERROR: %v\n", err)
		return nil
	}

	// --- leg3 ---
	if err := e.execLeg(ctx, 3, t.Legs[2], curAsset, 0); err != nil {
		e.handleExecErr(t.Legs[2].Symbol, err, 3)
		return nil
	}

	fmt.Fprintf(e.out, "  [REAL EXEC] done triangle %s\n", t.Name)
	return nil
}

// execLeg:
// - если leg.From == curAsset и leg.Dir >0 => SELL symbol qty=balance(curAsset)
// - если leg.From == curAsset и leg.Dir <0 => BUY symbol quoteOrderQty = fixedUSDT (только для 1-й ноги)
// - в остальных случаях считаем ошибкой (плохой треугольник)
func (e *RealExecutor) execLeg(ctx context.Context, n int, leg domain.Leg, curAsset string, fixedUSDT float64) error {
	if leg.From != curAsset {
		return fmt.Errorf("leg %d: expected FROM=%s got=%s", n, curAsset, leg.From)
	}

	// BUY/SELL определяется направлением
	if leg.Dir < 0 {
		// BUY: только по quoteOrderQty, иначе упрёмся в scale/step
		if fixedUSDT <= 0 {
			return fmt.Errorf("leg %d: BUY requires fixedUSDT>0", n)
		}
		fmt.Fprintf(e.out, "    [REAL EXEC] leg %d: BUY %s quoteOrderQty=%.6f\n", n, leg.Symbol, fixedUSDT)
		_, err := e.trader.PlaceMarketOrder(ctx, leg.Symbol, "BUY", 0, fixedUSDT)
		return err
	}

	// SELL по балансу
	free, err := e.trader.GetBalance(ctx, leg.From)
	if err != nil {
		return fmt.Errorf("leg %d: balance %s: %w", n, leg.From, err)
	}

	// анти Oversold запас
	free *= 0.995

	qty := e.normalizeQtyDown(leg.Symbol, free)
	if qty <= 0 {
		return fmt.Errorf("leg %d: qty<=0 after normalize (raw=%.12f) (SELL %s)", n, free, leg.Symbol)
	}

	fmt.Fprintf(e.out, "    [REAL EXEC] leg %d: SELL %s qty=%.12f\n", n, leg.Symbol, qty)
	_, err = e.trader.PlaceMarketOrder(ctx, leg.Symbol, "SELL", qty, 0)
	return err
}

func (e *RealExecutor) normalizeQtyDown(symbol string, qty float64) float64 {
	if qty <= 0 {
		return 0
	}
	f, ok := e.filters[symbol]
	if !ok || f.StepSize <= 0 {
		// грубо вниз на 8 знаков
		return floorToDecimals(qty, 8)
	}

	qty = math.Floor(qty/f.StepSize) * f.StepSize
	if f.MinQty > 0 && qty < f.MinQty {
		return 0
	}

	return trimFloat(qty, 12)
}

func floorToDecimals(v float64, dec int) float64 {
	if dec <= 0 {
		return math.Floor(v)
	}
	p := math.Pow(10, float64(dec))
	return math.Floor(v*p) / p
}

func trimFloat(v float64, dec int) float64 {
	return floorToDecimals(v, dec)
}

func (e *RealExecutor) waitBalance(ctx context.Context, asset string, tries int, sleep time.Duration) (float64, error) {
	var last float64
	var err error
	for i := 0; i < tries; i++ {
		last, err = e.trader.GetBalance(ctx, asset)
		if err == nil && last > 0 {
			return last, nil
		}
		time.Sleep(sleep)
	}
	return last, err
}

func (e *RealExecutor) handleExecErr(symbol string, err error, leg int) {
	// 10007 - баним
	if err == nil {
		return
	}
	msg := err.Error()
	if strings.Contains(msg, `"code":10007`) || strings.Contains(msg, "symbol not support api") {
		e.ban(symbol, "symbol not support api (10007)")
	}
	fmt.Fprintf(e.out, "    [REAL EXEC] leg %d ERROR: %v\n", leg, err)
}


package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"sync"
	"syscall"
	"time"

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

	// executor
	if cfg.Exchange == "MEXC" && cfg.TradeEnabled && cfg.APIKey != "" && cfg.APISecret != "" {
		caps, err := mexc.FetchSymbolCapsMEXC(ctx)
		if err != nil {
			log.Printf("WARN: cannot fetch exchangeInfo filters: %v", err)
		}
		filters := make(map[string]arb.SymbolFilter)
		for sym, c := range caps {
			filters[sym] = arb.SymbolFilter{StepSize: c.StepSize, MinQty: c.MinQty}
		}

		tr := mexc.NewTrader(cfg.APIKey, cfg.APISecret, cfg.Debug)
		consumer.Executor = arb.NewRealExecutor(tr, arbOut, filters)
		log.Printf("Executor: REAL")
	} else {
		consumer.Executor = arb.NewDryRunExecutor(arbOut)
		log.Printf("Executor: DRY-RUN")
	}

	consumer.Start(ctx, events, triangles, indexBySymbol, &wg)
	feed.Start(ctx, &wg, symbols, cfg.BookInterval, events)

	<-ctx.Done()
	log.Println("shutting down...")

	// НЕ закрываем events (иначе WS горутины могут паниковать)
	wg.Wait()
	log.Println("bye")
}



# ======================
# EXCHANGE
# ======================
EXCHANGE=MEXC
TRIANGLES_FILE=triangles_markets.csv
BOOK_INTERVAL=10ms

# ======================
# ARBITRAGE LOGIC
# ======================
# комиссия одной ноги (%)
FEE_PCT=0.04

# минимальная прибыль треугольника (%)
MIN_PROFIT_PCT=0.1

# минимальный старт (в USDT)
MIN_START_USDT=2

# доля от maxStart (1 = брать максимум)
START_FRACTION=1.0

# ======================
# TRADING
# ======================
# ВКЛ / ВЫКЛ реальную торговлю
TRADE_ENABLED=true

# ВСЕГДА торговать на эту сумму USDT
TRADE_AMOUNT_USDT=2

# анти-флуд между сделками (мс)
TRADE_COOLDOWN_MS=800

# ======================
# MEXC API
# ======================
MEXC_API_KEY=PASTE_YOUR_KEY_HERE
MEXC_API_SECRET=PASTE_YOUR_SECRET_HERE

# ======================
# DEBUG
# ======================
DEBUG=false









