package main

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

/* =========================  CONFIG  ========================= */

type Config struct {
	TrianglesFile string
	BookInterval  string
	FeePerLeg     float64 // комиссия за одну ногу, доля: 0.0004 = 0.04%
	MinProfit     float64 // минимальная прибыль за круг, доля: 0.005 = 0.5%
	Debug         bool
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

func loadConfig() Config {
	_ = godotenv.Load(".env")

	tf := os.Getenv("TRIANGLES_FILE")
	if tf == "" {
		tf = "triangles_markets.csv"
	}

	bi := os.Getenv("BOOK_INTERVAL")
	if bi == "" {
		bi = "100ms"
	}

	feePct := loadEnvFloat("FEE_PCT", 0.04)       // проценты
	minPct := loadEnvFloat("MIN_PROFIT_PCT", 0.5) // проценты

	debugFlag := strings.EqualFold(os.Getenv("DEBUG"), "true")

	cfg := Config{
		TrianglesFile: tf,
		BookInterval:  bi,
		FeePerLeg:     feePct / 100.0,
		MinProfit:     minPct / 100.0,
		Debug:         debugFlag,
	}

	log.Printf("Triangles file: %s", cfg.TrianglesFile)
	log.Printf("Book interval: %s", cfg.BookInterval)
	log.Printf("Fee per leg: %.4f %% (rate=%.6f)", feePct, cfg.FeePerLeg)
	log.Printf("Min profit per cycle: %.4f %% (rate=%.6f)", minPct, cfg.MinProfit)

	return cfg
}

/* =========================  LOGGING  ========================= */

var debug bool

func dlog(format string, args ...any) {
	if debug {
		log.Printf(format, args...)
	}
}
