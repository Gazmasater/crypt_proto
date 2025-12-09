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





BOOK_INTERVAL=100ms
SYMBOLS_FILE=triangles_markets.csv
DEBUG=false

# Комиссия (в процентах, одна нога)
FEE_PCT=0.1          # 0.1% = 0.001

# Минимальная прибыль по кругу (в процентах)
MIN_PROFIT_PCT=0.3   # 0.3% = 0.003



package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"google.golang.org/protobuf/proto"

	pb "crypt_proto/pb"
)

/* =========================  CONFIG  ========================= */

type Config struct {
	TrianglesFile  string  // файл с треугольниками (из Julia)
	BookInterval   string  // "100ms" | "10ms"
	FeePct         float64 // комиссия за сделку, в %
	MinProfitPct   float64 // минимальная прибыль по кругу, в %
	Debug          bool
	MaxSymbolsConn int // максимум символов на одно WS-подключение
}

func loadConfig() (Config, error) {
	_ = godotenv.Load(".env")

	cfg := Config{
		TrianglesFile:  os.Getenv("TRIANGLES_FILE"),
		BookInterval:   os.Getenv("BOOK_INTERVAL"),
		MaxSymbolsConn: 50,
	}

	if cfg.TrianglesFile == "" {
		cfg.TrianglesFile = "triangles_markets.csv"
	}
	if cfg.BookInterval == "" {
		cfg.BookInterval = "100ms"
	}

	// комиссия, % (например, 0.1 для 0.1%)
	if s := os.Getenv("FEE_PCT"); s != "" {
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return cfg, fmt.Errorf("parse FEE_PCT: %w", err)
		}
		cfg.FeePct = f
	} else {
		cfg.FeePct = 0.1
	}

	// минимальная прибыль по кругу, % (например, 0.5)
	if s := os.Getenv("MIN_PROFIT_PCT"); s != "" {
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return cfg, fmt.Errorf("parse MIN_PROFIT_PCT: %w", err)
		}
		cfg.MinProfitPct = f
	} else {
		cfg.MinProfitPct = 0.5
	}

	if strings.ToLower(os.Getenv("DEBUG")) == "true" {
		cfg.Debug = true
	}

	return cfg, nil
}

/* =========================  LOGGING / GLOBALЫ  ========================= */

var (
	debug          bool
	feeRate        float64 // 0.1% -> 0.001
	minProfitRate  float64 // 0.5% -> 0.005
)

func dlog(format string, args ...any) {
	if debug {
		log.Printf(format, args...)
	}
}

/* =========================  MARKET / TRIANGLE TYPES  ========================= */

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

// Треугольник вида:
//   (A/B), (B/C), (A/C)
// где market-символы: M1, M2, M3
type Triangle struct {
	A, B, C string // валюты
	M1      string // A/B -> строка "AB" типа BTCUSDT
	M2      string // B/C
	M3      string // A/C
}

type TriangleArb struct {
	T         Triangle
	Direction string  // "A-B-C-A" или "A-C-B-A"
	ProfitPct float64 // прибыль в %
}

/* =========================  ЗАГРУЗКА ТРЕУГОЛЬНИКОВ  ========================= */

// triangles_markets.csv c видами:
// base1,quote1,base2,quote2,base3,quote3
// пример: ETC,BTC,BTC,USDC,ETC,USDC
func loadTrianglesFromFile(path string) ([]Triangle, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open triangles file: %w", err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.TrimLeadingSpace = true

	var tris []Triangle

	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("read csv: %w", err)
		}
		if len(rec) < 6 {
			continue
		}

		a := strings.TrimSpace(rec[0]) // base1
		b1 := strings.TrimSpace(rec[1]) // quote1 = B
		b2 := strings.TrimSpace(rec[2]) // base2 = B
		c := strings.TrimSpace(rec[3])  // quote2 = C
		a2 := strings.TrimSpace(rec[4]) // base3 = A
		c2 := strings.TrimSpace(rec[5]) // quote3 = C

		if a == "" || b1 == "" || b2 == "" || c == "" || a2 == "" || c2 == "" {
			continue
		}

		// предполагаем, что это ровно A,B,B,C,A,C
		// и рынки:
		//  M1 = A/B
		//  M2 = B/C
		//  M3 = A/C
		m1 := a + b1
		m2 := b2 + c
		m3 := a2 + c2

		t := Triangle{
			A:  a,
			B:  b1,
			C:  c,
			M1: m1,
			M2: m2,
			M3: m3,
		}
		tris = append(tris, t)
	}

	return tris, nil
}

