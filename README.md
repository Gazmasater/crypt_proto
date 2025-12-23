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
	"fmt"
	"os"
	"strings"

	"crypt_proto/internal/collector"
	"crypt_proto/pkg/models"

	"github.com/joho/godotenv"
)

func main() {
	// Загружаем .env
	_ = godotenv.Load(".env")

	exchange := strings.ToLower(os.Getenv("EXCHANGE"))
	if exchange == "" {
		exchange = "mexc"
	}
	fmt.Println("EXCHANGE:", exchange)

	// Канал для получения рыночных данных
	marketDataCh := make(chan models.MarketData, 1000)

	var c collector.Collector
	switch exchange {
	case "mexc":
		c = collector.NewMEXCCollector([]string{"BTCUSDT", "ETHUSDT"})
	case "kucoin":
		c = collector.NewKuCoinCollector([]string{"BTC-USDT", "ETH-USDT"})
	case "okx":
		c = collector.NewOKXCollector() // у OKX нет аргументов
	default:
		panic("unknown exchange: " + exchange)
	}

	fmt.Println("Starting collector:", c.Name())
	if err := c.Start(marketDataCh); err != nil {
		panic(err)
	}
	defer c.Stop()

	// Consumer
	go func() {
		for data := range marketDataCh {
			fmt.Printf("[%s] %s bid=%.8f ask=%.8f\n",
				data.Exchange, data.Symbol, data.Bid, data.Ask)
		}
	}()

	select {}
}



gaz358@gaz358-BOD-WXX9:~/myprog/crypt_proto/cmd/arb$ go run .
EXCHANGE: kucoin
Starting collector: KuCoin
panic: dial tcp: lookup ws.kucoin.com on 127.0.0.53:53: no such host

goroutine 1 [running]:
main.main()
        /home/gaz358/myprog/crypt_proto/cmd/arb/main.go:41 +0x2c5
exit status 2







