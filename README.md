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




(pprof) gaz358@gaz358-BOD-WXX9:~/myprog/crypt_proto$    go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
Fetching profile over HTTP from http://localhost:6060/debug/pprof/profile?seconds=30
Saved profile in /home/gaz358/pprof/pprof.arb.samples.cpu.081.pb.gz
File: arb
Build ID: 00f359f630cea5d5eb1389920b6bee5aa91f0b5e
Type: cpu
Time: 2026-01-12 10:59:06 MSK
Duration: 30.04s, Total samples = 1.98s ( 6.59%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 1130ms, 57.07% of 1980ms total
Showing top 10 nodes out of 209
      flat  flat%   sum%        cum   cum%
     670ms 33.84% 33.84%      670ms 33.84%  internal/runtime/syscall.Syscall6
     100ms  5.05% 38.89%      100ms  5.05%  runtime.futex
      70ms  3.54% 42.42%       70ms  3.54%  aeshashbody
      60ms  3.03% 45.45%      130ms  6.57%  runtime.scanobject
      50ms  2.53% 47.98%       90ms  4.55%  github.com/tidwall/gjson.parseObject
      50ms  2.53% 50.51%      130ms  6.57%  runtime.mapassign_faststr
      40ms  2.02% 52.53%       50ms  2.53%  runtime.typePointers.next
      30ms  1.52% 54.04%      820ms 41.41%  bufio.(*Reader).fill
      30ms  1.52% 55.56%       30ms  1.52%  memeqbody
      30ms  1.52% 57.07%       60ms  3.03%  runtime.mapaccess1_faststr
(pprof) 



Fetching profile over HTTP from http://localhost:6060/debug/pprof/profile?seconds=30
Saved profile in /home/gaz358/pprof/pprof.arb.samples.cpu.082.pb.gz
File: arb
Build ID: 991d3b51d26d0a48852c28a66aa2039c318c2e53
Type: cpu
Time: 2026-01-12 11:20:48 MSK
Duration: 30s, Total samples = 1.55s ( 5.17%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 1100ms, 70.97% of 1550ms total
Showing top 10 nodes out of 127
      flat  flat%   sum%        cum   cum%
     760ms 49.03% 49.03%      760ms 49.03%  internal/runtime/syscall.Syscall6
     170ms 10.97% 60.00%      170ms 10.97%  runtime.futex
      30ms  1.94% 61.94%       40ms  2.58%  strings.Fields
      20ms  1.29% 63.23%       20ms  1.29%  crypto/internal/fips140/aes/gcm.gcmAesDec
      20ms  1.29% 64.52%      860ms 55.48%  crypto/tls.(*Conn).readRecordOrCCS
      20ms  1.29% 65.81%      950ms 61.29%  github.com/gorilla/websocket.(*Conn).ReadMessage
      20ms  1.29% 67.10%       20ms  1.29%  github.com/gorilla/websocket.(*messageReader).Read
      20ms  1.29% 68.39%       20ms  1.29%  runtime.(*mspan).base
      20ms  1.29% 69.68%       20ms  1.29%  runtime.execute
      20ms  1.29% 70.97%       30ms  1.94%  runtime.ifaceeq
(pprof) 



func (ws *kucoinWS) handle(c *KuCoinCollector, msg []byte) {
	// Проверяем тип сообщения
	if gjson.GetBytes(msg, "type").String() != "message" {
		return
	}

	// Проверяем топик
	topic := gjson.GetBytes(msg, "topic").String()
	if !strings.HasPrefix(topic, "/market/ticker:") {
		return
	}

	// Берём символ из топика, оставляем дефис
	symbol := strings.TrimPrefix(topic, "/market/ticker:") // пример: "MANA-USDT"

	// Получаем данные
	data := gjson.GetBytes(msg, "data")
	bid := data.Get("bestBid").Float()
	ask := data.Get("bestAsk").Float()
	bidSize := data.Get("bestBidSize").Float()
	askSize := data.Get("bestAskSize").Float()

	if bid == 0 || ask == 0 {
		return
	}

	// Проверка на изменение цены
	ws.mu.Lock()
	last := ws.last[symbol]
	if last[0] == bid && last[1] == ask {
		ws.mu.Unlock()
		return
	}
	ws.last[symbol] = [2]float64{bid, ask}
	ws.mu.Unlock()

	// Лог для отладки
	log.Printf("[WS %d] %s bid=%.6f ask=%.6f", ws.id, symbol, bid, ask)

	// Отправка в MemoryStore / канал калькулятора
	c.out <- &models.MarketData{
		Exchange:  "KuCoin",
		Symbol:    symbol, // совпадает с ключом в MemoryStore
		Bid:       bid,
		Ask:       ask,
		BidSize:   bidSize,
		AskSize:   askSize,
		Timestamp: time.Now().UnixMilli(),
	}
}



// legSymbol возвращает BASE-QUOTE из Leg, используется в калькуляторе
func legSymbol(leg string) string {
	parts := strings.Fields(strings.ToUpper(strings.TrimSpace(leg)))
	if len(parts) < 2 {
		return ""
	}
	p := strings.Split(parts[1], "/")
	if len(p) != 2 {
		return ""
	}
	return p[0] + "-" + p[1] // формат BASE-QUOTE
}






