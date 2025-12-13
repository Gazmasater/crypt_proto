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



1) Проверь формат топиков для bookTicker (он должен быть именно такой)

Для “Individual Symbol Book Ticker Streams” правильный параметр:

spot@public.aggre.bookTicker.v3.api.pb@100ms@BTCUSDT
или @10ms@ вместо @100ms@. 
MEXC

То есть генератор топиков должен быть примерно так:

topic := fmt.Sprintf("spot@public.aggre.bookTicker.v3.api.pb@100ms@%s", symbol) // symbol в UPPERCASE


Если ты сейчас подписываешься на старые/другие строки (например от wbs.mexc.com/ws), сервер их не признаёт → и через 30 сек режет.

Критично: символы должны быть UPPERCASE. 
MEXC

2) Логируй ACK подписки (это моментально покажет “валидна или нет”)

MEXC на успешную подписку отвечает JSON’ом вида:

{"id":0,"code":0,"msg":"<topic>"}
``` :contentReference[oaicite:3]{index=3}

Сделай так: после `WriteJSON(SUB)` в течение 2–3 секунд жди хотя бы один ACK и логируй `code/msg`.
- если ACK нет или `code != 0` → topic неправильный (или запрещён), это твой кейс.

---

## 3) Добавь keepalive по их протоколу (PING JSON), иначе будут отваливаться “тихие” пары
У них ping/pong именно в виде сообщения:

- запрос: `{"method":"PING"}`
- ответ: `{"id":0,"code":0,"msg":"PONG"}` :contentReference[oaicite:4]{index=4}

Иначе для пар, где минуту нет апдейтов, сервер может отключать через 1 минуту. :contentReference[oaicite:5]{index=5}

---

## 4) Лимит подписок ты уже учёл, но напомню
Один WS коннект поддерживает **максимум 30 подписок**. Ты уже сделал 25 — это ок. :contentReference[oaicite:6]{index=6}

---

### Что сделать прямо сейчас (быстрый чеклист)
1) Убедись, что **все 25 topics** в точности вида  
`spot@public.aggre.bookTicker.v3.api.pb@100ms@SYMBOL` :contentReference[oaicite:7]{index=7}  
2) После SUB дождись и залогируй **ACK** (`code==0`). :contentReference[oaicite:8]{index=8}  
3) Запусти goroutine с `{"method":"PING"}` раз в 15–20 секунд. :contentReference[oaicite:9]{index=9}  

Если хочешь — вставь сюда 5–10 строк кода, где ты **формируешь `params` (topics)** для подписки на bookTicker. По ним я сразу скажу, что именно у тебя не так (обычно там одна “лишняя/не та” часть строки).
::contentReference[oaicite:10]{index=10}










