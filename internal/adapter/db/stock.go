package db

import (
	"context"
	"errors"
	"iter"

	"github.com/dagulv/stock-api/internal/core/domain"
	"github.com/dagulv/stock-api/internal/core/port"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type stockStore struct {
	db *pgxpool.Pool
}

func NewStock(db *pgxpool.Pool) port.Stock {
	return stockStore{
		db: db,
	}
}

func (s stockStore) List(ctx context.Context) (_ iter.Seq[domain.Stock], err error) {
	rows, err := s.db.Query(
		ctx,
		`SELECT
			"companies"."symbol",
			"companies"."name",
			MAX("ohlcv"."close")
		FROM "companies"
		INNER JOIN "ohlcv" ON "ohlcv"."companyId" = "companies"."id"
		GROUP BY "companies"."id";`,
	)

	if err != nil {
		return
	}

	return Iter(rows, func(r pgx.Rows) (stock domain.Stock, err error) {
		if err = rows.Scan(&stock.Symbol, &stock.Name, &stock.Price); err != nil {
			return
		}

		return
	}), nil
}

func (s stockStore) Get(ctx context.Context, symbol string, dst *domain.Stock) (err error) {
	row := s.db.QueryRow(
		ctx,
		`SELECT
			"companies"."symbol",
			"companies"."name",
			MAX("ohlcv"."close")
		FROM "companies"
		INNER JOIN "ohlcv" ON "ohlcv"."companyId" = "companies"."id"
		WHERE "companies"."symbol" = $1
		GROUP BY "companies"."id";`, symbol,
	)

	if err = row.Scan(&dst.Symbol, &dst.Name, &dst.Price); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			//return no rows error
			return
		}

		return
	}

	return
}
