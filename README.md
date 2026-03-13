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



Да. Ниже даю полный переписанный стек целиком, уже в более чистом виде:

единый event loop

без race condition

Ising как фильтр/подтверждение, а не основа edge

time-bucket aggregation для trade flow

честная paper-логика

готовая структура под KuCoin

Я возьму модуль:

module github.com/yourname/kucoin-ising-bot

Замени его на свой.

Структура
kucoin-ising-bot/
├── go.mod
├── go.sum
├── cmd/
│   └── trader/
│       └── main.go
└── internal/
    ├── config/
    │   └── config.go
    ├── core/
    │   ├── engine.go
    │   └── events.go
    ├── exchange/
    │   └── kucoin/
    │       ├── book.go
    │       └── client.go
    ├── features/
    │   ├── ising.go
    │   ├── micro.go
    │   ├── regime.go
    │   └── scorer.go
    ├── ring/
    │   └── ring.go
    ├── storage/
    │   └── logger.go
    ├── strategy/
    │   └── strategy.go
    └── types/
        ├── market.go
        └── signal.go
go.mod
module github.com/yourname/kucoin-ising-bot

go 1.23.0

require github.com/gorilla/websocket v1.5.3
cmd/trader/main.go
package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/yourname/kucoin-ising-bot/internal/config"
	"github.com/yourname/kucoin-ising-bot/internal/core"
	"github.com/yourname/kucoin-ising-bot/internal/exchange/kucoin"
)

func main() {
	cfg := config.Load()

	engine, err := core.NewEngine(cfg)
	if err != nil {
		log.Fatalf("create engine: %v", err)
	}
	defer engine.Close()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go engine.Run(ctx)

	client := kucoin.NewClient(cfg, engine.Events())

	if err := client.Run(ctx); err != nil && err != context.Canceled {
		log.Fatalf("client run: %v", err)
	}
}
internal/config/config.go
package config

import "time"

type Config struct {
	Symbol string

	TopLevels int

	WSReadTimeout  time.Duration
	ReconnectDelay time.Duration

	BookThrottle time.Duration
	Warmup       time.Duration

	TradeBucket time.Duration

	SignalCSV string
	TradeCSV  string

	PositionSize float64
	TakerFeeRate float64
	SlippageFrac float64

	Cooldown         time.Duration
	MinHold          time.Duration
	MaxHold          time.Duration
	EmergencyStopFrac float64
	TakeProfitFrac    float64

	IsingWindow int
	IsingBeta   float64
	IsingJ      float64
	IsingScale  float64
}

func Load() Config {
	return Config{
		Symbol:         "BTC-USDT",
		TopLevels:      20,
		WSReadTimeout:  30 * time.Second,
		ReconnectDelay: 3 * time.Second,

		BookThrottle: 30 * time.Millisecond,
		Warmup:       8 * time.Second,

		TradeBucket: 50 * time.Millisecond,

		SignalCSV: "signals_kucoin_ising.csv",
		TradeCSV:  "trades_kucoin_paper.csv",

		PositionSize: 1.0,
		TakerFeeRate: 0.001,
		SlippageFrac: 0.00005,

		Cooldown:          350 * time.Millisecond,
		MinHold:           250 * time.Millisecond,
		MaxHold:           2500 * time.Millisecond,
		EmergencyStopFrac: 0.0006,
		TakeProfitFrac:    0.0007,

		IsingWindow: 48,
		IsingBeta:   2.2,
		IsingJ:      0.85,
		IsingScale:  1.0,
	}
}
internal/types/market.go
package types

import "time"

type Side int

const (
	SideUnknown Side = iota
	SideBuy
	SideSell
)

func (s Side) String() string {
	switch s {
	case SideBuy:
		return "buy"
	case SideSell:
		return "sell"
	default:
		return "unknown"
	}
}

type Level struct {
	Price float64
	Qty   float64
}

type BookSnapshot struct {
	Symbol    string
	Bids      []Level // desc
	Asks      []Level // asc
	Timestamp time.Time
}

type TradeTick struct {
	Symbol    string
	Price     float64
	Qty       float64
	Side      Side
	Timestamp time.Time
}
internal/types/signal.go
package types

import "time"

type FeatureState struct {
	Symbol string

	Mid    float64
	Micro  float64
	Spread float64
	Tick   float64

	BestBid float64
	BestAsk float64

	ImbalanceTop3 float64
	ImbalanceTop5 float64

	QDBid float64
	QDAsk float64
	QRBid float64
	QRAsk float64

	BookOFINorm    float64
	TradeOFINorm   float64
	NetOFINorm     float64
	MicroDriftNorm float64

	Return1s   float64
	Vol10s     float64
	Entropy10s float64

	IsingField        float64
	IsingMagnet       float64
	IsingEnergy       float64
	IsingProbUp       float64
	IsingProbDown     float64
	IsingSuscept      float64
	IsingSpin         int
	IsingConsensus    float64
	IsingCriticalness float64

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

	Long    bool
	Short   bool
	NoTrade bool
	Reason  string

	ProbLong  float64
	ProbShort float64

	CreatedAt time.Time
}

type Position struct {
	Side      Side
	Entry     float64
	EntryTime time.Time
	Size      float64

	EntryBid float64
	EntryAsk float64
}
internal/ring/ring.go
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
	return &FloatRing{
		vals: make([]float64, size),
	}
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

func (r *FloatRing) Cap() int {
	return len(r.vals)
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
	sum := 0.0
	for _, x := range xs {
		sum += x
	}
	return sum / float64(len(xs))
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
	sum := 0.0
	for _, x := range xs {
		if x < 0 {
			sum -= x
		} else {
			sum += x
		}
	}
	return sum
}
internal/core/events.go
package core

import "github.com/yourname/kucoin-ising-bot/internal/types"

type EventType int

