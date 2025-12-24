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

// NormalizeSymbol приводит символ биржи к общему виду BTC/USDT
func NormalizeSymbol(exchange, symbol string) string {
	s := strings.ToUpper(symbol)

	switch exchange {
	case "mexc":
		// BTCUSDT → BTC/USDT
		if len(s) > 4 {
			return s[:len(s)-4] + "/" + s[len(s)-4:]
		}

	case "okx":
		// BTC-USDT → BTC/USDT
		return strings.ReplaceAll(s, "-", "/")

	case "kucoin":
		// BTC-USDT → BTC/USDT
		return strings.ReplaceAll(s, "-", "/")
	}

	return ""
}



package market

import "strings"

// BuildKey формирует ключ для store / redis
// Пример: mexc:BTC/USDT
func BuildKey(exchange, normalizedSymbol string) string {
	return strings.ToLower(exchange) + ":" + normalizedSymbol
}




package market

import "strings"

func ParsePair(normalized string) Pair {
	parts := strings.Split(normalized, "/")
	if len(parts) != 2 {
		return Pair{}
	}
	if parts[0] == "" || parts[1] == "" {
		return Pair{}
	}

	return Pair{
		Base:  parts[0],
		Quote: parts[1],
	}
}




