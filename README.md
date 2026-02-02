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


package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
)

/*
KuCoin stat-arb (isolated, one file):

- Window: 120 x 1-min bars (2 hours) in ring buffers
- Warmup: REST candles (1min) for BTC-USDT and ETH-USDT
- Live: WS ticker -> mid -> minute bars (high/low/mean)

Signals (ASAP intra-minute):
- ENTER SHORT_SPREAD if z >= entryZ (sell BTC, buy ETH)
- ENTER LONG_SPREAD  if z <= -entryZ (buy BTC, sell ETH)
- EXIT when abs(z) <= exitZ

Time alignment fix:
- Stats computed using intersection of minute timestamps (MinuteMs) for BTC/ETH

Logging:
- Every 5 minutes (aligned to UTC boundary): ALWAYS prints either LOG_5M or LOG_5M_SKIP.
  Uses CLOSED-only window (120closed_aligned) to be stable.
- On ENTER/EXIT: logs immediately (uses intra-minute stats at that moment).
- Every minute (STRICT close for BOTH): prints PROFIT/LOSS (unrealized) and logs PNL_1M
  using CLOSED-only stats ("120closed_aligned").

PnL model (paper, log-space spread):
spread_t = ln(BTC_t) - beta * ln(ETH_t)
LONG_SPREAD  => pnl_log = +(spread_now - spread_entry)
SHORT_SPREAD => pnl_log = -(spread_now - spread_entry)
pnl_pct ≈ exp(pnl_log) - 1
*/

const (
	windowBars = 120
	tf         = "1min"
	logPath    = "statarb_5m.jsonl"
)

// ====================== Bars / Ring ======================

type Bar struct {
	MinuteMs int64
	High     float64
	Low      float64
	Mean     float64
	Count    int
}

type Ring struct {
	buf   []Bar
	cap   int
	head  int
	size  int
	mutex sync.Mutex
}

func NewRing(capacity int) *Ring { return &Ring{buf: make([]Bar, capacity), cap: capacity} }

func (r *Ring) Push(b Bar) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.buf[r.head] = b
	r.head = (r.head + 1) % r.cap
	if r.size < r.cap {
		r.size++
	}
}

func (r *Ring) Snapshot() []Bar {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	out := make([]Bar, 0, r.size)
	start := (r.head - r.size + r.cap) % r.cap
	for i := 0; i < r.size; i++ {
		out = append(out, r.buf[(start+i)%r.cap])
	}
	return out
}

func (r *Ring) Len() int {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	return r.size
}

// ====================== Ticks / Minute Aggregation ======================

type Tick struct {
	Symbol string
	TimeMs int64
	Bid    float64
	Ask    float64
}

func (t Tick) Mid() float64 {
	if t.Bid <= 0 || t.Ask <= 0 {
		return 0
	}
	return 0.5 * (t.Bid + t.Ask)
}

type MinuteAgg struct {
	symbol string
	curMin int64
	high   float64
	low    float64
	sum    float64
	cnt    int
}

func NewMinuteAgg(symbol string) *MinuteAgg { return &MinuteAgg{symbol: symbol} }

func minuteStartMs(tsMs int64) int64 { return (tsMs / 60000) * 60000 }

func (a *MinuteAgg) Update(t Tick) (Bar, bool) {
	mid := t.Mid()
	if mid <= 0 {
		return Bar{}, false
	}
	m := minuteStartMs(t.TimeMs)

	if a.curMin == 0 {
		a.curMin = m
		a.high, a.low = mid, mid
		a.sum, a.cnt = mid, 1
		return Bar{}, false
	}

	if m == a.curMin {
		if mid > a.high {
			a.high = mid
		}
		if mid < a.low {
			a.low = mid
		}
		a.sum += mid
		a.cnt++
		return Bar{}, false
	}

	closed := Bar{
		MinuteMs: a.curMin,
		High:     a.high,
		Low:      a.low,
		Mean:     a.sum / float64(a.cnt),
		Count:    a.cnt,
	}

	// start new minute
	a.curMin = m
	a.high, a.low = mid, mid
	a.sum, a.cnt = mid, 1
	return closed, true
}

func (a *MinuteAgg) Partial() (Bar, bool) {
	if a.curMin == 0 || a.cnt == 0 {
		return Bar{}, false
	}
	return Bar{
		MinuteMs: a.curMin,
		High:     a.high,
		Low:      a.low,
		Mean:     a.sum / float64(a.cnt),
		Count:    a.cnt,
	}, true
}

