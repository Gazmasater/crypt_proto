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



git pull --rebase origin privat
git push origin privat


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




go run -race main.go


GOMAXPROCS=8 go run -race main.go




package engine

import (
	"container/ring"
	"crypt_proto/internal/aring"
	"fmt"
	"math"
	"sort"
)

type Candle = struct {
	Ts    int64
	Open  float64
	High  float64
	Low   float64
	Close float64
	Vol   float64
}

type OIPoint = struct {
	Ts int64
	OI float64
}

type State int

const (
	FLAT State = iota
	WAIT_RETEST_LONG
	WAIT_RETEST_SHORT
	WAIT_HL_LONG
	WAIT_LH_SHORT
)

type Engine struct {
	Symbol string

	R1h  *aring.Ring[Candle]
	R15  *aring.Ring[Candle]
	R5   *aring.Ring[Candle]
	OI15 *aring.Ring[OIPoint]

	// Levels
	Rhigh, Rlow float64 // 1H 24h
	H15, L15    float64 // 15m 8h
	Mid15       float64

	DOI30 float64

	State State

	Level      float64
	SeenRetest bool
	RetestLow  float64
	RetestHigh float64
	Trigger    float64
}

func New(symbol string) *Engine {
	return &Engine{
		Symbol: symbol,
		R1h:    ring.New,
		R15:    ring.New,
		R5:     ring.New,
		OI15:   ring.New,
		State:  FLAT,
	}
}

func (e *Engine) Warmup(c1h, c15, c5 []Candle, oi []OIPoint) {
	for _, c := range c1h {
		e.R1h.Push(c)
	}
	for _, c := range c15 {
		e.R15.Push(c)
	}
	for _, c := range c5 {
		e.R5.Push(c)
	}
	for _, p := range oi {
		e.OI15.Push(p)
	}
	e.recalc()
}

func (e *Engine) OnClose1H(c Candle)  { e.R1h.Push(c) }
func (e *Engine) OnClose15m(c Candle) { e.R15.Push(c); e.recalc(); e.on15mClose(c) }
func (e *Engine) OnClose5m(c Candle)  { e.R5.Push(c); e.on5mClose(c) }
func (e *Engine) OnOI15m(p OIPoint)   { e.OI15.Push(p); e.recalcOI() }

func (e *Engine) recalc() {
	e.recalcLevels()
	e.recalcOI()
}

func (e *Engine) recalcLevels() {
	// 1H last 24 candles
	c1 := e.R1h.Snapshot()
	if len(c1) > 24 {
		c1 = c1[len(c1)-24:]
	}
	if len(c1) > 0 {
		e.Rhigh, e.Rlow = rangeHL(c1)
	}

	// 15m last 32 candles (8h)
	c15 := e.R15.Snapshot()
	if len(c15) > 32 {
		c15 = c15[len(c15)-32:]
	}
	if len(c15) > 0 {
		e.H15, e.L15 = rangeHL(c15)
		e.Mid15 = (e.H15 + e.L15) / 2
	}
}

func (e *Engine) recalcOI() {
	pts := e.OI15.Snapshot()
	if len(pts) < 3 {
		e.DOI30 = 0
		return
	}
	// ensure time order
	sort.Slice(pts, func(i, j int) bool { return pts[i].Ts < pts[j].Ts })
	last := pts[len(pts)-1]
	prev2 := pts[len(pts)-3]
	e.DOI30 = last.OI - prev2.OI
}

func rangeHL(cs []Candle) (hi, lo float64) {
	hi = -math.MaxFloat64
	lo = math.MaxFloat64
	for _, c := range cs {
		if c.High > hi {
			hi = c.High
		}
		if c.Low < lo {
			lo = c.Low
		}
	}
	return
}

func (e *Engine) lastPrice() float64 {
	c5 := e.R5.Snapshot()
	if len(c5) > 0 {
		return c5[len(c5)-1].Close
	}
	c15 := e.R15.Snapshot()
	if len(c15) > 0 {
		return c15[len(c15)-1].Close
	}
	c1 := e.R1h.Snapshot()
	if len(c1) > 0 {
		return c1[len(c1)-1].Close
	}
	return 0
}

func (e *Engine) Bias() string {
	price := e.lastPrice()
	if price == 0 || e.Rhigh == 0 || e.H15 == 0 {
		return "FLAT"
	}

	distToRes := (e.Rhigh - price) / price * 100
	distToSup := (price - e.Rlow) / price * 100
	const nearPct = 0.25

	longOK := distToRes >= nearPct && e.DOI30 > 0 && price > e.Mid15
	shortOK := distToSup >= nearPct && e.DOI30 < 0 && price < e.Mid15

	if longOK && !shortOK {
		return "LONG"
	}
	if shortOK && !longOK {
		return "SHORT"
	}
	return "FLAT"
}

