// metrics
package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
)

type Row struct {
	Timestamp int64
	Triangle  string
	ProfitPct float64
	Volume    float64
	AgeMax    float64
}

func main() {
	f, err := os.Open("arb_metrics.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	rows, _ := r.ReadAll()

	header := make(map[string]int)
	for i, h := range rows[0] {
		header[h] = i
	}

	var data []Row

	for _, row := range rows[1:] {
		ts, _ := strconv.ParseInt(row[header["timestamp"]], 10, 64)
		profit, _ := strconv.ParseFloat(row[header["profit_pct"]], 64)
		vol, _ := strconv.ParseFloat(row[header["volume_usdt"]], 64)
		age, _ := strconv.ParseFloat(row[header["age_max_ms"]], 64)

		data = append(data, Row{
			Timestamp: ts,
			Triangle:  row[header["triangle"]],
			ProfitPct: profit,
			Volume:    vol,
			AgeMax:    age,
		})
	}

	// фильтры
	minProfit := 0.0001
	minVolume := 50.0
	maxAge := 150.0

	total := len(data)
	valid := 0

	perTri := make(map[string]int)

	for _, d := range data {
		if d.ProfitPct > minProfit &&
			d.Volume > minVolume &&
			d.AgeMax < maxAge {

			valid++
			perTri[d.Triangle]++
		}
	}

	fmt.Println("=== FILTERED ===")
	fmt.Println("total:", total)
	fmt.Println("valid:", valid)

	if total > 0 {
		fmt.Printf("valid %%: %.6f\n", float64(valid)/float64(total)*100)
	}

	fmt.Println("\n=== TOP TRIANGLES ===")
	for tri, cnt := range perTri {
		if cnt > 10 {
			fmt.Println(tri, cnt)
		}
	}
}
