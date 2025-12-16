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



package arb

import (
	"context"
	"fmt"
	"io"

	"crypt_proto/domain"
)

type DryRunExecutor struct {
	out io.Writer
}

func NewDryRunExecutor(out io.Writer) *DryRunExecutor {
	return &DryRunExecutor{out: out}
}

func (e *DryRunExecutor) Name() string { return "DRY-RUN" }

func (e *DryRunExecutor) Execute(ctx context.Context, t domain.Triangle, quotes map[string]domain.Quote, startUSDT float64) error {
	fmt.Fprintf(e.out, "  [DRY RUN] start=%.6f USDT triangle=%s\n", startUSDT, t.Name)
	return nil
}



type Executor interface {
	Name() string
	Execute(ctx context.Context, t domain.Triangle, quotes map[string]domain.Quote, startUSDT float64) error
}


[ARB] +0.114%  USDT→USELESS→USDC→USDT  maxStart=9.9183 USDT (9.9183 USDT)  safeStart=9.9183 USDT (9.9183 USDT) (x1.00)  bottleneck=USELESSUSDT
  USELESSUSDT (USELESS/USDT): bid=0.0719700000 ask=0.0720810000  spread=0.0001110000 (0.15411%)  bidQty=514.1000 askQty=137.6000
  USELESSUSDC (USELESS/USDC): bid=0.0722500000 ask=0.0726600000  spread=0.0004100000 (0.56587%)  bidQty=343.2800 askQty=455.3900
  USDCUSDT (USDC/USDT): bid=1.0000000000 ask=1.0001000000  spread=0.0001000000 (0.01000%)  bidQty=164179.1900 askQty=262862.7200

  [REAL EXEC] start=2.000000 USDT triangle=USDT→USELESS→USDC→USDT
    [REAL EXEC] leg 1: BUY USELESSUSDT quoteOrderQty=2.000000
2025-12-16 06:33:06.550
[ARB] +0.216%  USDT→USELESS→USDC→USDT  maxStart=9.9083 USDT (9.9083 USDT)  safeStart=9.9083 USDT (9.9083 USDT) (x1.00)  bottleneck=USELESSUSDT
  USELESSUSDT (USELESS/USDT): bid=0.0718090000 ask=0.0720080000  spread=0.0001990000 (0.27674%)  bidQty=1054.8200 askQty=137.6000
  USELESSUSDC (USELESS/USDC): bid=0.0722500000 ask=0.0726600000  spread=0.0004100000 (0.56587%)  bidQty=343.2800 askQty=455.3900
  USDCUSDT (USDC/USDT): bid=1.0000000000 ask=1.0001000000  spread=0.0001000000 (0.01000%)  bidQty=164179.1900 askQty=262862.7200

2025-12-16 06:33:06.590
[ARB] +0.288%  USDT→USELESS→USDC→USDT  maxStart=24.7109 USDT (24.7109 USDT)  safeStart=24.7109 USDT (24.7109 USDT) (x1.00)  bottleneck=USELESSUSDC
  USELESSUSDT (USELESS/USDT): bid=0.0718090000 ask=0.0719560000  spread=0.0001470000 (0.20450%)  bidQty=1054.8200 askQty=514.1000
  USELESSUSDC (USELESS/USDC): bid=0.0722500000 ask=0.0726600000  spread=0.0004100000 (0.56587%)  bidQty=343.2800 askQty=455.3900
  USDCUSDT (USDC/USDT): bid=1.0000000000 ask=1.0001000000  spread=0.0001000000 (0.01000%)  bidQty=164179.1900 askQty=262862.7200

2025-12-16 06:33:06.600
[ARB] +0.288%  USDT→USELESS→USDC→USDT  maxStart=24.7109 USDT (24.7109 USDT)  safeStart=24.7109 USDT (24.7109 USDT) (x1.00)  bottleneck=USELESSUSDC
  USELESSUSDT (USELESS/USDT): bid=0.0718100000 ask=0.0719560000  spread=0.0001460000 (0.20311%)  bidQty=149.0000 askQty=514.1000
  USELESSUSDC (USELESS/USDC): bid=0.0722500000 ask=0.0726600000  spread=0.0004100000 (0.56587%)  bidQty=343.2800 askQty=455.3900
  USDCUSDT (USDC/USDT): bid=1.0000000000 ask=1.0001000000  spread=0.0001000000 (0.01000%)  bidQty=164179.1900 askQty=262862.7200

    [REAL EXEC] leg 1 ERROR: mexc order error: status=400 body={"code":700013,"msg":"Invalid content Type."}
2025-12-16 06:33:28.094
[ARB] +0.209%  USDT→NAKA→USDC→USDT  maxStart=5.7025 USDT (5.7025 USDT)  safeStart=5.7025 USDT (5.7025 USDT) (x1.00)  bottleneck=NAKAUSDC
  NAKAUSDT (NAKA/USDT): bid=0.0787500000 ask=0.0789400000  spread=0.0001900000 (0.24098%)  bidQty=1162.6600 askQty=206.4000
  NAKAUSDC (NAKA/USDC): bid=0.0792000000 ask=0.0793000000  spread=0.0001000000 (0.12618%)  bidQty=72.2100 askQty=36.6700
  USDCUSDT (USDC/USDT): bid=1.0000000000 ask=1.0001000000  spread=0.0001000000 (0.01000%)  bidQty=163678.6400 askQty=260653.7300

  [REAL EXEC] start=2.000000 USDT triangle=USDT→NAKA→USDC→USDT
    [REAL EXEC] leg 1: BUY NAKAUSDT quoteOrderQty=2.000000
    [REAL EXEC] leg 1 ERROR: mexc order error: status=400 body={"code":700013,"msg":"Invalid content Type."}
^C2025/12/16 06:33:38.709212 shutting down...
2025/12/16 06:33:38.713042 bye




