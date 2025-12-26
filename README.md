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



package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
)

// Структура под /api/v3/exchangeInfo (нужны только нужные поля)
type exchangeInfo struct {
	Symbols []struct {
		Symbol     string `json:"symbol"`
		BaseAsset  string `json:"baseAsset"`
		QuoteAsset string `json:"quoteAsset"`
		Status     string `json:"status"`
	} `json:"symbols"`
}

// Маркет (одна торговая пара)
type pairMarket struct {
	Symbol string
	Base   string
	Quote  string
}

// Ключ валютной пары без направления (min, max)
type pairKey struct {
	A, B string
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	// 1) Тянем exchangeInfo
	resp, err := http.Get("https://api.mexc.com/api/v3/exchangeInfo")
	if err != nil {
		log.Fatalf("get exchangeInfo: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		log.Fatalf("exchangeInfo status %d: %s", resp.StatusCode, string(b))
	}

	var info exchangeInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		log.Fatalf("decode exchangeInfo: %v", err)
	}

	log.Printf("total symbols from API: %d", len(info.Symbols))

	// 2) Фильтруем маркеты
	markets := make([]pairMarket, 0, len(info.Symbols))
	for _, s := range info.Symbols {
		base := s.BaseAsset
		quote := s.QuoteAsset
		if base == "" || quote == "" {
			continue
		}

		// Мягкий фильтр по статусу: берём всё "торгуемое"
		// (на MEXC это может быть "ENABLED", "TRADING", "1" и т.п.)
		if s.Status != "" &&
			s.Status != "ENABLED" &&
			s.Status != "TRADING" &&
			s.Status != "1" {
			continue
		}

		markets = append(markets, pairMarket{
			Symbol: s.Symbol,
			Base:   base,
			Quote:  quote,
		})
	}
	log.Printf("filtered markets: %d", len(markets))

	// 3) Строим pairmap и граф валют
	pairmap := make(map[pairKey][]pairMarket)
	adj := make(map[string]map[string]struct{})

	addEdge := func(a, b string) {
		if adj[a] == nil {
			adj[a] = make(map[string]struct{})
		}
		adj[a][b] = struct{}{}
	}

	for _, m := range markets {
		a, b := m.Base, m.Quote
		key := pairKey{A: a, B: b}
		if a > b {
			key = pairKey{A: b, B: a}
		}
		pairmap[key] = append(pairmap[key], m)

		addEdge(a, b)
		addEdge(b, a)
	}

	log.Printf("currencies (vertices): %d, pair keys: %d", len(adj), len(pairmap))

	// 4) Индексация валют
	coins := make([]string, 0, len(adj))
	for c := range adj {
		coins = append(coins, c)
	}
	sort.Strings(coins)

	idx := make(map[string]int, len(coins))
	for i, c := range coins {
		idx[c] = i
	}

	neighbors := make([]map[int]struct{}, len(coins))
	for i := range neighbors {
		neighbors[i] = make(map[int]struct{})
	}
	for c, neighs := range adj {
		i := idx[c]
		for nb := range neighs {
			j := idx[nb]
			neighbors[i][j] = struct{}{}
		}
	}

	// 5) Поиск валютных треугольников (A,B,C с A<B<C)
	type triangle struct {
		A, B, C string
	}

	triangles := make([]triangle, 0)

	for i := 0; i < len(coins); i++ {
		ni := neighbors[i]

		for j := range ni {
			if j <= i {
				continue
			}
			nj := neighbors[j]

			for k := range ni {
				if k <= j {
					continue
				}
				if _, ok := nj[k]; ok {
					triangles = append(triangles, triangle{
						A: coins[i],
						B: coins[j],
						C: coins[k],
					})
				}
			}
		}
	}
	log.Printf("found currency triangles: %d", len(triangles))

	// 6) Пишем в CSV реальные маркеты для пар (A,B), (B,C), (A,C)
	outFile := "triangles_markets.csv"
	f, err := os.Create(outFile)
	if err != nil {
		log.Fatalf("create %s: %v", outFile, err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	if err := w.Write([]string{"base1", "quote1", "base2", "quote2", "base3", "quote3"}); err != nil {
		log.Fatalf("write header: %v", err)
	}

	pick := func(x, y string) (pairMarket, bool) {
		key := pairKey{A: x, B: y}
		if x > y {
			key = pairKey{A: y, B: x}
		}
		list := pairmap[key]
		if len(list) == 0 {
			return pairMarket{}, false
		}
		return list[0], true
	}

	count := 0
	for _, t := range triangles {
		m1, ok1 := pick(t.A, t.B)
		m2, ok2 := pick(t.B, t.C)
		m3, ok3 := pick(t.A, t.C)
		if !ok1 || !ok2 || !ok3 {
			continue
		}

		rec := []string{
			m1.Base, m1.Quote,
			m2.Base, m2.Quote,
			m3.Base, m3.Quote,
		}
		if err := w.Write(rec); err != nil {
			log.Fatalf("write record: %v", err)
		}
		count++
	}

	log.Printf("written triangles to %s: %d", outFile, count)
	fmt.Println("Готово, файл:", outFile)
}