// ====================== REST Warmup ======================

type restKlinesResp struct {
	Code string     `json:"code"`
	Data [][]string `json:"data"`
}

type Kline struct {
	StartSec int64
	Open     float64
	Close    float64
	High     float64
	Low      float64
}

func fetchKlines(ctx context.Context, symbol string, startAtSec, endAtSec int64) ([]Kline, error) {
	u, _ := url.Parse("https://api.kucoin.com/api/v1/market/candles")
	q := u.Query()
	q.Set("symbol", symbol)
	q.Set("type", tf)
	q.Set("startAt", strconv.FormatInt(startAtSec, 10))
	q.Set("endAt", strconv.FormatInt(endAtSec, 10))
	u.RawQuery = q.Encode()

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	resp, err := (&http.Client{Timeout: 15 * time.Second}).Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		return nil, fmt.Errorf("rest http %d: %s", resp.StatusCode, string(b))
	}

	var r restKlinesResp
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}
	if r.Code != "200000" {
		return nil, fmt.Errorf("rest code=%s", r.Code)
	}

	out := make([]Kline, 0, len(r.Data))
	for _, row := range r.Data {
		if len(row) < 5 {
			continue
		}
		ts, _ := strconv.ParseInt(row[0], 10, 64)
		op, _ := strconv.ParseFloat(row[1], 64)
		cl, _ := strconv.ParseFloat(row[2], 64)
		hi, _ := strconv.ParseFloat(row[3], 64)
		lo, _ := strconv.ParseFloat(row[4], 64)
		if ts == 0 || op <= 0 || cl <= 0 || hi <= 0 || lo <= 0 {
			continue
		}
		out = append(out, Kline{StartSec: ts, Open: op, Close: cl, High: hi, Low: lo})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].StartSec < out[j].StartSec })
	return out, nil
}

// ====================== KuCoin WS ======================

type bulletResp struct {
	Code string `json:"code"`
	Data struct {
		Token           string `json:"token"`
		InstanceServers []struct {
			Endpoint     string `json:"endpoint"`
			PingInterval int    `json:"pingInterval"`
			PingTimeout  int    `json:"pingTimeout"`
		} `json:"instanceServers"`
	} `json:"data"`
}

func getWSEndpoint(ctx context.Context) (endpoint string, token string, pingInterval time.Duration, err error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.kucoin.com/api/v1/bullet-public", nil)
	resp, err := (&http.Client{Timeout: 15 * time.Second}).Do(req)
	if err != nil {
		return "", "", 0, err
	}
	defer resp.Body.Close()

	var r bulletResp
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return "", "", 0, err
	}
	if r.Code != "200000" || len(r.Data.InstanceServers) == 0 {
		return "", "", 0, fmt.Errorf("bullet code=%s servers=%d", r.Code, len(r.Data.InstanceServers))
	}

	endpoint = r.Data.InstanceServers[0].Endpoint
	token = r.Data.Token
	pi := r.Data.InstanceServers[0].PingInterval
	if pi <= 0 {
		pi = 18000
	}
	pingInterval = time.Duration(pi) * time.Millisecond
	return endpoint, token, pingInterval, nil
}

type wsEnvelope struct {
	Type    string          `json:"type"`
	Topic   string          `json:"topic"`
	Subject string          `json:"subject"`
	Data    json.RawMessage `json:"data"`
}

type wsTickerData struct {
	BestAsk string `json:"bestAsk"`
	BestBid string `json:"bestBid"`
	Time    int64  `json:"time"` // ms
	Time2   int64  `json:"Time"` // ms (fallback)
}

func parseTopicSymbol(topic string) string {
	const p = "/market/ticker:"
	if len(topic) >= len(p) && topic[:len(p)] == p {
		return topic[len(p):]
	}
	return ""
}

