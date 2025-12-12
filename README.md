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








Окей, давай сделаем нормальную архитектуру с интерфейсом биржи и сразу заложим KuCoin. Ниже — готовый каркас проекта: MarketDataFeed интерфейс, реализация для MEXC, заглушка под KuCoin, отдельные пакеты config, domain, arb, exchange, плюс cmd/cryptarb/main.go.
Я буду исходить из того, что корень репо — это твой crypt_proto.

Структура проекта
crypt_proto/
  cmd/
    cryptarb/
      main.go

  config/
    config.go

  domain/
    domain.go

  arb/
    consumer.go

  exchange/
    exchange.go

  mexc/
    feed.go

  kucoin/
    feed.go

  pb/
    ... (твои protobuf-генерации)

  triangles_markets.csv
  arbitrage.log   // создаётся автоматически
  .env


1. config/config.go
package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Exchange      string  // "MEXC" или "KUCOIN"
	TrianglesFile string
	BookInterval  string
	FeePerLeg     float64 // как доля, 0.001 = 0.1%
	MinProfit     float64 // как доля, 0.003 = 0.3%
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

func Load() Config {
	_ = godotenv.Load(".env")

	ex := strings.ToUpper(strings.TrimSpace(os.Getenv("EXCHANGE")))
	if ex == "" {
		ex = "MEXC"
	}

	tf := os.Getenv("TRIANGLES_FILE")
	if tf == "" {
		tf = "triangles_markets.csv"
	}
	bi := os.Getenv("BOOK_INTERVAL")
	if bi == "" {
		bi = "100ms"
	}

	// проценты
	feePct := loadEnvFloat("FEE_PCT", 0.04)
	minPct := loadEnvFloat("MIN_PROFIT_PCT", 0.5)

	debug := strings.ToLower(os.Getenv("DEBUG")) == "true"

	cfg := Config{
		Exchange:      ex,
		TrianglesFile: tf,
		BookInterval:  bi,
		FeePerLeg:     feePct / 100.0,
		MinProfit:     minPct / 100.0,
		Debug:         debug,
	}

	log.Printf("Exchange: %s", cfg.Exchange)
	log.Printf("Triangles file: %s", tf)
	log.Printf("Book interval: %s", bi)
	log.Printf("Fee per leg: %.4f %% (rate=%.6f)", feePct, cfg.FeePerLeg)
	log.Printf("Min profit per cycle: %.4f %% (rate=%.6f)", minPct, cfg.MinProfit)

	return cfg
}

.env например:
EXCHANGE=MEXC
TRIANGLES_FILE=triangles_markets.csv
BOOK_INTERVAL=10ms
FEE_PCT=0.04
MIN_PROFIT_PCT=0.5
DEBUG=true


2. domain/domain.go — общие типы и треугольники
package domain

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
)

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

