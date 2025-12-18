mx0vglmT3srN1IS19H
135bb7a7509e4421bad692415c53753b



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




export TRADE_AMOUNT_USDT=100
export FEE_PCT=0.04
export SELL_SAFETY=0.995

export TRIANGLES_FILE=triangles_markets.csv
export TRIANGLES_ENRICHED_FILE=triangles_markets_enriched.csv

go run ./cmd/triangles_enrich_mexc




package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type BookTicker struct {
	Symbol string `json:"symbol"`
	Bid    string `json:"bidPrice"`
	Ask    string `json:"askPrice"`
}

type Side int

const (
	Buy Side = iota
	Sell
)

type Leg struct {
	Symbol     string
	Side       Side
	Step       float64
	MinQty     float64
	MinNotional float64
	Base       string
	Quote      string
}

type Prices struct {
	Bid float64
	Ask float64
}

func main() {
	inPath := "triangles_markets.csv"
	outPath := "triangles_with_minstart.csv"

	rows, header, err := readCSV(inPath)
	if err != nil {
		fatal(err)
	}

	// Индексы нужных колонок
	idx := make(map[string]int)
	for i, h := range header {
		idx[h] = i
	}

	required := []string{
		"symbol1", "symbol2", "symbol3",
		"step1", "step2", "step3",
		"minQty1", "minQty2", "minQty3",
		"minNotional1", "minNotional2", "minNotional3",
		"base1", "quote1", "base2", "quote2", "base3", "quote3",
		"cycle_leg1", "cycle_leg2", "cycle_leg3",
	}
	for _, k := range required {
		if _, ok := idx[k]; !ok {
			fatal(fmt.Errorf("нет колонки %q в CSV", k))
		}
	}

	// Собираем список всех символов, чтобы одним запросом забрать bookTicker
	symbolSet := map[string]struct{}{}
	for _, r := range rows {
		symbolSet[r[idx["symbol1"]]] = struct{}{}
		symbolSet[r[idx["symbol2"]]] = struct{}{}
		symbolSet[r[idx["symbol3"]]] = struct{}{}
	}

	prices, err := fetchAllBookTickers()
	if err != nil {
		fatal(err)
	}

	// Добавим колонки
	outHeader := append(append([]string{}, header...), "ok_for_5", "min_start_usdt_calc")

	outRows := make([][]string, 0, len(rows))

	startUSDT := 5.0

	for _, r := range rows {
		leg1, err := parseLegFromCycle(r[idx["cycle_leg1"]], r, idx, 1)
		if err != nil {
			outRows = append(outRows, appendRowWithErr(r, err)...)
			continue
		}
		leg2, err := parseLegFromCycle(r[idx["cycle_leg2"]], r, idx, 2)
		if err != nil {
			outRows = append(outRows, appendRowWithErr(r, err)...)
			continue
		}
		leg3, err := parseLegFromCycle(r[idx["cycle_leg3"]], r, idx, 3)
		if err != nil {
			outRows = append(outRows, appendRowWithErr(r, err)...)
			continue
		}

		// цены
		p1, ok1 := prices[leg1.Symbol]
		p2, ok2 := prices[leg2.Symbol]
		p3, ok3 := prices[leg3.Symbol]
		if !ok1 || !ok2 || !ok3 {
			err := fmt.Errorf("нет цены для символа: %v%v%v",
				missingSym(leg1.Symbol, ok1),
				missingSym(leg2.Symbol, ok2),
				missingSym(leg3.Symbol, ok3),
			)
			outRows = append(outRows, appendRowWithErr(r, err)...)
			continue
		}

		okFor5 := canExecuteTriangle(startUSDT, leg1, leg2, leg3, p1, p2, p3)

		minStart, okMin := findMinStartUSDT(leg1, leg2, leg3, p1, p2, p3)
		if !okMin {
			// “неисполнимо” даже при большом старте (например, step слишком крупный)
			minStart = math.Inf(1)
		}

		out := append([]string{}, r...)
		out = append(out, boolTo01(okFor5))
		out = append(out, formatFloat(minStart))
		outRows = append(outRows, out)
	}

	if err := writeCSV(outPath, outHeader, outRows); err != nil {
		fatal(err)
	}

	fmt.Printf("OK: записал %s\n", outPath)
}

func readCSV(path string) (rows [][]string, header []string, err error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.FieldsPerRecord = -1

	header, err = r.Read()
	if err != nil {
		return nil, nil, err
	}

	for {
		rec, e := r.Read()
		if e == io.EOF {
			break
		}
		if e != nil {
			return nil, nil, e
		}
		rows = append(rows, rec)
	}
	return rows, header, nil
}

func writeCSV(path string, header []string, rows [][]string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	if err := w.Write(header); err != nil {
		return err
	}
	for _, r := range rows {
		if err := w.Write(r); err != nil {
			return err
		}
	}
	return nil
}

func fetchAllBookTickers() (map[string]Prices, error) {
	// MEXC: https://api.mexc.com/api/v3/ticker/bookTicker
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get("https://api.mexc.com/api/v3/ticker/bookTicker")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("bookTicker http %d", resp.StatusCode)
	}

	var raw []BookTicker
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}

	out := make(map[string]Prices, len(raw))
	for _, t := range raw {
		b, err1 := strconv.ParseFloat(t.Bid, 64)
		a, err2 := strconv.ParseFloat(t.Ask, 64)
		if err1 != nil || err2 != nil || b <= 0 || a <= 0 {
			continue
		}
		out[t.Symbol] = Prices{Bid: b, Ask: a}
	}
	return out, nil
}

