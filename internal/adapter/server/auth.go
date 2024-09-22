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
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/rs/xid"
)

const (
	RefreshTokenName = "refresh-token"
	BearerPrefix     = "Bearer "
	SessionUserKey   = "user"
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

func Auth(sessionStore port.Auth, parser paseto.Parser, secretKey paseto.V4AsymmetricSecretKey, publicKey paseto.V4AsymmetricPublicKey, exemptRoutes map[string][]string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Add("Vary", "Authorization")

			if check(c.Path(), c.Request().Method, exemptRoutes) {
				return next(c)
			}

			bearerToken, err := tokenFromBearer(c, parser, publicKey)

			if err != nil {
				return c.NoContent(http.StatusUnauthorized)
			}

			sessionId, err := SessionIdFromToken(bearerToken)

			if err != nil {
				return c.NoContent(http.StatusUnauthorized)
			}

			sessionUser, err := sessionStore.LazyGetSessionUser(c.Request().Context(), sessionId)

			if err != nil {
				//TODO: Move pgx to db and replace with custom error
				if errors.Is(err, pgx.ErrNoRows) {
					c.Response().Header().Set("WWW-Authenticate", "Bearer")
					return c.NoContent(http.StatusUnauthorized)
				}

				return err
			}

			c.Set(SessionUserKey, sessionUser)

			return next(c)
		}
	}
}

func tokenFromBearer(c echo.Context, parser paseto.Parser, publicKey paseto.V4AsymmetricPublicKey) (*paseto.Token, error) {
	authHeader := c.Request().Header.Get("Authorization")

	prefixLen := len(BearerPrefix)

	if len(authHeader) <= prefixLen || strings.EqualFold(authHeader[:prefixLen], BearerPrefix) {
		c.Response().Header().Set("WWW-Authenticate", "Bearer")
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

func SessionIdFromToken(token *paseto.Token) (sessionId xid.ID, err error) {
	sub, err := token.GetSubject()

	if err != nil {
		return
	}

	sessionId, err = xid.FromString(sub)

	return
}

func NewSession(ctx context.Context, sessionStore port.Auth, userId xid.ID, expired time.Time, scope int) (session domain.Session, err error) {
	session = domain.Session{
		Id:          xid.New(),
		UserId:      userId,
		Scope:       scope,
		TimeExpired: expired,
	}

	if err = sessionStore.InsertSession(ctx, session); err != nil {
		return
	}

	return session, nil
}

func NewToken(sessionId xid.ID, timeExpired time.Time, secretAuthKey paseto.V4AsymmetricSecretKey) string {
	token := paseto.NewToken()
	token.SetSubject(sessionId.String())
	token.SetExpiration(timeExpired)

	return token.V4Sign(secretAuthKey, nil)
}

func check(path string, method string, m map[string][]string) bool {
	for k, v := range m {
		if k == path {
			for _, i := range v {
				if i == "*" || method == i {
					return true
				}
			}
		}
	}
	return false
}
