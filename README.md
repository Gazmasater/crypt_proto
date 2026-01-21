Если оставить только нужное:

p99 execution latency
Micro-volatility (100 мс)
Fill ratio
Capture rate
Inventory drift




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
	// учитываем комиссию takerFee 0.1% при продаже BTC
	btc := roundDown(filledBTC*(1-0.001), stepBTC)
	log.Println("LEG2 FILLED (after fee):", btc)

	// ===== LEG3 =====
	if btc < stepBTC {
		log.Println("LEG3 skipped, BTC меньше минимального шага после комиссии")
		return
	}

	oid3 := uuid.NewString()
	ch3 := e.router.Register(oid3)
	log.Println("SEND LEG3")
	sendMarket(sym3, "sell", btc, oid3)

	filledUSDT := <-ch3
	usdtFinal := roundDown(filledUSDT, stepUSDT)
	log.Println("LEG3 FILLED:", usdtFinal)
	log.Println("PNL:", usdtFinal-startUSDT)
}


2026/01/19 15:00:40.355440 START TRIANGLE 12
2026/01/19 15:00:42.228118 SEND LEG1
2026/01/19 15:00:42.228270 WS MSG: {"id":"17lakE8wCgq","type":"welcome"}
2026/01/19 15:00:42.738053 WS MSG: {"topic":"/spotMarket/tradeOrders","type":"message","subject":"orderChange","userId":"693a99012d117700012b08db","channelType":"private","data":{"clientOid":"644c71f0-e2ff-4214-b70d-3174a2ba7366","feeType":"takerFee","filledSize":"0.1493","funds":"12","liquidity":"taker","matchPrice":"80.34","matchSize":"0.1493","orderId":"696e1ceab4b86900077be638","orderTime":1768824042544,"orderType":"market","pt":1768824042549,"remainFunds":"0.005238","side":"buy","status":"match","symbol":"DASH-USDT","tradeId":"19259769015523328","ts":1768824042548000000,"type":"match"}}
2026/01/19 15:00:42.738229 RESOLVE 644c71f0-e2ff-4214-b70d-3174a2ba7366 filled=0.149300
2026/01/19 15:00:42.738265 WS MSG: {"topic":"/spotMarket/tradeOrders","type":"message","subject":"orderChange","userId":"693a99012d117700012b08db","channelType":"private","data":{"clientOid":"644c71f0-e2ff-4214-b70d-3174a2ba7366","filledSize":"0.1493","funds":"12","orderId":"696e1ceab4b86900077be638","orderTime":1768824042544,"orderType":"market","pt":1768824042549,"remainFunds":"0","remainSize":"0","side":"buy","status":"done","symbol":"DASH-USDT","ts":1768824042548000000,"type":"canceled"}}
2026/01/19 15:00:42.738292 RESOLVE 644c71f0-e2ff-4214-b70d-3174a2ba7366 filled=0.149300
2026/01/19 15:00:42.738309 LEG1 FILLED: 0.1492
2026/01/19 15:00:42.738324 SEND LEG2
2026/01/19 15:00:43.147955 WS MSG: {"topic":"/spotMarket/tradeOrders","type":"message","subject":"orderChange","userId":"693a99012d117700012b08db","channelType":"private","data":{"clientOid":"4171e0db-d076-4801-b09a-ad84f15e8d8b","feeType":"takerFee","filledSize":"0.1492","liquidity":"taker","matchPrice":"0.0008613","matchSize":"0.1492","orderId":"696e1ceae3f47900079cd915","orderTime":1768824042918,"orderType":"market","pt":1768824042924,"remainSize":"0","side":"sell","size":"0.1492","status":"match","symbol":"DASH-BTC","tradeId":"16389537432358912","ts":1768824042923000000,"type":"match"}}
2026/01/19 15:00:43.148060 RESOLVE 4171e0db-d076-4801-b09a-ad84f15e8d8b filled=0.149200
2026/01/19 15:00:43.148084 WS MSG: {"topic":"/spotMarket/tradeOrders","type":"message","subject":"orderChange","userId":"693a99012d117700012b08db","channelType":"private","data":{"clientOid":"4171e0db-d076-4801-b09a-ad84f15e8d8b","filledSize":"0.1492","orderId":"696e1ceae3f47900079cd915","orderTime":1768824042918,"orderType":"market","pt":1768824042924,"remainFunds":"0","remainSize":"0","side":"sell","size":"0.1492","status":"done","symbol":"DASH-BTC","ts":1768824042923000000,"type":"filled"}}
2026/01/19 15:00:43.148109 RESOLVE 4171e0db-d076-4801-b09a-ad84f15e8d8b filled=0.149200
2026/01/19 15:00:43.148127 LEG2 FILLED (after fee): 0.14905000000000002
2026/01/19 15:00:43.148142 SEND LEG3
2026/01/19 15:01:07.621839 WS MSG: {"type":"pong","timestamp":1768824067407080}







