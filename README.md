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




Структура проекта
trade_f/
  go.mod
  main.go
  internal/
    ring/ring.go
    kucoin/kucoin.go
    engine/engine.go


Ниже — готовые файлы. Скопируй 1-в-1.

1) go.mod
module trade_f

go 1.22

2) internal/ring/ring.go
package ring

import "sync"

type Ring[T any] struct {
	mu    sync.RWMutex
	buf   []T
	size  int
	head  int
	count int
}

func New[T any](size int) *Ring[T] {
	return &Ring[T]{buf: make([]T, size), size: size}
}

func (r *Ring[T]) Push(v T) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.buf[r.head] = v
	r.head = (r.head + 1) % r.size
	if r.count < r.size {
		r.count++
	}
}

func (r *Ring[T]) Len() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.count
}

// Snapshot returns oldest->newest
func (r *Ring[T]) Snapshot() []T {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]T, r.count)
	if r.count == 0 {
		return out
	}
	start := (r.head - r.count + r.size) % r.size
	for i := 0; i < r.count; i++ {
		out[i] = r.buf[(start+i)%r.size]
	}
	return out
}

3) internal/kucoin/kucoin.go (REST: timestamp + klines + OI)
package kucoin

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	FuturesBase = "https://api-futures.kucoin.com"
	SpotBase    = "https://api.kucoin.com"
)

type Candle struct {
	Ts    int64   // close time ms
	Open  float64
	High  float64
	Low   float64
	Close float64
	Vol   float64
}

type OIPoint struct {
	Ts int64
	OI float64
}

type tsResp struct {
	Code string `json:"code"`
	Data int64  `json:"data"`
}

type klineResp struct {
	Code string  `json:"code"`
	Data [][]any `json:"data"` // [ts, open, close, high, low, vol, turnover]
}

type oiResp struct {
	Code string `json:"code"`
	Data []struct {
		OpenInterest string `json:"openInterest"`
		Ts           int64  `json:"ts"`
	} `json:"data"`
}

func ServerTimeMs() (int64, error) {
	resp, err := http.Get(SpotBase + "/api/v1/timestamp")
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	b, _ := io.ReadAll(resp.Body)
	var r tsResp
	if err := json.Unmarshal(b, &r); err != nil {
		return 0, err
	}
	if r.Code != "200000" {
		return 0, fmt.Errorf("timestamp bad code=%s body=%s", r.Code, string(b))
	}
	return r.Data, nil
}

