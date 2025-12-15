mx0vglmT3srN1IS19H
135bb7a7509e4421bad692415c53753b



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




Что сделать в arb/arb.go

Удали саму функцию:

func detectBaseQuote(symbol, a, b string) (base, quote string, ok bool) {
	if strings.HasPrefix(symbol, a) {
		return a, b, true
	}
	if strings.HasPrefix(symbol, b) {
		return b, a, true
	}
	return "", "", false
}


После этого, скорее всего, будет лишний импорт strings вверху файла.
В блоке import в arb.go убери "strings" если он больше нигде не нужен:

Было примерно так:

import (
    "bufio"
    "context"
    "fmt"
    "io"
    "log"
    "os"
    "strings"
    "sync"
    "time"

    "crypt_proto/domain"
)


Сделай так:

import (
    "bufio"
    "context"
    "fmt"
    "io"
    "log"
    "os"
    "sync"
    "time"

    "crypt_proto/domain"
)


После этого warning про detectBaseQuote исчезнет.

2️⃣ NewRealExecutor — не хватает аргумента

Сообщение:

not enough arguments in call to arb.NewRealExecutor
have (*mexc.Trader, io.Writer)
want (arb.SpotTrader, io.Writer, float64)


Это потому что мы добавили в RealExecutor третий параметр — fixedStartUSDT (твоя ручная сумма, типа 2 USDT), а в main.go всё ещё вызываем старую сигнатуру.

Что исправить в cmd/cryptarb/main.go

Найди кусок:

trader := mexc.NewTrader(cfg.APIKey, cfg.APISecret, cfg.Debug)
consumer.Executor = arb.NewRealExecutor(trader, arbOut)


и замени на:

trader := mexc.NewTrader(cfg.APIKey, cfg.APISecret, cfg.Debug)
consumer.Executor = arb.NewRealExecutor(trader, arbOut, cfg.TradeAmountUSDT)






