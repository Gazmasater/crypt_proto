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




1) Исправление main.go (чтобы компилилось при любой сигнатуре)
Заменяешь создание real executor вот так:
exec := arb.NewRealExecutor(tr, arbOut, cfg.TradeAmountUSDT) // <-- 3-й аргумент float64 (как у тебя сейчас)
exec.Filters = filters                               // <-- filters кладём в поле
exec.Cooldown = time.Duration(cfg.TradeCooldownMs) * time.Millisecond
consumer.Executor = exec

То есть filters больше не передаём в NewRealExecutor, а задаём после создания.

2) Тогда executor_real.go должен иметь поле Filters
В твоём executor_real.go должно быть:
type RealExecutor struct {
    trader SpotTrader
    out    io.Writer

    Filters map[string]SymbolFilter // <-- сюда кладём filters из main.go
    ...
}

А конструктор обязательно такой (3-й аргумент float64):
func NewRealExecutor(tr SpotTrader, out io.Writer, startUSDT float64) *RealExecutor {
    ...
}

⚠️ Это прямо соответствует твоей ошибке: компилятор хочет float64 третьим параметром.

3) Staticcheck SA6005 (EqualFold)
Это не ошибка компиляции, но поправь быстро:
Было:
if strings.ToUpper(a) != strings.ToUpper(b) { ... }

Стало:
if !strings.EqualFold(a, b) { ... }


Полный патч executor_real.go (только ключевые куски)
Вот минимально-правильные изменения, чтобы и компилилось, и filters работали:
// constructor: third arg is float64 (startUSDT)
func NewRealExecutor(tr SpotTrader, out io.Writer, startUSDT float64) *RealExecutor {
    if startUSDT <= 0 {
        startUSDT = 2
    }
    return &RealExecutor{
        trader:    tr,
        out:       out,
        StartUSDT: startUSDT,
        Cooldown:  300 * time.Millisecond,
        safety:    0.995,
        lastExec:  make(map[string]time.Time),
        Filters:   make(map[string]SymbolFilter),
    }
}

И в местах проверок:
if !strings.EqualFold(leg1.From, "USDT") { ... }
...
if !strings.EqualFold(leg3.To, "USDT") { ... }


Что сделать прямо сейчас (коротко по шагам)


В cmd/cryptarb/main.go перепиши строку, где ты вызываешь arb.NewRealExecutor(...), чтобы третьим параметром был float64, а filters присваивались после.


В arb/executor_real.go убедись, что:


NewRealExecutor(tr, out, startUSDT float64) именно так


есть поле Filters map[string]SymbolFilter




SA6005 — поменяй сравнения на strings.EqualFold.



Если ты скинешь текущую сигнатуру из твоего arb/executor_real.go (строку func NewRealExecutor(...)), я дам тебе ровно один готовый кусок main.go под неё без вариантов.













