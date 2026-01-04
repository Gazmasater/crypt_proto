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




func NewKuCoinCollectorFromCSV(csvFile string) (*KuCoinCollector, error) {
	// 1. Получаем пары из CSV
	pairs2D, err := readPairsFromCSV(csvFile)
	if err != nil {
		return nil, err
	}

	// Преобразуем пары [2]string -> "BASE-QUOTE"
	pairs := make([]string, 0, len(pairs2D))
	for _, p := range pairs2D {
		pairs = append(pairs, p[0]+"-"+p[1])
	}

	// 2. Получаем все символы KuCoin
	rawSymbols, err := FetchKuCoinSymbols()
	if err != nil {
		return nil, err
	}

	// Преобразуем в map для фильтра
	allSymbols := make(map[string]struct{ EnableTrading bool })
	for _, s := range rawSymbols {
		allSymbols[s] = struct{ EnableTrading bool }{EnableTrading: true}
	}

	// 3. Фильтруем
	symbols := FilterPairsForTriangles(pairs, allSymbols)
	if len(symbols) == 0 {
		return nil, fmt.Errorf("no valid symbols to subscribe")
	}

	// 4. Создаём коллектор
	ctx, cancel := context.WithCancel(context.Background())
	return &KuCoinCollector{
		ctx:      ctx,
		cancel:   cancel,
		symbols:  symbols,
		lastData: make(map[string]struct{ Bid, Ask, BidSize, AskSize float64 }),
	}, nil
}



