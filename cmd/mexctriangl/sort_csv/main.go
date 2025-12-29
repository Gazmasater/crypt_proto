package main

import (
	"encoding/csv"
	"os"
)

func contains(row []string, v string) bool {
	for i := 0; i < 3; i++ { // A, B, C
		if row[i] == v {
			return true
		}
	}
	return false
}

func main() {
	inFile, err := os.Open("triangles.csv")
	if err != nil {
		panic(err)
	}
	defer inFile.Close()

	reader := csv.NewReader(inFile)
	reader.FieldsPerRecord = -1

	rows, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	if len(rows) == 0 {
		return
	}

	header := rows[0]
	data := rows[1:]

	var withUSDT [][]string
	//var withUSDC [][]string
	var rest [][]string

	for _, row := range data {
		hasUSDT := contains(row, "USDT")
		//hasUSDC := contains(row, "USDC")

		switch {
		case hasUSDT:
			withUSDT = append(withUSDT, row)
		//case hasUSDC:
		//	withUSDC = append(withUSDC, row)
		default:
			rest = append(rest, row)
		}
	}

	outFile, err := os.Create("triangles_sorted.csv")
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	writer := csv.NewWriter(outFile)
	defer writer.Flush()

	_ = writer.Write(header)

	for _, r := range withUSDT {
		_ = writer.Write(r)
	}
	//for _, r := range withUSDC {
	//	_ = writer.Write(r)
	//}
	// при необходимости:
	// for _, r := range rest {
	//     _ = writer.Write(r)
	// }
}
