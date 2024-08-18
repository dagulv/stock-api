package http

import (
	"context"
	"net/http"

	"github.com/dagulv/stock-api/internal/adapter/server"
	"github.com/dagulv/stock-api/internal/core/service"
	"github.com/dagulv/stock-api/internal/env"
	"github.com/go-playground/validator/v10"
	"github.com/go-webauthn/webauthn/webauthn"
	jsoniter "github.com/json-iterator/go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	Json     jsoniter.API
	WebAuthn *webauthn.WebAuthn
	Tick     service.Ticker
	User     service.User
	Env      env.Env
}

func (s Server) StartServer(ctx context.Context) (err error) {
	e := echo.New()

	e.Validator = &server.Validate{
		Validator: validator.New(),
	}
	//TODO Add jsoniter de/serializer, error handling, etc

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		Skipper:          middleware.DefaultSkipper,
		AllowOrigins:     []string{"localhost:3000", s.Env.AppUrl},
		AllowMethods:     []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
		AllowCredentials: true,
	}))

	s.addRoutes(e)
	return e.Start(":3001")
}
