package server

import (
	"context"

	"github.com/dagulv/stock-api/internal/api/routes"
	"github.com/dagulv/stock-api/internal/json"
	"github.com/dagulv/stock-api/internal/services"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	UserService  *services.UserService
	StockService *services.StockService
}

func (s *Server) Start(ctx context.Context) (err error) {
	e := echo.New()
	e.JSONSerializer = json.Serializer{}

	e.Use(middleware.CORS())

	routes.Routes(&routes.UserRoutes{
		Service: s.UserService,
	}).CurrentRoutes(e.Group("/api"))

	routes.Routes(&routes.StockRoutes{
		Service: s.StockService,
	}).CurrentRoutes(e.Group("/api"))

	return e.Start(":8080")
}
