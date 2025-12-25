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



package market

import "strings"

var knownQuotes = []string{
	"USDT",
	"USDC",
	"USD",
	"EUR",
	"BTC",
	"ETH",
}

func NormalizeSymbol(s string) string {
	s = strings.TrimSpace(strings.ToUpper(s))
	if s == "" {
		return s
	}

	// если уже с разделителем
	if strings.ContainsAny(s, "-/") {
		s = strings.ReplaceAll(s, "-", "/")
		parts := strings.Split(s, "/")
		if len(parts) == 2 && parts[0] != "" && parts[1] != "" {
			return parts[0] + "/" + parts[1]
		}
		return s
	}

	// формат BTCUSDT → BTC/USDT
	for _, q := range knownQuotes {
		if strings.HasSuffix(s, q) && len(s) > len(q) {
			base := strings.TrimSuffix(s, q)
			return base + "/" + q
		}
	}

	// неизвестный формат → как есть
	return s
}




package market

func Key(exchange, symbol string) string {
	return exchange + ":" + NormalizeSymbol(symbol)
}



package market

import "strings"

type Pair struct {
	Base  string
	Quote string
}

func ParsePair(s string) Pair {
	if !strings.Contains(s, "/") {
		return Pair{}
	}

	parts := strings.Split(s, "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return Pair{}
	}

	return Pair{
		Base:  parts[0],
		Quote: parts[1],
	}
}


