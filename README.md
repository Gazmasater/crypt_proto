mx0vglmT3srN1IS19H
135bb7a7509e4421bad692415c53753b



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




[{
	"resource": "/home/gaz358/myprog/crypt_proto/arb/executor_real.go",
	"owner": "_generated_diagnostic_collection_name_#0",
	"code": {
		"value": "DuplicateDecl",
		"target": {
			"$mid": 1,
			"path": "/golang.org/x/tools/internal/typesinternal",
			"scheme": "https",
			"authority": "pkg.go.dev",
			"fragment": "DuplicateDecl"
		}
	},
	"severity": 8,
	"message": "SpotTrader redeclared in this block",
	"source": "compiler",
	"startLineNumber": 14,
	"startColumn": 6,
	"endLineNumber": 14,
	"endColumn": 16,
	"relatedInformation": [
		{
			"startLineNumber": 9,
			"startColumn": 6,
			"endLineNumber": 9,
			"endColumn": 16,
			"message": "other declaration of SpotTrader",
			"resource": "/home/gaz358/myprog/crypt_proto/arb/executor.go"
		}
	],
	"origin": "extHost1"
}]

[{
	"resource": "/home/gaz358/myprog/crypt_proto/arb/executor_real.go",
	"owner": "_generated_diagnostic_collection_name_#0",
	"code": {
		"value": "MissingFieldOrMethod",
		"target": {
			"$mid": 1,
			"path": "/golang.org/x/tools/internal/typesinternal",
			"scheme": "https",
			"authority": "pkg.go.dev",
			"fragment": "MissingFieldOrMethod"
		}
	},
	"severity": 8,
	"message": "e.trader.SmartMarketBuyUSDT undefined (type SpotTrader has no field or method SmartMarketBuyUSDT)",
	"source": "compiler",
	"startLineNumber": 84,
	"startColumn": 21,
	"endLineNumber": 84,
	"endColumn": 39,
	"origin": "extHost1"
}]

[{
	"resource": "/home/gaz358/myprog/crypt_proto/cmd/cryptarb/main.go",
	"owner": "_generated_diagnostic_collection_name_#0",
	"code": {
		"value": "InvalidIfaceAssign",
		"target": {
			"$mid": 1,
			"path": "/golang.org/x/tools/internal/typesinternal",
			"scheme": "https",
			"authority": "pkg.go.dev",
			"fragment": "InvalidIfaceAssign"
		}
	},
	"severity": 8,
	"message": "cannot use tr (variable of type *mexc.Trader) as arb.SpotTrader value in argument to arb.NewRealExecutor: *mexc.Trader does not implement arb.SpotTrader (missing method PlaceMarketOrder)",
	"source": "compiler",
	"startLineNumber": 83,
	"startColumn": 43,
	"endLineNumber": 83,
	"endColumn": 45,
	"origin": "extHost1"
}]

[{
	"resource": "/home/gaz358/myprog/crypt_proto/cmd/cryptarb/main.go",
	"owner": "_generated_diagnostic_collection_name_#0",
	"code": {
		"value": "IncompatibleAssign",
		"target": {
			"$mid": 1,
			"path": "/golang.org/x/tools/internal/typesinternal",
			"scheme": "https",
			"authority": "pkg.go.dev",
			"fragment": "IncompatibleAssign"
		}
	},
	"severity": 8,
	"message": "cannot use filters (variable of type map[string]arb.SymbolFilter) as float64 value in argument to arb.NewRealExecutor",
	"source": "compiler",
	"startLineNumber": 83,
	"startColumn": 55,
	"endLineNumber": 83,
	"endColumn": 62,
	"origin": "extHost1"
}]












