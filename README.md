arb.go

package calculator

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"crypt_proto/internal/queue"
	"crypt_proto/pkg/models"
)

const feeM = 0.9992

const (
	minLogProfitPct = 0.0
	minLogVolumeUSDT = 50.0
	arbDedupWindow = 300 * time.Millisecond
	lastPrintTTL = 5 * time.Second
)

type LegRule struct {
	Index int

	RawLeg      string
	Step        float64
	MinQty      float64
	MinNotional float64

	Symbol      string
	Side        string
	Base        string
	Quote       string
	QtyStep     float64
	QuoteStep   float64
	PriceStep   float64
	LegMinQty   float64
	LegMinQuote float64
	LegMinNotnl float64

	Key string
}

type Triangle struct {
	A, B, C string
	Legs    [3]LegRule
}

type Calculator struct {
	mem       *queue.MemoryStore
	bySymbol  map[string][]*Triangle
	fileLog   *log.Logger
	lastPrint map[string]time.Time
}

func NewCalculator(mem *queue.MemoryStore, triangles []*Triangle) *Calculator {
	f, err := os.OpenFile("arb_opportunities.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		log.Fatalf("failed to open log: %v", err)
	}

	bySymbol := make(map[string][]*Triangle, 1024)
	for _, t := range triangles {
		for _, leg := range t.Legs {
			if leg.Symbol == "" {
				continue
			}
			bySymbol[leg.Symbol] = append(bySymbol[leg.Symbol], t)
		}
	}

	log.Printf("[Calculator] indexed %d symbols\n", len(bySymbol))

	return &Calculator{
		mem:       mem,
		bySymbol:  bySymbol,
		fileLog:   log.New(f, "", log.LstdFlags),
		lastPrint: make(map[string]time.Time),
	}
}

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
		if quote.Bid <= 0 || quote.Ask <= 0 || quote.BidSize <= 0 || quote.AskSize <= 0 {
			return
		}
		q[i] = quote
	}

	maxStart, ok := computeMaxStartTopOfBook(tri, q)
	if !ok || maxStart <= 0 {
		return
	}

	finalAmount, diag, ok := simulateTriangle(maxStart, tri, q)
	if !ok || finalAmount <= 0 {
		return
	}

	profitUSDT := finalAmount - maxStart
	profitPct := profitUSDT / maxStart
	if profitPct <= minLogProfitPct || maxStart <= minLogVolumeUSDT {
		return
	}

	now := time.Now()
	dedupKey := fmt.Sprintf(
		"%s|%s|%s|%.4f|%.2f",
		tri.A, tri.B, tri.C,
		math.Round(profitPct*10000)/10000,
		math.Round(maxStart*100)/100,
	)
	if last, exists := c.lastPrint[dedupKey]; exists && now.Sub(last) < arbDedupWindow {
		return
	}
	c.lastPrint[dedupKey] = now
	c.cleanupLastPrint(now)

	msg := fmt.Sprintf(
		"[ARB] %s→%s→%s | %.4f%% | volume=%.2f USDT | profit=%.6f USDT | "+
			"l1=%s %s out=%.8f age=%dms | "+
			"l2=%s %s out=%.8f age=%dms | "+
			"l3=%s %s out=%.8f age=%dms",
		tri.A, tri.B, tri.C,
		profitPct*100, maxStart, profitUSDT,
		tri.Legs[0].Symbol, tri.Legs[0].Side, diag[0].Out, quoteAgeMs(now, q[0]),
		tri.Legs[1].Symbol, tri.Legs[1].Side, diag[1].Out, quoteAgeMs(now, q[1]),
		tri.Legs[2].Symbol, tri.Legs[2].Side, diag[2].Out, quoteAgeMs(now, q[2]),
	)
	log.Println(msg)
	c.fileLog.Println(msg)
}

func quoteAgeMs(now time.Time, q queue.Quote) int64 {
	if q.Timestamp == 0 {
		return -1
	}
	return now.UnixMilli() - q.Timestamp
}

func (c *Calculator) cleanupLastPrint(now time.Time) {
	if len(c.lastPrint) < 5000 {
		return
	}
	cutoff := now.Add(-lastPrintTTL)
	for k, ts := range c.lastPrint {
		if ts.Before(cutoff) {
			delete(c.lastPrint, k)
		}
	}
}

