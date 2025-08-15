package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"

	pb "crypt_proto/pb" // твой пакет со сгенерёнными *.pb.go
)

// хранит последний top по символу для дедупликации
type top struct{ bp, bq, ap, aq string } // bidPrice, bidQty, askPrice, askQty
var lastTop = map[string]top{}

func main() {
	const wsURL = "wss://wbs-api.mexc.com/ws"

	// 1) Подключение к публичному WS
	c, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	// 2) Подписка на публичные Protobuf-каналы bookTicker
	sub := map[string]any{
		"method": "SUBSCRIPTION",
		"params": []string{
			"spot@public.aggre.bookTicker.v3.api.pb@100ms@BTCUSDT",
			"spot@public.aggre.bookTicker.v3.api.pb@100ms@ETHUSDT",
			"spot@public.aggre.bookTicker.v3.api.pb@100ms@ETHBTC",
		},
	}
	if err := c.WriteJSON(sub); err != nil {
		log.Fatal("send sub:", err)
	}

	// 3) Пинги для поддержания соединения
	go func() {
		t := time.NewTicker(45 * time.Second)
		defer t.Stop()
		for range t.C {
			_ = c.WriteMessage(websocket.PingMessage, []byte("hb"))
		}
	}()

	// 4) Чтение сообщений
	for {
		mt, raw, err := c.ReadMessage()
		if err != nil {
			log.Fatal("read:", err)
		}

		// ACK/ошибки — текст/JSON
		if mt == websocket.TextMessage {
			var v any
			if json.Unmarshal(raw, &v) == nil {
				b, _ := json.MarshalIndent(v, "", "  ")
				fmt.Printf("ACK:\n%s\n", b)
			} else {
				fmt.Printf("TEXT:\n%s\n", string(raw))
			}
			continue
		}
		if mt != websocket.BinaryMessage {
			continue
		}

		// Декодируем обёртку Protobuf
		var w pb.PushDataV3ApiWrapper
		if err := proto.Unmarshal(raw, &w); err != nil {
			// не наше сообщение — пропускаем
			continue
		}

		// symbol/ts
		symbol := w.GetSymbol()
		if symbol == "" {
			if ch := w.GetChannel(); ch != "" {
				parts := strings.Split(ch, "@")
				symbol = parts[len(parts)-1]
			}
		}
		ts := time.Now()
		if t := w.GetSendTime(); t > 0 {
			ts = time.UnixMilli(t)
		}

		// Интересует PublicAggreBookTicker — выводим только при изменениях
		switch body := w.GetBody().(type) {
		case *pb.PushDataV3ApiWrapper_PublicAggreBookTicker:
			bt := body.PublicAggreBookTicker

			// строки для дедупа
			bpS, bqS := bt.GetBidPrice(), bt.GetBidQuantity()
			apS, aqS := bt.GetAskPrice(), bt.GetAskQuantity()
			cur := top{bpS, bqS, apS, aqS}

			// если значения идентичны предыдущим — пропускаем вывод
			if prev, ok := lastTop[symbol]; ok && prev == cur {
				continue // перейти к следующему сообщению цикла for
			}
			lastTop[symbol] = cur

			// парсим в числа и печатаем в исходном формате
			bid, _ := strconv.ParseFloat(bpS, 64)
			ask, _ := strconv.ParseFloat(apS, 64)
			bq, _ := strconv.ParseFloat(bqS, 64)
			aq, _ := strconv.ParseFloat(aqS, 64)

			fmt.Printf("%s  bid=%.8f (%.6f)  ask=%.8f (%.6f)  ts=%s\n",
				symbol, bid, bq, ask, aq, ts.Format(time.RFC3339Nano))

		default:
			// другие типы тел не выводим
		}
	}
}
