package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"google.golang.org/protobuf/proto"

	pb "crypt_proto/pb" // твои *.pb.go (MEXC v3)
)

/* =========================  CONFIG  ========================= */

type Config struct {
	BookInterval string // интервал книги: "100ms" / "10ms"
	Debug        bool
	SymbolsFile  string  // CSV с треугольниками (triangles_markets.csv)
	FeePct       float64 // комиссия за сделку, в %, например 0.1
	MinProfitPct float64 // минимальная прибыль по кругу, в %, например 0.3
}

func loadConfig() (Config, error) {
	_ = godotenv.Load(".env")

	cfg := Config{
		BookInterval: os.Getenv("BOOK_INTERVAL"),
		SymbolsFile:  os.Getenv("SYMBOLS_FILE"),
	}

	if cfg.BookInterval == "" {
		cfg.BookInterval = "100ms"
	}
	if cfg.SymbolsFile == "" {
		cfg.SymbolsFile = "triangles_markets.csv"
	}
	if strings.ToLower(os.Getenv("DEBUG")) == "true" {
		cfg.Debug = true
	}

	// Комиссия (в процентах), по умолчанию 0.1%
	if v := os.Getenv("FEE_PCT"); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			cfg.FeePct = f
		}
	}
	if cfg.FeePct <= 0 {
		cfg.FeePct = 0.1
	}

	// Минимальная прибыль по кругу (в процентах), по умолчанию 0.3%
	if v := os.Getenv("MIN_PROFIT_PCT"); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			cfg.MinProfitPct = f
		}
	}
	if cfg.MinProfitPct <= 0 {
		cfg.MinProfitPct = 0.3
	}

	return cfg, nil
}

/* =========================  LOGGING  ========================= */

var debug bool

func dlog(format string, args ...any) {
	if debug {
		log.Printf(format, args...)
	}
}

/* =========================  МАРКЕТЫ / ТРЕУГОЛЬНИКИ  ========================= */

type Market struct {
	Symbol string
	Base   string
	Quote  string
}

type Triangle struct {
	M [3]Market
}

// Загружаем треугольники из CSV и одновременно собираем все символы
func loadTriangles(path string) ([]Triangle, []string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, fmt.Errorf("open %s: %w", path, err)
	}
	defer f.Close()

	r := csv.NewReader(f)

	// заголовок
	if _, err := r.Read(); err != nil {
		return nil, nil, fmt.Errorf("read header: %w", err)
	}

	var tris []Triangle
	seen := make(map[string]struct{})

	for {
		rec, err := r.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, nil, fmt.Errorf("read record: %w", err)
		}
		if len(rec) < 6 {
			continue
		}

		b1, q1 := strings.TrimSpace(rec[0]), strings.TrimSpace(rec[1])
		b2, q2 := strings.TrimSpace(rec[2]), strings.TrimSpace(rec[3])
		b3, q3 := strings.TrimSpace(rec[4]), strings.TrimSpace(rec[5])

		if b1 == "" || q1 == "" || b2 == "" || q2 == "" || b3 == "" || q3 == "" {
			continue
		}

		m1 := Market{Base: b1, Quote: q1, Symbol: b1 + q1}
		m2 := Market{Base: b2, Quote: q2, Symbol: b2 + q2}
		m3 := Market{Base: b3, Quote: q3, Symbol: b3 + q3}

		tris = append(tris, Triangle{M: [3]Market{m1, m2, m3}})

		seen[m1.Symbol] = struct{}{}
		seen[m2.Symbol] = struct{}{}
		seen[m3.Symbol] = struct{}{}
	}

	symbols := make([]string, 0, len(seen))
	for s := range seen {
		symbols = append(symbols, s)
	}

	return tris, symbols, nil
}

/* =========================  CОБЫТИЯ КОТИРОВОК  ========================= */

type Event struct {
	Symbol   string
	Bid, Ask float64
	BidQty   float64
	AskQty   float64
}

type Quote struct {
	Bid, Ask float64
	BidQty   float64
	AskQty   float64
}

/* =========================  PROTO DECODER  ========================= */

// Пул для protobuf-структуры, чтобы не аллоцировать на каждый тик
var wrapperPool = sync.Pool{
	New: func() any { return new(pb.PushDataV3ApiWrapper) },
}

