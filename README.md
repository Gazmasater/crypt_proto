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
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"time"
)

// ===== Структуры =====
type Candle struct {
	Time  int64
	Close float64
}

type RingBuffer struct {
	Data  []Candle
	Size  int
	Index int
}

func (r *RingBuffer) Add(c Candle) {
	r.Data[r.Index] = c
	r.Index = (r.Index + 1) % r.Size
}

func (r *RingBuffer) GetAll() []Candle {
	result := make([]Candle, 0, r.Size)
	for i := 0; i < r.Size; i++ {
		idx := (r.Index + i) % r.Size
		result = append(result, r.Data[idx])
	}
	return result
}

// ====== KuCoin REST API ======
func getKlines(symbol string, interval string, startAt, endAt int64) ([]Candle, error) {
	url := fmt.Sprintf(
		"https://api.kucoin.com/api/v1/market/candles?symbol=%s&type=%s&startAt=%d&endAt=%d",
		symbol, interval, startAt, endAt,
	)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var data [][]string
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	candles := make([]Candle, 0, len(data))
	for _, item := range data {
		var ts int64
		var closePrice float64
		ts, _ = strconv.ParseInt(item[0], 10, 64)
		closePrice, _ = strconv.ParseFloat(item[2], 64)
		candles = append(candles, Candle{Time: ts, Close: closePrice})
	}
	return candles, nil
}

// Получаем последнюю цену (level1)
func getLastPrice(symbol string) (float64, error) {
	url := fmt.Sprintf("https://api.kucoin.com/api/v1/market/orderbook/level1?symbol=%s", symbol)
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var result struct {
		Code string `json:"code"`
		Data struct {
			Price string `json:"price"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return 0, err
	}

	price, err := strconv.ParseFloat(result.Data.Price, 64)
	if err != nil {
		return 0, err
	}

	return price, nil
}

// ====== Математика ======
func getCoef(btc, eth float64) float64 {
	return btc / eth
}

func getMinMaxCoef(btc, eth *RingBuffer) (float64, float64) {
	btcCandles := btc.GetAll()
	ethCandles := eth.GetAll()
	min := btcCandles[0].Close / ethCandles[0].Close
	max := min
	for i := 1; i < len(btcCandles); i++ {
		c := btcCandles[i].Close / ethCandles[i].Close
		if c > max {
			max = c
		}
		if c < min {
			min = c
		}
	}
	return min, max
}

func pearsonCorrelation(btc, eth *RingBuffer) float64 {
	btcCandles := btc.GetAll()
	ethCandles := eth.GetAll()
	n := float64(len(btcCandles))

	var sumX, sumY, sumXY, sumX2, sumY2 float64
	for i := 0; i < len(btcCandles); i++ {
		x := btcCandles[i].Close
		y := ethCandles[i].Close
		sumX += x
		sumY += y
		sumXY += x * y
		sumX2 += x * x
		sumY2 += y * y
	}

	numerator := sumXY - (sumX*sumY)/n
	denominator := math.Sqrt((sumX2 - (sumX*sumX)/n) * (sumY2 - (sumY*sumY)/n))
	if denominator == 0 {
		return 0
	}
	return numerator / denominator
}

// ====== Сигнал на вход ======
func checkSignal(btc, eth *RingBuffer, spread, minCorr float64) {
	lastBTC := btc.GetAll()[len(btc.Data)-1].Close
	lastETH := eth.GetAll()[len(eth.Data)-1].Close

	corr := pearsonCorrelation(btc, eth)
	if corr < minCorr {
		fmt.Printf("[%s] Корреляция %.2f ниже порога %.2f, сигнал пропущен\n", time.Now().Format("15:04"), corr, minCorr)
		return
	}

	currentCoef := getCoef(lastBTC, lastETH)
	minCoef, maxCoef := getMinMaxCoef(btc, eth)

	if currentCoef > maxCoef+spread {
		fmt.Printf("[%s] Сигнал: SELL BTC / BUY ETH | Coef=%.5f\n", time.Now().Format("15:04"), currentCoef)
	} else if currentCoef < minCoef-spread {
		fmt.Printf("[%s] Сигнал: BUY BTC / SELL ETH | Coef=%.5f\n", time.Now().Format("15:04"), currentCoef)
	} else {
		fmt.Printf("[%s] Сигнал: нет действия | Coef=%.5f\n", time.Now().Format("15:04"), currentCoef)
	}
}

// ====== Main ======
func main() {
	const windowSize = 120 // 2 часа × 1 мин
	btcBuffer := RingBuffer{Data: make([]Candle, windowSize), Size: windowSize}
	ethBuffer := RingBuffer{Data: make([]Candle, windowSize), Size: windowSize}

	// Загружаем историю 2 часа
	endAt := time.Now().Unix()
	startAt := endAt - int64(windowSize*60)

	btcCandles, err := getKlines("BTC-USDT", "1min", startAt, endAt)
	if err != nil {
		panic(err)
	}
	ethCandles, err := getKlines("ETH-USDT", "1min", startAt, endAt)
	if err != nil {
		panic(err)
	}

	for _, c := range btcCandles {
		btcBuffer.Add(c)
	}
	for _, c := range ethCandles {
		ethBuffer.Add(c)
	}

	spread := 0.001   // 0.1%
	minCorr := 0.85   // порог корреляции

	oneMinuteTicker := time.NewTicker(1 * time.Minute)
	fiveMinuteTicker := time.NewTicker(5 * time.Minute)

	for {
		select {
		case <-oneMinuteTicker.C:
			// Обновляем последнюю свечу из Level1
			btcPrice, err := getLastPrice("BTC-USDT")
			if err != nil {
				fmt.Println("Ошибка BTC:", err)
				continue
			}
			ethPrice, err := getLastPrice("ETH-USDT")
			if err != nil {
				fmt.Println("Ошибка ETH:", err)
				continue
			}
			now := time.Now().Unix()
			btcBuffer.Add(Candle{Time: now, Close: btcPrice})
			ethBuffer.Add(Candle{Time: now, Close: ethPrice})

		case <-fiveMinuteTicker.C:
			checkSignal(&btcBuffer, &ethBuffer, spread, minCorr)
		}
	}
}


gaz358@gaz358-BOD-WXX9:~/myprog/crypt_proto/cmd/stat_arb$ go run .
panic: json: cannot unmarshal object into Go value of type [][]string

goroutine 1 [running]:
main.main()
        /home/gaz358/myprog/crypt_proto/cmd/stat_arb/stat_arb.go:178 +0x60a
exit status 2

