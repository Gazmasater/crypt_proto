package kucoin

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	FuturesBase = "https://api-futures.kucoin.com"
	SpotBase    = "https://api.kucoin.com"
)

type Candle struct {
	Ts    int64 // close time ms
	Open  float64
	High  float64
	Low   float64
	Close float64
	Vol   float64
}

type OIPoint struct {
	Ts int64
	OI float64
}

type tsResp struct {
	Code string `json:"code"`
	Data int64  `json:"data"`
}

type klineResp struct {
	Code string  `json:"code"`
	Data [][]any `json:"data"` // [ts, open, close, high, low, vol, turnover]
}

type oiResp struct {
	Code string `json:"code"`
	Data []struct {
		OpenInterest string `json:"openInterest"`
		Ts           int64  `json:"ts"`
	} `json:"data"`
}

func ServerTimeMs() (int64, error) {
	resp, err := http.Get(SpotBase + "/api/v1/timestamp")
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

func FetchKlines(symbol string, granularity string, fromMs, toMs int64) ([]Candle, error) {
	u, _ := url.Parse(FuturesBase + "/api/v1/kline/query")
	q := u.Query()
	q.Set("symbol", symbol)
	q.Set("granularity", granularity) // futures: 60=1H, 15=15m, 5=5m
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

	out := make([]Candle, 0, len(r.Data))
	for _, row := range r.Data {
		if len(row) < 6 {
			continue
		}
		ts, _ := toInt64(row[0])
		open, _ := toFloat(row[1])
		closeV, _ := toFloat(row[2])
		high, _ := toFloat(row[3])
		low, _ := toFloat(row[4])
		vol, _ := toFloat(row[5])

		out = append(out, Candle{
			Ts: ts, Open: open, High: high, Low: low, Close: closeV, Vol: vol,
		})
	}
	return out, nil
}

func FetchOI15m(symbol string, startAt, endAt int64) ([]OIPoint, error) {
	u, _ := url.Parse(SpotBase + "/api/ua/v1/market/open-interest")
	q := u.Query()
	q.Set("symbol", symbol)
	q.Set("interval", "15min")
	q.Set("startAt", strconv.FormatInt(startAt, 10))
	q.Set("endAt", strconv.FormatInt(endAt, 10))
	q.Set("pageSize", "200")
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
	var r oiResp
	if err := json.Unmarshal(b, &r); err != nil {
		return nil, err
	}
	if r.Code != "200000" {
		return nil, fmt.Errorf("oi bad code=%s body=%s", r.Code, string(b))
	}

	out := make([]OIPoint, 0, len(r.Data))
	for _, it := range r.Data {
		v, err := strconv.ParseFloat(it.OpenInterest, 64)
		if err != nil {
			continue
		}
		out = append(out, OIPoint{Ts: it.Ts, OI: v})
	}
	return out, nil
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
