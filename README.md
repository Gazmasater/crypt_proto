sudo systemctl mask sleep.target suspend.target hibernate.target hybrid-sleep.target



wbs-api.mexc.com/ws 


[https://edis-global.vercel.app/ru/vps-hosting/singapore-singapore
](https://sg.edisglobal.com/)



git pull --rebase origin privat
git push origin privat


BOOK_INTERVAL=100ms
SYMBOLS_FILE=triangles_markets.csv
DEBUG=false


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



Ок, вот точный список файлов, которые нужно изменить у тебя в проекте, и полные переписанные версии (копипастишь 1-в-1).

1) config/config.go
package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

/* =========================  CONFIG  ========================= */

type Config struct {
	Exchange      string // "MEXC" или "KUCOIN"
	TrianglesFile string
	BookInterval  string
	FeePerLeg     float64 // как доля, 0.001 = 0.1%
	MinProfit     float64 // как доля, 0.003 = 0.3%

	// Минимальный стартовый объём (в валюте начала треугольника).
	// Обычно это USDT; если 0 - фильтр отключен.
	MinStart float64

	Debug bool
}

var debug bool

func SetDebug(v bool) {
	debug = v
}

func loadEnvFloat(name string, def float64) float64 {
	raw := strings.TrimSpace(os.Getenv(name))
	if raw == "" {
		return def
	}
	v, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		log.Printf("bad %s=%q: %v, using default %f", name, raw, err, def)
		return def
	}
	return v
}

func Load() Config {
	_ = godotenv.Load(".env")

	ex := strings.ToUpper(strings.TrimSpace(os.Getenv("EXCHANGE")))
	if ex == "" {
		ex = "MEXC"
	}

	tf := os.Getenv("TRIANGLES_FILE")
	if tf == "" {
		tf = "triangles_markets.csv"
	}
	bi := os.Getenv("BOOK_INTERVAL")
	if bi == "" {
		bi = "100ms"
	}

	// проценты
	feePct := loadEnvFloat("FEE_PCT", 0.04)
	minPct := loadEnvFloat("MIN_PROFIT_PCT", 0.5)

	// Минимальный стартовый объём (обычно USDT). Можно задавать как MIN_START_USDT
	// (предпочтительно), либо MIN_START. Если не задано - фильтр отключен.
	minStart := loadEnvFloat("MIN_START_USDT", -1)
	if minStart < 0 {
		minStart = loadEnvFloat("MIN_START", 0)
	}

	debug := strings.ToLower(os.Getenv("DEBUG")) == "true"

	cfg := Config{
		Exchange:      ex,
		TrianglesFile: tf,
		BookInterval:  bi,
		FeePerLeg:     feePct / 100.0,
		MinProfit:     minPct / 100.0,
		MinStart:      minStart,
		Debug:         debug,
	}

	log.Printf("Exchange: %s", cfg.Exchange)
	log.Printf("Triangles file: %s", tf)
	log.Printf("Book interval: %s", bi)
	log.Printf("Fee per leg: %.4f %% (rate=%.6f)", feePct, cfg.FeePerLeg)
	log.Printf("Min profit per cycle: %.4f %% (rate=%.6f)", minPct, cfg.MinProfit)
	log.Printf("Min start amount: %.4f", cfg.MinStart)

	return cfg
}

/* =========================  LOGGING  ========================= */

func Dlog(format string, args ...any) {
	if debug {
		log.Printf(format, args...)
	}
}

2) domain/domain.go

(в конце файла добавлены структуры и функция ComputeMaxStartTopOfBook, остальное оставлено как было — но чтобы было проще, даю весь файл целиком)

package domain

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
)

type Leg struct {
	From   string
	To     string
	Symbol string
	Dir    int8 // +1: From->To = base->quote; -1: From->To = quote->base
}

type Triangle struct {
	Legs [3]Leg
	Name string // A→B→C→A
}

type Quote struct {
	Bid    float64
	Ask    float64
	BidQty float64
	AskQty float64
}

type Event struct {
	Symbol string
	Bid    float64
	Ask    float64
	BidQty float64
	AskQty float64
}

type Pair struct {
	Base   string
	Quote  string
	Symbol string
}

