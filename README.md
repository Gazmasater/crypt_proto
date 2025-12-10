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
	"bufio"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

	pb "crypt_proto/pb"

	_ "net/http/pprof"
)

/* =========================  CONFIG  ========================= */

type Config struct {
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

	feePct := loadEnvFloat("FEE_PER_LEG", 0.1)   // проценты
	minPct := loadEnvFloat("MIN_PROFIT_PCT", 0.3) // проценты

	debug := strings.ToLower(os.Getenv("DEBUG")) == "true"

	cfg := Config{
		TrianglesFile: tf,
		BookInterval:  bi,
		FeePerLeg:     feePct / 100.0,
		MinProfit:     minPct / 100.0,
		Debug:         debug,
	}

	log.Printf("Triangles file: %s", tf)
	log.Printf("Book interval: %s", bi)
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

/* =========================  TRIANGLES  ========================= */

type Leg struct {
	From   string // валюта "от"
	To     string // валюта "к"
	Symbol string // символ рынка, например BDXUSDT
	Dir    int8   // +1: From->To = base->quote (продажа базовой по bid); -1: From->To = quote->base (покупка базовой по ask)
}

type Triangle struct {
	Legs [3]Leg
	Name string // удобное имя: A->B->C->A
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
	// собираем 3 разные валюты
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

	// перебираем перестановки валют и пар, ищем замкнутый цикл
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

// читаем файл треугольников, строим Triangle, множество символов и индекс "символ -> какие треугольники его используют"
func loadTriangles(path string) ([]Triangle, []string, map[string][]int, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, nil, err
	}
	defer f.Close()

	r := csv.NewReader(bufio.NewReader(f))
	r.TrimLeadingSpace = true
	// на всякий случай разрешим и запятую, и пробелы через кастомный сплит:
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
		// если кто-то руками правил файл — подчищаем
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
			// не удалось собрать реальный цикл валют – пропускаем
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

	// индекс по символам
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

/* =========================  PROTO DECODER  ========================= */

// Пул для wrapper'а
var wrapperPool = sync.Pool{
	New: func() any { return new(pb.PushDataV3ApiWrapper) },
}

// parsePBQuote: бинарное сообщение -> символ + Quote
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

	// 1) PublicBookTicker
	if b1, ok := w.GetBody().(*pb.PushDataV3ApiWrapper_PublicBookTicker); ok && b1.PublicBookTicker != nil {
		t := b1.PublicBookTicker

		bp := t.GetBidPrice()
		ap := t.GetAskPrice()
		if bp == "" || ap == "" {
			return "", Quote{}, false
		}
		bid, err1 := strconv.ParseFloat(bp, 64)
		ask, err2 := strconv.ParseFloat(ap, 64)
		if err1 != nil || err2 != nil || bid <= 0 || ask <= 0 {
			return "", Quote{}, false
		}
		return sym, Quote{
			Bid:    bid,
			Ask:    ask,
			BidQty: 0,
			AskQty: 0,
		}, true
	}

	// 2) PublicAggreBookTicker (у него как раз есть количество)
	if b2, ok := w.GetBody().(*pb.PushDataV3ApiWrapper_PublicAggreBookTicker); ok && b2.PublicAggreBookTicker != nil {
		t := b2.PublicAggreBookTicker

		bp := t.GetBidPrice()
		ap := t.GetAskPrice()
		bq := t.GetBidQuantity()
		aq := t.GetAskQuantity()

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

	return "", Quote{}, false
}

/* =========================  WS SUBSCRIBER  ========================= */

func runPublicBookTickerWS(ctx context.Context, wg *sync.WaitGroup, connID int, symbols []string, interval string, out chan<- Event) {
	defer wg.Done()

	const (
		baseRetry = 2 * time.Second
		maxRetry  = 30 * time.Second
	)

	urlWS := "wss://wbs-api.mexc.com/ws"

	// готовим список топиков
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
			dlog("[WS #%d] Pong через %v", connID, rtt)
			return conn.SetReadDeadline(time.Now().Add(90 * time.Second))
		})

		// keepalive
		stopPing := make(chan struct{})
		go func() {
			t := time.NewTicker(45 * time.Second)
			defer t.Stop()
			for {
				select {
				case <-t.C:
					lastPing = time.Now()
					if err := conn.WriteControl(websocket.PingMessage, []byte("hb"), time.Now().Add(5*time.Second)); err != nil {
						dlog("[WS #%d] ping error: %v", connID, err)
						return
					}
				case <-stopPing:
					return
				}
			}
		}()

		// подписка
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
				// ACK/ошибки подписки только в debug
				var tmp any
				if err := json.Unmarshal(raw, &tmp); err == nil {
					j, _ := json.Marshal(tmp)
					dlog("[WS #%d TEXT] %s", connID, string(j))
				} else {
					dlog("[WS #%d TEXT RAW] %s", connID, string(raw))
				}
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
					close(stopPing)
					_ = conn.Close()
					return
				}
			default:
				// игнор
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

/* =========================  TRIANGLE EVAL  ========================= */

// считаем один треугольник, если для всех 3-х символов есть котировки
// возвращаем (profit, ok)
func evalTriangle(t Triangle, quotes map[string]Quote, fee float64) (float64, bool) {
	amt := 1.0

	for _, leg := range t.Legs {
		q, ok := quotes[leg.Symbol]
		if !ok || q.Bid <= 0 || q.Ask <= 0 {
			return 0, false
		}

		// From->To
		if leg.Dir > 0 {
			// base -> quote, продаём базу по bid
			amt = amt * q.Bid
		} else {
			// quote -> base, покупаем базу за quote по ask
			amt = amt / q.Ask
		}

		// комиссия после каждой сделки
		amt = amt * (1 - fee)
		if amt <= 0 {
			return 0, false
		}
	}

	return amt - 1.0, true
}

func printTriangle(t Triangle, profit float64, quotes map[string]Quote) {
	fmt.Printf("\n[ARB] %+0.3f%%  %s\n", profit*100, t.Name)
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
		fmt.Printf("  %s (%s): bid=%.10f ask=%.10f  spread=%.10f (%.5f%%)  bidQty=%.4f askQty=%.4f\n",
			leg.Symbol, side,
			q.Bid, q.Ask,
			spreadAbs, spreadPct,
			q.BidQty, q.AskQty,
		)
	}
}

/* =========================  MAIN  ========================= */

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

	// чанкуем по 50 символов на одно WS
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
		go runPublicBookTickerWS(ctx, &wg, idx, chunk, cfg.BookInterval, events)
	}

	// консумер: на каждом тике обновляем котировку и считаем только те треугольники,
	// где этот символ участвует; выводим только прибыльные
	go func() {
		quotes := make(map[string]Quote)

		for {
			select {
			case ev, ok := <-events:
				if !ok {
					return
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

				for _, id := range trIDs {
					tr := triangles[id]
					prof, ok := evalTriangle(tr, quotes, cfg.FeePerLeg)
					if !ok {
						continue
					}
					if prof >= cfg.MinProfit {
						printTriangle(tr, prof, quotes)
					}
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	<-ctx.Done()
	log.Println("shutting down...")

	time.Sleep(200 * time.Millisecond)
	close(events)
	wg.Wait()
	log.Println("bye")
}



