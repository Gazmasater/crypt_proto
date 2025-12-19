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




export TRADE_AMOUNT_USDT=100
export FEE_PCT=0.04
export SELL_SAFETY=0.995

export TRIANGLES_FILE=triangles_markets.csv
export TRIANGLES_ENRICHED_FILE=triangles_markets_enriched.csv

go run ./cmd/triangles_enrich_mexc



package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"crypt_proto/arb"
	"crypt_proto/config"
	"crypt_proto/domain"
	"crypt_proto/exchange"
	"crypt_proto/kucoin"
	"crypt_proto/mexc"

	_ "net/http/pprof"
)

func main() {
	// pprof
	go func() {
		log.Println("pprof on http://localhost:6060/debug/pprof/")
		_ = http.ListenAndServe("localhost:6060", nil)
	}()

	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	cfg := config.Load()

	// Context / signals
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// Triangles
	triangles, symbols, indexBySymbol, err := domain.LoadTriangles(cfg.TrianglesFile)
	if err != nil {
		log.Fatalf("load triangles: %v", err)
	}
	if len(triangles) == 0 {
		log.Fatal("нет треугольников, нечего мониторить")
	}
	if len(symbols) == 0 {
		log.Fatal("нет символов для подписки")
	}
	log.Printf("треугольников: %d", len(triangles))
	log.Printf("символов для подписки: %d", len(symbols))

	// Exchange feed
	var feed exchange.MarketDataFeed
	switch cfg.Exchange {
	case "MEXC":
		feed = mexc.NewFeed(cfg.Debug)
	case "KUCOIN":
		feed = kucoin.NewFeed(cfg.Debug)
	default:
		log.Fatalf("unknown EXCHANGE=%q (expected MEXC or KUCOIN)", cfg.Exchange)
	}
	log.Printf("Using exchange: %s", feed.Name())

	// Log output
	logFile, logBuf, arbOut := arb.OpenLogWriter("arbitrage.log")
	defer logFile.Close()
	defer logBuf.Flush()

	// Events channel
	events := make(chan domain.Event, 8192)

	var wg sync.WaitGroup

	// Consumer
	consumer := arb.NewConsumer(cfg.FeePerLeg, cfg.MinProfit, cfg.MinStart, arbOut)
	consumer.StartFraction = cfg.StartFraction

	// Trading toggles
	consumer.TradeEnabled = cfg.TradeEnabled
	consumer.TradeAmountUSDT = cfg.TradeAmountUSDT
	consumer.TradeCooldown = time.Duration(cfg.TradeCooldownMs) * time.Millisecond

	log.Printf(
		"TRADE: enabled=%v amountUSDT=%.6f cooldown=%s feePerLeg=%.6f minProfit=%.6f minStart=%.6f startFraction=%.4f exchange=%s debug=%v",
		consumer.TradeEnabled,
		consumer.TradeAmountUSDT,
		consumer.TradeCooldown,
		cfg.FeePerLeg,
		cfg.MinProfit,
		cfg.MinStart,
		consumer.StartFraction,
		cfg.Exchange,
		cfg.Debug,
	)

	// Executor
	if cfg.Exchange == "MEXC" && cfg.TradeEnabled && cfg.APIKey != "" && cfg.APISecret != "" {
		tr := mexc.NewTrader(cfg.APIKey, cfg.APISecret, cfg.Debug)

		// startUSDT — фиксированная сумма на сделку
		startUSDT := cfg.TradeAmountUSDT
		if startUSDT <= 0 {
			startUSDT = 35.0
		}

		re := arb.NewRealExecutor(tr, arbOut, startUSDT)
		re.StopAfterOne = true
		re.SetStopFunc(cancel)

		consumer.Executor = re
		log.Printf("Executor: REAL (startUSDT=%.6f) STOP_AFTER_ONE=true", startUSDT)
	} else {
		consumer.Executor = arb.NewNoopExecutor()
		log.Printf("Executor: NOOP (trade disabled, non-MEXC exchange, or missing keys)")
	}

	// Start consumer
	consumer.Start(ctx, events, triangles, indexBySymbol, &wg)

	// Start feed
	feed.Start(ctx, &wg, symbols, cfg.BookInterval, events)

	// Wait stop
	<-ctx.Done()
	log.Println("shutting down...")

	// ВАЖНО: events не закрываем — WS-горутин(ы) могут ещё писать и словить panic
	wg.Wait()
	log.Println("bye")
}



