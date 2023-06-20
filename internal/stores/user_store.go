package stores

import (
	"context"

	"github.com/dagulv/stock-api/internal/models"
)

type UserStore interface {
	Post(ctx context.Context, user *models.User) error
}
