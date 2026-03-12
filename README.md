Если оставить только нужное:

p99 execution latency
Micro-volatility (100 мс)
Fill ratio
Capture rate
Inventory drift




Название API
9623527002

696935c42a6dcd00013273f2
b348b686-55ff-4290-897b-02d55f815f65




apikey = "4333ed4b-cd83-49f5-97d1-c399e2349748"
secretkey = "E3848531135EDB4CCFDA0F1BC14CD274"
IP = ""
Название API-ключа = "Arb"
Доступы = "Чтение"



sudo systemctl mask sleep.target suspend.target hibernate.target hybrid-sleep.target



wbs-api.mexc.com/ws 


[https://edis-global.vercel.app/ru/vps-hosting/singapore-singapore
](https://sg.edisglobal.com/)



# BTC/USDT microstructure model (>60% precision on selected signals)

Ниже — полный каркас системы для:

* WebSocket order book + trades
* признаки: OFI, queue depletion, queue replenishment, microprice drift, imbalance
* regime filter
* score model
* вход/выход
* paper/live ready architecture

---

## Архитектура

```text
cmd/
  trader/
    main.go

internal/
  config/
    config.go
  types/
    market.go
    signal.go
  ring/
    ring.go
  features/
    micro.go
    regime.go
    scorer.go
  strategy/
    strategy.go
  paper/
    engine.go
  storage/
    logger.go
```

### Поток данных

```text
WS depth + trades
    -> local market state
    -> feature engine
    -> regime filter
    -> score model
    -> enter/exit logic
    -> paper engine / live execution
```

---

## Файл: `internal/types/market.go`

```go
package types

import "time"

type Side int

const (
	SideUnknown Side = iota
	SideBuy
	SideSell
)

type Level struct {
	Price float64
	Qty   float64
}

type BookSnapshot struct {
	Symbol    string
	Bids      []Level // sorted desc
	Asks      []Level // sorted asc
	Timestamp time.Time
}

type TradeTick struct {
	Symbol    string
	Price     float64
	Qty       float64
	Side      Side // aggressive side
	Timestamp time.Time
}
```

---

## Файл: `internal/types/signal.go`

```go
package types

import "time"

type FeatureState struct {
	Symbol string

	Mid    float64
	Micro  float64
	Spread float64
	Tick   float64

	ImbalanceTop3 float64
	ImbalanceTop5 float64

	QDBid float64
	QDAsk float64
	QRBid float64
	QRAsk float64

	OFINorm        float64
	MicroDriftNorm float64
	Return1s       float64
	Vol10s         float64
	Entropy10s     float64

	UpdatedAt time.Time
}

type RegimeState struct {
	Tradable bool
	Reason   string
}

type Signal struct {
	Symbol string

	LongScore  float64
	ShortScore float64

	Long      bool
	Short     bool
	NoTrade   bool
	Reason    string
	CreatedAt time.Time
}

type Position struct {
	Side      Side
	Entry     float64
	EntryTime time.Time
	Size      float64
}
```

---

## Файл: `internal/ring/ring.go`

```go
package ring

import "math"

type FloatRing struct {
	vals []float64
	idx  int
	full bool
}

func NewFloatRing(size int) *FloatRing {
	if size <= 0 {
		size = 1
	}
	return &FloatRing{vals: make([]float64, size)}
}

func (r *FloatRing) Add(v float64) {
	r.vals[r.idx] = v
	r.idx = (r.idx + 1) % len(r.vals)
	if r.idx == 0 {
		r.full = true
	}
}

func (r *FloatRing) Len() int {
	if r.full {
		return len(r.vals)
	}
	return r.idx
}

func (r *FloatRing) Values() []float64 {
	if !r.full {
		out := make([]float64, r.idx)
		copy(out, r.vals[:r.idx])
		return out
	}
	out := make([]float64, len(r.vals))
	copy(out, r.vals[r.idx:])
	copy(out[len(r.vals)-r.idx:], r.vals[:r.idx])
	return out
}

func (r *FloatRing) Last() (float64, bool) {
	if !r.full && r.idx == 0 {
		return 0, false
	}
	p := r.idx - 1
	if p < 0 {
		p = len(r.vals) - 1
	}
	return r.vals[p], true
}

func Mean(xs []float64) float64 {
	if len(xs) == 0 {
		return 0
	}
	s := 0.0
	for _, x := range xs {
		s += x
	}
	return s / float64(len(xs))
}

func Std(xs []float64) float64 {
	if len(xs) < 2 {
		return 0
	}
	m := Mean(xs)
	ss := 0.0
	for _, x := range xs {
		d := x - m
		ss += d * d
	}
	return math.Sqrt(ss / float64(len(xs)))
}

func SumAbs(xs []float64) float64 {
	s := 0.0
	for _, x := range xs {
		if x < 0 {
			s -= x
		} else {
			s += x
		}
	}
	return s
}
```

---

## Файл: `internal/features/micro.go`

```go
package features

import (
	"math"
	"time"

	"your_project/internal/ring"
	"your_project/internal/types"
)

type Engine struct {
	Symbol string

	TopN int

	// current book
	bestBid float64
	bestAsk float64
	bid1Qty float64
	ask1Qty float64
	mid     float64
	micro   float64
	spread  float64
	tick    float64

	prevBestBid float64
	prevBestAsk float64
	prevBid1Qty float64
	prevAsk1Qty float64

	prevBidTop3 float64
	prevAskTop3 float64
	currBidTop3 float64
	currAskTop3 float64

	prevBidTop5 float64
	prevAskTop5 float64
	currBidTop5 float64
	currAskTop5 float64

	ofiWindow    *ring.FloatRing
	driftWindow  *ring.FloatRing
	returnWindow *ring.FloatRing

	initialized bool
}

func NewEngine(symbol string) *Engine {
	return &Engine{
		Symbol:       symbol,
		TopN:         5,
		ofiWindow:    ring.NewFloatRing(64),  // ~1s if high frequency
		driftWindow:  ring.NewFloatRing(64),
		returnWindow: ring.NewFloatRing(256), // longer vol/entropy context
	}
}

func (e *Engine) UpdateBook(book types.BookSnapshot) {
	if len(book.Bids) == 0 || len(book.Asks) == 0 {
		return
	}

	bestBid := book.Bids[0].Price
	bestAsk := book.Asks[0].Price
	bid1Qty := book.Bids[0].Qty
	ask1Qty := book.Asks[0].Qty
	mid := 0.5 * (bestBid + bestAsk)
	spread := bestAsk - bestBid
	micro := computeMicro(bestBid, bestAsk, bid1Qty, ask1Qty)

	bidTop3 := sumTopNQty(book.Bids, 3)
	askTop3 := sumTopNQty(book.Asks, 3)
	bidTop5 := sumTopNQty(book.Bids, 5)
	askTop5 := sumTopNQty(book.Asks, 5)

	if !e.initialized {
		e.bestBid, e.bestAsk = bestBid, bestAsk
		e.bid1Qty, e.ask1Qty = bid1Qty, ask1Qty
		e.mid, e.micro, e.spread = mid, micro, spread
		e.tick = estimateTick(book.Bids, book.Asks)

		e.prevBestBid, e.prevBestAsk = bestBid, bestAsk
		e.prevBid1Qty, e.prevAsk1Qty = bid1Qty, ask1Qty

		e.prevBidTop3, e.prevAskTop3 = bidTop3, askTop3
		e.currBidTop3, e.currAskTop3 = bidTop3, askTop3
		e.prevBidTop5, e.prevAskTop5 = bidTop5, askTop5
		e.currBidTop5, e.currAskTop5 = bidTop5, askTop5

		e.ofiWindow.Add(0)
		e.driftWindow.Add(0)
		e.returnWindow.Add(0)
		e.initialized = true
		return
	}

	prevMid := e.mid

	e.prevBidTop3, e.prevAskTop3 = e.currBidTop3, e.currAskTop3
	e.currBidTop3, e.currAskTop3 = bidTop3, askTop3
	
	e.prevBidTop5, e.prevAskTop5 = e.currBidTop5, e.currAskTop5
	e.currBidTop5, e.currAskTop5 = bidTop5, askTop5

	ofi := calcTopLevelOFI(
		e.prevBestBid, e.prevBid1Qty,
		e.prevBestAsk, e.prevAsk1Qty,
		bestBid, bid1Qty,
		bestAsk, ask1Qty,
	)
	e.ofiWindow.Add(ofi)
	e.driftWindow.Add(micro - e.micro)

	ret := 0.0
	if prevMid > 0 {
		ret = (mid - prevMid) / prevMid
	}
	e.returnWindow.Add(ret)

	e.bestBid, e.bestAsk = bestBid, bestAsk
	e.bid1Qty, e.ask1Qty = bid1Qty, ask1Qty
	e.mid, e.micro, e.spread = mid, micro, spread
	if e.tick == 0 {
		e.tick = estimateTick(book.Bids, book.Asks)
	}

	e.prevBestBid, e.prevBestAsk = bestBid, bestAsk
	e.prevBid1Qty, e.prevAsk1Qty = bid1Qty, ask1Qty
}

func (e *Engine) UpdateTrade(tr types.TradeTick) {
	if !e.initialized {
		return
	}
	switch tr.Side {
	case types.SideBuy:
		e.ofiWindow.Add(tr.Qty)
	case types.SideSell:
		e.ofiWindow.Add(-tr.Qty)
	}
}

func (e *Engine) Snapshot() types.FeatureState {
	imb3 := computeImbalance(e.currBidTop3, e.currAskTop3)
	imb5 := computeImbalance(e.currBidTop5, e.currAskTop5)

	qdBid := computeDepletion(e.prevBidTop5, e.currBidTop5)
	qdAsk := computeDepletion(e.prevAskTop5, e.currAskTop5)
	qrBid := computeReplenishment(e.prevBidTop5, e.currBidTop5)
	qrAsk := computeReplenishment(e.prevAskTop5, e.currAskTop5)

	ofiVals := e.ofiWindow.Values()
	ofiNorm := 0.0
	if sabs := ring.SumAbs(ofiVals); sabs > 0 {
		sum := 0.0
		for _, v := range ofiVals {
			sum += v
		}
		ofiNorm = sum / sabs
	}

	driftVals := e.driftWindow.Values()
	microDriftNorm := 0.0
	if len(driftVals) > 0 {
		last, _ := e.driftWindow.Last()
		den := e.mid
		if e.tick > 0 {
			den = math.Max(e.tick, e.mid*1e-6)
		}
		if den > 0 {
			microDriftNorm = last / den
		}
	}

	retVals := e.returnWindow.Values()
	ret1s := ring.Mean(lastN(retVals, 32))
	vol10s := ring.Std(lastN(retVals, 256))
	entropy := signEntropy(lastN(retVals, 256))

	return types.FeatureState{
		Symbol:          e.Symbol,
		Mid:             e.mid,
		Micro:           e.micro,
		Spread:          e.spread,
		Tick:            e.tick,
		ImbalanceTop3:   imb3,
		ImbalanceTop5:   imb5,
		QDBid:           qdBid,
		QDAsk:           qdAsk,
		QRBid:           qrBid,
		QRAsk:           qrAsk,
		OFINorm:         clamp(ofiNorm, -1, 1),
		MicroDriftNorm:  clamp(microDriftNorm, -1, 1),
		Return1s:        ret1s,
		Vol10s:          vol10s,
		Entropy10s:      entropy,
		UpdatedAt:       time.Now(),
	}
}

func computeMicro(bid, ask, bidQty, askQty float64) float64 {
	den := bidQty + askQty
	if den == 0 {
		return 0.5 * (bid + ask)
	}
	return (ask*bidQty + bid*askQty) / den
}

func sumTopNQty(levels []types.Level, n int) float64 {
	s := 0.0
	for i := 0; i < n && i < len(levels); i++ {
		s += levels[i].Qty
	}
	return s
}

func computeImbalance(bid, ask float64) float64 {
	den := bid + ask
	if den == 0 {
		return 0
	}
	return (bid - ask) / den
}

func computeDepletion(prev, curr float64) float64 {
	if prev <= 0 || curr >= prev {
		return 0
	}
	return (prev - curr) / prev
}

func computeReplenishment(prev, curr float64) float64 {
	if prev <= 0 || curr <= prev {
		return 0
	}
	return (curr - prev) / prev
}

func calcTopLevelOFI(prevBidPrice, prevBidQty, prevAskPrice, prevAskQty, bidPrice, bidQty, askPrice, askQty float64) float64 {
	var e float64
	switch {
	case bidPrice > prevBidPrice:
		e += bidQty
	case bidPrice == prevBidPrice:
		e += bidQty - prevBidQty
	case bidPrice < prevBidPrice:
		e -= prevBidQty
	}
	switch {
	case askPrice < prevAskPrice:
		e += askQty
	case askPrice == prevAskPrice:
		e -= (askQty - prevAskQty)
	case askPrice > prevAskPrice:
		e -= prevAskQty
	}
	return e
}

func estimateTick(bids, asks []types.Level) float64 {
	best := math.MaxFloat64
	for i := 1; i < len(bids); i++ {
		d := math.Abs(bids[i-1].Price - bids[i].Price)
		if d > 0 && d < best {
			best = d
		}
	}
	for i := 1; i < len(asks); i++ {
		d := math.Abs(asks[i].Price - asks[i-1].Price)
		if d > 0 && d < best {
			best = d
		}
	}
	if best == math.MaxFloat64 {
		return 0
	}
	return best
}

func signEntropy(xs []float64) float64 {
	if len(xs) == 0 {
		return 1
	}
	var pos, neg int
	for _, x := range xs {
		if x > 0 {
			pos++
		} else if x < 0 {
			neg++
		}
	}
	total := float64(pos + neg)
	if total == 0 {
		return 1
	}
	pp := float64(pos) / total
	pn := float64(neg) / total
	ent := 0.0
	if pp > 0 {
		ent -= pp * math.Log(pp)
	}
	if pn > 0 {
		ent -= pn * math.Log(pn)
	}
	return ent / math.Log(2)
}

func lastN(xs []float64, n int) []float64 {
	if len(xs) <= n {
		return xs
	}
	return xs[len(xs)-n:]
}

func clamp(x, lo, hi float64) float64 {
	if x < lo {
		return lo
	}
	if x > hi {
		return hi
	}
	return x
}
```

---

## Файл: `internal/features/regime.go`

```go
package features

import "your_project/internal/types"

func DetectRegime(f types.FeatureState) types.RegimeState {
	if f.Mid <= 0 {
		return types.RegimeState{Tradable: false, Reason: "mid<=0"}
	}

	spreadNorm := f.Spread / f.Mid

	if spreadNorm > 0.00008 {
		return types.RegimeState{Tradable: false, Reason: "spread_too_wide"}
	}
	if f.Vol10s < 0.000005 {
		return types.RegimeState{Tradable: false, Reason: "vol_too_low"}
	}
	if f.Vol10s > 0.0005 {
		return types.RegimeState{Tradable: false, Reason: "vol_too_high"}
	}
	if f.Entropy10s > 0.98 {
		return types.RegimeState{Tradable: false, Reason: "entropy_too_high"}
	}

	return types.RegimeState{Tradable: true}
}
```

---

## Файл: `internal/features/scorer.go`

```go
package features

import "your_project/internal/types"

func Score(f types.FeatureState, r types.RegimeState) types.Signal {
	s := types.Signal{
		Symbol:    f.Symbol,
		NoTrade:   true,
		Reason:    r.Reason,
		CreatedAt: f.UpdatedAt,
	}

	if !r.Tradable {
		return s
	}

	longScore :=
		0.30*f.QDAsk -
			0.20*f.QRAsk +
			0.25*f.OFINorm +
			0.15*f.MicroDriftNorm +
			0.10*f.ImbalanceTop5

	shortScore :=
		0.30*f.QDBid -
			0.20*f.QRBid +
			0.25*(-f.OFINorm) +
			0.15*(-f.MicroDriftNorm) +
			0.10*(-f.ImbalanceTop5)

	s.LongScore = longScore
	s.ShortScore = shortScore

	long :=
		f.QDAsk > 0.55 &&
			f.QRAsk < 0.15 &&
			f.OFINorm > 0.20 &&
			f.Micro > f.Mid &&
			f.MicroDriftNorm > 0 &&
			f.ImbalanceTop5 > 0.08 &&
			longScore > 0.25

	short :=
		f.QDBid > 0.55 &&
			f.QRBid < 0.15 &&
			f.OFINorm < -0.20 &&
			f.Micro < f.Mid &&
			f.MicroDriftNorm < 0 &&
			f.ImbalanceTop5 < -0.08 &&
			shortScore > 0.25

	if long {
		s.Long = true
		s.NoTrade = false
		s.Reason = "long_signal"
		return s
	}
	if short {
		s.Short = true
		s.NoTrade = false
		s.Reason = "short_signal"
		return s
	}

	s.Reason = "no_strong_signal"
	return s
}
```

---

## Файл: `internal/strategy/strategy.go`

```go
package strategy

import (
	"time"

	"your_project/internal/types"
)

type Strategy struct {
	Pos              types.Position
	CooldownUntil    time.Time
	MinHold          time.Duration
	MaxHold          time.Duration
	EmergencyStopBps float64
}

func NewStrategy() *Strategy {
	return &Strategy{
		MinHold:          200 * time.Millisecond,
		MaxHold:          2 * time.Second,
		EmergencyStopBps: 0.0006,
	}
}

func (s *Strategy) HasPosition() bool {
	return s.Pos.Side != types.SideUnknown
}

func (s *Strategy) CanEnter(now time.Time) bool {
	return !s.HasPosition() && now.After(s.CooldownUntil)
}

func (s *Strategy) Enter(sig types.Signal, f types.FeatureState, now time.Time) {
	if sig.Long {
		s.Pos = types.Position{Side: types.SideBuy, Entry: f.Mid, EntryTime: now, Size: 1}
		return
	}
	if sig.Short {
		s.Pos = types.Position{Side: types.SideSell, Entry: f.Mid, EntryTime: now, Size: 1}
	}
}

func (s *Strategy) ShouldExit(f types.FeatureState, sig types.Signal, now time.Time) (bool, string) {
	if !s.HasPosition() {
		return false, ""
	}

	held := now.Sub(s.Pos.EntryTime)
	if held > s.MaxHold {
		return true, "max_hold"
	}

	move := 0.0
	if s.Pos.Entry > 0 {
		move = (f.Mid - s.Pos.Entry) / s.Pos.Entry
	}

	switch s.Pos.Side {
	case types.SideBuy:
		if move <= -s.EmergencyStopBps {
			return true, "emergency_stop_long"
		}
		if held >= s.MinHold {
			if f.OFINorm < 0 || f.Micro <= f.Mid || sig.LongScore < 0.10 {
				return true, "model_exit_long"
			}
		}
		if move >= 0.0007 {
			return true, "take_profit_long"
		}
	case types.SideSell:
		if move >= s.EmergencyStopBps {
			return true, "emergency_stop_short"
		}
		if held >= s.MinHold {
			if f.OFINorm > 0 || f.Micro >= f.Mid || sig.ShortScore < 0.10 {
				return true, "model_exit_short"
			}
		}
		if move <= -0.0007 {
			return true, "take_profit_short"
		}
	}

	return false, ""
}

func (s *Strategy) Exit(now time.Time) {
	s.Pos = types.Position{}
	s.CooldownUntil = now.Add(300 * time.Millisecond)
}
```

---

## Файл: `internal/paper/engine.go`

```go
package paper

import (
	"fmt"
	"time"

	"your_project/internal/features"
	"your_project/internal/strategy"
	"your_project/internal/types"
)

type Engine struct {
	Features *features.Engine
	Strategy *strategy.Strategy
}

func New(symbol string) *Engine {
	return &Engine{
		Features: features.NewEngine(symbol),
		Strategy: strategy.NewStrategy(),
	}
}

func (e *Engine) OnBook(book types.BookSnapshot) {
	e.Features.UpdateBook(book)
	f := e.Features.Snapshot()
	r := features.DetectRegime(f)
	sig := features.Score(f, r)
	now := time.Now()

	if e.Strategy.CanEnter(now) {
		if sig.Long || sig.Short {
			e.Strategy.Enter(sig, f, now)
			fmt.Printf("ENTER %s mid=%.2f long=%.3f short=%.3f\n", sig.Reason, f.Mid, sig.LongScore, sig.ShortScore)
		}
		return
	}

	if exit, reason := e.Strategy.ShouldExit(f, sig, now); exit {
		fmt.Printf("EXIT %s mid=%.2f\n", reason, f.Mid)
		e.Strategy.Exit(now)
	}
}

func (e *Engine) OnTrade(tr types.TradeTick) {
	e.Features.UpdateTrade(tr)
}
```

---

## Файл: `internal/storage/logger.go`

```go
package storage

import (
	"encoding/csv"
	"os"
	"strconv"
	"time"

	"your_project/internal/types"
)

type CSVLogger struct {
	f *os.File
	w *csv.Writer
}

func NewCSVLogger(path string) (*CSVLogger, error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	w := csv.NewWriter(f)
	_ = w.Write([]string{
		"ts", "symbol", "mid", "micro", "spread", "imb5",
		"qd_bid", "qd_ask", "qr_bid", "qr_ask",
		"ofi", "micro_drift", "vol10s", "entropy10s",
		"long_score", "short_score", "long", "short", "reason",
	})
	w.Flush()
	return &CSVLogger{f: f, w: w}, nil
}

func (l *CSVLogger) LogFeatureSignal(f types.FeatureState, s types.Signal) {
	_ = l.w.Write([]string{
		time.Now().Format(time.RFC3339Nano),
		f.Symbol,
		ff(f.Mid), ff(f.Micro), ff(f.Spread), ff(f.ImbalanceTop5),
		ff(f.QDBid), ff(f.QDAsk), ff(f.QRBid), ff(f.QRAsk),
		ff(f.OFINorm), ff(f.MicroDriftNorm), ff(f.Vol10s), ff(f.Entropy10s),
		ff(s.LongScore), ff(s.ShortScore),
		strconv.FormatBool(s.Long), strconv.FormatBool(s.Short), s.Reason,
	})
	l.w.Flush()
}

func (l *CSVLogger) Close() error {
	l.w.Flush()
	return l.f.Close()
}

func ff(v float64) string {
	return strconv.FormatFloat(v, 'f', 8, 64)
}
```

---

## Файл: `internal/config/config.go`

```go
package config

type Config struct {
	Symbol string
}

func Load() Config {
	return Config{
		Symbol: "BTC-USDT",
	}
}
```

---

## Файл: `cmd/trader/main.go`

```go
package main

import (
	"math/rand"
	"time"

	"your_project/internal/config"
	"your_project/internal/paper"
	"your_project/internal/types"
)

func main() {
	cfg := config.Load()
	eng := paper.New(cfg.Symbol)

	mid := 70000.0

	for i := 0; i < 10000; i++ {
		mid += (rand.Float64() - 0.5) * 2.0
		bid := mid - 0.25
		ask := mid + 0.25

		// synthetic liquidity regime
		bid1 := 3.0 + rand.Float64()*2.0
		ask1 := 3.0 + rand.Float64()*2.0

		// occasionally create bullish queue collapse
		if i%200 == 50 || i%200 == 51 {
			ask1 = 0.8
		}
		if i%200 == 120 || i%200 == 121 {
			bid1 = 0.8
		}

		book := types.BookSnapshot{
			Symbol: cfg.Symbol,
			Bids: []types.Level{
				{Price: bid, Qty: bid1},
				{Price: bid - 0.5, Qty: 2.5},
				{Price: bid - 1.0, Qty: 2.0},
				{Price: bid - 1.5, Qty: 1.8},
				{Price: bid - 2.0, Qty: 1.5},
			},
			Asks: []types.Level{
				{Price: ask, Qty: ask1},
				{Price: ask + 0.5, Qty: 2.5},
				{Price: ask + 1.0, Qty: 2.0},
				{Price: ask + 1.5, Qty: 1.8},
				{Price: ask + 2.0, Qty: 1.5},
			},
			Timestamp: time.Now(),
		}
		eng.OnBook(book)

		trSide := types.SideBuy
		if rand.Intn(2) == 0 {
			trSide = types.SideSell
		}
		if i%200 == 50 || i%200 == 51 {
			trSide = types.SideBuy
		}
		if i%200 == 120 || i%200 == 121 {
			trSide = types.SideSell
		}

		eng.OnTrade(types.TradeTick{
			Symbol:    cfg.Symbol,
			Price:     mid,
			Qty:       0.5 + rand.Float64(),
			Side:      trSide,
			Timestamp: time.Now(),
		})

		time.Sleep(20 * time.Millisecond)
	}
}
```

---

## Как эту систему валидировать

### 1. Сначала только логи

Пишите на каждый тик:

* признаки
* score
* fired/not fired
* future return через 500ms / 1s / 2s

### 2. Разметка label

Для long:

* label=1, если в следующие 2 секунды был move вверх > 3–5 bps
* label=0 иначе

Для short — симметрично.

### 3. Дальше смотрите precision по bucket-ам

Например:

* score 0.20–0.25
* score 0.25–0.35
* score 0.35+

И торгуете только верхний bucket.

---

## Что дальше улучшать

1. Реальный KuCoin/Binance WS adapter
2. local book sync по sequence
3. time-based windows вместо event-count
4. separate cancel-vs-trade heuristics
5. 2-tick confirmation
6. cooldown / stale-book filter / slippage model
7. live execution layer

---

## Что это за модель по сути

Это не одна “магическая формула”, а селективная модель на базе:

* non-equilibrium statistical mechanics
* queue collapse / barrier breakdown
* order-flow confirmation
* local regime gating

Именно за счет **жесткой селекции сигналов** она может давать **precision >60% на выбранных входах**, а не на всём потоке тиков.
