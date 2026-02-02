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
KuCoin stat-arb (isolated, one file) — FIXED MODEL + REALISTIC PnL

Core fixes vs previous version:
✅ OLS with intercept: ln(BTC) = alpha + beta*ln(ETH) + e
✅ Spread = residual e (not lnBTC - beta*lnETH)
✅ z-score computed on residuals
✅ Entry alpha/beta are FIXED at entry (no "repainting" PnL)
✅ PnL is PORTFOLIO-based (returns), not exp(delta_spread)

Signals (ASAP intra-minute):
- ENTER SHORT_SPREAD if z >= entryZ  (short BTC, long beta*ETH)
- ENTER LONG_SPREAD  if z <= -entryZ (long BTC, short beta*ETH)
- EXIT when abs(z) <= exitZ

Logging:
- Every 5 minutes (UTC-aligned): CLOSED-only log (stable) always prints LOG_5M or LOG_5M_SKIP
- ENTER/EXIT logs immediately (intra-minute stats)
- Every minute, strictly when BOTH closed the same minute: print PROFIT/LOSS (unrealized)
  and log PNL_1M using CLOSED-only stats + portfolio pnl

Time alignment:
- all windows are built on intersection of MinuteMs (BTC/ETH)
- ring capacity has extra (windowBars + ringExtra) to avoid "common drops to 119" on boundary minutes

Notes:
- PnL is "paper" portfolio return with notionals:
  SHORT_SPREAD: short $1 BTC, long $beta ETH
  LONG_SPREAD:  long  $1 BTC, short $beta ETH
*/

const (
	windowBars = 120
	ringExtra  = 30
	tf         = "1min"
	logPath    = "statarb_5m.jsonl"
)

// ====================== Bars / Ring ======================

type Bar struct {
	MinuteMs int64
	High     float64
	Low      float64
	Mean     float64 // for us: minute price (warmup close, live mean mid)
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
		} `json:"instanceServers"`
	} `json:"data"`
}

