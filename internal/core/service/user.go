package service

import (
	"context"
	"iter"
	"time"

	"github.com/dagulv/stock-api/internal/core/domain"
	"github.com/dagulv/stock-api/internal/core/port"
	"github.com/rs/xid"
)

type User struct {
	Store port.User
}

func (s User) List(ctx context.Context) (_ iter.Seq[domain.User], err error) {
	return s.Store.List(ctx)
}

func (s User) Get(ctx context.Context, userId xid.ID, dst *domain.User) (err error) {
	return s.Store.Get(ctx, userId, dst)
}

func (s User) Create(ctx context.Context, user *domain.User) (err error) {
	user.Id = xid.NewWithTime(time.Now())
	user.TimeCreated = time.Now()
	user.TimeUpdated = user.TimeCreated

	return s.Store.Create(ctx, user)
}