// Возвращаем symbol, bid, ask (объёмы пока 0, позже можно добить из pb)
func parsePBWrapperQuote(raw []byte) (sym string, bid, ask, bidQty, askQty float64, ok bool) {
	w, _ := wrapperPool.Get().(*pb.PushDataV3ApiWrapper)
	defer func() {
		*w = pb.PushDataV3ApiWrapper{}
		wrapperPool.Put(w)
	}()

	if err := proto.Unmarshal(raw, w); err != nil {
		return "", 0, 0, 0, 0, false
	}

	sym = w.GetSymbol()
	if sym == "" {
		ch := w.GetChannel()
		if i := strings.LastIndex(ch, "@"); i >= 0 && i+1 < len(ch) {
			sym = ch[i+1:]
		}
	}
	if sym == "" {
		return "", 0, 0, 0, 0, false
	}

	// PublicBookTicker
	if b1, ok1 := w.GetBody().(*pb.PushDataV3ApiWrapper_PublicBookTicker); ok1 && b1.PublicBookTicker != nil {
		bp := b1.PublicBookTicker.GetBidPrice()
		ap := b1.PublicBookTicker.GetAskPrice()
		if bp == "" || ap == "" {
			return "", 0, 0, 0, 0, false
		}

		bid, err1 := strconv.ParseFloat(bp, 64)
		ask, err2 := strconv.ParseFloat(ap, 64)
		if err1 != nil || err2 != nil || bid <= 0 || ask <= 0 {
			return "", 0, 0, 0, 0, false
		}

		// объёмы пока не парсим — ставим 0
		return sym, bid, ask, 0, 0, true
	}

	// PublicAggreBookTicker
	if b2, ok2 := w.GetBody().(*pb.PushDataV3ApiWrapper_PublicAggreBookTicker); ok2 && b2.PublicAggreBookTicker != nil {
		bp := b2.PublicAggreBookTicker.GetBidPrice()
		ap := b2.PublicAggreBookTicker.GetAskPrice()
		if bp == "" || ap == "" {
			return "", 0, 0, 0, 0, false
		}

		bid, err1 := strconv.ParseFloat(bp, 64)
		ask, err2 := strconv.ParseFloat(ap, 64)
		if err1 != nil || err2 != nil || bid <= 0 || ask <= 0 {
			return "", 0, 0, 0, 0, false
		}

		// объёмы пока не парсим — ставим 0
		return sym, bid, ask, 0, 0, true
	}

	return "", 0, 0, 0, 0, false
}

/* =========================  WS RUNNER (PUBLIC)  ========================= */

