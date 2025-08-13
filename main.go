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
	"google.golang.org/protobuf/proto"

	// если не используешь .env, просто убери этот импорт и вызов godotenv.Load()
	"github.com/joho/godotenv"

	pb "crypt_proto/pb" // твой пакет со сгенерёнными *.pb.go
)

// ============ общие утилиты ============

func hmacSHA256Hex(secret, data string) string {
	m := hmac.New(sha256.New, []byte(secret))
	m.Write([]byte(data))
	return hex.EncodeToString(m.Sum(nil))
}

func httpClient() *http.Client {
	return &http.Client{Timeout: 10 * time.Second}
}

func doSigned(method, endpoint, apiKey, secret string, form url.Values) (int, []byte, error) {
	// добавим timestamp/recvWindow если не переданы
	if form.Get("timestamp") == "" {
		form.Set("timestamp", strconv.FormatInt(time.Now().UnixMilli(), 10))
	}
	if form.Get("recvWindow") == "" {
		form.Set("recvWindow", "5000")
	}
	// signature по всему body (form.Encode())
	sig := hmacSHA256Hex(secret, form.Encode())
	form.Set("signature", sig)

	req, _ := http.NewRequest(method, endpoint, strings.NewReader(form.Encode()))
	req.Header.Set("X-MEXC-APIKEY", apiKey)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := httpClient().Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, body, nil
}

// ============ listenKey helpers (теперь все запросы подписаны) ============

func createListenKey(apiKey, secret string) (string, error) {
	status, body, err := doSigned("POST", "https://api.mexc.com/api/v3/userDataStream", apiKey, secret, url.Values{})
	if err != nil {
		return "", fmt.Errorf("listenKey request: %w", err)
	}
	if status != http.StatusOK {
		return "", fmt.Errorf("listenKey http %d: %s", status, strings.TrimSpace(string(body)))
	}
	var out struct {
		ListenKey string `json:"listenKey"`
	}
	if err := json.Unmarshal(body, &out); err != nil {
		return "", fmt.Errorf("listenKey decode: %v; body=%s", err, body)
	}
	if out.ListenKey == "" {
		return "", fmt.Errorf("empty listenKey; body=%s", body)
	}
	return out.ListenKey, nil
}

func keepAliveListenKey(apiKey, secret, listenKey string, stop <-chan struct{}) {
	t := time.NewTicker(30 * time.Minute) // продлеваем до истечения 60м
	defer t.Stop()
	for {
		select {
		case <-t.C:
			form := url.Values{}
			form.Set("listenKey", listenKey)
			if status, body, err := doSigned("PUT",
				"https://api.mexc.com/api/v3/userDataStream", apiKey, secret, form); err != nil || status != http.StatusOK {
				log.Printf("keepAlive listenKey error: status=%d err=%v body=%s", status, err, strings.TrimSpace(string(body)))
			}
		case <-stop:
			return
		}
	}
}

func closeListenKey(apiKey, secret, listenKey string) {
	form := url.Values{}
	form.Set("listenKey", listenKey)
	if status, body, err := doSigned("DELETE",
		"https://api.mexc.com/api/v3/userDataStream", apiKey, secret, form); err != nil || status != http.StatusOK {
		log.Printf("close listenKey error: status=%d err=%v body=%s", status, err, strings.TrimSpace(string(body)))
	}
}

// ============ main ============

func main() {
	// подхватим .env (необязательно)
	_ = godotenv.Load(".env")

	const baseWS = "wss://wbs-api.mexc.com/ws"

	apiKey := os.Getenv("MEXC_API_KEY")
	secret := os.Getenv("MEXC_SECRET_KEY")
	if apiKey == "" || secret == "" {
		log.Fatal("MEXC_API_KEY / MEXC_SECRET_KEY пусты. Проверь .env или экспорт.")
	}

	// 1) создаём listenKey (подписанный POST)
	listenKey, err := createListenKey(apiKey, secret)
	if err != nil {
		log.Fatal("listenKey:", err)
	}
	defer closeListenKey(apiKey, secret, listenKey)

	// 2) продлеваем listenKey в фоне (подписанный PUT)
	stopKA := make(chan struct{})
	go keepAliveListenKey(apiKey, secret, listenKey, stopKA)
	defer close(stopKA)

	// 3) подключаемся к приватному ws и подписываемся на PUBLIC .pb каналы
	wsURL := baseWS + "?listenKey=" + listenKey
	c, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	sub := map[string]any{
		"method": "SUBSCRIPTION",
		"params": []string{
			"spot@public.aggre.bookTicker.v3.api.pb@100ms@BTCUSDT",
			"spot@public.aggre.bookTicker.v3.api.pb@100ms@ETHUSDT",
			"spot@public.aggre.bookTicker.v3.api.pb@100ms@ETHBTC",
		},
	}
	if err := c.WriteJSON(sub); err != nil {
		log.Fatal("send sub:", err)
	}

	// пингуем линк
	go func() {
		t := time.NewTicker(45 * time.Second)
		defer t.Stop()
		for range t.C {
			_ = c.WriteMessage(websocket.PingMessage, []byte("hb"))
		}
	}()

	// читаем и печатаем ТЕМ ЖЕ ФОРМАТОМ
	for {
		mt, raw, err := c.ReadMessage()
		if err != nil {
			log.Fatal("read:", err)
		}

		// ACK/ошибки — текст/JSON
		if mt == websocket.TextMessage {
			var v any
			if json.Unmarshal(raw, &v) == nil {
				b, _ := json.MarshalIndent(v, "", "  ")
				fmt.Printf("ACK:\n%s\n", b)
			} else {
				fmt.Printf("TEXT:\n%s\n", string(raw))
			}
			continue
		}
		if mt != websocket.BinaryMessage {
			continue
		}

		// 1) Декодируем ОБЁРТКУ
		var w pb.PushDataV3ApiWrapper
		if err := proto.Unmarshal(raw, &w); err != nil {
			// бинарь не нашей схемы — пропускаем
			continue
		}

		// 2) symbol/ts
		symbol := w.GetSymbol()
		if symbol == "" {
			ch := w.GetChannel()
			if ch != "" {
				parts := strings.Split(ch, "@")
				symbol = parts[len(parts)-1]
			}
		}
		ts := time.Now()
		if t := w.GetSendTime(); t > 0 {
			ts = time.UnixMilli(t)
		}

		// 3) интересует PublicAggreBookTicker — вывод НЕ МЕНЯЕМ
		switch body := w.GetBody().(type) {
		case *pb.PushDataV3ApiWrapper_PublicAggreBookTicker:
			bt := body.PublicAggreBookTicker

			bid, _ := strconv.ParseFloat(bt.GetBidPrice(), 64)
			ask, _ := strconv.ParseFloat(bt.GetAskPrice(), 64)
			bq, _ := strconv.ParseFloat(bt.GetBidQuantity(), 64)
			aq, _ := strconv.ParseFloat(bt.GetAskQuantity(), 64)

			fmt.Printf("%s  bid=%.8f (%.6f)  ask=%.8f (%.6f)  ts=%s\n",
				symbol, bid, bq, ask, aq, ts.Format(time.RFC3339Nano))

		default:
			// игнорим прочее
		}
	}
}
