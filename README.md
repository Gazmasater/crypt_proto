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
	"os"
	"strings"
	"sync"
	"time"

	"crypt_proto/pkg/models"

	"github.com/gorilla/websocket"
)

/* ================= STRUCT ================= */

type KuCoinCollector struct {
	ctx    context.Context
	cancel context.CancelFunc

	conn    *websocket.Conn
	wsURL   string
	symbols []string

	out chan<- *models.MarketData

	last map[string][2]float64
	mu   sync.Mutex

	ready bool
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

	return &KuCoinCollector{
		ctx:     ctx,
		cancel:  cancel,
		symbols: symbols,
		last:    make(map[string][2]float64),
	}, nil
}

/* ================= INTERFACE ================= */

func (c *KuCoinCollector) Name() string { return "KuCoin" }

func (c *KuCoinCollector) Start(out chan<- *models.MarketData) error {
	c.out = out

	if err := c.initWS(); err != nil {
		return err
	}

	go c.readLoop()
	go c.subscribeBatches(15, 400*time.Millisecond)

	log.Println("[KuCoin] started")
	return nil
}

func (c *KuCoinCollector) Stop() error {
	c.cancel()
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

/* ================= WS INIT ================= */

func (c *KuCoinCollector) initWS() error {
	req, _ := http.NewRequest(
		"POST",
		"https://api.kucoin.com/api/v1/bullet-public",
		nil,
	)

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

	c.wsURL = fmt.Sprintf(
		"%s?token=%s&connectId=%d",
		r.Data.InstanceServers[0].Endpoint,
		r.Data.Token,
		time.Now().UnixNano(),
	)

	conn, _, err := websocket.DefaultDialer.Dial(c.wsURL, nil)
	if err != nil {
		return err
	}

	// deadlines
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	c.conn = conn
	log.Println("[KuCoin] WS connected")
	return nil
}

/* ================= SUBSCRIBE ================= */

func (c *KuCoinCollector) subscribeBatches(batch int, delay time.Duration) {
	// ждём welcome
	for !c.ready {
		select {
		case <-c.ctx.Done():
			return
		default:
			time.Sleep(50 * time.Millisecond)
		}
	}

	log.Println("[KuCoin] subscribing symbols...")

	for i := 0; i < len(c.symbols); i += batch {
		end := i + batch
		if end > len(c.symbols) {
			end = len(c.symbols)
		}

		for _, s := range c.symbols[i:end] {
			topic := "/market/level2:" + s
			err := c.conn.WriteJSON(map[string]any{
				"id":       time.Now().UnixNano(),
				"type":     "subscribe",
				"topic":    topic,
				"response": true,
			})
			if err != nil {
				log.Println("[KuCoin] subscribe error:", err, topic)
			} else {
				log.Println("[KuCoin] subscribed:", s)
			}
		}

		time.Sleep(delay)
	}
}

/* ================= READ LOOP ================= */

func (c *KuCoinCollector) readLoop() {
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			_, msg, err := c.conn.ReadMessage()
			if err != nil {
				log.Println("[KuCoin] read error:", err)
				return
			}
			log.Println("[KuCoin] raw msg:", string(msg))
			c.handle(msg)
		}
	}
}

/* ================= HANDLE ================= */

