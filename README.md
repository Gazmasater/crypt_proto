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




gaz358@gaz358-BOD-WXX9:~/myprog/crypt_proto$    go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
Fetching profile over HTTP from http://localhost:6060/debug/pprof/profile?seconds=30
Saved profile in /home/gaz358/pprof/pprof.arb.samples.cpu.007.pb.gz
File: arb
Build ID: b7f6cbe195780e80f45cf9c0dc233b7b7862e62c
Type: cpu
Time: 2025-12-27 02:35:03 MSK
Duration: 30s, Total samples = 490ms ( 1.63%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 280ms, 57.14% of 490ms total
Showing top 10 nodes out of 117
      flat  flat%   sum%        cum   cum%
     130ms 26.53% 26.53%      130ms 26.53%  internal/runtime/syscall.Syscall6
      30ms  6.12% 32.65%       30ms  6.12%  runtime.futex
      20ms  4.08% 36.73%       60ms 12.24%  crypt_proto/internal/collector.(*MEXCCollector).handleWrapper
      20ms  4.08% 40.82%       60ms 12.24%  google.golang.org/protobuf/internal/impl.(*MessageInfo).initOneofFieldCoders.func1
      20ms  4.08% 44.90%       20ms  4.08%  runtime.nanotime
      20ms  4.08% 48.98%      150ms 30.61%  syscall.Syscall
      10ms  2.04% 51.02%       10ms  2.04%  crypto/internal/fips140/aes/gcm.gcmAesDec
      10ms  2.04% 53.06%       10ms  2.04%  gogo
      10ms  2.04% 55.10%       90ms 18.37%  google.golang.org/protobuf/internal/impl.(*MessageInfo).unmarshalPointerEager
      10ms  2.04% 57.14%       10ms  2.04%  internal/abi.(*Type).NumMethod
(pprof) 


gaz358@gaz358-BOD-WXX9:~/myprog/crypt_proto$    go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
Fetching profile over HTTP from http://localhost:6060/debug/pprof/profile?seconds=30
Saved profile in /home/gaz358/pprof/pprof.arb.samples.cpu.013.pb.gz
File: arb
Build ID: cb388e72beffdffa723ddf9def6383fa67349251
Type: cpu
Time: 2025-12-27 04:41:58 MSK
Duration: 30.08s, Total samples = 440ms ( 1.46%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 290ms, 65.91% of 440ms total
Showing top 10 nodes out of 89
      flat  flat%   sum%        cum   cum%
     160ms 36.36% 36.36%      160ms 36.36%  internal/runtime/syscall.Syscall6
      30ms  6.82% 43.18%       30ms  6.82%  runtime.futex
      20ms  4.55% 47.73%       30ms  6.82%  crypto/internal/fips140/aes/gcm.(*GCMForTLS13).Open
      20ms  4.55% 52.27%      200ms 45.45%  crypto/tls.(*Conn).readRecordOrCCS
      10ms  2.27% 54.55%       10ms  2.27%  bytes.(*Reader).Len (inline)
      10ms  2.27% 56.82%       20ms  4.55%  crypt_proto/internal/market.NormalizeSymbol_Full
      10ms  2.27% 59.09%       10ms  2.27%  crypto/internal/fips140.RecordApproved
      10ms  2.27% 61.36%      220ms 50.00%  github.com/gorilla/websocket.(*Conn).read
      10ms  2.27% 63.64%       30ms  6.82%  google.golang.org/protobuf/internal/impl.(*MessageInfo).unmarshal
      10ms  2.27% 65.91%       20ms  4.55%  google.golang.org/protobuf/internal/impl.(*MessageInfo).unmarshalPointerEager



gaz358@gaz358-BOD-WXX9:~/myprog/crypt_proto$ go tool pprof http://localhost:6060/debug/pprof/heap
Fetching profile over HTTP from http://localhost:6060/debug/pprof/heap
Saved profile in /home/gaz358/pprof/pprof.arb.alloc_objects.alloc_space.inuse_objects.inuse_space.002.pb.gz
File: arb
Build ID: b7f6cbe195780e80f45cf9c0dc233b7b7862e62c
Type: inuse_space
Time: 2025-12-27 02:36:50 MSK
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 2563.10kB, 100% of 2563.10kB total
Showing top 10 nodes out of 13
      flat  flat%   sum%        cum   cum%
    1539kB 60.04% 60.04%     1539kB 60.04%  runtime.allocm
  512.05kB 19.98% 80.02%   512.05kB 19.98%  runtime.main
  512.05kB 19.98%   100%   512.05kB 19.98%  runtime.acquireSudog
         0     0%   100%   512.05kB 19.98%  runtime.chanrecv
         0     0%   100%   512.05kB 19.98%  runtime.chanrecv1
         0     0%   100%     1539kB 60.04%  runtime.mcall
         0     0%   100%     1539kB 60.04%  runtime.newm
         0     0%   100%     1539kB 60.04%  runtime.park_m
         0     0%   100%     1539kB 60.04%  runtime.resetspinning
         0     0%   100%     1539kB 60.04%  runtime.schedule
(pprof) 




type MEXCCollector struct {
	ctx      context.Context
	cancel   context.CancelFunc
	conn     *websocket.Conn
	symbols  []string
	allowed  map[string]struct{} // whitelist
	lastData map[string]struct {
		Bid, Ask, BidSize, AskSize float64
	}
	mu  sync.Mutex
	buf []byte // отдельный буфер для NormalizeSymbol

	pool *sync.Pool
}

func NewMEXCCollector(symbols []string, whitelist []string, pool *sync.Pool) *MEXCCollector {
	ctx, cancel := context.WithCancel(context.Background())
	allowed := make(map[string]struct{}, len(whitelist))
	for _, s := range whitelist {
		allowed[market.NormalizeSymbol_Full(s)] = struct{}{}
	}
	return &MEXCCollector{
		ctx:      ctx,
		cancel:   cancel,
		symbols:  symbols,
		allowed:  allowed,
		lastData: make(map[string]struct{ Bid, Ask, BidSize, AskSize float64 }),
		pool:     pool,
		buf:      []buf,
	}
}


symbol = market.NormalizeSymbol_NoAlloc(symbol, &c.buf)
	if symbol == "" {
		return nil
	}
