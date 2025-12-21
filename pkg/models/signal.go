package models

// Signal представляет сигнал арбитража
type Signal struct {
	ExchangeStart string  `json:"exchange_start"` // биржа первой сделки
	ExchangeMid   string  `json:"exchange_mid"`   // биржа второй сделки
	ExchangeEnd   string  `json:"exchange_end"`   // биржа третьей сделки
	SymbolStart   string  `json:"symbol_start"`   // первая монета
	SymbolMid     string  `json:"symbol_mid"`     // вторая монета
	SymbolEnd     string  `json:"symbol_end"`     // третья монета
	ProfitPercent float64 `json:"profit_percent"` // ожидаемая прибыль в процентах
	Amount        float64 `json:"amount"`         // объём сделки
	Timestamp     int64   `json:"timestamp"`      // время создания сигнала
}
