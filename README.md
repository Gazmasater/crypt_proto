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



type Last struct {
	Bid float64
	Ask float64
}

type kucoinWS struct {
	id   int
	conn *websocket.Conn
	last map[string]Last
}


ws := &kucoinWS{
	id:   i,
	conn: conn,
	last: make(map[string]Last),
}


func NewKuCoinCollectorFromCSV(path string) (*KuCoinCollector, []string, error) {
	symbols, err := readPairsFromCSV(path)
	if err != nil {
		return nil, nil, err
	}
	if len(symbols) == 0 {
		return nil, nil, fmt.Errorf("no symbols")
	}

	ctx, cancel := context.WithCancel(context.Background())

	var wsList []*kucoinWS
	for i := 0; i < len(symbols); i += maxSubsPerWS {
		end := i + maxSubsPerWS
		if end > len(symbols) {
			end = len(symbols)
		}

		wsList = append(wsList, &kucoinWS{
			id:      len(wsList),
			symbols: symbols[i:end],
			last:    make(map[string]Last),
		})
	}

	c := &KuCoinCollector{
		ctx:    ctx,
		cancel: cancel,
		wsList: wsList,
	}

	return c, symbols, nil
}


[{
	"resource": "/home/gaz358/myprog/crypt_proto/internal/collector/kucoin_collector.go",
	"owner": "_generated_diagnostic_collection_name_#0",
	"code": {
		"value": "MissingLitField",
		"target": {
			"$mid": 1,
			"path": "/golang.org/x/tools/internal/typesinternal",
			"scheme": "https",
			"authority": "pkg.go.dev",
			"fragment": "MissingLitField"
		}
	},
	"severity": 8,
	"message": "unknown field symbols in struct literal of type kucoinWS",
	"source": "compiler",
	"startLineNumber": 64,
	"startColumn": 4,
	"endLineNumber": 64,
	"endColumn": 11,
	"modelVersionId": 3,
	"origin": "extHost1"
}]

[{
	"resource": "/home/gaz358/myprog/crypt_proto/internal/collector/kucoin_collector.go",
	"owner": "_generated_diagnostic_collection_name_#0",
	"code": {
		"value": "MissingFieldOrMethod",
		"target": {
			"$mid": 1,
			"path": "/golang.org/x/tools/internal/typesinternal",
			"scheme": "https",
			"authority": "pkg.go.dev",
			"fragment": "MissingFieldOrMethod"
		}
	},
	"severity": 8,
	"message": "ws.symbols undefined (type *kucoinWS has no field or method symbols)",
	"source": "compiler",
	"startLineNumber": 145,
	"startColumn": 23,
	"endLineNumber": 145,
	"endColumn": 30,
	"modelVersionId": 3,
	"origin": "extHost1"
}]




