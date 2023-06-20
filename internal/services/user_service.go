package services

import (
	"context"
	"time"

	"github.com/dagulv/stock-api/internal/models"
	"github.com/dagulv/stock-api/internal/stores"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/xid"
)

type UserService struct {
	store stores.UserStore
}

func (s UserService) Post(ctx context.Context, user *models.User) (err error) {
	user.Id = xid.NewWithTime(time.Now())
	user.TimeCreated = pgtype.Timestamptz{Time: time.Now()}
	user.TimeUpdated = user.TimeCreated

}