func runWS(ctx context.Context, tickOut chan<- Tick) error {
	endpoint, token, pingEvery, err := getWSEndpoint(ctx)
	if err != nil {
		return err
	}

	connectId := fmt.Sprintf("statarb-%d", time.Now().UnixNano())
	u, _ := url.Parse(endpoint)
	q := u.Query()
	q.Set("token", token)
	q.Set("connectId", connectId)
	u.RawQuery = q.Encode()

	c, _, err := websocket.DefaultDialer.DialContext(ctx, u.String(), nil)
	if err != nil {
		return err
	}
	defer c.Close()

	sub := func(topic string) error {
		msg := map[string]any{
			"id":       fmt.Sprintf("%d", time.Now().UnixNano()),
			"type":     "subscribe",
			"topic":    topic,
			"response": true,
		}
		return c.WriteJSON(msg)
	}
	if err := sub("/market/ticker:BTC-USDT"); err != nil {
		return err
	}
	if err := sub("/market/ticker:ETH-USDT"); err != nil {
		return err
	}

	done := make(chan struct{})
	go func() {
		t := time.NewTicker(pingEvery)
		defer t.Stop()
		for {
			select {
			case <-t.C:
				_ = c.WriteJSON(map[string]any{
					"id":   fmt.Sprintf("%d", time.Now().UnixNano()),
					"type": "ping",
				})
			case <-done:
				return
			case <-ctx.Done():
				return
			}
		}
	}()

	for {
		_, b, err := c.ReadMessage()
		if err != nil {
			close(done)
			return err
		}

		var env wsEnvelope
		if err := json.Unmarshal(b, &env); err != nil {
			continue
		}
		if env.Type != "message" || env.Subject != "trade.ticker" {
			continue
		}

		sym := parseTopicSymbol(env.Topic)
		if sym != "BTC-USDT" && sym != "ETH-USDT" {
			continue
		}

		var d wsTickerData
		if err := json.Unmarshal(env.Data, &d); err != nil {
			continue
		}

		ts := d.Time
		if ts <= 0 {
			ts = d.Time2
		}

		bid, _ := strconv.ParseFloat(d.BestBid, 64)
		ask, _ := strconv.ParseFloat(d.BestAsk, 64)
		if ts <= 0 || bid <= 0 || ask <= 0 {
			continue
		}

		select {
		case tickOut <- Tick{Symbol: sym, TimeMs: ts, Bid: bid, Ask: ask}:
		case <-ctx.Done():
			close(done)
			return ctx.Err()
		}
	}
}

// ====================== Stats helpers ======================

type Stats struct {
	Beta   float64
	Spread float64
	Mu     float64
	Sigma  float64
	Z      float64
	Mode   string
	MinMs  int64
}

func meanStd(a []float64) (mu, sd float64) {
	n := len(a)
	if n == 0 {
		return 0, 0
	}
	for _, v := range a {
		mu += v
	}
	mu /= float64(n)
	if n == 1 {
		return mu, 0
	}
	var ss float64
	for _, v := range a {
		d := v - mu
		ss += d * d
	}
	sd = math.Sqrt(ss / float64(n-1))
	return mu, sd
}

func olsBeta(x, y []float64) float64 {
	n := len(x)
	if n == 0 || n != len(y) {
		return 1.0
	}
	var sx, sy float64
	for i := 0; i < n; i++ {
		sx += x[i]
		sy += y[i]
	}
	mx := sx / float64(n)
	my := sy / float64(n)

	var cov, vx float64
	for i := 0; i < n; i++ {
		dx := x[i] - mx
		dy := y[i] - my
		cov += dx * dy
		vx += dx * dx
	}
	if vx == 0 {
		return 1.0
	}
	return cov / vx
}

func calcStatsFromPrices120(btc120, eth120 []float64) (beta, spread, mu, sigma, z float64, ok bool) {
	if len(btc120) != 120 || len(eth120) != 120 {
		return 0, 0, 0, 0, 0, false
	}
	lbtc := make([]float64, 120)
	leth := make([]float64, 120)
	for i := 0; i < 120; i++ {
		if btc120[i] <= 0 || eth120[i] <= 0 {
			return 0, 0, 0, 0, 0, false
		}
		lbtc[i] = math.Log(btc120[i])
		leth[i] = math.Log(eth120[i])
	}

	beta = olsBeta(leth, lbtc)

	spreads := make([]float64, 120)
	for i := 0; i < 120; i++ {
		spreads[i] = lbtc[i] - beta*leth[i]
	}
	mu, sigma = meanStd(spreads)

	spread = spreads[119]
	z = 0.0
	if sigma > 0 {
		z = (spread - mu) / sigma
	}
	return beta, spread, mu, sigma, z, true
}

