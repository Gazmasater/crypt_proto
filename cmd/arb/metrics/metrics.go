// metrics
package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

const (
	maxAgeMs      = 360.0
	maxSpreadMs   = 180.0
	minVolumeUSDT = 50.0
)

type Event struct {
	TSUnixMs            int64
	A                   string
	B                   string
	C                   string
	Triangle            string
	ProfitPct           float64
	ProfitUSDT          float64
	VolumeUSDT          float64
	FinalUSDT           float64
	OpportunityStrength float64
	AgeMinMs            float64
	AgeMaxMs            float64
	AgeSpreadMs         float64

	Leg1Symbol string
	Leg1Side   string
	Leg1AgeMs  float64

	Leg2Symbol string
	Leg2Side   string
	Leg2AgeMs  float64

	Leg3Symbol string
	Leg3Side   string
	Leg3AgeMs  float64
}

type Stats struct {
	Count int
	Min   float64
	Max   float64
	Mean  float64
	P50   float64
	P90   float64
	P95   float64
	P99   float64
}

type RejectStats struct {
	Total           int
	RejectAgeMax    int
	RejectAgeSpread int
	RejectVolume    int
	RejectProfit    int
	Usable          int
}

type TriangleAgg struct {
	Triangle        string
	Count           int
	MeanProfitPct   float64
	MedianProfitPct float64
	MaxProfitPct    float64
	MeanVolumeUSDT  float64
	MeanAgeMaxMs    float64
	MeanAgeSpreadMs float64
}

type SymbolQuality struct {
	Symbol    string
	Count     int
	MeanAgeMs float64
	P50AgeMs  float64
	P95AgeMs  float64
	MaxAgeMs  float64
}

func main() {
	input := "arb_metrics.csv"
	if len(os.Args) > 1 {
		input = os.Args[1]
	}

	events, err := loadEvents(input)
	if err != nil {
		log.Fatalf("load events: %v", err)
	}
	if len(events) == 0 {
		log.Fatalf("no events in %s", input)
	}

	printBasic(events)
	printProfitStats(events)
	printAgeStats(events)

	rejects, usable := classify(events)
	printRejects(rejects)

	printTopTrianglesByCount(events, 20)
	printTopTrianglesByMeanProfit(events, 20)
	printSymbolQuality(events, usable, 20)

	if err := exportReports("arb_reports", events, usable, rejects); err != nil {
		log.Fatalf("export reports: %v", err)
	}

	fmt.Println("\nreports saved to ./arb_reports")
}

func loadEvents(path string) ([]Event, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.FieldsPerRecord = -1

	rows, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(rows) < 2 {
		return nil, nil
	}

	header := make(map[string]int, len(rows[0]))
	for i, col := range rows[0] {
		header[strings.TrimSpace(col)] = i
	}

	events := make([]Event, 0, len(rows)-1)

	for _, row := range rows[1:] {
		if isRowEmpty(row) {
			continue
		}

		e := Event{
			TSUnixMs:            getInt64(row, header, "ts_unix_ms"),
			A:                   getString(row, header, "A"),
			B:                   getString(row, header, "B"),
			C:                   getString(row, header, "C"),
			ProfitPct:           getFloat(row, header, "profit_pct"),
			ProfitUSDT:          getFloat(row, header, "profit_usdt"),
			VolumeUSDT:          getFloat(row, header, "volume_usdt"),
			FinalUSDT:           getFloat(row, header, "final_usdt"),
			OpportunityStrength: getFloat(row, header, "opportunity_strength"),
			AgeMinMs:            getFloat(row, header, "age_min_ms"),
			AgeMaxMs:            getFloat(row, header, "age_max_ms"),
			AgeSpreadMs:         getFloat(row, header, "age_spread_ms"),

			Leg1Symbol: getString(row, header, "leg1_symbol"),
			Leg1Side:   getString(row, header, "leg1_side"),
			Leg1AgeMs:  getFloat(row, header, "leg1_age_ms"),

			Leg2Symbol: getString(row, header, "leg2_symbol"),
			Leg2Side:   getString(row, header, "leg2_side"),
			Leg2AgeMs:  getFloat(row, header, "leg2_age_ms"),

			Leg3Symbol: getString(row, header, "leg3_symbol"),
			Leg3Side:   getString(row, header, "leg3_side"),
			Leg3AgeMs:  getFloat(row, header, "leg3_age_ms"),
		}

		e.Triangle = e.A + "->" + e.B + "->" + e.C
		events = append(events, e)
	}

	return events, nil
}