func buildTriangleFromPairs(p1, p2, p3 Pair) (Triangle, bool) {
	set := map[string]struct{}{
		p1.Base:  {},
		p1.Quote: {},
		p2.Base:  {},
		p2.Quote: {},
		p3.Base:  {},
		p3.Quote: {},
	}
	if len(set) != 3 {
		return Triangle{}, false
	}
	currs := make([]string, 0, 3)
	for c := range set {
		currs = append(currs, c)
	}

	type edge struct{ From, To string }

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

// LoadTriangles читает CSV, строит треугольники и индекс по символам.
func LoadTriangles(path string) ([]Triangle, []string, map[string][]int, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, nil, err
	}
	defer f.Close()

	r := csv.NewReader(bufio.NewReader(f))
	r.TrimLeadingSpace = true
	r.Comma = ','

	var tris []Triangle
	symbolSet := make(map[string]struct{})

	for {
		rec, err := r.Read()
		if err != nil {
			if err.Error() == "EOF" {
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

// EvalTriangle считает доходность треугольника.
func EvalTriangle(t Triangle, quotes map[string]Quote, fee float64) (float64, bool) {
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


3. arb/consumer.go — потребитель и лог
package arb

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"crypt_proto/domain"
)

type Consumer struct {
	FeePerLeg float64
	MinProfit float64

	writer io.Writer
}

func NewConsumer(feePerLeg, minProfit float64, out io.Writer) *Consumer {
	return &Consumer{
		FeePerLeg: feePerLeg,
		MinProfit: minProfit,
		writer:    out,
	}
}

// Start запускает горутину-потребителя.
func (c *Consumer) Start(
	ctx context.Context,
	events <-chan domain.Event,
	triangles []domain.Triangle,
	indexBySymbol map[string][]int,
	wg *sync.WaitGroup,
) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		c.run(ctx, events, triangles, indexBySymbol)
	}()
}

func (c *Consumer) run(
	ctx context.Context,
	events <-chan domain.Event,
	triangles []domain.Triangle,
	indexBySymbol map[string][]int,
) {
	quotes := make(map[string]domain.Quote)

	const minPrintInterval = 5 * time.Millisecond
	lastPrint := make(map[int]time.Time)

	for {
		select {
		case ev, ok := <-events:
			if !ok {
				return
			}

			if prev, okPrev := quotes[ev.Symbol]; okPrev &&
				prev.Bid == ev.Bid &&
				prev.Ask == ev.Ask &&
				prev.BidQty == ev.BidQty &&
				prev.AskQty == ev.AskQty {
				continue
			}

			quotes[ev.Symbol] = domain.Quote{
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
				prof, ok := domain.EvalTriangle(tr, quotes, c.FeePerLeg)
				if !ok {
					continue
				}
				if prof >= c.MinProfit {
					if last, okLast := lastPrint[id]; okLast {
						if now.Sub(last) < minPrintInterval {
							continue
						}
					}
					lastPrint[id] = now
					c.printTriangle(now, tr, prof, quotes)
				}
			}

		case <-ctx.Done():
			return
		}
	}
}

func (c *Consumer) printTriangle(
	ts time.Time,
	t domain.Triangle,
	profit float64,
	quotes map[string]domain.Quote,
) {
	w := c.writer
	fmt.Fprintf(w, "%s\n", ts.Format("2006-01-02 15:04:05.000"))
	fmt.Fprintf(w, "[ARB] %+0.3f%%  %s\n", profit*100, t.Name)

	for _, leg := range t.Legs {
		q := quotes[leg.Symbol]
		mid := (q.Bid + q.Ask) / 2
		spreadAbs := q.Ask - q.Bid
		spreadPct := 0.0
		if mid > 0 {
			spreadPct = spreadAbs / mid * 100
		}
		side := ""
		if leg.Dir > 0 {
			side = fmt.Sprintf("%s/%s", leg.From, leg.To)
		} else {
			side = fmt.Sprintf("%s/%s", leg.To, leg.From)
		}
		fmt.Fprintf(w, "  %s (%s): bid=%.10f ask=%.10f  spread=%.10f (%.5f%%)  bidQty=%.4f askQty=%.4f\n",
			leg.Symbol, side,
			q.Bid, q.Ask,
			spreadAbs, spreadPct,
			q.BidQty, q.AskQty,
		)
	}
	fmt.Fprintln(w)
}

func OpenLogWriter(path string) (io.WriteCloser, *bufio.Writer, io.Writer) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		log.Fatalf("open %s: %v", path, err)
	}
	buf := bufio.NewWriter(f)
	out := io.MultiWriter(os.Stdout, buf)
	return f, buf, out
}


4. exchange/exchange.go — интерфейс биржи
package exchange

import (
	"context"
	"sync"

	"crypt_proto/domain"
)

// MarketDataFeed — общий интерфейс для биржи.
type MarketDataFeed interface {
	Name() string
	Start(ctx context.Context, wg *sync.WaitGroup, symbols []string, interval string, out chan<- domain.Event)
}


5. mexc/feed.go — реализация для MEXC
Это по сути твой текущий код, обёрнутый в Feed и использующий cfg.Debug как флаг.
package mexc

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"crypt_proto/domain"
	pb "crypt_proto/pb"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

type Feed struct {
	debug bool
}

func NewFeed(debug bool) *Feed {
	return &Feed{debug: debug}
}

func (f *Feed) Name() string { return "MEXC" }

/* ===== proto decoder ===== */

var wrapperPool = sync.Pool{
	New: func() any { return new(pb.PushDataV3ApiWrapper) },
}

