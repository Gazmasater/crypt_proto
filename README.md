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






package collector

import (
	"context"
	"encoding/json"
	"log"
	"strings"
	"time"

	"crypt_proto/pkg/models"

	"github.com/gorilla/websocket"
)

const mexcWS = "wss://wbs-api.mexc.com/ws"

type MEXCCollector struct {
	ctx     context.Context
	cancel  context.CancelFunc
	symbols []string
}

func NewMEXCCollector(symbols []string) *MEXCCollector {
	ctx, cancel := context.WithCancel(context.Background())

	up := make([]string, 0, len(symbols))
	for _, s := range symbols {
		up = append(up, strings.ToUpper(s))
	}

	return &MEXCCollector{
		ctx:     ctx,
		cancel:  cancel,
		symbols: up,
	}
}

func (c *MEXCCollector) Name() string {
	return "mexc"
}

func (c *MEXCCollector) Start(out chan<- models.MarketData) error {
	conn, _, err := websocket.DefaultDialer.Dial(mexcWS, nil)
	if err != nil {
		return err
	}

	log.Println("[MEXC] connected")

	// --- subscribe ---
	params := make([]string, 0, len(c.symbols))
	for _, s := range c.symbols {
		params = append(params, "spot@public.bookTicker.batch@"+s)
	}

	sub := map[string]interface{}{
		"method": "SUBSCRIPTION",
		"params": params,
	}

	if err := conn.WriteJSON(sub); err != nil {
		return err
	}

	log.Println("[MEXC] subscribed:", params)

	// --- ping loop ---
	go func() {
		t := time.NewTicker(20 * time.Second)
		defer t.Stop()
		for {
			select {
			case <-c.ctx.Done():
				return
			case <-t.C:
				_ = conn.WriteJSON(map[string]string{
					"method": "PING",
				})
			}
		}
	}()

	// --- read loop ---
	go func() {
		defer conn.Close()
		for {
			select {
			case <-c.ctx.Done():
				return
			default:
				_, msg, err := conn.ReadMessage()
				if err != nil {
					log.Println("[MEXC] read error:", err)
					return
				}
				c.handleMessage(msg, out)
			}
		}
	}()

	return nil
}

func (c *MEXCCollector) Stop() {
	c.cancel()
}




[{
	"resource": "/home/gaz358/myprog/crypt_proto/cmd/arb/main.go",
	"owner": "_generated_diagnostic_collection_name_#0",
	"code": {
		"value": "InvalidIfaceAssign",
		"target": {
			"$mid": 1,
			"path": "/golang.org/x/tools/internal/typesinternal",
			"scheme": "https",
			"authority": "pkg.go.dev",
			"fragment": "InvalidIfaceAssign"
		}
	},
	"severity": 8,
	"message": "cannot use collector.NewMEXCCollector(symbols) (value of type *collector.MEXCCollector) as collector.Collector value in assignment: *collector.MEXCCollector does not implement collector.Collector (wrong type for method Stop)\n\t\thave Stop()\n\t\twant Stop() error",
	"source": "compiler",
	"startLineNumber": 34,
	"startColumn": 7,
	"endLineNumber": 34,
	"endColumn": 42,
	"origin": "extHost1"
}]


[{
	"resource": "/home/gaz358/myprog/crypt_proto/internal/collector/mexc_collecter.go",
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
	"message": "c.handleMessage undefined (type *MEXCCollector has no field or method handleMessage)",
	"source": "compiler",
	"startLineNumber": 95,
	"startColumn": 7,
	"endLineNumber": 95,
	"endColumn": 20,
	"origin": "extHost1"
}]
