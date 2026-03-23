package calculator

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
	"strings"

	"crypt_proto/internal/queue"
	"crypt_proto/pkg/models"
)

const (
	DefaultFeeRate      = 0.0008 // 0.08%
	DefaultSafetyFactor = 0.7
	DefaultMinProfitPct = 0.001 // 0.1%
	DefaultMinStartUSDT = 50.0
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

type MarketRules struct {
	StepSize    float64
	MinQty      float64
	MinNotional float64
}

type LegResult struct {
	Symbol       string
	IsBuy        bool
	In           float64
	Out          float64
	Price        float64
	BookQty      float64
	FeePaid      float64
	Notional     float64
	Valid        bool
	InvalidCause string
}

type TriangleResult struct {
	StartUSDT  float64
	FinalUSDT  float64
	ProfitUSDT float64
	ProfitPct  float64
	Valid      bool
	Legs       [3]LegResult
}

type Calculator struct {
	mem          *queue.MemoryStore
	bySymbol     map[string][]*Triangle
	fileLog      *log.Logger
	rules        map[string]MarketRules
	feeRate      float64
	safetyFactor float64
	minProfitPct float64
	minStartUSDT float64
}

func NewCalculator(mem *queue.MemoryStore, triangles []*Triangle, rules map[string]MarketRules) *Calculator {
	f, err := os.OpenFile("arb_opportunities.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("failed to open log: %v", err)
	}

	bySymbol := make(map[string][]*Triangle, 1024)
	for _, t := range triangles {
		for _, leg := range t.Legs {
			bySymbol[leg.Symbol] = append(bySymbol[leg.Symbol], t)
		}
	}

	log.Printf("[Calculator] indexed %d symbols\n", len(bySymbol))

	return &Calculator{
		mem:          mem,
		bySymbol:     bySymbol,
		fileLog:      log.New(f, "", log.LstdFlags),
		rules:        rules,
		feeRate:      DefaultFeeRate,
		safetyFactor: DefaultSafetyFactor,
		minProfitPct: DefaultMinProfitPct,
		minStartUSDT: DefaultMinStartUSDT,
	}
}

// Необязательные сеттеры
func (c *Calculator) SetFeeRate(v float64) {
	if v >= 0 && v < 1 {
		c.feeRate = v
	}
}

func (c *Calculator) SetSafetyFactor(v float64) {
	if v > 0 && v <= 1 {
		c.safetyFactor = v
	}
}

func (c *Calculator) SetMinProfitPct(v float64) {
	if v >= 0 {
		c.minProfitPct = v
	}
}

func (c *Calculator) SetMinStartUSDT(v float64) {
	if v >= 0 {
		c.minStartUSDT = v
	}
}

// Run — считаем только нужные треугольники
func (c *Calculator) Run(in <-chan *models.MarketData) {
	for md := range in {
		c.mem.Push(md)

		tris := c.bySymbol[md.Symbol]
		if len(tris) == 0 {
			continue
		}

		for _, tri := range tris {
			c.calcTriangle(tri)
		}
	}
}

func (c *Calculator) calcTriangle(tri *Triangle) {
	var q [3]queue.Quote

	for i, leg := range tri.Legs {
		quote, ok := c.mem.Get("KuCoin", leg.Symbol)
		if !ok {
			return
		}
		q[i] = quote
	}

	maxStart, ok := c.computeMaxStartUSDT(tri, q)
	if !ok || maxStart <= 0 {
		return
	}

	safeStart := maxStart * c.safetyFactor
	if safeStart <= 0 {
		return
	}

	// Округляем старт вниз по шагу первой ноги, если это BUY.
	// safeStart у нас в USDT. Для BUY 1-й ноги шаг количества лежит в base-активе,
	// поэтому округление идёт через симуляцию.
	res := c.evaluateTriangleFromStart(tri, q, safeStart)
	if !res.Valid {
		return
	}

	if res.StartUSDT < c.minStartUSDT {
		return
	}

	if res.ProfitPct <= c.minProfitPct {
		return
	}

	msg := fmt.Sprintf(
		"[ARB] %s→%s→%s | %.4f%% | start=%.2f USDT | final=%.4f USDT | profit=%.4f USDT",
		tri.A, tri.B, tri.C,
		res.ProfitPct*100,
		res.StartUSDT,
		res.FinalUSDT,
		res.ProfitUSDT,
	)

	log.Println(msg)
	c.fileLog.Println(msg)
}

func (c *Calculator) computeMaxStartUSDT(tri *Triangle, q [3]queue.Quote) (float64, bool) {
	var usdtLimits [3]float64

	// LEG 1: всегда переводим лимит top-of-book в эквивалент USDT старта
	// Если BUY A/B: максимум quote = ask * askSize
	// Если SELL A/B: входной актив = A, но старт у нас USDT -> здесь предполагается,
	// что leg1 уже начинается из tri.A и в твоих треугольниках старт реально USDT.
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

	// LEG 2: приводим ограничение второй ноги к старту в USDT
	// Важно: здесь top-of-book модель, поэтому делаем приближение через leg1/leg2/leg3 цены.
	limit2, ok := c.computeLeg2StartLimitUSDT(tri, q)
	if !ok || limit2 <= 0 {
		return 0, false
	}
	usdtLimits[1] = limit2

	// LEG 3: сколько максимум может принять 3-я нога по своему bid/ask, тоже в USDT
	limit3, ok := c.computeLeg3StartLimitUSDT(tri, q)
	if !ok || limit3 <= 0 {
		return 0, false
	}
	usdtLimits[2] = limit3

	maxUSDT := usdtLimits[0]
	if usdtLimits[1] < maxUSDT {
		maxUSDT = usdtLimits[1]
	}
	if usdtLimits[2] < maxUSDT {
		maxUSDT = usdtLimits[2]
	}

	if maxUSDT <= 0 || math.IsNaN(maxUSDT) || math.IsInf(maxUSDT, 0) {
		return 0, false
	}

	return maxUSDT, true
}

// computeLeg2StartLimitUSDT — максимум стартового USDT, который не переполнит 2-ю ногу.
func (c *Calculator) computeLeg2StartLimitUSDT(tri *Triangle, q [3]queue.Quote) (float64, bool) {
	// После 1-й ноги получаем amount1.
	// Надо ограничить старт так, чтобы вход во 2-ю ногу поместился в top size второй ноги.

	if tri.Legs[0].IsBuy {
		if q[0].Ask <= 0 {
			return 0, false
		}
	} else {
		if q[0].Bid <= 0 {
			return 0, false
		}
	}

	// Максимально допустимый вход во 2-ю ногу в единицах её входного актива:
	var maxInputLeg2 float64
	if tri.Legs[1].IsBuy {
		// BUY: вход в quote, top capacity = ask * askSize
		if q[1].Ask <= 0 || q[1].AskSize <= 0 {
			return 0, false
		}
		maxInputLeg2 = q[1].Ask * q[1].AskSize
	} else {
		// SELL: вход в base, top capacity = bidSize
		if q[1].Bid <= 0 || q[1].BidSize <= 0 {
			return 0, false
		}
		maxInputLeg2 = q[1].BidSize
	}

	// Переводим допустимый старт так, чтобы output leg1 <= maxInputLeg2
	// leg1:
	// BUY  => out = start / ask * fee
	// SELL => out = start * bid * fee
	if tri.Legs[0].IsBuy {
		// start / ask * fee <= maxInputLeg2
		denom := c.feeRateMultiplier() / q[0].Ask
		if denom <= 0 {
			return 0, false
		}
		return maxInputLeg2 / denom, true
	}

	// start * bid * fee <= maxInputLeg2
	denom := q[0].Bid * c.feeRateMultiplier()
	if denom <= 0 {
		return 0, false
	}
	return maxInputLeg2 / denom, true
}

// computeLeg3StartLimitUSDT — максимум стартового USDT, который не переполнит 3-ю ногу.
func (c *Calculator) computeLeg3StartLimitUSDT(tri *Triangle, q [3]queue.Quote) (float64, bool) {
	// Ограничение по 3-й ноге: output leg2 должен поместиться во вход 3-й ноги.

	// Входная capacity 3-й ноги
	var maxInputLeg3 float64
	if tri.Legs[2].IsBuy {
		if q[2].Ask <= 0 || q[2].AskSize <= 0 {
			return 0, false
		}
		maxInputLeg3 = q[2].Ask * q[2].AskSize
	} else {
		if q[2].Bid <= 0 || q[2].BidSize <= 0 {
			return 0, false
		}
		maxInputLeg3 = q[2].BidSize
	}

	// Сначала выводим максимально допустимый output leg1,
	// чтобы leg2 output <= maxInputLeg3.
	var maxOutLeg1 float64
	if tri.Legs[1].IsBuy {
		// leg2 BUY: out2 = in1 / ask2 * fee <= maxInputLeg3
		if q[1].Ask <= 0 {
			return 0, false
		}
		maxOutLeg1 = maxInputLeg3 * q[1].Ask / c.feeRateMultiplier()
	} else {
		// leg2 SELL: out2 = in1 * bid2 * fee <= maxInputLeg3
		if q[1].Bid <= 0 {
			return 0, false
		}
		denom := q[1].Bid * c.feeRateMultiplier()
		if denom <= 0 {
			return 0, false
		}
		maxOutLeg1 = maxInputLeg3 / denom
	}

	// Теперь переводим maxOutLeg1 обратно в старт USDT через leg1
	if tri.Legs[0].IsBuy {
		// out1 = start / ask1 * fee <= maxOutLeg1
		if q[0].Ask <= 0 {
			return 0, false
		}
		return maxOutLeg1 * q[0].Ask / c.feeRateMultiplier(), true
	}

	// out1 = start * bid1 * fee <= maxOutLeg1
	denom := q[0].Bid * c.feeRateMultiplier()
	if denom <= 0 {
		return 0, false
	}
	return maxOutLeg1 / denom, true
}

func (c *Calculator) evaluateTriangleFromStart(tri *Triangle, q [3]queue.Quote, startUSDT float64) TriangleResult {
	var res TriangleResult
	res.StartUSDT = startUSDT
	res.Valid = false

	amount := startUSDT

	for i := 0; i < 3; i++ {
		lr := c.simulateLeg(tri.Legs[i], q[i], amount)
		res.Legs[i] = lr
		if !lr.Valid {
			return res
		}
		amount = lr.Out
	}

	res.FinalUSDT = amount
	res.ProfitUSDT = res.FinalUSDT - res.StartUSDT

	if res.StartUSDT > 0 {
		res.ProfitPct = res.ProfitUSDT / res.StartUSDT
	}

	res.Valid = true
	return res
}

func (c *Calculator) simulateLeg(leg LegIndex, q queue.Quote, amountIn float64) LegResult {
	rules := c.rules[leg.Symbol]
	feeMul := c.feeRateMultiplier()

	res := LegResult{
		Symbol: leg.Symbol,
		IsBuy:  leg.IsBuy,
		In:     amountIn,
		Valid:  false,
	}

	if amountIn <= 0 {
		res.InvalidCause = "non-positive input"
		return res
	}

	if leg.IsBuy {
		// BUY base with quote
		if q.Ask <= 0 || q.AskSize <= 0 {
			res.InvalidCause = "invalid ask/askSize"
			return res
		}

		maxQuoteAtTop := q.Ask * q.AskSize
		if amountIn > maxQuoteAtTop {
			res.InvalidCause = "top ask liquidity insufficient"
			return res
		}

		baseQty := amountIn / q.Ask
		baseQtyRounded := floorToStep(baseQty, rules.StepSize)
		if baseQtyRounded <= 0 {
			res.InvalidCause = "rounded qty <= 0"
			return res
		}

		if rules.MinQty > 0 && baseQtyRounded < rules.MinQty {
			res.InvalidCause = "qty < minQty"
			return res
		}

		notional := baseQtyRounded * q.Ask
		if rules.MinNotional > 0 && notional < rules.MinNotional {
			res.InvalidCause = "notional < minNotional"
			return res
		}

		if baseQtyRounded > q.AskSize {
			res.InvalidCause = "rounded qty > askSize"
			return res
		}

		out := baseQtyRounded * feeMul
		feePaid := baseQtyRounded - out

		res.Price = q.Ask
		res.BookQty = q.AskSize
		res.Notional = notional
		res.Out = out
		res.FeePaid = feePaid
		res.Valid = true
		return res
	}

	// SELL base for quote
	if q.Bid <= 0 || q.BidSize <= 0 {
		res.InvalidCause = "invalid bid/bidSize"
		return res
	}

	baseQtyRounded := floorToStep(amountIn, rules.StepSize)
	if baseQtyRounded <= 0 {
		res.InvalidCause = "rounded qty <= 0"
		return res
	}

	if rules.MinQty > 0 && baseQtyRounded < rules.MinQty {
		res.InvalidCause = "qty < minQty"
		return res
	}

	if baseQtyRounded > q.BidSize {
		res.InvalidCause = "top bid liquidity insufficient"
		return res
	}

	notional := baseQtyRounded * q.Bid
	if rules.MinNotional > 0 && notional < rules.MinNotional {
		res.InvalidCause = "notional < minNotional"
		return res
	}

	out := notional * feeMul
	feePaid := notional - out

	res.Price = q.Bid
	res.BookQty = q.BidSize
	res.Notional = notional
	res.Out = out
	res.FeePaid = feePaid
	res.Valid = true
	return res
}

func (c *Calculator) feeRateMultiplier() float64 {
	return 1.0 - c.feeRate
}

func floorToStep(v, step float64) float64 {
	if v <= 0 {
		return 0
	}
	if step <= 0 {
		return v
	}
	return math.Floor(v/step) * step
}

// CSV без изменений логики, но сразу сохраняем Symbol
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

	for _, row := range rows[1:] {
		if len(row) < 6 {
			continue
		}

		t := &Triangle{
			A: strings.TrimSpace(row[0]),
			B: strings.TrimSpace(row[1]),
			C: strings.TrimSpace(row[2]),
		}

		valid := true

		for i, leg := range []string{row[3], row[4], row[5]} {
			leg = strings.ToUpper(strings.TrimSpace(leg))
			parts := strings.Fields(leg)
			if len(parts) != 2 {
				valid = false
				break
			}

			isBuy := parts[0] == "BUY"
			pair := strings.Split(parts[1], "/")
			if len(pair) != 2 {
				valid = false
				break
			}

			symbol := pair[0] + "-" + pair[1]
			key := "KuCoin|" + symbol

			t.Legs[i] = LegIndex{
				Key:    key,
				Symbol: symbol,
				IsBuy:  isBuy,
			}
		}

		if valid {
			res = append(res, t)
		}
	}

	return res, nil
}
