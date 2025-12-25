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
	"USDT", "USDC", "USD", "EUR", "BTC", "ETH",
}

// NormalizeSymbol_Full приводит символ к виду BASE/QUOTE
// возвращает "" если формат неполный или некорректный
func NormalizeSymbol_Full(s string) string {
	s = strings.TrimSpace(strings.ToUpper(s))
	if s == "" {
		return ""
	}

	// уже с разделителем
	if strings.ContainsAny(s, "-/") {
		s = strings.ReplaceAll(s, "-", "/")
		parts := strings.Split(s, "/")
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return ""
		}
		return parts[0] + "/" + parts[1]
	}

	// слитные символы, пытаемся найти quote
	for _, q := range knownQuotes {
		if strings.HasSuffix(s, q) && len(s) > len(q) {
			base := strings.TrimSuffix(s, q)
			if base == "" {
				return ""
			}
			return base + "/" + q
		}
	}

	// неизвестный формат — возвращаем "" (т.к. без quote)
	return ""
}





package market

import "testing"

func TestNormalizeSymbol_Full(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"BTCUSDT", "BTC/USDT"},
		{"btcusdt", "BTC/USDT"},
		{"BTC-USDT", "BTC/USDT"},
		{"eth-btc", "ETH/BTC"},
		{"ETHBTC", "ETH/BTC"},
		{"BTC", ""},      // неполный символ
		{"BTC/", ""},     // неполный символ
		{"XYZABC", ""},   // неизвестный формат
		{"  btcusdt  ", "BTC/USDT"},
	}

	for _, tt := range tests {
		got := NormalizeSymbol_Full(tt.in)
		if got != tt.want {
			t.Errorf("NormalizeSymbol_Full(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestKey_Full(t *testing.T) {
	tests := []struct {
		exchange string
		symbol   string
		want     string
	}{
		{"MEXC", "BTCUSDT", "MEXC:BTC/USDT"},
		{"OKX", "BTC-USDT", "OKX:BTC/USDT"},
		{"KuCoin", "eth-btc", "KuCoin:ETH/BTC"},
		{"MEXC", "BTC", "MEXC:"},        // неполный символ
		{"KuCoin", "XYZABC", "KuCoin:"}, // неизвестный символ
	}

	for _, tt := range tests {
		got := Key_Full(tt.exchange, tt.symbol)
		if got != tt.want {
			t.Errorf("Key_Full(%q, %q) = %q, want %q", tt.exchange, tt.symbol, got, tt.want)
		}
	}
}



