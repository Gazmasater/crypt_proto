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
	"time"

	"github.com/tidwall/gjson"

	"crypt_proto/internal/calculator"
	"crypt_proto/pkg/models"

	"github.com/gorilla/websocket"
)

const (
	maxSubsPerWS = 90
	subRate      = 120 * time.Millisecond
	pingInterval = 20 * time.Second
)

/* ================= POOL ================= */

type KuCoinCollector struct {
	ctx       context.Context
	cancel    context.CancelFunc
	wsList    []*kucoinWS
	out       chan<- *models.MarketData
	triangles []calculator.Triangle
}

/* ================= WS ================= */

type kucoinWS struct {
	id      int
	conn    *websocket.Conn
	symbols []string
	last    map[string][2]float64 // ONLY readLoop touches this
}

/* ================= CONSTRUCTOR ================= */

func NewKuCoinCollectorFromCSV(path string) (*KuCoinCollector, error) {
	symbols, err := readPairsFromCSV(path)
	if err != nil {
		return nil, err
	}
	if len(symbols) == 0 {
		return nil, fmt.Errorf("no symbols")
	}

	triangles, _ := calculator.ParseTrianglesFromCSV(path)

	ctx, cancel := context.WithCancel(context.Background())
	var wsList []*kucoinWS

	for i := 0; i < len(symbols); i += maxSubsPerWS {
		end := i + maxSubsPerWS
		if end > len(symbols) {
			end = len(symbols)
		}
		wsList = append(wsList, &kucoinWS{
			id:      len(wsList),
			symbols: symbols[i:end],
			last:    make(map[string][2]float64),
		})
	}

	return &KuCoinCollector{
		ctx:       ctx,
		cancel:    cancel,
		wsList:    wsList,
		triangles: triangles,
	}, nil
}

func (kc *KuCoinCollector) Triangles() []calculator.Triangle {
	return kc.triangles
}

/* ================= INTERFACE ================= */

func (c *KuCoinCollector) Name() string {
	return "KuCoin"
}

func (c *KuCoinCollector) Start(out chan<- *models.MarketData) error {
	c.out = out

	for _, ws := range c.wsList {
		if err := ws.connect(); err != nil {
			return err
		}
		go ws.readLoop(c)
		go ws.subscribeLoop()
		go ws.pingLoop()
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

/* ================= CONNECT ================= */

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

/* ================= SUBSCRIBE ================= */

func (ws *kucoinWS) subscribeLoop() {
	time.Sleep(time.Second)

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
	}
}

/* ================= PING ================= */

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

/* ================= READ ================= */

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

/* ================= HANDLE ================= */

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
	bidSize := data.Get("bestBidSize").Float()
	askSize := data.Get("bestAskSize").Float()

	if bid == 0 || ask == 0 {
		return
	}

	last := ws.last[symbol]
	if last[0] == bid && last[1] == ask {
		return
	}
	ws.last[symbol] = [2]float64{bid, ask}

	c.out <- &models.MarketData{
		Exchange:  "KuCoin",
		Symbol:    symbol,
		Bid:       bid,
		Ask:       ask,
		BidSize:   bidSize,
		AskSize:   askSize,
		Timestamp: time.Now().UnixMilli(),
	}
}

/* ================= CSV ================= */

func readPairsFromCSV(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	rows, err := r.ReadAll()
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

/* ================= HELPERS ================= */

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



Fetching profile over HTTP from http://localhost:6060/debug/pprof/profile?seconds=30
Saved profile in /home/gaz358/pprof/pprof.arb.samples.cpu.104.pb.gz
File: arb
Build ID: 63302cd9d8087f6f82a53584f4286a4f35d4a4a6
Type: cpu
Time: 2026-01-13 12:28:29 MSK
Duration: 30s, Total samples = 1.11s ( 3.70%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 710ms, 63.96% of 1110ms total
Showing top 10 nodes out of 145
      flat  flat%   sum%        cum   cum%
     420ms 37.84% 37.84%      420ms 37.84%  internal/runtime/syscall.Syscall6
      90ms  8.11% 45.95%       90ms  8.11%  runtime.futex
      30ms  2.70% 48.65%       30ms  2.70%  aeshashbody
      30ms  2.70% 51.35%       70ms  6.31%  github.com/tidwall/gjson.parseObject
      30ms  2.70% 54.05%       30ms  2.70%  github.com/tidwall/gjson.parseSquash
      30ms  2.70% 56.76%       30ms  2.70%  runtime.mallocgcSmallNoscan
      20ms  1.80% 58.56%       90ms  8.11%  github.com/tidwall/gjson.Get
      20ms  1.80% 60.36%       20ms  1.80%  runtime.(*mcache).prepareForSweep
      20ms  1.80% 62.16%       20ms  1.80%  runtime.casgstatus
      20ms  1.80% 63.96%       20ms  1.80%  runtime.nanotime1
(pprof) 





