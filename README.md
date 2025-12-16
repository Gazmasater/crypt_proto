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




1) Добавь файл mexc/rules.go
package mexc

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type SymbolRules struct {
	Symbol string

	IsSpotTradingAllowed      bool
	QuoteOrderQtyMarketAllowed bool

	// baseSizePrecision приходит строкой "0.0001" — это step для quantity
	BaseStepStr string
	BaseStep    float64
	QtyDecimals int

	// точность для quoteOrderQty (amount)
	QuoteAssetPrecision int
	QuotePrecision      int

	// минималки (если есть)
	MinQty          float64
	MinNotional     float64
	MinOrderAmount  float64 // quoteAmountPrecision (по сути min amount), если пригодится
}

type exchangeInfoResp struct {
	Symbols []struct {
		Symbol string `json:"symbol"`

		Status string `json:"status"`

		IsSpotTradingAllowed       bool `json:"isSpotTradingAllowed"`
		QuoteOrderQtyMarketAllowed bool `json:"quoteOrderQtyMarketAllowed"`

		BaseSizePrecision string `json:"baseSizePrecision"`

		BaseAssetPrecision  int `json:"baseAssetPrecision"`
		QuoteAssetPrecision int `json:"quoteAssetPrecision"`
		QuotePrecision      int `json:"quotePrecision"`

		// Иногда присутствуют (зависит от версии ответа)
		MinQty             string `json:"minQty"`
		MinNotional        string `json:"minNotional"`
		QuoteAmountPrecision string `json:"quoteAmountPrecision"`
	} `json:"symbols"`
}

func LoadSymbolRules(ctx context.Context, baseURL string, client *http.Client) (map[string]SymbolRules, error) {
	if baseURL == "" {
		baseURL = "https://api.mexc.com"
	}
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+"/api/v3/exchangeInfo", nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, _ := io.ReadAll(resp.Body)
	if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("exchangeInfo error: status=%d body=%s", resp.StatusCode, string(b))
	}

	var root exchangeInfoResp
	if err := json.Unmarshal(b, &root); err != nil {
		return nil, err
	}

	out := make(map[string]SymbolRules, len(root.Symbols))
	for _, s := range root.Symbols {
		sym := strings.TrimSpace(s.Symbol)
		if sym == "" {
			continue
		}

		step := parseStep(s.BaseSizePrecision)
		dec := decimalsFromStepStr(s.BaseSizePrecision)

		r := SymbolRules{
			Symbol:                    sym,
			IsSpotTradingAllowed:      s.IsSpotTradingAllowed,
			QuoteOrderQtyMarketAllowed: s.QuoteOrderQtyMarketAllowed,
			BaseStepStr:               s.BaseSizePrecision,
			BaseStep:                  step,
			QtyDecimals:               dec,
			QuoteAssetPrecision:       s.QuoteAssetPrecision,
			QuotePrecision:            s.QuotePrecision,
			MinQty:                    parseFloatSafe(s.MinQty),
			MinNotional:               parseFloatSafe(s.MinNotional),
			MinOrderAmount:            parseFloatSafe(s.QuoteAmountPrecision),
		}

		out[sym] = r
	}

	return out, nil
}

func parseFloatSafe(s string) float64 {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return v
}

func parseStep(stepStr string) float64 {
	stepStr = strings.TrimSpace(stepStr)
	if stepStr == "" {
		return 0
	}
	v, err := strconv.ParseFloat(stepStr, 64)
	if err != nil {
		return 0
	}
	return v
}

func decimalsFromStepStr(step string) int {
	step = strings.TrimSpace(step)
	if step == "" || step == "1" {
		return 0
	}
	if i := strings.IndexByte(step, '.'); i >= 0 {
		frac := step[i+1:]
		frac = strings.TrimRight(frac, "0")
		return len(frac)
	}
	return 0
}

