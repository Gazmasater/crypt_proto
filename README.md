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
			if okPrev &&
				prev.Bid == ev.Bid &&
				prev.Ask == ev.Ask &&
				prev.BidQty == ev.BidQty &&
				prev.AskQty == ev.AskQty {
				continue
			}

			quotes[ev.Symbol] = domain.Quote{
				Bid:    ev.Bid,
				Ask:    ev.Ask,
				BidQty: ev.BidQty,
				AskQty: ev.AskQty,
			}

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

				// копируем ms, чтобы не гонять один и тот же указатель
				msCopy := ms

				// 5) Торговый исполнитель (если задан)
				if c.Executor != nil {
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

	// Если MaxStartInfo нет (ms == nil) — печатаем "короткий" формат и выходим.
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

	// ниже ms уже гарантированно не nil
	bneckSym := ""
	if ms.BottleneckLeg >= 0 && ms.BottleneckLeg < len(t.Legs) {
		bneckSym = t.Legs[ms.BottleneckLeg].Symbol
	}

	safeStart := ms.MaxStart * startFraction
	maxUSDT, okMax := convertToUSDT(ms.MaxStart, ms.StartAsset, quotes)
	safeUSDT, okSafe := convertToUSDT(safeStart, ms.StartAsset, quotes)

	maxUSDTStr := "?"
	safeUSDTStr := "?"
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
	// pprof сервер
	go func() {
		log.Println("pprof on http://localhost:6060/debug/pprof/")
		if err := http.ListenAndServe("localhost:6060", nil); err != nil {
			log.Printf("pprof server error: %v", err)
		}
	}()

	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	// --------- конфиг ---------
	cfg := config.Load()

	// --------- треугольники ---------
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
	log.Printf("символов для подписки всего: %d", len(symbols))

	// --------- выбор биржи (market data feed) ---------
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

	// --------- лог-файл для арбитража ---------
	logFile, logBuf, arbOut := arb.OpenLogWriter("arbitrage.log")
	defer logFile.Close()
	defer logBuf.Flush()

	// --------- контекст с остановкой по сигналу ---------
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	events := make(chan domain.Event, 8192)
	var wg sync.WaitGroup

	// --------- потребитель арбитража ---------
	consumer := arb.NewConsumer(cfg.FeePerLeg, cfg.MinProfit, cfg.MinStart, arbOut)
	consumer.StartFraction = cfg.StartFraction

	// --------- выбор исполнителя (DRY-RUN или REAL) ---------
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
		log.Printf(
			"[EXEC] DRY-RUN MODE (TRADE_ENABLED=%v, hasKeys=%v) — реальные ордера отправляться не будут",
			cfg.TradeEnabled, hasKeys,
		)
		consumer.Executor = arb.NewDryRunExecutor(arbOut)
	}

	// Стартуем потребителя
	consumer.Start(ctx, events, triangles, indexBySymbol, &wg)

	// --------- фид биржи (котировки) ---------
	feed.Start(ctx, &wg, symbols, cfg.BookInterval, events)

	// --------- ждём сигнал остановки ---------
	<-ctx.Done()
	log.Println("shutting down...")

	// немного времени на дообработку
	time.Sleep(200 * time.Millisecond)
	close(events)
	wg.Wait()
	log.Println("bye")
}