func buildTriangleFromPairs(p1, p2, p3 Pair) (Triangle, bool) {
	set := map[string]struct{}{
		p1.Base:  {},
		p1.Quote: {},
		p2.Base:  {},
		p2.Quote: {},
		p3.Base:  {},
		p3.Quote: {},
	}
	if len(set) != 3 {
		return Triangle{}, false
	}
	currs := make([]string, 0, 3)
	for c := range set {
		currs = append(currs, c)
	}

	type edge struct{ From, To string }

	pairs := []Pair{p1, p2, p3}
	perm3 := [][]int{
		{0, 1, 2},
		{0, 2, 1},
		{1, 0, 2},
		{1, 2, 0},
		{2, 0, 1},
		{2, 1, 0},
	}

	for _, order := range perm3 {
		c0, c1, c2 := currs[order[0]], currs[order[1]], currs[order[2]]
		edges := []edge{
			{From: c0, To: c1},
			{From: c1, To: c2},
			{From: c2, To: c0},
		}

		for _, pp := range perm3 {
			var legs [3]Leg
			okAll := true

			for i := 0; i < 3; i++ {
				e := edges[i]
				p := pairs[pp[i]]

				switch {
				case p.Base == e.From && p.Quote == e.To:
					legs[i] = Leg{From: e.From, To: e.To, Symbol: p.Symbol, Dir: +1}
				case p.Base == e.To && p.Quote == e.From:
					legs[i] = Leg{From: e.From, To: e.To, Symbol: p.Symbol, Dir: -1}
				default:
					okAll = false
				}
				if !okAll {
					break
				}
			}

			if okAll {
				name := fmt.Sprintf("%s→%s→%s→%s", edges[0].From, edges[1].From, edges[2].From, edges[0].From)
				return Triangle{Legs: legs, Name: name}, true
			}
		}
	}

	return Triangle{}, false
}

// LoadTriangles читает CSV, строит треугольники и индекс по символам.
func LoadTriangles(path string) ([]Triangle, []string, map[string][]int, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, nil, err
	}
	defer f.Close()

	r := csv.NewReader(bufio.NewReader(f))
	r.TrimLeadingSpace = true
	r.Comma = ','

	var tris []Triangle
	symbolSet := make(map[string]struct{})

	for {
		rec, err := r.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, nil, nil, err
		}

		var fields []string
		for _, v := range rec {
			v = strings.TrimSpace(v)
			if v != "" {
				fields = append(fields, v)
			}
		}
		if len(fields) == 0 {
			continue
		}
		if strings.HasPrefix(fields[0], "#") {
			continue
		}
		if len(fields) != 6 {
			log.Printf("skip line (need 6 fields): %v", fields)
			continue
		}

		p1 := Pair{Base: fields[0], Quote: fields[1], Symbol: fields[0] + fields[1]}
		p2 := Pair{Base: fields[2], Quote: fields[3], Symbol: fields[2] + fields[3]}
		p3 := Pair{Base: fields[4], Quote: fields[5], Symbol: fields[4] + fields[5]}

		t, ok := buildTriangleFromPairs(p1, p2, p3)
		if !ok {
			continue
		}

		tris = append(tris, t)
		for _, leg := range t.Legs {
			symbolSet[leg.Symbol] = struct{}{}
		}
	}

	symbols := make([]string, 0, len(symbolSet))
	for s := range symbolSet {
		symbols = append(symbols, s)
	}

	index := make(map[string][]int)
	for i, t := range tris {
		for _, leg := range t.Legs {
			index[leg.Symbol] = append(index[leg.Symbol], i)
		}
	}

	log.Printf("треугольников всего: %d", len(tris))
	log.Printf("символов в индексе треугольников: %d", len(symbols))

	return tris, symbols, index, nil
}

// EvalTriangle считает доходность треугольника.
func EvalTriangle(t Triangle, quotes map[string]Quote, fee float64) (float64, bool) {
	amt := 1.0

	for _, leg := range t.Legs {
		q, ok := quotes[leg.Symbol]
		if !ok || q.Bid <= 0 || q.Ask <= 0 {
			return 0, false
		}

		if leg.Dir > 0 {
			amt *= q.Bid
		} else {
			amt /= q.Ask
		}

		amt *= (1 - fee)
		if amt <= 0 {
			return 0, false
		}
	}

	return amt - 1.0, true
}

