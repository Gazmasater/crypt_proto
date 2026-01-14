Название API
9623527002

6966b78122ca320001d2acae
fa1e37ae-21ff-4257-844d-3dcd21d26ccd





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



Duration: 30.05s, Total samples = 1.99s ( 6.62%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 1130ms, 56.78% of 1990ms total
Showing top 10 nodes out of 182
      flat  flat%   sum%        cum   cum%
     570ms 28.64% 28.64%      570ms 28.64%  internal/runtime/syscall.Syscall6
     170ms  8.54% 37.19%      170ms  8.54%  runtime.futex
      70ms  3.52% 40.70%       70ms  3.52%  runtime.typePointers.next
      60ms  3.02% 43.72%       60ms  3.02%  runtime.memmove
      50ms  2.51% 46.23%       70ms  3.52%  github.com/tidwall/gjson.parseObject
      50ms  2.51% 48.74%      160ms  8.04%  runtime.scanobject
      40ms  2.01% 50.75%       40ms  2.01%  aeshashbody
      40ms  2.01% 52.76%       90ms  4.52%  runtime.concatstrings
      40ms  2.01% 54.77%      140ms  7.04%  runtime.mallocgcSmallScanNoHeader
      40ms  2.01% 56.78%       40ms  2.01%  runtime.nanotime
(pprof) 



gaz358@gaz358-BOD-WXX9:~/myprog/crypt_proto$ go tool pprof http://localhost:6060/debug/pprof/heap
Fetching profile over HTTP from http://localhost:6060/debug/pprof/heap
Saved profile in /home/gaz358/pprof/pprof.arb.alloc_objects.alloc_space.inuse_objects.inuse_space.006.pb.gz
File: arb
Build ID: d8095edd0e7b84fc1bd6776bcad2be691a5b7dcc
Type: inuse_space
Time: 2026-01-15 01:11:27 MSK
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 4075.42kB, 100% of 4075.42kB total
Showing top 10 nodes out of 40
      flat  flat%   sum%        cum   cum%
    1539kB 37.76% 37.76%     1539kB 37.76%  runtime.allocm
 1024.04kB 25.13% 62.89%  1024.04kB 25.13%  crypto/internal/fips140/nistec.NewP384Point (inline)
 1000.34kB 24.55% 87.44%  1000.34kB 24.55%  main.main
  512.05kB 12.56%   100%   512.05kB 12.56%  runtime.acquireSudog
         0     0%   100%  1024.04kB 25.13%  crypto/ecdsa.VerifyASN1
         0     0%   100%  1024.04kB 25.13%  crypto/ecdsa.verifyFIPS[go.shape.*crypto/internal/fips140/nistec.P384Point]
         0     0%   100%  1024.04kB 25.13%  crypto/internal/fips140/ecdsa.Verify[go.shape.*crypto/internal/fips140/nistec.P384Point]
         0     0%   100%  1024.04kB 25.13%  crypto/internal/fips140/ecdsa.verifyGeneric[go.shape.*crypto/internal/fips140/nistec.P384Point]
         0     0%   100%  1024.04kB 25.13%  crypto/internal/fips140/ecdsa.verify[go.shape.*crypto/internal/fips140/nistec.P384Point] (inline)
         0     0%   100%  1024.04kB 25.13%  crypto/internal/fips140/nistec.(*P384Point).ScalarBaseMult
(pprof) 


