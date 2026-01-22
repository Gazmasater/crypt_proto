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





package calculator

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"

	"crypt_proto/internal/queue"
	"crypt_proto/pkg/models"
)

const (
	minVolumeUSDT = 20.0
	feeMul        = 0.999
)

type LegIndex struct {
	Key    string
	Symbol string
	IsBuy  bool
}

type Triangle struct {
	A, B, C string
	Legs    [3]LegIndex
}

type Calculator struct {
	mem      *queue.MemoryStore
	bySymbol map[string][]*Triangle

	// roughMax для каждого треугольника, обновляемый при изменении котировок
	roughMax map[*Triangle]float64

	fileLog *log.Logger
}

// NewCalculator — строим индекс symbol -> triangles и создаем roughMax map
func NewCalculator(mem *queue.MemoryStore, triangles []*Triangle) *Calculator {
	f, err := os.OpenFile("arb_opportunities.log",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("failed to open log: %v", err)
	}

	bySymbol := make(map[string][]*Triangle, 1024)
	for _, t := range triangles {
		for _, leg := range t.Legs {
			bySymbol[leg.Symbol] = append(bySymbol[leg.Symbol], t)
		}
	}

	log.Printf("[Calculator] indexed %d symbols\n", len(bySymbol))

	return &Calculator{
		mem:      mem,
		bySymbol: bySymbol,
		roughMax: make(map[*Triangle]float64, len(triangles)),
		fileLog:  log.New(f, "", log.LstdFlags),
	}
}

// Run — обрабатываем поток котировок
func (c *Calculator) Run(in <-chan *models.MarketData) {
	for md := range in {
		c.mem.Push(md)

		tris := c.bySymbol[md.Symbol]
		if len(tris) == 0 {
			continue
		}

		// пересчитываем roughMax для всех треугольников с этой парой
		for _, tri := range tris {
			c.updateRoughMax(tri)
			c.calcTriangle(tri)
		}
	}
}

// updateRoughMax — грубая оценка прибыли (без объёмов и fee)
func (c *Calculator) updateRoughMax(tri *Triangle) {
	q0, ok0 := c.mem.Get("KuCoin", tri.Legs[0].Symbol)
	q1, ok1 := c.mem.Get("KuCoin", tri.Legs[1].Symbol)
	q2, ok2 := c.mem.Get("KuCoin", tri.Legs[2].Symbol)
	if !ok0 || !ok1 || !ok2 {
		c.roughMax[tri] = 0
		return
	}

	v0 := q0.Bid
	if tri.Legs[0].IsBuy {
		v0 = q0.Ask
	}
	v1 := q1.Bid
	if tri.Legs[1].IsBuy {
		v1 = q1.Ask
	}
	v2 := q2.Bid
	if tri.Legs[2].IsBuy {
		v2 = q2.Ask
	}

	// грубая проверка минимальной возможности прибыли
	if v0 <= v1*v2 || v0 < minVolumeUSDT || v1 < minVolumeUSDT || v2 < minVolumeUSDT {
		c.roughMax[tri] = 0
		return
	}

	c.roughMax[tri] = v0 - v1*v2
}

// calcTriangle — точный расчёт объёма и прибыли
func (c *Calculator) calcTriangle(tri *Triangle) {
	// если roughMax <= 0, треугольник сразу пропускаем
	if c.roughMax[tri] <= 0 {
		return
	}

	q0, ok0 := c.mem.Get("KuCoin", tri.Legs[0].Symbol)
	q1, ok1 := c.mem.Get("KuCoin", tri.Legs[1].Symbol)
	q2, ok2 := c.mem.Get("KuCoin", tri.Legs[2].Symbol)
	if !ok0 || !ok1 || !ok2 {
		return
	}

	var usdt0, usdt1, usdt2 float64

	if tri.Legs[0].IsBuy {
		if q0.Ask <= 0 || q0.AskSize <= 0 {
			return
		}
		usdt0 = q0.Ask * q0.AskSize
	} else {
		if q0.Bid <= 0 || q0.BidSize <= 0 {
			return
		}
		usdt0 = q0.Bid * q0.BidSize
	}

	if tri.Legs[1].IsBuy {
		if q1.Ask <= 0 || q1.AskSize <= 0 || q2.Bid <= 0 {
			return
		}
		usdt1 = q1.Ask * q1.AskSize * q2.Bid
	} else {
		if q1.Bid <= 0 || q1.BidSize <= 0 || q2.Bid <= 0 {
			return
		}
		usdt1 = q1.BidSize * q2.Bid
	}

	if q2.Bid <= 0 || q2.BidSize <= 0 {
		return
	}
	usdt2 = q2.Bid * q2.BidSize

	// минимальный объём
	maxUSDT := usdt0
	if usdt1 < maxUSDT {
		maxUSDT = usdt1
	}
	if usdt2 < maxUSDT {
		maxUSDT = usdt2
	}
	if maxUSDT < minVolumeUSDT {
		return
	}

	// расчет прибыли с учетом fee
	amount := maxUSDT
	if tri.Legs[0].IsBuy {
		amount = amount / q0.Ask * feeMul
	} else {
		amount = amount * q0.Bid * feeMul
	}
	if tri.Legs[1].IsBuy {
		amount = amount / q1.Ask * feeMul
	} else {
		amount = amount * q1.Bid * feeMul
	}
	if tri.Legs[2].IsBuy {
		amount = amount / q2.Ask * feeMul
	} else {
		amount = amount * q2.Bid * feeMul
	}

	profitUSDT := amount - maxUSDT
	profitPct := profitUSDT / maxUSDT

	if profitPct > 0.001 && profitUSDT > 0.02 {
		msg := fmt.Sprintf(
			"[ARB] %s→%s→%s | %.4f%% | volume=%.2f USDT | profit=%.4f USDT",
			tri.A, tri.B, tri.C,
			profitPct*100, maxUSDT, profitUSDT,
		)
		log.Println(msg)
		c.fileLog.Println(msg)
	}
}

