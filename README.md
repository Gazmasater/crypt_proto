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




package queue

import (
	"sync"
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
	m     sync.Map
	batch chan *models.MarketData
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		batch: make(chan *models.MarketData, 10_000),
	}
}

// Run — основной цикл стора
func (s *MemoryStore) Run() {
	for md := range s.batch {
		s.apply(md)
	}
}

// Push — приём данных от коллекторов
func (s *MemoryStore) Push(md *models.MarketData) {
	select {
	case s.batch <- md:
	default:
		// защита от переполнения
	}
}

// Get — lock-free чтение
func (s *MemoryStore) Get(exchange, symbol string) (Quote, bool) {
	key := exchange + "|" + symbol
	q, ok := s.m.Load(key)
	if !ok {
		return Quote{}, false
	}
	return q.(Quote), true
}

// apply — lock-free запись
func (s *MemoryStore) apply(md *models.MarketData) {
	key := md.Exchange + "|" + md.Symbol
	quote := Quote{
		Bid:       md.Bid,
		Ask:       md.Ask,
		BidSize:   md.BidSize,
		AskSize:   md.AskSize,
		Timestamp: time.Now().UnixMilli(),
	}
	s.m.Store(key, quote)
}




package calculator

import (
	"crypt_proto/internal/queue"
	"crypt_proto/pkg/models"
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

// fee
const fee = 0.001 // 0.1%

type Triangle struct {
	A, B, C string

	// заранее подготовленные LegIndex
	Legs [3]LegIndex
}

// LegIndex — информация по каждой ноге для быстрого расчёта
type LegIndex struct {
	Key   string // BASE-QUOTE
	IsBuy bool
}

type Calculator struct {
	mem       *queue.MemoryStore
	triangles []Triangle
	bySymbol  map[string][]*Triangle
	fileLog   *log.Logger
}

func NewCalculator(mem *queue.MemoryStore, triangles []Triangle) *Calculator {
	f, err := os.OpenFile("arb_opportunities.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("failed to open arb log file: %v", err)
	}

	c := &Calculator{
		mem:       mem,
		triangles: triangles,
		bySymbol:  make(map[string][]*Triangle),
		fileLog:   log.New(f, "", log.LstdFlags),
	}

	for i := range triangles {
		tri := &triangles[i]
		// подготавливаем LegIndex
		for j, leg := range [3]string{tri.Leg1, tri.Leg2, tri.Leg3} {
			tri.Legs[j] = parseLegIndex(leg)
		}
		// по символам для быстрого поиска
		for _, l := range tri.Legs {
			if l.Key != "" {
				c.bySymbol[l.Key] = append(c.bySymbol[l.Key], tri)
			}
		}
	}

	return c
}

func (c *Calculator) Run(in <-chan *models.MarketData) {
	for md := range in {
		c.mem.Push(md)

		tris := c.bySymbol[md.Symbol]
		for _, tri := range tris {
			c.calcTriangle(tri)
		}
	}
}

func (c *Calculator) calcTriangle(tri *Triangle) {
	q := make([]queue.Quote, 3)
	for i, leg := range tri.Legs {
		quote, ok := c.mem.Get("KuCoin", leg.Key)
		if !ok {
			return
		}
		q[i] = quote
	}

	// ===== 1. USDT LIMITS =====
	var usdt [3]float64

	for i, leg := range tri.Legs {
		switch i {
		case 0:
			if leg.IsBuy {
				if q[i].Ask <= 0 || q[i].AskSize <= 0 {
					return
				}
				usdt[i] = q[i].Ask * q[i].AskSize
			} else {
				if q[i].Bid <= 0 || q[i].BidSize <= 0 {
					return
				}
				usdt[i] = q[i].Bid * q[i].BidSize
			}
		case 1:
			if leg.IsBuy {
				if q[i].Ask <= 0 || q[i].AskSize <= 0 || q[2].Bid <= 0 {
					return
				}
				usdt[i] = q[i].Ask * q[i].AskSize * q[2].Bid
			} else {
				if q[i].Bid <= 0 || q[i].BidSize <= 0 || q[2].Bid <= 0 {
					return
				}
				usdt[i] = q[i].BidSize * q[2].Bid
			}
		case 2:
			if q[i].Bid <= 0 || q[i].BidSize <= 0 {
				return
			}
			usdt[i] = q[i].Bid * q[i].BidSize
		}
	}

	// ===== 2. MIN LIMIT =====
	maxUSDT := usdt[0]
	if usdt[1] < maxUSDT {
		maxUSDT = usdt[1]
	}
	if usdt[2] < maxUSDT {
		maxUSDT = usdt[2]
	}
	if maxUSDT <= 0 {
		return
	}

	// ===== 3. PROGON =====
	amount := maxUSDT
	for i, leg := range tri.Legs {
		if leg.IsBuy {
			amount = amount / q[i].Ask * (1 - fee)
		} else {
			amount = amount * q[i].Bid * (1 - fee)
		}
	}

	profitUSDT := amount - maxUSDT
	profitPct := profitUSDT / maxUSDT

	if profitPct > 0.001 && profitUSDT > 0.02 {
		msg := fmt.Sprintf("[ARB] %s→%s→%s | %.4f%% | vol=%.2f | profit=%.4f",
			tri.A, tri.B, tri.C, profitPct*100, maxUSDT, profitUSDT)
		log.Println(msg)
		c.fileLog.Println(msg)
	}
}

// ================= HELPERS =================

type LegIndex struct {
	Key   string
	IsBuy bool
}

func parseLegIndex(leg string) LegIndex {
	l := LegIndex{}
	if len(leg) < 4 {
		return l
	}
	if leg[:3] == "BUY" {
		l.IsBuy = true
	} else {
		l.IsBuy = false
	}
	parts := []rune(leg)
	// найти пробел
	for i, r := range parts {
		if r == ' ' {
			baseQuote := string(parts[i+1:])
			s := []rune(baseQuote)
			for j, c := range s {
				if c >= 'a' && c <= 'z' {
					s[j] = c - 'a' + 'A'
				}
			}
			p := string(s)
			// BASE/QUOTE -> BASE-QUOTE
			for k, c := range p {
				if c == '/' {
					p = p[:k] + "-" + p[k+1:]
					break
				}
			}
			l.Key = p
			break
		}
	}
	return l
}


