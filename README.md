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

type KuCoinCollector struct {
	ctx     context.Context
	cancel  context.CancelFunc
	conn    *websocket.Conn
	wsURL   string
	symbols []string
	out     chan<- *models.MarketData
	last    map[string][2]float64
	mu      sync.Mutex
	ready   bool
}

// ---------------- CONSTRUCTOR ----------------

func NewKuCoinCollectorFromCSV(path string) (*KuCoinCollector, error) {
	symbols, err := parseCSVSymbols(path)
	if err != nil {
		return nil, err
	}
	if len(symbols) == 0 {
		return nil, fmt.Errorf("no valid symbols")
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &KuCoinCollector{
		ctx:     ctx,
		cancel:  cancel,
		symbols: symbols,
		last:    make(map[string][2]float64),
	}, nil
}

// ---------------- PARSE CSV ----------------

func parseCSVSymbols(path string) ([]string, error) {
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
			symbol := extractSymbol(row[i])
			if symbol != "" {
				set[symbol] = struct{}{}
			}
		}
	}

	var res []string
	for k := range set {
		res = append(res, k)
	}
	return res, nil
}

// берем только валютную пару, убираем BUY/SELL
func extractSymbol(s string) string {
	parts := strings.Fields(strings.ToUpper(strings.TrimSpace(s)))
	if len(parts) < 2 {
		return ""
	}
	sym := parts[1]           // "LINK/USDT"
	return strings.ReplaceAll(sym, "/", "-") // "LINK-USDT"
}

// ---------------- INTERFACE ----------------

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

// ---------------- WS INIT ----------------

func (c *KuCoinCollector) initWS() error {
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

	c.wsURL = fmt.Sprintf("%s?token=%s&connectId=%d",
		r.Data.InstanceServers[0].Endpoint,
		r.Data.Token,
		time.Now().UnixNano(),
	)

	conn, _, err := websocket.DefaultDialer.Dial(c.wsURL, nil)
	if err != nil {
		return err
	}

	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	c.conn = conn
	log.Println("[KuCoin] WS connected")
	return nil
}

// ---------------- SUBSCRIBE ----------------

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
			_ = c.conn.WriteJSON(map[string]any{
				"id":       time.Now().UnixNano(),
				"type":     "subscribe",
				"topic":    "/market/ticker:" + s,
				"response": true,
			})
			log.Println("[KuCoin] subscribed:", s)
		}

		time.Sleep(delay)
	}
}

// ---------------- READ LOOP ----------------

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
			c.handle(msg)
		}
	}
}

// ---------------- HANDLE ----------------

func (c *KuCoinCollector) handle(msg []byte) {
	var raw map[string]any
	if err := json.Unmarshal(msg, &raw); err != nil {
		log.Println("[KuCoin] raw parse error:", err)
		log.Println("[KuCoin] raw msg:", string(msg))
		return
	}

	switch raw["type"] {
	case "welcome":
		c.ready = true
		log.Println("[KuCoin] welcome")
	case "ack":
	case "ping":
		_ = c.conn.WriteJSON(map[string]any{"id": raw["id"], "type": "pong"})
	case "message":
		data, ok := raw["data"].(map[string]any)
		if !ok {
			return
		}

		bid := parseFloat(data["bestBid"])
		ask := parseFloat(data["bestAsk"])
		if bid == 0 || ask == 0 {
			return
		}

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
	default:
		return
	}
}

// ---------------- HELPERS ----------------