// MaxStartInfo - диагностическая информация по максимальному стартовому объёму.
// Все значения рассчитаны по top-of-book (best bid/ask) без учёта округлений stepSize/minQty.
type MaxStartInfo struct {
	StartAsset    string
	MaxStart      float64    // в StartAsset
	BottleneckLeg int        // индекс ноги [0..2]
	LimitIn       [3]float64 // лимит на ВХОД каждой ноги (в единицах входного актива ноги)
	KIn           [3]float64 // сколько входного актива ноги получается из 1 StartAsset
	MaxStartByLeg [3]float64 // LimitIn / KIn
}

// ComputeMaxStartTopOfBook возвращает максимальный стартовый объём, который можно протащить через треугольник,
// не выходя за best bid/ask qty на каждой ноге. Комиссия учитывается как удержание из результата каждой ноги.
func ComputeMaxStartTopOfBook(t Triangle, quotes map[string]Quote, fee float64) (MaxStartInfo, bool) {
	var info MaxStartInfo
	info.StartAsset = t.Legs[0].From

	// kIn - сколько входной валюты текущей ноги получится из 1 единицы стартовой валюты.
	kIn := 1.0
	info.MaxStart = 0
	info.BottleneckLeg = -1

	// Инициализируем maxStart как +Inf, чтобы взять минимум по ногам.
	maxStart := 1e308

	for i, leg := range t.Legs {
		q, ok := quotes[leg.Symbol]
		if !ok || q.Bid <= 0 || q.Ask <= 0 {
			return MaxStartInfo{}, false
		}

		info.KIn[i] = kIn

		var limitIn float64
		var ratio float64 // out/in без комиссии

		if leg.Dir > 0 {
			// SELL base -> quote по bid, ограничение по количеству base на bid.
			if q.BidQty <= 0 {
				return MaxStartInfo{}, false
			}
			limitIn = q.BidQty
			ratio = q.Bid
		} else {
			// BUY base <- quote по ask, ограничение по объёму quote, который можно потратить: askQty*ask.
			if q.AskQty <= 0 {
				return MaxStartInfo{}, false
			}
			limitIn = q.AskQty * q.Ask
			ratio = 1.0 / q.Ask
		}

		info.LimitIn[i] = limitIn
		if kIn <= 0 {
			return MaxStartInfo{}, false
		}
		maxByThis := limitIn / kIn
		info.MaxStartByLeg[i] = maxByThis
		if maxByThis < maxStart {
			maxStart = maxByThis
			info.BottleneckLeg = i
		}

		// переход к следующей ноге
		kIn = kIn * ratio * (1 - fee)
		if kIn <= 0 {
			return MaxStartInfo{}, false
		}
	}

	if info.BottleneckLeg < 0 || maxStart <= 0 || maxStart > 1e307 {
		return MaxStartInfo{}, false
	}
	info.MaxStart = maxStart
	return info, true
}

3) arb/arb.go
package arb

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"crypt_proto/domain"
)

type Consumer struct {
	FeePerLeg float64
	MinProfit float64
	MinStart  float64

	writer io.Writer
}

func NewConsumer(feePerLeg, minProfit, minStart float64, out io.Writer) *Consumer {
	return &Consumer{
		FeePerLeg: feePerLeg,
		MinProfit: minProfit,
		MinStart:  minStart,
		writer:    out,
	}
}

// Start запускает горутину-потребителя.
func (c *Consumer) Start(
	ctx context.Context,
	events <-chan domain.Event,
	triangles []domain.Triangle,
	indexBySymbol map[string][]int,
	wg *sync.WaitGroup,
) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		c.run(ctx, events, triangles, indexBySymbol)
	}()
}

