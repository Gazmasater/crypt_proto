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
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"crypt_proto/pkg/models"

	"github.com/gorilla/websocket"
)

const (
	maxSubsPerWS = 50
	subRate      = 120 * time.Millisecond
	pingInterval = 20 * time.Second
)

/* ================= POOL ================= */

type KuCoinCollector struct {
	ctx    context.Context
	cancel context.CancelFunc

	wsList []*kucoinWS
	out    chan<- *models.MarketData
	pool   *sync.Pool
}

/* ================= WS ================= */

type kucoinWS struct {
	id      int
	conn    *websocket.Conn
	symbols []string

	// [bid, ask, bidSize, askSize]
	last map[string][4]float64
	mu   sync.Mutex
}

/* ================= CONSTRUCTOR ================= */

func NewKuCoinCollectorFromCSV(path string, pool *sync.Pool) (*KuCoinCollector, error) {
	symbols, err := readPairsFromCSV(path)
	if err != nil {
		return nil, err
	}
	if len(symbols) == 0 {
		return nil, nil
	}

	ctx, cancel := context.WithCancel(context.Background())

	wsList := make([]*kucoinWS, 0)
	for i := 0; i < len(symbols); i += maxSubsPerWS {
		end := i + maxSubsPerWS
		if end > len(symbols) {
			end = len(symbols)
		}
		wsList = append(wsList, &kucoinWS{
			id:      len(wsList),
			symbols: symbols[i:end],
			last:    make(map[string][4]float64),
		})
	}

	return &KuCoinCollector{
		ctx:    ctx,
		cancel: cancel,
		wsList: wsList,
		pool:   pool,
	}, nil
}

/* ================= INTERFACE ================= */

func (c *KuCoinCollector) Name() string { return "KuCoin" }

func (c *KuCoinCollector) Start(out chan<- *models.MarketData) error {
	c.out = out
	for _, ws := range c.wsList {
		if err := ws.connect(); err != nil {
			return err
		}
		go ws.readLoop(c)
		go ws.pingLoop()
		go ws.subscribeLoop()
	}
	log.Printf("[KuCoin] started with %d WS\n", len(c.wsList))
	return nil
}

func (c *KuCoinCollector) Stop() error {
	c.cancel()
	for _, ws := range c.wsList {
		if ws.conn != nil {
			ws.conn.Close()
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

	url := r.Data.InstanceServers[0].Endpoint + "?token=" + r.Data.Token + "&connectId=" + strconv.FormatInt(time.Now().UnixNano(), 10)
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
		topic := "/market/ticker:" + s
		_ = ws.conn.WriteJSON(map[string]any{
			"id":       time.Now().UnixNano(),
			"type":     "subscribe",
			"topic":    topic,
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

type kucoinMsg struct {
	Type  string `json:"type"`
	Topic string `json:"topic"`
	Data  struct {
		BestBid  string `json:"bestBid"`
		BestAsk  string `json:"bestAsk"`
		BuySize  string `json:"buySize"`
		SellSize string `json:"sellSize"`
	} `json:"data"`
}

func (ws *kucoinWS) handle(c *KuCoinCollector, msg []byte) {
	var m kucoinMsg
	if json.Unmarshal(msg, &m) != nil {
		return
	}
	if m.Type != "message" {
		return
	}

	// normalize symbol: BTC-USDT -> BTC/USDT
	sym := strings.Replace(m.Topic[strings.IndexByte(m.Topic, ':')+1:], "-", "/", 1)

	// parse floats
	bid, err1 := strconv.ParseFloat(m.Data.BestBid, 64)
	ask, err2 := strconv.ParseFloat(m.Data.BestAsk, 64)
	bidSize, err3 := strconv.ParseFloat(m.Data.BuySize, 64)
	askSize, err4 := strconv.ParseFloat(m.Data.SellSize, 64)
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil || bid == 0 || ask == 0 {
		return
	}

	// дедупликация
	ws.mu.Lock()
	last := ws.last[sym]
	if last[0] == bid && last[1] == ask && last[2] == bidSize && last[3] == askSize {
		ws.mu.Unlock()
		return
	}
	ws.last[sym] = [4]float64{bid, ask, bidSize, askSize}
	ws.mu.Unlock()

	// получаем MarketData из пула
	md := c.pool.Get().(*models.MarketData)
	md.Exchange = "KuCoin"
	md.Symbol = sym
	md.Bid = bid
	md.Ask = ask
	md.BidSize = bidSize
	md.AskSize = askSize
	md.Timestamp = time.Now().UnixMilli()

	c.out <- md
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






package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"crypt_proto/collector"
	"crypt_proto/pkg/models"
)

func main() {
	// ------------------- Пул для MarketData -------------------
	pool := &sync.Pool{
		New: func() any {
			return &models.MarketData{}
		},
	}

	// ------------------- Канал для данных -------------------
	out := make(chan *models.MarketData, 1000) // буфер для скорости

	// ------------------- Создаём KuCoinCollector -------------------
	kc, err := collector.NewKuCoinCollectorFromCSV("pairs.csv", pool)
	if err != nil {
		log.Fatal("KuCoinCollector init error:", err)
	}

	// ------------------- Старт коллектора -------------------
	if err := kc.Start(out); err != nil {
		log.Fatal("KuCoinCollector start error:", err)
	}
	log.Println("[Main] KuCoinCollector started")

	// ------------------- Обработка данных -------------------
	go func() {
		for md := range out {
			// Для теста просто выводим
			log.Printf("[%s] %s | Bid: %f | Ask: %f\n",
				md.Exchange, md.Symbol, md.Bid, md.Ask)

			// Возвращаем объект в пул
			pool.Put(md)
		}
	}()

	// ------------------- Ожидание завершения -------------------
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	log.Println("[Main] stopping...")
	if err := kc.Stop(); err != nil {
		log.Println("KuCoinCollector stop error:", err)
	}
	close(out)
	log.Println("[Main] stopped")
}









