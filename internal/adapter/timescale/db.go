package timescale

import (
	"context"
	"database/sql"

	"github.com/dagulv/stock-api/internal/env"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func Connect(ctx context.Context, env env.Env) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, env.DatabaseUrl)

	if err != nil {
		return nil, err
	}

	return pool, nil
}

func Open(ctx context.Context, env env.Env) (*sql.DB, error) {
	return sql.Open("pgx", env.DatabaseUrl)
}
