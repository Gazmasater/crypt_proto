package calculator

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"

	"crypt_proto/internal/queue"
	"crypt_proto/pkg/models"
)

const (
	feeM            = 0.9992
	defaultExchange = "KuCoin"
)

type LegIndex struct {
	Key    string
	Symbol string
	IsBuy  bool
}

type Triangle struct {
	A, B, C string
	Legs    [3]LegIndex
}

type Calculator struct {
	mem      *queue.MemoryStore
	bySymbol map[string][]*Triangle
	fileLog  *log.Logger
	debugLog *log.Logger
}

func NewCalculator(mem *queue.MemoryStore, triangles []*Triangle) *Calculator {
	f, err := os.OpenFile("arb_opportunities.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("failed to open arb_opportunities.log: %v", err)
	}

	df, err := os.OpenFile("arb_debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("failed to open arb_debug.log: %v", err)
	}

	bySymbol := make(map[string][]*Triangle, 1024)

	for _, t := range triangles {
		if t == nil {
			continue
		}

		for _, leg := range t.Legs {
			if leg.Symbol == "" {
				continue
			}
			bySymbol[leg.Symbol] = append(bySymbol[leg.Symbol], t)
		}
	}

	log.Printf("[Calculator] indexed %d symbols, loaded %d triangles", len(bySymbol), len(triangles))

	return &Calculator{
		mem:      mem,
		bySymbol: bySymbol,
		fileLog:  log.New(f, "", log.LstdFlags),
		debugLog: log.New(df, "", log.LstdFlags),
	}
}

// Run считает только те треугольники, где изменился один из символов.
func (c *Calculator) Run(in <-chan *models.MarketData) {
	for md := range in {
		if md == nil || md.Symbol == "" {
			continue
		}

		// Оставляем push как был — MemoryStore обновляется асинхронно.
		// Но в расчёте ниже текущий md подставляется напрямую, чтобы не было гонки.
		c.mem.Push(md)

		tris := c.bySymbol[md.Symbol]
		if len(tris) == 0 {
			continue
		}

		for _, tri := range tris {
			c.calcTriangleWithCurrent(md, tri)
		}
	}
}

func (c *Calculator) calcTriangleWithCurrent(md *models.MarketData, tri *Triangle) {
	var q [3]queue.Quote

	for i, leg := range tri.Legs {
		if leg.Symbol == "" {
			c.debugf("skip %s->%s->%s: empty leg at index %d", tri.A, tri.B, tri.C, i)
			return
		}

		quote, ok := c.getQuoteWithOverride(md, leg.Symbol)
		if !ok {
			c.debugf("skip %s->%s->%s: no quote for leg=%d symbol=%s", tri.A, tri.B, tri.C, i+1, leg.Symbol)
			return
		}

		q[i] = quote
	}

	maxUSDT, ok := calcMaxStartUSDT(tri, q)
	if !ok {
		c.debugf(
			"skip %s->%s->%s: invalid liquidity/price | q1=%+v q2=%+v q3=%+v",
			tri.A, tri.B, tri.C, q[0], q[1], q[2],
		)
		return
	}

	if maxUSDT <= 0 {
		c.debugf("skip %s->%s->%s: maxUSDT <= 0", tri.A, tri.B, tri.C)
		return
	}

	amount := maxUSDT

	for i, leg := range tri.Legs {
		if leg.IsBuy {
			if q[i].Ask <= 0 {
				c.debugf("skip %s->%s->%s: leg=%d ask <= 0", tri.A, tri.B, tri.C, i+1)
				return
			}
			amount = amount / q[i].Ask * feeM
		} else {
			if q[i].Bid <= 0 {
				c.debugf("skip %s->%s->%s: leg=%d bid <= 0", tri.A, tri.B, tri.C, i+1)
				return
			}
			amount = amount * q[i].Bid * feeM
		}
	}

	profitUSDT := amount - maxUSDT
	profitPct := profitUSDT / maxUSDT

	msg := fmt.Sprintf(
		"[ARB] %s→%s→%s | %.4f%% | volume=%.2f USDT | profit=%.6f USDT",
		tri.A, tri.B, tri.C,
		profitPct*100, maxUSDT, profitUSDT,
	)

	log.Println(msg)
	c.fileLog.Println(msg)
}

func (c *Calculator) getQuoteWithOverride(md *models.MarketData, symbol string) (queue.Quote, bool) {
	if md != nil && md.Exchange == defaultExchange && md.Symbol == symbol {
		return queue.Quote{
			Bid:       md.Bid,
			Ask:       md.Ask,
			BidSize:   md.BidSize,
			AskSize:   md.AskSize,
			Timestamp: md.Timestamp,
		}, true
	}

	return c.mem.Get(defaultExchange, symbol)
}

