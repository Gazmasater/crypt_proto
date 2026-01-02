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

// ===== OKX response =====

type okxInstrumentsResp struct {
	Code string          `json:"code"`
	Data []okxInstrument `json:"data"`
}

type okxInstrument struct {
	InstID   string `json:"instId"`
	BaseCcy  string `json:"baseCcy"`
	QuoteCcy string `json:"quoteCcy"`
	State    string `json:"state"`
}

// ===== internal structs =====

type pairMarket struct {
	Symbol string
	Base   string
	Quote  string
}

type pairKey struct {
	A, B string
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	url := "https://www.okx.com/api/v5/public/instruments?instType=SPOT"

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("get instruments: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		log.Fatalf("bad status %d: %s", resp.StatusCode, string(b))
	}

	var data okxInstrumentsResp
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Fatalf("decode json: %v", err)
	}

	if data.Code != "0" {
		log.Fatalf("api error code=%s", data.Code)
	}

	log.Printf("total instruments from OKX: %d", len(data.Data))

	// ============================
	// 1. фильтрация рынков
	// ============================
	markets := make([]pairMarket, 0, len(data.Data))

	for _, s := range data.Data {
		if s.BaseCcy == "" || s.QuoteCcy == "" {
			continue
		}
		if s.State != "live" {
			continue
		}

		markets = append(markets, pairMarket{
			Symbol: s.InstID,
			Base:   s.BaseCcy,
			Quote:  s.QuoteCcy,
		})
	}

	log.Printf("filtered spot markets: %d", len(markets))

	// ============================
	// 2. строим граф валют
	// ============================
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

	log.Printf("currencies: %d, pair-keys: %d", len(adj), len(pairmap))

	// ============================
	// 3. индексация валют
	// ============================
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

	for c, ns := range adj {
		i := idx[c]
		for n := range ns {
			j := idx[n]
			neighbors[i][j] = struct{}{}
		}
	}

	// ============================
	// 4. поиск треугольников
	// ============================
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

	// ============================
	// 5. запись CSV
	// ============================
	outFile := "triangles_markets_okx.csv"
	f, err := os.Create(outFile)
	if err != nil {
		log.Fatalf("create %s: %v", outFile, err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	if err := w.Write([]string{
		"base1", "quote1",
		"base2", "quote2",
		"base3", "quote3",
	}); err != nil {
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
			log.Fatalf("write row: %v", err)
		}
		count++
	}

	log.Printf("written triangles to %s: %d", outFile, count)
	fmt.Println("Готово:", outFile)
}
