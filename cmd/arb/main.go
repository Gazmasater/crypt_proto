package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"crypt_proto/internal/calculator"
	"crypt_proto/internal/collector"
	"crypt_proto/internal/executor"
	"crypt_proto/internal/queue"
	"crypt_proto/pkg/models"
)

func main() {
	go func() {
		log.Println("pprof on http://localhost:6060/debug/pprof/")
		if err := http.ListenAndServe("localhost:6060", nil); err != nil {
			log.Printf("pprof server stopped: %v", err)
		}
	}()

	out := make(chan *models.MarketData, 100_000)
	oppCh := make(chan *executor.Opportunity, 4096)

	mem := queue.NewMemoryStore()
	go mem.Run()

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

	calc := calculator.NewCalculator(mem, kc, triangles, oppCh)
	go calc.Run(out)

	exec := executor.NewExecutor(executor.DefaultConfig())
	go func() {
		for op := range oppCh {
			exec.Handle(op)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("[Main] shutting down...")

	if err := kc.Stop(); err != nil {
		log.Printf("[Main] collector stop error: %v", err)
	}

	close(out)
	close(oppCh)

	log.Println("[Main] exited")
}
