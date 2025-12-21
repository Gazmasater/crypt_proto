package collector

import "crypt_proto/pkg/models"

// Collector — интерфейс для любого коллектора биржи
type Collector interface {
	// Start запускает сбор данных и отправку в канал dataCh
	Start(dataCh chan<- models.MarketData) error

	// Stop останавливает сбор данных
	Stop() error
}