type legExecution struct {
	In            float64
	Out           float64
	Price         float64
	BookLimitIn   float64
	TradeQty      float64
	TradeNotional float64
}

func simulateTriangle(startUSDT float64, tri *Triangle, q [3]queue.Quote) (float64, [3]legExecution, bool) {
	var diag [3]legExecution
	amount := startUSDT
	for i := 0; i < 3; i++ {
		out, d, ok := executeLeg(amount, tri.Legs[i], q[i])
		if !ok {
			return 0, diag, false
		}
		diag[i] = d
		amount = out
	}
	return amount, diag, true
}

func executeLeg(in float64, leg LegRule, q queue.Quote) (float64, legExecution, bool) {
	side := strings.ToUpper(strings.TrimSpace(leg.Side))
	if side == "" {
		side = detectSideFromRawLeg(leg.RawLeg)
	}
	if side != "BUY" && side != "SELL" {
		return 0, legExecution{}, false
	}
	if in <= 0 {
		return 0, legExecution{}, false
	}

	qtyStep := firstPositive(leg.QtyStep, leg.Step)
	minQty := firstPositive(leg.LegMinQty, leg.MinQty)
	minQuote := leg.LegMinQuote
	minNotional := firstPositive(leg.LegMinNotnl, leg.MinNotional)

	switch side {
	case "BUY":
		if q.Ask <= 0 || q.AskSize <= 0 {
			return 0, legExecution{}, false
		}
		bookLimitIn := q.Ask * q.AskSize
		spendQuote := math.Min(in, bookLimitIn)
		if spendQuote <= 0 {
			return 0, legExecution{}, false
		}

		rawQty := spendQuote / q.Ask
		tradeQty := floorToStep(rawQty, qtyStep)
		if tradeQty <= 0 {
			return 0, legExecution{}, false
		}

		tradeNotional := tradeQty * q.Ask
		if tradeNotional <= 0 || tradeNotional > spendQuote+eps() {
			return 0, legExecution{}, false
		}
		if minQty > 0 && tradeQty+eps() < minQty {
			return 0, legExecution{}, false
		}
		if minQuote > 0 && tradeNotional+eps() < minQuote {
			return 0, legExecution{}, false
		}
		if minNotional > 0 && tradeNotional+eps() < minNotional {
			return 0, legExecution{}, false
		}

		outBase := tradeQty * feeM
		if outBase <= 0 {
			return 0, legExecution{}, false
		}

		return outBase, legExecution{In: in, Out: outBase, Price: q.Ask, BookLimitIn: bookLimitIn, TradeQty: tradeQty, TradeNotional: tradeNotional}, true

	case "SELL":
		if q.Bid <= 0 || q.BidSize <= 0 {
			return 0, legExecution{}, false
		}
		bookLimitIn := q.BidSize
		sellBase := math.Min(in, bookLimitIn)
		tradeQty := floorToStep(sellBase, qtyStep)
		if tradeQty <= 0 {
			return 0, legExecution{}, false
		}

		tradeNotional := tradeQty * q.Bid
		if tradeNotional <= 0 {
			return 0, legExecution{}, false
		}
		if minQty > 0 && tradeQty+eps() < minQty {
			return 0, legExecution{}, false
		}
		if minQuote > 0 && tradeNotional+eps() < minQuote {
			return 0, legExecution{}, false
		}
		if minNotional > 0 && tradeNotional+eps() < minNotional {
			return 0, legExecution{}, false
		}

		outQuote := tradeNotional * feeM
		if outQuote <= 0 {
			return 0, legExecution{}, false
		}
		return outQuote, legExecution{In: in, Out: outQuote, Price: q.Bid, BookLimitIn: bookLimitIn, TradeQty: tradeQty, TradeNotional: tradeNotional}, true
	}

	return 0, legExecution{}, false
}

