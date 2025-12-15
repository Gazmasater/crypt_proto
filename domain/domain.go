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
	Symbol                   string
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
