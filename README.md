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
	"sync"
	"time"

	"crypt_proto/pkg/models"

	"github.com/gorilla/websocket"
)

type KuCoinCollector struct {
	ctx        context.Context
	cancel     context.CancelFunc
	conn       *websocket.Conn
	wsURL      string
	symbols    []string
	whitelist  map[string]struct{}
	lastSeq    map[string]string
	lastData   map[string]*models.MarketData
	mu         sync.Mutex
	out        chan<- *models.MarketData
}

func NewKuCoinCollector(ctx context.Context, wsURL string, symbols []string, out chan<- *models.MarketData) *KuCoinCollector {
	c, cancel := context.WithCancel(ctx)
	return &KuCoinCollector{
		ctx:       c,
		cancel:    cancel,
		wsURL:     wsURL,
		symbols:   symbols,
		whitelist: make(map[string]struct{}),
		lastSeq:   make(map[string]string),
		lastData:  make(map[string]*models.MarketData),
		out:       out,
	}
}

func (k *KuCoinCollector) Connect() error {
	var err error
	k.conn, _, err = websocket.DefaultDialer.Dial(k.wsURL, nil)
	if err != nil {
		return err
	}
	log.Println("[KuCoin WS] connected")
	return nil
}

func (k *KuCoinCollector) Start() {
	go k.readLoop()
	k.subscribeSymbols()
	go k.logLoop()
}

func (k *KuCoinCollector) Stop() error {
	k.cancel()
	if k.conn != nil {
		return k.conn.Close()
	}
	return nil
}

func (k *KuCoinCollector) subscribeSymbols() {
	for _, sym := range k.symbols {
		msg := map[string]interface{}{
			"id":             "sub_" + sym,
			"type":           "subscribe",
			"topic":          "/market/ticker:" + sym,
			"privateChannel": false,
			"response":       true,
		}
		_ = k.conn.WriteJSON(msg)
	}
}

func (k *KuCoinCollector) readLoop() {
	for {
		select {
		case <-k.ctx.Done():
			return
		default:
			var msg map[string]json.RawMessage
			if err := k.conn.ReadJSON(&msg); err != nil {
				log.Println("[KuCoin WS ERROR]", err)
				time.Sleep(time.Second)
				continue
			}

			var topic string
			if t, ok := msg["topic"]; ok {
				_ = json.Unmarshal(t, &topic)
			}
			if topic == "" {
				continue
			}

			var data struct {
				Price    string `json:"price"`
				Size     string `json:"size"`
				Sequence string `json:"sequence"`
				Time     int64  `json:"time"`
			}
			if d, ok := msg["data"]; ok {
				_ = json.Unmarshal(d, &data)
			} else {
				continue
			}

			k.mu.Lock()
			// фильтруем повторные сообщения
			if seq, ok := k.lastSeq[topic]; ok && seq == data.Sequence {
				k.mu.Unlock()
				continue
			}
			k.lastSeq[topic] = data.Sequence

			// фильтр whitelist
			if len(k.whitelist) > 0 {
				if _, ok := k.whitelist[topic]; !ok {
					k.mu.Unlock()
					continue
				}
			}

			md := &models.MarketData{
				Symbol: topic,
				Price:  data.Price,
				Volume: data.Size,
				Time:   data.Time,
			}
			k.lastData[topic] = md
			k.mu.Unlock()

			// отправляем свежие данные наружу
			k.out <- md
		}
	}
}

// logLoop выводит актуальные данные каждые 1 сек
func (k *KuCoinCollector) logLoop() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-k.ctx.Done():
			return
		case <-ticker.C:
			k.mu.Lock()
			for sym, md := range k.lastData {
				log.Printf("[KuCoin] %s: Price=%s Volume=%s Time=%d\n",
					sym, md.Price, md.Volume, md.Time)
			}
			k.mu.Unlock()
		}
	}
}











