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

func main() {
	const wsURL = "wss://wbs-api.mexc.com/ws"

	c, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

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

	// поддерживаем линк
	go func() {
		t := time.NewTicker(45 * time.Second)
		defer t.Stop()
		for range t.C {
			_ = c.WriteMessage(websocket.PingMessage, []byte("hb"))
		}
	}()

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

		// 1) Декодируем ОБЁРТКУ
		var w pb.PushDataV3ApiWrapper
		if err := proto.Unmarshal(raw, &w); err != nil {
			// бинарь не нашей схемы — пропускаем
			continue
		}

		// 2) Вытаскиваем symbol/ts из обёртки
		symbol := w.GetSymbol()
		if symbol == "" {
			// если symbol пуст — берём из channel (последний сегмент после '@')
			ch := w.GetChannel()
			if ch != "" {
				parts := strings.Split(ch, "@")
				symbol = parts[len(parts)-1]
			}
		}
		ts := time.Now()
		if t := w.GetSendTime(); t > 0 {
			ts = time.UnixMilli(t)
		}

		// 3) oneof Body → нас интересует PublicAggreBookTicker
		switch body := w.GetBody().(type) {
		case *pb.PushDataV3ApiWrapper_PublicAggreBookTicker:
			bt := body.PublicAggreBookTicker

			bid, _ := strconv.ParseFloat(bt.GetBidPrice(), 64)
			ask, _ := strconv.ParseFloat(bt.GetAskPrice(), 64)
			bq, _ := strconv.ParseFloat(bt.GetBidQuantity(), 64)
			aq, _ := strconv.ParseFloat(bt.GetAskQuantity(), 64)

			fmt.Printf("%s  bid=%.8f (%.6f)  ask=%.8f (%.6f)  ts=%s\n",
				symbol, bid, bq, ask, aq, ts.Format(time.RFC3339Nano))

		// (не обязательно) если вдруг придёт другой кейс — можно игнорить
		default:
			// fmt.Printf("other body: %T\n", body)
		}
	}
}