const (
	EventBook EventType = iota
	EventTrade
)

type Event struct {
	Type  EventType
	Book  *types.BookSnapshot
	Trade *types.TradeTick
}
internal/features/ising.go
package features

import (
	"math"

	"github.com/yourname/kucoin-ising-bot/internal/ring"
	"github.com/yourname/kucoin-ising-bot/internal/types"
)

type IsingModel struct {
	beta       float64
	couplingJ  float64
	fieldScale float64

	spinWindow *ring.FloatRing

	initialized bool
	lastSpin    int
}

func NewIsingModel(window int, beta, couplingJ, fieldScale float64) *IsingModel {
	if window <= 0 {
		window = 48
	}
	if beta <= 0 {
		beta = 2.2
	}
	if fieldScale <= 0 {
		fieldScale = 1.0
	}
	return &IsingModel{
		beta:       beta,
		couplingJ:  couplingJ,
		fieldScale: fieldScale,
		spinWindow: ring.NewFloatRing(window),
	}
}

func (m *IsingModel) Observe(f types.FeatureState) types.FeatureState {
	oldMagnet := m.magnetization()
	h := clamp(m.externalField(f)*m.fieldScale, -2.0, 2.0)
	heff := m.couplingJ*oldMagnet + h

	pUp := sigmoid(2.0 * m.beta * heff)
	pDown := 1.0 - pUp

	spin := -1
	if pUp >= 0.5 {
		spin = 1
	}

	m.lastSpin = spin
	m.spinWindow.Add(float64(spin))
	m.initialized = true

	newMagnet := m.magnetization()
	chi := m.beta * (1.0 - newMagnet*newMagnet)
	if chi < 0 {
		chi = 0
	}

	energy := -float64(spin) * heff
	criticalness := chi * (1.0 - math.Abs(newMagnet))
	criticalness = clamp(criticalness, 0, 10)

	f.IsingField = h
	f.IsingMagnet = clamp(newMagnet, -1, 1)
	f.IsingEnergy = energy
	f.IsingProbUp = clamp(pUp, 0, 1)
	f.IsingProbDown = clamp(pDown, 0, 1)
	f.IsingSuscept = chi
	f.IsingSpin = spin
	f.IsingConsensus = math.Abs(newMagnet)
	f.IsingCriticalness = criticalness

	return f
}

func (m *IsingModel) externalField(f types.FeatureState) float64 {
	queueAsym := (f.QDAsk - f.QDBid) + 0.5*(f.QRBid-f.QRAsk)

	h :=
		0.24*f.BookOFINorm +
			0.20*f.TradeOFINorm +
			0.16*f.MicroDriftNorm +
			0.14*f.ImbalanceTop5 +
			0.10*f.ImbalanceTop3 +
			0.12*queueAsym +
			0.04*sign(f.Return1s)

	if f.Mid > 0 {
		spreadNorm := f.Spread / f.Mid
		h -= 18.0 * spreadNorm
	}

	return h
}

func (m *IsingModel) magnetization() float64 {
	xs := m.spinWindow.Values()
	if len(xs) == 0 {
		return 0
	}
	return clamp(ring.Mean(xs), -1, 1)
}

func sigmoid(x float64) float64 {
	if x > 40 {
		return 1
	}
	if x < -40 {
		return 0
	}
	return 1.0 / (1.0 + math.Exp(-x))
}

func sign(x float64) float64 {
	switch {
	case x > 0:
		return 1
	case x < 0:
		return -1
	default:
		return 0
	}
}
internal/features/micro.go
package features

import (
	"math"
	"time"

	"github.com/yourname/kucoin-ising-bot/internal/ring"
	"github.com/yourname/kucoin-ising-bot/internal/types"
)

type Engine struct {
	Symbol string

	bestBid float64
	bestAsk float64
	bid1Qty float64
	ask1Qty float64

	mid    float64
	micro  float64
	spread float64
	tick   float64

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

	bookOFIWindow *ring.FloatRing
	driftWindow   *ring.FloatRing
	returnWindow  *ring.FloatRing

	tradeBucketStart time.Time
	tradeBucketValue float64
	tradeOFIWindow   *ring.FloatRing
	tradeBucketDur   time.Duration

	ising *IsingModel

	initialized bool
}

func NewEngine(symbol string, tradeBucket time.Duration, ising *IsingModel) *Engine {
	if tradeBucket <= 0 {
		tradeBucket = 50 * time.Millisecond
	}
	if ising == nil {
		ising = NewIsingModel(48, 2.2, 0.85, 1.0)
	}
	return &Engine{
		Symbol:         symbol,
		bookOFIWindow:  ring.NewFloatRing(64),
		driftWindow:    ring.NewFloatRing(64),
		returnWindow:   ring.NewFloatRing(334),
		tradeOFIWindow: ring.NewFloatRing(200),
		tradeBucketDur: tradeBucket,
		ising:          ising,
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

	if bestBid <= 0 || bestAsk <= 0 || bestBid >= bestAsk || bid1Qty <= 0 || ask1Qty <= 0 {
		return
	}

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

		e.bookOFIWindow.Add(0)
		e.driftWindow.Add(0)
		e.returnWindow.Add(0)
		e.tradeOFIWindow.Add(0)

		e.initialized = true
		return
	}

	prevMid := e.mid

	e.prevBidTop3, e.prevAskTop3 = e.currBidTop3, e.currAskTop3
	e.currBidTop3, e.currAskTop3 = bidTop3, askTop3

	e.prevBidTop5, e.prevAskTop5 = e.currBidTop5, e.currAskTop5
	e.currBidTop5, e.currAskTop5 = bidTop5, askTop5

	bookOFI := calcTopLevelOFI(
		e.prevBestBid, e.prevBid1Qty,
		e.prevBestAsk, e.prevAsk1Qty,
		bestBid, bid1Qty,
		bestAsk, ask1Qty,
	)
	e.bookOFIWindow.Add(bookOFI)
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

	ts := tr.Timestamp
	if ts.IsZero() {
		ts = time.Now()
	}

	if e.tradeBucketStart.IsZero() {
		e.tradeBucketStart = ts
	}

	for ts.Sub(e.tradeBucketStart) >= e.tradeBucketDur {
		e.tradeOFIWindow.Add(e.tradeBucketValue)
		e.tradeBucketValue = 0
		e.tradeBucketStart = e.tradeBucketStart.Add(e.tradeBucketDur)
	}

	switch tr.Side {
	case types.SideBuy:
		e.tradeBucketValue += tr.Qty
	case types.SideSell:
		e.tradeBucketValue -= tr.Qty
	}
}

