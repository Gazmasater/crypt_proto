package arb

import (
	"context"

	"crypt_proto/domain"
)

// NoopExecutor — безопасный исполнитель "ничего не делаю".
// Используется когда трейд выключен или нет ключей API.
type NoopExecutor struct{}

func NewNoopExecutor() *NoopExecutor { return &NoopExecutor{} }

func (e *NoopExecutor) Name() string { return "NOOP" }

func (e *NoopExecutor) Execute(ctx context.Context, t domain.Triangle, quotes map[string]domain.Quote, startUSDT float64) error {
	// намеренно ничего не делаем
	_ = ctx
	_ = t
	_ = quotes
	_ = startUSDT
	return nil
}