func (c *Consumer) run(
	ctx context.Context,
	events <-chan domain.Event,
	triangles []domain.Triangle,
	indexBySymbol map[string][]int,
) {
	quotes := make(map[string]domain.Quote)

	const minPrintInterval = 5 * time.Millisecond
	lastPrint := make(map[int]time.Time)

	for {
		select {
		case ev, ok := <-events:
			if !ok {
				return
			}

			if prev, okPrev := quotes[ev.Symbol]; okPrev &&
				prev.Bid == ev.Bid &&
				prev.Ask == ev.Ask &&
				prev.BidQty == ev.BidQty &&
				prev.AskQty == ev.AskQty {
				continue
			}

			quotes[ev.Symbol] = domain.Quote{
				Bid:    ev.Bid,
				Ask:    ev.Ask,
				BidQty: ev.BidQty,
				AskQty: ev.AskQty,
			}

			trIDs := indexBySymbol[ev.Symbol]
			if len(trIDs) == 0 {
				continue
			}

			now := time.Now()

			for _, id := range trIDs {
				tr := triangles[id]

				prof, ok := domain.EvalTriangle(tr, quotes, c.FeePerLeg)
				if !ok {
					continue
				}
				if prof < c.MinProfit {
					continue
				}

				ms, okMS := domain.ComputeMaxStartTopOfBook(tr, quotes, c.FeePerLeg)
				if okMS && c.MinStart > 0 && ms.MaxStart < c.MinStart {
					continue
				}

				if last, okLast := lastPrint[id]; okLast {
					if now.Sub(last) < minPrintInterval {
						continue
					}
				}
				lastPrint[id] = now

				var msPtr *domain.MaxStartInfo
				if okMS {
					msCopy := ms
					msPtr = &msCopy
				}
				c.printTriangle(now, tr, prof, quotes, msPtr)
			}

		case <-ctx.Done():
			return
		}
	}
}

func (c *Consumer) printTriangle(
	ts time.Time,
	t domain.Triangle,
	profit float64,
	quotes map[string]domain.Quote,
	ms *domain.MaxStartInfo,
) {
	w := c.writer
	fmt.Fprintf(w, "%s\n", ts.Format("2006-01-02 15:04:05.000"))

	if ms != nil {
		bneckSym := ""
		if ms.BottleneckLeg >= 0 && ms.BottleneckLeg < len(t.Legs) {
			bneckSym = t.Legs[ms.BottleneckLeg].Symbol
		}
		fmt.Fprintf(w, "[ARB] %+0.3f%%  %s  maxStart=%.4f %s  bottleneck=%s\n",
			profit*100, t.Name,
			ms.MaxStart, ms.StartAsset,
			bneckSym,
		)
	} else {
		fmt.Fprintf(w, "[ARB] %+0.3f%%  %s\n", profit*100, t.Name)
	}

	for _, leg := range t.Legs {
		q := quotes[leg.Symbol]
		mid := (q.Bid + q.Ask) / 2
		spreadAbs := q.Ask - q.Bid
		spreadPct := 0.0
		if mid > 0 {
			spreadPct = spreadAbs / mid * 100
		}
		side := ""
		if leg.Dir > 0 {
			side = fmt.Sprintf("%s/%s", leg.From, leg.To)
		} else {
			side = fmt.Sprintf("%s/%s", leg.To, leg.From)
		}
		fmt.Fprintf(w, "  %s (%s): bid=%.10f ask=%.10f  spread=%.10f (%.5f%%)  bidQty=%.4f askQty=%.4f\n",
			leg.Symbol, side,
			q.Bid, q.Ask,
			spreadAbs, spreadPct,
			q.BidQty, q.AskQty,
		)
	}
	fmt.Fprintln(w)
}

func OpenLogWriter(path string) (io.WriteCloser, *bufio.Writer, io.Writer) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		log.Fatalf("open %s: %v", path, err)
	}
	buf := bufio.NewWriter(f)
	out := io.MultiWriter(os.Stdout, buf)
	return f, buf, out
}

4) cmd/cryptarb/main.go

(единственное — передаём cfg.MinStart в consumer)

package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"crypt_proto/arb"
	"crypt_proto/config"
	"crypt_proto/domain"
	"crypt_proto/exchange"
	"crypt_proto/kucoin"
	"crypt_proto/mexc"

	_ "net/http/pprof"
)

