package market

import "strings"

var knownQuotes = []string{
	"USDT", "USDC", "USD", "EUR", "BTC", "ETH",
}

// NormalizeSymbol_Full приводит символ к виду BASE/QUOTE
// возвращает "" если формат неполный или некорректный
func NormalizeSymbol_Full(s string) string {
	s = strings.TrimSpace(strings.ToUpper(s))
	if s == "" {
		return ""
	}

	// уже с разделителем
	if strings.ContainsAny(s, "-/") {
		s = strings.ReplaceAll(s, "-", "/")
		parts := strings.Split(s, "/")
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return ""
		}
		return parts[0] + "/" + parts[1]
	}

	// слитные символы, пытаемся найти quote
	for _, q := range knownQuotes {
		if strings.HasSuffix(s, q) && len(s) > len(q) {
			base := strings.TrimSuffix(s, q)
			if base == "" {
				return ""
			}
			return base + "/" + q
		}
	}

	// неизвестный формат — возвращаем "" (т.к. без quote)
	return ""
}
