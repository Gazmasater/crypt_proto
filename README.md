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



package builder

import "exchange/common"

// BuildTriangles строит треугольники из всех доступных рынков с учётом anchor.
// Пропускаются стейблкоины, кроме anchor.
// Возвращает все варианты: anchor → B → C → anchor и anchor → C → B → anchor
func BuildTriangles(
	markets map[string]common.Market,
	anchor string,
) []common.Triangle {

	var result []common.Triangle

	for _, m1 := range markets {
		if !m1.EnableTrading {
			continue
		}

		var B string
		if m1.Base == anchor {
			B = m1.Quote
		} else if m1.Quote == anchor {
			B = m1.Base
		} else {
			continue
		}

		if common.IsStable(B) && B != anchor {
			continue
		}

		for _, m2 := range markets {
			if !m2.EnableTrading {
				continue
			}

			var C string
			if m2.Base == B {
				C = m2.Quote
			} else if m2.Quote == B {
				C = m2.Base
			} else {
				continue
			}

			if C == anchor || C == B || (common.IsStable(C) && C != anchor) {
				continue
			}

			// Первая нога: anchor → B
			l1, ok1 := common.FindLeg(anchor, B, markets)
			// Вторая нога: B → C
			l2, ok2 := common.FindLeg(B, C, markets)
			// Третья нога: C → anchor
			l3, ok3 := common.FindLeg(C, anchor, markets)

			if ok1 && ok2 && ok3 {
				t := common.NewTriangle(anchor, B, C, l1, l2, l3)
				result = append(result, t)
			}

			// Вариант в обратном порядке: anchor → C → B → anchor
			l1r, ok1r := common.FindLeg(anchor, C, markets)
			l2r, ok2r := common.FindLeg(C, B, markets)
			l3r, ok3r := common.FindLeg(B, anchor, markets)

			if ok1r && ok2r && ok3r {
				t := common.NewTriangle(anchor, C, B, l1r, l2r, l3r)
				result = append(result, t)
			}
		}
	}

	return result
}



