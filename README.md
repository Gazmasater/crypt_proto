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
	triangles []Triangle
	mem       *queue.MemoryStore
	fileLog   *log.Logger
}


func NewCalculator(mem *queue.MemoryStore, triangles []Triangle) *Calculator {
	f, err := os.OpenFile(
		"arb_opportunities.log",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	if err != nil {
		log.Fatalf("failed to open arb log file: %v", err)
	}

	return &Calculator{
		mem:       mem,
		triangles: triangles,
		fileLog:   log.New(f, "", log.LstdFlags),
	}
}


if profit > 0 {
	msg := fmt.Sprintf(
		"[ARB] %s → %s → %s | profit=%.4f%% | volumes: [%.4f / %.4f / %.4f]",
		tri.A, tri.B, tri.C,
		profit*100,
		q1.BidSize, q2.BidSize, q3.BidSize,
	)

	// консоль
	log.Println(msg)

	// файл
	c.fileLog.Println(msg)
}