func getWSEndpoint(ctx context.Context) (endpoint, token string, pingInterval time.Duration, err error) {
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
	return endpoint, token, time.Duration(pi) * time.Millisecond, nil
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
				_ = c.WriteJSON(map[string]any{"id": fmt.Sprintf("%d", time.Now().UnixNano()), "type": "ping"})
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

// ====================== Stats (OLS alpha/beta + residual) ======================

type Stats struct {
	Alpha  float64
	Beta   float64
	Res    float64
	Mu     float64
	Sigma  float64
	Z      float64
	Corr   float64
	R2     float64
	Mode   string
	MinMs  int64
	Common int
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

func corr(x, y []float64) float64 {
	n := len(x)
	if n == 0 || n != len(y) {
		return 0
	}
	var sx, sy float64
	for i := 0; i < n; i++ {
		sx += x[i]
		sy += y[i]
	}
	mx := sx / float64(n)
	my := sy / float64(n)

	var sxx, syy, sxy float64
	for i := 0; i < n; i++ {
		dx := x[i] - mx
		dy := y[i] - my
		sxx += dx * dx
		syy += dy * dy
		sxy += dx * dy
	}
	if sxx == 0 || syy == 0 {
		return 0
	}
	return sxy / math.Sqrt(sxx*syy)
}

// OLS: y = alpha + beta*x
func olsAlphaBeta(x, y []float64) (alpha, beta, r2 float64) {
	n := len(x)
	if n == 0 || n != len(y) {
		return 0, 1, 0
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
		return my, 1, 0
	}
	beta = cov / vx
	alpha = my - beta*mx

	// R^2
	var sse, sst float64
	for i := 0; i < n; i++ {
		yhat := alpha + beta*x[i]
		err := y[i] - yhat
		sse += err * err
		dy := y[i] - my
		sst += dy * dy
	}
	if sst > 0 {
		r2 = 1 - sse/sst
		if r2 < 0 {
			r2 = 0
		}
		if r2 > 1 {
			r2 = 1
		}
	}
	return
}

// Compute stats on a window of prices (BTC, ETH)
func calcStatsFromPrices(btc, eth []float64, mode string, minMs int64, common int) (Stats, bool) {
	n := len(btc)
	if n == 0 || n != len(eth) {
		return Stats{}, false
	}

	lbtc := make([]float64, n)
	leth := make([]float64, n)
	for i := 0; i < n; i++ {
		if btc[i] <= 0 || eth[i] <= 0 {
			return Stats{}, false
		}
		lbtc[i] = math.Log(btc[i])
		leth[i] = math.Log(eth[i])
	}

	c := corr(leth, lbtc)
	alpha, beta, r2 := olsAlphaBeta(leth, lbtc)

	res := make([]float64, n)
	for i := 0; i < n; i++ {
		res[i] = lbtc[i] - (alpha + beta*leth[i])
	}
	mu, sigma := meanStd(res)

	last := res[n-1]
	z := 0.0
	if sigma > 0 {
		z = (last - mu) / sigma
	}

	return Stats{
		Alpha:  alpha,
		Beta:   beta,
		Res:    last,
		Mu:     mu,
		Sigma:  sigma,
		Z:      z,
		Corr:   c,
		R2:     r2,
		Mode:   mode,
		MinMs:  minMs,
		Common: common,
	}, true
}

// ====================== Alignment helpers (also return last close prices) ======================

// CLOSED-only window: last window from common minutes
func buildAlignedClosed(rBTC, rETH *Ring, window int) (btc, eth []float64, lastBTC, lastETH float64, minMs int64, common int, ok bool) {
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

	commonMs := make([]int64, 0, len(mB))
	for ms := range mB {
		if _, ok := mE[ms]; ok {
			commonMs = append(commonMs, ms)
		}
	}
	sort.Slice(commonMs, func(i, j int) bool { return commonMs[i] < commonMs[j] })
	common = len(commonMs)
	if common < window {
		return nil, nil, 0, 0, 0, common, false
	}
	commonMs = commonMs[common-window:]

	btc = make([]float64, 0, window)
	eth = make([]float64, 0, window)
	for _, ms := range commonMs {
		btc = append(btc, mB[ms])
		eth = append(eth, mE[ms])
	}
	lastBTC = btc[len(btc)-1]
	lastETH = eth[len(eth)-1]
	return btc, eth, lastBTC, lastETH, commonMs[len(commonMs)-1], common, true
}

// INTRA window: prefer (window-1) closed + same-minute partial
func buildAlignedIntra(rBTC, rETH *Ring, aggBTC, aggETH *MinuteAgg, window int) (btc, eth []float64, lastBTC, lastETH float64, mode string, minMs int64, common int, ok bool) {
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

	commonMs := make([]int64, 0, len(mB))
	for ms := range mB {
		if _, ok := mE[ms]; ok {
			commonMs = append(commonMs, ms)
		}
	}
	sort.Slice(commonMs, func(i, j int) bool { return commonMs[i] < commonMs[j] })
	common = len(commonMs)

	pbB, okB := aggBTC.Partial()
	pbE, okE := aggETH.Partial()
	usePartial := okB && okE && pbB.MinuteMs == pbE.MinuteMs && pbB.Mean > 0 && pbE.Mean > 0

	if usePartial {
		if common < window-1 {
			return nil, nil, 0, 0, "", 0, common, false
		}
		commonMs = commonMs[common-(window-1):]
		btc = make([]float64, 0, window)
		eth = make([]float64, 0, window)
		for _, ms := range commonMs {
			btc = append(btc, mB[ms])
			eth = append(eth, mE[ms])
		}
		btc = append(btc, pbB.Mean)
		eth = append(eth, pbE.Mean)
		lastBTC = pbB.Mean
		lastETH = pbE.Mean
		return btc, eth, lastBTC, lastETH, fmt.Sprintf("%d+partial_aligned", window-1), pbB.MinuteMs, common, true
	}

	// fallback to closed-only
	if common < window {
		return nil, nil, 0, 0, "", 0, common, false
	}
	commonMs = commonMs[common-window:]
	btc = make([]float64, 0, window)
	eth = make([]float64, 0, window)
	for _, ms := range commonMs {
		btc = append(btc, mB[ms])
		eth = append(eth, mE[ms])
	}
	lastBTC = btc[len(btc)-1]
	lastETH = eth[len(eth)-1]
	return btc, eth, lastBTC, lastETH, fmt.Sprintf("%dclosed_aligned", window), commonMs[len(commonMs)-1], common, true
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

	Alpha  float64 `json:"alpha"`
	Beta   float64 `json:"beta"`
	Res    float64 `json:"res"`
	Mu     float64 `json:"mu"`
	Sigma  float64 `json:"sigma"`
	Z      float64 `json:"z"`
	Corr   float64 `json:"corr"`
	R2     float64 `json:"r2"`
	Mode   string  `json:"mode"`
	Common int     `json:"common"`

	Pos         string  `json:"pos"`
	Action      string  `json:"action"` // LOG_5M, LOG_5M_SKIP, ENTER_LONG, ENTER_SHORT, EXIT, PNL_1M
	PnLPct      float64 `json:"pnl_pct,omitempty"`
	EntryBeta   float64 `json:"entry_beta,omitempty"`
	EntryAlpha  float64 `json:"entry_alpha,omitempty"`
	EntryBTC    float64 `json:"entry_btc,omitempty"`
	EntryETH    float64 `json:"entry_eth,omitempty"`
	NowBTC      float64 `json:"now_btc,omitempty"`
	NowETH      float64 `json:"now_eth,omitempty"`
	ClosedMinMs int64   `json:"closed_min_ms,omitempty"`
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

func writeLog(logg *JSONLLogger, action string, pos string, st Stats, trade *TradeState, pnlPct *float64, nowBTC, nowETH float64, closedMinMs int64) {
	row := LogRow{
		AtUTC:   time.Now().UTC().Format(time.RFC3339),
		Alpha:   st.Alpha,
		Beta:    st.Beta,
		Res:     st.Res,
		Mu:      st.Mu,
		Sigma:   st.Sigma,
		Z:       st.Z,
		Corr:    st.Corr,
		R2:      st.R2,
		Mode:    st.Mode,
		Common:  st.Common,
		Pos:     pos,
		Action:  action,
		NowBTC:  nowBTC,
		NowETH:  nowETH,
		ClosedMinMs: closedMinMs,
	}
	if trade != nil && trade.Pos != Flat {
		row.EntryBeta = trade.EntryBeta
		row.EntryAlpha = trade.EntryAlpha
		row.EntryBTC = trade.EntryBTC
		row.EntryETH = trade.EntryETH
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
	Pos Position

	EntryAlpha float64
	EntryBeta  float64

	EntryBTC float64
	EntryETH float64

	EntryAtMs int64
}

func pnlPctPortfolio(pos Position, entryBTC, entryETH, nowBTC, nowETH, beta float64) float64 {
	if entryBTC <= 0 || entryETH <= 0 || nowBTC <= 0 || nowETH <= 0 {
		return 0
	}
	rBTC := nowBTC/entryBTC - 1.0
	rETH := nowETH/entryETH - 1.0

	switch pos {
	case ShortSpread:
		// short $1 BTC, long $beta ETH
		return (-rBTC) + beta*(rETH)
	case LongSpread:
		// long $1 BTC, short $beta ETH
		return (rBTC) - beta*(rETH)
	default:
		return 0
	}
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

	// rings with EXTRA capacity
	rBTC := NewRing(windowBars + ringExtra)
	rETH := NewRing(windowBars + ringExtra)

	// warmup history (take a bit more than needed)
	nowSec := time.Now().Unix()
	startSec := nowSec - int64((windowBars+ringExtra+30)*60)

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

	// align warmup by startSec and push CLOSE into Mean
	mBTC := make(map[int64]Kline, len(btcK))
	for _, k := range btcK {
		mBTC[k.StartSec] = k
	}
	mETH := make(map[int64]Kline, len(ethK))
	for _, k := range ethK {
		mETH[k.StartSec] = k
	}

	secs := make([]int64, 0, len(mBTC))
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
		rBTC.Push(Bar{MinuteMs: s * 1000, High: kb.High, Low: kb.Low, Mean: kb.Close, Count: 1})
		rETH.Push(Bar{MinuteMs: s * 1000, High: ke.High, Low: ke.Low, Mean: ke.Close, Count: 1})
	}

	fmt.Printf("Warmup done: BTC=%d bars, ETH=%d bars | entryZ=%.2f exitZ=%.2f holdMs=%d freshMs=%d | ringCap=%d\n",
		rBTC.Len(), rETH.Len(), entryZ, exitZ, holdMs, freshMs, windowBars+ringExtra)

	// logger
	logg, err := NewJSONLLogger(logPath)
	if err != nil {
		fmt.Println("logger error:", err)
		return
	}
	defer logg.Close()
	fmt.Println("Logging: 5m(aligned, CLOSED) + ENTER/EXIT + PNL_1M(CLOSED) to:", logPath)

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
	computeIntra := func(nowMs int64) (Stats, float64, float64, bool) {
		if nowMs-lastCalcAtMs < calcThrottleMs {
			return Stats{}, 0, 0, false
		}
		lastCalcAtMs = nowMs

		if lastBTCms == 0 || lastETHms == 0 {
			return Stats{}, 0, 0, false
		}
		d := lastBTCms - lastETHms
		if d < 0 {
			d = -d
		}
		if d > freshMs {
			return Stats{}, 0, 0, false
		}

		btc, eth, lastBTC, lastETH, mode, minMs, common, ok := buildAlignedIntra(rBTC, rETH, aggBTC, aggETH, windowBars)
		if !ok {
			return Stats{}, 0, 0, false
		}
		st, ok := calcStatsFromPrices(btc, eth, mode, minMs, common)
		return st, lastBTC, lastETH, ok
	}

	// compute stats strictly on closed minutes (for 1m PnL + 5m logs)
	computeClosed := func() (Stats, float64, float64, int64, bool) {
		btc, eth, lastBTC, lastETH, minMs, common, ok := buildAlignedClosed(rBTC, rETH, windowBars)
		if !ok {
			return Stats{}, 0, 0, 0, false
		}
		st, ok := calcStatsFromPrices(btc, eth, "closed_aligned", minMs, common)
		return st, lastBTC, lastETH, minMs, ok
	}

	printSignal := func(tag string, st Stats) {
		t := time.Now().UTC().Format("2006-01-02 15:04:05")
		min := time.UnixMilli(st.MinMs).UTC().Format("2006-01-02 15:04")
		fmt.Printf("[%s] %s | minute=%s | z=%+.3f alpha=%.4f beta=%.4f corr=%.3f r2=%.3f res=%.6f mode=%s pos=%s\n",
			t, tag, min, st.Z, st.Alpha, st.Beta, st.Corr, st.R2, st.Res, st.Mode, trade.Pos.String())
	}

	printPnLClosed := func(closedMinuteMs int64) {
		if trade.Pos == Flat {
			return
		}

		st, nowBTC, nowETH, _, ok := computeClosed()
		if !ok {
			return
		}

		p := pnlPctPortfolio(trade.Pos, trade.EntryBTC, trade.EntryETH, nowBTC, nowETH, trade.EntryBeta)

		min := time.UnixMilli(closedMinuteMs).UTC().Format("2006-01-02 15:04")
		sign := "PROFIT"
		if p < 0 {
			sign = "LOSS"
		}
		fmt.Printf("[1m %s] %.4f%% | minute=%s | pos=%s | z=%+.3f entryBeta=%.4f liveBeta=%.4f\n",
			sign, p*100.0, min, trade.Pos.String(), st.Z, trade.EntryBeta, st.Beta)

		writeLog(logg, "PNL_1M", trade.Pos.String(), st, &trade, &p, nowBTC, nowETH, closedMinuteMs)
	}

	for {
		select {
		case <-ctx.Done():
			return

		// 5m log: ALWAYS print something, CLOSED-only
		case <-log5mCh:
			st, nowBTC, nowETH, minMs, ok := computeClosed()
			if !ok {
				fmt.Printf("[5m log] SKIP: not enough closed common minutes (need=%d)\n", windowBars)
				logg.Write(LogRow{
					AtUTC:   time.Now().UTC().Format(time.RFC3339),
					Mode:    "skip",
					Pos:     trade.Pos.String(),
					Action:  "LOG_5M_SKIP",
				})
				continue
			}
			writeLog(logg, "LOG_5M", trade.Pos.String(), st, func() *TradeState {
				if trade.Pos == Flat {
					return nil
				}
				return &trade
			}(), nil, nil, nowBTC, nowETH, minMs)

			fmt.Printf("[5m log] z=%+.3f alpha=%.4f beta=%.4f corr=%.3f r2=%.3f res=%.6f common=%d pos=%s\n",
				st.Z, st.Alpha, st.Beta, st.Corr, st.R2, st.Res, st.Common, trade.Pos.String())

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
			st, nowBTC, nowETH, ok := computeIntra(nowMs)
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
						trade.EntryAlpha = st.Alpha
						trade.EntryBeta = st.Beta
						trade.EntryBTC = nowBTC
						trade.EntryETH = nowETH
						trade.EntryAtMs = nowMs
						hold.shortSince, hold.longSince = 0, 0

						printSignal("ENTER SHORT_SPREAD (short BTC, long beta*ETH)", st)
						writeLog(logg, "ENTER_SHORT", trade.Pos.String(), st, &trade, nil, nil, nowBTC, nowETH, st.MinMs)
					}
				} else if st.Z <= -entryZ {
					if hold.longSince == 0 {
						hold.longSince = nowMs
					}
					hold.shortSince = 0
					if holdMs == 0 || nowMs-hold.longSince >= holdMs {
						trade.Pos = LongSpread
						trade.EntryAlpha = st.Alpha
						trade.EntryBeta = st.Beta
						trade.EntryBTC = nowBTC
						trade.EntryETH = nowETH
						trade.EntryAtMs = nowMs
						hold.shortSince, hold.longSince = 0, 0

						printSignal("ENTER LONG_SPREAD (long BTC, short beta*ETH)", st)
						writeLog(logg, "ENTER_LONG", trade.Pos.String(), st, &trade, nil, nil, nowBTC, nowETH, st.MinMs)
					}
				} else {
					hold.shortSince, hold.longSince = 0, 0
				}

			case LongSpread, ShortSpread:
				if absZ <= exitZ {
					// compute current closed stats for a more stable exit print is optional,
					// but we exit "ASAP" by intra-minute z.
					p := pnlPctPortfolio(trade.Pos, trade.EntryBTC, trade.EntryETH, nowBTC, nowETH, trade.EntryBeta)
					printSignal(fmt.Sprintf("EXIT (back to FLAT) | pnl=%.4f%%", p*100.0), st)
					writeLog(logg, "EXIT", trade.Pos.String(), st, &trade, &p, nowBTC, nowETH, st.MinMs)

					trade = TradeState{Pos: Flat}
					hold.shortSince, hold.longSince = 0, 0
				}
			}
		}
	}
}