func computeMaxStartTopOfBook(tri *Triangle, q [3]queue.Quote) (float64, bool) {
	kIn := 1.0
	maxStart := math.MaxFloat64

	for i := 0; i < 3; i++ {
		leg := tri.Legs[i]
		side := strings.ToUpper(strings.TrimSpace(leg.Side))
		if side == "" {
			side = detectSideFromRawLeg(leg.RawLeg)
		}
		if side != "BUY" && side != "SELL" {
			return 0, false
		}

		var limitIn float64
		switch side {
		case "BUY":
			if q[i].Ask <= 0 || q[i].AskSize <= 0 {
				return 0, false
			}
			limitIn = q[i].Ask * q[i].AskSize
			maxByThis := limitIn / kIn
			if maxByThis < maxStart {
				maxStart = maxByThis
			}
			kIn *= (1.0 / q[i].Ask) * feeM
		case "SELL":
			if q[i].Bid <= 0 || q[i].BidSize <= 0 {
				return 0, false
			}
			limitIn = q[i].BidSize
			maxByThis := limitIn / kIn
			if maxByThis < maxStart {
				maxStart = maxByThis
			}
			kIn *= q[i].Bid * feeM
		}

		if kIn <= 0 {
			return 0, false
		}
	}

	if maxStart <= 0 || !isFinite(maxStart) {
		return 0, false
	}
	return maxStart, true
}

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
	if len(rows) < 2 {
		return nil, nil
	}

	header := make(map[string]int, len(rows[0]))
	for i, col := range rows[0] {
		header[strings.TrimSpace(col)] = i
	}

	var res []*Triangle
	for _, row := range rows[1:] {
		if len(strings.TrimSpace(strings.Join(row, ""))) == 0 {
			continue
		}

		t := &Triangle{A: getString(row, header, "A"), B: getString(row, header, "B"), C: getString(row, header, "C")}
		for i := 1; i <= 3; i++ {
			idx := i - 1
			leg := LegRule{
				Index:       idx,
				RawLeg:      getString(row, header, fmt.Sprintf("Leg%d", i)),
				Step:        getFloat(row, header, fmt.Sprintf("Step%d", i)),
				MinQty:      getFloat(row, header, fmt.Sprintf("MinQty%d", i)),
				MinNotional: getFloat(row, header, fmt.Sprintf("MinNotional%d", i)),
				Symbol:      getString(row, header, fmt.Sprintf("Leg%dSymbol", i)),
				Side:        strings.ToUpper(getString(row, header, fmt.Sprintf("Leg%dSide", i))),
				Base:        getString(row, header, fmt.Sprintf("Leg%dBase", i)),
				Quote:       getString(row, header, fmt.Sprintf("Leg%dQuote", i)),
				QtyStep:     getFloat(row, header, fmt.Sprintf("Leg%dQtyStep", i)),
				QuoteStep:   getFloat(row, header, fmt.Sprintf("Leg%dQuoteStep", i)),
				PriceStep:   getFloat(row, header, fmt.Sprintf("Leg%dPriceStep", i)),
				LegMinQty:   getFloat(row, header, fmt.Sprintf("Leg%dMinQty", i)),
				LegMinQuote: getFloat(row, header, fmt.Sprintf("Leg%dMinQuote", i)),
				LegMinNotnl: getFloat(row, header, fmt.Sprintf("Leg%dMinNotional", i)),
			}
			if leg.Symbol == "" {
				leg.Symbol = symbolFromRawLeg(leg.RawLeg)
			}
			if leg.Side == "" {
				leg.Side = detectSideFromRawLeg(leg.RawLeg)
			}
			if leg.Symbol == "" || leg.Side == "" {
				continue
			}
			leg.Key = "KuCoin|" + leg.Symbol
			t.Legs[idx] = leg
		}
		if t.Legs[0].Symbol == "" || t.Legs[1].Symbol == "" || t.Legs[2].Symbol == "" {
			continue
		}
		res = append(res, t)
	}
	return res, nil
}

func symbolFromRawLeg(raw string) string {
	raw = strings.ToUpper(strings.TrimSpace(raw))
	if raw == "" {
		return ""
	}
	parts := strings.Fields(raw)
	if len(parts) != 2 {
		return ""
	}
	pair := strings.Split(parts[1], "/")
	if len(pair) != 2 {
		return ""
	}
	return strings.TrimSpace(pair[0]) + "-" + strings.TrimSpace(pair[1])
}

func detectSideFromRawLeg(raw string) string {
	raw = strings.ToUpper(strings.TrimSpace(raw))
	if strings.HasPrefix(raw, "BUY ") {
		return "BUY"
	}
	if strings.HasPrefix(raw, "SELL ") {
		return "SELL"
	}
	return ""
}

func getString(row []string, header map[string]int, key string) string {
	idx, ok := header[key]
	if !ok || idx < 0 || idx >= len(row) {
		return ""
	}
	return strings.TrimSpace(row[idx])
}

