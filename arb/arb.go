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

type Consumer struct {
	FeePerLeg float64
	MinProfit float64

	// MinStart — минимальный допустимый старт (обычно в USDT, но по факту это валюта старта треугольника).
	// ВАЖНО: фильтр применяется к safeStart (= maxStart * StartFraction).
	MinStart float64

	// StartFraction — доля от maxStart, которую считаем безопасной для исполнения (обычно 0.5).
	StartFraction float64

	writer io.Writer
}

func NewConsumer(feePerLeg, minProfit, minStart float64, out io.Writer) *Consumer {
	return &Consumer{
		FeePerLeg:     feePerLeg,
		MinProfit:     minProfit,
		MinStart:      minStart,
		StartFraction: 0.5, // дефолт: половина от maxStart
		writer:        out,
	}
}

// Start запускает горутину-потребителя.
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

	// анти-спам: один и тот же треугольник не печатаем чаще этого интервала
	const minPrintInterval = 5 * time.Millisecond
	lastPrint := make(map[int]time.Time)

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

			// если книга не изменилась — пропускаем
			if prev, okPrev := quotes[ev.Symbol]; okPrev &&
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

				// 1) прибыль (уже с учётом комиссии)
				prof, ok := domain.EvalTriangle(tr, quotes, c.FeePerLeg)
				if !ok || prof < c.MinProfit {
					continue
				}

				// 2) maxStart по top-of-book
				ms, okMS := domain.ComputeMaxStartTopOfBook(tr, quotes, c.FeePerLeg)
				if okMS {
					safeStart := ms.MaxStart * sf

					// фильтр применяем к safeStart (половина maxStart)
					if c.MinStart > 0 && safeStart < c.MinStart {
						continue
					}
				}

				// анти-спам
				if last, okLast := lastPrint[id]; okLast && now.Sub(last) < minPrintInterval {
					continue
				}
				lastPrint[id] = now

				var msPtr *domain.MaxStartInfo
				if okMS {
					msCopy := ms
					msPtr = &msCopy
				}

				c.printTriangle(now, tr, prof, quotes, msPtr, sf)
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

	if ms != nil {
		bneckSym := ""
		if ms.BottleneckLeg >= 0 && ms.BottleneckLeg < len(t.Legs) {
			bneckSym = t.Legs[ms.BottleneckLeg].Symbol
		}

		// safeStart = половина от maxStart (или другой startFraction)
		safeStart := ms.MaxStart * startFraction

		// конвертация maxStart/safeStart в USDT для вывода
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
	} else {
		fmt.Fprintf(w, "[ARB] %+0.3f%%  %s\n", profit*100, t.Name)
	}

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

func OpenLogWriter(path string) (io.WriteCloser, *bufio.Writer, io.Writer) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		log.Fatalf("open %s: %v", path, err)
	}
	buf := bufio.NewWriter(f)
	out := io.MultiWriter(os.Stdout, buf)
	return f, buf, out
}

// ==============================
// Конвертация для вывода maxStart в USDT
// ==============================

func convertToUSDT(amount float64, asset string, quotes map[string]domain.Quote) (float64, bool) {
	if amount <= 0 {
		return 0, false
	}
	if asset == "USDT" {
		return amount, true
	}

	// 1) Прямая пара: ASSETUSDT (продаём ASSET за USDT по bid)
	if q, ok := quotes[asset+"USDT"]; ok && q.Bid > 0 && q.BidQty > 0 {
		return amount * q.Bid, true
	}

	// 2) Прямая пара: USDTASSET (покупаем USDT за ASSET по ask)
	if q, ok := quotes["USDT"+asset]; ok && q.Ask > 0 && q.AskQty > 0 {
		// ask = ASSET per USDT => USDT = ASSET / ask
		return amount / q.Ask, true
	}

	// 3) Через USDC: ASSET -> USDC -> USDT
	amtUSDC, ok1 := convertViaQuote(amount, asset, "USDC", quotes)
	if ok1 {
		amtUSDT, ok2 := convertViaQuote(amtUSDC, "USDC", "USDT", quotes)
		if ok2 {
			return amtUSDT, true
		}
	}

	return 0, false
}

// convertViaQuote конвертирует amount из assetFrom в assetTo, используя только прямые пары.
// Правила:
// - FROMTO: продаём FROM за TO по bid => out = amount * bid
// - TOFROM: покупаем TO за FROM по ask => out = amount / ask
func convertViaQuote(amount float64, assetFrom, assetTo string, quotes map[string]domain.Quote) (float64, bool) {
	if amount <= 0 {
		return 0, false
	}
	if assetFrom == assetTo {
		return amount, true
	}

	// FROMTO: sell FROM -> TO
	if q, ok := quotes[assetFrom+assetTo]; ok && q.Bid > 0 && q.BidQty > 0 {
		return amount * q.Bid, true
	}

	// TOFROM: buy TO using FROM
	if q, ok := quotes[assetTo+assetFrom]; ok && q.Ask > 0 && q.AskQty > 0 {
		return amount / q.Ask, true
	}

	return 0, false
}
