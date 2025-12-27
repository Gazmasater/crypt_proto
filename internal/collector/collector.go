package collector

import "crypt_proto/pkg/models"

type Collector interface {
	Start(out chan<- *models.MarketData) error
	Stop() error
	Name() string
}
