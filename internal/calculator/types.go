package calculator

import (
	"fmt"
	"math"
	"strings"

	"crypt_proto/internal/queue"
)

const (
	defaultTakerFee      = 0.001
	defaultMinVolumeUSDT = 50.0
	defaultMinProfitPct  = 0.001
	defaultSearchStep    = 0.01
)

var triangleLegColumns = [3]int{3, 7, 11}

type LogMode int

const (
	LogSilent LogMode = iota
	LogNormal
	LogDebug
)

type Config struct {
	MinVolumeUSDT  float64
	MinProfitPct   float64
	SearchStepUSDT float64
	QuoteAgeMaxMS  int64
	StatsEverySec  int
	LogMode        LogMode
}

func DefaultConfig() Config {
	return Config{
		MinVolumeUSDT:  defaultMinVolumeUSDT,
		MinProfitPct:   defaultMinProfitPct,
		SearchStepUSDT: defaultSearchStep,
		QuoteAgeMaxMS:  2500,
		StatsEverySec:  5,
		LogMode:        LogNormal,
	}
}

type LegIndex struct {
	Key    string
	Symbol string
	IsBuy  bool
}

type LegRules struct {
	Symbol      string
	Side        string
	Base        string
	Quote       string
	QtyStep     float64
	QuoteStep   float64
	PriceStep   float64
	MinQty      float64
	MinQuote    float64
	MinNotional float64
	Fee         float64
}

type Triangle struct {
	A, B, C string
	Legs    [3]LegIndex
	Rules   [3]LegRules
}

type ScanCandidate struct {
	Triangle      *Triangle
	Quotes        [3]queue.Quote
	EstimatedPct  float64
	MaxStartUSDT  float64
	TriggeredBy   string
	TriggeredAtMS int64
}

type ExecutionResult struct {
	StartUSDT   float64
	MinStart    float64
	FinalUSDT   float64
	ProfitUSDT  float64
	ProfitPct   float64
	LegNotional [3]float64
	LegAmount   [3]float64
}

type ExecutableOpportunity struct {
	Triangle         *Triangle
	Quotes           [3]queue.Quote
	EstimatedPct     float64
	StartUSDT        float64
	MinStartUSDT     float64
	FinalUSDT        float64
	ProfitUSDT       float64
	ProfitPct        float64
	IdealFinalUSDT   float64
	IdealProfitPct   float64
	RoundedFinalUSDT float64
	RoundedProfitPct float64
	TriggeredBy      string
	TriggeredAtMS    int64
}

type ScanResult struct {
	Candidate ScanCandidate
	Reject    string
	OK        bool
}

type Stats struct {
	Ticks         int64
	TrianglesSeen int64
	Candidates    int64
	Opportunities int64

	Positive int64
	Negative int64
	Logged   int64

	ScanRejects map[string]int64
	ExecRejects map[string]int64
}

func feeMultiplier(fee float64) float64 {
	if fee > 0 && fee < 1 {
		return 1 - fee
	}
	return 1 - defaultTakerFee
}

func applyFloorStep(value, step float64) float64 {
	if value <= 0 {
		return 0
	}
	if step <= 0 {
		return value
	}
	return floorToStep(value, step)
}

func floorToStep(value, step float64) float64 {
	if step <= 0 {
		return value
	}
	units := math.Floor((value + 1e-12) / step)
	if units <= 0 {
		return 0
	}
	result := units * step
	precision := decimalsFromStep(step)
	pow := math.Pow10(precision)
	return math.Floor(result*pow+1e-9) / pow
}

func decimalsFromStep(step float64) int {
	if step <= 0 {
		return 8
	}
	s := strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.12f", step), "0"), ".")
	idx := strings.IndexByte(s, '.')
	if idx == -1 {
		return 0
	}
	return len(s) - idx - 1
}

func passesMinChecks(qty, notional float64, rules LegRules) bool {
	if rules.MinQty > 0 && qty+1e-12 < rules.MinQty {
		return false
	}
	if rules.MinQuote > 0 && notional+1e-12 < rules.MinQuote {
		return false
	}
	if rules.MinNotional > 0 && notional+1e-12 < rules.MinNotional {
		return false
	}
	return true
}
