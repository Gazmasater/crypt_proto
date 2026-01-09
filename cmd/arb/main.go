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
	"crypt_proto/internal/queue"
	"crypt_proto/pkg/models"
)

func main() {
	// ------------------- pprof -------------------
	go func() {
		log.Println("pprof on http://localhost:6060/debug/pprof/")
		_ = http.ListenAndServe("localhost:6060", nil)
	}()

	// ------------------- Канал данных от коллекторов -------------------
	out := make(chan *models.MarketData, 100_000)

	// ------------------- In-Memory Store -------------------
	mem := queue.NewMemoryStore()
	go mem.Run()

	// прокачка данных: out → mem
	go func() {
		for md := range out {
			mem.Push(md)
		}
	}()

	// ------------------- Коллектор -------------------
	kc, err := collector.NewKuCoinCollectorFromCSV(
		"../exchange/data/kucoin_triangles_usdt.csv",
	)
	if err != nil {
		log.Fatal(err)
	}

	if err := kc.Start(out); err != nil {
		log.Fatal(err)
	}
	log.Println("[Main] KuCoinCollector started")

	// ------------------- Треугольники -------------------
	triangles, err := calculator.ParseTrianglesFromCSV(
		"../exchange/data/kucoin_triangles_usdt.csv",
	)
	if err != nil {
		log.Fatal(err)
	}

	// ------------------- Калькулятор -------------------
	calc := calculator.NewCalculator(mem, triangles)
	go calc.Run()

	// ------------------- Graceful shutdown -------------------
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("[Main] shutting down...")

	kc.Stop()
	close(out)

	log.Println("[Main] exited")
}
