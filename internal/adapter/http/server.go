package http

import (
	"context"
	"net/http"

	"github.com/dagulv/stock-api/internal/adapter/server"
	"github.com/dagulv/stock-api/internal/core/service"
	"github.com/dagulv/stock-api/internal/env"
	"github.com/go-playground/validator/v10"
	jsoniter "github.com/json-iterator/go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	Json           jsoniter.API
	MiddlewareAuth echo.MiddlewareFunc
	Auth           service.Auth
	Tick           service.Ticker
	User           service.User
	Env            env.Env
}

func (s Server) StartServer(ctx context.Context) (err error) {
	e := echo.New()

	e.Validator = &server.Validate{
		Validator: validator.New(),
	}
	//TODO Add jsoniter de/serializer, error handling, etc

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		Skipper:          middleware.DefaultSkipper,
		AllowOrigins:     []string{"http://localhost:3000", s.Env.AppUrl},
		AllowMethods:     []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
		AllowCredentials: true,
	}))

	e.Use(server.Auth(s.Auth.Store, s.Auth.Parser, s.Auth.SecretAuthKey, s.Auth.SecretAuthKey.Public(), map[string][]string{
		"/login":        {http.MethodPost},
		"/access-token": {http.MethodGet},
		"/register":     {http.MethodPost},
	}))

	s.addRoutes(e)
	return e.Start(":3001")
}
