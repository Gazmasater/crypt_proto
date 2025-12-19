package domain

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
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
	Dir    int8 // +1 SELL, -1 BUY

	// Опционально (если есть в CSV)
	QtyDP      int // знаков qty
	QuoteDPMkt int // знаков quoteQty для MARKET
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

// ==============================
// Загрузка из CSV (routes формат)
// ==============================

// LoadTriangles читает triangles_usdt_routes_market.csv (и совместимые),
// где есть колонки start/mid1/mid2/end и leg{1..3}_*.
// ВАЖНО: пока используем BUY только за USDT -> фильтруем треугольники,
// где любая BUY-нога имеет From != USDT.
func LoadTriangles(path string) ([]Triangle, []string, map[string][]int, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, nil, err
	}
	defer f.Close()

	r := csv.NewReader(bufio.NewReader(f))
	r.TrimLeadingSpace = true
	r.Comma = ','

	// header
	header, err := r.Read()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("read header: %w", err)
	}
	h := make(map[string]int, len(header))
	for i, col := range header {
		h[strings.ToLower(strings.TrimSpace(col))] = i
	}

	need := []string{
		"start", "mid1", "mid2", "end",
		"leg1_symbol", "leg1_action", "leg1_from", "leg1_to",
		"leg2_symbol", "leg2_action", "leg2_from", "leg2_to",
		"leg3_symbol", "leg3_action", "leg3_from", "leg3_to",
	}
	for _, k := range need {
		if _, ok := h[k]; !ok {
			return nil, nil, nil, fmt.Errorf("CSV missing required column: %s", k)
		}
	}

	// optional dp columns
	opt := func(name string) int {
		if i, ok := h[name]; ok {
			return i
		}
		return -1
	}
	iL1Qty := opt("leg1_qty_dp")
	iL2Qty := opt("leg2_qty_dp")
	iL3Qty := opt("leg3_qty_dp")
	iL1Qm := opt("leg1_quote_dp_market")
	iL2Qm := opt("leg2_quote_dp_market")
	iL3Qm := opt("leg3_quote_dp_market")

	var tris []Triangle
	symbolSet := make(map[string]struct{})

	rowNum := 1 // header already read

	// stats
	var total, kept int
	var droppedBadLeg int
	var droppedNotUsdtCycle int
	var droppedBuyNotUSDT int

	for {
		rec, err := r.Read()
		if err != nil {
			break
		}
		rowNum++
		total++

		get := func(key string) string {
			idx := h[key]
			if idx < 0 || idx >= len(rec) {
				return ""
			}
			return strings.TrimSpace(rec[idx])
		}

		start := strings.ToUpper(get("start"))
		mid1 := strings.ToUpper(get("mid1"))
		mid2 := strings.ToUpper(get("mid2"))
		end := strings.ToUpper(get("end"))

		if start == "" || mid1 == "" || mid2 == "" || end == "" {
			continue
		}
		// строго USDT→...→USDT
		if start != BaseAsset || end != BaseAsset {
			droppedNotUsdtCycle++
			continue
		}

		legs := [3]Leg{
			parseLeg(rec, h, 1, iL1Qty, iL1Qm),
			parseLeg(rec, h, 2, iL2Qty, iL2Qm),
			parseLeg(rec, h, 3, iL3Qty, iL3Qm),
		}

		ok := true
		for i := 0; i < 3; i++ {
			if legs[i].Symbol == "" || legs[i].From == "" || legs[i].To == "" || legs[i].Dir == 0 {
				ok = false
				break
			}
		}
		if !ok {
			droppedBadLeg++
			continue
		}

		// связность по активам
		if legs[0].From != start || legs[0].To != mid1 ||
			legs[1].From != mid1 || legs[1].To != mid2 ||
			legs[2].From != mid2 || legs[2].To != end {
			droppedBadLeg++
			continue
		}

		// КРИТИЧНО: пока BUY только за USDT
		buyOk := true
		for i := 0; i < 3; i++ {
			if legs[i].Dir < 0 && legs[i].From != BaseAsset {
				buyOk = false
				break
			}
		}
		if !buyOk {
			droppedBuyNotUSDT++
			continue
		}

		name := fmt.Sprintf("%s→%s→%s→%s", start, mid1, mid2, end)
		t := Triangle{Legs: legs, Name: name}
		tris = append(tris, t)
		kept++

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

	log.Printf("triangles loaded: total=%d kept=%d droppedNotUsdtCycle=%d droppedBadLeg=%d droppedBuyNotUSDT=%d file=%s",
		total, kept, droppedNotUsdtCycle, droppedBadLeg, droppedBuyNotUSDT, path,
	)
	log.Printf("unique symbols: %d", len(symbols))
	return tris, symbols, index, nil
}

func parseLeg(rec []string, h map[string]int, n int, idxQtyDP, idxQuoteDPMkt int) Leg {
	idx := func(s string) int {
		return h[strings.ToLower(fmt.Sprintf("leg%d_%s", n, s))]
	}
	get := func(i int) string {
		if i < 0 || i >= len(rec) {
			return ""
		}
		return strings.TrimSpace(rec[i])
	}

	symbol := strings.ToUpper(get(idx("symbol")))
	action := strings.ToUpper(get(idx("action")))
	from := strings.ToUpper(get(idx("from")))
	to := strings.ToUpper(get(idx("to")))

	var dir int8
	switch action {
	case "BUY":
		dir = -1
	case "SELL":
		dir = +1
	default:
		dir = 0
	}

	leg := Leg{
		From:       from,
		To:         to,
		Symbol:     symbol,
		Dir:        dir,
		QtyDP:      -1,
		QuoteDPMkt: -1,
	}
	if idxQtyDP >= 0 && idxQtyDP < len(rec) {
		leg.QtyDP = parseIntSafe(get(idxQtyDP), -1)
	}
	if idxQuoteDPMkt >= 0 && idxQuoteDPMkt < len(rec) {
		leg.QuoteDPMkt = parseIntSafe(get(idxQuoteDPMkt), -1)
	}
	return leg
}

func parseIntSafe(s string, def int) int {
	s = strings.TrimSpace(s)
	if s == "" {
		return def
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return v
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
// Диагностика лимитов (Top-of-book)
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
