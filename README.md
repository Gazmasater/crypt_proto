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
	Code string        `json:"code"`
	Data [][]any       `json:"data"` // [ts, open, close, high, low, vol, turnover]
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
	q.Set("granularity", "3600") // 1h in seconds
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
	// API может отдавать string или number — обработаем оба
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

	for _, row := range klines {
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
		if high > hi {
			hi = high
		}
		if low < lo {
			lo = low
		}
	}

	fmt.Printf("1H range (last 24h) %s:\n", symbol)
	fmt.Printf("R_high = %.2f\n", hi)
	fmt.Printf("R_low  = %.2f\n", lo)
}