func buildAligned120(rBTC, rETH *Ring, aggBTC, aggETH *MinuteAgg) (btc120, eth120 []float64, mode string, minMs int64, ok bool) {
	sBTC := rBTC.Snapshot()
	sETH := rETH.Snapshot()

	mB := make(map[int64]float64, len(sBTC))
	for _, b := range sBTC {
		mB[b.MinuteMs] = b.Mean
	}
	mE := make(map[int64]float64, len(sETH))
	for _, b := range sETH {
		mE[b.MinuteMs] = b.Mean
	}

	common := make([]int64, 0, 256)
	for ms := range mB {
		if _, ok := mE[ms]; ok {
			common = append(common, ms)
		}
	}
	sort.Slice(common, func(i, j int) bool { return common[i] < common[j] })

	pbB, okB := aggBTC.Partial()
	pbE, okE := aggETH.Partial()
	usePartial := okB && okE && pbB.MinuteMs == pbE.MinuteMs && pbB.Mean > 0 && pbE.Mean > 0

	if usePartial {
		if len(common) < 119 {
			return nil, nil, "", 0, false
		}
		common = common[len(common)-119:]
		btc120 = make([]float64, 0, 120)
		eth120 = make([]float64, 0, 120)
		for _, ms := range common {
			btc120 = append(btc120, mB[ms])
			eth120 = append(eth120, mE[ms])
		}
		btc120 = append(btc120, pbB.Mean)
		eth120 = append(eth120, pbE.Mean)
		return btc120, eth120, "119+partial_aligned", pbB.MinuteMs, true
	}

	if len(common) < 120 {
		return nil, nil, "", 0, false
	}
	common = common[len(common)-120:]
	btc120 = make([]float64, 0, 120)
	eth120 = make([]float64, 0, 120)
	for _, ms := range common {
		btc120 = append(btc120, mB[ms])
		eth120 = append(eth120, mE[ms])
	}
	return btc120, eth120, "120closed_aligned", common[len(common)-1], true
}

func buildAligned120Closed(rBTC, rETH *Ring) (btc120, eth120 []float64, minMs int64, ok bool) {
	sBTC := rBTC.Snapshot()
	sETH := rETH.Snapshot()

	mB := make(map[int64]float64, len(sBTC))
	for _, b := range sBTC {
		mB[b.MinuteMs] = b.Mean
	}
	mE := make(map[int64]float64, len(sETH))
	for _, b := range sETH {
		mE[b.MinuteMs] = b.Mean
	}

	common := make([]int64, 0, 256)
	for ms := range mB {
		if _, ok := mE[ms]; ok {
			common = append(common, ms)
		}
	}
	sort.Slice(common, func(i, j int) bool { return common[i] < common[j] })

	if len(common) < 120 {
		return nil, nil, 0, false
	}
	common = common[len(common)-120:]

	btc120 = make([]float64, 0, 120)
	eth120 = make([]float64, 0, 120)
	for _, ms := range common {
		btc120 = append(btc120, mB[ms])
		eth120 = append(eth120, mE[ms])
	}
	return btc120, eth120, common[len(common)-1], true
}

// ====================== Aligned ticker ======================

func startAlignedTicker(ctx context.Context, period time.Duration) <-chan time.Time {
	ch := make(chan time.Time, 1)
	go func() {
		defer close(ch)

		now := time.Now().UTC()
		next := now.Truncate(period).Add(period)

		timer := time.NewTimer(time.Until(next))
		defer timer.Stop()

		select {
		case <-timer.C:
		case <-ctx.Done():
			return
		}

		// first aligned tick
		select {
		case ch <- time.Now().UTC():
		default:
		}

		ticker := time.NewTicker(period)
		defer ticker.Stop()

		for {
			select {
			case t := <-ticker.C:
				select {
				case ch <- t.UTC():
				default:
				}
			case <-ctx.Done():
				return
			}
		}
	}()
	return ch
}

// ====================== JSONL logger ======================

type LogRow struct {
	AtUTC string `json:"at_utc"`

	Beta   float64 `json:"beta"`
	Spread float64 `json:"spread"`
	Mu     float64 `json:"mu"`
	Sigma  float64 `json:"sigma"`
	Z      float64 `json:"z"`
	Mode   string  `json:"mode"`

	Pos    string  `json:"pos"`
	Action string  `json:"action"` // "LOG_5M", "LOG_5M_SKIP", "ENTER_LONG", "ENTER_SHORT", "EXIT", "PNL_1M"
	PnLPct float64 `json:"pnl_pct,omitempty"`
}