func getFloat(row []string, header map[string]int, key string) float64 {
	s := getString(row, header, key)
	if s == "" {
		return 0
	}
	v, err := strconv.ParseFloat(strings.ReplaceAll(s, ",", "."), 64)
	if err != nil {
		return 0
	}
	return v
}

func floorToStep(v, step float64) float64 {
	if v <= 0 {
		return 0
	}
	if step <= 0 {
		return v
	}
	n := math.Floor((v + eps()) / step)
	if n <= 0 {
		return 0
	}
	return n * step
}

func firstPositive(vals ...float64) float64 {
	for _, v := range vals {
		if v > 0 {
			return v
		}
	}
	return 0
}

func isFinite(v float64) bool { return !math.IsNaN(v) && !math.IsInf(v, 0) }
func eps() float64            { return 1e-12 }



collector


package collector

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"crypt_proto/pkg/models"

	"github.com/gorilla/websocket"
	"github.com/tidwall/gjson"
)

const (
	maxSubsPerWS   = 126
	subRate        = 120 * time.Millisecond
	pingInterval   = 20 * time.Second
	reconnectDelay = 3 * time.Second
)

type KuCoinCollector struct {
	ctx    context.Context
	cancel context.CancelFunc
	wsList []*kucoinWS
	out    chan<- *models.MarketData
}

type Last struct {
	Bid     float64
	Ask     float64
	BidSize float64
	AskSize float64
}

type kucoinWS struct {
	id      int
	conn    *websocket.Conn
	symbols []string

	last map[string]Last
}

func NewKuCoinCollectorFromCSV(path string) (*KuCoinCollector, []string, error) {
	symbols, err := readPairsFromCSV(path)
	if err != nil {
		return nil, nil, err
	}
	if len(symbols) == 0 {
		return nil, nil, fmt.Errorf("no symbols")
	}

	ctx, cancel := context.WithCancel(context.Background())

	var wsList []*kucoinWS
	for i := 0; i < len(symbols); i += maxSubsPerWS {
		end := i + maxSubsPerWS
		if end > len(symbols) {
			end = len(symbols)
		}

		wsList = append(wsList, &kucoinWS{
			id:      len(wsList),
			symbols: symbols[i:end],
			last:    make(map[string]Last),
		})
	}

	c := &KuCoinCollector{
		ctx:    ctx,
		cancel: cancel,
		wsList: wsList,
	}

	return c, symbols, nil
}

func (c *KuCoinCollector) Name() string { return "KuCoin" }

func (c *KuCoinCollector) Start(out chan<- *models.MarketData) error {
	c.out = out
	for _, ws := range c.wsList {
		go ws.run(c) // запускаем WS с авто-перезапуском
	}
	log.Printf("[KuCoin] started with %d WS\n", len(c.wsList))
	return nil
}

func (c *KuCoinCollector) Stop() error {
	c.cancel()
	for _, ws := range c.wsList {
		if ws.conn != nil {
			_ = ws.conn.Close()
		}
	}
	return nil
}

func (ws *kucoinWS) run(c *KuCoinCollector) {
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
		}

		// подключаемся
		if err := ws.connect(); err != nil {
			log.Printf("[KuCoin WS %d] connect error: %v, retry in %v\n", ws.id, err, reconnectDelay)
			time.Sleep(reconnectDelay)
			continue
		}

		// запускаем подписки и пинг
		done := make(chan struct{})
		go func() {
			ws.subscribeLoop()
			close(done)
		}()
		go ws.pingLoop()

		// читаем данные
		if err := ws.readLoop(c); err != nil {
			log.Printf("[KuCoin WS %d] readLoop ended: %v\n", ws.id, err)
		}

		// закрываем соединение и перезапускаем
		if ws.conn != nil {
			_ = ws.conn.Close()
			ws.conn = nil
		}

		log.Printf("[KuCoin WS %d] reconnecting in %v...\n", ws.id, reconnectDelay)
		time.Sleep(reconnectDelay)
	}
}

