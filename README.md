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
	retryDelay   = 5 * time.Second
)

/* ================= POOL ================= */
type KuCoinCollector struct {
	ctx    context.Context
	cancel context.CancelFunc
	wsList []*kucoinWS
	out    chan<- *models.MarketData
}

type kucoinWS struct {
	id      int
	conn    *websocket.Conn
	symbols []string
	last    map[string][2]float64
	mu      sync.Mutex
	ctx     context.Context
	cancel  context.CancelFunc
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

	ctx, cancel := context.WithCancel(context.Background())
	var wsList []*kucoinWS

	for i := 0; i < len(symbols); i += maxSubsPerWS {
		end := i + maxSubsPerWS
		if end > len(symbols) {
			end = len(symbols)
		}

		wsCtx, wsCancel := context.WithCancel(ctx)
		wsList = append(wsList, &kucoinWS{
			id:      len(wsList),
			symbols: symbols[i:end],
			last:    make(map[string][2]float64),
			ctx:     wsCtx,
			cancel:  wsCancel,
		})
	}

	return &KuCoinCollector{
		ctx:    ctx,
		cancel: cancel,
		wsList: wsList,
	}, nil
}

/* ================= INTERFACE ================= */
func (c *KuCoinCollector) Name() string {
	return "KuCoin"
}

func (c *KuCoinCollector) Start(out chan<- *models.MarketData) error {
	c.out = out
	for _, ws := range c.wsList {
		go ws.run(c)
	}
	log.Printf("[KuCoin] started with %d WS\n", len(c.wsList))
	return nil
}

func (c *KuCoinCollector) Stop() error {
	c.cancel()
	for _, ws := range c.wsList {
		ws.cancel()
		if ws.conn != nil {
			ws.conn.Close()
		}
	}
	return nil
}

/* ================= WS RUN ================= */
func (ws *kucoinWS) run(c *KuCoinCollector) {
	for {
		select {
		case <-ws.ctx.Done():
			return
		default:
			if err := ws.connect(); err != nil {
				log.Printf("[KuCoin WS %d] connect failed: %v, retrying in %v", ws.id, err, retryDelay)
				time.Sleep(retryDelay)
				continue
			}
			log.Printf("[KuCoin WS %d] connected", ws.id)

			wg := &sync.WaitGroup{}
			wg.Add(3)
			go ws.readLoop(c, wg)
			go ws.pingLoop(wg)
			go ws.subscribeLoop(wg)

			wg.Wait() // все горутины закончили работу (обычно из-за ошибки)
			log.Printf("[KuCoin WS %d] disconnected, reconnecting...", ws.id)
			if ws.conn != nil {
				ws.conn.Close()
			}
			time.Sleep(retryDelay)
		}
	}
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
	return nil
}

/* ================= SUBSCRIBE ================= */
func (ws *kucoinWS) subscribeLoop(wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(1 * time.Second)
	t := time.NewTicker(subRate)
	defer t.Stop()

	for _, s := range ws.symbols {
		select {
		case <-ws.ctx.Done():
			return
		case <-t.C:
			topic := "/market/ticker:" + s
			err := ws.conn.WriteJSON(map[string]any{
				"id":      time.Now().UnixNano(),
				"type":    "subscribe",
				"topic":   topic,
				"response": true,
			})
			if err != nil {
				log.Printf("[KuCoin WS %d] subscribe error %s: %v", ws.id, s, err)
			}
		}
	}
}

/* ================= PING ================= */
func (ws *kucoinWS) pingLoop(wg *sync.WaitGroup) {
	defer wg.Done()
	t := time.NewTicker(pingInterval)
	defer t.Stop()

	for {
		select {
		case <-ws.ctx.Done():
			return
		case <-t.C:
			_ = ws.conn.WriteJSON(map[string]any{
				"id":   time.Now().UnixNano(),
				"type": "ping",
			})
		}
	}
}

/* ================= READ ================= */
func (ws *kucoinWS) readLoop(c *KuCoinCollector, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ws.ctx.Done():
			return
		default:
			_, msg, err := ws.conn.ReadMessage()
			if err != nil {
				return
			}
			ws.handle(c, msg)
		}
	}
}

/* ================= HANDLE ================= */
func (ws *kucoinWS) handle(c *KuCoinCollector, msg []byte) {
	var raw map[string]any
	if json.Unmarshal(msg, &raw) != nil {
		return
	}
	if raw["type"] != "message" {
		return
	}

	topic := raw["topic"].(string)
	parts := strings.Split(topic, ":")
	if len(parts) != 2 {
		return
	}

	symbol := normalize(parts[1])
	data := raw["data"].(map[string]any)
	bid := parseFloat(data["bestBid"])
	ask := parseFloat(data["bestAsk"])
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

	select {
	case c.out <- &models.MarketData{
		Exchange: "KuCoin",
		Symbol:   symbol,
		Bid:      bid,
		Ask:      ask,
	}:
	default:
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
	p := strings.Split(s, "-")
	return p[0] + "/" + p[1]
}

func parseFloat(v any) float64 {
	switch t := v.(type) {
	case string:
		f, _ := strconv.ParseFloat(t, 64)
		return f
	case float64:
		return t
	default:
		return 0
	}
}







