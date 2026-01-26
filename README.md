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



package queue

import (
	"sync/atomic"
	"time"

	"crypt_proto/pkg/models"
)

type Quote struct {
	Bid, Ask         float64
	BidSize, AskSize float64
	Timestamp        int64
}

type MemoryStore struct {
	data  atomic.Value           // map[string]Quote
	batch chan *models.MarketData
}

// NewMemoryStore создаёт MemoryStore с буфером batch
func NewMemoryStore() *MemoryStore {
	s := &MemoryStore{
		batch: make(chan *models.MarketData, 10_000),
	}
	s.data.Store(make(map[string]Quote))
	return s
}

// Run обрабатывает входящие данные с минимальной задержкой
func (s *MemoryStore) Run() {
	const microBatch = 5 // применяем батч каждые 5 котировок

	batch := make([]*models.MarketData, 0, microBatch)

	for md := range s.batch {
		// фильтруем котировки с нулевым объемом
		if md.BidSize == 0 && md.AskSize == 0 {
			continue
		}

		batch = append(batch, md)

		// применяем батч
		if len(batch) >= microBatch {
			s.applyBatch(batch)
			batch = batch[:0]
		}
	}

	// оставшиеся котировки при завершении
	if len(batch) > 0 {
		s.applyBatch(batch)
	}
}

// Push добавляет MarketData в очередь
func (s *MemoryStore) Push(md *models.MarketData) {
	select {
	case s.batch <- md:
	default:
		// drop if full
	}
}

// Get возвращает snapshot lock-free
func (s *MemoryStore) Get(exchange, symbol string) (Quote, bool) {
	m := s.data.Load().(map[string]Quote)
	q, ok := m[exchange+"|"+symbol]
	return q, ok
}

// applyBatch применяет несколько котировок сразу
func (s *MemoryStore) applyBatch(batch []*models.MarketData) {
	oldMap := s.data.Load().(map[string]Quote)
	newMap := make(map[string]Quote, len(oldMap)+len(batch))

	// копируем старые данные
	for k, v := range oldMap {
		newMap[k] = v
	}

	// применяем новые котировки
	now := time.Now().UnixMilli()
	for _, md := range batch {
		key := md.Exchange + "|" + md.Symbol
		newMap[key] = Quote{
			Bid: md.Bid,
			Ask: md.Ask,
			BidSize: md.BidSize,
			AskSize: md.AskSize,
			Timestamp: now,
		}
	}

	s.data.Store(newMap)
}


gaz358@gaz358-BOD-WXX9:~/myprog/crypt_proto$    go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
Fetching profile over HTTP from http://localhost:6060/debug/pprof/profile?seconds=30
Saved profile in /home/gaz358/pprof/pprof.arb.samples.cpu.264.pb.gz
File: arb
Build ID: 92d6ab1c7fab5e78cad537019b70c12c249448f3
Type: cpu
Time: 2026-01-26 17:45:46 MSK
Duration: 30.01s, Total samples = 1.20s ( 4.00%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 820ms, 68.33% of 1200ms total
Showing top 10 nodes out of 143
      flat  flat%   sum%        cum   cum%
     490ms 40.83% 40.83%      490ms 40.83%  internal/runtime/syscall.Syscall6
     130ms 10.83% 51.67%      130ms 10.83%  runtime.futex
      40ms  3.33% 55.00%       40ms  3.33%  aeshashbody
      30ms  2.50% 57.50%       80ms  6.67%  github.com/tidwall/gjson.Get
      30ms  2.50% 60.00%       30ms  2.50%  runtime.nextFreeFast
      20ms  1.67% 61.67%      530ms 44.17%  github.com/gorilla/websocket.(*Conn).advanceFrame
      20ms  1.67% 63.33%       20ms  1.67%  github.com/tidwall/gjson.parseSquash
      20ms  1.67% 65.00%       20ms  1.67%  memeqbody
      20ms  1.67% 66.67%       20ms  1.67%  runtime.(*m).becomeSpinning (inline)
      20ms  1.67% 68.33%      290ms 24.17%  runtime.findRunnable
(pprof) gaz358@gaz358-BOD-WXX9:~/myprog/crypt_proto$    go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
Fetching profile over HTTP from http://localhost:6060/debug/pprof/profile?seconds=30
Saved profile in /home/gaz358/pprof/pprof.arb.samples.cpu.265.pb.gz
File: arb
Build ID: 0a7106a37f4339378b3764f072ac688e40be7fc5
Type: cpu
Time: 2026-01-26 17:49:23 MSK
Duration: 30s, Total samples = 1.08s ( 3.60%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 710ms, 65.74% of 1080ms total
Showing top 10 nodes out of 121
      flat  flat%   sum%        cum   cum%
     470ms 43.52% 43.52%      470ms 43.52%  internal/runtime/syscall.Syscall6
      60ms  5.56% 49.07%       60ms  5.56%  runtime.futex
      40ms  3.70% 52.78%       40ms  3.70%  github.com/tidwall/gjson.parseObjectPath
      20ms  1.85% 54.63%      490ms 45.37%  bufio.(*Reader).fill
      20ms  1.85% 56.48%       20ms  1.85%  runtime.(*timers).check
      20ms  1.85% 58.33%       60ms  5.56%  runtime.entersyscall_sysmon
      20ms  1.85% 60.19%      350ms 32.41%  runtime.findRunnable
      20ms  1.85% 62.04%       20ms  1.85%  runtime.nanotime (inline)
      20ms  1.85% 63.89%      160ms 14.81%  runtime.netpoll
      20ms  1.85% 65.74%       20ms  1.85%  runtime.nextFreeFast (inline)
(pprof) 