func classify(events []Event) (RejectStats, []Event) {
	var rs RejectStats
	usable := make([]Event, 0, len(events))

	for _, e := range events {
		rs.Total++

		rejected := false

		if e.A != "USDT" {
			rejected = true
		}

		if e.AgeMaxMs > maxAgeMs {
			rs.RejectAgeMax++
			rejected = true
		}
		if e.AgeSpreadMs > maxSpreadMs {
			rs.RejectAgeSpread++
			rejected = true
		}
		if e.VolumeUSDT < minVolumeUSDT {
			rs.RejectVolume++
			rejected = true
		}
		if e.ProfitPct <= 0 {
			rs.RejectProfit++
			rejected = true
		}

		if !rejected {
			rs.Usable++
			usable = append(usable, e)
		}
	}

	return rs, usable
}

func printBasic(events []Event) {
	fmt.Println("=== BASIC ===")
	fmt.Printf("rows: %d\n", len(events))

	minTS := int64(0)
	maxTS := int64(0)
	seenTS := false
	triSet := make(map[string]struct{})

	for _, e := range events {
		if e.TSUnixMs > 0 {
			if !seenTS {
				minTS = e.TSUnixMs
				maxTS = e.TSUnixMs
				seenTS = true
			} else {
				if e.TSUnixMs < minTS {
					minTS = e.TSUnixMs
				}
				if e.TSUnixMs > maxTS {
					maxTS = e.TSUnixMs
				}
			}
		}
		triSet[e.Triangle] = struct{}{}
	}

	if seenTS && maxTS > minTS {
		durationSec := float64(maxTS-minTS) / 1000.0
		eventsPerMin := float64(len(events)) / durationSec * 60.0
		fmt.Printf("duration_sec: %.2f\n", durationSec)
		fmt.Printf("events_per_min: %.2f\n", eventsPerMin)
	}

	fmt.Printf("unique_triangles: %d\n", len(triSet))
}

func printProfitStats(events []Event) {
	fmt.Println("\n=== PROFIT ===")
	printOneStats("profit_pct", collect(events, func(e Event) float64 { return e.ProfitPct }))
	printOneStats("profit_usdt", collect(events, func(e Event) float64 { return e.ProfitUSDT }))
	printOneStats("volume_usdt", collect(events, func(e Event) float64 { return e.VolumeUSDT }))
	printOneStats("opportunity_strength", collect(events, func(e Event) float64 { return e.OpportunityStrength }))

	total := float64(len(events))
	if total == 0 {
		return
	}

	var gt0, gt01, gt02 int
	for _, e := range events {
		if e.ProfitPct > 0 {
			gt0++
		}
		if e.ProfitPct > 0.001 {
			gt01++
		}
		if e.ProfitPct > 0.002 {
			gt02++
		}
	}

	fmt.Printf("\nshare(profit_pct > 0): %.6f\n", float64(gt0)/total)
	fmt.Printf("share(profit_pct > 0.001): %.6f\n", float64(gt01)/total)
	fmt.Printf("share(profit_pct > 0.002): %.6f\n", float64(gt02)/total)
}

func printAgeStats(events []Event) {
	fmt.Println("\n=== AGE ===")
	printOneStats("age_min_ms", collect(events, func(e Event) float64 { return e.AgeMinMs }))
	printOneStats("age_max_ms", collect(events, func(e Event) float64 { return e.AgeMaxMs }))
	printOneStats("age_spread_ms", collect(events, func(e Event) float64 { return e.AgeSpreadMs }))

	total := float64(len(events))
	if total == 0 {
		return
	}

	var max180, max360, maxOver360 int
	var spread50, spread180, spreadOver180 int

	for _, e := range events {
		if e.AgeMaxMs <= 180 {
			max180++
		}
		if e.AgeMaxMs <= 360 {
			max360++
		}
		if e.AgeMaxMs > 360 {
			maxOver360++
		}

		if e.AgeSpreadMs <= 50 {
			spread50++
		}
		if e.AgeSpreadMs <= 180 {
			spread180++
		}
		if e.AgeSpreadMs > 180 {
			spreadOver180++
		}
	}

	fmt.Printf("\nshare(age_max_ms <= 180): %.6f\n", float64(max180)/total)
	fmt.Printf("share(age_max_ms <= 360): %.6f\n", float64(max360)/total)
	fmt.Printf("share(age_max_ms > 360): %.6f\n", float64(maxOver360)/total)

	fmt.Printf("\nshare(age_spread_ms <= 50): %.6f\n", float64(spread50)/total)
	fmt.Printf("share(age_spread_ms <= 180): %.6f\n", float64(spread180)/total)
	fmt.Printf("share(age_spread_ms > 180): %.6f\n", float64(spreadOver180)/total)
}

