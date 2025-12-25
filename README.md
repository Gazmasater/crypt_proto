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


--- FAIL: TestNormalizeSymbol_Full (0.00s)
    market_test.go:36: NormalizeSymbol("ETHBTC") = "", want "ETHBTC"
    market_test.go:36: NormalizeSymbol("ethbtc") = "", want "ETHBTC"
    market_test.go:36: NormalizeSymbol("XYZABC") = "", want "XYZABC"
--- FAIL: TestKey_Full (0.00s)
    market_test.go:88: Key("MEXC", "ethbtc") = "MEXC:", want "MEXC:ETHBTC"
    market_test.go:88: Key("KuCoin", "XYZABC") = "KuCoin:", want "KuCoin:XYZABC"
FAIL
FAIL    crypt_proto/internal/market     0.003s
FAIL
gaz358@gaz358-BOD-WXX9:~/myprog/crypt_proto/internal/market$ 


