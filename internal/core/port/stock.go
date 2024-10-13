package port

import (
	"context"
	"iter"

	"github.com/dagulv/stock-api/internal/core/domain"
)

type Stock interface {
	List(ctx context.Context) (iter.Seq[domain.Stock], error)
	Get(ctx context.Context, symbol string, dst *domain.Stock) error
}
