package common

type Market struct {
	Symbol string

	Base  string
	Quote string

	EnableTrading bool

	BaseMinSize  float64
	QuoteMinSize float64

	BaseIncrement  float64
	QuoteIncrement float64
	PriceIncrement float64
}
