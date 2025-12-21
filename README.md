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
	"log"
	"strings"
	"sync"
	"time"

	"crypt_proto/domain"
	pb "crypt_proto/pb"
	"crypt_proto/pkg/models"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

type MEXCCollector struct {
	ctx     context.Context
	cancel  context.CancelFunc
	symbols []string
	out     chan<- models.MarketData
	debug   bool
}

func NewMEXCCollector(symbols []string, debug bool) *MEXCCollector {
	ctx, cancel := context.WithCancel(context.Background())
	return &MEXCCollector{
		ctx:     ctx,
		cancel:  cancel,
		symbols: symbols,
		debug:   debug,
	}
}

func (c *MEXCCollector) Name() string { return "MEXC" }

func (c *MEXCCollector) Start(out chan<- models.MarketData) error {
	c.out = out
	const maxPerConn = 25

	chunks := make([][]string, 0)
	for i := 0; i < len(c.symbols); i += maxPerConn {
		j := i + maxPerConn
		if j > len(c.symbols) {
			j = len(c.symbols)
		}
		chunks = append(chunks, c.symbols[i:j])
	}

	var wg sync.WaitGroup
	for idx, chunk := range chunks {
		wg.Add(1)
		go c.runFeed(idx, chunk, &wg)
	}
	go func() {
		wg.Wait()
	}()
	return nil
}

func (c *MEXCCollector) Stop() error {
	c.cancel()
	return nil
}

// -------------------- WS Feed --------------------

func (c *MEXCCollector) runFeed(connID int, symbols []string, wg *sync.WaitGroup) {
	defer wg.Done()

	urlWS := "wss://wbs-api.mexc.com/ws"
	topics := make([]string, 0, len(symbols))
	for _, s := range symbols {
		topics = append(topics, "spot@public.aggre.bookTicker.v3.api.pb@1s@"+strings.ToUpper(s))
	}

	const (
		baseRetry = 2 * time.Second
		maxRetry  = 30 * time.Second
	)
	retry := baseRetry

	for {
		select {
		case <-c.ctx.Done():
			return
		default:
		}

		conn, _, err := websocket.DefaultDialer.Dial(urlWS, nil)
		if err != nil {
			log.Printf("[MEXC WS #%d] dial error: %v (retry in %v)", connID, err, retry)
			time.Sleep(retry)
			if retry < maxRetry {
				retry *= 2
				if retry > maxRetry {
					retry = maxRetry
				}
			}
			continue
		}
		log.Printf("[MEXC WS #%d] connected (symbols: %d)", connID, len(symbols))
		retry = baseRetry

		conn.SetReadDeadline(time.Now().Add(90 * time.Second))
		lastPing := time.Now()
		stopPing := make(chan struct{})

		conn.SetPongHandler(func(string) error {
			conn.SetReadDeadline(time.Now().Add(90 * time.Second))
			return nil
		})

		go func() {
			t := time.NewTicker(45 * time.Second)
			defer t.Stop()
			for {
				select {
				case <-t.C:
					lastPing = time.Now()
					_ = conn.WriteControl(websocket.PingMessage, []byte("hb"), time.Now().Add(5*time.Second))
				case <-stopPing:
					return
				}
			}
		}()

		sub := map[string]interface{}{
			"method": "SUBSCRIPTION",
			"params": topics,
			"id":     time.Now().Unix(),
		}
		if err := conn.WriteJSON(sub); err != nil {
			log.Printf("[MEXC WS #%d] subscribe error: %v", connID, err)
			close(stopPing)
			conn.Close()
			time.Sleep(retry)
			continue
		}
		log.Printf("[MEXC WS #%d] SUB -> %d topics", connID, len(topics))

		for {
			_, raw, err := conn.ReadMessage()
			if err != nil {
				log.Printf("[MEXC WS #%d] read error: %v (reconnect)", connID, err)
				break
			}
			c.handleProtoMessage(raw)
		}

		close(stopPing)
		conn.Close()
		time.Sleep(retry)
		if retry < maxRetry {
			retry *= 2
			if retry > maxRetry {
				retry = maxRetry
			}
		}
	}
}

// -------------------- Protobuf parser --------------------

var wrapperPool = sync.Pool{
	New: func() any { return new(pb.PushDataV3ApiWrapper) },
}

func (c *MEXCCollector) handleProtoMessage(raw []byte) {
	w, _ := wrapperPool.Get().(*pb.PushDataV3ApiWrapper)
	defer func() {
		*w = pb.PushDataV3ApiWrapper{}
		wrapperPool.Put(w)
	}()
	if err := proto.Unmarshal(raw, w); err != nil {
		return
	}

	var sym string
	ch := w.GetChannel()
	if i := strings.LastIndex(ch, "@"); i >= 0 && i+1 < len(ch) {
		sym = ch[i+1:]
	}
	if sym == "" {
		sym = w.GetSymbol()
	}
	if sym == "" {
		return
	}

	var bid, ask float64
	if b1, ok := w.GetBody().(*pb.PushDataV3ApiWrapper_PublicAggreBookTicker); ok && b1.PublicAggreBookTicker != nil {
		t := b1.PublicAggreBookTicker
		bid, _ = strconv.ParseFloat(t.GetBidPrice(), 64)
		ask, _ = strconv.ParseFloat(t.GetAskPrice(), 64)
		if bid <= 0 || ask <= 0 {
			return
		}

		c.out <- models.MarketData{
			Exchange:  "MEXC",
			Symbol:    sym,
			Bid:       bid,
			Ask:       ask,
			Timestamp: time.Now().UnixMilli(),
		}
	}
}






package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"crypt_proto/collector"
	"crypt_proto/pkg/models"
)

func main() {
	exchange := strings.ToLower(os.Getenv("EXCHANGE"))
	if exchange == "" {
		exchange = "okx"
	}

	fmt.Println("EXCHANGE!!!!!!!!!", exchange)

	marketDataCh := make(chan models.MarketData, 1000)

	var c collector.Collector

	symbols := []string{"BTCUSDT", "ETHUSDT", "XRPUSDT"} // пример батча

	switch exchange {
	case "okx":
		c = collector.NewOKXCollector()
	case "mexc":
		c = collector.NewMEXCCollector(symbols, true)
	default:
		panic("unknown exchange")
	}

	fmt.Println("Starting collector:", c.Name())
	if err := c.Start(marketDataCh); err != nil {
		panic(err)
	}

	// consumer
	go func() {
		for data := range marketDataCh {
			fmt.Printf("[%s] %s bid=%.4f ask=%.4f\n", data.Exchange, data.Symbol, data.Bid, data.Ask)
		}
	}()

	// run forever
	for {
		time.Sleep(time.Hour)
	}
}


