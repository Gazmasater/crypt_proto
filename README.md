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
	"fmt"
	"math"
	"os"
	"time"
)

// ======== RingBuffer для котировок ========
type RingBuffer struct {
	data  []float64
	size  int
	index int
	full  bool
}

func NewRingBuffer(size int) *RingBuffer {
	return &RingBuffer{
		data: make([]float64, size),
		size: size,
	}
}

func (r *RingBuffer) Add(value float64) {
	r.data[r.index] = value
	r.index = (r.index + 1) % r.size
	if r.index == 0 {
		r.full = true
	}
}

func (r *RingBuffer) Len() int {
	if r.full {
		return r.size
	}
	return r.index
}

func (r *RingBuffer) GetAll() []float64 {
	if !r.full {
		return r.data[:r.index]
	}
	res := make([]float64, r.size)
	copy(res, r.data[r.index:])
	copy(res[r.size-r.index:], r.data[:r.index])
	return res
}

// ======== Расчёт коэффициента и спреда ========
func CoefBTCETH(btcBuf, ethBuf *RingBuffer) float64 {
	btc := btcBuf.GetAll()
	eth := ethBuf.GetAll()
	if len(btc) == 0 || len(eth) == 0 {
		return 0
	}
	return btc[len(btc)-1] / eth[len(eth)-1]
}

func MinMaxCoef(btcBuf, ethBuf *RingBuffer) (float64, float64) {
	btc := btcBuf.GetAll()
	eth := ethBuf.GetAll()
	if len(btc) == 0 || len(eth) == 0 {
		return 0, 0
	}
	min := btc[0] / eth[0]
	max := min
	for i := 1; i < len(btc) && i < len(eth); i++ {
		c := btc[i] / eth[i]
		if c > max {
			max = c
		}
		if c < min {
			min = c
		}
	}
	return min, max
}

func PearsonCorr(btcBuf, ethBuf *RingBuffer) float64 {
	btc := btcBuf.GetAll()
	eth := ethBuf.GetAll()
	n := float64(len(btc))
	if n == 0 {
		return 0
	}

	var sumX, sumY, sumXY, sumX2, sumY2 float64
	for i := 0; i < len(btc); i++ {
		x := btc[i]
		y := eth[i]
		sumX += x
		sumY += y
		sumXY += x * y
		sumX2 += x * x
		sumY2 += y * y
	}

	numerator := sumXY - (sumX*sumY)/n
	denominator := math.Sqrt((sumX2-(sumX*sumX)/n)*(sumY2-(sumY*sumY)/n))
	if denominator == 0 {
		return 0
	}
	return numerator / denominator
}

// ======== Проверка сигнала с порогом 0.6% ========
func CheckSignal(btcBuf, ethBuf *RingBuffer, spreadThresholdPct, minCorr float64, f *os.File) {
	currentCoef := CoefBTCETH(btcBuf, ethBuf)
	minCoef, maxCoef := MinMaxCoef(btcBuf, ethBuf)
	corr := PearsonCorr(btcBuf, ethBuf)

	midCoef := (minCoef + maxCoef) / 2
	spreadPct := (currentCoef - midCoef) / midCoef * 100 // спред в процентах

	fmt.Printf("[%s] Corr=%.5f | CurrentCoef=%.5f | MinCoef=%.5f | MaxCoef=%.5f | Spread=%.2f%%\n",
		time.Now().Format("15:04"), corr, currentCoef, minCoef, maxCoef, spreadPct)

	if corr < minCorr {
		return
	}

	var signal string
	if spreadPct > spreadThresholdPct {
		signal = fmt.Sprintf("[%s] SELL BTC / BUY ETH | Coef=%.5f | Corr=%.5f | Spread=%.2f%%\n",
			time.Now().Format("15:04"), currentCoef, corr, spreadPct)
	} else if spreadPct < -spreadThresholdPct {
		signal = fmt.Sprintf("[%s] BUY BTC / SELL ETH | Coef=%.5f | Corr=%.5f | Spread=%.2f%%\n",
			time.Now().Format("15:04"), currentCoef, corr, spreadPct)
	} else {
		signal = fmt.Sprintf("[%s] NO SIGNAL | Coef=%.5f | Corr=%.5f | Spread=%.2f%%\n",
			time.Now().Format("15:04"), currentCoef, corr, spreadPct)
	}

	fmt.Print(signal)
	if _, err := f.WriteString(signal); err != nil {
		fmt.Println("Ошибка записи в файл:", err)
	}
}

// ======== Лог метрик каждые 5 минут ========
func LogMetrics(btcBuf, ethBuf *RingBuffer, f *os.File) {
	currentCoef := CoefBTCETH(btcBuf, ethBuf)
	minCoef, maxCoef := MinMaxCoef(btcBuf, ethBuf)
	corr := PearsonCorr(btcBuf, ethBuf)
	midCoef := (minCoef + maxCoef) / 2
	spreadPct := (currentCoef - midCoef) / midCoef * 100

	line := fmt.Sprintf("[%s] COEF=%.5f | Min=%.5f | Max=%.5f | Corr=%.5f | Spread=%.2f%%\n",
		time.Now().Format("15:04"), currentCoef, minCoef, maxCoef, corr, spreadPct)
	fmt.Print(line)
	if _, err := f.WriteString(line); err != nil {
		fmt.Println("Ошибка записи в файл:", err)
	}
}



package main

import (
	"log"
	"os"
	"time"

	"crypt_proto/pkg/collector"
	"crypt_proto/pkg/models"
)

func main() {
	const windowSize = 120 // 2 часа
	btcBuf := NewRingBuffer(windowSize)
	ethBuf := NewRingBuffer(windowSize)

	f, err := os.OpenFile("signals.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	out := make(chan *models.MarketData, 1000)

	coll, _, err := collector.NewKuCoinCollectorFromCSV("pairs.csv")
	if err != nil {
		panic(err)
	}
	if err := coll.Start(out); err != nil {
		panic(err)
	}
	defer coll.Stop()

	spreadThresholdPct := 0.6
	minCorr := 0.85

	logTicker := time.NewTicker(5 * time.Minute)

	for {
		select {
		case md := <-out:
			switch md.Symbol {
			case "BTC-USDT":
				btcBuf.Add(md.Bid)
			case "ETH-USDT":
				ethBuf.Add(md.Bid)
			}

			if btcBuf.Len() >= windowSize && ethBuf.Len() >= windowSize {
				CheckSignal(btcBuf, ethBuf, spreadThresholdPct, minCorr, f)
			}

		case <-logTicker.C:
			if btcBuf.Len() >= windowSize && ethBuf.Len() >= windowSize {
				LogMetrics(btcBuf, ethBuf, f)
			}
		}
	}
}

