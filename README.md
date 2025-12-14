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

	// MinStart — минимальный допустимый старт в USDT (MIN_START_USDT). 0 = фильтр выключен.
	MinStart float64

	// StartFraction — доля от maxStart, которую считаем безопасной (обычно 0.5).
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

					// 3) ФИЛЬТР MIN_START_USDT — сравниваем safeStart В USDT
					if c.MinStart > 0 {
						safeUSDT, okConv := convertToUSDT(safeStart, ms.StartAsset, quotes)
						// если конвертацию сделать нельзя — пропускаем, раз порог задан в USDT
						if !okConv || safeUSDT < c.MinStart {
							continue
						}
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

		// ====== ДОП: объёмы всех ног в USDT (для safeStart) ======
		flows, okF := calcTriangleFlow(t, quotes, c.FeePerLeg, safeStart)
		if okF {
			// Покажем детально по каждой ноге
			for i := 0; i < 3; i++ {
				inUSDT, okIn := convertToUSDT(flows[i].InAmt, flows[i].InAsset, quotes)
				outUSDT, okOut := convertToUSDT(flows[i].OutAmt, flows[i].OutAsset, quotes)

				inStr := "?"
				outStr := "?"
				if okIn {
					inStr = fmt.Sprintf("%.6f", inUSDT)
				}
				if okOut {
					outStr = fmt.Sprintf("%.6f", outUSDT)
				}

				fmt.Fprintf(w, "  leg%d: %.6f %s (~%s USDT) -> %.6f %s (~%s USDT)\n",
					i+1,
					flows[i].InAmt, flows[i].InAsset, inStr,
					flows[i].OutAmt, flows[i].OutAsset, outStr,
				)
			}

			// И короткой строкой “только USDT”
			in1, okIn1 := convertToUSDT(flows[0].InAmt, flows[0].InAsset, quotes)
			out1, okOut1 := convertToUSDT(flows[0].OutAmt, flows[0].OutAsset, quotes)
			in2, okIn2 := convertToUSDT(flows[1].InAmt, flows[1].InAsset, quotes)
			out2, okOut2 := convertToUSDT(flows[1].OutAmt, flows[1].OutAsset, quotes)
			in3, okIn3 := convertToUSDT(flows[2].InAmt, flows[2].InAsset, quotes)
			out3, okOut3 := convertToUSDT(flows[2].OutAmt, flows[2].OutAsset, quotes)

			if okIn1 && okOut1 && okIn2 && okOut2 && okIn3 && okOut3 {
				fmt.Fprintf(w, "  legsUSDT: in1=%.6f out1=%.6f | in2=%.6f out2=%.6f | in3=%.6f out3=%.6f\n",
					in1, out1, in2, out2, in3, out3)
			} else {
				fmt.Fprintf(w, "  legsUSDT: (some legs can't convert to USDT with current subscribed books)\n")
			}
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
// Конвертация в USDT (для логов/фильтра)
// ==============================

func convertToUSDT(amount float64, asset string, quotes map[string]domain.Quote) (float64, bool) {
	if amount <= 0 {
		return 0, false
	}
	if asset == "USDT" {
		return amount, true
	}

	// 1) ASSETUSDT: продаём ASSET за USDT по bid
	if q, ok := quotes[asset+"USDT"]; ok && q.Bid > 0 && q.BidQty > 0 {
		return amount * q.Bid, true
	}

	// 2) USDTASSET: покупаем USDT за ASSET по ask (ask = ASSET per USDT)
	if q, ok := quotes["USDT"+asset]; ok && q.Ask > 0 && q.AskQty > 0 {
		return amount / q.Ask, true
	}

	// 3) Через USDC
	amtUSDC, ok1 := convertViaQuote(amount, asset, "USDC", quotes)
	if ok1 {
		amtUSDT, ok2 := convertViaQuote(amtUSDC, "USDC", "USDT", quotes)
		if ok2 {
			return amtUSDT, true
		}
	}

	return 0, false
}

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

	// TOFROM: buy TO using FROM (ask = FROM per TO)
	if q, ok := quotes[assetTo+assetFrom]; ok && q.Ask > 0 && q.AskQty > 0 {
		return amount / q.Ask, true
	}

	return 0, false
}

// ==============================
// Расчёт потока по ногам (для safeStart)
// ==============================

type legFlow struct {
	InAmt    float64
	InAsset  string
	OutAmt   float64
	OutAsset string
}

// calcTriangleFlow считает реальные вход/выход каждой ноги для заданного старта (start в стартовой валюте треугольника).
// Комиссию считаем как удержание из результата ноги (out *= 1-fee).
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