2025-12-19 22:07:14.609
[ARB] -0.334%  USDT→DOGE→USDE→USDT  maxStart=401.2095 USDT (401.2095 USDT)  safeStart=120.3628 USDT (120.3628 USDT) (x0.30)  bottleneck=DOGEUSDE
  DOGEUSDT (DOGE/USDT): bid=0.1313600000 ask=0.1313900000  spread=0.0000300000 (0.02284%)  bidQty=1909.7500 askQty=184010.4000
  DOGEUSDE (DOGE/USDE): bid=0.1312400000 ask=0.1316800000  spread=0.0004400000 (0.33470%)  bidQty=3052.0500 askQty=23.2200
  USDEUSDT (USDE/USDT): bid=0.9993000000 ask=0.9994000000  spread=0.0001000000 (0.01001%)  bidQty=1790026.8100 askQty=1843301.3600

2025-12-19 22:07:14.627
[ARB] -0.299%  USDT→XRP→BTC→USDT  maxStart=479.8713 USDT (479.8713 USDT)  safeStart=143.9614 USDT (143.9614 USDT) (x0.30)  bottleneck=XRPUSDT
  XRPUSDT (XRP/USDT): bid=1.8979000000 ask=1.8980000000  spread=0.0001000000 (0.00527%)  bidQty=2009.0200 askQty=252.8300
  XRPBTC (XRP/BTC): bid=0.0000217070 ask=0.0000217440  spread=0.0000000370 (0.17031%)  bidQty=297.1400 askQty=9.8500
  BTCUSDT (BTC/USDT): bid=87306.9100000000 ask=87308.0300000000  spread=1.1200000000 (0.00128%)  bidQty=0.3648 askQty=0.0050

2025-12-19 22:07:14.627
[ARB] -0.307%  USDT→XRP→ETH→USDT  maxStart=270.3155 USDT (270.3155 USDT)  safeStart=81.0946 USDT (81.0946 USDT) (x0.30)  bottleneck=XRPETH
  XRPUSDT (XRP/USDT): bid=1.8979000000 ask=1.8980000000  spread=0.0001000000 (0.00527%)  bidQty=2009.0200 askQty=252.8300
  XRPETH (XRP/ETH): bid=0.0006381000 ask=0.0006402000  spread=0.0000021000 (0.32856%)  bidQty=142.3500 askQty=0.6500
  ETHUSDT (ETH/USDT): bid=2969.7800000000 ask=2969.7900000000  spread=0.0100000000 (0.00034%)  bidQty=0.9745 askQty=0.1700

2025-12-19 22:07:14.627
[ARB] -0.229%  USDT→XRP→USD1→USDT  maxStart=479.8713 USDT (479.8713 USDT)  safeStart=143.9614 USDT (143.9614 USDT) (x0.30)  bottleneck=XRPUSDT
  XRPUSDT (XRP/USDT): bid=1.8979000000 ask=1.8980000000  spread=0.0001000000 (0.00527%)  bidQty=2009.0200 askQty=252.8300
  XRPUSD1 (XRP/USD1): bid=1.8982000000 ask=1.9006000000  spread=0.0024000000 (0.12636%)  bidQty=264.5800 askQty=260.7700
  USD1USDT (USD1/USDT): bid=0.9991000000 ask=0.9992000000  spread=0.0001000000 (0.01001%)  bidQty=4122.8900 askQty=252251.2100

2025-12-19 22:07:14.633
[ARB] -0.310%  USDT→ETH→USDE→USDT  maxStart=238.9797 USDT (238.9797 USDT)  safeStart=71.6939 USDT (71.6939 USDT) (x0.30)  bottleneck=ETHUSDE
  ETHUSDT (ETH/USDT): bid=2969.7800000000 ask=2969.7900000000  spread=0.0100000000 (0.00034%)  bidQty=0.9745 askQty=0.1700
  ETHUSDE (ETH/USDE): bid=2967.1000000000 ask=2976.3300000000  spread=9.2300000000 (0.31060%)  bidQty=0.0804 askQty=0.0170
  USDEUSDT (USDE/USDT): bid=0.9993000000 ask=0.9994000000  spread=0.0001000000 (0.01001%)  bidQty=1790026.8100 askQty=1843301.3600

