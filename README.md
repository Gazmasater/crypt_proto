sudo systemctl mask sleep.target suspend.target hibernate.target hybrid-sleep.target



wbs-api.mexc.com/ws 


[https://edis-global.vercel.app/ru/vps-hosting/singapore-singapore
](https://sg.edisglobal.com/)



git pull --rebase origin privat
git push origin privat


BOOK_INTERVAL=100ms
SYMBOLS_FILE=triangles_markets.csv
DEBUG=false


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









Структура каталога та же (crypt_proto), внутри:

main.go — только запуск и pprof

config.go — конфиг + debug-логгер

domain.go — типы и работа с треугольниками

proto_decoder.go — разбор protobuf от MEXC

ws.go — работа с WebSocket

arb.go — расчёт треугольников, логирование, консюмер

Все файлы ниже можно просто создать рядом и вставить как есть.

config.go
package main

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

/* =========================  CONFIG  ========================= */

type Config struct {
	TrianglesFile string
	BookInterval  string
	FeePerLeg     float64 // комиссия за одну ногу, доля: 0.0004 = 0.04%
	MinProfit     float64 // минимальная прибыль за круг, доля: 0.005 = 0.5%
	Debug         bool
}

func loadEnvFloat(name string, def float64) float64 {
	raw := strings.TrimSpace(os.Getenv(name))
	if raw == "" {
		return def
	}
	v, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		log.Printf("bad %s=%q: %v, using default %f", name, raw, err, def)
		return def
	}
	return v
}

func loadConfig() Config {
	_ = godotenv.Load(".env")

	tf := os.Getenv("TRIANGLES_FILE")
	if tf == "" {
		tf = "triangles_markets.csv"
	}

	bi := os.Getenv("BOOK_INTERVAL")
	if bi == "" {
		bi = "100ms"
	}

	feePct := loadEnvFloat("FEE_PCT", 0.04)       // проценты
	minPct := loadEnvFloat("MIN_PROFIT_PCT", 0.5) // проценты

	debugFlag := strings.EqualFold(os.Getenv("DEBUG"), "true")

	cfg := Config{
		TrianglesFile: tf,
		BookInterval:  bi,
		FeePerLeg:     feePct / 100.0,
		MinProfit:     minPct / 100.0,
		Debug:         debugFlag,
	}

	log.Printf("Triangles file: %s", cfg.TrianglesFile)
	log.Printf("Book interval: %s", cfg.BookInterval)
	log.Printf("Fee per leg: %.4f %% (rate=%.6f)", feePct, cfg.FeePerLeg)
	log.Printf("Min profit per cycle: %.4f %% (rate=%.6f)", minPct, cfg.MinProfit)

	return cfg
}

/* =========================  LOGGING  ========================= */

var debug bool

func dlog(format string, args ...any) {
	if debug {
		log.Printf(format, args...)
	}
}

domain.go
package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

/* =========================  DOMAIN TYPES  ========================= */

type Leg struct {
	From   string
	To     string
	Symbol string
	Dir    int8 // +1: From->To = base->quote; -1: From->To = quote->base
}

type Triangle struct {
	Legs [3]Leg
	Name string // A→B→C→A
}

type Quote struct {
	Bid    float64
	Ask    float64
	BidQty float64
	AskQty float64
}

type Event struct {
	Symbol string
	Bid    float64
	Ask    float64
	BidQty float64
	AskQty float64
}

type Pair struct {
	Base   string
	Quote  string
	Symbol string
}

/* =========================  TRIANGLES  ========================= */

