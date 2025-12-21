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





✅ internal/collector/okx_collector.go (адаптирован под models.MarketData)
package collector

import (
	"context"
	"crypt_proto/models"
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

const okxWS = "wss://ws.okx.com:8443/ws/v5/public"

type OKXCollector struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func NewOKXCollector() *OKXCollector {
	ctx, cancel := context.WithCancel(context.Background())
	return &OKXCollector{
		ctx:    ctx,
		cancel: cancel,
	}
}

func (c *OKXCollector) Name() string {
	return "OKX"
}

func (c *OKXCollector) Start(out chan<- models.MarketData) error {
	go c.run(out)
	return nil
}

func (c *OKXCollector) Stop() {
	c.cancel()
}

func (c *OKXCollector) run(out chan<- models.MarketData) {
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			log.Println("[OKX] connecting...")
			c.connectAndRead(out)
			log.Println("[OKX] reconnect in 1s...")
			time.Sleep(time.Second)
		}
	}
}

func (c *OKXCollector) connectAndRead(out chan<- models.MarketData) {
	conn, _, err := websocket.DefaultDialer.Dial(okxWS, nil)
	if err != nil {
		log.Println("[OKX] dial error:", err)
		return
	}
	defer conn.Close()

	subscribe := map[string]interface{}{
		"op": "subscribe",
		"args": []map[string]string{
			{
				"channel": "tickers",
				"instId":  "BTC-USDT",
			},
		},
	}

	if err := conn.WriteJSON(subscribe); err != nil {
		log.Println("[OKX] subscribe error:", err)
		return
	}

	// keepalive ping
	go func() {
		ticker := time.NewTicker(20 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-c.ctx.Done():
				return
			case <-ticker.C:
				_ = conn.WriteMessage(websocket.PingMessage, nil)
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
				log.Println("[OKX] read error:", err)
				return
			}
			c.handleMessage(msg, out)
		}
	}
}

func (c *OKXCollector) handleMessage(msg []byte, out chan<- models.MarketData) {
	var raw struct {
		Data []struct {
			InstId string `json:"instId"`
			BidPx  string `json:"bidPx"`
			AskPx  string `json:"askPx"`
		} `json:"data"`
	}

	if err := json.Unmarshal(msg, &raw); err != nil {
		return
	}

	for _, d := range raw.Data {
		bid, err1 := strconv.ParseFloat(d.BidPx, 64)
		ask, err2 := strconv.ParseFloat(d.AskPx, 64)
		if err1 != nil || err2 != nil {
			continue
		}

		out <- models.MarketData{
			Exchange:  "OKX",
			Symbol:    d.InstId,
			Bid:       bid,
			Ask:       ask,
			Timestamp: time.Now().UnixMilli(),
		}
	}
}

🔎 Что изменилось принципиально
Было	Стало
MarketTick	models.MarketData
BidPrice / AskPrice	Bid / Ask
time.Time	int64 (Unix ms)
ручной Stop	context.Cancel()
🧠 Архитектурно важно

Collector → только MarketData

Calculator → принимает MarketData[], отдаёт Signal

Executor → работает ТОЛЬКО с Signal

Это чистый SOLID:

SRP ✔

DIP ✔

Open/Closed ✔

✅ Проверка, что всё работает

В main.go:

for md := range marketDataCh {
	log.Printf(
		"%s %s bid=%.2f ask=%.2f",
		md.Exchange, md.Symbol, md.Bid, md.Ask,
	)
}


Если bid/ask обновляются → коллектор идеален.

Следующий логичный шаг

1️⃣ вынести Collector интерфейс
2️⃣ добавить MEXCCollector (почти копипаста)
3️⃣ сделать fan-in канал для калькулятора

Скажи — идём дальше по какому шагу.

перепиши main gjkyjcnm.

Понял 👍
Переписываю main.go полностью и аккуратно, под твою текущую архитектуру:

один процесс

без контейнеров

тест коллектора OKX

корректное завершение

готово к расширению (MEXC / KuCoin)

✅ cmd/arb/main.go (чистая, правильная версия)
package main

import (
	"crypt_proto/internal/collector"
	"crypt_proto/models"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// канал для рыночных данных
	marketDataCh := make(chan models.MarketData, 1000)

	// инициализация коллектора
	okxCollector := collector.NewOKXCollector()

	// старт коллектора
	if err := okxCollector.Start(marketDataCh); err != nil {
		log.Fatal("failed to start OKX collector:", err)
	}

	log.Println("OKX collector started")

	// consumer (пока просто логируем)
	go func() {
		for md := range marketDataCh {
			log.Printf(
				"[MARKET] %s %s bid=%.4f ask=%.4f",
				md.Exchange,
				md.Symbol,
				md.Bid,
				md.Ask,
			)
		}
	}()

	// graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	<-sigCh
	log.Println("shutdown signal received")

	okxCollector.Stop()

	log.Println("collector stopped, exit")
}

🔍 Что здесь важно архитектурно