2025-12-19 22:07:14.633
[ARB] -0.187%  USDT→ETH→BTC→USDT  maxStart=279.3009 USDT (279.3009 USDT)  safeStart=83.7903 USDT (83.7903 USDT) (x0.30)  bottleneck=ETHBTC
  ETHUSDT (ETH/USDT): bid=2969.7800000000 ask=2969.8000000000  spread=0.0200000000 (0.00067%)  bidQty=0.9745 askQty=1.1898
  ETHBTC (ETH/BTC): bid=0.0340030000 ask=0.0340170000  spread=0.0000140000 (0.04116%)  bidQty=0.0940 askQty=0.0140
  BTCUSDT (BTC/USDT): bid=87306.9100000000 ask=87308.0300000000  spread=1.1200000000 (0.00128%)  bidQty=0.3648 askQty=0.0050

2025-12-19 22:07:14.633
[ARB] -0.309%  USDT→ETH→EUR→USDT  maxStart=716.9712 USDT (716.9712 USDT)  safeStart=215.0914 USDT (215.0914 USDT) (x0.30)  bottleneck=ETHEUR
  ETHUSDT (ETH/USDT): bid=2969.7800000000 ask=2969.8000000000  spread=0.0200000000 (0.00067%)  bidQty=0.9745 askQty=1.1898
  ETHEUR (ETH/EUR): bid=2528.6300000000 ask=2536.4600000000  spread=7.8300000000 (0.30918%)  bidQty=0.2413 askQty=0.0510
  EURUSDT (EUR/USDT): bid=1.1726000000 ask=1.1727000000  spread=0.0001000000 (0.00853%)  bidQty=44509.6800 askQty=17901.9000

2025-12-19 22:07:14.633
[ARB] -0.198%  USDT→ETH→USD1→USDT  maxStart=288.5713 USDT (288.5713 USDT)  safeStart=86.5714 USDT (86.5714 USDT) (x0.30)  bottleneck=ETHUSD1
  ETHUSDT (ETH/USDT): bid=2969.7800000000 ask=2969.8000000000  spread=0.0200000000 (0.00067%)  bidQty=0.9745 askQty=1.1898
  ETHUSD1 (ETH/USD1): bid=2971.0500000000 ask=2974.3100000000  spread=3.2600000000 (0.10967%)  bidQty=0.0971 askQty=0.1731
  USD1USDT (USD1/USDT): bid=0.9991000000 ask=0.9992000000  spread=0.0001000000 (0.01001%)  bidQty=4122.8900 askQty=252251.2100

2025-12-19 22:07:14.633
[ARB] -0.498%  USDT→OXT→ETH→USDT  maxStart=223.2900 USDT (223.2900 USDT)  safeStart=66.9870 USDT (66.9870 USDT) (x0.30)  bottleneck=OXTUSDT
  OXTUSDT (OXT/USDT): bid=0.0247600000 ask=0.0248100000  spread=0.0000500000 (0.20173%)  bidQty=21191.3800 askQty=9000.0000
  OXTETH (OXT/ETH): bid=0.0000083250 ask=0.0000083680  spread=0.0000000430 (0.51519%)  bidQty=20767.5500 askQty=8820.0000
  ETHUSDT (ETH/USDT): bid=2969.7800000000 ask=2969.8000000000  spread=0.0200000000 (0.00067%)  bidQty=0.9745 askQty=1.1898

2025-12-19 22:07:14.633
[ARB] -0.424%  USDT→RSR→ETH→USDT  maxStart=146.1932 USDT (146.1932 USDT)  safeStart=43.8580 USDT (43.8580 USDT) (x0.30)  bottleneck=RSRUSDT
  RSRUSDT (RSR/USDT): bid=0.0025360000 ask=0.0025390000  spread=0.0000030000 (0.11823%)  bidQty=123677.3600 askQty=57579.0400
  RSRETH (RSR/ETH): bid=0.0000008526 ask=0.0000008552  spread=0.0000000026 (0.30449%)  bidQty=111387.8200 askQty=8900.1500
  ETHUSDT (ETH/USDT): bid=2969.7800000000 ask=2969.8000000000  spread=0.0200000000 (0.00067%)  bidQty=0.9745 askQty=1.1898

