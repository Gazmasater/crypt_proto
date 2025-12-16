package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Exchange      string
	TrianglesFile string
	BookInterval  string

	FeePerLeg float64 // доля, 0.0004 = 0.04%
	MinProfit float64 // доля

	MinStart      float64 // USDT, 0 = выключено
	StartFraction float64 // 0..1

	// Trading
	TradeEnabled    bool
	TradeAmountUSDT float64
	TradeCooldownMs int

	// API keys (for selected exchange)
	APIKey    string
	APISecret string

	Debug bool
}

var debug bool

func SetDebug(v bool) { debug = v }

func Dlog(format string, args ...any) {
	if debug {
		log.Printf(format, args...)
	}
}

func loadEnvFloat(name string, def float64) float64 {
	raw := strings.TrimSpace(os.Getenv(name))
	if raw == "" {
		return def
	}
	v, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		log.Printf("bad %s=%q: %v, using default %f", name, raw, err, def)
		return def
	}
	return v
}

func loadEnvInt(name string, def int) int {
	raw := strings.TrimSpace(os.Getenv(name))
	if raw == "" {
		return def
	}
	v, err := strconv.Atoi(raw)
	if err != nil {
		log.Printf("bad %s=%q: %v, using default %d", name, raw, err, def)
		return def
	}
	return v
}

func loadEnvBool(name string, def bool) bool {
	raw := strings.TrimSpace(os.Getenv(name))
	if raw == "" {
		return def
	}
	switch strings.ToLower(raw) {
	case "1", "true", "yes", "y", "on":
		return true
	case "0", "false", "no", "n", "off":
		return false
	default:
		log.Printf("bad %s=%q: using default %v", name, raw, def)
		return def
	}
}

func clamp01(v, def float64) float64 {
	if v <= 0 || v > 1 {
		return def
	}
	return v
}

func Load() Config {
	_ = godotenv.Load(".env")

	ex := strings.ToUpper(strings.TrimSpace(os.Getenv("EXCHANGE")))
	if ex == "" {
		ex = "MEXC"
	}

	tf := strings.TrimSpace(os.Getenv("TRIANGLES_FILE"))
	if tf == "" {
		tf = "triangles_markets.csv"
	}

	bi := strings.TrimSpace(os.Getenv("BOOK_INTERVAL"))
	if bi == "" {
		bi = "100ms"
	}

	// проценты -> доли
	feePct := loadEnvFloat("FEE_PCT", 0.04)
	minPct := loadEnvFloat("MIN_PROFIT_PCT", 0.1)

	// MIN_START_USDT (предпочтительно) или MIN_START
	minStart := loadEnvFloat("MIN_START_USDT", -1)
	if minStart < 0 {
		minStart = loadEnvFloat("MIN_START", 0)
	}

	startFraction := clamp01(loadEnvFloat("START_FRACTION", 0.5), 0.5)

	// trading flags
	tradeEnabled := loadEnvBool("TRADE_ENABLED", false)
	tradeAmount := loadEnvFloat("TRADE_AMOUNT_USDT", 2.0)
	tradeCooldown := loadEnvInt("TRADE_COOLDOWN_MS", 800)

	// debug
	dbg := loadEnvBool("DEBUG", false)
	SetDebug(dbg)

	// keys: exchange-specific first, then fallback
	apiKey, apiSecret := "", ""
	switch ex {
	case "MEXC":
		apiKey = strings.TrimSpace(os.Getenv("MEXC_API_KEY"))
		apiSecret = strings.TrimSpace(os.Getenv("MEXC_API_SECRET"))
	case "KUCOIN":
		// если позже добавишь KuCoin trader — будет удобно
		apiKey = strings.TrimSpace(os.Getenv("KUCOIN_API_KEY"))
		apiSecret = strings.TrimSpace(os.Getenv("KUCOIN_API_SECRET"))
	case "OKX":
		apiKey = strings.TrimSpace(os.Getenv("OKX_API_KEY"))
		apiSecret = strings.TrimSpace(os.Getenv("OKX_API_SECRET"))
	}

	// fallback для совместимости
	if apiKey == "" {
		apiKey = strings.TrimSpace(os.Getenv("API_KEY"))
	}
	if apiSecret == "" {
		apiSecret = strings.TrimSpace(os.Getenv("API_SECRET"))
	}

	cfg := Config{
		Exchange:        ex,
		TrianglesFile:   tf,
		BookInterval:    bi,
		FeePerLeg:       feePct / 100.0,
		MinProfit:       minPct / 100.0,
		MinStart:        minStart,
		StartFraction:   startFraction,
		TradeEnabled:    tradeEnabled,
		TradeAmountUSDT: tradeAmount,
		TradeCooldownMs: tradeCooldown,
		APIKey:          apiKey,
		APISecret:       apiSecret,
		Debug:           dbg,
	}

	log.Printf("Exchange: %s", cfg.Exchange)
	log.Printf("Triangles file: %s", cfg.TrianglesFile)
	log.Printf("Book interval: %s", cfg.BookInterval)
	log.Printf("Fee per leg: %.4f %% (rate=%.6f)", feePct, cfg.FeePerLeg)
	log.Printf("Min profit per cycle: %.4f %% (rate=%.6f)", minPct, cfg.MinProfit)
	log.Printf("Min start amount (USDT): %.4f", cfg.MinStart)
	log.Printf("Start fraction: %.4f", cfg.StartFraction)
	log.Printf("Trade enabled: %v", cfg.TradeEnabled)
	log.Printf("Trade amount (USDT): %.4f", cfg.TradeAmountUSDT)
	log.Printf("Trade cooldown (ms): %d", cfg.TradeCooldownMs)

	if cfg.TradeEnabled {
		if cfg.APIKey != "" && cfg.APISecret != "" {
			log.Printf("API key/secret: loaded for %s", cfg.Exchange)
		} else {
			log.Printf("API key/secret: MISSING for %s (will fall back to DRY-RUN if main.go checks keys)", cfg.Exchange)
		}
	}

	return cfg
}
