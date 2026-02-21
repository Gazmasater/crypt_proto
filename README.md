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




go run -race main.go


GOMAXPROCS=8 go run -race main.go


// +imports needed:
//   "encoding/json"
//   "fmt"
//   "os"
//   "strings"

func (s *Scanner) evaluate(triIdx int) {
	s.cooldownMu.Lock()
	if last, ok := s.cooldowns[triIdx]; ok && time.Since(last) < s.cooldownDuration {
		s.cooldownMu.Unlock()
		s.stats.cooldownSkips.Add(1)
		return
	}
	s.cooldownMu.Unlock()

	s.mu.RLock()

	tri := s.triangles[triIdx]
	s.stats.scanned.Add(1)

	fees, ok := s.fees[tri.Exchange]
	if !ok {
		s.mu.RUnlock()
		return
	}

	var books [3]exchange.OrderBook
	var spreads [3]float64
	now := time.Now()

	for i, edge := range tri.Edges {
		key := bookKey(tri.Exchange, edge.InstID)
		book, ok := s.books[key]
		if !ok || len(book.Asks) == 0 || len(book.Bids) == 0 {
			s.mu.RUnlock()
			return
		}

		if now.Sub(book.ReceivedAt) > s.maxBookAge {
			s.mu.RUnlock()
			s.stats.staleSkips.Add(1)
			return
		}

		spread := (book.Asks[0].Price - book.Bids[0].Price) / ((book.Asks[0].Price + book.Bids[0].Price) / 2)
		if spread > s.maxSpreadPct {
			s.mu.RUnlock()
			s.stats.spreadSkips.Add(1)
			return
		}
		spreads[i] = spread * 100
		books[i] = book
	}

	volumeUSDT := s.adaptiveVolume(tri, books)
	if volumeUSDT < s.minVolumeUSDT {
		s.mu.RUnlock()
		s.stats.volumeSkips.Add(1)
		return
	}
	if volumeUSDT > s.maxTradeSizeUSDT {
		volumeUSDT = s.maxTradeSizeUSDT
	}

	startAmount := s.fromUSDT(tri.Exchange, tri.Start, volumeUSDT)
	if startAmount <= 0 {
		s.mu.RUnlock()
		return
	}

	// дальше s.books не нужен
	s.mu.RUnlock()

	takerMul := 1.0 - fees.TakerPct
	makerMul := 1.0 - fees.MakerPct

	takerAmt := startAmount
	makerAmt := startAmount
	var prices, bestPrices, slippages [3]float64

	for i, edge := range tri.Edges {
		if edge.Buy {
			bestPrices[i] = books[i].Asks[0].Price
			out, ok := calcBuyOutput(books[i].Asks, takerAmt)
			if !ok {
				s.stats.noLiquidity.Add(1)
				return
			}
			prices[i] = takerAmt / out
			slippages[i] = (prices[i]/bestPrices[i] - 1.0) * 100
			takerAmt = out * takerMul

			outM, ok := calcBuyOutput(books[i].Asks, makerAmt)
			if !ok {
				return
			}
			makerAmt = outM * makerMul
		} else {
			bestPrices[i] = books[i].Bids[0].Price
			out, ok := calcSellOutput(books[i].Bids, takerAmt)
			if !ok {
				s.stats.noLiquidity.Add(1)
				return
			}
			prices[i] = out / takerAmt
			slippages[i] = (1.0 - prices[i]/bestPrices[i]) * 100
			takerAmt = out * takerMul

			outM, ok := calcSellOutput(books[i].Bids, makerAmt)
			if !ok {
				return
			}
			makerAmt = outM * makerMul
		}
	}

	takerProfit := (takerAmt/startAmount - 1.0) * 100
	makerProfit := (makerAmt/startAmount - 1.0) * 100
	totalSlippage := slippages[0] + slippages[1] + slippages[2]

	// В ФАЙЛ — только положительный профит
	if takerProfit > 0 {
		edgeStr := func(e Edge) string {
			side := "SELL"
			if e.Buy {
				side = "BUY"
			}
			return fmt.Sprintf("%s %s (%s→%s)", side, e.InstID, e.From, e.To)
		}

		type line struct {
			AtUTC          string     `json:"at_utc"`
			Exchange       string     `json:"exchange"`
			TriID          string     `json:"tri_id"`
			Start          string     `json:"start"`
			Legs           string     `json:"legs"`
			VolumeUSDT     float64    `json:"volume_usdt"`
			TakerProfitPct float64    `json:"taker_profit_pct"`
			MakerProfitPct float64    `json:"maker_profit_pct"`
			ProfitUSDT     float64    `json:"profit_usdt"`
			SpreadsPct     [3]float64 `json:"spreads_pct"`
			SlippagesPct   [3]float64 `json:"slippages_pct"`
			TotalSlippage  float64    `json:"total_slippage_pct"`
		}

		path := os.Getenv("TRI_LOG_FILE")
		if path == "" {
			path = "triangles_profit.jsonl"
		}

		legs := strings.Join([]string{
			edgeStr(tri.Edges[0]),
			edgeStr(tri.Edges[1]),
			edgeStr(tri.Edges[2]),
		}, " | ")

		l := line{
			AtUTC:          now.UTC().Format(time.RFC3339Nano),
			Exchange:       tri.Exchange,
			TriID:          tri.ID,
			Start:          tri.Start,
			Legs:           legs,
			VolumeUSDT:     volumeUSDT,
			TakerProfitPct: takerProfit,
			MakerProfitPct: makerProfit,
			ProfitUSDT:     volumeUSDT * takerProfit / 100,
			SpreadsPct:     spreads,
			SlippagesPct:   slippages,
			TotalSlippage:  totalSlippage,
		}

		if b, err := json.Marshal(l); err == nil {
			if f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644); err == nil {
				_, _ = f.Write(append(b, '\n'))
				_ = f.Close()
			}
		}
	}

	if takerProfit >= s.minProfitPct {
		s.stats.opportunities.Add(1)

		s.cooldownMu.Lock()
		s.cooldowns[triIdx] = now
		s.cooldownMu.Unlock()

		s.onOpportunity(Opportunity{
			Type:           OpTriangle,
			Exchange:       tri.Exchange,
			Triangle:       tri,
			TakerProfitPct: takerProfit,
			MakerProfitPct: makerProfit,
			Prices:         prices,
			BestPrices:     bestPrices,
			Slippages:      slippages,
			TotalSlippage:  totalSlippage,
			VolumeUSDT:     volumeUSDT,
			ProfitUSDT:     volumeUSDT * takerProfit / 100,
			Spreads:        spreads,
			Timestamp:      now,
		})
	}
}
