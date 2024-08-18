package port

import (
	"context"
	"iter"

	"github.com/dagulv/stock-api/internal/core/domain"
	"github.com/rs/xid"
)

type User interface {
	List(ctx context.Context) (iter.Seq[domain.User], error)
	Get(ctx context.Context, userId xid.ID, dst *domain.User) error
	Create(ctx context.Context, user *domain.User) error
}
