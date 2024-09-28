package service

import (
	"context"
	"fmt"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/dagulv/stock-api/internal/adapter/mailer"
	"github.com/dagulv/stock-api/internal/adapter/server"
	"github.com/dagulv/stock-api/internal/core/domain"
	"github.com/dagulv/stock-api/internal/core/port"
	"github.com/dagulv/stock-api/internal/env"
	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	SecretAuthKey paseto.V4AsymmetricSecretKey
	Parser        paseto.Parser
	Mailer        *mailer.Mailer
	Env           env.Env
	Store         port.Auth
	UserStore     port.User
}

// TODO: Change to Login WithPassword func as parameter that returns bool
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

	session, err := server.NewSession(ctx, s.Store, existingCredentials.UserId, time.Now().Add(time.Hour*24), domain.ScopeAuthentication)

	if err != nil {
		return
	}

	signedRefreshToken = server.NewToken(session.Id, session.TimeExpired, s.SecretAuthKey)
	signedAccessToken = server.NewToken(session.Id, time.Now().Add(time.Hour), s.SecretAuthKey)

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

func (s Auth) VerifyEmail(ctx context.Context, email string) (err error) {
	var existingCredentials domain.Credentials

	err = s.Store.GetCredentialsByEmail(ctx, email, &existingCredentials)

	if err != nil {
		return
	}

	if !existingCredentials.UserId.IsNil() {
		//return already exists error
		return
	}

	user := domain.User{
		Email: email,
	}

	if err = s.UserStore.Create(ctx, &user); err != nil {
		return
	}

	//Create session from created user
	session, err := server.NewSession(ctx, s.Store, user.Id, time.Now().Add(time.Hour*24), domain.ScopeRegister)

	if err != nil {
		return
	}

	signedRefreshToken := server.NewToken(session.Id, session.TimeExpired, s.SecretAuthKey)

	template := fmt.Sprintf(`
	%s
	`, s.Env.AppUrl+"/recover/"+signedRefreshToken)

	mail := mailer.Email{
		Email:        email,
		Subject:      "Sign in to " + s.Env.RecipientEmail,
		HTMLTemplate: &template,
	}

	s.Mailer.Send(mail)

	return
}

func (s Auth) SessionUserFromToken(ctx context.Context, token *paseto.Token) (_ *domain.SessionUser, err error) {
	sessionId, err := server.SessionIdFromToken(token)

	if err != nil {
		return
	}

	return s.Store.LazyGetSessionUser(ctx, sessionId)
}
