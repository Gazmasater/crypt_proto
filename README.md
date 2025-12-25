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
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type MarketData struct {
	Symbol string
	Bid    float64
	Ask    float64
}

// NormalizeSymbol для KuCoin (делаем дефис)
func NormalizeSymbol(s string) string {
	s = strings.ToUpper(strings.TrimSpace(s))
	if strings.ContainsAny(s, "-/") {
		return strings.ReplaceAll(s, "/", "-")
	}
	if strings.HasSuffix(s, "USDT") && len(s) > 4 {
		return s[:len(s)-4] + "-USDT"
	}
	return s
}

func main() {
	symbols := []string{"BTCUSDT", "ETHUSDT", "ETHBTC"}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 1) Получаем endpoint и token
	wsURL, err := getKuCoinWS()
	if err != nil {
		log.Fatal(err)
	}

	// 2) Подключаемся
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	log.Println("Connected to KuCoin WS")

	// 3) Подписка на Level2/Ticker
	for _, s := range symbols {
		topic := "level2/ticker:" + NormalizeSymbol(s)
		sub := map[string]interface{}{
			"id":       time.Now().UnixNano(),
			"type":     "subscribe",
			"topic":    topic,
			"response": true,
		}
		if err := conn.WriteJSON(sub); err != nil {
			log.Println("subscribe error:", err)
		} else {
			log.Println("subscribed:", topic)
		}
	}

	// 4) Read loop с выводом первых данных
	go func() {
		for {
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
				body, ok := data["data"].(map[string]interface{})
				if !ok {
					continue
				}
				topic := data["topic"].(string)
				symbol := strings.Split(topic, ":")[1]

				bid, _ := parseFloat(body["bestBid"])
				ask, _ := parseFloat(body["bestAsk"])

				md := MarketData{
					Symbol: symbol,
					Bid:    bid,
					Ask:    ask,
				}
				fmt.Println("MarketData:", md)
			}
		}
	}()

	// держим соединение 20 секунд для теста
	select {
	case <-ctx.Done():
	case <-time.After(20 * time.Second):
	}
}

// =================== helpers ===================

func getKuCoinWS() (string, error) {
	resp, err := http.Post("https://api.kucoin.com/api/v1/bullet-public", "application/json", strings.NewReader("{}"))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var r struct {
		Code string `json:"code"`
		Data struct {
			Token           string `json:"token"`
			InstanceServers []struct {
				Endpoint string `json:"endpoint"`
			} `json:"instanceServers"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return "", err
	}
	if len(r.Data.InstanceServers) == 0 {
		return "", fmt.Errorf("no KuCoin WS endpoints returned")
	}

	endpoint := r.Data.InstanceServers[0].Endpoint
	wsURL := fmt.Sprintf("%s?token=%s&connectId=%d", endpoint, r.Data.Token, time.Now().UnixNano())
	return wsURL, nil
}

func parseFloat(v interface{}) (float64, error) {
	switch t := v.(type) {
	case string:
		return strconv.ParseFloat(t, 64)
	case float64:
		return t, nil
	default:
		return 0, fmt.Errorf("unknown type")
	}
}