func FetchKlines(symbol string, granularity string, fromMs, toMs int64) ([]Candle, error) {
	u, _ := url.Parse(FuturesBase + "/api/v1/kline/query")
	q := u.Query()
	q.Set("symbol", symbol)
	q.Set("granularity", granularity) // futures: 60=1H, 15=15m, 5=5m
	q.Set("from", strconv.FormatInt(fromMs, 10))
	q.Set("to", strconv.FormatInt(toMs, 10))
	u.RawQuery = q.Encode()

	req, _ := http.NewRequest("GET", u.String(), nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, _ := io.ReadAll(resp.Body)
	var r klineResp
	if err := json.Unmarshal(b, &r); err != nil {
		return nil, err
	}
	if r.Code != "200000" {
		return nil, fmt.Errorf("kline bad code=%s body=%s", r.Code, string(b))
	}

	out := make([]Candle, 0, len(r.Data))
	for _, row := range r.Data {
		if len(row) < 6 {
			continue
		}
		ts, _ := toInt64(row[0])
		open, _ := toFloat(row[1])
		closeV, _ := toFloat(row[2])
		high, _ := toFloat(row[3])
		low, _ := toFloat(row[4])
		vol, _ := toFloat(row[5])

		out = append(out, Candle{
			Ts: ts, Open: open, High: high, Low: low, Close: closeV, Vol: vol,
		})
	}
	return out, nil
}

func FetchOI15m(symbol string, startAt, endAt int64) ([]OIPoint, error) {
	u, _ := url.Parse(SpotBase + "/api/ua/v1/market/open-interest")
	q := u.Query()
	q.Set("symbol", symbol)
	q.Set("interval", "15min")
	q.Set("startAt", strconv.FormatInt(startAt, 10))
	q.Set("endAt", strconv.FormatInt(endAt, 10))
	q.Set("pageSize", "200")
	u.RawQuery = q.Encode()

	req, _ := http.NewRequest("GET", u.String(), nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, _ := io.ReadAll(resp.Body)
	var r oiResp
	if err := json.Unmarshal(b, &r); err != nil {
		return nil, err
	}
	if r.Code != "200000" {
		return nil, fmt.Errorf("oi bad code=%s body=%s", r.Code, string(b))
	}

	out := make([]OIPoint, 0, len(r.Data))
	for _, it := range r.Data {
		v, err := strconv.ParseFloat(it.OpenInterest, 64)
		if err != nil {
			continue
		}
		out = append(out, OIPoint{Ts: it.Ts, OI: v})
	}
	return out, nil
}

func toFloat(v any) (float64, error) {
	switch t := v.(type) {
	case string:
		return strconv.ParseFloat(t, 64)
	case float64:
		return t, nil
	case json.Number:
		return t.Float64()
	default:
		return 0, fmt.Errorf("unexpected type %T", v)
	}
}

func toInt64(v any) (int64, error) {
	switch t := v.(type) {
	case string:
		return strconv.ParseInt(t, 10, 64)
	case float64:
		return int64(t), nil
	case json.Number:
		return t.Int64()
	default:
		return 0, fmt.Errorf("unexpected type %T", v)
	}
}

4) internal/engine/engine.go (ringbuffers + FSM ожидания сигнала)
package engine

import (
	"fmt"
	"math"
	"sort"

	"trade_f/internal/ring"
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

	R1h  *ring.Ring[Candle]
	R15  *ring.Ring[Candle]
	R5   *ring.Ring[Candle]
	OI15 *ring.Ring[OIPoint]

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

5) main.go (загрузка initial bars → затем “живой” цикл ожидания)

Сначала сделаем простой вариант: каждые 5 минут подтягиваем последние свечи REST’ом и если появилась новая закрытая свеча — пушим в engine. Это уже будет “живой бот”, который ждёт сигнал.

package main

import (
	"fmt"
	"sort"
	"time"

	"trade_f/internal/engine"
	"trade_f/internal/kucoin"
)

const symbol = "XBTUSDTM"

func main() {
	e := engine.New(symbol)

	// ===== initial warmup =====
	nowMs, err := kucoin.ServerTimeMs()
	if err != nil {
		panic(err)
	}

	// берём чуть с запасом и потом отсортируем
	c1h, _ := kucoin.FetchKlines(symbol, "60", nowMs-30*60*60*1000, nowMs)
	c15, _ := kucoin.FetchKlines(symbol, "15", nowMs-12*60*60*1000, nowMs)
	c5, _ := kucoin.FetchKlines(symbol, "5", nowMs-3*60*60*1000, nowMs)
	oi, _ := kucoin.FetchOI15m(symbol, nowMs-20*60*60*1000, nowMs)

	// sort by ts asc
	sort.Slice(c1h, func(i, j int) bool { return c1h[i].Ts < c1h[j].Ts })
	sort.Slice(c15, func(i, j int) bool { return c15[i].Ts < c15[j].Ts })
	sort.Slice(c5, func(i, j int) bool { return c5[i].Ts < c5[j].Ts })
	sort.Slice(oi, func(i, j int) bool { return oi[i].Ts < oi[j].Ts })

	// convert to engine types
	convC := func(in []kucoin.Candle) []engine.Candle {
		out := make([]engine.Candle, 0, len(in))
		for _, c := range in {
			out = append(out, engine.Candle(c))
		}
		return out
	}
	convOI := func(in []kucoin.OIPoint) []engine.OIPoint {
		out := make([]engine.OIPoint, 0, len(in))
		for _, p := range in {
			out = append(out, engine.OIPoint(p))
		}
		return out
	}

	e.Warmup(convC(c1h), convC(c15), convC(c5), convOI(oi))

	fmt.Println("engine started; waiting for signals...")

	// ===== live loop (REST polling) =====
	var last5Ts int64
	var last15Ts int64
	var last1hTs int64
	if s := e.R5.Snapshot(); len(s) > 0 {
		last5Ts = s[len(s)-1].Ts
	}
	if s := e.R15.Snapshot(); len(s) > 0 {
		last15Ts = s[len(s)-1].Ts
	}
	if s := e.R1h.Snapshot(); len(s) > 0 {
		last1hTs = s[len(s)-1].Ts
	}

	ticker := time.NewTicker(10 * time.Second) // часто, но запросы делай экономно
	defer ticker.Stop()

	for range ticker.C {
		nowMs, err := kucoin.ServerTimeMs()
		if err != nil {
			continue
		}

		// 5m: берём последние ~30 минут
		c5n, err := kucoin.FetchKlines(symbol, "5", nowMs-40*60*1000, nowMs)
		if err == nil && len(c5n) > 0 {
			sort.Slice(c5n, func(i, j int) bool { return c5n[i].Ts < c5n[j].Ts })
			for _, c := range c5n {
				if c.Ts > last5Ts {
					e.OnClose5m(engine.Candle(c))
					last5Ts = c.Ts
				}
			}
		}

		// 15m: последние ~4 часа
		c15n, err := kucoin.FetchKlines(symbol, "15", nowMs-5*60*60*1000, nowMs)
		if err == nil && len(c15n) > 0 {
			sort.Slice(c15n, func(i, j int) bool { return c15n[i].Ts < c15n[j].Ts })
			for _, c := range c15n {
				if c.Ts > last15Ts {
					e.OnClose15m(engine.Candle(c))
					last15Ts = c.Ts
				}
			}
		}

		// 1h: последние 30 часов
		c1n, err := kucoin.FetchKlines(symbol, "60", nowMs-35*60*60*1000, nowMs)
		if err == nil && len(c1n) > 0 {
			sort.Slice(c1n, func(i, j int) bool { return c1n[i].Ts < c1n[j].Ts })
			for _, c := range c1n {
				if c.Ts > last1hTs {
					e.OnClose1H(engine.Candle(c))
					last1hTs = c.Ts
				}
			}
		}

		// OI: раз в минуту достаточно (тут просто берём последние 2 часа и пушим новые точки)
		oin, err := kucoin.FetchOI15m(symbol, nowMs-2*60*60*1000, nowMs)
		if err == nil && len(oin) > 0 {
			sort.Slice(oin, func(i, j int) bool { return oin[i].Ts < oin[j].Ts })
			// пушим всё, ring сам ограничит; дубликаты можно отфильтровать позже
			for _, p := range oin {
				e.OnOI15m(engine.OIPoint(p))
			}
		}
	}
}





