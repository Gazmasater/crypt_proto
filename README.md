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
	var withUSDC [][]string
	var rest [][]string

	for _, row := range data {
		hasUSDT := contains(row, "USDT")
		hasUSDC := contains(row, "USDC")

		switch {
		case hasUSDT:
			withUSDT = append(withUSDT, row)
		case hasUSDC:
			withUSDC = append(withUSDC, row)
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
	for _, r := range withUSDC {
		_ = writer.Write(r)
	}
	// при необходимости:
	// for _, r := range rest {
	//     _ = writer.Write(r)
	// }
}

