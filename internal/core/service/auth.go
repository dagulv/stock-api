package service

import (
	"context"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/dagulv/stock-api/internal/adapter/server"
	"github.com/dagulv/stock-api/internal/core/domain"
	"github.com/dagulv/stock-api/internal/core/port"
	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	SecretAuthKey paseto.V4AsymmetricSecretKey
	Store         port.Auth
}

func (s Auth) Auth(ctx context.Context, sessionUserId xid.ID) (sessionUser *domain.SessionUser, err error) {
	return s.Store.LazyGetSessionUser(ctx, sessionUserId)
}

func (s Auth) LoginWithPassword(ctx context.Context, credentials domain.Credentials) (signedRefreshToken string, signedAccessToken string, err error) {
	var existingCredentials domain.Credentials

	err = s.Store.GetCredentialsByEmail(ctx, credentials.Email, &existingCredentials)

	if err != nil {
		bcrypt.CompareHashAndPassword(nil, nil)

		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(existingCredentials.Password), []byte(credentials.Password))

	if err != nil {
		return
	}

	refreshSession, err := server.NewSession(ctx, s.Store, existingCredentials.UserId, time.Now().Add(time.Hour*24))

	if err != nil {
		return
	}

	accessSession, err := server.NewSession(ctx, s.Store, existingCredentials.UserId, time.Now().Add(time.Hour))

	if err != nil {
		return
	}

	signedRefreshToken = server.NewToken(refreshSession, s.SecretAuthKey)
	signedAccessToken = server.NewToken(accessSession, s.SecretAuthKey)

	return
}

func (s Auth) NewPassword(ctx context.Context, credentials domain.Credentials) (err error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(credentials.Password), bcrypt.DefaultCost)

	if err != nil {
		return
	}

	var existingCredentials domain.Credentials

	err = s.Store.GetCredentialsByEmail(ctx, credentials.Email, &existingCredentials)

	if err != nil {
		return
	}

	return s.Store.UpdatePassword(ctx, existingCredentials.UserId, hashedPassword)
}
