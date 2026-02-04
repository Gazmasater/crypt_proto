Если оставить только нужное:

p99 execution latency
Micro-volatility (100 мс)
Fill ratio
Capture rate
Inventory drift




Название API
9623527002

696935c42a6dcd00013273f2
b348b686-55ff-4290-897b-02d55f815f65




apikey = "4333ed4b-cd83-49f5-97d1-c399e2349748"
secretkey = "E3848531135EDB4CCFDA0F1BC14CD274"
IP = ""
Название API-ключа = "Arb"
Доступы = "Чтение"



sudo systemctl mask sleep.target suspend.target hibernate.target hybrid-sleep.target



wbs-api.mexc.com/ws 


[https://edis-global.vercel.app/ru/vps-hosting/singapore-singapore
](https://sg.edisglobal.com/)



git pull --rebase origin privat
git push origin privat


import (
    // ...
    "net/http"
    _ "net/http/pprof"
)


   // pprof HTTP-сервер
    go func() {
        log.Println("pprof on http://localhost:6060/debug/pprof/")
        if err := http.ListenAndServe("localhost:6060", nil); err != nil {
            log.Printf("pprof server error: %v", err)
        }
    }()


	go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30


(pprof) top        # показать топ функций по CPU
(pprof) top10
(pprof) list parsePBWrapperMid   # подробный разбор одной функции
(pprof) quit


go tool pprof http://localhost:6060/debug/pprof/heap


(pprof) top
(pprof) top -cum
(pprof) list parsePBWrapperMid
(pprof) quit




go run -race main.go


GOMAXPROCS=8 go run -race main.go




package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"time"
)

const (
	symbol      = "XBTUSDTM"
	futuresBase = "https://api-futures.kucoin.com"
	spotBase    = "https://api.kucoin.com"
)

type tsResp struct {
	Code string `json:"code"`
	Data int64  `json:"data"` // ms
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

func getServerTimeMs() (int64, error) {
	resp, err := http.Get(spotBase + "/api/v1/timestamp")
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

func fetchKlines(granularity string, fromMs, toMs int64) ([][]any, error) {
	u, _ := url.Parse(futuresBase + "/api/v1/kline/query")
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
	return r.Data, nil
}

func fetchOI15m(startAt, endAt int64) ([]struct {
	OpenInterest string `json:"openInterest"`
	Ts           int64  `json:"ts"`
}, error) {
	u, _ := url.Parse(spotBase + "/api/ua/v1/market/open-interest")
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
	fmt.Println("=== MAIN: 1H(24h) + 15m(8h) + OI15m(20h) ===") // маркер, чтобы точно видеть что это этот файл

	nowMs, err := getServerTimeMs()
	if err != nil {
		panic(err)
	}

	// ===== 1H range (24h) =====
	kl1h, err := fetchKlines("60", nowMs-24*60*60*1000, nowMs)
	if err != nil {
		panic(err)
	}
	if len(kl1h) == 0 {
		panic("no 1h klines")
	}

	hi1 := -math.MaxFloat64
	lo1 := math.MaxFloat64
	var latestTs int64 = -1
	var latestClose float64

	for _, row := range kl1h {
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
		closeV, err := toFloat(row[2])
		if err != nil {
			panic(err)
		}

		if high > hi1 {
			hi1 = high
		}
		if low < lo1 {
			lo1 = low
		}
		if ts > latestTs {
			latestTs = ts
			latestClose = closeV
		}
	}

	mid1 := (hi1 + lo1) / 2
	rangeAbs1 := hi1 - lo1
	rangePct1 := (rangeAbs1 / mid1) * 100
	distToResPct := (hi1 - latestClose) / latestClose * 100
	distToSupPct := (latestClose - lo1) / latestClose * 100

	fmt.Printf("1H range (last 24h) %s:\n", symbol)
	fmt.Printf("R_high = %.2f\nR_low  = %.2f\nR_mid  = %.2f\nRange  = %.2f (%.2f%%)\n",
		hi1, lo1, mid1, rangeAbs1, rangePct1)
	fmt.Printf("Now    = %.2f (latest 1H close, ts=%d)\n", latestClose, latestTs)
	fmt.Printf("distToRes = %.3f%% | distToSup = %.3f%%\n", distToResPct, distToSupPct)
	const nearPct = 0.25
	fmt.Printf("Near resistance (<%.2f%%)? %v\n", nearPct, distToResPct < nearPct)
	fmt.Printf("Near support     (<%.2f%%)? %v\n\n", nearPct, distToSupPct < nearPct)

	// ===== 15m range (8h) =====
	kl15, err := fetchKlines("15", nowMs-8*60*60*1000, nowMs)
	if err != nil {
		panic(err)
	}
	if len(kl15) == 0 {
		panic("no 15m klines")
	}

	hi15 := -math.MaxFloat64
	lo15 := math.MaxFloat64
	for _, row := range kl15 {
		if len(row) < 5 {
			continue
		}
		high, err := toFloat(row[3])
		if err != nil {
			panic(err)
		}
		low, err := toFloat(row[4])
		if err != nil {
			panic(err)
		}
		if high > hi15 {
			hi15 = high
		}
		if low < lo15 {
			lo15 = low
		}
	}
	fmt.Printf("15m range (last 8h):\nH15 = %.2f\nL15 = %.2f\n\n", hi15, lo15)

	// ===== OI 15m (last 20h) =====
	oi, err := fetchOI15m(nowMs-20*60*60*1000, nowMs)
	if err != nil {
		panic(err)
	}
	if len(oi) < 3 {
		fmt.Printf("OI(15m) points returned: %d (need >=3 for ΔOI_30m)\n", len(oi))
		return
	}

	sort.Slice(oi, func(i, j int) bool { return oi[i].Ts < oi[j].Ts })

	last := oi[len(oi)-1]
	prev := oi[len(oi)-2]
	prev2 := oi[len(oi)-3]

	oiLast, _ := strconv.ParseFloat(last.OpenInterest, 64)
	oiPrev, _ := strconv.ParseFloat(prev.OpenInterest, 64)
	oiPrev2, _ := strconv.ParseFloat(prev2.OpenInterest, 64)

	dOI15 := oiLast - oiPrev
	dOI30 := oiLast - oiPrev2

	fmt.Printf("OI(15m) last points:\n")
	fmt.Printf("t-30m ts=%d oi=%.0f\n", prev2.Ts, oiPrev2)
	fmt.Printf("t-15m ts=%d oi=%.0f\n", prev.Ts, oiPrev)
	fmt.Printf("t     ts=%d oi=%.0f\n", last.Ts, oiLast)
	fmt.Printf("ΔOI_15m = %.0f | ΔOI_30m = %.0f\n", dOI15, dOI30)
}




