package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"google.golang.org/protobuf/proto"

	pb "crypt_proto/pb" // твои *.pb.go из схем MEXC v3
)

func main() {
	_ = godotenv.Load(".env")

	apiKey := os.Getenv("MEXC_API_KEY")
	secret := os.Getenv("MEXC_SECRET_KEY")
	if apiKey == "" || secret == "" {
		log.Fatal("MEXC_API_KEY / MEXC_SECRET_KEY пусты в .env")
	}

	// --- параметры ---
	symbol := os.Getenv("SYMBOL")
	if symbol == "" {
		symbol = "BTCUSDT"
	}
	interval := os.Getenv("BOOK_INTERVAL")
	if interval == "" {
		interval = "100ms" // варианты обычно: 100ms / 300ms / 500ms / 1s
	}

	// --- вычислим дрейф времени, чтобы не ловить 700003 ---
	drift, err := mexcDriftMs()
	if err != nil {
		log.Printf("⚠️ Не удалось получить server time, продолжаю без дрейфа: %v", err)
		drift = 0
	}

	// --- создаём listenKey (SIGNED с учетом дрейфа и широким recvWindow) ---
	lk, err := createListenKey(apiKey, secret, drift)
	if err != nil {
		log.Fatal(err)
	}
	defer closeListenKey(apiKey, secret, lk, drift)

	// --- коннект к приватному WS ---
	wsURL := "wss://wbs.mexc.com/ws?listenKey=" + lk
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer conn.Close()

	// --- подписка на protobuf bookTicker ---
	channel := "spot@public.bookTicker.v3.api.pb@" + interval + "@" + symbol
	sub := map[string]any{
		"method": "SUBSCRIPTION",
		"params": []string{channel},
		"id":     time.Now().Unix(),
	}
	if err := conn.WriteJSON(sub); err != nil {
		log.Fatal("send sub:", err)
	}
	log.Printf("✅ SUB → %s", channel)

	// --- heartbeat ---
	go func() {
		t := time.NewTicker(45 * time.Second)
		defer t.Stop()
		for range t.C {
			_ = conn.WriteMessage(websocket.PingMessage, []byte("hb"))
		}
	}()

	// --- чтение ---
	for {
		mt, raw, err := conn.ReadMessage()
		if err != nil {
			log.Fatal("read:", err)
		}

		// ACK / ошибки — текст
		if mt == websocket.TextMessage {
			// красиво распечатаем
			var v any
			if json.Unmarshal(raw, &v) == nil {
				b, _ := json.MarshalIndent(v, "", "  ")
				fmt.Printf("ACK:\n%s\n", b)
			} else {
				fmt.Println("TEXT:", string(raw))
			}
			continue
		}
		if mt != websocket.BinaryMessage {
			continue
		}

		// бинарь → protobuf PushDataV3ApiWrapper
		if out, ok := parsePBWrapperToSP(raw); ok {
			// {"s":"SYMBOL","p":"PRICE"} — mid=(bid+ask)/2
			fmt.Println(string(out))
		}
	}
}

/* -------------------------- REST utils -------------------------- */

// смещение серверных часов MEXC относительно локальных (ms)
func mexcDriftMs() (int64, error) {
	resp, err := (&http.Client{Timeout: 5 * time.Second}).Get("https://api.mexc.com/api/v3/time")
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

func hmacSHA256Hex(secret, data string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

// добавляет timestamp(+drift), recvWindow=60000 и signature в query
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

func rest(method, endpoint, apiKey string, v url.Values) (int, []byte, error) {
	if len(v) > 0 {
		endpoint += "?" + v.Encode()
	}
	req, _ := http.NewRequest(method, endpoint, nil)
	req.Header.Set("X-MEXC-APIKEY", apiKey)
	resp, err := (&http.Client{Timeout: 10 * time.Second}).Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, b, nil
}

func createListenKey(apiKey, secret string, driftMs int64) (string, error) {
	st, body, err := rest("POST", "https://api.mexc.com/api/v3/userDataStream", apiKey, signParams(secret, nil, driftMs))
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

func closeListenKey(apiKey, secret, listenKey string, driftMs int64) {
	v := url.Values{}
	v.Set("listenKey", listenKey)
	_, _, _ = rest("DELETE", "https://api.mexc.com/api/v3/userDataStream", apiKey, signParams(secret, v, driftMs))
}

/* ---------------------- Protobuf parsing ------------------------ */

// Поддержим два типа: PublicBookTicker и PublicAggreBookTicker.
// На выход даём {"s","p"} где p = mid(bid, ask).
func parsePBWrapperToSP(raw []byte) ([]byte, bool) {
	var w pb.PushDataV3ApiWrapper
	if err := proto.Unmarshal(raw, &w); err != nil {
		return nil, false
	}

	sym := w.GetSymbol()
	if sym == "" {
		ch := w.GetChannel()
		if i := strings.LastIndex(ch, "@"); i >= 0 && i+1 < len(ch) {
			sym = ch[i+1:]
		}
	}
	if sym == "" {
		return nil, false
	}

	// 1) Обычный bookTicker
	if b1, ok := w.GetBody().(*pb.PushDataV3ApiWrapper_PublicBookTicker); ok && b1.PublicBookTicker != nil {
		bp := b1.PublicBookTicker.GetBidPrice()
		ap := b1.PublicBookTicker.GetAskPrice()
		if bp == "" || ap == "" {
			return nil, false
		}
		bid, err1 := strconv.ParseFloat(bp, 64)
		ask, err2 := strconv.ParseFloat(ap, 64)
		if err1 != nil || err2 != nil || bid <= 0 || ask <= 0 {
			return nil, false
		}
		mid := (bid + ask) / 2
		out := fmt.Sprintf(`{"s":"%s","p":"%.10f"}`, sym, mid)
		return []byte(out), true
	}

	// 2) AggreBookTicker (если вдруг подпишешься на aggre)
	if b2, ok := w.GetBody().(*pb.PushDataV3ApiWrapper_PublicAggreBookTicker); ok && b2.PublicAggreBookTicker != nil {
		bp := b2.PublicAggreBookTicker.GetBidPrice()
		ap := b2.PublicAggreBookTicker.GetAskPrice()
		if bp == "" || ap == "" {
			return nil, false
		}
		bid, err1 := strconv.ParseFloat(bp, 64)
		ask, err2 := strconv.ParseFloat(ap, 64)
		if err1 != nil || err2 != nil || bid <= 0 || ask <= 0 {
			return nil, false
		}
		mid := (bid + ask) / 2
		out := fmt.Sprintf(`{"s":"%s","p":"%.10f"}`, sym, mid)
		return []byte(out), true
	}

	return nil, false
}
