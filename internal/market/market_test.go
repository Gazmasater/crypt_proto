package market

import "testing"

func TestNormalizeSymbol(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"BTCUSDT", "BTC/USDT"},
		{"btcusdt", "BTC/USDT"},
		{"BTC-USDT", "BTC/USDT"},
		{"eth-btc", "ETH/BTC"},
		{"ETHBTC", "ETHBTC"}, // неизвестный формат — не ломаем
		{"  btcusdt  ", "BTC/USDT"},
	}

	for _, tt := range tests {
		got := NormalizeSymbol(tt.in)
		if got != tt.want {
			t.Errorf("NormalizeSymbol(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestParsePair(t *testing.T) {
	tests := []struct {
		in        string
		wantBase  string
		wantQuote string
	}{
		{"BTC/USDT", "BTC", "USDT"},
		{"ETH/BTC", "ETH", "BTC"},
		{"INVALID", "", ""},
		{"BTC/", "", ""},
	}

	for _, tt := range tests {
		p := ParsePair(tt.in)
		if p.Base != tt.wantBase || p.Quote != tt.wantQuote {
			t.Errorf(
				"ParsePair(%q) = %+v, want Base=%q Quote=%q",
				tt.in, p, tt.wantBase, tt.wantQuote,
			)
		}
	}
}

func TestKey(t *testing.T) {
	tests := []struct {
		exchange string
		symbol   string
		want     string
	}{
		{"MEXC", "BTCUSDT", "MEXC:BTC/USDT"},
		{"OKX", "BTC-USDT", "OKX:BTC/USDT"},
		{"KuCoin", "eth-btc", "KuCoin:ETH/BTC"},
	}

	for _, tt := range tests {
		got := Key(tt.exchange, tt.symbol)
		if got != tt.want {
			t.Errorf("Key(%q, %q) = %q, want %q",
				tt.exchange, tt.symbol, got, tt.want)
		}
	}
}
