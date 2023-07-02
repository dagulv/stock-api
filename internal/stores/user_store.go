package stores

import (
	"context"

	"github.com/dagulv/stock-api/internal/models"
	"github.com/rs/xid"
)

type UserStore interface {
	Create(ctx context.Context, user *models.User) error
	SetPassword(ctx context.Context, passwordHash string, userId xid.ID) error
	GetByEmail(ctx context.Context, email string, dst *models.User) error
	List(ctx context.Context, set func(*models.User) error) error
	Get(ctx context.Context, userId xid.ID, dst *models.User) error
}