func main() {
	// pprof
	go func() {
		log.Println("pprof on http://localhost:6060/debug/pprof/")
		if err := http.ListenAndServe("localhost:6060", nil); err != nil {
			log.Printf("pprof server error: %v", err)
		}
	}()

	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	cfg := config.Load()

	triangles, symbols, indexBySymbol, err := domain.LoadTriangles(cfg.TrianglesFile)
	if err != nil {
		log.Fatalf("load triangles: %v", err)
	}
	if len(triangles) == 0 {
		log.Fatal("нет треугольников, нечего мониторить")
	}
	if len(symbols) == 0 {
		log.Fatal("нет символов для подписки")
	}
	log.Printf("символов для подписки всего: %d", len(symbols))

	var feed exchange.MarketDataFeed

	switch cfg.Exchange {
	case "MEXC":
		feed = mexc.NewFeed(cfg.Debug)
	case "KUCOIN":
		feed = kucoin.NewFeed(cfg.Debug)
	default:
		log.Fatalf("unknown EXCHANGE=%q (expected MEXC or KUCOIN)", cfg.Exchange)
	}

	log.Printf("Using exchange: %s", feed.Name())

	// лог-файл для арбитража
	logFile, logBuf, arbOut := arb.OpenLogWriter("arbitrage.log")
	defer logFile.Close()
	defer logBuf.Flush()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	events := make(chan domain.Event, 8192)

	var wg sync.WaitGroup

	// запускаем потребителя
	consumer := arb.NewConsumer(cfg.FeePerLeg, cfg.MinProfit, cfg.MinStart, arbOut)
	consumer.Start(ctx, events, triangles, indexBySymbol, &wg)

	// запускаем фид биржи
	feed.Start(ctx, &wg, symbols, cfg.BookInterval, events)

	// ждём сигнал
	<-ctx.Done()
	log.Println("shutting down...")

	time.Sleep(200 * time.Millisecond)
	close(events)
	wg.Wait()
	log.Println("bye")
}

5) cmd/cryptarb/.env

(добавляем параметр минимального старта)

EXCHANGE=MEXC
TRIANGLES_FILE=triangles_markets.csv
BOOK_INTERVAL=10ms
FEE_PCT=0.04
MIN_PROFIT_PCT=0.5

# Минимальный стартовый объём (в валюте начала треугольника, обычно USDT).
# 0 = фильтр отключен. Пример: 0.5
MIN_START_USDT=0

DEBUG=false

Как пользоваться

Хочешь торговать/видеть арбы “по полдоллара” — ставь:

MIN_START_USDT=0.5

Хочешь просто видеть всё (как раньше) — оставь 0.

Если скажешь, как у тебя сейчас определяется “валюта старта” для треугольника (всегда USDT или бывает другое), я допилю, чтобы MIN_START_USDT работал корректно даже если старт не USDT (через автоконвертацию).

как расчитывается MIN_START_USDT  

MIN_START_USDT не “рассчитывается” из чего-то сам по себе — это порог, с которым мы сравниваем maxStart, вычисленный по стаканам.

Что такое maxStart

maxStart — это максимальная сумма стартовой валюты треугольника, которую можно прогнать по первой строке стакана (best bid/ask) на всех 3 ногах так, чтобы ни на одной ноге не превысить доступный объём (bidQty/askQty).

Если maxStart < MIN_START_USDT, то мы такой треугольник не выводим/не торгуем.

Как считается maxStart (формула)

Мы идём по ногам 1→2→3 и держим коэффициент kIn:

kIn[i] = сколько входной валюты ноги i получится из 1 единицы стартовой валюты (с учётом предыдущих конверсий и комиссий).

Для каждой ноги считаем ограничение на вход limitIn:

Если нога SELL (base → quote), Dir > 0

исполняемся по bid

вход ноги = base

максимум на вход:
limitIn = bidQty (сколько base можно продать по лучшему bid)

коэффициент конверсии:
ratio = bid

Если нога BUY (quote → base), Dir < 0

исполняемся по ask

вход ноги = quote

максимум на вход:
limitIn = askQty * ask (сколько quote можно потратить на лучший ask)

коэффициент конверсии:
ratio = 1 / ask

Дальше получаем ограничение на старт по этой ноге:

maxStartByLeg = limitIn / kIn

И общий максимум старта:

maxStart = min(maxStartByLeg[0], maxStartByLeg[1], maxStartByLeg[2])

Обновление kIn после ноги (учитываем комиссию как удержание из результата ноги):

kIn = kIn * ratio * (1 - feeRate)

Важно про название MIN_START_USDT

В моём коде maxStart считается в стартовой валюте треугольника:

StartAsset = t.Legs[0].From

Поэтому MIN_START_USDT корректен только если стартовая валюта реально USDT (например, все треугольники у тебя начинаются с USDT).
Если бывают треугольники, где старт не USDT — этот порог будет интерпретироваться в другой валюте, и это уже неверно по смыслу.

Если хочешь, я сделаю правильно: MIN_START_USDT всегда в USDT, а maxStart будет конвертироваться в USDT по текущей книге (или будем принудительно нормализовать треугольники так, чтобы старт всегда USDT).

