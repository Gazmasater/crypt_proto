git rm --cached cmd/arb/metrics/arb_metrics.csv

echo "cmd/arb/metrics/*.csv" >> .gitignore

git add .gitignore
git commit --amend --no-edit


git push origin new_arh --force



git filter-branch --force --index-filter \
'git rm --cached --ignore-unmatch cmd/arb/metrics/arb_metrics.csv' \
--prune-empty --tag-name-filter cat -- new_arh


rm -rf .git/refs/original/
git reflog expire --expire=now --all
git gc --prune=now --aggressive


git push origin new_arh --force






internal/collector/book_snapshot.go — новый
internal/collector/kucoin_collector.go
internal/calculator/arb.go
internal/calculator/depth.go — новый
internal/executor/executor.go
cmd/arb/main.go



book_snapshot.go


package collector

import "sort"

type BookLevel struct {
	Price float64
	Size  float64
}

type BookSnapshot struct {
	Symbol      string
	Bids        []BookLevel
	Asks        []BookLevel
	Sequence    int64
	TimestampMS int64
}

type DepthBookSource interface {
	GetBookSnapshot(symbol string, depth int) (BookSnapshot, bool)
}

func buildSnapshot(symbol string, book *bookState, depth int) (BookSnapshot, bool) {
	if book == nil {
		return BookSnapshot{}, false
	}

	book.mu.Lock()
	defer book.mu.Unlock()

	if !book.ready || book.bestBid <= 0 || book.bestAsk <= 0 {
		return BookSnapshot{}, false
	}

	bids := make([]BookLevel, 0, len(book.bids))
	for p, s := range book.bids {
		if p > 0 && s > 0 {
			bids = append(bids, BookLevel{Price: p, Size: s})
		}
	}
	asks := make([]BookLevel, 0, len(book.asks))
	for p, s := range book.asks {
		if p > 0 && s > 0 {
			asks = append(tasks, BookLevel{Price: p, Size: s})
		}
	}
	if len(bids) == 0 || len(tasks) == 0 {
		return BookSnapshot{}, false
	}

	sort.Slice(bids, func(i, j int) bool { return bids[i].Price > bids[j].Price })
	sort.Slice(asks, func(i, j int) bool { return asks[i].Price < asks[j].Price })

	if depth > 0 {
		if len(bids) > depth {
			bids = bids[:depth]
		}
		if len(asks) > depth {
			asks = asks[:depth]
		}
	}

	return BookSnapshot{
		Symbol:      symbol,
		Bids:        bids,
		Asks:        asks,
		Sequence:    book.sequence,
		TimestampMS: book.lastEventTS,
	}, true
}


collector.go

package collector

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
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
	readTimeout    = 45 * time.Second

	snapshotTimeout = 10 * time.Second
	httpTimeout     = 10 * time.Second
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

type priceLevel struct {
	Price float64
	Size  float64
}

type bookDelta struct {
	SequenceStart int64
	SequenceEnd   int64
	Bids          []priceLevel
	Asks          []priceLevel
	TimestampMS   int64
}

type bookState struct {
	mu sync.Mutex

	bids map[float64]float64
	asks map[float64]float64

	bestBid     float64
	bestBidSize float64
	bestAsk     float64
	bestAskSize float64

	sequence int64
	last     Last

	ready       bool
	needResync  bool
	buffering   bool
	buffered    []bookDelta
	lastEventTS int64
}

type kucoinWS struct {
	id      int
	conn    *websocket.Conn
	symbols []string
	books   map[string]*bookState
	writeMu sync.Mutex
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

		ws := &kucoinWS{
			id:      len(wsList),
			symbols: symbols[i:end],
			books:   make(map[string]*bookState, end-i),
		}
		for _, symbol := range symbols[i:end] {
			ws.books[symbol] = &bookState{
				bids:      make(map[float64]float64),
				asks:      make(map[float64]float64),
				buffering: true,
				buffered:  make([]bookDelta, 0, 64),
			}
		}

		wsList = append(wsList, ws)
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
		go ws.run(c)
	}
	log.Printf("[KuCoin] started with %d WS\n", len(c.wsList))
	return nil
}

func (c *KuCoinCollector) Stop() error {
	c.cancel()
	for _, ws := range c.wsList {
		ws.closeConn()
	}
	return nil
}

func (c *KuCoinCollector) GetBookSnapshot(symbol string, depth int) (BookSnapshot, bool) {
	for _, ws := range c.wsList {
		if book, ok := ws.books[symbol]; ok {
			return buildSnapshot(symbol, book, depth)
		}
	}
	return BookSnapshot{}, false
}

func (ws *kucoinWS) run(c *KuCoinCollector) {
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
		}

		if err := ws.connect(); err != nil {
			log.Printf("[KuCoin WS %d] connect error: %v, retry in %v\n", ws.id, err, reconnectDelay)
			time.Sleep(reconnectDelay)
			continue
		}

		ws.resetBooksForReconnect()

		connDone := make(chan struct{})
		go ws.pingLoop(c.ctx, connDone)

		readErrCh := make(chan error, 1)
		go func() {
			readErrCh <- ws.readLoop(c)
		}()

		if err := ws.subscribeAll(c.ctx, connDone); err != nil {
			close(connDone)
			ws.closeConn()
			log.Printf("[KuCoin WS %d] subscribe error: %v\n", ws.id, err)
			time.Sleep(reconnectDelay)
			continue
		}

		if err := ws.bootstrapAllBooks(c.ctx); err != nil {
			close(connDone)
			ws.closeConn()
			log.Printf("[KuCoin WS %d] bootstrap error: %v\n", ws.id, err)
			time.Sleep(reconnectDelay)
			continue
		}

		err := <-readErrCh
		close(connDone)
		ws.closeConn()

		if err != nil {
			log.Printf("[KuCoin WS %d] readLoop ended: %v\n", ws.id, err)
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
	if len(r.Data.InstanceServers) == 0 {
		return fmt.Errorf("no instance servers")
	}

	url := fmt.Sprintf("%s?token=%s&connectId=%d", r.Data.InstanceServers[0].Endpoint, r.Data.Token, time.Now().UnixNano())
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return err
	}

	_ = conn.SetReadDeadline(time.Now().Add(readTimeout))
	conn.SetPongHandler(func(string) error {
		return conn.SetReadDeadline(time.Now().Add(readTimeout))
	})

	ws.conn = conn
	log.Printf("[KuCoin WS %d] connected\n", ws.id)
	return nil
}

func (ws *kucoinWS) subscribeAll(ctx context.Context, connDone <-chan struct{}) error {
	t := time.NewTicker(subRate)
	defer t.Stop()

	for _, s := range ws.symbols {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-connDone:
			return fmt.Errorf("connection closed")
		case <-t.C:
			if err := ws.writeJSON(map[string]any{
				"id":       time.Now().UnixNano(),
				"type":     "subscribe",
				"topic":    "/market/level2:" + s,
				"response": true,
			}); err != nil {
				return err
			}
		}
	}

	return nil
}

func (ws *kucoinWS) pingLoop(ctx context.Context, connDone <-chan struct{}) {
	t := time.NewTicker(pingInterval)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-connDone:
			return
		case <-t.C:
			if err := ws.writeJSON(map[string]any{
				"id":   time.Now().UnixNano(),
				"type": "ping",
			}); err != nil {
				return
			}
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

		_ = ws.conn.SetReadDeadline(time.Now().Add(readTimeout))
		ws.handle(c, msg)
	}
}

