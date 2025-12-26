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




csvPath := "triangles_usdt_routes.csv" // твой CSV с белым списком
symbols, err := readSymbolsFromCSV(csvPath)
if err != nil {
	log.Fatal("read CSV symbols:", err)
}
log.Printf("Subscribing to %d symbols", len(symbols))


csvPath := "triangles_usdt_routes.csv" // твой CSV с белым списком
symbols, err := readSymbolsFromCSV(csvPath)
if err != nil {
	log.Fatal("read CSV symbols:", err)
}
log.Printf("Subscribing to %d symbols", len(symbols))








func readSymbolsFromCSV(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	header, err := r.Read()
	if err != nil {
		return nil, err
	}

	// находим индексы колонок leg1/2/3
	idx := map[string]int{}
	for i, h := range header {
		idx[strings.ToLower(strings.TrimSpace(h))] = i
	}

	legs := []string{"leg1_symbol", "leg2_symbol", "leg3_symbol"}
	var indices []int
	for _, l := range legs {
		if i, ok := idx[l]; ok {
			indices = append(indices, i)
		} else {
			return nil, &csv.ParseError{StartLine: 0, Err: csv.ErrFieldCount}
		}
	}

	set := map[string]struct{}{}
	for {
		row, err := r.Read()
		if err != nil {
			break
		}
		for _, i := range indices {
			s := strings.TrimSpace(row[i])
			if s != "" {
				set[s] = struct{}{}
			}
		}
	}

	var symbols []string
	for s := range set {
		symbols = append(symbols, s)
	}

	return symbols, nil
}



