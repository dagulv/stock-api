package db

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

func StockStore(db *pgxpool.Pool) *stockStore {
	s := &stockStore{
		db: db,
	}

	return s
}

type stockStore struct {
	db *pgxpool.Pool
}
