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

type Consumer struct {
	FeePerLeg float64
	MinProfit float64

	// MinStart — минимальный старт в USDT (MIN_START_USDT). 0 = фильтр выключен.
	MinStart float64

	// StartFraction — доля от maxStart, которую считаем безопасной (0..1). Например 0.5.
	StartFraction float64

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

				// 1) прибыль (уже с учетом комиссии)
				prof, ok := domain.EvalTriangle(tr, quotes, c.FeePerLeg)
				if !ok || prof < c.MinProfit {
					continue
				}

				// 2) maxStart по top-of-book
				ms, okMS := domain.ComputeMaxStartTopOfBook(tr, quotes, c.FeePerLeg)

				// 3) фильтр MIN_START_USDT применяем к safeStart, но в USDT
				if okMS && c.MinStart > 0 {
					safeStart := ms.MaxStart * sf
					safeUSDT, okConv := convertToUSDT(safeStart, ms.StartAsset, quotes)
					// если не можем конвертировать в USDT — лучше пропустить, раз порог в USDT
					if !okConv || safeUSDT < c.MinStart {
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

		// ====== ДОП: объёмы 1-й и 2-й ноги для safeStart ======
		if flows, okF := calcTriangleFlow(t, quotes, c.FeePerLeg, safeStart); okF {
			fmt.Fprintf(w, "  leg1: %s -> %s\n",
				fmtAmtWithUSDT(flows[0].InAmt, flows[0].InAsset, quotes),
				fmtAmtWithUSDT(flows[0].OutAmt, flows[0].OutAsset, quotes),
			)
			fmt.Fprintf(w, "  leg2: %s -> %s\n",
				fmtAmtWithUSDT(flows[1].InAmt, flows[1].InAsset, quotes),
				fmtAmtWithUSDT(flows[1].OutAmt, flows[1].OutAsset, quotes),
			)
		}
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

// ==============================
// Вывод "реальных" объёмов 1-й и 2-й ноги для safeStart
// ==============================

type legFlow struct {
	InAmt    float64
	InAsset  string
	OutAmt   float64
	OutAsset string
}

func calcTriangleFlow(t domain.Triangle, quotes map[string]domain.Quote, fee float64, start float64) ([3]legFlow, bool) {
	var flows [3]legFlow
	amt := start

	for i, leg := range t.Legs {
		q, ok := quotes[leg.Symbol]
		if !ok || q.Bid <= 0 || q.Ask <= 0 {
			return flows, false
		}

		inAmt := amt
		outAmt := 0.0

		// ВАЖНО: in/out asset тут берём из leg.From/leg.To — они уже заданы направлением треугольника.
		// Dir>0: From(base) -> To(quote) SELL по bid
		if leg.Dir > 0 {
			if q.BidQty <= 0 {
				return flows, false
			}
			outAmt = inAmt * q.Bid
		} else {
			// Dir<0: From(quote) -> To(base) BUY по ask
			if q.AskQty <= 0 {
				return flows, false
			}
			outAmt = inAmt / q.Ask
		}

		// комиссия удерживается из результата ноги
		outAmt *= (1 - fee)

		flows[i] = legFlow{
			InAmt:    inAmt,
			InAsset:  leg.From,
			OutAmt:   outAmt,
			OutAsset: leg.To,
		}

		amt = outAmt
		if amt <= 0 {
			return flows, false
		}
	}

	return flows, true
}

func fmtAmtWithUSDT(amount float64, asset string, quotes map[string]domain.Quote) string {
	usdt, ok := convertToUSDT(amount, asset, quotes)
	if ok {
		return fmt.Sprintf("%.6f %s (%.6f USDT)", amount, asset, usdt)
	}
	return fmt.Sprintf("%.6f %s (? USDT)", amount, asset)
}