func (ws *kucoinWS) handle(c *KuCoinCollector, msg []byte) {
	const prefix = "/market/level2:"
	const prefixLen = len(prefix)

	msgType := gjson.GetBytes(msg, "type").String()
	if msgType != "message" {
		return
	}

	topicRes := gjson.GetBytes(msg, "topic")
	if !topicRes.Exists() {
		return
	}

	raw := topicRes.Raw
	if len(raw) <= prefixLen+2 {
		return
	}
	if raw[1:1+prefixLen] != prefix {
		return
	}

	symbol := raw[1+prefixLen : len(raw)-1]
	book := ws.books[symbol]
	if book == nil {
		return
	}

	data := gjson.GetBytes(msg, "data")
	if !data.Exists() {
		return
	}

	seqStart := data.Get("sequenceStart").Int()
	seqEnd := data.Get("sequenceEnd").Int()
	if seqStart == 0 || seqEnd == 0 {
		return
	}

	delta := bookDelta{
		SequenceStart: seqStart,
		SequenceEnd:   seqEnd,
		Bids:          parseLevels(data.Get("changes.bids")),
		Asks:          parseLevels(data.Get("changes.asks")),
		TimestampMS:   extractKuCoinTimestampMS(msg, time.Now().UnixMilli()),
	}

	var emit *models.MarketData

	book.mu.Lock()
	book.lastEventTS = delta.TimestampMS

	if book.buffering {
		book.buffered = append(book.buffered, delta)
		book.mu.Unlock()
		return
	}

	if book.needResync || !book.ready {
		book.mu.Unlock()
		return
	}

	if delta.SequenceEnd <= book.sequence {
		book.mu.Unlock()
		return
	}

	if delta.SequenceStart > book.sequence+1 {
		book.ready = false
		book.needResync = true
		book.mu.Unlock()

		log.Printf("[KuCoin WS %d] gap detected symbol=%s localSeq=%d seqStart=%d seqEnd=%d\n",
			ws.id, symbol, book.sequence, delta.SequenceStart, delta.SequenceEnd)

		go ws.resyncBook(c.ctx, symbol)
		return
	}

	applyDelta(book, delta)
	book.sequence = delta.SequenceEnd

	top := Last{
		Bid:     book.bestBid,
		Ask:     book.bestAsk,
		BidSize: book.bestBidSize,
		AskSize: book.bestAskSize,
	}

	if top.Bid > 0 && top.Ask > 0 && top.BidSize > 0 && top.AskSize > 0 && top != book.last {
		book.last = top
		emit = &models.MarketData{
			Exchange:  "KuCoin",
			Symbol:    symbol,
			Bid:       top.Bid,
			Ask:       top.Ask,
			BidSize:   top.BidSize,
			AskSize:   top.AskSize,
			Timestamp: delta.TimestampMS,
		}
	}

	book.mu.Unlock()

	if emit != nil {
		select {
		case c.out <- emit:
		case <-c.ctx.Done():
		}
	}
}

func (ws *kucoinWS) bootstrapAllBooks(ctx context.Context) error {
	const workers = 8

	start := time.Now()
	client := &http.Client{Timeout: httpTimeout}

	jobs := make(chan string, len(ws.symbols))
	doneCh := make(chan string, len(ws.symbols))
	errCh := make(chan error, 1)

	var wg sync.WaitGroup

	worker := func() {
		defer wg.Done()
		for symbol := range jobs {
			if ctx.Err() != nil {
				return
			}
			if err := ws.bootstrapBook(ctx, client, symbol); err != nil {
				select {
				case errCh <- fmt.Errorf("bootstrap %s: %w", symbol, err):
				default:
				}
				return
			}
			select {
			case doneCh <- symbol:
			case <-ctx.Done():
				return
			}
		}
	}

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go worker()
	}

	for _, symbol := range ws.symbols {
		jobs <- symbol
	}
	close(jobs)

	total := len(ws.symbols)
	completed := 0

	waitCh := make(chan struct{})
	go func() {
		wg.Wait()
		close(waitCh)
	}()

	for completed < total {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-errCh:
			return err
		case symbol := <-doneCh:
			completed++
			if completed%10 == 0 || completed == total {
				log.Printf("[KuCoin WS %d] bootstrap progress %d/%d last=%s\n",
					ws.id, completed, total, symbol)
			}
		case <-waitCh:
			if completed == total {
				log.Printf("[KuCoin WS %d] bootstrap finished %d/%d in %v\n",
					ws.id, completed, total, time.Since(start))
				return nil
			}
			return fmt.Errorf("bootstrap stopped early: %d/%d", completed, total)
		}
	}

	log.Printf("[KuCoin WS %d] bootstrap finished %d/%d in %v\n",
		ws.id, completed, total, time.Since(start))
	return nil
}

func (ws *kucoinWS) bootstrapBook(ctx context.Context, client *http.Client, symbol string) error {
	book := ws.books[symbol]
	if book == nil {
		return fmt.Errorf("book not found for symbol=%s", symbol)
	}

	snapshot, err := fetchSnapshot(ctx, client, symbol)
	if err != nil {
		return err
	}

	book.mu.Lock()
	defer book.mu.Unlock()

	book.bids = make(map[float64]float64, len(snapshot.Bids))
	book.asks = make(map[float64]float64, len(snapshot.Asks))

	for _, lvl := range snapshot.Bids {
		if lvl.Size > 0 {
			book.bids[lvl.Price] = lvl.Size
		}
	}
	for _, lvl := range snapshot.Asks {
		if lvl.Size > 0 {
			book.asks[lvl.Price] = lvl.Size
		}
	}

	book.sequence = snapshot.Sequence

	replayed := make([]bookDelta, 0, len(book.buffered))
	for _, d := range book.buffered {
		if d.SequenceEnd <= snapshot.Sequence {
			continue
		}
		replayed = append(replayed, d)
	}

	sort.Slice(replayed, func(i, j int) bool {
		if replayed[i].SequenceStart == replayed[j].SequenceStart {
			return replayed[i].SequenceEnd < replayed[j].SequenceEnd
		}
		return replayed[i].SequenceStart < replayed[j].SequenceStart
	})

	recomputeBestBid(book)
	recomputeBestAsk(book)

	for _, d := range replayed {
		if d.SequenceEnd <= book.sequence {
			continue
		}
		if d.SequenceStart > book.sequence+1 {
			return fmt.Errorf("bootstrap replay gap symbol=%s localSeq=%d seqStart=%d seqEnd=%d",
				symbol, book.sequence, d.SequenceStart, d.SequenceEnd)
		}
		applyDelta(book, d)
		book.sequence = d.SequenceEnd
	}

	book.last = Last{
		Bid:     book.bestBid,
		Ask:     book.bestAsk,
		BidSize: book.bestBidSize,
		AskSize: book.bestAskSize,
	}
	book.ready = book.last.Bid > 0 && book.last.Ask > 0 && book.last.BidSize > 0 && book.last.AskSize > 0
	book.needResync = false
	book.buffering = false
	book.buffered = nil

	return nil
}

func (ws *kucoinWS) resyncBook(ctx context.Context, symbol string) {
	client := &http.Client{Timeout: httpTimeout}
	book := ws.books[symbol]
	if book == nil {
		return
	}

	book.mu.Lock()
	if book.buffering {
		book.mu.Unlock()
		return
	}
	book.buffering = true
	book.buffered = book.buffered[:0]
	book.mu.Unlock()

	resyncCtx, cancel := context.WithTimeout(ctx, snapshotTimeout)
	defer cancel()

	if err := ws.bootstrapBook(resyncCtx, client, symbol); err != nil {
		log.Printf("[KuCoin WS %d] resync failed symbol=%s: %v\n", ws.id, symbol, err)

		book.mu.Lock()
		book.ready = false
		book.needResync = true
		book.buffering = false
		book.mu.Unlock()
		return
	}

	log.Printf("[KuCoin WS %d] resync complete symbol=%s\n", ws.id, symbol)
}

