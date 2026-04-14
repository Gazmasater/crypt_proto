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
