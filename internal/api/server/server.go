package server

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
}

func (s *Server) Start(ctx context.Context) (err error) {
	e := echo.New()

	e.Use(middleware.CORS())

	// e.Group("/api")

	return e.Start(":8080")
}
