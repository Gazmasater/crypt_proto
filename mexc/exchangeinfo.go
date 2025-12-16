package mexc

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type SymbolCaps struct {
	Symbol      string
	Status      string
	HasMarket   bool
	StepSize    float64
	MinQty      float64
	MinNotional float64
}

func FetchSymbolCapsMEXC(ctx context.Context) (map[string]SymbolCaps, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.mexc.com/api/v3/exchangeInfo", nil)
	if err != nil {
		return nil, err
	}

	cl := &http.Client{Timeout: 12 * time.Second}
	resp, err := cl.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("exchangeInfo status=%d", resp.StatusCode)
	}

	var root map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&root); err != nil {
		return nil, err
	}

	rawSyms, _ := root["symbols"].([]any)
	out := make(map[string]SymbolCaps, len(rawSyms))

	for _, item := range rawSyms {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}

		symbol, _ := m["symbol"].(string)
		if symbol == "" {
			continue
		}

		status, _ := m["status"].(string)

		// safe default: если orderTypes нет — не режем
		hasMarket := true
		if otsAny, ok := m["orderTypes"]; ok {
			hasMarket = false
			if ots, ok := otsAny.([]any); ok {
				for _, v := range ots {
					if s, ok := v.(string); ok && strings.EqualFold(s, "MARKET") {
						hasMarket = true
						break
					}
				}
			} else {
				hasMarket = true
			}
		}

		stepSize, minQty, minNotional := 0.0, 0.0, 0.0
		if flt, ok := m["filters"].([]any); ok {
			for _, f := range flt {
				fm, ok := f.(map[string]any)
				if !ok {
					continue
				}
				ft, _ := fm["filterType"].(string)
				switch ft {
				case "LOT_SIZE":
					stepSize = readFloatAny(fm["stepSize"])
					minQty = readFloatAny(fm["minQty"])
				case "MIN_NOTIONAL":
					minNotional = readFloatAny(fm["minNotional"])
				}
			}
		}

		out[symbol] = SymbolCaps{
			Symbol:      symbol,
			Status:      status,
			HasMarket:   hasMarket,
			StepSize:    stepSize,
			MinQty:      minQty,
			MinNotional: minNotional,
		}
	}

	return out, nil
}

func readFloatAny(v any) float64 {
	switch t := v.(type) {
	case string:
		f, _ := strconv.ParseFloat(t, 64)
		return f
	case float64:
		return t
	default:
		return 0
	}
}
