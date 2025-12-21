package collector

import (
	"crypt_proto/pkg/models"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
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
				"channel":  "tickers",
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