func printRejects(rs RejectStats) {
	fmt.Println("\n=== REJECTS ===")
	fmt.Printf("total: %d\n", rs.Total)
	fmt.Printf("reject_age_max: %d\n", rs.RejectAgeMax)
	fmt.Printf("reject_age_spread: %d\n", rs.RejectAgeSpread)
	fmt.Printf("reject_volume: %d\n", rs.RejectVolume)
	fmt.Printf("reject_profit: %d\n", rs.RejectProfit)
	fmt.Printf("usable: %d\n", rs.Usable)

	if rs.Total > 0 {
		fmt.Printf("usable_ratio: %.6f\n", float64(rs.Usable)/float64(rs.Total))
	}
}

func printTopTrianglesByCount(events []Event, topN int) {
	fmt.Println("\n=== TOP TRIANGLES BY COUNT ===")

	cnt := make(map[string]int)
	for _, e := range events {
		if e.A != "USDT" {
			continue
		}
		cnt[e.Triangle]++
	}

	type kv struct {
		Key   string
		Value int
	}

	arr := make([]kv, 0, len(cnt))
	for k, v := range cnt {
		arr = append(arr, kv{Key: k, Value: v})
	}

	sort.Slice(arr, func(i, j int) bool {
		if arr[i].Value == arr[j].Value {
			return arr[i].Key < arr[j].Key
		}
		return arr[i].Value > arr[j].Value
	})

	limit := min(topN, len(arr))
	for i := 0; i < limit; i++ {
		fmt.Printf("%2d. %-32s %d\n", i+1, arr[i].Key, arr[i].Value)
	}
}

func printTopTrianglesByMeanProfit(events []Event, topN int) {
	fmt.Println("\n=== TOP TRIANGLES BY MEAN PROFIT ===")

	group := make(map[string][]Event)
	for _, e := range events {
		if e.A != "USDT" {
			continue
		}
		group[e.Triangle] = append(group[e.Triangle], e)
	}

	aggs := make([]TriangleAgg, 0, len(group))

	for tri, rows := range group {
		profits := make([]float64, 0, len(rows))
		var sumProfit, sumVolume, sumAgeMax, sumAgeSpread float64
		maxProfit := -math.MaxFloat64

		for _, e := range rows {
			profits = append(profits, e.ProfitPct)
			sumProfit += e.ProfitPct
			sumVolume += e.VolumeUSDT
			sumAgeMax += e.AgeMaxMs
			sumAgeSpread += e.AgeSpreadMs
			if e.ProfitPct > maxProfit {
				maxProfit = e.ProfitPct
			}
		}

		sort.Float64s(profits)

		aggs = append(aggs, TriangleAgg{
			Triangle:        tri,
			Count:           len(rows),
			MeanProfitPct:   sumProfit / float64(len(rows)),
			MedianProfitPct: percentileSorted(profits, 0.50),
			MaxProfitPct:    maxProfit,
			MeanVolumeUSDT:  sumVolume / float64(len(rows)),
			MeanAgeMaxMs:    sumAgeMax / float64(len(rows)),
			MeanAgeSpreadMs: sumAgeSpread / float64(len(rows)),
		})
	}

	sort.Slice(aggs, func(i, j int) bool {
		if aggs[i].MeanProfitPct == aggs[j].MeanProfitPct {
			return aggs[i].Count > aggs[j].Count
		}
		return aggs[i].MeanProfitPct > aggs[j].MeanProfitPct
	})

	limit := min(topN, len(aggs))
	for i := 0; i < limit; i++ {
		a := aggs[i]
		fmt.Printf(
			"%2d. %-32s count=%d mean_profit_pct=%.6f median=%.6f max=%.6f mean_vol=%.2f mean_age_max=%.2f mean_age_spread=%.2f\n",
			i+1,
			a.Triangle,
			a.Count,
			a.MeanProfitPct,
			a.MedianProfitPct,
			a.MaxProfitPct,
			a.MeanVolumeUSDT,
			a.MeanAgeMaxMs,
			a.MeanAgeSpreadMs,
		)
	}
}

