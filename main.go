package main

import (
	"crypt_proto/pb"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

func main() {
	header := http.Header{}
	header.Set("Sec-WebSocket-Protocol", "protobuf")

	conn, _, err := websocket.DefaultDialer.Dial("wss://wbs.mexc.com/ws", header)
	if err != nil {
		log.Fatal("❌ dial:", err)
	}
	defer conn.Close()

	// подписка на стакан BTCUSDT
	sub := map[string]interface{}{
		"method": "SUBSCRIPTION",
		"params": []string{"spot@public.depth.v3.api@BTCUSDT"},
		"id":     time.Now().Unix(),
	}
	if err := conn.WriteJSON(sub); err != nil {
		log.Fatal("❌ send:", err)
	}

	log.Println("🟢 Subscribed to depth. Waiting for protobuf messages...")

	for {
		mt, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("❌ read:", err)
			break
		}
		if mt != websocket.BinaryMessage {
			log.Printf("⚠️  Skip non-binary message: %s", message)
			continue
		}

		var depth pb.PublicAggreDepthsV3Api
		if err := proto.Unmarshal(message, &depth); err != nil {
			log.Println("❌ proto.Unmarshal:", err)
			continue
		}

		log.Printf("📊 Depth update: %d asks / %d bids | type: %s | version: %s → %s",
			len(depth.Asks), len(depth.Bids), depth.EventType, depth.FromVersion, depth.ToVersion)

		// Вывод первых 3 ask/bid
		for i := 0; i < 3 && i < len(depth.Asks); i++ {
			log.Printf("🟢 ASK %s @ %s", depth.Asks[i].Quantity, depth.Asks[i].Price)
		}
		for i := 0; i < 3 && i < len(depth.Bids); i++ {
			log.Printf("🔴 BID %s @ %s", depth.Bids[i].Quantity, depth.Bids[i].Price)
		}
	}
}
