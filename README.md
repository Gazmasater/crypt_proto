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




func main() {

	go func() {
		log.Println("pprof on http://localhost:6060/debug/pprof/")
		_ = http.ListenAndServe("localhost:6060", nil)
	}()

	out := make(chan *models.MarketData, 100_000)

	// === IN-MEMORY STORE ===
	mem := store.NewMemoryStore()
	go mem.Run(out)

	kc, err := collector.NewKuCoinCollectorFromCSV("../exchange/data/kucoin_triangles_usdt.csv")
	if err != nil {
		log.Fatal(err)
	}

	if err := kc.Start(out); err != nil {
		log.Fatal(err)
	}

	log.Println("[Main] KuCoinCollector started")

	// === пример reader-а (НЕ логировать всё!) ===
	go func() {
		for {
			snap := mem.Snapshot()
			log.Printf("[Store] quotes=%d", len(snap))
			time.Sleep(5 * time.Second)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	kc.Stop()
	close(out)
}


gaz358@gaz358-BOD-WXX9:~/myprog/crypt_proto$    go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
Fetching profile over HTTP from http://localhost:6060/debug/pprof/profile?seconds=30
Saved profile in /home/gaz358/pprof/pprof.arb.samples.cpu.036.pb.gz
File: arb
Build ID: d5397b492681dad75216f5b44ef5c4498887c007
Type: cpu
Time: 2026-01-08 10:26:52 MSK
Duration: 30.07s, Total samples = 2.66s ( 8.85%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 1460ms, 54.89% of 2660ms total
Dropped 104 nodes (cum <= 13.30ms)
Showing top 10 nodes out of 151
      flat  flat%   sum%        cum   cum%
     900ms 33.83% 33.83%      900ms 33.83%  internal/runtime/syscall.Syscall6
     180ms  6.77% 40.60%      180ms  6.77%  runtime.futex
      60ms  2.26% 42.86%      160ms  6.02%  github.com/tidwall/gjson.parseObject
      60ms  2.26% 45.11%       60ms  2.26%  runtime.(*mspan).base (inline)
      50ms  1.88% 46.99%      110ms  4.14%  runtime.mapassign_faststr
      50ms  1.88% 48.87%      230ms  8.65%  runtime.scanobject
      40ms  1.50% 50.38%       40ms  1.50%  aeshashbody
      40ms  1.50% 51.88%       40ms  1.50%  github.com/tidwall/gjson.parseSquash
      40ms  1.50% 53.38%       50ms  1.88%  runtime.(*mcache).prepareForSweep
      40ms  1.50% 54.89%       40ms  1.50%  runtime.memclrNoHeapPointers
(pprof) 



type MemoryStore struct {
	m sync.Map // key -> Quote
}

func (s *MemoryStore) Run(in <-chan *models.MarketData) {
	for md := range in {
		s.m.Store(md.Exchange+"|"+md.Symbol, Quote{
			Bid: md.Bid, Ask: md.Ask,
			BidSize: md.BidSize, AskSize: md.AskSize,
			Timestamp: md.Timestamp,
		})
	}
}

func (s *MemoryStore) Len() int {
	n := 0
	s.m.Range(func(_, _ any) bool {
		n++
		return true
	})
	return n
}