2025-12-19 22:07:14.633
[ARB] -0.307%  USDT→XRP→ETH→USDT  maxStart=270.3155 USDT (270.3155 USDT)  safeStart=81.0946 USDT (81.0946 USDT) (x0.30)  bottleneck=XRPETH
  XRPUSDT (XRP/USDT): bid=1.8979000000 ask=1.8980000000  spread=0.0001000000 (0.00527%)  bidQty=2009.0200 askQty=252.8300
  XRPETH (XRP/ETH): bid=0.0006381000 ask=0.0006402000  spread=0.0000021000 (0.32856%)  bidQty=142.3500 askQty=0.6500
  ETHUSDT (ETH/USDT): bid=2969.7800000000 ask=2969.8000000000  spread=0.0200000000 (0.00067%)  bidQty=0.9745 askQty=1.1898

  [REAL EXEC] SKIP: queue full (triangle=USDT→OXT→ETH→USDT)
  [REAL EXEC] SKIP: queue full (triangle=USDT→RSR→ETH→USDT)
2025-12-19 22:07:14.647
[ARB] -0.317%  USDT→SOL→BTC→USDT  maxStart=291.8619 USDT (291.8619 USDT)  safeStart=87.5586 USDT (87.5586 USDT) (x0.30)  bottleneck=SOLBTC
  SOLUSDT (SOL/USDT): bid=125.1800000000 ask=125.2000000000  spread=0.0200000000 (0.01598%)  bidQty=2.5890 askQty=819.7300
  SOLBTC (SOL/BTC): bid=0.0014316200 ask=0.0014415900  spread=0.0000099700 (0.69400%)  bidQty=2.3300 askQty=0.0500
  BTCUSDT (BTC/USDT): bid=87306.9100000000 ask=87308.0300000000  spread=1.1200000000 (0.00128%)  bidQty=0.3648 askQty=0.0050

2025-12-19 22:07:14.649
[ARB] -0.444%  USDT→0G→USDC→USDT  maxStart=490.2289 USDT (490.2289 USDT)  safeStart=147.0687 USDT (147.0687 USDT) (x0.30)  bottleneck=0GUSDC
  0GUSDT (0G/USDT): bid=0.7630000000 ask=0.7640000000  spread=0.0010000000 (0.13098%)  bidQty=780.9700 askQty=1614.7100
  0GUSDC (0G/USDC): bid=0.7616000000 ask=0.7650000000  spread=0.0034000000 (0.44543%)  bidQty=641.3400 askQty=1789.9200
  USDCUSDT (USDC/USDT): bid=1.0002000000 ask=1.0003000000  spread=0.0001000000 (0.01000%)  bidQty=29726.2000 askQty=175671.1300

2025-12-19 22:07:14.649
[ARB] -0.580%  USDT→1INCH→USDC→USDT  maxStart=637.3045 USDT (637.3045 USDT)  safeStart=191.1913 USDT (191.1913 USDT) (x0.30)  bottleneck=1INCHUSDC
  1INCHUSDT (1INCH/USDT): bid=0.1548000000 ask=0.1552000000  spread=0.0004000000 (0.25806%)  bidQty=5130.3600 askQty=15743.4400
  1INCHUSDC (1INCH/USDC): bid=0.1545000000 ask=0.1555000000  spread=0.0010000000 (0.64516%)  bidQty=4104.2900 askQty=12594.7500
  USDCUSDT (USDC/USDT): bid=1.0002000000 ask=1.0003000000  spread=0.0001000000 (0.01000%)  bidQty=29726.2000 askQty=175671.1300

2025-12-19 22:07:14.649
[ARB] -0.630%  USDT→ACH→USDC→USDT  maxStart=241.0288 USDT (241.0288 USDT)  safeStart=72.3086 USDT (72.3086 USDT) (x0.30)  bottleneck=ACHUSDC
  ACHUSDT (ACH/USDT): bid=0.0077580000 ask=0.0077830000  spread=0.0000250000 (0.32173%)  bidQty=64074.6400 askQty=507420.1500
  ACHUSDC (ACH/USDC): bid=0.0077440000 ask=0.0077940000  spread=0.0000500000 (0.64358%)  bidQty=30953.1400 askQty=159.5700
  USDCUSDT (USDC/USDT): bid=1.0002000000 ask=1.0003000000  spread=0.0001000000 (0.01000%)  bidQty=29726.2000 askQty=175671.1300