[{
	"resource": "/home/gaz358/myprog/crypt_proto/cmd/stat_arb/stat_arb.go",
	"owner": "_generated_diagnostic_collection_name_#1",
	"code": {
		"value": "WrongArgCount",
		"target": {
			"$mid": 1,
			"path": "/golang.org/x/tools/internal/typesinternal",
			"scheme": "https",
			"authority": "pkg.go.dev",
			"fragment": "WrongArgCount"
		}
	},
	"severity": 8,
	"message": "too many arguments in call to writeLog\n\thave (*JSONLLogger, string, string, Stats, *TradeState, nil, nil, float64, float64, int64)\n\twant (*JSONLLogger, string, string, Stats, *TradeState, *float64, float64, float64, int64)",
	"source": "compiler",
	"startLineNumber": 1049,
	"startColumn": 35,
	"endLineNumber": 1049,
	"endColumn": 40,
	"modelVersionId": 4,
	"origin": "extHost1"
}]


[{
	"resource": "/home/gaz358/myprog/crypt_proto/cmd/stat_arb/stat_arb.go",
	"owner": "_generated_diagnostic_collection_name_#1",
	"code": {
		"value": "WrongArgCount",
		"target": {
			"$mid": 1,
			"path": "/golang.org/x/tools/internal/typesinternal",
			"scheme": "https",
			"authority": "pkg.go.dev",
			"fragment": "WrongArgCount"
		}
	},
	"severity": 8,
	"message": "too many arguments in call to writeLog\n\thave (*JSONLLogger, string, string, Stats, *TradeState, nil, nil, float64, float64, int64)\n\twant (*JSONLLogger, string, string, Stats, *TradeState, *float64, float64, float64, int64)",
	"source": "compiler",
	"startLineNumber": 1114,
	"startColumn": 95,
	"endLineNumber": 1114,
	"endColumn": 103,
	"modelVersionId": 4,
	"origin": "extHost1"
}]


[{
	"resource": "/home/gaz358/myprog/crypt_proto/cmd/stat_arb/stat_arb.go",
	"owner": "_generated_diagnostic_collection_name_#1",
	"code": {
		"value": "WrongArgCount",
		"target": {
			"$mid": 1,
			"path": "/golang.org/x/tools/internal/typesinternal",
			"scheme": "https",
			"authority": "pkg.go.dev",
			"fragment": "WrongArgCount"
		}
	},
	"severity": 8,
	"message": "too many arguments in call to writeLog\n\thave (*JSONLLogger, string, string, Stats, *TradeState, nil, nil, float64, float64, int64)\n\twant (*JSONLLogger, string, string, Stats, *TradeState, *float64, float64, float64, int64)",
	"source": "compiler",
	"startLineNumber": 1131,
	"startColumn": 94,
	"endLineNumber": 1131,
	"endColumn": 102,
	"modelVersionId": 4,
	"origin": "extHost1"
}]





