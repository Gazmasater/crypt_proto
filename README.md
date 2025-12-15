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


# --- Основное ---
EXCHANGE=MEXC
TRIANGLES_FILE=triangles_markets.csv
BOOK_INTERVAL=10ms

# --- Арбитраж ---
FEE_PCT=0.04
MIN_PROFIT_PCT=0.1
MIN_START_USDT=2
START_FRACTION=0.5
DEBUG=false

# --- Торговля ---
TRADE_ENABLED=false   # <<< ГЛАВНЫЙ ФЛАГ

# --- API ---
MEXC_API_KEY=
MEXC_API_SECRET=



package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Exchange      string // EXCHANGE=MEXC / KUCOIN / OKX ...
	TrianglesFile string // TRIANGLES_FILE=triangles_markets.csv
	BookInterval  string // BOOK_INTERVAL=10ms

	// Комиссия и минимальная прибыль (в долях, а не в процентах).
	// Пример: 0.0004 = 0.04%
	FeePerLeg float64 // FEE_PCT
	MinProfit float64 // MIN_PROFIT_PCT

	// Минимальный старт (обычно в USDT). 0 = фильтр выключен.
	MinStart float64 // MIN_START_USDT / MIN_START

	// Доля от maxStart, которую считаем безопасной (0..1). Например 0.5.
	StartFraction float64 // START_FRACTION

	// Логика
	Debug        bool // DEBUG
	TradeEnabled bool // TRADE_ENABLED

	// API ключи: для текущей биржи и/или глобальные API_KEY / API_SECRET
	APIKey    string
	APISecret string
}

// package-level флаг для Dlog
var debug bool

func SetDebug(v bool) { debug = v }

func Dlog(format string, args ...any) {
	if debug {
		log.Printf(format, args...)
	}
}

// loadEnvFloat читает float из ENV с дефолтом.
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

// clamp01 ограничивает значение [0,1], иначе возвращает def.
func clamp01(v, def float64) float64 {
	if v <= 0 || v > 1 {
		return def
	}
	return v
}

func Load() Config {
	// подхватываем .env, если есть
	_ = godotenv.Load(".env")

	// Биржа
	ex := strings.ToUpper(strings.TrimSpace(os.Getenv("EXCHANGE")))
	if ex == "" {
		ex = "MEXC"
	}

	// Файл треугольников
	tf := strings.TrimSpace(os.Getenv("TRIANGLES_FILE"))
	if tf == "" {
		tf = "triangles_markets.csv"
	}

	// Интервал обновления книги
	bi := strings.TrimSpace(os.Getenv("BOOK_INTERVAL"))
	if bi == "" {
		bi = "10ms"
	}

	// Проценты из ENV: FEE_PCT, MIN_PROFIT_PCT
	// Например: FEE_PCT=0.04 => 0.04% => 0.0004
	feePct := loadEnvFloat("FEE_PCT", 0.04)        // в процентах
	minPct := loadEnvFloat("MIN_PROFIT_PCT", 0.1)  // в процентах

	feePerLeg := feePct / 100.0
	minProfit := minPct / 100.0

	// MIN_START_USDT (предпочтительно) или MIN_START
	minStart := loadEnvFloat("MIN_START_USDT", -1)
	if minStart < 0 {
		minStart = loadEnvFloat("MIN_START", 0)
	}

	// START_FRACTION (0..1)
	startFraction := clamp01(loadEnvFloat("START_FRACTION", 0.5), 0.5)

	// DEBUG
	debugFlag := strings.ToLower(strings.TrimSpace(os.Getenv("DEBUG"))) == "true"
	debug = debugFlag // чтобы Dlog() сразу работал

	// TRADE_ENABLED
	tradeEnabled := strings.ToLower(strings.TrimSpace(os.Getenv("TRADE_ENABLED"))) == "true"

	// API-ключи: сперва EXCHANGE_API_KEY/SECRET, потом API_KEY/SECRET
	apiKey := strings.TrimSpace(os.Getenv(ex + "_API_KEY"))
	if apiKey == "" {
		apiKey = strings.TrimSpace(os.Getenv("API_KEY"))
	}

	apiSecret := strings.TrimSpace(os.Getenv(ex + "_API_SECRET"))
	if apiSecret == "" {
		apiSecret = strings.TrimSpace(os.Getenv("API_SECRET"))
	}

	cfg := Config{
		Exchange:      ex,
		TrianglesFile: tf,
		BookInterval:  bi,
		FeePerLeg:     feePerLeg,
		MinProfit:     minProfit,
		MinStart:      minStart,
		StartFraction: startFraction,
		Debug:         debugFlag,
		TradeEnabled:  tradeEnabled,
		APIKey:        apiKey,
		APISecret:     apiSecret,
	}

	log.Printf("Exchange: %s", cfg.Exchange)
	log.Printf("Triangles file: %s", cfg.TrianglesFile)
	log.Printf("Book interval: %s", cfg.BookInterval)
	log.Printf("Fee per leg: %.4f %% (rate=%.6f)", feePct, cfg.FeePerLeg)
	log.Printf("Min profit per cycle: %.4f %% (rate=%.6f)", minPct, cfg.MinProfit)
	log.Printf("Min start amount: %.4f", cfg.MinStart)
	log.Printf("Start fraction: %.4f", cfg.StartFraction)
	log.Printf("Debug: %v", cfg.Debug)
	log.Printf("Trade enabled: %v", cfg.TradeEnabled)

	if cfg.APIKey == "" || cfg.APISecret == "" {
		log.Printf("API key/secret: NOT SET (торговля по API невозможна)")
	} else {
		log.Printf("API key/secret: loaded for %s", cfg.Exchange)
	}

	return cfg
}