2025-12-19 22:07:14.649
[ARB] -0.572%  USDT→AEVO→USDC→USDT  maxStart=158.0108 USDT (158.0108 USDT)  safeStart=47.4032 USDT (47.4032 USDT) (x0.30)  bottleneck=AEVOUSDC
  AEVOUSDT (AEVO/USDT): bid=0.0360300000 ask=0.0361200000  spread=0.0000900000 (0.24948%)  bidQty=5465.5300 askQty=38021.5800
  AEVOUSDC (AEVO/USDC): bid=0.0359600000 ask=0.0361700000  spread=0.0002100000 (0.58228%)  bidQty=4372.4200 askQty=15924.7000
  USDCUSDT (USDC/USDT): bid=1.0002000000 ask=1.0003000000  spread=0.0001000000 (0.01000%)  bidQty=29726.2000 askQty=175671.1300

2025-12-19 22:07:14.649
[ARB] -0.505%  USDT→ALLO→USDC→USDT  maxStart=165.7683 USDT (165.7683 USDT)  safeStart=49.7305 USDT (49.7305 USDT) (x0.30)  bottleneck=ALLOUSDT
  ALLOUSDT (ALLO/USDT): bid=0.1063000000 ask=0.1065000000  spread=0.0002000000 (0.18797%)  bidQty=19496.0000 askQty=1556.5100
  ALLOUSDC (ALLO/USDC): bid=0.1061000000 ask=0.1067000000  spread=0.0006000000 (0.56391%)  bidQty=3436.5800 askQty=644.2700
  USDCUSDT (USDC/USDT): bid=1.0002000000 ask=1.0003000000  spread=0.0001000000 (0.01000%)  bidQty=29726.2000 askQty=175671.1300

2025-12-19 22:07:14.649
[ARB] -0.445%  USDT→ARB→USDC→USDT  maxStart=941.2188 USDT (941.2188 USDT)  safeStart=282.3656 USDT (282.3656 USDT) (x0.30)  bottleneck=ARBUSDC
  ARBUSDT (ARB/USDT): bid=0.1900000000 ask=0.1902000000  spread=0.0002000000 (0.10521%)  bidQty=25566.5700 askQty=39708.1800
  ARBUSDC (ARB/USDC): bid=0.1896000000 ask=0.1905000000  spread=0.0009000000 (0.47356%)  bidQty=4946.1000 askQty=6.3000
  USDCUSDT (USDC/USDT): bid=1.0002000000 ask=1.0003000000  spread=0.0001000000 (0.01000%)  bidQty=29726.2000 askQty=175671.1300

2025-12-19 22:07:14.649
[ARB] -0.543%  USDT→BANANA→USDC→USDT  maxStart=129.3336 USDT (129.3336 USDT)  safeStart=38.8001 USDT (38.8001 USDT) (x0.30)  bottleneck=BANANAUSDT
  BANANAUSDT (BANANA/USDT): bid=6.5170000000 ask=6.5320000000  spread=0.0150000000 (0.22990%)  bidQty=45.5600 askQty=19.8000
  BANANAUSDC (BANANA/USDC): bid=6.5050000000 ask=6.5410000000  spread=0.0360000000 (0.55189%)  bidQty=36.7700 askQty=15.8400
  USDCUSDT (USDC/USDT): bid=1.0002000000 ask=1.0003000000  spread=0.0001000000 (0.01000%)  bidQty=29726.2000 askQty=175671.1300

2025-12-19 22:07:14.649
[ARB] -0.667%  USDT→BERA→USDC→USDT  maxStart=1538.7955 USDT (1538.7955 USDT)  safeStart=461.6387 USDT (461.6387 USDT) (x0.30)  bottleneck=BERAUSDT
  BERAUSDT (BERA/USDT): bid=0.5740000000 ask=0.5760000000  spread=0.0020000000 (0.34783%)  bidQty=5125.3100 askQty=2671.5200
  BERAUSDC (BERA/USDC): bid=0.5729000000 ask=0.5738000000  spread=0.0009000000 (0.15697%)  bidQty=4100.2500 askQty=147.9200
  USDCUSDT (USDC/USDT): bid=1.0002000000 ask=1.0003000000  spread=0.0001000000 (0.01000%)  bidQty=29726.2000 askQty=175671.1300

