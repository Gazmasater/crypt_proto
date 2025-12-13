package domain

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
)

type Leg struct {
	From   string
	To     string
	Symbol string
	Dir    int8 // +1: From->To = base->quote; -1: From->To = quote->base
}

type Triangle struct {
	Legs [3]Leg
	Name string // A→B→C→A
}

type Quote struct {
	Bid    float64
	Ask    float64
	BidQty float64
	AskQty float64
}

type Event struct {
	Symbol string
	Bid    float64
	Ask    float64
	BidQty float64
	AskQty float64
}

type Pair struct {
	Base   string
	Quote  string
	Symbol string
}

func buildTriangleFromPairs(p1, p2, p3 Pair) (Triangle, bool) {
	set := map[string]struct{}{
		p1.Base:  {},
		p1.Quote: {},
		p2.Base:  {},
		p2.Quote: {},
		p3.Base:  {},
		p3.Quote: {},
	}
	if len(set) != 3 {
		return Triangle{}, false
	}
	currs := make([]string, 0, 3)
	for c := range set {
		currs = append(currs, c)
	}

	type edge struct{ From, To string }

	pairs := []Pair{p1, p2, p3}
	perm3 := [][]int{
		{0, 1, 2},
		{0, 2, 1},
		{1, 0, 2},
		{1, 2, 0},
		{2, 0, 1},
		{2, 1, 0},
	}

	for _, order := range perm3 {
		c0, c1, c2 := currs[order[0]], currs[order[1]], currs[order[2]]
		edges := []edge{
			{From: c0, To: c1},
			{From: c1, To: c2},
			{From: c2, To: c0},
		}

		for _, pp := range perm3 {
			var legs [3]Leg
			okAll := true

			for i := 0; i < 3; i++ {
				e := edges[i]
				p := pairs[pp[i]]

				switch {
				case p.Base == e.From && p.Quote == e.To:
					legs[i] = Leg{From: e.From, To: e.To, Symbol: p.Symbol, Dir: +1}
				case p.Base == e.To && p.Quote == e.From:
					legs[i] = Leg{From: e.From, To: e.To, Symbol: p.Symbol, Dir: -1}
				default:
					okAll = false
				}
				if !okAll {
					break
				}
			}

			if okAll {
				name := fmt.Sprintf("%s→%s→%s→%s", edges[0].From, edges[1].From, edges[2].From, edges[0].From)
				return Triangle{Legs: legs, Name: name}, true
			}
		}
	}

	return Triangle{}, false
}

// LoadTriangles читает CSV, строит треугольники и индекс по символам.
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
			if err.Error() == "EOF" {
				break
			}
			return nil, nil, nil, err
		}

		var fields []string
		for _, v := range rec {
			v = strings.TrimSpace(v)
			if v != "" {
				fields = append(fields, v)
			}
		}
		if len(fields) == 0 {
			continue
		}
		if strings.HasPrefix(fields[0], "#") {
			continue
		}
		if len(fields) != 6 {
			log.Printf("skip line (need 6 fields): %v", fields)
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

	log.Printf("треугольников всего: %d", len(tris))
	log.Printf("символов в индексе треугольников: %d", len(symbols))

	return tris, symbols, index, nil
}

// EvalTriangle считает доходность треугольника.
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

	return amt - 1.0, true
}

// MaxStartInfo - диагностическая информация по максимальному стартовому объёму.
// Все значения рассчитаны по top-of-book (best bid/ask) без учёта округлений stepSize/minQty.
type MaxStartInfo struct {
	StartAsset    string
	MaxStart      float64    // в StartAsset
	BottleneckLeg int        // индекс ноги [0..2]
	LimitIn       [3]float64 // лимит на ВХОД каждой ноги (в единицах входного актива ноги)
	KIn           [3]float64 // сколько входного актива ноги получается из 1 StartAsset
	MaxStartByLeg [3]float64 // LimitIn / KIn
}

// ComputeMaxStartTopOfBook возвращает максимальный стартовый объём, который можно протащить через треугольник,
// не выходя за best bid/ask qty на каждой ноге. Комиссия учитывается как удержание из результата каждой ноги.
func ComputeMaxStartTopOfBook(t Triangle, quotes map[string]Quote, fee float64) (MaxStartInfo, bool) {
	var info MaxStartInfo
	info.StartAsset = t.Legs[0].From

	// kIn - сколько входной валюты текущей ноги получится из 1 единицы стартовой валюты.
	kIn := 1.0
	info.MaxStart = 0
	info.BottleneckLeg = -1

	// Инициализируем maxStart как +Inf, чтобы взять минимум по ногам.
	maxStart := 1e308

	for i, leg := range t.Legs {
		q, ok := quotes[leg.Symbol]
		if !ok || q.Bid <= 0 || q.Ask <= 0 {
			return MaxStartInfo{}, false
		}

		info.KIn[i] = kIn

		var limitIn float64
		var ratio float64 // out/in без комиссии

		if leg.Dir > 0 {
			// SELL base -> quote по bid, ограничение по количеству base на bid.
			if q.BidQty <= 0 {
				return MaxStartInfo{}, false
			}
			limitIn = q.BidQty
			ratio = q.Bid
		} else {
			// BUY base <- quote по ask, ограничение по объёму quote, который можно потратить: askQty*ask.
			if q.AskQty <= 0 {
				return MaxStartInfo{}, false
			}
			limitIn = q.AskQty * q.Ask
			ratio = 1.0 / q.Ask
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

		// переход к следующей ноге
		kIn = kIn * ratio * (1 - fee)
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
