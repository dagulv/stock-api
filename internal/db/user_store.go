package db

import (
	"context"
	"time"

	"github.com/dagulv/stock-api/internal/models"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/xid"
)

func UserStore(db *pgxpool.Pool) *userStore {
	s := &userStore{
		db: db,
	}

	return s
}

type userStore struct {
	db *pgxpool.Pool
}

func (s *userStore) List(ctx context.Context, set func(*models.User) error) (err error) {
	rows, err := s.db.Query(
		ctx,
		`SELECT id, active, name, email, "timeCreated", "timeUpdated"
		FROM users`,
	)

	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		var user models.User

		if err = rows.Scan(&user.Id, &user.Active, &user.Name, &user.Email, &user.TimeCreated, &user.TimeUpdated); err != nil {
			return
		}

		if err = set(&user); err != nil {
			return
		}
	}

	return rows.Err()
}

func (s *userStore) Get(ctx context.Context, userId xid.ID, dst *models.User) error {
	row := s.db.QueryRow(
		ctx,
		`SELECT id, active, name, email, "timeCreated", "timeUpdated"
		FROM users
		WHERE id = $1`,
		userId,
	)

	return row.Scan(&dst.Id, &dst.Active, &dst.Name, &dst.Email, &dst.TimeCreated, &dst.TimeUpdated)
}

func (s *userStore) Create(ctx context.Context, user *models.User) (err error) {
	_, err = s.db.Exec(
		ctx,
		`INSERT INTO users(id, active, name, email, "timeCreated", "timeUpdated") VALUES($1, $2, $3, $4, $5, $6)`,
		user.Id,
		user.Active,
		user.Name,
		user.Email,
		user.TimeCreated,
		user.TimeUpdated,
	)

	return
}

func (s *userStore) SetPassword(ctx context.Context, passwordHash string, userId xid.ID) (err error) {
	_, err = s.db.Exec(
		ctx,
		`UPDATE users SET password = $1, "timeUpdated" = $2 
		WHERE id = $3`,
		passwordHash,
		pgtype.Timestamptz{Time: time.Now(), Valid: true},
		userId,
	)

	return
}

func (s *userStore) GetByEmail(ctx context.Context, email string, dst *models.User) (err error) {
	row := s.db.QueryRow(
		ctx,
		`SELECT id, active, name, email, "timeCreated", "timeUpdated"
		FROM users
		WHERE email = $1`,
		email,
	)

	return row.Scan(&dst.Id, &dst.Active, &dst.Name, &dst.Email, &dst.TimeCreated, &dst.TimeUpdated)
}
