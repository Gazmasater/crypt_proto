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
	"strings"
	"syscall"
)

func main() {
	// читаем EXCHANGE из .env или среды
	exchange := strings.ToLower(os.Getenv("EXCHANGE"))
	if exchange == "" {
		exchange = "mexc" // значение по умолчанию
	}

	log.Println("EXCHANGE:", exchange)

	var c collector.Collector

	// создаём нужный коллектор по EXCHANGE
	switch exchange {
	case "mexc":
		c = collector.NewMEXCCollector([]string{
			"BTCUSDT",
			"ETHUSDT",
			"ETHBTC",
		})
	case "okx":
		c = collector.NewOKXCollector()
	case "kucoin":
		c = collector.NewKuCoinCollector([]string{
			"BTC-USDT",
			"ETH-USDT",
			"ETH-BTC",
		})
	default:
		log.Fatal("unknown exchange: ", exchange)
	}

	log.Println("Starting collector:", c.Name())

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







