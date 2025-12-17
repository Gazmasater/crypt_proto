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

type tradeReq struct {
	id        int
	tr        domain.Triangle
	quotesSnap map[string]domain.Quote // снапшот ТОЛЬКО нужных символов треугольника
	enqueuedAt time.Time
}

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

	// котировки (живые, обновляются из events)
	quotes map[string]domain.Quote

	// очередь сделок (FIFO)
	tradeQ chan tradeReq

	mu sync.Mutex
	// cooldown per triangle (по времени ПОСЛЕ исполнения)
	lastTrade map[int]time.Time
	// чтобы не ставить один и тот же треугольник в очередь много раз
	queued map[int]bool

	// анти-спам печати
	lastPrint map[int]time.Time
	minPrint  time.Duration
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
		Executor:        nil,
		writer:          out,

		quotes:    make(map[string]domain.Quote),
		tradeQ:    make(chan tradeReq, 256),
		lastTrade: make(map[int]time.Time),
		queued:    make(map[int]bool),

		lastPrint: make(map[int]time.Time),
		minPrint:  5 * time.Millisecond,
	}
}

func (c *Consumer) Start(
	ctx context.Context,
	events <-chan domain.Event,
	triangles []domain.Triangle,
	indexBySymbol map[string][]int,
	wg *sync.WaitGroup,
) {
	// 1) worker (последовательное исполнение)
	wg.Add(1)
	go func() {
		defer wg.Done()
		c.tradeWorker(ctx)
	}()

	// 2) обработка market-data
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

			// обновляем quote, если изменился
			prev, okPrev := c.quotes[ev.Symbol]
			if okPrev &&
				prev.Bid == ev.Bid && prev.Ask == ev.Ask &&
				prev.BidQty == ev.BidQty && prev.AskQty == ev.AskQty {
				continue
			}
			c.quotes[ev.Symbol] = domain.Quote{Bid: ev.Bid, Ask: ev.Ask, BidQty: ev.BidQty, AskQty: ev.AskQty}

			trIDs := indexBySymbol[ev.Symbol]
			if len(trIDs) == 0 {
				continue
			}

			now := time.Now()

			for _, id := range trIDs {
				tr := triangles[id]

				prof, ok := domain.EvalTriangle(tr, c.quotes, c.FeePerLeg)
				if !ok || prof < c.MinProfit {
					continue
				}

				ms, okMS := domain.ComputeMaxStartTopOfBook(tr, c.quotes, c.FeePerLeg)
				if !okMS {
					continue
				}

				safeStart := ms.MaxStart * sf

				// фильтр по MIN_START_USDT (по safeStart)
				if c.MinStart > 0 {
					safeUSDT, okConv := convertToUSDT(safeStart, ms.StartAsset, c.quotes)
					if !okConv || safeUSDT < c.MinStart {
						continue
					}
				}

				// анти-спам печати
				if last, okLast := c.lastPrint[id]; !okLast || now.Sub(last) >= c.minPrint {
					c.lastPrint[id] = now
					msCopy := ms
					c.printTriangle(now, tr, prof, c.quotes, &msCopy, sf)
				}

				// торговля: поставить в очередь (снапшот)
				c.tryEnqueueTrade(ctx, now, id, tr)
			}
		}
	}
}

func (c *Consumer) tryEnqueueTrade(ctx context.Context, now time.Time, id int, tr domain.Triangle) {
	if !c.TradeEnabled || c.Executor == nil {
		return
	}

	// проверим cooldown и queued/inflight
	c.mu.Lock()
	if last, ok := c.lastTrade[id]; ok && now.Sub(last) < c.TradeCooldown {
		c.mu.Unlock()
		return
	}
	if c.queued[id] {
		c.mu.Unlock()
		return
	}
	// помечаем queued сразу (чтобы не спамить)
	c.queued[id] = true

	// делаем снапшот котировок только нужных символов
	snap := make(map[string]domain.Quote, len(tr.Legs))
	okSnap := true
	for _, leg := range tr.Legs {
		q, ok := c.quotes[leg.Symbol]
		if !ok || q.Bid <= 0 || q.Ask <= 0 {
			okSnap = false
			break
		}
		snap[leg.Symbol] = q
	}
	c.mu.Unlock()

	if !okSnap {
		// снимем флаг queued, т.к. снапшот не удался
		c.mu.Lock()
		c.queued[id] = false
		c.mu.Unlock()
		return
	}

	req := tradeReq{
		id:         id,
		tr:         tr,
		quotesSnap: snap,
		enqueuedAt: now,
	}

	// не блокируем навсегда: если очередь забита — просто отпускаем queued и пропускаем
	select {
	case c.tradeQ <- req:
		// ok
	case <-ctx.Done():
		c.mu.Lock()
		c.queued[id] = false
		c.mu.Unlock()
	default:
		// очередь переполнена
		c.mu.Lock()
		c.queued[id] = false
		c.mu.Unlock()
	}
}

func (c *Consumer) tradeWorker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return

		case req := <-c.tradeQ:
			// снимаем queued (теперь "в работе")
			c.mu.Lock()
			c.queued[req.id] = false
			c.mu.Unlock()

			// исполняем строго последовательно (мы в одном воркере)
			_ = c.Executor.Execute(ctx, req.tr, req.quotesSnap, c.TradeAmountUSDT)

			// cooldown считаем от факта исполнения (после Execute)
			c.mu.Lock()
			c.lastTrade[req.id] = time.Now()
			c.mu.Unlock()
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




