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

	var filledDash float64
	select {
	case filledDash = <-ch1:
	case <-time.After(5 * time.Second):
		log.Println("LEG1 timeout!")
		return
	}
	time.Sleep(e.delayLeg1)
	dash := roundDown(filledDash, stepDash)
	log.Println("LEG1 FILLED:", dash)

	// ===== LEG2 =====
	oid2 := uuid.NewString()
	ch2 := e.router.Register(oid2)
	log.Println("SEND LEG2")
	sendMarket(sym2, "sell", dash, oid2)

	var filledBTC float64
	select {
	case filledBTC = <-ch2:
	case <-time.After(5 * time.Second):
		log.Println("LEG2 timeout!")
		return
	}
	time.Sleep(e.delayLeg2)
	btc := roundDown(filledBTC*fee, stepBTC)
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

	var filledUSDT float64
	select {
	case filledUSDT = <-ch3:
	case <-time.After(5 * time.Second):
		log.Println("LEG3 timeout!")
		return
	}
	time.Sleep(e.delayLeg3)
	usdtFinal := roundDown(filledUSDT, stepUSDT)
	log.Println("LEG3 FILLED:", usdtFinal)
	log.Println("PNL:", usdtFinal-startUSDT)
}


az358@gaz358-BOD-WXX9:~/myprog/crypt_proto/test$ go run .
2026/01/19 13:45:07.698820 START TRIANGLE 12
2026/01/19 13:45:09.899016 SEND LEG1
2026/01/19 13:45:11.609389 LEG1 FILLED: 0.1482
2026/01/19 13:45:11.609434 SEND LEG2
2026/01/19 13:45:13.049410 LEG2 FILLED: 0.06130000000000001
2026/01/19 13:45:13.049462 SEND LEG3
2026/01/19 13:45:18.050497 LEG3 timeout!





