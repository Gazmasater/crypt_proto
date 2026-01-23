Если оставить только нужное:

p99 execution latency
Micro-volatility (100 мс)
Fill ratio
Capture rate
Inventory drift




Название API
9623527002

696935c42a6dcd00013273f2
b348b686-55ff-4290-897b-02d55f815f65




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




package collector

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"crypt_proto/pkg/models"

	"github.com/gorilla/websocket"
	"github.com/tidwall/gjson"
)

const (
	maxSubsPerWS   = 126
	subRate        = 120 * time.Millisecond
	pingInterval   = 20 * time.Second
	deadAfter      = 120 * time.Second
	unsubInterval  = 60 * time.Second
)

type KuCoinCollector struct {
	ctx    context.Context
	cancel context.CancelFunc
	wsList []*kucoinWS
	out    chan<- *models.MarketData
}

type quoteState struct {
	bid  float64
	ask  float64
	ts   int64
}

type kucoinWS struct {
	id      int
	conn    *websocket.Conn
	symbols []string

	mu        sync.RWMutex
	last      map[string]quoteState
	subscribed map[string]struct{}
}

func NewKuCoinCollectorFromCSV(path string) (*KuCoinCollector, []string, error) {
	symbols, err := readPairsFromCSV(path)
	if err != nil {
		return nil, nil, err
	}
	if len(symbols) == 0 {
		return nil, nil, fmt.Errorf("no symbols")
	}

	ctx, cancel := context.WithCancel(context.Background())
	var wsList []*kucoinWS

	for i := 0; i < len(symbols); i += maxSubsPerWS {
		end := i + maxSubsPerWS
		if end > len(symbols) {
			end = len(symbols)
		}

		wsList = append(wsList, &kucoinWS{
			id:         len(wsList),
			symbols:    symbols[i:end],
			last:       make(map[string]quoteState),
			subscribed: make(map[string]struct{}),
		})
	}

	return &KuCoinCollector{
		ctx:    ctx,
		cancel: cancel,
		wsList: wsList,
	}, symbols, nil
}

func (c *KuCoinCollector) Name() string { return "KuCoin" }

func (c *KuCoinCollector) Start(out chan<- *models.MarketData) error {
	c.out = out

	for _, ws := range c.wsList {
		if err := ws.connect(); err != nil {
			return err
		}
		go ws.readLoop(c)
		go ws.subscribeLoop()
		go ws.pingLoop()
		go ws.unsubscribeLoop()
	}

	log.Printf("[KuCoin] started with %d WS\n", len(c.wsList))
	return nil
}

func (c *KuCoinCollector) Stop() error {
	c.cancel()
	for _, ws := range c.wsList {
		if ws.conn != nil {
			_ = ws.conn.Close()
		}
	}
	return nil
}

/* ================= WS ================= */