func normalize(s string) string {
	return strings.ReplaceAll(s, "-", "/")
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




2026/01/05 11:39:38 [KuCoin] WS connected
2026/01/05 11:39:38 [KuCoin] started
2026/01/05 11:39:38 [Main] KuCoinCollector started. Listening for data...
2026/01/05 11:39:38 [KuCoin] welcome
2026/01/05 11:39:38 [KuCoin] subscribing symbols...
2026/01/05 11:39:38 [KuCoin] subscribed: WIN-USDT
2026/01/05 11:39:38 [KuCoin] subscribed: BCHSV-BTC
2026/01/05 11:39:38 [KuCoin] subscribed: MANA-ETH
2026/01/05 11:39:38 [KuCoin] subscribed: ENJ-USDT
2026/01/05 11:39:38 [KuCoin] subscribed: ERG-USDT
2026/01/05 11:39:38 [KuCoin] subscribed: STX-USDT
2026/01/05 11:39:38 [KuCoin] subscribed: RSR-BTC
2026/01/05 11:39:38 [KuCoin] subscribed: XRP-KCS
2026/01/05 11:39:38 [KuCoin] subscribed: XYO-ETH
2026/01/05 11:39:38 [KuCoin] subscribed: BAX-ETH
2026/01/05 11:39:38 [KuCoin] subscribed: XLM-USDT
2026/01/05 11:39:38 [KuCoin] subscribed: EGLD-USDT
2026/01/05 11:39:38 [KuCoin] subscribed: PAXG-BTC
2026/01/05 11:39:38 [KuCoin] subscribed: AAVE-USDT
2026/01/05 11:39:38 [KuCoin] subscribed: PEPE-KCS
2026/01/05 11:39:38 [KuCoin] subscribed: NFT-TRX
2026/01/05 11:39:38 [KuCoin] subscribed: A-USDT
2026/01/05 11:39:38 [KuCoin] subscribed: REQ-USDT
2026/01/05 11:39:38 [KuCoin] subscribed: ANKR-BTC
2026/01/05 11:39:38 [KuCoin] subscribed: VET-USDT
2026/01/05 11:39:38 [KuCoin] subscribed: KLV-USDT
2026/01/05 11:39:38 [KuCoin] subscribed: HYPE-KCS
2026/01/05 11:39:38 [KuCoin] subscribed: ICX-USDT
2026/01/05 11:39:38 [KuCoin] subscribed: LTC-KCS
2026/01/05 11:39:38 [KuCoin] subscribed: KLV-BTC
2026/01/05 11:39:38 [KuCoin] subscribed: ATOM-ETH
2026/01/05 11:39:38 [KuCoin] subscribed: VET-BTC
2026/01/05 11:39:38 [KuCoin] subscribed: ADA-USDT
2026/01/05 11:39:38 [KuCoin] subscribed: OM-USDT
2026/01/05 11:39:38 [KuCoin] subscribed: LYX-USDT
2026/01/05 11:39:39 [KuCoin] subscribed: WAVES-BTC
2026/01/05 11:39:39 [KuCoin] subscribed: TRX-USDT
2026/01/05 11:39:39 [KuCoin] subscribed: RSR-USDT
2026/01/05 11:39:39 [KuCoin] subscribed: BNB-KCS
2026/01/05 11:39:39 [KuCoin] subscribed: PEPE-USDT
2026/01/05 11:39:39 [KuCoin] subscribed: VRA-BTC
2026/01/05 11:39:39 [KuCoin] subscribed: ERG-BTC
2026/01/05 11:39:39 [KuCoin] subscribed: KCS-BTC
2026/01/05 11:39:39 [KuCoin] subscribed: EWT-USDT
2026/01/05 11:39:39 [KuCoin] subscribed: GAS-USDT
2026/01/05 11:39:39 [KuCoin] subscribed: STORJ-USDT
2026/01/05 11:39:39 [KuCoin] subscribed: WBTC-BTC
2026/01/05 11:39:39 [KuCoin] subscribed: SHIB-USDT
2026/01/05 11:39:39 [KuCoin] subscribed: KCS-USDT
2026/01/05 11:39:39 [KuCoin] subscribed: DAG-USDT
2026/01/05 11:39:39 [KuCoin] subscribed: VSYS-BTC
2026/01/05 11:39:39 [KuCoin] subscribed: RUNE-BTC
2026/01/05 11:39:39 [KuCoin] subscribed: XRP-USDT
2026/01/05 11:39:39 [KuCoin] subscribed: ASTR-USDT
2026/01/05 11:39:39 [KuCoin] subscribed: CHZ-BTC
2026/01/05 11:39:39 [KuCoin] subscribed: HYPE-USDT
2026/01/05 11:39:39 [KuCoin] subscribed: KNC-USDT
2026/01/05 11:39:39 [KuCoin] subscribed: CFX-USDT
2026/01/05 11:39:39 [KuCoin] subscribed: TRVL-USDT
2026/01/05 11:39:39 [KuCoin] subscribed: SNX-ETH
2026/01/05 11:39:39 [KuCoin] subscribed: LINK-USDT
2026/01/05 11:39:39 [KuCoin] subscribed: DGB-BTC
2026/01/05 11:39:39 [KuCoin] subscribed: USDT-EUR
2026/01/05 11:39:39 [KuCoin] subscribed: CHZ-USDT
2026/01/05 11:39:39 [KuCoin] subscribed: IOST-USDT
2026/01/05 11:39:40 [KuCoin] subscribed: ZIL-USDT
2026/01/05 11:39:40 [KuCoin] subscribed: BDX-USDT
2026/01/05 11:39:40 [KuCoin] subscribed: SCRT-USDT
2026/01/05 11:39:40 [KuCoin] subscribed: KNC-ETH
2026/01/05 11:39:40 [KuCoin] subscribed: ALGO-ETH
2026/01/05 11:39:40 [KuCoin] subscribed: DGB-USDT
2026/01/05 11:39:40 [KuCoin] subscribed: CFX-BTC
2026/01/05 11:39:40 [KuCoin] subscribed: BTC-EUR
2026/01/05 11:39:40 [KuCoin] subscribed: TRX-BTC
2026/01/05 11:39:40 [KuCoin] subscribed: XCN-BTC
2026/01/05 11:39:40 [KuCoin] subscribed: WAN-BTC
2026/01/05 11:39:40 [KuCoin] subscribed: XYO-BTC
2026/01/05 11:39:40 [KuCoin] subscribed: DOT-USDT
2026/01/05 11:39:40 [KuCoin] subscribed: A-BTC
2026/01/05 11:39:40 [KuCoin] subscribed: AVA-BTC
2026/01/05 11:39:40 [KuCoin] subscribed: ATOM-USDT
2026/01/05 11:39:40 [KuCoin] subscribed: WBTC-USDT
2026/01/05 11:39:40 [KuCoin] subscribed: NEAR-USDT
2026/01/05 11:39:40 [KuCoin] subscribed: CKB-BTC
2026/01/05 11:39:40 [KuCoin] subscribed: ONE-BTC
2026/01/05 11:39:40 [KuCoin] subscribed: EGLD-BTC
2026/01/05 11:39:40 [KuCoin] subscribed: XDC-BTC
2026/01/05 11:39:40 [KuCoin] subscribed: DGB-ETH
2026/01/05 11:39:40 [KuCoin] subscribed: AVAX-BTC
2026/01/05 11:39:40 [KuCoin] subscribed: IOTX-USDT
2026/01/05 11:39:40 [KuCoin] subscribed: PERP-USDT
2026/01/05 11:39:40 [KuCoin] subscribed: DOT-KCS
2026/01/05 11:39:40 [KuCoin] subscribed: INJ-USDT
2026/01/05 11:39:40 [KuCoin] subscribed: LTC-USDT
2026/01/05 11:39:40 [KuCoin] subscribed: SUI-USDT
2026/01/05 11:39:40 [KuCoin] subscribed: DOGE-USDT
2026/01/05 11:39:40 [KuCoin] subscribed: ELA-USDT
2026/01/05 11:39:40 [KuCoin] subscribed: OM-BTC
2026/01/05 11:39:40 [KuCoin] subscribed: ONT-ETH
2026/01/05 11:39:40 [KuCoin] subscribed: ETC-BTC
2026/01/05 11:39:40 [KuCoin] subscribed: IOTA-BTC
2026/01/05 11:39:40 [KuCoin] subscribed: SXP-BTC
2026/01/05 11:39:40 [KuCoin] subscribed: NEO-BTC
2026/01/05 11:39:40 [KuCoin] subscribed: XDC-USDT
2026/01/05 11:39:40 [KuCoin] subscribed: LTC-ETH
2026/01/05 11:39:40 [KuCoin] subscribed: TRAC-USDT
2026/01/05 11:39:40 [KuCoin] subscribed: WAVES-USDT
2026/01/05 11:39:40 [KuCoin] subscribed: ASTR-BTC
2026/01/05 11:39:40 [KuCoin] subscribed: MANA-USDT
2026/01/05 11:39:40 [KuCoin] subscribed: NKN-USDT
2026/01/05 11:39:41 [KuCoin] read error: websocket: close 1000 (normal): Bye
2026/01/05 11:39:41 [KuCoin] subscribed: TRVL-BTC
2026/01/05 11:39:41 [KuCoin] subscribed: RLC-BTC
2026/01/05 11:39:41 [KuCoin] subscribed: ALGO-BTC
2026/01/05 11:39:41 [KuCoin] subscribed: MOVR-ETH
2026/01/05 11:39:41 [KuCoin] subscribed: ETH-DAI
2026/01/05 11:39:41 [KuCoin] subscribed: AR-USDT
2026/01/05 11:39:41 [KuCoin] subscribed: XCN-USDT
2026/01/05 11:39:41 [KuCoin] subscribed: KLV-TRX
2026/01/05 11:39:41 [KuCoin] subscribed: TRAC-BTC
2026/01/05 11:39:41 [KuCoin] subscribed: PERP-BTC
2026/01/05 11:39:41 [KuCoin] subscribed: BCHSV-USDT
2026/01/05 11:39:41 [KuCoin] subscribed: SUI-KCS
2026/01/05 11:39:41 [KuCoin] subscribed: KCS-ETH
2026/01/05 11:39:41 [KuCoin] subscribed: IOTA-USDT
2026/01/05 11:39:41 [KuCoin] subscribed: PAXG-USDT
2026/01/05 11:39:41 [KuCoin] subscribed: ZIL-ETH
2026/01/05 11:39:41 [KuCoin] subscribed: ENJ-ETH
2026/01/05 11:39:41 [KuCoin] subscribed: TRX-ETH
2026/01/05 11:39:41 [KuCoin] subscribed: NEAR-BTC
2026/01/05 11:39:41 [KuCoin] subscribed: BCH-USDT
2026/01/05 11:39:41 [KuCoin] subscribed: OGN-BTC
2026/01/05 11:39:41 [KuCoin] subscribed: ONT-USDT
2026/01/05 11:39:41 [KuCoin] subscribed: DYP-ETH
2026/01/05 11:39:41 [KuCoin] subscribed: SKL-BTC
2026/01/05 11:39:41 [KuCoin] subscribed: ICP-USDT
2026/01/05 11:39:41 [KuCoin] subscribed: AAVE-BTC
2026/01/05 11:39:41 [KuCoin] subscribed: COTI-USDT
2026/01/05 11:39:41 [KuCoin] subscribed: IOTX-ETH
2026/01/05 11:39:41 [KuCoin] subscribed: BDX-BTC
2026/01/05 11:39:41 [KuCoin] subscribed: NKN-BTC
2026/01/05 11:39:42 [KuCoin] subscribed: XRP-BTC
2026/01/05 11:39:42 [KuCoin] subscribed: ZEC-BTC
2026/01/05 11:39:42 [KuCoin] subscribed: BTC-BRL
2026/01/05 11:39:42 [KuCoin] subscribed: AVA-USDT
2026/01/05 11:39:42 [KuCoin] subscribed: GAS-BTC
2026/01/05 11:39:42 [KuCoin] subscribed: INJ-BTC
2026/01/05 11:39:42 [KuCoin] subscribed: CRO-BTC
2026/01/05 11:39:42 [KuCoin] subscribed: CRO-USDT
2026/01/05 11:39:42 [KuCoin] subscribed: TRAC-ETH
2026/01/05 11:39:42 [KuCoin] subscribed: A-ETH
2026/01/05 11:39:42 [KuCoin] subscribed: ETC-ETH
2026/01/05 11:39:42 [KuCoin] subscribed: ETH-EUR
2026/01/05 11:39:42 [KuCoin] subscribed: ETC-USDT
2026/01/05 11:39:42 [KuCoin] subscribed: BNB-USDT
2026/01/05 11:39:42 [KuCoin] subscribed: ONT-BTC
2026/01/05 11:39:42 [KuCoin] subscribed: BTC-USDT
2026/01/05 11:39:42 [KuCoin] subscribed: AVAX-USDT
2026/01/05 11:39:42 [KuCoin] subscribed: SCRT-BTC
2026/01/05 11:39:42 [KuCoin] subscribed: ALGO-USDT
2026/01/05 11:39:42 [KuCoin] subscribed: SKL-USDT
2026/01/05 11:39:42 [KuCoin] subscribed: DOGE-BTC
2026/01/05 11:39:42 [KuCoin] subscribed: TEL-USDT
2026/01/05 11:39:42 [KuCoin] subscribed: KRL-USDT
2026/01/05 11:39:42 [KuCoin] subscribed: COTI-BTC
2026/01/05 11:39:42 [KuCoin] subscribed: CSPR-USDT
2026/01/05 11:39:42 [KuCoin] subscribed: BAX-BTC
2026/01/05 11:39:42 [KuCoin] subscribed: ANKR-USDT
2026/01/05 11:39:42 [KuCoin] subscribed: VRA-USDT
2026/01/05 11:39:42 [KuCoin] subscribed: DOGE-KCS
2026/01/05 11:39:42 [KuCoin] subscribed: AR-BTC
2026/01/05 11:39:42 [KuCoin] subscribed: CKB-USDT
2026/01/05 11:39:42 [KuCoin] subscribed: STORJ-ETH
2026/01/05 11:39:42 [KuCoin] subscribed: VET-ETH
2026/01/05 11:39:42 [KuCoin] subscribed: XRP-ETH
2026/01/05 11:39:42 [KuCoin] subscribed: RLC-USDT
2026/01/05 11:39:42 [KuCoin] subscribed: XLM-BTC
2026/01/05 11:39:42 [KuCoin] subscribed: XMR-USDT
2026/01/05 11:39:42 [KuCoin] subscribed: SOL-KCS
2026/01/05 11:39:42 [KuCoin] subscribed: DAG-ETH
2026/01/05 11:39:42 [KuCoin] subscribed: TWT-BTC
2026/01/05 11:39:42 [KuCoin] subscribed: HBAR-USDT
2026/01/05 11:39:42 [KuCoin] subscribed: VSYS-USDT
2026/01/05 11:39:42 [KuCoin] subscribed: BCH-BTC
2026/01/05 11:39:42 [KuCoin] subscribed: KAS-BTC
2026/01/05 11:39:42 [KuCoin] subscribed: SOL-USDT
2026/01/05 11:39:43 [KuCoin] subscribed: ETH-USDT
2026/01/05 11:39:43 [KuCoin] subscribed: XTZ-USDT
2026/01/05 11:39:43 [KuCoin] subscribed: SUPER-USDT
2026/01/05 11:39:43 [KuCoin] subscribed: SNX-USDT
2026/01/05 11:39:43 [KuCoin] subscribed: ATOM-BTC
2026/01/05 11:39:43 [KuCoin] subscribed: XDC-ETH
2026/01/05 11:39:43 [KuCoin] subscribed: BCHSV-ETH
2026/01/05 11:39:43 [KuCoin] subscribed: POND-USDT
2026/01/05 11:39:43 [KuCoin] subscribed: NEO-USDT
2026/01/05 11:39:43 [KuCoin] subscribed: XMR-BTC
2026/01/05 11:39:43 [KuCoin] subscribed: NFT-USDT
2026/01/05 11:39:43 [KuCoin] subscribed: MOVR-USDT
2026/01/05 11:39:43 [KuCoin] subscribed: BAX-USDT
2026/01/05 11:39:43 [KuCoin] subscribed: REQ-BTC
2026/01/05 11:39:43 [KuCoin] subscribed: SUPER-BTC
2026/01/05 11:39:43 [KuCoin] subscribed: WIN-TRX
2026/01/05 11:39:43 [KuCoin] subscribed: ETH-BRL
2026/01/05 11:39:43 [KuCoin] subscribed: ONE-USDT
2026/01/05 11:39:43 [KuCoin] subscribed: POND-BTC
2026/01/05 11:39:43 [KuCoin] subscribed: HBAR-BTC
2026/01/05 11:39:43 [KuCoin] subscribed: KRL-BTC
2026/01/05 11:39:43 [KuCoin] subscribed: IOST-ETH
2026/01/05 11:39:43 [KuCoin] subscribed: ADA-BTC
2026/01/05 11:39:43 [KuCoin] subscribed: STX-BTC
2026/01/05 11:39:43 [KuCoin] subscribed: DOT-BTC
2026/01/05 11:39:43 [KuCoin] subscribed: ELA-BTC
2026/01/05 11:39:43 [KuCoin] subscribed: LTC-BTC
2026/01/05 11:39:43 [KuCoin] subscribed: USDT-DAI
2026/01/05 11:39:43 [KuCoin] subscribed: KAS-USDT
2026/01/05 11:39:43 [KuCoin] subscribed: TEL-ETH
2026/01/05 11:39:44 [KuCoin] subscribed: TWT-USDT
2026/01/05 11:39:44 [KuCoin] subscribed: RUNE-USDT
2026/01/05 11:39:44 [KuCoin] subscribed: BNB-BTC
2026/01/05 11:39:44 [KuCoin] subscribed: USDT-BRL
2026/01/05 11:39:44 [KuCoin] subscribed: DYP-USDT
2026/01/05 11:39:44 [KuCoin] subscribed: XLM-ETH
2026/01/05 11:39:44 [KuCoin] subscribed: LINK-BTC
2026/01/05 11:39:44 [KuCoin] subscribed: WIN-BTC
2026/01/05 11:39:44 [KuCoin] subscribed: IOTX-BTC
2026/01/05 11:39:44 [KuCoin] subscribed: TEL-BTC
2026/01/05 11:39:44 [KuCoin] subscribed: ETH-BTC
2026/01/05 11:39:44 [KuCoin] subscribed: KNC-BTC
2026/01/05 11:39:44 [KuCoin] subscribed: FET-BTC
2026/01/05 11:39:44 [KuCoin] subscribed: EWT-BTC
2026/01/05 11:39:44 [KuCoin] subscribed: DASH-BTC
2026/01/05 11:39:44 [KuCoin] subscribed: ADA-KCS
2026/01/05 11:39:44 [KuCoin] subscribed: OGN-USDT
2026/01/05 11:39:44 [KuCoin] subscribed: DASH-USDT
2026/01/05 11:39:44 [KuCoin] subscribed: WAN-USDT
2026/01/05 11:39:44 [KuCoin] subscribed: LYX-ETH
2026/01/05 11:39:44 [KuCoin] subscribed: XTZ-BTC
2026/01/05 11:39:44 [KuCoin] subscribed: ICP-BTC
2026/01/05 11:39:44 [KuCoin] subscribed: XYO-USDT
2026/01/05 11:39:44 [KuCoin] subscribed: AVA-ETH
2026/01/05 11:39:44 [KuCoin] subscribed: SXP-USDT
2026/01/05 11:39:44 [KuCoin] subscribed: SHIB-DOGE
2026/01/05 11:39:44 [KuCoin] subscribed: DASH-ETH
2026/01/05 11:39:44 [KuCoin] subscribed: SNX-BTC
2026/01/05 11:39:44 [KuCoin] subscribed: ZEC-USDT
2026/01/05 11:39:44 [KuCoin] subscribed: BTC-DAI
2026/01/05 11:39:44 [KuCoin] subscribed: XMR-ETH
2026/01/05 11:39:44 [KuCoin] subscribed: ICX-ETH
2026/01/05 11:39:44 [KuCoin] subscribed: FET-ETH
2026/01/05 11:39:44 [KuCoin] subscribed: DAG-BTC
2026/01/05 11:39:44 [KuCoin] subscribed: FET-USDT
2026/01/05 11:39:44 [KuCoin] subscribed: CSPR-ETH


