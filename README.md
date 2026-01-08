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




type Calculator struct {
	mem       *store.MemoryStore
	triangles []Triangle
}


func NewCalculator(
	mem *store.MemoryStore,
	triangles []Triangle,
) *Calculator {
	return &Calculator{
		mem:       mem,
		triangles: triangles,
	}
}


package calculator

import (
	"encoding/csv"
	"os"
	"strings"
)

type Triangle struct {
	A, B, C       string
	Leg1, Leg2, Leg3 string
}

func ParseTrianglesFromCSV(path string) ([]Triangle, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	rows, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	var res []Triangle
	for _, row := range rows[1:] {
		if len(row) < 6 {
			continue
		}

		res = append(res, Triangle{
			A:    strings.TrimSpace(row[0]),
			B:    strings.TrimSpace(row[1]),
			C:    strings.TrimSpace(row[2]),
			Leg1: strings.TrimSpace(row[3]),
			Leg2: strings.TrimSpace(row[4]),
			Leg3: strings.TrimSpace(row[5]),
		})
	}

	return res, nil
}




func main() {
	out := make(chan *models.MarketData, 100_000)

	mem := store.NewMemoryStore()
	go mem.Run(out)

	kc, err := collector.NewKuCoinCollectorFromCSV(
		"../exchange/data/kucoin_triangles_usdt.csv",
	)
	if err != nil {
		log.Fatal(err)
	}

	if err := kc.Start(out); err != nil {
		log.Fatal(err)
	}

	triangles, err := calculator.ParseTrianglesFromCSV(
		"../exchange/data/kucoin_triangles_usdt.csv",
	)
	if err != nil {
		log.Fatal(err)
	}

	calc := calculator.NewCalculator(mem, triangles)
	go calc.Run()

	select {}
}




