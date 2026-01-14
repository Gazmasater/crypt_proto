Название API
9623527002

6966b78122ca320001d2acae
fa1e37ae-21ff-4257-844d-3dcd21d26ccd





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



func ParseTrianglesFromCSV(path string) ([]*Triangle, error) {
    f, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer f.Close()

    rows, err := csv.NewReader(f).ReadAll()
    if err != nil {
        return nil, err
    }

    var res []*Triangle
    for _, row := range rows[1:] {
        if len(row) < 6 {
            continue
        }

        tri := &Triangle{
            A: strings.TrimSpace(row[0]),
            B: strings.TrimSpace(row[1]),
            C: strings.TrimSpace(row[2]),
        }

        legs := []string{row[3], row[4], row[5]}
        for i, leg := range legs {
            leg = strings.ToUpper(strings.TrimSpace(leg))
            parts := strings.Fields(leg)
            if len(parts) != 2 {
                continue
            }
            isBuy := parts[0] == "BUY"
            symbolParts := strings.Split(parts[1], "/")
            if len(symbolParts) != 2 {
                continue
            }
            key := "KuCoin|" + symbolParts[0] + "-" + symbolParts[1]
            tri.Legs[i] = LegIndex{Key: key, IsBuy: isBuy}
        }

        res = append(res, tri)
    }
    return res, nil
}

