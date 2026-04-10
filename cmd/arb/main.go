package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"crypt_proto/internal/calculator"
	"crypt_proto/internal/collector"
	"crypt_proto/internal/queue"
	"crypt_proto/pkg/models"
)

func main() {
	logFile, cleanup, err := setupLogging()
	if err != nil {
		log.Fatalf("setup logging: %v", err)
	}
	defer cleanup()
	log.Printf("[Main] log file: %s", logFile)

	go func() {
		log.Println("pprof on http://localhost:6060/debug/pprof/")
		if err := http.ListenAndServe("localhost:6060", nil); err != nil {
			log.Printf("pprof server error: %v", err)
		}
	}()

	out := make(chan *models.MarketData, 100_000)
	mem := queue.NewMemoryStore()

	kc, _, err := collector.NewKuCoinCollectorFromCSV("../exchange/data/kucoin_triangles_usdt.csv")
	if err != nil {
		log.Fatal(err)
	}
	if err := kc.Start(out); err != nil {
		log.Fatal(err)
	}
	log.Println("[Main] KuCoinCollector started")

	triangles, err := calculator.ParseTrianglesFromCSV("../exchange/data/kucoin_triangles_usdt.csv")
	if err != nil {
		log.Fatal(err)
	}

	cfg := calculator.DefaultConfig()
	cfg.LogMode = calculator.LogDebug
	cfg.MinVolumeUSDT = 10
	cfg.MinProfitPct = 0
	cfg.QuoteAgeMaxMS = 400
	cfg.StatsEverySec = 5

	calc := calculator.NewCalculator(mem, triangles, cfg)
	go calc.Run(out)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("[Main] shutting down...")
	kc.Stop()
	close(out)
	log.Println("[Main] exited")
}

func setupLogging() (string, func(), error) {
	logDir := "logs"
	if err := os.MkdirAll(logDir, 0o755); err != nil {
		return "", func() {}, err
	}
	filename := filepath.Join(logDir, fmt.Sprintf("arb_%s.log", time.Now().Format("20060110_150405")))
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return "", func() {}, err
	}
	mw := io.MultiWriter(os.Stdout, f)
	log.SetOutput(mw)
	log.SetFlags(log.LstdFlags)
	cleanup := func() {
		_ = f.Close()
	}
	return filename, cleanup, nil
}
