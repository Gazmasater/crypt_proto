mx0vglmT3srN1IS19H
135bb7a7509e4421bad692415c53753b



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





1. Добавляем ключи в Config (cmd/cryptarb/config.go)

Расширяем структуру:

type Config struct {
	Exchange      string
	TrianglesFile string
	BookInterval  time.Duration

	FeePerLeg     float64
	MinProfit     float64
	MinStartUSDT  float64
	StartFraction float64

	Debug bool

	// ДОБАВЛЕНО:
	APIKey    string
	APISecret string
}


И в LoadConfig() добавляем чтение ключей.

Я сделаю так:

сначала читаем EXCHANGE (MEXC/OKX/KUCOIN);

пытаемся взять <EXCHANGE>_API_KEY и <EXCHANGE>_API_SECRET;

если их нет — падаем назад на API_KEY / API_SECRET.

func LoadConfig() Config {
	exch := getenvDefault("EXCHANGE", "MEXC")

	apiKey := getenvDefault(exch+"_API_KEY", "")
	if apiKey == "" {
		apiKey = getenvDefault("API_KEY", "")
	}

	apiSecret := getenvDefault(exch+"_API_SECRET", "")
	if apiSecret == "" {
		apiSecret = getenvDefault("API_SECRET", "")
	}

	return Config{
		Exchange:      exch,
		TrianglesFile: getenvDefault("TRIANGLES_FILE", "triangles_markets.csv"),
		BookInterval:  parseDurationEnv("BOOK_INTERVAL", 10*time.Millisecond),

		FeePerLeg:     parseFloatEnvPercent("FEE_PCT", 0.0004),
		MinProfit:     parseFloatEnvPercent("MIN_PROFIT_PCT", 0.001),
		MinStartUSDT:  parseFloatEnv("MIN_START_USDT", 2),
		StartFraction: parseFloatEnv("START_FRACTION", 0.5),

		Debug: parseBoolEnv("DEBUG", false),

		APIKey:    apiKey,
		APISecret: apiSecret,
	}
}


Теперь ты можешь в .env записать, например для MEXC:

EXCHANGE=MEXC
MEXC_API_KEY=xxx
MEXC_API_SECRET=yyy


или просто:

EXCHANGE=MEXC
API_KEY=xxx
API_SECRET=yyy


— оба варианта заработают.

2. Подцепляем ключи в main.go

В cmd/cryptarb/main.go после cfg := LoadConfig() добавь проверку и заготовку для трейдера/исполнителя.

Пример:

func main() {
	cfg := LoadConfig()
	log.Printf("config: %+v", cfg)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	triangles, symbols, indexBySymbol, err := domain.LoadTriangles(cfg.TrianglesFile)
	if err != nil {
		log.Fatalf("load triangles: %v", err)
	}

	logFile, logBuf, out := arb.OpenLogWriter("arbitrage.log")
	defer logFile.Close()
	defer logBuf.Flush()

	consumer := arb.NewConsumer(cfg.FeePerLeg, cfg.MinProfit, cfg.MinStartUSDT, out)
	consumer.StartFraction = cfg.StartFraction

	// === ТУТ УЧЁТ API КЛЮЧЕЙ ===

	if cfg.APIKey == "" || cfg.APISecret == "" {
		log.Printf("[WARN] API_KEY/API_SECRET не заданы — работаем только в режиме логирования (без реальной торговли)")
		// consumer.Executor = arb.NewDryRunExecutor(...) // если сделаешь dry-run
	} else {
		log.Printf("[INFO] API-ключи для %s загружены, можно подключать торгового исполнителя", cfg.Exchange)

		// здесь, когда напишешь трейдер, будет что-то вроде:
		//
		// var exec arb.TriangleExecutor
		// switch cfg.Exchange {
		// case "MEXC":
		//     trader := mexc.NewTrader(cfg.APIKey, cfg.APISecret, cfg.Debug)
		//     exec = arb.NewTriangleExecutor(trader, cfg.FeePerLeg, cfg.MinProfit, cfg.MinStartUSDT)
		// case "OKX":
		//     ...
		// }
		// consumer.Executor = exec
	}

	events := make(chan domain.Event, 1024)

	var wg sync.WaitGroup
	consumer.Start(ctx, events, triangles, indexBySymbol, &wg)

	// тут твой фид по стаканам (MEXC/OKX/KuCoin) пишет в events...

	wg.Wait()
}


👉 Ключевая мысль:
Ключи живут только в Config и дальше передаются в “трейдер” (обертка над API биржи).
Ни domain, ни arb.Consumer про них знать не должны — они бирже-независимые.

3. Где будут реально использоваться ключи

Сейчас мы только:

читаем ключи из ENV,

прокидываем их до main.go,

показываем, куда их передать.

Реальное использование будет в специфичном адаптере, например mexc:

// pseudo: crypt_proto/exchange/mexc/trader.go

type Trader struct {
	apiKey    string
	apiSecret string
	debug     bool
	// httpClient и т.п.
}

func NewTrader(apiKey, apiSecret string, debug bool) *Trader {
	return &Trader{
		apiKey:    apiKey,
		apiSecret: apiSecret,
		debug:     debug,
	}
}

func (t *Trader) PlaceMarket(ctx context.Context, symbol string, side arb.OrderSide, amountBase float64) (arb.OrderResult, error) {
	// здесь используешь apiKey/apiSecret для подписи запросов к MEXC
}


И этот Trader потом оборачиваешь в твой TriangleExecutor, который уже знает, как идти по трём ногам треугольника.

4. Коротко по сути вопроса

а где апи ключ и секретный ключ добавлен

До этого момента — нигде: мы занимались только аналитикой и логами.
Сейчас:

в Config добавили APIKey и APISecret;

читаем их из ENV (MEXC_API_KEY / MEXC_API_SECRET или API_KEY / API_SECRET);

в main.go они доступны, и туда подключается твой будущий трейдер.

Если хочешь, следующим шагом могу прямо набросать:

минимальный MexcTrader с подписями запросов (timestamp, sign),

и DryRunTriangleExecutor, чтобы ты мог проверить интеграцию без риска денег.

