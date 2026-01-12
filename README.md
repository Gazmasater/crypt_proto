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




(pprof) gaz358@gaz358-BOD-WXX9:~/myprog/crypt_proto$    go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
Fetching profile over HTTP from http://localhost:6060/debug/pprof/profile?seconds=30
Saved profile in /home/gaz358/pprof/pprof.arb.samples.cpu.081.pb.gz
File: arb
Build ID: 00f359f630cea5d5eb1389920b6bee5aa91f0b5e
Type: cpu
Time: 2026-01-12 10:59:06 MSK
Duration: 30.04s, Total samples = 1.98s ( 6.59%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 1130ms, 57.07% of 1980ms total
Showing top 10 nodes out of 209
      flat  flat%   sum%        cum   cum%
     670ms 33.84% 33.84%      670ms 33.84%  internal/runtime/syscall.Syscall6
     100ms  5.05% 38.89%      100ms  5.05%  runtime.futex
      70ms  3.54% 42.42%       70ms  3.54%  aeshashbody
      60ms  3.03% 45.45%      130ms  6.57%  runtime.scanobject
      50ms  2.53% 47.98%       90ms  4.55%  github.com/tidwall/gjson.parseObject
      50ms  2.53% 50.51%      130ms  6.57%  runtime.mapassign_faststr
      40ms  2.02% 52.53%       50ms  2.53%  runtime.typePointers.next
      30ms  1.52% 54.04%      820ms 41.41%  bufio.(*Reader).fill
      30ms  1.52% 55.56%       30ms  1.52%  memeqbody
      30ms  1.52% 57.07%       60ms  3.03%  runtime.mapaccess1_faststr
(pprof) 



Fetching profile over HTTP from http://localhost:6060/debug/pprof/profile?seconds=30
Saved profile in /home/gaz358/pprof/pprof.arb.samples.cpu.082.pb.gz
File: arb
Build ID: 991d3b51d26d0a48852c28a66aa2039c318c2e53
Type: cpu
Time: 2026-01-12 11:20:48 MSK
Duration: 30s, Total samples = 1.55s ( 5.17%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 1100ms, 70.97% of 1550ms total
Showing top 10 nodes out of 127
      flat  flat%   sum%        cum   cum%
     760ms 49.03% 49.03%      760ms 49.03%  internal/runtime/syscall.Syscall6
     170ms 10.97% 60.00%      170ms 10.97%  runtime.futex
      30ms  1.94% 61.94%       40ms  2.58%  strings.Fields
      20ms  1.29% 63.23%       20ms  1.29%  crypto/internal/fips140/aes/gcm.gcmAesDec
      20ms  1.29% 64.52%      860ms 55.48%  crypto/tls.(*Conn).readRecordOrCCS
      20ms  1.29% 65.81%      950ms 61.29%  github.com/gorilla/websocket.(*Conn).ReadMessage
      20ms  1.29% 67.10%       20ms  1.29%  github.com/gorilla/websocket.(*messageReader).Read
      20ms  1.29% 68.39%       20ms  1.29%  runtime.(*mspan).base
      20ms  1.29% 69.68%       20ms  1.29%  runtime.execute
      20ms  1.29% 70.97%       30ms  1.94%  runtime.ifaceeq
(pprof) 



type Calculator struct {
	mem       *queue.MemoryStore
	triangles []Triangle
	bySymbol  map[string][]Triangle
	fileLog   *log.Logger
}



func NewCalculator(mem *queue.MemoryStore, triangles []Triangle) *Calculator {
	f, err := os.OpenFile(
		"arb_opportunities.log",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	if err != nil {
		log.Fatalf("failed to open arb log file: %v", err)
	}

	c := &Calculator{
		mem:       mem,
		triangles: triangles,
		bySymbol:  make(map[string][]Triangle),
		fileLog:   log.New(f, "", log.LstdFlags),
	}

	for _, tri := range triangles {
		s1 := legSymbol(tri.Leg1)
		s2 := legSymbol(tri.Leg2)
		s3 := legSymbol(tri.Leg3)

		if s1 != "" {
			c.bySymbol[s1] = append(c.bySymbol[s1], tri)
		}
		if s2 != "" {
			c.bySymbol[s2] = append(c.bySymbol[s2], tri)
		}
		if s3 != "" {
			c.bySymbol[s3] = append(c.bySymbol[s3], tri)
		}
	}

	return c
}



func (c *Calculator) Run(in <-chan *models.MarketData) {
	for md := range in {

		// сохраняем маркет
		c.mem.Set(md)

		tris := c.bySymbol[md.Symbol]
		if len(tris) == 0 {
			continue
		}

		for _, tri := range tris {
			c.calcTriangle(&tri)
		}
	}
}



func (c *Calculator) calcTriangle(tri *Triangle) {

	s1 := legSymbol(tri.Leg1)
	s2 := legSymbol(tri.Leg2)
	s3 := legSymbol(tri.Leg3)

	q1, ok1 := c.mem.Get("KuCoin", s1)
	q2, ok2 := c.mem.Get("KuCoin", s2)
	q3, ok3 := c.mem.Get("KuCoin", s3)

	if !ok1 || !ok2 || !ok3 {
		return
	}

	// ===== 1. USDT LIMITS =====

	usdtLimits := make([]float64, 0, 3)

	// LEG 1
	if strings.HasPrefix(tri.Leg1, "BUY") {
		if q1.Ask <= 0 || q1.AskSize <= 0 {
			return
		}
		usdtLimits = append(usdtLimits, q1.Ask*q1.AskSize)
	} else {
		if q1.Bid <= 0 || q1.BidSize <= 0 {
			return
		}
		usdtLimits = append(usdtLimits, q1.Bid*q1.BidSize)
	}

	// LEG 2
	if strings.HasPrefix(tri.Leg2, "BUY") {
		if q2.Ask <= 0 || q2.AskSize <= 0 || q3.Bid <= 0 {
			return
		}
		usdtLimits = append(usdtLimits, q2.Ask*q2.AskSize*q3.Bid)
	} else {
		if q2.Bid <= 0 || q2.BidSize <= 0 || q3.Bid <= 0 {
			return
		}
		usdtLimits = append(usdtLimits, q2.BidSize*q3.Bid)
	}

	// LEG 3
	if q3.Bid <= 0 || q3.BidSize <= 0 {
		return
	}
	usdtLimits = append(usdtLimits, q3.Bid*q3.BidSize)

	// ===== 2. MIN LIMIT =====

	maxUSDT := usdtLimits[0]
	for _, v := range usdtLimits {
		if v < maxUSDT {
			maxUSDT = v
		}
	}
	if maxUSDT <= 0 {
		return
	}

	// ===== 3. PROGON =====

	amount := maxUSDT

	if strings.HasPrefix(tri.Leg1, "BUY") {
		amount = amount / q1.Ask * (1 - fee)
	} else {
		amount = amount * q1.Bid * (1 - fee)
	}

	if strings.HasPrefix(tri.Leg2, "BUY") {
		amount = amount / q2.Ask * (1 - fee)
	} else {
		amount = amount * q2.Bid * (1 - fee)
	}

	if strings.HasPrefix(tri.Leg3, "BUY") {
		amount = amount / q3.Ask * (1 - fee)
	} else {
		amount = amount * q3.Bid * (1 - fee)
	}

	profitUSDT := amount - maxUSDT
	profitPct := profitUSDT / maxUSDT

	if profitPct > 0.001 && profitUSDT > 0.02 {
		// лог включишь когда надо
	}
}


[{
	"resource": "/home/gaz358/myprog/crypt_proto/cmd/arb/main.go",
	"owner": "_generated_diagnostic_collection_name_#0",
	"code": {
		"value": "WrongArgCount",
		"target": {
			"$mid": 1,
			"path": "/golang.org/x/tools/internal/typesinternal",
			"scheme": "https",
			"authority": "pkg.go.dev",
			"fragment": "WrongArgCount"
		}
	},
	"severity": 8,
	"message": "not enough arguments in call to calc.Run\n\thave ()\n\twant (<-chan *models.MarketData)",
	"source": "compiler",
	"startLineNumber": 61,
	"startColumn": 14,
	"endLineNumber": 61,
	"endColumn": 14,
	"origin": "extHost1"
}]

[{
	"resource": "/home/gaz358/myprog/crypt_proto/internal/calculator/arb.go",
	"owner": "_generated_diagnostic_collection_name_#0",
	"code": {
		"value": "MissingFieldOrMethod",
		"target": {
			"$mid": 1,
			"path": "/golang.org/x/tools/internal/typesinternal",
			"scheme": "https",
			"authority": "pkg.go.dev",
			"fragment": "MissingFieldOrMethod"
		}
	},
	"severity": 8,
	"message": "c.mem.Set undefined (type *queue.MemoryStore has no field or method Set)",
	"source": "compiler",
	"startLineNumber": 70,
	"startColumn": 9,
	"endLineNumber": 70,
	"endColumn": 12,
	"origin": "extHost1"
}]



marketCh := make(chan *models.MarketData, 10000)

collector.Start(marketCh)

calc := calculator.NewCalculator(mem, collector.Triangles())
go calc.Run(marketCh)








