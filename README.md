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





package domain

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
)

const BaseAsset = "USDT"

// ==============================
// Базовые структуры
// ==============================

type Leg struct {
	From   string
	To     string
	Symbol string
	Dir    int8
}

type Triangle struct {
	Legs [3]Leg
	Name string // USDT→A→B→USDT
}

type Quote struct {
	Bid, Ask, BidQty, AskQty float64
}

type Event struct {
	Symbol string
	Bid, Ask, BidQty, AskQty float64
}

type Pair struct {
	Base, Quote, Symbol string
}

// ==============================
// Построение треугольников
// ==============================

func buildTriangleFromPairs(p1, p2, p3 Pair) (Triangle, bool) {
	pairs := []Pair{p1, p2, p3}
	set := map[string]struct{}{
		p1.Base: {}, p1.Quote: {},
		p2.Base: {}, p2.Quote: {},
		p3.Base: {}, p3.Quote: {},
	}
	if len(set) != 3 {
		return Triangle{}, false
	}

	// Берём все валюты
	var assets []string
	for k := range set {
		assets = append(assets, k)
	}

	// если в тройке нет USDT — пропускаем
	if _, ok := set[BaseAsset]; !ok {
		return Triangle{}, false
	}

	// Две прочие валюты
	var others []string
	for _, a := range assets {
		if a != BaseAsset {
			others = append(others, a)
		}
	}
	if len(others) != 2 {
		return Triangle{}, false
	}
	A, B := others[0], others[1]

	// Пробуем два направления: USDT→A→B→USDT и USDT→B→A→USDT
	orders := [][]string{
		{BaseAsset, A, B, BaseAsset},
		{BaseAsset, B, A, BaseAsset},
	}

	for _, order := range orders {
		var legs [3]Leg
		okAll := true
		for i := 0; i < 3; i++ {
			from, to := order[i], order[i+1]
			var found bool
			for _, p := range pairs {
				switch {
				case p.Base == from && p.Quote == to:
					legs[i] = Leg{From: from, To: to, Symbol: p.Symbol, Dir: +1}
					found = true
				case p.Base == to && p.Quote == from:
					legs[i] = Leg{From: from, To: to, Symbol: p.Symbol, Dir: -1}
					found = true
				}
			}
			if !found {
				okAll = false
				break
			}
		}
		if okAll {
			name := fmt.Sprintf("%s→%s→%s→%s", order[0], order[1], order[2], order[3])
			return Triangle{Legs: legs, Name: name}, true
		}
	}
	return Triangle{}, false
}

// ==============================
// Загрузка из CSV
// ==============================

func LoadTriangles(path string) ([]Triangle, []string, map[string][]int, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, nil, err
	}
	defer f.Close()

	r := csv.NewReader(bufio.NewReader(f))
	r.TrimLeadingSpace = true
	r.Comma = ','

	var tris []Triangle
	symbolSet := make(map[string]struct{})

	for {
		rec, err := r.Read()
		if err != nil {
			break
		}
		var fields []string
		for _, v := range rec {
			v = strings.TrimSpace(v)
			if v != "" {
				fields = append(fields, v)
			}
		}
		if len(fields) != 6 || strings.HasPrefix(fields[0], "#") {
			continue
		}

		p1 := Pair{Base: fields[0], Quote: fields[1], Symbol: fields[0] + fields[1]}
		p2 := Pair{Base: fields[2], Quote: fields[3], Symbol: fields[2] + fields[3]}
		p3 := Pair{Base: fields[4], Quote: fields[5], Symbol: fields[4] + fields[5]}

		t, ok := buildTriangleFromPairs(p1, p2, p3)
		if !ok {
			continue
		}

		tris = append(tris, t)
		for _, leg := range t.Legs {
			symbolSet[leg.Symbol] = struct{}{}
		}
	}

	symbols := make([]string, 0, len(symbolSet))
	for s := range symbolSet {
		symbols = append(symbols, s)
	}

	index := make(map[string][]int)
	for i, t := range tris {
		for _, leg := range t.Legs {
			index[leg.Symbol] = append(index[leg.Symbol], i)
		}
	}

	log.Printf("треугольников (USDT→...→USDT): %d", len(tris))
	log.Printf("уникальных символов: %d", len(symbols))
	return tris, symbols, index, nil
}

// ==============================
// Расчёт доходности
// ==============================

func EvalTriangle(t Triangle, quotes map[string]Quote, fee float64) (float64, bool) {
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
	return amt - 1, true
}