func printSymbolQuality(events, usable []Event, topN int) {
	fmt.Println("\n=== SYMBOL QUALITY (ALL EVENTS) ===")
	printOneSymbolQuality("leg1_symbol", collectSymbolQuality(events, func(e Event) (string, float64) {
		return e.Leg1Symbol, e.Leg1AgeMs
	}), topN)
	printOneSymbolQuality("leg2_symbol", collectSymbolQuality(events, func(e Event) (string, float64) {
		return e.Leg2Symbol, e.Leg2AgeMs
	}), topN)
	printOneSymbolQuality("leg3_symbol", collectSymbolQuality(events, func(e Event) (string, float64) {
		return e.Leg3Symbol, e.Leg3AgeMs
	}), topN)

	fmt.Println("\n=== SYMBOL QUALITY (USABLE EVENTS) ===")
	printOneSymbolQuality("leg1_symbol", collectSymbolQuality(usable, func(e Event) (string, float64) {
		return e.Leg1Symbol, e.Leg1AgeMs
	}), topN)
	printOneSymbolQuality("leg2_symbol", collectSymbolQuality(usable, func(e Event) (string, float64) {
		return e.Leg2Symbol, e.Leg2AgeMs
	}), topN)
	printOneSymbolQuality("leg3_symbol", collectSymbolQuality(usable, func(e Event) (string, float64) {
		return e.Leg3Symbol, e.Leg3AgeMs
	}), topN)
}

func collectSymbolQuality(events []Event, f func(Event) (string, float64)) []SymbolQuality {
	m := make(map[string][]float64)

	for _, e := range events {
		if e.A != "USDT" {
			continue
		}

		symbol, age := f(e)
		if symbol == "" {
			continue
		}

		m[symbol] = append(m[symbol], age)
	}

	out := make([]SymbolQuality, 0, len(m))
	for sym, ages := range m {
		sort.Float64s(ages)

		sum := 0.0
		for _, a := range ages {
			sum += a
		}

		out = append(out, SymbolQuality{
			Symbol:    sym,
			Count:     len(ages),
			MeanAgeMs: sum / float64(len(ages)),
			P50AgeMs:  percentileSorted(ages, 0.50),
			P95AgeMs:  percentileSorted(ages, 0.95),
			MaxAgeMs:  ages[len(ages)-1],
		})
	}

	sort.Slice(out, func(i, j int) bool {
		if out[i].Count == out[j].Count {
			return out[i].MeanAgeMs < out[j].MeanAgeMs
		}
		return out[i].Count > out[j].Count
	})

	return out
}

func printOneSymbolQuality(name string, rows []SymbolQuality, topN int) {
	fmt.Printf("\n%s:\n", name)

	limit := min(topN, len(rows))
	for i := 0; i < limit; i++ {
		r := rows[i]
		fmt.Printf(
			"%2d. %-20s count=%d mean_age=%.2f p50=%.2f p95=%.2f max=%.2f\n",
			i+1,
			r.Symbol,
			r.Count,
			r.MeanAgeMs,
			r.P50AgeMs,
			r.P95AgeMs,
			r.MaxAgeMs,
		)
	}
}

func exportReports(dir string, events, usable []Event, rejects RejectStats) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	if err := writeEventsCSV(filepath.Join(dir, "usable_subset.csv"), usable); err != nil {
		return err
	}

	if err := writeTriangleCountCSV(filepath.Join(dir, "top_triangles_by_count.csv"), events); err != nil {
		return err
	}

	if err := writeTriangleSummaryCSV(filepath.Join(dir, "usable_triangles_summary.csv"), usable); err != nil {
		return err
	}

	if err := writeRejectsCSV(filepath.Join(dir, "reject_summary.csv"), rejects); err != nil {
		return err
	}

	if err := writeSymbolQualityCSV(
		filepath.Join(dir, "symbol_quality_leg1.csv"),
		collectSymbolQuality(events, func(e Event) (string, float64) {
			return e.Leg1Symbol, e.Leg1AgeMs
		}),
	); err != nil {
		return err
	}

	if err := writeSymbolQualityCSV(
		filepath.Join(dir, "symbol_quality_leg2.csv"),
		collectSymbolQuality(events, func(e Event) (string, float64) {
			return e.Leg2Symbol, e.Leg2AgeMs
		}),
	); err != nil {
		return err
	}

	if err := writeSymbolQualityCSV(
		filepath.Join(dir, "symbol_quality_leg3.csv"),
		collectSymbolQuality(events, func(e Event) (string, float64) {
			return e.Leg3Symbol, e.Leg3AgeMs
		}),
	); err != nil {
		return err
	}

	return nil
}