2025-12-19 22:07:14.649
[ARB] -0.585%  USDT→BOME→USDC→USDT  maxStart=143.4160 USDT (143.4160 USDT)  safeStart=43.0248 USDT (43.0248 USDT) (x0.30)  bottleneck=BOMEUSDT
  BOMEUSDT (BOME/USDT): bid=0.0005916000 ask=0.0005927000  spread=0.0000011000 (0.18576%)  bidQty=602354.4000 askQty=241970.6900
  BOMEUSDC (BOME/USDC): bid=0.0005900000 ask=0.0005929000  spread=0.0000029000 (0.49032%)  bidQty=266607.6500 askQty=27992.9500
  USDCUSDT (USDC/USDT): bid=1.0002000000 ask=1.0003000000  spread=0.0001000000 (0.01000%)  bidQty=29726.2000 askQty=175671.1300

2025-12-19 22:07:14.649
[ARB] -0.188%  USDT→BTC→USDC→USDT  maxStart=436.5401 USDT (436.5401 USDT)  safeStart=130.9620 USDT (130.9620 USDT) (x0.30)  bottleneck=BTCUSDT
  BTCUSDT (BTC/USDT): bid=87306.9100000000 ask=87308.0300000000  spread=1.1200000000 (0.00128%)  bidQty=0.3648 askQty=0.0050
  BTCUSDC (BTC/USDC): bid=87256.8700000000 ask=87277.0300000000  spread=20.1600000000 (0.02310%)  bidQty=0.0124 askQty=0.0142
  USDCUSDT (USDC/USDT): bid=1.0002000000 ask=1.0003000000  spread=0.0001000000 (0.01000%)  bidQty=29726.2000 askQty=175671.1300

2025-12-19 22:07:14.649
[ARB] -0.470%  USDT→C→USDC→USDT  maxStart=188.1628 USDT (188.1628 USDT)  safeStart=56.4488 USDT (56.4488 USDT) (x0.30)  bottleneck=CUSDT
  CUSDT (C/USDT): bid=0.0880900000 ask=0.0882300000  spread=0.0001400000 (0.15880%)  bidQty=3787.6800 askQty=2132.6400
  CUSDC (C/USDC): bid=0.0879300000 ask=0.0883500000  spread=0.0004200000 (0.47651%)  bidQty=3030.1400 askQty=1706.1100
  USDCUSDT (USDC/USDT): bid=1.0002000000 ask=1.0003000000  spread=0.0001000000 (0.01000%)  bidQty=29726.2000 askQty=175671.1300

2025-12-19 22:07:14.649
[ARB] -0.671%  USDT→COOKIE→USDC→USDT  maxStart=170.4375 USDT (170.4375 USDT)  safeStart=51.1312 USDT (51.1312 USDT) (x0.30)  bottleneck=COOKIEUSDC
  COOKIEUSDT (COOKIE/USDT): bid=0.0404900000 ask=0.0406300000  spread=0.0001400000 (0.34517%)  bidQty=5240.9600 askQty=86138.5800
  COOKIEUSDC (COOKIE/USDC): bid=0.0404100000 ask=0.0406900000  spread=0.0002800000 (0.69051%)  bidQty=4192.7700 askQty=27328.5500
  USDCUSDT (USDC/USDT): bid=1.0002000000 ask=1.0003000000  spread=0.0001000000 (0.01000%)  bidQty=29726.2000 askQty=175671.1300

2025-12-19 22:07:14.649
[ARB] -0.406%  USDT→EIGEN→USDC→USDT  maxStart=120.1980 USDT (120.1980 USDT)  safeStart=36.0594 USDT (36.0594 USDT) (x0.30)  bottleneck=EIGENUSDC
  EIGENUSDT (EIGEN/USDT): bid=0.3976000000 ask=0.3979000000  spread=0.0003000000 (0.07542%)  bidQty=377.4100 askQty=377.4100
  EIGENUSDC (EIGEN/USDC): bid=0.3968000000 ask=0.3986000000  spread=0.0018000000 (0.45260%)  bidQty=301.9300 askQty=301.9300
  USDCUSDT (USDC/USDT): bid=1.0002000000 ask=1.0003000000  spread=0.0001000000 (0.01000%)  bidQty=29726.2000 askQty=175671.1300




