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


2025/12/26 00:33:40 ✅ connected
{"id":"1766698419918079443","type":"welcome"}
{"id":"1766698420962179455","type":"ack"}
{"topic":"/market/ticker:BTC-USDT","type":"message","subject":"trade.ticker","data":{"bestAsk":"87956.1","bestAskSize":"0.26876662","bestBid":"87956","bestBidSize":"0.74174027","price":"87956","sequence":"25178512486","size":"0.00002263","time":1766698420140}}
{"topic":"/market/ticker:BTC-USDT","type":"message","subject":"trade.ticker","data":{"bestAsk":"87956.1","bestAskSize":"0.26876662","bestBid":"87956","bestBidSize":"0.74174027","price":"87956","sequence":"25178512491","size":"0.00002263","time":1766698420140}}
{"topic":"/market/ticker:BTC-USDT","type":"message","subject":"trade.ticker","data":{"bestAsk":"87956.1","bestAskSize":"0.26876662","bestBid":"87956","bestBidSize":"0.74174027","price":"87956","sequence":"25178512501","size":"0.00002263","time":1766698420140}}
{"topic":"/market/ticker:BTC-USDT","type":"message","subject":"trade.ticker","data":{"bestAsk":"87956.1","bestAskSize":"0.26876662","bestBid":"87956","bestBidSize":"0.74174027","price":"87956","sequence":"25178512505","size":"0.00002263","time":1766698420140}}
{"topic":"/market/ticker:BTC-USDT","type":"message","subject":"trade.ticker","data":{"bestAsk":"87956.1","bestAskSize":"0.26876662","bestBid":"87956","bestBidSize":"0.74174027","price":"87956","sequence":"25178512509","size":"0.00002263","time":1766698420140}}
{"topic":"/market/ticker:BTC-USDT","type":"message","subject":"trade.ticker","data":{"bestAsk":"87956.1","bestAskSize":"0.26876662","bestBid":"87956","bestBidSize":"0.74174027","price":"87956","sequence":"25178512519","size":"0.00002263","time":1766698420140}}
{"topic":"/market/ticker:BTC-USDT","type":"message","subject":"trade.ticker","data":{"bestAsk":"87956.1","bestAskSize":"0.26876662","bestBid":"87956","bestBidSize":"0.73628248","price":"87956","sequence":"25178512543","size":"0.00002263","time":1766698420140}}
{"topic":"/market/ticker:BTC-USDT","type":"message","subject":"trade.ticker","data":{"bestAsk":"87956.1","bestAskSize":"0.26876662","bestBid":"87956","bestBidSize":"0.73628248","price":"87956","sequence":"25178512547","size":"0.00002263","time":1766698420140}}
{"topic":"/market/ticker:BTC-USDT","type":"message","subject":"trade.ticker","data":{"bestAsk":"87956.1","bestAskSize":"0.27311087","bestBid":"87956","bestBidSize":"0.73628248","price":"87956","sequence":"25178512549","size":"0.00002263","time":1766698420140}}
{"topic":"/market/ticker:BTC-USDT","type":"message","subject":"trade.ticker","data":{"bestAsk":"87956.1","bestAskSize":"0.27311087","bestBid":"87956","bestBidSize":"0.73628248","price":"87956","sequence":"25178512559","size":"0.00002263","time":1766698420140}}
{"topic":"/market/ticker:BTC-USDT","type":"message","subject":"trade.ticker","data":{"bestAsk":"87956.1","bestAskSize":"0.27311087","bestBid":"87956","bestBidSize":"0.73628248","price":"87956","sequence":"25178512569","size":"0.00002263","time":1766698420140}}
{"topic":"/market/ticker:BTC-USDT","type":"message","subject":"trade.ticker","data":{"bestAsk":"87956.1","bestAskSize":"0.27311087","bestBid":"87956","bestBidSize":"0.73628248","price":"87956","sequence":"25178512570","size":"0.00002263","time":1766698420140}}
{"topic":"/market/ticker:BTC-USDT","type":"message","subject":"trade.ticker","data":{"bestAsk":"87956.1","bestAskSize":"0.27311087","bestBid":"87956","bestBidSize":"0.73628248","price":"87956","sequence":"25178512575","size":"0.00002263","time":1766698420140}}
{"topic":"/market/ticker:BTC-USDT","type":"message","subject":"trade.ticker","data":{"bestAsk":"87956.1","bestAskSize":"0.27311087","bestBid":"87956","bestBidSize":"0.73628248","price":"87956","sequence":"25178512583","size":"0.00002263","time":1766698420140}}
{"topic":"/market/ticker:BTC-USDT","type":"message","subject":"trade.ticker","data":{"bestAsk":"87956.1","bestAskSize":"0.27311087","bestBid":"87956","bestBidSize":"0.73628248","price":"87956","sequence":"25178512595","size":"0.00002263","time":1766698420140}}
{"topic":"/market/ticker:BTC-USDT","type":"message","subject":"trade.ticker","data":{"bestAsk":"87956.1","bestAskSize":"0.27311087","bestBid":"87956","bestBidSize":"0.73628248","price":"87956","sequence":"25178512602","size":"0.00002263","time":1766698420140}}
{"topic":"/market/ticker:BTC-USDT","type":"message","subject":"trade.ticker","data":{"bestAsk":"87956.1","bestAskSize":"0.27311087","bestBid":"87956","bestBidSize":"0.73628248","price":"87956","sequence":"25178512605","size":"0.00002263","time":1766698420140}}
{"topic":"/market/ticker:BTC-USDT","type":"message","subject":"trade.ticker","data":{"bestAsk":"87956.1","bestAskSize":"0.27311087","bestBid":"87956","bestBidSize":"0.73628248","price":"87956","sequence":"25178512612","size":"0.00002263","time":1766698420140}}
{"topic":"/market/ticker:BTC-USDT","type":"message","subject":"trade.ticker","data":{"bestAsk":"87956.1","bestAskSize":"0.27311087","bestBid":"87956","bestBidSize":"0.73628248","price":"87956","sequence":"25178512617","size":"0.00002263","time":1766698420140}}
{"topic":"/market/ticker:BTC-USDT","type":"message","subject":"trade.ticker","data":{"bestAsk":"87956.1","bestAskSize":"0.27301401","bestBid":"87956","bestBidSize":"0.73628248","price":"87956.1","sequence":"25178512645","size":"0.00009686","time":1766698424245}}
{"topic":"/market/ticker:BTC-USDT","type":"message","subject":"trade.ticker","data":{"bestAsk":"87956.1","bestAskSize":"0.27301401","bestBid":"87956","bestBidSize":"0.73628248","price":"87956.1","sequence":"25178512652","size":"0.00009686","time":1766698424245}}
{"topic":"/market/ticker:BTC-USDT","type":"message","subject":"trade.ticker","data":{"bestAsk":"87956.1","bestAskSize":"0.27301401","bestBid":"87956","bestBidSize":"0.73628248","price":"87956.1","sequence":"25178512662","size":"0.00009686","time":1766698424245}}
^Csignal: interrupt