// CSV без изменений логики, но сразу сохраняем Symbol
func ParseTrianglesFromCSV(path string) ([]*Triangle, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rows, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return nil, err
	}

	var res []*Triangle

	for _, row := range rows[1:] {
		if len(row) < 6 {
			continue
		}

		t := &Triangle{
			A: strings.TrimSpace(row[0]),
			B: strings.TrimSpace(row[1]),
			C: strings.TrimSpace(row[2]),
		}

		for i, leg := range []string{row[3], row[4], row[5]} {
			leg = strings.ToUpper(strings.TrimSpace(leg))
			parts := strings.Fields(leg)
			if len(parts) != 2 {
				continue
			}

			isBuy := parts[0] == "BUY"
			pair := strings.Split(parts[1], "/")
			if len(pair) != 2 {
				continue
			}

			symbol := pair[0] + "-" + pair[1]
			key := "KuCoin|" + symbol

			t.Legs[i] = LegIndex{
				Key:    key,
				Symbol: symbol,
				IsBuy:  isBuy,
			}
		}

		res = append(res, t)
	}

	return res, nil
}



gaz358@gaz358-BOD-WXX9:~/myprog/crypt_proto$ go tool pprof http://localhost:6060/debug/pprof/heap
Fetching profile over HTTP from http://localhost:6060/debug/pprof/heap
Saved profile in /home/gaz358/pprof/pprof.arb.alloc_objects.alloc_space.inuse_objects.inuse_space.007.pb.gz
File: arb
Build ID: 2c883718df17e94f03754d4a86aacb6161e54ba4
Type: inuse_space
Time: 2026-01-23 00:53:42 MSK
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 4078.44kB, 100% of 4078.44kB total
Showing top 10 nodes out of 44
      flat  flat%   sum%        cum   cum%
    1539kB 37.74% 37.74%     1539kB 37.74%  runtime.allocm
 1000.34kB 24.53% 62.26%  1515.34kB 37.15%  main.main
     515kB 12.63% 74.89%      515kB 12.63%  bytes.growSlice
  512.05kB 12.55% 87.45%   512.05kB 12.55%  crypto/x509.(*CertPool).AppendCertsFromPEM
  512.05kB 12.55%   100%   512.05kB 12.55%  runtime.acquireSudog
         0     0%   100%      515kB 12.63%  bytes.(*Buffer).Grow (inline)
         0     0%   100%      515kB 12.63%  bytes.(*Buffer).grow
         0     0%   100%      515kB 12.63%  crypt_proto/internal/collector.(*KuCoinCollector).Start
         0     0%   100%      515kB 12.63%  crypt_proto/internal/collector.(*kucoinWS).connect
         0     0%   100%  1027.05kB 25.18%  crypto/tls.(*Conn).HandshakeContext (inline)
(pprof) 




(pprof) gaz358@gaz358-BOD-WXX9:~/myprog/crypt_proto$    go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
Fetching profile over HTTP from http://localhost:6060/debug/pprof/profile?seconds=30
Saved profile in /home/gaz358/pprof/pprof.arb.samples.cpu.211.pb.gz
File: arb
Build ID: 2c883718df17e94f03754d4a86aacb6161e54ba4
Type: cpu
Time: 2026-01-23 00:52:29 MSK
Duration: 30s, Total samples = 1.37s ( 4.57%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 870ms, 63.50% of 1370ms total
Showing top 10 nodes out of 165
      flat  flat%   sum%        cum   cum%
     520ms 37.96% 37.96%      520ms 37.96%  internal/runtime/syscall.Syscall6
     130ms  9.49% 47.45%      130ms  9.49%  runtime.futex
      50ms  3.65% 51.09%       60ms  4.38%  runtime.typePointers.next
      30ms  2.19% 53.28%       90ms  6.57%  crypto/tls.(*halfConn).decrypt
      30ms  2.19% 55.47%       30ms  2.19%  runtime.memclrNoHeapPointers
      30ms  2.19% 57.66%       30ms  2.19%  runtime.nextFreeFast
      20ms  1.46% 59.12%      720ms 52.55%  bufio.(*Reader).fill
      20ms  1.46% 60.58%       20ms  1.46%  crypto/internal/fips140/aes/gcm.gcmAesData
      20ms  1.46% 62.04%       20ms  1.46%  crypto/tls.(*halfConn).explicitNonceLen
      20ms  1.46% 63.50%      750ms 54.74%  github.com/gorilla/websocket.(*Conn).NextReader
(pprof) 







