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

func NormalizeSymbol_NoAlloc(s string, buf *[]byte) string {
	b := (*buf)[:0]

	// конвертируем в верхний регистр и убираем пробелы
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == ' ' {
			continue
		}
		if c >= 'a' && c <= 'z' {
			c -= 'a' - 'A'
		}
		b = append(b, c)
	}

	if len(b) < 2 {
		return ""
	}

	// ищем разделитель '-' или '/'
	for i := 0; i < len(b)-1; i++ {
		if b[i] == '-' || b[i] == '/' {
			base := b[:i]
			quote := b[i+1:]
			if len(base) == 0 || len(quote) == 0 {
				return ""
			}
			*buf = append(base, '/')
			*buf = append(*buf, quote...)
			return string(*buf)
		}
	}

	// ищем слитные символы с известными quote
	for _, q := range knownQuotes {
		lq := len(q)
		if len(b) > lq && string(b[len(b)-lq:]) == string(q) {
			base := b[:len(b)-lq]
			if len(base) == 0 {
				return ""
			}
			*buf = append(base, '/')
			*buf = append(*buf, q...)
			return string(*buf)
		}
	}

	return ""
}
