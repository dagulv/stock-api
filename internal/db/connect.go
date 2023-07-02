package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(ctx context.Context) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig("postgresql://admin:admin@localhost:5432/stock?sslmode=disable")

	if err != nil {
		return nil, err
	}

	return pgxpool.NewWithConfig(ctx, config)
}