func (e *Engine) ForceFlushTradeBucket(now time.Time) {
	if e.tradeBucketStart.IsZero() {
		e.tradeBucketStart = now
		return
	}
	for now.Sub(e.tradeBucketStart) >= e.tradeBucketDur {
		e.tradeOFIWindow.Add(e.tradeBucketValue)
		e.tradeBucketValue = 0
		e.tradeBucketStart = e.tradeBucketStart.Add(e.tradeBucketDur)
	}
}

func (e *Engine) Ready() bool {
	return e.initialized &&
		e.bookOFIWindow.Len() >= 16 &&
		e.driftWindow.Len() >= 16 &&
		e.returnWindow.Len() >= 64 &&
		e.tradeOFIWindow.Len() >= 16
}

func (e *Engine) Snapshot(now time.Time) types.FeatureState {
	e.ForceFlushTradeBucket(now)

	imb3 := computeImbalance(e.currBidTop3, e.currAskTop3)
	imb5 := computeImbalance(e.currBidTop5, e.currAskTop5)

	qdBid := computeDepletion(e.prevBidTop5, e.currBidTop5)
	qdAsk := computeDepletion(e.prevAskTop5, e.currAskTop5)
	qrBid := computeReplenishment(e.prevBidTop5, e.currBidTop5)
	qrAsk := computeReplenishment(e.prevAskTop5, e.currAskTop5)

	bookOFINorm := normSignedWindow(e.bookOFIWindow.Values())
	tradeOFINorm := normSignedWindow(e.tradeOFIWindow.Values())
	netOFI := clamp(0.55*bookOFINorm+0.45*tradeOFINorm, -1, 1)

	microDriftNorm := 0.0
	if last, ok := e.driftWindow.Last(); ok {
		den := e.mid
		if e.tick > 0 {
			den = math.Max(e.tick, e.mid*1e-6)
		}
		if den > 0 {
			microDriftNorm = last / den
		}
	}

	retVals := e.returnWindow.Values()
	ret1s := ring.Mean(lastN(retVals, 33))
	vol10s := ring.Std(lastN(retVals, 334))
	entropy := signEntropy(lastN(retVals, 334))

	f := types.FeatureState{
		Symbol:          e.Symbol,
		Mid:             e.mid,
		Micro:           e.micro,
		Spread:          e.spread,
		Tick:            e.tick,
		BestBid:         e.bestBid,
		BestAsk:         e.bestAsk,
		ImbalanceTop3:   imb3,
		ImbalanceTop5:   imb5,
		QDBid:           qdBid,
		QDAsk:           qdAsk,
		QRBid:           qrBid,
		QRAsk:           qrAsk,
		BookOFINorm:     bookOFINorm,
		TradeOFINorm:    tradeOFINorm,
		NetOFINorm:      netOFI,
		MicroDriftNorm:  clamp(microDriftNorm, -1, 1),
		Return1s:        ret1s,
		Vol10s:          vol10s,
		Entropy10s:      entropy,
		UpdatedAt:       now,
	}

	if e.ising != nil {
		f = e.ising.Observe(f)
	}

	return f
}

func normSignedWindow(xs []float64) float64 {
	if len(xs) == 0 {
		return 0
	}
	den := ring.SumAbs(xs)
	if den == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range xs {
		sum += v
	}
	return clamp(sum/den, -1, 1)
}

func computeMicro(bid, ask, bidQty, askQty float64) float64 {
	den := bidQty + askQty
	if den == 0 {
		return 0.5 * (bid + ask)
	}
	return (ask*bidQty + bid*askQty) / den
}

