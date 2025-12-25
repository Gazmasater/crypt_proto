package market

func Key(exchange, symbol string) string {
	return exchange + ":" + NormalizeSymbol_Full(symbol)
}
