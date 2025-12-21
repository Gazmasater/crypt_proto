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
	"crypt_proto/pkg/models"
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

const mexcWS = "wss://wbs.mexc.com/ws"

type MEXCCollector struct {
	ctx     context.Context
	cancel  context.CancelFunc
	symbols []string
}

func NewMEXCCollector(symbols []string) *MEXCCollector {
	ctx, cancel := context.WithCancel(context.Background())
	upperSymbols := make([]string, len(symbols))
	for i, s := range symbols {
		upperSymbols[i] = strings.ToUpper(strings.ReplaceAll(s, "-", "_"))
	}
	return &MEXCCollector{
		ctx:     ctx,
		cancel:  cancel,
		symbols: upperSymbols,
	}
}

func (c *MEXCCollector) Name() string {
	return "MEXC"
}

func (c *MEXCCollector) Start(out chan<- models.MarketData) error {
	go c.run(out)
	return nil
}

func (c *MEXCCollector) Stop() error {
	c.cancel()
	return nil
}

func (c *MEXCCollector) run(out chan<- models.MarketData) {
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			log.Println("[MEXC] connecting...")
			c.connectAndRead(out)
			log.Println("[MEXC] reconnect in 1s...")
			time.Sleep(time.Second)
		}
	}
}

func (c *MEXCCollector) connectAndRead(out chan<- models.MarketData) {
	conn, _, err := websocket.DefaultDialer.Dial(mexcWS, nil)
	if err != nil {
		log.Println("[MEXC] dial error:", err)
		return
	}
	defer conn.Close()

	// Формируем батч подписки (до 25 символов за раз)
	params := make([]string, len(c.symbols))
	for i, s := range c.symbols {
		params[i] = "spot@public.ticker.v3.api@" + s
	}

	subscribe := map[string]interface{}{
		"method": "sub.ticker",
		"params": params,
		"id":     1,
	}

	if err := conn.WriteJSON(subscribe); err != nil {
		log.Println("[MEXC] subscribe error:", err)
		return
	}

	// Ping каждые 20 секунд
	go func() {
		ticker := time.NewTicker(20 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-c.ctx.Done():
				return
			case <-ticker.C:
				_ = conn.WriteJSON(map[string]interface{}{"method": "ping"})
			}
		}
	}()

	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println("[MEXC] read error:", err)
				return
			}
			c.handleMessage(msg, out)
		}
	}
}

func (c *MEXCCollector) handleMessage(msg []byte, out chan<- models.MarketData) {
	var raw struct {
		Method string `json:"method"`
		Params []struct {
			Symbol string `json:"s"`
			Bid    string `json:"b"`
			Ask    string `json:"a"`
		} `json:"params"`
	}

	if err := json.Unmarshal(msg, &raw); err != nil {
		return
	}

	if raw.Method != "ticker.update" || len(raw.Params) == 0 {
		return
	}

	for _, d := range raw.Params {
		bid, err1 := strconv.ParseFloat(d.Bid, 64)
		ask, err2 := strconv.ParseFloat(d.Ask, 64)
		if err1 != nil || err2 != nil {
			continue
		}

		out <- models.MarketData{
			Exchange:  "MEXC",
			Symbol:    d.Symbol,
			Bid:       bid,
			Ask:       ask,
			Timestamp: time.Now().UnixMilli(),
		}
	}
}




package main

import (
	"crypt_proto/pkg/collector"
	"crypt_proto/pkg/models"
	"fmt"
	"os"
	"strings"
)

func main() {
	exchange := strings.ToLower(os.Getenv("EXCHANGE"))
	if exchange == "" {
		exchange = "okx"
	}

	symbolEnv := os.Getenv("SYMBOLS") // допустим, SYMBOLS=BTCUSDT,ETHUSDT,XRPUSDT
	if symbolEnv == "" {
		symbolEnv = "BTCUSDT"
	}
	symbols := strings.Split(symbolEnv, ",")

	fmt.Println("EXCHANGE!!!!!!!!!", exchange)
	fmt.Println("SYMBOLS!!!!!!!!!", symbols)

	marketDataCh := make(chan models.MarketData, 1000)

	var c collector.Collector

	switch exchange {
	case "okx":
		c = collector.NewOKXCollector() // старый OKX
	case "mexc":
		c = collector.NewMEXCCollector(symbols)
	default:
		panic("unsupported exchange")
	}

	fmt.Println("Starting collector:", c.Name())
	if err := c.Start(marketDataCh); err != nil {
		panic(err)
	}

	// простой вывод данных
	for md := range marketDataCh {
		fmt.Printf("%s %s bid=%.6f ask=%.6f\n", md.Exchange, md.Symbol, md.Bid, md.Ask)
	}
}


