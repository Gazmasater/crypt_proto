apikey = "4333ed4b-cd83-49f5-97d1-c399e2349748"
secretkey = "E3848531135EDB4CCFDA0F1BC14CD274"
IP = ""
Название API-ключа = "Arb"
Доступы = "Чтение"



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





package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type BulletResponse struct {
	Code string `json:"code"`
	Data struct {
		Token           string `json:"token"`
		InstanceServers []struct {
			Endpoint     string `json:"endpoint"`
			Encrypt      bool   `json:"encrypt"`
			Protocol     string `json:"protocol"`
			PingInterval int    `json:"pingInterval"`
		} `json:"instanceServers"`
	} `json:"data"`
}

type MarketData struct {
	Symbol string
	Bid    float64
	Ask    float64
}

func getKuCoinToken() (string, string, int, error) {
	reqBody := []byte("{}")
	resp, err := http.Post("https://api.kucoin.com/api/v1/bullet-public", "application/json", bytes.NewReader(reqBody))
	if err != nil {
		return "", "", 0, err
	}
	defer resp.Body.Close()

	var r BulletResponse
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return "", "", 0, err
	}

	if len(r.Data.InstanceServers) == 0 {
		return "", "", 0, fmt.Errorf("no WS endpoints returned")
	}

	endpoint := r.Data.InstanceServers[0].Endpoint
	pingInterval := r.Data.InstanceServers[0].PingInterval
	return endpoint, r.Data.Token, pingInterval, nil
}

func normalizeSymbol(s string) string {
	return strings.ToUpper(strings.ReplaceAll(s, "-", "/"))
}

func subscribeSymbols(conn *websocket.Conn, symbols []string) error {
	for _, s := range symbols {
		topic := "level2/ticker:" + s
		sub := map[string]interface{}{
			"id":       time.Now().UnixNano(),
			"type":     "subscribe",
			"topic":    topic,
			"response": true,
		}
		if err := conn.WriteJSON(sub); err != nil {
			return err
		}
		log.Println("subscribed:", topic)
	}
	return nil
}

func readLoop(ctx context.Context, conn *websocket.Conn) {
	defer conn.Close()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println("read error:", err)
				return
			}

			var data map[string]interface{}
			if err := json.Unmarshal(msg, &data); err != nil {
				continue
			}

			if data["type"] == "message" {
				topic := data["topic"].(string)
				symbol := strings.Split(topic, ":")[1]
				body := data["data"].(map[string]interface{})

				bidStr, ok1 := body["bestBid"].(string)
				askStr, ok2 := body["bestAsk"].(string)
				if !ok1 || !ok2 {
					continue
				}

				bid, _ := strconv.ParseFloat(bidStr, 64)
				ask, _ := strconv.ParseFloat(askStr, 64)

				log.Printf("%s -> Bid: %f, Ask: %f\n", normalizeSymbol(symbol), bid, ask)
			}
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	symbols := []string{"BTC-USDT", "ETH-USDT"}

	endpoint, token, pingInterval, err := getKuCoinToken()
	if err != nil {
		log.Fatal(err)
	}

	wsURL := fmt.Sprintf("%s?token=%s&connectId=%d", endpoint, token, time.Now().UnixNano())
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to KuCoin WS")

	// Подписка
	if err := subscribeSymbols(conn, symbols); err != nil {
		log.Fatal(err)
	}

	// Ping loop
	go func() {
		ticker := time.NewTicker(time.Duration(pingInterval) * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				_ = conn.WriteMessage(websocket.PingMessage, nil)
			}
		}
	}()

	// Read loop
	readLoop(ctx, conn)
}