func floorToStep(x, step float64) float64 {
	if step <= 0 {
		return x
	}
	return math.Floor(x/step) * step
}

func truncToDecimals(x float64, decimals int) float64 {
	if decimals <= 0 {
		return math.Floor(x)
	}
	p := math.Pow10(decimals)
	return math.Floor(x*p) / p
}


2) Перепиши mexc/trader.go (полностью)
Ключевое: трейдер хранит rules и умеет “умно” ставить market BUY/SELL.
package mexc

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Trader struct {
	apiKey    string
	apiSecret string
	debug     bool
	baseURL   string
	client    *http.Client

	rules map[string]SymbolRules
}

func NewTrader(apiKey, apiSecret string, debug bool) *Trader {
	return &Trader{
		apiKey:    strings.TrimSpace(apiKey),
		apiSecret: strings.TrimSpace(apiSecret),
		debug:     debug,
		baseURL:   "https://api.mexc.com",
		client:    &http.Client{Timeout: 10 * time.Second},
		rules:     map[string]SymbolRules{},
	}
}

func (t *Trader) SetRules(r map[string]SymbolRules) {
	if r == nil {
		t.rules = map[string]SymbolRules{}
		return
	}
	t.rules = r
}

func (t *Trader) Rules(symbol string) (SymbolRules, bool) {
	r, ok := t.rules[strings.TrimSpace(symbol)]
	return r, ok
}

// SmartMarketBuyUSDT:
// Покупка по USDT (quote), нормализует amount/qty по правилам symbol.
func (t *Trader) SmartMarketBuyUSDT(ctx context.Context, symbol string, usdt float64, ask float64) (string, error) {
	symbol = strings.TrimSpace(symbol)
	if usdt <= 0 {
		return "", fmt.Errorf("usdt<=0")
	}
	r, ok := t.Rules(symbol)
	if !ok {
		return "", fmt.Errorf("no rules for symbol=%s", symbol)
	}
	if !r.IsSpotTradingAllowed {
		return "", fmt.Errorf("symbol not allowed for spot/api: %s", symbol)
	}

	// если биржа разрешает quoteOrderQty — используем amount (и режем precision)
	if r.QuoteOrderQtyMarketAllowed {
		dec := r.QuoteAssetPrecision
		if dec <= 0 {
			// fallback
			dec = 2
		}
		amount := truncToDecimals(usdt, dec)
		if amount <= 0 {
			return "", fmt.Errorf("amount<=0 after trunc (dec=%d)", dec)
		}
		return t.placeMarket(ctx, symbol, "BUY", 0, amount)
	}

	// иначе — покупаем quantity по ask и режем step
	if ask <= 0 {
		return "", fmt.Errorf("ask<=0 for %s", symbol)
	}

	qtyRaw := usdt / ask
	qty := qtyRaw
	if r.BaseStep > 0 {
		qty = floorToStep(qtyRaw, r.BaseStep)
	} else if r.QtyDecimals >= 0 {
		qty = truncToDecimals(qtyRaw, r.QtyDecimals)
	}

	if r.MinQty > 0 && qty < r.MinQty {
		return "", fmt.Errorf("qty<minQty (qty=%.12f minQty=%.12f)", qty, r.MinQty)
	}
	if qty <= 0 {
		return "", fmt.Errorf("qty<=0 after normalize (raw=%.12f)", qtyRaw)
	}

	return t.placeMarket(ctx, symbol, "BUY", qty, 0)
}

