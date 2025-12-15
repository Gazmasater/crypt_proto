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

	// Минимальный старт (обычно USDT). 0 = фильтр выключен.
	MinStart float64

	// Доля от maxStart, которую считаем безопасной (0..1). Например 0.5.
	StartFraction float64

	Debug     bool
	APIKey    string
	APISecret string
}

var debug bool

func SetDebug(v bool) { debug = v }

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

	tf := os.Getenv("TRIANGLES_FILE")
	if tf == "" {
		tf = "triangles_markets.csv"
	}

	bi := os.Getenv("BOOK_INTERVAL")
	if bi == "" {
		bi = "100ms"
	}

	feePct := loadEnvFloat("FEE_PCT", 0.04)
	minPct := loadEnvFloat("MIN_PROFIT_PCT", 0.5)

	// MIN_START_USDT (предпочтительно) или MIN_START
	minStart := loadEnvFloat("MIN_START_USDT", -1)
	if minStart < 0 {
		minStart = loadEnvFloat("MIN_START", 0)
	}

	startFraction := clamp01(loadEnvFloat("START_FRACTION", 0.5), 0.5)

	debug := strings.ToLower(os.Getenv("DEBUG")) == "true"

	cfg := Config{
		Exchange:      ex,
		TrianglesFile: tf,
		BookInterval:  bi,
		FeePerLeg:     feePct / 100.0,
		MinProfit:     minPct / 100.0,
		MinStart:      minStart,
		StartFraction: startFraction,
		Debug:         debug,
	}

	log.Printf("Exchange: %s", cfg.Exchange)
	log.Printf("Triangles file: %s", cfg.TrianglesFile)
	log.Printf("Book interval: %s", cfg.BookInterval)
	log.Printf("Fee per leg: %.4f %% (rate=%.6f)", feePct, cfg.FeePerLeg)
	log.Printf("Min profit per cycle: %.4f %% (rate=%.6f)", minPct, cfg.MinProfit)
	log.Printf("Min start amount: %.4f", cfg.MinStart)
	log.Printf("Start fraction: %.4f", cfg.StartFraction)

	return cfg
}

func Dlog(format string, args ...any) {
	if debug {
		log.Printf(format, args...)
	}
}
