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








1. Пакет mexc: свой debug-флаг

В начале файла mexc/ws.go (после import) добавь:

package mexc

import (
    // ...
)

// глобальный флаг для отладочных логов в пакете mexc
var debug bool

// публичный сеттер, чтобы main мог включать/выключать
func SetDebug(v bool) {
    debug = v
}


И оставь твой dlog примерно таким:

func dlog(format string, args ...any) {
    if debug {
        log.Printf(format, args...)
    }
}


(если dlog уже есть – просто убедись, что он смотрит на этот debug.)

2. В cmd/cryptarb/main.go убрать старый debug

Скорее всего, у тебя там где-то есть строка вроде:

debug = cfg.Debug


и при этом переменная debug больше не объявлена в main.
Эту строку нужно заменить на вызов для mexc:

package main

import (
    // ...
    "crypt_proto/mexc"
)

func main() {
    cfg := config.Load()

    // включаем/выключаем отладку в пакете mexc
    mexc.SetDebug(cfg.Debug)

    // дальше твой код...
}


Если в main больше нигде debug не используется – просто больше его не трогаем, он теперь «живёт» внутри mexc.