// SmartMarketSellQty:
// Продажа base quantity, нормализует по step/precision.
func (t *Trader) SmartMarketSellQty(ctx context.Context, symbol string, qtyRaw float64) (string, error) {
	symbol = strings.TrimSpace(symbol)
	if qtyRaw <= 0 {
		return "", fmt.Errorf("qty<=0")
	}
	r, ok := t.Rules(symbol)
	if !ok {
		return "", fmt.Errorf("no rules for symbol=%s", symbol)
	}
	if !r.IsSpotTradingAllowed {
		return "", fmt.Errorf("symbol not allowed for spot/api: %s", symbol)
	}

	qty := qtyRaw
	if r.BaseStep > 0 {
		qty = floorToStep(qtyRaw, r.BaseStep)
	} else if r.QtyDecimals >= 0 {
		qty = truncToDecimals(qtyRaw, r.QtyDecimals)
	}

	if r.MinQty > 0 && qty < r.MinQty {
		return "", fmt.Errorf("qty<minQty (qty=%.12f minQty=%.12f)", qty, r.MinQty)
	}
	if qty <= 0 {
		return "", fmt.Errorf("qty<=0 after normalize (raw=%.12f)", qtyRaw)
	}

	return t.placeMarket(ctx, symbol, "SELL", qty, 0)
}

func (t *Trader) placeMarket(ctx context.Context, symbol, side string, quantity, quoteOrderQty float64) (string, error) {
	side = strings.ToUpper(strings.TrimSpace(side))

	params := url.Values{}
	params.Set("symbol", symbol)
	params.Set("side", side)
	params.Set("type", "MARKET")
	params.Set("timestamp", strconv.FormatInt(time.Now().UnixMilli(), 10))

	if side == "BUY" {
		if quoteOrderQty > 0 {
			params.Set("quoteOrderQty", stripZeros(quoteOrderQty))
		} else {
			params.Set("quantity", stripZeros(quantity))
		}
	} else {
		params.Set("quantity", stripZeros(quantity))
	}

	queryToSign := params.Encode()
	params.Set("signature", t.sign(queryToSign))

	reqURL := t.baseURL + "/api/v3/order" + "?" + params.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("X-MEXC-APIKEY", t.apiKey)

	resp, err := t.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	b, _ := io.ReadAll(resp.Body)
	if resp.StatusCode/100 != 2 {
		return "", fmt.Errorf("mexc order error: status=%d body=%s", resp.StatusCode, string(b))
	}

	var m map[string]any
	_ = json.Unmarshal(b, &m)
	if v, ok := m["orderId"]; ok {
		return fmt.Sprintf("%v", v), nil
	}
	if v, ok := m["orderIdStr"]; ok {
		return fmt.Sprintf("%v", v), nil
	}
	return "", nil
}

func (t *Trader) GetBalance(ctx context.Context, asset string) (float64, error) {
	asset = strings.ToUpper(strings.TrimSpace(asset))
	if asset == "" {
		return 0, fmt.Errorf("empty asset")
	}

	params := url.Values{}
	params.Set("timestamp", strconv.FormatInt(time.Now().UnixMilli(), 10))

	queryToSign := params.Encode()
	params.Set("signature", t.sign(queryToSign))

	reqURL := t.baseURL + "/api/v3/account" + "?" + params.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("X-MEXC-APIKEY", t.apiKey)

	resp, err := t.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	b, _ := io.ReadAll(resp.Body)
	if resp.StatusCode/100 != 2 {
		return 0, fmt.Errorf("mexc account error: status=%d body=%s", resp.StatusCode, string(b))
	}

	var root map[string]any
	if err := json.Unmarshal(b, &root); err != nil {
		return 0, err
	}

	balAny, _ := root["balances"].([]any)
	for _, it := range balAny {
		m, ok := it.(map[string]any)
		if !ok {
			continue
		}
		a, _ := m["asset"].(string)
		if strings.ToUpper(strings.TrimSpace(a)) != asset {
			continue
		}

		if s, ok := m["free"].(string); ok {
			v, _ := strconv.ParseFloat(s, 64)
			return v, nil
		}
		if f, ok := m["free"].(float64); ok {
			return f, nil
		}
		return 0, nil
	}

	return 0, nil
}

func (t *Trader) sign(query string) string {
	mac := hmac.New(sha256.New, []byte(t.apiSecret))
	_, _ = mac.Write([]byte(query))
	return hex.EncodeToString(mac.Sum(nil))
}