// достаём уникальные символы рынка из треугольников
func extractSymbols(tris []Triangle) []string {
	m := make(map[string]struct{})
	for _, t := range tris {
		if t.M1 != "" {
			m[t.M1] = struct{}{}
		}
		if t.M2 != "" {
			m[t.M2] = struct{}{}
		}
		if t.M3 != "" {
			m[t.M3] = struct{}{}
		}
	}
	out := make([]string, 0, len(m))
	for s := range m {
		out = append(out, s)
	}
	sort.Strings(out)
	return out
}

// индекс: symbol -> какие треугольники его используют
func buildTriangleIndex(tris []Triangle) map[string][]int {
	idx := make(map[string][]int)
	for i, t := range tris {
		for _, s := range []string{t.M1, t.M2, t.M3} {
			if s == "" {
				continue
			}
			idx[s] = append(idx[s], i)
		}
	}
	return idx
}

/* =========================  ПРОТОBUF PARSER (MEXC PB v3)  ========================= */

// Достаём символ и bid/ask из PB-обёртки.
// Поддерживаем PublicBookTicker и PublicAggreBookTicker.
func parsePBWrapperQuote(raw []byte) (sym string, bid, ask float64, ok bool) {
	var w pb.PushDataV3ApiWrapper
	if err := proto.Unmarshal(raw, &w); err != nil {
		return "", 0, 0, false
	}

	sym = w.GetSymbol()
	if sym == "" {
		ch := w.GetChannel()
		if i := strings.LastIndex(ch, "@"); i >= 0 && i+1 < len(ch) {
			sym = ch[i+1:]
		}
	}
	if sym == "" {
		return "", 0, 0, false
	}

	// PublicBookTicker
	if b1, ok1 := w.GetBody().(*pb.PushDataV3ApiWrapper_PublicBookTicker); ok1 && b1.PublicBookTicker != nil {
		bp := b1.PublicBookTicker.GetBidPrice()
		ap := b1.PublicBookTicker.GetAskPrice()
		if bp == "" || ap == "" {
			return "", 0, 0, false
		}
		bidF, err1 := strconv.ParseFloat(bp, 64)
		askF, err2 := strconv.ParseFloat(ap, 64)
		if err1 != nil || err2 != nil || bidF <= 0 || askF <= 0 {
			return "", 0, 0, false
		}
		return sym, bidF, askF, true
	}

	// PublicAggreBookTicker
	if b2, ok2 := w.GetBody().(*pb.PushDataV3ApiWrapper_PublicAggreBookTicker); ok2 && b2.PublicAggreBookTicker != nil {
		bp := b2.PublicAggreBookTicker.GetBidPrice()
		ap := b2.PublicAggreBookTicker.GetAskPrice()
		if bp == "" || ap == "" {
			return "", 0, 0, false
		}
		bidF, err1 := strconv.ParseFloat(bp, 64)
		askF, err2 := strconv.ParseFloat(ap, 64)
		if err1 != nil || err2 != nil || bidF <= 0 || askF <= 0 {
			return "", 0, 0, false
		}
		return sym, bidF, askF, true
	}

	return "", 0, 0, false
}

/* =========================  TRIANGLE EVAL  ========================= */

// Считает прибыль по треугольнику t на основе котировок quotes.
// Используем mid = (bid+ask)/2, комиссия в feeRate, порог в minProfitRate.
// Возвращаем nil, если нет арба.
func evalTriangle(t Triangle, quotes map[string]Quote) *TriangleArb {
	q1, ok1 := quotes[t.M1]
	q2, ok2 := quotes[t.M2]
	q3, ok3 := quotes[t.M3]
	if !ok1 || !ok2 || !ok3 {
		return nil
	}
	if q1.Bid <= 0 || q1.Ask <= 0 || q2.Bid <= 0 || q2.Ask <= 0 || q3.Bid <= 0 || q3.Ask <= 0 {
		return nil
	}

	p1 := (q1.Bid + q1.Ask) / 2 // A/B
	p2 := (q2.Bid + q2.Ask) / 2 // B/C
	p3 := (q3.Bid + q3.Ask) / 2 // A/C

	if p1 <= 0 || p2 <= 0 || p3 <= 0 {
		return nil
	}

	f := 1.0 - feeRate

	// путь 1: A -> B -> C -> A
	a1 := 1.0
	b1Amt := a1 * p1 * f        // A->B
	c1Amt := b1Amt * p2 * f     // B->C
	aFinal1 := c1Amt / p3 * f   // C->A  (1 A = p3 C => 1 C = 1/p3 A)
	profit1 := aFinal1 - 1.0    // в долях

	// путь 2: A -> C -> B -> A
	a2 := 1.0
	c2Amt := a2 * p3 * f        // A->C
	b2Amt := c2Amt / p2 * f     // C->B  (1 B = p2 C => 1 C = 1/p2 B)
	aFinal2 := b2Amt / p1 * f   // B->A  (1 A = p1 B => 1 B = 1/p1 A)
	profit2 := aFinal2 - 1.0

	bestProfit := profit1
	dir := "A-B-C-A"
	if profit2 > bestProfit {
		bestProfit = profit2
		dir = "A-C-B-A"
	}

	if bestProfit <= minProfitRate {
		return nil
	}
	return &TriangleArb{
		T:         t,
		Direction: dir,
		ProfitPct: bestProfit * 100,
	}
}

