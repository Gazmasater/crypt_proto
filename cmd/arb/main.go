package main

import (
	"crypt_proto/internal/collector"
	"crypt_proto/pkg/models"
	"fmt"
	"time"
)

func main() {
	dataCh := make(chan models.MarketData, 100)
	okxCollector := collector.NewOKXCollector()

	err := okxCollector.Start(dataCh)
	if err != nil {
		fmt.Println("Ошибка запуска Collector:", err)
		return
	}

	// Вывод данных из канала
	go func() {
		for md := range dataCh {
			fmt.Printf("Collector: %s %s Bid=%.2f Ask=%.2f\n",
				md.Exchange, md.Symbol, md.Bid, md.Ask)
		}
	}()

	time.Sleep(10 * time.Second) // пусть Collector поработает
	okxCollector.Stop()
}