func stripZeros(v float64) string {
	// много знаков не надо, всё равно выше мы нормализуем
	s := strconv.FormatFloat(v, 'f', 12, 64)
	s = strings.TrimRight(s, "0")
	s = strings.TrimRight(s, ".")
	if s == "" {
		return "0"
	}
	return s
}


3) Перепиши arb/executor_real.go (полностью)
Это версия, которая:


BUY делает через SmartMarketBuyUSDT (amount/qty нормализуются)


после BUY берёт реальный баланс купленного актива и продаёт его на следующей ноге (чтобы не ловить Oversold)


фильтрует пары по rules (isSpotTradingAllowed)


package arb

import (
	"context"
	"fmt"
	"io"
	"math"
	"strings"
	"time"

	"crypt_proto/domain"
)

type SpotTrader interface {
	SmartMarketBuyUSDT(ctx context.Context, symbol string, usdt float64, ask float64) (string, error)
	SmartMarketSellQty(ctx context.Context, symbol string, qty float64) (string, error)
	GetBalance(ctx context.Context, asset string) (float64, error)
}

type RealExecutor struct {
	trader SpotTrader
	out    io.Writer

	// фиксированный старт в USDT (например 2)
	StartUSDT float64

	// safety чтобы не словить Oversold на SELL
	SellSafety float64

	// анти-спам, чтобы один и тот же треугольник не пытался исполняться 50 раз/сек
	Cooldown time.Duration
	lastExec map[string]time.Time
}

func NewRealExecutor(tr SpotTrader, out io.Writer, startUSDT float64) *RealExecutor {
	return &RealExecutor{
		trader:     tr,
		out:        out,
		StartUSDT:  startUSDT,
		SellSafety: 0.995,
		Cooldown:   500 * time.Millisecond,
		lastExec:   make(map[string]time.Time),
	}
}

func (e *RealExecutor) Name() string { return "REAL" }

