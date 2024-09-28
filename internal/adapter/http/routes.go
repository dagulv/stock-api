package http

import (
	"github.com/dagulv/stock-api/internal/adapter/http/routes"
	"github.com/labstack/echo/v4"
)

func (s Server) addRoutes(e *echo.Echo) {
	// routes.Routes(e, s.Tick, s.Json)
	routes.UserRoutes(e, s.User, s.Json)
	routes.AuthRoutes(e, s.Auth, s.Json)
}
