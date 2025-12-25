package market

import "strings"

var knownQuotes = []string{
	"USDT",
	"USDC",
	"USD",
	"EUR",
	"BTC",
	"ETH",
}

func NormalizeSymbol(s string) string {
	s = strings.TrimSpace(strings.ToUpper(s))
	if s == "" {
		return s
	}

	// если уже с разделителем
	if strings.ContainsAny(s, "-/") {
		s = strings.ReplaceAll(s, "-", "/")
		parts := strings.Split(s, "/")
		if len(parts) == 2 && parts[0] != "" && parts[1] != "" {
			return parts[0] + "/" + parts[1]
		}
		return s
	}

	// формат BTCUSDT → BTC/USDT
	for _, q := range knownQuotes {
		if strings.HasSuffix(s, q) && len(s) > len(q) {
			base := strings.TrimSuffix(s, q)
			return base + "/" + q
		}
	}

	// неизвестный формат → как есть
	return s
}
