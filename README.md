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






package main

import (
	"crypt_proto/internal/collector"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.Println("EXCHANGE: mexc")
	log.Println("Starting collector: mexc")
	var c collector.Collector

	c = collector.NewMEXCCollector([]string{
		"BTCUSDT",
		"ETHUSDT",
		"ETHBTC",
	})

	if err := c.Start(); err != nil {
		log.Fatal(err)
	}

	// корректное завершение по Ctrl+C
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	log.Println("Stopping collector...")
	c.Stop()
}



[{
	"resource": "/home/gaz358/myprog/crypt_proto/cmd/arb/main.go",
	"owner": "_generated_diagnostic_collection_name_#0",
	"code": {
		"value": "InvalidIfaceAssign",
		"target": {
			"$mid": 1,
			"path": "/golang.org/x/tools/internal/typesinternal",
			"scheme": "https",
			"authority": "pkg.go.dev",
			"fragment": "InvalidIfaceAssign"
		}
	},
	"severity": 8,
	"message": "cannot use collector.NewMEXCCollector([]string{…}) (value of type *collector.MEXCCollector) as collector.Collector value in assignment: *collector.MEXCCollector does not implement collector.Collector (missing method Name)",
	"source": "compiler",
	"startLineNumber": 16,
	"startColumn": 6,
	"endLineNumber": 20,
	"endColumn": 4,
	"origin": "extHost1"
}]


[{
	"resource": "/home/gaz358/myprog/crypt_proto/cmd/arb/main.go",
	"owner": "_generated_diagnostic_collection_name_#0",
	"code": {
		"value": "WrongArgCount",
		"target": {
			"$mid": 1,
			"path": "/golang.org/x/tools/internal/typesinternal",
			"scheme": "https",
			"authority": "pkg.go.dev",
			"fragment": "WrongArgCount"
		}
	},
	"severity": 8,
	"message": "not enough arguments in call to c.Start\n\thave ()\n\twant (chan<- models.MarketData)",
	"source": "compiler",
	"startLineNumber": 22,
	"startColumn": 20,
	"endLineNumber": 22,
	"endColumn": 20,
	"origin": "extHost1"
}]