type JSONLLogger struct {
	mu  sync.Mutex
	f   *os.File
	enc *json.Encoder
}

func NewJSONLLogger(path string) (*JSONLLogger, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	return &JSONLLogger{f: f, enc: json.NewEncoder(f)}, nil
}
func (l *JSONLLogger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.f.Close()
}
func (l *JSONLLogger) Write(row LogRow) {
	l.mu.Lock()
	defer l.mu.Unlock()
	_ = l.enc.Encode(row)
}

func writeLog(logg *JSONLLogger, action string, pos string, st Stats, pnlPct *float64) {
	row := LogRow{
		AtUTC:  time.Now().UTC().Format(time.RFC3339),
		Beta:   st.Beta,
		Spread: st.Spread,
		Mu:     st.Mu,
		Sigma:  st.Sigma,
		Z:      st.Z,
		Mode:   st.Mode,
		Pos:    pos,
		Action: action,
	}
	if pnlPct != nil {
		row.PnLPct = *pnlPct
	}
	logg.Write(row)
}

// ====================== Trading state + PnL ======================

type Position int

const (
	Flat Position = iota
	LongSpread
	ShortSpread
)

func (p Position) String() string {
	switch p {
	case LongSpread:
		return "LONG_SPREAD"
	case ShortSpread:
		return "SHORT_SPREAD"
	default:
		return "FLAT"
	}
}

type HoldState struct {
	longSince  int64
	shortSince int64
}

type TradeState struct {
	Pos         Position
	EntrySpread float64
	EntryAtMs   int64
}

func pnlPctFromSpreads(pos Position, entrySpread, curSpread float64) float64 {
	d := curSpread - entrySpread
	if pos == ShortSpread {
		d = -d
	}
	return math.Exp(d) - 1.0
}

// ====================== Main ======================

