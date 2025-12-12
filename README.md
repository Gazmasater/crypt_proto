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








Да, я посмотрел твой код — причина, почему debug не работает, довольно простая:

У тебя сейчас два разных флага debug:

В пакете config

В пакете mexc

В mexc/ws.go есть свой var debug bool, но нигде не устанавливается в true, поэтому:

func handleTextMessage(connID int, raw []byte) {
    if !debug {      // <- всегда false
        return
    }
    ...
}


И даже если в .env ты ставишь DEBUG=true, это влияет только на config.debug, а до mexc.debug это не доходит.

Давай сделаем один общий флаг debug (в пакете config) и всё к нему привяжем.

1. Пакет config — ОСТАВИТЬ как есть

У тебя уже нормально:

package config

import (
    "log"
    "os"
    "strconv"
    "strings"

    "github.com/joho/godotenv"
)

type Config struct {
    ...
    Debug bool
}

var debug bool

func SetDebug(v bool) {
    debug = v
}

func Dlog(format string, args ...any) {
    if !debug {
        return
    }
    log.Printf(format, args...)
}


LoadConfig() читает DEBUG из .env и кладёт в cfg.Debug — это ок.

2. cmd/cryptarb/main.go — включаем глобальный debug

В начале main() обязательно после cfg := config.LoadConfig() добавь:

package main

import (
    ...
    "crypt_proto/config"
    "crypt_proto/mexc"
)

func main() {
    // pprof и т.п...

    cfg := config.LoadConfig()

    // включаем глобальный debug-флаг
    config.SetDebug(cfg.Debug)

    // (если хочешь, можно пробросить в mexc, но после правки ниже это уже не нужно)
    // mexc.SetDebug(cfg.Debug)

    ...
}


И не нужно больше никаких глобальных debug в main-пакете.

3. mexc/ws.go — убираем свой debug

Сделай так:

Было (упрощённо)
package mexc

import (
    "context"
    "crypt_proto/config"
    "crypt_proto/domain"
    "encoding/json"
    "log"
    "sync"
    "time"

    "github.com/gorilla/websocket"
)

var debug bool

...

func handleTextMessage(connID int, raw []byte) {
    if !debug {
        return
    }
    var tmp any
    if err := json.Unmarshal(raw, &tmp); err == nil {
        j, _ := json.Marshal(tmp)
        config.Dlog("[WS #%d TEXT] %s", connID, string(j))
    } else {
        config.Dlog("[WS #%d TEXT RAW] %s", connID, string(raw))
    }
}

Нужно сделать так

Удалить строку:

var debug bool


И в handleTextMessage убрать проверку локального флага:

func handleTextMessage(connID int, raw []byte) {
    // локальную проверку снимаем, т.к. Dlog сам смотрит на config.debug
    var tmp any
    if err := json.Unmarshal(raw, &tmp); err == nil {
        j, _ := json.Marshal(tmp)
        config.Dlog("[WS #%d TEXT] %s", connID, string(j))
    } else {
        config.Dlog("[WS #%d TEXT RAW] %s", connID, string(raw))
    }
}


Теперь логика такая:

.env → DEBUG=true

LoadConfig() читает это в cfg.Debug

main вызывает config.SetDebug(cfg.Debug)

config.Dlog(...) в любом пакете печатает логи только если DEBUG=true

4. Как проверить, что всё работает

Убедись, что в .env есть строка:

DEBUG=true


Пересобери/перезапусти:

go run ./cmd/cryptarb


При подключении к MEXC ты должен начать видеть в логах JSON-ответы / ACK’и от TEXT сообщений вебсокета типа:

[WS #0 TEXT] {"method":"SUBSCRIPTION","code":"0",...}


Если поставишь DEBUG=false (или уберёшь переменную из .env), эти строки пропадут, а обычные INFO-логи останутся.

Если хочешь, дальше можем:

Развести два уровня логов: DEBUG и, например, TRACE (для прямого дампа сырого protobuf/JSON).

Сделать флаг включения логов по бирже: отдельно для MEXC, отдельно для KuCoin/OKX.