type snapshotBook struct {
	Sequence int64
	Bids     []priceLevel
	Asks     []priceLevel
}

func fetchSnapshot(ctx context.Context, client *http.Client, symbol string) (*snapshotBook, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"https://api.kucoin.com/api/v1/market/orderbook/level2_100?symbol="+symbol,
		nil,
	)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("snapshot request symbol=%s: %w", symbol, err)
	}
	defer resp.Body.Close()

	var r struct {
		Code string `json:"code"`
		Data struct {
			Sequence string     `json:"sequence"`
			Bids     [][]string `json:"bids"`
			Asks     [][]string `json:"asks"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, fmt.Errorf("decode snapshot symbol=%s: %w", symbol, err)
	}
	if r.Code != "" && r.Code != "200000" {
		return nil, fmt.Errorf("snapshot bad code symbol=%s code=%s", symbol, r.Code)
	}

	seq, err := strconv.ParseInt(r.Data.Sequence, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("parse snapshot sequence symbol=%s: %w", symbol, err)
	}

	out := &snapshotBook{
		Sequence: seq,
		Bids:     make([]priceLevel, 0, len(r.Data.Bids)),
		Asks:     make([]priceLevel, 0, len(r.Data.Asks)),
	}

	for _, row := range r.Data.Bids {
		if len(row) < 2 {
			continue
		}
		price, err1 := strconv.ParseFloat(row[0], 64)
		size, err2 := strconv.ParseFloat(row[1], 64)
		if err1 != nil || err2 != nil || price <= 0 || size <= 0 {
			continue
		}
		out.Bids = append(out.Bids, priceLevel{Price: price, Size: size})
	}

	for _, row := range r.Data.Asks {
		if len(row) < 2 {
			continue
		}
		price, err1 := strconv.ParseFloat(row[0], 64)
		size, err2 := strconv.ParseFloat(row[1], 64)
		if err1 != nil || err2 != nil || price <= 0 || size <= 0 {
			continue
		}
		out.Asks = append(out.Asks, priceLevel{Price: price, Size: size})
	}

	return out, nil
}

func parseLevels(arr gjson.Result) []priceLevel {
	if !arr.Exists() || !arr.IsArray() {
		return nil
	}

	raw := arr.Array()
	out := make([]priceLevel, 0, len(raw))

	for _, lvl := range raw {
		parts := lvl.Array()
		if len(parts) < 2 {
			continue
		}

		price, err1 := strconv.ParseFloat(parts[0].String(), 64)
		size, err2 := strconv.ParseFloat(parts[1].String(), 64)
		if err1 != nil || err2 != nil || price <= 0 {
			continue
		}

		out = append(out, priceLevel{
			Price: price,
			Size:  size,
		})
	}

	return out
}

func applyDelta(book *bookState, delta bookDelta) {
	for _, lvl := range delta.Bids {
		applyBidLevel(book, lvl.Price, lvl.Size)
	}
	for _, lvl := range delta.Asks {
		applyAskLevel(book, lvl.Price, lvl.Size)
	}
}

func applyBidLevel(book *bookState, price, size float64) {
	if price <= 0 {
		return
	}

	if size <= 0 {
		delete(book.bids, price)
		if price == book.bestBid {
			recomputeBestBid(book)
		}
		return
	}

	book.bids[price] = size

	if price > book.bestBid || book.bestBid == 0 {
		book.bestBid = price
		book.bestBidSize = size
		return
	}

	if price == book.bestBid {
		book.bestBidSize = size
	}
}

func applyAskLevel(book *bookState, price, size float64) {
	if price <= 0 {
		return
	}

	if size <= 0 {
		delete(book.asks, price)
		if price == book.bestAsk {
			recomputeBestAsk(book)
		}
		return
	}

	book.asks[price] = size

	if book.bestAsk == 0 || price < book.bestAsk {
		book.bestAsk = price
		book.bestAskSize = size
		return
	}

	if price == book.bestAsk {
		book.bestAskSize = size
	}
}

func recomputeBestBid(book *bookState) {
	bestPrice := 0.0
	bestSize := 0.0

	for price, size := range book.bids {
		if size <= 0 {
			continue
		}
		if price > bestPrice {
			bestPrice = price
			bestSize = size
		}
	}

	book.bestBid = bestPrice
	book.bestBidSize = bestSize
}

func recomputeBestAsk(book *bookState) {
	bestPrice := 0.0
	bestSize := 0.0

	for price, size := range book.asks {
		if size <= 0 {
			continue
		}
		if bestPrice == 0 || price < bestPrice {
			bestPrice = price
			bestSize = size
		}
	}

	book.bestAsk = bestPrice
	book.bestAskSize = bestSize
}

func (ws *kucoinWS) resetBooksForReconnect() {
	for _, symbol := range ws.symbols {
		book := ws.books[symbol]
		if book == nil {
			continue
		}

		book.mu.Lock()
		book.bids = make(map[float64]float64)
		book.asks = make(map[float64]float64)
		book.bestBid = 0
		book.bestBidSize = 0
		book.bestAsk = 0
		book.bestAskSize = 0
		book.sequence = 0
		book.last = Last{}
		book.ready = false
		book.needResync = false
		book.buffering = true
		book.buffered = book.buffered[:0]
		book.lastEventTS = 0
		book.mu.Unlock()
	}
}

func extractKuCoinTimestampMS(msg []byte, fallback int64) int64 {
	for _, path := range []string{"data.time", "data.ts", "data.timestamp", "ts", "time"} {
		v := gjson.GetBytes(msg, path)
		if !v.Exists() {
			continue
		}
		ts := parseTSMillis(v)
		if ts > 0 {
			return ts
		}
	}
	return fallback
}

func parseTSMillis(v gjson.Result) int64 {
	s := strings.TrimSpace(v.String())
	if s == "" {
		return 0
	}
	n, err := timeFromNumericString(s)
	if err == nil && n > 0 {
		return n
	}
	return 0
}

func timeFromNumericString(s string) (int64, error) {
	if strings.Contains(s, ".") {
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return 0, err
		}
		return normalizeMillis(int64(f)), nil
	}
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return normalizeMillis(n), nil
}

func normalizeMillis(ts int64) int64 {
	switch {
	case ts <= 0:
		return 0
	case ts >= 1_000_000_000_000_000_000:
		return ts / 1_000_000
	case ts >= 1_000_000_000_000_000:
		return ts / 1_000
	case ts >= 1_000_000_000_000:
		return ts
	case ts >= 1_000_000_000:
		return ts * 1000
	default:
		return 0
	}
}

func (ws *kucoinWS) writeJSON(v any) error {
	ws.writeMu.Lock()
	defer ws.writeMu.Unlock()

	if ws.conn == nil {
		return fmt.Errorf("connection closed")
	}
	return ws.conn.WriteJSON(v)
}

func (ws *kucoinWS) closeConn() {
	ws.writeMu.Lock()
	defer ws.writeMu.Unlock()

	if ws.conn != nil {
		_ = ws.conn.Close()
		ws.conn = nil
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

		foundInRow := false
		for _, key := range []string{"Leg1Symbol", "Leg2Symbol", "Leg3Symbol"} {
			if idx, ok := header[key]; ok && idx < len(row) {
				symbol := strings.ToUpper(strings.TrimSpace(row[idx]))
				if symbol != "" {
					set[symbol] = struct{}{}
					foundInRow = true
				}
			}
		}

		if !foundInRow {
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
	sort.Strings(res)
	return res, nil
}

func parseLeg(raw string) string {
	raw = strings.ToUpper(strings.TrimSpace(raw))
	if raw == "" {
		return ""
	}
	parts := strings.Fields(raw)
	if len(parts) < 2 {
		return ""
	}
	pair := strings.Split(parts[1], "/")
	if len(pair) != 2 {
		return ""
	}
	return strings.TrimSpace(pair[0]) + "-" + strings.TrimSpace(pair[1])
}




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
	"sync"
	"time"

	"crypt_proto/internal/collector"
	"crypt_proto/internal/executor"
	"crypt_proto/internal/queue"
	"crypt_proto/pkg/models"
)

const (
	feeM                 = 0.9992
	defaultMaxQuoteAgeMS = int64(300)
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

type metricsConfig struct {
	MinProfitPct     float64
	MinVolumeUSDT    float64
	MaxAgeMS         int64
	NearProfitPct    float64
	SummaryEvery     time.Duration
	DedupWindow      time.Duration
	FlushEveryWrites int
}

type metricsWriter struct {
	mu          sync.Mutex
	file        *os.File
	csv         *csv.Writer
	enabled     bool
	writeCount  int
	lastKey     string
	lastWriteAt time.Time
	cfg         metricsConfig
}

type metricsSummary struct {
	mu         sync.Mutex
	startedAt  time.Time
	lastLogAt  time.Time
	checked    uint64
	written    uint64
	profitable uint64
	bestPct    float64
	bestUSDT   float64
	bestTri    string
}

func newMetricsSummary() *metricsSummary {
	now := time.Now()
	return &metricsSummary{
		startedAt: now,
		lastLogAt: now,
		bestPct:   math.Inf(-1),
		bestUSDT:  math.Inf(-1),
	}
}

func (ms *metricsSummary) Observe(tri string, profitPct, profitUSDT float64, written bool) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	ms.checked++
	if written {
		ms.written++
	}
	if profitPct > 0 {
		ms.profitable++
	}
	if profitPct > ms.bestPct || (profitPct == ms.bestPct && profitUSDT > ms.bestUSDT) {
		ms.bestPct = profitPct
		ms.bestUSDT = profitUSDT
		ms.bestTri = tri
	}
}

func (ms *metricsSummary) MaybeLog(now time.Time, every time.Duration) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	if every <= 0 || now.Sub(ms.lastLogAt) < every {
		return
	}
	ms.lastLogAt = now

	bestPct := 0.0
	bestUSDT := 0.0
	bestTri := ""
	if ms.checked > 0 && !math.IsInf(ms.bestPct, -1) {
		bestPct = ms.bestPct * 100
		bestUSDT = ms.bestUSDT
		bestTri = ms.bestTri
	}

	log.Printf("[Calculator] summary checked=%d written=%d profitable=%d best_pct=%.4f%% best_usdt=%.6f best_tri=%s\n",
		ms.checked, ms.written, ms.profitable, bestPct, bestUSDT, bestTri)
}

func defaultMetricsConfig() metricsConfig {
	return metricsConfig{
		MinProfitPct:     0.0,
		MinVolumeUSDT:    50.0,
		MaxAgeMS:         defaultMaxQuoteAgeMS,
		NearProfitPct:    -0.0002, // -0.02%
		SummaryEvery:     10 * time.Second,
		DedupWindow:      250 * time.Millisecond,
		FlushEveryWrites: 1,
	}
}

func newMetricsWriter(path string, cfg metricsConfig) (*metricsWriter, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return nil, err
	}

	mw := &metricsWriter{
		file:    f,
		csv:     csv.NewWriter(f),
		enabled: true,
		cfg:     cfg,
	}

	stat, err := f.Stat()
	if err != nil {
		_ = f.Close()
		return nil, err
	}
	if stat.Size() == 0 {
		header := []string{
			"ts_unix_ms", "tri", "A", "B", "C",
			"profit_pct", "profit_usdt", "volume_usdt", "final_usdt",
			"opportunity_strength", "age_min_ms", "age_max_ms", "age_spread_ms",
			"leg1_symbol", "leg1_side", "leg1_bid", "leg1_ask", "leg1_bid_size", "leg1_ask_size", "leg1_age_ms", "leg1_in", "leg1_out", "leg1_trade_qty", "leg1_trade_notional", "leg1_book_limit_in",
			"leg2_symbol", "leg2_side", "leg2_bid", "leg2_ask", "leg2_bid_size", "leg2_ask_size", "leg2_age_ms", "leg2_in", "leg2_out", "leg2_trade_qty", "leg2_trade_notional", "leg2_book_limit_in",
			"leg3_symbol", "leg3_side", "leg3_bid", "leg3_ask", "leg3_bid_size", "leg3_ask_size", "leg3_age_ms", "leg3_in", "leg3_out", "leg3_trade_qty", "leg3_trade_notional", "leg3_book_limit_in",
		}
		if err := mw.csv.Write(header); err != nil {
			_ = f.Close()
			return nil, err
		}
		mw.csv.Flush()
		if err := mw.csv.Error(); err != nil {
			_ = f.Close()
			return nil, err
		}
	}
	return mw, nil
}

func (mw *metricsWriter) Close() error {
	if mw == nil || !mw.enabled {
		return nil
	}
	mw.mu.Lock()
	defer mw.mu.Unlock()

	mw.csv.Flush()
	if err := mw.csv.Error(); err != nil {
		_ = mw.file.Close()
		return err
	}
	return mw.file.Close()
}

// Прибыльные сигналы пишем всегда.
// Для near-profit сигналов применяем фильтры и dedup.
func (mw *metricsWriter) ShouldWrite(anchorTS int64, tri string, profitPct, volumeUSDT float64, maxAgeMS int64) bool {
	if mw == nil || !mw.enabled {
		return false
	}

	// Всегда пишем прибыльные окна.
	if profitPct >= mw.cfg.MinProfitPct {
		return true
	}

	// Ниже — только почти-прибыльные кандидаты.
	if volumeUSDT < mw.cfg.MinVolumeUSDT {
		return false
	}
	if maxAgeMS < 0 || maxAgeMS > mw.cfg.MaxAgeMS {
		return false
	}
	if profitPct < mw.cfg.NearProfitPct {
		return false
	}

	now := time.UnixMilli(anchorTS)
	key := fmt.Sprintf("%s|%.4f|%.2f", tri, profitPct*100, volumeUSDT)

	mw.mu.Lock()
	defer mw.mu.Unlock()

	if mw.lastKey == key && now.Sub(mw.lastWriteAt) < mw.cfg.DedupWindow {
		return false
	}
	mw.lastKey = key
	mw.lastWriteAt = now
	return true
}

func (mw *metricsWriter) Write(record []string) error {
	if mw == nil || !mw.enabled {
		return nil
	}

	mw.mu.Lock()
	defer mw.mu.Unlock()

	if err := mw.csv.Write(record); err != nil {
		return err
	}
	mw.writeCount++

	if mw.cfg.FlushEveryWrites <= 1 || mw.writeCount%mw.cfg.FlushEveryWrites == 0 {
		mw.csv.Flush()
		return mw.csv.Error()
	}
	return nil
}

type Calculator struct {
	mem           *queue.MemoryStore
	bookSource    collector.DepthBookSource
	bySymbol      map[string][]*Triangle
	fileLog       *log.Logger
	metrics       *metricsWriter
	summary       *metricsSummary
	maxQuoteAgeMS int64
	metricsCfg    metricsConfig
	oppOut        chan<- *executor.Opportunity
}

func NewCalculator(mem *queue.MemoryStore, bookSource collector.DepthBookSource, triangles []*Triangle, oppOut chan<- *executor.Opportunity) *Calculator {
	f, err := os.OpenFile("arb_opportunities.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		log.Fatalf("failed to open log: %v", err)
	}

	cfg := defaultMetricsConfig()
	metrics, err := newMetricsWriter("arb_metrics.csv", cfg)
	if err != nil {
		log.Fatalf("failed to open metrics file: %v", err)
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
		mem:           mem,
		bookSource:    bookSource,
		bySymbol:      bySymbol,
		fileLog:       log.New(f, "", log.LstdFlags),
		metrics:       metrics,
		summary:       newMetricsSummary(),
		maxQuoteAgeMS: defaultMaxQuoteAgeMS,
		metricsCfg:    cfg,
		oppOut:        oppOut,
	}
}

func (c *Calculator) Run(in <-chan *models.MarketData) {
	defer func() {
		if c.metrics != nil {
			_ = c.metrics.Close()
		}
	}()

	for md := range in {
		if md == nil {
			continue
		}
		c.mem.Push(md)

		if c.summary != nil {
			c.summary.MaybeLog(time.Now(), c.metricsCfg.SummaryEvery)
		}

		tris := c.bySymbol[md.Symbol]
		if len(tris) == 0 {
			continue
		}

		for _, tri := range tris {
			c.calcTriangle(md, tri)
		}
	}
}

func (c *Calculator) calcTriangle(md *models.MarketData, tri *Triangle) {
	anchorTS := md.Timestamp
	if anchorTS <= 0 {
		return
	}

	var q [3]queue.Quote
	for i, leg := range tri.Legs {
		quote, ok := c.mem.GetLatestBefore(md.Exchange, leg.Symbol, anchorTS, c.maxQuoteAgeMS)
		if !ok {
			return
		}
		if quote.Bid <= 0 || quote.Ask <= 0 || quote.BidSize <= 0 || quote.AskSize <= 0 {
			return
		}
		q[i] = quote
	}

	ages := [3]int64{
		quoteLagMS(anchorTS, q[0]),
		quoteLagMS(anchorTS, q[1]),
		quoteLagMS(anchorTS, q[2]),
	}
	minAge, maxAge, spreadAge := minMaxSpread(ages)

	maxStart, ok := computeMaxStartTopOfBook(tri, q)
	if !ok || maxStart <= 0 {
		return
	}

	if maxStart < 50 {
		return
	}

	finalAmount, diag, ok := simulateTriangle(maxStart, tri, q)
	if !ok || finalAmount <= 0 {
		return
	}

	profitUSDT := finalAmount - maxStart
	profitPct := profitUSDT / maxStart

	if profitPct < 0 {
		return
	}

	var books [3]collector.BookSnapshot
	if c.bookSource != nil {
		for i, leg := range tri.Legs {
			b, ok := c.bookSource.GetBookSnapshot(leg.Symbol, 32)
			if !ok {
				return
			}
			books[i] = b
		}

		depthMaxStart, ok := computeMaxStartByDepth(tri, books)
		if !ok || depthMaxStart <= 0 {
			return
		}
		if depthMaxStart < maxStart {
			maxStart = depthMaxStart
		}
		if maxStart < 50 {
			return
		}

		depthFinal, depthDiag, ok := simulateTriangleDepth(maxStart, tri, books)
		if !ok || depthFinal <= 0 {
			return
		}
		depthProfitUSDT := depthFinal - maxStart
		depthProfitPct := depthProfitUSDT / maxStart
		if depthProfitPct <= 0 {
			return
		}

		finalAmount = depthFinal
		diag = depthDiag
		profitUSDT = depthProfitUSDT
		profitPct = depthProfitPct
	}

	strength := computeOpportunityStrength(profitPct, maxStart, spreadAge, maxAge)
	triName := fmt.Sprintf("%s->%s->%s", tri.A, tri.B, tri.C)

	if c.oppOut != nil {
		op := &executor.Opportunity{
			Exchange: md.Exchange,
			Triangle: toExecutorTriangle(tri),
			AnchorTS: anchorTS,
			Quotes: [3]queue.Quote{
				q[0],
				q[1],
				q[2],
			},
			AgesMS: [3]int64{
				ages[0],
				ages[1],
				ages[2],
			},
			MaxStart:   maxStart,
			Books:      books,
			ProfitPct:  profitPct,
			ProfitUSDT: profitUSDT,
			FinalUSDT:  finalAmount,
		}

		select {
		case c.oppOut <- op:
		default:
		}
	}

	written := false

	shouldWrite := c.metrics != nil && c.metrics.ShouldWrite(anchorTS, triName, profitPct, maxStart, maxAge)
	if shouldWrite {
		record := []string{
			strconv.FormatInt(anchorTS, 10),
			triName,
			tri.A, tri.B, tri.C,
			fmtFloat(profitPct), fmtFloat(profitUSDT), fmtFloat(maxStart), fmtFloat(finalAmount),
			fmtFloat(strength), strconv.FormatInt(minAge, 10), strconv.FormatInt(maxAge, 10), strconv.FormatInt(spreadAge, 10),

			tri.Legs[0].Symbol, tri.Legs[0].Side,
			fmtFloat(q[0].Bid), fmtFloat(q[0].Ask), fmtFloat(q[0].BidSize), fmtFloat(q[0].AskSize),
			strconv.FormatInt(ages[0], 10),
			fmtFloat(diag[0].In), fmtFloat(diag[0].Out), fmtFloat(diag[0].TradeQty), fmtFloat(diag[0].TradeNotional), fmtFloat(diag[0].BookLimitIn),

			tri.Legs[1].Symbol, tri.Legs[1].Side,
			fmtFloat(q[1].Bid), fmtFloat(q[1].Ask), fmtFloat(q[1].BidSize), fmtFloat(q[1].AskSize),
			strconv.FormatInt(ages[1], 10),
			fmtFloat(diag[1].In), fmtFloat(diag[1].Out), fmtFloat(diag[1].TradeQty), fmtFloat(diag[1].TradeNotional), fmtFloat(diag[1].BookLimitIn),

			tri.Legs[2].Symbol, tri.Legs[2].Side,
			fmtFloat(q[2].Bid), fmtFloat(q[2].Ask), fmtFloat(q[2].BidSize), fmtFloat(q[2].AskSize),
			strconv.FormatInt(ages[2], 10),
			fmtFloat(diag[2].In), fmtFloat(diag[2].Out), fmtFloat(diag[2].TradeQty), fmtFloat(diag[2].TradeNotional), fmtFloat(diag[2].BookLimitIn),
		}

		if err := c.metrics.Write(record); err != nil {
			log.Printf("[Calculator] metrics write error: %v", err)
		} else {
			written = true
		}
	}

	if c.summary != nil {
		c.summary.Observe(triName, profitPct, profitUSDT, written)
	}

	if profitPct > 0.0 && maxStart > 50 {
		msg := fmt.Sprintf(
			"[ARB] %s→%s→%s | %.4f%% | volume=%.2f USDT | profit=%.6f USDT | "+
				"anchor=%d | l1=%s %s out=%.8f age=%dms | "+
				"l2=%s %s out=%.8f age=%dms | "+
				"l3=%s %s out=%.8f age=%dms",
			tri.A, tri.B, tri.C,
			profitPct*100, maxStart, profitUSDT,
			anchorTS,
			tri.Legs[0].Symbol, tri.Legs[0].Side, diag[0].Out, ages[0],
			tri.Legs[1].Symbol, tri.Legs[1].Side, diag[1].Out, ages[1],
			tri.Legs[2].Symbol, tri.Legs[2].Side, diag[2].Out, ages[2],
		)
		log.Println(msg)
		c.fileLog.Println(msg)
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

		return outBase, legExecution{
			In:            in,
			Out:           outBase,
			Price:         q.Ask,
			BookLimitIn:   bookLimitIn,
			TradeQty:      tradeQty,
			TradeNotional: tradeNotional,
		}, true

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
		return outQuote, legExecution{
			In:            in,
			Out:           outQuote,
			Price:         q.Bid,
			BookLimitIn:   bookLimitIn,
			TradeQty:      tradeQty,
			TradeNotional: tradeNotional,
		}, true
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

		t := &Triangle{
			A: getString(row, header, "A"),
			B: getString(row, header, "B"),
			C: getString(row, header, "C"),
		}

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

func quoteLagMS(anchorTS int64, q queue.Quote) int64 {
	if q.Timestamp <= 0 || anchorTS <= 0 {
		return -1
	}
	if q.Timestamp > anchorTS {
		return 0
	}
	return anchorTS - q.Timestamp
}

func minMaxSpread(ages [3]int64) (int64, int64, int64) {
	minAge := ages[0]
	maxAge := ages[0]
	for _, age := range ages[1:] {
		if age < minAge {
			minAge = age
		}
		if age > maxAge {
			maxAge = age
		}
	}
	return minAge, maxAge, maxAge - minAge
}

func computeOpportunityStrength(profitPct, volumeUSDT float64, ageSpreadMS, maxAgeMS int64) float64 {
	base := profitPct * math.Log1p(math.Max(volumeUSDT, 0))
	freshPenalty := 1.0 / (1.0 + math.Max(float64(maxAgeMS), 0)/180.0)
	spreadPenalty := 1.0 / (1.0 + math.Max(float64(ageSpreadMS), 0)/180.0)
	return base * freshPenalty * spreadPenalty
}

func fmtFloat(v float64) string {
	return strconv.FormatFloat(v, 'f', 12, 64)
}

func isFinite(v float64) bool { return !math.IsNaN(v) && !math.IsInf(v, 0) }
func eps() float64            { return 1e-12 }

func toExecutorTriangle(t *Triangle) *executor.Triangle {
	if t == nil {
		return nil
	}

	out := &executor.Triangle{
		A: t.A,
		B: t.B,
		C: t.C,
	}

	for i, leg := range t.Legs {
		out.Legs[i] = executor.LegRule{
			Index:       leg.Index,
			RawLeg:      leg.RawLeg,
			Step:        leg.Step,
			MinQty:      leg.MinQty,
			MinNotional: leg.MinNotional,

			Symbol:      leg.Symbol,
			Side:        leg.Side,
			Base:        leg.Base,
			Quote:       leg.Quote,
			QtyStep:     leg.QtyStep,
			QuoteStep:   leg.QuoteStep,
			PriceStep:   leg.PriceStep,
			LegMinQty:   leg.LegMinQty,
			LegMinQuote: leg.LegMinQuote,
			LegMinNotnl: leg.LegMinNotnl,

			Key: leg.Key,
		}
	}

	return out
}



depth.go

package calculator

import (
	"math"
	"strings"

	"crypt_proto/internal/collector"
)

func simulateTriangleDepth(startUSDT float64, tri *Triangle, books [3]collector.BookSnapshot) (float64, [3]legExecution, bool) {
	var diag [3]legExecution
	amount := startUSDT
	for i := 0; i < 3; i++ {
		out, d, ok := executeLegByDepth(amount, tri.Legs[i], books[i])
		if !ok {
			return 0, diag, false
		}
		diag[i] = d
		amount = out
	}
	return amount, diag, true
}

func executeLegByDepth(in float64, leg LegRule, book collector.BookSnapshot) (float64, legExecution, bool) {
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
		if len(book.Asks) == 0 {
			return 0, legExecution{}, false
		}
		remainingQuote := in
		totalQty := 0.0
		totalNotional := 0.0
		bookLimitIn := 0.0
		for _, lvl := range book.Asks {
			if lvl.Price <= 0 || lvl.Size <= 0 {
				continue
			}
			bookLimitIn += lvl.Price * lvl.Size
			if remainingQuote <= eps() {
				break
			}
			maxSpend := lvl.Price * lvl.Size
			spend := math.Min(remainingQuote, maxSpend)
			qty := spend / lvl.Price
			totalQty += qty
			totalNotional += qty * lvl.Price
			remainingQuote -= qty * lvl.Price
		}
		tradeQty := floorToStep(totalQty, qtyStep)
		if tradeQty <= 0 {
			return 0, legExecution{}, false
		}
		avgPrice := totalNotional / totalQty
		tradeNotional := tradeQty * avgPrice
		if tradeNotional <= 0 || tradeNotional > in+1e-9 {
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
		return outBase, legExecution{In: in, Out: outBase, Price: avgPrice, BookLimitIn: bookLimitIn, TradeQty: tradeQty, TradeNotional: tradeNotional}, true
	case "SELL":
		if len(book.Bids) == 0 {
			return 0, legExecution{}, false
		}
		remainingBase := in
		totalQty := 0.0
		totalNotional := 0.0
		bookLimitIn := 0.0
		for _, lvl := range book.Bids {
			if lvl.Price <= 0 || lvl.Size <= 0 {
				continue
			}
			bookLimitIn += lvl.Size
			if remainingBase <= eps() {
				break
			}
			qty := math.Min(remainingBase, lvl.Size)
			totalQty += qty
			totalNotional += qty * lvl.Price
			remainingBase -= qty
		}
		tradeQty := floorToStep(totalQty, qtyStep)
		if tradeQty <= 0 {
			return 0, legExecution{}, false
		}
		avgPrice := totalNotional / totalQty
		tradeNotional := tradeQty * avgPrice
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
		return outQuote, legExecution{In: in, Out: outQuote, Price: avgPrice, BookLimitIn: bookLimitIn, TradeQty: tradeQty, TradeNotional: tradeNotional}, true
	}
	return 0, legExecution{}, false
}

func computeMaxStartByDepth(tri *Triangle, books [3]collector.BookSnapshot) (float64, bool) {
	low, high := 0.0, math.MaxFloat64
	for i := 0; i < 3; i++ {
		leg := tri.Legs[i]
		side := strings.ToUpper(strings.TrimSpace(leg.Side))
		if side == "" {
			side = detectSideFromRawLeg(leg.RawLeg)
		}
		var bookCap float64
		switch side {
		case "BUY":
			for _, lvl := range books[i].Asks {
				bookCap += lvl.Price * lvl.Size
			}
		case "SELL":
			for _, lvl := range books[i].Bids {
				bookCap += lvl.Size
			}
		default:
			return 0, false
		}
		if bookCap <= 0 {
			return 0, false
		}
		if bookCap < high {
			high = bookCap
		}
	}
	if !isFinite(high) || high <= 0 {
		return 0, false
	}
	for i := 0; i < 20; i++ {
		mid := (low + high) / 2
		if mid <= 0 {
			break
		}
		if _, _, ok := simulateTriangleDepth(mid, tri, books); ok {
			low = mid
		} else {
			high = mid
		}
	}
	if low <= 0 {
		return 0, false
	}
	return low, true
}


executer.go

package executor

import (
	"fmt"
	"log"
	"math"
	"os"
	"strings"
	"sync"
	"time"

	"crypt_proto/internal/collector"
	"crypt_proto/internal/queue"
)

const feeM = 0.9992

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

type Opportunity struct {
	Exchange string
	Triangle *Triangle
	AnchorTS int64

	Quotes   [3]queue.Quote
	Books    [3]collector.BookSnapshot
	AgesMS   [3]int64
	MaxStart float64

	ProfitPct  float64
	ProfitUSDT float64
	FinalUSDT  float64
}

type RejectReason string

const (
	RejectDuplicate          RejectReason = "duplicate"
	RejectNilTriangle        RejectReason = "nil_triangle"
	RejectBadAnchor          RejectReason = "bad_anchor"
	RejectBadInputProfit     RejectReason = "bad_input_profit"
	RejectNoQuotes           RejectReason = "no_quotes"
	RejectStale              RejectReason = "stale"
	RejectDesync             RejectReason = "desync"
	RejectSmallMaxStart      RejectReason = "small_max_start"
	RejectSmallSafeStart     RejectReason = "small_safe_start"
	RejectSimulationFailed   RejectReason = "simulation_failed"
	RejectNonPositiveProfit  RejectReason = "non_positive_profit"
	RejectTooSmallProfitUSDT RejectReason = "too_small_profit_usdt"
)

type Config struct {
	SafeStartFactor float64
	MinSafeStart    float64
	MaxQuoteAgeMS   int64
	MaxAgeSpreadMS  int64
	MinProfitUSDT   float64

	// Что вообще пускать в executor из calculator.
	// Всё что ниже этого значения - сразу reject без лишней работы.
	MinInputProfitPct float64

	DedupWindow  time.Duration
	SummaryEvery time.Duration
}

func DefaultConfig() Config {
	return Config{
		SafeStartFactor:   0.70,
		MinSafeStart:      20.0,
		MaxQuoteAgeMS:     300,
		MaxAgeSpreadMS:    200,
		MinProfitUSDT:     0.001,
		MinInputProfitPct: 0.0, // executor по умолчанию ждёт только неотрицательные окна
		DedupWindow:       250 * time.Millisecond,
		SummaryEvery:      10 * time.Second,
	}
}

type Executor struct {
	cfg Config

	mu          sync.Mutex
	lastKey     string
	lastSeenAt  time.Time
	fileLog     *log.Logger
	rejects     map[RejectReason]uint64
	accepted    uint64
	lastSummary time.Time
}

func NewExecutor(cfg Config) *Executor {
	f, err := os.OpenFile("executor_paper.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		log.Fatalf("failed to open executor log: %v", err)
	}

	now := time.Now()

	return &Executor{
		cfg:         cfg,
		fileLog:     log.New(f, "", log.LstdFlags),
		rejects:     make(map[RejectReason]uint64),
		lastSummary: now,
	}
}

func (e *Executor) Handle(op *Opportunity) {
	if op == nil || op.Triangle == nil {
		e.reject(RejectNilTriangle, nil, "")
		return
	}
	if op.AnchorTS <= 0 {
		e.reject(RejectBadAnchor, op, "")
		return
	}

	// Главная правка: executor не должен тратить время на явный минус.
	if op.ProfitPct < e.cfg.MinInputProfitPct {
		e.reject(RejectBadInputProfit, op, fmt.Sprintf("input_profit_pct=%.8f", op.ProfitPct))
		return
	}

	if !e.hasValidQuotes(op) {
		e.reject(RejectNoQuotes, op, "")
		return
	}

	minAge, maxAge, spreadAge := minMaxSpread(op.AgesMS)
	if maxAge < 0 || maxAge > e.cfg.MaxQuoteAgeMS {
		e.reject(RejectStale, op, fmt.Sprintf("max_age=%d", maxAge))
		return
	}
	if spreadAge > e.cfg.MaxAgeSpreadMS {
		e.reject(RejectDesync, op, fmt.Sprintf("spread_age=%d", spreadAge))
		return
	}

	if op.MaxStart <= 0 {
		e.reject(RejectSmallMaxStart, op, "")
		return
	}

	safeStart := floorToReasonable(op.MaxStart*e.cfg.SafeStartFactor, 1e-8)
	if safeStart < e.cfg.MinSafeStart {
		e.reject(RejectSmallSafeStart, op, fmt.Sprintf("safe_start=%.8f", safeStart))
		return
	}

	if e.isDuplicate(op, safeStart) {
		e.reject(RejectDuplicate, op, "")
		return
	}

	finalAmount, diag, ok := simulateTriangle(safeStart, op.Triangle, op.Quotes)
	if ok && hasAllBooks(op.Books) {
		if depthFinal, depthDiag, depthOK := simulateTriangleDepth(safeStart, op.Triangle, op.Books); depthOK && depthFinal > 0 {
			finalAmount = depthFinal
			diag = depthDiag
		} else {
			ok = false
		}
	}
	if !ok || finalAmount <= 0 {
		e.reject(RejectSimulationFailed, op, fmt.Sprintf("safe_start=%.8f", safeStart))
		return
	}

	profitUSDT := finalAmount - safeStart
	profitPct := profitUSDT / safeStart

	if profitPct <= 0 {
		e.reject(RejectNonPositiveProfit, op, fmt.Sprintf("safe_start=%.8f profit_pct=%.8f", safeStart, profitPct))
		return
	}
	if profitUSDT < e.cfg.MinProfitUSDT {
		e.reject(RejectTooSmallProfitUSDT, op, fmt.Sprintf("profit_usdt=%.8f", profitUSDT))
		return
	}

	e.accept(op, safeStart, finalAmount, profitPct, profitUSDT, diag, minAge, maxAge, spreadAge)
}

func hasAllBooks(books [3]collector.BookSnapshot) bool {
	for _, b := range books {
		if len(b.Bids) == 0 || len(b.Asks) == 0 {
			return false
		}
	}
	return true
}

func (e *Executor) hasValidQuotes(op *Opportunity) bool {
	for _, q := range op.Quotes {
		if q.Bid <= 0 || q.Ask <= 0 || q.BidSize <= 0 || q.AskSize <= 0 {
			return false
		}
	}
	return true
}

func (e *Executor) isDuplicate(op *Opportunity, safeStart float64) bool {
	key := fmt.Sprintf("%s->%s->%s|%d|%.4f",
		op.Triangle.A, op.Triangle.B, op.Triangle.C,
		op.AnchorTS/100, // корзина 100мс
		round4(safeStart),
	)

	now := time.UnixMilli(op.AnchorTS)

	e.mu.Lock()
	defer e.mu.Unlock()

	if e.lastKey == key && now.Sub(e.lastSeenAt) < e.cfg.DedupWindow {
		return true
	}
	e.lastKey = key
	e.lastSeenAt = now
	return false
}

func (e *Executor) accept(
	op *Opportunity,
	safeStart float64,
	finalAmount float64,
	profitPct float64,
	profitUSDT float64,
	diag [3]legExecution,
	minAge, maxAge, spreadAge int64,
) {
	msg := fmt.Sprintf(
		"[EXEC PAPER OK] %s→%s→%s | safe=%.2f USDT | final=%.6f | profit=%.6f USDT | %.4f%% | age[min=%d max=%d spread=%d] | l1=%s %s out=%.8f | l2=%s %s out=%.8f | l3=%s %s out=%.8f",
		op.Triangle.A, op.Triangle.B, op.Triangle.C,
		safeStart, finalAmount, profitUSDT, profitPct*100,
		minAge, maxAge, spreadAge,
		op.Triangle.Legs[0].Symbol, op.Triangle.Legs[0].Side, diag[0].Out,
		op.Triangle.Legs[1].Symbol, op.Triangle.Legs[1].Side, diag[1].Out,
		op.Triangle.Legs[2].Symbol, op.Triangle.Legs[2].Side, diag[2].Out,
	)

	log.Println(msg)
	e.fileLog.Println(msg)

	e.mu.Lock()
	e.accepted++
	e.maybeSummaryLocked()
	e.mu.Unlock()
}

func (e *Executor) reject(reason RejectReason, op *Opportunity, extra string) {
	e.mu.Lock()
	e.rejects[reason]++
	e.maybeSummaryLocked()
	e.mu.Unlock()

	if op == nil || op.Triangle == nil {
		return
	}

	// Не шумим в stdout по каждому reject, только в файл.
	msg := fmt.Sprintf("[EXEC PAPER REJECT] reason=%s tri=%s->%s->%s %s",
		reason, op.Triangle.A, op.Triangle.B, op.Triangle.C, extra)
	e.fileLog.Println(msg)
}

func (e *Executor) maybeSummaryLocked() {
	now := time.Now()
	if now.Sub(e.lastSummary) < e.cfg.SummaryEvery {
		return
	}
	e.lastSummary = now

	msg := fmt.Sprintf(
		"[Executor] summary accepted=%d reject_duplicate=%d reject_bad_input_profit=%d reject_stale=%d reject_desync=%d reject_small_safe=%d reject_non_positive=%d reject_too_small_profit=%d reject_sim_failed=%d",
		e.accepted,
		e.rejects[RejectDuplicate],
		e.rejects[RejectBadInputProfit],
		e.rejects[RejectStale],
		e.rejects[RejectDesync],
		e.rejects[RejectSmallSafeStart],
		e.rejects[RejectNonPositiveProfit],
		e.rejects[RejectTooSmallProfitUSDT],
		e.rejects[RejectSimulationFailed],
	)

	log.Println(msg)
	e.fileLog.Println(msg)
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

		return outBase, legExecution{
			In:            in,
			Out:           outBase,
			Price:         q.Ask,
			BookLimitIn:   bookLimitIn,
			TradeQty:      tradeQty,
			TradeNotional: tradeNotional,
		}, true

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

		return outQuote, legExecution{
			In:            in,
			Out:           outQuote,
			Price:         q.Bid,
			BookLimitIn:   bookLimitIn,
			TradeQty:      tradeQty,
			TradeNotional: tradeNotional,
		}, true
	}

	return 0, legExecution{}, false
}

func firstPositive(vals ...float64) float64 {
	for _, v := range vals {
		if v > 0 {
			return v
		}
	}
	return 0
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

func floorToReasonable(v, step float64) float64 {
	if v <= 0 {
		return 0
	}
	if step <= 0 {
		return v
	}
	n := math.Floor(v / step)
	if n <= 0 {
		return 0
	}
	return n * step
}

func minMaxSpread(ages [3]int64) (int64, int64, int64) {
	minAge := ages[0]
	maxAge := ages[0]
	for _, age := range ages[1:] {
		if age < minAge {
			minAge = age
		}
		if age > maxAge {
			maxAge = age
		}
	}
	return minAge, maxAge, maxAge - minAge
}

func round4(v float64) float64 {
	return math.Round(v*1e4) / 1e4
}

func eps() float64 { return 1e-12 }

func simulateTriangleDepth(startUSDT float64, tri *Triangle, books [3]collector.BookSnapshot) (float64, [3]legExecution, bool) {
	var diag [3]legExecution
	amount := startUSDT
	for i := 0; i < 3; i++ {
		out, d, ok := executeLegByDepth(amount, tri.Legs[i], books[i])
		if !ok {
			return 0, diag, false
		}
		diag[i] = d
		amount = out
	}
	return amount, diag, true
}

func executeLegByDepth(in float64, leg LegRule, book collector.BookSnapshot) (float64, legExecution, bool) {
	side := strings.ToUpper(strings.TrimSpace(leg.Side))
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
		if len(book.Asks) == 0 {
			return 0, legExecution{}, false
		}
		remainingQuote := in
		totalQty := 0.0
		totalNotional := 0.0
		bookLimitIn := 0.0
		for _, lvl := range book.Asks {
			if lvl.Price <= 0 || lvl.Size <= 0 {
				continue
			}
			bookLimitIn += lvl.Price * lvl.Size
			if remainingQuote <= eps() {
				break
			}
			maxSpend := lvl.Price * lvl.Size
			spend := math.Min(remainingQuote, maxSpend)
			qty := spend / lvl.Price
			totalQty += qty
			totalNotional += qty * lvl.Price
			remainingQuote -= qty * lvl.Price
		}
		tradeQty := floorToStep(totalQty, qtyStep)
		if tradeQty <= 0 {
			return 0, legExecution{}, false
		}
		avgPrice := totalNotional / totalQty
		tradeNotional := tradeQty * avgPrice
		if tradeNotional <= 0 || tradeNotional > in+1e-9 {
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
		return outBase, legExecution{In: in, Out: outBase, Price: avgPrice, BookLimitIn: bookLimitIn, TradeQty: tradeQty, TradeNotional: tradeNotional}, true
	case "SELL":
		if len(book.Bids) == 0 {
			return 0, legExecution{}, false
		}
		remainingBase := in
		totalQty := 0.0
		totalNotional := 0.0
		bookLimitIn := 0.0
		for _, lvl := range book.Bids {
			if lvl.Price <= 0 || lvl.Size <= 0 {
				continue
			}
			bookLimitIn += lvl.Size
			if remainingBase <= eps() {
				break
			}
			qty := math.Min(remainingBase, lvl.Size)
			totalQty += qty
			totalNotional += qty * lvl.Price
			remainingBase -= qty
		}
		tradeQty := floorToStep(totalQty, qtyStep)
		if tradeQty <= 0 {
			return 0, legExecution{}, false
		}
		avgPrice := totalNotional / totalQty
		tradeNotional := tradeQty * avgPrice
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
		return outQuote, legExecution{In: in, Out: outQuote, Price: avgPrice, BookLimitIn: bookLimitIn, TradeQty: tradeQty, TradeNotional: tradeNotional}, true
	}
	return 0, legExecution{}, false
}



main.go\



package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"crypt_proto/internal/calculator"
	"crypt_proto/internal/collector"
	"crypt_proto/internal/executor"
	"crypt_proto/internal/queue"
	"crypt_proto/pkg/models"
)

func main() {
	go func() {
		log.Println("pprof on http://localhost:6060/debug/pprof/")
		if err := http.ListenAndServe("localhost:6060", nil); err != nil {
			log.Printf("pprof server stopped: %v", err)
		}
	}()

	out := make(chan *models.MarketData, 100_000)
	oppCh := make(chan *executor.Opportunity, 4096)

	mem := queue.NewMemoryStore()
	go mem.Run()

	kc, _, err := collector.NewKuCoinCollectorFromCSV("../exchange/data/kucoin_triangles_usdt.csv")
	if err != nil {
		log.Fatal(err)
	}
	if err := kc.Start(out); err != nil {
		log.Fatal(err)
	}
	log.Println("[Main] KuCoinCollector started")

	triangles, err := calculator.ParseTrianglesFromCSV("../exchange/data/kucoin_triangles_usdt.csv")
	if err != nil {
		log.Fatal(err)
	}

	calc := calculator.NewCalculator(mem, kc, triangles, oppCh)
	go calc.Run(out)

	exec := executor.NewExecutor(executor.DefaultConfig())
	go func() {
		for op := range oppCh {
			exec.Handle(op)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("[Main] shutting down...")

	if err := kc.Stop(); err != nil {
		log.Printf("[Main] collector stop error: %v", err)
	}

	close(out)
	close(oppCh)

	log.Println("[Main] exited")
}



[{
	"resource": "/home/gaz358/myprog/crypt_proto/internal/collector/book_snapshot.go",
	"owner": "_generated_diagnostic_collection_name_#1",
	"code": {
		"value": "UndeclaredName",
		"target": {
			"$mid": 1,
			"path": "/golang.org/x/tools/internal/typesinternal",
			"scheme": "https",
			"authority": "pkg.go.dev",
			"fragment": "UndeclaredName"
		}
	},
	"severity": 8,
	"message": "undefined: tasks",
	"source": "compiler",
	"startLineNumber": 43,
	"startColumn": 18,
	"endLineNumber": 43,
	"endColumn": 23,
	"modelVersionId": 3,
	"origin": "extHost1"
}]
