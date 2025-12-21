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




1️⃣ Интерфейс Collector

Судя по ошибкам, у тебя в коде есть что-то вроде:

var c collector.Collector


и ты хочешь присвоить туда:

c = collector.NewOKXCollector()
c = collector.NewMEXCCollector("BTCUSDT")


Ошибка:

MEXCCollector не реализует метод Stop() error, а только Stop().

Name() тоже не входит в интерфейс.

✅ Исправление: сделать интерфейс так:

package collector

import "crypt_proto/pkg/models"

type Collector interface {
	Start(out chan<- models.MarketData) error
	Stop() error
	Name() string
}

2️⃣ Реализация OKXCollector
func (c *OKXCollector) Stop() error {
	c.cancel()
	return nil
}

3️⃣ Реализация MEXCCollector
func (c *MEXCCollector) Stop() error {
	c.cancel()
	return nil
}


Теперь и OKX, и MEXC реализуют интерфейс Collector.

4️⃣ Конструкторы

Ошибка WrongArgCount говорит о том, что конструктор NewOKXCollector() в интерфейсе ожидает аргументы?

Сейчас у нас NewOKXCollector() без аргументов.

В main.go используй без аргументов:

c = collector.NewOKXCollector()


Для MEXC нужно указать символ:

c = collector.NewMEXCCollector("BTCUSDT")

5️⃣ main.go

После исправлений:

var c collector.Collector

switch exchange {
case "okx":
	c = collector.NewOKXCollector()
case "mexc":
	c = collector.NewMEXCCollector("BTCUSDT")
default:
	log.Fatalf("unknown exchange: %s", exchange)
}

log.Printf("Starting collector: %s\n", c.Name())


Теперь c.Name() и c.Stop() корректно работают через интерфейс.