// Только по подмножеству треугольников (индексы idxs)
func findProfitableTrianglesForSymbol(
	tris []Triangle,
	idxs []int,
	quotes map[string]Quote,
) []TriangleArb {
	res := make([]TriangleArb, 0, len(idxs))
	for _, i := range idxs {
		if arb := evalTriangle(tris[i], quotes); arb != nil {
			res = append(res, *arb)
		}
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].ProfitPct > res[j].ProfitPct
	})
	return res
}

func printTriangleWithDetails(arb TriangleArb, quotes map[string]Quote) {
	t := arb.T
	q1 := quotes[t.M1]
	q2 := quotes[t.M2]
	q3 := quotes[t.M3]

	fmt.Printf("[ARB] %+0.3f%%  %s→%s→%s→%s\n",
		arb.ProfitPct, t.A, t.B, t.C, t.A)

	printQuote := func(label, mkt string, q Quote) {
		if q.Bid <= 0 || q.Ask <= 0 {
			fmt.Printf("  %s (%s): no quote\n", label, mkt)
			return
		}
		mid := (q.Bid + q.Ask) / 2
		spreadAbs := q.Ask - q.Bid
		spreadPct := spreadAbs / mid * 100
		fmt.Printf("  %s (%s): bid=%.*f ask=%.*f  spread=%.*f (%0.5f%%)  bidQty=%0.4f askQty=%0.4f\n",
			label, mkt,
			10, q.Bid,
			10, q.Ask,
			10, spreadAbs,
			spreadPct,
			q.BidQty, q.AskQty,
		)
	}

	printQuote(fmt.Sprintf("%s/%s", t.A, t.B), t.M1, q1)
	printQuote(fmt.Sprintf("%s/%s", t.B, t.C), t.M2, q2)
	printQuote(fmt.Sprintf("%s/%s", t.A, t.C), t.M3, q3)
	fmt.Println()
}

/* =========================  WS RUNNER (PUBLIC, MULTI-SYMBOL)  ========================= */

