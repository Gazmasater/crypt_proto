package arb

import (
	"context"

	"crypt_proto/domain"
)

type SymbolFilter struct {
	StepSize float64
	MinQty   float64
}

type Executor interface {
	Name() string
	Execute(ctx context.Context, t domain.Triangle, quotes map[string]domain.Quote, startUSDT float64) error
}
