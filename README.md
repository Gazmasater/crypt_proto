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

const KuCoinPublicAPI = "https://api.kucoin.com/api/v1/bullet-public"

type KuCoinWS struct {
	ctx          context.Context
	cancel       context.CancelFunc
	conn         *websocket.Conn
	wsURL        string
	symbols      []string
	pingInterval time.Duration
	noTickTimeout time.Duration
	lastTick     time.Time
}

type wsTokenResponse struct {
	Code string `json:"code"`
	Data struct {
		Token           string `json:"token"`
		InstanceServers []struct {
			Endpoint     string `json:"endpoint"`
			Encrypt      bool   `json:"encrypt"`
			PingInterval int    `json:"pingInterval"`
		} `json:"instanceServers"`
	} `json:"data"`
}

func NewKuCoinWS(symbols []string) *KuCoinWS {
	ctx, cancel := context.WithCancel(context.Background())
	return &KuCoinWS{
		ctx:           ctx,
		cancel:        cancel,
		symbols:       symbols,
		noTickTimeout: 30 * time.Second,
	}
}

// Получаем токен и endpoint
func (k *KuCoinWS) init() error {
	resp, err := http.Post(KuCoinPublicAPI, "application/json", strings.NewReader("{}"))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var r wsTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return err
	}

	if len(r.Data.InstanceServers) == 0 {
		return fmt.Errorf("no WS endpoints returned")
	}

	server := r.Data.InstanceServers[0]
	k.pingInterval = time.Duration(server.PingInterval) * time.Millisecond
	k.wsURL = fmt.Sprintf("%s?token=%s&connectId=%d", server.Endpoint, r.Data.Token, time.Now().UnixNano())
	return nil
}

func (k *KuCoinWS) Start() error {
	if err := k.init(); err != nil {
		return err
	}

	conn, _, err := websocket.DefaultDialer.Dial(k.wsURL, nil)
	if err != nil {
		return err
	}
	k.conn = conn
	log.Println("[KuCoin] Connected to WS")

	// Подписываемся на тикеры
	for _, s := range k.symbols {
		topic := fmt.Sprintf("level2/ticker:%s", s)
		msg := map[string]interface{}{
			"id":       time.Now().UnixNano(),
			"type":     "subscribe",
			"topic":    topic,
			"response": true,
		}
		if err := conn.WriteJSON(msg); err != nil {
			return err
		}
		log.Println("[KuCoin] Subscribed:", s)
	}

	// Ping loop
	go func() {
		ticker := time.NewTicker(k.pingInterval)
		defer ticker.Stop()
		for {
			select {
			case <-k.ctx.Done():
				return
			case <-ticker.C:
				_ = conn.WriteMessage(websocket.PingMessage, nil)
			}
		}
	}()

	// Read loop
	go k.readLoop()

	return nil
}

func (k *KuCoinWS) readLoop() {
	k.lastTick = time.Now()
	for {
		select {
		case <-k.ctx.Done():
			return
		default:
			_, msg, err := k.conn.ReadMessage()
			if err != nil {
				log.Println("[KuCoin] read error:", err)
				return
			}

			var data map[string]interface{}
			if err := json.Unmarshal(msg, &data); err != nil {
				continue
			}

			if data["type"] == "message" {
				k.lastTick = time.Now()
				topic, _ := data["topic"].(string)
				body, _ := data["data"].(map[string]interface{})
				log.Printf("[KuCoin] Tick %s: %+v\n", topic, body)
			}
		}
	}
}

func (k *KuCoinWS) MonitorNoTicks() {
	for {
		time.Sleep(5 * time.Second)
		if time.Since(k.lastTick) > k.noTickTimeout {
			log.Println("[KuCoin] No more ticks received in 30s, exiting")
			k.Stop()
			return
		}
	}
}

func (k *KuCoinWS) Stop() {
	k.cancel()
	if k.conn != nil {
		k.conn.Close()
	}
	log.Println("[KuCoin] WS closed normally")
}

func main() {
	symbols := []string{"BTC-USDT", "ETH-USDT", "XRP-USDT", "DOGE-USDT"}
	k := NewKuCoinWS(symbols)
	if err := k.Start(); err != nil {
		log.Fatal(err)
	}

	k.MonitorNoTicks()
}




