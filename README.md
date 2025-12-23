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



internal/market/
  normalize.go
  pair.go
  key.go


package market

import "strings"

// NormalizeSymbol приводит символ к виду BASE/QUOTE
// BTCUSDT   -> BTC/USDT
// BTC-USDT  -> BTC/USDT
func NormalizeSymbol(symbol string) string {
	s := strings.ToUpper(strings.TrimSpace(symbol))

	if strings.Contains(s, "-") {
		parts := strings.Split(s, "-")
		if len(parts) == 2 {
			return parts[0] + "/" + parts[1]
		}
	}

	if strings.HasSuffix(s, "USDT") {
		return strings.TrimSuffix(s, "USDT") + "/USDT"
	}

	return s
}



package market

import "strings"

type Pair struct {
	Base  string
	Quote string
}

func ParsePair(normalized string) Pair {
	parts := strings.Split(normalized, "/")
	if len(parts) != 2 {
		return Pair{}
	}

	return Pair{
		Base:  parts[0],
		Quote: parts[1],
	}
}




package market

// Key формирует единый ключ хранения
// MEXC:BTC/USDT
func Key(exchange, symbol string) string {
	return exchange + ":" + NormalizeSymbol(symbol)
}





📁 Что тестируем

Пакет:

internal/market/
  normalize.go
  pair.go
  key.go


Создаём рядом:

internal/market/market_test.go

✅ Тест NormalizeSymbol
package market

import "testing"

func TestNormalizeSymbol(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"BTCUSDT", "BTC/USDT"},
		{"btcusdt", "BTC/USDT"},
		{"BTC-USDT", "BTC/USDT"},
		{"eth-btc", "ETH/BTC"},
		{"ETHBTC", "ETHBTC"}, // неизвестный формат — не ломаем
		{"  btcusdt  ", "BTC/USDT"},
	}

	for _, tt := range tests {
		got := NormalizeSymbol(tt.in)
		if got != tt.want {
			t.Errorf("NormalizeSymbol(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

✅ Тест ParsePair
func TestParsePair(t *testing.T) {
	tests := []struct {
		in        string
		wantBase string
		wantQuote string
	}{
		{"BTC/USDT", "BTC", "USDT"},
		{"ETH/BTC", "ETH", "BTC"},
		{"INVALID", "", ""},
		{"BTC/", "", ""},
	}

	for _, tt := range tests {
		p := ParsePair(tt.in)
		if p.Base != tt.wantBase || p.Quote != tt.wantQuote {
			t.Errorf(
				"ParsePair(%q) = %+v, want Base=%q Quote=%q",
				tt.in, p, tt.wantBase, tt.wantQuote,
			)
		}
	}
}

✅ Тест Key
func TestKey(t *testing.T) {
	tests := []struct {
		exchange string
		symbol   string
		want     string
	}{
		{"MEXC", "BTCUSDT", "MEXC:BTC/USDT"},
		{"OKX", "BTC-USDT", "OKX:BTC/USDT"},
		{"KuCoin", "eth-btc", "KuCoin:ETH/BTC"},
	}

	for _, tt := range tests {
		got := Key(tt.exchange, tt.symbol)
		if got != tt.want {
			t.Errorf("Key(%q, %q) = %q, want %q",
				tt.exchange, tt.symbol, got, tt.want)
		}
	}
}

▶️ Как запускать
go test ./internal/market


или всё сразу:

go test ./...


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


[{
	"resource": "/home/gaz358/myprog/crypt_proto/internal/market/pair.go",
	"owner": "_generated_diagnostic_collection_name_#0",
	"code": {
		"value": "UndeclaredName",
		"target": {
			"$mid": 1,
			"path": "/golang.org/x/tools/internal/typesinternal",
			"scheme": "https",
			"authority": "pkg.go.dev",
			"fragment": "UndeclaredName"
		}
	},
	"severity": 8,
	"message": "undefined: Pair",
	"source": "compiler",
	"startLineNumber": 5,
	"startColumn": 35,
	"endLineNumber": 5,
	"endColumn": 39,
	"origin": "extHost1"
}]

