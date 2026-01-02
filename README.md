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



✅ 1. common/stable.go

Список стейблов — используется как фильтр

package common

var StableCoins = map[string]bool{
	"USDT": true,
	"USDC": true,
	"DAI":  true,
	"BUSD": true,
	"TUSD": true,
	"EUR":  true,
}

✅ 2. common/market.go

(оставляем как есть, минимально)

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
}

✅ 3. common/resolver.go

Поиск пары и направление BUY / SELL

package common

func FindLeg(a, b string, markets map[string]Market) (Market, bool) {
	if m, ok := markets[a+"_"+b]; ok {
		return m, true
	}
	if m, ok := markets[b+"_"+a]; ok {
		return m, true
	}
	return Market{}, false
}

func ResolveSide(from, to string, m Market) string {
	if m.Base == to && m.Quote == from {
		return "BUY"
	}
	if m.Base == from && m.Quote == to {
		return "SELL"
	}
	return ""
}

✅ 4. common/triangle.go
package common

type Triangle struct {
	A string
	B string
	C string

	Leg1 string
	Leg2 string
	Leg3 string
}

✅ 5. ДЕДУПЛИКАТОР (ключ треугольника)

Чтобы
USDT → A → B → USDT
и
USDT → B → A → USDT

не дублировались.

package common

import "sort"

func TriangleKey(a, b, c string) string {
	x := []string{b, c}
	sort.Strings(x)
	return a + "|" + x[0] + "|" + x[1]
}

✅ 6. ГЛАВНОЕ — универсальный генератор

📄 builder/triangles.go

package builder

import (
	"exchange/common"
)

func BuildTriangles(
	markets map[string]common.Market,
	anchor string,
) []common.Triangle {

	result := []common.Triangle{}
	seen := map[string]bool{}

	for _, m1 := range markets {
		if !m1.EnableTrading {
			continue
		}

		// шаг 1: A -> B
		var B string
		if m1.Base == anchor {
			B = m1.Quote
		} else if m1.Quote == anchor {
			B = m1.Base
		} else {
			continue
		}

		// ❌ нельзя стейб в середине
		if common.StableCoins[B] {
			continue
		}

		for _, m2 := range markets {
			if !m2.EnableTrading {
				continue
			}

			// шаг 2: B -> C
			var C string
			if m2.Base == B {
				C = m2.Quote
			} else if m2.Quote == B {
				C = m2.Base
			} else {
				continue
			}

			if C == anchor || C == B {
				continue
			}

			if common.StableCoins[C] {
				continue
			}

			// шаг 3: C -> A должен существовать
			l3, ok := common.FindLeg(C, anchor, markets)
			if !ok {
				continue
			}

			l1, ok1 := common.FindLeg(anchor, B, markets)
			l2, ok2 := common.FindLeg(B, C, markets)
			if !ok1 || !ok2 {
				continue
			}

			// 🔒 дедупликация A-X-Y-A / A-Y-X-A
			key := common.TriangleKey(anchor, B, C)
			if seen[key] {
				continue
			}
			seen[key] = true

			t := common.Triangle{
				A: anchor,
				B: B,
				C: C,

				Leg1: common.ResolveSide(anchor, B, l1) + " " + l1.Base + "/" + l1.Quote,
				Leg2: common.ResolveSide(B, C, l2) + " " + l2.Base + "/" + l2.Quote,
				Leg3: common.ResolveSide(C, anchor, l3) + " " + l3.Base + "/" + l3.Quote,
			}

			result = append(result, t)
		}
	}

	return result
}

✅ 7. CSV сохранение
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

	w.Write([]string{"A", "B", "C", "Leg1", "Leg2", "Leg3"})

	for _, t := range list {
		w.Write([]string{
			t.A, t.B, t.C,
			t.Leg1, t.Leg2, t.Leg3,
		})
	}

	return nil
}

✅ 8. main.go
package main

import (
	"exchange/builder"
	"exchange/common"
	"exchange/kucoin"
)

func main() {
	markets := kucoin.LoadMarkets()

	triangles := builder.BuildTriangles(markets, "USDT")

	common.SaveTrianglesCSV("data/kucoin_triangles.csv", triangles)
}