// cycle_legX пример: "BDXUSDT:BUY USDT→BDX"
func parseLegFromCycle(cycle string, row []string, idx map[string]int, legNum int) (Leg, error) {
	parts := strings.Split(cycle, ":")
	if len(parts) < 2 {
		return Leg{}, fmt.Errorf("cycle_leg%d кривой формат: %q", legNum, cycle)
	}
	symbol := strings.TrimSpace(parts[0])

	sideStr := strings.TrimSpace(parts[1])
	sideStr = strings.Split(sideStr, " ")[0] // BUY/SELL

	var side Side
	switch strings.ToUpper(sideStr) {
	case "BUY":
		side = Buy
	case "SELL":
		side = Sell
	default:
		return Leg{}, fmt.Errorf("cycle_leg%d неизвестная сторона: %q", legNum, sideStr)
	}

	step := mustFloat(row[idx[fmt.Sprintf("step%d", legNum)]])
	minQty := mustFloat(row[idx[fmt.Sprintf("minQty%d", legNum)]])
	minNotional := mustFloat(row[idx[fmt.Sprintf("minNotional%d", legNum)]])
	base := row[idx[fmt.Sprintf("base%d", legNum)]]
	quote := row[idx[fmt.Sprintf("quote%d", legNum)]]

	// step может быть "0" в твоих данных — это плохо, ставим минимально разумное,
	// но лучше исправить генератор, чтобы step всегда был > 0.
	if step <= 0 {
		step = 1e-12
	}

	return Leg{
		Symbol:      symbol,
		Side:        side,
		Step:        step,
		MinQty:      minQty,
		MinNotional: minNotional,
		Base:        base,
		Quote:       quote,
	}, nil
}

func canExecuteTriangle(startUSDT float64, l1, l2, l3 Leg, p1, p2, p3 Prices) bool {
	amt := startUSDT
	var ok bool

	amt, ok = execLeg(amt, l1, p1)
	if !ok {
		return false
	}
	amt, ok = execLeg(amt, l2, p2)
	if !ok {
		return false
	}
	amt, ok = execLeg(amt, l3, p3)
	return ok && amt > 0
}

// вход amt — это количество "входного актива" для ноги (в терминах cycle: слева от стрелки)
// BUY: вход = quote, получаем base по ASK
// SELL: вход = base, получаем quote по BID
func execLeg(in float64, leg Leg, pr Prices) (out float64, ok bool) {
	switch leg.Side {
	case Buy:
		price := pr.Ask
		// qty base = inQuote / ask, округление вниз по step
		qty := floorToStep(in/price, leg.Step)
		if qty <= 0 {
			return 0, false
		}
		if leg.MinQty > 0 && qty+1e-18 < leg.MinQty {
			return 0, false
		}
		notional := qty * price
		if leg.MinNotional > 0 && notional+1e-12 < leg.MinNotional {
			return 0, false
		}
		return qty, true

	case Sell:
		price := pr.Bid
		qty := floorToStep(in, leg.Step)
		if qty <= 0 {
			return 0, false
		}
		if leg.MinQty > 0 && qty+1e-18 < leg.MinQty {
			return 0, false
		}
		notional := qty * price
		if leg.MinNotional > 0 && notional+1e-12 < leg.MinNotional {
			return 0, false
		}
		return qty * price, true
	default:
		return 0, false
	}
}

func findMinStartUSDT(l1, l2, l3 Leg, p1, p2, p3 Prices) (min float64, ok bool) {
	// если даже на огромном старте не исполняется — выходим
	hi := 1e6 // 1,000,000 USDT как верхняя граница
	if !canExecuteTriangle(hi, l1, l2, l3, p1, p2, p3) {
		return 0, false
	}

	lo := 0.0
	for i := 0; i < 60; i++ { // достаточно для double
		mid := (lo + hi) / 2
		if canExecuteTriangle(mid, l1, l2, l3, p1, p2, p3) {
			hi = mid
		} else {
			lo = mid
		}
	}
	return hi, true
}

func floorToStep(x, step float64) float64 {
	if step <= 0 {
		return x
	}
	n := math.Floor(x / step)
	return n * step
}

func mustFloat(s string) float64 {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return v
}

func boolTo01(b bool) string {
	if b {
		return "1"
	}
	return "0"
}

func formatFloat(v float64) string {
	if math.IsInf(v, 1) {
		return "inf"
	}
	return strconv.FormatFloat(v, 'f', 6, 64)
}

func appendRowWithErr(r []string, err error) [][]string {
	out := append([]string{}, r...)
	out = append(out, "0")
	out = append(out, "inf")
	_ = err // если хочешь — можно добавить ещё колонку "err" и писать сюда
	return [][]string{out}
}

func missingSym(sym string, ok bool) string {
	if ok {
		return ""
	}
	return " " + sym
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, "ERR:", err)
	os.Exit(1)
}

// (на будущее) если захочешь — можно валидировать, что symbolX из CSV совпадает с symbol в cycle_legX.
var _ = errors.New


