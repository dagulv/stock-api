package routes

import (
	"github.com/dagulv/stock-api/internal/services"
	"github.com/labstack/echo/v4"
)

type StockRoutes struct {
	Service *services.StockService
}

func (r *StockRoutes) CurrentRoutes(e *echo.Group) {
	e.GET("/stock", r.get)
}

func (r *StockRoutes) get(ctx echo.Context) (err error) {
	if err = r.Service.Get(ctx.Request().Context()); err != nil {
		return
	}

	return
}
