package stores

import "context"

type StockStore interface {
	Get(ctx context.Context) error
}
