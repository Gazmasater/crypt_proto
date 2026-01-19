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





func (e *Executor) execute(usdt float64) {
	// ===== LEG1 =====
	oid1 := uuid.NewString()
	ch1 := e.router.Register(oid1)
	log.Println("SEND LEG1")
	sendMarket(sym1, "buy", usdt, oid1)

	filledDash := <-ch1
	dash := roundDown(filledDash, stepDash)
	log.Println("LEG1 FILLED:", dash)

	// ===== LEG2 =====
	oid2 := uuid.NewString()
	ch2 := e.router.Register(oid2)
	log.Println("SEND LEG2")
	sendMarket(sym2, "sell", dash, oid2)

	filledBTC := <-ch2
	btc := roundDown(filledBTC, stepBTC) // WS уже возвращает с комиссией
	log.Println("LEG2 FILLED:", btc)

	// ===== LEG3 =====
	if btc < stepBTC {
		log.Println("LEG3 skipped, BTC меньше минимального шага")
		return
	}

	oid3 := uuid.NewString()
	ch3 := e.router.Register(oid3)
	log.Println("SEND LEG3")
	sendMarket(sym3, "sell", btc, oid3)

	// ждем fill через WS
	filledUSDT := <-ch3
	usdtFinal := roundDown(filledUSDT, stepUSDT)
	log.Println("LEG3 FILLED:", usdtFinal)
	log.Println("PNL:", usdtFinal-startUSDT)
}




go func() {
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Fatal(err)
		}

		// лог всех сообщений для дебага
		log.Println("WS MSG:", string(msg))

		if !bytes.Contains(msg, []byte("tradeOrders")) {
			continue
		}

		var evt struct {
			Topic string `json:"topic"`
			Data  struct {
				ClientOid   string `json:"clientOid"`
				FilledSize  string `json:"filledSize"`
				FilledFunds string `json:"filledFunds"`
			} `json:"data"`
		}

		if err := json.Unmarshal(msg, &evt); err != nil {
			log.Println("WS unmarshal error:", err)
			continue
		}

		if evt.Data.ClientOid == "" {
			continue
		}

		var filled float64
		if evt.Data.FilledFunds != "" {
			filled, _ = strconv.ParseFloat(evt.Data.FilledFunds, 64)
		} else if evt.Data.FilledSize != "" {
			filled, _ = strconv.ParseFloat(evt.Data.FilledSize, 64)
		}

		if filled > 0 {
			log.Printf("RESOLVE %s filled=%f\n", evt.Data.ClientOid, filled)
			router.Resolve(evt.Data.ClientOid, filled)
		}
	}
}()




2026/01/19 14:44:34.456953 START TRIANGLE 12
2026/01/19 14:44:35.960020 SEND LEG1
2026/01/19 14:44:35.960085 WS MSG: {"id":"17lZasYykka","type":"welcome"}
2026/01/19 14:44:36.472296 WS MSG: {"topic":"/spotMarket/tradeOrders","type":"message","subject":"orderChange","userId":"693a99012d117700012b08db","channelType":"private","data":{"clientOid":"673ace3e-4e31-4d91-bd10-5cdbad2c22fa","feeType":"takerFee","filledSize":"0.1481","funds":"12","liquidity":"taker","matchPrice":"80.99","matchSize":"0.1481","orderId":"696e1924b4b86900076f34dd","orderTime":1768823076278,"orderType":"market","pt":1768823076297,"remainFunds":"0.005381","side":"buy","status":"match","symbol":"DASH-USDT","tradeId":"19259671366621184","ts":1768823076296000000,"type":"match"}}
2026/01/19 14:44:36.472542 RESOLVE 673ace3e-4e31-4d91-bd10-5cdbad2c22fa filled=0.148100
2026/01/19 14:44:36.472573 WS MSG: {"topic":"/spotMarket/tradeOrders","type":"message","subject":"orderChange","userId":"693a99012d117700012b08db","channelType":"private","data":{"clientOid":"673ace3e-4e31-4d91-bd10-5cdbad2c22fa","filledSize":"0.1481","funds":"12","orderId":"696e1924b4b86900076f34dd","orderTime":1768823076278,"orderType":"market","pt":1768823076297,"remainFunds":"0","remainSize":"0","side":"buy","status":"done","symbol":"DASH-USDT","ts":1768823076296000000,"type":"canceled"}}
2026/01/19 14:44:36.472598 RESOLVE 673ace3e-4e31-4d91-bd10-5cdbad2c22fa filled=0.148100
2026/01/19 14:44:36.472613 LEG1 FILLED: 0.1481
2026/01/19 14:44:36.472630 SEND LEG2
2026/01/19 14:44:36.881906 WS MSG: {"topic":"/spotMarket/tradeOrders","type":"message","subject":"orderChange","userId":"693a99012d117700012b08db","channelType":"private","data":{"clientOid":"078db25f-c00a-4399-8485-415e9c6b09a9","feeType":"takerFee","filledSize":"0.1481","liquidity":"taker","matchPrice":"0.0008676","matchSize":"0.1481","orderId":"696e19248ae5f90007130e26","orderTime":1768823076637,"orderType":"market","pt":1768823076643,"remainSize":"0","side":"sell","size":"0.1481","status":"match","symbol":"DASH-BTC","tradeId":"16389436815462400","ts":1768823076642000000,"type":"match"}}
2026/01/19 14:44:36.881984 RESOLVE 078db25f-c00a-4399-8485-415e9c6b09a9 filled=0.148100
2026/01/19 14:44:36.882010 WS MSG: {"topic":"/spotMarket/tradeOrders","type":"message","subject":"orderChange","userId":"693a99012d117700012b08db","channelType":"private","data":{"clientOid":"078db25f-c00a-4399-8485-415e9c6b09a9","filledSize":"0.1481","orderId":"696e19248ae5f90007130e26","orderTime":1768823076637,"orderType":"market","pt":1768823076643,"remainFunds":"0","remainSize":"0","side":"sell","size":"0.1481","status":"done","symbol":"DASH-BTC","ts":1768823076642000000,"type":"filled"}}
2026/01/19 14:44:36.882033 RESOLVE 078db25f-c00a-4399-8485-415e9c6b09a9 filled=0.148100
2026/01/19 14:44:36.882074 LEG2 FILLED: 0.1481
2026/01/19 14:44:36.882095 SEND LEG3
2026/01/19 14:45:01.254483 WS MSG: {"type":"pong","timestamp":1768823101106020}