func (e *Engine) on15mClose(c Candle) {
	const buf = 0.0005 // 0.05%
	breakUp := c.Close > e.H15*(1+buf)
	breakDn := c.Close < e.L15*(1-buf)

	b := e.Bias()
	fmt.Printf("[15m] close=%.2f H15=%.2f L15=%.2f ΔOI30=%.0f bias=%s state=%d\n",
		c.Close, e.H15, e.L15, e.DOI30, b, e.State)

	if e.State != FLAT {
		return
	}

	if breakUp && b == "LONG" {
		e.State = WAIT_RETEST_LONG
		e.Level = e.H15
		e.SeenRetest = false
		e.RetestLow = math.MaxFloat64
		e.RetestHigh = -math.MaxFloat64
		e.Trigger = 0
		fmt.Printf("SETUP: BREAKUP -> WAIT_RETEST_LONG level=%.2f\n", e.Level)
		return
	}

	if breakDn && b == "SHORT" {
		e.State = WAIT_RETEST_SHORT
		e.Level = e.L15
		e.SeenRetest = false
		e.RetestLow = math.MaxFloat64
		e.RetestHigh = -math.MaxFloat64
		e.Trigger = 0
		fmt.Printf("SETUP: BREAKDOWN -> WAIT_RETEST_SHORT level=%.2f\n", e.Level)
		return
	}
}

func (e *Engine) on5mClose(c Candle) {
	switch e.State {
	case WAIT_RETEST_LONG:
		e.handleRetestLong(c)
	case WAIT_RETEST_SHORT:
		e.handleRetestShort(c)
	case WAIT_HL_LONG:
		e.handleHLTriggerLong(c)
	case WAIT_LH_SHORT:
		e.handleLHTriggerShort(c)
	}
}

func (e *Engine) handleRetestLong(c Candle) {
	const buf = 0.0007
	levelHi := e.Level * (1 + buf)
	levelLo := e.Level * (1 - buf)

	touched := c.Low <= levelHi && c.High >= levelLo
	if touched {
		e.SeenRetest = true
		if c.Low < e.RetestLow {
			e.RetestLow = c.Low
		}
		if c.High > e.RetestHigh {
			e.RetestHigh = c.High
		}
	}

	// ретест ок: после касания закрылись выше уровня
	if e.SeenRetest && c.Close > e.Level {
		e.State = WAIT_HL_LONG
		fmt.Printf("RETEST OK -> WAIT_HL_LONG retestLow=%.2f\n", e.RetestLow)
	}
}

func (e *Engine) handleRetestShort(c Candle) {
	const buf = 0.0007
	levelHi := e.Level * (1 + buf)
	levelLo := e.Level * (1 - buf)

	touched := c.Low <= levelHi && c.High >= levelLo
	if touched {
		e.SeenRetest = true
		if c.High > e.RetestHigh {
			e.RetestHigh = c.High
		}
		if c.Low < e.RetestLow {
			e.RetestLow = c.Low
		}
	}

	if e.SeenRetest && c.Close < e.Level {
		e.State = WAIT_LH_SHORT
		fmt.Printf("RETEST OK -> WAIT_LH_SHORT retestHigh=%.2f\n", e.RetestHigh)
	}
}

func (e *Engine) handleHLTriggerLong(c Candle) {
	// упрощённо: после ретеста ждём рост и формируем trigger как max(high) в восстановлении
	if c.Close > e.Level {
		if c.High > e.Trigger {
			e.Trigger = c.High
		}
	}

	// вход: пробой trigger (по close)
	if e.Trigger > 0 && c.Close > e.Trigger {
		stop := e.RetestLow
		fmt.Printf("SIGNAL: ENTER_LONG price=%.2f stop=%.2f level=%.2f\n", c.Close, stop, e.Level)
		e.reset()
	}
}

func (e *Engine) handleLHTriggerShort(c Candle) {
	// упрощённо: trigger как min(low) после отката
	if c.Close < e.Level {
		if e.Trigger == 0 || c.Low < e.Trigger {
			e.Trigger = c.Low
		}
	}

	if e.Trigger > 0 && c.Close < e.Trigger {
		stop := e.RetestHigh
		fmt.Printf("SIGNAL: ENTER_SHORT price=%.2f stop=%.2f level=%.2f\n", c.Close, stop, e.Level)
		e.reset()
	}
}

func (e *Engine) reset() {
	e.State = FLAT
	e.Level = 0
	e.SeenRetest = false
	e.RetestLow = 0
	e.RetestHigh = 0
	e.Trigger = 0
}


[{
	"resource": "/home/gaz358/myprog/crypt_proto/internal/engine/engine.go",
	"owner": "_generated_diagnostic_collection_name_#0",
	"code": {
		"value": "IncompatibleAssign",
		"target": {
			"$mid": 1,
			"path": "/golang.org/x/tools/internal/typesinternal",
			"scheme": "https",
			"authority": "pkg.go.dev",
			"fragment": "IncompatibleAssign"
		}
	},
	"severity": 8,
	"message": "cannot use ring.New (value of type func(n int) *ring.Ring) as *aring.Ring[Candle] value in struct literal",
	"source": "compiler",
	"startLineNumber": 62,
	"startColumn": 11,
	"endLineNumber": 62,
	"endColumn": 19,
	"modelVersionId": 102,
	"origin": "extHost1"
}]