func (e *RealExecutor) Execute(ctx context.Context, t domain.Triangle, quotes map[string]domain.Quote, startUSDT float64) error {
	start := e.StartUSDT
	if startUSDT > 0 {
		start = startUSDT
	}
	if start <= 0 {
		return fmt.Errorf("start<=0")
	}

	// cooldown по имени треугольника
	if e.Cooldown > 0 {
		if last, ok := e.lastExec[t.Name]; ok && time.Since(last) < e.Cooldown {
			return nil
		}
		e.lastExec[t.Name] = time.Now()
	}

	fmt.Fprintf(e.out, "  [REAL EXEC] start=%.6f USDT triangle=%s\n", start, t.Name)

	// ожидаем, что стартовая валюта реально USDT
	curAsset := "USDT"
	curAmount := start

	// LEG 1: BUY (USDT -> A)
	leg1 := t.Legs[0]
	q1, ok := quotes[leg1.Symbol]
	if !ok || q1.Ask <= 0 {
		return fmt.Errorf("no quote/ask for %s", leg1.Symbol)
	}

	if !legMatchesFlow(leg1, curAsset) {
		return fmt.Errorf("leg1 flow mismatch: have=%s leg=%s", curAsset, leg1.Symbol)
	}

	// покупаем на 2 USDT
	fmt.Fprintf(e.out, "    [REAL EXEC] leg 1: BUY %s by USDT=%.6f\n", leg1.Symbol, curAmount)
	_, err := e.trader.SmartMarketBuyUSDT(ctx, leg1.Symbol, curAmount, q1.Ask)
	if err != nil {
		return fmt.Errorf("leg1 error: %w", err)
	}

	// после покупки: узнаём что реально купили (баланс base актива leg1)
	nextAsset1 := leg1.To // USDT->A по идее
	if nextAsset1 == "USDT" {
		nextAsset1 = leg1.From
	}
	bal1, err := e.trader.GetBalance(ctx, nextAsset1)
	if err != nil {
		return fmt.Errorf("leg1 balance error: %w", err)
	}
	curAsset = nextAsset1
	curAmount = bal1
	if curAmount <= 0 {
		return fmt.Errorf("leg1 result balance=0 asset=%s", curAsset)
	}

	// LEG 2: SELL/BUY в зависимости от Dir, но мы делаем проще:
	// Мы всегда хотим перейти curAsset -> nextAsset по leg.Dir.
	leg2 := t.Legs[1]
	if !legMatchesFlow(leg2, curAsset) {
		return fmt.Errorf("leg2 flow mismatch: have=%s leg=%s", curAsset, leg2.Symbol)
	}

	// если leg.Dir>0: From->To это SELL base->quote
	// если leg.Dir<0: From->To это BUY base<-quote (то есть curAsset=quote), мы должны BUY base за quote qty
	// У нас в SpotTrader только:
	// - SmartMarketSellQty(symbol, qty)
	// - SmartMarketBuyUSDT(symbol, usdt, ask) (это только когда quote=USDT)
	//
	// Поэтому для универсальности:
	// - Если мы на leg2 должны ПРОДАТЬ текущий актив -> используем SELL qty.
	// - Если должны КУПИТЬ base за quote и quote != USDT — пока НЕ делаем “BUY quantity” в Executor,
	//   а делаем SELL на другом направлении через правильный symbol (в домене Dir уже отражает направление).
	//
	// Практически для твоих треугольников USDT→X→USDC→USDT:
	// leg2 обычно SELL XUSDC, т.е. SELL qty — подходит.

	fmt.Fprintf(e.out, "    [REAL EXEC] leg 2: SELL %s qty=%.12f\n", leg2.Symbol, curAmount)
	sell2 := curAmount * e.SellSafety
	if sell2 <= 0 {
		return fmt.Errorf("leg2 qty<=0 after safety")
	}
	_, err = e.trader.SmartMarketSellQty(ctx, leg2.Symbol, sell2)
	if err != nil {
		return fmt.Errorf("leg2 error: %w", err)
	}

	// после продажи: баланс следующего актива
	nextAsset2 := leg2.To
	if strings.ToUpper(nextAsset2) == strings.ToUpper(curAsset) {
		nextAsset2 = leg2.From
	}
	bal2, err := e.trader.GetBalance(ctx, nextAsset2)
	if err != nil {
		return fmt.Errorf("leg2 balance error: %w", err)
	}
	curAsset = nextAsset2
	curAmount = bal2
	if curAmount <= 0 {
		return fmt.Errorf("leg2 result balance=0 asset=%s", curAsset)
	}

	// LEG 3: обычно SELL USDCUSDT
	leg3 := t.Legs[2]
	if !legMatchesFlow(leg3, curAsset) {
		return fmt.Errorf("leg3 flow mismatch: have=%s leg=%s", curAsset, leg3.Symbol)
	}

	fmt.Fprintf(e.out, "    [REAL EXEC] leg 3: SELL %s qty=%.12f\n", leg3.Symbol, curAmount)
	sell3 := curAmount * e.SellSafety
	if sell3 <= 0 {
		return fmt.Errorf("leg3 qty<=0 after safety")
	}
	_, err = e.trader.SmartMarketSellQty(ctx, leg3.Symbol, sell3)
	if err != nil {
		return fmt.Errorf("leg3 error: %w", err)
	}

	// финальный USDT баланс (опционально)
	usdtBal, _ := e.trader.GetBalance(ctx, "USDT")
	fmt.Fprintf(e.out, "  [REAL EXEC] done triangle %s  USDT_balance=%.6f\n", t.Name, math.Max(0, usdtBal))
	return nil
}

