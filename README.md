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





package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
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
	BookInterval string
	Debug        bool
	SymbolsFile  string
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

// Загружаем треугольники из CSV и одновременно собираем уникальные символы
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
			if err.Error() == "EOF" {
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

/* =========================  EVENTS  ========================= */

type Event struct {
	Symbol string
	Mid    float64
}

/* =========================  PROTO DECODER  ========================= */

// Пул для protobuf-структуры, чтобы не аллоцировать на каждый тик
var wrapperPool = sync.Pool{
	New: func() any { return new(pb.PushDataV3ApiWrapper) },
}

// Возвращаем символ и mid=(bid+ask)/2, если это (aggre.)bookTicker
func parsePBWrapperMid(raw []byte) (sym string, mid float64, ok bool) {
	w, _ := wrapperPool.Get().(*pb.PushDataV3ApiWrapper)
	defer func() {
		*w = pb.PushDataV3ApiWrapper{} // очищаем
		wrapperPool.Put(w)
	}()

	if err := proto.Unmarshal(raw, w); err != nil {
		return "", 0, false
	}

	sym = w.GetSymbol()
	if sym == "" {
		ch := w.GetChannel()
		if i := strings.LastIndex(ch, "@"); i >= 0 && i+1 < len(ch) {
			sym = ch[i+1:]
		}
	}
	if sym == "" {
		return "", 0, false
	}

	// PublicBookTicker
	if b1, ok1 := w.GetBody().(*pb.PushDataV3ApiWrapper_PublicBookTicker); ok1 && b1.PublicBookTicker != nil {
		bp := b1.PublicBookTicker.GetBidPrice()
		ap := b1.PublicBookTicker.GetAskPrice()
		if bp == "" || ap == "" {
			return "", 0, false
		}
		bid, err1 := strconv.ParseFloat(bp, 64)
		ask, err2 := strconv.ParseFloat(ap, 64)
		if err1 != nil || err2 != nil || bid <= 0 || ask <= 0 {
			return "", 0, false
		}
		return sym, (bid + ask) / 2, true
	}

	// PublicAggreBookTicker
	if b2, ok2 := w.GetBody().(*pb.PushDataV3ApiWrapper_PublicAggreBookTicker); ok2 && b2.PublicAggreBookTicker != nil {
		bp := b2.PublicAggreBookTicker.GetBidPrice()
		ap := b2.PublicAggreBookTicker.GetAskPrice()
		if bp == "" || ap == "" {
			return "", 0, false
		}
		bid, err1 := strconv.ParseFloat(bp, 64)
		ask, err2 := strconv.ParseFloat(ap, 64)
		if err1 != nil || err2 != nil || bid <= 0 || ask <= 0 {
			return "", 0, false
		}
		return sym, (bid + ask) / 2, true
	}

	return "", 0, false
}

/* =========================  WS RUNNER (PUBLIC)  ========================= */

func runPublicBookTicker(ctx context.Context, wg *sync.WaitGroup, symbols []string, interval string, out chan<- Event) {
	defer wg.Done()

	if len(symbols) == 0 {
		log.Println("[PUB] нет символов для подписки")
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
	log.Printf("[PUB] всего символов для подписки в этом conn: %d", len(params))

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
				// ACK/ошибки подписки и т.п. — только в debug
				var tmp any
				if err := json.Unmarshal(raw, &tmp); err == nil {
					j, _ := json.Marshal(tmp)
					dlog("[PUB TEXT] %s", string(j))
				} else {
					dlog("[PUB TEXT RAW] %s", string(raw))
				}
			case websocket.BinaryMessage:
				if sym, mid, ok := parsePBWrapperMid(raw); ok {
					ev := Event{Symbol: sym, Mid: mid}
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

const (
	tradeFee  = 0.001  // 0.1% комиссия за сделку
	minProfit = 0.0005 // 0.05% минимальная прибыль по кругу
)

type TriangleArb struct {
	Triangle Triangle
	StartCur string
	Path     [3]string
	Profit   float64 // множитель, например 1.001 = +0.1%
}

// считаем прибыль по одному треугольнику, если есть все котировки
func evalTriangle(t Triangle, mids map[string]float64) *TriangleArb {
	// набираем уникальные валюты
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

	// строим маленький граф курсов
	type key struct{ From, To string }
	rates := make(map[key]float64)

	ff := 1.0 - tradeFee

	for _, m := range t.M {
		price, ok := mids[m.Symbol]
		if !ok || price <= 0 {
			return nil
		}
		// BASE -> QUOTE
		rates[key{m.Base, m.Quote}] = price * ff
		// QUOTE -> BASE
		rates[key{m.Quote, m.Base}] = (1.0 / price) * ff
	}

	// выберем просто первый порядок curs[0]->curs[1]->curs[2]->curs[0]
	a, b, c := curs[0], curs[1], curs[2]

	get := func(from, to string) (float64, bool) {
		r, ok := rates[key{from, to}]
		return r, ok
	}

	r1, ok1 := get(a, b)
	r2, ok2 := get(b, c)
	r3, ok3 := get(c, a)

	if !(ok1 && ok2 && ok3) {
		// попробуем другой порядок: a->c->b->a
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
		return &TriangleArb{
			Triangle: t,
			StartCur: a,
			Path:     [3]string{a + "→" + c, c + "→" + b, b + "→" + a},
			Profit:   profit,
		}
	}

	profit := r1 * r2 * r3
	if profit <= 1.0+minProfit {
		return nil
	}
	return &TriangleArb{
		Triangle: t,
		StartCur: a,
		Path:     [3]string{a + "→" + b, b + "→" + c, c + "→" + a},
		Profit:   profit,
	}
}

// находим все прибыльные треугольники за текущий тик
func findProfitableTriangles(tris []Triangle, mids map[string]float64) []TriangleArb {
	var res []TriangleArb
	for _, t := range tris {
		if arb := evalTriangle(t, mids); arb != nil {
			res = append(res, *arb)
		}
	}
	return res
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

	// грузим треугольники и список символов
	tris, symbols, err := loadTriangles(cfg.SymbolsFile)
	if err != nil {
		log.Fatalf("load triangles: %v", err)
	}
	if len(symbols) == 0 {
		log.Fatal("нет символов в файле ", cfg.SymbolsFile)
	}
	log.Printf("треугольников: %d, символов для подписки всего: %d", len(tris), len(symbols))

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

	// Консумер: хранит последний mid по символу, раз в секунду ищет прибыльные треугольники
	go func(tris []Triangle) {
		lastMid := make(map[string]float64)
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case ev, ok := <-events:
				if !ok {
					return
				}
				lastMid[ev.Symbol] = ev.Mid
			case <-ticker.C:
				prof := findProfitableTriangles(tris, lastMid)
				fmt.Printf("known mids: %d symbols, profitable triangles: %d\n", len(lastMid), len(prof))

				// выводим максимум 5 лучших за тик
				for i, a := range prof {
					if i >= 5 {
						break
					}
					perc := (a.Profit - 1.0) * 100.0
					m1, m2, m3 := a.Triangle.M[0], a.Triangle.M[1], a.Triangle.M[2]
					fmt.Printf("[ARB] %+0.3f%%  %s  markets: %s(%s/%s), %s(%s/%s), %s(%s/%s)\n",
						perc,
						strings.Join(a.Path[:], "  "),
						m1.Symbol, m1.Base, m1.Quote,
						m2.Symbol, m2.Base, m2.Quote,
						m3.Symbol, m3.Base, m3.Quote,
					)
				}
			}
		}
	}(tris)

	<-ctx.Done()

	time.Sleep(300 * time.Millisecond)
	close(events)
	wg.Wait()
	log.Println("bye")
}