package arb

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"crypt_proto/domain"
)

const BaseAsset = "USDT"

// Исполнитель треугольника (DRY-RUN или реальный трейдер).
type TriangleExecutor interface {
	ExecuteTriangle(
		ctx context.Context,
		t domain.Triangle,
		quotes map[string]domain.Quote,
		ms *domain.MaxStartInfo,
		startFraction float64,
	)
}

type Consumer struct {
	FeePerLeg     float64
	MinProfit     float64
	MinStart      float64
	StartFraction float64

	// Если nil — только логируем, без попыток торговать.
	Executor TriangleExecutor

	writer io.Writer
}

func NewConsumer(feePerLeg, minProfit, minStart float64, out io.Writer) *Consumer {
	return &Consumer{
		FeePerLeg:     feePerLeg,
		MinProfit:     minProfit,
		MinStart:      minStart,
		StartFraction: 0.5,
		writer:        out,
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
	lastPrint := make(map[int]time.Time)
	const minPrintInterval = 5 * time.Millisecond

	sf := c.StartFraction
	if sf <= 0 || sf > 1 {
		sf = 0.5
	}

	for {
		select {
		case ev, ok := <-events:
			if !ok {
				return
			}

			prev, okPrev := quotes[ev.Symbol]
			if okPrev && prev.Bid == ev.Bid && prev.Ask == ev.Ask &&
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

				// 1) прибыль по треугольнику (с учётом комиссии FeePerLeg)
				prof, ok := domain.EvalTriangle(tr, quotes, c.FeePerLeg)
				if !ok || prof < c.MinProfit {
					continue
				}

				// 2) maxStart по top-of-book
				ms, okMS := domain.ComputeMaxStartTopOfBook(tr, quotes, c.FeePerLeg)
				if !okMS {
					continue
				}

				// safeStart = maxStart * StartFraction
				safeStart := ms.MaxStart * sf

				// 3) фильтр по MIN_START (в USDT по safeStart)
				if c.MinStart > 0 {
					safeUSDT, okConv := convertToUSDT(safeStart, ms.StartAsset, quotes)
					if !okConv || safeUSDT < c.MinStart {
						continue
					}
				}

				// 4) анти-спам: не чаще, чем minPrintInterval на один и тот же треугольник
				if last, okLast := lastPrint[id]; okLast && now.Sub(last) < minPrintInterval {
					continue
				}
				lastPrint[id] = now

				// копируем ms, чтобы не гонять указатель на одну и ту же структуру
				msCopy := ms

				// 5) Торговый исполнитель (если задан)
				if c.Executor != nil {
					// отдельная горутина, чтобы не тормозить обработку тиков
					go c.Executor.ExecuteTriangle(ctx, tr, quotes, &msCopy, sf)
				}

				// 6) Лог
				c.printTriangle(now, tr, prof, quotes, &msCopy, sf)
			}

		case <-ctx.Done():
			return
		}
	}
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

	// Если MaxStartInfo нет (ms == nil) — печатаем "короткий" формат и уходим.
	if ms == nil {
		fmt.Fprintf(w, "[ARB] %+0.3f%%  %s\n", profit*100, t.Name)
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
		return
	}

	// ----- ниже ms уже гарантированно не nil -----

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

	if c.FeePerLeg > 0 {
		execs, okExec := simulateTriangleExecution(t, quotes, ms.StartAsset, safeStart, c.FeePerLeg)
		if okExec {
			fmt.Fprintln(w, "  Legs execution with fees:")
			for i, e := range execs {
				fmt.Fprintf(w,
					"    leg %d: %s  %.6f %s → %.6f %s  fee=%.8f %s\n",
					i+1, e.Symbol,
					e.AmountIn, e.From,
					e.AmountOut, e.To,
					e.FeeAmount, e.FeeAsset,
				)
			}
		}
	}

	fmt.Fprintln(w)
}