func legMatchesFlow(leg domain.Leg, have string) bool {
	// leg.From/To — это уже “логическая” цепочка в Triangle
	return strings.EqualFold(leg.From, have) || strings.EqualFold(leg.To, have)
}


Важно: этот executor сейчас покрывает твой основной кейс USDT→X→USDC→USDT (где 2-я и 3-я ноги — SELL). Для экзотики вида USDT→BTC→ALT→USDT (где появляется BUY за BTC/EUR) — добавим SmartMarketBuyByQty позже.


4) Перепиши cmd/cryptarb/main.go (кусок с инициализацией трейдера/правил)
Добавляем загрузку rules и SetRules(), чтобы “scale invalid” исчезло до ордера.
// ... внутри main после cfg := config.Load()

// ... triangles load etc

// создаём фид как было

// лог-файл
logFile, logBuf, arbOut := arb.OpenLogWriter("arbitrage.log")
defer logFile.Close()
defer logBuf.Flush()

ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
defer cancel()

events := make(chan domain.Event, 8192)
var wg sync.WaitGroup

consumer := arb.NewConsumer(cfg.FeePerLeg, cfg.MinProfit, cfg.MinStart, arbOut)
consumer.StartFraction = cfg.StartFraction

// ---- ВОТ ТУТ: реальная торговля ----
if cfg.TradeEnabled && cfg.APIKey != "" && cfg.APISecret != "" && cfg.TradeAmountUSDT > 0 {
	tr := mexc.NewTrader(cfg.APIKey, cfg.APISecret, cfg.Debug)

	// грузим правила символов
	rules, err := mexc.LoadSymbolRules(ctx, "https://api.mexc.com", nil)
	if err != nil {
		log.Fatalf("load mexc exchangeInfo: %v", err)
	}
	tr.SetRules(rules)
	log.Printf("[MEXC] loaded rules: %d symbols", len(rules))

	consumer.Executor = arb.NewRealExecutor(tr, arbOut, cfg.TradeAmountUSDT)
	log.Printf("[TRADE] enabled: startUSDT=%.4f", cfg.TradeAmountUSDT)
} else {
	consumer.Executor = arb.NewDryRunExecutor(arbOut)
	log.Printf("[TRADE] disabled (dry-run)")
}

consumer.Start(ctx, events, triangles, indexBySymbol, &wg)
feed.Start(ctx, &wg, symbols, cfg.BookInterval, events)

<-ctx.Done()
log.Println("shutting down...")
time.Sleep(200 * time.Millisecond)
close(events)
wg.Wait()
log.Println("bye")


5) Конфиг / env (чтобы ты мог задавать 2 USDT нормально)
В config.Config добавь:
TradeEnabled    bool
TradeAmountUSDT float64

В config.Load() добавь:
tradeEnabled := strings.ToLower(os.Getenv("TRADE_ENABLED")) == "true"
tradeAmount := loadEnvFloat("TRADE_AMOUNT_USDT", 2)

cfg := Config{
   ...
   TradeEnabled: tradeEnabled,
   TradeAmountUSDT: tradeAmount,
}
log.Printf("Trade enabled: %v", cfg.TradeEnabled)
log.Printf("Trade amount USDT: %.4f", cfg.TradeAmountUSDT)

.env:
TRADE_ENABLED=true
TRADE_AMOUNT_USDT=2


Что это тебе даст сразу


уйдёт amount scale is invalid и quantity scale is invalid (в рамках того, что exchangeInfo даёт step/precision)


уйдёт большая часть Oversold (мы продаём по реальному балансу, а не по симуляции)


10007 станет редким/управляемым (пары будут отфильтрованы правилами в трейдере)



Если хочешь — следующим шагом я добавлю:


жёсткий фильтр треугольников: оставить только те, где все 3 symbol isSpotTradingAllowed=true и leg2/leg3 действительно SELL (чтобы executor был 100% безопасен)


один активный трейд за раз (mutex), чтобы не было параллельных исполнений на одном балансе.










