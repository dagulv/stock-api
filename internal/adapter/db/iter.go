package db

import (
	"iter"

	"github.com/jackc/pgx/v5"
)

func Iter[T any](rows pgx.Rows, scan func(pgx.Rows) (T, error)) iter.Seq[T] {
	return func(yield func(T) bool) {
		defer rows.Close()

		for rows.Next() {
			v, err := scan(rows)

			if err != nil {
				return
			}

			if !yield(v) {
				return
			}
		}
	}
}
