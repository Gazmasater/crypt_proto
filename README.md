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




1️⃣ Ошибка quantity scale is invalid — формат quantity

Файл: mexc/trader.go
Функция: PlaceMarket

Найди строку с quantity:

params.Set("quantity", fmt.Sprintf("%.8f", quantity))


Замени на:

qtyStr := strconv.FormatFloat(quantity, 'f', -1, 64)
params.Set("quantity", qtyStr)


И сверху в импорты добавь:

import (
    // ...
    "strconv"
    // ...
)


Это уберёт лишние нули (2.00000000 → 2), и биржа больше не будет ругаться на scale.

2️⃣ panic: send on closed channel — не закрываем канал под писателями

Файл: cmd/cryptarb/main.go

В самом конце main() у тебя сейчас что-то типа:

<-ctx.Done()
log.Println("shutting down...")

time.Sleep(200 * time.Millisecond)
close(events)
wg.Wait()
log.Println("bye")


Сделай так:

<-ctx.Done()
log.Println("shutting down...")

wg.Wait()
log.Println("bye")


Просто убери time.Sleep и close(events) — канал закрывать не нужно, процесс всё равно завершится, а паника пропадёт.

Сделай эти две правки, пересобери и запусти — дальше посмотрим, что ответит MEXC: если формат ок, ошибка сменится на что-нибудь про баланс/лимиты, а не про quantity scale.





