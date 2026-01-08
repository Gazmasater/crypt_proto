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




package store

import (
	"sync/atomic"
	"time"

	"crypt_proto/pkg/models"
)

type Quote struct {
	Bid       float64
	Ask       float64
	BidSize   float64
	AskSize   float64
	Timestamp int64
}

type MemoryStore struct {
	data atomic.Value // map[string]Quote
}

func NewMemoryStore() *MemoryStore {
	s := &MemoryStore{}
	s.data.Store(make(map[string]Quote))
	return s
}


func (s *MemoryStore) Run(in <-chan *models.MarketData) {
	for md := range in {
		old := s.data.Load().(map[string]Quote)

		// copy-on-write (дёшево, т.к. map маленькая)
		next := make(map[string]Quote, len(old)+1)
		for k, v := range old {
			next[k] = v
		}

		next[md.Exchange+"|"+md.Symbol] = Quote{
			Bid:       md.Bid,
			Ask:       md.Ask,
			BidSize:   md.BidSize,
			AskSize:   md.AskSize,
			Timestamp: md.Timestamp,
		}

		s.data.Store(next)
	}
}


func (s *MemoryStore) Get(exchange, symbol string) (Quote, bool) {
	m := s.data.Load().(map[string]Quote)
	q, ok := m[exchange+"|"+symbol]
	return q, ok
}

func (s *MemoryStore) Snapshot() map[string]Quote {
	return s.data.Load().(map[string]Quote)
}




В handle():

c.out <- &models.MarketData{
	Exchange:  "KuCoin",
	Symbol:    symbol,
	Bid:       bid,
	Ask:       ask,
	BidSize:   bidSize,
	AskSize:   askSize,
	Timestamp: time.Now().UnixMilli(),
}






