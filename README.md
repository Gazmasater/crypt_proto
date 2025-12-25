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
	"time"

	"github.com/gorilla/websocket"
)

const bulletURL = "https://api.kucoin.com/api/v1/bullet-public"

type BulletResp struct {
	Data struct {
		Token string `json:"token"`
		InstanceServers []struct {
			Endpoint     string `json:"endpoint"`
			PingInterval int    `json:"pingInterval"`
		} `json:"instanceServers"`
	} `json:"data"`
}

func getWS() (string, time.Duration, error) {
	resp, err := http.Post(bulletURL, "application/json", bytes.NewBuffer([]byte("{}")))
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	var r BulletResp
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return "", 0, err
	}

	s := r.Data.InstanceServers[0]
	url := fmt.Sprintf("%s?token=%s&connectId=%d",
		s.Endpoint,
		r.Data.Token,
		time.Now().UnixNano(),
	)

	return url, time.Duration(s.PingInterval) * time.Millisecond, nil
}

func main() {
	_ = context.Background()

	wsURL, pingInterval, err := getWS()
	if err != nil {
		log.Fatal(err)
	}

	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer conn.Close()

	log.Println("✅ connected")

	// подписка
	sub := map[string]any{
		"id":       time.Now().UnixNano(),
		"type":     "subscribe",
		"topic":    "/market/ticker:BTC-USDT",
		"response": true,
	}

	if err := conn.WriteJSON(sub); err != nil {
		log.Fatal("subscribe error:", err)
	}

	// ping loop (ВАЖНО: JSON ping!)
	go func() {
		t := time.NewTicker(pingInterval)
		for range t.C {
			_ = conn.WriteJSON(map[string]any{
				"id":   time.Now().UnixNano(),
				"type": "ping",
			})
		}
	}()

	// read loop
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("read error:", err)
			return
		}

		fmt.Println(string(msg))
	}
}


gaz358@gaz358-BOD-WXX9:~/myprog/crypt_proto/cmd/arb/arb_test$ go run .
2025/12/26 00:27:53 ✅ connected
{"id":"1766698072350895818","type":"welcome"}
{"id":"1766698073681334730","type":"ack"}
{"topic":"/market/ticker:BTC-USDT","type":"message","subject":"trade.ticker","data":{"bestAsk":"87953.9","bestAskSize":"0.59800291","bestBid":"87953.8","bestBidSize":"0.31142342","price":"87953.9","sequence":"25178475757","size":"0.00056846","time":1766698072551}}
{"topic":"/market/ticker:BTC-USDT","type":"message","subject":"trade.ticker","data":{"bestAsk":"87953.9","bestAskSize":"0.59800291","bestBid":"87953.8","bestBidSize":"0.31142342","price":"87953.9","sequence":"25178475765","size":"0.00056846","time":1766698072551}}