func writeEventsCSV(path string, events []Event) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	header := []string{
		"ts_unix_ms", "A", "B", "C", "triangle",
		"profit_pct", "profit_usdt", "volume_usdt", "final_usdt", "opportunity_strength",
		"age_min_ms", "age_max_ms", "age_spread_ms",
		"leg1_symbol", "leg1_side", "leg1_age_ms",
		"leg2_symbol", "leg2_side", "leg2_age_ms",
		"leg3_symbol", "leg3_side", "leg3_age_ms",
	}
	if err := w.Write(header); err != nil {
		return err
	}

	for _, e := range events {
		row := []string{
			strconv.FormatInt(e.TSUnixMs, 10),
			e.A,
			e.B,
			e.C,
			e.Triangle,
			ff(e.ProfitPct),
			ff(e.ProfitUSDT),
			ff(e.VolumeUSDT),
			ff(e.FinalUSDT),
			ff(e.OpportunityStrength),
			ff(e.AgeMinMs),
			ff(e.AgeMaxMs),
			ff(e.AgeSpreadMs),
			e.Leg1Symbol,
			e.Leg1Side,
			ff(e.Leg1AgeMs),
			e.Leg2Symbol,
			e.Leg2Side,
			ff(e.Leg2AgeMs),
			e.Leg3Symbol,
			e.Leg3Side,
			ff(e.Leg3AgeMs),
		}

		if err := w.Write(row); err != nil {
			return err
		}
	}

	return w.Error()
}

func writeRejectsCSV(path string, rs RejectStats) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	if err := w.Write([]string{"metric", "value"}); err != nil {
		return err
	}

	rows := [][]string{
		{"total", strconv.Itoa(rs.Total)},
		{"reject_age_max", strconv.Itoa(rs.RejectAgeMax)},
		{"reject_age_spread", strconv.Itoa(rs.RejectAgeSpread)},
		{"reject_volume", strconv.Itoa(rs.RejectVolume)},
		{"reject_profit", strconv.Itoa(rs.RejectProfit)},
		{"usable", strconv.Itoa(rs.Usable)},
	}

	for _, row := range rows {
		if err := w.Write(row); err != nil {
			return err
		}
	}

	return w.Error()
}

func writeTriangleCountCSV(path string, events []Event) error {
	cnt := make(map[string]int)
	for _, e := range events {
		if e.A != "USDT" {
			continue
		}
		cnt[e.Triangle]++
	}

	type row struct {
		Triangle string
		Count    int
	}

	rows := make([]row, 0, len(cnt))
	for k, v := range cnt {
		rows = append(rows, row{Triangle: k, Count: v})
	}

	sort.Slice(rows, func(i, j int) bool {
		if rows[i].Count == rows[j].Count {
			return rows[i].Triangle < rows[j].Triangle
		}
		return rows[i].Count > rows[j].Count
	})

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	if err := w.Write([]string{"triangle", "count"}); err != nil {
		return err
	}

	for _, r := range rows {
		if err := w.Write([]string{r.Triangle, strconv.Itoa(r.Count)}); err != nil {
			return err
		}
	}

	return w.Error()
}

func writeTriangleSummaryCSV(path string, events []Event) error {
	group := make(map[string][]Event)
	for _, e := range events {
		if e.A != "USDT" {
			continue
		}
		group[e.Triangle] = append(group[e.Triangle], e)
	}

	rows := make([]TriangleAgg, 0, len(group))
	for tri, evs := range group {
		profits := make([]float64, 0, len(evs))
		var sumProfit, sumVol, sumAgeMax, sumAgeSpread float64
		maxProfit := -math.MaxFloat64

		for _, e := range evs {
			profits = append(profits, e.ProfitPct)
			sumProfit += e.ProfitPct
			sumVol += e.VolumeUSDT
			sumAgeMax += e.AgeMaxMs
			sumAgeSpread += e.AgeSpreadMs
			if e.ProfitPct > maxProfit {
				maxProfit = e.ProfitPct
			}
		}

		sort.Float64s(profits)

		rows = append(rows, TriangleAgg{
			Triangle:        tri,
			Count:           len(evs),
			MeanProfitPct:   sumProfit / float64(len(evs)),
			MedianProfitPct: percentileSorted(profits, 0.50),
			MaxProfitPct:    maxProfit,
			MeanVolumeUSDT:  sumVol / float64(len(evs)),
			MeanAgeMaxMs:    sumAgeMax / float64(len(evs)),
			MeanAgeSpreadMs: sumAgeSpread / float64(len(evs)),
		})
	}

	sort.Slice(rows, func(i, j int) bool {
		if rows[i].MeanProfitPct == rows[j].MeanProfitPct {
			return rows[i].Count > rows[j].Count
		}
		return rows[i].MeanProfitPct > rows[j].MeanProfitPct
	})

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	header := []string{
		"triangle",
		"count",
		"mean_profit_pct",
		"median_profit_pct",
		"max_profit_pct",
		"mean_volume_usdt",
		"mean_age_max_ms",
		"mean_age_spread_ms",
	}
	if err := w.Write(header); err != nil {
		return err
	}

	for _, r := range rows {
		row := []string{
			r.Triangle,
			strconv.Itoa(r.Count),
			ff(r.MeanProfitPct),
			ff(r.MedianProfitPct),
			ff(r.MaxProfitPct),
			ff(r.MeanVolumeUSDT),
			ff(r.MeanAgeMaxMs),
			ff(r.MeanAgeSpreadMs),
		}
		if err := w.Write(row); err != nil {
			return err
		}
	}

	return w.Error()
}

