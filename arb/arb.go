package arb

import (
	"bufio"
	"context"
	"crypt_proto/domain"
	"crypt_proto/mexc"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

/* =========================  TRIANGLE EVAL + PRINT  ========================= */

func evalTriangle(t domain.Triangle, quotes map[string]domain.Quote, fee float64) (float64, bool) {
	amt := 1.0

	for _, leg := range t.Legs {
		q, ok := quotes[leg.Symbol]
		if !ok || q.Bid <= 0 || q.Ask <= 0 {
			return 0, false
		}

		if leg.Dir > 0 {
			amt *= q.Bid
		} else {
			amt /= q.Ask
		}

		amt *= (1 - fee)
		if amt <= 0 {
			return 0, false
		}
	}

	return amt - 1.0, true
}

func printTriangle(w io.Writer, t domain.Triangle, profit float64, quotes map[string]domain.Quote) {
	ts := time.Now().Format("2006-01-02 15:04:05.000")
	fmt.Fprintf(w, "%s\n", ts)
	fmt.Fprintf(w, "[ARB] %+0.3f%%  %s\n", profit*100, t.Name)

	for _, leg := range t.Legs {
		q := quotes[leg.Symbol]
		mid := (q.Bid + q.Ask) / 2
		spreadAbs := q.Ask - q.Bid
		spreadPct := 0.0
		if mid > 0 {
			spreadPct = spreadAbs / mid * 100
		}
		side := fmt.Sprintf("%s/%s", leg.From, leg.To)
		if leg.Dir < 0 {
			side = fmt.Sprintf("%s/%s", leg.To, leg.From)
		}

		fmt.Fprintf(
			w,
			"  %s (%s): bid=%.10f ask=%.10f  spread=%.10f (%.5f%%)  bidQty=%.4f askQty=%.4f\n",
			leg.Symbol, side,
			q.Bid, q.Ask,
			spreadAbs, spreadPct,
			q.BidQty, q.AskQty,
		)
	}
	fmt.Fprintln(w)
}

/* =========================  ARB LOGGING + PIPELINE  ========================= */

func InitArbLogger(path string) (io.Writer, func()) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		log.Fatalf("open %s: %v", path, err)
	}
	buf := bufio.NewWriter(f)

	out := io.MultiWriter(os.Stdout, buf)

	cleanup := func() {
		_ = buf.Flush()
		_ = f.Close()
	}

	return out, cleanup
}

func StartWSWorkers(
	ctx context.Context,
	wg *sync.WaitGroup,
	symbols []string,
	interval string,
	out chan<- domain.Event,
) {
	const maxPerConn = 50

	var chunks [][]string
	for i := 0; i < len(symbols); i += maxPerConn {
		j := i + maxPerConn
		if j > len(symbols) {
			j = len(symbols)
		}
		chunks = append(chunks, symbols[i:j])
	}

	log.Printf("будем использовать %d WS-подключений", len(chunks))

	for idx, chunk := range chunks {
		wg.Add(1)
		go mexc.RunPublicBookTickerWS(ctx, wg, idx, chunk, interval, out)
	}
}

func ConsumeEvents(
	ctx context.Context,
	events <-chan domain.Event,
	triangles []domain.Triangle,
	indexBySymbol map[string][]int,
	feePerLeg, minProfit float64,
	out io.Writer,
) {
	quotes := make(map[string]domain.Quote)

	const minPrintInterval = 5 * time.Millisecond
	lastPrint := make(map[int]time.Time)

	for {
		select {
		case ev, ok := <-events:
			if !ok {
				return
			}

			// если котировка по символу не изменилась, не пересчитываем
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
				prof, ok := evalTriangle(tr, quotes, feePerLeg)
				if !ok || prof < minProfit {
					continue
				}

				if last, okLast := lastPrint[id]; okLast && now.Sub(last) < minPrintInterval {
					continue
				}
				lastPrint[id] = now

				printTriangle(out, tr, prof, quotes)
			}
		case <-ctx.Done():
			return
		}
	}
}
