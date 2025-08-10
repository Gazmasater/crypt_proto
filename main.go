package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	// 👇 ЗАМЕНИ на реальные пакеты/типы из твоих *_pb.go
	bookpb "crypt_proto/pb"    // из PublicAggreBookTickerV3Api.proto
	wrapperpb "crypt_proto/pb" // из PushDataV3ApiWrapper.proto
)

func main() {
	const wsURL = "wss://wbs-api.mexc.com/ws"

	c, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()
	log.Println("connected")

	// подписка на три топика
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
	log.Println("subscription sent")

	// пингуем периодически
	go func() {
		t := time.NewTicker(45 * time.Second)
		defer t.Stop()
		for range t.C {
			_ = c.WriteMessage(websocket.PingMessage, []byte("hb"))
		}
	}()

	// читаем бесконечно
	for {
		mt, msg, err := c.ReadMessage()
		if err != nil {
			log.Fatal("read:", err)
		}

		// ACK/ошибки приходят как TEXT/JSON — распечатаем красиво
		if mt == websocket.TextMessage {
			var v any
			if json.Unmarshal(msg, &v) == nil {
				pre, _ := json.MarshalIndent(v, "", "  ")
				fmt.Printf("ACK:\n%s\n\n", pre)
			} else {
				fmt.Printf("TEXT:\n%s\n\n", string(msg))
			}
			continue
		}

		if mt != websocket.BinaryMessage {
			continue
		}

		// 1) Декодируем внешнюю обёртку
		//    Тип возьми из PushDataV3ApiWrapper.proto (например, PushDataV3ApiWrapper / PushData)
		var w wrapperpb.PushDataV3ApiWrapper // <-- подставь точное имя типа
		if err := proto.Unmarshal(msg, &w); err != nil {
			log.Printf("wrapper unmarshal: %v", err)
			continue
		}

		// (необязательно) покажем метаданные обёртки
		wJSON, _ := protojson.MarshalOptions{EmitUnpopulated: true}.Marshal(&w)
		fmt.Printf("WRAPPER: %s\n", wJSON)

		// 2) Декодируем полезную нагрузку как BookTicker
		//    Тип возьми из PublicAggreBookTickerV3Api.proto (например, PublicAggreBookTickerV3Api / BookTicker)
		var bt bookpb.PublicAggreBookTickerV3Api // <-- подставь точное имя типа
		if err := proto.Unmarshal(w.GetD(), &bt); err != nil {
			log.Printf("bookTicker unmarshal: %v", err)
			continue
		}

		// 3) Выведем в JSON (универсально, без знания имён полей)
		out, _ := protojson.MarshalOptions{EmitUnpopulated: true}.Marshal(&bt)
		fmt.Printf("BOOK_TICKER: %s\n\n", out)

		// Если хочешь конкретные поля (symbol/bid/ask/ts), используй геттеры:
		// fmt.Println(bt.GetS(), bt.GetBp(), bt.GetAp(), bt.GetT())
	}
}
