package market

import "testing"

func TestNormalizeSymbol_Full(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"BTCUSDT", "BTC/USDT"},
		{"btcusdt", "BTC/USDT"},
		{"BTC-USDT", "BTC/USDT"},
		{"eth-btc", "ETH/BTC"},
		{"ETHBTC", "ETH/BTC"},
		{"BTC", ""},    // неполный символ
		{"BTC/", ""},   // неполный символ
		{"XYZABC", ""}, // неизвестный формат
		{"  btcusdt  ", "BTC/USDT"},
	}

	for _, tt := range tests {
		got := NormalizeSymbol_Full(tt.in)
		if got != tt.want {
			t.Errorf("NormalizeSymbol_Full(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestKey_Full(t *testing.T) {
	tests := []struct {
		exchange string
		symbol   string
		want     string
	}{
		{"MEXC", "BTCUSDT", "MEXC:BTC/USDT"},
		{"OKX", "BTC-USDT", "OKX:BTC/USDT"},
		{"KuCoin", "eth-btc", "KuCoin:ETH/BTC"},
		{"MEXC", "BTC", "MEXC:"},        // неполный символ
		{"KuCoin", "XYZABC", "KuCoin:"}, // неизвестный символ
	}

	for _, tt := range tests {
		got := Key(tt.exchange, tt.symbol)
		if got != tt.want {
			t.Errorf("Key_Full(%q, %q) = %q, want %q", tt.exchange, tt.symbol, got, tt.want)
		}
	}
}