// ==============================
// DRY-RUN исполнитель треугольника
// ==============================

type DryRunExecutor struct {
	out io.Writer
}

func NewDryRunExecutor(out io.Writer) *DryRunExecutor {
	return &DryRunExecutor{out: out}
}

func (e *DryRunExecutor) ExecuteTriangle(
	ctx context.Context,
	t domain.Triangle,
	quotes map[string]domain.Quote,
	ms *domain.MaxStartInfo,
	startFraction float64,
) {
	if ms == nil {
		return
	}
	safeStart := ms.MaxStart * startFraction
	if safeStart <= 0 {
		return
	}

	execs, ok := simulateTriangleExecution(t, quotes, ms.StartAsset, safeStart, 0)
	if !ok || len(execs) == 0 {
		return
	}

	fmt.Fprintf(e.out, "  [DRY-RUN EXEC] start=%.6f %s (safeStart)\n", safeStart, ms.StartAsset)
	for i, lg := range execs {
		fmt.Fprintf(e.out,
			"    leg %d: %s  %.6f %s -> %.6f %s\n",
			i+1,
			lg.Symbol,
			lg.AmountIn, lg.From,
			lg.AmountOut, lg.To,
		)
	}
}

// ==============================
// Симуляция исполнения треугольника
// ==============================

type legExec struct {
	Symbol    string
	From      string
	To        string
	AmountIn  float64
	AmountOut float64
	FeeAmount float64
	FeeAsset  string
}