maxStart как считается?

maxStart — это максимальная сумма стартовой валюты, которую можно прогнать через 3 ноги по лучшим bid/ask, не превысив доступный объём (bidQty/askQty) ни на одной ноге.

Он считается как минимум из трёх ограничений (по каждой ноге отдельно).

Идея на пальцах

Ты стартуешь, например, с S USDT.

На 1-й ноге ты либо покупаешь (берёшь ask и ограничен askQty), либо продаёшь (берёшь bid и ограничен bidQty).

После 1-й ноги из S получается уже другая валюта (и меньше из-за комиссии).

На 2-й ноге вход — уже эти деньги. Там тоже есть лимит по top-of-book.

maxStart — такой S, чтобы на каждой ноге вход не превышал лимит стакана.

Формулы

Держим коэффициент kIn — сколько “входной валюты текущей ноги” получается из 1 единицы стартовой валюты.

Идём по ногам i=0..2:

1) Лимит входа ноги limitIn[i]

Если нога SELL (base → quote) (в моём коде это Dir > 0):

исполняемся по bid

вход ноги = base

максимум входа:
limitIn = bidQty

Если нога BUY (quote → base) (Dir < 0):

исполняемся по ask

вход ноги = quote

максимум входа:
limitIn = askQty * ask
(потому что покупка base по askQty требует quote = askQty×ask)

2) Ограничение на старт по этой ноге

Если стартовая сумма = S, то вход на ногу i равен S * kIn.

Чтобы не превысить стакан:
S * kIn <= limitIn

Отсюда:
maxStartByLeg[i] = limitIn / kIn

3) Общий maxStart

maxStart = min(maxStartByLeg[0], maxStartByLeg[1], maxStartByLeg[2])

4) Как обновляется kIn

Нужен коэффициент конверсии ratio (out/in) по лучшей цене:

SELL: ratio = bid

BUY: ratio = 1/ask

Комиссия на ногу feeRate (у тебя 0.0004 для 0.04%):
kIn = kIn * ratio * (1 - feeRate)

Мини-пример (приближённо, как в твоём логе)

Треугольник: USDT -> ZEC -> USDC -> USDT

BUY ZECUSDT (покупаем ZEC за USDT):

limitIn1 = askQty * ask (в USDT)

SELL ZECUSDC (продаём ZEC за USDC):

вход = ZEC

limitIn2 = bidQty (в ZEC)

SELL USDCUSDT:

вход = USDC

limitIn3 = bidQty (в USDC)

Если на 2-й ноге bidQty = 0.003 ZEC, то именно она часто и будет bottleneck: она задаёт маленький maxStart.

Если хочешь, я могу на твоём конкретном одном [ARB] из лога (вставь 3 строки с bid/ask/qty) посчитать maxStart руками и показать, какая нога ограничивает и почему.

но maxstart лучше брать половину от минимального объема

Да, это здраво. maxStart считается по первой строке стакана, а она “хрупкая”: объём может исчезнуть, цена сдвинуться на тик, плюс будут округления stepSize. Поэтому в реальной торговле обычно берут коэффициент безопасности.

Как правильно

Сначала считаем maxStart (как мы сделали).

Берём рабочий старт:

start = maxStart * SAFETY_FACTOR

Если ты хочешь “половину”, то:

SAFETY_FACTOR = 0.5

Почему половина лучше, чем 100%

меньше шанс, что best-объём “съедят” до твоего ордера

меньше проскальзывание

меньше риск, что из-за округления 2–3 нога не закроется

меньше “хвостов” в промежуточной валюте

Практическая схема

Обычно делают так:

start = min(userTarget, balanceStart, maxStart * SAFETY_FACTOR)

И фильтр:

maxStart * SAFETY_FACTOR >= MIN_START_USDT

(то есть сравнивать лучше уже “безопасный” старт, а не теоретический максимум)

Какой фактор ставить

0.5 — хороший дефолт.

0.3 — если пары очень тонкие/скачут.

0.7–0.8 — если пары ликвидные и ты видишь стабильные книги.

Если хочешь, я внесу это прямо в правки, которые я тебе давал:

добавлю env START_FRACTION=0.5

и в расчёте буду использовать safeStart = ms.MaxStart * START_FRACTION (и его логировать).







