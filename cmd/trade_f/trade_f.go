package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	symbol      = "XBTUSDTM"
	futuresBase = "https://api-futures.kucoin.com"
)

type tsResp struct {
	Code string `json:"code"`
	Data int64  `json:"data"` // ms
}

type klineResp struct {
	Code string  `json:"code"`
	Data [][]any `json:"data"` // [ts, open, close, high, low, vol, turnover]
}

func getServerTimeMs() (int64, error) {
	resp, err := http.Get("https://api.kucoin.com/api/v1/timestamp")
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	b, _ := io.ReadAll(resp.Body)
	var r tsResp
	if err := json.Unmarshal(b, &r); err != nil {
		return 0, err
	}
	if r.Code != "200000" {
		return 0, fmt.Errorf("timestamp bad code=%s body=%s", r.Code, string(b))
	}
	return r.Data, nil
}

func fetchKlines1H(fromMs, toMs int64) ([][]any, error) {
	u, _ := url.Parse(futuresBase + "/api/v1/kline/query")
	q := u.Query()
	q.Set("symbol", symbol)
	q.Set("granularity", "60") // ✅ 1H for KuCoin Futures kline/query
	q.Set("from", strconv.FormatInt(fromMs, 10))
	q.Set("to", strconv.FormatInt(toMs, 10))
	u.RawQuery = q.Encode()

	req, _ := http.NewRequest("GET", u.String(), nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, _ := io.ReadAll(resp.Body)
	var r klineResp
	if err := json.Unmarshal(b, &r); err != nil {
		return nil, err
	}
	if r.Code != "200000" {
		return nil, fmt.Errorf("kline bad code=%s body=%s", r.Code, string(b))
	}
	return r.Data, nil
}

func toFloat(v any) (float64, error) {
	switch t := v.(type) {
	case string:
		return strconv.ParseFloat(t, 64)
	case float64:
		return t, nil
	case json.Number:
		return t.Float64()
	default:
		return 0, fmt.Errorf("unexpected type %T", v)
	}
}

func toInt64(v any) (int64, error) {
	switch t := v.(type) {
	case string:
		return strconv.ParseInt(t, 10, 64)
	case float64:
		return int64(t), nil
	case json.Number:
		return t.Int64()
	default:
		return 0, fmt.Errorf("unexpected type %T", v)
	}
}

func main() {
	nowMs, err := getServerTimeMs()
	if err != nil {
		panic(err)
	}

	fromMs := nowMs - 24*60*60*1000 // 24 часа назад
	toMs := nowMs

	klines, err := fetchKlines1H(fromMs, toMs)
	if err != nil {
		panic(err)
	}
	if len(klines) == 0 {
		panic("no klines returned")
	}

	// KuCoin futures kline: [ts, open, close, high, low, vol, turnover]
	hi := -math.MaxFloat64
	lo := math.MaxFloat64

	// найдём самую "новую" свечу по ts, чтобы взять её close как текущую цену
	var latestTs int64 = -1
	var latestClose float64

	for _, row := range klines {
		if len(row) < 5 {
			continue
		}

		ts, err := toInt64(row[0])
		if err != nil {
			panic(err)
		}

		high, err := toFloat(row[3])
		if err != nil {
			panic(err)
		}
		low, err := toFloat(row[4])
		if err != nil {
			panic(err)
		}
		if high > hi {
			hi = high
		}
		if low < lo {
			lo = low
		}

		// close
		if len(row) > 2 {
			closeV, err := toFloat(row[2])
			if err != nil {
				panic(err)
			}
			if ts > latestTs {
				latestTs = ts
				latestClose = closeV
			}
		}
	}

	if hi == -math.MaxFloat64 || lo == math.MaxFloat64 {
		panic("failed to compute range: empty/invalid rows")
	}

	mid := (hi + lo) / 2
	rangeAbs := hi - lo
	rangePct := (rangeAbs / mid) * 100

	// расстояние до стен (в % от текущей цены)
	distToResPct := (hi - latestClose) / latestClose * 100
	distToSupPct := (latestClose - lo) / latestClose * 100

	fmt.Printf("1H range (last 24h) %s:\n", symbol)
	fmt.Printf("R_high = %.2f\n", hi)
	fmt.Printf("R_low  = %.2f\n", lo)
	fmt.Printf("R_mid  = %.2f\n", mid)
	fmt.Printf("Range  = %.2f (%.2f%%)\n", rangeAbs, rangePct)
	fmt.Printf("Now    = %.2f (latest 1H close, ts=%d)\n", latestClose, latestTs)
	fmt.Printf("distToRes = %.3f%% | distToSup = %.3f%%\n", distToResPct, distToSupPct)

	// “впритык” фильтр (стартовый)
	const nearPct = 0.25
	fmt.Printf("Near resistance (<%.2f%%)? %v\n", nearPct, distToResPct < nearPct)
	fmt.Printf("Near support     (<%.2f%%)? %v\n", nearPct, distToSupPct < nearPct)
}