func buildTriangleFromPairs(p1, p2, p3 Pair) (Triangle, bool) {
	currencies := map[string]struct{}{
		p1.Base:  {},
		p1.Quote: {},
		p2.Base:  {},
		p2.Quote: {},
		p3.Base:  {},
		p3.Quote: {},
	}
	if len(currencies) != 3 {
		return Triangle{}, false
	}

	currs := make([]string, 0, 3)
	for c := range currencies {
		currs = append(currs, c)
	}

	type edge struct {
		From, To string
	}

	pairs := []Pair{p1, p2, p3}
	perm3 := [][]int{
		{0, 1, 2},
		{0, 2, 1},
		{1, 0, 2},
		{1, 2, 0},
		{2, 0, 1},
		{2, 1, 0},
	}

	for _, order := range perm3 {
		c0, c1, c2 := currs[order[0]], currs[order[1]], currs[order[2]]
		edges := []edge{
			{From: c0, To: c1},
			{From: c1, To: c2},
			{From: c2, To: c0},
		}

		for _, pp := range perm3 {
			var legs [3]Leg
			okAll := true

			for i := 0; i < 3; i++ {
				e := edges[i]
				p := pairs[pp[i]]

				switch {
				case p.Base == e.From && p.Quote == e.To:
					legs[i] = Leg{From: e.From, To: e.To, Symbol: p.Symbol, Dir: +1}
				case p.Base == e.To && p.Quote == e.From:
					legs[i] = Leg{From: e.From, To: e.To, Symbol: p.Symbol, Dir: -1}
				default:
					okAll = false
				}
				if !okAll {
					break
				}
			}

			if okAll {
				name := fmt.Sprintf("%s→%s→%s→%s", edges[0].From, edges[1].From, edges[2].From, edges[0].From)
				return Triangle{Legs: legs, Name: name}, true
			}
		}
	}

	return Triangle{}, false
}

func loadTriangles(path string) ([]Triangle, []string, map[string][]int, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, nil, err
	}
	defer f.Close()

	r := csv.NewReader(bufio.NewReader(f))
	r.TrimLeadingSpace = true
	r.Comma = ','

	var (
		tris      []Triangle
		symbolSet = make(map[string]struct{})
	)

	for {
		rec, err := r.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, nil, nil, err
		}

		var fields []string
		for _, v := range rec {
			v = strings.TrimSpace(v)
			if v != "" {
				fields = append(fields, v)
			}
		}
		if len(fields) == 0 {
			continue
		}
		if strings.HasPrefix(fields[0], "#") {
			continue
		}
		if len(fields) != 6 {
			log.Printf("skip line (need 6 fields): %v", fields)
			continue
		}

		p1 := Pair{Base: fields[0], Quote: fields[1], Symbol: fields[0] + fields[1]}
		p2 := Pair{Base: fields[2], Quote: fields[3], Symbol: fields[2] + fields[3]}
		p3 := Pair{Base: fields[4], Quote: fields[5], Symbol: fields[4] + fields[5]}

		t, ok := buildTriangleFromPairs(p1, p2, p3)
		if !ok {
			continue
		}

		tris = append(tris, t)
		for _, leg := range t.Legs {
			symbolSet[leg.Symbol] = struct{}{}
		}
	}

	symbols := make([]string, 0, len(symbolSet))
	for s := range symbolSet {
		symbols = append(symbols, s)
	}

	index := make(map[string][]int)
	for i, t := range tris {
		for _, leg := range t.Legs {
			index[leg.Symbol] = append(index[leg.Symbol], i)
		}
	}

	log.Printf("треугольников всего: %d", len(tris))
	log.Printf("символов в индексе треугольников: %d", len(symbols))

	return tris, symbols, index, nil
}

proto_decoder.go
package main

import (
	"strconv"
	"strings"
	"sync"

	"google.golang.org/protobuf/proto"

	pb "crypt_proto/pb"
)

/* =========================  PROTO DECODER  ========================= */

var wrapperPool = sync.Pool{
	New: func() any { return new(pb.PushDataV3ApiWrapper) },
}

func parsePBQuote(raw []byte) (string, Quote, bool) {
	w, _ := wrapperPool.Get().(*pb.PushDataV3ApiWrapper)
	defer func() {
		*w = pb.PushDataV3ApiWrapper{}
		wrapperPool.Put(w)
	}()

	if err := proto.Unmarshal(raw, w); err != nil {
		return "", Quote{}, false
	}

	sym := w.GetSymbol()
	if sym == "" {
		ch := w.GetChannel()
		if i := strings.LastIndex(ch, "@"); i >= 0 && i+1 < len(ch) {
			sym = ch[i+1:]
		}
	}
	if sym == "" {
		return "", Quote{}, false
	}

	// PublicBookTicker
	if b1, ok := w.GetBody().(*pb.PushDataV3ApiWrapper_PublicBookTicker); ok && b1.PublicBookTicker != nil {
		t := b1.PublicBookTicker
		return parseQuoteFromStrings(sym, t.GetBidPrice(), t.GetAskPrice(), "", "")
	}

	// PublicAggreBookTicker
	if b2, ok := w.GetBody().(*pb.PushDataV3ApiWrapper_PublicAggreBookTicker); ok && b2.PublicAggreBookTicker != nil {
		t := b2.PublicAggreBookTicker
		return parseQuoteFromStrings(
			sym,
			t.GetBidPrice(),
			t.GetAskPrice(),
			t.GetBidQuantity(),
			t.GetAskQuantity(),
		)
	}

	return "", Quote{}, false
}

