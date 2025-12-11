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








crypt_proto/
  pb/                     // как было
  triangles_markets.csv   // как было
  arbitrage.log           // лог — не код

  cmd/
    cryptarb/
      main.go

  config/
    config.go

  domain/
    domain.go

  arb/
    arb.go

  mexc/
    ws.go
    proto_decoder.go


	[{
	"resource": "/home/gaz358/myprog/crypt_proto/cmd/cryptarb/main.go",
	"owner": "_generated_diagnostic_collection_name_#1",
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
	"message": "undefined: debug",
	"source": "compiler",
	"startLineNumber": 32,
	"startColumn": 2,
	"endLineNumber": 32,
	"endColumn": 7,
	"origin": "extHost1"
}]

[{
	"resource": "/home/gaz358/myprog/crypt_proto/mexc/ws.go",
	"owner": "_generated_diagnostic_collection_name_#1",
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
	"message": "undefined: debug",
	"source": "compiler",
	"startLineNumber": 171,
	"startColumn": 6,
	"endLineNumber": 171,
	"endColumn": 11,
	"origin": "extHost1"
}]