// ==============================
// Диагностика лимитов
// ==============================

type MaxStartInfo struct {
	StartAsset    string
	MaxStart      float64
	BottleneckLeg int
	LimitIn       [3]float64
	KIn           [3]float64
	MaxStartByLeg [3]float64
}

func ComputeMaxStartTopOfBook(t Triangle, quotes map[string]Quote, fee float64) (MaxStartInfo, bool) {
	var info MaxStartInfo
	info.StartAsset = BaseAsset

	kIn := 1.0
	maxStart := 1e308
	info.BottleneckLeg = -1

	for i, leg := range t.Legs {
		q, ok := quotes[leg.Symbol]
		if !ok || q.Bid <= 0 || q.Ask <= 0 {
			return MaxStartInfo{}, false
		}
		info.KIn[i] = kIn

		var limitIn, ratio float64
		if leg.Dir > 0 {
			if q.BidQty <= 0 {
				return MaxStartInfo{}, false
			}
			limitIn = q.BidQty
			ratio = q.Bid
		} else {
			if q.AskQty <= 0 {
				return MaxStartInfo{}, false
			}
			limitIn = q.AskQty * q.Ask
			ratio = 1 / q.Ask
		}
		info.LimitIn[i] = limitIn
		if kIn <= 0 {
			return MaxStartInfo{}, false
		}
		maxByThis := limitIn / kIn
		info.MaxStartByLeg[i] = maxByThis
		if maxByThis < maxStart {
			maxStart = maxByThis
			info.BottleneckLeg = i
		}
		kIn *= ratio * (1 - fee)
		if kIn <= 0 {
			return MaxStartInfo{}, false
		}
	}

	if info.BottleneckLeg < 0 || maxStart <= 0 || maxStart > 1e307 {
		return MaxStartInfo{}, false
	}
	info.MaxStart = maxStart
	return info, true
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

type Consumer struct {
	FeePerLeg     float64
	MinProfit     float64
	MinStart      float64
	StartFraction float64
	writer        io.Writer
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

				prof, ok := domain.EvalTriangle(tr, quotes, c.FeePerLeg)
				if !ok || prof < c.MinProfit {
					continue
				}

				ms, okMS := domain.ComputeMaxStartTopOfBook(tr, quotes, c.FeePerLeg)
				if !okMS {
					continue
				}

				safeStart := ms.MaxStart * sf
				if c.MinStart > 0 {
					safeUSDT, okConv := convertToUSDT(safeStart, ms.StartAsset, quotes)
					if !okConv || safeUSDT < c.MinStart {
						continue
					}
				}

				if last, okLast := lastPrint[id]; okLast && now.Sub(last) < minPrintInterval {
					continue
				}
				lastPrint[id] = now

				msCopy := ms
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

	bneckSym := ""
	if ms != nil && ms.BottleneckLeg >= 0 && ms.BottleneckLeg < len(t.Legs) {
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
			q.BidQty, q.AskQty)
	}

	if ms != nil && c.FeePerLeg > 0 {
		execs, okExec := simulateTriangleExecution(t, quotes, ms.StartAsset, safeStart, c.FeePerLeg)
		if okExec {
			fmt.Fprintln(w, "  Legs execution with fees:")
			for i, e := range execs {
				fmt.Fprintf(w, "    leg %d: %s  %.6f %s → %.6f %s  fee=%.8f %s\n",
					i+1, e.Symbol, e.AmountIn, e.From, e.AmountOut, e.To, e.FeeAmount, e.FeeAsset)
			}
		}
	}
	fmt.Fprintln(w)
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

		var amountOut, feeAmount float64
		var feeAsset string

		switch {
		case curAsset == base:
			gross := curAmount * q.Bid
			feeAmount = gross * feePerLeg
			amountOut = gross - feeAmount
			feeAsset = quote
			curAsset, curAmount = quote, amountOut
		case curAsset == quote:
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
			AmountIn:  curAmount + feeAmount,
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



[{
	"resource": "/home/gaz358/myprog/crypt_proto/arb/arb.go",
	"owner": "go-staticcheck",
	"severity": 4,
	"message": "possible nil pointer dereference (SA5011)\n\tarb.go:181:5: this check suggests that the pointer can be nil",
	"source": "go-staticcheck",
	"startLineNumber": 139,
	"startColumn": 18,
	"endLineNumber": 139,
	"endColumn": 42,
	"origin": "extHost1"
}]