func parsePBQuote(raw []byte) (string, domain.Quote, bool) {
	w, _ := wrapperPool.Get().(*pb.PushDataV3ApiWrapper)
	defer func() {
		*w = pb.PushDataV3ApiWrapper{}
		wrapperPool.Put(w)
	}()

	if err := proto.Unmarshal(raw, w); err != nil {
		return "", domain.Quote{}, false
	}

	sym := w.GetSymbol()
	if sym == "" {
		ch := w.GetChannel()
		if i := strings.LastIndex(ch, "@"); i >= 0 && i+1 < len(ch) {
			sym = ch[i+1:]
		}
	}
	if sym == "" {
		return "", domain.Quote{}, false
	}

	if b1, ok := w.GetBody().(*pb.PushDataV3ApiWrapper_PublicBookTicker); ok && b1.PublicBookTicker != nil {
		t := b1.PublicBookTicker
		bp := t.GetBidPrice()
		ap := t.GetAskPrice()
		if bp == "" || ap == "" {
			return "", domain.Quote{}, false
		}
		bid, err1 := strconv.ParseFloat(bp, 64)
		ask, err2 := strconv.ParseFloat(ap, 64)
		if err1 != nil || err2 != nil || bid <= 0 || ask <= 0 {
			return "", domain.Quote{}, false
		}
		return sym, domain.Quote{
			Bid:    bid,
			Ask:    ask,
			BidQty: 0,
			AskQty: 0,
		}, true
	}

	if b2, ok := w.GetBody().(*pb.PushDataV3ApiWrapper_PublicAggreBookTicker); ok && b2.PublicAggreBookTicker != nil {
		t := b2.PublicAggreBookTicker

		bp := t.GetBidPrice()
		ap := t.GetAskPrice()
		bq := t.GetBidQuantity()
		aq := t.GetAskQuantity()

		if bp == "" || ap == "" {
			return "", domain.Quote{}, false
		}
		bid, err1 := strconv.ParseFloat(bp, 64)
		ask, err2 := strconv.ParseFloat(ap, 64)
		if err1 != nil || err2 != nil || bid <= 0 || ask <= 0 {
			return "", domain.Quote{}, false
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

		return sym, domain.Quote{
			Bid:    bid,
			Ask:    ask,
			BidQty: bidQty,
			AskQty: askQty,
		}, true
	}

	return "", domain.Quote{}, false
}

func (f *Feed) dlog(format string, args ...any) {
	if f.debug {
		log.Printf(format, args...)
	}
}

/* ===== WS ===== */

func (f *Feed) runPublicBookTickerWS(
	ctx context.Context,
	wg *sync.WaitGroup,
	connID int,
	symbols []string,
	interval string,
	out chan<- domain.Event,
) {
	defer wg.Done()

	const (
		baseRetry = 2 * time.Second
		maxRetry  = 30 * time.Second
	)

	urlWS := "wss://wbs-api.mexc.com/ws"

	topics := make([]string, 0, len(symbols))
	for _, s := range symbols {
		topics = append(topics, "spot@public.aggre.bookTicker.v3.api.pb@"+interval+"@"+s)
	}

	retry := baseRetry

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		conn, _, err := websocket.DefaultDialer.Dial(urlWS, nil)
		if err != nil {
			log.Printf("[MEXC WS #%d] dial err: %v (retry in %v)", connID, err, retry)
			time.Sleep(retry)
			if retry < maxRetry {
				retry *= 2
				if retry > maxRetry {
					retry = maxRetry
				}
			}
			continue
		}
		log.Printf("[MEXC WS #%d] connected to %s (symbols: %d)", connID, urlWS, len(symbols))
		retry = baseRetry

		_ = conn.SetReadDeadline(time.Now().Add(90 * time.Second))

		var lastPing time.Time
		conn.SetPongHandler(func(appData string) error {
			rtt := time.Since(lastPing)
			f.dlog("[MEXC WS #%d] Pong через %v", connID, rtt)
			return conn.SetReadDeadline(time.Now().Add(90 * time.Second))
		})

		stopPing := make(chan struct{})
		go func() {
			t := time.NewTicker(45 * time.Second)
			defer t.Stop()
			for {
				select {
				case <-t.C:
					lastPing = time.Now()
					if err := conn.WriteControl(websocket.PingMessage, []byte("hb"), time.Now().Add(5*time.Second)); err != nil {
						f.dlog("[MEXC WS #%d] ping error: %v", connID, err)
						return
					}
				case <-stopPing:
					return
				}
			}
		}()

		sub := map[string]any{
			"method": "SUBSCRIPTION",
			"params": topics,
			"id":     time.Now().Unix(),
		}
		if err := conn.WriteJSON(sub); err != nil {
			log.Printf("[MEXC WS #%d] subscribe send err: %v", connID, err)
			close(stopPing)
			_ = conn.Close()
			time.Sleep(retry)
			continue
		}
		log.Printf("[MEXC WS #%d] SUB -> %d topics", connID, len(topics))

		for {
			mt, raw, err := conn.ReadMessage()
			if err != nil {
				log.Printf("[MEXC WS #%d] read err: %v (reconnect)", connID, err)
				break
			}

			switch mt {
			case websocket.TextMessage:
				if f.debug {
					var tmp any
					if err := json.Unmarshal(raw, &tmp); err == nil {
						j, _ := json.Marshal(tmp)
						f.dlog("[MEXC #%d TEXT] %s", connID, string(j))
					} else {
						f.dlog("[MEXC #%d TEXT RAW] %s", connID, string(raw))
					}
				}
			case websocket.BinaryMessage:
				sym, q, ok := parsePBQuote(raw)
				if !ok {
					continue
				}
				ev := domain.Event{
					Symbol: sym,
					Bid:    q.Bid,
					Ask:    q.Ask,
					BidQty: q.BidQty,
					AskQty: q.AskQty,
				}
				select {
				case out <- ev:
				case <-ctx.Done():
					close(stopPing)
					_ = conn.Close()
					return
				}
			default:
			}
		}

		close(stopPing)
		_ = conn.Close()
		time.Sleep(retry)
		if retry < maxRetry {
			retry *= 2
			if retry > maxRetry {
				retry = maxRetry
			}
		}
	}
}

