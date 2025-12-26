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

// Ответ KuCoin
type kucoinExchangeInfo struct {
	Code string `json:"code"`
	Data []struct {
		Symbol        string `json:"symbol"`
		BaseCurrency  string `json:"baseCurrency"`
		QuoteCurrency string `json:"quoteCurrency"`
		EnableTrading bool   `json:"enableTrading"`
	} `json:"data"`
}

// Маркет
type pairMarket struct {
	Symbol string
	Base   string
	Quote  string
}

// ключ без направления
type pairKey struct {
	A, B string
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	resp, err := http.Get("https://api.kucoin.com/api/v2/symbols")
	if err != nil {
		log.Fatalf("get symbols: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		log.Fatalf("status %d: %s", resp.StatusCode, string(b))
	}

	var info kucoinExchangeInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		log.Fatalf("decode: %v", err)
	}

	log.Printf("symbols from API: %d", len(info.Data))

	// 1) фильтрация маркетов
	markets := make([]pairMarket, 0, len(info.Data))
	for _, s := range info.Data {
		if !s.EnableTrading {
			continue
		}
		if s.BaseCurrency == "" || s.QuoteCurrency == "" {
			continue
		}

		markets = append(markets, pairMarket{
			Symbol: s.Symbol,
			Base:   s.BaseCurrency,
			Quote:  s.QuoteCurrency,
		})
	}

	log.Printf("filtered markets: %d", len(markets))

	// 2) строим граф валют
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

	log.Printf("currencies: %d, pair keys: %d", len(adj), len(pairmap))

	// 3) индексация валют
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

	// 4) поиск треугольников валют
	type triangle struct {
		A, B, C string
	}

	var triangles []triangle

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

	log.Printf("currency triangles: %d", len(triangles))

	// 5) запись CSV
	outFile := "triangles_markets_kucoin.csv"
	f, err := os.Create(outFile)
	if err != nil {
		log.Fatalf("create %s: %v", outFile, err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	_ = w.Write([]string{
		"base1", "quote1",
		"base2", "quote2",
		"base3", "quote3",
	})

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
	fmt.Println("Готово:", outFile)
}
