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
	go func() {
		log.Println("pprof on http://localhost:6060/debug/pprof/")
		_ = http.ListenAndServe("localhost:6060", nil)
	}()

	// ------------------- Канал для данных -------------------
	out := make(chan *models.MarketData, 100_000)

	// ------------------- In-Memory Store -------------------
	mem := queue.NewMemoryStore()
	go mem.Run(out)

	// ------------------- Коллектор -------------------
	kc, err := collector.NewKuCoinCollectorFromCSV("../exchange/data/kucoin_triangles_usdt.csv")
	if err != nil {
		log.Fatal(err)
	}

	if err := kc.Start(out); err != nil {
		log.Fatal(err)
	}
	log.Println("[Main] KuCoinCollector started")

	triangles, err := calculator.ParseTrianglesFromCSV(
		"../exchange/data/kucoin_triangles_usdt.csv",
	)
	if err != nil {
		log.Fatal(err)
	}

	calc := calculator.NewCalculator(mem, triangles)
	go calc.Run()

	// ------------------- Вывод snapshot для проверки -------------------
	//go func() {
	//	for {
	//		snap := mem.Snapshot()
	//		log.Printf("[Store] quotes=%d", len(snap))
	//		time.Sleep(5 * time.Second)
	//	}
	//}()

	// ------------------- Завершение при SIGINT / SIGTERM -------------------
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	kc.Stop()
	close(out)
	log.Println("[Main] Exited.")
}
