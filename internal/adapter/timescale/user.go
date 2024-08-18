package timescale

import (
	"context"
	"errors"
	"iter"

	"github.com/dagulv/stock-api/internal/core/domain"
	"github.com/dagulv/stock-api/internal/core/port"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/xid"
)

type userStore struct {
	db *pgxpool.Pool
}

func NewUser(db *pgxpool.Pool) port.User {
	return userStore{
		db: db,
	}
}

func (s userStore) List(ctx context.Context) (_ iter.Seq[domain.User], err error) {
	rows, err := s.db.Query(
		ctx,
		`SELECT
			"id",
			"tenantId",
			"firstName",
			"lastName",
			"email",
			"active",
			"timeCreated",
			"timeUpdated"
		FROM "user"`,
	)

	if err != nil {
		return
	}

	return Iter(rows, func(r pgx.Rows) (user domain.User, err error) {
		if err = rows.Scan(&user.Id, &user.TenantId, &user.FirstName, &user.LastName, &user.Email, &user.Active, &user.TimeCreated, &user.TimeUpdated); err != nil {
			return
		}

		return
	}), nil
}

func (s userStore) Get(ctx context.Context, userId xid.ID, dst *domain.User) (err error) {
	row := s.db.QueryRow(
		ctx,
		`SELECT
			"user"."id"
		FROM "user"`,
	)

	if err = row.Scan(&dst.Id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			//return no rows error
			return
		}

		return
	}

	return
}

func (s userStore) Create(ctx context.Context, user *domain.User) (err error) {
	_, err = s.db.Exec(
		ctx,
		`INSERT INTO "user" (
			"id",
			"tenantId",
			"firstName",
			"lastName",
			"email",
			"active",
			"timeCreated",
			"timeUpdated"
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		user.Id,
		user.TenantId,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Active,
		user.TimeCreated,
		user.TimeUpdated,
	)

	return
}
