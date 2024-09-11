package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/dagulv/stock-api/internal/core/domain"
	"github.com/dagulv/stock-api/internal/core/port"
	"github.com/dagulv/stock-api/internal/env"
	"github.com/labstack/echo/v4"
	"github.com/rs/xid"
)

const (
	RefreshTokenName = "refresh-token"
	BearerPrefix     = "Bearer "
)

var errHeaderValueInvalid = errors.New("invalid value in request header")

func GenerateSecretKeyFromEnv(env env.Env) (secret paseto.V4AsymmetricSecretKey) {
	secret, err := paseto.NewV4AsymmetricSecretKeyFromHex(env.AuthSecretKey)

	if err == nil {
		return secret
	}

	secret = paseto.NewV4AsymmetricSecretKey()

	fmt.Printf("New secret key: %s", secret.ExportHex())

	return secret
}

func Auth(sessionStore port.Auth, secretKey paseto.V4AsymmetricSecretKey, publicKey paseto.V4AsymmetricPublicKey, exemptRoutes map[string][]string) echo.MiddlewareFunc {
	parser := paseto.NewParser()
	parser.AddRule(paseto.Subject("id"), paseto.NotExpired())

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if check(c.Path(), c.Request().Method, exemptRoutes) {
				return next(c)
			}

			bearerToken, err := tokenFromBearer(c, parser, publicKey)

			if err == nil {
				userId, err := getUserIdFromToken(bearerToken)

				if err != nil {
					return c.NoContent(http.StatusUnauthorized)
				}

				c.Set("userId", userId)

				return next(c)
			}

			cookie, err := c.Cookie(RefreshTokenName)

			if err != nil {
				return c.NoContent(http.StatusUnauthorized)
			}

			token, err := parser.ParseV4Public(publicKey, cookie.Value, nil)

			if err != nil {
				UnsetCookie(c)
				return c.NoContent(http.StatusUnauthorized)
			}

			userId, err := getUserIdFromToken(token)

			if err != nil {
				return c.NoContent(http.StatusUnauthorized)
			}

			refreshSession, err := NewSession(c.Request().Context(), sessionStore, userId, time.Now().Add(time.Hour*24))

			if err != nil {
				return err
			}

			accessSession, err := NewSession(c.Request().Context(), sessionStore, userId, time.Now().Add(time.Hour))

			if err != nil {
				return err
			}

			signedRefreshToken := NewToken(refreshSession, secretKey)
			signedAccessToken := NewToken(accessSession, secretKey)

			SetCookie(c, signedRefreshToken)
			SetAuthHeader(c, signedAccessToken)

			c.Set("userId", userId)

			return nil
		}
	}
}

func tokenFromBearer(c echo.Context, parser paseto.Parser, publicKey paseto.V4AsymmetricPublicKey) (*paseto.Token, error) {
	authHeader := c.Request().Header.Get("Authorization")

	prefixLen := len(BearerPrefix)

	if len(authHeader) <= prefixLen || strings.EqualFold(authHeader[:prefixLen], BearerPrefix) {
		return nil, errHeaderValueInvalid
	}

	return parser.ParseV4Public(publicKey, authHeader[prefixLen:], nil)
}

func UnsetCookie(c echo.Context) {
	c.SetCookie(&http.Cookie{
		Name:     RefreshTokenName,
		Value:    "",
		Path:     "/auth",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
	})
}

func SetCookie(c echo.Context, signedToken string) {
	http.SetCookie(c.Response(), &http.Cookie{
		Name:     RefreshTokenName,
		Value:    signedToken,
		Path:     "/auth",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

func SetAuthHeader(c echo.Context, signedToken string) {
	c.Response().Header().Set("Authorization", signedToken)
}

func getUserIdFromToken(token *paseto.Token) (userId xid.ID, err error) {
	sub, err := token.GetSubject()

	if err != nil {
		return
	}

	userId, err = xid.FromString(sub)

	return
}

func NewSession(ctx context.Context, sessionStore port.Auth, userId xid.ID, expired time.Time) (session domain.Session, err error) {
	session = domain.Session{
		Id:          xid.New(),
		UserId:      userId,
		TimeExpired: expired,
	}

	return session, sessionStore.InsertSession(ctx, session)
}

func NewToken(session domain.Session, secretAuthKey paseto.V4AsymmetricSecretKey) string {
	token := paseto.NewToken()
	token.SetSubject(session.UserId.String())
	token.SetExpiration(session.TimeExpired)
	token.SetString("id", session.Id.String())

	return token.V4Sign(secretAuthKey, nil)
}

func check(path string, method string, m map[string][]string) bool {
	for k, v := range m {
		if k == path {
			for _, i := range v {
				if "*" == i || method == i {
					return true
				}
			}
		}
	}
	return false
}
