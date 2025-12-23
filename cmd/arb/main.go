package main

import (
	"crypt_proto/internal/collector"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.Println("EXCHANGE: mexc")
	log.Println("Starting collector: mexc")

	mexc := collector.NewMEXCCollector([]string{
		"BTCUSDT",
		"ETHUSDT",
		"ETHBTC",
	})

	if err := mexc.Start(); err != nil {
		log.Fatal(err)
	}

	// корректное завершение по Ctrl+C
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	log.Println("Stopping collector...")
	mexc.Stop()
}
