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
	BookInterval string // "100ms" | "10ms"
	Debug        bool
	SymbolsFile  string // файл с треугольниками (или парами), по умолч. triangles_markets.csv
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

/* =========================  ЗАГРУЗКА СИМВОЛОВ ИЗ CSV  ========================= */

// Ожидаем файл вида:
// base1,quote1,base2,quote2,base3,quote3
// ETC,BTC,BTC,USDC,ETC,USDC
// ...
// Берём все base/quote, собираем уникальные символы BASE+QUOTE.
func loadSymbolsFromTriangles(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", path, err)
	}
	defer f.Close()

	r := csv.NewReader(f)

	// читаем заголовок
	if _, err := r.Read(); err != nil {
		return nil, fmt.Errorf("read header: %w", err)
	}

	seen := make(map[string]struct{})

	for {
		rec, err := r.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, fmt.Errorf("read record: %w", err)
		}
		if len(rec) < 6 {
			continue
		}

		// три рынка (base, quote)
		bases := []string{rec[0], rec[2], rec[4]}
		quotes := []string{rec[1], rec[3], rec[5]}

		for i := 0; i < 3; i++ {
			base := strings.TrimSpace(bases[i])
			quote := strings.TrimSpace(quotes[i])
			if base == "" || quote == "" {
				continue
			}
			sym := base + quote // формат MEXC: BTC + USDT = BTCUSDT
			seen[sym] = struct{}{}
		}
	}

	symbols := make([]string, 0, len(seen))
	for s := range seen {
		symbols = append(symbols, s)
	}
	return symbols, nil
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

/* =========================  MAIN  ========================= */

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	cfg, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}
	debug = cfg.Debug

	// грузим символы из triangles_markets.csv (или другого файла)
	symbols, err := loadSymbolsFromTriangles(cfg.SymbolsFile)
	if err != nil {
		log.Fatalf("load symbols: %v", err)
	}
	if len(symbols) == 0 {
		log.Fatal("нет символов в файле ", cfg.SymbolsFile)
	}
	log.Printf("символов для подписки всего: %d", len(symbols))

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

	// Консумер: хранит последний mid по символу, печатает агрегированно раз в секунду
	go func() {
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
				// тут потом будет движок треугольников
				fmt.Printf("known mids: %d symbols\n", len(lastMid))
				// можно вывести несколько для контроля
				i := 0
				for sym, mid := range lastMid {
					if i >= 5 {
						break
					}
					fmt.Printf("[MID] %s = %.10f\n", sym, mid)
					i++
				}
			}
		}
	}()

	<-ctx.Done()

	time.Sleep(300 * time.Millisecond)
	close(events)
	wg.Wait()
	log.Println("bye")
}







