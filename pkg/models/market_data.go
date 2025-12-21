package models

// MarketData хранит данные с биржи для одного инструмента
type MarketData struct {
	Exchange  string  `json:"exchange"`  // название биржи
	Symbol    string  `json:"symbol"`    // торговая пара, например BTC-USDT
	Bid       float64 `json:"bid"`       // лучшая цена покупки
	Ask       float64 `json:"ask"`       // лучшая цена продажи
	Timestamp int64   `json:"timestamp"` // метка времени в миллисекундах
}
