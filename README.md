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









Структура каталога та же (crypt_proto), внутри:

main.go — только запуск и pprof

config.go — конфиг + debug-логгер

domain.go — типы и работа с треугольниками

proto_decoder.go — разбор protobuf от MEXC

ws.go — работа с WebSocket

arb.go — расчёт треугольников, логирование, консюмер

Все файлы ниже можно просто создать рядом и вставить как есть.



package main

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

/* =========================  WS SUBSCRIBER  ========================= */

func buildTopics(symbols []string, interval string) []string {
	topics := make([]string, 0, len(symbols))
	for _, s := range symbols {
		topics = append(topics, "spot@public.aggre.bookTicker.v3.api.pb@"+interval+"@"+s)
	}
	return topics
}

func runPublicBookTickerWS(
	ctx context.Context,
	wg *sync.WaitGroup, // <-- обычный WaitGroup
	connID int,
	symbols []string,
	interval string,
	out chan<- Event,
) {
	defer wg.Done()

	const (
		baseRetry = 2 * time.Second
		maxRetry  = 30 * time.Second
	)

	urlWS := "wss://wbs-api.mexc.com/ws"
	topics := buildTopics(symbols, interval)
	retry := baseRetry

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		conn, _, err := websocket.DefaultDialer.Dial(urlWS, nil)
		if err != nil {
			log.Printf("[WS #%d] dial err: %v (retry in %v)", connID, err, retry)
			time.Sleep(retry)
			retry = nextRetry(retry, maxRetry)
			continue
		}

		log.Printf("[WS #%d] connected to %s (symbols: %d)", connID, urlWS, len(symbols))
		retry = baseRetry

		_ = conn.SetReadDeadline(time.Now().Add(90 * time.Second))

		var lastPing time.Time
		conn.SetPongHandler(func(appData string) error {
			rtt := time.Since(lastPing)
			dlog("[WS #%d] Pong через %v", connID, rtt)
			return conn.SetReadDeadline(time.Now().Add(90 * time.Second))
		})

		stopPing := make(chan struct{})
		go pingLoop(connID, conn, &lastPing, stopPing)

		if err := sendSubscription(conn, topics, connID); err != nil {
			close(stopPing)
			_ = conn.Close()
			time.Sleep(retry)
			retry = nextRetry(retry, maxRetry)
			continue
		}

		if !readLoop(ctx, connID, conn, out) {
			close(stopPing)
			_ = conn.Close()
			time.Sleep(retry)
			retry = nextRetry(retry, maxRetry)
			continue
		}
	}
}

func nextRetry(cur, max time.Duration) time.Duration {
	cur *= 2
	if cur > max {
		return max
	}
	return cur
}

func pingLoop(connID int, conn *websocket.Conn, lastPing *time.Time, stop <-chan struct{}) {
	t := time.NewTicker(45 * time.Second)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			*lastPing = time.Now()
			if err := conn.WriteControl(
				websocket.PingMessage,
				[]byte("hb"),
				time.Now().Add(5*time.Second),
			); err != nil {
				dlog("[WS #%d] ping error: %v", connID, err)
				return
			}
		case <-stop:
			return
		}
	}
}

func sendSubscription(conn *websocket.Conn, topics []string, connID int) error {
	sub := map[string]any{
		"method": "SUBSCRIPTION",
		"params": topics,
		"id":     time.Now().Unix(),
	}
	if err := conn.WriteJSON(sub); err != nil {
		log.Printf("[WS #%d] subscribe send err: %v", connID, err)
		return err
	}
	log.Printf("[WS #%d] SUB -> %d topics", connID, len(topics))
	return nil
}

func readLoop(ctx context.Context, connID int, conn *websocket.Conn, out chan<- Event) bool {
	for {
		mt, raw, err := conn.ReadMessage()
		if err != nil {
			log.Printf("[WS #%d] read err: %v (reconnect)", connID, err)
			return false
		}

		switch mt {
		case websocket.TextMessage:
			handleTextMessage(connID, raw)
		case websocket.BinaryMessage:
			sym, q, ok := parsePBQuote(raw)
			if !ok {
				continue
			}
			ev := Event{
				Symbol: sym,
				Bid:    q.Bid,
				Ask:    q.Ask,
				BidQty: q.BidQty,
				AskQty: q.AskQty,
			}
			select {
			case out <- ev:
			case <-ctx.Done():
				return true
			}
		default:
			// игнор
		}
	}
}

func handleTextMessage(connID int, raw []byte) {
	if !debug {
		return
	}
	var tmp any
	if err := json.Unmarshal(raw, &tmp); err == nil {
		j, _ := json.Marshal(tmp)
		dlog("[WS #%d TEXT] %s", connID, string(j))
	} else {
		dlog("[WS #%d TEXT RAW] %s", connID, string(raw))
	}
}




func runPublicBookTickerWS(ctx context.Context, wg *sync.WaitGroup, ...)