func simulateTriangleExecution(
	t domain.Triangle,
	quotes map[string]domain.Quote,
	startAsset string,
	startAmount float64,
	feePerLeg float64,
) ([]legExec, bool) {
	if startAmount <= 0 {
		return nil, false
	}

	curAsset := startAsset
	curAmount := startAmount
	var res []legExec

	for _, leg := range t.Legs {
		q, ok := quotes[leg.Symbol]
		if !ok || q.Bid <= 0 || q.Ask <= 0 {
			return nil, false
		}

		var from, to string
		if leg.Dir > 0 {
			from, to = leg.From, leg.To
		} else {
			from, to = leg.To, leg.From
		}

		if curAsset != from {
			return nil, false
		}

		base, quote, okPQ := detectBaseQuote(leg.Symbol, from, to)
		if !okPQ {
			return nil, false
		}

		prevAmount := curAmount
		var amountOut, feeAmount float64
		var feeAsset string

		switch {
		case curAsset == base:
			// продаём base → получаем quote по bid
			gross := curAmount * q.Bid
			feeAmount = gross * feePerLeg
			amountOut = gross - feeAmount
			feeAsset = quote
			curAsset, curAmount = quote, amountOut

		case curAsset == quote:
			// покупаем base за quote по ask
			gross := curAmount / q.Ask
			feeAmount = gross * feePerLeg
			amountOut = gross - feeAmount
			feeAsset = base
			curAsset, curAmount = base, amountOut

		default:
			return nil, false
		}

		res = append(res, legExec{
			Symbol:    leg.Symbol,
			From:      from,
			To:        to,
			AmountIn:  prevAmount,
			AmountOut: amountOut,
			FeeAmount: feeAmount,
			FeeAsset:  feeAsset,
		})
	}
	return res, true
}

func detectBaseQuote(symbol, a, b string) (base, quote string, ok bool) {
	if strings.HasPrefix(symbol, a) {
		return a, b, true
	}
	if strings.HasPrefix(symbol, b) {
		return b, a, true
	}
	return "", "", false
}

// ==============================
// Конвертация для вывода maxStart в USDT
// ==============================

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

// ==============================
// Работа с логом
// ==============================

func OpenLogWriter(path string) (io.WriteCloser, *bufio.Writer, io.Writer) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		log.Fatalf("open %s: %v", path, err)
	}
	buf := bufio.NewWriter(f)
	out := io.MultiWriter(os.Stdout, buf)
	return f, buf, out
}






1. mexc/trader.go

Создай файл crypt_proto/mexc/trader.go:

package mexc

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Trader struct {
	apiKey    string
	apiSecret string
	debug     bool

	client  *http.Client
	baseURL string
}

func NewTrader(apiKey, apiSecret string, debug bool) *Trader {
	return &Trader{
		apiKey:    strings.TrimSpace(apiKey),
		apiSecret: strings.TrimSpace(apiSecret),
		debug:     debug,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
		baseURL: "https://api.mexc.com",
	}
}

func (t *Trader) dlog(format string, args ...any) {
	if t.debug {
		log.Printf("[MEXC TRADER] "+format, args...)
	}
}

// PlaceMarket отправляет MARKET-ордер на MEXC Spot.
// quantity — в БАЗОВОЙ валюте символа (как в Spot API MEXC).
// side: "BUY" или "SELL".
func (t *Trader) PlaceMarket(
	ctx context.Context,
	symbol string,
	side string,
	quantity float64,
) error {
	if quantity <= 0 {
		return fmt.Errorf("quantity must be > 0, got %f", quantity)
	}
	side = strings.ToUpper(strings.TrimSpace(side))
	if side != "BUY" && side != "SELL" {
		return fmt.Errorf("invalid side %q", side)
	}

	endpoint := "/api/v3/order"

	params := url.Values{}
	params.Set("symbol", symbol)
	params.Set("side", side)
	params.Set("type", "MARKET")
	params.Set("quantity", fmt.Sprintf("%.8f", quantity))

	// Рекомендуется ставить recvWindow, чтобы избежать проблем с задержкой.
	params.Set("recvWindow", "5000")
	params.Set("timestamp", fmt.Sprintf("%d", time.Now().UnixMilli()))

	// Подпись HMAC SHA256(secret, totalParams)
	queryString := params.Encode()
	signature := t.sign(queryString)
	params.Set("signature", signature)

	body := params.Encode()

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		t.baseURL+endpoint,
		strings.NewReader(body),
	)
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}

	req.Header.Set("X-MEXC-APIKEY", t.apiKey)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	t.dlog("PlaceMarket %s %s qty=%.8f body=%s", side, symbol, quantity, body)

	resp, err := t.client.Do(req)
	if err != nil {
		return fmt.Errorf("http do: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("mexc order error: status=%d body=%s", resp.StatusCode, string(respBody))
	}

	t.dlog("PlaceMarket OK: %s", string(respBody))
	return nil
}

