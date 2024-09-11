package routes

import (
	"net/http"

	"github.com/dagulv/stock-api/internal/adapter/server"
	"github.com/dagulv/stock-api/internal/core/domain"
	"github.com/dagulv/stock-api/internal/core/service"
	jsoniter "github.com/json-iterator/go"
	"github.com/labstack/echo/v4"
	"github.com/rs/xid"
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
	e.PUT("/auth/password", r.newPassword)
}

func (r authRoutes) auth(c echo.Context) (err error) {
	var sessionUser *domain.SessionUser

	sessionUserId := c.Get("sessionUserId").(xid.ID)

	if sessionUser, err = r.Service.Auth(c.Request().Context(), sessionUserId); err != nil {
		return
	}

	return c.JSON(http.StatusOK, *sessionUser)
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
	server.SetAuthHeader(c, signedAccessToken)

	return c.NoContent(http.StatusOK)
}

func (r authRoutes) logout(c echo.Context) (err error) {
	server.UnsetCookie(c)

	return c.NoContent(http.StatusOK)
}

func (r authRoutes) newPassword(c echo.Context) (err error) {
	return
}
