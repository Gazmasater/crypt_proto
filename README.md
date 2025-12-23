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



parts := strings.Split(ch, "@")
symbol = parts[len(parts)-1]





package collector

import (
	"context"
	"encoding/json"
	"log"
	"strings"
	"time"

	"crypt_proto/configs"
	"crypt_proto/pkg/models"

	"github.com/gorilla/websocket"
)

type KuCoinCollector struct {
	ctx     context.Context
	cancel  context.CancelFunc
	symbols []string
	conn    *websocket.Conn
}

func NewKuCoinCollector(symbols []string) *KuCoinCollector {
	ctx, cancel := context.WithCancel(context.Background())
	return &KuCoinCollector{ctx: ctx, cancel: cancel, symbols: symbols}
}

func (c *KuCoinCollector) Name() string { return "KuCoin" }

func (c *KuCoinCollector) Start(out chan<- models.MarketData) error {
	conn, _, err := websocket.DefaultDialer.Dial(configs.KUCOIN_WS, nil)
	if err != nil {
		return err
	}
	c.conn = conn
	log.Println("[KuCoin] connected")

	// subscribe
	for _, s := range c.symbols {
		sub := map[string]interface{}{
			"id":       time.Now().Unix(),
			"type":     "subscribe",
			"topic":    "level2/ticker:" + s,
			"response": true,
		}
		if err := conn.WriteJSON(sub); err != nil {
			return err
		}
		log.Println("[KuCoin] subscribed:", s)
	}

	// ping loop
	go func() {
		t := time.NewTicker(configs.KUCOIN_PING_INTERVAL)
		defer t.Stop()
		for {
			select {
			case <-c.ctx.Done():
				return
			case <-t.C:
				_ = conn.WriteMessage(websocket.PingMessage, nil)
			}
		}
	}()

	// read loop
	go func() {
		defer conn.Close()
		for {
			select {
			case <-c.ctx.Done():
				return
			default:
				_, msg, err := conn.ReadMessage()
				if err != nil {
					log.Println("[KuCoin] read error:", err)
					return
				}

				var data map[string]interface{}
				if err := json.Unmarshal(msg, &data); err != nil {
					continue
				}

				// проверяем, что это update
				if data["type"] == "message" {
					topic := data["topic"].(string)
					symbol := strings.Split(topic, ":")[1]

					body := data["data"].(map[string]interface{})
					bid, _ := parseStringToFloat(body["bestBid"].(string))
					ask, _ := parseStringToFloat(body["bestAsk"].(string))

					out <- models.MarketData{
						Exchange: "KuCoin",
						Symbol:   symbol,
						Bid:      bid,
						Ask:      ask,
					}
				}
			}
		}
	}()

	return nil
}

func (c *KuCoinCollector) Stop() error {
	c.cancel()
	return nil
}

// вспомогательная функция
func parseStringToFloat(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}




package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"crypt_proto/internal/collector"
	"crypt_proto/pkg/models"
)

func main() {
	exchange := os.Getenv("EXCHANGE")
	if exchange == "" {
		log.Fatal("EXCHANGE не задан. Используй MEXC, OKX или KuCoin")
	}
	log.Println("EXCHANGE:", exchange)

	marketDataCh := make(chan models.MarketData, 1000)

	var c collector.Collector

	switch exchange {
	case "MEXC":
		c = collector.NewMEXCCollector([]string{
			"BTCUSDT",
			"ETHUSDT",
			"ETHBTC",
		})
	case "OKX":
		c = collector.NewOKXCollector([]string{
			"BTC-USDT",
			"ETH-USDT",
			"ETH-BTC",
		})
	case "KuCoin":
		c = collector.NewKuCoinCollector([]string{
			"BTC-USDT",
			"ETH-USDT",
			"ETH-BTC",
		})
	default:
		log.Fatal("Неподдерживаемый EXCHANGE:", exchange)
	}

	if err := c.Start(marketDataCh); err != nil {
		log.Fatal("Start error:", err)
	}

	// читаем данные в фоне
	go func() {
		for data := range marketDataCh {
			log.Printf("[%s] %s bid=%.6f ask=%.6f\n",
				data.Exchange, data.Symbol, data.Bid, data.Ask)
		}
	}()

	// корректное завершение по Ctrl+C
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	log.Println("Stopping collector...")
	if err := c.Stop(); err != nil {
		log.Println("Stop error:", err)
	}
}