func (t *Trader) sign(payload string) string {
	mac := hmac.New(sha256.New, []byte(t.apiSecret))
	mac.Write([]byte(payload))
	return hex.EncodeToString(mac.Sum(nil))
}


Это уже полноценный клиент под POST /api/v3/order с подписью и MARKET-ордером по количеству в базовой валюте.

2. arb/executor_real.go

Создай файл crypt_proto/arb/executor_real.go:

package arb

import (
	"context"
	"fmt"
	"io"
	"time"

	"crypt_proto/domain"
)

// SpotTrader — минимальный интерфейс для трейдера биржи.
type SpotTrader interface {
	PlaceMarket(ctx context.Context, symbol, side string, quantity float64) error
}

// RealExecutor — реальный исполнитель треугольника через SpotTrader.
type RealExecutor struct {
	trader SpotTrader
	out    io.Writer
}

func NewRealExecutor(trader SpotTrader, out io.Writer) *RealExecutor {
	return &RealExecutor{
		trader: trader,
		out:    out,
	}
}

func (e *RealExecutor) log(format string, args ...any) {
	if e.out == nil {
		return
	}
	fmt.Fprintf(e.out, format+"\n", args...)
}

// ExecuteTriangle запускает реальную торговлю по треугольнику.
// ВАЖНО: предполагаем, что StartAsset = USDT (BaseAsset) и что
// фильтры по MinStart и profitability уже прошли в Consumer.
func (e *RealExecutor) ExecuteTriangle(
	ctx context.Context,
	t domain.Triangle,
	quotes map[string]domain.Quote,
	ms *domain.MaxStartInfo,
	startFraction float64,
) {
	if e.trader == nil || ms == nil {
		return
	}

	// Сейчас для простоты торгуем только треугольники с USDT-стартом.
	if ms.StartAsset != BaseAsset {
		e.log("  [REAL EXEC] skip triangle %s: StartAsset=%s != %s",
			t.Name, ms.StartAsset, BaseAsset)
		return
	}

	safeStart := ms.MaxStart * startFraction
	if safeStart <= 0 {
		return
	}

	// считаем путь исполнения без комиссии (fee=0), чтобы получить объёмы по ногам
	execs, ok := simulateTriangleExecution(t, quotes, ms.StartAsset, safeStart, 0)
	if !ok || len(execs) == 0 {
		e.log("  [REAL EXEC] simulate failed for %s", t.Name)
		return
	}

	e.log("  [REAL EXEC] start=%.6f %s triangle=%s", safeStart, ms.StartAsset, t.Name)

	// Локальный таймаут на всю цепочку
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	for i, lg := range execs {
		// По legExec:
		//  - AmountIn — сколько у нас ВХОДА на этой ноге (в lg.From)
		//  - AmountOut — сколько выйдет (в lg.To)
		// Нужно определить side и quantity в базе символа.

		base, quote, ok := detectBaseQuote(lg.Symbol, lg.From, lg.To)
		if !ok {
			e.log("    [REAL EXEC] leg %d: cannot detect base/quote for %s (%s->%s)",
				i+1, lg.Symbol, lg.From, lg.To)
			return
		}

		side := ""
		var qty float64

		switch {
		case lg.From == base && lg.To == quote:
			// Продаём base -> получаем quote
			side = "SELL"
			qty = lg.AmountIn // количество base, которое продаём
		case lg.From == quote && lg.To == base:
			// Покупаем base за quote
			side = "BUY"
			qty = lg.AmountOut // количество base, которое хотим купить
		default:
			e.log("    [REAL EXEC] leg %d: inconsistent path %s (%s->%s, base=%s, quote=%s)",
				i+1, lg.Symbol, lg.From, lg.To, base, quote)
			return
		}

		if qty <= 0 {
			e.log("    [REAL EXEC] leg %d: qty<=0, skip (%s %s)", i+1, side, lg.Symbol)
			return
		}

		e.log("    [REAL EXEC] leg %d: %s %s qty=%.8f", i+1, side, lg.Symbol, qty)

		if err := e.trader.PlaceMarket(ctx, lg.Symbol, side, qty); err != nil {
			e.log("    [REAL EXEC] leg %d ERROR: %v", i+1, err)
			// На первой же ошибке останавливаем весь треугольник.
			return
		}
	}

	e.log("  [REAL EXEC] done triangle %s", t.Name)
}


