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
	mu        sync.Mutex
	inFlight  bool
	lastTrade map[int]time.Time
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
