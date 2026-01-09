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
	subRate      = 120 * time.Millisecond
	pingInterval = 20 * time.Second
)

// ================= POOL =================
type KuCoinCollector struct {
	ctx        context.Context
	cancel     context.CancelFunc
	wsList     []*kucoinWS
	mem        *queue.MemoryStore
	triangles  []calculator.Triangle
	symbolMap  map[string]string // map csvSymbol -> normalizedSymbol
}

// ================= WS =================
type kucoinWS struct {
	id      int
	conn    *websocket.Conn
	symbols []string
	last    map[string][2]float64
	mu      sync.Mutex
}

// ================= CONSTRUCTOR =================
func NewKuCoinCollectorFromCSV(path string) (*KuCoinCollector, error) {
	symbols, symbolMap, err := readPairsFromCSV(path)
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
		symbolMap: symbolMap,
	}, nil
}

func (kc *KuCoinCollector) Name() string {
	return "KuCoin"
}

// Start принимает MemoryStore
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

// ================= CONNECT =================
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

// ================= SUBSCRIBE =================
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
			log.Printf("[KuCoin WS %d] subscribe error %s\n", ws.id, s)
		}
	}
}

// ================= PING =================
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

// ================= READ =================
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

// ================= HANDLE =================
func (ws *kucoinWS) handle(kc *KuCoinCollector, msg []byte) {
	if gjson.GetBytes(msg, "type").String() != "message" {
		return
	}

	topic := gjson.GetBytes(msg, "topic").String()
	if !strings.HasPrefix(topic, "/market/ticker:") {
		return
	}

	csvSymbol := strings.TrimPrefix(topic, "/market/ticker:")
	symbol, ok := kc.symbolMap[csvSymbol]
	if !ok {
		return
	}

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

	kc.mem.Push(&models.MarketData{
		Exchange:  "KuCoin",
		Symbol:    symbol,
		Bid:       bid,
		Ask:       ask,
		BidSize:   data.Get("bestBidSize").Float(),
		AskSize:   data.Get("bestAskSize").Float(),
		Timestamp: time.Now().UnixMilli(),
	})
}

// ================= CSV =================
func readPairsFromCSV(path string) ([]string, map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	rows, err := r.ReadAll()
	if err != nil {
		return nil, nil, err
	}

	set := make(map[string]struct{})
	symbolMap := make(map[string]string)
	for _, row := range rows[1:] {
		for i := 3; i <= 5 && i < len(row); i++ {
			csvSym := parseLeg(row[i])
			if csvSym != "" {
				set[csvSym] = struct{}{}
				symbolMap[csvSym] = normalize(csvSym)
			}
		}
	}

	var res []string
	for k := range set {
		res = append(res, k)
	}
	return res, symbolMap, nil
}

// ================= HELPERS =================
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