func writeSymbolQualityCSV(path string, rows []SymbolQuality) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	if err := w.Write([]string{
		"symbol",
		"count",
		"mean_age_ms",
		"p50_age_ms",
		"p95_age_ms",
		"max_age_ms",
	}); err != nil {
		return err
	}

	for _, r := range rows {
		row := []string{
			r.Symbol,
			strconv.Itoa(r.Count),
			ff(r.MeanAgeMs),
			ff(r.P50AgeMs),
			ff(r.P95AgeMs),
			ff(r.MaxAgeMs),
		}
		if err := w.Write(row); err != nil {
			return err
		}
	}

	return w.Error()
}

func printOneStats(name string, values []float64) {
	s := calcStats(values)
	fmt.Printf(
		"\n%s:\ncount=%d min=%.6f max=%.6f mean=%.6f p50=%.6f p90=%.6f p95=%.6f p99=%.6f\n",
		name,
		s.Count,
		s.Min,
		s.Max,
		s.Mean,
		s.P50,
		s.P90,
		s.P95,
		s.P99,
	)
}

func calcStats(values []float64) Stats {
	if len(values) == 0 {
		return Stats{}
	}

	cp := append([]float64(nil), values...)
	sort.Float64s(cp)

	sum := 0.0
	for _, v := range cp {
		sum += v
	}

	return Stats{
		Count: len(cp),
		Min:   cp[0],
		Max:   cp[len(cp)-1],
		Mean:  sum / float64(len(cp)),
		P50:   percentileSorted(cp, 0.50),
		P90:   percentileSorted(cp, 0.90),
		P95:   percentileSorted(cp, 0.95),
		P99:   percentileSorted(cp, 0.99),
	}
}

func percentileSorted(sorted []float64, p float64) float64 {
	if len(sorted) == 0 {
		return 0
	}
	if len(sorted) == 1 {
		return sorted[0]
	}
	if p <= 0 {
		return sorted[0]
	}
	if p >= 1 {
		return sorted[len(sorted)-1]
	}

	pos := p * float64(len(sorted)-1)
	lo := int(math.Floor(pos))
	hi := int(math.Ceil(pos))

	if lo == hi {
		return sorted[lo]
	}

	w := pos - float64(lo)
	return sorted[lo]*(1-w) + sorted[hi]*w
}

func collect(events []Event, f func(Event) float64) []float64 {
	out := make([]float64, 0, len(events))
	for _, e := range events {
		out = append(out, f(e))
	}
	return out
}

func getString(row []string, header map[string]int, key string) string {
	idx, ok := header[key]
	if !ok || idx >= len(row) {
		return ""
	}
	return strings.TrimSpace(row[idx])
}

func getFloat(row []string, header map[string]int, key string) float64 {
	s := getString(row, header, key)
	if s == "" {
		return 0
	}

	v, err := strconv.ParseFloat(strings.ReplaceAll(s, ",", "."), 64)
	if err != nil {
		return 0
	}
	return v
}

func getInt64(row []string, header map[string]int, key string) int64 {
	s := getString(row, header, key)
	if s == "" {
		return 0
	}

	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}
	return v
}

func isRowEmpty(row []string) bool {
	for _, s := range row {
		if strings.TrimSpace(s) != "" {
			return false
		}
	}
	return true
}

func ff(v float64) string {
	return strconv.FormatFloat(v, 'f', 6, 64)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