func runPublicBookTicker(ctx context.Context, wg *sync.WaitGroup, symbols []string, interval string, out chan<- Event) {
	defer wg.Done()

	if len(symbols) == 0 {
		log.Println("[PUB] нет символов для подписки в этом conn")
		return
	}

	const (
		baseRetry = 2 * time.Second
		maxRetry  = 30 * time.Second
	)

	urlWS := "wss://wbs-api.mexc.com/ws"

	// формируем список топиков
	params := make([]string, 0, len(symbols))
	for _, sym := range symbols {
		topic := "spot@public.aggre.bookTicker.v3.api.pb@" + interval + "@" + sym
		params = append(params, topic)
	}
	log.Printf("[PUB] symbols in this conn: %d", len(params))

	retry := baseRetry

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		conn, _, err := websocket.DefaultDialer.Dial(urlWS, nil)
		if err != nil {
			log.Printf("[PUB] dial err: %v (retry in %v)", err, retry)
			time.Sleep(retry)
			if retry < maxRetry {
				retry *= 2
				if retry > maxRetry {
					retry = maxRetry
				}
			}
			continue
		}
		log.Printf("[PUB] connected to %s", urlWS)
		retry = baseRetry

		_ = conn.SetReadDeadline(time.Now().Add(90 * time.Second))

		var lastPing time.Time
		conn.SetPongHandler(func(appData string) error {
			rtt := time.Since(lastPing)
			dlog("[PING] Pong от %s через %v", urlWS, rtt)
			return conn.SetReadDeadline(time.Now().Add(90 * time.Second))
		})

		// keepalive (PING)
		stopPing := make(chan struct{})
		go func() {
			t := time.NewTicker(45 * time.Second)
			defer t.Stop()
			for {
				select {
				case <-t.C:
					lastPing = time.Now()
					if err := conn.WriteControl(websocket.PingMessage, []byte("hb"), time.Now().Add(5*time.Second)); err != nil {
						dlog("⚠️ [PING] send error: %v", err)
						return
					}
					dlog("[PING] Sent at %s", lastPing.Format("15:04:05.000"))
				case <-stopPing:
					return
				}
			}
		}()

		// подписка на все топики одним запросом
		sub := map[string]any{
			"method": "SUBSCRIPTION",
			"params": params,
			"id":     time.Now().Unix(),
		}
		if err := conn.WriteJSON(sub); err != nil {
			log.Printf("[PUB] subscribe send err: %v", err)
			close(stopPing)
			_ = conn.Close()
			time.Sleep(retry)
			continue
		}
		log.Printf("[PUB] SUB → %d топиков", len(params))

		// цикл чтения
		for {
			mt, raw, err := conn.ReadMessage()
			if err != nil {
				log.Printf("[PUB] read err: %v (reconnect)", err)
				break
			}

			switch mt {
			case websocket.TextMessage:
				// ACK / ошибки подписки — в debug
				var tmp any
				if err := json.Unmarshal(raw, &tmp); err == nil {
					j, _ := json.Marshal(tmp)
					dlog("[PUB TEXT] %s", string(j))
				} else {
					dlog("[PUB TEXT RAW] %s", string(raw))
				}
			case websocket.BinaryMessage:
				if sym, bid, ask, bidQty, askQty, ok := parsePBWrapperQuote(raw); ok {
					ev := Event{
						Symbol: sym,
						Bid:    bid,
						Ask:    ask,
						BidQty: bidQty,
						AskQty: askQty,
					}
					select {
					case out <- ev:
					case <-ctx.Done():
						close(stopPing)
						_ = conn.Close()
						return
					}
				}
			default:
				// игнорируем прочие типы
			}
		}

		// cleanup + реконнект
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

/* =========================  АРБИТРАЖ ПО ТРЕУГОЛЬНИКАМ  ========================= */

var (
	tradeFee  float64 // доля комиссии, например 0.001
	minProfit float64 // минимальная прибыль, доля, например 0.003
)

type TriangleArb struct {
	Triangle Triangle
	StartCur string
	Path     [3]string
	Profit   float64 // множитель, например 1.001 = +0.1%
}

func evalTriangle(t Triangle, quotes map[string]Quote) *TriangleArb {
	// 3 валюты
	curSet := make(map[string]struct{})
	for _, m := range t.M {
		curSet[m.Base] = struct{}{}
		curSet[m.Quote] = struct{}{}
	}
	if len(curSet) != 3 {
		return nil
	}
	curs := make([]string, 0, 3)
	for c := range curSet {
		curs = append(curs, c)
	}

	type key struct{ From, To string }
	rates := make(map[key]float64)

	ff := 1.0 - tradeFee

	for _, m := range t.M {
		q, ok := quotes[m.Symbol]
		if !ok || q.Bid <= 0 || q.Ask <= 0 {
			return nil
		}

		// BASE -> QUOTE: продаём базу, ударяем по bid
		rates[key{m.Base, m.Quote}] = q.Bid * ff

		// QUOTE -> BASE: покупаем базу за котировку, платим ask
		rates[key{m.Quote, m.Base}] = (1.0 / q.Ask) * ff
	}

	a, b, c := curs[0], curs[1], curs[2]

	get := func(from, to string) (float64, bool) {
		r, ok := rates[key{from, to}]
		return r, ok
	}

	// вариант a->b->c->a
	r1, ok1 := get(a, b)
	r2, ok2 := get(b, c)
	r3, ok3 := get(c, a)

	var path [3]string

	if ok1 && ok2 && ok3 {
		profit := r1 * r2 * r3
		if profit > 1.0+minProfit {
			path = [3]string{a + "→" + b, b + "→" + c, c + "→" + a}
			return &TriangleArb{
				Triangle: t,
				StartCur: a,
				Path:     path,
				Profit:   profit,
			}
		}
	}

	// вариант a->c->b->a
	r1, ok1 = get(a, c)
	r2, ok2 = get(c, b)
	r3, ok3 = get(b, a)

	if !(ok1 && ok2 && ok3) {
		return nil
	}

	profit := r1 * r2 * r3
	if profit <= 1.0+minProfit {
		return nil
	}
	path = [3]string{a + "→" + c, c + "→" + b, b + "→" + a}
	return &TriangleArb{
		Triangle: t,
		StartCur: a,
		Path:     path,
		Profit:   profit,
	}
}

func findProfitableTriangles(tris []Triangle, quotes map[string]Quote) []TriangleArb {
	var res []TriangleArb
	for _, t := range tris {
		if arb := evalTriangle(t, quotes); arb != nil {
			res = append(res, *arb)
		}
	}
	return res
}

func printTriangleWithDetails(a TriangleArb, quotes map[string]Quote) {
	perc := (a.Profit - 1.0) * 100.0
	m1, m2, m3 := a.Triangle.M[0], a.Triangle.M[1], a.Triangle.M[2]

	q1 := quotes[m1.Symbol]
	q2 := quotes[m2.Symbol]
	q3 := quotes[m3.Symbol]

	fmt.Printf("\n[ARB] %+0.3f%%  %s\n", perc, strings.Join(a.Path[:], "  "))

	printMarket := func(m Market, q Quote) {
		mid := (q.Bid + q.Ask) / 2
		spreadAbs := q.Ask - q.Bid
		var spreadPct float64
		if mid > 0 {
			spreadPct = spreadAbs / mid * 100
		}
		fmt.Printf("  %s (%s/%s): bid=%.10f ask=%.10f  spread=%.10f (%.5f%%)  bidQty=%.4f askQty=%.4f\n",
			m.Symbol, m.Base, m.Quote,
			q.Bid, q.Ask,
			spreadAbs, spreadPct,
			q.BidQty, q.AskQty,
		)
	}

	printMarket(m1, q1)
	printMarket(m2, q2)
	printMarket(m3, q3)
	fmt.Println()
}

/* =========================  MAIN  ========================= */

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	// pprof HTTP-сервер
	go func() {
		log.Println("pprof on http://localhost:6060/debug/pprof/")
		if err := http.ListenAndServe("localhost:6060", nil); err != nil {
			log.Printf("pprof server error: %v", err)
		}
	}()

	cfg, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}
	debug = cfg.Debug

	tradeFee = cfg.FeePct / 100.0
	minProfit = cfg.MinProfitPct / 100.0

	log.Printf("fee=%.4f%%, minProfit=%.4f%%", cfg.FeePct, cfg.MinProfitPct)

	// грузим треугольники и список символов
	tris, symbols, err := loadTriangles(cfg.SymbolsFile)
	if err != nil {
		log.Fatalf("load triangles: %v", err)
	}
	if len(symbols) == 0 {
		log.Fatal("нет символов в файле ", cfg.SymbolsFile)
	}
	log.Printf("треугольников: %d, символов для подписки: %d", len(tris), len(symbols))

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	events := make(chan Event, 8192)

	var wg sync.WaitGroup

	// ---- Чанкуем по 50 символов на одно WS-подключение ----
	const maxPerConn = 50

	chunks := make([][]string, 0)
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
		go func(i int, syms []string) {
			log.Printf("[WS #%d] symbols in this conn: %d", i, len(syms))
			runPublicBookTicker(ctx, &wg, syms, cfg.BookInterval, events)
		}(idx, chunk)
	}

	go func(tris []Triangle, ctx context.Context) {
		last := make(map[string]Quote)

		// минимальный интервал между выводами, чтобы не зафлудить stdout
		const minLogGap = 200 * time.Millisecond
		nextLogTime := time.Now()

		for {
			select {
			case ev, ok := <-events:
				if !ok {
					return
				}
				// обновили котировку по символу
				last[ev.Symbol] = Quote{
					Bid:    ev.Bid,
					Ask:    ev.Ask,
					BidQty: ev.BidQty,
					AskQty: ev.AskQty,
				}

				// ограничиваем частоту логов
				if time.Now().Before(nextLogTime) {
					continue
				}
				nextLogTime = time.Now().Add(minLogGap)

				prof := findProfitableTriangles(tris, last)
				if len(prof) == 0 {
					// прибыльных нет – молчим
					continue
				}

				fmt.Printf("\nquotes known: %d symbols, profitable triangles: %d\n",
					len(last), len(prof))

				maxShow := 5
				if len(prof) < maxShow {
					maxShow = len(prof)
				}
				for i := 0; i < maxShow; i++ {
					printTriangleWithDetails(prof[i], last)
				}

			case <-ctx.Done():
				return
			}
		}
	}(tris, ctx)

	<-ctx.Done()
	log.Println("ctx done, waiting ws goroutines...")

	// ждём завершения всех WS-подключений
	wg.Wait()
	close(events)

	log.Println("bye")
}
