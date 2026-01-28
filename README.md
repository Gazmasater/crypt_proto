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



type kucoinWS struct {
	id      int
	conn    *websocket.Conn
	symbols []string

	last map[string]Last
}

type Last struct {
	Bid float64
	Ask float64
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



         0      140ms (flat, cum) 12.84% of Total
         .          .    177:func (ws *kucoinWS) handle(c *KuCoinCollector, msg []byte) {
         .       70ms    178:   if gjson.GetBytes(msg, "type").String() != "message" {
         .          .    179:           return
         .          .    180:   }
         .          .    181:
         .       10ms    182:   topic := gjson.GetBytes(msg, "topic").String()
         .          .    183:   const prefix = "/market/ticker:"
         .          .    184:   if len(topic) <= len(prefix) || topic[:len(prefix)] != prefix {
         .          .    185:           return
         .          .    186:   }
         .          .    187:   symbol := topic[len(prefix):]
         .          .    188:
         .       30ms    189:   bid := gjson.GetBytes(msg, "data.bestBid").Float()
         .          .    190:   ask := gjson.GetBytes(msg, "data.bestAsk").Float()
         .          .    191:   if bid == 0 || ask == 0 {
         .          .    192:           return
         .          .    193:   }
         .          .    194:
         .          .    195:   last := ws.last[symbol]
         .          .    196:   if last.Bid == bid && last.Ask == ask {
         .          .    197:           return
         .          .    198:   }
         .          .    199:
         .          .    200:   // если реально нужны
         .       10ms    201:   bidSize := gjson.GetBytes(msg, "data.bestBidSize").Float()
         .       10ms    202:   askSize := gjson.GetBytes(msg, "data.bestAskSize").Float()
         .          .    203:
         .          .    204:   ws.last[symbol] = Last{Bid: bid, Ask: ask}
         .          .    205:
         .       10ms    206:   c.out <- &models.MarketData{
         .          .    207:           Exchange: "KuCoin",
         .          .    208:           Symbol:   symbol,
         .          .    209:           Bid:      bid,
         .          .    210:           Ask:      ask,
         .          .    211:           BidSize:  bidSize,



