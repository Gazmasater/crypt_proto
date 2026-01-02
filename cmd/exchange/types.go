package exchange

// ==========================
// Универсальная модель рынка
// ==========================
type Market struct {
	Symbol string
	Base   string
	Quote  string

	EnableTrading bool

	BaseMinSize    string
	QuoteMinSize   string
	BaseIncrement  string
	QuoteIncrement string
	PriceIncrement string
}
