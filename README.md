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
	"container/heap"
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"sync"
	"time"

	"crypt_proto/domain"
)

type Consumer struct {
	FeePerLeg float64
	MinProfit float64

	writer io.Writer
}

func NewConsumer(feePerLeg, minProfit float64, out io.Writer) *Consumer {
	return &Consumer{
		FeePerLeg: feePerLeg,
		MinProfit: minProfit,
		writer:    out,
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

	// граф конвертации (для maxStart -> USDT, если USDT нет в самом треугольнике)
	adj := buildAdjFromTriangles(triangles)

	const minPrintInterval = 5 * time.Millisecond
	lastPrint := make(map[int]time.Time)

	for {
		select {
		case ev, ok := <-events:
			if !ok {
				return
			}

			// дедуп по одинаковому стакану
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

				prof, ok := domain.EvalTriangle(tr, quotes, c.FeePerLeg)
				if !ok {
					continue
				}
				if prof < c.MinProfit {
					continue
				}

				if last, okLast := lastPrint[id]; okLast {
					if now.Sub(last) < minPrintInterval {
						continue
					}
				}
				lastPrint[id] = now

				maxStart, bottleneck, okMS := calcMaxStart(tr, quotes, c.FeePerLeg)
				var maxStartUSDT float64
				if okMS {
					maxStartUSDT = calcMaxStartUSDTAlways(tr, maxStart, quotes, c.FeePerLeg, adj)
				} else {
					maxStartUSDT = math.NaN()
				}

				c.printTriangle(now, tr, prof, maxStart, maxStartUSDT, bottleneck, quotes)
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
	maxStart float64,
	maxStartUSDT float64,
	bottleneck string,
	quotes map[string]domain.Quote,
) {
	w := c.writer

	startAsset := t.Legs[0].From
	if bottleneck == "" {
		bottleneck = "-"
	}

	fmt.Fprintf(w, "%s\n", ts.Format("2006-01-02 15:04:05.000"))

	// ВСЕГДА печатаем "(... USDT)"
	if maxStart > 0 {
		fmt.Fprintf(w,
			"[ARB] %+0.3f%%  %s  maxStart=%.4f %s (%.2f USDT)  bottleneck=%s\n",
			profit*100, t.Name,
			maxStart, startAsset, maxStartUSDT,
			bottleneck,
		)
	} else {
		// даже если maxStart не посчитался/0 — скобки остаются
		fmt.Fprintf(w,
			"[ARB] %+0.3f%%  %s  maxStart=%.4f %s (%.2f USDT)  bottleneck=%s\n",
			profit*100, t.Name,
			0.0, startAsset, maxStartUSDT,
			bottleneck,
		)
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

		fmt.Fprintf(w,
			"  %s (%s): bid=%.10f ask=%.10f  spread=%.10f (%.5f%%)  bidQty=%.4f askQty=%.4f\n",
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

// -------------------- maxStart --------------------

// calcMaxStart оценивает максимальный стартовый объём в валюте t.Legs[0].From,
// исходя из ограничений ликвидности (bidQty/askQty) по каждой ноге.
func calcMaxStart(t domain.Triangle, quotes map[string]domain.Quote, fee float64) (maxStart float64, bottleneck string, ok bool) {
	start := t.Legs[0].From
	_ = start

	prod := 1.0                    // сколько валюты текущей ноги.From получится из 1 единицы старта (после предыдущих ног)
	maxStart = math.Inf(1)         // минимизируем
	bottleneck = ""

	for i := 0; i < 3; i++ {
		leg := t.Legs[i]
		q, okQ := quotes[leg.Symbol]
		if !okQ || q.Bid <= 0 || q.Ask <= 0 {
			return 0, "", false
		}

		// лимит в валюте leg.From (до комиссии на этой ноге)
		var capFrom float64
		var k float64 // мультипликатор перехода From->To (с комиссией)

		if leg.Dir > 0 {
			// base->quote: продаём base, ограничение bidQty (в base)
			capFrom = q.BidQty
			k = q.Bid * (1 - fee)
		} else {
			// quote->base: покупаем base за quote
			// askQty в base => максимум quote, который можно потратить: askQty * ask
			capFrom = q.AskQty * q.Ask
			k = (1.0 / q.Ask) * (1 - fee)
		}

		if capFrom <= 0 || k <= 0 || prod <= 0 {
			return 0, "", false
		}

		limitStart := capFrom / prod
		if limitStart < maxStart {
			maxStart = limitStart
			bottleneck = leg.Symbol
		}

		prod *= k
	}

	if !math.IsInf(maxStart, 1) && maxStart > 0 {
		return maxStart, bottleneck, true
	}
	return 0, bottleneck, false
}

// calcMaxStartUSDTAlways ВСЕГДА возвращает число для печати в скобках "(... USDT)".
// 1) если USDT встречается по пути треугольника — конвертим по ногам треугольника;
// 2) иначе — пытаемся найти путь через общий граф по всем парам;
// 3) если не нашли — NaN.
func calcMaxStartUSDTAlways(
	t domain.Triangle,
	maxStart float64,
	quotes map[string]domain.Quote,
	fee float64,
	adj map[string][]convEdge,
) float64 {
	if maxStart <= 0 {
		return math.NaN()
	}

	startAsset := t.Legs[0].From
	if startAsset == "USDT" {
		return maxStart
	}

	// (1) пробуем по ногам этого треугольника до первого USDT
	cur := startAsset
	amt := maxStart

	for i := 0; i < 3; i++ {
		leg := t.Legs[i]
		if leg.From != cur {
			break
		}

		q, ok := quotes[leg.Symbol]
		if !ok || q.Bid <= 0 || q.Ask <= 0 {
			break
		}

		if leg.Dir > 0 {
			amt *= q.Bid
		} else {
			amt /= q.Ask
		}
		amt *= (1 - fee)

		cur = leg.To
		if cur == "USDT" {
			return amt
		}
	}

	// (2) общий граф
	if v, ok := convertByGraph(startAsset, "USDT", maxStart, quotes, adj, fee, 6); ok {
		return v
	}

	// (3) не нашли
	return math.NaN()
}

// -------------------- граф конвертации --------------------

type convEdge struct {
	to     string
	symbol string
	mode   int8 // +1: base->quote (sell base at bid), -1: quote->base (buy base at ask)
}

func buildAdjFromTriangles(tris []domain.Triangle) map[string][]convEdge {
	adj := make(map[string][]convEdge)
	seen := make(map[string]struct{})

	add := func(from string, e convEdge) {
		key := from + "|" + e.to + "|" + e.symbol + fmt.Sprintf("|%d", e.mode)
		if _, ok := seen[key]; ok {
			return
		}
		seen[key] = struct{}{}
		adj[from] = append(adj[from], e)
	}

	for _, t := range tris {
		for _, leg := range t.Legs {
			// base/quote для symbol
			var base, quote string
			if leg.Dir > 0 {
				base, quote = leg.From, leg.To
			} else {
				base, quote = leg.To, leg.From
			}

			add(base, convEdge{to: quote, symbol: leg.Symbol, mode: +1})
			add(quote, convEdge{to: base, symbol: leg.Symbol, mode: -1})
		}
	}
	return adj
}

func convertByGraph(
	fromAsset, toAsset string,
	amount float64,
	quotes map[string]domain.Quote,
	adj map[string][]convEdge,
	fee float64,
	maxHops int,
) (float64, bool) {
	if amount <= 0 {
		return 0, false
	}
	if fromAsset == toAsset {
		return amount, true
	}
	if maxHops <= 0 {
		return 0, false
	}

	// best[asset][hops] = лучший мультипликатор
	best := make(map[stateKey]float64)

	pq := &maxPQ{}
	heap.Init(pq)

	k0 := stateKey{asset: fromAsset, hops: 0}
	best[k0] = 1.0
	heap.Push(pq, pqState{asset: fromAsset, hops: 0, mult: 1.0})

	for pq.Len() > 0 {
		cur := heap.Pop(pq).(pqState)

		key := stateKey{asset: cur.asset, hops: cur.hops}
		if b, ok := best[key]; ok && cur.mult < b {
			continue
		}

		if cur.asset == toAsset {
			return amount * cur.mult, true
		}
		if cur.hops >= maxHops {
			continue
		}

		for _, e := range adj[cur.asset] {
			q, ok := quotes[e.symbol]
			if !ok || q.Bid <= 0 || q.Ask <= 0 {
				continue
			}

			var em float64
			if e.mode > 0 {
				em = q.Bid * (1 - fee)
			} else {
				em = (1.0 / q.Ask) * (1 - fee)
			}
			if em <= 0 || math.IsNaN(em) || math.IsInf(em, 0) {
				continue
			}

			nmult := cur.mult * em
			nhops := cur.hops + 1
			nkey := stateKey{asset: e.to, hops: nhops}

			if prev, ok := best[nkey]; !ok || nmult > prev {
				best[nkey] = nmult
				heap.Push(pq, pqState{asset: e.to, hops: nhops, mult: nmult})
			}
		}
	}

	return 0, false
}

type stateKey struct {
	asset string
	hops  int
}

type pqState struct {
	asset string
	hops  int
	mult  float64
}

type maxPQ []pqState

func (h maxPQ) Len() int            { return len(h) }
func (h maxPQ) Less(i, j int) bool  { return h[i].mult > h[j].mult } // max-heap
func (h maxPQ) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *maxPQ) Push(x interface{}) { *h = append(*h, x.(pqState)) }
func (h *maxPQ) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}







