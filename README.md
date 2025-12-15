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




type Config struct {
	Exchange      string
	TrianglesFile string
	BookInterval  string

	FeePerLeg float64 // доля, 0.0004 = 0.04%
	MinProfit float64 // доля
	MinStart  float64 // MIN_START_USDT / MIN_START

	StartFraction   float64 // 0..1
	Debug           bool
	TradeEnabled    bool    // из TRADE_ENABLED
	TradeAmountUSDT float64 // из TRADE_AMOUNT_USDT

	APIKey    string
	APISecret string
}



func Load() Config {
	_ = godotenv.Load(".env")

	// Биржа
	ex := strings.ToUpper(strings.TrimSpace(os.Getenv("EXCHANGE")))
	if ex == "" {
		ex = "MEXC"
	}

	// Файл треугольников
	tf := strings.TrimSpace(os.Getenv("TRIANGLES_FILE"))
	if tf == "" {
		tf = "triangles_markets.csv"
	}

	// Интервал книги
	bi := strings.TrimSpace(os.Getenv("BOOK_INTERVAL"))
	if bi == "" {
		bi = "100ms"
	}

	// Проценты (в ENV в процентах, в cfg в долях)
	feePct := loadEnvFloat("FEE_PCT", 0.04)       // 0.04% по умолчанию
	minPct := loadEnvFloat("MIN_PROFIT_PCT", 0.5) // 0.5% по умолчанию

	feePerLeg := feePct / 100.0
	minProfit := minPct / 100.0

	// MIN_START_USDT (предпочтительно) или MIN_START
	minStart := loadEnvFloat("MIN_START_USDT", -1)
	if minStart < 0 {
		minStart = loadEnvFloat("MIN_START", 0)
	}

	// START_FRACTION (0..1)
	startFraction := clamp01(loadEnvFloat("START_FRACTION", 0.5), 0.5)

	// DEBUG
	debugFlag := strings.ToLower(strings.TrimSpace(os.Getenv("DEBUG"))) == "true"
	debug = debugFlag // чтобы Dlog() знал про debug

	// TRADE_ENABLED
	tradeEnabled := strings.ToLower(strings.TrimSpace(os.Getenv("TRADE_ENABLED"))) == "true"

	// ФИКСИРОВАННЫЙ ОБЪЁМ ТОРГОВЛИ В USDT (0 = не ограничиваем, берём safeStart)
	tradeAmountUSDT := loadEnvFloat("TRADE_AMOUNT_USDT", 0)

	// API-ключи: сперва EXCHANGE_API_KEY/SECRET, потом API_KEY/SECRET
	apiKey := strings.TrimSpace(os.Getenv(ex + "_API_KEY"))
	if apiKey == "" {
		apiKey = strings.TrimSpace(os.Getenv("API_KEY"))
	}
	apiSecret := strings.TrimSpace(os.Getenv(ex + "_API_SECRET"))
	if apiSecret == "" {
		apiSecret = strings.TrimSpace(os.Getenv("API_SECRET"))
	}

	cfg := Config{
		Exchange:        ex,
		TrianglesFile:   tf,
		BookInterval:    bi,
		FeePerLeg:       feePerLeg,
		MinProfit:       minProfit,
		MinStart:        minStart,
		StartFraction:   startFraction,
		Debug:           debugFlag,
		TradeEnabled:    tradeEnabled,
		TradeAmountUSDT: tradeAmountUSDT,
		APIKey:          apiKey,
		APISecret:       apiSecret,
	}

	log.Printf("Exchange: %s", cfg.Exchange)
	log.Printf("Triangles file: %s", cfg.TrianglesFile)
	log.Printf("Book interval: %s", cfg.BookInterval)
	log.Printf("Fee per leg: %.4f %% (rate=%.6f)", feePct, cfg.FeePerLeg)
	log.Printf("Min profit per cycle: %.4f %% (rate=%.6f)", minPct, cfg.MinProfit)
	log.Printf("Min start amount: %.4f", cfg.MinStart)
	log.Printf("Start fraction: %.4f", cfg.StartFraction)
	log.Printf("Debug: %v", cfg.Debug)
	log.Printf("Trade enabled: %v", cfg.TradeEnabled)
	log.Printf("Trade amount (USDT): %.4f", cfg.TradeAmountUSDT)

	if cfg.APIKey == "" || cfg.APISecret == "" {
		log.Printf("API key/secret: NOT SET (реальная торговля невозможна)")
	} else {
		log.Printf("API key/secret: loaded for %s", cfg.Exchange)
	}

	return cfg
}



EXCHANGE=MEXC
TRIANGLES_FILE=triangles_markets.csv
BOOK_INTERVAL=10ms

FEE_PCT=0.04
MIN_PROFIT_PCT=0.1
MIN_START_USDT=2
START_FRACTION=0.5

DEBUG=false
TRADE_ENABLED=false       # пока false, чтобы не торговать
TRADE_AMOUNT_USDT=2       # твои 2 USDT на круг

MEXC_API_KEY=
MEXC_API_SECRET=




