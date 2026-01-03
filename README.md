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




func readSymbolsFromCSV(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.TrimLeadingSpace = true

	// --- читаем header ---
	header, err := r.Read()
	if err != nil {
		return nil, err
	}

	// map: column name -> index (в lowercase)
	colIndex := make(map[string]int, len(header))
	for i, h := range header {
		colIndex[strings.ToLower(strings.TrimSpace(h))] = i
	}

	// обязательные колонки
	required := []string{"leg1", "leg2", "leg3"}
	idx := make([]int, 0, 3)

	for _, name := range required {
		i, ok := colIndex[name]
		if !ok {
			return nil, fmt.Errorf("CSV: missing required column %q", name)
		}
		idx = append(idx, i)
	}

	// --- собираем уникальные символы ---
	uniq := make(map[string]struct{})

	for {
		row, err := r.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		for _, i := range idx {
			if i >= len(row) {
				continue
			}

			cell := strings.TrimSpace(row[i])
			if cell == "" {
				continue
			}

			// "BUY PEPE/USDT" → "PEPE/USDT"
			parts := strings.Fields(cell)
			if len(parts) < 2 {
				continue
			}

			symbol := parts[len(parts)-1]
			uniq[symbol] = struct{}{}
		}
	}

	// --- map → slice ---
	out := make([]string, 0, len(uniq))
	for s := range uniq {
		out = append(out, s)
	}

	return out, nil
}



