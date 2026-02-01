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
Isolated stat-arb scaffold for KuCoin:
- window: 120 x 1-min bars (2 hours) stored in ring buffer
- warmup history via REST candles (1min)
- live updates via WS ticker events -> mid-price -> minute bars (high/low/mean)
- every 5 minutes log beta(OLS), spread, z-score into JSONL

REST:
GET /api/v1/market/candles?symbol=BTC-USDT&type=1min&startAt=...&endAt=...

WS:
POST /api/v1/bullet-public -> endpoint/token/pingInterval
subscribe topics:
  /market/ticker:BTC-USDT
  /market/ticker:ETH-USDT
*/

const (
	windowBars = 120
	tf         = "1min"
	logPath    = "statarb_5m.jsonl"
)

// ====================== Types ======================

type Bar struct {
	MinuteMs int64
	High     float64
	Low      float64
	Mean     float64 // avg of mid in this minute
	Count    int
}

type Ring struct {
	buf   []Bar
	cap   int
	head  int
	size  int
	mutex sync.Mutex
}

func NewRing(capacity int) *Ring {
	return &Ring{buf: make([]Bar, capacity), cap: capacity}
}
func (r *Ring) Push(b Bar) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.buf[r.head] = b
	r.head = (r.head + 1) % r.cap
	if r.size < r.cap {
		r.size++
	}
}
func (r *Ring) Len() int {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	return r.size
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

// Update returns (closedBar, ok) when minute switches.
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

// Partial bar for the current (not yet closed) minute.
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

// ====================== REST warmup ======================

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
		// [time, open, close, high, low, volume, turnover]
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
	Time    int64  `json:"time"` // ms (some payloads use "time")
	Time2   int64  `json:"Time"` // ms (some payloads use "Time")
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

// ====================== Stats + 5m logging ======================

type Stats struct {
	Beta   float64
	Spread float64
	Mu     float64
	Sigma  float64
	Z      float64
	Mode   string // "119+partial" or "120closed"
}

type LogRow struct {
	AtUTC string `json:"at_utc"`

	Beta   float64 `json:"beta"`
	Spread float64 `json:"spread"`
	Mu     float64 `json:"mu"`
	Sigma  float64 `json:"sigma"`
	Z      float64 `json:"z"`
	Mode   string  `json:"mode"`
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

func buildSeries120(rBTC, rETH *Ring, aggBTC, aggETH *MinuteAgg) (btc120, eth120 []float64, mode string, ok bool) {
	sBTC := rBTC.Snapshot()
	sETH := rETH.Snapshot()

	// try 119 closed + partial
	pbBTC, okB := aggBTC.Partial()
	pbETH, okE := aggETH.Partial()
	if okB && okE && pbBTC.MinuteMs == pbETH.MinuteMs && len(sBTC) >= 119 && len(sETH) >= 119 {
		baseBTC := sBTC[len(sBTC)-119:]
		baseETH := sETH[len(sETH)-119:]
		btc120 = make([]float64, 0, 120)
		eth120 = make([]float64, 0, 120)
		for i := 0; i < 119; i++ {
			btc120 = append(btc120, baseBTC[i].Mean)
			eth120 = append(eth120, baseETH[i].Mean)
		}
		btc120 = append(btc120, pbBTC.Mean)
		eth120 = append(eth120, pbETH.Mean)
		return btc120, eth120, "119+partial", true
	}

	// fallback: 120 closed
	if len(sBTC) >= 120 && len(sETH) >= 120 {
		baseBTC := sBTC[len(sBTC)-120:]
		baseETH := sETH[len(sETH)-120:]
		btc120 = make([]float64, 0, 120)
		eth120 = make([]float64, 0, 120)
		for i := 0; i < 120; i++ {
			btc120 = append(btc120, baseBTC[i].Mean)
			eth120 = append(eth120, baseETH[i].Mean)
		}
		return btc120, eth120, "120closed", true
	}

	return nil, nil, "", false
}

func calcStats(btc120, eth120 []float64) (Stats, bool) {
	if len(btc120) != 120 || len(eth120) != 120 {
		return Stats{}, false
	}
	lbtc := make([]float64, 120)
	leth := make([]float64, 120)
	for i := 0; i < 120; i++ {
		if btc120[i] <= 0 || eth120[i] <= 0 {
			return Stats{}, false
		}
		lbtc[i] = math.Log(btc120[i])
		leth[i] = math.Log(eth120[i])
	}

	beta := olsBeta(leth, lbtc)

	spreads := make([]float64, 120)
	for i := 0; i < 120; i++ {
		spreads[i] = lbtc[i] - beta*leth[i]
	}
	mu, sigma := meanStd(spreads)

	cur := spreads[119]
	z := 0.0
	if sigma > 0 {
		z = (cur - mu) / sigma
	}

	return Stats{Beta: beta, Spread: cur, Mu: mu, Sigma: sigma, Z: z}, true
}

func nextLogTimeUTC(now time.Time) time.Time {
	return now.Truncate(5 * time.Minute).Add(5 * time.Minute)
}

// ====================== Main ======================

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// graceful stop
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sig
		cancel()
	}()

	// rings
	rBTC := NewRing(windowBars)
	rETH := NewRing(windowBars)

	// warmup history (take a bit more than 120 minutes)
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
		// history minute mean proxy: (high+low)/2
		rBTC.Push(Bar{MinuteMs: s * 1000, High: kb.High, Low: kb.Low, Mean: 0.5 * (kb.High + kb.Low), Count: 1})
		rETH.Push(Bar{MinuteMs: s * 1000, High: ke.High, Low: ke.Low, Mean: 0.5 * (ke.High + ke.Low), Count: 1})
	}
	fmt.Printf("Warmup done: BTC=%d bars, ETH=%d bars\n", rBTC.Len(), rETH.Len())

	// logger
	logg, err := NewJSONLLogger(logPath)
	if err != nil {
		fmt.Println("logger error:", err)
		return
	}
	defer logg.Close()
	fmt.Println("Logging every 5 minutes to:", logPath)

	// 5-minute aligned timer
	logTimer := time.NewTimer(time.Until(nextLogTimeUTC(time.Now().UTC())))
	defer logTimer.Stop()

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

	// pending to align closed minutes
	pendingBTC := make(map[int64]Bar)
	pendingETH := make(map[int64]Bar)

	for {
		select {
		case <-ctx.Done():
			return

		case <-logTimer.C:
			// reset timer first (no drift)
			logTimer.Reset(time.Until(nextLogTimeUTC(time.Now().UTC())))

			btc120, eth120, mode, ok := buildSeries120(rBTC, rETH, aggBTC, aggETH)
			if !ok {
				fmt.Println("[5m log] not enough data yet")
				continue
			}
			st, ok := calcStats(btc120, eth120)
			if !ok {
				fmt.Println("[5m log] calc failed")
				continue
			}
			st.Mode = mode

			row := LogRow{
				AtUTC:  time.Now().UTC().Format(time.RFC3339),
				Beta:   st.Beta,
				Spread: st.Spread,
				Mu:     st.Mu,
				Sigma:  st.Sigma,
				Z:      st.Z,
				Mode:   st.Mode,
			}
			logg.Write(row)
			fmt.Printf("[5m log] beta=%.4f spread=%.6f z=%+.3f mode=%s\n", st.Beta, st.Spread, st.Z, st.Mode)

		case t := <-ticks:
			switch t.Symbol {
			case "BTC-USDT":
				if bar, ok := aggBTC.Update(t); ok {
					rBTC.Push(bar)
					pendingBTC[bar.MinuteMs] = bar
					if e, ok2 := pendingETH[bar.MinuteMs]; ok2 {
						// minute closed for both -> cleanup
						delete(pendingBTC, bar.MinuteMs)
						delete(pendingETH, bar.MinuteMs)
						_ = e // (hook for future: onPairBar)
					}
				}
			case "ETH-USDT":
				if bar, ok := aggETH.Update(t); ok {
					rETH.Push(bar)
					pendingETH[bar.MinuteMs] = bar
					if b, ok2 := pendingBTC[bar.MinuteMs]; ok2 {
						// minute closed for both -> cleanup
						delete(pendingBTC, bar.MinuteMs)
						delete(pendingETH, bar.MinuteMs)
						_ = b // (hook for future: onPairBar)
					}
				}
			}
		}
	}
}



