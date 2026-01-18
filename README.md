Название API
9623527002

696935c42a6dcd00013273f2
b348b686-55ff-4290-897b-02d55f815f65




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





gaz358@gaz358-BOD-WXX9:~/myprog/crypt_proto/test$ go run .
2026/01/18 19:11:12.540545 START TRIANGLE 12.00 USDT
2026/01/18 19:11:15.204222 LEG1 USDT→DASH | DASH=0.286400 | total=818.177249ms
2026/01/18 19:11:16.023274 LEG2 DASH→BTC | BTC=0.00026000 | total=819.003358ms
2026/01/18 19:11:16.434442 LEG3 BTC→USDT | BTC sold=0.00026000
2026/01/18 19:11:16.843108 LEG3 BTC→USDT | USDT=34.1307 | total=819.808459ms
2026/01/18 19:11:16.843143 ====== TRIANGLE SUMMARY ======
2026/01/18 19:11:16.843148 LEG1 time: 818.177249ms
2026/01/18 19:11:16.843152 LEG2 time: 819.003358ms
2026/01/18 19:11:16.843156 LEG3 time: 819.808459ms
2026/01/18 19:11:16.843160 TOTAL time: 4.302597573s
2026/01/18 19:11:16.843164 PNL: 22.130750 USDT




gaz358@gaz358-BOD-WXX9:~/myprog/crypt_proto/step$ go run .
2026/01/18 19:28:47 Getting steps for triangle: USDT → DASH → BTC → USDT
====== TRADING STEPS ======
DASH-USDT step: 0.00010000
DASH-BTC step: 0.00010000
BTC-USDT step: 0.00001000





func subscribeOrders(conn *websocket.Conn) {
	sub := map[string]interface{}{
		"id":             1,
		"type":           "subscribe",
		"topic":          "/spot/tradeOrders", // правильный канал для своих ордеров
		"privateChannel": true,
		"response":       true,
	}
	raw, _ := json.Marshal(sub)
	conn.WriteMessage(websocket.TextMessage, raw)

	// Ждём подтверждения подписки
	_, msg, err := conn.ReadMessage()
	if err != nil {
		log.Fatal("WS subscribe response error:", err)
	}
	log.Println("WS subscription response:", string(msg))

	// Короткая пауза, чтобы WS успел начать получать события
	time.Sleep(100 * time.Millisecond)
}