func runPublicBookTicker(
	ctx context.Context,
	wg *sync.WaitGroup,
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
			if retry < maxRetry {
				retry *= 2
				if retry > maxRetry {
					retry = maxRetry
				}
			}
			continue
		}
		log.Printf("[WS #%d] connected to %s (symbols: %d)", connID, urlWS, len(symbols))
		retry = baseRetry

		_ = conn.SetReadDeadline(time.Now().Add(90 * time.Second))

		var lastPing time.Time
		conn.SetPongHandler(func(appData string) error {
			rtt := time.Since(lastPing)
			dlog("[WS #%d] PONG via %s in %v", connID, urlWS, rtt)
			return conn.SetReadDeadline(time.Now().Add(90 * time.Second))
		})

		// ping goroutine
		stopPing := make(chan struct{})
		go func(id int, c *websocket.Conn, stop <-chan struct{}) {
			t := time.NewTicker(45 * time.Second)
			defer t.Stop()
			for {
				select {
				case <-t.C:
					lastPing = time.Now()
					if err := c.WriteControl(websocket.PingMessage, []byte("hb"), time.Now().Add(5*time.Second)); err != nil {
						dlog("[WS #%d] ping error: %v", id, err)
						return
					}
				case <-stop:
					return
				}
			}
		}(connID, conn, stopPing)

		// подписка: пачка топиков
		topics := make([]string, 0, len(symbols))
		for _, s := range symbols {
			topics = append(topics,
				fmt.Sprintf("spot@public.aggre.bookTicker.v3.api.pb@%s@%s", interval, s))
		}

		sub := map[string]any{
			"method": "SUBSCRIPTION",
			"params": topics,
			"id":     time.Now().Unix(),
		}
		if err := conn.WriteJSON(sub); err != nil {
			log.Printf("[WS #%d] subscribe send err: %v", connID, err)
			close(stopPing)
			_ = conn.Close()
			time.Sleep(retry)
			continue
		}
		log.Printf("[WS #%d] SUB -> %d topics", connID, len(topics))

		// цикл чтения
		for {
			mt, raw, err := conn.ReadMessage()
			if err != nil {
				log.Printf("[WS #%d] read err: %v (reconnect)", connID, err)
				break
			}

			switch mt {
			case websocket.TextMessage:
				var tmp any
				if err := json.Unmarshal(raw, &tmp); err == nil {
					j, _ := json.Marshal(tmp)
					dlog("[WS #%d TEXT] %s", connID, string(j))
				} else {
					dlog("[WS #%d TEXT RAW] %s", connID, string(raw))
				}

			case websocket.BinaryMessage:
				if sym, bid, ask, ok := parsePBWrapperQuote(raw); ok {
					ev := Event{Symbol: sym, Bid: bid, Ask: ask}
					select {
					case out <- ev:
					case <-ctx.Done():
						close(stopPing)
						_ = conn.Close()
						return
					}
				}
			default:
				// игнор прочих типов
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
	feeRate = cfg.FeePct / 100.0
	minProfitRate = cfg.MinProfitPct / 100.0

	log.Printf("Triangles file: %s", cfg.TrianglesFile)
	log.Printf("Book interval: %s", cfg.BookInterval)
	log.Printf("Fee per leg: %0.4f %% (rate=%f)", cfg.FeePct, feeRate)
	log.Printf("Min profit per cycle: %0.4f %% (rate=%f)", cfg.MinProfitPct, minProfitRate)

	// грузим треугольники
	tris, err := loadTrianglesFromFile(cfg.TrianglesFile)
	if err != nil {
		log.Fatal(err)
	}
	if len(tris) == 0 {
		log.Fatal("нет треугольников в файле ", cfg.TrianglesFile)
	}
	log.Printf("треугольников всего: %d", len(tris))

	// индекс: symbol -> индексы треугольников
	triBySymbol := buildTriangleIndex(tris)
	log.Printf("символов в индексе треугольников: %d", len(triBySymbol))

	// извлекаем список всех символов для подписки
	symbols := extractSymbols(tris)
	log.Printf("символов для подписки всего: %d", len(symbols))

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	events := make(chan Event, 8192)
	var wg sync.WaitGroup

	// чанкуем символы по cfg.MaxSymbolsConn
	maxPerConn := cfg.MaxSymbolsConn
	if maxPerConn <= 0 {
		maxPerConn = 50
	}
	chunks := make([][]string, 0)
	for i := 0; i < len(symbols); i += maxPerConn {
		j := i + maxPerConn
		if j > len(symbols) {
			j = len(symbols)
		}
		chunks = append(chunks, symbols[i:j])
	}
	log.Printf("будем использовать %d WS-подключений", len(chunks))

	// стартуем WS-подключения
	for i, chunk := range chunks {
		wg.Add(1)
		go runPublicBookTicker(ctx, &wg, i, chunk, cfg.BookInterval, events)
	}

	// Консумер: на КАЖДОМ тике считает только треугольники,
	// где участвует обновившийся символ.
	go func(tris []Triangle, ctx context.Context, triBySymbol map[string][]int) {
		last := make(map[string]Quote)

		for {
			select {
			case ev, ok := <-events:
				if !ok {
					return
				}

				// обновили последнюю котировку по символу
				last[ev.Symbol] = Quote{
					Bid:    ev.Bid,
					Ask:    ev.Ask,
					BidQty: ev.BidQty,
					AskQty: ev.AskQty,
				}

				// смотрим только треугольники, где этот symbol участвует
				idxs := triBySymbol[ev.Symbol]
				if len(idxs) == 0 {
					continue
				}

				prof := findProfitableTrianglesForSymbol(tris, idxs, last)
				if len(prof) == 0 {
					continue
				}

				fmt.Printf(
					"\nquotes known: %d symbols, profitable triangles (on %s update): %d\n",
					len(last), ev.Symbol, len(prof),
				)

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
	}(tris, ctx, triBySymbol)

	<-ctx.Done()
	log.Println("shutting down...")

	time.Sleep(300 * time.Millisecond)
	close(events)
	wg.Wait()
	log.Println("bye")
}



uotes known: 360 symbols, profitable triangles (on DOTUSDC update): 2
[ARB] +178290041384.829%  DOT→BTC→USDC→DOT
  DOT/BTC (DOTBTC): bid=0.0000236000 ask=0.0000236900  spread=0.0000000900 (0.38063%)  bidQty=0.0000 askQty=0.0000
  BTC/USDC (DOTUSDC): bid=2.1500000000 ask=2.1600000000  spread=0.0100000000 (0.46404%)  bidQty=0.0000 askQty=0.0000
  DOT/USDC (BTCUSDC): bid=91110.1300000000 ask=91131.3400000000  spread=21.2100000000 (0.02328%)  bidQty=0.0000 askQty=0.0000


