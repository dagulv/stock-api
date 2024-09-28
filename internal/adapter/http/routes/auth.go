package routes

import (
	"net/http"
	"time"

	"github.com/dagulv/stock-api/internal/adapter/server"
	"github.com/dagulv/stock-api/internal/core/domain"
	"github.com/dagulv/stock-api/internal/core/service"
	jsoniter "github.com/json-iterator/go"
	"github.com/labstack/echo/v4"
)

type authRoutes struct {
	Json    jsoniter.API
	Service service.Auth
}

func AuthRoutes(e *echo.Echo, s service.Auth, jsonApi jsoniter.API) {
	r := authRoutes{
		Json:    jsonApi,
		Service: s,
	}

	e.GET("/auth", r.auth)
	e.POST("/auth", r.login)
	e.DELETE("/auth", r.logout)
	e.GET("/auth/access-token", r.accessToken)
	e.PUT("/auth/password", r.newPassword)
	e.POST("/auth/sign-up/email", r.signUpEmail)
}

func (r authRoutes) auth(c echo.Context) (err error) {
	sessionUser, ok := c.Get(server.SessionUserKey).(*domain.SessionUser)

	if !ok {
		return c.NoContent(http.StatusUnauthorized)
	}

	return c.JSON(http.StatusOK, sessionUser)
}

func (r authRoutes) login(c echo.Context) (err error) {
	var credentials domain.Credentials

	if err = c.Bind(&credentials); err != nil {
		return
	}

	signedRefreshToken, signedAccessToken, err := r.Service.LoginWithPassword(c.Request().Context(), credentials)

	if err != nil {
		return
	}

	server.SetCookie(c, signedRefreshToken)

	return c.JSON(http.StatusOK, map[string]string{
		"accessToken": signedAccessToken,
	})
}

func (r authRoutes) logout(c echo.Context) (err error) {
	server.UnsetCookie(c)

	return c.NoContent(http.StatusOK)
}

func (r authRoutes) accessToken(c echo.Context) (err error) {
	ctx := c.Request().Context()

	cookie, err := c.Cookie(server.RefreshTokenName)

	if err != nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	token, err := r.Service.Parser.ParseV4Public(r.Service.SecretAuthKey.Public(), cookie.Value, nil)

	if err != nil {
		server.UnsetCookie(c)
		return c.NoContent(http.StatusUnauthorized)
	}

	sessionUser, err := r.Service.SessionUserFromToken(ctx, token)

	if err != nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	session, err := server.NewSession(ctx, r.Service.Store, sessionUser.Id, time.Now().Add(time.Hour*24), domain.ScopeAuthentication)

	if err != nil {
		return err
	}

	signedRefreshToken := server.NewToken(session.Id, session.TimeExpired, r.Service.SecretAuthKey)
	signedAccessToken := server.NewToken(session.Id, time.Now().Add(time.Hour), r.Service.SecretAuthKey)

	server.SetCookie(c, signedRefreshToken)

	return c.JSON(http.StatusOK, map[string]string{
		"accessToken": signedAccessToken,
	})
}

func (r authRoutes) newPassword(c echo.Context) (err error) {
	return
}

func (r authRoutes) signUpEmail(c echo.Context) (err error) {
	var cred struct {
		Email string `json:"email" validate:"required,email"`
	}

	if err = c.Bind(&cred); err != nil {
		return
	}

	if err = r.Service.VerifyEmail(c.Request().Context(), cred.Email); err != nil {
		return
	}

	return c.NoContent(http.StatusOK)
}
