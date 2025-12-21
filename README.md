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




const mexcWS = "wss://wbs-api.mexc.com/ws"


{
  "method": "sub.deals",
  "params": ["spot@public.deals.v3.api@BTCUSDT"],
  "id": 1
}


{
  "method": "sub.ticker",
  "params": ["spot@public.ticker.v3.api@BTCUSDT"],
  "id": 1
}




package collector

import (
	"context"
	"crypt_proto/pkg/models"
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

const mexcWS = "wss://wbs-api.mexc.com/ws"

type MEXCCollector struct {
	ctx    context.Context
	cancel context.CancelFunc
	symbol string
}

func NewMEXCCollector(symbol string) *MEXCCollector {
	ctx, cancel := context.WithCancel(context.Background())
	return &MEXCCollector{
		ctx:    ctx,
		cancel: cancel,
		symbol: symbol, // ожидаем формат: BTCUSDT
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

	// подписка на ticker через v3 API каналы
	subscribe := map[string]interface{}{
		"method": "sub.ticker",
		"params": []string{"spot@public.ticker.v3.api@" + c.symbol},
		"id":     1,
	}

	if err := conn.WriteJSON(subscribe); err != nil {
		log.Println("[MEXC] subscribe error:", err)
		return
	}

	// ping каждые 20s, иначе MEXC может разорвать
	go func() {
		ticker := time.NewTicker(20 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-c.ctx.Done():
				return
			case <-ticker.C:
				conn.WriteMessage(websocket.TextMessage, []byte(`{"method":"ping"}`))
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
			Symbol string `json:"symbol"`
			Bid    string `json:"bid"`
			Ask    string `json:"ask"`
		} `json:"params"`
	}

	if err := json.Unmarshal(msg, &raw); err != nil {
		return
	}

	for _, p := range raw.Params {
		bid, err1 := strconv.ParseFloat(p.Bid, 64)
		ask, err2 := strconv.ParseFloat(p.Ask, 64)
		if err1 != nil || err2 != nil {
			continue
		}

		out <- models.MarketData{
			Exchange:  "MEXC",
			Symbol:    p.Symbol,
			Bid:       bid,
			Ask:       ask,
			Timestamp: time.Now().UnixMilli(),
		}
	}
}


gaz358@gaz358-BOD-WXX9:~/myprog/crypt_proto/cmd/arb$ go run .
EXCHANGE!!!!!!!!! mexc
2025/12/21 23:46:17 Starting collector: MEXC
2025/12/21 23:46:17 [MEXC] connecting...
2025/12/21 23:46:53 [MEXC] read error: websocket: close 1005 (no status)
2025/12/21 23:46:53 [MEXC] reconnect in 1s...
2025/12/21 23:46:54 [MEXC] connecting...
2025/12/21 23:47:28 [MEXC] read error: websocket: close 1005 (no status)
2025/12/21 23:47:28 [MEXC] reconnect in 1s...
2025/12/21 23:47:29 [MEXC] connecting...