func parseQuoteFromStrings(sym, bp, ap, bq, aq string) (string, Quote, bool) {
	if bp == "" || ap == "" {
		return "", Quote{}, false
	}

	bid, err1 := strconv.ParseFloat(bp, 64)
	ask, err2 := strconv.ParseFloat(ap, 64)
	if err1 != nil || err2 != nil || bid <= 0 || ask <= 0 {
		return "", Quote{}, false
	}

	var bidQty, askQty float64
	if bq != "" {
		if v, err := strconv.ParseFloat(bq, 64); err == nil {
			bidQty = v
		}
	}
	if aq != "" {
		if v, err := strconv.ParseFloat(aq, 64); err == nil {
			askQty = v
		}
	}

	return sym, Quote{
		Bid:    bid,
		Ask:    ask,
		BidQty: bidQty,
		AskQty: askQty,
	}, true
}

ws.go
package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

/* =========================  WS SUBSCRIBER  ========================= */

func buildTopics(symbols []string, interval string) []string {
	topics := make([]string, 0, len(symbols))
	for _, s := range symbols {
		topics = append(topics, "spot@public.aggre.bookTicker.v3.api.pb@"+interval+"@"+s)
	}
	return topics
}

func runPublicBookTickerWS(
	ctx context.Context,
	wg WaitGroupLike,
	connID int,
	symbols []string,
	interval string,
	out chan<- Event,
) {
	defer wg.Done()

	const (
		baseRetry = 2 * time.Second
		maxRetry  = 30 * time.Second
	)

	urlWS := "wss://wbs-api.mexc.com/ws"
	topics := buildTopics(symbols, interval)
	retry := baseRetry

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		conn, _, err := websocket.DefaultDialer.Dial(urlWS, nil)
		if err != nil {
			log.Printf("[WS #%d] dial err: %v (retry in %v)", connID, err, retry)
			time.Sleep(retry)
			retry = nextRetry(retry, maxRetry)
			continue
		}

		log.Printf("[WS #%d] connected to %s (symbols: %d)", connID, urlWS, len(symbols))
		retry = baseRetry

		_ = conn.SetReadDeadline(time.Now().Add(90 * time.Second))

		var lastPing time.Time
		conn.SetPongHandler(func(appData string) error {
			rtt := time.Since(lastPing)
			dlog("[WS #%d] Pong через %v", connID, rtt)
			return conn.SetReadDeadline(time.Now().Add(90 * time.Second))
		})

		stopPing := make(chan struct{})
		go pingLoop(connID, conn, &lastPing, stopPing)

		if err := sendSubscription(conn, topics, connID); err != nil {
			close(stopPing)
			_ = conn.Close()
			time.Sleep(retry)
			retry = nextRetry(retry, maxRetry)
			continue
		}

		if !readLoop(ctx, connID, conn, out) {
			close(stopPing)
			_ = conn.Close()
			time.Sleep(retry)
			retry = nextRetry(retry, maxRetry)
			continue
		}
	}
}

func nextRetry(cur, max time.Duration) time.Duration {
	cur *= 2
	if cur > max {
		return max
	}
	return cur
}

func pingLoop(connID int, conn *websocket.Conn, lastPing *time.Time, stop <-chan struct{}) {
	t := time.NewTicker(45 * time.Second)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			*lastPing = time.Now()
			if err := conn.WriteControl(
				websocket.PingMessage,
				[]byte("hb"),
				time.Now().Add(5*time.Second),
			); err != nil {
				dlog("[WS #%d] ping error: %v", connID, err)
				return
			}
		case <-stop:
			return
		}
	}
}

