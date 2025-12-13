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
Saved profile in /home/gaz358/pprof/pprof.cryptarb.samples.cpu.022.pb.gz
File: cryptarb
Build ID: ec57481fe3c993ad5c0ecc311bfddcb4dc62557e
Type: cpu
Time: 2025-12-13 04:00:23 MSK
Duration: 30s, Total samples = 70ms ( 0.23%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 70ms, 100% of 70ms total
Showing top 10 nodes out of 53
      flat  flat%   sum%        cum   cum%
      10ms 14.29% 14.29%       10ms 14.29%  crypto/internal/fips140/bigmod.addMulVVW2048
      10ms 14.29% 28.57%       10ms 14.29%  crypto/internal/fips140/mlkem.inverseNTT
      10ms 14.29% 42.86%       10ms 14.29%  crypto/internal/fips140/mlkem.polyByteEncode[go.shape.[256]crypto/internal/fips140/mlkem.fieldElement]
      10ms 14.29% 57.14%       10ms 14.29%  runtime.futex
      10ms 14.29% 71.43%       10ms 14.29%  runtime.memclrNoHeapPointers
      10ms 14.29% 85.71%       10ms 14.29%  vendor/golang.org/x/crypto/cryptobyte.(*Builder).addLengthPrefixed
      10ms 14.29%   100%       10ms 14.29%  vendor/golang.org/x/net/dns/dnsmessage.(*Parser).AnswerHeader
         0     0%   100%       50ms 71.43%  crypt_proto/mexc.(*Feed).runPublicBookTickerWS
         0     0%   100%       10ms 14.29%  crypto/internal/fips140/bigmod.(*Nat).ExpShortVarTime
         0     0%   100%       10ms 14.29%  crypto/internal/fips140/bigmod.(*Nat).montgomeryMul
(pprof) 



25/12/13 04:03:46.636752 [MEXC WS #7] connected to wss://wbs-api.mexc.com/ws (symbols: 50)
2025/12/13 04:03:46.636950 [MEXC WS #7] SUB -> 50 topics
2025/12/13 04:03:46.866113 [MEXC WS #1] read err: websocket: close 1005 (no status) (reconnect)
2025/12/13 04:03:47.788372 [MEXC WS #5] read err: websocket: close 1005 (no status) (reconnect)
2025/12/13 04:03:48.095745 [MEXC WS #3] connected to wss://wbs-api.mexc.com/ws (symbols: 50)
2025/12/13 04:03:48.095925 [MEXC WS #3] SUB -> 50 topics
2025/12/13 04:03:49.792474 [MEXC WS #1] connected to wss://wbs-api.mexc.com/ws (symbols: 50)
2025/12/13 04:03:49.792667 [MEXC WS #1] SUB -> 50 topics
2025/12/13 04:03:50.656911 [MEXC WS #5] connected to wss://wbs-api.mexc.com/ws (symbols: 50)
2025/12/13 04:03:50.657065 [MEXC WS #5] SUB -> 50 topics






