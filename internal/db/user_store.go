package db

import (
	"context"

	"github.com/dagulv/stock-api/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userStore struct {
	db *pgxpool.Pool
}

func (s *userStore) Post(ctx context.Context, user *models.User) (err error) {
	_, err = s.db.Exec(
		ctx,
		`INSERT INTO users`
	)
}