// calcMaxStartUSDT оценивает максимально возможный стартовый объём в USDT,
// который можно реально провести через все 3 ноги по top-of-book.
func calcMaxStartUSDT(tri *Triangle, q [3]queue.Quote) (float64, bool) {
	var usdtLimits [3]float64

	// LEG 1
	if tri.Legs[0].IsBuy {
		if q[0].Ask <= 0 || q[0].AskSize <= 0 {
			return 0, false
		}
		usdtLimits[0] = q[0].Ask * q[0].AskSize
	} else {
		if q[0].Bid <= 0 || q[0].BidSize <= 0 {
			return 0, false
		}
		usdtLimits[0] = q[0].Bid * q[0].BidSize
	}

	// LEG 2
	if tri.Legs[1].IsBuy {
		// Чтобы купить на 2-й ноге, нужен ask и askSize на 2-й ноге.
		// Чтобы перевести лимит в USDT, используем 3-ю ногу как возврат в USDT.
		if q[1].Ask <= 0 || q[1].AskSize <= 0 || q[2].Bid <= 0 {
			return 0, false
		}
		usdtLimits[1] = q[1].Ask * q[1].AskSize * q[2].Bid
	} else {
		// На продаже 2-й ноги объём ограничен bidSize этой ноги,
		// а в USDT переводим через 3-ю ногу.
		if q[1].Bid <= 0 || q[1].BidSize <= 0 || q[2].Bid <= 0 {
			return 0, false
		}
		usdtLimits[1] = q[1].BidSize * q[2].Bid
	}

	// LEG 3
	if q[2].Bid <= 0 || q[2].BidSize <= 0 {
		return 0, false
	}
	usdtLimits[2] = q[2].Bid * q[2].BidSize

	maxUSDT := usdtLimits[0]
	if usdtLimits[1] < maxUSDT {
		maxUSDT = usdtLimits[1]
	}
	if usdtLimits[2] < maxUSDT {
		maxUSDT = usdtLimits[2]
	}

	return maxUSDT, true
}

func (c *Calculator) debugf(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	log.Println("[ARB-DEBUG]", msg)
	c.debugLog.Println(msg)
}

// ParseTrianglesFromCSV поддерживает два формата:
//  1. короткий: A,B,C,Leg1,Leg2,Leg3
//  2. расширенный:
//     A,B,C,Leg1,Step1,MinQty1,MinNotional1,Leg2,Step2,MinQty2,MinNotional2,Leg3,Step3,MinQty3,MinNotional3
func ParseTrianglesFromCSV(path string) ([]*Triangle, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rows, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return nil, err
	}

	var res []*Triangle

	for rowIdx, row := range rows[1:] {
		if len(row) < 6 {
			continue
		}

		t := &Triangle{
			A: strings.TrimSpace(row[0]),
			B: strings.TrimSpace(row[1]),
			C: strings.TrimSpace(row[2]),
		}

		legs, ok := extractLegColumns(row)
		if !ok {
			log.Printf("[ParseTrianglesFromCSV] skip row %d: unsupported csv format, cols=%d", rowIdx+2, len(row))
			continue
		}

		validLegs := 0
		for i, rawLeg := range legs {
			leg, ok := parseLeg(rawLeg)
			if !ok {
				log.Printf("[ParseTrianglesFromCSV] row %d: invalid leg %d: %q", rowIdx+2, i+1, rawLeg)
				continue
			}
			t.Legs[i] = leg
			validLegs++
		}

		if validLegs != 3 {
			log.Printf("[ParseTrianglesFromCSV] skip row %d: parsed only %d/3 legs", rowIdx+2, validLegs)
			continue
		}

		res = append(res, t)
	}

	return res, nil
}

func extractLegColumns(row []string) ([3]string, bool) {
	var legs [3]string

	switch {
	case len(row) >= 15:
		// A,B,C,Leg1,Step1,MinQty1,MinNotional1,Leg2,...,Leg3,...
		legs[0] = row[3]
		legs[1] = row[7]
		legs[2] = row[11]
		return legs, true

	case len(row) >= 6:
		// A,B,C,Leg1,Leg2,Leg3
		legs[0] = row[3]
		legs[1] = row[4]
		legs[2] = row[5]
		return legs, true

	default:
		return legs, false
	}
}

func parseLeg(raw string) (LegIndex, bool) {
	raw = strings.ToUpper(strings.TrimSpace(raw))
	parts := strings.Fields(raw)
	if len(parts) != 2 {
		return LegIndex{}, false
	}

	side := parts[0]
	if side != "BUY" && side != "SELL" {
		return LegIndex{}, false
	}

	pair := strings.Split(parts[1], "/")
	if len(pair) != 2 {
		return LegIndex{}, false
	}

	base := strings.TrimSpace(pair[0])
	quote := strings.TrimSpace(pair[1])
	if base == "" || quote == "" {
		return LegIndex{}, false
	}

	symbol := base + "-" + quote
	key := defaultExchange + "|" + symbol

	return LegIndex{
		Key:    key,
		Symbol: symbol,
		IsBuy:  side == "BUY",
	}, true
}