func sendSubscription(conn *websocket.Conn, topics []string, connID int) error {
	sub := map[string]any{
		"method": "SUBSCRIPTION",
		"params": topics,
		"id":     time.Now().Unix(),
	}
	if err := conn.WriteJSON(sub); err != nil {
		log.Printf("[WS #%d] subscribe send err: %v", connID, err)
		return err
	}
	log.Printf("[WS #%d] SUB -> %d topics", connID, len(topics))
	return nil
}

func readLoop(ctx context.Context, connID int, conn *websocket.Conn, out chan<- Event) bool {
	for {
		mt, raw, err := conn.ReadMessage()
		if err != nil {
			log.Printf("[WS #%d] read err: %v (reconnect)", connID, err)
			return false
		}

		switch mt {
		case websocket.TextMessage:
			handleTextMessage(connID, raw)
		case websocket.BinaryMessage:
			sym, q, ok := parsePBQuote(raw)
			if !ok {
				continue
			}
			ev := Event{
				Symbol: sym,
				Bid:    q.Bid,
				Ask:    q.Ask,
				BidQty: q.BidQty,
				AskQty: q.AskQty,
			}
			select {
			case out <- ev:
			case <-ctx.Done():
				return true
			}
		default:
			// игнорируем прочие типы
		}
	}
}

func handleTextMessage(connID int, raw []byte) {
	if !debug {
		return
	}
	var tmp any
	if err := json.Unmarshal(raw, &tmp); err == nil {
		j, _ := json.Marshal(tmp)
		dlog("[WS #%d TEXT] %s", connID, string(j))
	} else {
		dlog("[WS #%d TEXT RAW] %s", connID, string(raw))
	}
}


⚠️ В ws.go я использовал интерфейс WaitGroupLike, чтобы можно было тестировать без sync.WaitGroup. Для простоты можешь убрать интерфейс и вернуться к *sync.WaitGroup. Ниже в arb.go будет реализация, которая ждёт именно *sync.WaitGroup. Если не хочешь усложнять, поменяй сигнатуру runPublicBookTickerWS обратно на wg *sync.WaitGroup.

Если не хочешь интерфейс, сразу смотри упрощённый вариант ниже в arb.go и просто исправь сигнатуру как раньше.

arb.go
package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

/* =========================  TRIANGLE EVAL + PRINT  ========================= */

func evalTriangle(t Triangle, quotes map[string]Quote, fee float64) (float64, bool) {
	amt := 1.0

	for _, leg := range t.Legs {
		q, ok := quotes[leg.Symbol]
		if !ok || q.Bid <= 0 || q.Ask <= 0 {
			return 0, false
		}

		if leg.Dir > 0 {
			amt *= q.Bid
		} else {
			amt /= q.Ask
		}

		amt *= (1 - fee)
		if amt <= 0 {
			return 0, false
		}
	}

	return amt - 1.0, true
}

func printTriangle(w io.Writer, t Triangle, profit float64, quotes map[string]Quote) {
	ts := time.Now().Format("2006-01-02 15:04:05.000")
	fmt.Fprintf(w, "%s\n", ts)
	fmt.Fprintf(w, "[ARB] %+0.3f%%  %s\n", profit*100, t.Name)

	for _, leg := range t.Legs {
		q := quotes[leg.Symbol]
		mid := (q.Bid + q.Ask) / 2
		spreadAbs := q.Ask - q.Bid
		spreadPct := 0.0
		if mid > 0 {
			spreadPct = spreadAbs / mid * 100
		}
		side := fmt.Sprintf("%s/%s", leg.From, leg.To)
		if leg.Dir < 0 {
			side = fmt.Sprintf("%s/%s", leg.To, leg.From)
		}

		fmt.Fprintf(
			w,
			"  %s (%s): bid=%.10f ask=%.10f  spread=%.10f (%.5f%%)  bidQty=%.4f askQty=%.4f\n",
			leg.Symbol, side,
			q.Bid, q.Ask,
			spreadAbs, spreadPct,
			q.BidQty, q.AskQty,
		)
	}
	fmt.Fprintln(w)
}

/* =========================  ARB LOGGING + PIPELINE  ========================= */

func initArbLogger(path string) (io.Writer, func()) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		log.Fatalf("open %s: %v", path, err)
	}
	buf := bufio.NewWriter(f)

	out := io.MultiWriter(os.Stdout, buf)

	cleanup := func() {
		_ = buf.Flush()
		_ = f.Close()
	}

	return out, cleanup
}