⚠️ Тут важно:

Реальный запуск упирается в твой баланс — при нуле будут ошибки insufficient balance.

Мы специально делаем лог и немедленный выход при ошибке любой ноги, чтобы не продолжать сломанную цепочку.

3. Обновление main.go — выбор между DRY-RUN и REAL

Теперь в твоём main.go вот этот кусок:

	// запускаем потребителя
	consumer := arb.NewConsumer(cfg.FeePerLeg, cfg.MinProfit, cfg.MinStart, arbOut)
	consumer.StartFraction = cfg.StartFraction

	// пока только DRY-RUN: считаем, логируем, но не шлём реальные ордера
	consumer.Executor = arb.NewDryRunExecutor(arbOut)

	// позже сюда можно будет повесить реального исполнителя:
	// if cfg.TradeEnabled && cfg.APIKey != "" && cfg.APISecret != "" {
	//     trader := mexc.NewTrader(cfg.APIKey, cfg.APISecret, cfg.Debug)
	//     consumer.Executor = arb.NewRealExecutor(trader, cfg.FeePerLeg, cfg.MinStart)
	// }


замени на:

	// запускаем потребителя
	consumer := arb.NewConsumer(cfg.FeePerLeg, cfg.MinProfit, cfg.MinStart, arbOut)
	consumer.StartFraction = cfg.StartFraction

	// Выбор исполнителя
	hasKeys := cfg.APIKey != "" && cfg.APISecret != ""

	if cfg.TradeEnabled && hasKeys {
		log.Printf("[EXEC] REAL TRADING ENABLED on %s", cfg.Exchange)

		switch cfg.Exchange {
		case "MEXC":
			trader := mexc.NewTrader(cfg.APIKey, cfg.APISecret, cfg.Debug)
			consumer.Executor = arb.NewRealExecutor(trader, arbOut)
		default:
			log.Printf("[EXEC] Real trading not implemented for %s, fallback to DRY-RUN", cfg.Exchange)
			consumer.Executor = arb.NewDryRunExecutor(arbOut)
		}
	} else {
		log.Printf("[EXEC] DRY-RUN MODE (TRADE_ENABLED=%v, hasKeys=%v) — реальные ордера отправляться не будут",
			cfg.TradeEnabled, hasKeys)
		consumer.Executor = arb.NewDryRunExecutor(arbOut)
	}


Остальная часть main.go остаётся как есть.

4. Напоминание про .env

Чтобы НИЧЕГО не ушло на биржу случайно, оставь пока так:

TRADE_ENABLED=false
MEXC_API_KEY=
MEXC_API_SECRET=


Когда будешь готов реально торговать:

Задаёшь:

TRADE_ENABLED=true
MEXC_API_KEY=твой_ключ
MEXC_API_SECRET=твой_секрет


Убедись, что:

баланс на споте есть,

MIN_START_USDT и START_FRACTION разумные.

Если хочешь, следующим шагом можем:

добавить проверку баланса USDT перед сделкой (GET /api/v3/account или аналог),

или сделать режим: log-only реальных ошибок ордеров, чтобы видеть, что отваливается (min notional, step size и т.п.).


