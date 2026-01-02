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



📁 Итоговая структура проекта
exchange/
├── common/
│   ├── market.go
│   ├── stable.go
│   ├── leg.go
│   ├── triangle.go
│   ├── resolver.go
│   └── csv.go
│
├── builder/
│   └── triangles.go
│
├── kucoin/
│   └── markets.go   // LoadMarkets()
│
└── main.go

✅ common/market.go
package common

type Market struct {
	Symbol string

	Base  string
	Quote string

	EnableTrading bool

	BaseMinSize  float64
	QuoteMinSize float64

	BaseIncrement  float64
	QuoteIncrement float64
	PriceIncrement float64

	OrderTypes []string
}

func (m Market) HasMarketOrder() bool {
	for _, t := range m.OrderTypes {
		if t == "MARKET" {
			return true
		}
	}
	return false
}

✅ common/stable.go
package common

var StableCoins = map[string]bool{
	"USDT": true,
	"USDC": true,
	"BUSD": true,
	"DAI":  true,
	"TUSD": true,
	"FDUSD": true,
}

✅ common/leg.go
package common

type Leg struct {
	From   string
	To     string
	Symbol string
	Side   string // BUY or SELL
}

✅ common/resolver.go
package common

func ResolveLeg(from, to string, markets map[string]Market) (*Leg, bool) {

	// BUY → to/from
	if m, ok := markets[to+"-"+from]; ok {
		if m.EnableTrading && m.HasMarketOrder() {
			return &Leg{
				From:   from,
				To:     to,
				Symbol: m.Symbol,
				Side:   "BUY",
			}, true
		}
	}

	// SELL → from/to
	if m, ok := markets[from+"-"+to]; ok {
		if m.EnableTrading && m.HasMarketOrder() {
			return &Leg{
				From:   from,
				To:     to,
				Symbol: m.Symbol,
				Side:   "SELL",
			}, true
		}
	}

	return nil, false
}

✅ common/triangle.go
package common

import "strings"

type Triangle struct {
	A string
	X string
	Y string

	Legs []Leg
}

func (t Triangle) Key() string {
	parts := make([]string, 0, 3)
	for _, l := range t.Legs {
		parts = append(parts, l.Side+":"+l.Symbol)
	}
	return strings.Join(parts, "|")
}

✅ common/csv.go
package common

import (
	"encoding/csv"
	"os"
)

func SaveTrianglesCSV(path string, list []Triangle) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	w.Write([]string{
		"A", "X", "Y",
		"LEG1", "LEG2", "LEG3",
	})

	for _, t := range list {
		w.Write([]string{
			t.A,
			t.X,
			t.Y,
			t.Legs[0].Side + " " + t.Legs[0].Symbol,
			t.Legs[1].Side + " " + t.Legs[1].Symbol,
			t.Legs[2].Side + " " + t.Legs[2].Symbol,
		})
	}

	return nil
}

✅ builder/triangles.go
package builder

import (
	"exchange/common"
)

func buildTriangle(a, x, y string, markets map[string]common.Market) (*common.Triangle, bool) {

	l1, ok := common.ResolveLeg(a, x, markets)
	if !ok {
		return nil, false
	}

	l2, ok := common.ResolveLeg(x, y, markets)
	if !ok {
		return nil, false
	}

	l3, ok := common.ResolveLeg(y, a, markets)
	if !ok {
		return nil, false
	}

	return &common.Triangle{
		A: a,
		X: x,
		Y: y,
		Legs: []common.Leg{*l1, *l2, *l3},
	}, true
}

✅ builder/triangles.go (основной генератор)
package builder

import "exchange/common"

func BuildTriangles(
	markets map[string]common.Market,
) []common.Triangle {

	var result []common.Triangle
	seen := map[string]bool{}

	assets := map[string]bool{}

	for _, m := range markets {
		assets[m.Base] = true
		assets[m.Quote] = true
	}

	for a := range assets {
		if !common.StableCoins[a] {
			continue
		}

		for x := range assets {
			for y := range assets {

				if x == y || x == a || y == a {
					continue
				}

				// A-X-Y-A
				if t, ok := buildTriangle(a, x, y, markets); ok {
					key := t.Key()
					if !seen[key] {
						seen[key] = true
						result = append(result, *t)
					}
				}

				// A-Y-X-A
				if t, ok := buildTriangle(a, y, x, markets); ok {
					key := t.Key()
					if !seen[key] {
						seen[key] = true
						result = append(result, *t)
					}
				}
			}
		}
	}

	return result
}

✅ main.go
package main

import (
	"exchange/builder"
	"exchange/common"
	"exchange/kucoin"
)

func main() {

	markets := kucoin.LoadMarkets()

	triangles := builder.BuildTriangles(markets)

	common.SaveTrianglesCSV("data/kucoin_triangles.csv", triangles)
}