func main() {
	// ---- params ----
	entryZ := 2.2
	exitZ := 0.8
	freshMs := int64(300)        // require BTC/ETH tick timestamps within 300ms
	holdMs := int64(1000)        // condition must hold for 1000ms before entering (0 disables)
	calcThrottleMs := int64(150) // don’t calc faster than every 150ms on ticks
	// --------------

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// graceful stop
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	go func() { <-sig; cancel() }()

	// rings
	rBTC := NewRing(windowBars)
	rETH := NewRing(windowBars)

	// warmup history (a bit more than 120 mins)
	nowSec := time.Now().Unix()
	startSec := nowSec - int64((windowBars+30)*60)

	btcK, err := fetchKlines(ctx, "BTC-USDT", startSec, nowSec)
	if err != nil {
		fmt.Println("history BTC error:", err)
		return
	}
	ethK, err := fetchKlines(ctx, "ETH-USDT", startSec, nowSec)
	if err != nil {
		fmt.Println("history ETH error:", err)
		return
	}

	// align warmup by startSec
	mBTC := make(map[int64]Kline, len(btcK))
	for _, k := range btcK {
		mBTC[k.StartSec] = k
	}
	mETH := make(map[int64]Kline, len(ethK))
	for _, k := range ethK {
		mETH[k.StartSec] = k
	}

	secs := make([]int64, 0, 256)
	for s := range mBTC {
		if _, ok := mETH[s]; ok {
			secs = append(secs, s)
		}
	}
	sort.Slice(secs, func(i, j int) bool { return secs[i] < secs[j] })
	if len(secs) > windowBars {
		secs = secs[len(secs)-windowBars:]
	}
	for _, s := range secs {
		kb := mBTC[s]
		ke := mETH[s]
		rBTC.Push(Bar{MinuteMs: s * 1000, High: kb.High, Low: kb.Low, Mean: 0.5 * (kb.High + kb.Low), Count: 1})
		rETH.Push(Bar{MinuteMs: s * 1000, High: ke.High, Low: ke.Low, Mean: 0.5 * (ke.High + ke.Low), Count: 1})
	}
	fmt.Printf("Warmup done: BTC=%d bars, ETH=%d bars | entryZ=%.2f exitZ=%.2f holdMs=%d freshMs=%d\n",
		rBTC.Len(), rETH.Len(), entryZ, exitZ, holdMs, freshMs)

	// logger
	logg, err := NewJSONLLogger(logPath)
	if err != nil {
		fmt.Println("logger error:", err)
		return
	}
	defer logg.Close()
	fmt.Println("Logging: 5m(aligned, closed-only) + ENTER/EXIT + PNL_1M(closed-only) to:", logPath)

	// aligned 5m ticker
	log5mCh := startAlignedTicker(ctx, 5*time.Minute)

	// WS ticks
	ticks := make(chan Tick, 20000)
	go func() {
		for {
			if ctx.Err() != nil {
				return
			}
			if err := runWS(ctx, ticks); err != nil && ctx.Err() == nil {
				fmt.Println("ws error:", err, "reconnect in 1s")
				time.Sleep(1 * time.Second)
			}
		}
	}()

	aggBTC := NewMinuteAgg("BTC-USDT")
	aggETH := NewMinuteAgg("ETH-USDT")

	// last tick times for freshness
	var lastBTCms int64
	var lastETHms int64

	// state
	hold := HoldState{}
	trade := TradeState{Pos: Flat}
	var lastCalcAtMs int64

	// pending closed bars (to detect same-minute close)
	pendingBTC := make(map[int64]Bar)
	pendingETH := make(map[int64]Bar)

	// compute stats (intra-minute) with throttle + freshness
	computeIntra := func(nowMs int64) (Stats, bool) {
		if nowMs-lastCalcAtMs < calcThrottleMs {
			return Stats{}, false
		}
		lastCalcAtMs = nowMs

		if lastBTCms == 0 || lastETHms == 0 {
			return Stats{}, false
		}
		d := lastBTCms - lastETHms
		if d < 0 {
			d = -d
		}
		if d > freshMs {
			return Stats{}, false
		}

		btc120, eth120, mode, minMs, ok := buildAligned120(rBTC, rETH, aggBTC, aggETH)
		if !ok {
			return Stats{}, false
		}
		beta, spread, mu, sigma, z, ok := calcStatsFromPrices120(btc120, eth120)
		if !ok {
			return Stats{}, false
		}
		return Stats{Beta: beta, Spread: spread, Mu: mu, Sigma: sigma, Z: z, Mode: mode, MinMs: minMs}, true
	}

	// compute stats strictly on closed minutes (for 1m PnL + 5m logs)
	computeClosed := func() (Stats, bool) {
		btc120, eth120, minMs, ok := buildAligned120Closed(rBTC, rETH)
		if !ok {
			return Stats{}, false
		}
		beta, spread, mu, sigma, z, ok := calcStatsFromPrices120(btc120, eth120)
		if !ok {
			return Stats{}, false
		}
		return Stats{Beta: beta, Spread: spread, Mu: mu, Sigma: sigma, Z: z, Mode: "120closed_aligned", MinMs: minMs}, true
	}

	printSignal := func(tag string, st Stats) {
		t := time.Now().UTC().Format("2006-01-02 15:04:05")
		min := time.UnixMilli(st.MinMs).UTC().Format("2006-01-02 15:04")
		fmt.Printf("[%s] %s | minute=%s | z=%+.3f beta=%.4f spread=%.6f mode=%s pos=%s\n",
			t, tag, min, st.Z, st.Beta, st.Spread, st.Mode, trade.Pos.String())
	}

	printPnLClosed := func(closedMinuteMs int64) {
		if trade.Pos == Flat {
			return
		}
		st, ok := computeClosed()
		if !ok {
			return
		}
		p := pnlPctFromSpreads(trade.Pos, trade.EntrySpread, st.Spread)
		min := time.UnixMilli(closedMinuteMs).UTC().Format("2006-01-02 15:04")
		sign := "PROFIT"
		if p < 0 {
			sign = "LOSS"
		}
		fmt.Printf("[1m %s] %.4f%% | minute=%s | pos=%s | z=%+.3f mode=%s\n",
			sign, p*100.0, min, trade.Pos.String(), st.Z, st.Mode)
		writeLog(logg, "PNL_1M", trade.Pos.String(), st, &p)
	}

	for {
		select {
		case <-ctx.Done():
			return

		// 5m log: ALWAYS print something
		case <-log5mCh:
			st, ok := computeClosed()
			if !ok {
				fmt.Println("[5m log] SKIP: not enough closed aligned bars yet")
				// optional: also write skip into file (minimal info)
				logg.Write(LogRow{
					AtUTC:   time.Now().UTC().Format(time.RFC3339),
					Mode:    "skip",
					Pos:     trade.Pos.String(),
					Action:  "LOG_5M_SKIP",
					Beta:    0,
					Spread:  0,
					Mu:      0,
					Sigma:   0,
					Z:       0,
					PnLPct:  0,
				})
				continue
			}
			writeLog(logg, "LOG_5M", trade.Pos.String(), st, nil)
			fmt.Printf("[5m log] beta=%.4f spread=%.6f z=%+.3f mode=%s pos=%s\n",
				st.Beta, st.Spread, st.Z, st.Mode, trade.Pos.String())

		case t := <-ticks:
			// update freshness timestamps
			if t.Symbol == "BTC-USDT" {
				lastBTCms = t.TimeMs
			} else if t.Symbol == "ETH-USDT" {
				lastETHms = t.TimeMs
			}

			// update aggregators and push closed bars, detect same-minute close for 1m PnL
			if t.Symbol == "BTC-USDT" {
				if bar, ok := aggBTC.Update(t); ok {
					rBTC.Push(bar)
					pendingBTC[bar.MinuteMs] = bar
					if _, ok2 := pendingETH[bar.MinuteMs]; ok2 {
						delete(pendingBTC, bar.MinuteMs)
						delete(pendingETH, bar.MinuteMs)
						printPnLClosed(bar.MinuteMs)
					}
				}
			} else if t.Symbol == "ETH-USDT" {
				if bar, ok := aggETH.Update(t); ok {
					rETH.Push(bar)
					pendingETH[bar.MinuteMs] = bar
					if _, ok2 := pendingBTC[bar.MinuteMs]; ok2 {
						delete(pendingBTC, bar.MinuteMs)
						delete(pendingETH, bar.MinuteMs)
						printPnLClosed(bar.MinuteMs)
					}
				}
			}

			// compute stats and fire signals ASAP (intra-minute)
			nowMs := time.Now().UnixMilli()
			st, ok := computeIntra(nowMs)
			if !ok {
				continue
			}

			absZ := st.Z
			if absZ < 0 {
				absZ = -absZ
			}

			switch trade.Pos {
			case Flat:
				if st.Z >= entryZ {
					if hold.shortSince == 0 {
						hold.shortSince = nowMs
					}
					hold.longSince = 0
					if holdMs == 0 || nowMs-hold.shortSince >= holdMs {
						trade.Pos = ShortSpread
						trade.EntrySpread = st.Spread
						trade.EntryAtMs = nowMs
						hold.shortSince, hold.longSince = 0, 0

						printSignal("ENTER SHORT_SPREAD (sell BTC, buy ETH)", st)
						writeLog(logg, "ENTER_SHORT", trade.Pos.String(), st, nil)
					}
				} else if st.Z <= -entryZ {
					if hold.longSince == 0 {
						hold.longSince = nowMs
					}
					hold.shortSince = 0
					if holdMs == 0 || nowMs-hold.longSince >= holdMs {
						trade.Pos = LongSpread
						trade.EntrySpread = st.Spread
						trade.EntryAtMs = nowMs
						hold.shortSince, hold.longSince = 0, 0

						printSignal("ENTER LONG_SPREAD (buy BTC, sell ETH)", st)
						writeLog(logg, "ENTER_LONG", trade.Pos.String(), st, nil)
					}
				} else {
					hold.shortSince, hold.longSince = 0, 0
				}

			case LongSpread, ShortSpread:
				if absZ <= exitZ {
					p := pnlPctFromSpreads(trade.Pos, trade.EntrySpread, st.Spread)
					printSignal(fmt.Sprintf("EXIT (back to FLAT) | pnl=%.4f%%", p*100.0), st)
					writeLog(logg, "EXIT", trade.Pos.String(), st, &p)

					trade.Pos = Flat
					trade.EntrySpread = 0
					trade.EntryAtMs = 0
					hold.shortSince, hold.longSince = 0, 0
				}
			}
		}
	}
}


az358@gaz358-BOD-WXX9:~/myprog/crypt_proto/cmd/stat_arb$ go run .
Warmup done: BTC=120 bars, ETH=120 bars | entryZ=2.20 exitZ=0.80 holdMs=1000 freshMs=300
Logging: 5m(aligned, closed-only) + ENTER/EXIT + PNL_1M(closed-only) to: statarb_5m.jsonl
[5m log] beta=0.3506 spread=8.542199 z=+1.085 mode=120closed_aligned pos=FLAT
[5m log] SKIP: not enough closed aligned bars yet





