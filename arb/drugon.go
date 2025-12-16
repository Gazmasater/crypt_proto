package arb

import (
	"context"
	"fmt"
	"io"

	"crypt_proto/domain"
)

type DryRunExecutor struct {
	out io.Writer
}

func NewDryRunExecutor(out io.Writer) *DryRunExecutor {
	return &DryRunExecutor{out: out}
}

func (e *DryRunExecutor) Name() string { return "DRY-RUN" }

func (e *DryRunExecutor) Execute(ctx context.Context, t domain.Triangle, quotes map[string]domain.Quote, startUSDT float64) error {
	fmt.Fprintf(e.out, "  [DRY RUN] start=%.6f USDT triangle=%s\n", startUSDT, t.Name)
	return nil
}
