Если оставить только нужное:

p99 execution latency
Micro-volatility (100 мс)
Fill ratio
Capture rate
Inventory drift




Название API
9623527002

696935c42a6dcd00013273f2
b348b686-55ff-4290-897b-02d55f815f65




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




go run -race main.go


GOMAXPROCS=8 go run -race main.go


case <-log5mTimer.C:
    log5mTimer.Reset(time.Until(nextLogTimeUTC(time.Now().UTC())))

    // 5m лог делаем СТРОГО по закрытым данным, без freshness/hold/throttle
    st, ok := computeClosed()
    if !ok {
        // чтобы ты видел, что таймер живой, можно печатать пропуск
        fmt.Println("[5m log] skip: not enough closed aligned bars yet")
        // (опционально) можно логировать SKIP в jsonl
        // writeLog(logg, "LOG_5M_SKIP", trade.Pos.String(), Stats{Mode:"skip"}, nil)
        continue
    }

    writeLog(logg, "LOG_5M", trade.Pos.String(), st, nil)
    fmt.Printf("[5m log] beta=%.4f spread=%.6f z=%+.3f mode=%s pos=%s\n",
        st.Beta, st.Spread, st.Z, st.Mode, trade.Pos.String())


