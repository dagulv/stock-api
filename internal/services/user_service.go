package services

import (
	"context"
	"time"

	"github.com/dagulv/stock-api/internal/models"
	"github.com/dagulv/stock-api/internal/stores"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	Store    stores.UserStore
	Validate *validator.Validate
}

func (s UserService) List(ctx context.Context, set func(*models.User) error) (err error) {
	return s.Store.List(ctx, set)
}

func (s UserService) Create(ctx context.Context, user *models.User) (err error) {
	user.Id = xid.NewWithTime(time.Now())
	user.TimeCreated = pgtype.Timestamptz{Time: time.Now(), Valid: true}
	user.TimeUpdated = user.TimeCreated

	s.Store.Create(ctx, user)

	return
}

func (s UserService) Get(ctx context.Context, userId xid.ID, dst *models.User) error {
	return s.Store.Get(ctx, userId, dst)
}

func (s UserService) GetByEmail(ctx context.Context, email string, dst *models.User) error {
	return s.Store.GetByEmail(ctx, email, dst)
}

func (s UserService) SetPassword(ctx context.Context, password string, userId xid.ID) (err error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 10)

	if err != nil {
		return
	}

	return s.Store.SetPassword(ctx, string(passwordHash[:]), userId)
}

func (s UserService) Put(ctx context.Context, user *models.User) error {
	user.TimeUpdated = pgtype.Timestamptz{Valid: true, Time: time.Now()}

	return s.Store.Put(ctx, user)
}

func (s UserService) Delete(ctx context.Context, userId xid.ID) error {
	return s.Store.Delete(ctx, userId)
}
