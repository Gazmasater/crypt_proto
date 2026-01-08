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




package collector

import "crypt_proto/pkg/calculator"

// ------------------ В KuCoinCollector ------------------

// Добавляем поле для треугольников
type KuCoinCollector struct {
	ctx       context.Context
	cancel    context.CancelFunc
	wsList    []*kucoinWS
	out       chan<- *models.MarketData
	triangles []calculator.Triangle // <- сюда сохраняем треугольники из CSV
}

// В конструкторе NewKuCoinCollectorFromCSV после чтения CSV:
func NewKuCoinCollectorFromCSV(path string) (*KuCoinCollector, error) {
	symbols, err := readPairsFromCSV(path)
	if err != nil {
		return nil, err
	}
	if len(symbols) == 0 {
		return nil, fmt.Errorf("no symbols")
	}

	// --- формируем треугольники ---
	triangles := calculator.ParseTrianglesFromCSV(path) // <- создаём функцию, которая вернёт []Triangle

	ctx, cancel := context.WithCancel(context.Background())
	var wsList []*kucoinWS
	for i := 0; i < len(symbols); i += maxSubsPerWS {
		end := i + maxSubsPerWS
		if end > len(symbols) {
			end = len(symbols)
		}
		wsList = append(wsList, &kucoinWS{
			id:      len(wsList),
			symbols: symbols[i:end],
			last:    make(map[string][2]float64),
		})
	}

	return &KuCoinCollector{
		ctx:       ctx,
		cancel:    cancel,
		wsList:    wsList,
		triangles: triangles,
	}, nil
}

// ------------------ Метод для калькулятора ------------------
func (kc *KuCoinCollector) Triangles() []calculator.Triangle {
	return kc.triangles
}



[{
	"resource": "/home/gaz358/myprog/crypt_proto/internal/collector/kucoin_collector.go",
	"owner": "_generated_diagnostic_collection_name_#0",
	"code": {
		"value": "UndeclaredImportedName",
		"target": {
			"$mid": 1,
			"path": "/golang.org/x/tools/internal/typesinternal",
			"scheme": "https",
			"authority": "pkg.go.dev",
			"fragment": "UndeclaredImportedName"
		}
	},
	"severity": 8,
	"message": "undefined: calculator.ParseTrianglesFromCSV",
	"source": "compiler",
	"startLineNumber": 58,
	"startColumn": 26,
	"endLineNumber": 58,
	"endColumn": 47,
	"origin": "extHost1"
}]



package calculator

import (
	"encoding/csv"
	"os"
	"strings"
)

// Triangle — структура треугольника
type Triangle struct {
	A, B, C string
	Leg1    string
	Leg2    string
	Leg3    string
}

// ParseTrianglesFromCSV парсит CSV и возвращает треугольники
func ParseTrianglesFromCSV(path string) []Triangle {
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()

	r := csv.NewReader(f)
	rows, err := r.ReadAll()
	if err != nil {
		return nil
	}

	var triangles []Triangle
	for _, row := range rows[1:] {
		if len(row) < 6 {
			continue
		}
		triangles = append(triangles, Triangle{
			A:    row[0],
			B:    row[1],
			C:    row[2],
			Leg1: row[3],
			Leg2: row[4],
			Leg3: row[5],
		})
	}
	return triangles
}





triangles := kc.Triangles()
calc := calculator.NewCalculator(triangles, mem)
go calc.Run()


