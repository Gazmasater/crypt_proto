package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"sync"
	"syscall"
	"time"

	_ "net/http/pprof"
)

/* =========================  MAIN / APP  ========================= */

func main() {
	// pprof
	go func() {
		log.Println("pprof on http://localhost:6060/debug/pprof/")
		if err := http.ListenAndServe("localhost:6060", nil); err != nil {
			log.Printf("pprof server error: %v", err)
		}
	}()

	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	cfg := loadConfig()
	debug = cfg.Debug

	arbOut, closeArb := initArbLogger("arbitrage.log")
	defer closeArb()

	triangles, symbols, indexBySymbol, err := loadTriangles(cfg.TrianglesFile)
	if err != nil {
		log.Fatalf("load triangles: %v", err)
	}
	if len(triangles) == 0 {
		log.Fatal("нет треугольников, нечего мониторить")
	}
	if len(symbols) == 0 {
		log.Fatal("нет символов для подписки")
	}
	log.Printf("символов для подписки всего: %d", len(symbols))

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	events := make(chan Event, 8192)

	var wg sync.WaitGroup
	startWSWorkers(ctx, &wg, symbols, cfg.BookInterval, events)

	go consumeEvents(ctx, events, triangles, indexBySymbol, cfg.FeePerLeg, cfg.MinProfit, arbOut)

	<-ctx.Done()
	log.Println("shutting down...")

	time.Sleep(200 * time.Millisecond)
	close(events)
	wg.Wait()
	log.Println("bye")
}
