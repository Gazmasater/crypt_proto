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




package calculator

import (
	"fmt"
	"log"
	"os"

	"crypt_proto/internal/queue"
)

const fee = 0.001

type TriangleFast struct {
	A, B, C                   string
	Leg1Idx, Leg2Idx, Leg3Idx string
	Buy1, Buy2, Buy3          bool
}

type CalculatorFast struct {
	mem       *queue.MemoryStore
	triangles []TriangleFast
	fileLog   *log.Logger
}

func NewCalculatorFast(mem *queue.MemoryStore, triangles []TriangleFast) *CalculatorFast {
	f, err := os.OpenFile("arb_opportunities.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("failed to open arb log file: %v", err)
	}
	return &CalculatorFast{
		mem:       mem,
		triangles: triangles,
		fileLog:   log.New(f, "", log.LstdFlags),
	}
}

func (c *CalculatorFast) CalcTriangleFast(tri TriangleFast) {
	q1, ok1 := c.mem.Get("KuCoin", tri.Leg1Idx)
	q2, ok2 := c.mem.Get("KuCoin", tri.Leg2Idx)
	q3, ok3 := c.mem.Get("KuCoin", tri.Leg3Idx)
	if !ok1 || !ok2 || !ok3 {
		return
	}

	// Рассчёт лимитов
	usdt1 := q1.Ask * q1.AskSize
	if !tri.Buy1 {
		usdt1 = q1.Bid * q1.BidSize
	}
	usdt2 := q2.Ask * q2.AskSize
	if !tri.Buy2 {
		usdt2 = q2.Bid * q2.BidSize
	}
	usdt3 := q3.Ask * q3.AskSize
	if !tri.Buy3 {
		usdt3 = q3.Bid * q3.BidSize
	}

	maxUSDT := usdt1
	if usdt2 < maxUSDT {
		maxUSDT = usdt2
	}
	if usdt3 < maxUSDT {
		maxUSDT = usdt3
	}
	if maxUSDT <= 0 {
		return
	}

	amount := maxUSDT
	if tri.Buy1 {
		amount = amount / q1.Ask * (1 - fee)
	} else {
		amount = amount * q1.Bid * (1 - fee)
	}
	if tri.Buy2 {
		amount = amount / q2.Ask * (1 - fee)
	} else {
		amount = amount * q2.Bid * (1 - fee)
	}
	if tri.Buy3 {
		amount = amount / q3.Ask * (1 - fee)
	} else {
		amount = amount * q3.Bid * (1 - fee)
	}

	profitUSDT := amount - maxUSDT
	profitPct := profitUSDT / maxUSDT

	//if profitPct > 0.001 && profitUSDT > 0.02 {
	msg := fmt.Sprintf("[ARB] %s → %s → %s | %.4f%% | volume=%.2f USDT | profit=%.4f USDT",
		tri.A, tri.B, tri.C, profitPct*100, maxUSDT, profitUSDT)
	log.Println(msg)
	c.fileLog.Println(msg)
	//}
}





