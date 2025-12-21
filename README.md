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
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"arb_project/models"
)

// OKXCollector реализует Collector для OKX
type OKXCollector struct {
	wsConn *websocket.Conn
	done   chan struct{}
}

// NewOKXCollector создает новый экземпляр OKXCollector
func NewOKXCollector() *OKXCollector {
	return &OKXCollector{
		done: make(chan struct{}),
	}
}

// Start подключается к WebSocket OKX и начинает отправку MarketData в dataCh
func (c *OKXCollector) Start(dataCh chan<- models.MarketData) error {
	var err error
	c.wsConn, _, err = websocket.DefaultDialer.Dial("wss://ws.okx.com:8443/ws/v5/public", nil)
	if err != nil {
		return fmt.Errorf("OKX websocket dial error: %v", err)
	}

	// Подписка на тикеры спот-рынка
	sub := map[string]interface{}{
		"op": "subscribe",
		"args": []map[string]string{
			{
				"channel": "tickers",
				"instType": "SPOT",
			},
		},
	}
	if err := c.wsConn.WriteJSON(sub); err != nil {
		return fmt.Errorf("OKX subscribe error: %v", err)
	}

	// Запуск цикла чтения
	go c.readLoop(dataCh)
	return nil
}

// readLoop читает сообщения с WebSocket и отправляет их в канал MarketData
func (c *OKXCollector) readLoop(dataCh chan<- models.MarketData) {
	for {
		select {
		case <-c.done:
			return
		default:
			_, msg, err := c.wsConn.ReadMessage()
			if err != nil {
				log.Println("OKX read error:", err)
				time.Sleep(time.Second)
				continue
			}

			var resp map[string]interface{}
			if err := json.Unmarshal(msg, &resp); err != nil {
				log.Println("OKX unmarshal error:", err)
				continue
			}

			data, ok := resp["data"].([]interface{})
			if !ok {
				continue
			}

			for _, d := range data {
				item, ok := d.(map[string]interface{})
				if !ok {
					continue
				}

				md := models.MarketData{
					Exchange:  "OKX",
					Symbol:    item["instId"].(string),
					Bid:       parseFloat(item["bidPx"]),
					Ask:       parseFloat(item["askPx"]),
					Timestamp: time.Now().UnixMilli(),
				}
				dataCh <- md
			}
		}
	}
}

// Stop закрывает WebSocket соединение
func (c *OKXCollector) Stop() error {
	close(c.done)
	if c.wsConn != nil {
		return c.wsConn.Close()
	}
	return nil
}

// parseFloat — вспомогательная функция для конвертации интерфейса в float64
func parseFloat(val interface{}) float64 {
	switch v := val.(type) {
	case string:
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return 0
		}
		return f
	case float64:
		return v
	default:
		return 0
	}
}





package main

import (
	"fmt"
	"time"
	"arb_project/collector"
	"arb_project/models"
)

func main() {
	dataCh := make(chan models.MarketData, 100)
	okxCollector := collector.NewOKXCollector()
	
	err := okxCollector.Start(dataCh)
	if err != nil {
		fmt.Println("Ошибка запуска Collector:", err)
		return
	}

	// Вывод данных из канала
	go func() {
		for md := range dataCh {
			fmt.Printf("Collector: %s %s Bid=%.2f Ask=%.2f\n",
				md.Exchange, md.Symbol, md.Bid, md.Ask)
		}
	}()

	time.Sleep(10 * time.Second) // пусть Collector поработает
	okxCollector.Stop()
}