func startWSWorkers(
	ctx context.Context,
	wg *sync.WaitGroup,
	symbols []string,
	interval string,
	out chan<- Event,
) {
	const maxPerConn = 50

	var chunks [][]string
	for i := 0; i < len(symbols); i += maxPerConn {
		j := i + maxPerConn
		if j > len(symbols) {
			j = len(symbols)
		}
		chunks = append(chunks, symbols[i:j])
	}

	log.Printf("будем использовать %d WS-подключений", len(chunks))

	for idx, chunk := range chunks {
		wg.Add(1)
		go runPublicBookTickerWS(ctx, wg, idx, chunk, interval, out)
	}
}

func consumeEvents(
	ctx context.Context,
	events <-chan Event,
	triangles []Triangle,
	indexBySymbol map[string][]int,
	feePerLeg, minProfit float64,
	out io.Writer,
) {
	quotes := make(map[string]Quote)

	const minPrintInterval = 5 * time.Millisecond
	lastPrint := make(map[int]time.Time)

	for {
		select {
		case ev, ok := <-events:
			if !ok {
				return
			}

			// если котировка по символу не изменилась, не пересчитываем
			if prev, okPrev := quotes[ev.Symbol]; okPrev &&
				prev.Bid == ev.Bid &&
				prev.Ask == ev.Ask &&
				prev.BidQty == ev.BidQty &&
				prev.AskQty == ev.AskQty {
				continue
			}

			quotes[ev.Symbol] = Quote{
				Bid:    ev.Bid,
				Ask:    ev.Ask,
				BidQty: ev.BidQty,
				AskQty: ev.AskQty,
			}

			trIDs := indexBySymbol[ev.Symbol]
			if len(trIDs) == 0 {
				continue
			}

			now := time.Now()
			for _, id := range trIDs {
				tr := triangles[id]
				prof, ok := evalTriangle(tr, quotes, feePerLeg)
				if !ok || prof < minProfit {
					continue
				}

				if last, okLast := lastPrint[id]; okLast && now.Sub(last) < minPrintInterval {
					continue
				}
				lastPrint[id] = now

				printTriangle(out, tr, prof, quotes)
			}
		case <-ctx.Done():
			return
		}
	}
}


Обрати внимание: здесь runPublicBookTickerWS ожидает wg *sync.WaitGroup, поэтому в ws.go сигнатуру функции тоже делай с wg *sync.WaitGroup, как у тебя было изначально.

main.go
package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"sync"
	"syscall"
	"time"

	_ "net/http/pprof"
)

/* =========================  MAIN / APP  ========================= */

func main() {
	// pprof
	go func() {
		log.Println("pprof on http://localhost:6060/debug/pprof/")
		if err := http.ListenAndServe("localhost:6060", nil); err != nil {
			log.Printf("pprof server error: %v", err)
		}
	}()

	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	cfg := loadConfig()
	debug = cfg.Debug

	arbOut, closeArb := initArbLogger("arbitrage.log")
	defer closeArb()

	triangles, symbols, indexBySymbol, err := loadTriangles(cfg.TrianglesFile)
	if err != nil {
		log.Fatalf("load triangles: %v", err)
	}
	if len(triangles) == 0 {
		log.Fatal("нет треугольников, нечего мониторить")
	}
	if len(symbols) == 0 {
		log.Fatal("нет символов для подписки")
	}
	log.Printf("символов для подписки всего: %d", len(symbols))

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	events := make(chan Event, 8192)

	var wg sync.WaitGroup
	startWSWorkers(ctx, &wg, symbols, cfg.BookInterval, events)

	go consumeEvents(ctx, events, triangles, indexBySymbol, cfg.FeePerLeg, cfg.MinProfit, arbOut)

	<-ctx.Done()
	log.Println("shutting down...")

	time.Sleep(200 * time.Millisecond)
	close(events)
	wg.Wait()
	log.Println("bye")
}


Если хочешь, дальше можем:

добавить флаг -json, чтобы вместо текста писать JSON-лог арбитражей;

сделать лимит по количеству одновременных «живых» треугольников;

добавить метрики (кол-во тиков, найденных арбов, средняя прибыль и т.д.).




