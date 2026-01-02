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



package common

import "sort"

func CanonicalKey(a, b, c string) string {
	arr := []string{a, b, c}
	sort.Strings(arr)
	return arr[0] + "|" + arr[1] + "|" + arr[2]
}




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

		// A -> B
		var B string
		if m1.Base == anchor {
			B = m1.Quote
		} else if m1.Quote == anchor {
			B = m1.Base
		} else {
			continue
		}

		if common.IsStable(B) {
			continue
		}

		for _, m2 := range markets {
			if !m2.EnableTrading {
				continue
			}

			// B -> C
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

			if common.IsStable(C) {
				continue
			}

			// проверяем замыкание C -> A
			l3, ok := common.FindLeg(C, anchor, markets)
			if !ok {
				continue
			}

			l1, ok1 := common.FindLeg(anchor, B, markets)
			l2, ok2 := common.FindLeg(B, C, markets)
			if !ok1 || !ok2 {
				continue
			}

			// ===== дедуп =====
			key := common.CanonicalKey(anchor, B, C)
			if seen[key] {
				continue
			}
			seen[key] = true

			result = append(result, common.NewTriangle(
				anchor,
				B,
				C,
				l1,
				l2,
				l3,
			))
		}
	}

	return result
}




