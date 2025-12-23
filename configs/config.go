package configs

import "time"

// --- MEXC ---
const (
	MEXC_WS            = "wss://wbs-api.mexc.com/ws"
	MEXC_READ_TIMEOUT  = 30 * time.Second
	MEXC_PING_INTERVAL = 10 * time.Second
	MEXC_RECONNECT_DUR = time.Second
)

// --- KuCoin ---
const (
	KUCOIN_WS            = "wss://ws.kucoin.com/endpoint"
	KUCOIN_PING_INTERVAL = 30 * time.Second
	KUCOIN_READ_TIMEOUT  = 60 * time.Second
)

// --- OKX ---
const (
	OKX_WS            = "wss://ws.okx.com:8443/ws/v5/public"
	OKX_PING_INTERVAL = 20 * time.Second
	OKX_READ_TIMEOUT  = 60 * time.Second
)
