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

func TestParsePair(t *testing.T) {
	tests := []struct {
		in        string
		wantBase  string
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



[{
	"resource": "/home/gaz358/myprog/crypt_proto/internal/market/market_test.go",
	"owner": "_generated_diagnostic_collection_name_#0",
	"code": {
		"value": "WrongArgCount",
		"target": {
			"$mid": 1,
			"path": "/golang.org/x/tools/internal/typesinternal",
			"scheme": "https",
			"authority": "pkg.go.dev",
			"fragment": "WrongArgCount"
		}
	},
	"severity": 8,
	"message": "not enough arguments in call to NormalizeSymbol\n\thave (string)\n\twant (string, string)",
	"source": "compiler",
	"startLineNumber": 19,
	"startColumn": 31,
	"endLineNumber": 19,
	"endColumn": 31,
	"origin": "extHost1"
}]

[{
	"resource": "/home/gaz358/myprog/crypt_proto/internal/market/market_test.go",
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
	"message": "undefined: Key",
	"source": "compiler",
	"startLineNumber": 61,
	"startColumn": 10,
	"endLineNumber": 61,
	"endColumn": 13,
	"origin": "extHost1"
}]









