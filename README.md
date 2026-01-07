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




import (
	"strings"
	"sync"

	"github.com/tidwall/gjson"
	"crypt_proto/pkg/models"
)

type kucoinWS struct {
	id      int
	conn    *websocket.Conn
	symbols []string
	last    map[string][2]float64
	mu      sync.Mutex
}

func (ws *kucoinWS) handle(c *KuCoinCollector, msg []byte) {
	// Быстрая проверка type
	if gjson.GetBytes(msg, "type").String() != "message" {
		return
	}

	// Быстрая проверка topic
	topic := gjson.GetBytes(msg, "topic").String()
	parts := strings.Split(topic, ":")
	if len(parts) != 2 {
		return
	}
	symbol := normalize(parts[1])

	// Извлечение цен через gjson
	data := gjson.GetBytes(msg, "data")
	bid := data.Get("bestBid").Float()
	ask := data.Get("bestAsk").Float()
	if bid == 0 || ask == 0 {
		return
	}

	// Проверяем, изменилась ли цена
	ws.mu.Lock()
	last := ws.last[symbol]
	if last[0] == bid && last[1] == ask {
		ws.mu.Unlock()
		return
	}
	ws.last[symbol] = [2]float64{bid, ask}
	ws.mu.Unlock()

	// Отправка данных дальше
	c.out <- &models.MarketData{
		Exchange: "KuCoin",
		Symbol:   symbol,
		Bid:      bid,
		Ask:      ask,
	}
}

// normalize оставляем прежним
func normalize(s string) string {
	parts := strings.Split(s, "-")
	return parts[0] + "/" + parts[1]
}



[{
	"resource": "/home/gaz358/myprog/crypt_proto/internal/collector/kucoin_collector.go",
	"owner": "_generated_diagnostic_collection_name_#0",
	"code": {
		"value": "IncompatibleAssign",
		"target": {
			"$mid": 1,
			"path": "/golang.org/x/tools/internal/typesinternal",
			"scheme": "https",
			"authority": "pkg.go.dev",
			"fragment": "IncompatibleAssign"
		}
	},
	"severity": 8,
	"message": "cannot use make(map[string]lastData) (value of type map[string]lastData) as map[string][2]float64 value in struct literal",
	"source": "compiler",
	"startLineNumber": 70,
	"startColumn": 13,
	"endLineNumber": 70,
	"endColumn": 38,
	"origin": "extHost1"
}]