func (ws *kucoinWS) connect() error {
	req, _ := http.NewRequest("POST", "https://api.kucoin.com/api/v1/bullet-public", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var r struct {
		Data struct {
			Token           string `json:"token"`
			InstanceServers []struct {
				Endpoint string `json:"endpoint"`
			} `json:"instanceServers"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return err
	}

	url := fmt.Sprintf(
		"%s?token=%s&connectId=%d",
		r.Data.InstanceServers[0].Endpoint,
		r.Data.Token,
		time.Now().UnixNano(),
	)

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return err
	}

	ws.conn = conn
	log.Printf("[KuCoin WS %d] connected\n", ws.id)
	return nil
}

func (ws *kucoinWS) subscribeLoop() {
	t := time.NewTicker(subRate)
	defer t.Stop()
	for _, s := range ws.symbols {
		<-t.C
		_ = ws.conn.WriteJSON(map[string]any{
			"id":       time.Now().UnixNano(),
			"type":     "subscribe",
			"topic":    "/market/ticker:" + s,
			"response": true,
		})
	}
}

func (ws *kucoinWS) pingLoop() {
	t := time.NewTicker(pingInterval)
	defer t.Stop()
	for range t.C {
		if ws.conn != nil {
			_ = ws.conn.WriteJSON(map[string]any{"id": time.Now().UnixNano(), "type": "ping"})
		}
	}
}

func (ws *kucoinWS) readLoop(c *KuCoinCollector) error {
	for {
		if ws.conn == nil {
			return fmt.Errorf("connection closed")
		}
		_, msg, err := ws.conn.ReadMessage()
		if err != nil {
			return err
		}
		ws.handle(c, msg)
	}
}

func (ws *kucoinWS) handle(c *KuCoinCollector, msg []byte) {
	const prefix = "/market/ticker:"
	const prefixLen = len(prefix)

	topicRes := gjson.GetBytes(msg, "topic")
	if !topicRes.Exists() {
		return
	}

	raw := topicRes.Raw // строка JSON: "/market/ticker:BTC-USDT"
	if len(raw) <= prefixLen+2 {
		return
	}

	if raw[1:1+prefixLen] != prefix { // пропускаем первую кавычку
		return
	}

	symbol := raw[1+prefixLen : len(raw)-1]

	data := gjson.GetBytes(msg, "data")
	bid := data.Get("bestBid").Float()
	ask := data.Get("bestAsk").Float()
	if bid == 0 || ask == 0 {
		return
	}

	bidSize := data.Get("bestBidSize").Float()
	askSize := data.Get("bestAskSize").Float()
	if bidSize == 0 || askSize == 0 {
		return
	}

	if last, ok := ws.last[symbol]; ok && last.Bid == bid && last.Ask == ask && last.BidSize == bidSize && last.AskSize == askSize {
		return
	}

	ws.last[symbol] = Last{Bid: bid, Ask: ask, BidSize: bidSize, AskSize: askSize}

	c.out <- &models.MarketData{
		Exchange:  "KuCoin",
		Symbol:    symbol,
		Bid:       bid,
		Ask:       ask,
		BidSize:   bidSize,
		AskSize:   askSize,
		Timestamp: time.Now().UnixMilli(),
	}
}

func readPairsFromCSV(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rows, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return nil, err
	}
	if len(rows) < 2 {
		return nil, nil
	}

	header := make(map[string]int, len(rows[0]))
	for i, col := range rows[0] {
		header[strings.TrimSpace(col)] = i
	}

	set := make(map[string]struct{})
	for _, row := range rows[1:] {
		if len(strings.TrimSpace(strings.Join(row, ""))) == 0 {
			continue
		}

		for _, key := range []string{"Leg1Symbol", "Leg2Symbol", "Leg3Symbol"} {
			if idx, ok := header[key]; ok && idx < len(row) {
				symbol := strings.ToUpper(strings.TrimSpace(row[idx]))
				if symbol != "" {
					set[symbol] = struct{}{}
				}
			}
		}

		if len(set) == 0 {
			for _, key := range []string{"Leg1", "Leg2", "Leg3"} {
				if idx, ok := header[key]; ok && idx < len(row) {
					if p := parseLeg(row[idx]); p != "" {
						set[p] = struct{}{}
					}
				}
			}
		}
	}

	res := make([]string, 0, len(set))
	for k := range set {
		res = append(res, k)
	}
	return res, nil
}

func parseLeg(s string) string {
	parts := strings.Fields(strings.ToUpper(strings.TrimSpace(s)))
	if len(parts) < 2 {
		return ""
	}
	p := strings.Split(parts[1], "/")
	if len(p) != 2 {
		return ""
	}
	return p[0] + "-" + p[1]
}