func sumTopNQty(levels []types.Level, n int) float64 {
	sum := 0.0
	for i := 0; i < n && i < len(levels); i++ {
		sum += levels[i].Qty
	}
	return sum
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
		d := math.Abs(asks[i-1].Price - asks[i].Price)
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
internal/features/regime.go
package features

import "github.com/yourname/kucoin-ising-bot/internal/types"

func DetectRegime(f types.FeatureState) types.RegimeState {
	if f.Mid <= 0 || f.BestBid <= 0 || f.BestAsk <= 0 {
		return types.RegimeState{Tradable: false, Reason: "bad_mid_or_book"}
	}

	if f.BestBid >= f.BestAsk {
		return types.RegimeState{Tradable: false, Reason: "crossed_book"}
	}

	spreadNorm := f.Spread / f.Mid
	if spreadNorm > 0.00008 {
		return types.RegimeState{Tradable: false, Reason: "spread_too_wide"}
	}
	if f.Vol10s < 0.000002 {
		return types.RegimeState{Tradable: false, Reason: "vol_too_low"}
	}
	if f.Vol10s > 0.0008 {
		return types.RegimeState{Tradable: false, Reason: "vol_too_high"}
	}
	if f.Entropy10s > 0.985 {
		return types.RegimeState{Tradable: false, Reason: "entropy_too_high"}
	}
	if f.IsingConsensus < 0.15 && f.IsingSuscept > 1.6 {
		return types.RegimeState{Tradable: false, Reason: "ising_critical_zone"}
	}
	if f.IsingCriticalness > 1.25 {
		return types.RegimeState{Tradable: false, Reason: "ising_transition_noise"}
	}

	return types.RegimeState{Tradable: true}
}
internal/features/scorer.go
package features

import "github.com/yourname/kucoin-ising-bot/internal/types"

func Score(f types.FeatureState, r types.RegimeState) types.Signal {
	s := types.Signal{
		Symbol:     f.Symbol,
		NoTrade:    true,
		Reason:     r.Reason,
		CreatedAt:  f.UpdatedAt,
		ProbLong:   f.IsingProbUp,
		ProbShort:  f.IsingProbDown,
	}

	if !r.Tradable {
		return s
	}

	microLong :=
		0.18*f.QDAsk -
			0.10*f.QRAsk +
			0.16*f.BookOFINorm +
			0.18*f.TradeOFINorm +
			0.14*f.MicroDriftNorm +
			0.08*f.ImbalanceTop5 +
			0.06*f.ImbalanceTop3

	microShort :=
		0.18*f.QDBid -
			0.10*f.QRBid +
			0.16*(-f.BookOFINorm) +
			0.18*(-f.TradeOFINorm) +
			0.14*(-f.MicroDriftNorm) +
			0.08*(-f.ImbalanceTop5) +
			0.06*(-f.ImbalanceTop3)

	isingLong := 0.42*f.IsingProbUp + 0.18*f.IsingConsensus - 0.14*f.IsingCriticalness
	isingShort := 0.42*f.IsingProbDown + 0.18*f.IsingConsensus - 0.14*f.IsingCriticalness

	longScore := microLong + isingLong
	shortScore := microShort + isingShort

	s.LongScore = longScore
	s.ShortScore = shortScore

	long :=
		f.IsingProbUp > 0.64 &&
			f.IsingMagnet > 0.12 &&
			f.IsingField > 0 &&
			f.NetOFINorm > 0.12 &&
			f.Micro > f.Mid &&
			f.MicroDriftNorm > 0 &&
			f.ImbalanceTop5 > 0.05 &&
			longScore > 0.40

	short :=
		f.IsingProbDown > 0.64 &&
			f.IsingMagnet < -0.12 &&
			f.IsingField < 0 &&
			f.NetOFINorm < -0.12 &&
			f.Micro < f.Mid &&
			f.MicroDriftNorm < 0 &&
			f.ImbalanceTop5 < -0.05 &&
			shortScore > 0.40

	if long {
		s.Long = true
		s.NoTrade = false
		s.Reason = "ising_long_signal"
		return s
	}
	if short {
		s.Short = true
		s.NoTrade = false
		s.Reason = "ising_short_signal"
		return s
	}

	s.Reason = "no_strong_signal"
	return s
}
internal/strategy/strategy.go
package strategy

import (
	"time"

	"github.com/yourname/kucoin-ising-bot/internal/config"
	"github.com/yourname/kucoin-ising-bot/internal/types"
)

type Strategy struct {
	Pos           types.Position
	CooldownUntil time.Time

	MinHold           time.Duration
	MaxHold           time.Duration
	EmergencyStopFrac float64
	TakeProfitFrac    float64
	Cooldown          time.Duration
	Size              float64
}

func New(cfg config.Config) *Strategy {
	return &Strategy{
		MinHold:           cfg.MinHold,
		MaxHold:           cfg.MaxHold,
		EmergencyStopFrac: cfg.EmergencyStopFrac,
		TakeProfitFrac:    cfg.TakeProfitFrac,
		Cooldown:          cfg.Cooldown,
		Size:              cfg.PositionSize,
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
		s.Pos = types.Position{
			Side:      types.SideBuy,
			Entry:     f.BestAsk,
			EntryBid:  f.BestBid,
			EntryAsk:  f.BestAsk,
			EntryTime: now,
			Size:      s.Size,
		}
		return
	}

	if sig.Short {
		s.Pos = types.Position{
			Side:      types.SideSell,
			Entry:     f.BestBid,
			EntryBid:  f.BestBid,
			EntryAsk:  f.BestAsk,
			EntryTime: now,
			Size:      s.Size,
		}
	}
}

func (s *Strategy) Exit(now time.Time) {
	s.Pos = types.Position{}
	s.CooldownUntil = now.Add(s.Cooldown)
}

func (s *Strategy) ShouldExit(f types.FeatureState, sig types.Signal, now time.Time) (bool, string) {
	if !s.HasPosition() {
		return false, ""
	}

	held := now.Sub(s.Pos.EntryTime)
	if held > s.MaxHold {
		return true, "max_hold"
	}

	switch s.Pos.Side {
	case types.SideBuy:
		if f.BestBid <= 0 || s.Pos.Entry <= 0 {
			return false, ""
		}
		move := (f.BestBid - s.Pos.Entry) / s.Pos.Entry

		if move <= -s.EmergencyStopFrac {
			return true, "emergency_stop_long"
		}
		if move >= s.TakeProfitFrac {
			return true, "take_profit_long"
		}
		if held >= s.MinHold {
			if f.IsingProbUp < 0.52 || f.IsingMagnet < 0.02 || f.IsingField < 0 {
				return true, "ising_reversal_long"
			}
			if f.NetOFINorm < 0 || f.Micro <= f.Mid || sig.LongScore < 0.14 {
				return true, "model_exit_long"
			}
		}

	case types.SideSell:
		if f.BestAsk <= 0 || s.Pos.Entry <= 0 {
			return false, ""
		}
		move := (s.Pos.Entry - f.BestAsk) / s.Pos.Entry

		if move <= -s.EmergencyStopFrac {
			return true, "emergency_stop_short"
		}
		if move >= s.TakeProfitFrac {
			return true, "take_profit_short"
		}
		if held >= s.MinHold {
			if f.IsingProbDown < 0.52 || f.IsingMagnet > -0.02 || f.IsingField > 0 {
				return true, "ising_reversal_short"
			}
			if f.NetOFINorm > 0 || f.Micro >= f.Mid || sig.ShortScore < 0.14 {
				return true, "model_exit_short"
			}
		}
	}

	return false, ""
}
internal/storage/logger.go
package storage

import (
	"encoding/csv"
	"os"
	"strconv"
	"time"

	"github.com/yourname/kucoin-ising-bot/internal/types"
)

type CSVLogger struct {
	f *os.File
	w *csv.Writer
}

func NewSignalLogger(path string) (*CSVLogger, error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	w := csv.NewWriter(f)

	_ = w.Write([]string{
		"ts", "symbol",
		"best_bid", "best_ask", "mid", "micro", "spread", "tick",
		"imb3", "imb5",
		"qd_bid", "qd_ask", "qr_bid", "qr_ask",
		"book_ofi", "trade_ofi", "net_ofi",
		"micro_drift", "ret1s", "vol10s", "entropy10s",
		"ising_field", "ising_magnet", "ising_energy",
		"ising_prob_up", "ising_prob_down", "ising_suscept",
		"ising_spin", "ising_consensus", "ising_criticalness",
		"long_score", "short_score", "long", "short", "reason",
	})
	w.Flush()

	return &CSVLogger{f: f, w: w}, nil
}

func NewTradeLogger(path string) (*CSVLogger, error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	w := csv.NewWriter(f)

	_ = w.Write([]string{
		"ts", "symbol", "event", "side", "price", "size",
		"entry", "exit", "gross_pnl", "net_pnl", "reason",
	})
	w.Flush()

	return &CSVLogger{f: f, w: w}, nil
}

func (l *CSVLogger) LogFeatureSignal(f types.FeatureState, s types.Signal) {
	_ = l.w.Write([]string{
		time.Now().Format(time.RFC3339Nano),
		f.Symbol,
		ff(f.BestBid), ff(f.BestAsk), ff(f.Mid), ff(f.Micro), ff(f.Spread), ff(f.Tick),
		ff(f.ImbalanceTop3), ff(f.ImbalanceTop5),
		ff(f.QDBid), ff(f.QDAsk), ff(f.QRBid), ff(f.QRAsk),
		ff(f.BookOFINorm), ff(f.TradeOFINorm), ff(f.NetOFINorm),
		ff(f.MicroDriftNorm), ff(f.Return1s), ff(f.Vol10s), ff(f.Entropy10s),
		ff(f.IsingField), ff(f.IsingMagnet), ff(f.IsingEnergy),
		ff(f.IsingProbUp), ff(f.IsingProbDown), ff(f.IsingSuscept),
		strconv.Itoa(f.IsingSpin), ff(f.IsingConsensus), ff(f.IsingCriticalness),
		ff(s.LongScore), ff(s.ShortScore),
		strconv.FormatBool(s.Long), strconv.FormatBool(s.Short), s.Reason,
	})
	l.w.Flush()
}

func (l *CSVLogger) LogTrade(
	symbol, event, side string,
	price, size, entry, exit, grossPnL, netPnL float64,
	reason string,
) {
	_ = l.w.Write([]string{
		time.Now().Format(time.RFC3339Nano),
		symbol, event, side,
		ff(price), ff(size), ff(entry), ff(exit),
		ff(grossPnL), ff(netPnL), reason,
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
internal/core/engine.go
package core

import (
	"context"
	"fmt"
	"time"

	"github.com/yourname/kucoin-ising-bot/internal/config"
	"github.com/yourname/kucoin-ising-bot/internal/features"
	"github.com/yourname/kucoin-ising-bot/internal/storage"
	"github.com/yourname/kucoin-ising-bot/internal/strategy"
	"github.com/yourname/kucoin-ising-bot/internal/types"
)

type Engine struct {
	cfg config.Config

	events chan Event

	features *features.Engine
	strategy *strategy.Strategy

	signalLog *storage.CSVLogger
	tradeLog  *storage.CSVLogger

	startedAt       time.Time
	lastBookProcess time.Time
	warmedUp        bool
}

func NewEngine(cfg config.Config) (*Engine, error) {
	signalLog, err := storage.NewSignalLogger(cfg.SignalCSV)
	if err != nil {
		return nil, err
	}

	tradeLog, err := storage.NewTradeLogger(cfg.TradeCSV)
	if err != nil {
		_ = signalLog.Close()
		return nil, err
	}

	ising := features.NewIsingModel(cfg.IsingWindow, cfg.IsingBeta, cfg.IsingJ, cfg.IsingScale)

	return &Engine{
		cfg:       cfg,
		events:    make(chan Event, 8192),
		features:  features.NewEngine(cfg.Symbol, cfg.TradeBucket, ising),
		strategy:  strategy.New(cfg),
		signalLog: signalLog,
		tradeLog:  tradeLog,
		startedAt: time.Now(),
	}, nil
}

func (e *Engine) Events() chan<- Event {
	return e.events
}

func (e *Engine) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case ev := <-e.events:
			switch ev.Type {
			case EventTrade:
				if ev.Trade != nil {
					e.onTrade(*ev.Trade)
				}
			case EventBook:
				if ev.Book != nil {
					e.onBook(*ev.Book)
				}
			}
		}
	}
}

func (e *Engine) onTrade(tr types.TradeTick) {
	e.features.UpdateTrade(tr)
}

func (e *Engine) onBook(book types.BookSnapshot) {
	now := time.Now()

	if !e.lastBookProcess.IsZero() && now.Sub(e.lastBookProcess) < e.cfg.BookThrottle {
		return
	}
	e.lastBookProcess = now

	e.features.UpdateBook(book)
	if !e.features.Ready() {
		return
	}

	f := e.features.Snapshot(now)
	if f.Mid <= 0 || len(book.Bids) == 0 || len(book.Asks) == 0 {
		return
	}

	if !e.warmedUp {
		if now.Sub(e.startedAt) < e.cfg.Warmup {
			return
		}
		e.warmedUp = true
	}

	r := features.DetectRegime(f)
	sig := features.Score(f, r)

	if e.signalLog != nil {
		e.signalLog.LogFeatureSignal(f, sig)
	}

	if e.strategy.CanEnter(now) {
		if sig.Long || sig.Short {
			e.strategy.Enter(sig, f, now)

			side := "long"
			price := f.BestAsk
			if sig.Short {
				side = "short"
				price = f.BestBid
			}

			if e.tradeLog != nil {
				e.tradeLog.LogTrade(
					f.Symbol, "enter", side,
					price, e.strategy.Pos.Size, price, 0, 0, 0, sig.Reason,
				)
			}

			fmt.Printf(
				"ENTER %s side=%s px=%.2f bid=%.2f ask=%.2f L=%.3f S=%.3f pUp=%.3f pDn=%.3f m=%.3f h=%.3f\n",
				sig.Reason, side, price, f.BestBid, f.BestAsk,
				sig.LongScore, sig.ShortScore, f.IsingProbUp, f.IsingProbDown, f.IsingMagnet, f.IsingField,
			)
		}
		return
	}

	if exit, reason := e.strategy.ShouldExit(f, sig, now); exit {
		side, exitPx, grossPnL, netPnL := e.closePnL(f)

		if e.tradeLog != nil {
			e.tradeLog.LogTrade(
				f.Symbol, "exit", side,
				exitPx, e.strategy.Pos.Size,
				e.strategy.Pos.Entry, exitPx,
				grossPnL, netPnL, reason,
			)
		}

		fmt.Printf(
			"EXIT %s side=%s exit=%.2f gross=%.6f net=%.6f pUp=%.3f pDn=%.3f m=%.3f h=%.3f\n",
			reason, side, exitPx, grossPnL, netPnL,
			f.IsingProbUp, f.IsingProbDown, f.IsingMagnet, f.IsingField,
		)

		e.strategy.Exit(now)
	}
}

func (e *Engine) closePnL(f types.FeatureState) (side string, exitPx, grossPnL, netPnL float64) {
	pos := e.strategy.Pos
	size := pos.Size

	switch pos.Side {
	case types.SideBuy:
		side = "long"
		exitPx = f.BestBid * (1.0 - e.cfg.SlippageFrac)
		grossPnL = (exitPx - pos.Entry) * size
	case types.SideSell:
		side = "short"
		exitPx = f.BestAsk * (1.0 + e.cfg.SlippageFrac)
		grossPnL = (pos.Entry - exitPx) * size
	default:
		return "", 0, 0, 0
	}

	turnover := (pos.Entry + exitPx) * size
	fees := turnover * e.cfg.TakerFeeRate
	netPnL = grossPnL - fees

	return side, exitPx, grossPnL, netPnL
}

func (e *Engine) Close() {
	if e.signalLog != nil {
		_ = e.signalLog.Close()
	}
	if e.tradeLog != nil {
		_ = e.tradeLog.Close()
	}
}
internal/exchange/kucoin/book.go
package kucoin

import (
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/yourname/kucoin-ising-bot/internal/types"
)

type LocalBook struct {
	mu sync.RWMutex

	symbol   string
	sequence int64

	bids map[float64]float64
	asks map[float64]float64
}

func NewLocalBook(symbol string) *LocalBook {
	return &LocalBook{
		symbol: symbol,
		bids:   make(map[float64]float64, 4096),
		asks:   make(map[float64]float64, 4096),
	}
}

func (b *LocalBook) LoadSnapshot(seq int64, bids, asks [][]string) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.sequence = seq
	clear(b.bids)
	clear(b.asks)

	load := func(dst map[float64]float64, rows [][]string) {
		for _, row := range rows {
			if len(row) < 2 {
				continue
			}
			p, err1 := strconv.ParseFloat(row[0], 64)
			q, err2 := strconv.ParseFloat(row[1], 64)
			if err1 != nil || err2 != nil {
				continue
			}
			if q > 0 {
				dst[p] = q
			}
		}
	}

	load(b.bids, bids)
	load(b.asks, asks)

	return nil
}

