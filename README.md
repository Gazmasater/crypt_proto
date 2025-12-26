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
Saved profile in /home/gaz358/pprof/pprof.arb.samples.cpu.001.pb.gz
File: arb
Build ID: 2f2a4b34fa41455b1a30bee46dd74b5e51f355d0
Type: cpu
Time: 2025-12-26 03:57:21 MSK
Duration: 30s, Total samples = 150ms (  0.5%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top 
Showing nodes accounting for 150ms, 100% of 150ms total
Showing top 10 nodes out of 53
      flat  flat%   sum%        cum   cum%
      80ms 53.33% 53.33%       80ms 53.33%  internal/runtime/syscall.Syscall6
      10ms  6.67% 60.00%       10ms  6.67%  crypto/internal/fips140/aes.encryptBlock
      10ms  6.67% 66.67%       10ms  6.67%  gogo
      10ms  6.67% 73.33%       10ms  6.67%  google.golang.org/protobuf/internal/impl.offset.IsValid
      10ms  6.67% 80.00%       10ms  6.67%  os.(*File).write
      10ms  6.67% 86.67%       10ms  6.67%  reflect.(*rtype).Elem
      10ms  6.67% 93.33%       10ms  6.67%  runtime.newobject
      10ms  6.67%   100%       10ms  6.67%  sync.(*Pool).pin
         0     0%   100%       80ms 53.33%  bufio.(*Reader).Peek
         0     0%   100%       80ms 53.33%  bufio.(*Reader).fill
(pprof) 



gaz358@gaz358-BOD-WXX9:~/myprog/crypt_proto$    go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
Fetching profile over HTTP from http://localhost:6060/debug/pprof/profile?seconds=30
Saved profile in /home/gaz358/pprof/pprof.arb.samples.cpu.002.pb.gz
File: arb
Build ID: 2f2a4b34fa41455b1a30bee46dd74b5e51f355d0
Type: cpu
Time: 2025-12-26 03:59:12 MSK
Duration: 30s, Total samples = 100ms ( 0.33%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 100ms, 100% of 100ms total
Showing top 10 nodes out of 56
      flat  flat%   sum%        cum   cum%
      60ms 60.00% 60.00%       60ms 60.00%  internal/runtime/syscall.Syscall6
      10ms 10.00% 70.00%       10ms 10.00%  encoding/json.checkValid
      10ms 10.00% 80.00%       10ms 10.00%  runtime.futex
      10ms 10.00% 90.00%       10ms 10.00%  runtime.write1
      10ms 10.00%   100%       10ms 10.00%  strconv.atof64
         0     0%   100%       40ms 40.00%  bufio.(*Reader).Peek
         0     0%   100%       40ms 40.00%  bufio.(*Reader).fill
         0     0%   100%       40ms 40.00%  bytes.(*Buffer).ReadFrom
         0     0%   100%       70ms 70.00%  crypt_proto/internal/collector.(*OKXCollector).Start.func2
         0     0%   100%       40ms 40.00%  crypto/tls.(*Conn).Read
(pprof) 




gaz358@gaz358-BOD-WXX9:~/myprog/crypt_proto$    go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
Fetching profile over HTTP from http://localhost:6060/debug/pprof/profile?seconds=30
Saved profile in /home/gaz358/pprof/pprof.arb.samples.cpu.003.pb.gz
File: arb
Build ID: 2f2a4b34fa41455b1a30bee46dd74b5e51f355d0
Type: cpu
Time: 2025-12-26 04:00:51 MSK
Duration: 30s, Total samples = 110ms ( 0.37%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 110ms, 100% of 110ms total
Showing top 10 nodes out of 63
      flat  flat%   sum%        cum   cum%
      30ms 27.27% 27.27%       30ms 27.27%  internal/runtime/syscall.Syscall6
      10ms  9.09% 36.36%       60ms 54.55%  bufio.(*Reader).fill
      10ms  9.09% 45.45%       10ms  9.09%  runtime.(*itabTableType).find
      10ms  9.09% 54.55%       10ms  9.09%  runtime.futex
      10ms  9.09% 63.64%       10ms  9.09%  runtime.gogo
      10ms  9.09% 72.73%       10ms  9.09%  runtime.nextFreeFast
      10ms  9.09% 81.82%       10ms  9.09%  runtime.typePointers.next
      10ms  9.09% 90.91%       10ms  9.09%  runtime.write1
      10ms  9.09%   100%       10ms  9.09%  strings.Replace
         0     0%   100%       60ms 54.55%  bufio.(*Reader).Peek
(pprof) 







