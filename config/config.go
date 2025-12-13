package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

/* =========================  CONFIG  ========================= */

type Config struct {
	Exchange      string // "MEXC" или "KUCOIN"
	TrianglesFile string
	BookInterval  string
	FeePerLeg     float64 // как доля, 0.001 = 0.1%
	MinProfit     float64 // как доля, 0.003 = 0.3%

	// Минимальный стартовый объём (в валюте начала треугольника).
	// Обычно это USDT; если 0 - фильтр отключен.
	MinStart float64

	Debug bool
}

var debug bool

func SetDebug(v bool) {
	debug = v
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

func Load() Config {
	_ = godotenv.Load(".env")

	ex := strings.ToUpper(strings.TrimSpace(os.Getenv("EXCHANGE")))
	if ex == "" {
		ex = "MEXC"
	}

	tf := os.Getenv("TRIANGLES_FILE")
	if tf == "" {
		tf = "triangles_markets.csv"
	}
	bi := os.Getenv("BOOK_INTERVAL")
	if bi == "" {
		bi = "100ms"
	}

	// проценты
	feePct := loadEnvFloat("FEE_PCT", 0.04)
	minPct := loadEnvFloat("MIN_PROFIT_PCT", 0.5)

	// Минимальный стартовый объём (обычно USDT). Можно задавать как MIN_START_USDT
	// (предпочтительно), либо MIN_START. Если не задано - фильтр отключен.
	minStart := loadEnvFloat("MIN_START_USDT", -1)
	if minStart < 0 {
		minStart = loadEnvFloat("MIN_START", 0)
	}

	debug := strings.ToLower(os.Getenv("DEBUG")) == "true"

	cfg := Config{
		Exchange:      ex,
		TrianglesFile: tf,
		BookInterval:  bi,
		FeePerLeg:     feePct / 100.0,
		MinProfit:     minPct / 100.0,
		MinStart:      minStart,
		Debug:         debug,
	}

	log.Printf("Exchange: %s", cfg.Exchange)
	log.Printf("Triangles file: %s", tf)
	log.Printf("Book interval: %s", bi)
	log.Printf("Fee per leg: %.4f %% (rate=%.6f)", feePct, cfg.FeePerLeg)
	log.Printf("Min profit per cycle: %.4f %% (rate=%.6f)", minPct, cfg.MinProfit)
	log.Printf("Min start amount: %.4f", cfg.MinStart)

	return cfg
}

/* =========================  LOGGING  ========================= */

func Dlog(format string, args ...any) {
	if debug {
		log.Printf(format, args...)
	}
}