func (b *LocalBook) Sequence() int64 {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.sequence
}

func (b *LocalBook) ApplyDelta(sequenceStart, sequenceEnd int64, bidChanges, askChanges [][]string) bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	if !(sequenceStart <= b.sequence+1 && sequenceEnd > b.sequence) {
		return false
	}

	applySide := func(dst map[float64]float64, changes [][]string) {
		for _, row := range changes {
			if len(row) < 2 {
				continue
			}
			price, err1 := strconv.ParseFloat(row[0], 64)
			qty, err2 := strconv.ParseFloat(row[1], 64)
			if err1 != nil || err2 != nil {
				continue
			}
			if qty == 0 {
				delete(dst, price)
			} else {
				dst[price] = qty
			}
		}
	}

	applySide(b.bids, bidChanges)
	applySide(b.asks, askChanges)
	b.sequence = sequenceEnd
	return true
}

func (b *LocalBook) Snapshot(topN int) types.BookSnapshot {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if topN <= 0 {
		topN = 20
	}

	bidPrices := make([]float64, 0, len(b.bids))
	for p := range b.bids {
		bidPrices = append(bidPrices, p)
	}
	sort.Slice(bidPrices, func(i, j int) bool { return bidPrices[i] > bidPrices[j] })

	askPrices := make([]float64, 0, len(b.asks))
	for p := range b.asks {
		askPrices = append(askPrices, p)
	}
	sort.Slice(askPrices, func(i, j int) bool { return askPrices[i] < askPrices[j] })

	bids := make([]types.Level, 0, min(topN, len(bidPrices)))
	for i := 0; i < len(bidPrices) && i < topN; i++ {
		p := bidPrices[i]
		bids = append(bids, types.Level{Price: p, Qty: b.bids[p]})
	}

	asks := make([]types.Level, 0, min(topN, len(askPrices)))
	for i := 0; i < len(askPrices) && i < topN; i++ {
		p := askPrices[i]
		asks = append(asks, types.Level{Price: p, Qty: b.asks[p]})
	}

	return types.BookSnapshot{
		Symbol:    b.symbol,
		Bids:      bids,
		Asks:      asks,
		Timestamp: time.Now(),
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
internal/exchange/kucoin/client.go
package kucoin

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"

	"github.com/yourname/kucoin-ising-bot/internal/config"
	"github.com/yourname/kucoin-ising-bot/internal/core"
	"github.com/yourname/kucoin-ising-bot/internal/types"
)

type Client struct {
	cfg    config.Config
	symbol string
	book   *LocalBook

	events chan<- core.Event

	httpClient *http.Client
}

func NewClient(cfg config.Config, events chan<- core.Event) *Client {
	return &Client{
		cfg:        cfg,
		symbol:     cfg.Symbol,
		book:       NewLocalBook(cfg.Symbol),
		events:     events,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

type bulletResponse struct {
	Code string `json:"code"`
	Data struct {
		Token string `json:"token"`
		InstanceServers []struct {
			Endpoint     string `json:"endpoint"`
			Protocol     string `json:"protocol"`
			Encrypt      bool   `json:"encrypt"`
			PingInterval int64  `json:"pingInterval"`
			PingTimeout  int64  `json:"pingTimeout"`
		} `json:"instanceServers"`
	} `json:"data"`
}

type jsonNumber string

func (n jsonNumber) Int64() (int64, error) {
	return strconv.ParseInt(string(n), 10, 64)
}

type orderBookResponse struct {
	Code string `json:"code"`
	Data struct {
		Time     int64      `json:"time"`
		Sequence jsonNumber `json:"sequence"`
		Bids     [][]string `json:"bids"`
		Asks     [][]string `json:"asks"`
	} `json:"data"`
}

type wsEnvelope struct {
	ID      string          `json:"id"`
	Type    string          `json:"type"`
	Topic   string          `json:"topic"`
	Subject string          `json:"subject"`
	Data    json.RawMessage `json:"data"`
}

type l2Data struct {
	Changes struct {
		Asks [][]string `json:"asks"`
		Bids [][]string `json:"bids"`
	} `json:"changes"`
	SequenceEnd   int64  `json:"sequenceEnd"`
	SequenceStart int64  `json:"sequenceStart"`
	Symbol        string `json:"symbol"`
	Time          int64  `json:"time"`
}

type tradeData struct {
	Price    string `json:"price"`
	Sequence string `json:"sequence"`
	Side     string `json:"side"`
	Size     string `json:"size"`
	Symbol   string `json:"symbol"`
	Time     string `json:"time"`
}

func (c *Client) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if err := c.runOnce(ctx); err != nil && !errors.Is(err, context.Canceled) {
			fmt.Printf("kucoin reconnect after error: %v\n", err)
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(c.cfg.ReconnectDelay):
		}
	}
}

func (c *Client) runOnce(ctx context.Context) error {
	wsURL, pingInterval, err := c.getWSURL(ctx)
	if err != nil {
		return fmt.Errorf("get ws url: %w", err)
	}

	conn, _, err := websocket.DefaultDialer.DialContext(ctx, wsURL, nil)
	if err != nil {
		return fmt.Errorf("dial ws: %w", err)
	}
	defer conn.Close()

	if err := c.subscribe(conn); err != nil {
		return fmt.Errorf("subscribe: %w", err)
	}

	deltaBuf := make([]l2Data, 0, 256)

	snapDone := make(chan error, 1)
	go func() {
		snapDone <- c.loadInitialSnapshot(ctx)
	}()

	pingEvery := time.Duration(pingInterval) * time.Millisecond
	if pingEvery <= 0 {
		pingEvery = 18 * time.Second
	}

	errCh := make(chan error, 2)

	go func() {
		ticker := time.NewTicker(maxDuration(5*time.Second, pingEvery/2))
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				errCh <- ctx.Err()
				return
			case <-ticker.C:
				msg := map[string]any{
					"id":   strconv.FormatInt(time.Now().UnixNano(), 10),
					"type": "ping",
				}
				if err := conn.WriteJSON(msg); err != nil {
					errCh <- err
					return
				}
			}
		}
	}()

	go func() {
		snapshotReady := false

		for {
			select {
			case <-ctx.Done():
				errCh <- ctx.Err()
				return
			default:
			}

			_ = conn.SetReadDeadline(time.Now().Add(c.cfg.WSReadTimeout))

			_, payload, err := conn.ReadMessage()
			if err != nil {
				errCh <- err
				return
			}

			var env wsEnvelope
			if err := json.Unmarshal(payload, &env); err != nil {
				continue
			}

			switch env.Type {
			case "welcome", "ack", "pong":
				continue
			case "message":
			default:
				continue
			}

			if strings.HasPrefix(env.Topic, "/market/level2:") {
				var d l2Data
				if err := json.Unmarshal(env.Data, &d); err != nil {
					continue
				}

				if !snapshotReady {
					deltaBuf = append(deltaBuf, d)

					select {
					case err := <-snapDone:
						if err != nil {
							errCh <- err
							return
						}
						snapshotReady = true
						deltaBuf = c.playbackBufferedDeltas(ctx, deltaBuf)
					default:
					}
					continue
				}

				ok := c.book.ApplyDelta(d.SequenceStart, d.SequenceEnd, d.Changes.Bids, d.Changes.Asks)
				if !ok {
					if err := c.loadInitialSnapshot(ctx); err != nil {
						errCh <- fmt.Errorf("reload snapshot after gap: %w", err)
						return
					}
					continue
				}

				c.emitBook(c.book.Snapshot(c.cfg.TopLevels))
				continue
			}

			if strings.HasPrefix(env.Topic, "/market/match:") {
				var d tradeData
				if err := json.Unmarshal(env.Data, &d); err != nil {
					continue
				}

				price, err1 := strconv.ParseFloat(d.Price, 64)
				qty, err2 := strconv.ParseFloat(d.Size, 64)
				if err1 != nil || err2 != nil {
					continue
				}

				ts := time.Now()
				if ms, err := strconv.ParseInt(d.Time, 10, 64); err == nil {
					switch {
					case ms > 1e15:
						ts = time.Unix(0, ms)
					case ms > 1e12:
						ts = time.UnixMilli(ms)
					case ms > 1e9:
						ts = time.Unix(ms, 0)
					}
				}

				side := types.SideUnknown
				switch strings.ToLower(d.Side) {
				case "buy":
					side = types.SideBuy
				case "sell":
					side = types.SideSell
				}

				c.emitTrade(types.TradeTick{
					Symbol:    d.Symbol,
					Price:     price,
					Qty:       qty,
					Side:      side,
					Timestamp: ts,
				})
			}
		}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errCh:
		return err
	}
}

