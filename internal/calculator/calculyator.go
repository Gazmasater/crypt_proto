package calculator

import (
	"io"
	"log"
	"os"

	"crypt_proto/internal/queue"
	"crypt_proto/pkg/models"
)

type Calculator struct {
	mem     *queue.MemoryStore
	scanner *Scanner
	filter  *ExecutorFilter
	oppLog  *log.Logger
}

func NewCalculator(mem *queue.MemoryStore, triangles []*Triangle) *Calculator {
	f, err := os.OpenFile("arb_opportunities.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		log.Fatalf("failed to open log: %v", err)
	}

	mw := io.MultiWriter(os.Stdout, f)
	return &Calculator{
		mem:     mem,
		scanner: NewScanner(mem, triangles),
		filter:  NewExecutorFilter(),
		oppLog:  log.New(mw, "", log.LstdFlags),
	}
}

func (c *Calculator) Run(in <-chan *models.MarketData) {
	for md := range in {
		c.mem.Push(md)

		cands := c.scanner.CandidatesFor(md.Symbol, md.Timestamp)
		if len(cands) == 0 {
			continue
		}

		for _, cand := range cands {
			opp, ok := c.filter.Evaluate(cand)
			if !ok {
				continue
			}
			c.logOpportunity(opp)
		}
	}
}

func (c *Calculator) logOpportunity(opp ExecutableOpportunity) {
	c.oppLog.Printf(
		"[ARB] %s→%s→%s | est=%.4f%% | real=%.4f%% | start=%.2f USDT | minStart=%.2f USDT | final=%.4f USDT | profit=%.4f USDT | symbol=%s",
		opp.Triangle.A,
		opp.Triangle.B,
		opp.Triangle.C,
		opp.EstimatedPct*100,
		opp.ProfitPct*100,
		opp.StartUSDT,
		opp.MinStartUSDT,
		opp.FinalUSDT,
		opp.ProfitUSDT,
		opp.TriggeredBy,
	)
}
