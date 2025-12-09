package main

import (
	"context"
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
	Symbol       string // напр. BTCUSDT
	BookInterval string // "100ms" | "10ms"
	Debug        bool
}

func loadConfig() (Config, error) {
	_ = godotenv.Load(".env")

	cfg := Config{
		Symbol:       os.Getenv("SYMBOL"),
		BookInterval: os.Getenv("BOOK_INTERVAL"),
	}

	if cfg.Symbol == "" {
		cfg.Symbol = "BTCUSDT"
	}
	if cfg.BookInterval == "" {
		cfg.BookInterval = "100ms"
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

/* =========================  WS RUNNER (PUBLIC)  ========================= */

func runPublicBookTicker(ctx context.Context, wg *sync.WaitGroup, symbol, interval string, out chan<- Event) {
	defer wg.Done()

	const (
		baseRetry = 2 * time.Second
		maxRetry  = 30 * time.Second
	)

	urlWS := "wss://wbs-api.mexc.com/ws"
	topic := "spot@public.aggre.bookTicker.v3.api.pb@" + interval + "@" + symbol

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

		// подписка
		sub := map[string]any{
			"method": "SUBSCRIPTION",
			"params": []string{topic},
			"id":     time.Now().Unix(),
		}
		if err := conn.WriteJSON(sub); err != nil {
			log.Printf("[PUB] subscribe send err: %v", err)
			close(stopPing)
			_ = conn.Close()
			time.Sleep(retry)
			continue
		}
		log.Printf("[PUB] SUB → %s", topic)

		// цикл чтения
		for {
			mt, raw, err := conn.ReadMessage()
			if err != nil {
				log.Printf("[PUB] read err: %v (reconnect)", err)
				break
			}

			switch mt {
			case websocket.TextMessage:
				// ACK/ошибки подписки — только в debug
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
				// игнор
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

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	events := make(chan Event, 4096)

	var wg sync.WaitGroup
	wg.Add(1)
	go runPublicBookTicker(ctx, &wg, cfg.Symbol, cfg.BookInterval, events)

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
				for sym, mid := range lastMid {
					fmt.Printf("[MID] %s = %.10f\n", sym, mid)
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