func (ws *kucoinWS) connect() error {
	req, _ := http.NewRequest("POST", "https://api.kucoin.com/api/v1/bullet-public", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var r struct {
		Data struct {
			Token           string `json:"token"`
			InstanceServers []struct {
				Endpoint string `json:"endpoint"`
			} `json:"instanceServers"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return err
	}

	url := fmt.Sprintf(
		"%s?token=%s&connectId=%d",
		r.Data.InstanceServers[0].Endpoint,
		r.Data.Token,
		time.Now().UnixNano(),
	)

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return err
	}

	ws.conn = conn
	log.Printf("[KuCoin WS %d] connected\n", ws.id)
	return nil
}

func (ws *kucoinWS) subscribeLoop() {
	t := time.NewTicker(subRate)
	defer t.Stop()

	for _, s := range ws.symbols {
		<-t.C
		_ = ws.conn.WriteJSON(map[string]any{
			"id":       time.Now().UnixNano(),
			"type":     "subscribe",
			"topic":    "/market/ticker:" + s,
			"response": true,
		})

		ws.mu.Lock()
		ws.subscribed[s] = struct{}{}
		ws.mu.Unlock()
	}
}

func (ws *kucoinWS) unsubscribeLoop() {
	t := time.NewTicker(unsubInterval)
	defer t.Stop()

	for range t.C {
		now := time.Now().UnixMilli()
		deadBefore := now - deadAfter.Milliseconds()

		ws.mu.Lock()
		for sym := range ws.subscribed {
			qs, ok := ws.last[sym]
			if !ok || qs.ts < deadBefore {
				_ = ws.conn.WriteJSON(map[string]any{
					"id":    time.Now().UnixNano(),
					"type":  "unsubscribe",
					"topic": "/market/ticker:" + sym,
				})
				delete(ws.subscribed, sym)
				delete(ws.last, sym)
				log.Printf("[KuCoin WS %d] unsubscribe %s\n", ws.id, sym)
			}
		}
		ws.mu.Unlock()
	}
}

func (ws *kucoinWS) pingLoop() {
	t := time.NewTicker(pingInterval)
	defer t.Stop()
	for range t.C {
		_ = ws.conn.WriteJSON(map[string]any{
			"id":   time.Now().UnixNano(),
			"type": "ping",
		})
	}
}

func (ws *kucoinWS) readLoop(c *KuCoinCollector) {
	for {
		_, msg, err := ws.conn.ReadMessage()
		if err != nil {
			log.Printf("[KuCoin WS %d] read error: %v\n", ws.id, err)
			return
		}
		ws.handle(c, msg)
	}
}

func (ws *kucoinWS) handle(c *KuCoinCollector, msg []byte) {
	if gjson.GetBytes(msg, "type").String() != "message" {
		return
	}

	topic := gjson.GetBytes(msg, "topic").String()
	if !strings.HasPrefix(topic, "/market/ticker:") {
		return
	}

	symbol := strings.TrimPrefix(topic, "/market/ticker:")
	data := gjson.GetBytes(msg, "data")

	bid := data.Get("bestBid").Float()
	ask := data.Get("bestAsk").Float()
	if bid == 0 || ask == 0 {
		return
	}

	bidSize := data.Get("bestBidSize").Float()
	askSize := data.Get("bestAskSize").Float()

	now := time.Now().UnixMilli()

	ws.mu.Lock()
	prev, ok := ws.last[symbol]
	if ok && prev.bid == bid && prev.ask == ask {
		ws.last[symbol] = quoteState{bid, ask, now}
		ws.mu.Unlock()
		return
	}
	ws.last[symbol] = quoteState{bid, ask, now}
	ws.mu.Unlock()

	c.out <- &models.MarketData{
		Exchange:  "KuCoin",
		Symbol:    symbol,
		Bid:       bid,
		Ask:       ask,
		BidSize:   bidSize,
		AskSize:   askSize,
		Timestamp: now,
	}
}

/* ================= CSV ================= */

func readPairsFromCSV(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rows, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return nil, err
	}

	set := make(map[string]struct{})
	for _, row := range rows[1:] {
		for i := 3; i <= 5 && i < len(row); i++ {
			if p := parseLeg(row[i]); p != "" {
				set[p] = struct{}{}
			}
		}
	}

	res := make([]string, 0, len(set))
	for k := range set {
		res = append(res, k)
	}
	return res, nil
}

func parseLeg(s string) string {
	parts := strings.Fields(strings.ToUpper(strings.TrimSpace(s)))
	if len(parts) < 2 {
		return ""
	}
	p := strings.Split(parts[1], "/")
	if len(p) != 2 {
		return ""
	}
	return p[0] + "-" + p[1]
}


gaz358@gaz358-BOD-WXX9:~/myprog/crypt_proto/cmd/arb$ go run .
2026/01/24 00:00:01 pprof on http://localhost:6060/debug/pprof/
2026/01/24 00:00:07 [KuCoin WS 0] connected
2026/01/24 00:00:09 [KuCoin WS 1] connected
2026/01/24 00:00:09 [KuCoin] started with 2 WS
2026/01/24 00:00:09 [Main] KuCoinCollector started
2026/01/24 00:00:09 [Calculator] indexed 214 symbols
panic: concurrent write to websocket connection

goroutine 71 [running]:
github.com/gorilla/websocket.(*messageWriter).flushFrame(0xc000394930, 0x1, {0x0?, 0x0?, 0x0?})
        /home/gaz358/go/pkg/mod/github.com/gorilla/websocket@v1.5.3/conn.go:617 +0x4af
github.com/gorilla/websocket.(*messageWriter).Close(0x4267bc?)
        /home/gaz358/go/pkg/mod/github.com/gorilla/websocket@v1.5.3/conn.go:731 +0x35
github.com/gorilla/websocket.(*Conn).beginMessage(0xc00016b760, 0xc0005d03f0, 0x1)
        /home/gaz358/go/pkg/mod/github.com/gorilla/websocket@v1.5.3/conn.go:480 +0x37
github.com/gorilla/websocket.(*Conn).NextWriter(0xc00016b760, 0x1)
        /home/gaz358/go/pkg/mod/github.com/gorilla/websocket@v1.5.3/conn.go:520 +0x3f
github.com/gorilla/websocket.(*Conn).WriteJSON(0x85f200?, {0x85f200, 0xc0005d03c0})
        /home/gaz358/go/pkg/mod/github.com/gorilla/websocket@v1.5.3/json.go:24 +0x34
crypt_proto/internal/collector.(*kucoinWS).unsubscribeLoop(0xc0001d23c0)
        /home/gaz358/myprog/crypt_proto/internal/collector/kucoin_collector.go:185 +0x398
created by crypt_proto/internal/collector.(*KuCoinCollector).Start in goroutine 1
        /home/gaz358/myprog/crypt_proto/internal/collector/kucoin_collector.go:97 +0x5d
exit status 2