func (c *KuCoinCollector) handle(msg []byte) {
	var raw map[string]any
	if err := json.Unmarshal(msg, &raw); err != nil {
		return
	}

	switch raw["type"] {
	case "welcome":
		c.ready = true
		log.Println("[KuCoin] welcome")
		return
	case "ack":
		return
	case "ping":
		_ = c.conn.WriteJSON(map[string]any{
			"id":   raw["id"],
			"type": "pong",
		})
		return
	case "message":
		// continue
	default:
		return
	}

	data, ok := raw["data"].(map[string]any)
	if !ok {
		return
	}

	// level2 data
	bids, _ := data["bids"].([]any)
	asks, _ := data["asks"].([]any)

	if len(bids) == 0 || len(asks) == 0 {
		return
	}

	bid := parseFloat(bids[0].([]any)[0])
	ask := parseFloat(asks[0].([]any)[0])

	sym, ok := data["symbol"].(string)
	if !ok {
		return
	}
	symbol := normalize(sym)

	c.mu.Lock()
	last := c.last[symbol]
	if last[0] == bid && last[1] == ask {
		c.mu.Unlock()
		return
	}
	c.last[symbol] = [2]float64{bid, ask}
	c.mu.Unlock()

	c.out <- &models.MarketData{
		Exchange: c.Name(),
		Symbol:   symbol,
		Bid:      bid,
		Ask:      ask,
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

/* ================= HELPERS ================= */

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



2026/01/06 09:57:55 [KuCoin] raw msg: {"topic":"/market/level2:DAG-USDT","type":"message","subject":"trade.l2update","data":{"changes":{"asks":[],"bids":[["0.014142","0","3428888568"]]},"sequenceEnd":3428888568,"sequenceStart":3428888568,"symbol":"DAG-USDT","time":1767682674903}}
2026/01/06 09:57:55 [KuCoin] raw msg: {"topic":"/market/level2:BCH-USDT","type":"message","subject":"trade.l2update","data":{"changes":{"asks":[],"bids":[["645.34","0.0154","4359956823"]]},"sequenceEnd":4359956823,"sequenceStart":4359956823,"symbol":"BCH-USDT","time":1767682674903}}
2026/01/06 09:57:55 [KuCoin] raw msg: {"topic":"/market/level2:LINK-BTC","type":"message","subject":"trade.l2update","data":{"changes":{"asks":[],"bids":[["0.00014805","0","2339568791"]]},"sequenceEnd":2339568791,"sequenceStart":2339568791,"symbol":"LINK-BTC","time":1767682674903}}
2026/01/06 09:57:55 [KuCoin] raw msg: {"topic":"/market/level2:LINK-BTC","type":"message","subject":"trade.l2update","data":{"changes":{"asks":[["0.00014868","6.2437","2339568792"]],"bids":[]},"sequenceEnd":2339568792,"sequenceStart":2339568792,"symbol":"LINK-BTC","time":1767682674904}}
2026/01/06 09:57:55 [KuCoin] raw msg: {"id":"1767682674776323797","type":"ack"}
2026/01/06 09:57:55 [KuCoin] raw msg: {"id":"1767682674776409489","type":"ack"}
2026/01/06 09:57:55 [KuCoin] raw msg: {"id":"1767682674776423020","type":"ack"}
2026/01/06 09:57:55 [KuCoin] raw msg: {"id":"1767682674776442132","type":"ack"}
2026/01/06 09:57:55 [KuCoin] raw msg: {"id":"1767682674776458933","type":"ack"}
2026/01/06 09:57:55 [KuCoin] raw msg: {"id":"1767682674776469443","type":"ack"}
2026/01/06 09:57:55 [KuCoin] raw msg: {"topic":"/market/level2:BCH-USDT","type":"message","subject":"trade.l2update","data":{"changes":{"asks":[["646.67","0.0308","4359956824"]],"bids":[]},"sequenceEnd":4359956824,"sequenceStart":4359956824,"symbol":"BCH-USDT","time":1767682674906}}
2026/01/06 09:57:55 [KuCoin] raw msg: {"topic":"/market/level2:BNB-USDT","type":"message","subject":"trade.l2update","data":{"changes":{"asks":[],"bids":[["908.809","1.2841","10489462323"]]},"sequenceEnd":10489462323,"sequenceStart":10489462323,"symbol":"BNB-USDT","time":1767682674906}}
2026/01/06 09:57:55 [KuCoin] raw msg: {"topic":"/market/level2:AAVE-USDT","type":"message","subject":"trade.l2update","data":{"changes":{"asks":[["173.877","0.8097","8387547760"]],"bids":[]},"sequenceEnd":8387547760,"sequenceStart":8387547760,"symbol":"AAVE-USDT","time":1767682674906}}
2026/01/06 09:57:55 [KuCoin] raw msg: {"id":"1767682674776642917","type":"ack"}
2026/01/06 09:57:55 [KuCoin] raw msg: {"id":"1767682674776656361","type":"ack"}
2026/01/06 09:57:55 [KuCoin] raw msg: {"id":"1767682674776667111","type":"ack"}
2026/01/06 09:57:55 [KuCoin] raw msg: {"id":"1767682674776677054","type":"ack"}
2026/01/06 09:57:55 [KuCoin] raw msg: {"id":"695cb272a02344126eb870a6","type":"error","code":509,"data":"exceed max permits per second"}
2026/01/06 09:57:55 [KuCoin] read error: websocket: close 1000 (normal): Bye
2026/01/06 09:57:55 [KuCoin] subscribe error: websocket: close sent /market/level2:ETH-BTC
2026/01/06 09:57:55 [KuCoin] subscribe error: websocket: close sent /market/level2:STORJ-ETH
2026/01/06 09:57:55 [KuCoin] subscribe error: websocket: close sent /market/level2:BNB-BTC
2026/01/06 09:57:55 [KuCoin] subscribe error: websocket: close sent /market/level2:EWT-USDT
2026/01/06 09:57:55 [KuCoin] subscribe error: websocket: close sent /market/level2:WIN-TRX
2026/01/06 09:57:55 [KuCoin] subscribe error: websocket: close sent /market/level2:VSYS-USDT
2026/01/06 09:57:55 [KuCoin] subscribe error: websocket: close sent /market/level2:RUNE-BTC
2026/01/06 09:57:55 [KuCoin] subscribe error: websocket: close sent /market/level2:IOTX-ETH
2026/01/06 09:57:55 [KuCoin] subscribe error: websocket: close sent /market/level2:TEL-ETH
2026/01/06 09:57:55 [KuCoin] subscribe error: websocket: close sent /market/level2:ANKR-BTC
2026/01/06 09:57:55 [KuCoin] subscribe error: websocket: close sent /market/level2:IOTA-USDT
2026/01/06 09:57:55 [KuCoin] subscribe error: websocket: close sent /market/level2:NEAR-USDT
2026/01/06 09:57:55 [KuCoin] subscribe error: websocket: close sent /market/level2:AAVE-BTC
2026/01/06 09:57:55 [KuCoin] subscribe error: websocket: close sent /market/level2:STX-USDT
2026/01/06 09:57:55 [KuCoin] subscribe error: websocket: close sent /market/level2:FET-BTC
2026/01/06 09:57:55 [KuCoin] subscribe error: websocket: close sent /market/level2:ALGO-USDT
2026/01/06 09:57:55 [KuCoin] subscribe error: websocket: close sent /market/level2:VET-BTC
2026/01/06 09:57:55 [KuCoin] subscribe error: websocket: close sent /market/level2:REQ-BTC
2026/01/06 09:57:55 [KuCoin] subscribe error: websocket: close sent /market/level2:HBAR-USDT
2026/01/06 09:57:55 [KuCoin] subscribe error: websocket: close sent /market/level2:KRL-USDT
2026/01/06 09:57:55 [KuCoin] subscribe error: websocket: close sent /market/level2:CHZ-BTC
2026/01/06 09:57:55 [KuCoin] subscribe error: websocket: close sent /market/level2:CHZ-USDT
2026/01/06 09:57:55 [KuCoin] subscribe error: websocket: close sent /market/level2:KCS-ETH
2026/01/06 09:57:55 [KuCoin] subscribe error: websocket: close sent /market/level2:ZIL-USDT
2026/01/06 09:57:55 [KuCoin] subscribe error: websocket: close sent /market/level2:DAG-ETH
2026/01/06 09:57:55 [KuCoin] subscribe error: websocket: close sent /market/level2:NEO-BTC
2026/01/06 09:57:55 [KuCoin] subscribe error: websocket: close sent /market/level2:WBTC-USDT
2026/01/06 09:57:55 [KuCoin] subscribe error: websocket: close sent /market/level2:KCS-BTC
2026/01/06 09:57:55 [KuCoin] subscribe error: websocket: close sent /market/level2:CFX-USDT
2026/01/06 09:57:55 [KuCoin] subscribe error: websocket: close sent /market/level2:ADA-KCS
2026/01/06 09:57:55 [KuCoin] subscribe error: websocket: close sent /market/level2:ENJ-USDT
2026/01/06 09:57:55 [KuCoin] subscribe error: websocket: close sent /market/level2:ONT-BTC
2026/01/06 09:57:55 [KuCoin] subscribe error: websocket: close sent /market/level2:DYP-USDT
2026/01/06 09:57:55 [KuCoin] subscribe error: websocket: close sent /market/level2:AVAX-BTC
2026/01/06 09:57:55 [KuCoin] subscribe error: websocket: close sent /market/level2:KRL-BTC
2026/01/06 09:57:55 [KuCoin] subscribe error: websocket: close sent /market/level2:DASH-BTC
2026/01/06 09:57:55 [KuCoin] subscribe error: websocket: close sent /market/level2:XCN-USDT
2026/01/06 09:57:55 [KuCoin] subscribe error: websocket: close sent /market/level2:ETC-ETH
2026/01/06 09:57:55 [KuCoin] subscribe error: websocket: close sent /market/level2:ICP-USDT
2026/01/06 09:57:55 [KuCoin] subscribe error: websocket: close sent /market/level2:AVAX-USDT
2026/01/06 09:57:55 [KuCoin] subscribe error: websocket: close sent /market/level2:XLM-BTC
2026/01/06 09:57:55 [KuCoin] subscribe error: websocket: close sent /market/level2:ELA-USDT
2026/01/06 09:57:55 [KuCoin] subscribe error: websocket: close sent /market/level2:TRX-USDT
2026/01/06 09:57:55 [KuCoin] subscribe error: websocket: close sent /market/level2:AVA-USDT
2026/01/06 09:57:55 [KuCoin] subscribe error: websocket: close sent /market/level2:XYO-USDT
2026/01/06 09:57:56 [KuCoin] subscribe error: websocket: close sent /market/level2:XLM-ETH
2026/01/06 09:57:56 [KuCoin] subscribe error: websocket: close sent /market/level2:SNX-BTC
2026/01/06 09:57:56 [KuCoin] subscribe error: websocket: close sent /market/level2:NEO-USDT
2026/01/06 09:57:56 [KuCoin] subscribe error: websocket: close sent /market/level2:MOVR-ETH
2026/01/06 09:57:56 [KuCoin] subscribe error: websocket: close sent /market/level2:ENJ-ETH
2026/01/06 09:57:56 [KuCoin] subscribe error: websocket: close sent /market/level2:KLV-BTC
2026/01/06 09:57:56 [KuCoin] subscribe error: websocket: close sent /market/level2:RUNE-USDT
2026/01/06 09:57:56 [KuCoin] subscribe error: websocket: close sent /market/level2:IOST-ETH
2026/01/06 09:57:56 [KuCoin] subscribe error: websocket: close sent /market/level2:DGB-ETH
2026/01/06 09:57:56 [KuCoin] subscribe error: websocket: close sent /market/level2:SUPER-BTC
2026/01/06 09:57:56 [KuCoin] subscribe error: websocket: close sent /market/level2:WAVES-USDT
2026/01/06 09:57:56 [KuCoin] subscribe error: websocket: close sent /market/level2:ONT-USDT
2026/01/06 09:57:56 [KuCoin] subscribe error: websocket: close sent /market/level2:KAS-BTC
2026/01/06 09:57:56 [KuCoin] subscribe error: websocket: close sent /market/level2:TRVL-BTC
2026/01/06 09:57:56 [KuCoin] subscribe error: websocket: close sent /market/level2:SOL-KCS
2026/01/06 09:57:56 [KuCoin] subscribe error: websocket: close sent /market/level2:IOST-USDT
2026/01/06 09:57:56 [KuCoin] subscribe error: websocket: close sent /market/level2:STORJ-USDT
2026/01/06 09:57:56 [KuCoin] subscribe error: websocket: close sent /market/level2:ICP-BTC
2026/01/06 09:57:56 [KuCoin] subscribe error: websocket: close sent /market/level2:BCH-BTC
2026/01/06 09:57:56 [KuCoin] subscribe error: websocket: close sent /market/level2:XMR-BTC
2026/01/06 09:57:56 [KuCoin] subscribe error: websocket: close sent /market/level2:ATOM-USDT
2026/01/06 09:57:56 [KuCoin] subscribe error: websocket: close sent /market/level2:BCHSV-BTC
2026/01/06 09:57:56 [KuCoin] subscribe error: websocket: close sent /market/level2:BAX-BTC
2026/01/06 09:57:56 [KuCoin] subscribe error: websocket: close sent /market/level2:TEL-BTC
2026/01/06 09:57:56 [KuCoin] subscribe error: websocket: close sent /market/level2:PERP-BTC
2026/01/06 09:57:56 [KuCoin] subscribe error: websocket: close sent /market/level2:TEL-USDT
2026/01/06 09:57:56 [KuCoin] subscribe error: websocket: close sent /market/level2:XMR-USDT
2026/01/06 09:57:56 [KuCoin] subscribe error: websocket: close sent /market/level2:ANKR-USDT
2026/01/06 09:57:56 [KuCoin] subscribe error: websocket: close sent /market/level2:IOTX-BTC
2026/01/06 09:57:56 [KuCoin] subscribe error: websocket: close sent /market/level2:BTC-USDT
2026/01/06 09:57:57 [KuCoin] subscribe error: websocket: close sent /market/level2:CKB-BTC
2026/01/06 09:57:57 [KuCoin] subscribe error: websocket: close sent /market/level2:WBTC-BTC
2026/01/06 09:57:57 [KuCoin] subscribe error: websocket: close sent /market/level2:DOGE-BTC
2026/01/06 09:57:57 [KuCoin] subscribe error: websocket: close sent /market/level2:KCS-USDT
2026/01/06 09:57:57 [KuCoin] subscribe error: websocket: close sent /market/level2:HYPE-KCS
2026/01/06 09:57:57 [KuCoin] subscribe error: websocket: close sent /market/level2:XDC-USDT
2026/01/06 09:57:57 [KuCoin] subscribe error: websocket: close sent /market/level2:ALGO-BTC
2026/01/06 09:57:57 [KuCoin] subscribe error: websocket: close sent /market/level2:AVA-ETH
2026/01/06 09:57:57 [KuCoin] subscribe error: websocket: close sent /market/level2:SUPER-USDT
2026/01/06 09:57:57 [KuCoin] subscribe error: websocket: close sent /market/level2:SCRT-USDT
2026/01/06 09:57:57 [KuCoin] subscribe error: websocket: close sent /market/level2:VRA-BTC
2026/01/06 09:57:57 [KuCoin] subscribe error: websocket: close sent /market/level2:KNC-BTC
2026/01/06 09:57:57 [KuCoin] subscribe error: websocket: close sent /market/level2:A-USDT
2026/01/06 09:57:57 [KuCoin] subscribe error: websocket: close sent /market/level2:WAN-BTC
2026/01/06 09:57:57 [KuCoin] subscribe error: websocket: close sent /market/level2:WIN-BTC
2026/01/06 09:57:57 [KuCoin] subscribe error: websocket: close sent /market/level2:RLC-BTC
2026/01/06 09:57:57 [KuCoin] subscribe error: websocket: close sent /market/level2:EWT-BTC
2026/01/06 09:57:57 [KuCoin] subscribe error: websocket: close sent /market/level2:CRO-BTC
2026/01/06 09:57:57 [KuCoin] subscribe error: websocket: close sent /market/level2:AR-USDT
2026/01/06 09:57:57 [KuCoin] subscribe error: websocket: close sent /market/level2:RSR-BTC
2026/01/06 09:57:57 [KuCoin] subscribe error: websocket: close sent /market/level2:COTI-USDT
2026/01/06 09:57:57 [KuCoin] subscribe error: websocket: close sent /market/level2:DOT-KCS
2026/01/06 09:57:57 [KuCoin] subscribe error: websocket: close sent /market/level2:ETH-DAI
2026/01/06 09:57:57 [KuCoin] subscribe error: websocket: close sent /market/level2:ERG-BTC
2026/01/06 09:57:57 [KuCoin] subscribe error: websocket: close sent /market/level2:GAS-BTC
2026/01/06 09:57:57 [KuCoin] subscribe error: websocket: close sent /market/level2:TRAC-ETH
2026/01/06 09:57:57 [KuCoin] subscribe error: websocket: close sent /market/level2:SNX-ETH
2026/01/06 09:57:57 [KuCoin] subscribe error: websocket: close sent /market/level2:ETH-BRL
2026/01/06 09:57:57 [KuCoin] subscribe error: websocket: close sent /market/level2:BDX-USDT
2026/01/06 09:57:57 [KuCoin] subscribe error: websocket: close sent /market/level2:LYX-ETH
2026/01/06 09:57:57 [KuCoin] subscribe error: websocket: close sent /market/level2:MANA-ETH
2026/01/06 09:57:57 [KuCoin] subscribe error: websocket: close sent /market/level2:MANA-USDT
2026/01/06 09:57:57 [KuCoin] subscribe error: websocket: close sent /market/level2:REQ-USDT
2026/01/06 09:57:57 [KuCoin] subscribe error: websocket: close sent /market/level2:DAG-BTC
2026/01/06 09:57:57 [KuCoin] subscribe error: websocket: close sent /market/level2:CFX-BTC
2026/01/06 09:57:57 [KuCoin] subscribe error: websocket: close sent /market/level2:XTZ-USDT
2026/01/06 09:57:57 [KuCoin] subscribe error: websocket: close sent /market/level2:TRAC-USDT
2026/01/06 09:57:57 [KuCoin] subscribe error: websocket: close sent /market/level2:IOTA-BTC
2026/01/06 09:57:57 [KuCoin] subscribe error: websocket: close sent /market/level2:BTC-BRL
2026/01/06 09:57:57 [KuCoin] subscribe error: websocket: close sent /market/level2:COTI-BTC
2026/01/06 09:57:57 [KuCoin] subscribe error: websocket: close sent /market/level2:SHIB-DOGE
2026/01/06 09:57:57 [KuCoin] subscribe error: websocket: close sent /market/level2:ELA-BTC
2026/01/06 09:57:57 [KuCoin] subscribe error: websocket: close sent /market/level2:XCN-BTC
2026/01/06 09:57:57 [KuCoin] subscribe error: websocket: close sent /market/level2:MOVR-USDT
2026/01/06 09:57:57 [KuCoin] subscribe error: websocket: close sent /market/level2:WIN-USDT
2026/01/06 09:57:58 [KuCoin] subscribe error: websocket: close sent /market/level2:VSYS-BTC
2026/01/06 09:57:58 [KuCoin] subscribe error: websocket: close sent /market/level2:CRO-USDT
2026/01/06 09:57:58 [KuCoin] subscribe error: websocket: close sent /market/level2:ETH-USDT
2026/01/06 09:57:58 [KuCoin] subscribe error: websocket: close sent /market/level2:SNX-USDT
2026/01/06 09:57:58 [KuCoin] subscribe error: websocket: close sent /market/level2:SXP-USDT
2026/01/06 09:57:58 [KuCoin] subscribe error: websocket: close sent /market/level2:ASTR-BTC
2026/01/06 09:57:58 [KuCoin] subscribe error: websocket: close sent /market/level2:PEPE-KCS
2026/01/06 09:57:58 [KuCoin] subscribe error: websocket: close sent /market/level2:NEAR-BTC
2026/01/06 09:57:58 [KuCoin] subscribe error: websocket: close sent /market/level2:SHIB-USDT
2026/01/06 09:57:58 [KuCoin] subscribe error: websocket: close sent /market/level2:BCHSV-ETH
2026/01/06 09:57:58 [KuCoin] subscribe error: websocket: close sent /market/level2:SKL-USDT
2026/01/06 09:57:58 [KuCoin] subscribe error: websocket: close sent /market/level2:OM-BTC
2026/01/06 09:57:58 [KuCoin] subscribe error: websocket: close sent /market/level2:TRX-ETH
2026/01/06 09:57:58 [KuCoin] subscribe error: websocket: close sent /market/level2:VRA-USDT
2026/01/06 09:57:58 [KuCoin] subscribe error: websocket: close sent /market/level2:ETC-USDT
2026/01/06 09:57:58 [KuCoin] subscribe error: websocket: close sent /market/level2:OGN-USDT
2026/01/06 09:57:58 [KuCoin] subscribe error: websocket: close sent /market/level2:KNC-ETH
2026/01/06 09:57:58 [KuCoin] subscribe error: websocket: close sent /market/level2:XRP-BTC
2026/01/06 09:57:58 [KuCoin] subscribe error: websocket: close sent /market/level2:XRP-USDT
2026/01/06 09:57:58 [KuCoin] subscribe error: websocket: close sent /market/level2:A-ETH
2026/01/06 09:57:58 [KuCoin] subscribe error: websocket: close sent /market/level2:XYO-ETH



