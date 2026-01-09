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
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/tidwall/gjson"

	"crypt_proto/internal/calculator"
	"crypt_proto/internal/queue"
	"crypt_proto/pkg/models"
)

const (
	maxSubsPerWS = 90
	subRate      = 120 * time.Millisecond // ~8 подписок/сек
	pingInterval = 20 * time.Second
)

/* ================= KUCOIN COLLECTOR ================= */
type KuCoinCollector struct {
	ctx       context.Context
	cancel    context.CancelFunc
	wsList    []*kucoinWS
	triangles []calculator.Triangle
	mem       *queue.MemoryStore

	OnUpdate func(symbol string) // callback на апдейт котировки
}

/* ================= WS ================= */
type kucoinWS struct {
	id      int
	conn    *websocket.Conn
	symbols []string
	last    map[string][2]float64
	mu      sync.Mutex
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

func (kc *KuCoinCollector) Name() string {
	return "KuCoin"
}

func (kc *KuCoinCollector) Triangles() []calculator.Triangle {
	return kc.triangles
}

/* ================= START / STOP ================= */
func (kc *KuCoinCollector) Start(mem *queue.MemoryStore) error {
	kc.mem = mem

	for _, ws := range kc.wsList {
		if err := ws.connect(); err != nil {
			return err
		}
		go ws.readLoop(kc)
		go ws.pingLoop()
		go ws.subscribeLoop()
	}

	log.Printf("[KuCoin] started with %d WS\n", len(kc.wsList))
	return nil
}

func (kc *KuCoinCollector) Stop() error {
	kc.cancel()
	for _, ws := range kc.wsList {
		if ws.conn != nil {
			ws.conn.Close()
		}
	}
	return nil
}

/* ================= WS CONNECT ================= */
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

	url := fmt.Sprintf("%s?token=%s&connectId=%d",
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
	time.Sleep(1 * time.Second) // ждём welcome
	t := time.NewTicker(subRate)
	defer t.Stop()

	for _, s := range ws.symbols {
		<-t.C
		topic := "/market/ticker:" + s
		err := ws.conn.WriteJSON(map[string]any{
			"id":       time.Now().UnixNano(),
			"type":     "subscribe",
			"topic":    topic,
			"response": true,
		})
		if err != nil {
			log.Printf("[KuCoin WS %d] subscribe error %s: %v\n", ws.id, s, err)
		} else {
			log.Printf("[KuCoin WS %d] subscribed %s\n", ws.id, s)
		}
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
func (ws *kucoinWS) readLoop(kc *KuCoinCollector) {
	for {
		_, msg, err := ws.conn.ReadMessage()
		if err != nil {
			log.Printf("[KuCoin WS %d] read error: %v\n", ws.id, err)
			return
		}
		ws.handle(kc, msg)
	}
}

/* ================= HANDLE ================= */
func (ws *kucoinWS) handle(kc *KuCoinCollector, msg []byte) {
	if gjson.GetBytes(msg, "type").String() != "message" {
		return
	}

	topic := gjson.GetBytes(msg, "topic").String()
	if !strings.HasPrefix(topic, "/market/ticker:") {
		return
	}

	symbol := normalize(strings.TrimPrefix(topic, "/market/ticker:"))

	data := gjson.GetBytes(msg, "data")
	bid := data.Get("bestBid").Float()
	ask := data.Get("bestAsk").Float()
	if bid == 0 || ask == 0 {
		return
	}

	ws.mu.Lock()
	last := ws.last[symbol]
	if last[0] == bid && last[1] == ask {
		ws.mu.Unlock()
		return
	}
	ws.last[symbol] = [2]float64{bid, ask}
	ws.mu.Unlock()

	// ---- Push в MemoryStore ----
	kc.mem.Push(&models.MarketData{
		Exchange:  "KuCoin",
		Symbol:    symbol,
		Bid:       bid,
		Ask:       ask,
		BidSize:   data.Get("bestBidSize").Float(),
		AskSize:   data.Get("bestAskSize").Float(),
		Timestamp: time.Now().UnixMilli(),
	})

	// ---- Callback на апдейт ----
	if kc.OnUpdate != nil {
		kc.OnUpdate(symbol)
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
			p := parseLeg(row[i])
			if p != "" {
				set[p] = struct{}{}
			}
		}
	}

	var res []string
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

func normalize(s string) string {
	parts := strings.Split(s, "-")
	return parts[0] + "/" + parts[1]
}



gaz358@gaz358-BOD-WXX9:~/myprog/crypt_proto$    go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
Fetching profile over HTTP from http://localhost:6060/debug/pprof/profile?seconds=30
Saved profile in /home/gaz358/pprof/pprof.arb.samples.cpu.058.pb.gz
File: arb
Build ID: 93785d6bc3dd1aee89cd9819484725d29d217b9e
Type: cpu
Time: 2026-01-09 18:09:31 MSK
Duration: 30.15s, Total samples = 4.09s (13.56%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 2060ms, 50.37% of 4090ms total
Dropped 133 nodes (cum <= 20.45ms)
Showing top 10 nodes out of 137
      flat  flat%   sum%        cum   cum%
     920ms 22.49% 22.49%      920ms 22.49%  internal/runtime/syscall.Syscall6
     230ms  5.62% 28.12%      620ms 15.16%  strings.Fields
     210ms  5.13% 33.25%      210ms  5.13%  runtime.futex
     130ms  3.18% 36.43%      130ms  3.18%  strings.ToUpper
     120ms  2.93% 39.36%      480ms 11.74%  runtime.scanobject
     100ms  2.44% 41.81%      120ms  2.93%  runtime.findObject
      90ms  2.20% 44.01%       90ms  2.20%  aeshashbody
      90ms  2.20% 46.21%       90ms  2.20%  runtime.memclrNoHeapPointers
      90ms  2.20% 48.41%       90ms  2.20%  runtime.nextFreeFast (inline)
      80ms  1.96% 50.37%       80ms  1.96%  runtime.(*mspan).base (inline)
(pprof) 