func (c *Client) emitBook(book types.BookSnapshot) {
	select {
	case c.events <- core.Event{Type: core.EventBook, Book: &book}:
	default:
	}
}

func (c *Client) emitTrade(tr types.TradeTick) {
	select {
	case c.events <- core.Event{Type: core.EventTrade, Trade: &tr}:
	default:
	}
}

func (c *Client) playbackBufferedDeltas(ctx context.Context, deltas []l2Data) []l2Data {
	if len(deltas) == 0 {
		c.emitBook(c.book.Snapshot(c.cfg.TopLevels))
		return nil
	}

	currentSeq := c.book.Sequence()
	remaining := deltas[:0]

	for _, d := range deltas {
		select {
		case <-ctx.Done():
			return remaining
		default:
		}

		if d.SequenceEnd <= currentSeq {
			continue
		}

		ok := c.book.ApplyDelta(d.SequenceStart, d.SequenceEnd, d.Changes.Bids, d.Changes.Asks)
		if !ok {
			remaining = append(remaining, d)
			break
		}
		currentSeq = d.SequenceEnd
	}

	c.emitBook(c.book.Snapshot(c.cfg.TopLevels))
	return remaining
}

func (c *Client) getWSURL(ctx context.Context) (string, int64, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.kucoin.com/api/v1/bullet-public", nil)
	if err != nil {
		return "", 0, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	var br bulletResponse
	if err := json.NewDecoder(resp.Body).Decode(&br); err != nil {
		return "", 0, err
	}
	if br.Code != "200000" {
		return "", 0, fmt.Errorf("unexpected code: %s", br.Code)
	}
	if len(br.Data.InstanceServers) == 0 {
		return "", 0, errors.New("no instanceServers")
	}

	srv := br.Data.InstanceServers[0]
	u, err := url.Parse(srv.Endpoint)
	if err != nil {
		return "", 0, err
	}

	q := u.Query()
	q.Set("token", br.Data.Token)
	q.Set("connectId", strconv.FormatInt(time.Now().UnixNano(), 10))
	u.RawQuery = q.Encode()

	return u.String(), srv.PingInterval, nil
}

func (c *Client) loadInitialSnapshot(ctx context.Context) error {
	u := "https://api.kucoin.com/api/v3/market/orderbook/level2?symbol=" + url.QueryEscape(c.symbol)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var ob orderBookResponse
	if err := json.NewDecoder(resp.Body).Decode(&ob); err != nil {
		return err
	}
	if ob.Code != "200000" {
		return fmt.Errorf("unexpected code: %s", ob.Code)
	}

	seq, err := ob.Data.Sequence.Int64()
	if err != nil {
		return err
	}

	return c.book.LoadSnapshot(seq, ob.Data.Bids, ob.Data.Asks)
}

func (c *Client) subscribe(conn *websocket.Conn) error {
	topics := []string{
		"/market/level2:" + c.symbol,
		"/market/match:" + c.symbol,
	}

	for _, topic := range topics {
		msg := map[string]any{
			"id":       strconv.FormatInt(time.Now().UnixNano(), 10),
			"type":     "subscribe",
			"topic":    topic,
			"response": true,
		}
		if err := conn.WriteJSON(msg); err != nil {
			return err
		}
	}

	return nil
}

func maxDuration(a, b time.Duration) time.Duration {
	if a > b {
		return a
	}
	return b
}
Как запускать
go mod tidy
go run ./cmd/trader
Что здесь уже исправлено по сравнению с прежним вариантом

нет параллельной записи в features.Engine

book и trade идут через один event loop

trade flow стал time-based, а не message-based

окна Return1s и Vol10s выровнены честнее

есть проверка готовности окон

Ising встроен как верхний слой подтверждения

структура уже годится для дальнейшего перехода к backtest и live execution

Что я бы делал следующим шагом

Сразу после этого стека я бы добавил:

Makefile

unit-тесты для calcTopLevelOFI, regime, strategy

отдельный backtest runner по CSV

online stats: winrate, pnl, avg hold time, avg adverse excursion

Если хочешь, следующим сообщением я дам ещё и полный стек с Makefile, .gitignore, README и тестами, чтобы проект уже был совсем как нормальный репозиторий.
