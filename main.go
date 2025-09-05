package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
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
	APIKey       string
	APISecret    string
	Symbol       string // напр. BTCUSDT
	BookInterval string // "100ms" | "10ms"
}

func loadConfig() (Config, error) {
	_ = godotenv.Load(".env")
	cfg := Config{
		APIKey:       os.Getenv("MEXC_API_KEY"),
		APISecret:    os.Getenv("MEXC_SECRET_KEY"),
		Symbol:       os.Getenv("SYMBOL"),
		BookInterval: os.Getenv("BOOK_INTERVAL"),
	}
	if cfg.Symbol == "" {
		cfg.Symbol = "BTCUSDT"
	}
	if cfg.BookInterval == "" {
		cfg.BookInterval = "100ms"
	}
	if cfg.APIKey == "" || cfg.APISecret == "" {
		return cfg, errors.New("MEXC_API_KEY / MEXC_SECRET_KEY пусты")
	}
	return cfg, nil
}

/* =========================  REST UTILS  ========================= */

func hmacSHA256Hex(secret, data string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func mexcDriftMs(ctx context.Context) (int64, error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.mexc.com/api/v3/time", nil)
	resp, err := (&http.Client{Timeout: 5 * time.Second}).Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	var v struct {
		ServerTime int64 `json:"serverTime"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return 0, err
	}
	local := time.Now().UnixMilli()
	return v.ServerTime - local, nil
}

func signParams(secret string, v url.Values, driftMs int64) url.Values {
	if v == nil {
		v = url.Values{}
	}
	v.Set("timestamp", strconv.FormatInt(time.Now().UnixMilli()+driftMs, 10))
	if v.Get("recvWindow") == "" {
		v.Set("recvWindow", "60000")
	}
	v.Set("signature", hmacSHA256Hex(secret, v.Encode()))
	return v
}

func rest(ctx context.Context, method, endpoint, apiKey string, v url.Values) (int, []byte, error) {
	if len(v) > 0 {
		endpoint += "?" + v.Encode()
	}
	req, _ := http.NewRequestWithContext(ctx, method, endpoint, nil)
	req.Header.Set("X-MEXC-APIKEY", apiKey)
	resp, err := (&http.Client{Timeout: 10 * time.Second}).Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, b, nil
}

func createListenKey(ctx context.Context, apiKey, secret string, driftMs int64) (string, error) {
	st, body, err := rest(ctx, http.MethodPost, "https://api.mexc.com/api/v3/userDataStream", apiKey, signParams(secret, nil, driftMs))
	if err != nil {
		return "", fmt.Errorf("listenKey req: %w", err)
	}
	if st != http.StatusOK {
		return "", fmt.Errorf("listenKey %d: %s", st, strings.TrimSpace(string(body)))
	}
	var out struct {
		ListenKey string `json:"listenKey"`
	}
	if err := json.Unmarshal(body, &out); err != nil || out.ListenKey == "" {
		return "", fmt.Errorf("listenKey decode: %v; body=%s", err, body)
	}
	return out.ListenKey, nil
}

func closeListenKey(ctx context.Context, apiKey, secret, listenKey string, driftMs int64) {
	v := url.Values{}
	v.Set("listenKey", listenKey)
	_, _, _ = rest(ctx, http.MethodDelete, "https://api.mexc.com/api/v3/userDataStream", apiKey, signParams(secret, v, driftMs))
}

/* =========================  PROTO DECODER  ========================= */

// Возвращаем символ и mid=(bid+ask)/2 если это (aggre.)bookTicker
func parsePBWrapperMid(raw []byte) (sym string, mid float64, ok bool) {
	var w pb.PushDataV3ApiWrapper
	if err := proto.Unmarshal(raw, &w); err != nil {
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

/* =========================  WS RUNNERS  ========================= */

func runPublicBookTicker(ctx context.Context, wg *sync.WaitGroup, symbol, interval string, out chan<- string) {
	defer wg.Done()

	const (
		baseRetry = 2 * time.Second
		maxRetry  = 30 * time.Second
	)

	urlWS := "wss://wbs-api.mexc.com/ws" // актуальный публичный WS
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
		retry = baseRetry // сбрасываем бэкофф после успешного коннекта

		// дедлайн для чтения + обработчик PONG с печатью RTT
		_ = conn.SetReadDeadline(time.Now().Add(90 * time.Second))

		var lastPing time.Time
		conn.SetPongHandler(func(appData string) error {
			rtt := time.Since(lastPing)
			log.Printf("[PING] Pong от %s через %v", urlWS, rtt)
			return conn.SetReadDeadline(time.Now().Add(90 * time.Second))
		})

		// keepalive (PING с логом времени отправки)
		stopPing := make(chan struct{})
		go func() {
			t := time.NewTicker(45 * time.Second)
			defer t.Stop()
			for {
				select {
				case <-t.C:
					lastPing = time.Now()
					if err := conn.WriteControl(websocket.PingMessage, []byte("hb"), time.Now().Add(5*time.Second)); err != nil {
						log.Printf("⚠️ [PING] send error: %v", err)
						return
					}
					log.Printf("[PING] Sent at %s", lastPing.Format("15:04:05.000"))
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
				var v any
				if json.Unmarshal(raw, &v) == nil {
					b, _ := json.MarshalIndent(v, "", "  ")
					log.Printf("[PUB ACK]\n%s", b)
				} else {
					log.Printf("[PUB TEXT] %s", string(raw))
				}
			case websocket.BinaryMessage:
				if sym, mid, ok := parsePBWrapperMid(raw); ok {
					out <- fmt.Sprintf(`{"type":"bookTicker","s":"%s","mid":%.10f}`, sym, mid)
				}
			default:
				// игнорируем другие типы
			}
		}

		// cleanup и реконнект
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

func runPrivateUserStream(ctx context.Context, wg *sync.WaitGroup, apiKey, secret string, out chan<- string) {
	defer wg.Done()

	retry := time.Second * 2

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		// drift + listenKey
		drift, err := mexcDriftMs(ctx)
		if err != nil {
			log.Printf("[PRIV] drift err: %v (continue w/0)", err)
			drift = 0
		}
		lk, err := createListenKey(ctx, apiKey, secret, drift)
		if err != nil {
			log.Printf("[PRIV] create listenKey err: %v", err)
			time.Sleep(retry)
			continue
		}
		log.Printf("[PRIV] listenKey OK")

		urlWS := "wss://wbs-api.mexc.com/ws?listenKey=" + lk

		conn, _, err := websocket.DefaultDialer.Dial(urlWS, nil)
		if err != nil {
			log.Printf("[PRIV] dial err: %v", err)
			time.Sleep(retry)
			continue
		}
		log.Printf("[PRIV] connected")

		// read deadline + PONG
		_ = conn.SetReadDeadline(time.Now().Add(90 * time.Second))
		conn.SetPongHandler(func(string) error {
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
					_ = conn.WriteControl(websocket.PingMessage, []byte("hb"), time.Now().Add(5*time.Second))
				case <-stopPing:
					return
				}
			}
		}()

		// read loop
		readErr := false
		for {
			mt, raw, err := conn.ReadMessage()
			if err != nil {
				log.Printf("[PRIV] read err: %v (reconnect)", err)
				readErr = true
				break
			}
			if mt == websocket.TextMessage {
				// событие аккаунта / ACK
				out <- fmt.Sprintf(`{"type":"user","raw":%s}`, safeText(raw))
				continue
			}
			if mt == websocket.BinaryMessage {
				// у MEXC user stream тоже бывают pb — просто выведем длину
				out <- fmt.Sprintf(`{"type":"user.pb","bytes":%d}`, len(raw))
				continue
			}
		}

		// cleanup
		close(stopPing)
		_ = conn.Close()
		// закрыть listenKey
		closeListenKey(context.Background(), apiKey, secret, lk, drift)

		// задержка и повторная попытка
		if readErr {
			time.Sleep(retry)
		}
	}
}

func safeText(b []byte) string {
	// если это JSON — оставим как есть, иначе экранируем строкой
	var v any
	if json.Unmarshal(b, &v) == nil {
		return string(b)
	}
	j, _ := json.Marshal(string(b))
	return string(j)
}

/* =========================  MAIN  ========================= */

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	cfg, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// канал событий (книга/юзер)
	events := make(chan string, 1024)

	var wg sync.WaitGroup
	wg.Add(2)
	go runPublicBookTicker(ctx, &wg, cfg.Symbol, cfg.BookInterval, events)
	go runPrivateUserStream(ctx, &wg, cfg.APIKey, cfg.APISecret, events)

	// консумер событий
	go func() {
		for ev := range events {
			fmt.Println(ev)
		}
	}()

	<-ctx.Done()
	// даём горутинам корректно завершиться
	time.Sleep(300 * time.Millisecond)
	close(events)
	wg.Wait()
	log.Println("bye")
}
