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



func (ws *kucoinWS) connect() error {
    req, _ := http.NewRequest("POST", "https://api.kucoin.com/api/v1/bullet-public", nil)
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    var r struct {
        Data struct {
            Token           string `json:"token"`
            InstanceServers []struct {
                Endpoint string `json:"endpoint"`
            } `json:"instanceServers"`
        } `json:"data"`
    }

    if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
        return err
    }

    url := fmt.Sprintf(
        "%s?token=%s&connectId=%d",
        r.Data.InstanceServers[0].Endpoint,
        r.Data.Token,
        time.Now().UnixNano(),
    )

    conn, _, err := websocket.DefaultDialer.Dial(url, nil)
    if err != nil {
        return err
    }

    ws.conn = conn
    log.Printf("[KuCoin WS %d] connected\n", ws.id)
    return nil
}



