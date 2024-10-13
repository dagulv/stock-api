package service

import (
	"context"
	"iter"

	"github.com/dagulv/stock-api/internal/core/domain"
	"github.com/dagulv/stock-api/internal/core/port"
)

type Stock struct {
	Store port.Stock
}

func (s Stock) List(ctx context.Context) (_ iter.Seq[domain.Stock], err error) {
	return s.Store.List(ctx)
}

func (s Stock) Get(ctx context.Context, symbol string, dst *domain.Stock) (err error) {
	return s.Store.Get(ctx, symbol, dst)
}
