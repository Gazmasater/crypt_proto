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




go run -race main.go


GOMAXPROCS=8 go run -race main.go



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

	"crypt_proto/pkg/models"

	"github.com/gorilla/websocket"
	"github.com/tidwall/gjson"
)

const (
	maxSubsPerWS = 126
	subRate      = 120 * time.Millisecond
	pingInterval = 20 * time.Second
	reconnectDelay = 3 * time.Second
)

type KuCoinCollector struct {
	ctx    context.Context
	cancel context.CancelFunc
	wsList []*kucoinWS
	out    chan<- *models.MarketData
}

type Last struct {
	Bid float64
	Ask float64
}

type kucoinWS struct {
	id      int
	conn    *websocket.Conn
	symbols []string

	last map[string]Last
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
			id:      len(wsList),
			symbols: symbols[i:end],
			last:    make(map[string]Last),
		})
	}

	c := &KuCoinCollector{
		ctx:    ctx,
		cancel: cancel,
		wsList: wsList,
	}

	return c, symbols, nil
}

func (c *KuCoinCollector) Name() string { return "KuCoin" }

func (c *KuCoinCollector) Start(out chan<- *models.MarketData) error {
	c.out = out
	for _, ws := range c.wsList {
		go ws.run(c) // запускаем WS с авто-перезапуском
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

func (ws *kucoinWS) run(c *KuCoinCollector) {
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
		}

		// подключаемся
		if err := ws.connect(); err != nil {
			log.Printf("[KuCoin WS %d] connect error: %v, retry in %v\n", ws.id, err, reconnectDelay)
			time.Sleep(reconnectDelay)
			continue
		}

		// запускаем подписки и пинг
		done := make(chan struct{})
		go func() {
			ws.subscribeLoop()
			close(done)
		}()
		go ws.pingLoop()

		// читаем данные
		if err := ws.readLoop(c); err != nil {
			log.Printf("[KuCoin WS %d] readLoop ended: %v\n", ws.id, err)
		}

		// закрываем соединение и перезапускаем
		if ws.conn != nil {
			_ = ws.conn.Close()
			ws.conn = nil
		}

		log.Printf("[KuCoin WS %d] reconnecting in %v...\n", ws.id, reconnectDelay)
		time.Sleep(reconnectDelay)
	}
}

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
	}
}

func (ws *kucoinWS) pingLoop() {
	t := time.NewTicker(pingInterval)
	defer t.Stop()
	for range t.C {
		if ws.conn != nil {
			_ = ws.conn.WriteJSON(map[string]any{"id": time.Now().UnixNano(), "type": "ping"})
		}
	}
}

func (ws *kucoinWS) readLoop(c *KuCoinCollector) error {
	for {
		if ws.conn == nil {
			return fmt.Errorf("connection closed")
		}
		_, msg, err := ws.conn.ReadMessage()
		if err != nil {
			return err
		}
		ws.handle(c, msg)
	}
}

func (ws *kucoinWS) handle(c *KuCoinCollector, msg []byte) {
	const prefix = "/market/ticker:"
	const prefixLen = len(prefix)

	topicRes := gjson.GetBytes(msg, "topic")
	if !topicRes.Exists() {
		return
	}

	raw := topicRes.Raw // строка JSON: "/market/ticker:BTC-USDT"
	if len(raw) <= prefixLen+2 {
		return
	}

	if raw[1:1+prefixLen] != prefix { // пропускаем первую кавычку
		return
	}

	symbol := raw[1+prefixLen : len(raw)-1]

	data := gjson.GetBytes(msg, "data")
	bid := data.Get("bestBid").Float()
	ask := data.Get("bestAsk").Float()
	if bid == 0 || ask == 0 {
		return
	}

	if last, ok := ws.last[symbol]; ok && last.Bid == bid && last.Ask == ask {
		return
	}

	ws.last[symbol] = Last{Bid: bid, Ask: ask}

	c.out <- &models.MarketData{
		Exchange: "KuCoin",
		Symbol:   symbol,
		Bid:      bid,
		Ask:      ask,
		BidSize:  data.Get("bestBidSize").Float(),
		AskSize:  data.Get("bestAskSize").Float(),
	}
}

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