// Start реализует интерфейс MarketDataFeed.
func (f *Feed) Start(
	ctx context.Context,
	wg *sync.WaitGroup,
	symbols []string,
	interval string,
	out chan<- domain.Event,
) {
	const maxPerConn = 50
	chunks := make([][]string, 0)
	for i := 0; i < len(symbols); i += maxPerConn {
		j := i + maxPerConn
		if j > len(symbols) {
			j = len(symbols)
		}
		chunks = append(chunks, symbols[i:j])
	}
	log.Printf("[MEXC] будем использовать %d WS-подключений", len(chunks))

	for idx, chunk := range chunks {
		wg.Add(1)
		go f.runPublicBookTickerWS(ctx, wg, idx, chunk, interval, out)
	}
}


6. kucoin/feed.go — заглушка под KuCoin
Тут только каркас. Позже ты добавишь реальную подписку по их протоколу.
package kucoin

import (
	"context"
	"log"
	"sync"

	"crypt_proto/domain"
)

type Feed struct {
	debug bool
}

func NewFeed(debug bool) *Feed { return &Feed{debug: debug} }

func (f *Feed) Name() string { return "KuCoin" }

func (f *Feed) Start(
	ctx context.Context,
	wg *sync.WaitGroup,
	symbols []string,
	interval string,
	out chan<- domain.Event,
) {
	// TODO: здесь нужно реализовать KuCoin WebSocket подписку
	// Сейчас просто логируем, чтобы не падало.
	log.Printf("[KuCoin] Start called, but not implemented yet. symbols=%d interval=%s", len(symbols), interval)
}


7. cmd/cryptarb/main.go — входная точка
package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"crypt_proto/arb"
	"crypt_proto/config"
	"crypt_proto/domain"
	"crypt_proto/exchange"
	"crypt_proto/kucoin"
	"crypt_proto/mexc"

	_ "net/http/pprof"
)

func main() {
	// pprof
	go func() {
		log.Println("pprof on http://localhost:6060/debug/pprof/")
		if err := http.ListenAndServe("localhost:6060", nil); err != nil {
			log.Printf("pprof server error: %v", err)
		}
	}()

	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	cfg := config.Load()

	triangles, symbols, indexBySymbol, err := domain.LoadTriangles(cfg.TrianglesFile)
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

	var feed exchange.MarketDataFeed

	switch cfg.Exchange {
	case "MEXC":
		feed = mexc.NewFeed(cfg.Debug)
	case "KUCOIN":
		feed = kucoin.NewFeed(cfg.Debug)
	default:
		log.Fatalf("unknown EXCHANGE=%q (expected MEXC or KUCOIN)", cfg.Exchange)
	}

	log.Printf("Using exchange: %s", feed.Name())

	// лог-файл для арбитража
	logFile, logBuf, arbOut := arb.OpenLogWriter("arbitrage.log")
	defer logFile.Close()
	defer logBuf.Flush()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	events := make(chan domain.Event, 8192)

	var wg sync.WaitGroup

	// запускаем потребителя
	consumer := arb.NewConsumer(cfg.FeePerLeg, cfg.MinProfit, arbOut)
	consumer.Start(ctx, events, triangles, indexBySymbol, &wg)

	// запускаем фид биржи
	feed.Start(ctx, &wg, symbols, cfg.BookInterval, events)

	// ждём сигнал
	<-ctx.Done()
	log.Println("shutting down...")

	time.Sleep(200 * time.Millisecond)
	close(events)
	wg.Wait()
	log.Println("bye")
}


Что это даёт


Интерфейс MarketDataFeed — main стал полностью отвязан от конкретной биржи.


MEXC реализован в mexc.Feed, KuCoin уже имеет каркас.


Ошибка undefined: debug исчезает, потому что debug теперь поле Feed, а глобальной переменной нет.


Логика поиска треугольников и логирования выделена в domain и arb, код стал “чище”.



Если хочешь, дальше можем:


Вместо заглушки в kucoin.Feed.Start реально подключить KuCoin WS (у них другой протокол — авторизация на /api/v1/bullet-public, topic’и типа /market/ticker:BTC-USDT и т. п.).


Добавить ещё один интерфейс для “исполнителя сделок” (Trader) под приватные ордера.







